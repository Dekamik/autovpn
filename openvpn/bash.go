//go:build linux || darwin

package openvpn

import "os/exec"

func ovpnConnect(ovpnConfig string) error {
	fmt.Println("Connecting...")

	cmd := exec.Command("sudo", "openvpn", "--config", ovpnConfig, "--verb", "0", "--script-security", "2",
		"--up", "/bin/echo \"Connected! Press CTRL+C to disconnect\n\"")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}

	var waiting = true

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	go func() {
		<-sigc
		_ = cmd.Process.Kill()
		fmt.Println("\nDisconnected")
		waiting = false
	}()

	for waiting {
	}

	return nil
}
