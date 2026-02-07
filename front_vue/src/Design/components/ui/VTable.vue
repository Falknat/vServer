<script setup>
defineProps({
  columns: { type: Array, default: () => [] },
  data: { type: Array, default: () => [] },
})
</script>

<template>
  <div class="v-table-container">
    <table class="v-table">
      <thead>
        <tr>
          <th v-for="col in columns" :key="col.key">{{ col.label }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(row, index) in data" :key="index">
          <slot name="row" :row="row" :index="index" />
        </tr>
        <tr v-if="data.length === 0">
          <td :colspan="columns.length" class="v-table-empty">
            <slot name="empty">{{ $t('common.loading') }}</slot>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.v-table-container {
  overflow-x: auto;
  border-radius: var(--radius);
  border: 1px solid var(--glass-border);
}

.v-table {
  width: 100%;
  border-collapse: collapse;
}

.v-table th {
  padding: var(--space-sm) var(--space-md);
  text-align: left;
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  background: var(--glass-bg-dark);
  border-bottom: 1px solid var(--glass-border);
}

.v-table td {
  padding: var(--space-sm) var(--space-md);
  font-size: var(--text-md);
  color: var(--text-primary);
  border-bottom: 1px solid var(--glass-border);
}

.v-table tbody tr {
  transition: background var(--transition-fast);
}

.v-table tbody tr:hover {
  background: var(--glass-bg-light);
}

.v-table tbody tr:last-child td {
  border-bottom: none;
}

.v-table-empty {
  text-align: center;
  padding: var(--space-xl) !important;
  color: var(--text-muted);
}
</style>
