package ssh

import (
	"bytes"
	"golang.org/x/crypto/ssh"
)

func Command(client *ssh.Client, cmd string) (stdout string, err error) {
	session, err := client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	err = session.Run(cmd)
	if err != nil {
		return
	}

	stdout = string(stdoutBuf.Bytes())

	return
}
