//go:build windows

package openvpn

import (
	"fmt"
	"os/exec"
)

func ovpnConnect(ovpnConfig string) error {
	fmt.Println("Connected! Press CTRL+C to disconnect")

	cmd := exec.Command("Powershell")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer stdin.Close()
		cmdStr := fmt.Sprintf("C:\\Program` Files\\OpenVPN\\bin\\openvpn.exe --config %s --verb 0", ovpnConfig)
		fmt.Fprintln(stdin, cmdStr)
	}()
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		return err
	}
	fmt.Printf("%s\n", out)
	return nil
}
