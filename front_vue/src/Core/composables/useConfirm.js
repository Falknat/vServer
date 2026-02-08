const isOpen = ref(false)
const config = ref({
  icon: 'fas fa-question-circle',
  iconColor: 'var(--accent-red)',
  title: '',
  message: '',
  buttons: [],
})
let resolvePromise = null

export function useConfirm() {
  const open = (options = {}) => {
    config.value = {
      icon: options.icon || 'fas fa-exclamation-triangle',
      iconColor: options.iconColor || 'var(--accent-red)',
      title: options.title || '',
      message: options.message || '',
      buttons: options.buttons || [],
    }
    isOpen.value = true
  }

  const close = () => {
    isOpen.value = false
    if (resolvePromise) {
      resolvePromise(false)
      resolvePromise = null
    }
  }

  const confirm = (options = {}) => {
    return new Promise((resolve) => {
      resolvePromise = resolve
      open({
        ...options,
        buttons: (options.buttons || []).map(btn => ({
          ...btn,
          action: () => {
            resolvePromise = null
            isOpen.value = false
            resolve(btn.value !== undefined ? btn.value : btn.label)
            if (btn.action) btn.action()
          },
        })),
      })
    })
  }

  const confirmDelete = (options = {}) => {
    return confirm({
      icon: options.icon || 'fas fa-trash-alt',
      iconColor: options.iconColor || 'var(--accent-red)',
      title: options.title || '',
      message: options.message || '',
      buttons: [
        { label: options.cancelText || 'Отмена', variant: 'default', value: false },
        { label: options.confirmText || 'Удалить', variant: 'danger', icon: 'fas fa-trash', value: true },
      ],
    })
  }

  return {
    isOpen: readonly(isOpen),
    config: readonly(config),
    open,
    close,
    confirm,
    confirmDelete,
  }
}
