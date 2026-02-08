export function useCertLookup() {
  const certsStore = useCertsStore()

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

  return { findCertForDomain }
}
