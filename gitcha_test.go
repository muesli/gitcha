package gitcha

import (
	"path/filepath"
	"testing"
)

func TestGitRepoForPath(t *testing.T) {
	abs, _ := filepath.Abs(".")
	tt := []struct {
		path string
		exp  string
	}{
		{"/", ""},
		{".", abs},
		{"gitcha.go", abs},
	}

	for _, test := range tt {
		r, err := GitRepoForPath(test.path)
		if err != nil {
			t.Error(err)
		}

		if test.exp != r {
			t.Errorf("Expected %v, got %v for %s", test.exp, r, test.path)
		}
	}
}

func TestFindFiles(t *testing.T) {
	tt := []struct {
		path string
		list []string
		exp  string
	}{
		{"../", []string{"gitcha.go"}, "gitcha.go"},
		{".", []string{"gitcha_test.go"}, "gitcha_test.go"},
		{".", []string{"README.MD"}, "README.md"},
		{".", []string{"*.md"}, "README.md"},
		{".", []string{"*.MD"}, "README.md"},
	}

	for _, test := range tt {
		ch, err := FindFiles(test.path, test.list)
		if err != nil {
			t.Fatal(err)
		}

		for v := range ch {
			var err error
			test.exp, err = filepath.Abs(test.exp)
			if err != nil {
				t.Fatal(err)
			}
			if test.exp != v.Path {
				t.Errorf("Expected %v, got %v for %s", test.exp, v, test.path)
			}
		}
	}
}

func TestFindFirstFile(t *testing.T) {
	tt := []struct {
		path   string
		list   []string
		exp    string
		expErr bool
	}{
		{"../", []string{"gitcha.go"}, "gitcha.go", false},
		{".", []string{"gitcha_test.go"}, "gitcha_test.go", false},
		{".", []string{"README.MD"}, "README.md", false},
		{".", []string{"*.md"}, "README.md", false},
		{".", []string{"*.MD"}, "README.md", false},
	}

	for _, test := range tt {
		r, err := FindFirstFile(test.path, test.list)
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
		if test.exp != r.Path {
			t.Errorf("Expected %v, got %v for %s", test.exp, r, test.path)
		}
	}
}
