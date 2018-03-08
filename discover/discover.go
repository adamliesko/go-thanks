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
	DiscoverRepositories() (map[string]Repository, error)
}

// ThankableRepositories produces a slice of repositories that can be later thanked to using a set of known discoverers
// exploring the package managers.
func ThankableRepositories(discoverers []Discoverer) ([]Repository, error) {
	repos := []Repository{}

	for _, d := range discoverers {
		if inUse, err := d.InUse(); !inUse || err != nil {
			continue
		}
		repoMap, err := d.DiscoverRepositories()
		if err != nil {
			return nil, err
		}

		repos = append(repos, reposSlice(repoMap)...)
	}

	return repos, nil
}

func reposSlice(repoMap map[string]Repository) []Repository {
	repos := []Repository{}
	for _, repo := range repoMap {
		repos = append(repos, repo)
	}

	return repos
}
