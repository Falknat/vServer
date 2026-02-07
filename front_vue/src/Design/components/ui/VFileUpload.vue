<script setup>
const props = defineProps({
  label: { type: String, default: '' },
  accept: { type: String, default: '' },
  required: { type: Boolean, default: false },
  hint: { type: String, default: '' },
})

const emit = defineEmits(['select'])
const fileName = ref('')

const onFileSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    fileName.value = file.name
    emit('select', file)
  }
}
</script>

<template>
  <div class="v-file-group">
    <label v-if="label" class="v-file-label">
      {{ label }} <span v-if="required" class="v-file-required">*</span>
    </label>
    <div class="v-file-wrapper">
      <input type="file" class="v-file-input" :accept="accept" @change="onFileSelect">
      <div class="v-file-btn">
        <i class="fas fa-file-upload"></i>
        <span>{{ fileName || $t('certs.selectFile') }}</span>
      </div>
    </div>
    <small v-if="hint" class="v-file-hint">
      <i class="fas fa-info-circle"></i> {{ hint }}
    </small>
  </div>
</template>

<style scoped>
.v-file-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.v-file-label {
  font-size: var(--text-md);
  color: var(--text-secondary);
  font-weight: var(--font-medium);
}

.v-file-required { color: var(--accent-red); }

.v-file-wrapper {
  position: relative;
}

.v-file-input {
  position: absolute;
  inset: 0;
  opacity: 0;
  cursor: pointer;
}

.v-file-btn {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-sm) var(--space-md);
  background: var(--glass-bg);
  border: 1px dashed var(--glass-border);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  transition: all var(--transition-fast);
}

.v-file-wrapper:hover .v-file-btn {
  border-color: var(--accent-purple);
  color: var(--text-primary);
}

.v-file-hint {
  font-size: var(--text-sm);
  color: var(--text-muted);
}
</style>
