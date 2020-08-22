package gitcha

import (
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

// SearchResult combines the absolute path of a file with a FileInfo struct.
type SearchResult struct {
	Path string
	Info os.FileInfo
}

// GitRepoForPath returns the directory of the git repository path is a member
// of, or an error.
func GitRepoForPath(path string) (string, error) {
	dir, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	for {
		st, err := os.Stat(filepath.Join(dir, ".git"))
		if err == nil && st.IsDir() {
			return dir, nil
		}

		// reached root?
		if dir == filepath.Dir(dir) {
			return "", nil
		}

		// check parent dir
		dir = filepath.Dir(dir)
	}
}

// IsPathInGit returns true when a path is part of a git repository.
func IsPathInGit(path string) bool {
	p, err := GitRepoForPath(path)
	if err != nil {
		return false
	}

	return len(p) > 0
}

// FindFiles finds files from list in path. It respects all .gitignores it finds
// while traversing paths.
func FindFiles(path string, list []string) (chan SearchResult, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !st.IsDir() {
		return nil, err
	}

	ch := make(chan SearchResult)
	go func() {
		defer close(ch)

		var lastGit string
		var gi *ignore.GitIgnore

		_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			git, _ := GitRepoForPath(path)
			if git != "" && git != path {
				if lastGit != git {
					lastGit = git
					gi, err = ignore.CompileIgnoreFile(filepath.Join(git, ".gitignore"))
				}

				if err == nil && gi != nil && gi.MatchesPath(path) {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}

			for _, v := range list {
				matched := strings.EqualFold(filepath.Base(path), v)
				if !matched {
					matched, _ = filepath.Match(strings.ToLower(v), strings.ToLower(filepath.Base(path)))
				}

				if matched {
					res, _ := filepath.Abs(path)
					ch <- SearchResult{
						Path: res,
						Info: info,
					}

					// only match each path once
					return nil
				}
			}
			return nil
		})
	}()

	return ch, nil
}

func FindFirstFile(path string, list []string) (SearchResult, error) {
	ch, err := FindFiles(path, list)
	if err != nil {
		return SearchResult{}, err
	}

	for v := range ch {
		return v, nil
	}

	return SearchResult{}, nil
}
