package helpers

import (
	"fmt"
	"time"
)

// WaitPrint is a goroutine that prints the message and prints dots while waiting.
// It finishes when it receives true in the finish bool channel.
// To await method completion, simply call <-exited.
// This command uses VT100 escape codes to erase the line in terminal.
func WaitPrint(message string, finish chan bool, exited chan bool) {
	i := 0
	for {
		select {
		case <-finish:
			fmt.Printf("\033[2K\r%s [%ds] OK\n", message, i)
			exited <- true
			return
		default:
			fmt.Printf("\033[2K\r%s [%ds]", message, i)
			time.Sleep(time.Second)
			i++
		}
	}
}
