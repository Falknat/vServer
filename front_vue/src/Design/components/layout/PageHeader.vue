<script setup>
defineProps({
  icon: { type: String, default: '' },
  title: { type: String, default: '' },
  subtitle: { type: String, default: '' },
})
</script>

<template>
  <div class="page-header">
    <div class="page-title-block">
      <h2 class="page-title">
        <i v-if="icon" :class="icon"></i>
        <span>{{ title }}</span>
      </h2>
      <p v-if="subtitle" class="page-subtitle">{{ subtitle }}</p>
    </div>
    <div v-if="$slots.actions" class="page-actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--space-md);
  padding: var(--space-lg);
  background: rgba(var(--accent-rgb), 0.03);
  border-radius: var(--radius-xl);
  border: 1px solid var(--glass-border);
}

.page-title-block {
  flex: 1;
}

.page-title {
  font-size: var(--text-3xl);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin: 0 0 var(--space-sm) 0;
  display: flex;
  align-items: center;
  gap: var(--space-lg);
}

.page-title i {
  color: var(--accent-purple-light);
  font-size: 24px;
}

.page-subtitle {
  font-size: var(--text-md);
  color: var(--text-muted);
  margin: 0;
}

.page-actions {
  display: flex;
  gap: var(--space-sm);
  align-items: center;
}
</style>
