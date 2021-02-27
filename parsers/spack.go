package parsers

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alecbcs/cuppa/results"
	"github.com/alecbcs/cuppa/version"
	"github.com/autamus/go-parspack"
	"github.com/autamus/go-parspack/pkg"
)

// Spack is a wrapper struct for the Spack Parser
type Spack struct {
}

func init() {
	registerParser(Spack{}, "*.py")
}

// Decode decodes a Spack Spec using go-parspack
func (s Spack) Decode(content string) (pkg Package, err error) {
	internal := SpackPackage{}
	internal.Raw = content
	internal.Data, err = parspack.Decode(string(content))
	return internal, err
}

// Encode encodes an updated Spack Spec using go-parspack
func (s Spack) Encode(pkg Package) (result string, err error) {
	internal, ok := pkg.(SpackPackage)
	if !ok {
		return "", errors.New("cannot encode not a spackPackage")
	}
	return parspack.PatchVersion(internal.Data, internal.Raw)
}

// SpackPackage is a wrapper struct for a Spack Package
type SpackPackage struct {
	Data pkg.Package
	Raw  string
}

// AddVersion is a wrapper for interacting with a spack package
func (p SpackPackage) AddVersion(input results.Result) (err error) {
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
func (p SpackPackage) GetLatestVersion() (result version.Version) {
	return p.Data.LatestVersion.Value
}

// GetURL is a wrapper for getting the latest url from a spack
// package.
func (p SpackPackage) GetURL() (result string) {
	result = p.Data.URL
	if p.Data.LatestVersion.URL != "" {
		result = p.Data.LatestVersion.URL
	}
	return result
}

// GetName is a wrapper which returns the name of a package
func (p SpackPackage) GetName() string {
	return p.Data.Name
}
