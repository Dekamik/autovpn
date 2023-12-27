//go:build linux || darwin

package openvpn

import (
	"os/exec"
)

func IsInstalled(executable string) bool {
	// If the version command returns an ExitError, we know it's installed.
	// If it's not installed, it returns a regular Error type.
	cmd := exec.Command(executable, "--version")
	err := cmd.Run()
	return err == nil
}

func GetExecutable(executableOverride string) string {
	if len(executableOverride) == 0 {
		return "openvpn"
	} else {
		return executableOverride
	}
}

func ovpnConnect(executable string, ovpnConfig string) *exec.Cmd {
	return exec.Command(executable, "--config", ovpnConfig, "--verb", "0")
}
