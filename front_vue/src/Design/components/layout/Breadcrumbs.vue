<script setup>
const { t } = useI18n()
const router = useRouter()

defineProps({
  items: { type: Array, default: () => [] },
})

const goBack = () => {
  router.back()
}
</script>

<template>
  <div class="breadcrumbs">
    <div class="breadcrumbs-left">
      <button class="breadcrumb-item" @click="goBack">
        <i class="fas fa-arrow-left"></i> {{ t('common.back') }}
      </button>
      <template v-for="(item, index) in items" :key="index">
        <span class="breadcrumb-separator">/</span>
        <span class="breadcrumb-item" :class="{ active: index === items.length - 1 }">
          {{ item }}
        </span>
      </template>
    </div>
    <div v-if="$slots.tabs" class="breadcrumbs-tabs">
      <slot name="tabs" />
    </div>
  </div>
</template>

<style scoped>
.breadcrumbs {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-lg);
  margin-bottom: var(--space-md);
  padding: var(--space-md) 20px;
  background: rgba(var(--accent-rgb), 0.05);
  border-radius: var(--radius);
  border: 1px solid var(--glass-border);
}

.breadcrumbs-left {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.breadcrumbs-tabs {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.breadcrumb-item {
  font-size: var(--text-md);
  color: var(--text-muted);
  background: none;
  border: none;
  padding: var(--space-sm) var(--space-lg);
  border-radius: var(--radius);
  cursor: pointer;
  transition: all var(--transition-base);
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

button.breadcrumb-item:hover {
  background: rgba(var(--accent-rgb), 0.1);
  color: var(--accent-purple-light);
}

.breadcrumb-item.active {
  color: var(--text-primary);
  font-weight: var(--font-medium);
  cursor: default;
}

.breadcrumb-separator {
  color: var(--text-muted);
  opacity: 0.3;
}
</style>
