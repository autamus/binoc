package update

import (
	"regexp"
	"strings"
	"sync"

	"github.com/DataDrake/cuppa/results"
	lookout "github.com/alecbcs/lookout/update"
	repository "github.com/autamus/binoc/repo"
	"github.com/autamus/binoc/upstream"
)

// Init initializes the lookout library
func Init(token string) {
	lookout.Init(token)
}

// RunPollWorker checks for an upstream update to the
// provided package on the input channel.
func RunPollWorker(
	wg *sync.WaitGroup,
	repo *repository.Repo,
	upstreamTemplatePath string,
	token string,
	upstreamOnly bool,
	input <-chan repository.Result,
	output chan<- repository.Result,
) {
	for app := range input {
		var outOfDate bool
		var result results.Result

		if !upstreamOnly {
			outOfDate, result = app.Package.CheckUpdate()
			if outOfDate {
				app.LookOutput = result
			}
		}
		if upstreamTemplatePath != "" {
			pkg, remoteModified, err := upstream.GetPackage(upstreamTemplatePath, toHyphenCase(app.Package.GetName()), token)
			if err != nil {
				goto END
			}
			if remoteModified.After(app.Modified) &&
				!app.Equals(repository.Result{Package: pkg, Parser: app.Parser}) {
				for _, version := range app.Package.GetAllVersions() {
					pkg.AddVersion(version)
				}
				app.LookOutput = results.Result{
					Name:     "spack/upstream",
					Location: app.Package.GetURL(),
				}
				app.Package = pkg
				outOfDate = true
			}
		}
	END:
		if outOfDate {
			output <- app
		}
	}
	wg.Done()
}

// toHypenCase converts a string to a hyphenated version appropriate
// for the commandline.
func toHyphenCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	return strings.ToLower(snake)
}
