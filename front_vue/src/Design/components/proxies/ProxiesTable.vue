<script setup>

const { t } = useI18n()
const router = useRouter()
const proxiesStore = useProxiesStore()
const certsStore = useCertsStore()

const togglingStatus = ref('')
const dragIdx = ref(null)
const overIdx = ref(null)

const onDragStart = (i, e) => { dragIdx.value = i; e.dataTransfer.effectAllowed = 'move' }
const onDragOver = (i, e) => { e.preventDefault(); overIdx.value = i }
const onDragLeave = () => { overIdx.value = null }
const onDragEnd = () => { dragIdx.value = null; overIdx.value = null }

const onDrop = async (i) => {
  if (dragIdx.value === null || dragIdx.value === i) { dragIdx.value = null; overIdx.value = null; return }
  const items = [...proxiesStore.list]
  const [moved] = items.splice(dragIdx.value, 1)
  items.splice(i, 0, moved)
  proxiesStore.list = items
  dragIdx.value = null
  overIdx.value = null
  const config = await api.getConfig()
  config.Proxy_Service = items.map(p => config.Proxy_Service.find(c => c.ExternalDomain === p.ExternalDomain)).filter(Boolean)
  await api.saveConfig(JSON.stringify(config))
}

const toggleStatus = async (proxy) => {
  togglingStatus.value = proxy.ExternalDomain
  const config = await api.getConfig()
  const idx = config.Proxy_Service.findIndex(p => p.ExternalDomain === proxy.ExternalDomain)
  if (idx >= 0) {
    config.Proxy_Service[idx].Enable = !proxy.Enable
    await api.saveConfig(JSON.stringify(config))
    await api.reloadConfig()
    await proxiesStore.load()
  }
  togglingStatus.value = ''
}

const openUrl = (host) => {
  if (window.runtime?.BrowserOpenURL) {
    window.runtime.BrowserOpenURL('http://' + host)
  } else {
    window.open('http://' + host, '_blank')
  }
}

const findCertForDomain = (domain) => {
  const direct = certsStore.list.find(c => c.domain === domain && c.has_cert)
  if (direct) return direct

  const parts = domain.split('.')
  if (parts.length >= 2) {
    const wildcard = '*.' + parts.slice(1).join('.')
    const wc = certsStore.list.find(c => c.domain === wildcard && c.has_cert)
    if (wc) return wc
  }

  for (const cert of certsStore.list) {
    if (cert.has_cert && cert.dns_names) {
      for (const dns of cert.dns_names) {
        if (dns === domain) return cert
        if (dns.startsWith('*.')) {
          const base = dns.slice(2)
          const dParts = domain.split('.')
          if (dParts.length >= 2 && dParts.slice(1).join('.') === base) return cert
        }
      }
    }
  }
  return null
}
</script>

<template>
  <section class="section">
    <VSectionHeader :title="t('proxies.title')" addable @add="router.push('/proxies/create')" />
    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th><i class="fas fa-tag th-icon"></i> {{ t('sites.name') }}</th>
            <th><i class="fas fa-globe th-icon"></i> {{ t('proxies.externalDomain') }}</th>
            <th><i class="fas fa-server th-icon"></i> {{ t('proxies.localAddress') }}</th>
            <th><i class="fas fa-lock th-icon"></i> HTTPS</th>
            <th><i class="fas fa-circle-check th-icon"></i> {{ t('proxies.status') }}</th>
            <th class="th-actions">{{ t('proxies.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(proxy, i) in proxiesStore.list"
            :key="proxy.ExternalDomain"
            draggable="true"
            :class="{ 'drag-over': overIdx === i, 'dragging': dragIdx === i }"
            @dragstart="onDragStart(i, $event)"
            @dragover="onDragOver(i, $event)"
            @dragleave="onDragLeave"
            @drop="onDrop(i)"
            @dragend="onDragEnd"
          >
            <td>
              <i class="fas fa-grip-vertical drag-grip"></i>
              <span class="cert-icon" :class="findCertForDomain(proxy.ExternalDomain) ? (findCertForDomain(proxy.ExternalDomain).is_expired ? 'cert-expired' : 'cert-valid') : 'cert-none'">
                <i class="fas fa-shield-alt"></i>
              </span>
              {{ proxy.Name || '—' }}
            </td>
            <td>
              <code class="clickable-link" @click="openUrl(proxy.ExternalDomain)">{{ proxy.ExternalDomain }} <i class="fas fa-external-link-alt"></i></code>
            </td>
            <td><code>{{ proxy.LocalAddress }}:{{ proxy.LocalPort }}</code></td>
            <td>
              <VBadge :variant="proxy.ServiceHTTPSuse ? 'yes' : 'no'">
                {{ proxy.ServiceHTTPSuse ? 'HTTPS' : 'HTTP' }}
              </VBadge>
            </td>
            <td>
              <VBadge
                class="status-toggle"
                :variant="togglingStatus === proxy.ExternalDomain ? 'pending' : (proxy.Enable ? 'online' : 'offline')"
                @click="toggleStatus(proxy)"
              >
                {{ togglingStatus === proxy.ExternalDomain ? '...' : (proxy.Enable ? 'active' : 'disabled') }}
              </VBadge>
            </td>
            <td>
              <button class="icon-btn" :title="t('sites.editVaccess')" @click="router.push(`/vaccess/${proxy.ExternalDomain}?proxy=true`)">
                <i class="fas fa-user-lock"></i>
              </button>
              <button class="icon-btn" :title="t('certs.title')" @click="router.push(`/certs/${proxy.ExternalDomain}`)">
                <i class="fas fa-shield-alt"></i>
              </button>
              <button class="icon-btn" :title="t('sites.edit')" @click="router.push(`/proxies/edit/${proxy.ExternalDomain}`)">
                <i class="fas fa-edit"></i>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
