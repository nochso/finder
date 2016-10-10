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
