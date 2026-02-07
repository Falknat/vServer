<script setup>
const servicesStore = useServicesStore()
const sitesStore = useSitesStore()
const proxiesStore = useProxiesStore()
const certsStore = useCertsStore()

onMounted(async () => {
  await Promise.all([
    servicesStore.load(),
    sitesStore.load(),
    proxiesStore.load(),
    certsStore.loadAll(),
  ])
})
</script>

<template>
  <div class="dashboard-view">
    <ServicesGrid />
    <SitesTable />
    <ProxiesTable />
  </div>
</template>

<style scoped>
.dashboard-view {
  display: flex;
  flex-direction: column;
  gap: var(--space-xl);
}
</style>
