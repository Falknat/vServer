<script setup>
const { t } = useI18n()

defineProps({
  text: { type: String, default: '' },
})
</script>

<template>
  <div class="v-loader">
    <div class="v-loader-content">
      <div class="v-loader-icon">🚀</div>
      <div class="v-loader-text">{{ text || t('app.loading') }}</div>
      <div class="v-loader-spinner"></div>
    </div>
  </div>
</template>

<style scoped>
.v-loader {
  position: fixed;
  inset: 0;
  background: var(--bg-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-loader);
}

.v-loader-content {
  text-align: center;
}

.v-loader-icon {
  font-size: 48px;
  margin-bottom: var(--space-md);
  animation: pulse 2s infinite;
}

.v-loader-text {
  font-size: var(--text-lg);
  color: var(--text-secondary);
  margin-bottom: var(--space-lg);
}

.v-loader-spinner {
  width: 40px;
  height: 40px;
  margin: 0 auto;
  border: 3px solid var(--glass-border);
  border-top-color: var(--accent-purple);
  border-radius: var(--radius-full);
  animation: spin 1s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }
@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.5; } }
</style>
