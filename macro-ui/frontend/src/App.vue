<template>
  <div class="app">
    <!-- 顶部状态栏 -->
    <header class="header">
      <div class="logo">
        <span class="logo-icon">⚡</span>
        <span class="logo-text">Macro Recorder</span>
      </div>
      <div class="header-status">
        <div class="status-badge" :class="{ active: connected }">
          <span class="status-dot"></span>
          {{ connected ? '引擎已连接' : '未连接' }}
        </div>
      </div>
    </header>

    <div class="main-layout">
      <!-- 左侧边栏 -->
      <aside class="sidebar">
        <!-- 录制控制 -->
        <div class="card record-card">
          <div class="card-header">
            <h3>录制控制</h3>
          </div>
          <div class="record-controls">
            <button
              class="btn btn-record"
              :class="{ recording: isRecording }"
              @click="toggleRecording"
              :disabled="isPlaying"
            >
              <span class="record-icon"></span>
              {{ isRecording ? '停止录制' : '开始录制' }}
            </button>
            <div class="record-info" v-if="isRecording">
              <span class="pulse"></span>
              正在录制... {{ eventCount }} 个事件
            </div>
          </div>
        </div>

        <!-- 宏列表 -->
        <div class="card macro-card">
          <div class="card-header">
            <h3>宏列表</h3>
            <button class="btn-icon" @click="loadMacros" title="刷新">
              <svg viewBox="0 0 24 24" width="16" height="16"><path fill="currentColor" d="M17.65 6.35A7.958 7.958 0 0012 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0112 18c-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"/></svg>
            </button>
          </div>
          <div class="macro-list">
            <div
              v-for="macro in macros"
              :key="macro.name"
              class="macro-item"
              :class="{ active: selectedMacro === macro.name }"
              @click="selectMacro(macro)"
            >
              <div class="macro-info">
                <span class="macro-name">{{ macro.name }}</span>
                <span class="macro-events">{{ macro.event_count }} 事件</span>
              </div>
              <div class="macro-actions">
                <button class="btn-icon" @click.stop="playMacro(macro.name)" title="播放">
                  <svg viewBox="0 0 24 24" width="14" height="14"><path fill="currentColor" d="M8 5v14l11-7z"/></svg>
                </button>
                <button class="btn-icon danger" @click.stop="deleteMacro(macro.name)" title="删除">
                  <svg viewBox="0 0 24 24" width="14" height="14"><path fill="currentColor" d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"/></svg>
                </button>
              </div>
            </div>
            <div v-if="macros.length === 0" class="empty-state">
              暂无保存的宏
            </div>
          </div>
        </div>
      </aside>

      <!-- 右侧主内容 -->
      <main class="main-content">
        <!-- 脚本编辑区 -->
        <div class="card script-card">
          <div class="card-header">
            <h3>脚本预览</h3>
            <div class="script-actions">
              <button class="btn btn-secondary" @click="playCurrentScript" :disabled="!script || isRecording || isPlaying">
                <svg viewBox="0 0 24 24" width="14" height="14"><path fill="currentColor" d="M8 5v14l11-7z"/></svg>
                播放
              </button>
              <button class="btn btn-secondary" @click="saveMacro" :disabled="!script">
                <svg viewBox="0 0 24 24" width="14" height="14"><path fill="currentColor" d="M17 3H5c-1.11 0-2 .9-2 2v14c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V7l-4-4zm-5 16c-1.66 0-3-1.34-3-3s1.34-3 3-3 3 1.34 3 3-1.34 3-3 3zm3-10H5V5h10v4z"/></svg>
                保存
              </button>
            </div>
          </div>
          <div class="script-editor">
            <textarea
              class="script-area"
              v-model="script"
              readonly
              placeholder="录制事件后将自动生成脚本..."
            ></textarea>
          </div>
          <!-- 保存输入框（显示时覆盖） -->
          <div class="save-overlay" v-if="showSaveInput">
            <input
              class="save-input"
              v-model="macroName"
              placeholder="输入宏名称..."
              @keyup.enter="confirmSave"
              @keyup.esc="cancelSave"
              ref="saveInput"
            />
            <div class="save-actions">
              <button class="btn btn-small" @click="cancelSave">取消</button>
              <button class="btn btn-primary btn-small" @click="confirmSave">保存</button>
            </div>
          </div>
        </div>

        <!-- 日志面板 -->
        <div class="card log-card">
          <div class="card-header">
            <h3>运行日志</h3>
            <button class="btn-icon" @click="clearLogs" title="清空日志">
              <svg viewBox="0 0 24 24" width="16" height="16"><path fill="currentColor" d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/></svg>
            </button>
          </div>
          <div class="log-panel" ref="logPanel">
            <div v-for="log in logs" :key="log.id" class="log-entry" :class="log.type">
              <span class="log-time">{{ log.time }}</span>
              <span class="log-message">{{ log.message }}</span>
            </div>
            <div v-if="logs.length === 0" class="empty-state small">
              暂无日志
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      isRecording: false,
      isPlaying: false,
      connected: false,
      eventCount: 0,
      script: '',
      macroName: '',
      selectedMacro: null,
      showSaveInput: false,
      macros: [],
      logs: [],
      statusTimer: null
    }
  },

  mounted() {
    this.init()
  },

  beforeUnmount() {
    this.stopStatusPolling()
  },

  methods: {
    async init() {
      this.addLog('应用启动', 'info')
      await this.$nextTick()
      await this.checkConnection()
      this.loadMacros()
      this.startStatusPolling()
    },

    addLog(message, type = 'success') {
      const now = new Date()
      const time = now.toLocaleTimeString()
      this.logs.unshift({
        id: Date.now() + Math.random(),
        time,
        message,
        type
      })
      if (this.logs.length > 50) {
        this.logs.pop()
      }
    },

    clearLogs() {
      this.logs = []
    },

    async checkConnection() {
      try {
        await window.go.main.App.IsRecording()
        this.connected = true
        this.addLog('引擎服务连接成功', 'success')
      } catch (e) {
        this.connected = false
        this.addLog('无法连接到引擎服务', 'error')
      }
    },

    startStatusPolling() {
      this.statusTimer = setInterval(async () => {
        if (this.isRecording) {
          try {
            const count = await window.go.main.App.GetEventCount()
            this.eventCount = count
            if (count > 0 && !this.script) {
              await this.refreshScript()
            }
          } catch (e) {
            console.error('获取状态失败:', e)
          }
        }
      }, 500)
    },

    stopStatusPolling() {
      if (this.statusTimer) {
        clearInterval(this.statusTimer)
      }
    },

    async toggleRecording() {
      if (this.isRecording) {
        await this.stopRecording()
      } else {
        await this.startRecording()
      }
    },

    async startRecording() {
      try {
        await window.go.main.App.StartRecording()
        this.isRecording = true
        this.eventCount = 0
        this.script = ''
        this.selectedMacro = null
        this.addLog('开始录制', 'info')
      } catch (e) {
        this.addLog('开始录制失败: ' + e.message, 'error')
      }
    },

    async stopRecording() {
      try {
        await window.go.main.App.StopRecording()
        this.isRecording = false
        this.addLog(`录制完成，共 ${this.eventCount} 个事件`, 'success')
        await this.refreshScript()
      } catch (e) {
        this.addLog('停止录制失败: ' + e.message, 'error')
      }
    },

    async refreshScript() {
      try {
        this.script = await window.go.main.App.GenerateCurrentScript()
        this.addLog('脚本已生成', 'success')
      } catch (e) {
        this.addLog('生成脚本失败: ' + e.message, 'error')
      }
    },

    async playCurrentScript() {
      if (!this.script) {
        this.addLog('没有可播放的脚本', 'warn')
        return
      }
      try {
        this.isPlaying = true
        this.addLog('开始播放脚本', 'info')
        await window.go.main.App.PlayScript(this.script)
        this.addLog('脚本播放完成', 'success')
      } catch (e) {
        this.addLog('播放失败: ' + e.message, 'error')
      } finally {
        this.isPlaying = false
      }
    },

    selectMacro(macro) {
      this.selectedMacro = macro.name
      this.script = macro.script || ''
    },

    async playMacro(name) {
      try {
        const script = await window.go.main.App.LoadMacro(name)
        if (!script) {
          this.addLog('宏 ' + name + ' 没有脚本', 'warn')
          return
        }
        this.script = script
        this.selectedMacro = name
        this.isPlaying = true
        this.addLog('播放宏: ' + name, 'info')
        await window.go.main.App.PlayScript(script)
        this.addLog('宏播放完成', 'success')
      } catch (e) {
        this.addLog('播放宏失败: ' + e.message, 'error')
      } finally {
        this.isPlaying = false
      }
    },

    saveMacro() {
      if (!this.script) {
        this.addLog('没有可保存的脚本', 'warn')
        return
      }
      this.showSaveInput = true
      this.macroName = this.selectedMacro || ''
      this.$nextTick(() => {
        this.$refs.saveInput?.focus()
      })
    },

    cancelSave() {
      this.showSaveInput = false
      this.macroName = ''
    },

    async confirmSave() {
      if (!this.macroName.trim()) {
        this.addLog('请输入宏名称', 'warn')
        return
      }
      try {
        await window.go.main.App.SaveMacro(this.macroName.trim(), this.script)
        this.addLog(`宏 "${this.macroName}" 已保存`, 'success')
        this.selectedMacro = this.macroName.trim()
        this.macroName = ''
        this.showSaveInput = false
        await this.loadMacros()
      } catch (e) {
        this.addLog('保存宏失败: ' + e.message, 'error')
      }
    },

    async deleteMacro(name) {
      if (!confirm(`确定要删除宏 "${name}" 吗？`)) return
      try {
        await window.go.main.App.DeleteMacro(name)
        this.addLog(`宏 "${name}" 已删除`, 'success')
        if (this.selectedMacro === name) {
          this.selectedMacro = null
        }
        await this.loadMacros()
      } catch (e) {
        this.addLog('删除宏失败: ' + e.message, 'error')
      }
    },

    async loadMacros() {
      try {
        this.macros = await window.go.main.App.ListMacros()
      } catch (e) {
        this.addLog('加载宏列表失败: ' + e.message, 'error')
      }
    }
  }
}
</script>
