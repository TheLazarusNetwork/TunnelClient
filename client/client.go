package client

import (
	"os"

	"golang.org/x/crypto/ssh"
)

func Conn() *ssh.ClientConfig {
	// refer to https://godoc.org/golang.org/x/crypto/ssh for other authentication types
	sshConfig := &ssh.ClientConfig{
		// SSH connection username
		User: os.Getenv("SSH_USER"),
		Auth: []ssh.AuthMethod{
			// publicKeyFile("/home/operatore/.ssh/id_rsa"),
			ssh.Password("access@2021"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return sshConfig
}
