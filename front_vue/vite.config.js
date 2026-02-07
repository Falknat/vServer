import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { resolve } from 'path'

export default defineConfig({
  plugins: [
    vue(),

    AutoImport({
      imports: [
        'vue',
        'vue-router',
        'vue-i18n',
        'pinia',
      ],
      dirs: [
        'src/Core/composables',
        'src/Core/stores',
      ],
      vueTemplate: true,
      dts: 'src/auto-imports.d.ts',
    }),

    Components({
      dirs: [
        'src/Design/components/ui',
        'src/Design/components/layout',
        'src/Design/components/services',
        'src/Design/components/sites',
        'src/Design/components/proxies',
        'src/Design/components/vaccess',
        'src/Design/components/certs',
      ],
      dts: 'src/components.d.ts',
    }),
  ],

  server: {
    port: 4444,
  },

  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@core': resolve(__dirname, 'src/Core'),
      '@design': resolve(__dirname, 'src/Design'),
    },
  },
})
