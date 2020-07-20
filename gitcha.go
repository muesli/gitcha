package gitcha

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func IsPathInGit(path string) bool {
	p, err := GitRepoForPath(path)
	if err != nil {
		return false
	}

	return len(p) > 0
}

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

		if dir == filepath.Dir(dir) {
			// reached root
			return "", nil
		}
		dir = filepath.Dir(dir)
	}
}

func FindFirstInList(path string, list []string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	st, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !st.IsDir() {
		return "", errors.New("not a directory")
	}

	var res string
	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		for _, v := range list {
			matched := strings.EqualFold(filepath.Base(path), v)
			if !matched {
				matched, _ = filepath.Match(strings.ToLower(v), strings.ToLower(filepath.Base(path)))
			}

			if matched {
				res, _ = filepath.Abs(path)

				// abort filepath.Walk
				return errors.New("source found")
			}
		}
		return nil
	})

	if res != "" {
		return res, nil
	}

	return "", errors.New("none found")
}

func FindFileFromList(path string, list []string) chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)

		path, err := filepath.Abs(path)
		if err != nil {
			return
		}
		st, err := os.Stat(path)
		if err != nil {
			return
		}
		if !st.IsDir() {
			return
		}

		var res string
		_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			for _, v := range list {
				matched := strings.EqualFold(filepath.Base(path), v)
				if !matched {
					matched, _ = filepath.Match(strings.ToLower(v), strings.ToLower(filepath.Base(path)))
				}

				if matched {
					res, _ = filepath.Abs(path)
					ch <- res

					// only match each path once
					continue
				}
			}
			return nil
		})
	}()

	return ch
}
