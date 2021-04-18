package parsers

import (
	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/version"
	"gopkg.in/yaml.v2"
)

type SHPC struct {
}

func init() {
	registerParser(SHPC{}, "container.yaml")
}

// Decode decodes a Container YAML Spec using a yaml parser.
func (s SHPC) Decode(content string) (pkg Package, err error) {
	// Parse YAML file
	internal := &ContainerSpec{}
	err = yaml.Unmarshal([]byte(content), &internal)
	// Generate name from URI
	if internal.Docker != "" {
		internal.Name = internal.Docker
	}
	if internal.Gh != "" {
		internal.Name = internal.Gh
	}
	return internal, err
}

// Encode encodes an updated container.yml using a yaml parser.
func (s SHPC) Encode(pkg Package) (result string, err error) {
	internal := pkg.(*ContainerSpec)
	internal.Name = ""
	output, err := yaml.Marshal(internal)
	return string(output), err
}

// ContainerSpec is a wrapper struct for a container.yaml
type ContainerSpec struct {
	Name        string            `yaml:"name,omitempty"`
	Docker      string            `yaml:"docker,omitempty"`
	Gh          string            `yaml:"gh,omitempty"`
	Url         string            `yaml:"url"`
	Maintainer  string            `yaml:"maintainer"`
	Description string            `yaml:"description"`
	Latest      map[string]string `yaml:"latest"`
	Versions    map[string]string `yaml:"tags"`
	Aliases     map[string]string `yaml:"aliases"`
}

// AddVersion adds a tagged version to a container.
func (s *ContainerSpec) AddVersion(input results.Result) (err error) {
	// Add version to versions.
	s.Versions[input.Version.String()] = input.Name
	s.Latest = map[string]string{}
	// Presume that added version is latest and make latest.
	s.Latest[input.Version.String()] = input.Name
	return nil
}

// GetLatestVersion returns the latest known tag of the container.
func (s *ContainerSpec) GetLatestVersion() (result version.Version) {
	for k, v := range s.Latest {
		return version.Version{k + "@" + v}
	}
	return
}

// GetURL returns the location of a container for Lookout
func (s *ContainerSpec) GetURL() (result string) {
	if s.Docker != "" {
		return "docker://" + s.Docker
	}
	return "https://github.com/" + s.Gh
}

// GetGitURL just returns the normal url for a container
func (s *ContainerSpec) GetGitURL() (result string) {
	return s.GetURL()
}

// GetName is a wrapper which returns the name of a container
func (s *ContainerSpec) GetName() string {
	return s.Name
}

// GetDependencies for containers doesn't do anything.
func (s *ContainerSpec) GetDependencies() []string {
	return []string{}
}

// CompareResult for containers compares the sha's to
// see if they are the same or not.
func (s *ContainerSpec) CompareResult(input results.Result) int {
	new := version.Version{input.Version.String() + "@" + input.Name}
	if s.GetLatestVersion().String() != new.String() {
		return -1
	}
	return 0
}
