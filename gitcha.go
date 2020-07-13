package gitcha

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func IsPathInGit(path string) (bool, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}

	for {
		dir := filepath.Dir(absPath)

		st, err := os.Stat(filepath.Join(dir, ".git"))
		if err == nil && st.IsDir() {
			return true, nil
		}

		if dir == absPath {
			// reached root
			return false, nil
		}
		absPath = dir
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
			if strings.EqualFold(filepath.Base(path), v) {
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
