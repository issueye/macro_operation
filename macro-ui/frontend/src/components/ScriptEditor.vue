<template>
    <div class="card script-card">
        <div class="card-header">
            <h3>脚本预览</h3>
            <div class="script-actions">
                <button
                    class="btn btn-secondary"
                    @click="emit('play')"
                    :disabled="!script || isRecording || isPlaying"
                >
                    <svg viewBox="0 0 24 24" width="14" height="14">
                        <path fill="currentColor" d="M8 5v14l11-7z" />
                    </svg>
                    播放
                </button>
                <button
                    class="btn btn-secondary"
                    @click="handleSave"
                    :disabled="!script"
                >
                    <svg viewBox="0 0 24 24" width="14" height="14">
                        <path
                            fill="currentColor"
                            d="M17 3H5c-1.11 0-2 .9-2 2v14c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V7l-4-4zm-5 16c-1.66 0-3-1.34-3-3s1.34-3 3-3 3 1.34 3 3-1.34 3-3 3zm3-10H5V5h10v4z"
                        />
                    </svg>
                    保存
                </button>
                <button
                    class="btn btn-secondary"
                    @click="handleSaveAsNew"
                    :disabled="!script"
                    title="另存为新宏"
                >
                    <svg viewBox="0 0 24 24" width="14" height="14">
                        <path
                            fill="currentColor"
                            d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"
                        />
                    </svg>
                    新建
                </button>
            </div>
        </div>
        <div class="script-editor">
            <MonacoEditor
                ref="monacoEditor"
                :model-value="script"
                @update:model-value="emit('update:script', $event)"
                @save="emit('save', $event)"
                language="javascript"
                theme="vs-dark"
                :read-only="false"
            />
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
                <button class="btn btn-primary btn-small" @click="confirmSave">
                    保存
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, nextTick } from "vue";
import MonacoEditor from "./MonacoEditor.vue";

const props = defineProps({
    script: {
        type: String,
        default: "",
    },
    isRecording: {
        type: Boolean,
        default: false,
    },
    isPlaying: {
        type: Boolean,
        default: false,
    },
    selectedMacro: {
        type: String,
        default: null,
    },
});

const emit = defineEmits(["update:script", "play", "save", "show-message"]);

const showSaveInput = ref(false);
const macroName = ref("");
const saveInput = ref(null);

const handleSave = () => {
    showSaveInput.value = true;
    macroName.value = props.selectedMacro || "";
    nextTick(() => {
        saveInput.value?.focus();
    });
};

const handleSaveAsNew = () => {
    showSaveInput.value = true;
    macroName.value = "";
    nextTick(() => {
        saveInput.value?.focus();
    });
};

const cancelSave = () => {
    showSaveInput.value = false;
    macroName.value = "";
};

const confirmSave = () => {
    if (!macroName.value.trim()) {
        emit("show-message", "请输入宏名称", "warn");
        return;
    }
    emit("save", macroName.value.trim());
    macroName.value = "";
    showSaveInput.value = false;
};
</script>
