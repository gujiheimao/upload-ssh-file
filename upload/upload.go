package upload

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"path/filepath"
)

// 上传单个文件
func UploadFile(client *ssh.Client, localFile, remotePath string) error {
	// 创建 SFTP 客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("创建 SFTP 客户端失败: %v", err)
	}
	defer sftpClient.Close()

	// 确保远程路径存在
	err = ensureRemoteDirExists(client, getLastPartBeforeSlash(remotePath))
	if err != nil {
		return fmt.Errorf("确保远程目录存在失败: %v", err)
	}

	// 检查远程文件是否存在，如果存在则删除
	_, err = sftpClient.Stat(remotePath)
	if err == nil {
		// 文件已存在，进行删除
		fmt.Println("文件已存在，正在覆盖...")
		err := sftpClient.Remove(remotePath)
		if err != nil {
			return fmt.Errorf("删除远程文件失败: %v", err)
		}
	}

	// 打开本地文件
	file, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("打开本地文件失败: %v", err)
	}
	defer file.Close()

	// 创建远程文件
	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("创建远程文件失败: %v %s", err, remotePath)
	}
	defer remoteFile.Close()

	// 将本地文件内容复制到远程文件
	_, err = io.Copy(remoteFile, file)
	if err != nil {
		return fmt.Errorf("文件上传失败: %v", err)
	}

	fmt.Println("文件已成功上传到:", remotePath)
	return nil
}

// 根据反斜杠截取最后一个斜杠前的字符串
func getLastPartBeforeSlash(path string) string {
	// 找到最后一个反斜杠的位置
	lastSlashIndex := strings.LastIndex(path, "/")

	// 如果没有找到反斜杠，返回原始字符串
	if lastSlashIndex == -1 {
		return path
	}

	// 截取反斜杠前的部分
	return path[:lastSlashIndex]
}

// UploadFolder 上传文件夹
func UploadFolder(client *ssh.Client, localDir, remotePath string) error {
	// 创建 SFTP 客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("创建 SFTP 客户端失败: %v", err)
	}
	defer sftpClient.Close()

	// 检查远程文件夹是否存在，如果存在则删除
	_, err = sftpClient.Stat(remotePath)
	if err == nil {
		// 文件夹已存在，进行删除
		fmt.Println("文件夹已存在，正在覆盖...")
		err := sftpClient.RemoveDirectory(remotePath)
		if err != nil {
			return fmt.Errorf("删除远程文件夹失败: %v", err)
		}
	}

	// 获取本地文件夹中的所有文件
	files, err := ioutil.ReadDir(localDir)
	if err != nil {
		return fmt.Errorf("读取文件夹 %s 内容失败: %v", localDir, err)
	}

	// 创建远程文件夹
	err = sftpClient.Mkdir(remotePath)
	if err != nil {
		return fmt.Errorf("创建远程文件夹失败: %v", err)
	}

	// 遍历本地文件夹中的所有文件和文件夹
	for _, file := range files {
		// 拼接文件路径
		localFilePath := filepath.Join(localDir, file.Name())
		remoteFilePath := filepath.Join(remotePath, file.Name())

		if file.IsDir() {
			// 如果是文件夹，递归上传文件夹内容
			err := UploadFolder(client, localFilePath, remoteFilePath)
			if err != nil {
				return fmt.Errorf("上传文件夹内容失败: %v", err)
			}
		} else {
			// 上传文件
			err := UploadFile(client, localFilePath, remoteFilePath)
			if err != nil {
				return fmt.Errorf("上传文件失败: %v", err)
			}
		}
	}
	return nil
}

// 上传多个文件或文件夹
func UploadFiles(client *ssh.Client, files []string, remotePath string) error {
	// 处理通配符路径，例如 "./*"
	var allFiles []string
	for _, file := range files {
		// 使用 filepath.Glob 来处理通配符
		matches, err := filepath.Glob(file)
		if err != nil {
			return fmt.Errorf("处理通配符失败: %v", err)
		}
		allFiles = append(allFiles, matches...)
	}

	// 确保远程目录存在
	err := ensureRemoteDirExists(client, remotePath)
	if err != nil {
		return fmt.Errorf("确保远程目录存在失败: %v", err)
	}

	// 上传每个文件
	for _, file := range allFiles {
		// 检查文件或文件夹
		info, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("检查文件/文件夹 %s 失败: %v", file, err)
		}

		if info.IsDir() {
			// 如果是文件夹，递归上传文件夹内容
			err := UploadFolder(client, file, remotePath)
			if err != nil {
				return fmt.Errorf("上传文件夹失败: %v", err)
			}
		} else {
			// 上传单个文件
			remoteFile := strings.ReplaceAll(remotePath+"/"+file, "\\", "/") //filepath.Join(remotePath, filepath.Base(file))
			err := UploadFile(client, file, remoteFile)
			if err != nil {
				return fmt.Errorf("上传文件失败: %v", err)
			}
		}
	}
	return nil
}

// 判断并确保远程目录存在（通过 SSH 创建目录）
func ensureRemoteDirExists(client *ssh.Client, remotePath string) error {
	// 创建 SSH 会话
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建 SSH 会话失败: %v", err)
	}
	defer session.Close()

	// 检查远程路径是否存在
	// 如果路径不存在，则创建该目录
	cmd := fmt.Sprintf("mkdir -p %s", remotePath)
	_, err = session.CombinedOutput(cmd)
	if err != nil {
		return fmt.Errorf("创建远程目录失败: %v", err)
	}
	fmt.Printf("远程目录创建成功: %s \n", remotePath)
	return nil
}
