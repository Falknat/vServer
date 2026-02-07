import { mockApi } from './mock.js'
import { wailsApi } from './wails.js'

const isWails = typeof window !== 'undefined' && window?.go?.admin?.App

export const api = isWails ? wailsApi : mockApi
