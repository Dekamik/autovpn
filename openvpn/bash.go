//go:build linux || darwin

package openvpn

import "os/exec"

func ovpnConnect(ovpnConfig string) *exec.Cmd {
	return exec.Command("sudo", "openvpn", ovpnConfig)
}
