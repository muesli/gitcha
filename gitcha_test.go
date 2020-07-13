package gitcha

import (
	"path/filepath"
	"testing"
)

func TestIsPathInGit(t *testing.T) {
	tt := []struct {
		path string
		exp  bool
	}{
		{"/", false},
		{".", false},
		{"gitcha.go", true},
	}

	for _, test := range tt {
		r, err := IsPathInGit(test.path)
		if err != nil {
			t.Error(err)
		}

		if test.exp != r {
			t.Errorf("Expected %v, got %v for %s", test.exp, r, test.path)
		}
	}
}

func TestFindFirstInList(t *testing.T) {
	tt := []struct {
		path   string
		list   []string
		exp    string
		expErr bool
	}{
		{"../", []string{"gitcha.go"}, "gitcha.go", false},
		{".", []string{"gitcha_test.go"}, "gitcha_test.go", false},
		{".", []string{"exist.not"}, "", true},
	}

	for _, test := range tt {
		r, err := FindFirstInList(test.path, test.list)
		if err != nil && !test.expErr {
			t.Error(err)
		}
		if err == nil && test.expErr {
			t.Errorf("Expected error, got none for %s", test.path)
		}

		if err != nil && test.expErr {
			continue
		}

		test.exp, err = filepath.Abs(test.exp)
		if err != nil {
			t.Fatal(err)
		}
		if test.exp != r {
			t.Errorf("Expected %v, got %v for %s", test.exp, r, test.path)
		}
	}
}
