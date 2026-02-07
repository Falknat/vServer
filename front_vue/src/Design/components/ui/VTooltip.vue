<script setup>
defineProps({
  text: { type: String, default: '' },
  items: { type: Array, default: () => [] },
})

const show = ref(false)
</script>

<template>
  <div class="v-tooltip-wrap" @mouseenter="show = true" @mouseleave="show = false">
    <span class="v-tooltip-trigger">
      <i class="fas fa-info-circle"></i>
    </span>
    <transition name="tooltip">
      <div v-if="show" class="v-tooltip-popup">
        <p v-if="text">{{ text }}</p>
        <div v-if="items.length" class="v-tooltip-list">
          <div v-for="(item, i) in items" :key="i" class="v-tooltip-item">{{ item }}</div>
        </div>
      </div>
    </transition>
  </div>
</template>

<style scoped>
.v-tooltip-wrap {
  position: relative;
  display: inline-flex;
  align-items: center;
}

.v-tooltip-trigger {
  color: var(--text-muted);
  font-size: 12px;
  cursor: help;
  transition: color var(--transition-fast);
}

.v-tooltip-trigger:hover {
  color: var(--accent-purple-light);
}

.v-tooltip-popup {
  position: absolute;
  bottom: calc(100% + 8px);
  left: 0;
  min-width: 220px;
  max-width: 360px;
  padding: 10px 14px;
  background: var(--bg-secondary);
  border: 1px solid rgba(var(--accent-rgb), 0.3);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  z-index: 1000;
  font-size: var(--text-sm);
  color: var(--text-secondary);
  line-height: 1.5;
}

.v-tooltip-popup::before {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 8px;
  transform: rotate(45deg);
  width: 8px;
  height: 8px;
  background: var(--bg-secondary);
  border-bottom: 1px solid rgba(var(--accent-rgb), 0.3);
  border-right: 1px solid rgba(var(--accent-rgb), 0.3);
}

.v-tooltip-popup p {
  margin: 0;
}

.v-tooltip-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.v-tooltip-item {
  padding: 4px 0;
  border-bottom: 1px solid rgba(var(--accent-rgb), 0.05);
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  color: var(--text-primary);
}

.v-tooltip-item:last-child {
  border-bottom: none;
}

.tooltip-enter-active,
.tooltip-leave-active {
  transition: all var(--transition-fast);
}

.tooltip-enter-from,
.tooltip-leave-to {
  opacity: 0;
  transform: translateY(4px);
}
</style>
