package sshexec

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// ssh session

type HostSession struct {
	Username string
	Password string
	IP       string
	Signers  []ssh.Signer
	Port     int
	Auths    []ssh.AuthMethod
}

// result of the command execution

type ExecResult struct {
	Id             int
	IP             string
	Command        string
	LocalFilePath  string
	RemoteFilePath string
	Result         string
	StartTime      time.Time
	EndTime        time.Time
	Error          error
	ErrorInfo      string
}

// execute the command and return a result structure

func (exec *HostSession) Exec(id int, command string, config ssh.ClientConfig) *ExecResult {

	result := &ExecResult{
		Id:      id,
		IP:      exec.IP,
		Command: command,
	}

	client, err := ssh.Dial("tcp", exec.IP+":"+strconv.Itoa(exec.Port), &config)
	if err != nil {
		result.Error = err
		result.ErrorInfo = fmt.Sprintf("%v连接异常：%v", exec.IP, err.Error())
		return result
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		result.Error = err
		result.ErrorInfo = fmt.Sprintf("%v连接异常：%v", exec.IP, err.Error())
		return result
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	var b1 bytes.Buffer
	session.Stderr = &b1
	start := time.Now()
	if err := session.Run(command); err != nil {
		result.Error = err
		result.ErrorInfo = fmt.Sprintf("%v执行异常：%v", exec.IP, err.Error())
		result.Result = b1.String()
		return result
	}
	end := time.Now()
	result.Result = b.String()
	result.StartTime = start
	result.EndTime = end
	return result
}

// execute the command and return a result structure

func (exec *HostSession) Transfer(id int, localFilePath string, remoteFilePath string, config ssh.ClientConfig) *ExecResult {

	result := &ExecResult{
		Id:             id,
		IP:             exec.IP,
		LocalFilePath:  localFilePath,
		RemoteFilePath: remoteFilePath,
	}
	start := time.Now()
	result.StartTime = start
	client, err := ssh.Dial("tcp", exec.IP+":"+strconv.Itoa(exec.Port), &config)
	if err != nil {
		result.Error = err
		result.ErrorInfo = fmt.Sprintf("%v连接异常：%v", exec.IP, err.Error())
		return result
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		result.Error = err
		result.ErrorInfo = fmt.Sprintf("%v连接异常：%v", exec.IP, err.Error())
		return result
	}
	defer session.Close()

	var fileSize int64
	if s, err := os.Stat(localFilePath); err != nil {
		result.Error = err
		result.ErrorInfo = fmt.Sprintf("%v文件不存在：%v", localFilePath, err.Error())
		return result
	} else {
		fileSize = s.Size()
	}
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		result.Error = err
		result.ErrorInfo = fmt.Sprintf("%v文件打开异常：%v", localFilePath, err.Error())
		return result
	}

	defer srcFile.Close()

	sftpClient, err := sftp.NewClient(client)
	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	// create sftp client
	if err != nil {
		result.Error = err
		result.ErrorInfo = fmt.Sprintf("%v连接异常：%v", exec.IP, err.Error())
		return result
	}
	defer sftpClient.Close()

	dstFile, err := sftpClient.Create(remoteFilePath)
	if err != nil {
		result.Error = err
		result.ErrorInfo = fmt.Sprintf("%v创建文件%v异常：%v", exec.IP, remoteFilePath, err.Error())
		return result
	}
	defer dstFile.Close()
	//todo 这里使用buff池 或io.Copy 性能没有测试出差距
	//n, err := Copy(dstFile, io.LimitReader(srcFile, fileSize))
	n, err := io.Copy(dstFile, io.LimitReader(srcFile, fileSize))
	if err != nil {
		result.Error = err
		result.ErrorInfo = fmt.Sprintf("%v文件传输异常：%v", exec.IP, err.Error())
		return result
	}
	if n != fileSize {
		result.Error = errors.New(fmt.Sprintf("copy: expected %v bytes, got %d", fileSize, n))
		result.ErrorInfo = fmt.Sprintf("%v文件传输异常：%v", exec.IP, err.Error())
		return result
	}
	end := time.Now()
	result.EndTime = end
	return result
}

func (exec *HostSession) GenerateConfig() ssh.ClientConfig {
	var auths []ssh.AuthMethod

	if len(exec.Password) != 0 {
		auths = append(auths, ssh.Password(exec.Password))
	} else {
		if len(exec.Auths) > 0 {
			auths = exec.Auths
		} else {
			auths = append(auths, ssh.PublicKeys(exec.Signers...))
		}

	}

	config := ssh.ClientConfig{
		User: exec.Username,
		Auth: auths,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	config.Ciphers = []string{"aes128-cbc", "3des-cbc"}

	return config
}

func readKey(filename string) (ssh.Signer, error) {
	f, _ := os.Open(filename)

	bytes, _ := ioutil.ReadAll(f)
	return generateKey(bytes)
}

func generateKey(keyData []byte) (ssh.Signer, error) {
	return ssh.ParsePrivateKey(keyData)
}
