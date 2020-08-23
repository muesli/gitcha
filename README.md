# gitcha

[![Latest Release](https://img.shields.io/github/release/muesli/gitcha.svg)](https://github.com/muesli/gitcha/releases)
[![Build Status](https://github.com/muesli/gitcha/workflows/build/badge.svg)](https://github.com/muesli/gitcha/actions)
[![Coverage Status](https://coveralls.io/repos/github/muesli/gitcha/badge.svg?branch=master)](https://coveralls.io/github/muesli/gitcha?branch=master)
[![Go ReportCard](http://goreportcard.com/badge/muesli/gitcha)](http://goreportcard.com/report/muesli/gitcha)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/muesli/gitcha)

Go helpers to work with git repositories

## Examples

```go
import "github.com/muesli/gitcha"

// returns the directory of the git repository path is a member of:
repo, err := gitcha.GitRepoForPath(path)

// finds files from list in path. It respects all .gitignores it finds while
// traversing paths:
ch, err := gitcha.FindFiles(path, []string{"*.md"})

for v := range ch {
    fmt.Println(v.Path)
}

// just for convenience:
ok := gitcha.IsPathInGit(path)
```
