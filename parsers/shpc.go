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
	if err != nil {
		return result, err
	}
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
	Name            string            `yaml:"name,omitempty"`
	Oras            string            `yaml:"oras,omitempty"`
	Docker          string            `yaml:"docker,omitempty"`
	Gh              string            `yaml:"gh,omitempty"`
	Url             string            `yaml:"url,omitempty"`
	Maintainer      string            `yaml:"maintainer"`
	Description     string            `yaml:"description"`
	Latest          map[string]string `yaml:"latest"`
	Versions        map[string]string `yaml:"tags"`
	Filter          []string          `yaml:"filter,omitempty"`
	Aliases         map[string]string `yaml:"-"`
	AliasesStruct   []Alias           `yaml:"-"`
	Features        map[string]bool   `yaml:"features,omitempty"`
	SingularityOpts string            `yaml:"singularity_options,omitempty"`
	DockerOpts      string            `yaml:"docker_options,omitempty"`
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
func (s *ContainerSpec) GetLatestVersion() (result results.Result) {
	for k := range s.Latest {
		return results.Result{
			Version:  version.Version{k},
			Location: s.Url,
		}
	}
	return
}

func (s *ContainerSpec) GetAllVersions() (result []results.Result) {
	for v := range s.Versions {
		result = append(result, results.Result{
			Version:  version.Version{v},
			Location: s.Url,
		})
	}
	return result
}

// GetURL returns the location of a container for Lookout
func (s *ContainerSpec) GetURL() (result string) {

	// Empty string means we do not know how to parse (yet)
	var result string
	if s.Docker != "" {
		result = "docker://" + s.Docker
		if len(s.Filter) > 0 {
			result = result + ":" + s.Filter[0]
		}
	} else if s.Gh != "" {
		result = "https://github.com/" + s.Gh
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

// GetDescription returns the package's description.
func (s *ContainerSpec) GetDescription() string {
	return s.Description
}

// CheckUpdate checks for an update to the container
func (s *ContainerSpec) CheckUpdate() (outOfDate bool, output results.Result) {
	outOfDate = false
	url := s.GetURL()
	docker := strings.HasPrefix(url, "docker://")

	// Check for new latest version
	out, found := lookout.CheckUpdate(url)
	if found {
		result := *out
		latestKey := s.GetLatestVersion().Version.String()
		latest := version.Version{latestKey + "@" + s.Latest[latestKey]}
		var new version.Version
		if docker {
			new = version.Version{result.Version.String() + "@" + result.Name}
		} else {
			new = version.Version{latestKey + "@" + result.Version.String()}

			// A gh release expects the "tag" as the recipe extension
			result.Name = result.Version.String()

			// And the digest as the release version
			result.Version = version.NewVersion(latestKey)
		}
		if latest.String() != new.String() {
			outOfDate = true
			s.AddVersion(result)
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
			out, found := lookout.CheckUpdate(baseUrl + ":" + tag)
			if found {
				result := *out
				if digest != result.Name {
					outOfDate = true
					s.Versions[tag] = result.Name
					if output.Location == "" {
						output = result
					}
					if s.Latest[tag] != "" {
						s.Latest[tag] = result.Name
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
