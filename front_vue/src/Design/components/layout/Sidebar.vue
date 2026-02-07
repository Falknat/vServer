<script setup>
const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const navItems = [
  { name: 'dashboard', icon: 'fas fa-home', route: '/' },
  { name: 'settings', icon: 'fas fa-cog', route: '/settings' },
]

const isActive = (item) => {
  if (item.route === '/') return route.path === '/'
  return route.path.startsWith(item.route)
}

const navigate = (item) => {
  router.push(item.route)
}
</script>

<template>
  <aside class="sidebar">
    <nav class="sidebar-nav">
      <button
        v-for="item in navItems"
        :key="item.name"
        class="nav-item"
        :class="{ active: isActive(item) }"
        :title="t(`nav.${item.name}`)"
        @click="navigate(item)"
      >
        <i :class="item.icon"></i>
      </button>
    </nav>
  </aside>
</template>

<style scoped>
.sidebar {
  width: var(--sidebar-width);
  background: var(--sidebar-bg);
  border-right: 1px solid var(--sidebar-border);
  display: flex;
  flex-direction: column;
  padding: 20px 0;
}

.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  padding: 0 var(--space-md);
}

.nav-item {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: var(--radius-lg);
  color: var(--nav-color);
  font-size: 20px;
  cursor: pointer;
  transition: all var(--transition-base);
  position: relative;
}

.nav-item::before {
  content: '';
  position: absolute;
  left: -16px;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 0;
  background: var(--nav-indicator);
  border-radius: 0 2px 2px 0;
  transition: height var(--transition-base);
}

.nav-item:hover {
  background: var(--nav-hover-bg);
  color: var(--nav-hover-color);
}

.nav-item.active {
  background: var(--nav-active-bg);
  color: var(--nav-active-color);
}

.nav-item.active::before {
  height: 24px;
}
</style>
