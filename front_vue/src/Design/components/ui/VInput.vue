<script setup>
defineProps({
  modelValue: { type: [String, Number], default: '' },
  label: { type: String, default: '' },
  placeholder: { type: String, default: '' },
  hint: { type: String, default: '' },
  type: { type: String, default: 'text' },
  required: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
})

defineEmits(['update:modelValue'])
</script>

<template>
  <div class="v-input-group">
    <div v-if="label || hint" class="v-input-header">
      <VTooltip v-if="hint" :text="hint" />
      <label v-if="label" class="v-input-label">
        {{ label }} <span v-if="required" class="v-input-required">*</span>
      </label>
    </div>
    <input
      class="v-input"
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      @input="$emit('update:modelValue', $event.target.value)"
    >
  </div>
</template>

<style scoped>
.v-input-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.v-input-header {
  display: flex;
  align-items: center;
  gap: var(--space-xs);
}

.v-input-label {
  font-size: 12px;
  font-weight: var(--font-semibold);
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  flex-shrink: 0;
}

.v-input-required { color: var(--accent-red); }

.v-input {
  padding: 10px 14px;
  background: var(--glass-bg-dark);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius);
  color: var(--text-primary);
  font-size: var(--text-base);
  transition: all var(--transition-fast);
  outline: none;
}

.v-input::placeholder { color: var(--text-muted); }

.v-input:focus {
  border-color: var(--accent-purple);
  box-shadow: var(--input-focus-shadow);
}

.v-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

</style>
