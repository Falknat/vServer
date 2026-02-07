<script setup>
const { theme, isDark, toggleTheme, initTheme } = useTheme()
const appStore = useAppStore()

onMounted(() => {
  initTheme()
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
  padding: var(--space-lg);
}
</style>
