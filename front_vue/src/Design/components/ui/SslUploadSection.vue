<script setup>
const { t } = useI18n()

const props = defineProps({
  modelValue: { type: String, default: 'none' },
})

const emit = defineEmits(['update:modelValue', 'cert-file', 'key-file', 'ca-file'])

const certModeOptions = computed(() => [
  { value: 'none', label: t('certs.modeNone') },
  { value: 'auto', label: t('certs.modeAuto') },
  { value: 'upload', label: t('certs.modeUpload') },
])
</script>

<template>
  <div class="ssl-section">
    <h3 class="form-subsection-title">
      <i class="fas fa-lock"></i> {{ t('common.sslCerts') }}
    </h3>

    <VSelect
      :label="t('certs.mode')"
      :model-value="modelValue"
      :options="certModeOptions"
      @update:model-value="emit('update:modelValue', $event)"
    />

    <template v-if="modelValue === 'upload'">
      <VFileUpload
        :label="t('certs.certFile')"
        accept=".crt,.pem"
        required
        @select="emit('cert-file', $event)"
      />
      <VFileUpload
        :label="t('certs.keyFile')"
        accept=".key,.pem"
        required
        @select="emit('key-file', $event)"
      />
      <VFileUpload
        :label="t('certs.caFile')"
        accept=".crt,.pem"
        :hint="t('certs.caHint')"
        @select="emit('ca-file', $event)"
      />
    </template>
  </div>
</template>

<style>
.ssl-section {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}
</style>
