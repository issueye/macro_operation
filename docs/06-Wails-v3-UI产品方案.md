# Wails v3 UI 产品需求分析与方案设计

> **文档版本**: v1.0.0
> **创建日期**: 2025-12-31
> **产品经理**: Claude (Product Manager Agent)
> **项目**: 操作宏录制和回放工具 - GUI 版本

---

## 目录

1. [项目概述](#1-项目概述)
2. [Wails v3 技术分析](#2-wails-v3-技术分析)
3. [用户需求分析](#3-用户需求分析)
4. [UI 界面设计](#4-ui-界面设计)
5. [前后端 API 设计](#5-前后端-api-设计)
6. [功能模块拆解](#6-功能模块拆解)
7. [开发路线图](#7-开发路线图)
8. [技术风险评估](#8-技术风险评估)
9. [附录](#9-附录)

---

## 1. 项目概述

### 1.1 产品定位

**操作宏录制和回放工具** 是一款桌面自动化工具，旨在帮助用户录制键盘和鼠标操作，生成可重复执行的宏脚本，从而提高工作效率。

**当前状态**: 已完成 CLI 版本，具备完整的核心功能。
**目标**: 开发基于 Wails v3 的图形用户界面，提供更直观、易用的交互体验。

### 1.2 核心价值主张

| 用户痛点 | 解决方案 | 价值体现 |
|---------|---------|---------|
| 重复性操作繁琐 | 一键录制并回放操作 | 节省时间，提高效率 |
| CLI 操作门槛高 | 可视化界面操作 | 降低使用门槛 |
| 脚本管理不便 | 图形化宏管理 | 方便组织和管理 |
| 实时反馈不足 | 录制/播放状态可视化 | 提升用户体验 |

### 1.3 目标用户

| 用户类型 | 特征 | 需求优先级 |
|---------|------|-----------|
| **办公人员** | 重复性数据录入、表单填写 | 录制、回放、简单编辑 |
| **测试人员** | 自动化测试、回归测试 | 高级编辑、循环播放 |
| **开发者** | 快速原型、工具集成 | API 导出、脚本编辑 |
| **普通用户** | 简单任务自动化 | 易用性、快速上手 |

### 1.4 产品目标

**短期目标 (1-2 个月)**
- 完成核心 UI 界面开发
- 实现录制、回放、管理基本功能
- 达到 MVP 可用状态

**中期目标 (3-6 个月)**
- 完善高级功能（脚本编辑、调试）
- 优化性能和用户体验
- 发布正式版本 v1.0

**长期目标 (6-12 个月)**
- 插件系统支持
- 云同步功能
- 社区分享平台

---

## 2. Wails v3 技术分析

### 2.1 Wails v3 核心特性

基于官方文档分析，Wails v3 相比 v2 带来了以下重大改进：

#### 2.1.1 架构变革

**从声明式到过程式**
```go
// v2 声明式 API
app.Options{
  Title: "My App",
  Width: 1024,
  Height: 768,
}

// v3 过程式 API (更灵活)
app := application.New(application.Options{
  Name: "My App",
  Services: []application.Service{
    application.NewService(&MyService{}),
  },
})
window := app.NewWebviewWindow(webview.WindowOptions{...})
```

#### 2.1.2 多窗口支持

| 功能 | v2 | v3 | 应用场景 |
|-----|----|----|---------|
| 窗口数量 | 单窗口 | 多窗口 | 主界面 + 设置 + 预览 |
| 窗口管理 | 有限配置 | 独立配置 | 灵活布局 |
| 窗口通信 | 复杂 | 事件系统 | 数据同步 |

**应用价值**
- **主窗口**: 录制控制、宏列表
- **设置窗口**: 配置管理
- **预览窗口**: 实时显示捕获事件
- **编辑器窗口**: 脚本编辑

#### 2.1.3 系统托盘集成

```go
// 系统托盘功能
tray := app.NewSystemTray(
  application.SystemTrayOptions{
    Icon: []byte(trayIcon),
  }
)

// 快速菜单
tray.AddMenuItem(&application.MenuItem{
  Label: "开始录制",
  OnClick: func(*application.MenuItem) {
    // 快速启动录制
  },
})
```

**用户价值**
- 后台运行，不占用任务栏
- 快速访问核心功能
- 状态指示（录制中/空闲）

#### 2.1.4 改进的绑定系统

**优势**
- 使用静态分析器生成绑定
- 保留注释和参数名
- 一条命令生成: `wails3 generate bindings`
- TypeScript 类型自动生成

**前后端通信**
```typescript
// 前端调用 (自动生成类型定义)
import { RecordService } from './bindings/record'

async function startRecording() {
  try {
    await RecordService.Start()
    updateUIState('recording')
  } catch (error) {
    showError(error.message)
  }
}
```

#### 2.1.5 事件系统

```go
// Go 后端发送事件
app.Emit(events.WailsEvent{
  Name: "recording:progress",
  Data: RecordingProgress{
    EventCount: 150,
    Duration:   4500,
  },
})

// 前端监听
import { On } from '@wailsio/runtime'

On('recording:progress', (data) => {
  updateProgressBar(data.eventCount)
})
```

**应用场景**
- 录制进度实时更新
- 播放状态通知
- 错误和警告提示

### 2.2 技术栈选择

#### 2.2.1 前端技术栈

| 技术 | 版本 | 用途 | 理由 |
|-----|------|-----|-----|
| **React** | 18.x | UI 框架 | 生态成熟，组件化开发 |
| **TypeScript** | 5.x | 类型安全 | 与 Wails 生成的绑定完美配合 |
| **Tailwind CSS** | 3.x | 样式框架 | 快速开发，一致的设计系统 |
| **Zustand** | 4.x | 状态管理 | 轻量级，适合中小型应用 |
| **React Query** | 5.x | 数据获取 | 自动缓存、重试、实时更新 |

#### 2.2.2 后端技术栈

| 技术 | 版本 | 用途 | 理由 |
|-----|------|-----|-----|
| **Go** | 1.24+ | 后端语言 | Wails v3 最低要求 |
| **Wails v3** | alpha | GUI 框架 | 项目要求 |
| **现有代码** | - | 核心功能 | 复用已实现的引擎 |

#### 2.2.3 UI 组件库

| 组件库 | 优势 | 适用场景 |
|-------|-----|---------|
| **Headless UI** | 无样式、可定制 | 基础组件 |
| **shadcn/ui** | Tailwind 集成好 | 完整 UI 系统 |
| **Lucide React** | 图标库 | 统一图标风格 |

**推荐方案**: shadcn/ui + Headless UI + Tailwind CSS
- 现代化设计
- TypeScript 友好
- 可定制性强
- 社区活跃

### 2.3 技术架构图

```
┌─────────────────────────────────────────────────────────┐
│                      Wails v3 Application                │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────────┐          ┌──────────────┐            │
│  │   Window 1   │          │   Window 2   │            │
│  │  (Main UI)   │          │  (Settings)  │            │
│  │              │          │              │            │
│  │  ┌────────┐  │          │  ┌────────┐  │            │
│  │  │ React  │  │          │  │ React  │  │            │
│  │  │   +    │  │          │  │   +    │  │            │
│  │  │  Tailwind│ │          │  │  Tailwind│ │            │
│  │  └────────┘  │          │  └────────┘  │            │
│  └──────┬───────┘          └──────┬───────┘            │
│         │                         │                     │
└─────────┼─────────────────────────┼─────────────────────┘
          │                         │
          │  Wails Bindings (TypeScript)
          │
┌─────────┼─────────────────────────┼─────────────────────┐
│         ▼                         ▼                     │
│  ┌─────────────────────────────────────────────────┐    │
│  │              Go Backend Services                │    │
│  ├─────────────────────────────────────────────────┤    │
│  │  ┌─────────────┐  ┌─────────────┐              │    │
│  │  │RecordService│  │PlayService  │              │    │
│  │  │             │  │             │              │    │
│  │  └──────┬──────┘  └──────┬──────┘              │    │
│  │         │                │                      │    │
│  │  ┌──────▼──────┐  ┌──────▼──────┐              │    │
│  │  │MacroService │  │ConfigService│              │    │
│  │  └─────────────┘  └─────────────┘              │    │
│  └──────────────────────────┬──────────────────────┘    │
│                             │                           │
│  ┌──────────────────────────▼──────────────────────┐    │
│  │          Core Engine (复用现有代码)              │    │
│  ├─────────────────────────────────────────────────┤    │
│  │  ┌─────────┐  ┌──────────┐  ┌──────────┐      │    │
│  │  │ Events  │  │Generator │  │Executor  │      │    │
│  │  │ Capture │  │          │  │          │      │    │
│  │  └─────────┘  └──────────┘  └──────────┘      │    │
│  └─────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

---

## 3. 用户需求分析

### 3.1 用户故事

#### US1: 快速录制操作
**作为** 办公人员
**我想要** 一键开始录制我的操作
**以便** 快速创建重复任务的自动化脚本

**验收标准**
- [ ] 点击"开始录制"按钮后立即开始捕获
- [ ] 录制过程中显示实时事件计数
- [ ] 停止后显示录制摘要（事件数、时长）
- [ ] 可以直接保存为宏

#### US2: 管理已保存的宏
**作为** 用户
**我想要** 查看和管理所有保存的宏
**以便** 快速找到并执行我需要的任务

**验收标准**
- [ ] 显示宏列表（名称、描述、创建时间）
- [ ] 支持搜索和过滤
- [ ] 可以播放、编辑、删除宏
- [ ] 显示宏的详细信息

#### US3: 编辑脚本
**作为** 高级用户
**我想要** 编辑生成的脚本
**以便** 添加自定义逻辑或优化操作

**验收标准**
- [ ] 提供代码编辑器（语法高亮）
- [ ] 支持保存修改
- [ ] 语法错误检查
- [ ] 可以测试运行修改后的脚本

#### US4: 配置录制选项
**作为** 用户
**我想要** 配置录制参数
**以便** 满足不同场景的需求

**验收标准**
- [ ] 可以设置录制热键
- [ ] 可以过滤鼠标移动事件
- [ ] 可以设置自动保存
- [ ] 可以配置播放速度

#### US5: 系统托盘快速访问
**作为** 用户
**我想要** 通过系统托盘快速访问功能
**以便** 在不影响其他工作的情况下使用

**验收标准**
- [ ] 托盘图标显示当前状态
- [ ] 右键菜单提供快速操作
- [ ] 点击图标显示/隐藏主窗口
- [ ] 录制时图标有视觉反馈

### 3.2 功能优先级 (RICE 框架)

| 功能 | Reach (触达) | Impact (影响) | Confidence (信心) | Effort (工作量) | RICE 分数 | 优先级 |
|-----|-------------|-------------|-----------------|---------------|----------|-------|
| 基础录制/停止 | 100% | 高 | 100% | 中 | **高** | P0 |
| 宏列表/播放 | 100% | 高 | 100% | 中 | **高** | P0 |
| 实时状态显示 | 100% | 中 | 100% | 低 | **高** | P0 |
| 脚本编辑器 | 60% | 高 | 90% | 高 | 中 | P1 |
| 系统托盘 | 80% | 中 | 100% | 中 | 中 | P1 |
| 热键支持 | 70% | 中 | 100% | 中 | 中 | P1 |
| 脚本调试 | 30% | 高 | 80% | 高 | 低 | P2 |
| 导入/导出 | 40% | 中 | 90% | 中 | 低 | P2 |
| 云同步 | 20% | 中 | 70% | 高 | 低 | P3 |

### 3.3 用户流程图

#### 流程 1: 快速录制和保存

```
┌─────────┐    ┌────────────┐    ┌──────────┐    ┌──────────┐
│ 启动应用 │ -> │ 点击录制按钮 │ -> │ 执行操作  │ -> │ 停止录制  │
└─────────┘    └────────────┘    └──────────┘    └──────────┘
                                                    │
                                                    ▼
                                           ┌──────────────┐
                                           │ 查看录制摘要  │
                                           └──────────────┘
                                                    │
                                                    ▼
                                           ┌──────────────┐
                                           │ 输入宏名称    │
                                           └──────────────┘
                                                    │
                                                    ▼
                                           ┌──────────────┐
                                           │ 保存宏       │
                                           └──────────────┘
```

#### 流程 2: 播放已保存的宏

```
┌─────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ 启动应用 │ -> │ 宏列表   │ -> │ 选择宏   │ -> │ 点击播放  │
└─────────┘    └──────────┘    └──────────┘    └──────────┘
                                               │
                                               ▼
                                      ┌──────────────┐
                                      │ 显示播放进度  │
                                      └──────────────┘
                                               │
                                               ▼
                                      ┌──────────────┐
                                      │ 播放完成提示  │
                                      └──────────────┘
```

### 3.4 非功能性需求

#### 3.4.1 性能要求

| 指标 | 目标 | 测量方法 |
|-----|-----|---------|
| 应用启动时间 | < 2 秒 | 冷启动计时 |
| 录制响应延迟 | < 50ms | 事件捕获时间戳 |
| UI 更新频率 | 60 FPS | 帧率监控 |
| 内存占用 | < 100 MB | 空闲时测量 |

#### 3.4.2 可用性要求

- **易学性**: 新用户 5 分钟内完成首次录制
- **效率**: 常用操作 < 3 次点击
- **容错**: 误操作可撤销或提示确认
- **反馈**: 所有操作有明确的视觉反馈

#### 3.4.3 兼容性要求

| 平台 | 最低版本 | 测试覆盖 |
|-----|---------|---------|
| Windows | Windows 10 21H2 | ✅ 必须 |
| macOS | macOS 12 Monterey | ✅ 必须 |
| Linux | Ubuntu 22.04 | ⚠️ 可选 |

---

## 4. UI 界面设计

### 4.1 整体布局

#### 4.1.1 主窗口结构

```
┌─────────────────────────────────────────────────────────┐
│  操作宏录制器                            ─ □ ×           │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌───────────────────────────────────────────────────┐  │
│  │  🎬 录制控制                                      │  │
│  │  ┌────────────┐    ┌────────────┐                │  │
│  │  │  开始录制   │    │  停止录制   │    事件数: 0  │  │
│  │  └────────────┘    └────────────┘    时长: 00:00 │  │
│  │                                                   │  │
│  │  ┌─────────────────────────────────────────┐     │  │
│  │  │ 实时事件日志                            │     │  │
│  │  │ [15:30:01] 鼠标移动: (100, 200)         │     │  │
│  │  │ [15:30:02] 鼠标左键点击                 │     │  │
│  │  │ [15:30:03] 键盘输入: 'a'                │     │  │
│  │  │ ...                                     │     │  │
│  │  └─────────────────────────────────────────┘     │  │
│  └───────────────────────────────────────────────────┘  │
│                                                         │
│  ┌───────────────────────────────────────────────────┐  │
│  │  📚 我的宏 (5)                          [+ 新建]  │  │
│  │                                                   │  │
│  │  ┌─────────────────────────────────────────────┐ │  │
│  │  │ 🔍 搜索宏...                                │ │  │
│  │  ├─────────────────────────────────────────────┤ │  │
│  │  │ 📄 数据录入自动化    2025-12-31  [▶] [✏] [🗑]│ │  │
│  │  │ 📄 测试脚本 v2        2025-12-30  [▶] [✏] [🗑]│ │  │
│  │  │ 📄 表单填充          2025-12-29  [▶] [✏] [🗑]│ │  │
│  │  │ ...                                       │ │  │
│  │  └─────────────────────────────────────────────┘ │  │
│  └───────────────────────────────────────────────────┘  │
│                                                         │
│  ┌───────────────────────────────────────────────────┐  │
│  │  ⚙️ 快速设置                                       │  │
│  │  播放速度: [1.0x ▼]  循环次数: [1 ▼]  [详细设置]  │  │
│  └───────────────────────────────────────────────────┘  │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

#### 4.1.2 窗口尺寸规范

| 窗口类型 | 最小尺寸 | 推荐尺寸 | 最大尺寸 |
|---------|---------|---------|---------|
| 主窗口 | 900×600 | 1200×800 | 无限制 |
| 设置窗口 | 600×400 | 700×500 | 无限制 |
| 编辑器窗口 | 800×600 | 1000×700 | 无限制 |

### 4.2 页面设计

#### 4.2.1 主页 (Home Page)

**区域划分**
1. **顶部栏**: 应用标题、窗口控制
2. **录制区**: 大按钮、状态指示、实时日志
3. **宏列表区**: 搜索、列表项、操作按钮
4. **底部栏**: 快速设置、状态信息

**颜色方案**
```css
/* 主色调 - 专业蓝 */
--primary-50: #eff6ff;
--primary-500: #3b82f6;
--primary-600: #2563eb;
--primary-700: #1d4ed8;

/* 状态色 */
--success: #10b981;  /* 录制中 */
--warning: #f59e0b;  /* 播放中 */
--danger: #ef4444;   /* 错误 */
--neutral: #6b7280;  /* 空闲 */

/* 背景和文本 */
--bg-primary: #ffffff;
--bg-secondary: #f9fafb;
--bg-tertiary: #f3f4f6;
--text-primary: #111827;
--text-secondary: #6b7280;
--text-tertiary: #9ca3af;
```

#### 4.2.2 录制状态设计

**空闲状态**
```
┌──────────────────────────────────┐
│   🎯 准备就绪                     │
│   点击"开始录制"捕获操作          │
│                                  │
│   [ 开始录制 ]                    │
└──────────────────────────────────┘
```

**录制中状态**
```
┌──────────────────────────────────┐
│   🔴 正在录制...                  │
│   ● 事件数: 127    ⏱ 00:00:15   │
│                                  │
│   [ 停止录制 ]                    │
│                                  │
│   ┌────────────────────────┐     │
│   │ 15:30:01  鼠标移动      │     │
│   │ 15:30:02  左键点击      │     │
│   │ 15:30:03  输入 'a'      │     │
│   │ ...                    │     │
│   └────────────────────────┘     │
└──────────────────────────────────┘
```

**录制完成状态**
```
┌──────────────────────────────────┐
│   ✅ 录制完成                     │
│   共捕获 127 个事件               │
│   用时 15.3 秒                    │
│                                  │
│   [ 保存为宏 ]  [ 放弃 ]         │
└──────────────────────────────────┘
```

#### 4.2.3 宏列表设计

**列表项结构**
```
┌────────────────────────────────────────────┐
│  📄 数据录入自动化              2025-12-31 │
│     自动化填充用户信息                     │
│     版本: 1.0.0  |  操作数: 15            │
│                                            │
│  [ ▶ 播放 ]  [ ✏ 编辑 ]  [ 🗑 删除 ]       │
└────────────────────────────────────────────┘
```

**空状态**
```
┌────────────────────────────────────────────┐
│                                            │
│         📭 还没有任何宏                     │
│                                            │
│      点击"开始录制"创建你的第一个宏         │
│                                            │
└────────────────────────────────────────────┘
```

#### 4.2.4 设置页面

**布局结构**
```
┌──────────────────────────────────────────┐
│  ⚙️ 设置                       [ 保存 ]  │
├──────────────────────────────────────────┤
│                                          │
│  📝 录制设置                             │
│  ├─ 过滤鼠标移动事件     [✓]             │
│  ├─ 录制热键             [F9 ▼]          │
│  ├─ 停止热键             [F10 ▼]         │
│  └─ 自动保存             [ ]             │
│                                          │
│  ▶️ 播放设置                             │
│  ├─ 默认播放速度         [1.0x ▼]        │
│  ├─ 播放前确认           [✓]             │
│  ├─ 出错时停止           [✓]             │
│  └─ 播放完成后提示       [✓]             │
│                                          │
│  💾 存储设置                             │
│  ├─ 宏文件目录           [.../macros]    │
│  ├─ 启用备份             [✓]             │
│  └─ 备份数量             [10]            │
│                                          │
│  🔔 通知设置                             │
│  ├─ 播放完成通知         [✓]             │
│  ├─ 错误通知             [✓]             │
│  └─ 系统托盘通知         [✓]             │
│                                          │
└──────────────────────────────────────────┘
```

### 4.3 交互设计

#### 4.3.1 按钮交互规范

| 按钮类型 | 样式 | 悬停 | 点击 | 禁用 |
|---------|-----|-----|-----|-----|
| 主要按钮 | 蓝色填充 | 深蓝色 | 按下效果 | 灰色 + 禁用光标 |
| 危险按钮 | 红色填充 | 深红色 | 按下效果 | 灰色 + 禁用光标 |
| 次要按钮 | 灰色边框 | 深灰色 | 按下效果 | 灰色 + 禁用光标 |
| 图标按钮 | 透明图标 | 背景色 | 按下效果 | 半透明 |

#### 4.3.2 动画效果

**场景 1: 录制按钮**
```css
/* 录制中: 脉冲动画 */
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.recording-indicator {
  animation: pulse 2s ease-in-out infinite;
}
```

**场景 2: 列表项悬停**
```css
/* 平滑过渡 */
.macro-item {
  transition: all 0.2s ease;
}

.macro-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}
```

**场景 3: 模态对话框**
```css
/* 淡入淡出 */
.modal-overlay {
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
```

#### 4.3.3 快捷键

| 功能 | 快捷键 | 说明 |
|-----|-------|-----|
| 开始/停止录制 | F9 | 全局热键 |
| 播放选中的宏 | Ctrl+Enter | 仅在主窗口 |
| 新建宏 | Ctrl+N | 仅在主窗口 |
| 删除宏 | Delete | 选中宏时 |
| 打开设置 | Ctrl+, | 全局 |
| 搜索宏 | Ctrl+F | 仅在主窗口 |

### 4.4 系统托盘设计

#### 4.4.1 图标设计

| 状态 | 图标样式 |
|-----|---------|
| 空闲 | 灰色圆点 ⚪ |
| 录制中 | 红色圆点 🔴 |
| 播放中 | 绿色圆点 🟢 |

#### 4.4.2 右键菜单

```
┌─────────────────┐
│ 显示主窗口       │
├─────────────────┤
│ 开始录制         │
│ 停止录制         │
├─────────────────┤
│ 最近播放         │
│  ├─ 数据录入...  │
│  └─ 测试脚本...  │
├─────────────────┤
│ 设置...          │
├─────────────────┤
│ 退出             │
└─────────────────┘
```

### 4.5 响应式设计

#### 4.5.1 窗口大小适配

| 宽度范围 | 布局调整 |
|---------|---------|
| < 900px | 单列布局，录制区在上，列表在下 |
| 900-1200px | 标准两列布局 |
| > 1200px | 宽屏优化，增加侧边栏 |

#### 4.5.2 高 DPI 支持

- 使用 Wails 的自动 DPI 缩放
- 图标资源提供多尺寸 (16x16, 32x32, 64x64, 128x128)
- 测试 125%, 150%, 200% 缩放

---

## 5. 前后端 API 设计

### 5.1 后端服务定义

基于现有的服务结构，我们需要将它们暴露为 Wails v3 Services。

#### 5.1.1 RecordService

```go
// internal/app/record_service.go
package app

import (
    "context"
    "time"

    "github.com/issueye/macro_operation/internal/service"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type RecordService struct {
    service *service.RecordService
    app     *application.Application
}

// RecordingStatus 录制状态
type RecordingStatus struct {
    IsStarted  bool      `json:"is_started"`
    EventCount int       `json:"event_count"`
    Duration   int64     `json:"duration_ms"`
    StartTime  time.Time `json:"start_time"`
}

// NewRecordService 创建录制服务
func NewRecordService(app *application.Application) *RecordService {
    return &RecordService{
        service: service.NewRecordService(),
        app:     app,
    }
}

// Start 开始录制
func (s *RecordService) Start(ctx context.Context) error {
    if err := s.service.Start(); err != nil {
        return err
    }

    // 发送事件
    s.app.Emit("recording:started", nil)

    // 启动状态更新
    go s.sendStatusUpdates()

    return nil
}

// Stop 停止录制
func (s *RecordService) Stop(ctx context.Context) (*RecordingStatus, error) {
    if err := s.service.Stop(); err != nil {
        return nil, err
    }

    status := s.getStatus()
    s.app.Emit("recording:stopped", status)

    return status, nil
}

// GetStatus 获取录制状态
func (s *RecordService) GetStatus(ctx context.Context) (*RecordingStatus, error) {
    return s.getStatus(), nil
}

// GenerateScript 生成脚本
func (s *RecordService) GenerateScript(ctx context.Context, name string) (string, error) {
    return s.service.GenerateScript(name)
}

// 内部方法
func (s *RecordService) getStatus() *RecordingStatus {
    return &RecordingStatus{
        IsStarted:  s.service.IsStarted(),
        EventCount: s.service.GetEventCount(),
    }
}

func (s *RecordService) sendStatusUpdates() {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    for s.service.IsStarted() {
        select {
        case <-ticker.C:
            status := s.getStatus()
            s.app.Emit("recording:progress", status)
        }
    }
}
```

#### 5.1.2 PlayService

```go
// internal/app/play_service.go
package app

import (
    "context"

    "github.com/issueye/macro_operation/internal/service"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type PlayService struct {
    service *service.PlayService
    app     *application.Application
}

// PlaybackStatus 播放状态
type PlaybackStatus struct {
    IsPlaying bool   `json:"is_playing"`
    Progress  int    `json:"progress"` // 0-100
    Error     string `json:"error,omitempty"`
}

// NewPlayService 创建播放服务
func NewPlayService(app *application.Application) *PlayService {
    return &PlayService{
        service: service.NewPlayService(),
        app:     app,
    }
}

// PlayScript 播放脚本
func (s *PlayService) PlayScript(ctx context.Context, script string) error {
    s.app.Emit("playback:started", nil)

    go func() {
        if err := s.service.PlayScript(script); err != nil {
            s.app.Emit("playback:error", map[string]string{
                "error": err.Error(),
            })
            return
        }
        s.app.Emit("playback:completed", nil)
    }()

    return nil
}

// IsPlaying 检查是否正在播放
func (s *PlayService) IsPlaying(ctx context.Context) (bool, error) {
    return s.service.IsPlaying(), nil
}

// Stop 停止播放 (需要扩展 PlayService)
func (s *PlayService) Stop(ctx context.Context) error {
    // TODO: 在 PlayService 中添加 Stop 方法
    s.app.Emit("playback:stopped", nil)
    return nil
}
```

#### 5.1.3 MacroService

```go
// internal/app/macro_service.go
package app

import (
    "context"

    "github.com/issueye/macro_operation/internal/model"
    "github.com/issueye/macro_operation/internal/repository"
    "github.com/issueye/macro_operation/internal/service"
)

type MacroService struct {
    service *service.MacroService
}

// MacroListItem 宏列表项
type MacroListItem struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Version     string `json:"version"`
    CreatedAt   string `json:"created_at"`
    EventCount  int    `json:"event_count"`
}

// MacroDetail 宏详情
type MacroDetail struct {
    MacroListItem
    ScriptCode string `json:"script_code"`
    Config     Config `json:"config"`
}

type Config struct {
    PlaybackSpeed float64 `json:"playback_speed"`
    LoopCount     int     `json:"loop_count"`
    AutoSave      bool    `json:"auto_save"`
}

// NewMacroService 创建宏服务
func NewMacroService(repo repository.MacroRepository) *MacroService {
    return &MacroService{
        service: service.NewMacroService(repo),
    }
}

// List 列出所有宏
func (s *MacroService) List(ctx context.Context) ([]*MacroListItem, error) {
    macros, err := s.service.List()
    if err != nil {
        return nil, err
    }

    items := make([]*MacroListItem, len(macros))
    for i, macro := range macros {
        items[i] = &MacroListItem{
            Name:        macro.Meta.Name,
            Description: macro.Meta.Description,
            Version:     macro.Meta.Version,
            CreatedAt:   macro.Meta.CreatedAt.Format("2006-01-02 15:04:05"),
            EventCount:  len(macro.Operations),
        }
    }

    return items, nil
}

// Load 加载宏详情
func (s *MacroService) Load(ctx context.Context, name string) (*MacroDetail, error) {
    macro, err := s.service.Load(name)
    if err != nil {
        return nil, err
    }

    return &MacroDetail{
        MacroListItem: MacroListItem{
            Name:        macro.Meta.Name,
            Description: macro.Meta.Description,
            Version:     macro.Meta.Version,
            CreatedAt:   macro.Meta.CreatedAt.Format("2006-01-02 15:04:05"),
        },
        ScriptCode: macro.Script.Code,
        Config: Config{
            PlaybackSpeed: macro.Config.PlaybackSpeed,
            LoopCount:     macro.Config.LoopCount,
            AutoSave:      macro.Config.AutoSave,
        },
    }, nil
}

// Save 保存宏
func (s *MacroService) Save(ctx context.Context, name, description, script string) error {
    macro := model.NewMacro(name)
    macro.Meta.Description = description
    macro.Script.Code = script

    return s.service.Save(macro)
}

// Delete 删除宏
func (s *MacroService) Delete(ctx context.Context, name string) error {
    return s.service.Delete(name)
}

// Exists 检查宏是否存在
func (s *MacroService) Exists(ctx context.Context, name string) (bool, error) {
    return s.service.Exists(name), nil
}

// Rename 重命名宏
func (s *MacroService) Rename(ctx context.Context, oldName, newName string) error {
    return s.service.Rename(oldName, newName)
}
```

#### 5.1.4 ConfigService

```go
// internal/app/config_service.go
package app

import (
    "context"

    "github.com/issueye/macro_operation/configs"
)

type ConfigService struct {
    config *configs.Config
}

// AppConfig 应用配置
type AppConfig struct {
    Name    string `json:"name"`
    Version string `json:"version"`
    Debug   bool   `json:"debug"`
}

// RecordConfig 录制配置
type RecordConfig struct {
    MaxDuration        int  `json:"max_duration"`
    EnableScreenshot   bool `json:"enable_screenshot"`
    FilterMouseMove    bool `json:"filter_mouse_move"`
    RecordHotkey       string `json:"record_hotkey"`
    StopHotkey         string `json:"stop_hotkey"`
    AutoSave           bool `json:"auto_save"`
}

// PlaybackConfig 播放配置
type PlaybackConfig struct {
    DefaultSpeed  float64 `json:"default_speed"`
    EnableSound   bool    `json:"enable_sound"`
    StopOnError   bool    `json:"stop_on_error"`
    ConfirmBefore bool    `json:"confirm_before"`
}

// StorageConfig 存储配置
type StorageConfig struct {
    MacrosDir       string `json:"macros_dir"`
    BackupEnabled   bool   `json:"backup_enabled"`
    BackupDir       string `json:"backup_dir"`
    MaxBackupCount  int    `json:"max_backup_count"`
}

// NewConfigService 创建配置服务
func NewConfigService(config *configs.Config) *ConfigService {
    return &ConfigService{config: config}
}

// GetAppConfig 获取应用配置
func (s *ConfigService) GetAppConfig(ctx context.Context) (*AppConfig, error) {
    return &AppConfig{
        Name:    s.config.App.Name,
        Version: s.config.App.Version,
        Debug:   s.config.App.Debug,
    }, nil
}

// GetRecordConfig 获取录制配置
func (s *ConfigService) GetRecordConfig(ctx context.Context) (*RecordConfig, error) {
    return &RecordConfig{
        MaxDuration:      s.config.Record.MaxDuration,
        EnableScreenshot: s.config.Record.EnableScreenshot,
        FilterMouseMove:  s.config.Record.FilterMouseMove,
        AutoSave:         s.config.Record.AutoSave,
    }, nil
}

// UpdateRecordConfig 更新录制配置
func (s *ConfigService) UpdateRecordConfig(ctx context.Context, config *RecordConfig) error {
    s.config.Record.MaxDuration = config.MaxDuration
    s.config.Record.EnableScreenshot = config.EnableScreenshot
    s.config.Record.FilterMouseMove = config.FilterMouseMove
    s.config.Record.AutoSave = config.AutoSave

    // 保存到文件
    return configs.Save(configs.GetConfigPath(), s.config)
}

// GetPlaybackConfig 获取播放配置
func (s *ConfigService) GetPlaybackConfig(ctx context.Context) (*PlaybackConfig, error) {
    return &PlaybackConfig{
        DefaultSpeed: s.config.Playback.DefaultSpeed,
        EnableSound:  s.config.Playback.EnableSound,
        StopOnError:  s.config.Playback.StopOnError,
    }, nil
}

// UpdatePlaybackConfig 更新播放配置
func (s *ConfigService) UpdatePlaybackConfig(ctx context.Context, config *PlaybackConfig) error {
    s.config.Playback.DefaultSpeed = config.DefaultSpeed
    s.config.Playback.EnableSound = config.EnableSound
    s.config.Playback.StopOnError = config.StopOnError

    return configs.Save(configs.GetConfigPath(), s.config)
}

// GetStorageConfig 获取存储配置
func (s *ConfigService) GetStorageConfig(ctx context.Context) (*StorageConfig, error) {
    return &StorageConfig{
        MacrosDir:      s.config.Storage.MacrosDir,
        BackupEnabled:  s.config.Storage.BackupEnabled,
        BackupDir:      s.config.Storage.BackupDir,
        MaxBackupCount: s.config.Storage.MaxBackupCount,
    }, nil
}
```

### 5.2 前端 API 调用

#### 5.2.1 自动生成的绑定

Wails v3 会自动生成 TypeScript 绑定：

```typescript
// bindings/recordService.ts
// 自动生成，无需手动编写

export interface RecordingStatus {
  is_started: boolean
  event_count: number
  duration_ms: number
  start_time: string
}

export class RecordService {
  Start(): Promise<void>
  Stop(): Promise<RecordingStatus>
  GetStatus(): Promise<RecordingStatus>
  GenerateScript(name: string): Promise<string>
}
```

#### 5.2.2 自定义 Hooks

```typescript
// src/hooks/useRecording.ts
import { useEffect, useState } from 'react'
import { RecordService, RecordingStatus } from '../bindings/record'
import { On } from '@wailsio/runtime'

export function useRecording() {
  const [status, setStatus] = useState<RecordingStatus | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // 开始录制
  const startRecording = async () => {
    setIsLoading(true)
    setError(null)
    try {
      await RecordService.Start()
    } catch (err) {
      setError(err.message)
    } finally {
      setIsLoading(false)
    }
  }

  // 停止录制
  const stopRecording = async () => {
    setIsLoading(true)
    setError(null)
    try {
      const result = await RecordService.Stop()
      setStatus(result)
      return result
    } catch (err) {
      setError(err.message)
      return null
    } finally {
      setIsLoading(false)
    }
  }

  // 刷新状态
  const refreshStatus = async () => {
    try {
      const result = await RecordService.GetStatus()
      setStatus(result)
    } catch (err) {
      setError(err.message)
    }
  }

  // 监听事件
  useEffect(() => {
    const unlisten = On('recording:progress', (data: RecordingStatus) => {
      setStatus(data)
    })

    return () => unlisten()
  }, [])

  return {
    status,
    isLoading,
    error,
    startRecording,
    stopRecording,
    refreshStatus,
  }
}
```

```typescript
// src/hooks/useMacros.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { MacroService, MacroListItem, MacroDetail } from '../bindings/macro'

export function useMacros() {
  const queryClient = useQueryClient()

  // 获取宏列表
  const macros = useQuery({
    queryKey: ['macros'],
    queryFn: () => MacroService.List(),
  })

  // 获取宏详情
  const useMacroDetail = (name: string) => {
    return useQuery({
      queryKey: ['macros', name],
      queryFn: () => MacroService.Load(name),
      enabled: !!name,
    })
  }

  // 保存宏
  const saveMacro = useMutation({
    mutationFn: ({ name, description, script }: {
      name: string
      description: string
      script: string
    }) => MacroService.Save(name, description, script),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['macros'] })
    },
  })

  // 删除宏
  const deleteMacro = useMutation({
    mutationFn: (name: string) => MacroService.Delete(name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['macros'] })
    },
  })

  // 重命名宏
  const renameMacro = useMutation({
    mutationFn: ({ oldName, newName }: { oldName: string; newName: string }) =>
      MacroService.Rename(oldName, newName),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['macros'] })
    },
  })

  return {
    macros,
    useMacroDetail,
    saveMacro,
    deleteMacro,
    renameMacro,
  }
}
```

### 5.3 事件系统

#### 5.3.1 后端事件

```go
// 事件常量
const (
    EventRecordingStarted   = "recording:started"
    EventRecordingStopped   = "recording:stopped"
    EventRecordingProgress  = "recording:progress"
    EventPlaybackStarted    = "playback:started"
    EventPlaybackStopped    = "playback:stopped"
    EventPlaybackCompleted  = "playback:completed"
    EventPlaybackError      = "playback:error"
    EventMacroSaved         = "macro:saved"
    EventMacroDeleted       = "macro:deleted"
)
```

#### 5.3.2 前端事件监听

```typescript
// src/events/index.ts
import { On, Emit } from '@wailsio/runtime'
import { RecordingStatus, PlaybackStatus } from '../bindings/types'

// 录制事件
export const onRecordingStarted = (callback: () => void) =>
  On('recording:started', callback)

export const onRecordingStopped = (callback: (status: RecordingStatus) => void) =>
  On('recording:stopped', callback)

export const onRecordingProgress = (callback: (status: RecordingStatus) => void) =>
  On('recording:progress', callback)

// 播放事件
export const onPlaybackStarted = (callback: () => void) =>
  On('playback:started', callback)

export const onPlaybackCompleted = (callback: () => void) =>
  On('playback:completed', callback)

export const onPlaybackError = (callback: (error: { error: string }) => void) =>
  On('playback:error', callback)

// 宏事件
export const onMacroSaved = (callback: (macro: { name: string }) => void) =>
  On('macro:saved', callback)

export const onMacroDeleted = (callback: (name: string) => void) =>
  On('macro:deleted', name)
```

---

## 6. 功能模块拆解

### 6.1 模块架构

```
┌─────────────────────────────────────────────────────────┐
│                    Frontend Modules                     │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐        │
│  │   录制     │  │  宏管理    │  │  播放控制  │        │
│  │  模块      │  │  模块      │  │  模块      │        │
│  └────────────┘  └────────────┘  └────────────┘        │
│         │                │                │             │
│         └────────────────┼────────────────┘             │
│                          │                              │
│                  ┌───────▼────────┐                     │
│                  │  共享组件/状态  │                     │
│                  └────────────────┘                     │
│                                                          │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                    Backend Modules                      │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐        │
│  │RecordService│  │MacroService│  │PlayService │        │
│  └────────────┘  └────────────┘  └────────────┘        │
│         │                │                │             │
│         └────────────────┼────────────────┘             │
│                          │                              │
│                  ┌───────▼────────┐                     │
│                  │  核心引擎层     │                     │
│                  │  (events,      │                     │
│                  │   generator,   │                     │
│                  │   executor)    │                     │
│                  └────────────────┘                     │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 6.2 前端模块详细拆解

#### 模块 1: 录制控制模块 (Recording)

**职责**
- 提供录制控制界面
- 显示录制状态和实时日志
- 处理录制相关用户操作

**组件列表**
```
src/features/recording/
├── components/
│   ├── RecordingControl.tsx        # 录制控制按钮组
│   ├── RecordingStatus.tsx         # 状态显示
│   ├── EventLog.tsx                # 事件日志列表
│   └── RecordingSummary.tsx        # 录制完成摘要
├── hooks/
│   └── useRecording.ts             # 录制状态 Hook
├── types/
│   └── recording.ts                # 类型定义
└── index.ts                        # 导出
```

**关键代码示例**

```typescript
// src/features/recording/components/RecordingControl.tsx
import { useRecording } from '../../hooks/useRecording'

export function RecordingControl() {
  const { status, isLoading, startRecording, stopRecording } = useRecording()

  return (
    <div className="recording-control">
      <button
        onClick={startRecording}
        disabled={isLoading || status?.is_started}
        className="btn btn-primary"
      >
        {status?.is_started ? '录制中...' : '开始录制'}
      </button>

      <button
        onClick={stopRecording}
        disabled={!status?.is_started}
        className="btn btn-danger"
      >
        停止录制
      </button>

      {status && (
        <div className="recording-stats">
          <span>事件: {status.event_count}</span>
          <span>时长: {formatDuration(status.duration_ms)}</span>
        </div>
      )}
    </div>
  )
}
```

**任务清单**
- [ ] 创建录制控制组件
- [ ] 实现状态显示
- [ ] 添加事件日志列表
- [ ] 实现录制摘要对话框
- [ ] 添加加载和错误状态

#### 模块 2: 宏管理模块 (Macro Management)

**职责**
- 显示宏列表
- 提供宏的增删改查操作
- 支持搜索和过滤

**组件列表**
```
src/features/macros/
├── components/
│   ├── MacroList.tsx               # 宏列表
│   ├── MacroListItem.tsx           # 列表项
│   ├── MacroSearch.tsx             # 搜索框
│   ├── MacroDetail.tsx             # 宏详情
│   └── ConfirmDeleteDialog.tsx     # 删除确认对话框
├── hooks/
│   └── useMacros.ts                # 宏数据 Hook
├── types/
│   └── macros.ts                   # 类型定义
└── index.ts
```

**任务清单**
- [ ] 创建宏列表组件
- [ ] 实现列表项渲染
- [ ] 添加搜索和过滤功能
- [ ] 实现删除确认对话框
- [ ] 添加空状态和加载状态
- [ ] 实现宏详情查看

#### 模块 3: 播放控制模块 (Playback)

**职责**
- 控制宏的播放
- 显示播放进度
- 处理播放错误

**组件列表**
```
src/features/playback/
├── components/
│   ├── PlaybackControl.tsx         # 播放控制
│   ├── PlaybackProgress.tsx        # 进度显示
│   └── PlaybackSettings.tsx        # 播放设置
├── hooks/
│   └── usePlayback.ts              # 播放状态 Hook
├── types/
│   └── playback.ts                 # 类型定义
└── index.ts
```

**任务清单**
- [ ] 创建播放控制组件
- [ ] 实现进度显示
- [ ] 添加播放设置（速度、循环）
- [ ] 处理播放错误
- [ ] 实现播放完成提示

#### 模块 4: 脚本编辑器模块 (Script Editor)

**职责**
- 提供代码编辑界面
- 语法高亮和提示
- 保存和验证

**组件列表**
```
src/features/editor/
├── components/
│   ├── ScriptEditor.tsx            # 代码编辑器
│   ├── EditorToolbar.tsx           # 工具栏
│   └── SaveConfirmDialog.tsx       # 保存确认
├── hooks/
│   └── useScriptEditor.ts          # 编辑器 Hook
└── index.ts
```

**技术选型**
- **CodeMirror 6**: 轻量级、可定制
- **Monaco Editor**: 功能丰富（VS Code 同款）

**推荐**: CodeMirror 6（更轻量，足够使用）

**任务清单**
- [ ] 集成 CodeMirror
- [ ] 添加 JavaScript 语法高亮
- [ ] 实现保存功能
- [ ] 添加基本的语法检查
- [ ] 实现撤销/重做

#### 模块 5: 设置模块 (Settings)

**职责**
- 提供设置界面
- 管理应用配置

**组件列表**
```
src/features/settings/
├── components/
│   ├── SettingsPage.tsx            # 设置页面
│   ├── RecordSettings.tsx          # 录制设置
│   ├── PlaybackSettings.tsx        # 播放设置
│   └── StorageSettings.tsx         # 存储设置
├── hooks/
│   └── useSettings.ts              # 设置 Hook
└── index.ts
```

**任务清单**
- [ ] 创建设置页面布局
- [ ] 实现录制设置表单
- [ ] 实现播放设置表单
- [ ] 实现存储设置
- [ ] 添加保存和重置功能

#### 模块 6: 系统托盘模块 (System Tray)

**职责**
- 管理系统托盘图标
- 处理托盘菜单

**组件列表**
```
src/features/tray/
├── hooks/
│   └── useTray.ts                  # 托盘 Hook
├── icons/
│   ├── icon-idle.ts                # 空闲图标
│   ├── icon-recording.ts           # 录制图标
│   └── icon-playing.ts             # 播放图标
└── index.ts
```

**任务清单**
- [ ] 创建托盘菜单结构
- [ ] 实现图标状态切换
- [ ] 添加托盘事件处理
- [ ] 实现快速操作

#### 模块 7: 共享组件 (Shared Components)

**职责**
- 提供可复用的 UI 组件

**组件列表**
```
src/components/shared/
├── Button.tsx                      # 按钮
├── Dialog.tsx                      # 对话框
├── Input.tsx                       # 输入框
├── Select.tsx                      # 下拉选择
├── Switch.tsx                      # 开关
├── Icon.tsx                        # 图标
└── Loading.tsx                     # 加载指示器
```

**任务清单**
- [ ] 创建基础按钮组件
- [ ] 实现对话框组件
- [ ] 创建表单组件
- [ ] 添加图标库
- [ ] 实现加载状态组件

### 6.3 后端模块详细拆解

#### 模块 1: 应用初始化 (App Bootstrap)

**文件**: `cmd/gui/main.go`

```go
package main

import (
    "log"

    "github.com/issueye/macro_operation/configs"
    "github.com/issueye/macro_operation/internal/app"
    "github.com/issueye/macro_operation/internal/repository"
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    // 加载配置
    config, err := configs.Load(configs.GetConfigPath())
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 创建应用
    app := application.New(application.Options{
        Name:        "Macro Recorder",
        Description: "操作宏录制和回放工具",
        MacDocker:   application.MacDockerOptions{},
    })

    // 初始化仓库
    macroRepo, err := repository.NewFileMacroRepository(config.Storage.MacrosDir)
    if err != nil {
        log.Fatalf("Failed to initialize repository: %v", err)
    }

    // 注册服务
    recordService := app.NewService(app.NewRecordService(app))
    playService := app.NewService(app.NewPlayService(app))
    macroService := app.NewService(app.NewMacroService(macroRepo))
    configService := app.NewService(app.NewConfigService(config))

    // 创建主窗口
    mainWindow := app.NewWebviewWindow(application.WebviewWindowOptions{
        Title:  "操作宏录制器",
        Width:  1200,
        Height: 800,
        URL:    "http://localhost:5173", // Vite dev server
    })

    // 创建系统托盘
    tray := createSystemTray(app, mainWindow)

    // 启动应用
    err = app.Run()

    if err != nil {
        log.Fatalf("Failed to run application: %v", err)
    }
}

func createSystemTray(app *application.Application, mainWindow *application.WebviewWindow) *application.SystemTray {
    // TODO: 实现系统托盘
    return nil
}
```

**任务清单**
- [ ] 创建 Wails v3 应用结构
- [ ] 初始化配置系统
- [ ] 注册所有服务
- [ ] 创建主窗口
- [ ] 实现系统托盘
- [ ] 配置前端资源路径

#### 模块 2: 服务层 (Services)

**任务清单**
- [ ] 将现有 Service 包装为 Wails Services
- [ ] 添加事件发送逻辑
- [ ] 实现状态查询方法
- [ ] 添加错误处理
- [ ] 编写单元测试

#### 模块 3: 前端构建 (Frontend Build)

**任务清单**
- [ ] 初始化 Vite + React + TypeScript 项目
- [ ] 配置 Tailwind CSS
- [ ] 安装 UI 组件库
- [ ] 配置 Wails 绑定生成
- [ ] 设置开发环境热更新
- [ ] 配置生产构建

---

## 7. 开发路线图

### 7.1 迭代计划

#### Sprint 1: 项目初始化 (Week 1)

**目标**: 搭建基础开发环境

| 任务 | 优先级 | 预估时间 | 负责人 |
|-----|-------|---------|-------|
| 初始化 Wails v3 项目 | P0 | 0.5天 | 后端 |
| 设置前端开发环境 | P0 | 0.5天 | 前端 |
| 配置 Tailwind CSS | P0 | 0.5天 | 前端 |
| 创建基础组件库 | P0 | 1天 | 前端 |
| 实现应用启动和窗口 | P0 | 0.5天 | 后端 |
| 配置开发热更新 | P1 | 0.5天 | 前端 |

**交付物**
- 可运行的开发环境
- 基础 UI 组件
- 主窗口显示

#### Sprint 2: 核心功能 MVP (Week 2-3)

**目标**: 实现基本的录制和播放功能

| 任务 | 优先级 | 预估时间 | 负责人 |
|-----|-------|---------|-------|
| 实现 RecordService | P0 | 1天 | 后端 |
| 实现 PlayService | P0 | 1天 | 后端 |
| 实现 MacroService | P0 | 1天 | 后端 |
| 创建录制控制界面 | P0 | 1.5天 | 前端 |
| 创建宏列表界面 | P0 | 1.5天 | 前端 |
| 实现播放控制 | P0 | 1天 | 前端 |
| 集成前后端 | P0 | 1天 | 全栈 |

**交付物**
- 可用的录制功能
- 可用的播放功能
- 宏列表显示

#### Sprint 3: 完善体验 (Week 4)

**目标**: 优化用户界面和交互

| 任务 | 优先级 | 预估时间 | 负责人 |
|-----|-------|---------|-------|
| 添加实时状态更新 | P0 | 1天 | 全栈 |
| 实现事件日志显示 | P0 | 0.5天 | 前端 |
| 添加加载和错误状态 | P0 | 0.5天 | 前端 |
| 优化动画和过渡 | P1 | 1天 | 前端 |
| 实现确认对话框 | P0 | 0.5天 | 前端 |
| 添加键盘快捷键 | P1 | 0.5天 | 全栈 |

**交付物**
- 流畅的用户体验
- 完整的状态反馈

#### Sprint 4: 高级功能 (Week 5-6)

**目标**: 实现脚本编辑和设置

| 任务 | 优先级 | 预估时间 | 负责人 |
|-----|-------|---------|-------|
| 实现脚本编辑器 | P1 | 2天 | 前端 |
| 实现设置页面 | P1 | 1.5天 | 前端 |
| 实现 ConfigService | P1 | 1天 | 后端 |
| 添加热键支持 | P1 | 1天 | 全栈 |
| 实现搜索和过滤 | P1 | 1天 | 前端 |

**交付物**
- 脚本编辑功能
- 完整的设置界面
- 热键支持

#### Sprint 5: 系统托盘和优化 (Week 7)

**目标**: 完成系统托盘和性能优化

| 任务 | 优先级 | 预估时间 | 负责人 |
|-----|-------|---------|-------|
| 实现系统托盘 | P1 | 1.5天 | 全栈 |
| 性能优化 | P1 | 1.5天 | 全栈 |
| 代码重构和优化 | P1 | 1天 | 全栈 |
| 错误处理完善 | P0 | 1天 | 全栈 |

**交付物**
- 系统托盘功能
- 优化的性能

#### Sprint 6: 测试和发布 (Week 8)

**目标**: 全面测试和发布准备

| 任务 | 优先级 | 预估时间 | 负责人 |
|-----|-------|---------|-------|
| 编写单元测试 | P0 | 2天 | 全栈 |
| 集成测试 | P0 | 1.5天 | 全栈 |
| Bug 修复 | P0 | 2天 | 全栈 |
| 文档编写 | P1 | 1天 | 全栈 |
| 打包和发布 | P0 | 1天 | 后端 |

**交付物**
- 可发布的版本
- 完整的文档

### 7.2 里程碑

| 里程碑 | 日期 | 标准 |
|-------|-----|-----|
| **M1: 环境就绪** | Week 1 结束 | 开发环境可运行 |
| **M2: MVP 可用** | Week 3 结束 | 核心功能可用 |
| **M3: 功能完整** | Week 6 结束 | 所有计划功能完成 |
| **M4: 发布就绪** | Week 8 结束 | 可正式发布 |

### 7.3 资源分配

**团队配置建议**
- 1 名后端工程师 (Go)
- 1 名前端工程师 (React/TypeScript)
- 1 名全栈工程师 (协调和集成)
- 1 名产品经理 (需求管理和测试)
- 1 名 UI/UX 设计师 (前期 2 周)

**时间线**
```
Week 1        Week 2-3      Week 4       Week 5-6      Week 7       Week 8
├──────┤      ├────────┤    ├──────┤     ├────────┤    ├──────┤    ├──────┤
初始化         MVP开发      体验优化      高级功能      托盘优化     测试发布
```

---

## 8. 技术风险评估

### 8.1 风险识别

| 风险 | 影响 | 概率 | 等级 | 缓解措施 |
|-----|-----|-----|-----|---------|
| **Wails v3 API 变更** | 高 | 中 | 🔴高 | 密切关注官方更新，及时适配 |
| **性能问题** | 中 | 低 | 🟡中 | 提前进行性能测试和优化 |
| **跨平台兼容性** | 高 | 中 | 🔴高 | 优先支持主流平台，预留测试时间 |
| **热键冲突** | 低 | 中 | 🟢低 | 允许用户自定义热键 |
| **前端构建复杂** | 中 | 低 | 🟡中 | 使用成熟的工具链，参考官方示例 |
| **状态同步** | 中 | 中 | 🟡中 | 使用事件系统，保持一致性 |
| **第三方依赖** | 中 | 低 | 🟡中 | 选择稳定、活跃的库 |

### 8.2 技术难点

#### 难点 1: 录制状态实时同步

**问题**: 后端录制状态需要实时反映到前端

**解决方案**
```go
// 后端: 定期发送状态更新
func (s *RecordService) startStatusUpdater() {
    ticker := time.NewTicker(100 * time.Millisecond)
    for range ticker.C {
        if !s.service.IsStarted() {
            return
        }
        status := s.getStatus()
        s.app.Emit("recording:progress", status)
    }
}
```

```typescript
// 前端: 监听状态更新
useEffect(() => {
  const unlisten = On('recording:progress', (status) => {
    setRecordingStatus(status)
  })
  return unlisten
}, [])
```

#### 难点 2: 全局热键

**问题**: 在后台也能响应热键

**解决方案**
- 使用 Wails v3 的全局快捷键 API
- 或使用 Go 库 `github.com/go-vgo/robotgo` 的 hook 功能

#### 难点 3: 事件日志性能

**问题**: 大量事件可能导致 UI 卡顿

**解决方案**
```typescript
// 使用虚拟滚动
import { useVirtualizer } from '@tanstack/react-virtual'

function EventLog({ events }: { events: Event[] }) {
  const parentRef = useRef<HTMLDivElement>(null)

  const virtualizer = useVirtualizer({
    count: events.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 30,
  })

  return (
    <div ref={parentRef} style={{ height: '400px', overflow: 'auto' }}>
      <div style={{ height: `${virtualizer.getTotalSize()}px` }}>
        {virtualizer.getVirtualItems().map((virtualRow) => (
          <div
            key={virtualRow.key}
            style={{
              position: 'absolute',
              top: 0,
              left: 0,
              width: '100%',
              height: `${virtualRow.size}px`,
              transform: `translateY(${virtualRow.start}px)`,
            }}
          >
            {events[virtualRow.index].message}
          </div>
        ))}
      </div>
    </div>
  )
}
```

### 8.3 回退方案

| 场景 | 回退方案 |
|-----|---------|
| **Wails v3 不稳定** | 降级到 Wails v2 |
| **性能不达标** | 简化 UI，减少动画 |
| **跨平台问题** | 优先支持 Windows |
| **热键冲突** | 仅在窗口激活时响应 |

---

## 9. 附录

### 9.1 参考资料

#### 官方文档
- [Wails v3 官方文档](https://v3alpha.wails.io/)
- [Wails v3 API 参考](https://v3alpha.wails.io/reference/overview/)
- [Wails v3 迁移指南](https://v3alpha.wails.io/migration/v2-to-v3/)
- [Wails v3 示例项目](https://github.com/wailsapp/wails/tree/master/examples)

#### 技术文档
- [React 文档](https://react.dev/)
- [TypeScript 文档](https://www.typescriptlang.org/docs/)
- [Tailwind CSS 文档](https://tailwindcss.com/docs)
- [Zustand 文档](https://zustand-demo.pmnd.rs/)
- [React Query 文档](https://tanstack.com/query/latest)

#### UI 组件库
- [shadcn/ui](https://ui.shadcn.com/)
- [Headless UI](https://headlessui.com/)
- [Lucide Icons](https://lucide.dev/)

#### Go 库
- [gohook](https://github.com/robotn/gohook) - 键盘监听
- [robotgo](https://github.com/go-vgo/robotgo) - 鼠标控制

### 9.2 术语表

| 术语 | 说明 |
|-----|-----|
| **宏 (Macro)** | 一系列可重复执行的操作序列 |
| **录制 (Recording)** | 捕获用户键盘和鼠标操作的过程 |
| **回放 (Playback)** | 重新执行录制的操作序列 |
| **事件 (Event)** | 键盘或鼠标的单次操作 |
| **脚本 (Script)** | 录制操作生成的可执行代码 |
| **热键 (Hotkey)** | 快速触发功能的键盘快捷键 |
| **系统托盘 (System Tray)** | 操作系统任务栏的通知区域 |
| **绑定 (Bindings)** | 前后端通信的自动生成代码 |

### 9.3 设计规范

#### 颜色规范

```css
/* 品牌色 */
--brand-primary: #3b82f6;
--brand-secondary: #8b5cf6;

/* 功能色 */
--color-success: #10b981;
--color-warning: #f59e0b;
--color-danger: #ef4444;
--color-info: #3b82f6;

/* 中性色 */
--gray-50: #f9fafb;
--gray-100: #f3f4f6;
--gray-200: #e5e7eb;
--gray-300: #d1d5db;
--gray-400: #9ca3af;
--gray-500: #6b7280;
--gray-600: #4b5563;
--gray-700: #374151;
--gray-800: #1f2937;
--gray-900: #111827;
```

#### 字体规范

```css
/* 字体家族 */
--font-sans: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
--font-mono: 'JetBrains Mono', 'Fira Code', monospace;

/* 字体大小 */
--text-xs: 0.75rem;    /* 12px */
--text-sm: 0.875rem;   /* 14px */
--text-base: 1rem;     /* 16px */
--text-lg: 1.125rem;   /* 18px */
--text-xl: 1.25rem;    /* 20px */
--text-2xl: 1.5rem;    /* 24px */
--text-3xl: 1.875rem;  /* 30px */
```

#### 间距规范

```css
/* 间距单位 (基于 4px 网格) */
--spacing-1: 0.25rem;  /* 4px */
--spacing-2: 0.5rem;   /* 8px */
--spacing-3: 0.75rem;  /* 12px */
--spacing-4: 1rem;     /* 16px */
--spacing-5: 1.25rem;  /* 20px */
--spacing-6: 1.5rem;   /* 24px */
--spacing-8: 2rem;     /* 32px */
--spacing-10: 2.5rem;  /* 40px */
--spacing-12: 3rem;    /* 48px */
```

#### 圆角规范

```css
--radius-sm: 0.25rem;   /* 4px */
--radius-md: 0.375rem;  /* 6px */
--radius-lg: 0.5rem;    /* 8px */
--radius-xl: 0.75rem;   /* 12px */
--radius-2xl: 1rem;     /* 16px */
--radius-full: 9999px;
```

### 9.4 用户反馈收集计划

#### 反馈渠道
1. **应用内反馈按钮**
2. **GitHub Issues**
3. **用户调研**
4. **数据分析**

#### 关键指标
- DAU (日活跃用户)
- 录制成功率
- 宏使用频率
- 平均会话时长
- 错误率

### 9.5 后续版本规划

#### v1.1 (发布后 1 个月)
- 性能优化
- Bug 修复
- 小功能改进

#### v1.2 (发布后 2-3 个月)
- 脚本调试功能
- 导入/导出功能
- 更多自定义选项

#### v2.0 (发布后 6 个月)
- 插件系统
- 云同步
- 社区分享平台
- AI 辅助录制

---

## 总结

本文档提供了 **Wails v3 宏录制界面** 的完整产品需求分析和方案设计，涵盖：

1. **技术分析**: Wails v3 的核心特性和最佳实践
2. **用户需求**: 用户故事、功能优先级、用户流程
3. **UI 设计**: 界面布局、交互设计、视觉规范
4. **API 设计**: 前后端接口定义和事件系统
5. **模块拆解**: 详细的功能模块和开发任务
6. **开发路线**: 8 周迭代计划和里程碑
7. **风险评估**: 技术风险识别和缓解措施

**核心建议**

1. **优先级**: 先实现核心录制和播放功能，再添加高级功能
2. **技术选型**: 使用成熟稳定的库，避免实验性技术
3. **用户体验**: 注重实时反馈和流畅交互
4. **性能**: 提前考虑性能优化，避免后期重构
5. **可扩展性**: 预留接口和扩展点，支持未来功能

**下一步行动**

1. 评审本文档，确认需求和方案
2. 组建开发团队
3. 启动 Sprint 1: 项目初始化
4. 建立定期评审机制

---

**文档维护**
- 创建者: Claude (Product Manager Agent)
- 最后更新: 2025-12-31
- 版本: v1.0.0

---

**Sources:**
- [What's New in Wails v3](https://v3alpha.wails.io/whats-new/)
- [API Reference - Wails v3](https://v3alpha.wails.io/reference/overview/)
- [Wails v3 Migration Guide](https://v3alpha.wails.io/migration/v2-to-v3/)
