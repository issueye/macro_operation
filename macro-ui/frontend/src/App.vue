<template>
    <div class="app">
        <!-- 顶部状态栏 -->
        <Header
            :connected="connected"
            :is-recording="isRecording"
            :is-playing="isPlaying"
            :recording-duration="recordingDuration"
            :event-count="eventCount"
        />

        <div class="main-layout">
            <!-- 左侧边栏 -->
            <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
                <div class="sidebar-header">
                    <span class="sidebar-title">{{ sidebarCollapsed ? '' : 'My Macros' }}</span>
                    <button
                        class="sidebar-toggle"
                        @click="toggleSidebar"
                        :title="sidebarCollapsed ? '展开' : '收起'"
                    >
                        <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                            <path v-if="sidebarCollapsed" d="M9 18l6-6-6-6" />
                            <path v-else d="M15 18l-6-6 6-6" />
                        </svg>
                    </button>
                </div>

                <div class="sidebar-content">
                    <!-- 录制控制 -->
                    <RecordControl
                        :is-recording="isRecording"
                        :is-playing="isPlaying"
                        :event-count="eventCount"
                        :recording-duration="recordingDuration"
                        @toggle-recording="toggleRecording"
                    />

                    <!-- 宏列表 -->
                    <MacroList
                        :macros="macros"
                        :selected-macro="selectedMacro"
                        :collapsed="sidebarCollapsed"
                        @refresh="loadMacros"
                        @select="selectMacro"
                        @play="playMacro"
                        @delete="deleteMacro"
                    />
                </div>
            </aside>

            <!-- 右侧主内容 -->
            <main class="main-content">
                <!-- 脚本编辑区 -->
                <div class="script-section">
                    <ScriptEditor
                        v-model:script="script"
                        :is-recording="isRecording"
                        :is-playing="isPlaying"
                        :selected-macro="selectedMacro"
                        @play="playCurrentScript"
                        @save="saveMacro"
                        @show-message="addLog"
                    />
                </div>

                <!-- 右侧面板：事件监控和日志（左右布局） -->
                <div class="right-panel">
                    <EventMonitor
                        ref="eventMonitor"
                        :is-recording="isRecording"
                        class="event-card"
                    />
                    <LogPanel :logs="logs" @clear="clearLogs" class="log-card" />
                </div>
            </main>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick, watch } from "vue";
import Header from "./components/Header.vue";
import RecordControl from "./components/RecordControl.vue";
import MacroList from "./components/MacroList.vue";
import EventMonitor from "./components/EventMonitor.vue";
import ScriptEditor from "./components/ScriptEditor.vue";
import LogPanel from "./components/LogPanel.vue";

// 状态
const isRecording = ref(false);
const isPlaying = ref(false);
const connected = ref(false);
const eventCount = ref(0);
const recordingDuration = ref(0);
const script = ref("");
const selectedMacro = ref(null);
const macros = ref([]);
const logs = ref([]);
const statusTimer = ref(null);
const durationTimer = ref(null);
const eventMonitor = ref(null);

// 侧边栏折叠状态
const sidebarCollapsed = ref(false);

// 侧边栏折叠切换
const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value;
    // 保存设置到本地存储
    localStorage.setItem('sidebarCollapsed', sidebarCollapsed.value);
};

// 加载侧边栏设置
const loadSidebarSetting = () => {
    const saved = localStorage.getItem('sidebarCollapsed');
    if (saved !== null) {
        sidebarCollapsed.value = saved === 'true';
    }
};

// 日志管理
const addLog = (message, type = "success") => {
    const now = new Date();
    const time = now.toLocaleTimeString();
    logs.value.unshift({
        id: Date.now() + Math.random(),
        time,
        message,
        type,
    });
    if (logs.value.length > 100) {
        logs.value.pop();
    }
};

const clearLogs = () => {
    logs.value = [];
};

// 连接检查
const checkConnection = async () => {
    try {
        await window.go.main.App.IsRecording();
        connected.value = true;
        addLog("引擎服务连接成功", "success");
    } catch (e) {
        connected.value = false;
        addLog("无法连接到引擎服务", "error");
    }
};

// 状态轮询
const startStatusPolling = () => {
    statusTimer.value = setInterval(async () => {
        if (isRecording.value) {
            try {
                const count = await window.go.main.App.GetEventCount();
                eventCount.value = count;
            } catch (e) {
                console.error("获取状态失败:", e);
            }
        }
    }, 500);
};

const stopStatusPolling = () => {
    if (statusTimer.value) {
        clearInterval(statusTimer.value);
    }
};

// 录制时长计时器
const startDurationTimer = () => {
    recordingDuration.value = 0;
    durationTimer.value = setInterval(() => {
        recordingDuration.value++;
    }, 1000);
};

const stopDurationTimer = () => {
    if (durationTimer.value) {
        clearInterval(durationTimer.value);
        durationTimer.value = null;
    }
};

// 格式化时长
const formatDuration = (seconds) => {
    const hrs = Math.floor(seconds / 3600);
    const mins = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;
    if (hrs > 0) {
        return `${String(hrs).padStart(2, '0')}:${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
    }
    return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
};

// 录制控制
const toggleRecording = async () => {
    if (isRecording.value) {
        await stopRecording();
    } else {
        await startRecording();
    }
};

const startRecording = async () => {
    try {
        await window.go.main.App.StartRecording();
        isRecording.value = true;
        eventCount.value = 0;
        recordingDuration.value = 0;
        script.value = "";
        selectedMacro.value = null;
        startDurationTimer();
        addLog("开始录制", "info");
    } catch (e) {
        addLog("开始录制失败: " + e.message, "error");
    }
};

const stopRecording = async () => {
    try {
        await window.go.main.App.StopRecording();
        isRecording.value = false;
        stopDurationTimer();
        addLog(`录制完成，共 ${eventCount.value} 个事件`, "success");
        await refreshScript();
    } catch (e) {
        addLog("停止录制失败: " + e.message, "error");
    }
};

const refreshScript = async () => {
    try {
        script.value = await window.go.main.App.GenerateCurrentScript();
        addLog("脚本已生成", "success");
    } catch (e) {
        addLog("生成脚本失败: " + e.message, "error");
    }
};

// 脚本播放
const playCurrentScript = async () => {
    if (!script.value) {
        addLog("没有可播放的脚本", "warn");
        return;
    }
    try {
        isPlaying.value = true;
        addLog("开始播放脚本", "info");
        await window.go.main.App.PlayScript(script.value);
        addLog("脚本播放完成", "success");
    } catch (e) {
        addLog("播放失败: " + e.message, "error");
    } finally {
        isPlaying.value = false;
    }
};

// 宏管理
const selectMacro = (macro) => {
    selectedMacro.value = macro.name;
    script.value = macro.script || "";
};

const playMacro = async (name) => {
    try {
        const macroScript = await window.go.main.App.LoadMacro(name);
        if (!macroScript) {
            addLog("宏 " + name + " 没有脚本", "warn");
            return;
        }
        script.value = macroScript;
        selectedMacro.value = name;
        isPlaying.value = true;
        addLog("播放宏: " + name, "info");
        await window.go.main.App.PlayScript(macroScript);
        addLog("宏播放完成", "success");
    } catch (e) {
        addLog("播放宏失败: " + e.message, "error");
    } finally {
        isPlaying.value = false;
    }
};

const saveMacro = async (name) => {
    try {
        await window.go.main.App.SaveMacro(name, script.value);
        addLog(`宏 "${name}" 已保存`, "success");
        selectedMacro.value = name;
        await loadMacros();
    } catch (e) {
        addLog("保存宏失败: " + e.message, "error");
    }
};

const deleteMacro = async (name) => {
    if (!confirm(`确定要删除宏 "${name}" 吗？`)) return;
    try {
        await window.go.main.App.DeleteMacro(name);
        addLog(`宏 "${name}" 已删除`, "success");
        if (selectedMacro.value === name) {
            selectedMacro.value = null;
        }
        await loadMacros();
    } catch (e) {
        addLog("删除宏失败: " + e.message, "error");
    }
};

const loadMacros = async () => {
    try {
        macros.value = await window.go.main.App.ListMacros();
    } catch (e) {
        addLog("加载宏列表失败: " + e.message, "error");
    }
};

// 初始化
onMounted(async () => {
    addLog("应用启动", "info");
    loadSidebarSetting();
    await nextTick();
    await checkConnection();
    await loadMacros();
    startStatusPolling();
});

onBeforeUnmount(() => {
    stopStatusPolling();
    stopDurationTimer();
});

// 监听录制状态变化，控制事件监控
watch(isRecording, (newVal) => {
    if (newVal) {
        eventMonitor.value?.startPolling();
    } else {
        eventMonitor.value?.stopPolling();
    }
});

// 暴露格式化时长函数给模板
defineExpose({
    formatDuration
});
</script>
