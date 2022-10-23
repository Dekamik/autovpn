package openvpn

import (
	"autovpn/data"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
	"os/signal"
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

func runCommands(sshClient *ssh.Client, commands []string) error {
	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()

	err = session.Shell()
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			return err
		}
	}

	err = session.Wait()
	if err != nil {
		return err
	}

	return nil
}

func Install(instance data.Instance, installScriptUrl string, setupTimeout []string) (*string, error) {
	configPath := "client.ovpn"
	config := &ssh.ClientConfig{
		User:            instance.User,
		Auth:            []ssh.AuthMethod{ssh.Password(instance.Pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := dial("tcp", fmt.Sprintf("%s:%d", instance.IpAddress, instance.SshPort), config, 10, 0)
	if err != nil {
		return nil, err
	}
	defer sshClient.Close()

	commands := []string{
		fmt.Sprintf("curl %s -o openvpn-install.sh", installScriptUrl),
		"chmod +x openvpn-install.sh",
		"export AUTO_INSTALL=y; ./openvpn-install.sh",
		"sed -i 's/^verb [0-9]*$/verb 0/g' /etc/openvpn/server.conf",
	}
	if setupTimeout != nil {
		commands = append(commands, setupTimeout...)
	}
	commands = append(commands, "exit")

	err = runCommands(sshClient, commands)
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

func Connect(executable string, ovpnConfig string, setupTimeout []string, instance data.Instance) error {
	fmt.Println("Connecting...")

	cmd := ovpnConnect(executable, ovpnConfig)
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

	fmt.Println("Connected! Press CTRL+C to disconnect")
	fmt.Printf("Server timeout: approx. %s\n", time.Now().Add(time.Hour))
	lastTimeoutSetup := time.Now()

	for waiting {
		time.Sleep(time.Millisecond * 200)

		if len(setupTimeout) != 0 && time.Since(lastTimeoutSetup) >= time.Minute*30 {
			config := &ssh.ClientConfig{
				User:            instance.User,
				Auth:            []ssh.AuthMethod{ssh.Password(instance.Pass)},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}
			sshClient, err := dial("tcp", fmt.Sprintf("%s:%d", instance.IpAddress, instance.SshPort), config, 10, 0)
			if err != nil {
				return err
			}

			setupTimeout = append(setupTimeout, "exit")
			err = runCommands(sshClient, setupTimeout)
			sshClient.Close()
			if err != nil {
				return err
			}

			lastTimeoutSetup = time.Now()
			fmt.Printf("Server timeout updated: %s\n", time.Now().Add(time.Hour))
		}
	}

	return nil
}
