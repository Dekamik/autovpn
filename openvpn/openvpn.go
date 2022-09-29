package openvpn

import (
    "autovpn/providers"
    "fmt"
    "github.com/pkg/sftp"
    "golang.org/x/crypto/ssh"
    "os"
    "time"
)

func dial(network string, addr string, config *ssh.ClientConfig, maxTries int, currentTry int) (*ssh.Client, error) {
	sshClient, err := ssh.Dial(network, addr, config)
	if err != nil {
		if currentTry < maxTries {
			time.Sleep(time.Second * 3)
			return dial(network, addr, config, maxTries, currentTry+1)
		}
		return nil, err
	}
	return sshClient, nil
}

func Install(instance providers.Instance, installScriptUrl string) (*string, error) {
	configPath := "client.ovpn"
	config := &ssh.ClientConfig{
		User:            instance.RootUser,
		Auth:            []ssh.AuthMethod{ssh.Password(instance.RootPass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := dial("tcp", fmt.Sprintf("%s:%d", instance.IpAddress, instance.SshPort), config, 10, 0)
	if err != nil {
		return nil, err
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()

	err = session.Shell()
	if err != nil {
		return nil, err
	}

	commands := []string{
		fmt.Sprintf("curl %s -o openvpn-install.sh", installScriptUrl),
		"chmod +x openvpn-install.sh",
		"export AUTO_INSTALL=y; ./openvpn-install.sh",
		"sed -i 's/^verb [0-9]*$/verb 0/g' /etc/openvpn/server.conf",
		"exit",
	}

	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			return nil, err
		}
	}

	err = session.Wait()
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}
	defer sftpClient.Close()

	localConfig, err := os.Create(configPath)
	if err != nil {
		return nil, err
	}
	defer localConfig.Close()

	remoteConfig, err := sftpClient.Open("/root/client.ovpn")
	if err != nil {
		return nil, err
	}
	defer remoteConfig.Close()

	if _, err := localConfig.ReadFrom(remoteConfig); err != nil {
		return nil, err
	}

	return &configPath, nil
}

func Connect(ovpnConfig string) error {
	return ovpnConnect(ovpnConfig)
}
