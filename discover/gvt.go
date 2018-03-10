package discover

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// Gvt is a discoverer for gvt Go package manager.
type Gvt struct{}

// Name returns name of Gvt discoverer.
func (g Gvt) Name() string {
	return "gvt"
}

// InUse checks whether project is using gvt as a package manager.
func (g Gvt) InUse(projectPath string) (bool, error) {
	_, err := os.Stat(gvtFilePath(projectPath))
	if err != nil {
		return false, err
	}

	return true, nil
}

func gvtFilePath(projectPath string) string {
	return path.Join(projectPath, "vendor", "manifest")
}

// Repositories discovers repositories belonging to packages imported inside a project managed by Gvt.
func (g Gvt) Repositories(projectPath string) (RepoMap, error) {
	list, err := packageListGvt(projectPath)
	if err != nil {
		return nil, err
	}

	repoMap := RepoMap{}
	for _, p := range list.Dependencies {
		repoMap.add(p.Importpath)
	}

	return repoMap, nil
}

type gvtList struct {
	Dependencies []struct {
		Importpath string `json:"importpath"`
	} `json:"dependencies"`
}

func packageListGvt(projectPath string) (gvtList, error) {
	list := gvtList{}

	content, err := ioutil.ReadFile(gvtFilePath(projectPath))
	if err != nil {
		return gvtList{}, err
	}
	err = json.Unmarshal(content, &list)
	if err != nil {
		return gvtList{}, err
	}

	return list, nil
}
