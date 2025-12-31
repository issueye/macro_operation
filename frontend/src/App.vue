<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

// State
const recording = ref(false)
const eventCount = ref(0)
const macroName = ref('')
const macros = ref([])
const selectedMacro = ref('')
const currentScript = ref('')

// Initialize
onMounted(() => {
  loadMacros()
})

// Recording Controls
async function startRecording() {
  try {
    await window.go.main.App.StartRecording()
    recording.value = true

    // Update event count and script every 500ms
    const interval = setInterval(async () => {
      if (recording.value) {
        eventCount.value = await window.go.main.App.GetEventCount()
        // 获取当前脚本并更新到编辑框
        try {
          const script = await window.go.main.App.GenerateCurrentScript()
          currentScript.value = script
        } catch (err) {
          console.error('获取脚本失败:', err)
        }
      } else {
        clearInterval(interval)
      }
    }, 500)
  } catch (error) {
    alert('开始录制失败: ' + error)
  }
}

async function stopRecording() {
  try {
    await window.go.main.App.StopRecording()
    recording.value = false
    eventCount.value = await window.go.main.App.GetEventCount()

    // 获取最终脚本并更新到编辑框
    try {
      const script = await window.go.main.App.GenerateCurrentScript()
      currentScript.value = script
    } catch (err) {
      console.error('获取脚本失败:', err)
    }

    alert('停止录制成功')
  } catch (error) {
    alert('停止录制失败: ' + error)
  }
}

async function saveMacro() {
  if (!macroName.value.trim()) {
    alert('请输入宏名称')
    return
  }

  try {
    const script = await window.go.main.App.SaveMacro(macroName.value.trim())
    alert('宏保存成功')
    macroName.value = ''
    eventCount.value = 0
    loadMacros()
  } catch (error) {
    alert('宏保存失败: ' + error)
  }
}

async function recordAndSave() {
  alert('请使用开始录制和停止录制功能')
}

// Macro Management
async function loadMacros() {
  try {
    macros.value = await window.go.main.App.ListMacros()
  } catch (error) {
    alert('加载宏失败: ' + error)
  }
}

async function playMacro(name) {
  try {
    await window.go.main.App.PlayMacro(name)
    alert('宏播放中...')
  } catch (error) {
    alert('宏播放失败: ' + error)
  }
}

async function deleteMacro(name) {
  try {
    await window.go.main.App.DeleteMacro(name)
    alert('宏删除成功')
    loadMacros()
    if (selectedMacro.value === name) {
      selectedMacro.value = ''
    }
  } catch (error) {
    alert('宏删除失败: ' + error)
  }
}

// Log Management (placeholder)
const logs = ref(['应用已启动'])
function addLog(log) {
  logs.value.unshift(`${new Date().toLocaleTimeString()}: ${log}`)
  if (logs.value.length > 100) {
    logs.value.pop()
  }
}
</script>

<template>
  <div class="app-container">
    <header class="app-header">
      <h1>操作宏录制和回放工具</h1>
    </header>

    <main class="app-main">
      <div class="control-panel">
        <!-- Recording Controls -->
        <section class="recording-controls">
          <h2>录制控制</h2>
          <div class="button-group">
            <button
              class="btn btn-primary"
              :disabled="recording"
              @click="startRecording"
            >
              开始录制
            </button>
            <button
              class="btn btn-danger"
              :disabled="!recording"
              @click="stopRecording"
            >
              停止录制
            </button>
            <button
              class="btn btn-success"
              :disabled="recording || eventCount === 0"
              @click="saveMacro"
            >
              保存宏
            </button>
            <button
              class="btn btn-secondary"
              :disabled="recording"
              @click="recordAndSave"
            >
              录制并保存
            </button>
          </div>

          <div class="status-info">
            <div>
              <span class="label">录制状态: </span>
              <span class="status" :class="recording ? 'recording' : 'stopped'">
                {{ recording ? '录制中' : '未录制' }}
              </span>
            </div>
            <div>
              <span class="label">事件数: </span>
              <span class="event-count">{{ eventCount }}</span>
            </div>
          </div>

          <div class="input-group">
            <input
              type="text"
              v-model="macroName"
              placeholder="请输入宏名称"
              :disabled="recording"
            />
          </div>
        </section>

        <!-- Macro Management -->
        <section class="macro-management">
          <h2>宏管理</h2>

          <div class="macro-list">
            <div class="list-header">
              <span>宏名称</span>
              <span>操作</span>
            </div>
            <div
              class="list-item"
              v-for="macro in macros"
              :key="macro.meta.name"
              :class="{ selected: selectedMacro === macro.meta.name }"
            >
              <span>{{ macro.meta.name }}</span>
              <div class="button-group">
                <button
                  class="btn btn-small btn-primary"
                  @click="playMacro(macro.meta.name)"
                >
                  播放
                </button>
                <button
                  class="btn btn-small btn-danger"
                  @click="deleteMacro(macro.meta.name)"
                >
                  删除
                </button>
              </div>
            </div>
            <div v-if="macros.length === 0" class="empty-list">
              暂无保存的宏
            </div>
          </div>
        </section>
      </div>

      <!-- Content Panels -->
      <div class="content-panels">
        <!-- Script Editor -->
        <section class="script-editor">
          <h2>当前脚本</h2>
          <textarea
            v-model="currentScript"
            placeholder="录制期间将实时显示生成的脚本..."
            readonly
          ></textarea>
        </section>

        <!-- Logs Panel -->
        <section class="logs-panel">
          <h2>操作日志</h2>
          <div class="logs-container">
            <div class="log-item" v-for="log in logs" :key="log">
              {{ log }}
            </div>
          </div>
        </section>
      </div>
    </main>

    <footer class="app-footer">
      <p>操作宏录制和回放工具 © 2024</p>
    </footer>
  </div>
</template>

<style scoped>
/* Global Styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  background-color: #1b2636;
  color: #fff;
}

.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

/* Header */
.app-header {
  background: linear-gradient(135deg, #2c3e50 0%, #34495e 100%);
  padding: 1.5rem 2rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4);
  border-bottom: 1px solid #34495e;
}

.app-header h1 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #ecf0f1;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

/* Main Content */
.app-main {
  flex: 1;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}

.control-panel {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 2rem;
  min-height: 350px;
}

@media (max-width: 968px) {
  .control-panel {
    grid-template-columns: 1fr;
  }
}

/* Recording Controls */
.recording-controls {
  background-color: #2c3e50;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
  border: 1px solid #34495e;
  transition: all 0.3s ease;
}

.recording-controls:hover {
  box-shadow: 0 6px 25px rgba(0, 0, 0, 0.4);
  border-color: #3498db;
}

.recording-controls h2,
.macro-management h2,
.script-editor h2,
.logs-panel h2 {
  margin-bottom: 1.5rem;
  font-size: 1.2rem;
  color: #3498db;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.button-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 600;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  letter-spacing: 0.5px;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background-color: #3498db;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #2980b9;
}

.btn-danger {
  background-color: #e74c3c;
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background-color: #c0392b;
}

.btn-success {
  background-color: #2ecc71;
  color: white;
}

.btn-success:hover:not(:disabled) {
  background-color: #27ae60;
}

.btn-secondary {
  background-color: #95a5a6;
  color: white;
}

.btn-secondary:hover:not(:disabled) {
  background-color: #7f8c8d;
}

.btn-small {
  padding: 0.5rem 0.75rem;
  font-size: 0.9rem;
}

.status-info {
  margin-bottom: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.label {
  font-weight: 500;
}

.status {
  font-weight: bold;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
}

.status.recording {
  background-color: #e74c3c;
  color: white;
}

.status.stopped {
  background-color: #2ecc71;
  color: white;
}

.event-count {
  font-weight: bold;
  color: #3498db;
}

.input-group {
  margin-top: 1rem;
}

.input-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #34495e;
  border-radius: 4px;
  background-color: #1b2636;
  color: white;
  font-size: 1rem;
}

.input-group input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
}

/* Macro Management */
.macro-management {
  flex: 1;
  background-color: #2c3e50;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
  border: 1px solid #34495e;
  transition: all 0.3s ease;
  min-width: 300px;
  display: flex;
  flex-direction: column;
}

.macro-management:hover {
  box-shadow: 0 6px 25px rgba(0, 0, 0, 0.4);
  border-color: #3498db;
}

.macro-management h2 {
  margin-bottom: 1.5rem;
  font-size: 1.2rem;
  color: #3498db;
}

.macro-list {
  flex: 1;
  min-height: 200px;
  max-height: 400px;
  overflow-y: auto;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 0;
  border-bottom: 1px solid #34495e;
  font-weight: bold;
}

.list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 0;
  border-bottom: 1px solid #34495e;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.list-item:hover {
  background-color: rgba(52, 152, 219, 0.1);
}

.list-item.selected {
  background-color: rgba(52, 152, 219, 0.2);
}

.empty-list {
  text-align: center;
  padding: 2rem;
  color: #95a5a6;
  font-style: italic;
}

/* Logs Panel */
.logs-panel {
  background-color: #2c3e50;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
  max-height: 250px;
}

.logs-panel h2 {
  margin-bottom: 1.5rem;
  font-size: 1.2rem;
  color: #3498db;
}

.logs-container {
  max-height: 150px;
  overflow-y: auto;
  background-color: #1b2636;
  padding: 1rem;
  border-radius: 4px;
}

.log-item {
  margin-bottom: 0.5rem;
  color: #bdc3c7;
  font-size: 0.9rem;
}

/* Content Panels */
.content-panels {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 2rem;
  min-height: 300px;
}

.script-editor,
.logs-panel {
  transition: all 0.3s ease;
}

.script-editor:hover,
.logs-panel:hover {
  box-shadow: 0 6px 25px rgba(0, 0, 0, 0.4);
  border-color: #3498db;
}

/* Script Editor */
.script-editor {
  background-color: #2c3e50;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
}

.script-editor h2 {
  margin-bottom: 1.5rem;
  font-size: 1.2rem;
  color: #3498db;
  flex-shrink: 0;
}

.script-editor textarea {
  flex: 1;
  width: 100%;
  padding: 1rem;
  border: 1px solid #34495e;
  border-radius: 4px;
  background-color: #1b2636;
  color: #fff;
  font-size: 0.9rem;
  font-family: 'Courier New', Courier, monospace;
  resize: vertical;
  min-height: 250px;
  overflow-y: auto;
}

/* Logs Panel */
.logs-panel {
  background-color: #2c3e50;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
}

.logs-panel h2 {
  margin-bottom: 1.5rem;
  font-size: 1.2rem;
  color: #3498db;
  flex-shrink: 0;
}

.logs-container {
  flex: 1;
  max-height: 100%;
  overflow-y: auto;
  background-color: #1b2636;
  padding: 1rem;
  border-radius: 4px;
}

@media (max-width: 1280px) {
  .content-panels {
    grid-template-columns: 1fr;
  }

  .script-editor,
  .logs-panel {
    min-height: 250px;
  }
}

/* Footer */
.app-footer {
  background-color: #2c3e50;
  padding: 1rem;
  text-align: center;
  color: #95a5a6;
  margin-top: auto;
}

/* Responsive Design */
@media (max-width: 768px) {
  .control-panel {
    flex-direction: column;
  }

  .recording-controls {
    width: 100%;
  }
}
</style>
