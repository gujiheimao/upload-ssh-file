package runshell

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
	"uploadFile/upload"
)

// CreateSSHClient 创建 SSH 客户端连接
func CreateSSHClient(host string, port int, username, password string) (*ssh.Client, error) {
	// 配置 SSH 客户端
	sshConfig := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略主机密钥验证
		Timeout:         5 * time.Second,             // 设置连接超时
	}

	// 建立与服务器的连接
	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("无法连接到 SSH 服务器: %v", err)
	}

	return client, nil
}

// 执行脚本
func ExecuteScript(client *ssh.Client, scriptContent string, scriptPath string, executeScript bool) error {
	var cmd string

	if executeScript {
		// 如果 executeScript 为 true，直接执行内联脚本
		cmd = scriptContent
	} else {
		// 否则，上传脚本到服务器并执行
		err := upload.UploadFile(client, scriptPath, "/tmp/remote-script.sh")
		if err != nil {
			return fmt.Errorf("上传脚本失败: %v", err)
		}

		// 执行上传的脚本
		cmd = fmt.Sprintf("bash /tmp/remote-script.sh")
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return fmt.Errorf("执行脚本失败: %v, 输出: %s", err, output)
	}

	fmt.Println("脚本执行成功，输出:", string(output))
	return nil
}
