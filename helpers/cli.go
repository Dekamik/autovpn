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
	startTime := time.Now().Local()
	wheel := "-\\|/"
	i := 0
	eraseCode := map[string]string{
		"bash":    "\u001B[2K\r",
		"windows": "ESC[K",
	}
	for {
		select {
		case <-finish:
			fmt.Printf("%s%-35s %-6s OK\n", eraseCode["bash"], message, time.Since(startTime).Truncate(time.Second))
			exited <- true
			return
		default:
			fmt.Printf("%s%-35s %-6s %s ", eraseCode["bash"], message, time.Since(startTime).Truncate(time.Second), string(wheel[i]))
			time.Sleep(time.Millisecond * 100)

			if i == len(wheel)-1 {
				i = 0
			} else {
				i++
			}
		}
	}
}
