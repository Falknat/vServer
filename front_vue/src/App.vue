<script setup>
import MainLayout from '@design/layouts/MainLayout.vue'

const loading = ref(true)
const servicesStore = useServicesStore()

const unwatch = watch(
  () => servicesStore.list,
  (list) => {
    const hasRunning = list.some(s => s.status === true)
    if (hasRunning) {
      setTimeout(() => { loading.value = false }, 300)
      unwatch()
    }
  },
  { deep: true }
)

onMounted(() => {
  setTimeout(() => { loading.value = false }, 15000)
})
</script>

<template>
  <transition name="loader-fade">
    <div v-if="loading" class="app-loader">
      <div class="loader-content">
        <div class="loader-icon">🚀</div>
        <div class="loader-text">Запуск vServer...</div>
        <div class="loader-spinner"></div>
      </div>
    </div>
  </transition>
  <MainLayout />
</template>

<style>
.app-loader {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, var(--bg-primary) 0%, var(--bg-secondary, #0f1729) 50%, var(--bg-tertiary, #1a1040) 100%);
  z-index: 10000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loader-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.loader-icon {
  font-size: 64px;
  animation: rocket-bounce 1.5s ease-in-out infinite;
}

.loader-text {
  font-size: var(--text-xl, 20px);
  font-weight: var(--font-semibold, 600);
  background: linear-gradient(135deg, #7c3aed 0%, #a78bfa 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.loader-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(139, 92, 246, 0.2);
  border-top-color: var(--accent-purple, #7c3aed);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes rocket-bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-20px); }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loader-fade-leave-active {
  transition: opacity 0.5s ease;
}

.loader-fade-leave-to {
  opacity: 0;
}
</style>
