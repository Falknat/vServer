import { useAppStore } from '@core/stores/app.js'

export function useTheme() {
  const appStore = useAppStore()

  const initTheme = () => {
    document.documentElement.setAttribute('data-theme', appStore.theme)
  }

  return {
    theme: computed(() => appStore.theme),
    isDark: computed(() => appStore.isDark),
    toggleTheme: () => appStore.toggleTheme(),
    initTheme,
  }
}
