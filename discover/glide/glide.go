package glide

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/adamliesko/go-thanks/discover"
	yaml "gx/ipfs/QmNNARAR2ncSEDKGVAjsE77VYAtNE6qMMPEC7hfWMwMdF9/yaml.v2"
)

const glideFilePath = "glide.yaml"

type Discoverer struct{}

func (d Discoverer) InUse() (bool, error) {
	_, err := os.Stat(glideFilePath)
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
	for _, p := range list.Import {
		splits := strings.SplitAfterN(p.Package, "/", 4)
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
	Import []struct {
		Package string `yaml:"package"`
	} `yaml:"import"`
}

func readPackageList() (List, error) {
	content, err := ioutil.ReadFile(glideFilePath)
	if err != nil {
		return List{}, err
	}

	list := List{}
	err = yaml.Unmarshal(content, &list)
	if err != nil {
		return List{}, err
	}

	return list, err
}
