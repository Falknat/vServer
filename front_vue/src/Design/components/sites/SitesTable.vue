<script setup>
const { t } = useI18n()
const router = useRouter()
const sitesStore = useSitesStore()
const { success, error } = useNotification()
const { confirmDelete: showConfirm } = useConfirm()
const { findCertForDomain } = useCertLookup()

const openingFolder = ref('')
const togglingStatus = ref('')

const siteList = computed({
  get: () => sitesStore.list,
  set: (v) => { sitesStore.list = v },
})

const { dragIndex, dragOverIndex, onDragStart, onDragOver, onDragLeave, onDrop, onDragEnd } = useDraggable(siteList, {
  onReorder: async (items) => {
    const config = await api.getConfig()
    config.Site_www = items.map(s => config.Site_www.find(c => c.host === s.host)).filter(Boolean)
    await api.saveConfig(JSON.stringify(config))
  },
})

const toggleStatus = async (site) => {
  togglingStatus.value = site.host
  const newStatus = site.status === 'active' ? 'inactive' : 'active'
  const config = await api.getConfig()
  const idx = config.Site_www.findIndex(s => s.host === site.host)
  if (idx >= 0) {
    config.Site_www[idx].status = newStatus
    await api.saveConfig(JSON.stringify(config))
    await api.updateSiteCache()
    await sitesStore.load()
  }
  togglingStatus.value = ''
}

const openFolder = async (host) => {
  openingFolder.value = host
  await sitesStore.openFolder(host)
  setTimeout(() => { openingFolder.value = '' }, 800)
}

const confirmDelete = async (site) => {
  const result = await showConfirm({
    title: t('sites.deleteTitle'),
    message: `${site.name} (${site.host})`,
  })
  if (result) {
    const res = await sitesStore.remove(site.host)
    if (isSuccess(res)) success(t('notify.siteDeleted'))
    else error(String(res))
  }
}
</script>

<template>
  <section class="section">
    <VSectionHeader :title="t('sites.title')" addable @add="router.push('/sites/create')" />
    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th><i class="fas fa-tag th-icon"></i> {{ t('sites.name') }}</th>
            <th><i class="fas fa-globe th-icon"></i> {{ t('sites.host') }}</th>
            <th><i class="fas fa-link th-icon"></i> {{ t('sites.alias') }}</th>
            <th><i class="fas fa-circle-check th-icon"></i> {{ t('sites.status') }}</th>
            <th><i class="fas fa-file th-icon"></i> {{ t('sites.rootFile') }}</th>
            <th class="th-actions">{{ t('sites.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(site, i) in sitesStore.list"
            :key="site.host"
            draggable="true"
            :class="{ 'drag-over': dragOverIndex === i, 'dragging': dragIndex === i }"
            @dragstart="onDragStart(i, $event)"
            @dragover="onDragOver(i, $event)"
            @dragleave="onDragLeave"
            @drop="onDrop(i)"
            @dragend="onDragEnd($event)"
          >
            <td>
              <i class="fas fa-grip-vertical drag-grip"></i>
              <span class="cert-icon" :class="findCertForDomain(site.host, site.alias) ? (findCertForDomain(site.host, site.alias).is_expired ? 'cert-expired' : 'cert-valid') : 'cert-none'" :title="findCertForDomain(site.host, site.alias) ? (findCertForDomain(site.host, site.alias).is_expired ? 'SSL истёк' : `SSL (${findCertForDomain(site.host, site.alias).days_left} дн.)`) : 'Нет сертификата'">
                <i class="fas fa-shield-alt"></i>
              </span>
              {{ site.name }}
            </td>
            <td>
              <code class="clickable-link" @click="openUrl('http://' + site.host)">{{ site.host }} <i class="fas fa-external-link-alt"></i></code>
            </td>
            <td><code>{{ site.alias?.join(', ') || '—' }}</code></td>
            <td>
              <VBadge
                class="status-toggle"
                :variant="togglingStatus === site.host ? 'pending' : (site.status === 'active' ? 'online' : 'offline')"
                @click="toggleStatus(site)"
              >
                {{ togglingStatus === site.host ? '...' : site.status }}
              </VBadge>
            </td>
            <td><code>{{ site.root_file }}</code></td>
            <td>
              <button class="icon-btn" :title="t('sites.openFolder')" :disabled="openingFolder === site.host" @click="openFolder(site.host)">
                <i v-if="openingFolder === site.host" class="fas fa-spinner icon-spin"></i>
                <i v-else class="fas fa-folder-open"></i>
              </button>
              <button class="icon-btn" :title="t('sites.editVaccess')" @click="router.push(`/vaccess/${site.host}`)">
                <i class="fas fa-user-lock"></i>
              </button>
              <button class="icon-btn" :title="t('certs.title')" @click="router.push(`/certs/${site.host}`)">
                <i class="fas fa-shield-alt"></i>
              </button>
              <button class="icon-btn" :title="t('sites.edit')" @click="router.push(`/sites/edit/${site.host}`)">
                <i class="fas fa-edit"></i>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
