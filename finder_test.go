package finder

import (
	"fmt"
	"regexp"
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
	files := tc.finder.ToSlice()
	if tc.finder.Error() != nil {
		t.Fatal(tc.finder.Error())
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
		_, exists := m[file.RelPath()]
		if !exists {
			return fmt.Errorf("unexpected path: %s", file.RelPath())
		}
		m[file.RelPath()] = true
	}
	for path, found := range m {
		if !found {
			return fmt.Errorf("missing path: %s", path)
		}
	}
	return nil
}

func TestFinder_In(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures"),
			[]string{
				"1.log",
				"1.txt",
				"CVS",
				"CVS/1",
				"CVS/.config",
			},
		},
		{
			New().In("test-fixtures/CVS"),
			[]string{
				"1",
				".config",
			},
		},
	}
	test(t, tests)
}

func TestFinder_Path(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").Path("CVS"),
			[]string{
				"CVS",
				"CVS/1",
				"CVS/.config",
			},
		},
		{
			New().In("test-fixtures").Path("x"),
			[]string{},
		},
		{
			New().In("test-fixtures").Path("*"),
			[]string{
				"1.log",
				"1.txt",
				"CVS",
				"CVS/1",
				"CVS/.config",
			},
		},
	}
	test(t, tests)
}

func TestFinder_Error(t *testing.T) {
	f := New().Path("[").Name("[").NameRegex("[")
	if f.Error() == nil {
		t.Fatal("expecting error, got nil")
	}
	re := regexp.MustCompile(`(?s)1\..*2\..*3\.`)
	if !re.MatchString(f.Error().Error()) {
		t.Fatalf("expected output to match regexp %#v; got %#v", re.String(), f.Error().Error())
	}
}

func TestFinder_NotPath(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").NotPath("CVS"),
			[]string{
				"1.log",
				"1.txt",
			},
		},
	}
	test(t, tests)
}

func TestFinder_Name(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").Name("*.log"),
			[]string{"1.log"},
		},
		{
			New().In("test-fixtures").Name("?.log"),
			[]string{"1.log"},
		},
		{
			New().In("test-fixtures").Name("?.{log,txt}"),
			[]string{
				"1.log",
				"1.txt",
			},
		},
	}
	test(t, tests)
}

func TestFinder_NameRegex(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").NameRegex("\\.log$"),
			[]string{"1.log"},
		},
	}
	test(t, tests)
}

func TestFinder_NotName(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").NotName("*.log").NotName("*.txt"),
			[]string{
				"CVS/1",
				"CVS",
				"CVS/.config",
			},
		},
	}
	test(t, tests)
}

func TestFinder_NotNameRegex(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").NotNameRegex("\\.(log|txt)$"),
			[]string{
				"CVS/1",
				"CVS/.config",
				"CVS",
			},
		},
	}
	test(t, tests)
}

func TestFinder_Files(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").Files(),
			[]string{
				"1.log",
				"1.txt",
				"CVS/1",
				"CVS/.config",
			},
		},
	}
	test(t, tests)
}

func TestFinder_Dirs(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").Dirs(),
			[]string{
				"CVS",
			},
		},
	}
	test(t, tests)
}

func TestFinder_Depth(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").Depth(2, -1),
			[]string{
				"CVS/1",
				"CVS/.config",
			},
		}, {
			New().In("test-fixtures").Depth(1, 1),
			[]string{
				"1.log",
				"1.txt",
				"CVS",
			},
		}, {
			New().In("test-fixtures").Depth(0, -1),
			[]string{
				"1.log",
				"1.txt",
				"CVS",
				"CVS/1",
				"CVS/.config",
			},
		},
	}
	test(t, tests)
}

func TestFinder_Size(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").Files().Size(1, 1),
			[]string{
				"1.txt",
			},
		},
		{
			New().In("test-fixtures").Files().Size(0, 0),
			[]string{
				"1.log",
				"CVS/.config",
				"CVS/1",
			},
		},
		{
			New().In("test-fixtures").Files().Size(0, -1),
			[]string{
				"1.log",
				"1.txt",
				"CVS/1",
				"CVS/.config",
			},
		},
	}
	test(t, tests)
}

func TestFinder_IgnoreVCS(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").IgnoreVCS(),
			[]string{
				"1.txt",
				"1.log",
			},
		},
	}
	test(t, tests)
}

func TestFinder_IgnoreDots(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").IgnoreDots(),
			[]string{
				"1.log",
				"1.txt",
				"CVS",
				"CVS/1",
			},
		},
	}
	test(t, tests)
}

func TestFinder_Filter(t *testing.T) {
	tests := []testCase{
		{
			New().In("test-fixtures").Filter(
				func(i Item) bool {
					return i.Depth() == 2
				},
			),
			[]string{
				"CVS/.config",
				"CVS/1",
			},
		},
	}
	test(t, tests)
}
