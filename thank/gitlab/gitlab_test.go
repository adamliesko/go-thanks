package gitlab

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamliesko/go-thanks/discover"
)

func TestThanker_Auth(t *testing.T) {
	t.Parallel()

	const goodPrivateToken = "good-auth"

	tcs := []struct {
		name      string
		apiToken  string
		wantError bool
	}{
		{
			name:      "ok",
			apiToken:  goodPrivateToken,
			wantError: false,
		},
		{
			name:      "bad auth token",
			apiToken:  "bad-auth",
			wantError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokens, ok := r.URL.Query()["private_token"]
		if !ok || len(tokens) != 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if tokens[0] != goodPrivateToken {
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
			t.Errorf("%s: Gitlab Auth() unexpected result, wantErr: %v err: %v ", tc.name, tc.wantError, err)
		}
	}
}

func TestThanker_Thank(t *testing.T) {
	t.Parallel()

	const goodPrivateToken = "good-auth"

	tcs := []struct {
		name      string
		repo      discover.Repository
		apiToken  string
		wantError bool
	}{
		{
			name: "ok",
			repo: discover.Repository{
				URL:   "gitlab.com/user/repo",
				Owner: "user",
				Name:  "repo",
			},
			apiToken:  goodPrivateToken,
			wantError: false,
		},
		{
			name: "repo not found",
			repo: discover.Repository{
				URL:   "github.com/another/repo",
				Owner: "another",
				Name:  "repo",
			},
			apiToken:  goodPrivateToken,
			wantError: true,
		},
		{
			name: "bad auth token",
			repo: discover.Repository{
				URL:   "gitlab.com/user/repo",
				Owner: "user",
				Name:  "repo",
			},
			apiToken:  "bad-auth",
			wantError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokens, ok := r.URL.Query()["private_token"]
		if !ok || len(tokens) != 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if tokens[0] != goodPrivateToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.URL.Path != "/projects/user/repo/star" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	for _, tc := range tcs {
		th := New(tc.apiToken)
		th.apiEndpoint = ts.URL
		err := th.Thank(tc.repo)
		if (err != nil) != tc.wantError {
			t.Errorf("%s: Gitlab Thank() unexpected result, wantErr: %v err: %v ", tc.name, tc.wantError, err)
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
				URL:   "gitlab.com/Masterminds/semver",
			},
			wantResult: true,
		},
		{
			repo: discover.Repository{
				Name:  "semver",
				Owner: "Masterminds",
				URL:   "github.com/Masterminds/semver",
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
