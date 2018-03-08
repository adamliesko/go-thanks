package govendor

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/adamliesko/go-thanks/discover"
)

type Discoverer struct{}

func (d Discoverer) InUse() (bool, error) {
	_, err := os.Stat(vendorFilePath())
	if err != nil {
		return false, err
	}

	return true, nil
}

func vendorFilePath() string {
	return path.Join("vendor", "vendor.json")
}

func (d Discoverer) DiscoverRepositories() (map[string]discover.Repository, error) {
	list, err := readPackageList()
	if err != nil {
		return nil, err
	}

	repoMap := map[string]discover.Repository{}
	for _, p := range list.Package {
		splits := strings.SplitAfterN(p.Path, "/", 4)
		// TODO: extract logic out of govendor and deal with gopkgin
		if len(splits) < 3 {
			continue
		}
		repoURL := strings.Join(splits[0:3], "")

		if _, ok := repoMap[repoURL]; !ok {
			repoMap[repoURL] = discover.Repository{Owner: strings.TrimSuffix(splits[1], "/"), Name: strings.TrimSuffix(splits[2], "/"), URL: repoURL}
		}
	}

	return repoMap, nil
}

type List struct {
	Package []struct {
		Path string `json:"path"`
	} `json:"package"`
}

func readPackageList() (List, error) {
	list := List{}

	content, err := ioutil.ReadFile(vendorFilePath())
	if err != nil {
		return List{}, err
	}
	err = json.Unmarshal(content, &list)
	if err != nil {
		return List{}, err
	}

	return list, nil
}
