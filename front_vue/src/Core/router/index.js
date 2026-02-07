import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'dashboard',
    component: () => import('@design/views/DashboardView.vue'),
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('@design/views/SettingsView.vue'),
  },
  {
    path: '/sites/create',
    name: 'site-create',
    component: () => import('@design/views/SiteCreateView.vue'),
  },
  {
    path: '/proxies/create',
    name: 'proxy-create',
    component: () => import('@design/views/ProxyCreateView.vue'),
  },
  {
    path: '/sites/edit/:host',
    name: 'site-edit',
    component: () => import('@design/views/SiteEditView.vue'),
    props: true,
  },
  {
    path: '/proxies/edit/:domain',
    name: 'proxy-edit',
    component: () => import('@design/views/ProxyEditView.vue'),
    props: true,
  },
  {
    path: '/vaccess/:host',
    name: 'vaccess',
    component: () => import('@design/views/VAccessView.vue'),
    props: true,
  },
  {
    path: '/certs/:host',
    name: 'certs',
    component: () => import('@design/views/CertManagerView.vue'),
    props: true,
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
