package config

import (
	"os"
	"reflect"
	"strings"
)

// Config defines the configuration struct for importing settings from ENV Variables
type Config struct {
	General general
	Git     git
	Repo    repo
	Parsers parsers
	Branch  branch
	PR      pr
}

type general struct {
	Version string
	Action  string
}

type repo struct {
	Path              string
	SpackUpstreamLink string
}

type git struct {
	Name     string
	Username string
	Email    string
	Token    string
}

type parsers struct {
	Loaded string
}

type branch struct {
	Prefix string
}

type pr struct {
	IgnoreLabel string
	Skip        string
}

var (
	// Global is the configuration struct for the application.
	Global Config
)

func init() {
	defaultConfig()
	parseConfigEnv()
}

func defaultConfig() {
	Global.General.Version = "0.3.3"
	Global.Parsers.Loaded = "spack,shpc"
	Global.Branch.Prefix = "binoc/"
	Global.PR.IgnoreLabel = "binoc-blacklist"
	Global.PR.Skip = "true"
}

func parseConfigEnv() {
	numSubStructs := reflect.ValueOf(&Global).Elem().NumField()
	for i := 0; i < numSubStructs; i++ {
		iter := reflect.ValueOf(&Global).Elem().Field(i)
		subStruct := strings.ToUpper(iter.Type().Name())

		structType := iter.Type()
		for j := 0; j < iter.NumField(); j++ {
			fieldVal := iter.Field(j).String()
			if fieldVal != "Version" {
				fieldName := structType.Field(j).Name
				for _, prefix := range []string{"BINOC", "INPUT"} {
					evName := prefix + "_" + subStruct + "_" + strings.ToUpper(fieldName)
					evVal, evExists := os.LookupEnv(evName)
					if evExists && evVal != fieldVal {
						iter.FieldByName(fieldName).SetString(evVal)
					}
				}
			}
		}
	}
}
