export function useDraggable(list, { onReorder } = {}) {
  const dragIndex = ref(null)
  const dragOverIndex = ref(null)

  const onDragStart = (index, event) => {
    dragIndex.value = index
    event.dataTransfer.effectAllowed = 'move'
    event.target.closest('tr, .draggable-item')?.classList.add('dragging')
  }

  const onDragOver = (index, event) => {
    event.preventDefault()
    event.dataTransfer.dropEffect = 'move'
    dragOverIndex.value = index
  }

  const onDragEnter = (index, event) => {
    event.preventDefault()
    dragOverIndex.value = index
  }

  const onDragLeave = () => {
    dragOverIndex.value = null
  }

  const onDrop = async (index) => {
    if (dragIndex.value === null || dragIndex.value === index) {
      dragIndex.value = null
      dragOverIndex.value = null
      return
    }

    const items = [...list.value]
    const [moved] = items.splice(dragIndex.value, 1)
    items.splice(index, 0, moved)
    list.value = items

    dragIndex.value = null
    dragOverIndex.value = null

    if (onReorder) await onReorder(items)
  }

  const onDragEnd = (event) => {
    event.target.closest('tr, .draggable-item')?.classList.remove('dragging')
    dragIndex.value = null
    dragOverIndex.value = null
  }

  return {
    dragIndex,
    dragOverIndex,
    onDragStart,
    onDragOver,
    onDragEnter,
    onDragLeave,
    onDrop,
    onDragEnd,
  }
}
