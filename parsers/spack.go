package parsers

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/version"
	lookout "github.com/alecbcs/lookout/update"
	"github.com/autamus/go-parspack"
	"github.com/autamus/go-parspack/pkg"
)

var (
	SpackUpstreamLink string
)

// Spack is a wrapper struct for the Spack Parser
type Spack struct {
}

func init() {
	registerParser(Spack{}, "package.py")
}

// Decode decodes a Spack Spec using go-parspack
func (s Spack) Decode(content string) (pkg Package, err error) {
	internal := &SpackPackage{}
	internal.Raw = content
	internal.Data, err = parspack.Decode(string(content))
	return internal, err
}

// Encode encodes an updated Spack Spec using go-parspack
func (s Spack) Encode(pkg Package) (result string, err error) {
	internal := pkg.(*SpackPackage)
	return parspack.PatchVersion(internal.Data, internal.Raw)
}

// SpackPackage is a wrapper struct for a Spack Package
type SpackPackage struct {
	Data pkg.Package
	Raw  string
}

// AddVersion is a wrapper for interacting with a spack package
func (p *SpackPackage) AddVersion(input results.Result) (err error) {
	resp, err := http.Get(input.Location)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	sha256 := fmt.Sprintf("%x", sha256.Sum256(bytes))

	err = resp.Body.Close()
	if err != nil {
		return nil
	}

	p.Data.AddVersion(pkg.Version{
		Value:    input.Version,
		Checksum: "sha256='" + sha256 + "'",
		URL:      input.Location,
	})
	p.Data.URL = input.Location
	return nil
}

// GetLatestVersion is a wrapper for getting the latest version from
// a spack package.
func (p *SpackPackage) GetLatestVersion() (result version.Version) {
	return p.Data.LatestVersion.Value
}

// GetURL is a wrapper for getting the latest url from a spack
// package.
func (p *SpackPackage) GetURL() (result string) {
	result = p.Data.URL
	if p.Data.LatestVersion.URL != "" {
		result = p.Data.LatestVersion.URL
	}
	return result
}

// GetGitURL is a wrapper for getting the latest url from a spack
// package git repository.
func (p *SpackPackage) GetGitURL() (result string) {
	return p.Data.GitURL
}

// GetName is a wrapper which returns the name of a package
func (p *SpackPackage) GetName() string {
	return p.Data.Name
}

// GetDependencies is a wrapper which returns the dependencies of a package
func (p *SpackPackage) GetDependencies() []string {
	return p.Data.Dependencies
}

// GetDescription returns the package's description.
func (p *SpackPackage) GetDescription() string {
	return p.Data.Description
}

// CheckUpdate checks for an update to source code project
// of the current Spack package.
func (p *SpackPackage) CheckUpdate() (outofDate bool, result results.Result) {
	url := p.GetURL()
	out, found := lookout.CheckUpdate(url)
	if !found {
		out, found = lookout.CheckUpdate(p.GetGitURL())
		if found {
			result.Location, found = patchGitURL(url, result.Version)
		}
	}
	if found {
		result = *out
	}

	// Check for update from Spack
	if SpackUpstreamLink != "" {
		concreteLink := strings.ReplaceAll(SpackUpstreamLink, "{{package}}", toHyphenCase(p.GetName()))
		resp, err := http.Get(concreteLink)
		if err != nil || resp.StatusCode != 200 {
			goto END
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			goto END
		}
		spackParser := Spack{}
		newPkg, err := spackParser.Decode(string(body))
		if err != nil {
			goto END
		}
		if newPkg.GetLatestVersion().Compare(p.GetLatestVersion()) < 0 {
			// Test encode and redecode to skip packages which would introduce
			// errors to the system.
			testOutput, err := spackParser.Encode(newPkg)
			if err != nil {
				goto END
			}
			_, err = spackParser.Decode(testOutput)
			if err != nil {
				goto END
			}
			p.Data = newPkg.(*SpackPackage).Data
			p.Raw = newPkg.(*SpackPackage).Raw

			// Setup fake result for new package.py
			if !found {
				result.Location = newPkg.GetURL()
				result.Version = newPkg.GetLatestVersion()
				result.Published = time.Now()
				result.Name = "spack/upstream"
				return true, result
			}
		}
	}
END:
	outOfDate := found && result.Version.Less(p.Data.LatestVersion.Value)
	return outOfDate, result
}

func (p *SpackPackage) UpdatePackage(input results.Result) (err error) {
	if input.Name != "spack/upstream" {
		return p.AddVersion(input)
	}
	return nil
}

// patchGitURL attempts to find an updated release url based on the version from the git url.
func patchGitURL(url string, input version.Version) (output string, found bool) {
	vexp := regexp.MustCompile(`([0-9]{1,4}[.])+[0-9,a-d]{1,4}`)
	output = vexp.ReplaceAllString(url, strings.Join(input, "."))

	resp, err := http.Head(output)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return "", false
	}
	return output, true
}

// ToHypenCase converts a string to a hyphenated version appropriate
// for the commandline.
func toHyphenCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	return strings.ToLower(snake)
}
