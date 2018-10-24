package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
	"regexp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"strings"
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
	// ssh := &Config{
	// 	User:    "root",
	// 	Server:  "10.111.1.12",
	// 	Port:    "22",
	// 	KeyPath: "/home/keith/public_key/haochang-admin-key",
	// 	Proxy: ProxyConfig{
	// 		User:    "zengming",
	// 		Server:  "123.57.80.54",
	// 		Port:    "22",
	// 		KeyPath: "/home/keith/id_rsa",
	// 	},
	// }

	ssh := &Config{
		User:    "root",
		Server:  "192.168.1.173",
		Port:    "22",
		KeyPath: "/home/zengm/public_key/local",
	}

	stdout, err := ssh.Command("/sbin/ifconfig")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(strings.Split(stdout, " "))
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

	// 创建伪终端, 执行sudo命令时需设置
	modes := ssh.TerminalModes{
		ssh.ECHO: 53,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return
	}

	var stdOutBuf bytes.Buffer
	session.Stdout = &stdOutBuf

	err = session.Run(command)
	if err != nil {
		return
	}

	// 去掉输出结果中末尾换行符
	stdout = strings.TrimSuffix(string(stdOutBuf.Bytes()), "\n")

	// 如果执行ls命令去掉结果中的多余空格
	if f, _:= regexp.MatchString(".*ls.*", command); f {
		stdout = strings.TrimSuffix(replaceSpace(stdout), " ")
	}
	return
}

// 替换字符串中连续多个空格为一个
func replaceSpace(str string) string {
	if str == "" {
		return ""
	}
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, " ")
}