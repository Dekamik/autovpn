package helpers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
	"time"
)

var ErrUnsupportedOS = errors.New("OS is not supported")

// IsAdmin checks if the program is running with administrative privileges.
// This is required when using a VPN, since those programs requires higher privileges than usual.
func IsAdmin() (bool, error) {
	switch runtime.GOOS {
	case "windows":
		_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		if err != nil {
			return false, nil
		}
		return true, nil

	case "darwin":
	case "linux":
		currentUser, err := user.Current()
		if err != nil {
			return false, err
		}
		if currentUser.Name != "root" {
			return false, nil
		}
		return true, nil
	}

	log.Printf("current OS = %s", runtime.GOOS)
	return false, ErrUnsupportedOS
}

// WaitPrint is a goroutine that prints the message and prints dots while waiting.
// It finishes when it receives true in the finish bool channel.
// To await method completion, simply call <-exited.
// This command uses VT100 escape codes to erase the line in terminal.
func WaitPrint(message string, finish chan bool, exited chan bool) {
	startTime := time.Now().Local()
	wheel := "-\\|/"
	i := 0
	eraseCode := "\u001B[2K\r"
	for {
		select {
		case <-finish:
			fmt.Printf("%s%-35s %-6s OK\n", eraseCode, message, time.Since(startTime).Truncate(time.Second))
			exited <- true
			return
		default:
			fmt.Printf("%s%-35s %-6s %s  ", eraseCode, message, time.Since(startTime).Truncate(time.Second), string(wheel[i]))
			time.Sleep(time.Millisecond * 100)

			if i == len(wheel)-1 {
				i = 0
			} else {
				i++
			}
		}
	}
}
