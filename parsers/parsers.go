package parsers

import (
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/DataDrake/cuppa/results"
)

// Parser is a universal parser interface implemented by
// all parsers in Binoc
type Parser interface {
	Decode(content string, modified time.Time) (pkg Package, err error)
	Encode(pkg Package) (result string, err error)
}

// Package is a universal package interface for working
// with packages in Binoc.
type Package interface {
	AddVersion(results.Result) (err error)
	GetLatestVersion() (result results.Result)
	GetAllVersions() (result []results.Result)
	GetURL() (result string)
	GetName() (result string)
	GetDependencies() (results []string)
	GetGitURL() (result string)
	GetDescription() (result string)
	CheckUpdate() (outOfDate bool, result results.Result)
	UpdatePackage(input results.Result) (err error)
}

// bundle is a struct for handling parsers within Binoc.
type bundle struct {
	FileExt string
	Name    string
	Parser  Parser
}

var (
	// AvailableParsers is a map containing all available parsing engines.
	AvailableParsers map[string]bundle
)

func registerParser(parser Parser, fileExt string) {
	// Create parser's struct if not already created.
	if AvailableParsers == nil {
		AvailableParsers = make(map[string]bundle)
	}

	// Add parser to known parsers.
	name := strings.ToLower(reflect.ValueOf(parser).Type().Name())
	AvailableParsers[name] = bundle{Parser: parser, FileExt: fileExt}

	// Load parser config
	parser = parserConfEnv(name, parser)
}

func parserConfEnv(name string, in interface{}) Parser {
	val := reflect.ValueOf(in)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}
	// Get type of the input interface.
	refType := val.Type()

	// Create new interface to write changes too.
	new := reflect.Indirect(reflect.New(refType))

	for j := 0; j < refType.NumField(); j++ {
		fieldName := refType.Field(j).Name
		new.Field(j).SetString(val.Field(j).String())
		for _, prefix := range []string{"BINOC", "INPUT"} {
			evName := prefix + "_" + "PARSER" + "_" + strings.ToUpper(name) + "_" + strings.ToUpper(fieldName)
			evVal, evExists := os.LookupEnv(evName)
			if evExists {
				new.Field(j).SetString(evVal)
			}
		}
	}
	return new.Interface().(Parser)
}
