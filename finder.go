package finder

import (
	"os"
	"path/filepath"
)

type Finder struct {
	dirs []string
}

func New() *Finder {
	return &Finder{}
}

func (f *Finder) In(directories ...string) *Finder {
	f.dirs = append(f.dirs, directories...)
	return f
}

func (f *Finder) Each(fn func(Item)) []error {
	var errs []error
	var dir string
	walker := func(path string, info os.FileInfo, err error) error {
		if dir == path {
			return nil
		}
		if err != nil {
			return err
		}
		fn(newItem(info, path))
		return nil
	}
	for _, dir = range f.dirs {
		err := filepath.Walk(dir, walker)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (f *Finder) ToSlice() ([]Item, []error) {
	var l []Item
	errs := f.Each(func(file Item) {
		l = append(l, file)
	})
	return l, errs
}
