package discover

import (
	"sort"
	"strings"
)

// Repository represents a VCS repository hosted on a remote site, containing one or more Go packages.
type Repository struct {
	Name  string
	Owner string
	URL   string
}

type RepoMap map[string]Repository

func (rm RepoMap) add(pkgPath string) {
	if strings.HasPrefix(pkgPath, "gopkg.in") {
		pkgPath = resolveGoPkgIn(pkgPath)
	}
	splits := strings.SplitN(pkgPath, "/", 4)
	if len(splits) < 3 {
		return
	}
	repoURL := strings.Join(splits[0:3], "/")

	if _, ok := rm[repoURL]; !ok {
		repo := Repository{
			Owner: splits[1],
			Name:  splits[2],
			URL:   repoURL,
		}
		rm[repoURL] = repo
	}
}

// gopkg.in/pkg.v3      → github.com/go-pkg/pkg
// gopkg.in/user/pkg.v3 → github.com/user/pkg
func resolveGoPkgIn(pkgPath string) string {
	switch strings.Count(pkgPath, "/") {
	case 1:
		splits := strings.SplitN(pkgPath, "/", 2)
		dotIndex := strings.LastIndex(splits[1], ".")
		return "github.com/go-pkg/" + splits[1][:dotIndex]
	case 2:
		replaced := strings.Replace(pkgPath, "gopkg.in", "github.com", 1)
		dotIndex := strings.LastIndex(replaced, ".")
		return replaced[:dotIndex]
	default:
		return pkgPath
	}
}

func (rm RepoMap) toSortedSlice() []Repository {
	repos := []Repository{}
	for _, repo := range rm {
		repos = append(repos, repo)
	}
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].URL < repos[j].URL
	})
	return repos
}
