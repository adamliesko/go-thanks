package discover

import (
	"sort"
	"strings"
)

// Repository represents a VCS repository hosted on a remote site, containing on or more Go packages.
type Repository struct {
	Name  string
	Owner string
	URL   string
}

// Strings implements Stringer interface for Go.
func (r Repository) String() string {
	return r.URL
}

type repoMap map[string]Repository

func (rm repoMap) add(pkgPath string) {
	if rm == nil {
		rm = repoMap{}
	}

	splits := strings.SplitAfterN(pkgPath, "/", 4)
	if len(splits) < 3 {
		return
	}
	repoURL := strings.Join(splits[0:3], "")

	if _, ok := rm[repoURL]; !ok {
		repo := Repository{
			Owner: strings.TrimSuffix(splits[1], "/"),
			Name:  strings.TrimSuffix(splits[2], "/"),
			URL:   repoURL,
		}
		rm[repoURL] = repo
	}
}

func (rm repoMap) toSortedSlice() []Repository {
	repos := []Repository{}
	for _, repo := range rm {
		repos = append(repos, repo)
	}
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].URL < repos[j].URL
	})
	return repos

}
