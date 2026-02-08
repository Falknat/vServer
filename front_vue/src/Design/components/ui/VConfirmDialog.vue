<script setup>

const { isOpen, config, close } = useConfirm()
</script>

<template>
  <Teleport to="body">
    <transition name="fade">
      <div v-if="isOpen" class="confirm-overlay" @click.self="close()">
        <transition name="scale">
          <div v-if="isOpen" class="confirm-dialog">
            <div class="confirm-icon" :style="{ color: config.iconColor }">
              <i :class="config.icon"></i>
            </div>
            <h3 class="confirm-title">{{ config.title }}</h3>
            <p v-if="config.message" class="confirm-message">{{ config.message }}</p>
            <div class="confirm-buttons">
              <button
                v-for="(btn, i) in config.buttons"
                :key="i"
                class="confirm-btn"
                :class="`confirm-btn--${btn.variant || 'default'}`"
                @click="btn.action?.()"
              >
                <i v-if="btn.icon" :class="btn.icon"></i>
                {{ btn.label }}
              </button>
            </div>
          </div>
        </transition>
      </div>
    </transition>
  </Teleport>
</template>

<style scoped>
.confirm-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-modal);
  backdrop-filter: blur(6px);
}

.confirm-dialog {
  background: transparent;
  border: none;
  padding: 40px;
  max-width: 420px;
  width: 90%;
  text-align: center;
}

.confirm-icon {
  font-size: 42px;
  margin-bottom: 16px;
  filter: drop-shadow(0 0 16px currentColor);
}

.confirm-title {
  font-size: var(--text-xl);
  font-weight: var(--font-bold);
  color: var(--text-primary);
  margin: 0 0 8px 0;
}

.confirm-message {
  font-size: var(--text-md);
  color: var(--text-muted);
  margin: 0 0 28px 0;
  line-height: 1.5;
}

.confirm-buttons {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.confirm-btn {
  width: 100%;
  padding: 12px 20px;
  border-radius: var(--radius);
  font-size: var(--text-md);
  font-weight: var(--font-semibold);
  cursor: pointer;
  transition: all var(--transition-base);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border: 1px solid transparent;
}

.confirm-btn--default {
  background: rgba(var(--muted-rgb), 0.15);
  border-color: rgba(var(--muted-rgb), 0.3);
  color: var(--text-secondary);
}

.confirm-btn--default:hover {
  background: rgba(var(--muted-rgb), 0.25);
}

.confirm-btn--primary {
  background: rgba(var(--accent-rgb), 0.2);
  border-color: rgba(var(--accent-rgb), 0.4);
  color: var(--accent-purple-light);
}

.confirm-btn--primary:hover {
  background: rgba(var(--accent-rgb), 0.3);
}

.confirm-btn--danger {
  background: rgba(var(--danger-rgb), 0.2);
  border-color: rgba(var(--danger-rgb), 0.4);
  color: var(--accent-red);
}

.confirm-btn--danger:hover {
  background: rgba(var(--danger-rgb), 0.3);
}

.confirm-btn--success {
  background: rgba(var(--success-rgb), 0.2);
  border-color: rgba(var(--success-rgb), 0.4);
  color: var(--accent-green);
}

.confirm-btn--success:hover {
  background: rgba(var(--success-rgb), 0.3);
}

.scale-enter-active { transition: all 0.25s cubic-bezier(0.34, 1.56, 0.64, 1); }
.scale-leave-active { transition: all 0.15s ease-in; }
.scale-enter-from { opacity: 0; transform: scale(0.9); }
.scale-leave-to { opacity: 0; transform: scale(0.95); }
</style>
