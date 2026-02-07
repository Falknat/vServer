<script setup>
defineProps({
  variant: { type: String, default: 'default' },
})
</script>

<template>
  <span class="v-badge" :class="`v-badge--${variant}`">
    <slot />
  </span>
</template>

<style scoped>
.v-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: var(--text-xs);
  font-weight: var(--font-bold);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border: 1px solid transparent;
}

.v-badge--online { background: linear-gradient(135deg, rgba(var(--success-rgb), 0.25), rgba(var(--success-rgb), 0.15)); color: var(--accent-green); border-color: rgba(var(--success-rgb), 0.4); }
.v-badge--offline { background: linear-gradient(135deg, rgba(var(--danger-rgb), 0.25), rgba(var(--danger-rgb), 0.15)); color: var(--accent-red); border-color: rgba(var(--danger-rgb), 0.4); }
.v-badge--pending { background: linear-gradient(135deg, rgba(var(--warning-rgb), 0.25), rgba(var(--warning-rgb), 0.15)); color: var(--accent-yellow); border-color: rgba(var(--warning-rgb), 0.4); }
.v-badge--yes { background: linear-gradient(135deg, rgba(var(--success-rgb), 0.25), rgba(var(--success-rgb), 0.15)); color: var(--accent-green); border-color: rgba(var(--success-rgb), 0.4); }
.v-badge--no { background: linear-gradient(135deg, rgba(var(--muted-rgb), 0.2), rgba(var(--muted-rgb), 0.1)); color: var(--text-muted); border-color: rgba(var(--muted-rgb), 0.3); }
.v-badge--default { background: var(--glass-bg); color: var(--text-secondary); border-color: var(--glass-border); }
</style>
