<template>
    <div ref="editorContainer" class="simple-editor-container">
        <!-- 行号 -->
        <div class="line-numbers">
            <div
                v-for="(line, index) in lines"
                :key="index"
                class="line-number"
            >
                {{ index + 1 }}
            </div>
        </div>

        <!-- 代码编辑区 -->
        <textarea
            ref="textareaRef"
            class="code-textarea"
            :value="modelValue"
            @input="handleInput"
            @keydown="handleKeydown"
            @scroll="syncScroll"
            @click="updateCursorPosition"
            @keyup="updateCursorPosition"
            :readonly="readOnly"
            spellcheck="false"
            autocomplete="off"
            autocorrect="off"
            autocapitalize="off"
            lang="en"
        ></textarea>

        <!-- 光标状态显示 -->
        <div class="cursor-status" v-if="cursorPosition">
            Ln {{ cursorPosition.line }}, Col {{ cursorPosition.column }}
        </div>
    </div>
</template>

<script setup>
import {
    ref,
    computed,
    watch,
    onMounted,
    onBeforeUnmount,
    nextTick,
} from "vue";

const props = defineProps({
    modelValue: {
        type: String,
        default: "",
    },
    language: {
        type: String,
        default: "javascript",
    },
    theme: {
        type: String,
        default: "vs-dark",
    },
    readOnly: {
        type: Boolean,
        default: false,
    },
    options: {
        type: Object,
        default: () => ({}),
    },
});

const emit = defineEmits(["update:modelValue", "change", "save"]);

const textareaRef = ref(null);
const cursorPosition = ref(null);
const isComposing = ref(false); // 输入法 composing 状态

// 计算行数
const lines = computed(() => {
    if (!props.modelValue) return [1];
    return props.modelValue.split("\n");
});

const handleInput = (e) => {
    // 输入法正在输入时，不处理
    if (isComposing.value) {
        return;
    }

    const value = e.target.value;
    emit("update:modelValue", value);
    emit("change", value);
    updateCursorPosition();
};

// 输入法开始
const handleCompositionStart = () => {
    isComposing.value = true;
};

// 输入法结束
const handleCompositionEnd = (e) => {
    isComposing.value = false;
    const value = e.target.value;
    emit("update:modelValue", value);
    emit("change", value);
    updateCursorPosition();
};

const handleKeydown = (e) => {
    // 输入法正在输入时，不处理特殊键
    if (isComposing.value) {
        return;
    }

    // Ctrl+S 保存
    if ((e.ctrlKey || e.metaKey) && e.key === "s") {
        e.preventDefault();
        emit("save", textareaRef.value.value);
        return;
    }

    // Tab 键支持
    if (e.key === "Tab" && !props.readOnly) {
        e.preventDefault();
        const textarea = e.target;
        const start = textarea.selectionStart;
        const end = textarea.selectionEnd;
        const value = textarea.value;

        // 插入两个空格
        const newValue =
            value.substring(0, start) + "  " + value.substring(end);
        textarea.value = newValue;

        // 恢复光标位置
        textarea.selectionStart = textarea.selectionEnd = start + 2;

        // 触发更新
        emit("update:modelValue", newValue);
        emit("change", newValue);
        return;
    }

    // 自动括号匹配（仅在英文输入状态下）
    if (!props.readOnly && !e.ctrlKey && !e.metaKey && !e.altKey) {
        const pairs = {
            "(": ")",
            "[": "]",
            "{": "}",
            '"': '"',
            "'": "'",
        };

        if (pairs[e.key]) {
            const textarea = e.target;
            const start = textarea.selectionStart;
            const end = textarea.selectionEnd;
            const value = textarea.value;
            const selectedText = value.substring(start, end);
            const closeChar = pairs[e.key];

            e.preventDefault();
            const newValue =
                value.substring(0, start) +
                e.key +
                selectedText +
                closeChar +
                value.substring(end);
            textarea.value = newValue;

            // 光标放在插入的字符后面
            textarea.selectionStart = textarea.selectionEnd = start + 1;

            emit("update:modelValue", newValue);
            emit("change", newValue);
        }
    }
};

const syncScroll = (e) => {
    const lineNumbers = e.target.previousElementSibling;
    if (lineNumbers) {
        lineNumbers.scrollTop = e.target.scrollTop;
    }
};

const updateCursorPosition = () => {
    if (!textareaRef.value) return;

    const textarea = textareaRef.value;
    const text = textarea.value.substring(0, textarea.selectionStart);
    const lines = text.split("\n");

    cursorPosition.value = {
        line: lines.length,
        column: lines[lines.length - 1].length + 1,
    };
};

const getValue = () => {
    return textareaRef.value ? textareaRef.value.value : "";
};

const setValue = (value) => {
    if (textareaRef.value) {
        textareaRef.value.value = value;
    }
};

const focus = () => {
    if (textareaRef.value) {
        textareaRef.value.focus();
    }
};

defineExpose({
    getValue,
    setValue,
    focus,
});

// 监听外部值变化
watch(
    () => props.modelValue,
    (newValue) => {
        if (textareaRef.value && textareaRef.value.value !== newValue) {
            textareaRef.value.value = newValue;
        }
    },
);

onMounted(() => {
    nextTick(() => {
        updateCursorPosition();

        // 添加输入法事件监听
        if (textareaRef.value) {
            textareaRef.value.addEventListener(
                "compositionstart",
                handleCompositionStart,
            );
            textareaRef.value.addEventListener(
                "compositionend",
                handleCompositionEnd,
            );
        }
    });
});

onBeforeUnmount(() => {
    if (textareaRef.value) {
        textareaRef.value.removeEventListener(
            "compositionstart",
            handleCompositionStart,
        );
        textareaRef.value.removeEventListener(
            "compositionend",
            handleCompositionEnd,
        );
    }
});
</script>

<style scoped>
.simple-editor-container {
    width: 100%;
    height: 100%;
    min-height: 200px;
    position: relative;
    display: flex;
    background: #1e1e1e;
    border-radius: 4px;
    overflow: hidden;
}

.line-numbers {
    width: 45px;
    background: #1e1e1e;
    border-right: 1px solid #3c3c3c;
    color: #858585;
    font-family: "Consolas", "Monaco", "Courier New", monospace;
    font-size: 14px;
    line-height: 1.6;
    padding: 12px 0;
    text-align: right;
    padding-right: 10px;
    overflow: hidden;
    user-select: none;
    flex-shrink: 0;
}

.line-number {
    height: 22.4px;
    padding-right: 8px;
    box-sizing: border-box;
}

.code-textarea {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    color: #d4d4d4;
    font-family: "Consolas", "Monaco", "Courier New", monospace;
    font-size: 14px;
    line-height: 1.6;
    padding: 12px;
    resize: none;
    tab-size: 2;
    white-space: pre;
    overflow-wrap: normal;
    overflow-x: auto;
    min-width: 0;
}

.code-textarea::selection {
    background: #264f78;
}

.code-textarea:focus {
    outline: none;
}

.cursor-status {
    position: absolute;
    bottom: 0;
    right: 0;
    background: rgba(0, 0, 0, 0.8);
    color: #858585;
    font-size: 11px;
    padding: 3px 8px;
    font-family: "Consolas", "Monaco", monospace;
    pointer-events: none;
    border-top-left-radius: 4px;
    z-index: 1;
}

/* 滚动条样式 */
.code-textarea::-webkit-scrollbar {
    width: 14px;
    height: 14px;
}

.code-textarea::-webkit-scrollbar-track {
    background: #1e1e1e;
}

.code-textarea::-webkit-scrollbar-thumb {
    background: #424242;
    border-radius: 7px;
    border: 3px solid #1e1e1e;
}

.code-textarea::-webkit-scrollbar-thumb:hover {
    background: #4f4f4f;
}

.line-numbers::-webkit-scrollbar {
    width: 0;
}
</style>
