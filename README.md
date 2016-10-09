`github.com/nochso/finder`
==========================

[![Released under MIT license](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GitHub release](https://img.shields.io/github/release/nochso/finder.svg?maxAge=2592000)](https://github.com/nochso/finder/releases)
[![SemVer 2.0.0](https://img.shields.io/badge/SemVer-2.0.0-blue.svg)][semver]

A fluent interface around [path/filepath.Walk].

- Inspired by [symfony/finder]
- Single external dependency on [gobwas/glob] for glob path matching

Installation
------------

```
go get github.com/nochso/finder
```

Usage
-----
See [documentation on godoc][godoc].

To do
-----

- Limit by folder depth
- Limit by modification date
- Limit by file content
- Limit by file size
- Ignore dot files
- Ignore common VCS folders
- Sort by..?
- Filter using user supplied func
- Command line utility as sane replacement for `find`

Versioning
----------

This package adheres to [semantic versioning 2.0.0][semver].

*TL;DR* If you use version `1.*` you should never have problems using future
`1.*` versions. Only a major release e.g. `2.0.0` will break backwards
compatibility.

License
-------

This package is released under the MIT license. See [LICENSE] for the full
license.

[path/filepath.Walk]: https://golang.org/pkg/path/filepath/#Walk
[symfony/finder]: https://symfony.com/doc/current/components/finder.html
[gobwas/glob]: https://github.com/gobwas/glob
[LICENSE]: LICENSE
[semver]: http://semver.org/spec/v2.0.0.html
[godoc]: https://godoc.org/github.com/nochso/finder