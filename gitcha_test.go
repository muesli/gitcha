package gitcha

import (
	"io/ioutil"
	"os"
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

func TestFindAllFiles(t *testing.T) {
	tmp := t.TempDir()

	gitignore, err := os.Create(filepath.Join(tmp, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	defer gitignore.Close()
	_, err = gitignore.WriteString("*.test")
	if err != nil {
		t.Fatal(err)
	}
	tt := []struct {
		path string
		list []string
		exp  string
	}{
		{tmp, []string{"*.test"}, "ignore.test"},
	}

	for _, test := range tt {
		f, err := os.Create(filepath.Join(tmp, test.exp))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		ch, err := FindAllFiles(test.path, test.list)
		if err != nil {
			t.Fatal(err)
		}

		var counter int
		for v := range ch {
			counter++
			if test.exp != v.Info.Name() {
				t.Errorf("Expected %v, got %v for %s", test.exp, v.Path, test.path)
			}
		}

		if counter != 1 {
			t.Errorf("Expected 1 file found, got %d for %s", counter, test.path)
		}
	}
}

func TestFindFiles(t *testing.T) {
	tlink, err := tempLink(".")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tlink)

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
		{tlink, []string{"gitcha.go"}, "gitcha.go"},
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

// tempLink creates a temporary symbolic link pointing to "dest".
func tempLink(dest string) (string, error) {
	// Use tempfile just to create a name and remove it afterwards.
	tmp, err := ioutil.TempFile("", "gitcha_test")
	if err != nil {
		return "", err
	}
	tmp.Close()
	if err := os.Remove(tmp.Name()); err != nil {
		return "", err
	}

	d, err := filepath.Abs(dest)
	if err != nil {
		return "", err
	}
	if err := os.Symlink(d, tmp.Name()); err != nil {
		return "", err
	}
	return tmp.Name(), nil
}
