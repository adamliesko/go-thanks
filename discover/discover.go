// Package discover host all the discovereres which can explore user project and identify package manager in use and
// discovers Go packages used, with respective repositories.
package discover

// Discoverer can explore the workspace and discover Go packages that are used in project by specific package manager.
type Discoverer interface {
	Name() string
	InUse(path string) (bool, error)
	Repositories(path string) (RepoMap, error)
}

// Repositories produces a slice of repositories extracted from the passed in discoverers within one's Go project.
func Repositories(path string) ([]Repository, error) {
	discoverers := []Discoverer{Dep{}, Glide{}, Govendor{}}

	repoMap := make(RepoMap)
	for _, d := range discoverers {
		if inUse, err := d.InUse(path); !inUse || err != nil {
			continue
		}

		rs, err := d.Repositories(path)
		if err != nil {
			return nil, err
		}

		for k, v := range rs {
			if _, ok := repoMap[k]; ok {
				continue
			}
			repoMap[k] = v
		}

	}

	// predictive sorting of resulting repos
	repos := repoMap.toSortedSlice()

	return repos, nil
}
