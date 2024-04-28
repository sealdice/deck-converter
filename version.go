package deckconverter

import "github.com/Masterminds/semver/v3"

// VersionStr is the string representation of version of this software
//
// Following the [SemVer 2.0.0 specification](https://semver.org).
var VersionStr = "1.0.0"

// Version is a parsed SemVer 2.0.0 object
var Version = semver.MustParse(VersionStr)
