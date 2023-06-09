//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a process
	proc := MockProcess{}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, os.Kill, syscall.SIGINT)

	doneChan := make(chan struct{})
	go func() {
		i := 0
		for sig := range signals {
			log.Printf("receive %s\n", sig)
			if i == 0 {
				go runStop(&proc)
				i++
			} else {
				doneChan <- struct{}{}
			}

		}
	}()
	// Run the process (blocking)
	go proc.Run()
	<-doneChan
}

func runStop(mock *MockProcess) {
	mock.Stop()
}
