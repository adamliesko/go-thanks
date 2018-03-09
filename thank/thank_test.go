package thank

import (
	"errors"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/adamliesko/go-thanks/discover"
)

const goodAPIToken = "pass123"

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestAuthThankers(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		token     string
		wantError bool
	}{
		{
			token:     goodAPIToken,
			wantError: false,
		},
		{
			token:     "bad",
			wantError: true,
		},
	}

	for _, tc := range tcs {
		ft := &FakeThanker{
			thanked:  map[string]bool{},
			apiToken: tc.token,
		}
		_, err := AuthThankers([]Thanker{ft})
		if (err == nil) == tc.wantError {
			t.Errorf("unexpected error response %v for token %s wantError %v", err, tc.token, tc.wantError)
		}
	}
}

func TestThank(t *testing.T) {
	t.Parallel()

	ft := &FakeThanker{
		thanked: map[string]bool{},
	}
	repos := []discover.Repository{
		{URL: "fake-thanker.com/me/em", Owner: "me", Name: "em"},
		{URL: "fake-thanker.com/me/ok", Owner: "me", Name: "ok"},
		{URL: "xx-thanker.com/me/ok", Owner: "me", Name: "ok"},
	}
	thanked, err := Thank([]Thanker{ft}, repos)
	if err != nil {
		t.Fatal("unexpected error from thanker")
	}
	if thanked != 2 {
		t.Errorf("got bad thanked count, got %d want 2", thanked)
	}
	if len(ft.thanked) != 2 {
		t.Errorf("got bad thanked state, got %d want 2 thanked", len(ft.thanked))
	}
}

func TestThankPartiallyFaulty(t *testing.T) {
	t.Parallel()

	ft := &FakeThanker{
		thanked: map[string]bool{},
	}
	repos := []discover.Repository{
		{URL: "fake-thanker.com/me/em", Owner: "me", Name: "em"},
		{URL: "fake-thanker.com/me/error", Owner: "me", Name: "error"},
		{URL: "fake-thanker.com/me/eme", Owner: "me", Name: "em"},
		{URL: "xxx-thanker.com/me/eme", Owner: "me", Name: "em"},
	}
	thanked, err := Thank([]Thanker{ft}, repos)
	if err == nil {
		t.Fatal("missing error from bad thanker")
	}
	if thanked != 1 {
		t.Errorf("got bad thanked count, got %d want 1", thanked)
	}
	if len(ft.thanked) != 1 {
		t.Errorf("got bad thanked state, got %d want 1 thanked", len(ft.thanked))
	}
}

type FakeThanker struct {
	apiToken string
	thanked  map[string]bool
}

func (ft *FakeThanker) CanThank(r discover.Repository) bool {
	return strings.HasPrefix(r.URL, "fake-thanker.com")
}

func (ft *FakeThanker) Thank(r discover.Repository) error {
	if r.URL == "fake-thanker.com/me/error" {
		return errors.New("forced error")
	}
	ft.thanked[r.URL] = true
	return nil
}

func (ft *FakeThanker) Auth() error {
	if ft.apiToken == goodAPIToken {
		return nil
	}
	return errors.New("bad auth")
}
