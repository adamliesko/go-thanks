package discover

import (
	"testing"
)

func TestRepoMapAdditionDuplicate(t *testing.T) {
	t.Parallel()

	rm := RepoMap{}
	rm.add("github.com/owner/repo/pkg")
	rm.add("github.com/owner/repo/pkg")
	if len(rm) != 1 {
		t.Errorf("didn't add package to repoMap")
	}
}

func TestRepoMapAdditionGoPkgIn(t *testing.T) {
	t.Parallel()

	rm := RepoMap{}
	rm.add("gopkg.in/yaml.v3")
	rm.add("gopkg.in/user/hope.v3")

	ss := rm.toSortedSlice()
	if len(ss) != 2 {
		t.Fatal("bad repos count")
	}
	want0 := "github.com/go-yaml/yaml"
	if ss[0].URL != want0 {
		t.Errorf("bad url on index 0, got %s want %s", ss[0].URL, want0)
	}
	want1 := "github.com/user/hope"
	if ss[1].URL != want1 {
		t.Errorf("bad url on index 1, got %s want %s", ss[1].URL, want1)
	}

}

func TestRepoMapAdditionMultiple(t *testing.T) {
	t.Parallel()

	rm := RepoMap{}
	rm.add("github.com/owner/repo/pkg")
	rm.add("github.com/owner/another_repo/pkg")
	if len(rm) != 2 {
		t.Errorf("didn't add package to repoMap")
	}
}

func TestRepoMapSorting(t *testing.T) {
	t.Parallel()

	rm := RepoMap{
		"github.com/owner/xrepo": Repository{
			URL: "github.com/owner/xrepo",
		},
		"github.com/owner/arepo": Repository{
			URL: "github.com/owner/arepo",
		},
		"github.com/owner/repo/pkg": Repository{
			URL: "github.com/owner/repo",
		},
	}

	ss := rm.toSortedSlice()
	if len(ss) != 3 {
		t.Fatal("repos count mismatch")
	}
	if ss[0].URL != "github.com/owner/arepo" {
		t.Error("bad order on 0 index", ss)
	}
	if ss[1].URL != "github.com/owner/repo" {
		t.Error("bad order on 1 index", ss)
	}
	if ss[2].URL != "github.com/owner/xrepo" {
		t.Error("bad order on 2 index", ss)
	}
}
