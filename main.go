package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/adamliesko/go-thanks/discover"
	"github.com/adamliesko/go-thanks/discover/govendor"
	"github.com/adamliesko/go-thanks/thank"
	"github.com/pkg/errors"
)

var (
	githubToken = flag.String("github-token", os.Getenv("GITHUB_API_TOKEN"), "Github API token. Default env variable GITHUB_API_TOKEN.")
	gitlabToken = flag.String("gitlab-token", os.Getenv("GITLAB_API_TOKE "), "Gitlab API token. Default env variable GITLAB_API_TOKEN.")
)

func thankGiants() error {
	log.SetPrefix("go-thanks: ")
	ts, err := thank.Thankers(*githubToken, *gitlabToken)
	if err != nil {
		return fmt.Errorf("error getting available thankers: %v", err)
	}
	if len(ts) == 0 {
		return errors.New("none capable thankers found")
	}

	repos, err := discover.ThankableRepositories([]discover.Discoverer{govendor.Discoverer{}})
	if err != nil {
		return fmt.Errorf("error getting thankable repositories: %v", err)
	}

	thanked, err := thank.Thank(ts, repos)
	if err != nil {
		return fmt.Errorf("error thanking: %v", err)
	}
	log.Printf("Thanked to %d repositories of your giants' packages", thanked)

	return nil
}

func main() {
	if err := thankGiants(); err != nil {
		log.Fatal(err)
	}
}
