package finder

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestItemSlice_Size(t *testing.T) {
	sl, err := New().In("test-fixtures").ToSlice()
	if err != nil {
		t.Fatal(err)
	}
	if sl.Size() != 1 {
		t.Fatalf("expecting size 1, got %d", sl.Size())
	}
}

func testSort(t *testing.T, sortFn func(Item, Item) bool) ItemSlice {
	sl, err := New().In("test-fixtures").ToSlice()
	if err != nil {
		t.Fatal(err)
		return nil
	}
	sl.Sort(sortFn)
	return sl
}

func TestItemSlice_Sort_BySize(t *testing.T) {
	sl := testSort(t, BySize)
	if sl[0].Size() >= sl[len(sl)-1].Size() {
		t.Fatalf("expecting order to ascend by size; got first = %#v; last = %#v", sl[0].Size(), sl[len(sl)-1].Size())
	}
}

func TestItemSlice_Sort_ByModified(t *testing.T) {
	sl := testSort(t, ByModified)
	if sl[0].ModTime().After(sl[len(sl)-1].ModTime()) {
		t.Fatalf("expecting order to ascend by modification time; got first = %#v; last = %#v", sl[0].ModTime(), sl[len(sl)-1].ModTime())
	}
}

func TestItemSlice_Sort_ByExtension(t *testing.T) {
	sl := testSort(t, ByExtension)
	first := filepath.Ext(sl[0].Name())
	last := filepath.Ext(sl[len(sl)-1].Name())
	if first >= last {
		t.Fatalf("expecting order to ascend by extension; got first = %#v; last = %#v", first, last)
	}
}

func TestItemSlice_Sort_ByPath(t *testing.T) {
	sl := testSort(t, ByPath)
	if sl[0].Path() >= sl[len(sl)-1].Path() {
		t.Fatalf("expecting order to ascend by path; got first = %#v; last = %#v", sl[0].Path(), sl[len(sl)-1].Path())
	}
}

func TestItemSlice_Sort_ByName(t *testing.T) {
	sl := testSort(t, ByName)
	if sl[0].Name() >= sl[len(sl)-1].Name() {
		t.Fatalf("expecting order to ascend by name; got first = %#v; last = %#v", sl[0].Name(), sl[len(sl)-1].Name())
	}
}
