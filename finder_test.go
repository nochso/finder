package finder

import (
	"fmt"
	"testing"
)

type testCase struct {
	finder *Finder
	paths  []string
}

func TestIn(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures"),
			[]string{
				"test-fixtures/1.log",
				"test-fixtures/1.txt",
				"test-fixtures/d1",
				"test-fixtures/d1/1",
			},
		},
		{
			New().In("test-fixtures/d1"),
			[]string{"test-fixtures/d1/1"},
		},
	}
	for _, tc := range tests {
		tc.Test(t)
	}
}

func (tc testCase) Test(t *testing.T) {
	files, errs := tc.finder.ToSlice()
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	err := matchFiles(t, tc.paths, files)
	if err != nil {
		t.Fatal(err)
	}
}

func matchFiles(t *testing.T, paths []string, files []Item) error {
	m := make(map[string]bool)
	for _, path := range paths {
		m[path] = false
	}
	for _, file := range files {
		_, exists := m[file.Path()]
		if !exists {
			return fmt.Errorf("unexpected path: %s", file.Path())
		}
		m[file.Path()] = true
	}
	for path, found := range m {
		if !found {
			return fmt.Errorf("missing path: %s", path)
		}
	}
	return nil
}
