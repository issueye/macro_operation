<template>
  <div class="card log-card">
    <div class="card-header">
      <h3>运行日志</h3>
      <button class="btn-icon" @click="emit('clear')" title="清空日志">
        <svg viewBox="0 0 24 24" width="16" height="16"><path fill="currentColor" d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/></svg>
      </button>
    </div>
    <div class="log-panel" ref="logPanel">
      <div v-for="log in logs" :key="log.id" class="log-entry" :class="log.type">
        <span class="log-time">{{ log.time }}</span>
        <span class="log-message">{{ log.message }}</span>
      </div>
      <div v-if="logs.length === 0" class="empty-state small">
        暂无日志
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'

const props = defineProps({
  logs: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['clear'])

const logPanel = ref(null)

watch(() => props.logs, () => {
  nextTick(() => {
    scrollToBottom()
  })
})

const scrollToBottom = () => {
  if (logPanel.value) {
    logPanel.value.scrollTop = logPanel.value.scrollHeight
  }
}
</script>
