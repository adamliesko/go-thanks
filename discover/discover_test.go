package discover

import (
	"errors"
	"io/ioutil"
	"log"
	"path"
	"reflect"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestDiscoverRepositories(t *testing.T) {
	tcs := []struct {
		name      string
		path      string
		wantRepos []Repository
		wantErr   error
	}{
		{
			name: "use dep",
			path: path.Join("test", "project_dep"),
			wantRepos: []Repository{
				{
					Name:  "semver",
					Owner: "Masterminds",
					URL:   "github.com/Masterminds/semver",
				},
				{
					Name:  "vcs",
					Owner: "Masterminds",
					URL:   "github.com/Masterminds/vcs",
				},
				{
					Name:  "go-toml",
					Owner: "pelletier",
					URL:   "github.com/pelletier/go-toml",
				},
				{
					Name:  "go-toml",
					Owner: "pelletier",
					URL:   "gitlab.com/pelletier/go-toml",
				},
			},
			wantErr: nil,
		},
		{
			name: "use glide",
			path: path.Join("test", "project_glide"),
			wantRepos: []Repository{
				{
					Name:  "vcs",
					Owner: "Masterminds",
					URL:   "github.com/Masterminds/vcs",
				},
				{
					Name:  "cli",
					Owner: "codegangsta",
					URL:   "github.com/codegangsta/cli",
				},
				{
					Name:  "yaml",
					Owner: "go-yaml",
					URL:   "github.com/go-yaml/yaml",
				},
			}, wantErr: nil,
		},
		{
			name: "use govendor",
			path: path.Join("test", "project_govendor"),
			wantRepos: []Repository{
				{
					Name:  "prompt",
					Owner: "Bowery",
					URL:   "github.com/Bowery/prompt",
				},
				{
					Name:  "safefile",
					Owner: "dchest",
					URL:   "github.com/dchest/safefile",
				},
			},
			wantErr: nil,
		},
		{
			name: "use dep and glide - crazy people right",
			path: path.Join("test", "project_dep_glide"),
			wantRepos: []Repository{
				{
					Name:  "semver",
					Owner: "Masterminds",
					URL:   "github.com/Masterminds/semver",
				},
				{
					Name:  "vcs",
					Owner: "Masterminds",
					URL:   "github.com/Masterminds/vcs",
				},
				{
					Name:  "cli",
					Owner: "codegangsta",
					URL:   "github.com/codegangsta/cli",
				},
				{
					Name:  "yaml",
					Owner: "go-yaml",
					URL:   "github.com/go-yaml/yaml",
				},
				{
					Name:  "go-toml",
					Owner: "pelletier",
					URL:   "github.com/pelletier/go-toml",
				},
			},
			wantErr: nil,
		},
		{
			name:      "nothing to be found",
			path:      path.Join("test", "project_none"),
			wantRepos: []Repository{},
			wantErr:   nil,
		},
		{
			name:    "errors from all discoverers",
			path:    path.Join("test", "project_all_errors"),
			wantErr: errors.New(""),
		},
		{
			name:    "errors from govendor",
			path:    path.Join("test", "project_errors_govendor"),
			wantErr: errors.New(""),
		},
		{
			name:    "errors from all except dep still errors out",
			path:    path.Join("test", "project_errors_dep_ok"),
			wantErr: errors.New(""),
		},
	}

	for _, tc := range tcs {
		repos, err := Repositories(tc.path)
		if (err == nil) != (tc.wantErr == nil) {
			t.Errorf("%s: error missmatch, got %v want %v", tc.name, err, tc.wantErr)
		}
		if len(repos) != len(tc.wantRepos) {
			t.Errorf("%s: bad repo count, got %d want %d", tc.name, len(repos), len(tc.wantRepos))
		}
		if !reflect.DeepEqual(repos, tc.wantRepos) {
			t.Errorf("%s: bad repos, got %v want %v", tc.name, repos, tc.wantRepos)
		}
	}
}
