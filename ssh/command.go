package ssh

import (
	"bytes"
	_ "fmt"
	"log"
)

func Command(user, password, key, host, port, cmd string) bytes.Buffer {
	sshClient, err := Connect(user, password, key, host, port)
	if err != nil {
		log.Fatal("Authentication faild: ", err)
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	// cmd := "ls " + src

	if err := session.Run(cmd); err != nil {
		log.Fatalf("Faild to run: %s\nError: %v", cmd, err)
	}

	return stdoutBuf
}
