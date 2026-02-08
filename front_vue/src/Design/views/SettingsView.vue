<script setup>
const { t } = useI18n()
const configStore = useConfigStore()
const servicesStore = useServicesStore()
const { success } = useNotification()

const form = reactive({
  mysqlHost: '',
  mysqlPort: 3306,
  phpHost: '',
  phpPort: 8000,
  proxyEnabled: false,
  acmeEnabled: false,
})

const saving = ref(false)

onMounted(async () => {
  await configStore.load()
  const s = configStore.softSettings
  form.mysqlHost = s.mysql_host || '127.0.0.1'
  form.mysqlPort = s.mysql_port || 3306
  form.phpHost = s.php_host || 'localhost'
  form.phpPort = s.php_port || 8000
  form.proxyEnabled = s.proxy_enabled || false
  form.acmeEnabled = s.ACME_enabled || false
})

const saveSettings = async () => {
  saving.value = true
  const configData = {
    ...configStore.data,
    Soft_Settings: {
      mysql_host: form.mysqlHost,
      mysql_port: Number(form.mysqlPort),
      php_host: form.phpHost,
      php_port: Number(form.phpPort),
      proxy_enabled: form.proxyEnabled,
      ACME_enabled: form.acmeEnabled,
    },
  }
  await configStore.save(configData)
  await servicesStore.restartAll()
  saving.value = false
  success(t('notify.settingsSaved'))
}

const toggleProxy = async () => {
  form.proxyEnabled = !form.proxyEnabled
  if (form.proxyEnabled) await servicesStore.enableProxy()
  else await servicesStore.disableProxy()
  success(form.proxyEnabled ? t('notify.proxyEnabled') : t('notify.proxyDisabled'))
}

const toggleAcme = async () => {
  form.acmeEnabled = !form.acmeEnabled
  if (form.acmeEnabled) await servicesStore.enableACME()
  else await servicesStore.disableACME()
  success(form.acmeEnabled ? t('notify.certManagerEnabled') : t('notify.certManagerDisabled'))
}
</script>

<template>
  <div class="settings-view">
    <div class="settings-header">
      <VSectionHeader :title="t('settings.title')" />
      <VButton variant="success" icon="fas fa-save" :loading="saving" @click="saveSettings">
        {{ saving ? t('settings.saving') : t('settings.save') }}
      </VButton>
    </div>

    <div class="settings-grid">
      <div class="settings-card">
        <h3 class="settings-card-title"><i class="fas fa-database"></i> {{ t('settings.mysql') }}</h3>
        <div class="settings-form">
          <VInput v-model="form.mysqlHost" :label="t('settings.hostAddr')" placeholder="127.0.0.1" />
          <VInput v-model="form.mysqlPort" :label="t('settings.port')" type="number" placeholder="3306" />
        </div>
      </div>

      <div class="settings-card">
        <h3 class="settings-card-title"><i class="fab fa-php"></i> {{ t('settings.php') }}</h3>
        <div class="settings-form">
          <VInput v-model="form.phpHost" :label="t('settings.hostAddr')" placeholder="localhost" />
          <VInput v-model="form.phpPort" :label="t('settings.port')" type="number" placeholder="8000" />
        </div>
      </div>

      <div class="settings-card">
        <h3 class="settings-card-title">
          <VTooltip :text="t('settings.proxyHint')" />
          <i class="fas fa-network-wired"></i> {{ t('settings.proxyManager') }}
        </h3>
        <div class="settings-form">
          <VToggle v-model="form.proxyEnabled" :label="t('settings.proxyManager')" @update:model-value="toggleProxy" />
        </div>
      </div>

      <div class="settings-card">
        <h3 class="settings-card-title">
          <VTooltip :text="t('settings.certHint')" />
          <i class="fas fa-certificate"></i> {{ t('settings.certManager') }}
        </h3>
        <div class="settings-form">
          <VToggle v-model="form.acmeEnabled" :label="t('settings.certManager')" @update:model-value="toggleAcme" />
        </div>
      </div>
    </div>
  </div>
</template>

<style>
.settings-view {
  animation: fadeIn var(--transition-slow);
}

.settings-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-lg);
}

.settings-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--space-lg);
}

.settings-card {
  background: rgba(var(--accent-rgb), 0.02);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius);
  padding: var(--space-lg);
  transition: all var(--transition-base);
}

.settings-card:hover {
  background: rgba(var(--accent-rgb), 0.04);
  border-color: rgba(var(--accent-rgb), 0.2);
}

.settings-card-title {
  font-size: var(--text-md);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin: 0 0 20px 0;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(var(--accent-rgb), 0.1);
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.settings-card-title i {
  color: var(--accent-purple-light);
}

@media (max-width: 900px) {
  .settings-grid { grid-template-columns: 1fr; }
}
</style>
