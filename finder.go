package finder

import (
	"os"
	"path/filepath"
	"regexp"
)

type Finder struct {
	dirs        []string
	names       []Matcher
	setupErrors []error
}

type Matcher func(Item) bool

func New() *Finder {
	return &Finder{}
}

func (f *Finder) In(directories ...string) *Finder {
	f.dirs = append(f.dirs, directories...)
	return f
}

func (f *Finder) Name(n string) *Finder {
	_, err := filepath.Match(n, "foo")
	if err != nil {
		f.setupErrors = append(f.setupErrors, err)
		return f
	}
	f.names = append(f.names, func(i Item) bool {
		ok, _ := filepath.Match(n, i.Name())
		return ok
	})
	return f
}

func (f *Finder) NameRegex(n string) *Finder {
	re, err := regexp.Compile(n)
	if err != nil {
		f.setupErrors = append(f.setupErrors, err)
		return f
	}
	f.names = append(f.names, func(i Item) bool {
		return re.MatchString(i.Name())
	})
	return f
}

func (f *Finder) match(i Item) bool {
	match := true
	if len(f.names) > 0 {
		match = false
		for _, n := range f.names {
			if n(i) {
				match = true
				break
			}
		}
	}
	return match
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
		item := newItem(info, path)
		if f.match(item) {
			fn(item)
		}
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
