package discover

// Packages ...
type Packages []Package

// Package ...
type Package struct {
	Name      string
	Thankable bool
}

// Repository
type Repository struct {
	Name  string
	Owner string
	URL   string
}

// Discoverer ...
type Discoverer interface {
	InUse() (bool, error)
	DiscoverRepositories() ([]Repository, error)
}

// Thankable
func ThankableRepositories(discoverers []Discoverer) ([]Repository, error) {
	repos := []Repository{}

	for _, d := range discoverers {
		if inUse, err := d.InUse(); !inUse || err != nil {
			continue
		}
		rs, err := d.DiscoverRepositories()
		if err != nil {
			return nil, err
		}

		repos = append(repos, rs...)
	}

	return repos, nil
}
