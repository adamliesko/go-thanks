package thank

import (
	"log"

	"github.com/adamliesko/go-thanks/discover"
	"github.com/adamliesko/go-thanks/thank/github"
	"github.com/adamliesko/go-thanks/thank/gitlab"
)

// Thanker ...
type Thanker interface {
	CanThank(discover.Repository) (bool, error)
	Thank(discover.Repository) error
	Auth() error
}

// Thankers producers a list of verified and authenticated thankers.
func Thankers(githubToken, gitlabToken string) ([]Thanker, error) {
	var ts []Thanker
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

// Thank thanks to all the repositories and their owners using one of the passed in thankers, by starring the repositories.
func Thank(ts []Thanker, repos []discover.Repository) (int, error) {
	var thankedCount int

	for _, r := range repos {
		for _, s := range ts {
			if can, err := s.CanThank(r); !can || err != nil {
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
