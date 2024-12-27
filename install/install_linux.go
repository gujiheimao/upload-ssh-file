//go:build linux

package install

// 将程序复制到系统的全局目录（/usr/local/bin）
func InstallProgram() error {
	// 检查文件是否存在
	if _, err := os.Stat("uf"); os.IsNotExist(err) {
		return fmt.Errorf("编译后的程序 uf 不存在")
	}

	// 将二进制文件复制到 /usr/local/bin
	fmt.Println("正在将程序复制到 /usr/local/bin ...")
	cmd := exec.Command("sudo", "cp", "uf", "/usr/local/bin/")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("将程序复制到 /usr/local/bin 失败: %v", err)
	}
	fmt.Println("程序成功安装到 /usr/local/bin/uf")
	return nil
}

// 取消安装程序
func UninstallProgram() error {
	// 删除 /usr/local/bin/uf 二进制文件
	ufPath := "/usr/local/bin/uf"
	fmt.Println("正在删除程序：", ufPath)

	// 检查程序文件是否存在
	if _, err := os.Stat(ufPath); os.IsNotExist(err) {
		return fmt.Errorf("程序 uf 未安装，无法删除")
	}

	// 删除 uf 文件
	err := os.Remove(ufPath)
	if err != nil {
		return fmt.Errorf("删除程序失败: %v", err)
	}
	fmt.Println("程序 uf 已成功删除")

	// 清除 PATH 环境变量中可能添加的安装路径
	err = removeFromPath("/usr/local/bin")
	if err != nil {
		return fmt.Errorf("清除环境变量失败: %v", err)
	}
	fmt.Println("/usr/local/bin 从 PATH 中已删除")

	return nil
}

// 从 PATH 环境变量中删除指定路径
func removeFromPath(pathToRemove string) error {
	// 获取当前用户的 PATH 环境变量
	path := os.Getenv("PATH")

	// 检查是否包含该路径
	if !strings.Contains(path, pathToRemove) {
		return nil // 如果没有该路径，直接返回
	}

	// 删除该路径
	paths := strings.Split(path, ":")
	newPaths := []string{}
	for _, p := range paths {
		if p != pathToRemove {
			newPaths = append(newPaths, p)
		}
	}

	// 将新的 PATH 环境变量写回
	newPath := strings.Join(newPaths, ":")
	err := os.Setenv("PATH", newPath)
	if err != nil {
		return fmt.Errorf("更新 PATH 环境变量失败: %v", err)
	}

	// 在系统中更新 PATH（全局生效）
	cmd := exec.Command("export", fmt.Sprintf("PATH=%s", newPath))
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("更新系统 PATH 环境变量失败: %v", err)
	}

	return nil
}
