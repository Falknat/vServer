<script setup>
const { t } = useI18n()

defineProps({
  service: { type: Object, required: true },
})

const serviceIcons = {
  HTTP: 'fas fa-globe',
  HTTPS: 'fas fa-lock',
  MySQL: 'fas fa-database',
  PHP: 'fab fa-php',
  Proxy: 'fas fa-exchange-alt',
}
</script>

<template>
  <div class="service-item" :class="{ 'service-active': service.status, 'service-pending': service.pending }">
    <span class="service-dot"></span>
    <i :class="serviceIcons[service.name] || 'fas fa-server'" class="service-icon"></i>
    <span class="service-label">{{ service.name }}</span>
    <span v-if="service.pending" class="service-status-text pending">{{ service.pending }}</span>
  </div>
</template>

<style scoped>
.service-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 18px;
  border-radius: var(--radius);
  border: 1px solid var(--glass-border);
  background: var(--glass-bg-light);
  transition: all var(--transition-base);
  cursor: default;
}

.service-item:hover {
  border-color: var(--card-border-hover);
}

.service-dot {
  width: 8px;
  height: 8px;
  border-radius: var(--radius-full);
  background: var(--accent-red);
  flex-shrink: 0;
  transition: all var(--transition-base);
}

.service-active .service-dot {
  background: var(--accent-green);
  box-shadow: 0 0 8px var(--accent-green);
}

.service-pending .service-dot {
  background: var(--accent-yellow, #f0ad4e);
  box-shadow: 0 0 8px var(--accent-yellow, #f0ad4e);
  animation: dot-blink 1s ease-in-out infinite;
}

.service-icon {
  font-size: var(--text-md);
  color: var(--text-muted);
  width: 18px;
  text-align: center;
  transition: color var(--transition-base);
}

.service-active .service-icon {
  color: var(--accent-purple-light);
}

.service-label {
  font-size: var(--text-md);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
}

.service-status-text.pending {
  margin-left: auto;
  font-size: var(--text-sm);
  color: var(--accent-yellow, #f0ad4e);
}

@keyframes dot-blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}
</style>
