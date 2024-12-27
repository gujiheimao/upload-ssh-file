package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

// Config 配置结构体，用来存储从配置文件读取的值
type Config struct {
	Server struct {
		Host         string   `yaml:"host"`
		Port         int      `yaml:"port"`
		Username     string   `yaml:"username"`
		Password     string   `yaml:"password"`
		UploadTarget string   `yaml:"upload_target"`
		UploadFiles  []string `yaml:"upload_files"` // 支持多个文件或文件夹
		Script       struct {
			ExecuteScript bool   `yaml:"executeScript"` // 是否直接执行脚本
			ScriptContent string `yaml:"scriptContent"` // 内联脚本内容
			ScriptPath    string `yaml:"scriptPath"`    // 脚本上传路径
		} `yaml:"script"`
	} `yaml:"server"`
}

// LoadConfig 从配置文件读取配置信息
func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("加载配置文件失败: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

// GenerateConfigTemplate 生成配置文件模板
func GenerateConfigTemplate() {
	// 配置文件模板内容，包含占位符，供用户填写
	template := `
server:
  host: "<your_host>"
  port: 22
  username: "<your_username>"
  password: "<your_password>"
  upload_target: "<upload_target>"
  upload_files:
    - "buildGo.bat"
    - "./*.go"
    - "./**/*.go"
  script:
	# 设置为 true 执行内联脚本(scriptContent)，false 则上传并执行(scriptPath)
    executeScript: true         
    scriptContent: |
      echo "Hello, World!"
      ls -l /home/user/uploads
    scriptPath: "<remote_script_path>" # 如果 executeScript 为 false，则使用该路径上传脚本
`
	// 将模板写入文件 uf.yml
	err := os.WriteFile("uf.yml", []byte(template), 0644)
	if err != nil {
		fmt.Printf("生成配置模板失败: %v\n", err)
		return
	}
	fmt.Println("配置模板 (uf.yml) 已生成，请编辑此文件。")
}
