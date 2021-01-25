package update

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/autamus/go-parspack/pkg"

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
		url := app.Data.URL
		if app.Data.LatestVersion.URL != "" {
			url = app.Data.LatestVersion.URL
		}

		result, found := lookout.CheckUpdate(url)
		if found {
			if result.Version.Compare(app.Data.LatestVersion.Value) != 0 {
				resp, err := http.Get(url)
				if err != nil {
					log.Fatal(err)
				}
				bytes, err := ioutil.ReadAll(resp.Body)
				sha256 := fmt.Sprintf("%x", sha256.Sum256(bytes))

				resp.Body.Close()

				app.Data.AddVersion(pkg.Version{
					Value:    result.Version,
					Checksum: "sha256='" + sha256 + "'",
					URL:      result.Location,
				})
				app.Data.URL = result.Location
				output <- app
			}
		}
	}
	wg.Done()
}
