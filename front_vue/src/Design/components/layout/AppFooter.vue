<script setup>
const { t } = useI18n()

const currentYear = new Date().getFullYear()

const openSite = () => {
  if (window.runtime?.BrowserOpenURL) {
    window.runtime.BrowserOpenURL('https://vserf.ru')
  } else {
    window.open('https://vserf.ru', '_blank')
  }
}
</script>

<template>
  <footer class="footer">
    <p>vServer Admin Panel © 2025-{{ currentYear }} | <a href="#" class="footer-link" @click.prevent="openSite">https://vserf.ru</a> | {{ t('app.footerAuthor') }}</p>
  </footer>
</template>

<style scoped>
.footer {
  padding: var(--space-sm) var(--space-lg);
  text-align: center;
  color: var(--text-muted);
  font-size: var(--text-sm);
  border-top: 1px solid var(--glass-border);
  background: var(--glass-bg-dark);
}

.footer-link {
  color: var(--text-secondary);
  text-decoration: underline;
  transition: color var(--transition-fast);
}

.footer-link:hover {
  color: var(--text-primary);
}
</style>
