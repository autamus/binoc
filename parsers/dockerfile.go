package parsers

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/version"
//	lookout "github.com/alecbcs/lookout/update"
)

type Dockerfile struct {}

func init() {
	registerParser(Dockerfile{}, "Dockerfile")
}

// Dockerfile parser does allow a prefix
func (s Dockerfile) AllowsPrefix() bool {
	return true
}

// Decode decodes a Dockerfile using a yaml parser.
func (s Dockerfile) Decode(content string) (pkg Package, err error) {

	// Prepare a Dockerfile, which implements required functions for
	// a package interface, plus additional fields
	internal := &DockerfilePackage{}
	internal.Raw = content

	// Read through each line looking for FROM
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineno := 0
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")

		// Did we find a FROM?
		if strings.HasPrefix(line, "FROM") {

			// Parse the name, version, isVaraible
			isVariable := strings.Contains(line, "$")
			
			// Create a new FROM entry
			newfrom := From{Raw: line, LineNo: lineno, IsVariable: isVariable}

			// Now get rid of the FROM and trim the edges
			line = strings.Replace(line, "FROM", "", 1)
			line = strings.Trim(line, " ")

			// Do we have an AS something?
			hasAs := strings.Contains(strings.ToLower(line), " as ")
			if hasAs {
				// Find the index of where " as " starts to split it
				asIndex := strings.Index(strings.ToLower(line), " as ")
				line = line[0:asIndex]
			}
			// The container includes the name and version			
			newfrom.Container = line

			// If we have a : then there is a version, otherwise latest
			containerName := newfrom.Container
			version := "latest"
			if strings.Contains(line, ":") {
				parts := strings.SplitN(line, ":", 2)
				containerName = parts[0]	
				version = parts[1]			
			}
			newfrom.Name = containerName
			newfrom.Version = version
			
			// Add the parsed FROM to our list
			internal.Froms = append(internal.Froms, newfrom)
		}
		lineno += 1
	}
	return internal, err
}

// Encode encodes an updated Dockerfile
func (s Dockerfile) Encode(pkg Package) (result string, err error) {
	internal := pkg.(*DockerfilePackage)
	internal.Name = ""
	fmt.Println(internal)

	// TODO how to replace current?
	return result, err
}

// A FROM stateent
type From struct {
	Raw		string
	Container	string
	Name		string
	Version	string

	// Is the FROM a variable (meaning we shouldn't change it)
	IsVariable	bool
	LineNo		int
}

// A Dockerfile is a wrapper struct for a Dockerfile
type DockerfilePackage struct {
	Name		string
	Raw		string
	Froms		[]From
}


// AddVersion adds a tagged version to a container.
func (s *DockerfilePackage) AddVersion(input results.Result) (err error) {
	// Add version to versions.
	//s.Versions[input.Version.String()] = input.Name
	//s.Latest = map[string]string{}
	// Presume that added version is latest and make latest.
	//s.Latest[input.Version.String()] = input.Name
	return nil
}

// GetLatestVersion returns the latest known tag of the container.
func (s *DockerfilePackage) GetLatestVersion() (result version.Version) {
	//for k := range s.Latest {
	//	return version.Version{k}
	//}
	return
}

// GetURL returns the location of a container for Lookout
func (s *DockerfilePackage) GetURL() (result string) {
	result = "docker://busybox"
	//if s.Docker != "" {
	//	result = "docker://" + s.Docker
	//	if len(s.Filter) > 0 {
	//		result = result + ":" + s.Filter[0]
	//	}
	//} else {
	//	result = "https://github.com/" + s.Gh
	//}
	return result
}

// GetGitURL just returns the normal url for a container
func (s *DockerfilePackage) GetGitURL() (result string) {
	return s.GetURL()
}

// GetName is a wrapper which returns the name of a container
func (s *DockerfilePackage) GetName() string {
	return s.Name
}

// GetDependencies for containers doesn't do anything.
func (s *DockerfilePackage) GetDependencies() []string {
	return []string{}
}

// GetDescription returns the package's description.
func (s *DockerfilePackage) GetDescription() string {
	return "Dockerfile"
}

// CheckUpdate checks for an update to the container
func (s *DockerfilePackage) CheckUpdate() (outOfDate bool, output results.Result) {
	outOfDate = false
	//url := s.GetURL()
	//docker := strings.HasPrefix(url, "docker://")

	// Check for new latest version
	//out, found := lookout.CheckUpdate(url)
	//if found {
	//	result := *out
	//	latestKey := s.GetLatestVersion().String()
	//	latest := version.Version{latestKey + "@" + s.Latest[latestKey]}
	//	var new version.Version
	//	if docker {
	//		new = version.Version{result.Version.String() + "@" + result.Name}
	//	} else {
	//		new = version.Version{latestKey + "@" + result.Version.String()}
//
	//		// A gh release expects the "tag" as the recipe extension
	//		result.Name = result.Version.String()
//
	//		// And the digest as the release version
	//		result.Version = version.NewVersion(latestKey)
	//	}
	//	if latest.String() != new.String() {
	//		outOfDate = true
	//		s.AddVersion(result)
	//		output = result
	//	}
	//}

	//if docker {
	//	// Iteratively check previous labels for digest updates
	//	var baseUrl string
	//	if len(s.Filter) > 0 {
	//		baseUrl = strings.TrimSuffix(url, ":"+s.Filter[0])
	//	} else {
	//		baseUrl = url
	//	}
	//	for tag, digest := range s.Versions {
	//		out, found := lookout.CheckUpdate(baseUrl + ":" + tag)
	//		if found {
	//			result := *out
	//			if digest != result.Name {
	//				outOfDate = true
	//				s.Versions[tag] = result.Name
	//				if output.Location == "" {
	//					output = result
	//				}
	//				if s.Latest[tag] != "" {
	//					s.Latest[tag] = result.Name
	//				}
	//			}
	//		}
	//	}
	//}

	return outOfDate, output
}

func (s *DockerfilePackage) UpdatePackage(input results.Result) (err error) {
	return nil
}
