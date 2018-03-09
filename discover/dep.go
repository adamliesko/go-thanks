package discover

import (
	"io/ioutil"
	"os"
	"path"

	toml "github.com/pelletier/go-toml"
)

const depFilePath = "Gopkg.toml"

// Dep is a discoverer for dep Go package manager.
type Dep struct{}

// Name returns name of Dep discoverer.
func (d Dep) Name() string {
	return "dep"
}

// InUse checks whether project is using dep as a package manager.
func (d Dep) InUse(projectPath string) (bool, error) {
	_, err := os.Stat(path.Join(projectPath, depFilePath))
	if err != nil {
		return false, err
	}

	return true, nil
}

// Repositories discovers repositories belonging to packages imported inside a project managed by Dep.
func (d Dep) Repositories(projectPath string) (RepoMap, error) {
	list, err := packageListDep(projectPath)
	if err != nil {
		return nil, err
	}

	repoMap := RepoMap{}
	for _, c := range list.Constraint {
		repoMap.add(c.Name)
	}

	return repoMap, nil
}

type depList struct {
	Constraint []struct {
		Name string `toml:"name"`
	} `toml:"constraint"`
}

func packageListDep(projectPath string) (depList, error) {
	content, err := ioutil.ReadFile(path.Join(projectPath, depFilePath))
	if err != nil {
		return depList{}, err
	}

	list := depList{}
	err = toml.Unmarshal(content, &list)
	if err != nil {
		return depList{}, err
	}

	return list, nil
}
