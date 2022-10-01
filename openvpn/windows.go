//go:build windows

package openvpn

import (
	"fmt"
	"os/exec"
)

func IsInstalled(executable string) bool {
	cmd := exec.Command("Powershell", "Start", executable, "-ArgumentList", "--version", "-NoNewWindow")
	err := cmd.Run()
	return err == nil
}

func GetExecutable(executableOverride string) string {
	if len(executableOverride) == 0 {
		return "C:\\Program` Files\\OpenVPN\\bin\\openvpn.exe"
	} else {
		return executableOverride
	}
}

func ovpnConnect(executable string, ovpnConfig string) *exec.Cmd {
	return exec.Command("Powershell", "Start", executable, "-ArgumentList",
		fmt.Sprintf("--config,%s,--verb,0", ovpnConfig), "-NoNewWindow")
}
