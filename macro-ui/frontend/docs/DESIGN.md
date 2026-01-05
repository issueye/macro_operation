# 桌面宏录制工具 - UI/UX 设计文档

> 版本：v1.0
> 设计师：Claude UI Designer
> 日期：2026-01-05

---

## 一、设计概述

### 1.1 设计定位

**产品定位**：桌面自动化工具，用于录制和回放键盘/鼠标操作作为JavaScript宏。

**目标用户**：
| 用户类型 | 需求特点 | 设计侧重 |
|---------|---------|---------|
| 开发者 | 熟悉快捷键，关注代码编辑 | 完整的快捷键支持，VS Code风格 |
| 办公白领 | 重复性操作多，技术背景弱 | 一键录制，简化流程，清晰反馈 |
| 数据录入人员 | 批量操作，高频使用 | 快速回放，宏管理便捷 |

### 1.2 设计风格

- **风格**：现代极简 + 轻量级深色主题
- **关键词**：高效、专注、清晰、可控、专业
- **设计理念**：工具型应用定位，用户专注操作，视觉简洁不干扰

---

## 二、布局设计

### 2.1 整体布局结构

```
+-----------------------------------------------------------------------+
|                           Header (56px)                               |
+----------------+------------------------------------------------------+
|                |                                                      |
|  Left Panel    |                    Main Area                         |
|  (可收起/240px |                                                      |
|   ~64px)       |  +--------------------------------------------------+ |
|                |  |                                                  | |
|  +----------+  |  |          Script Editor                          | |
|  | ○ Record |  |  |          (flex: 1)                              | |
|  | (主按钮)  |  |  |                                                  | |
|  +----------+  |  +--------------------------------------------------+ |
|                |                                                    |
|  +----------+  |  +--------------------------------------------------+ |
|  | Macros   |  |  |  Event Monitor       Log Panel                  | |
|  |          |  |  |  (flex: 1)            (flex: 0.8)               | |
|  |          |  |  |  [可折叠]                                    | |
|  |          |  |  +--------------------------------------------------+ |
|  +----------+  |                                                    |
|                |                                                    |
+----------------+------------------------------------------------------+
```

### 2.2 布局特性

| 区域 | 宽度策略 | 交互行为 |
|-----|---------|---------|
| 左侧面板 | 240px (可折叠至64px) | hover展开 / 双击收起 / 设置记忆 |
| 脚本编辑区 | flex: 1，自适应剩余空间 | 无固定限制 |
| 底部面板 | 高度可拖拽调整 | 拖拽分隔线 |
| 整体 | 100% - 侧边栏 - 16px gaps | 响应式断点 |

---

## 三、视觉设计规范

### 3.1 色彩系统

#### 主色系（Indigo）

```css
:root {
  --color-primary: #6366f1;        /* Indigo 500 - 主要操作 */
  --color-primary-hover: #818cf8;  /* Indigo 400 */
  --color-primary-active: #4f46e5; /* Indigo 600 */
}
```

#### 功能色

```css
:root {
  --color-success: #10b981;   /* 成功 - 绿色 */
  --color-warning: #f59e0b;   /* 警告 - 琥珀色 */
  --color-error: #ef4444;     /* 错误 - 红色 */
  --color-info: #3b82f6;      /* 信息 - 蓝色 */
}
```

#### 录制状态色

```css
:root {
  --color-recording: #ef4444;
  --color-recording-bg: rgba(239, 68, 68, 0.1);
  --color-recording-glow: rgba(239, 68, 68, 0.4);

  --color-playing: #10b981;
  --color-playing-bg: rgba(16, 185, 129, 0.1);

  --color-connected: #10b981;
  --color-disconnected: #6b7280;
}
```

#### 背景色系（层次递进）

```css
:root {
  --bg-base: #0f172a;      /* Slate 900 - 最底层 */
  --bg-surface: #1e293b;   /* Slate 800 - 卡片/面板 */
  --bg-element: #334155;   /* Slate 700 - 可交互元素 */
  --bg-hover: #475569;     /* Slate 600 - hover状态 */

  --text-primary: #f1f5f9;   /* Slate 100 - 主要文字 */
  --text-secondary: #94a3b8; /* Slate 400 - 次要文字 */
  --text-muted: #64748b;     /* Slate 500 - 禁用/提示 */

  --border-subtle: rgba(148, 163, 184, 0.1);
  --border-default: rgba(148, 163, 184, 0.2);
}
```

### 3.2 字体系统

```css
:root {
  --font-sans: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
  --font-mono: 'JetBrains Mono', 'Fira Code', 'Consolas', monospace;

  --text-xs: 11px;    /* 标签、小字 */
  --text-sm: 12px;    /* 正文辅助 */
  --text-base: 13px;  /* 常规正文 */
  --text-lg: 14px;    /* 强调文字 */
  --text-xl: 16px;    /* 标题 */
  --text-2xl: 20px;   /* 大标题 */
}
```

### 3.3 间距系统（8px基准）

```css
:root {
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-5: 20px;
  --space-6: 24px;
  --space-8: 32px;
}
```

### 3.4 圆角规范

```css
:root {
  --radius-sm: 4px;      /* 小按钮、标签 */
  --radius-md: 6px;      /* 输入框、小卡片 */
  --radius-lg: 8px;      /* 主要按钮、中卡片 */
  --radius-xl: 12px;     /* 大卡片、模态框 */
  --radius-full: 9999px; /* 徽章、状态点 */
}
```

### 3.5 阴影与层次

```css
:root {
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.3);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.3);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.3);
  --shadow-glow-recording: 0 0 20px rgba(239, 68, 68, 0.5);
  --shadow-glow-playing: 0 0 20px rgba(16, 185, 129, 0.5);
}
```

---

## 四、组件设计

### 4.1 录制控制按钮

**设计理念**：最高频操作，视觉聚焦、一键触达、状态清晰

```
┌─────────────────────────────────────────────┐
│                                             │
│              ┌───────────┐                  │
│              │    ◉      │  ← 动态光晕      │
│              │   Start   │  ← 录制中变Stop  │
│              │           │                  │
│              └───────────┘                  │
│                                             │
│         [ 00:00:12 ]                        │
│          256 events                         │
│                                             │
└─────────────────────────────────────────────┘
```

**状态设计**：

| 状态 | 视觉表现 | 交互效果 |
|-----|---------|---------|
| 默认 | 主色填充，白色文字 | hover: 亮度提升 + 上移2px |
| 录制中 | 红色填充，脉动光晕 | 点击即停止 |
| 禁用（播放中） | 50%透明度，禁止图标 | 无hover效果 |
| 成功（录制完成） | 绿色边框 + Check图标 | 2秒后自动消失 |

**CSS要点**：

```css
.btn-record.recording {
  background: var(--color-recording);
  box-shadow: var(--shadow-md), var(--shadow-glow-recording);
  animation: pulse-recording 2s ease-in-out infinite;
}

@keyframes pulse-recording {
  0%, 100% { box-shadow: var(--shadow-md), var(--shadow-glow-recording); }
  50% { box-shadow: var(--shadow-md), 0 0 0 8px rgba(239, 68, 68, 0); }
}
```

### 4.2 宏列表组件

```
┌─────────────────────────────────────────┐
│  My Macros                      [+ New] │
├─────────────────────────────────────────┤
│  ┌─────────────────────────────────────┐│
│  │ 📁 daily-report              ▶ 🗑    ││
│  │     12 events · Updated 2h ago      ││
│  └─────────────────────────────────────┘│
│  ┌─────────────────────────────────────┐│
│  │ 📁 fill-form               ▶ 🗑 ★   ││
│  │     89 events · Updated yesterday    ││
│  └─────────────────────────────────────┘│
└─────────────────────────────────────────┘
```

**设计要点**：
- 卡片背景：`var(--bg-surface)`
- hover效果：左侧出现2px主色边框
- 选中状态：`var(--bg-element)` + 主色左侧边框

### 4.3 代码编辑器区域

**Monaco Editor 配置建议**：

```javascript
{
  base: 'vs-dark',
  inherit: true,
  rules: [
    { token: 'keyword', foreground: 'c792ea' },
    { token: 'string', foreground: 'c3e88d' },
    { token: 'number', foreground: 'f78c6c' },
    { token: 'function', foreground: '82aaff' },
    { token: 'comment', foreground: '676e95' },
  ],
  colors: {
    'editor.background': '#0f172a',
    'editor.foreground': '#f1f5f9',
    'editor.lineHighlightBackground': '#1e293b',
    'editor.selectionBackground': '#334155',
    'editorCursor.foreground': '#6366f1',
  }
}
```

### 4.4 事件监控面板

```
┌─────────────────────────────────────────┐
│  Event Monitor              [Pause] [Clear] │
├─────────────────────────────────────────┤
│  [All] [Keyboard] [Mouse] [Scroll]      │
├─────────────────────────────────────────┤
│  ⌨ keyDown  "a"                          │
│  ⌨ keyUp    "a"                          │
│  🖱 move    1200, 450                    │
│  🖱 click   left (1200, 450)             │
└─────────────────────────────────────────┘
```

**事件类型颜色编码**：

| 事件类型 | 图标 | 颜色 |
|---------|-----|------|
| 键盘按下 | ⌨ | --color-primary |
| 键盘释放 | ⌨ | --color-secondary |
| 鼠标移动 | 🖱 | --color-info |
| 鼠标点击 | 🖱 | --color-success |
| 鼠标滚轮 | ↕ | --color-warning |

### 4.5 状态指示系统

| 状态 | 样式 | 位置 |
|-----|------|------|
| 引擎连接 | 绿色圆点 + "已连接" | Header右侧 |
| 录制中 | 红色脉动圆点 + 时长 + 事件数 | Header右侧 |
| 播放中 | 绿色播放图标 + 进度 | Header右侧 |
| 错误 | 红色警告图标 + 提示文字 | Toast通知 |

---

## 五、交互设计

### 5.1 快捷键设计

| 快捷键 | 功能 | 设计理由 |
|-------|------|---------|
| `Ctrl/Cmd + R` | 开始/停止录制 | R = Record，最常用 |
| `Ctrl/Cmd + Enter` | 运行当前脚本 | 快速测试 |
| `Ctrl/Cmd + S` | 保存宏 | 符合用户习惯 |
| `Ctrl/Cmd + N` | 新建宏 | 符合用户习惯 |
| `Ctrl/Cmd + ,` | 打开设置 | 通用快捷键 |
| `F1` | 显示快捷键列表 | 辅助功能 |
| `Esc` | 取消当前操作 | 通用取消 |

### 5.2 关键交互反馈

| 操作 | 即时反馈 | 完成反馈 | 异常反馈 |
|-----|---------|---------|---------|
| 点击录制 | 按钮变红 + 脉动动画 | Toast: "录制完成，128个事件" | Toast: "录制失败" |
| 点击停止 | 红色按钮变回主色 | Toast: "脚本已生成" | Toast: "停止失败" |
| 点击播放 | 按钮显示播放图标 | Toast: "播放完成" | Toast: "播放失败" |
| 保存宏 | 按钮短暂绿色高亮 | Toast: "宏已保存" | Toast: "保存失败" |
| 删除宏 | 弹出确认对话框 | Toast: "宏已删除" | - |

---

## 六、主界面示意

### 6.1 默认状态

```
+-----------------------------------------------------------------------+
|  ⚡ Macro Recorder                                     ● Connected   |
+-----------------------------------------------------------------------+
|                               |                                        |
|  [☰] Collapse                 |  +---------------------------------+  |
|                               |  | script.js             [▶ Run]   |  |
|  +-------------------------+  |  +---------------------------------+  |
|  |  ● Start Recording      |  |  |                                 |  |
|  |    00:12:34 · 256 evts  |  |  |  1  │ function main() {         |  |
|  +-------------------------+  |  |  2  │     mouseMove(100, 200);   |  |
|                               |  |  3  │     mouseClick("left");    |  |
|  My Macros (3)                |  |  4  │     keyType("Hello");      |  |
|  +-------------------------+  |  |  5  │ }                          |  |
|  | 📁 daily-report    ▶ 🗑 |  |  |  6  │                             |  |
|  |     12 events · 2h ago |  |  |  7  │                             |  |
|  +-------------------------+  |  |  ...│                             |  |
|  | 📁 data-entry      ▶ 🗑 |  |  |                                 |  |
|  |     89 events · yesterday|  |  +---------------------------------+  |
|  +-------------------------+  |                                        |
|  | 📁 backup-data     ▶ 🗑 |  |  +---------------------------------+  |
|  |     234 events · 3d ago |  |  | [Keyboard] [Mouse]     [Clear]  |  |
|  +-------------------------+  |  +---------------------------------+  |
|                               |  |  ⌨ keyDown  "Ctrl"              |  |
|                               |  |  ⌨ keyUp    "Ctrl"              |  |
|                               |  |  🖱 move    850, 320             |  |
|                               |  +---------------------------------+  |
+-------------------------------+----------------------------------------+
```

### 6.2 录制状态

```
+-----------------------------------------------------------------------+
|  ⚡ Macro Recorder                                     ● Recording...  |
+-----------------------------------------------------------------------+
|                               |                                        |
|  [☰] Expand                   |  +---------------------------------+  |
|                               |  | script.js             [⏹ Stop]  |  |
|  +-------------------------+  |  +---------------------------------+  |
|  |  ⏹ Stop Recording       |  |  |  Recording...                   |  |
|  |    00:05:23 · 67 evts   |  |  +---------------------------------+  |
|  +-------------------------+  |                                        |
|                               |  +---------------------------------+  |
|  My Macros                    |  | [Keyboard] [Mouse]     [Pause]  |  |
|  +-------------------------+  |  +---------------------------------+  |
|  | 📁 daily-report    ▶ 🗑 |  |  ⌨ keyDown  "a"                   |  |
|  +-------------------------+  |  ⌨ keyUp    "a"                   |  |
|                               |  🖱 move    1200, 450             |  |
|                               |  🖱 click   left (1200, 450)      |  |
|                               +-------------------------------------+
+-------------------------------+
```

---

## 七、实施建议

### 阶段一：基础改造（1-2天）
1. 替换CSS变量系统
2. 重构整体布局为弹性布局
3. 更新Header和状态系统

### 阶段二：组件优化（2-3天）
1. 重绘录制按钮（状态动画）
2. 优化宏列表样式
3. 更新事件监控面板

### 阶段三：交互增强（1-2天）
1. 添加快捷键支持
2. 完善空状态和加载态
3. 添加微交互动效

### 验收清单
- [ ] 布局支持任意窗口尺寸
- [ ] 录制按钮状态清晰可辨
- [ ] Monaco Editor与整体风格统一
- [ ] 所有操作有明确反馈
- [ ] 支持常用快捷键

---

## 八、文件结构

```
macro-ui/frontend/
├── src/
│   ├── assets/
│   │   └── variables.css      # CSS变量定义
│   ├── components/
│   │   ├── Header.vue         # 顶部栏（状态指示）
│   │   ├── RecordControl.vue  # 录制控制按钮
│   │   ├── MacroList.vue      # 宏列表
│   │   ├── ScriptEditor.vue   # 脚本编辑区
│   │   ├── MonacoEditor.vue   # Monaco编辑器封装
│   │   ├── EventMonitor.vue   # 事件监控
│   │   └── LogPanel.vue       # 日志面板
│   ├── App.vue                # 主布局
│   └── style.css              # 全局样式
```

---

## 附录：CSS变量速查表

```css
:root {
  /* Colors - Primary */
  --color-primary: #6366f1;
  --color-primary-hover: #818cf8;
  --color-primary-active: #4f46e5;

  /* Colors - Recording */
  --color-recording: #ef4444;
  --color-recording-bg: rgba(239, 68, 68, 0.1);
  --color-recording-glow: rgba(239, 68, 68, 0.5);

  /* Colors - Playing */
  --color-playing: #10b981;
  --color-playing-bg: rgba(16, 185, 129, 0.1);

  /* Colors - Status */
  --color-success: #10b981;
  --color-warning: #f59e0b;
  --color-error: #ef4444;
  --color-info: #3b82f6;

  /* Background - Dark Theme */
  --bg-base: #0f172a;
  --bg-surface: #1e293b;
  --bg-element: #334155;
  --bg-hover: #475569;

  /* Text */
  --text-primary: #f1f5f9;
  --text-secondary: #94a3b8;
  --text-muted: #64748b;

  /* Border */
  --border-subtle: rgba(148, 163, 184, 0.1);
  --border-default: rgba(148, 163, 184, 0.2);

  /* Typography */
  --font-sans: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
  --font-mono: 'JetBrains Mono', 'Fira Code', monospace;

  /* Spacing */
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-5: 20px;
  --space-6: 24px;

  /* Border Radius */
  --radius-sm: 4px;
  --radius-md: 6px;
  --radius-lg: 8px;
  --radius-xl: 12px;
  --radius-full: 9999px;

  /* Shadows */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.3);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.3);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.3);
  --shadow-glow-recording: 0 0 20px rgba(239, 68, 68, 0.5);

  /* Transitions */
  --transition-fast: 150ms ease;
  --transition-normal: 200ms ease;
  --transition-slow: 300ms ease;
}
```
