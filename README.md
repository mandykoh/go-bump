# go-version

Simple utility to generate version constants for Go libraries.


## Usage

Add `go-version` to your `Gopkg.toml` external dependencies for `dep`:

```toml
required = [
  "github.com/mandykoh/go-version"
]
```

Ensure the new dependency is imported and available:

```
$ dep ensure
```

Tag your library if it hasn’t already been tagged:

```
$ git tag 0.1.5
```

Run `go-version` from your library root:

```
$ go run vendor/github.com/mandykoh/go-version/main.go mypackage
```

This will generate a `version.go` with the following constants:

```go
mypackage.VersionMajor     // Major version string (eg "0")
mypackage.VersionMinor     // Minor version string (eg "1")
mypackage.VersionRevision  // Revision string (eg "5")
mypackage.Version          // Full version string (eg "0.1.5")
```

The version is determined by using `git describe --tags` to find the revision relative to the latest Git tag.
 
 To update the version in future, simply run `go generate` from the same location.
 