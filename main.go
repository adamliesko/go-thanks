package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/adamliesko/go-thanks/discover"
	"github.com/adamliesko/go-thanks/thank"
)

var (
	githubToken = flag.String("github-token", os.Getenv("GITHUB_TOKEN"), "Github API token. Defaults to env variable GITHUB_TOKEN.")
	gitlabToken = flag.String("gitlab-token", os.Getenv("GITLAB_TOKEN"), "Gitlab API token. Defaults to env variable GITLAB_TOKEN.")
	path        = flag.String("path", ".", "Path to Go project.")
)

func thankGiants() error {
	log.Println("==== Discovering ====")

	repos, err := discover.DiscoverRepositories(*path)
	if err != nil {
		return fmt.Errorf("error getting thankable repositories: %v", err)
	}

	log.Printf("Discovered %d repositories", len(repos))

	log.Println("====== Thanking =====")

	ts, err := thank.Thankers(*githubToken, *gitlabToken)
	if err != nil {
		return fmt.Errorf("error getting available thankers: %v", err)
	}
	if len(ts) == 0 {
		return errors.New("none capable thankers found")
	}

	thanked, err := thank.Thank(ts, repos)
	if err != nil {
		return fmt.Errorf("error thanking: %v", err)
	}

	log.Println("======== Done =======")
	log.Printf("Thanked to %d repositories 🙏", thanked)
	return nil
}

func main() {
	flag.Parse()

	if err := thankGiants(); err != nil {
		log.Fatal(err)
	}
}
