<template>
  <div class="card macro-card">
    <div class="card-header">
      <h3>宏列表</h3>
      <button class="btn-icon" @click="emit('refresh')" title="刷新">
        <svg viewBox="0 0 24 24" width="16" height="16"><path fill="currentColor" d="M17.65 6.35A7.958 7.958 0 0012 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0112 18c-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"/></svg>
      </button>
    </div>
    <div class="macro-list">
      <div
        v-for="macro in macros"
        :key="macro.name"
        class="macro-item"
        :class="{ active: selectedMacro === macro.name }"
        @click="emit('select', macro)"
      >
        <div class="macro-info">
          <span class="macro-name">{{ macro.name }}</span>
          <span class="macro-events">{{ macro.event_count }} 事件</span>
        </div>
        <div class="macro-actions">
          <button class="btn-icon" @click.stop="emit('play', macro.name)" title="播放">
            <svg viewBox="0 0 24 24" width="14" height="14"><path fill="currentColor" d="M8 5v14l11-7z"/></svg>
          </button>
          <button class="btn-icon danger" @click.stop="emit('delete', macro.name)" title="删除">
            <svg viewBox="0 0 24 24" width="14" height="14"><path fill="currentColor" d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"/></svg>
          </button>
        </div>
      </div>
      <div v-if="macros.length === 0" class="empty-state">
        暂无保存的宏
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  macros: {
    type: Array,
    default: () => []
  },
  selectedMacro: {
    type: String,
    default: null
  }
})

const emit = defineEmits(['refresh', 'select', 'play', 'delete'])
</script>
