<script setup>
const { t } = useI18n()
const router = useRouter()
const proxiesStore = useProxiesStore()
const { success, error } = useNotification()

const form = reactive({
  domain: '',
  localAddr: '127.0.0.1',
  localPort: '',
  serviceHttps: false,
  autoHttps: true,
  certMode: 'none',
})

const creating = ref(false)

const createProxy = async () => {
  if (!form.domain || !form.localPort) return
  creating.value = true
  const result = await proxiesStore.load()
  creating.value = false
  success(t('notify.dataSaved'))
  router.push('/')
}
</script>

<template>
  <div class="vaccess-page">
    <Breadcrumbs :items="[t('proxies.create')]" />

    <PageHeader icon="fas fa-plus-circle" :title="t('proxies.create')" :subtitle="t('proxies.createDesc')">
      <template #actions>
        <VButton variant="success" icon="fas fa-check" :loading="creating" @click="createProxy">
          {{ t('proxies.createBtn') }}
        </VButton>
      </template>
    </PageHeader>

    <div class="form-section">
      <div class="settings-form">
        <h3 class="form-subsection-title"><i class="fas fa-info-circle"></i> {{ t('common.basicInfo') }}</h3>

        <VInput v-model="form.domain" :label="t('proxies.formDomain')" :placeholder="t('proxies.formDomainPlaceholder')" :hint="t('proxies.formDomainHint')" required />

        <div class="form-row">
          <VInput v-model="form.localAddr" :label="t('proxies.formLocalAddr')" placeholder="127.0.0.1" required />
          <VInput v-model="form.localPort" :label="t('proxies.formLocalPort')" placeholder="3000" required />
        </div>

        <div class="form-row">
          <div class="form-group">
            <div class="form-label-row">
              <VTooltip :text="t('proxies.formServiceHttpsHint')" />
              <label class="form-label">{{ t('proxies.formServiceHttps') }}</label>
            </div>
            <VToggle v-model="form.serviceHttps" :label="t('common.enabled')" />
          </div>
          <div class="form-group">
            <div class="form-label-row">
              <VTooltip :text="t('proxies.formAutoHttpsHint')" />
              <label class="form-label">{{ t('proxies.formAutoHttps') }}</label>
            </div>
            <VToggle v-model="form.autoHttps" :label="t('common.enabled')" />
          </div>
        </div>

        <SslUploadSection v-model="form.certMode" />
      </div>
    </div>
  </div>
</template>
