# go-thanks
[![Build Status](https://secure.travis-ci.org/adamliesko/go-thanks.svg)](http://travis-ci.org/adamliesko/go-thanks)
[![Go Report Card](https://goreportcard.com/badge/github.com/adamliesko/go-thanks)](https://goreportcard.com/report/github.com/adamliesko/go-thanks)
[![GoDoc](https://godoc.org/github.com/adamliesko/go-thanks?status.svg)](https://godoc.org/github.com/adamliesko/go-thanks) 
[![Coverage Status](https://img.shields.io/coveralls/adamliesko/go-thanks.svg)](https://coveralls.io/r/adamliesko/go-thanks?branch=master)

`go-thanks` is a cmd line utility to show some love to all the hardworking developers and organzaitions, from whose work you profit daily by using their OSS.
It  automatically detects imported packages from your package managers ([dep](https://github.com/golang/dep), [Govendor](https://github.com/kardianos/govendor) and [Glide](https://github.com/Masterminds/glide) and stars repositories of imported packages.
on Github and Gitlab.

### Installation
```
go get -u github.com/adamliesko/go-thanks
```

### Usage

```
go-thanks --github-token GITHUB_TOKEN
``` 

As an alternative, a project path can be specified
by adding `--path PATH` argument. If no tokens are provided from the command line, `go-thanks` falls back to reading
respective variables.

```
Usage of ./go-thanks:
  -github-token string
    	Github API token. Defaults to env variable GITHUB_TOKEN.
  -gitlab-token string
    	Gitlab API token. Defaults to env variable GITLAB_TOKEN.
  -path string
    	Path to Go project. (default ".")

```

### Access Tokens

`go-thanks` requires personal access tokens, to be able to perform the thank action (starring a repository).

For Github follow their [creating-a-personal-access-token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/) 
guide and check only `public_repo` access.

For Gitlab follow their [personal_access_tokens](https://docs.gitlab.com/ce/user/profile/personal_access_tokens.html) 
guide and use scope `api`.
