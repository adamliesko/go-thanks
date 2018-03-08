package discover

import (
	"testing"
)

func TestRepoMapAdditionEmpty(t *testing.T) {
	var rm RepoMap
	rm.add("github.com/owner/repo/pkg")
	if len(rm) != 1 {
		t.Errorf("didn't add package to nil repoMap")
	}
}

func TestRepoMapAdditionDuplicate(t *testing.T) {
	rm := RepoMap{}
	rm.add("github.com/owner/repo/pkg")
	rm.add("github.com/owner/repo/pkg")
	if len(rm) != 1 {
		t.Errorf("didn't add package to repoMap")
	}
}

func TestRepoMapAdditionMultiple(t *testing.T) {
	rm := RepoMap{}
	rm.add("github.com/owner/repo/pkg")
	rm.add("github.com/owner/another_repo/pkg")
	if len(rm) != 2 {
		t.Errorf("didn't add package to repoMap")
	}
}

func TestRepoMapSorting(t *testing.T) {
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
		t.Fatal("missing repos")
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
