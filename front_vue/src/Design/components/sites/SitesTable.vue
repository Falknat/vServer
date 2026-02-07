<script setup>
const { t } = useI18n()
const router = useRouter()
const sitesStore = useSitesStore()
const certsStore = useCertsStore()
const { success, error } = useNotification()
const modal = useModal()

const openingFolder = ref('')

const openUrl = (host) => {
  if (window.runtime?.BrowserOpenURL) {
    window.runtime.BrowserOpenURL('http://' + host)
  } else {
    window.open('http://' + host, '_blank')
  }
}

const openFolder = async (host) => {
  openingFolder.value = host
  await sitesStore.openFolder(host)
  setTimeout(() => { openingFolder.value = '' }, 800)
}

const findCertForDomain = (domain, aliases = []) => {
  const allDomains = [domain, ...aliases.filter(a => !a.includes('*'))]
  for (const d of allDomains) {
    const direct = certsStore.list.find(c => c.domain === d && c.has_cert)
    if (direct) return direct

    const parts = d.split('.')
    if (parts.length >= 2) {
      const wildcard = '*.' + parts.slice(1).join('.')
      const wc = certsStore.list.find(c => c.domain === wildcard && c.has_cert)
      if (wc) return wc
    }

    for (const cert of certsStore.list) {
      if (cert.has_cert && cert.dns_names) {
        for (const dns of cert.dns_names) {
          if (dns === d) return cert
          if (dns.startsWith('*.')) {
            const base = dns.slice(2)
            const dParts = d.split('.')
            if (dParts.length >= 2 && dParts.slice(1).join('.') === base) return cert
          }
        }
      }
    }
  }
  return null
}

const confirmDelete = (site) => {
  modal.open({
    title: t('sites.deleteTitle'),
    message: t('sites.deleteConfirm', { name: site.name, host: site.host }),
    warning: t('sites.deleteWarning'),
    onConfirm: async () => {
      const result = await sitesStore.remove(site.host)
      if (result && !String(result).startsWith('Error')) success(t('notify.siteDeleted'))
      else error(String(result))
      modal.close()
    },
  })
}
</script>

<template>
  <section class="section">
    <div class="section-header">
      <h2 class="section-title">{{ t('sites.title') }}</h2>
      <VButton variant="primary" icon="fas fa-plus" @click="router.push('/sites/create')">
        {{ t('sites.add') }}
      </VButton>
    </div>
    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>{{ t('sites.name') }}</th>
            <th>{{ t('sites.host') }}</th>
            <th>{{ t('sites.alias') }}</th>
            <th>{{ t('sites.status') }}</th>
            <th>{{ t('sites.rootFile') }}</th>
            <th>{{ t('sites.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="site in sitesStore.list" :key="site.host">
            <td>
              <span class="cert-icon" :class="findCertForDomain(site.host, site.alias) ? (findCertForDomain(site.host, site.alias).is_expired ? 'cert-expired' : 'cert-valid') : 'cert-none'" :title="findCertForDomain(site.host, site.alias) ? (findCertForDomain(site.host, site.alias).is_expired ? 'SSL истёк' : `SSL (${findCertForDomain(site.host, site.alias).days_left} дн.)`) : 'Нет сертификата'">
                <i class="fas fa-shield-alt"></i>
              </span>
              {{ site.name }}
            </td>
            <td>
              <code class="clickable-link" @click="openUrl(site.host)">{{ site.host }} <i class="fas fa-external-link-alt"></i></code>
            </td>
            <td><code>{{ site.alias?.join(', ') || '—' }}</code></td>
            <td>
              <VBadge :variant="site.status === 'active' ? 'online' : 'offline'">
                {{ site.status }}
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
