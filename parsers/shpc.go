package parsers

import (
	"strings"

	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/version"
	lookout "github.com/alecbcs/lookout/update"
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
	if err != nil {
		return internal, err
	}
	// Attempt to decode aliases by map
	aMap := AliasMap{}
	err = yaml.Unmarshal([]byte(content), &aMap)
	if err != nil {
		aStruct := AliasStruct{}
		err = yaml.Unmarshal([]byte(content), &aStruct)
		if err != nil {
			return internal, err
		}
		internal.AliasesStruct = aStruct.Aliases
	} else {
		internal.Aliases = aMap.Aliases
	}

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
	if err != nil {
		return result, err
	}
	result = result + string(output)

	// encode aliases
	aliasesMap, err := yaml.Marshal(&AliasMap{internal.Aliases})
	if string(aliasesMap) != "" && string(aliasesMap) != "{}\n" {
		result = result + string(aliasesMap)
	}
	aliasesStruct, err := yaml.Marshal(&AliasStruct{internal.AliasesStruct})
	if string(aliasesStruct) != "" && string(aliasesStruct) != "{}\n" {
		result = result + string(aliasesStruct)
	}

	return result, err
}

// ContainerSpec is a wrapper struct for a container.yaml
type ContainerSpec struct {
	Name          string            `yaml:"name,omitempty"`
	Docker        string            `yaml:"docker,omitempty"`
	Gh            string            `yaml:"gh,omitempty"`
	Url           string            `yaml:"url,omitempty"`
	Maintainer    string            `yaml:"maintainer"`
	Description   string            `yaml:"description"`
	Latest        map[string]string `yaml:"latest"`
	Versions      map[string]string `yaml:"tags"`
	Filter        []string          `yaml:"filter,omitempty"`
	Aliases       map[string]string `yaml:"-"`
	AliasesStruct []Alias           `yaml:"-"`
}

type AliasMap struct {
	Aliases map[string]string `yaml:"aliases,omitempty"`
}

type AliasStruct struct {
	Aliases []Alias `yaml:"aliases,omitempty"`
}

type Alias struct {
	Name    string `yaml:"name,omitempty"`
	Command string `yaml:"command,omitempty"`
	Options string `yaml:"options,omitempty"`
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
	for k := range s.Latest {
		return version.Version{k}
	}
	return
}

// GetURL returns the location of a container for Lookout
func (s *ContainerSpec) GetURL() (result string) {
	if s.Docker != "" {
		result = "docker://" + s.Docker
	} else {
		result = "https://github.com/" + s.Gh
	}
	if len(s.Filter) > 0 {
		result = result + ":" + s.Filter[0]
	}
	return result
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

// CheckUpdate checks for an update to the container
func (s *ContainerSpec) CheckUpdate() (outOfDate bool, output *results.Result) {
	outOfDate = false
	url := s.GetURL()
	docker := strings.HasPrefix(url, "docker://")

	// Check for new latest version
	result, found := lookout.CheckUpdate(url)
	if found && docker {
		latestKey := s.GetLatestVersion().String()
		latest := version.Version{latestKey + "@" + s.Latest[latestKey]}
		var new version.Version
		if docker {
			new = version.Version{result.Version.String() + "@" + result.Name}
		} else {
			new = version.Version{result.Version.String() + "@" + s.Name}
		}
		if latest.String() != new.String() {
			outOfDate = true
			s.AddVersion(*result)
			output = result
		}
	}

	if docker {
		// Iteratively check previous labels for digest updates
		var baseUrl string
		if len(s.Filter) > 0 {
			baseUrl = strings.TrimSuffix(url, ":"+s.Filter[0])
		} else {
			baseUrl = url
		}
		for tag, digest := range s.Versions {
			result, found := lookout.CheckUpdate(baseUrl + ":" + tag)
			if found {
				if digest != result.Name {
					outOfDate = true
					s.Versions[tag] = result.Name
					if output == nil {
						output = result
					}
				}
			}
		}
	}

	return outOfDate, output
}

func (s *ContainerSpec) UpdatePackage(input results.Result) (err error) {
	return nil
}
