package ssh

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"time"
	"strconv"
)

// Connect ssh连接,参数authentication为用户密码或用户私钥
func Connect(user, host string, port int, authentication ...string) (client *ssh.Client, err error) {
	auth := make([]ssh.AuthMethod,0)
	for _, a := range authentication {
		buffer, err := ioutil.ReadFile(a)
		if err == nil {
			signer, err := ssh.ParsePrivateKey(buffer)
			if err != nil {
				return nil, err
			}
			auth = append(auth, ssh.PublicKeys(signer))
		} else {
			auth = append(auth, ssh.Password(a))
		}
	}

	clientConfig := &ssh.ClientConfig{
		User:            user,
		Auth:    	     auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := host + ":" + strconv.Itoa(port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	return client, nil
}
