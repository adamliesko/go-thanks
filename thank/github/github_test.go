package github

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamliesko/go-thanks/discover"
)

func TestThanker_Auth(t *testing.T) {
	t.Parallel()

	const goodAccessToken = "good-auth"

	tcs := []struct {
		name      string
		apiToken  string
		wantError bool
	}{
		{
			name:      "ok",
			apiToken:  goodAccessToken,
			wantError: false,
		},
		{
			name:      "bad authr",
			apiToken:  "bad-auth",
			wantError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokens, ok := r.URL.Query()["access_token"]
		if !ok || len(tokens) != 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if tokens[0] != goodAccessToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	for _, tc := range tcs {
		th := New(tc.apiToken)
		th.apiEndpoint = ts.URL
		err := th.Auth()
		if (err != nil) != tc.wantError {
			t.Errorf("%s: Github Auth() unexpected result, wantErr: %v err: %v ", tc.name, tc.wantError, err)
		}
	}
}

func TestThanker_Thank(t *testing.T) {
	t.Parallel()

	const goodAccessToken = "good-auth"

	tcs := []struct {
		name      string
		repo      discover.Repository
		apiToken  string
		wantError bool
	}{
		{
			name: "ok",
			repo: discover.Repository{
				URL:   "github.com/user/repo",
				Owner: "user",
				Name:  "repo",
			},
			apiToken:  goodAccessToken,
			wantError: false,
		},
		{
			name: "repo not found",
			repo: discover.Repository{
				URL:   "github.com/another/repo",
				Owner: "another",
				Name:  "repo",
			},
			apiToken:  goodAccessToken,
			wantError: true,
		},
		{
			name: "bad auth token",
			repo: discover.Repository{
				URL:   "github.com/user/repo",
				Owner: "user",
				Name:  "repo",
			},
			apiToken:  "bad-auth",
			wantError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokens, ok := r.URL.Query()["access_token"]
		if !ok || len(tokens) != 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if tokens[0] != goodAccessToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.URL.Path != "/user/starred/user/repo" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	for _, tc := range tcs {
		th := New(tc.apiToken)
		th.apiEndpoint = ts.URL
		err := th.Thank(tc.repo)
		if (err != nil) != tc.wantError {
			t.Errorf("%s: Github Thank() unexpected result, wantErr: %v err: %v ", tc.name, tc.wantError, err)
		}
	}
}

func TestThanker_CanThank(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		repo       discover.Repository
		wantResult bool
	}{
		{
			repo: discover.Repository{
				Name:  "semver",
				Owner: "Masterminds",
				URL:   "github.com/Masterminds/semver",
			},
			wantResult: true,
		},
		{
			repo: discover.Repository{
				Name:  "semver",
				Owner: "Masterminds",
				URL:   "gitlab.com/Masterminds/semver",
			},
			wantResult: false,
		},
	}

	th := New("")

	for _, tc := range tcs {
		if got := th.CanThank(tc.repo); got != tc.wantResult {
			t.Errorf("got bad CanThank() result for repo %s", tc.repo.URL)
		}
	}
}

func TestThankURL(t *testing.T) {
	t.Parallel()

	th := New("token")
	repo := discover.Repository{
		Name:  "semver",
		Owner: "Masterminds",
		URL:   "gitlab.com/Masterminds/semver",
	}
	url, err := th.thankURL(repo)
	if err != nil {
		t.Fatal(err)
	}
	wantURL := "https://api.github.com/user/starred/Masterminds/semver?access_token=token"
	if url != wantURL {
		t.Errorf("bad url, got %s, want %s", url, wantURL)
	}
}
