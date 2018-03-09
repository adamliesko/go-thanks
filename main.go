package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"time"

	"github.com/adamliesko/go-thanks/discover"
	"github.com/adamliesko/go-thanks/thank"
	"github.com/adamliesko/go-thanks/thank/github"
	"github.com/adamliesko/go-thanks/thank/gitlab"
)

var (
	githubToken = flag.String("github-token", os.Getenv("GITHUB_TOKEN"), "Github API token. Defaults to env variable GITHUB_TOKEN.")
	gitlabToken = flag.String("gitlab-token", os.Getenv("GITLAB_TOKEN"), "Gitlab API token. Defaults to env variable GITLAB_TOKEN.")
	projectPath = flag.String("project-path", ".", "Path to Go project.")
)

func thankGiants(thankers []thank.Thanker, path string) error {
	startedAt := time.Now()

	log.Println("==== Discovering ====")
	repos, err := discover.DiscoverRepositories(path)
	if err != nil {
		return fmt.Errorf("error getting thankable repositories: %v", err)
	}
	log.Printf("Discovered %d repositories\n", len(repos))
	if len(repos) == 0 {
		return nil
	}

	log.Println("====== Thanking =====")
	ts, err := thank.AuthThankers(thankers)
	if err != nil {
		return fmt.Errorf("error authenticating available thankers: %v", err)
	}
	if len(ts) == 0 {
		return errors.New("none authenticated thankers found")
	}

	thanked, err := thank.Thank(ts, repos)
	if err != nil {
		return fmt.Errorf("error thanking: %v", err)
	}

	log.Println("======== Done =======")
	took := time.Now().Sub(startedAt)
	log.Printf("Thanked to %d repositories 🙏! (took %v)\n", thanked, took.Round(time.Millisecond))
	return nil
}

func run() error {
	var ts []thank.Thanker
	if *githubToken != "" {
		gt := github.New(*githubToken)
		ts = append(ts, gt)
	}
	if *gitlabToken != "" {
		gt := gitlab.New(*gitlabToken)
		ts = append(ts, gt)
	}
	return thankGiants(ts, *projectPath)
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
