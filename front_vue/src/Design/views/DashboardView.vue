<script setup>
const servicesStore = useServicesStore()
const sitesStore = useSitesStore()
const proxiesStore = useProxiesStore()
const certsStore = useCertsStore()

onMounted(async () => {
  await Promise.all([
    sitesStore.load(),
    proxiesStore.load(),
    certsStore.loadAll(),
  ])
  if (!servicesStore.loaded) {
    await servicesStore.load()
  }
})
</script>

<template>
  <div class="dashboard-view">
    <ServicesGrid />
    <div class="dashboard-tables">
      <SitesTable />
      <ProxiesTable />
    </div>
  </div>
</template>

<style scoped>
.dashboard-view {
  display: flex;
  flex-direction: column;
  gap: var(--space-xl);
}

.dashboard-tables {
  display: flex;
  flex-direction: column;
  gap: var(--space-xl);
}

@media (min-width: 2200px) {
  .dashboard-tables {
    flex-direction: row;
    gap: var(--space-xl);
  }

  .dashboard-tables > * {
    flex: 1;
    min-width: 0;
  }
}
</style>
