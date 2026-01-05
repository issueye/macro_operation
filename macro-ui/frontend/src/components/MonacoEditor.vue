<template>
    <div
        ref="editorContainer"
        class="monaco-editor-container"
        :style="{ height: '100%' }"
    ></div>
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount, nextTick } from "vue";

// Monaco Editor will be dynamically imported
let monaco = null;
let editor = null;

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
        default: "macro-dark",
    },
    readOnly: {
        type: Boolean,
        default: false,
    },
    height: {
        type: Number,
        default: 400,
    },
    options: {
        type: Object,
        default: () => ({}),
    },
});

const emit = defineEmits([
    "update:modelValue",
    "change",
    "save",
    "editor-mounted",
]);

const editorContainer = ref(null);

// Define custom dark theme
const defineMacroTheme = (monaco) => {
    monaco.editor.defineTheme("macro-dark", {
        base: "vs-dark",
        inherit: true,
        rules: [
            { token: "keyword", foreground: "c792ea", fontStyle: "bold" },
            { token: "string", foreground: "c3e88d" },
            { token: "number", foreground: "f78c6c" },
            { token: "function", foreground: "82aaff" },
            { token: "variable", foreground: "f8f8f2" },
            { token: "comment", foreground: "676e95", fontStyle: "italic" },
            { token: "type", foreground: "ffcb6b" },
            { token: "operator", foreground: "89ddff" },
            { token: "delimiter", foreground: "8992a6" },
            { token: "property", foreground: "ffcb6b" },
        ],
        colors: {
            "editor.background": "#0f172a",
            "editor.foreground": "#f1f5f9",
            "editor.lineHighlightBackground": "#1e293b",
            "editor.selectionBackground": "#334155",
            "editorCursor.foreground": "#6366f1",
            "editorLineNumber.foreground": "#64748b",
            "editorLineNumber.activeForeground": "#94a3b8",
            "editor.inactiveSelectionBackground": "#1e293b",
            "editor.selectionHighlightBackground": "#33415580",
            "editorIndentGuide.background": "#1e293b",
            "editorIndentGuide.activeBackground": "#334155",
            "editorBracketMatch.background": "#334155",
            "editorBracketMatch.border": "#6366f1",
        },
    });
};

// Initialize editor
const initEditor = async () => {
    if (!editorContainer.value) return;

    // Dynamically import Monaco Editor
    try {
        const monacoModule = await import("monaco-editor");
        monaco = monacoModule;

        // Define custom theme
        defineMacroTheme(monaco);

        // Merge options
        const editorOptions = {
            value: props.modelValue,
            language: props.language,
            theme: props.theme,
            readOnly: props.readOnly,
            automaticLayout: true,
            fontFamily: "'JetBrains Mono', 'Fira Code', 'Consolas', monospace",
            fontSize: 14,
            lineHeight: 22,
            minimap: {
                enabled: true,
                scale: 0.8,
                maxColumn: 100,
            },
            scrollBeyondLastLine: false,
            wordWrap: "off",
            tabSize: 2,
            insertSpaces: true,
            renderLineHighlight: "all",
            cursorBlinking: "smooth",
            cursorSmoothCaretAnimation: "on",
            smoothScrolling: true,
            padding: {
                top: 12,
                bottom: 12,
            },
            bracketPairColorization: {
                enabled: true,
            },
            guides: {
                bracketPairs: true,
                indentation: true,
            },
            ...props.options,
        };

        // Create editor
        editor = monaco.editor.create(editorContainer.value, editorOptions);

        // Listen for content changes
        editor.onDidChangeModelContent(() => {
            const value = editor.getValue();
            emit("update:modelValue", value);
            emit("change", value);
        });

        // Listen for save (Ctrl+S)
        editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, () => {
            emit("save", editor.getValue());
        });

        // Focus on mount
        editor.focus();

        // Emit editor mounted event
        emit("editor-mounted", editor);
    } catch (error) {
        console.error("Failed to load Monaco Editor:", error);
    }
};

// Update editor value
watch(
    () => props.modelValue,
    (newValue) => {
        if (editor && editor.getValue() !== newValue) {
            editor.setValue(newValue);
        }
    },
);

// Update language
watch(
    () => props.language,
    (newLanguage) => {
        if (editor) {
            const model = editor.getModel();
            if (model) {
                monaco.editor.setModelLanguage(model, newLanguage);
            }
        }
    },
);

// Update readOnly
watch(
    () => props.readOnly,
    (readOnly) => {
        if (editor) {
            editor.updateOptions({ readOnly });
        }
    },
);

// Update theme
watch(
    () => props.theme,
    (newTheme) => {
        if (editor && monaco) {
            monaco.editor.setTheme(newTheme);
        }
    },
);

// Get editor value
const getValue = () => {
    return editor ? editor.getValue() : "";
};

// Set editor value
const setValue = (value) => {
    if (editor) {
        editor.setValue(value);
    }
};

// Focus editor
const focus = () => {
    if (editor) {
        editor.focus();
    }
};

// Resize editor
const resize = () => {
    if (editor) {
        editor.layout();
    }
};

// Expose methods
defineExpose({
    getValue,
    setValue,
    focus,
    resize,
    getEditor: () => editor,
});

onMounted(async () => {
    await nextTick();
    await initEditor();
});

onBeforeUnmount(() => {
    if (editor) {
        editor.dispose();
        editor = null;
    }
});
</script>

<style scoped>
.monaco-editor-container {
    width: 100%;
    height: 100%;
    min-height: 150px;
    border-radius: var(--radius-md);
    overflow: hidden;
    background: var(--bg-base);
}

.monaco-editor-container :deep(.monaco-editor) {
    padding-top: 12px;
    padding-bottom: 12px;
}

.monaco-editor-container :deep(.minimap) {
    opacity: 0.8;
}

.monaco-editor-container
    :deep(.monaco-scrollable-element > .scrollbar > .slider) {
    background: var(--bg-element);
    border-radius: var(--radius-full);
}

.monaco-editor-container
    :deep(.monaco-scrollable-element > .scrollbar > .slider:hover) {
    background: var(--bg-hover);
}
</style>
