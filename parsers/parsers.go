package parsers

import (
	"reflect"
	"strings"

	"github.com/alecbcs/cuppa/results"
	"github.com/alecbcs/cuppa/version"
)

// Parser is a universal parser interface implemented by
// all parsers in Binoc
type Parser interface {
	Decode(content string) (pkg Package, err error)
	Encode(pkg Package) (result string, err error)
}

// Package is a universal package interface for working
// with packages in Binoc.
type Package interface {
	AddVersion(results.Result) (err error)
	GetLatestVersion() (result version.Version)
	GetURL() (result string)
	GetName() (result string)
}

type entry struct {
	FileExt string
	Parser  Parser
}

var (
	// AvailableParsers is a map containing all available parsing engines.
	AvailableParsers map[string]entry
)

func registerParser(parser Parser, fileExt string) {
	if AvailableParsers == nil {
		AvailableParsers = make(map[string]entry)
	}
	name := strings.ToLower(reflect.ValueOf(parser).Type().Name())
	AvailableParsers[name] = entry{FileExt: fileExt, Parser: parser}
}
