package govendor

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"strings"

	"path"

	"github.com/adamliesko/go-thanks/discover"
)

type List struct {
	Package []struct {
		Path string `json:"path"`
	} `json:"package"`
}

type Discoverer struct{}

func (d Discoverer) InUse() (bool, error) {
	_, err := os.Stat(vendorFilePath())
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d Discoverer) DiscoverRepositories() ([]discover.Repository, error) {
	content, err := ioutil.ReadFile(vendorFilePath())
	if err != nil {
		return nil, err
	}
	list := List{}
	if err := json.Unmarshal(content, &list); err != nil {
		return nil, err
	}

	repoMap := map[string]discover.Repository{}
	for _, p := range list.Package {
		splits := strings.SplitAfterN(p.Path, "/", 4)
		// TODO: extract logic out of goendor and deal with gopkgin
		if len(splits) < 3 {
			continue
		}
		repoURL := strings.Join(splits[0:3], "")

		if _, ok := repoMap[repoURL]; !ok {
			repoMap[repoURL] = discover.Repository{Owner: strings.TrimSuffix(splits[1], "/"), Name: strings.TrimSuffix(splits[2], "/"), URL: repoURL}
		}
	}

	repos := []discover.Repository{}
	for _, repo := range repoMap {
		repos = append(repos, repo)
	}

	return repos, nil
}

func vendorFilePath() string {
	return path.Join("vendor", "vendor.json")
}
