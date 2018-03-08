package dep

import "github.com/adamliesko/go-thanks/discover"

type Discoverer struct{}

func (d Discoverer) InUse() (bool, error) {
	return false, nil
}

func (d Discoverer) DiscoverRepositories() ([]discover.Repository, error) {
	return nil, nil
}
