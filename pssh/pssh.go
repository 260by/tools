package pssh

import (
	"bytes"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"regexp"
	"strings"
	// "sync"
	"time"
)

// Server ssh配置信息
type Server struct {
	Addr     string
	Port     string
	User     string
	Key      string
	KeyFile  string
	Password string
	Timeout  time.Duration
	Proxy    ProxyServer
}

// ProxyServer ssh代理配置
type ProxyServer struct {
	Addr     string
	Port     string
	User     string
	Key      string
	KeyFile  string
	Password string
	Timeout  time.Duration
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

func getSSHConfig(s ProxyServer) (*ssh.ClientConfig, io.Closer) {
	var sshAgent io.Closer

	// auths holds the detected ssh auth methods
	auths := []ssh.AuthMethod{}

	// figure out what auths are requested, what is supported
	if s.Password != "" {
		auths = append(auths, ssh.Password(s.Password))
	}
	if s.KeyFile != "" {
		if pubkey, err := getPrivateKeyFile(s.KeyFile); err != nil {
			log.Printf("Get private key file error: %v\n", err)
		} else {
			auths = append(auths, ssh.PublicKeys(pubkey))
		}
	}

	if s.Key != "" {
		if signer, err := ssh.ParsePrivateKey([]byte(s.Key)); err != nil {
			log.Printf("Parse private key error: %v\n", err)
		} else {
			auths = append(auths, ssh.PublicKeys(signer))
		}
	}

	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers))
	}

	return &ssh.ClientConfig{
		Timeout:         s.Timeout,
		User:            s.User,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, sshAgent
}

// Connect 连接到远程服务器并返回*ssh.Session
func (s *Server) Connect() (*ssh.Client, error) {
	var client *ssh.Client
	var err error

	targetServer, closer := getSSHConfig(ProxyServer{
		User:     s.User,
		Key:      s.Key,
		KeyFile:  s.KeyFile,
		Password: s.Password,
		Timeout:  s.Timeout,
	})
	if closer != nil {
		defer closer.Close()
	}

	// 开启ssh代理
	if s.Proxy.Addr != "" {
		proxyServer, closer := getSSHConfig(ProxyServer{
			User:     s.Proxy.User,
			Key:      s.Proxy.Key,
			KeyFile:  s.Proxy.KeyFile,
			Password: s.Proxy.Password,
			Timeout:  s.Proxy.Timeout,
		})
		if closer != nil {
			defer closer.Close()
		}

		proxyClient, err := ssh.Dial("tcp", net.JoinHostPort(s.Proxy.Addr, s.Proxy.Port), proxyServer)
		if err != nil {
			return nil, err
		}

		conn, err := proxyClient.Dial("tcp", net.JoinHostPort(s.Addr, s.Port))
		if err != nil {
			return nil, err
		}

		ncc, chans, reqs, err := ssh.NewClientConn(conn, net.JoinHostPort(s.Addr, s.Port), targetServer)
		if err != nil {
			return nil, err
		}

		client = ssh.NewClient(ncc, chans, reqs)
	} else {
		client, err = ssh.Dial("tcp", net.JoinHostPort(s.Addr, s.Port), targetServer)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

// Command 执行命令
func (s *Server) Command(command string) (stdout string, err error) {
	client, err := s.Connect()
	if err != nil {
		return stdout, err
	}
	session, err := client.NewSession()
	if err != nil {
		return stdout, err
	}

	// 创建伪终端, 执行sudo命令时需设置
	modes := ssh.TerminalModes{
		ssh.ECHO:          53,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return stdout, err
	}

	var stdOutBuf bytes.Buffer
	session.Stdout = &stdOutBuf

	err = session.Run(command)
	if err != nil {
		return stdout, err
	}

	// 去掉输出结果中末尾换行符
	stdout = strings.TrimSuffix(string(stdOutBuf.Bytes()), "\n")

	// 如果执行ls命令去掉结果中的多余空格,并返回以空格为分割符的字符串
	if f, _ := regexp.MatchString(".*ls.*", command); f {
		stdout = strings.TrimSuffix(replaceSpace(stdout), " ")
	}
	return stdout, err
}

// 替换字符串中连续多个空格为一个
func replaceSpace(str string) string {
	if str == "" {
		return ""
	}
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, " ")
}

// Get 使用sftp从远程服务器下载文件
func (s *Server) Get(src, dst string) (result bool, err error) {
	cmd := fmt.Sprintf("ls %s", src)
	stdout, err := s.Command(cmd)
	if err != nil {
		return false, err
	}
	remoteFiles := strings.Split(stdout, " ")

	sshClient, err := s.Connect()
	if err != nil {
		return false, err
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return false, err
	}

	// var wg sync.WaitGroup
	for _, file := range remoteFiles {
		// wg.Add(1)

		// go func(file string) {
		// 	defer wg.Add(-1)

		var localDir = dst
		var remoteFilePath string
		if strings.HasPrefix(file, "/") {
			remoteFilePath = file
		} else {
			remoteFilePath = path.Join(src, file)
		}

		srcFile, err := sftpClient.Open(remoteFilePath)
		if err != nil {
			return false, err
		}
		defer srcFile.Close()

		var localFileName = path.Base(remoteFilePath)
		dstFile, err := os.Create(path.Join(localDir, localFileName))
		if err != nil {
			return false, err
		}
		defer dstFile.Close()

		if _, err = srcFile.WriteTo(dstFile); err != nil {
			return false, err
		}
		// }(file)
	}
	// wg.Wait()

	defer sftpClient.Close()
	return true, nil
}

// Put 使用sftp从本地上传文件到远程服务器
func (s *Server) Put(src, dst string) (result bool, err error) {
	sshClient, err := s.Connect()
	if err != nil {
		return false, err
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return false, err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return false, err
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(src)
	dstFile, err := sftpClient.Create(path.Join(dst, remoteFileName))
	if err != nil {
		return false, err
	}
	defer dstFile.Close()

	f, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return false, err
	}

	dstFile.Write(f)
	// fmt.Printf("%s Upload file to remote finished!", src)

	return true, nil
}
