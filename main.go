package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/adamliesko/go-thanks/discover"
	"github.com/adamliesko/go-thanks/thank"
	"github.com/pkg/errors"
)

var (
	githubToken = flag.String("github-token", os.Getenv("GITHUB_API_TOKEN"), "Github API token. Defaults to env variable GITHUB_API_TOKEN.")
	gitlabToken = flag.String("gitlab-token", os.Getenv("GITLAB_API_TOKEN"), "Gitlab API token. Defaults to env variable GITLAB_API_TOKEN.")
	path        = flag.String("path", ".", "Path to Go project.")
)

func thankGiants() error {
	ts, err := thank.Thankers(*githubToken, *gitlabToken)
	if err != nil {
		return fmt.Errorf("error getting available thankers: %v", err)
	}
	if len(ts) == 0 {
		return errors.New("none capable thankers found")
	}

	repos, err := discover.DiscoverRepositories(*path)
	if err != nil {
		return fmt.Errorf("error getting thankable repositories: %v", err)
	}

	thanked, err := thank.Thank(ts, repos)
	if err != nil {
		return fmt.Errorf("error thanking: %v", err)
	}

	log.Printf("Thanked to %d repositories.", thanked)
	return nil
}

func main() {
	if err := thankGiants(); err != nil {
		log.Fatal(err)
	}
}
