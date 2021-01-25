package config

import (
	"os"
)

// Config defines the configuration struct for importing settings from ENV Variables
type Config struct {
	General general
	Git     git
	Repos   repos
}

type general struct {
}

type repos struct {
	Path string
}

type git struct {
	Name     string
	Username string
	Email    string
	Token    string
}

var (
	// Global is the configuration struct for the application.
	Global Config
)

func init() {
	Global.Git.Name = os.Getenv("BINOC_GIT_NAME")
	Global.Git.Username = os.Getenv("BINOC_GIT_USERNAME")
	Global.Git.Email = os.Getenv("BINOC_GIT_EMAIL")
	Global.Git.Token = os.Getenv("BINOC_GIT_TOKEN")
	Global.Repos.Path = ".binoc/sources/"
}
