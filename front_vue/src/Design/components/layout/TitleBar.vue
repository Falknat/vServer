<script setup>
const { t, locale } = useI18n()
const appStore = useAppStore()
const { isDark, toggleTheme } = useTheme()

const toggleLocale = () => {
  const next = locale.value === 'ru' ? 'en' : 'ru'
  locale.value = next
  localStorage.setItem('vserver-locale', next)
}

const serverToggle = async () => {
  if (appStore.serverRunning) {
    appStore.setServerRunning(false)
  } else {
    appStore.setServerRunning(true)
  }
}
</script>

<template>
  <div class="title-bar" style="--wails-draggable: drag">
    <div class="title-bar-left">
      <div class="app-logo">
        <span class="logo-icon">🚀</span>
        <span class="logo-text">{{ t('app.logo') }}</span>
      </div>
      <div class="server-status">
        <span class="status-indicator" :class="appStore.serverRunning ? 'status-online' : 'status-offline'"></span>
        <span class="status-text">{{ appStore.serverRunning ? t('server.running') : t('server.stopped') }}</span>
      </div>
    </div>
    <div class="title-bar-right" style="--wails-draggable: no-drag">
      <button class="server-control-btn" @click="serverToggle">
        <i class="fas fa-power-off"></i>
        <span class="btn-text">{{ appStore.serverRunning ? t('server.stop') : t('server.start') }}</span>
      </button>
      <button class="locale-btn" @click="toggleLocale" :title="locale === 'ru' ? 'English' : 'Русский'">
        {{ locale === 'ru' ? 'EN' : 'RU' }}
      </button>
      <button class="theme-btn" @click="toggleTheme" :title="isDark ? t('theme.light') : t('theme.dark')">
        <i :class="isDark ? 'fas fa-sun' : 'fas fa-moon'"></i>
      </button>
      <button class="window-btn minimize-btn" :title="t('window.minimize')">
        <i class="fas fa-window-minimize"></i>
      </button>
      <button class="window-btn maximize-btn" :title="t('window.maximize')">
        <i class="far fa-window-maximize"></i>
      </button>
      <button class="window-btn close-btn" :title="t('window.close')">
        <i class="fas fa-times"></i>
      </button>
    </div>
  </div>
</template>

<style scoped>
.title-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: var(--header-height);
  padding: 0 var(--space-md);
  background: var(--titlebar-bg);
  border-bottom: 1px solid var(--titlebar-border);
  user-select: none;
}

.title-bar-left,
.title-bar-right {
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

.app-logo {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.logo-icon { font-size: var(--text-xl); }

.logo-text {
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
  color: var(--text-primary);
}

.server-status {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-xs) var(--space-md);
  border-radius: var(--radius-lg);
  background: var(--glass-bg);
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: var(--radius-full);
}

.status-online {
  background: var(--accent-green);
  box-shadow: 0 0 8px var(--accent-green);
}

.status-offline {
  background: var(--accent-red);
  box-shadow: 0 0 8px var(--accent-red);
}

.status-text {
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.server-control-btn {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-xs) var(--space-md);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-md);
  background: var(--glass-bg);
  color: var(--text-primary);
  cursor: pointer;
  transition: all var(--transition-base);
  font-size: var(--text-sm);
}

.server-control-btn:hover {
  border-color: var(--glass-border-hover);
  background: var(--glass-bg-light);
}

.theme-btn,
.window-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.locale-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
  font-size: 11px;
  font-weight: var(--font-bold);
  letter-spacing: 0.5px;
}

.locale-btn:hover { color: var(--text-primary); }

.theme-btn:hover { color: var(--accent-yellow); }
.minimize-btn:hover { color: var(--accent-green); }
.maximize-btn:hover { color: var(--accent-cyan); }
.close-btn:hover { color: var(--accent-red); background: rgba(var(--danger-rgb), 0.15); }
</style>
