<script setup>
const { t } = useI18n()
const modal = useModal()
</script>

<template>
  <Teleport to="body">
    <transition name="fade">
      <div v-if="modal.isOpen.value" class="v-modal-overlay" @click.self="modal.close()">
        <transition name="scale">
          <div v-if="modal.isOpen.value" class="v-modal-window">
            <div class="v-modal-header">
              <h3 class="v-modal-title">{{ modal.title.value }}</h3>
              <button class="v-modal-close" @click="modal.close()">
                <i class="fas fa-times"></i>
              </button>
            </div>
            <div class="v-modal-content">
              <component v-if="modal.content.value" :is="modal.content.value" />
            </div>
            <div class="v-modal-footer">
              <VButton v-if="modal.hasDelete.value" variant="danger" icon="fas fa-trash" @click="modal.remove()">
                {{ t('common.delete') }}
              </VButton>
              <VButton @click="modal.close()">{{ t('common.cancel') }}</VButton>
              <VButton variant="primary" icon="fas fa-save" @click="modal.save()">
                {{ t('common.save') }}
              </VButton>
            </div>
          </div>
        </transition>
      </div>
    </transition>
  </Teleport>
</template>

<style scoped>
.v-modal-overlay {
  position: fixed;
  inset: 0;
  background: var(--overlay-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-modal);
  backdrop-filter: blur(4px);
}

.v-modal-window {
  background: var(--bg-secondary);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius);
  width: 90%;
  max-width: 700px;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow-lg);
}

.v-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-md) var(--space-lg);
  border-bottom: 1px solid var(--glass-border);
}

.v-modal-title {
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
}

.v-modal-close {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: var(--radius);
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.v-modal-close:hover {
  color: var(--accent-red);
  background: rgba(var(--danger-rgb), 0.1);
}

.v-modal-content {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-lg);
}

.v-modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-sm);
  padding: var(--space-md) var(--space-lg);
  border-top: 1px solid var(--glass-border);
}

.v-modal-footer .v-btn--danger {
  margin-right: auto;
}
</style>
