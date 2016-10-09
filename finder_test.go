package finder

import (
	"fmt"
	"testing"
)

// test runs a slice of test cases.
func test(t *testing.T, cases []testCase) {
	for _, tc := range cases {
		tc.Test(t)
	}
}

type testCase struct {
	finder *Finder
	paths  []string
}

// Test runs a single test case.
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

// matchFiles asserts that only the expected paths occur in the given files.
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

func TestIn(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures"),
			[]string{
				"test-fixtures/1.log",
				"test-fixtures/1.txt",
				"test-fixtures/CVS",
				"test-fixtures/CVS/1",
				"test-fixtures/CVS/.config",
			},
		},
		{
			New().In("test-fixtures/CVS"),
			[]string{
				"test-fixtures/CVS/1",
				"test-fixtures/CVS/.config",
			},
		},
	}
	test(t, tests)
}

func TestPath(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").Path("CVS"),
			[]string{
				"test-fixtures/CVS",
				"test-fixtures/CVS/1",
				"test-fixtures/CVS/.config",
			},
		},
		{
			New().In("test-fixtures").Path("x"),
			[]string{},
		},
		{
			New().In("test-fixtures").Path("*"),
			[]string{
				"test-fixtures/1.log",
				"test-fixtures/1.txt",
				"test-fixtures/CVS",
				"test-fixtures/CVS/1",
				"test-fixtures/CVS/.config",
			},
		},
	}
	test(t, tests)
}

func TestNotPath(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").NotPath("CVS"),
			[]string{
				"test-fixtures/1.log",
				"test-fixtures/1.txt",
			},
		},
	}
	test(t, tests)
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
	test(t, tests)
}

func TestNameRegex(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").NameRegex("\\.log$"),
			[]string{"test-fixtures/1.log"},
		},
	}
	test(t, tests)
}

func TestNotName(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").NotName("*.log").NotName("*.txt"),
			[]string{
				"test-fixtures/CVS/1",
				"test-fixtures/CVS",
				"test-fixtures/CVS/.config",
			},
		},
	}
	test(t, tests)
}

func TestNotNameRegex(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").NotNameRegex("\\.(log|txt)$"),
			[]string{
				"test-fixtures/CVS/1",
				"test-fixtures/CVS/.config",
				"test-fixtures/CVS",
			},
		},
	}
	test(t, tests)
}

func TestFiles(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").Files(),
			[]string{
				"test-fixtures/1.log",
				"test-fixtures/1.txt",
				"test-fixtures/CVS/1",
				"test-fixtures/CVS/.config",
			},
		},
	}
	test(t, tests)
}

func TestDirs(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").Dirs(),
			[]string{
				"test-fixtures/CVS",
			},
		},
	}
	test(t, tests)
}

func TestMinDepth(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").MinDepth(2),
			[]string{
				"test-fixtures/CVS/1",
				"test-fixtures/CVS/.config",
			},
		},
	}
	test(t, tests)
}

func TestMaxDepth(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").MaxDepth(1),
			[]string{
				"test-fixtures/1.log",
				"test-fixtures/1.txt",
				"test-fixtures/CVS",
			},
		},
	}
	test(t, tests)
}

func TestSize(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").Files().Size(func(size int64) bool { return size > 0 }),
			[]string{
				"test-fixtures/1.txt",
			},
		},
	}
	test(t, tests)
}

func TestIgnoreVCS(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").IgnoreVCS(),
			[]string{
				"test-fixtures/1.txt",
				"test-fixtures/1.log",
			},
		},
	}
	test(t, tests)
}

func TestIgnoreDots(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").IgnoreDots(),
			[]string{
				"test-fixtures/1.log",
				"test-fixtures/1.txt",
				"test-fixtures/CVS",
				"test-fixtures/CVS/1",
			},
		},
	}
	test(t, tests)
}
