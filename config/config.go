package config

import (
	"os"
)

// Config defines the configuration struct for importing settings from ENV Variables
type Config struct {
	General general
	GitHub  github
	Repos   repos
}

type general struct {
}

type repos struct {
	Path string
}

type github struct {
	Token string
}

var (
	// Global is the configuration struct for the application.
	Global Config
)

func init() {
	Global.GitHub.Token = os.Getenv("BINOC_GITHUB_TOKEN")
	Global.Repos.Path = ".binoc/sources/"
}
