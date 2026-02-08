<script setup>
const { t } = useI18n()
const router = useRouter()
const sitesStore = useSitesStore()
const { success, error } = useNotification()
const { confirmDelete: showConfirm } = useConfirm()

const props = defineProps({
  host: { type: String, required: true },
})

const form = reactive({
  name: '',
  host: '',
  alias: [],
  rootFile: 'index.html',
  status: 'active',
  routing: true,
  compression: true,
  autoSSL: false,
})


const saving = ref(false)
const aliasInput = ref('')

const rootFileOptions = [
  { value: 'index.html', label: 'index.html' },
  { value: 'index.php', label: 'index.php' },
]

onMounted(async () => {
  if (!sitesStore.loaded) await sitesStore.load()
  const site = sitesStore.list.find(s => s.host === props.host)
  if (site) {
    form.name = site.name
    form.host = site.host
    form.alias = [...(site.alias || [])]
    form.rootFile = site.root_file
    form.status = site.status
    form.routing = site.root_file_routing
    form.compression = site.Compression !== false
    form.autoSSL = site.AutoCreateSSL || false
  }
})

const addAlias = () => {
  const val = aliasInput.value.trim()
  if (val && !form.alias.includes(val)) {
    form.alias.push(val)
    aliasInput.value = ''
  }
}

const removeAlias = (index) => {
  form.alias.splice(index, 1)
}

const saveSite = async () => {
  saving.value = true
  const config = await api.getConfig()
  const idx = config.Site_www.findIndex(s => s.host === props.host)
  if (idx >= 0) {
    config.Site_www[idx].host = form.host
    config.Site_www[idx].name = form.name
    config.Site_www[idx].alias = form.alias
    config.Site_www[idx].root_file = form.rootFile
    config.Site_www[idx].status = form.status
    config.Site_www[idx].root_file_routing = form.routing
    config.Site_www[idx].Compression = form.compression
    config.Site_www[idx].AutoCreateSSL = form.autoSSL
    const result = await api.saveConfig(JSON.stringify(config))
    if (isSuccess(result)) {
      await sitesStore.load()
      success(t('notify.dataSaved'))
      router.push('/')
    } else {
      error(result)
    }
  } else {
    error('Site not found in config')
  }
  saving.value = false
}

const confirmDelete = async () => {
  const result = await showConfirm({
    title: t('sites.deleteTitle'),
    message: `${form.name} (${form.host})`,
  })
  if (result) {
    const res = await sitesStore.remove(form.host)
    if (isSuccess(res)) {
      success(t('notify.siteDeleted'))
      router.push('/')
    } else {
      error(String(res))
    }
  }
}
</script>

<template>
  <div class="vaccess-page">
    <Breadcrumbs :items="[host]" />

    <PageHeader icon="fas fa-edit" :title="`${t('sites.edit')} — ${host}`">
      <template #actions>
        <VButton variant="danger" icon="fas fa-trash" @click="confirmDelete">{{ t('common.delete') }}</VButton>
        <VButton icon="fas fa-times" @click="router.push('/')">{{ t('common.cancel') }}</VButton>
        <VButton variant="success" icon="fas fa-save" :loading="saving" @click="saveSite">{{ t('common.save') }}</VButton>
      </template>
    </PageHeader>

    <div class="form-section">
      <div class="settings-form">
        <!-- Статус -->
        <div class="form-group">
          <label class="form-label">{{ t('sites.formStatus') }}:</label>
          <div class="status-toggle">
            <button class="status-btn" :class="{ active: form.status === 'active' }" @click="form.status = 'active'">
              <i class="fas fa-check-circle"></i> Active
            </button>
            <button class="status-btn" :class="{ active: form.status === 'inactive' }" @click="form.status = 'inactive'">
              <i class="fas fa-times-circle"></i> Inactive
            </button>
          </div>
        </div>

        <!-- Основная информация -->
        <VInput v-model="form.host" :label="t('sites.host')" required />
        <VInput v-model="form.name" :label="t('sites.formName')" required />

        <!-- Alias с тегами -->
        <div class="form-group">
          <label class="form-label">{{ t('sites.formAlias') }}:</label>
          <div class="tag-input-row">
            <input v-model="aliasInput" class="form-input" :placeholder="t('sites.formAliasPlaceholder')" @keydown.enter.prevent="addAlias">
            <button class="action-btn" @click="addAlias"><i class="fas fa-plus"></i> {{ t('common.add') }}</button>
          </div>
          <div v-if="form.alias.length" class="tags-container">
            <span v-for="(alias, i) in form.alias" :key="alias" class="tag">
              {{ alias }}
              <button class="tag-remove" @click="removeAlias(i)"><i class="fas fa-times"></i></button>
            </span>
          </div>
        </div>

        <!-- Root File -->
        <VSelect v-model="form.rootFile" :label="t('sites.formRootFile')" :options="rootFileOptions" />

        <!-- Toggles -->
        <div class="form-row form-row-3">
          <div class="form-group">
            <label class="form-label">{{ t('sites.formCompression') }}:</label>
            <VToggle v-model="form.compression" :label="t('common.enabled')" />
          </div>
          <div class="form-group">
            <label class="form-label">{{ t('sites.formRouting') }}:</label>
            <VToggle v-model="form.routing" :label="t('common.enabled')" />
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

<style>
.tag-input-row {
  display: flex;
  gap: var(--space-sm);
}

.form-input {
  flex: 1;
  padding: 10px 14px;
  background: var(--glass-bg-dark);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius);
  color: var(--text-primary);
  font-size: var(--text-base);
  outline: none;
  transition: all var(--transition-base);
}

.form-input:focus {
  border-color: rgba(var(--accent-rgb), 0.5);
  box-shadow: 0 0 12px rgba(var(--accent-rgb), 0.2);
}

.form-input::placeholder {
  color: var(--text-muted);
  opacity: 0.5;
}

.action-btn {
  padding: var(--space-sm) var(--space-md);
  background: rgba(var(--accent-rgb), 0.15);
  border: 1px solid rgba(var(--accent-rgb), 0.3);
  border-radius: var(--radius);
  color: var(--accent-purple-light);
  font-size: var(--text-base);
  font-weight: var(--font-semibold);
  cursor: pointer;
  transition: all var(--transition-base);
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  white-space: nowrap;
}

.action-btn:hover {
  background: rgba(var(--accent-rgb), 0.25);
  border-color: rgba(var(--accent-rgb), 0.5);
}

/* Tags */
.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-sm);
  padding: 12px;
  background: var(--glass-bg-dark);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius);
  min-height: 48px;
}

.tag {
  display: inline-flex;
  align-items: center;
  gap: var(--space-sm);
  padding: 4px 10px;
  background: rgba(var(--accent-rgb), 0.2);
  border: 1px solid rgba(var(--accent-rgb), 0.4);
  border-radius: 16px;
  color: var(--text-primary);
  font-size: 12px;
  font-weight: var(--font-medium);
}

.tag-remove {
  background: transparent;
  border: none;
  color: var(--accent-red);
  cursor: pointer;
  padding: 0;
  width: 14px;
  height: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-full);
  transition: all var(--transition-base);
  font-size: 10px;
}

.tag-remove:hover {
  background: rgba(var(--danger-rgb), 0.2);
}
</style>
