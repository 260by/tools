package ssh

import (
	"bytes"
	"golang.org/x/crypto/ssh"
)

// Command 通过ssh连接执行远程命令,支持sudo
func Command(client *ssh.Client, cmd string) (stdout string, err error) {
	session, err := client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	// 执行sudo命令时需设置session.RequestPty
	modes := ssh.TerminalModes{
		ssh.ECHO: 0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return
	}

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	err = session.Run(cmd)
	if err != nil {
		return
	}

	stdout = string(stdoutBuf.Bytes())

	return
}
