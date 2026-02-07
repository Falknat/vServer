const notifications = ref([])
let nextId = 0

export function useNotification() {
  const show = (message, type = 'info', duration = 2000) => {
    const id = nextId++
    notifications.value.push({ id, message, type })

    if (duration > 0) {
      setTimeout(() => {
        remove(id)
      }, duration)
    }

    return id
  }

  const remove = (id) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index !== -1) {
      notifications.value.splice(index, 1)
    }
  }

  const success = (message, duration = 2000) => show(message, 'success', duration)
  const error = (message, duration = 3000) => show(message, 'error', duration)
  const info = (message, duration = 2000) => show(message, 'info', duration)

  return {
    notifications: readonly(notifications),
    show,
    success,
    error,
    info,
    remove,
  }
}
