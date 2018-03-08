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
	githubToken = flag.String("github-token", os.Getenv("GITHUB_TOKEN"), "Github API token. Defaults to env variable GITHUB_TOKEN.")
	gitlabToken = flag.String("gitlab-token", os.Getenv("GITLAB_TOKEN"), "Gitlab API token. Defaults to env variable GITLAB_TOKEN.")
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
	flag.Parse()

	if err := thankGiants(); err != nil {
		log.Fatal(err)
	}
}
