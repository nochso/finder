package finder

import "testing"

func TestItemSlice_Size(t *testing.T) {
	sl, err := New().In("test-fixtures").ToSlice()
	if err != nil {
		t.Fatal(err)
	}
	if sl.Size() != 1 {
		t.Fatalf("expecting size 1, got %d", sl.Size())
	}
}

func TestItemSlice_Sort(t *testing.T) {
	sl, err := New().In("test-fixtures").ToSlice()
	if err != nil {
		t.Fatal(err)
	}
	sl.Sort(func(i, j Item) bool {
		return i.Size() < j.Size()
	})
	if sl[0].Size() >= sl[len(sl)-1].Size() {
		t.Fatalf("expecting order to ascend by size; got first = %#v; last = %#v", sl[0].Size(), sl[len(sl)-1].Size())
	}
}
