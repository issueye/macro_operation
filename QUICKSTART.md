# 快速开始指南

## 项目已成功编译！

可执行文件位置: `bin/macro-recorder.exe`

## 如何使用

### 1. 运行程序

在 Windows 命令行中执行：

```bash
cd E:\code\issueye\test\macro_operation
.\bin\macro-recorder.exe
```

### 2. 菜单操作

程序启动后会显示以下菜单：

```
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
```

### 3. 快速测试（推荐使用选项7）

1. 输入 `7` 并按 Enter
2. 输入宏名称（例如：`test_macro`）并按 Enter
3. 程序显示"开始录制..."
4. 执行一些操作（移动鼠标、点击、键盘输入）
5. 按 Enter 停止录制
6. 程序会自动保存并显示生成的JavaScript脚本

### 4. 查看录制的宏

录制的宏保存在 `macros/` 目录下，文件格式为 JSON：

```json
{
  "meta": {
    "name": "test_macro",
    "version": "1.0.0",
    "created_at": "2024-12-31T15:00:00Z"
  },
  "script": {
    "language": "javascript",
    "code": "function main() {\n  mouseMove(100, 200);\n  ...\n}"
  }
}
```

### 5. 回放宏

1. 输入 `4` 并按 Enter
2. 选择或输入要播放的宏名称
3. 程序会自动执行录制的操作

## 技术特性

### 支持的操作

#### 鼠标操作
- ✅ 鼠标移动 (mousemove)
- ⚠️ 鼠标点击 (待完善)
- ✅ 鼠标拖拽
- ✅ 鼠标滚动

#### 键盘操作
- ✅ 按键按下/释放
- ✅ 文本输入
- ✅ 组合键

### JavaScript API

生成的脚本支持以下API：

```javascript
// 鼠标操作
mouseMove(x, y)
mouseClick("left"|"right"|"middle")
mouseDrag(x, y)
mouseScroll(delta)

// 键盘操作
keyDown("key")
keyUp("key")
keyType("text")
keyTap("key")

// 系统操作
sleep(ms)
log("message")
```

## 故障排查

### 问题1: 程序无法启动
- 检查是否有杀毒软件阻止
- 确保在 Windows 系统上运行
- 尝试以管理员身份运行

### 问题2: 录制没有反应
- 确保选择了正确的录制选项
- 检查是否有权限问题
- 尝试重新启动程序

### 问题3: 回放不准确
- 检查屏幕分辨率是否与录制时相同
- 确保没有其他程序干扰
- 尝试重新录制

## 开发者信息

### 重新编译

```bash
# 清理旧的构建
make clean

# 重新编译
make build
```

### 运行测试

```bash
# 运行所有测试
make test

# 查看代码覆盖率
make test-coverage
```

### 代码检查

```bash
# 格式化代码
make fmt

# 代码检查
make lint
```

## 项目文件说明

```
macro_operation/
├── bin/
│   └── macro-recorder.exe    # 可执行文件 (32MB)
├── macros/                    # 宏文件存储目录
├── configs/
│   ├── config.go             # 配置管理代码
│   └── config.yaml           # 配置文件
├── cmd/macro/
│   └── main.go               # 程序入口
├── internal/
│   ├── engine/               # 核心引擎
│   │   ├── events/           # 事件捕获
│   │   ├── generator/        # 脚本生成
│   │   └── executor/         # 脚本执行
│   ├── service/              # 业务服务
│   ├── repository/           # 数据存储
│   └── model/                # 数据模型
├── pkg/bindings/             # JavaScript API绑定
├── docs/                     # 文档
├── Makefile                  # 构建脚本
├── go.mod                    # Go模块定义
└── README.md                 # 项目说明
```

## 下一步

1. **测试功能**: 尝试录制并回放一些简单操作
2. **查看文档**: 阅读 `docs/` 目录下的详细文档
3. **二次开发**: 根据需求修改和扩展功能
4. **贡献代码**: 提交 Issue 或 Pull Request

## 联系方式

- 项目地址: E:\code\issueye\test\macro_operation
- 开发者: Golang 后端开发工程师
- 版本: v1.0.0 (MVP)
