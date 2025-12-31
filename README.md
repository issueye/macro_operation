# 操作宏录制和回放

一个强大的桌面自动化工具，可以录制键盘和鼠标操作并生成可重复执行的宏脚本。

## 功能特性

- **录制功能**: 实时捕获键盘和鼠标操作
- **脚本生成**: 自动生成可读的JavaScript脚本
- **回放执行**: 精准回放录制的操作
- **宏管理**: 保存、加载、编辑和管理宏
- **跨平台**: 支持 Windows、macOS 和 Linux

## 技术栈

- **语言**: Go 1.21+
- **键盘监听**: [gohook](https://github.com/robotn/gohook)
- **鼠标监听**: [robotgo](https://github.com/go-vgo/robotgo)
- **脚本引擎**: [goja](https://github.com/dop251/goja) (JavaScript)
- **配置管理**: [viper](https://github.com/spf13/viper)

## 快速开始

### 环境要求

- Go 1.21 或更高版本
- Windows / macOS / Linux

### 安装

\`\`\`bash
# 克隆仓库
git clone https://github.com/issueye/macro_operation.git
cd macro_operation

# 安装依赖
make deps

# 构建
make build
\`\`\`

### 运行

\`\`\`bash
# 直接运行
make run

# 或运行构建的二进制
./bin/macro-recorder
\`\`\`

## 使用说明

### 命令行界面

程序启动后会显示菜单：

\`\`\`
========== 操作宏录制和回放 ==========
1. 开始录制
2. 停止录制
3. 保存当前录制的宏
4. 播放宏
5. 列出所有宏
6. 删除宏
7. 录制并保存 (一键操作)
0. 退出
======================================
\`\`\`

### 基本使用流程

1. **开始录制**: 选择菜单 `1` 开始录制操作
2. **执行操作**: 执行你想要录制的键盘和鼠标操作
3. **停止录制**: 选择菜单 `2` 停止录制
4. **保存宏**: 选择菜单 `3` 保存为宏，输入宏名称
5. **播放宏**: 选择菜单 `4` 选择并播放已保存的宏

### 快速录制

使用菜单 `7` 可以一键完成录制和保存：
1. 输入宏名称
2. 执行操作
3. 按 Enter 停止
4. 自动保存

## JavaScript API

录制的操作会转换为JavaScript脚本，支持以下API：

### 鼠标操作

\`\`\`javascript
mouseMove(x, y)           // 移动鼠标到指定位置
mouseClick(button)        // 鼠标点击 (left/right/middle)
mouseDrag(x, y)           // 鼠标拖拽
mouseScroll(delta)        // 鼠标滚轮
\`\`\`

### 键盘操作

\`\`\`javascript
keyDown(key)              // 按下按键
keyUp(key)                // 释放按键
keyType(text)             // 输入文本
keyTap(key)               // 按键点击（按下+释放）
\`\`\`

### 系统操作

\`\`\`javascript
sleep(ms)                 // 延迟等待
screenshot(path)          // 屏幕截图
log(message)              // 日志输出
\`\`\`

### 示例脚本

\`\`\`javascript
// 宏名称: 数据录入
// 录制时间: 2024-12-31 14:30:00
// 操作数量: 5

function main() {
  // 鼠标移动到 (100, 200)
  mouseMove(100, 200);
  sleep(150);
  // 鼠标点击 (100, 200) 按钮: left
  mouseClick('left');
  sleep(500);
  // 输入文本: Hello
  keyTap('h');
  keyTap('e');
  keyTap('l');
  keyTap('l');
  keyTap('o');
}

// 执行宏
main();
\`\`\`

## 配置文件

配置文件位于 `configs/config.yaml`：

\`\`\`yaml
app:
  name: "操作宏录制和回放"
  version: "1.0.0"
  debug: true
  log_path: "./logs"

record:
  max_duration: 3600
  enable_screenshot: false
  filter_mouse_move: true

playback:
  default_speed: 1.0
  enable_sound: true
  stop_on_error: true

storage:
  macros_dir: "./macros"
  backup_enabled: true
  backup_dir: "./backups"
  max_backup_count: 10
\`\`\`

## 项目结构

\`\`\`
macro_operation/
├── cmd/
│   └── macro/              # 程序入口
│       └── main.go
├── internal/
│   ├── engine/             # 核心引擎
│   │   ├── events/         # 事件捕获
│   │   ├── generator/      # 脚本生成器
│   │   └── executor/       # 脚本执行器
│   ├── service/            # 业务服务
│   │   ├── record_service.go
│   │   ├── play_service.go
│   │   └── macro_service.go
│   ├── repository/         # 数据存储
│   │   └── macro_repo.go
│   └── model/              # 数据模型
│       ├── event.go
│       └── macro.go
├── pkg/
│   └── bindings/           # JS API绑定
│       └── api_bindings.go
├── configs/                # 配置文件
│   ├── config.go
│   └── config.yaml
├── macros/                 # 宏存储目录
├── docs/                   # 文档
├── Makefile
├── go.mod
└── README.md
\`\`\`

## 开发

\`\`\`bash
# 运行测试
make test

# 代码检查
make check

# 格式化代码
make fmt

# 构建所有平台
make build-all
\`\`\`

## 文档

- [产品需求分析](docs/01-产品需求分析.md)
- [技术实现方案](docs/02-技术实现方案.md)
- [核心功能拆解](docs/03-核心功能拆解.md)
- [项目启动清单](docs/04-项目启动清单.md)

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 作者

issueye

## 更新日志

### v1.0.0 (2024-12-31)

- 初始版本发布
- 支持键盘和鼠标事件录制
- 支持JavaScript脚本生成和执行
- 支持宏的保存和管理
- 提供命令行界面
