package main

import (
	"bytes"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
	"sync"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	// "strings"
)

// Config Contains main authority information
type Config struct {
	User     string
	Server   string
	Key      string
	KeyPath  string
	Port     string
	Password string
	Timeout  time.Duration
	Proxy    ProxyConfig
}

// ProxyConfig for ssh proxy config
type ProxyConfig struct {
	User     string
	Server   string
	Key      string
	KeyPath  string
	Port     string
	Password string
	Timeout  time.Duration
}


func main()  {
	ssh := &Config{
		User:    "root",
		Server:  "10.111.1.12",
		Port:    "22",
		KeyPath: "/home/keith/public_key/haochang-admin-key",
		Proxy: ProxyConfig{
			User:    "zengming",
			Server:  "123.57.80.54",
			Port:    "22",
			KeyPath: "/home/keith/id_rsa",
		},
	}

	// ssh := &Config{
	// 	User:    "root",
	// 	Server:  "192.168.1.113",
	// 	Port:    "22",
	// 	KeyPath: "/home/keith/public_key/local",
	// }


	// stdout, stderr, isTimeout, err := ssh.Run("ls /root", 60)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// if !isTimeout {
	// 	log.Fatalln("time out")
	// }
	// if stderr != "" {
	// 	log.Fatalln(stderr)
	// }
	// // s := strings.Split(stdout, "\n")
	// fmt.Println(stdout)

	stdout, err := ssh.Command("ls /root")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(stdout)
}

func getPrivateKeyFile(file string) (ssh.Signer, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func getSSHConfig(config ProxyConfig) (*ssh.ClientConfig, io.Closer) {
	var sshAgent io.Closer

	// auths holds the detected ssh auth methods
	auths := []ssh.AuthMethod{}

	// figure out what auths are requested, what is supported
	if config.Password != "" {
		auths = append(auths, ssh.Password(config.Password))
	}
	if config.KeyPath != "" {
		if pubkey, err := getPrivateKeyFile(config.KeyPath); err != nil {
			log.Printf("Get private key file error: %v\n", err)
		} else {
			auths = append(auths, ssh.PublicKeys(pubkey))
		}
	}

	if config.Key != "" {
		if signer, err := ssh.ParsePrivateKey([]byte(config.Key)); err != nil {
			log.Printf("Parse private key error: %v\n", err)
		} else {
			auths = append(auths, ssh.PublicKeys(signer))
		}
	}

	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers))
	}

	return &ssh.ClientConfig{
		Timeout:         config.Timeout,
		User:            config.User,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, sshAgent
}

// Connect 连接到远程服务器并返回*ssh.Session
func (ssh_conf *Config) Connect() (*ssh.Session, error) {
	var client *ssh.Client
	var err error

	targetConfig, closer := getSSHConfig(ProxyConfig{
		User:     ssh_conf.User,
		Key:      ssh_conf.Key,
		KeyPath:  ssh_conf.KeyPath,
		Password: ssh_conf.Password,
		Timeout:  ssh_conf.Timeout,
	})
	if closer != nil {
		defer closer.Close()
	}

	// 开启ssh代理
	if ssh_conf.Proxy.Server != "" {
		proxyConfig, closer := getSSHConfig(ProxyConfig{
			User:     ssh_conf.Proxy.User,
			Key:      ssh_conf.Proxy.Key,
			KeyPath:  ssh_conf.Proxy.KeyPath,
			Password: ssh_conf.Proxy.Password,
			Timeout:  ssh_conf.Proxy.Timeout,
		})
		if closer != nil {
			defer closer.Close()
		}

		proxyClient, err := ssh.Dial("tcp", net.JoinHostPort(ssh_conf.Proxy.Server, ssh_conf.Proxy.Port), proxyConfig)
		if err != nil {
			return nil, err
		}

		conn, err := proxyClient.Dial("tcp", net.JoinHostPort(ssh_conf.Server, ssh_conf.Port))
		if err != nil {
			return nil, err
		}

		ncc, chans, reqs, err := ssh.NewClientConn(conn, net.JoinHostPort(ssh_conf.Server, ssh_conf.Port), targetConfig)
		if err != nil {
			return nil, err
		}

		client = ssh.NewClient(ncc, chans, reqs)
	} else {
		client, err = ssh.Dial("tcp", net.JoinHostPort(ssh_conf.Server, ssh_conf.Port), targetConfig)
		if err != nil {
			return nil, err
		}
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Command 执行命令
func (ssh_conf *Config) Command(command string) (stdout string, err error) {
	session, err := ssh_conf.Connect()
	if err != nil {
		return
	}

	// 创建伪终端, 执行sudo命令时需设要
	modes := ssh.TerminalModes{
		ssh.ECHO: 53,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err = session.RequestPty("vt100", 80, 40, modes)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run("ls /root")
	if err != nil {
		fmt.Println(err)
	}
	stdout = string(buf.Bytes())
	return
}

func (ssh_conf *Config) Stream(command string, timeout time.Duration) (<-chan string, <-chan string, <-chan bool, <-chan error, error) {
	// continuously send the command's output over the channel
	stdoutChan := make(chan string)
	stderrChan := make(chan string)
	doneChan := make(chan bool)
	errChan := make(chan error)

	// connect to remote host
	session, err := ssh_conf.Connect()
	if err != nil {
		return stdoutChan, stderrChan, doneChan, errChan, err
	}

	// 创建伪终端, 执行sudo命令时需设要
	modes := ssh.TerminalModes{
		ssh.ECHO: 53,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err = session.RequestPty("vt100", 80, 40, modes)
	if err != nil {
		return stdoutChan, stderrChan, doneChan, errChan, err
	}

	// defer session.Close()
	// connect to both outputs (they are of type io.Reader)
	outReader, err := session.StdoutPipe()
	if err != nil {
		return stdoutChan, stderrChan, doneChan, errChan, err
	}
	errReader, err := session.StderrPipe()
	if err != nil {
		return stdoutChan, stderrChan, doneChan, errChan, err
	}
	err = session.Start(command)
	if err != nil {
		return stdoutChan, stderrChan, doneChan, errChan, err
	}

	// combine outputs, create a line-by-line scanner
	stdoutReader := io.MultiReader(outReader)
	stderrReader := io.MultiReader(errReader)
	stdoutScanner := bufio.NewScanner(stdoutReader)
	stderrScanner := bufio.NewScanner(stderrReader)

	go func(stdoutScanner, stderrScanner *bufio.Scanner, stdoutChan, stderrChan chan string, doneChan chan bool, errChan chan error) {
		defer close(stdoutChan)
		defer close(stderrChan)
		defer close(doneChan)
		defer close(errChan)
		defer session.Close()

		timeoutChan := time.After(timeout * time.Second)
		res := make(chan struct{}, 1)
		var resWg sync.WaitGroup
		resWg.Add(2)

		go func() {
			for stdoutScanner.Scan() {
				stdoutChan <- stdoutScanner.Text()
			}
			resWg.Done()
		}()

		go func() {
			for stderrScanner.Scan() {
				stderrChan <- stderrScanner.Text()
			}
			resWg.Done()
		}()

		go func() {
			resWg.Wait()
			// close all of our open resources
			res <- struct{}{}
		}()

		select {
		case <-res:
			errChan <- session.Wait()
			doneChan <- true
		case <-timeoutChan:
			stderrChan <- "Run Command Timeout!"
			errChan <- nil
			doneChan <- false
		}
	}(stdoutScanner, stderrScanner, stdoutChan, stderrChan, doneChan, errChan)

	return stdoutChan, stderrChan, doneChan, errChan, err
}

// Run command on remote machine and returns its stdout as a string
func (ssh_conf *Config) Run(command string, timeout time.Duration) (outStr string, errStr string, isTimeout bool, err error) {
	stdoutChan, stderrChan, doneChan, errChan, err := ssh_conf.Stream(command, timeout)
	if err != nil {
		return outStr, errStr, isTimeout, err
	}
	// read from the output channel until the done signal is passed
loop:
	for {
		select {
		case isTimeout = <-doneChan:
			break loop
		case outline := <-stdoutChan:
			if outline != "" {
				fmt.Println("Line:",outline)
				outStr += outline + "\n"
			}
		case errline := <-stderrChan:
			if errline != "" {
				errStr += errline + "\n"
			}
		case err = <-errChan:
		}
	}
	// return the concatenation of all signals from the output channel
	return outStr, errStr, isTimeout, err
}