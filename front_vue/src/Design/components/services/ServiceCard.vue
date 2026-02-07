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

const serviceInfoLabel = {
  HTTP: 'services.port',
  HTTPS: 'services.port',
  MySQL: 'services.port',
  PHP: 'services.ports',
  Proxy: 'services.rules',
}
</script>

<template>
  <div class="service-card">
    <div class="service-header">
      <h3 class="service-name">
        <i :class="serviceIcons[service.name] || 'fas fa-server'"></i>
        {{ service.name }}
      </h3>
      <VBadge :variant="service.status ? 'online' : 'offline'">
        {{ service.status ? t('common.enabled') : t('common.disabled') }}
      </VBadge>
    </div>
    <div class="service-info">
      <div class="info-row">
        <span class="info-label">{{ t(serviceInfoLabel[service.name] || 'services.port') }}:</span>
        <span class="info-value">{{ service.port || service.info || '—' }}</span>
      </div>
      <div v-if="service.info && service.port" class="info-row">
        <span class="info-label">{{ t('services.rules') }}:</span>
        <span class="info-value">{{ service.info }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.service-card {
  background: var(--glass-bg-light);
  backdrop-filter: var(--backdrop-blur);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-xl);
  padding: var(--space-lg);
  transition: all var(--transition-bounce);
  position: relative;
  overflow: hidden;
}

.service-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: var(--service-card-gradient);
  opacity: 0;
  transition: opacity var(--transition-slow);
}

.service-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--card-hover-shadow);
  border-color: var(--card-border-hover);
}

.service-card:hover::before {
  opacity: 1;
}

.service-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
}

.service-name {
  font-size: var(--text-md);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  margin: 0;
}

.service-name i {
  color: var(--accent-purple-light);
  font-size: var(--text-lg);
}

.service-info {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-size: var(--text-sm);
  color: var(--text-secondary);
  font-weight: var(--font-medium);
}

.info-value {
  font-size: 12px;
  color: var(--text-primary);
  font-weight: var(--font-semibold);
}
</style>
