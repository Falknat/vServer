<script setup>
import { useWailsEvents } from '@core/composables/useWailsEvents.js'

const { theme, isDark, toggleTheme, initTheme } = useTheme()
const appStore = useAppStore()
const servicesStore = useServicesStore()
const { on, off } = useWailsEvents()

onMounted(() => {
  initTheme()

  on('service:changed', (status) => {
    if (status && !servicesStore.isOperating) {
      servicesStore.list = Object.values(status)
      servicesStore.loaded = true
      const hasRunning = status.http?.status || status.https?.status
      appStore.setServerRunning(hasRunning)
    }
  })

  on('server:already_running', () => {
    appStore.setServerRunning(false)
  })
})

onUnmounted(() => {
  off('service:changed')
  off('server:already_running')
})
</script>

<template>
  <div class="app-layout">
    <TitleBar />
    <div class="app-body">
      <Sidebar />
      <div class="app-content">
        <main class="main-content">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </main>
        <AppFooter />
      </div>
    </div>
    <VNotification />
    <VModal />
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
  background: var(--bg-primary);
}

.app-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.app-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 40px var(--space-3xl);
}
</style>
