<script setup>
const { t } = useI18n()
const certsStore = useCertsStore()
const { success, error } = useNotification()

const props = defineProps({
  host: { type: String, required: true },
})

const certs = ref([])
const loading = ref(true)
const issuing = ref('')
const renewing = ref('')
const deleting = ref('')

const sleep = (ms) => new Promise(r => setTimeout(r, ms))

const refreshCerts = async () => {
  await certsStore.loadAll()
  certs.value = certsStore.list.filter(c =>
    c.dns_names?.some(d => d === props.host || d === `*.${props.host}` || props.host.match(new RegExp(d.replace('*.', '.*\\.'))))
  )
  if (certs.value.length === 0) {
    const info = await certsStore.getInfo(props.host)
    if (info && info.has_cert) certs.value = [info]
  }
}

onMounted(async () => {
  await refreshCerts()
  loading.value = false
})

const issueCert = async (domain) => {
  issuing.value = domain
  const [result] = await Promise.all([certsStore.issue(domain), sleep(1000)])
  if (result && !String(result).startsWith('Error')) {
    success(t('notify.certIssued'))
    await refreshCerts()
  } else {
    error(String(result))
  }
  issuing.value = ''
}

const renewCert = async (domain) => {
  renewing.value = domain
  const [result] = await Promise.all([certsStore.renew(domain), sleep(1000)])
  if (result && !String(result).startsWith('Error')) {
    success(t('notify.certRenewed'))
    await refreshCerts()
  } else {
    error(String(result))
  }
  renewing.value = ''
}

const deleteCert = async (domain) => {
  deleting.value = domain
  const [result] = await Promise.all([certsStore.remove(domain), sleep(1000)])
  if (result && !String(result).startsWith('Error')) {
    success(t('notify.certDeleted'))
    await refreshCerts()
  } else {
    error(String(result))
  }
  deleting.value = ''
}
</script>

<template>
  <div class="vaccess-page">
    <Breadcrumbs :items="[host, t('certs.title')]" />

    <PageHeader icon="fas fa-shield-alt" :title="t('certs.title')" :subtitle="`${t('certs.subtitle')} — ${host}`" />

    <div class="cert-manager-content">
      <!-- Нет сертификатов -->
      <div v-if="!loading && certs.length === 0" class="cert-empty">
        <i class="fas fa-certificate"></i>
        <h3>{{ t('certs.noCert') }}</h3>
        <p>{{ t('certs.subtitle') }}</p>
        <VButton icon="fas fa-plus" :loading="issuing === host" @click="issueCert(host)">{{ t('certs.issue') }}</VButton>
      </div>

      <!-- Карточки сертификатов -->
      <div v-for="cert in certs" :key="cert.domain" class="cert-card" :class="{ 'cert-card-empty': !cert.has_cert }">
        <div class="cert-card-header">
          <div class="cert-card-title" :class="{ expired: cert.is_expired }">
            <i class="fas fa-shield-alt"></i>
            <h3>{{ cert.domain }}</h3>
          </div>
          <div class="cert-card-actions">
            <VButton v-if="cert.has_cert && !cert.is_expired" variant="success" icon="fas fa-sync" :loading="renewing === cert.domain" @click="renewCert(cert.domain)">
              {{ t('certs.renew') }}
            </VButton>
            <VButton v-if="!cert.has_cert || cert.is_expired" icon="fas fa-plus" :loading="issuing === cert.domain" @click="issueCert(cert.domain)">
              {{ t('certs.issue') }}
            </VButton>
            <VButton v-if="cert.has_cert" variant="danger" icon="fas fa-trash" :loading="deleting === cert.domain" @click="deleteCert(cert.domain)">
              {{ t('certs.delete') }}
            </VButton>
          </div>
        </div>

        <div v-if="cert.has_cert" class="cert-info-grid">
          <div class="cert-info-item">
            <div class="cert-info-label">{{ t('certs.status') }}</div>
            <div class="cert-info-value" :class="cert.is_expired ? 'expired' : 'valid'">
              {{ cert.is_expired ? t('certs.expired') : `${t('certs.active')} (${t('certs.daysLeft', { n: cert.days_left })})` }}
            </div>
          </div>
          <div class="cert-info-item">
            <div class="cert-info-label">{{ t('certs.issuer') }}</div>
            <div class="cert-info-value">{{ cert.issuer }}</div>
          </div>
          <div class="cert-info-item">
            <div class="cert-info-label">{{ t('certs.issued') }}</div>
            <div class="cert-info-value">{{ cert.not_before }}</div>
          </div>
          <div class="cert-info-item">
            <div class="cert-info-label">{{ t('certs.expires') }}</div>
            <div class="cert-info-value" :class="cert.is_expired ? 'expired' : ''">{{ cert.not_after }}</div>
          </div>
        </div>

        <div v-if="cert.dns_names?.length" class="cert-domains-list">
          <span v-for="dns in cert.dns_names" :key="dns" class="cert-domain-tag">{{ dns }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
.cert-manager-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.cert-card {
  background: rgba(var(--accent-rgb), 0.03);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-xl);
  padding: var(--space-lg);
  transition: all var(--transition-base);
}

.cert-card-empty {
  border-style: dashed;
}

.cert-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-md);
  padding-bottom: var(--space-md);
  border-bottom: 1px solid var(--glass-border);
}

.cert-card-title {
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

.cert-card-title i {
  font-size: 24px;
  color: var(--accent-green);
}

.cert-card-title.expired i {
  color: var(--accent-red);
}

.cert-card-title h3 {
  margin: 0;
  font-size: var(--text-xl);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
}

.cert-card-actions {
  display: flex;
  gap: var(--space-sm);
}

.cert-info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--space-md);
}

.cert-info-item {
  padding: var(--space-md);
  background: var(--subtle-overlay);
  border-radius: var(--radius-md);
}

.cert-info-label {
  font-size: var(--text-sm);
  color: var(--text-muted);
  margin-bottom: var(--space-xs);
}

.cert-info-value {
  font-size: var(--text-md);
  color: var(--text-primary);
  font-weight: var(--font-medium);
}

.cert-info-value.valid {
  color: var(--accent-green);
}

.cert-info-value.expired {
  color: var(--accent-red);
}

.cert-domains-list {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-sm);
  margin-top: var(--space-md);
}

.cert-domain-tag {
  padding: 4px 12px;
  background: rgba(var(--accent-rgb), 0.15);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--accent-purple-light);
  font-family: var(--font-mono);
}

.cert-empty {
  text-align: center;
  padding: 60px 40px;
  color: var(--text-muted);
}

.cert-empty i {
  font-size: 48px;
  margin-bottom: var(--space-lg);
  opacity: 0.3;
}

.cert-empty h3 {
  font-size: var(--text-xl);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin-bottom: var(--space-sm);
}

.cert-empty p {
  font-size: var(--text-md);
  margin-bottom: var(--space-lg);
}
</style>
