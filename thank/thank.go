// Package thank thanks (usually by starring) repositories providing Go packages.
package thank

import (
	"log"

	"github.com/adamliesko/go-thanks/discover"
	"github.com/adamliesko/go-thanks/thank/github"
	"github.com/adamliesko/go-thanks/thank/gitlab"
)

// Thanker abstracts away a mean of thanking for a Go OSS contribution, usually by starring a repository.
type Thanker interface {
	Auth() error
	CanThank(discover.Repository) bool
	Thank(discover.Repository) error
}

// Thankers produces a list of verified and authenticated thankers.
func Thankers(ts []Thanker, githubToken, gitlabToken string) ([]Thanker, error) {
	if githubToken != "" {
		gt := github.New(githubToken)
		if err := gt.Auth(); err != nil {
			return nil, err
		}
		ts = append(ts, gt)
	}
	if gitlabToken != "" {
		gt := gitlab.New(gitlabToken)
		if err := gt.Auth(); err != nil {
			return nil, err
		}
		ts = append(ts, gt)
	}
	return ts, nil
}

// Thank thanks to all the repositories and their owners using one of the passed in thankers, usually by starring the
// repositories.
func Thank(ts []Thanker, repos []discover.Repository) (int, error) {
	var thankedCount int

	for _, r := range repos {
		for _, s := range ts {
			if !s.CanThank(r) {
				continue
			}
			err := s.Thank(r)
			if err != nil {
				return thankedCount, err
			}
			log.Printf("Thanked to repository %s by %s\n", r.URL, r.Owner)
			thankedCount++
			break
		}
	}

	return thankedCount, nil
}
