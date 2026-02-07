<script setup>
const { t } = useI18n()
const servicesStore = useServicesStore()
</script>

<template>
  <section class="section">
    <div class="services-grid">
      <ServiceCard
        v-for="service in servicesStore.list"
        :key="service.name"
        :service="service"
      />
    </div>
  </section>
</template>

<style>
.services-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: var(--space-lg);
}

@media (max-width: 1200px) {
  .services-grid { grid-template-columns: repeat(3, 1fr); }
}

@media (max-width: 900px) {
  .services-grid { grid-template-columns: repeat(2, 1fr); }
}

@media (max-width: 600px) {
  .services-grid { grid-template-columns: 1fr; }
}
</style>
