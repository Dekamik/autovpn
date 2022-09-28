package openvpn

import (
	"autovpn/providers"
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"os/signal"
)

func Install(instance providers.Instance, installScriptUrl string) (*string, error) {
	configPath := "./client.ovpn"
	config := &ssh.ClientConfig{
		User:            instance.RootUser,
		Auth:            []ssh.AuthMethod{ssh.Password(instance.RootPass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", instance.IpAddress, instance.SshPort), config)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Run(fmt.Sprintf("curl %s -o openvpn-install.sh", installScriptUrl)); err != nil {
		return nil, err
	}
	if err := session.Run("chmod +x openvpn-install.sh"); err != nil {
		return nil, err
	}
	if err := session.Run("export AUTO_INSTALL=y; ./openvpn-install.sh"); err != nil {
		return nil, err
	}
	if err := session.Run("sed -i 's/^verb [0-9]*$/verb 0/g' /etc/openvpn/server.conf"); err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	session.Stdout = &buffer
	if err := session.Run("cat /root/client.ovpn"); err != nil {
		return nil, err
	}

	f, err := os.Create(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = f.Write(buffer.Bytes())
	if err != nil {
		return nil, err
	}

	return &configPath, nil
}

func Connect(ovpnConfig string) error {
	cmd := ovpnConnect(ovpnConfig)
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
		waiting = false
	}()

	for waiting {
	}

	return nil
}
