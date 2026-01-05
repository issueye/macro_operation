<template>
    <div class="card event-monitor">
        <div class="card-header">
            <h3>事件监控</h3>
            <div class="event-stats">
                <span class="stat">总计: {{ totalEvents }}</span>
                <span class="stat keydown">KeyDown: {{ keyDownCount }}</span>
                <span class="stat keyup">KeyUp: {{ keyUpCount }}</span>
            </div>
        </div>
        <div class="event-list">
            <div
                v-for="event in displayedEvents"
                :key="event.index"
                class="event-item"
                :class="getEventClass(event)"
            >
                <span class="event-index">#{{ event.index }}</span>
                <span class="event-type">{{ event.type }}</span>
                <span class="event-detail">{{ getEventDetail(event) }}</span>
            </div>
            <div v-if="displayedEvents.length === 0" class="empty-state">
                暂无事件数据
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
    isRecording: {
        type: Boolean,
        default: false
    }
})

const events = ref([])
const pollTimer = ref(null)

// 统计数据
const totalEvents = computed(() => events.value.length)
const keyDownCount = computed(() => events.value.filter(e => e.type === 'keydown').length)
const keyUpCount = computed(() => events.value.filter(e => e.type === 'keyup').length)

// 显示最近的事件（最多 50 条）
const displayedEvents = computed(() => {
    return events.value.slice(-50).reverse()
})

// 获取事件详情
const getEventDetail = (event) => {
    if (event.type === 'keydown' || event.type === 'keyup') {
        if (event.keyCode > 0) {
            return `KeyCode: ${event.keyCode}`
        }
        return ''
    }
    if (event.type === 'mousemove' || event.type === 'mousedown' || event.type === 'mouseup') {
        return `(${event.x}, ${event.y})`
    }
    if (event.type === 'wheel') {
        return `Delta: ${event.delta}`
    }
    return ''
}

// 获取事件样式类
const getEventClass = (event) => {
    return {
        'event-keydown': event.type === 'keydown',
        'event-keyup': event.type === 'keyup',
        'event-mouse': event.type.startsWith('mouse'),
        'event-wheel': event.type === 'wheel'
    }
}

// 轮询获取事件
const pollEvents = async () => {
    try {
        const recentEvents = await window.go.main.App.GetRecentEvents(50)
        if (recentEvents && recentEvents.length > 0) {
            events.value = recentEvents
        }
    } catch (e) {
        console.error('获取事件失败:', e)
    }
}

// 开始轮询
const startPolling = () => {
    if (pollTimer.value) return
    pollTimer.value = setInterval(pollEvents, 200) // 每 200ms 轮询一次
}

// 停止轮询
const stopPolling = () => {
    if (pollTimer.value) {
        clearInterval(pollTimer.value)
        pollTimer.value = null
    }
}

// 监听录制状态
const startWatchingRecording = () => {
    if (props.isRecording) {
        startPolling()
    } else {
        stopPolling()
    }
}

onMounted(() => {
    startPolling()
})

onBeforeUnmount(() => {
    stopPolling()
})

// 暴露方法供父组件调用
defineExpose({
    startPolling,
    stopPolling
})
</script>

<style scoped>
.event-monitor {
    height: 300px;
    display: flex;
    flex-direction: column;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    background: rgba(255, 255, 255, 0.05);
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.event-stats {
    display: flex;
    gap: 16px;
    font-size: 12px;
}

.stat {
    color: #858585;
}

.stat.keydown {
    color: #4ec9b0;
}

.stat.keyup {
    color: #ce9178;
}

.event-list {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
}

.event-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 6px 8px;
    border-radius: 4px;
    font-size: 12px;
    font-family: 'Consolas', 'Monaco', monospace;
    transition: background 0.2s;
}

.event-item:hover {
    background: rgba(255, 255, 255, 0.05);
}

.event-index {
    color: #858585;
    min-width: 40px;
}

.event-type {
    font-weight: bold;
    min-width: 80px;
}

.event-item.event-keydown .event-type {
    color: #4ec9b0;
}

.event-item.event-keyup .event-type {
    color: #ce9178;
}

.event-item.event-mouse .event-type {
    color: #569cd6;
}

.event-item.event-wheel .event-type {
    color: #dcdcaa;
}

.event-detail {
    color: #d4d4d4;
    flex: 1;
}

.empty-state {
    text-align: center;
    padding: 40px 20px;
    color: #858585;
}

/* 滚动条样式 */
.event-list::-webkit-scrollbar {
    width: 8px;
}

.event-list::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.05);
}

.event-list::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.2);
    border-radius: 4px;
}

.event-list::-webkit-scrollbar-thumb:hover {
    background: rgba(255, 255, 255, 0.3);
}
</style>
