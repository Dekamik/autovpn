//go:build windows

package openvpn

import (
	"fmt"
	"os/exec"
)

func ovpnConnect(ovpnConfig string) *exec.Cmd {
	fmt.Println("Connected! Press CTRL+C to disconnect")
	return exec.Command("Powershell", "Start", "C:\\Program` Files\\OpenVPN\\bin\\openvpn.exe",
		"-ArgumentList", fmt.Sprintf("--config,%s,--verb,0", ovpnConfig), "-NoNewWindow")
}
