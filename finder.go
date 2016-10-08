package finder

import (
	"os"
	"path/filepath"
	"regexp"
)

type itemType uint8

const (
	typeAll itemType = iota
	typeFile
	typeDir
)

type Finder struct {
	dirs        []string
	names       []Matcher
	notNames    []Matcher
	setupErrors []error
	itype       itemType
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
	matcher := f.name(n)
	if matcher != nil {
		f.names = append(f.names, matcher)
	}
	return f
}

func (f *Finder) name(n string) Matcher {
	_, err := filepath.Match(n, "foo")
	if err != nil {
		f.setupErrors = append(f.setupErrors, err)
		return nil
	}
	return func(i Item) bool {
		ok, _ := filepath.Match(n, i.Name())
		return ok
	}
}

func (f *Finder) NameRegex(n string) *Finder {
	matcher := f.nameRegex(n)
	if matcher != nil {
		f.names = append(f.names, matcher)
	}
	return f
}

func (f *Finder) nameRegex(n string) Matcher {
	re, err := regexp.Compile(n)
	if err != nil {
		f.setupErrors = append(f.setupErrors, err)
		return nil
	}
	return func(i Item) bool {
		return re.MatchString(i.Name())
	}
}

func (f *Finder) NotName(n string) *Finder {
	matcher := f.name(n)
	if matcher != nil {
		f.notNames = append(f.notNames, matcher)
	}
	return f
}

func (f *Finder) NotNameRegex(n string) *Finder {
	matcher := f.nameRegex(n)
	if matcher != nil {
		f.notNames = append(f.notNames, matcher)
	}
	return f
}

func (f *Finder) Files() *Finder {
	f.itype = typeFile
	return f
}

func (f *Finder) Dirs() *Finder {
	f.itype = typeDir
	return f
}

func (f *Finder) match(i Item) bool {
	if (f.itype == typeDir && !i.IsDir()) || (f.itype == typeFile && i.IsDir()) {
		return false
	}
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
	if !match {
		return match
	}
	if len(f.notNames) > 0 {
		for _, matcher := range f.notNames {
			if matcher(i) {
				match = false
				continue
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
	return append(f.setupErrors, errs...)
}

func (f *Finder) ToSlice() ([]Item, []error) {
	var l []Item
	errs := f.Each(func(file Item) {
		l = append(l, file)
	})
	return l, errs
}
