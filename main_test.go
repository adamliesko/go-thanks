package main

import (
	"errors"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"testing"

	"github.com/adamliesko/go-thanks/discover"
	"github.com/adamliesko/go-thanks/thank"
)

const goodApiToken = "pass123"

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestThankingTheGiants(t *testing.T) {
	*projectPath = path.Join("discover", "test", "project_dep")
	ft := &FakeGithubThanker{
		apiToken: goodApiToken,
		thanked:  map[string]bool{},
	}
	err := thankGiants([]thank.Thanker{ft}, *projectPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(ft.thanked) != 3 {
		t.Errorf("bad thanked count, got %d want %d", len(ft.thanked), 3)
	}
}

func TestThankingTheGiantsWithBadTokenFails(t *testing.T) {
	*projectPath = path.Join("discover", "test", "project_dep")
	ft := &FakeGithubThanker{
		apiToken: "bad-token",
		thanked:  map[string]bool{},
	}
	err := thankGiants([]thank.Thanker{ft}, *projectPath)
	if err == nil {
		t.Fatal(err)
	}

	if len(ft.thanked) != 0 {
		t.Errorf("bad thanked count, got %d want %d", len(ft.thanked), 0)
	}
}

type FakeGithubThanker struct {
	apiToken string
	thanked  map[string]bool
}

func (ft *FakeGithubThanker) CanThank(r discover.Repository) bool {
	return strings.HasPrefix(r.URL, "github.com")
}

func (ft *FakeGithubThanker) Thank(r discover.Repository) error {
	ft.thanked[r.URL] = true
	return nil
}

func (ft *FakeGithubThanker) Auth() error {
	if ft.apiToken == goodApiToken {
		return nil
	}
	return errors.New("bad auth")
}
