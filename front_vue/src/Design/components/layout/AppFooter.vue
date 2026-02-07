<script setup>
const { t } = useI18n()
</script>

<template>
  <footer class="footer">
    <p>{{ t('app.footer') }}</p>
  </footer>
</template>

<style scoped>
.footer {
  padding: var(--space-sm) var(--space-lg);
  text-align: center;
  color: var(--text-muted);
  font-size: var(--text-xs);
  border-top: 1px solid var(--glass-border);
  background: var(--glass-bg-dark);
}
</style>
