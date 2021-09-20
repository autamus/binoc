package update

import (
	"sync"

	lookout "github.com/alecbcs/lookout/update"
	"github.com/autamus/binoc/repo"
)

// Init initializes the lookout library
func Init(token string) {
	lookout.Init(token)
}

// RunPollWorker checks for an upstream update to the
// provided package on the input channel.
func RunPollWorker(
	wg *sync.WaitGroup,
	repo *repo.Repo,
	token string,
	input <-chan repo.Result,
	output chan<- repo.Result,
) {
	for app := range input {
		outOfDate, result := app.Package.CheckUpdate()
		if outOfDate {
			app.LookOutput = result
		}

		if outOfDate {
			output <- app
		}
	}
	wg.Done()
}
