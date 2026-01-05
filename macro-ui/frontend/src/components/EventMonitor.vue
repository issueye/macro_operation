<template>
    <div class="card event-monitor">
        <div class="card-header">
            <div class="header-left">
                <h3>事件监控</h3>
                <span class="event-count" v-if="displayedEvents.length > 0">
                    {{ displayedEvents.length }}
                </span>
            </div>
            <div class="header-right">
                <!-- 事件类型筛选 -->
                <div class="filter-tabs">
                    <button
                        v-for="filter in filters"
                        :key="filter.value"
                        class="filter-tab"
                        :class="{ active: activeFilter === filter.value }"
                        @click="activeFilter = filter.value"
                    >
                        <span class="filter-icon" v-html="filter.icon"></span>
                        <span class="filter-label">{{ filter.label }}</span>
                    </button>
                </div>
                <button class="btn-icon" @click="clearEvents" title="清空">
                    <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6" />
                        <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2" />
                    </svg>
                </button>
            </div>
        </div>

        <!-- 事件列表 -->
        <div class="event-list" ref="eventListRef">
            <div
                v-for="event in displayedEvents"
                :key="event.index"
                class="event-item"
                :class="getEventClass(event)"
            >
                <!-- 事件图标 -->
                <span class="event-icon" v-html="getEventIcon(event)"></span>

                <!-- 事件信息 -->
                <span class="event-info">
                    <span class="event-type">{{ getEventTypeName(event) }}</span>
                    <span class="event-detail">{{ getEventDetail(event) }}</span>
                </span>

                <!-- 事件索引 -->
                <span class="event-index">#{{ event.index }}</span>
            </div>

            <!-- 空状态 -->
            <div v-if="displayedEvents.length === 0" class="empty-state">
                <div class="empty-icon">
                    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="1.5">
                        <circle cx="12" cy="12" r="10" />
                        <line x1="12" y1="8" x2="12" y2="12" />
                        <line x1="12" y1="16" x2="12.01" y2="16" />
                    </svg>
                </div>
                <p class="empty-text">
                    {{ isRecording ? '等待事件输入...' : '开始录制后显示事件' }}
                </p>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick } from 'vue';

const props = defineProps({
    isRecording: {
        type: Boolean,
        default: false
    }
});

const events = ref([]);
const pollTimer = ref(null);
const eventListRef = ref(null);
const activeFilter = ref('all');

// 事件类型筛选
const filters = [
    { value: 'all', label: '全部', icon: '<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/></svg>' },
    { value: 'keyboard', label: '键盘', icon: '<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="4" width="20" height="16" rx="2"/><path d="M6 8h.01M10 8h.01M14 8h.01M18 8h.01M8 12h8M6 16h.01M10 16h.01M14 16h.01M18 16h.01M8 20h4"/></svg>' },
    { value: 'mouse', label: '鼠标', icon: '<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2c-4 0-7 3-7 7v6c0 1.1.9 2 2 2h10c1.1 0 2-.9 2-2V9c0-4-3-7-7-7z"/><circle cx="12" cy="18" r="2"/></svg>' },
];

// 筛选后的事件
const displayedEvents = computed(() => {
    let filtered = events.value;

    if (activeFilter.value === 'keyboard') {
        filtered = filtered.filter(e => e.type === 'keydown' || e.type === 'keyup');
    } else if (activeFilter.value === 'mouse') {
        filtered = filtered.filter(e => e.type.startsWith('mouse') || e.type === 'wheel');
    }

    return filtered.slice(-100).reverse();
});

// 获取事件图标
const getEventIcon = (event) => {
    const icons = {
        keydown: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="4" width="20" height="16" rx="2"/><path d="M6 8h.01M10 8h.01M14 8h.01M18 8h.01M8 12h8"/></svg>',
        keyup: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-opacity="0.6"><rect x="2" y="4" width="20" height="16" rx="2"/><path d="M6 8h.01M10 8h.01M14 8h.01M18 8h.01M8 12h8"/></svg>',
        mousemove: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M12 2v4M12 18v4M2 12h4M18 12h4" stroke-dasharray="2 2"/></svg>',
        mousedown: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M12 2v4M12 18v4M2 12h4M18 12h4" stroke-dasharray="2 2"/><path d="M8 8l8 8" stroke-width="2"/></svg>',
        mouseup: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M12 2v4M12 18v4M2 12h4M18 12h4" stroke-dasharray="2 2"/></svg>',
        click: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M12 2v4M12 18v4M2 12h4M18 12h4"/></svg>',
        wheel: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M12 2v4M12 18v4M2 12h4M18 12h4"/><path d="M12 8l4-4M12 8l-4-4M12 8l4 4M12 8l-4 4" stroke-width="1.5"/></svg>',
    };
    return icons[event.type] || icons.mousemove;
};

// 获取事件类型名称
const getEventTypeName = (event) => {
    const names = {
        keydown: 'KeyDown',
        keyup: 'KeyUp',
        mousemove: 'Move',
        mousedown: 'MouseDown',
        mouseup: 'MouseUp',
        click: 'Click',
        wheel: 'Scroll',
    };
    return names[event.type] || event.type;
};

// 获取事件详情
const getEventDetail = (event) => {
    if (event.type === 'keydown' || event.type === 'keyup') {
        if (event.key) {
            return `"${event.key}"`;
        }
        if (event.keyCode > 0) {
            return `KeyCode: ${event.keyCode}`;
        }
        return '';
    }
    if (event.type === 'mousemove' || event.type === 'mousedown' || event.type === 'mouseup' || event.type === 'click') {
        return `(${event.x}, ${event.y})`;
    }
    if (event.type === 'wheel') {
        const direction = event.delta > 0 ? 'Down' : 'Up';
        return `Delta: ${event.delta} ${direction}`;
    }
    return '';
};

// 获取事件样式类
const getEventClass = (event) => {
    return {
        'event-keydown': event.type === 'keydown',
        'event-keyup': event.type === 'keyup',
        'event-mousemove': event.type === 'mousemove',
        'event-click': event.type === 'click' || event.type === 'mousedown' || event.type === 'mouseup',
        'event-wheel': event.type === 'wheel'
    };
};

// 清空事件
const clearEvents = () => {
    events.value = [];
};

// 轮询获取事件
const pollEvents = async () => {
    try {
        const recentEvents = await window.go.main.App.GetRecentEvents(100);
        if (recentEvents && recentEvents.length > 0) {
            events.value = recentEvents;
            // 自动滚动到底部
            nextTick(() => {
                if (eventListRef.value) {
                    eventListRef.value.scrollTop = eventListRef.value.scrollHeight;
                }
            });
        }
    } catch (e) {
        console.error('获取事件失败:', e);
    }
};

// 开始轮询
const startPolling = () => {
    if (pollTimer.value) return;
    pollTimer.value = setInterval(pollEvents, 200);
};

// 停止轮询
const stopPolling = () => {
    if (pollTimer.value) {
        clearInterval(pollTimer.value);
        pollTimer.value = null;
    }
};

onMounted(() => {
    if (props.isRecording) {
        startPolling();
    }
});

onBeforeUnmount(() => {
    stopPolling();
});

// 监听录制状态
watch(() => props.isRecording, (newVal) => {
    if (newVal) {
        startPolling();
    } else {
        stopPolling();
    }
});

// 暴露方法
defineExpose({
    startPolling,
    stopPolling
});
</script>

<style scoped>
.event-monitor {
    display: flex;
    flex-direction: column;
    height: 100%;
    min-height: 0;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: var(--space-2) var(--space-3);
    background: var(--bg-element);
    border-bottom: 1px solid var(--border-subtle);
    flex-shrink: 0;
}

.header-left {
    display: flex;
    align-items: center;
    gap: var(--space-2);
}

.header-left h3 {
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--text-primary);
}

.event-count {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 20px;
    height: 20px;
    padding: 0 6px;
    background: var(--color-primary);
    border-radius: var(--radius-full);
    font-size: 11px;
    font-weight: 600;
    color: var(--text-primary);
}

.header-right {
    display: flex;
    align-items: center;
    gap: var(--space-2);
}

/* 筛选标签 */
.filter-tabs {
    display: flex;
    gap: var(--space-1);
    padding: 2px;
    background: var(--bg-base);
    border-radius: var(--radius-md);
}

.filter-tab {
    display: flex;
    align-items: center;
    gap: var(--space-1);
    padding: var(--space-1) var(--space-2);
    background: transparent;
    border: none;
    border-radius: var(--radius-sm);
    font-size: var(--text-xs);
    color: var(--text-secondary);
    cursor: pointer;
    transition: all var(--transition-fast);
}

.filter-tab:hover {
    color: var(--text-primary);
    background: var(--bg-hover);
}

.filter-tab.active {
    color: var(--text-primary);
    background: var(--color-primary);
}

.filter-icon {
    display: flex;
    align-items: center;
}

/* 事件列表 */
.event-list {
    flex: 1;
    overflow-y: auto;
    padding: var(--space-2);
}

.event-item {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-2) var(--space-3);
    border-radius: var(--radius-sm);
    font-size: var(--text-sm);
    font-family: var(--font-mono);
    transition: all var(--transition-fast);
    animation: slide-in var(--transition-normal);
}

@keyframes slide-in {
    from {
        opacity: 0;
        transform: translateX(-10px);
    }
    to {
        opacity: 1;
        transform: translateX(0);
    }
}

.event-item:hover {
    background: var(--bg-element);
}

/* 事件图标 */
.event-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    flex-shrink: 0;
}

/* 事件信息 */
.event-info {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    flex: 1;
    min-width: 0;
}

.event-type {
    font-weight: 600;
    min-width: 80px;
}

.event-detail {
    color: var(--text-secondary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

/* 事件索引 */
.event-index {
    font-size: var(--text-xs);
    color: var(--text-muted);
    flex-shrink: 0;
}

/* 键盘事件 */
.event-keydown .event-icon {
    color: var(--color-event-keydown);
}

.event-keydown .event-type {
    color: var(--color-event-keydown);
}

.event-keyup .event-icon {
    color: var(--color-event-keyup);
}

.event-keyup .event-type {
    color: var(--color-event-keyup);
}

/* 鼠标移动事件 */
.event-mousemove .event-icon {
    color: var(--color-event-mouse);
}

.event-mousemove .event-type {
    color: var(--color-event-mouse);
}

/* 点击事件 */
.event-click .event-icon {
    color: var(--color-event-click);
}

.event-click .event-type {
    color: var(--color-event-click);
}

/* 滚轮事件 */
.event-wheel .event-icon {
    color: var(--color-event-scroll);
}

.event-wheel .event-type {
    color: var(--color-event-scroll);
}

/* 空状态 */
.empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: var(--space-6);
    text-align: center;
}

.empty-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    background: var(--bg-element);
    border-radius: var(--radius-lg);
    color: var(--text-muted);
    margin-bottom: var(--space-3);
}

.empty-text {
    font-size: var(--text-sm);
    color: var(--text-muted);
}

/* 滚动条样式 */
.event-list::-webkit-scrollbar {
    width: 6px;
}

.event-list::-webkit-scrollbar-track {
    background: transparent;
}

.event-list::-webkit-scrollbar-thumb {
    background: var(--bg-element);
    border-radius: var(--radius-full);
}

.event-list::-webkit-scrollbar-thumb:hover {
    background: var(--bg-hover);
}
</style>
