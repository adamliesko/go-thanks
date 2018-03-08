package dep

import (
	"io/ioutil"
	"os"
	"strings"

	toml "github.com/pelletier/go-toml"

	"github.com/adamliesko/go-thanks/discover"
)

const depFilePath = "Gopkg.toml"

type Discoverer struct{}

func (d Discoverer) InUse() (bool, error) {
	_, err := os.Stat(depFilePath)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d Discoverer) DiscoverRepositories() (map[string]discover.Repository, error) {
	list, err := readPackageList()
	if err != nil {
		return nil, err
	}

	repoMap := map[string]discover.Repository{}
	for _, c := range list.Constraint {
		splits := strings.SplitAfterN(c.Name, "/", 4)
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
	Constraint []struct {
		Name string `toml:"name"`
	} `toml:"constraint"`
}

func readPackageList() (List, error) {
	content, err := ioutil.ReadFile(depFilePath)
	if err != nil {
		return List{}, err
	}

	list := List{}
	err = toml.Unmarshal(content, &list)
	if err != nil {
		return List{}, err
	}

	return list, nil
}
