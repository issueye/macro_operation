<template>
    <div class="macro-list-container" :class="{ collapsed }">
        <!-- 标题栏 -->
        <div class="macro-header" v-if="!collapsed">
            <span class="macro-title">宏列表</span>
            <span class="macro-count">({{ macros.length }})</span>
            <div class="macro-actions">
                <button class="btn-icon" @click="emit('refresh')" title="刷新">
                    <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M23 4v6h-6" />
                        <path d="M1 20v-6h6" />
                        <path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15" />
                    </svg>
                </button>
                <button class="btn-icon" @click="handleNewMacro" title="新建宏">
                    <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                        <line x1="12" y1="5" x2="12" y2="19" />
                        <line x1="5" y1="12" x2="19" y2="12" />
                    </svg>
                </button>
            </div>
        </div>

        <!-- 宏列表 -->
        <div class="macro-list">
            <template v-if="macros.length > 0">
                <div
                    v-for="macro in macros"
                    :key="macro.name"
                    class="macro-item"
                    :class="{ active: selectedMacro === macro.name }"
                    @click="handleSelect(macro)"
                >
                    <!-- 折叠模式只显示图标 -->
                    <template v-if="collapsed">
                        <div class="collapsed-icon" :class="{ active: selectedMacro === macro.name }">
                            <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2">
                                <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
                                <polyline points="14 2 14 8 20 8" />
                            </svg>
                        </div>
                    </template>

                    <!-- 展开模式显示完整信息 -->
                    <template v-else>
                        <div class="macro-icon">
                            <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2">
                                <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
                                <polyline points="14 2 14 8 20 8" />
                                <line x1="16" y1="13" x2="8" y2="13" />
                                <line x1="16" y1="17" x2="8" y2="17" />
                                <polyline points="10 9 9 9 8 9" />
                            </svg>
                        </div>
                        <div class="macro-content">
                            <div class="macro-name">{{ macro.name }}</div>
                            <div class="macro-meta">
                                <span class="meta-item">
                                    <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2">
                                        <circle cx="12" cy="12" r="10" />
                                        <polyline points="12 6 12 12 16 14" />
                                    </svg>
                                    {{ macro.event_count }} 事件
                                </span>
                                <span class="meta-item">
                                    {{ formatTime(macro.updated_at) }}
                                </span>
                            </div>
                        </div>
                        <div class="macro-actions">
                            <button
                                class="btn-icon btn-icon-sm"
                                @click.stop="emit('play', macro.name)"
                                title="播放"
                            >
                                <svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor">
                                    <polygon points="5 3 19 12 5 21 5 3" />
                                </svg>
                            </button>
                            <button
                                class="btn-icon btn-icon-sm danger"
                                @click.stop="emit('delete', macro.name)"
                                title="删除"
                            >
                                <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2">
                                    <polyline points="3 6 5 6 21 6" />
                                    <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2" />
                                </svg>
                            </button>
                        </div>
                    </template>
                </div>
            </template>

            <!-- 空状态 -->
            <div v-else class="empty-state">
                <div class="empty-icon">
                    <svg viewBox="0 0 24 24" width="32" height="32" fill="none" stroke="currentColor" stroke-width="1.5">
                        <path d="M13 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V9z" />
                        <polyline points="13 2 13 9 20 9" />
                    </svg>
                </div>
                <p class="empty-text" v-if="!collapsed">暂无保存的宏</p>
                <p class="empty-text" v-else>-</p>
            </div>
        </div>
    </div>
</template>

<script setup>
import { formatDistanceToNow } from 'date-fns';

const props = defineProps({
    macros: {
        type: Array,
        default: () => []
    },
    selectedMacro: {
        type: String,
        default: null
    },
    collapsed: {
        type: Boolean,
        default: false
    }
});

const emit = defineEmits(['refresh', 'select', 'play', 'delete']);

// 选择宏
const handleSelect = (macro) => {
    emit('select', macro);
};

// 新建宏（选择当前空脚本）
const handleNewMacro = () => {
    // 触发一个特殊事件，告诉父组件创建新宏
    emit('select', { name: null, isNew: true });
};

// 格式化时间
const formatTime = (timestamp) => {
    if (!timestamp) return '未知';

    try {
        const date = new Date(timestamp * 1000);
        const now = new Date();
        const diff = date.getTime() - now.getTime();
        const absDiff = Math.abs(diff);

        // 小于1小时显示分钟
        if (absDiff < 3600000) {
            const mins = Math.floor(absDiff / 60000);
            if (mins < 1) return '刚刚';
            return diff > 0 ? `${mins}后` : `${mins}分钟前`;
        }

        // 小于24小时显示小时
        if (absDiff < 86400000) {
            const hours = Math.floor(absDiff / 3600000);
            return diff > 0 ? `${hours}小时后` : `${hours}小时前`;
        }

        // 小于7天显示天
        if (absDiff < 604800000) {
            const days = Math.floor(absDiff / 86400000);
            return diff > 0 ? `${days}天后` : `${days}天前`;
        }

        // 否则显示日期
        return date.toLocaleDateString('zh-CN', {
            month: 'short',
            day: 'numeric'
        });
    } catch (e) {
        return '未知';
    }
};
</script>

<style scoped>
.macro-list-container {
    display: flex;
    flex-direction: column;
    flex: 1;
    min-height: 0;
}

.macro-list-container.collapsed {
    align-items: center;
}

/* 标题栏 */
.macro-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--space-2) var(--space-2);
    margin-bottom: var(--space-2);
}

.macro-title {
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--text-primary);
}

.macro-count {
    font-size: var(--text-xs);
    color: var(--text-muted);
    margin-left: var(--space-1);
}

.macro-actions {
    display: flex;
    gap: var(--space-1);
}

/* 宏列表 */
.macro-list {
    flex: 1;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
}

/* 宏项 */
.macro-item {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-3);
    background: var(--bg-base);
    border: 1px solid transparent;
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: all var(--transition-fast);
    position: relative;
}

.macro-item::before {
    content: '';
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 3px;
    height: 0;
    background: var(--color-primary);
    border-radius: 0 2px 2px 0;
    transition: height var(--transition-fast);
}

.macro-item:hover {
    background: var(--bg-element);
}

.macro-item:hover::before {
    height: 60%;
}

.macro-item.active {
    background: var(--bg-element);
    border-color: var(--border-active);
}

.macro-item.active::before {
    height: 80%;
}

/* 宏图标 */
.macro-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    background: var(--color-primary);
    background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-hover) 100%);
    border-radius: var(--radius-md);
    color: var(--text-primary);
    flex-shrink: 0;
}

.macro-item:hover .macro-icon {
    transform: scale(1.05);
}

.macro-item.active .macro-icon {
    box-shadow: 0 0 12px var(--color-primary-glow);
}

/* 折叠模式图标 */
.collapsed-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    background: var(--bg-element);
    border-radius: var(--radius-md);
    color: var(--text-secondary);
    transition: all var(--transition-fast);
}

.collapsed-icon.active {
    background: var(--color-primary);
    color: var(--text-primary);
}

/* 宏内容 */
.macro-content {
    flex: 1;
    min-width: 0;
}

.macro-name {
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.macro-item:hover .macro-name {
    color: var(--text-primary);
}

.macro-item.active .macro-name {
    color: var(--color-primary);
}

/* 宏元信息 */
.macro-meta {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    margin-top: var(--space-1);
}

.meta-item {
    display: flex;
    align-items: center;
    gap: var(--space-1);
    font-size: var(--text-xs);
    color: var(--text-muted);
}

.meta-item svg {
    opacity: 0.6;
}

/* 宏操作按钮 */
.macro-item .macro-actions {
    opacity: 0;
    transition: opacity var(--transition-fast);
}

.macro-item:hover .macro-actions {
    opacity: 1;
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
    width: 48px;
    height: 48px;
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
.macro-list::-webkit-scrollbar {
    width: 6px;
}

.macro-list::-webkit-scrollbar-track {
    background: transparent;
}

.macro-list::-webkit-scrollbar-thumb {
    background: var(--bg-element);
    border-radius: var(--radius-full);
}

.macro-list::-webkit-scrollbar-thumb:hover {
    background: var(--bg-hover);
}
</style>
