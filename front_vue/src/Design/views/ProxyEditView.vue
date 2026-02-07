<script setup>
const { t } = useI18n()
const router = useRouter()
const proxiesStore = useProxiesStore()
const { success, error } = useNotification()
const modal = useModal()

const props = defineProps({
  domain: { type: String, required: true },
})

const form = reactive({
  domain: '',
  localAddr: '127.0.0.1',
  localPort: '',
  enabled: true,
  serviceHttps: false,
  autoHttps: true,
  autoSSL: false,
})

const saving = ref(false)

onMounted(async () => {
  if (!proxiesStore.loaded) await proxiesStore.load()
  const proxy = proxiesStore.list.find(p => p.ExternalDomain === props.domain)
  if (proxy) {
    form.domain = proxy.ExternalDomain
    form.localAddr = proxy.LocalAddress
    form.localPort = proxy.LocalPort
    form.enabled = proxy.Enable
    form.serviceHttps = proxy.ServiceHTTPSuse
    form.autoHttps = proxy.AutoHTTPS
    form.autoSSL = proxy.AutoCreateSSL || false
  }
})

const saveProxy = async () => {
  saving.value = true
  success(t('notify.dataSaved'))
  saving.value = false
  router.push('/')
}

const confirmDelete = () => {
  modal.open({
    title: t('proxies.deleteTitle'),
    message: t('proxies.deleteConfirm', { domain: form.domain }),
    onConfirm: async () => {
      success(t('notify.proxyDeleted'))
      modal.close()
      router.push('/')
    },
  })
}
</script>

<template>
  <div class="vaccess-page">
    <Breadcrumbs :items="[domain]" />

    <PageHeader icon="fas fa-edit" :title="`${t('sites.edit')} — ${domain}`">
      <template #actions>
        <VButton variant="danger" icon="fas fa-trash" @click="confirmDelete">{{ t('common.delete') }}</VButton>
        <VButton icon="fas fa-times" @click="router.push('/')">{{ t('common.cancel') }}</VButton>
        <VButton variant="success" icon="fas fa-save" :loading="saving" @click="saveProxy">{{ t('common.save') }}</VButton>
      </template>
    </PageHeader>

    <div class="form-section">
      <div class="settings-form">
        <!-- Статус -->
        <div class="form-group">
          <label class="form-label">{{ t('proxies.status') }}:</label>
          <div class="status-toggle">
            <button class="status-btn" :class="{ active: form.enabled }" @click="form.enabled = true">
              <i class="fas fa-check-circle"></i> {{ t('common.enabled') }}
            </button>
            <button class="status-btn" :class="{ active: !form.enabled }" @click="form.enabled = false">
              <i class="fas fa-times-circle"></i> {{ t('common.disabled') }}
            </button>
          </div>
        </div>

        <!-- Адрес и порт -->
        <div class="form-row">
          <VInput v-model="form.localAddr" :label="t('proxies.formLocalAddr')" required />
          <VInput v-model="form.localPort" :label="t('proxies.formLocalPort')" required />
        </div>

        <!-- Toggles -->
        <div class="form-row form-row-3">
          <div class="form-group">
            <label class="form-label">{{ t('proxies.formServiceHttps') }}:</label>
            <VToggle v-model="form.serviceHttps" :label="t('common.enabled')" />
          </div>
          <div class="form-group">
            <label class="form-label">{{ t('proxies.formAutoHttps') }}:</label>
            <VToggle v-model="form.autoHttps" :label="t('common.enabled')" />
          </div>
          <div class="form-group">
            <label class="form-label">Auto SSL:</label>
            <VToggle v-model="form.autoSSL" :label="t('common.enabled')" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
