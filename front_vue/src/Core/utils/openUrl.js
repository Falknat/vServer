export function openUrl(url) {
  if (window.runtime?.BrowserOpenURL) {
    window.runtime.BrowserOpenURL(url)
  } else {
    window.open(url, '_blank')
  }
}
