package update

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/DataDrake/cuppa/version"
	lookout "github.com/alecbcs/lookout/update"
	"github.com/autamus/binoc/repo"
)

// Init initializes the lookout library
func Init(token string) {
	lookout.Init(token)
}

// RunPollWorker checks for an upstream update to the
// provided package on the input channel.
func RunPollWorker(wg *sync.WaitGroup, input <-chan repo.Result, output chan<- repo.Result) {
	for app := range input {
		url := app.Package.GetURL()
		result, found := lookout.CheckUpdate(url)
		if !found {
			result, found = lookout.CheckUpdate(app.Package.GetGitURL())
			if found {
				result.Location, found = patchGitURL(url, result.Version)
			}
		}
		if found && app.Package.CompareResult(*result) < 0 {
			fmt.Println("NOT UP TO DATE")
			app.LookOutput = *result
			output <- app
		}
	}
	wg.Done()
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
