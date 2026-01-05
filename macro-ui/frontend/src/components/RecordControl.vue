<template>
  <div class="card record-card">
    <div class="card-header">
      <h3>录制控制</h3>
    </div>
    <div class="record-controls">
      <button
        class="btn btn-record"
        :class="{ recording: isRecording }"
        @click="emit('toggle-recording')"
        :disabled="isPlaying"
      >
        <span class="record-icon"></span>
        {{ isRecording ? '停止录制' : '开始录制' }}
      </button>
      <div class="record-info" v-if="isRecording">
        <span class="pulse"></span>
        正在录制... {{ eventCount }} 个事件
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
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
  }
})

const emit = defineEmits(['toggle-recording'])
</script>
