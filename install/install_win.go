//go:build windows

package install

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 安装程序，将二进制文件复制到指定路径，并更新系统 PATH 环境变量
func InstallProgram() error {
	// 检查是否存在编译后的二进制文件
	if _, err := os.Stat("uf.exe"); os.IsNotExist(err) {
		return fmt.Errorf("编译后的程序 uf.exe 不存在")
	}

	// 设置目标安装路径
	installDir := "C:\\Program Files\\uf"
	installFilePath := filepath.Join(installDir, "uf.exe")

	// 创建目标目录
	if err := os.MkdirAll(installDir, os.ModePerm); err != nil {
		return fmt.Errorf("创建安装目录失败: %v", err)
	}

	// 使用 os.Copy 来复制文件（跨磁盘）
	err := copyFile("uf.exe", installFilePath)
	if err != nil {
		return fmt.Errorf("无法将 uf.exe 复制到目标目录: %v", err)
	}
	fmt.Println("程序成功安装到", installFilePath)

	// 更新系统的 PATH 环境变量
	err = addToSystemPath(installDir)
	if err != nil {
		return fmt.Errorf("无法更新系统 PATH: %v", err)
	}

	fmt.Println("系统 PATH 环境变量已更新，uf 已添加到 PATH 中。")
	return nil
}

// 复制文件
func copyFile(src, dst string) error {
	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %v", err)
	}
	defer sourceFile.Close()

	// 创建目标文件
	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer destinationFile.Close()

	// 使用 io.Copy 来复制文件内容
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("复制文件失败: %v", err)
	}

	// 设置目标文件的权限（可选）
	err = os.Chmod(dst, 0755)
	if err != nil {
		return fmt.Errorf("设置目标文件权限失败: %v", err)
	}

	return nil
}

// 使用管理员权限更新系统 PATH 环境变量
//func addToSystemPath(installDir string) error {
//	// 通过 `setx` 命令更新系统级的 PATH
//	cmd := exec.Command("setx", "PATH", "%PATH%;"+installDir, "/M") // /M 表示更新系统级 PATH
//	err := cmd.Run()
//	if err != nil {
//		return fmt.Errorf("更新 PATH 环境变量失败: %v", err)
//	}
//	fmt.Println("系统 PATH 环境变量更新成功，请重启命令行窗口以使更改生效。")
//	return nil
//}

// 将 uf 程序的安装目录添加到系统的 PATH 环境变量
func addToSystemPath(installDir string) error {
	// 获取当前系统的 PATH 环境变量
	key, err := registry.OpenKey(registry.CLASSES_ROOT, `Environment`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("打开注册表失败: %v", err)
	}
	defer key.Close()

	// 查询现有的 PATH 环境变量
	path, _, err := key.GetStringValue("PATH")
	if err != nil {
		return fmt.Errorf("获取 PATH 环境变量失败: %v", err)
	}

	// 检查 uf 程序的安装目录是否已经在 PATH 中
	if strings.Contains(path, installDir) {
		fmt.Println("安装目录已经在 PATH 环境变量中，无需修改。")
		return nil
	}

	// 将安装路径添加到 PATH
	newPath := path + ";" + installDir
	err = key.SetStringValue("PATH", newPath)
	if err != nil {
		return fmt.Errorf("更新 PATH 环境变量失败: %v", err)
	}

	// 使环境变量生效（立即生效，重启命令行窗口即可）
	cmd := exec.Command("setx", "PATH", newPath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("更新 PATH 环境变量失败: %v", err)
	}

	fmt.Println("PATH 环境变量更新成功。请重新启动命令行窗口以使更改生效。")
	return nil
}

// 卸载程序
func UninstallProgram() error {
	// 删除安装的程序文件
	ufPath := "C:\\Program Files\\uf\\uf.exe"
	fmt.Println("正在删除程序：", ufPath)

	// 检查程序文件是否存在
	if _, err := os.Stat(ufPath); os.IsNotExist(err) {
		return fmt.Errorf("程序 uf 未安装，无法删除")
	}

	// 删除程序文件
	err := os.Remove(ufPath)
	if err != nil {
		return fmt.Errorf("删除程序失败: %v", err)
	}
	fmt.Println("程序 uf 已成功删除")

	// 清除 PATH 环境变量中的安装路径
	err = removeFromSystemPath("C:\\Program Files\\uf")
	if err != nil {
		return fmt.Errorf("清除环境变量失败: %v", err)
	}
	fmt.Println("从 PATH 中移除了 C:\\Program Files\\uf")

	return nil
}

// 从 PATH 环境变量中删除指定路径
func removeFromSystemPath(pathToRemove string) error {
	// 获取当前系统的 PATH 环境变量
	key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("打开注册表失败: %v", err)
	}
	defer key.Close()

	// 查询现有的 PATH 环境变量
	path, _, err := key.GetStringValue("PATH")
	if err != nil {
		return fmt.Errorf("获取 PATH 环境变量失败: %v", err)
	}

	// 检查路径是否存在于 PATH 中
	if !strings.Contains(path, pathToRemove) {
		return nil // 如果没有该路径，直接返回
	}

	// 删除路径
	paths := strings.Split(path, ";")
	var newPaths []string
	for _, p := range paths {
		if p != pathToRemove {
			newPaths = append(newPaths, p)
		}
	}

	// 将新的 PATH 环境变量写回
	newPath := strings.Join(newPaths, ";")
	err = key.SetStringValue("PATH", newPath)
	if err != nil {
		return fmt.Errorf("更新 PATH 环境变量失败: %v", err)
	}

	// 使环境变量生效（全局生效）
	cmd := exec.Command("setx", "PATH", newPath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("更新系统 PATH 环境变量失败: %v", err)
	}

	return nil
}
