package main

import (
	"bytes"
	// "bufio"
	"fmt"
	// "io"
	"io/ioutil"
	// "log"
	"net"
	"os"
	// "time"
	// "sync"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"strings"
)

func main()  {
	// sshConfig := &ssh.ClientConfig{
	// 	User: "root",
	// 	Auth: []ssh.AuthMethod{
	// 		ssh.Password("XTkj123!@#")
	// 	},
	// }

	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			PublicKeyFile("/home/keith/public_key/local"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// sshConfig := &ssh.ClientConfig{
	// 	User: "root",
	// 	Auth: []ssh.AuthMethod{
	// 		SSHAgent(),
	// 	},
	// }

	connection, err := ssh.Dial("tcp", "192.168.1.113:22", sshConfig)
	if err != nil {
		fmt.Errorf("Failed to dial: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		fmt.Errorf("Failed to create session: %s", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		fmt.Errorf("request for pseudo terminal failed: %s", err)
	}
	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run("ls /root")
	if err != nil {
		fmt.Println(err)
	}
	out := string(buf.Bytes())
	fmt.Println(strings.TrimSuffix(out, "\n"))
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}