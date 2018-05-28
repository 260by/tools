package sftp

import (
	"github.com/260by/tools/ssh"
	"github.com/pkg/sftp"
)

func Connect(user, password, key, host, port string) (*sftp.Client, error) {
	var sftpClient *sftp.Client
	var err error

	sshClient, _ := ssh.Connect(user, password, key, host, port)
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}
