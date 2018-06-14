package ssh

import (
	_ "fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"time"
)

func Connect(user, password, key, host, port string) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	if password != "" {
		auth = append(auth, ssh.Password(password))
	}
	if key != "" {
		buffer, err := ioutil.ReadFile(key)
		if err != nil {
			return nil, err
		}

		signer, err := ssh.ParsePrivateKey(buffer)
		if err != nil {
			return nil, err
		}

		auth = append(auth, ssh.PublicKeys(signer))
	}

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr = host + ":" + port

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}


	return sshClient, nil
}
