<script setup>
defineProps({
  variant: { type: String, default: 'default' },
  icon: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
})

defineEmits(['click'])
</script>

<template>
  <button
    class="v-btn"
    :class="[`v-btn--${variant}`, { 'v-btn--loading': loading }]"
    :disabled="disabled || loading"
    @click="$emit('click', $event)"
  >
    <i v-if="loading" class="fas fa-spinner icon-spin"></i>
    <i v-else-if="icon" :class="icon"></i>
    <span v-if="$slots.default"><slot /></span>
  </button>
</template>

<style scoped>
.v-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-sm) var(--space-md);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-md);
  background: var(--glass-bg);
  color: var(--text-primary);
  cursor: pointer;
  transition: all var(--transition-base);
  font-size: var(--text-md);
  font-weight: var(--font-medium);
  white-space: nowrap;
}

.v-btn:hover:not(:disabled) {
  border-color: var(--glass-border-hover);
  background: var(--glass-bg-light);
}

.v-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.v-btn--primary {
  background: var(--btn-primary-bg);
  border-color: var(--btn-primary-border);
  color: var(--btn-primary-color);
}

.v-btn--primary:hover:not(:disabled) {
  background: var(--btn-primary-hover-bg);
  border-color: var(--btn-primary-hover-border);
}

.v-btn--danger {
  border-color: var(--accent-red);
  color: var(--accent-red);
}

.v-btn--danger:hover:not(:disabled) {
  background: rgba(var(--danger-rgb), 0.15);
  border-color: rgba(var(--danger-rgb), 0.5);
}

.v-btn--success {
  border-color: var(--accent-green);
  color: var(--accent-green);
}

.v-btn--success:hover:not(:disabled) {
  background: rgba(var(--success-rgb), 0.15);
  border-color: rgba(var(--success-rgb), 0.5);
}

.v-btn--icon {
  padding: var(--space-xs);
  width: 32px;
  height: 32px;
  justify-content: center;
  border: none;
  background: transparent;
}

.v-btn--icon:hover:not(:disabled) {
  background: var(--btn-icon-bg);
  color: var(--btn-icon-color);
}
</style>
