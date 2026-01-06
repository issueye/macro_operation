# 快捷键功能文档

## 概述

宏录制引擎支持全局快捷键控制，可以在录制过程中使用快捷键快速执行常见操作。

## 内置快捷键

### RecordService 默认快捷键

| 快捷键 | 功能 | 说明 |
|--------|------|------|
| `Ctrl+R` | 开始/停止录制 | 切换录制状态 |
| `Ctrl+S` | 保存 | 触发保存回调（需要通过 `SetActionCallback` 设置） |
| `Ctrl+Shift+C` | 清空录制 | 清空当前录制的事件 |
| `Ctrl+Shift+R` | 停止并保存 | 停止录制并触发保存回调 |

## 常用编辑快捷键

以下快捷键会被正常记录到脚本中：

| 快捷键 | 功能 |
|--------|------|
| `Ctrl+C` | 复制 |
| `Ctrl+V` | 粘贴 |
| `Ctrl+X` | 剪切 |
| `Ctrl+Z` | 撤销 |
| `Ctrl+Y` | 重做 |
| `Ctrl+A` | 全选 |

## 使用示例

### 基础使用

```go
package main

import (
    "macro-engine/internal/service"
)

func main() {
    // 创建录制服务（自动注册默认快捷键）
    recordService := service.NewRecordService()

    // 设置快捷键回调
    recordService.SetActionCallback(service.HotkeyActionSave, func() {
        // 保存逻辑
        script, _ := recordService.GenerateScript("my_macro")
        println("脚本已保存:", script)
    })

    // 启动事件监听
    capture := recordService.GetCapture()
    capture.Start()

    // 现在可以使用 Ctrl+R 开始/停止录制
    // 使用 Ctrl+S 保存
    // 使用 Ctrl+Shift+C 清空
}
```

### 注册自定义快捷键

```go
recordService := service.NewRecordService()
capture := recordService.GetCapture()

// 注册播放快捷键
capture.RegisterHotkey("ctrl+p", func() {
    println("播放宏")
    // 播放逻辑
})

// 注册删除快捷键
capture.RegisterHotkey("ctrl+d", func() {
    println("删除最后一个事件")
    // 删除逻辑
})

capture.Start()
```

### 支持的快捷键格式

快捷键字符串格式：`修饰键1+修饰键2+主键`

**修饰键（按顺序）：**
- `ctrl` - Ctrl 键
- `shift` - Shift 键
- `alt` - Alt 键
- `win` - Windows 键

**主键：**
- 字母键：`a` 到 `z`
- 数字键：`0` 到 `9`
- 功能键：`f1` 到 `f12`
- 方向键：`up`, `down`, `left`, `right`
- 特殊键：`enter`, `space`, `escape`, `tab`, `backspace`, `delete`

**示例：**
- `ctrl+r` - Ctrl + R
- `ctrl+shift+r` - Ctrl + Shift + R
- `alt+f4` - Alt + F4
- `ctrl+shift+alt+f1` - Ctrl + Shift + Alt + F1
- `ctrl+shift+s` - Ctrl + Shift + S

## 快捷键动作回调

可以通过 `SetActionCallback` 方法为特定动作设置回调函数：

```go
type HotkeyAction string

const (
    HotkeyActionToggleRecording HotkeyAction = "toggle_recording"
    HotkeyActionSave            HotkeyAction = "save"
    HotkeyActionClear           HotkeyAction = "clear"
)

// 设置回调
recordService.SetActionCallback(service.HotkeyActionSave, func() {
    // 当按下 Ctrl+S 时执行
})
```

## 注意事项

1. **全局钩子**：快捷键使用全局键盘钩子，会在整个系统中生效
2. **快捷键冲突**：避免与其他应用的全局快捷键冲突
3. **权限要求**：在某些系统上可能需要管理员权限
4. **记录机制**：快捷键本身也会被记录到事件中，除非在回调中明确过滤
5. **线程安全**：回调函数在独立的 goroutine 中执行，需要注意线程安全

## 实现细节

### 事件捕获流程

1. `gohook.Start()` 启动全局钩子
2. 每个键盘事件触发 `updateModifierKeys()` 更新修饰键状态
3. KeyDown 事件触发 `checkHotkey()` 检测快捷键组合
4. 匹配成功时在独立 goroutine 中执行回调函数
5. 事件继续被记录到事件列表

### 修饰键处理

修饰键使用 keycode 映射：
```go
16, 160  // Shift (左右)
17, 162  // Ctrl (左右)
18, 164  // Alt (左右)
91, 92   // Win (左右)
```

## 测试

运行示例代码测试快捷键：

```bash
cd macro-engine
go run internal/service/record_service_hotkey_example.go
```

## 参考

- 事件捕获器实现：`internal/engine/events/capture.go`
- 录制服务实现：`internal/service/record_service.go`
- 示例代码：`internal/service/record_service_hotkey_example.go`
