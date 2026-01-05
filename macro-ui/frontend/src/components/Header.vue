<template>
  <header class="header">
    <div class="logo">
      <span class="logo-icon">⚡</span>
      <span class="logo-text">Macro Recorder</span>
    </div>
    <div class="header-status">
      <div class="status-badge" :class="{ active: connected }">
        <span class="status-dot"></span>
        {{ connected ? '引擎已连接' : '未连接' }}
      </div>
    </div>
  </header>
</template>

<script setup>
defineProps({
  connected: {
    type: Boolean,
    default: false
  }
})
</script>
