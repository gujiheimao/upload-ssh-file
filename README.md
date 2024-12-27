# `uf` - 一个简单的配置驱动文件上传工具

`uf` 是一个命令行工具，用于根据 YAML 配置文件上传文件和执行指定的脚本。它支持多个文件或文件夹上传，并允许在远程服务器上执行脚本。该工具能够自动配置环境，简化服务器管理。

## 主要功能

- **提示说明**：通过 `uf help` 获取文件功能和配置文件的说明。
- **生成配置文件模板**：通过 `uf template` 生成配置文件模板。
- **上传文件和执行脚本**：通过 `uf`或者 `uf <xxx.yml>` 执行文件上传和脚本执行。
- **第一次执行**:如果第一次执行uf时，会自动生成模板文件。
- **SSH 执行脚本**：支持通过 SSH 执行远程脚本，并通过开关来选择是否删除上传的脚本。

## PS：环境变量配置功能暂时未完成，请不要尝试

## 使用说明

### 1. 初始化配置文件模板

在项目目录下运行以下命令，生成配置文件模板：

```bash
uf template
```

该命令将在当前目录下创建一个 `uf.yml` 配置文件。

配置好配置文件之后，就可以通过 `uf` 命令执行上传操作。


### 2. 配置文件说明

`uf.yml` 文件采用 YAML 格式，包含以下字段：

```yaml
# 配置文件示例
server:
  host: "<your_host>" # 服务器地址
  port: 22
  username: "<your_username>" # 用户名
  password: "<your_password>" # 密码
  upload_target: "<upload_target>" #文件上传到的目标路径
  upload_files: # 你需要上传的文件列表，这里可以使用通配符
    - "buildGo.bat"
    - "./*.go"
    - "./**/*.go"
  script:
    # 设置为 true 执行内联脚本(scriptContent)，false 则上传路径下的文件并执行(scriptPath)
    executeScript: true
    scriptContent: |
      echo "Hello, World!" > /opt/dev/hello.txt
    scriptPath: "<remote_script_path>" # 如果 executeScript 为 false，则使用该路径上传脚本
```

### 3. 根据配置上传文件并执行脚本

在配置文件准备好之后，可以通过以下命令上传文件并执行配置中的脚本：

```bash
uf <config_url>
```

### 4. 安装并配置环境

如果您没有将 `uf` 添加到环境变量中，执行 `uf` 命令时，它将自动进行安装并配置环境：

```bash
uf
```

## 配置文件示例

```yaml
files:
  - source: "./data/*"
    destination: "/home/user/data"
  - source: "./logs"
    destination: "/home/user/logs"

scripts:
  - script_path: "./scripts/setup.sh"
    remote_path: "/home/user/setup.sh"
    execute_on_ssh: true
    delete_after_execution: true
  - script_path: "./scripts/cleanup.sh"
    remote_path: "/home/user/cleanup.sh"
    execute_on_ssh: false
    delete_after_execution: false
```

### 说明
- 本地路径 `./data/*` 表示将本地 `data` 文件夹中的所有文件上传到远程服务器。
- 配置文件中的脚本可以选择在上传后通过 SSH 执行。
- 执行后的脚本如果配置为 `delete_after_execution: true`，则在执行完后会被删除。

## 示例操作

1. 初始化配置文件模板：
   ```bash
   uf template
   ```

2. 配置文件示例：
   编辑生成的 `uf.yml` 配置文件，添加要上传的文件和要执行的脚本。

3. 上传文件并执行脚本：
   ```bash
   uf uf.yml
   ```