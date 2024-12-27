package main

import (
	"fmt"
	"os"
	"strings"
	"uploadFile/config"
	"uploadFile/runshell"
	"uploadFile/upload"
)

var DefConfigPath = "uf.yml"

func main() {
	var conf *config.Config
	args := os.Args
	if len(args) == 1 {
		//判断当前目录下配置文件是否存在
		file, err := checkYAMLFile()
		if err != nil {
			panic(err)
		}
		if file {
			conf, err = config.LoadConfig(DefConfigPath)
			if err != nil {
				panic(err)
			}
			//默认配置
			Upload(conf)
		} else {
			config.GenerateConfigTemplate()
		}
		return
	}
	t := args[1]
	switch t {
	//case "install":
	//	err := install.InstallProgram()
	//	if err != nil {
	//		panic(err)
	//	}
	//case "uninstall":
	//	err := install.UninstallProgram()
	//	if err != nil {
	//		panic(err)
	//	}
	case "template":
		config.GenerateConfigTemplate()
	case "help":
		//	fmt.Println(`Usage: uf <command>:
		//install    | 安装uf程序，包括配置环境变量
		//uninstall  | 取消uf安装程序，包括环境变量
		//help       | 获取帮助提示
		//template   | 在执行命令的目录下生成uf.yml模板文件
		//*.yml/yaml | 读取其他名称的配置文件`)
		fmt.Println(`Usage: uf <command>:
	help       | 获取帮助提示
	template   | 在执行命令的目录下生成uf.yml模板文件
	*.yml/yaml | 读取其他名称的配置文件`)
	default:
		//如果是配置文件
		if strings.HasSuffix(t, "yml") || strings.HasSuffix(t, ".yaml") {
			conf, _ = config.LoadConfig(t)
			//执行命令
			Upload(conf)
			return
		}
		fmt.Println(`Usage: uf <command>:
	help       | 获取帮助提示
	template   | 在执行命令的目录下生成uf.yml模板文件
	*.yml/yaml | 读取其他名称的配置文件`)
	}
}

func Upload(conf *config.Config) {
	client, err := runshell.CreateSSHClient(conf.Server.Host, conf.Server.Port, conf.Server.Username, conf.Server.Password)
	if err != nil {
		panic(err)
	}

	err = upload.UploadFiles(client, conf.Server.UploadFiles, conf.Server.UploadTarget)
	if err != nil {
		panic(err)
	}
	fmt.Println("执行脚本")
	//执行脚本
	err = runshell.ExecuteScript(client, conf.Server.Script.ScriptContent, conf.Server.Script.ScriptPath, conf.Server.Script.ExecuteScript)
}

// 判断当前目录下的 uf.yml 或 uf.yaml 文件是否存在
func checkYAMLFile() (bool, error) {
	// 检查 uf.yml 文件是否存在
	if _, err := os.Stat("uf.yml"); err == nil {
		fmt.Println("uf.yml 文件存在")
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("检查 uf.yml 时出错: %v", err)
	}

	// 检查 uf.yaml 文件是否存在
	if _, err := os.Stat("uf.yaml"); err == nil {
		fmt.Println("uf.yaml 文件存在")
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("检查 uf.yaml 时出错: %v", err)
	}

	// 如果两个文件都不存在
	fmt.Println("uf.yml 和 uf.yaml 文件都不存在")
	return false, nil
}
