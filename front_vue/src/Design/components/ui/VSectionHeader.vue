<script setup>
defineProps({
  title: { type: String, required: true },
  addable: { type: Boolean, default: false },
})

defineEmits(['add'])
</script>

<template>
  <div class="v-section-header">
    <span class="v-section-bar"></span>
    <h2 class="v-section-title">{{ title }}</h2>
    <button v-if="addable" class="v-section-add" @click="$emit('add')">
      <svg width="10" height="10" viewBox="0 0 10 10">
        <line x1="5" y1="0" x2="5" y2="10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
        <line x1="0" y1="5" x2="10" y2="5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
      </svg>
    </button>
  </div>
</template>

<style scoped>
.v-section-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: var(--space-md);
}

.v-section-bar {
  width: 3px;
  height: 14px;
  background: var(--accent-purple);
  border-radius: 2px;
  flex-shrink: 0;
}

.v-section-title {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--text-muted);
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 1.5px;
}

.v-section-add {
  width: 20px;
  height: 20px;
  display: grid;
  place-items: center;
  background: rgba(var(--accent-rgb), 0.1);
  border: 1px solid rgba(var(--accent-rgb), 0.2);
  border-radius: var(--radius-full);
  color: var(--accent-purple-light);
  cursor: pointer;
  transition: all var(--transition-base);
  padding: 0;
  margin: 0;
  flex-shrink: 0;
}

.v-section-add:hover {
  background: rgba(var(--accent-rgb), 0.2);
  border-color: rgba(var(--accent-rgb), 0.4);
}

.v-section-add svg {
  display: block;
}
</style>
