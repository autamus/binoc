package parsers

import (
	"github.com/alecbcs/cuppa/results"
	"github.com/alecbcs/cuppa/version"
)

// Parser is a universal parser interface implemented by
// all parsers in Binoc
type Parser interface {
	Decode(path string) (pkg Package, err error)
	Encode(path string, pkg Package) (err error)
}

// Package is a universal package interface for working
// with packages in Binoc.
type Package interface {
	AddVersion(results.Result) (err error)
	GetLatestVersion() (result version.Version)
}
