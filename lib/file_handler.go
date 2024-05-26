package lib

import (
	"os"
	"time"
)

func WaitForFile(path string, timeout time.Duration) bool {
	timeoutChan := time.After(timeout)
	tick := time.NewTicker(10 * time.Millisecond)

	for {
		select {
		case <-timeoutChan:
			return false
		case <-tick.C:
			if _, err := os.Stat(path); err == nil {
				return true
			}
		}
	}
}
