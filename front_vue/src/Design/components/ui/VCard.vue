<script setup>
defineProps({
  title: { type: String, default: '' },
  icon: { type: String, default: '' },
})
</script>

<template>
  <div class="v-card">
    <div v-if="title || $slots.header" class="v-card-header">
      <h3 v-if="title" class="v-card-title">
        <i v-if="icon" :class="icon"></i> {{ title }}
      </h3>
      <slot name="header" />
    </div>
    <div class="v-card-body">
      <slot />
    </div>
    <div v-if="$slots.footer" class="v-card-footer">
      <slot name="footer" />
    </div>
  </div>
</template>

<style scoped>
.v-card {
  background: var(--glass-bg);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius);
  backdrop-filter: var(--backdrop-blur);
  transition: all var(--transition-base);
}

.v-card:hover {
  border-color: var(--glass-border-hover);
}

.v-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-md);
  border-bottom: 1px solid var(--glass-border);
}

.v-card-title {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  font-size: var(--text-md);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
}

.v-card-body {
  padding: var(--space-md);
}

.v-card-footer {
  padding: var(--space-md);
  border-top: 1px solid var(--glass-border);
}
</style>
