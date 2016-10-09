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

func TestName(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").Name("*.log"),
			[]string{"test-fixtures/1.log"},
		},
		{
			New().In("test-fixtures").Name("?.log"),
			[]string{"test-fixtures/1.log"},
		},
		{
			New().In("test-fixtures").Name("?.{log,txt}"),
			[]string{
				"test-fixtures/1.log",
				"test-fixtures/1.txt",
			},
		},
	}
	for _, tc := range tests {
		tc.Test(t)
	}
}

func TestNameRegex(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").NameRegex("\\.log$"),
			[]string{"test-fixtures/1.log"},
		},
	}
	for _, tc := range tests {
		tc.Test(t)
	}
}

func TestNotName(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").NotName("*.log").NotName("*.txt"),
			[]string{
				"test-fixtures/d1/1",
				"test-fixtures/d1",
			},
		},
	}
	for _, tc := range tests {
		tc.Test(t)
	}
}

func TestNotNameRegex(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").NotNameRegex("\\.(log|txt)$"),
			[]string{
				"test-fixtures/d1/1",
				"test-fixtures/d1",
			},
		},
	}
	for _, tc := range tests {
		tc.Test(t)
	}
}

func TestFiles(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").Files(),
			[]string{
				"test-fixtures/1.log",
				"test-fixtures/1.txt",
				"test-fixtures/d1/1",
			},
		},
	}
	for _, tc := range tests {
		tc.Test(t)
	}
}

func TestDirs(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").Dirs(),
			[]string{
				"test-fixtures/d1",
			},
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
