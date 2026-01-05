<template>
    <div class="record-card" :class="{ collapsed }">
        <!-- 录制按钮区域 -->
        <div class="record-button-wrapper">
            <button
                class="btn-record"
                :class="{
                    recording: isRecording,
                    playing: isPlaying,
                    disabled: isPlaying && !isRecording
                }"
                @click="emit('toggle-recording')"
                :disabled="isPlaying && !isRecording"
                :title="isRecording ? '点击停止录制' : '点击开始录制'"
            >
                <!-- 录制中状态 -->
                <template v-if="isRecording">
                    <span class="record-icon stop-icon">
                        <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor">
                            <rect x="6" y="6" width="12" height="12" rx="2" />
                        </svg>
                    </span>
                    <span class="btn-text">停止录制</span>
                </template>

                <!-- 默认状态 -->
                <template v-else>
                    <span class="record-icon" :class="{ disabled: isPlaying }">
                        <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor">
                            <circle cx="12" cy="12" r="8" />
                        </svg>
                    </span>
                    <span class="btn-text">{{ isPlaying ? '播放中...' : '开始录制' }}</span>
                </template>

                <!-- 脉动光晕效果 -->
                <span class="pulse-ring" v-if="isRecording"></span>
            </button>
        </div>

        <!-- 录制信息 -->
        <div class="record-info" v-if="isRecording || isPlaying">
            <div class="info-row duration">
                <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10" />
                    <polyline points="12 6 12 12 16 14" />
                </svg>
                <span class="duration-text">{{ formattedDuration }}</span>
            </div>
            <div class="info-row events">
                <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
                </svg>
                <span class="event-count">{{ eventCount }} 事件</span>
            </div>
        </div>

        <!-- 快捷键提示 -->
        <div class="shortcut-hint" v-if="!isRecording && !isPlaying && !collapsed">
            <kbd>Ctrl</kbd> + <kbd>R</kbd>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
    isRecording: {
        type: Boolean,
        default: false
    },
    isPlaying: {
        type: Boolean,
        default: false
    },
    eventCount: {
        type: Number,
        default: 0
    },
    recordingDuration: {
        type: Number,
        default: 0
    },
    collapsed: {
        type: Boolean,
        default: false
    }
});

const emit = defineEmits(['toggle-recording']);

// 格式化时长显示
const formattedDuration = computed(() => {
    const seconds = props.recordingDuration;
    const hrs = Math.floor(seconds / 3600);
    const mins = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;

    if (hrs > 0) {
        return `${String(hrs).padStart(2, '0')}:${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
    }
    return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
});
</script>

<style scoped>
.record-card {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
}

/* 折叠状态下的样式 */
.record-card.collapsed {
    align-items: center;
}

.record-card.collapsed .record-button-wrapper {
    width: 100%;
}

.record-card.collapsed .btn-record {
    width: 44px;
    height: 44px;
    padding: 0;
    border-radius: 50%;
}

.record-card.collapsed .btn-text,
.record-card.collapsed .record-info,
.record-card.collapsed .shortcut-hint {
    display: none;
}

/* 录制按钮容器 */
.record-button-wrapper {
    display: flex;
    justify-content: center;
}

/* 录制按钮 */
.btn-record {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--space-2);
    width: 100%;
    height: 48px;
    padding: 0 var(--space-4);
    border: none;
    border-radius: var(--radius-lg);
    font-size: var(--text-base);
    font-weight: 600;
    cursor: pointer;
    transition: all var(--transition-normal);
    overflow: hidden;
}

/* 默认状态 - 主色 */
.btn-record:not(.recording):not(.playing):not(.disabled) {
    background: var(--color-primary);
    color: var(--text-primary);
    box-shadow: var(--shadow-md), 0 0 0 0 transparent;
}

.btn-record:not(.recording):not(.playing):not(.disabled):hover {
    background: var(--color-primary-hover);
    transform: translateY(-2px);
    box-shadow: var(--shadow-lg), var(--shadow-glow-primary);
}

.btn-record:not(.recording):not(.playing):not(.disabled):active {
    background: var(--color-primary-active);
    transform: translateY(0);
}

/* 录制中状态 - 红色 */
.btn-record.recording {
    background: var(--color-recording);
    color: var(--text-primary);
    box-shadow: var(--shadow-md), var(--shadow-glow-recording);
    animation: pulse-recording 2s ease-in-out infinite;
}

.btn-record.recording:hover {
    animation-play-state: paused;
}

/* 播放中状态（禁用） */
.btn-record.playing,
.btn-record.disabled {
    background: var(--bg-element);
    color: var(--text-muted);
    cursor: not-allowed;
    opacity: 0.6;
}

/* 图标样式 */
.record-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    flex-shrink: 0;
}

.record-icon circle {
    fill: currentColor;
}

.record-icon.disabled circle {
    opacity: 0.5;
}

/* 停止图标 */
.stop-icon {
    color: var(--text-primary);
}

/* 脉动光晕 */
.pulse-ring {
    position: absolute;
    top: 50%;
    left: 50%;
    width: 100%;
    height: 100%;
    border-radius: inherit;
    transform: translate(-50%, -50%) scale(1);
    pointer-events: none;
}

.btn-record.recording .pulse-ring {
    background: var(--color-recording);
    animation: ring-pulse 2s ease-out infinite;
}

@keyframes ring-pulse {
    0% {
        transform: translate(-50%, -50%) scale(1);
        opacity: 0.4;
    }
    100% {
        transform: translate(-50%, -50%) scale(1.8);
        opacity: 0;
    }
}

/* 录制信息 */
.record-info {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
    padding: var(--space-3);
    background: var(--bg-base);
    border-radius: var(--radius-md);
    animation: fade-in var(--transition-normal);
}

.info-row {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    font-size: var(--text-sm);
    color: var(--text-secondary);
}

.info-row svg {
    opacity: 0.7;
}

.info-row.duration {
    color: var(--color-recording);
}

.duration-text {
    font-family: var(--font-mono);
    font-weight: 500;
    font-variant-numeric: tabular-nums;
}

.info-row.events {
    color: var(--text-muted);
}

.event-count {
    font-family: var(--font-mono);
    font-variant-numeric: tabular-nums;
}

/* 快捷键提示 */
.shortcut-hint {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--space-1);
    font-size: var(--text-xs);
    color: var(--text-muted);
}

.shortcut-hint kbd {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 20px;
    height: 20px;
    padding: 0 6px;
    background: var(--bg-element);
    border: 1px solid var(--border-default);
    border-radius: var(--radius-sm);
    font-family: var(--font-sans);
    font-size: 10px;
    font-weight: 500;
}

/* 脉动动画关键帧 */
@keyframes pulse-recording {
    0%, 100% {
        box-shadow: var(--shadow-md), var(--shadow-glow-recording);
    }
    50% {
        box-shadow: var(--shadow-md), 0 0 0 8px var(--color-recording-glow-outer);
    }
}

/* 淡入动画 */
@keyframes fade-in {
    0% {
        opacity: 0;
        transform: translateY(-4px);
    }
    100% {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>
