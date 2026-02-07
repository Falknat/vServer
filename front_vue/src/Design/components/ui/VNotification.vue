<script setup>
const { notifications, remove } = useNotification()
</script>

<template>
  <Teleport to="body">
    <div class="v-notifications">
      <TransitionGroup name="notification">
        <div
          v-for="n in notifications"
          :key="n.id"
          class="v-notify"
          :class="`v-notify--${n.type}`"
          @click="remove(n.id)"
        >
          <i :class="{
            'fas fa-check-circle': n.type === 'success',
            'fas fa-exclamation-circle': n.type === 'error',
            'fas fa-info-circle': n.type === 'info',
          }"></i>
          <span>{{ n.message }}</span>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.v-notifications {
  position: fixed;
  top: calc(var(--header-height) + var(--space-md));
  right: var(--space-md);
  z-index: var(--z-notification);
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  max-width: 400px;
}

.v-notify {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius);
  font-size: var(--text-md);
  cursor: pointer;
  backdrop-filter: var(--backdrop-blur);
  border: 1px solid var(--glass-border);
}

.v-notify--success {
  background: rgba(var(--success-rgb), 0.15);
  color: var(--accent-green);
  border-color: rgba(var(--success-rgb), 0.3);
}

.v-notify--error {
  background: rgba(var(--danger-rgb), 0.15);
  color: var(--accent-red);
  border-color: rgba(var(--danger-rgb), 0.3);
}

.v-notify--info {
  background: rgba(var(--info-rgb), 0.15);
  color: var(--accent-cyan);
  border-color: rgba(var(--info-rgb), 0.3);
}
</style>
