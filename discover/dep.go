package discover

import (
	"io/ioutil"
	"os"
	"path"

	toml "github.com/pelletier/go-toml"
)

const depFilePath = "Gopkg.toml"

type Dep struct{}

func (d Dep) Name() string {
	return "dep"
}

func (d Dep) InUse(projectPath string) (bool, error) {
	_, err := os.Stat(path.Join(projectPath, depFilePath))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d Dep) DiscoverRepositories(projectPath string) (RepoMap, error) {
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
