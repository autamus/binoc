package display

import (
	"fmt"
	"sync"
	"time"
)

// Spinner is an array of the progression of the spinner.
var Spinner = []string{"|", "/", "-", "\\"}

// SpinnerWait displays the actual spinner
func SpinnerWait(done chan int, message string, wg *sync.WaitGroup) {
	ticker := time.Tick(time.Millisecond * 120)
	frameCounter := 0
	for {
		select {
		case _ = <-done:
			wg.Done()
			return
		default:
			<-ticker
			ind := frameCounter % len(Spinner)
			fmt.Printf("\r[%v] "+message, Spinner[ind])
			frameCounter++
		}
	}
}
