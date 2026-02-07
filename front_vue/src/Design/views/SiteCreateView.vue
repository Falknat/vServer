<script setup>
const { t } = useI18n()
const router = useRouter()
const sitesStore = useSitesStore()
const { success, error } = useNotification()

const form = reactive({
  name: '',
  host: '',
  alias: '',
  rootFile: 'index.html',
  status: 'active',
  routing: true,
  certMode: 'none',
})

const creating = ref(false)

const rootFileOptions = [
  { value: 'index.html', label: 'index.html' },
  { value: 'index.php', label: 'index.php' },
]

const statusOptions = computed(() => [
  { value: 'active', label: `Active (${t('common.enabled')})` },
  { value: 'inactive', label: `Inactive (${t('common.disabled')})` },
])

const createSite = async () => {
  if (!form.name || !form.host) return
  creating.value = true
  const siteData = {
    name: form.name,
    host: form.host,
    alias: form.alias ? form.alias.split(',').map(a => a.trim()).filter(Boolean) : [],
    root_file: form.rootFile,
    status: form.status,
    root_file_routing: form.routing,
    AutoCreateSSL: form.certMode === 'auto',
  }
  const result = await sitesStore.create(siteData)
  creating.value = false
  if (result === 'OK') {
    success(t('notify.siteCreated'))
    router.push('/')
  } else {
    error(String(result))
  }
}
</script>

<template>
  <div class="vaccess-page">
    <Breadcrumbs :items="[t('sites.create')]" />

    <PageHeader icon="fas fa-plus-circle" :title="t('sites.create')" :subtitle="t('sites.createDesc')">
      <template #actions>
        <VButton variant="success" icon="fas fa-check" :loading="creating" @click="createSite">
          {{ t('sites.createBtn') }}
        </VButton>
      </template>
    </PageHeader>

    <div class="form-section">
      <div class="settings-form">
        <h3 class="form-subsection-title"><i class="fas fa-info-circle"></i> {{ t('common.basicInfo') }}</h3>

        <VInput v-model="form.name" :label="t('sites.formName')" :placeholder="t('sites.formNamePlaceholder')" required />
        <VInput v-model="form.host" :label="t('sites.formHost')" :placeholder="t('sites.formHostPlaceholder')" :hint="t('sites.formHostHint')" required />
        <VInput v-model="form.alias" :label="t('sites.formAlias')" :placeholder="t('sites.formAliasPlaceholder')" :hint="t('sites.formAliasHint')" />

        <div class="form-row">
          <VSelect v-model="form.rootFile" :label="t('sites.formRootFile')" :options="rootFileOptions" required />
          <VSelect v-model="form.status" :label="t('sites.formStatus')" :options="statusOptions" />
        </div>

        <div class="form-group">
          <div class="form-label-row">
            <VTooltip :text="t('sites.formRoutingHint')" />
            <label class="form-label">{{ t('sites.formRouting') }}</label>
          </div>
          <VToggle v-model="form.routing" :label="t('common.enabled')" />
        </div>

        <SslUploadSection v-model="form.certMode" />
      </div>
    </div>
  </div>
</template>
