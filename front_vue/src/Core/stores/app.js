import { defineStore } from 'pinia'
import { STORAGE_KEYS, THEME, LOCALE } from '@core/constants.js'

export const useAppStore = defineStore('app', {
  state: () => ({
    theme: localStorage.getItem(STORAGE_KEYS.THEME) || THEME.DARK,
    locale: localStorage.getItem(STORAGE_KEYS.LOCALE) || LOCALE.RU,
    loading: true,
    serverRunning: false,
  }),

  getters: {
    isDark: (state) => state.theme === THEME.DARK,
  },

  actions: {
    toggleTheme() {
      this.theme = this.isDark ? THEME.LIGHT : THEME.DARK
      localStorage.setItem(STORAGE_KEYS.THEME, this.theme)
      document.documentElement.setAttribute('data-theme', this.theme)
    },

    setLocale(locale) {
      this.locale = locale
      localStorage.setItem(STORAGE_KEYS.LOCALE, locale)
    },

    setLoading(value) {
      this.loading = value
    },

    setServerRunning(value) {
      this.serverRunning = value
    },
  },
})
