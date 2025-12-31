# 操作宏录制和回放

## 功能描述
操作宏录制和回放功能允许用户录制一系列操作，并将它们保存为宏，以便在需要时重复执行。这可以大大提高工作效率，尤其是在需要重复执行相同操作的情况下。

## 操作步骤
1. 打开操作宏录制和回放功能。
2. 点击“录制”按钮开始录制操作。
3. 执行一系列操作。
4. 点击“停止”按钮停止录制操作。
5. 点击“保存”按钮将录制的操作保存为宏。
6. 点击“播放”按钮执行保存的宏。

## 技术栈

- golang
  - https://github.com/robotn/gohook (键盘监听)
  - https://github.com/go-vgo/robotgo （鼠标监听）
  - https://github.com/go-vgo/robotgo （屏幕截图）
  - https://github.com/dop251/goja （JavaScript引擎）
  
## 技术说明

1. 将 gohook robotgo api 封装提供给 goja 中
2. 将操作动态生成 js 脚本
3. 通过执行 js 脚本回放录制的宏

## 技术实现

1. 使用 gohook 监听键盘事件，使用 robotgo 监听鼠标事件，使用 robotgo 截取屏幕，使用 goja 执行 js 脚本。
2. 将操作动态生成 js 脚本，通过 goja 执行 js 脚本回放录制的宏。
3. 将录制的操作保存为宏，以便在需要时重复执行。
