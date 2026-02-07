<script setup>
const props = defineProps({
  modelValue: { type: Array, default: () => [] },
  placeholder: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue'])

const input = ref('')
const editing = ref(false)

const addItem = () => {
  const raw = input.value.trim()
  if (!raw) return
  const items = raw.split(',').map(s => s.trim()).filter(s => s && !props.modelValue.includes(s))
  if (items.length) {
    emit('update:modelValue', [...props.modelValue, ...items])
  }
  input.value = ''
}

const removeItem = (index) => {
  const updated = [...props.modelValue]
  updated.splice(index, 1)
  emit('update:modelValue', updated)
}

const onKeydown = (e) => {
  if (e.key === 'Enter') {
    e.preventDefault()
    addItem()
  }
}
</script>

<template>
  <div class="field-editor" @click="editing = true">
    <div v-if="!editing && modelValue.length === 0" class="empty-field" @click="editing = true">—</div>
    <div v-else-if="!editing" class="mini-tags" @click="editing = true">
      <code v-for="(item, i) in modelValue" :key="i">{{ item }}</code>
      <span class="edit-hint"><i class="fas fa-pen"></i></span>
    </div>
    <div v-else class="field-editor-active" @click.stop>
      <div class="field-input-row">
        <input
          v-model="input"
          class="field-input"
          :placeholder="placeholder"
          @keydown="onKeydown"
          @blur="input || (editing = false)"
        >
        <button v-if="input" class="field-add-btn" @click="addItem"><i class="fas fa-plus"></i></button>
      </div>
      <div v-if="modelValue.length" class="field-tags">
        <span v-for="(item, i) in modelValue" :key="i" class="field-tag">
          {{ item }}
          <button class="field-tag-remove" @click="removeItem(i)"><i class="fas fa-times"></i></button>
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.field-editor {
  cursor: pointer;
  min-height: 28px;
}

.edit-hint {
  opacity: 0;
  font-size: 10px;
  color: var(--text-muted);
  transition: opacity var(--transition-fast);
}

.field-editor:hover .edit-hint {
  opacity: 0.6;
}

.field-editor-active {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.field-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  background: rgba(var(--accent-rgb), 0.2);
  border: 1px solid rgba(var(--accent-rgb), 0.4);
  border-radius: 12px;
  font-size: 11px;
  color: var(--text-primary);
  font-family: var(--font-mono);
}

.field-tag-remove {
  background: none;
  border: none;
  color: var(--accent-red);
  cursor: pointer;
  padding: 0;
  font-size: 9px;
  opacity: 0.6;
  transition: opacity var(--transition-fast);
}

.field-tag-remove:hover {
  opacity: 1;
}

.field-input-row {
  display: flex;
  gap: 4px;
}

.field-input {
  flex: 1;
  padding: 4px 8px;
  background: var(--glass-bg-dark);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius);
  color: var(--text-primary);
  font-size: 12px;
  font-family: var(--font-mono);
  outline: none;
  min-width: 80px;
}

.field-input:focus {
  border-color: rgba(var(--accent-rgb), 0.5);
}

.field-input::placeholder {
  color: var(--text-muted);
  opacity: 0.4;
}

.field-add-btn {
  padding: 4px 8px;
  background: rgba(var(--accent-rgb), 0.15);
  border: 1px solid rgba(var(--accent-rgb), 0.3);
  border-radius: var(--radius);
  color: var(--accent-purple-light);
  cursor: pointer;
  font-size: 11px;
}

.field-add-btn:hover {
  background: rgba(var(--accent-rgb), 0.25);
}
</style>
