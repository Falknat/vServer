<script setup>
const { t } = useI18n()
const router = useRouter()
const proxiesStore = useProxiesStore()
const certsStore = useCertsStore()

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
    <div class="section-header">
      <h2 class="section-title">{{ t('proxies.title') }}</h2>
      <VButton variant="primary" icon="fas fa-plus" @click="router.push('/proxies/create')">
        {{ t('proxies.add') }}
      </VButton>
    </div>
    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>{{ t('sites.name') }}</th>
            <th>{{ t('proxies.externalDomain') }}</th>
            <th>{{ t('proxies.localAddress') }}</th>
            <th>HTTPS</th>
            <th>{{ t('proxies.status') }}</th>
            <th>{{ t('proxies.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="proxy in proxiesStore.list" :key="proxy.ExternalDomain">
            <td>
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
              <VBadge :variant="proxy.Enable ? 'online' : 'offline'">
                {{ proxy.Enable ? 'active' : 'disabled' }}
              </VBadge>
            </td>
            <td>
              <button class="icon-btn" :title="t('sites.editVaccess')" @click="router.push(`/vaccess/${proxy.ExternalDomain}`)">
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
