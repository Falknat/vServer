<script setup>
const { t } = useI18n()
const servicesStore = useServicesStore()
</script>

<template>
  <section class="section">
    <div class="services-row">
      <ServiceCard
        v-for="service in servicesStore.list"
        :key="service.name"
        :service="service"
      />
    </div>
  </section>
</template>

<style>
.services-row {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.services-row > * {
  flex: 1;
  min-width: 140px;
}
</style>
