package discover

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// Govendor is a discoverer for govendor Go package manager.
type Govendor struct{}

// Name returns name of Govendor discoverer.
func (g Govendor) Name() string {
	return "govendor"
}

// InUse checks whether project is using govendor as a package manager.
func (g Govendor) InUse(projectPath string) (bool, error) {
	_, err := os.Stat(govendorFilePath(projectPath))
	if err != nil {
		return false, err
	}

	return true, nil
}

func govendorFilePath(projectPath string) string {
	return path.Join(projectPath, "vendor", "vendor.json")
}

// Repositories discovers repositories belonging to packages imported inside a project managed by Govendor.
func (g Govendor) Repositories(projectPath string) (RepoMap, error) {
	list, err := packageListGovendor(projectPath)
	if err != nil {
		return nil, err
	}

	repoMap := RepoMap{}
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
