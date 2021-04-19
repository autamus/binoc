package parsers

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/version"
	lookout "github.com/alecbcs/lookout/update"
	"github.com/autamus/go-parspack"
	"github.com/autamus/go-parspack/pkg"
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

// CheckUpdate checks for an update to source code project
// of the current Spack package.
func (p *SpackPackage) CheckUpdate() (outofDate bool, result *results.Result) {
	url := p.GetURL()
	result, found := lookout.CheckUpdate(url)
	if !found {
		result, found = lookout.CheckUpdate(p.GetGitURL())
		if found {
			result.Location, found = patchGitURL(url, result.Version)
		}
	}
	outOfDate := found && p.Data.LatestVersion.Value.Less(result.Version)
	if outOfDate {
		p.AddVersion(*result)
	}
	return outOfDate, result
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
