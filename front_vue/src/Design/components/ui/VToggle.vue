<script setup>
const props = defineProps({
  modelValue: { type: Boolean, default: false },
  label: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue'])

const toggle = () => {
  emit('update:modelValue', !props.modelValue)
}
</script>

<template>
  <div class="v-toggle-wrapper" @click="toggle">
    <div class="v-toggle">
      <input type="checkbox" :checked="modelValue">
      <span class="v-toggle-slider"></span>
    </div>
    <span v-if="label" class="v-toggle-label">{{ label }}</span>
  </div>
</template>

<style scoped>
.v-toggle-wrapper {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  cursor: pointer;
}

.v-toggle {
  position: relative;
  width: 40px;
  height: 22px;
  flex-shrink: 0;
}

.v-toggle input {
  opacity: 0;
  width: 0;
  height: 0;
  position: absolute;
}

.v-toggle-slider {
  position: absolute;
  inset: 0;
  background: var(--glass-bg);
  border: 1px solid var(--glass-border);
  border-radius: 11px;
  transition: all var(--transition-base);
}

.v-toggle-slider::before {
  content: '';
  position: absolute;
  width: 16px;
  height: 16px;
  left: 2px;
  top: 2px;
  background: var(--text-muted);
  border-radius: var(--radius-full);
  transition: all var(--transition-base);
}

.v-toggle input:checked + .v-toggle-slider {
  background: rgba(var(--accent-rgb), 0.3);
  border-color: var(--accent-purple);
}

.v-toggle input:checked + .v-toggle-slider::before {
  transform: translateX(18px);
  background: var(--accent-purple-light);
}

.v-toggle-label {
  font-size: var(--text-md);
  color: var(--text-secondary);
}
</style>
