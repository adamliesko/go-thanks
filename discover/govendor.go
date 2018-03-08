package discover

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type Govendor struct{}

func (d Govendor) InUse(projectPath string) (bool, error) {
	_, err := os.Stat(govendorFilePath(projectPath))
	if err != nil {
		return false, err
	}

	return true, nil
}

func govendorFilePath(projectPath string) string {
	return path.Join(projectPath, "vendor", "vendor.json")
}

func (d Govendor) DiscoverRepositories(projectPath string) (repoMap, error) {
	list, err := packageListGovendor(projectPath)
	if err != nil {
		return nil, err
	}

	repoMap := repoMap{}
	for _, p := range list.Package {
		repoMap.add(p.Path)
	}

	return repoMap, nil
}

type govendorList struct {
	Package []struct {
		Path string `json:"path"`
	} `json:"package"`
}

func packageListGovendor(projectPath string) (govendorList, error) {
	list := govendorList{}

	content, err := ioutil.ReadFile(govendorFilePath(projectPath))
	if err != nil {
		return govendorList{}, err
	}
	err = json.Unmarshal(content, &list)
	if err != nil {
		return govendorList{}, err
	}

	return list, nil
}
