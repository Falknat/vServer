<script setup>

const { t } = useI18n()
const router = useRouter()
const proxiesStore = useProxiesStore()
const { findCertForDomain } = useCertLookup()

const togglingStatus = ref('')

const proxyList = computed({
  get: () => proxiesStore.list,
  set: (v) => { proxiesStore.list = v },
})

const { dragIndex, dragOverIndex, onDragStart, onDragOver, onDragLeave, onDrop, onDragEnd } = useDraggable(proxyList, {
  onReorder: async (items) => {
    const config = await api.getConfig()
    config.Proxy_Service = items.map(p => config.Proxy_Service.find(c => c.ExternalDomain === p.ExternalDomain)).filter(Boolean)
    await api.saveConfig(JSON.stringify(config))
  },
})

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
            :class="{ 'drag-over': dragOverIndex === i, 'dragging': dragIndex === i }"
            @dragstart="onDragStart(i, $event)"
            @dragover="onDragOver(i, $event)"
            @dragleave="onDragLeave"
            @drop="onDrop(i)"
            @dragend="onDragEnd($event)"
          >
            <td>
              <i class="fas fa-grip-vertical drag-grip"></i>
              <span class="cert-icon" :class="findCertForDomain(proxy.ExternalDomain) ? (findCertForDomain(proxy.ExternalDomain).is_expired ? 'cert-expired' : 'cert-valid') : 'cert-none'">
                <i class="fas fa-shield-alt"></i>
              </span>
              {{ proxy.Name || '—' }}
            </td>
            <td>
              <code class="clickable-link" @click="openUrl('http://' + proxy.ExternalDomain)">{{ proxy.ExternalDomain }} <i class="fas fa-external-link-alt"></i></code>
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
