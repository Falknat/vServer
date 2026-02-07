export function useWailsEvents() {
  const isWails = typeof window !== 'undefined' && window?.runtime?.EventsOn

  const on = (eventName, callback) => {
    if (isWails) {
      window.runtime.EventsOn(eventName, callback)
    }
  }

  const off = (eventName) => {
    if (isWails && window.runtime.EventsOff) {
      window.runtime.EventsOff(eventName)
    }
  }

  return { on, off, isWails }
}
