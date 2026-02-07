const isOpen = ref(false)
const title = ref('')
const content = ref(null)
const onSave = ref(null)
const onDelete = ref(null)

export function useModal() {
  const open = (options = {}) => {
    title.value = options.title || ''
    content.value = options.content || null
    onSave.value = options.onSave || null
    onDelete.value = options.onDelete || null
    isOpen.value = true
  }

  const close = () => {
    isOpen.value = false
    title.value = ''
    content.value = null
    onSave.value = null
    onDelete.value = null
  }

  const save = async () => {
    if (onSave.value) await onSave.value()
  }

  const remove = async () => {
    if (onDelete.value) await onDelete.value()
  }

  return {
    isOpen: readonly(isOpen),
    title: readonly(title),
    content: readonly(content),
    hasDelete: computed(() => !!onDelete.value),
    open,
    close,
    save,
    remove,
  }
}
