package discover

import (
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

const glideFilePath = "glide.yaml"

type Glide struct{}

func (d Glide) InUse(projectPath string) (bool, error) {
	_, err := os.Stat(path.Join(projectPath, glideFilePath))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d Glide) DiscoverRepositories(projectPath string) (repoMap, error) {
	list, err := packageListGlide(projectPath)
	if err != nil {
		return nil, err
	}

	repoMap := repoMap{}
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
