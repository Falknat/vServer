import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from '@core/router/index.js'
import i18n from '@core/i18n/index.js'
import App from './App.vue'

import '@design/assets/css/variables.css'
import '@design/assets/css/themes/dark.css'
import '@design/assets/css/themes/light.css'
import '@design/assets/css/base.css'
import '@design/assets/css/forms.css'
import '@design/assets/css/tables.css'
import '@design/assets/css/transitions.css'
import '@fortawesome/fontawesome-free/css/all.min.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(i18n)

app.mount('#app')
