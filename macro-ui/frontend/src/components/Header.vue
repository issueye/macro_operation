<template>
    <header class="header">
        <div class="logo">
            <span class="logo-icon">
                <svg viewBox="0 0 24 24" width="24" height="24" fill="currentColor">
                    <path d="M13 2.05v2.02c3.95.49 7 3.85 7 7.93 0 3.21-1.92 6-4.72 7.28L13 17v5h5l-1.22-1.22C19.91 19.07 22 15.76 22 12c0-5.18-3.95-9.45-9-9.95zM11 2.05C5.94 2.55 2 6.81 2 12c0 3.76 2.09 7.07 5.22 8.78L6 22h5v-5l-2.28 2.28C6.92 18 5 15.21 5 12c0-4.08 3.05-7.44 7-7.93V2.05z"/>
                </svg>
            </span>
            <span class="logo-text">Macro Recorder</span>
        </div>

        <div class="header-status">
            <!-- 录制中状态 -->
            <div v-if="isRecording" class="status-badge recording">
                <span class="status-dot"></span>
                <span class="status-text">
                    录制中
                    <span class="status-duration">{{ formattedDuration }}</span>
                    <span class="status-events">| {{ eventCount }} 事件</span>
                </span>
            </div>

            <!-- 播放中状态 -->
            <div v-else-if="isPlaying" class="status-badge playing">
                <span class="status-dot"></span>
                <span class="status-text">播放中...</span>
            </div>

            <!-- 引擎连接状态 -->
            <div v-else class="status-badge" :class="{ active: connected }">
                <span class="status-dot"></span>
                <span class="status-text">{{ connected ? '引擎已连接' : '未连接' }}</span>
            </div>
        </div>
    </header>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
    connected: {
        type: Boolean,
        default: false
    },
    isRecording: {
        type: Boolean,
        default: false
    },
    isPlaying: {
        type: Boolean,
        default: false
    },
    recordingDuration: {
        type: Number,
        default: 0
    },
    eventCount: {
        type: Number,
        default: 0
    }
});

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
.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: var(--header-height);
    padding: 0 var(--space-4);
    background: var(--bg-surface);
    border-bottom: 1px solid var(--border-subtle);
    flex-shrink: 0;
}

.logo {
    display: flex;
    align-items: center;
    gap: var(--space-3);
}

.logo-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    background: var(--color-primary);
    border-radius: var(--radius-md);
    color: var(--text-primary);
}

.logo-text {
    font-size: var(--text-lg);
    font-weight: 600;
    color: var(--text-primary);
    letter-spacing: -0.02em;
}

.header-status {
    display: flex;
    align-items: center;
    gap: var(--space-3);
}

/* 状态徽章基础样式 */
.status-badge {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    padding: var(--space-2) var(--space-4);
    background: var(--bg-element);
    border-radius: var(--radius-full);
    font-size: var(--text-sm);
    transition: all var(--transition-normal);
}

/* 状态点 */
.status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--text-muted);
    flex-shrink: 0;
    transition: all var(--transition-normal);
}

/* 引擎已连接状态 */
.status-badge.active {
    color: var(--color-success);
    background: var(--color-playing-bg);
}

.status-badge.active .status-dot {
    background: var(--color-success);
    box-shadow: 0 0 8px var(--color-playing-glow);
    animation: status-pulse 2s ease-in-out infinite;
}

/* 录制状态 */
.status-badge.recording {
    color: var(--color-recording);
    background: var(--color-recording-bg);
    animation: fade-in var(--transition-normal);
}

.status-badge.recording .status-dot {
    background: var(--color-recording);
    box-shadow: 0 0 8px var(--color-recording-glow);
    animation: status-pulse 1s ease-in-out infinite;
}

/* 播放状态 */
.status-badge.playing {
    color: var(--color-playing);
    background: var(--color-playing-bg);
}

.status-badge.playing .status-dot {
    background: var(--color-playing);
    box-shadow: 0 0 8px var(--color-playing-glow);
    animation: status-pulse 1.5s ease-in-out infinite;
}

/* 状态文字 */
.status-text {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    font-weight: 500;
}

.status-duration {
    font-family: var(--font-mono);
    font-variant-numeric: tabular-nums;
}

.status-events {
    font-family: var(--font-mono);
    opacity: 0.8;
}

/* 动画关键帧 */
@keyframes status-pulse {
    0% {
        box-shadow: 0 0 0 0 currentColor;
        opacity: 1;
    }
    70% {
        box-shadow: 0 0 0 8px transparent;
        opacity: 0.6;
    }
    100% {
        box-shadow: 0 0 0 0 transparent;
        opacity: 1;
    }
}

@keyframes fade-in {
    0% {
        opacity: 0;
        transform: scale(0.95);
    }
    100% {
        opacity: 1;
        transform: scale(1);
    }
}
</style>
