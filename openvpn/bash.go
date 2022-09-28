//go:build linux || darwin

package openvpn

import "os/exec"

func ovpnConnect(ovpnConfig string) *exec.Cmd {
    return exec.Command("sudo", "openvpn", "--config", ovpnConfig, "--verb", "0", "--script-security", "2",
        "--up", "/bin/echo \"Connected! Press CTRL+C to disconnect\n\"")
}
