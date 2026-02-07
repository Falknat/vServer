import { createI18n } from 'vue-i18n'
import ru from './ru.json'
import en from './en.json'
import { STORAGE_KEYS, LOCALE } from '@core/constants.js'

const savedLocale = localStorage.getItem(STORAGE_KEYS.LOCALE) || LOCALE.RU

const i18n = createI18n({
  legacy: false,
  locale: savedLocale,
  fallbackLocale: LOCALE.RU,
  messages: { ru, en },
})

export default i18n
