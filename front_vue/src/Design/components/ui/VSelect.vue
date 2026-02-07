<script setup>
const props = defineProps({
  modelValue: { type: [String, Number], default: '' },
  label: { type: String, default: '' },
  options: { type: Array, default: () => [] },
  required: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue'])

const isOpen = ref(false)
const selectRef = ref(null)

const selectedLabel = computed(() => {
  const opt = props.options.find(o => o.value === props.modelValue)
  return opt ? opt.label : ''
})

const select = (value) => {
  emit('update:modelValue', value)
  isOpen.value = false
}

const toggle = () => {
  isOpen.value = !isOpen.value
}

const onDocClick = (e) => {
  if (selectRef.value && !selectRef.value.contains(e.target)) {
    isOpen.value = false
  }
}

onMounted(() => document.addEventListener('click', onDocClick))
onUnmounted(() => document.removeEventListener('click', onDocClick))
</script>

<template>
  <div class="v-select-group">
    <label v-if="label" class="v-select-label">
      {{ label }} <span v-if="required" class="v-select-required">*</span>
    </label>
    <div ref="selectRef" class="v-select" :class="{ open: isOpen }">
      <div class="v-select-trigger" @click="toggle">
        <span class="v-select-value">{{ selectedLabel }}</span>
        <i class="fas fa-chevron-down v-select-arrow"></i>
      </div>
      <transition name="dropdown">
        <div v-if="isOpen" class="v-select-dropdown">
          <div
            v-for="opt in options"
            :key="opt.value"
            class="v-select-option"
            :class="{ selected: opt.value === modelValue }"
            @click="select(opt.value)"
          >
            {{ opt.label }}
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<style scoped>
.v-select-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.v-select-label {
  font-size: 12px;
  font-weight: var(--font-semibold);
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.v-select-required { color: var(--accent-red); }

.v-select {
  position: relative;
  width: 100%;
}

.v-select-trigger {
  padding: 10px 14px;
  background: var(--glass-bg-dark);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-md);
  color: var(--text-primary);
  font-size: var(--text-base);
  cursor: pointer;
  transition: all var(--transition-base);
  display: flex;
  align-items: center;
  justify-content: space-between;
  user-select: none;
}

.v-select-trigger:hover {
  border-color: rgba(var(--accent-rgb), 0.4);
}

.v-select.open .v-select-trigger {
  border-color: rgba(var(--accent-rgb), 0.6);
  box-shadow: var(--input-focus-shadow);
}

.v-select-value {
  flex: 1;
}

.v-select-arrow {
  color: var(--text-muted);
  font-size: 12px;
  transition: transform var(--transition-base);
}

.v-select.open .v-select-arrow {
  transform: rotate(180deg);
}

.v-select-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  background: var(--bg-secondary);
  border: 1px solid rgba(var(--accent-rgb), 0.3);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  max-height: 300px;
  overflow-y: auto;
  z-index: 1000;
}

.v-select-option {
  padding: 10px 14px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all var(--transition-base);
  font-size: var(--text-base);
}

.v-select-option:hover {
  background: rgba(var(--accent-rgb), 0.1);
  color: var(--text-primary);
}

.v-select-option.selected {
  background: rgba(var(--accent-rgb), 0.2);
  color: var(--accent-purple-light);
  font-weight: var(--font-semibold);
}

.v-select-option.selected::before {
  content: '✓ ';
  margin-right: 8px;
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: all var(--transition-fast);
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
