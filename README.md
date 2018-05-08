# go-bump

Simple utility to generate and maintain version constants for Go libraries.


## Usage

Add `go-bump` to your `Gopkg.toml` external dependencies for `dep`:

```toml
required = [
  "github.com/mandykoh/go-bump"
]
```

Ensure the new dependency is imported and available:

```
$ dep ensure
```

Run `go-bump` from your library root:

```
$ go run vendor/github.com/mandykoh/go-bump/main.go mypackage
```

This will generate a `version.go` with the following constants:

```go
mypackage.VersionMajor     // Major version string (eg "0")
mypackage.VersionMinor     // Minor version string (eg "1")
mypackage.VersionRevision  // Revision string (eg "5")
mypackage.Version          // Full version string (eg "0.1.5")
```

The version is determined by using `git describe --tags` to find the revision relative to the latest Git tag.

Tag your library with the version:

```
$ git commit -m "Bump version." version.go
$ git tag 0.1.5
```

To update the version in future, simply run `go generate` from the same location. The revision will be automatically bumped based on the nearest Git tag, so re-tag with the new version.

```
$ go generate
Generated version 0.1.6 for package mypackage
$ git commit -m "Bump version." version.go
$ git tag 0.1.6
```

To bump the major or minor version numbers, specify the new full version when running `go-bump`:

```
$ go run vendor/github.com/mandykoh/go-bump/main.go mypackage 0.2.0
Generated version 0.2.0 for package mypackage
$ git commit -m "Bump version." version.go
$ git tag 0.2.0
```
