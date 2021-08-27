package parsers

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/version"
	lookout "github.com/alecbcs/lookout/update"
)

type Dockerfile struct{}

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
				newfrom.Extra = line[asIndex:]
			
			// otherwise, the extra is just whatever is beyond the version
			} else {
				parts := strings.SplitN(line, " ", 2)
				if len(parts) > 1 {
					newfrom.Extra = parts[1]
				} else {
					newfrom.Extra = ""
				}
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

	// Start with the original Dockerfile
	lines := strings.Split(internal.Raw, "\n")

	// For each Update, replace exact line with new version
	for _, from := range internal.Updates {
		fmt.Printf("\nUpdating %s:%s to %s on line %d\n", from.Name, from.Version, from.Updated, from.LineNo)
		lines[from.LineNo] = "FROM " + from.Updated + " " + from.Extra
	}
	dockerfile := strings.Join(lines, "\n")
	return dockerfile, err
}

// A FROM statement
type From struct {
	Raw       string
	Container string
	Name      string
	Version   string
	Updated   string

	// Extra content in the version string
	Extra	  string

	// Is the FROM a variable (meaning we shouldn't change it)
	IsVariable bool
	LineNo     int
}

// A Dockerfile is a wrapper struct for a Dockerfile
type DockerfilePackage struct {
	Name  string
	Raw   string

	// original list of Froms to check or parse
	Froms []From

	// Final Updates (also froms) that will be used to update!
	Updates []From
}

// GetAllVersions returns all versions (we don't use this)
func (s *DockerfilePackage) GetAllVersions() (result []results.Result) {
	return result
}

// AddVersion adds a tagged version to a container (we don't use this)
func (s *DockerfilePackage) AddVersion(input results.Result) (err error) {
	// TODO need to store entire FROM here as an update to do,
	// full name, sha/tag, and line number
	// to some other attribute on the class
	return nil
}

// GetLatestVersion returns the latest known tag of the container (we don't use this)
func (s *DockerfilePackage) GetLatestVersion() results.Result {
	return results.Result{Version: version.Version{}, Location: ""}
}

// We don't use this function - we use the GetURL that takes a docker URI
func (s *DockerfilePackage) GetURL() (result string) {
	result = "docker://busybox"
	return result
}

// We don't use this function - we use the GetURL that takes a docker URI
func (s *DockerfilePackage) GetNamedURL(name string) (result string) {
	return "docker://" + name
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

	// For each FROM, check if it's out of date
	for _, from := range s.Froms {

		// If there is a variable, we can't easily parse
		if from.IsVariable {
			continue
		}

		// This doesn't have a tag, and is always docker://
		url := s.GetNamedURL(from.Container)

		// Check for new latest version
		out, found := lookout.CheckUpdate(url)

		// If we find a result, get the latest and compare to current
		if found {
			result := *out
			latestKey := from.Version

			// Latest key is going to be a sha256sum, not a tag
			// TODO: question - can we keep an associated tag?
			latest := version.Version{from.Name + "@" + latestKey}
			newVersion := version.Version{from.Name + "@" + result.Name}

			if latest.String() != newVersion.String() {
				outOfDate = true
				
				// Updated from with updated version
				from.Updated = newVersion.String()
				s.Updates = append(s.Updates, from)
				output = result
			}
		}

	}
	return outOfDate, output
}

func (s *DockerfilePackage) UpdatePackage(input results.Result) (err error) {
	return nil
}
