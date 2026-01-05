<template>
    <div class="app">
        <!-- 顶部状态栏 -->
        <Header :connected="connected" />

        <div class="main-layout">
            <!-- 左侧边栏 -->
            <aside class="sidebar">
                <!-- 录制控制 -->
                <RecordControl
                    :is-recording="isRecording"
                    :is-playing="isPlaying"
                    :event-count="eventCount"
                    @toggle-recording="toggleRecording"
                />

                <!-- 宏列表 -->
                <MacroList
                    :macros="macros"
                    :selected-macro="selectedMacro"
                    @refresh="loadMacros"
                    @select="selectMacro"
                    @play="playMacro"
                    @delete="deleteMacro"
                />
            </aside>

            <!-- 右侧主内容 -->
            <main class="main-content">
                <!-- 内容区域：左右布局 -->
                <div class="content-area">
                    <!-- 左侧：脚本编辑器 -->
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

                    <!-- 右侧：事件监控和日志 -->
                    <div class="right-panel">
                        <!-- 事件监控面板 -->
                        <EventMonitor
                            ref="eventMonitor"
                            :is-recording="isRecording"
                        />

                        <!-- 日志面板 -->
                        <LogPanel :logs="logs" @clear="clearLogs" />
                    </div>
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
const script = ref("");
const selectedMacro = ref(null);
const macros = ref([]);
const logs = ref([]);
const statusTimer = ref(null);
const eventMonitor = ref(null);

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
    if (logs.value.length > 50) {
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
        script.value = "";
        selectedMacro.value = null;
        addLog("开始录制", "info");
    } catch (e) {
        addLog("开始录制失败: " + e.message, "error");
    }
};

const stopRecording = async () => {
    try {
        await window.go.main.App.StopRecording();
        isRecording.value = false;
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
    await nextTick();
    await checkConnection();
    await loadMacros();
    startStatusPolling();
});

onBeforeUnmount(() => {
    stopStatusPolling();
});

// 监听录制状态变化，控制事件监控
watch(isRecording, (newVal) => {
    if (newVal) {
        eventMonitor.value?.startPolling();
    } else {
        eventMonitor.value?.stopPolling();
    }
});
</script>
