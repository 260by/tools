package sftp

import (
	"github.com/260by/tools/ssh"
	"github.com/pkg/sftp"
)

func Connect(user, host string, port int, authentication ...string) (client *sftp.Client, err error) {
	sshClient, err := ssh.Connect(user, host, port, authentication...)
	if err != nil {
		return
	}
	if client, err = sftp.NewClient(sshClient); err != nil {
		return
	}
	return client, nil
}
