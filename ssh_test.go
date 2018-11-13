package sshexec

import (
	"testing"
	"time"
)

func Test_SshHost(t *testing.T) {
	sshParms := []SSHParm{}
	sshParms = append(sshParms, SSHParm{
		IP:        "172.16.20.45",
		Port:      22,
		Username:  "root",
		LoginType: SSH_LOGIN_PASS,
		Password:  "123456",
	})
	sshParms = append(sshParms, SSHParm{
		IP:        "172.16.20.45",
		Port:      22,
		Username:  "root",
		LoginType: SSH_LOGIN_KEY,
		Password: `
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`,
	})
	sshExecAgent := SSHExecAgent{}
	sshExecAgent.Worker = 10
	sshExecAgent.TimeOut = time.Duration(120) * time.Second
	rs, err := sshExecAgent.SshHost(sshParms, "echo $GOPATH")
	t.Log(rs, err)
}

func Test_SftpHost(t *testing.T) {
	sshParms := []SSHParm{}
	sshParms = append(sshParms, SSHParm{
		IP:        "172.16.20.45",
		Port:      22,
		Username:  "root",
		LoginType: SSH_LOGIN_PASS,
		Password:  "123456",
	})
	sshParms = append(sshParms, SSHParm{
		IP:        "172.16.20.46",
		Port:      22,
		Username:  "root",
		LoginType: SSH_LOGIN_KEY,
		Password: `
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`,
	})
	sshExecAgent := SSHExecAgent{}
	sshExecAgent.Worker = 10
	sshExecAgent.TimeOut = time.Duration(120) * time.Second
	t.Log(sshParms)
	rs, err := sshExecAgent.SftpHost(sshParms, "/data/bbbb.txt", "/data/bbbb.txt")
	t.Log(rs, err)
}
