package discover

import (
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

const glideFilePath = "glide.yaml"

// Glide is a discoverer for glide Go package manager.
type Glide struct{}

// Name returns name of Glide discoverer.
func (g Glide) Name() string {
	return "glide"
}

// InUse checks whether project is using glide as a package manager.
func (g Glide) InUse(projectPath string) (bool, error) {
	_, err := os.Stat(path.Join(projectPath, glideFilePath))
	if err != nil {
		return false, err
	}

	return true, nil
}

// Repositories discovers repositories belonging to packages imported inside a project managed by Glide.
func (g Glide) Repositories(projectPath string) (RepoMap, error) {
	list, err := packageListGlide(projectPath)
	if err != nil {
		return nil, err
	}

	repoMap := RepoMap{}
	for _, p := range list.Import {
		repoMap.add(p.Package)
	}

	return repoMap, nil
}

type glideList struct {
	Import []struct {
		Package string `yaml:"package"`
	} `yaml:"import"`
}

func packageListGlide(projectPath string) (glideList, error) {
	content, err := ioutil.ReadFile(path.Join(projectPath, glideFilePath))
	if err != nil {
		return glideList{}, err
	}

	list := glideList{}
	err = yaml.Unmarshal(content, &list)
	if err != nil {
		return glideList{}, err
	}

	return list, err
}
