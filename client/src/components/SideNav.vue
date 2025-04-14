<template>
  <nav class="side-nav">
    <a 
      class="nav-item"
      :class="{ 'nav-active': $route.path === '/' }"
      @click.prevent="$emit('open-tab', '/', 'home', '首页')"
    >
      <span class="icon">
        <q-icon name="home" size="sm"/>
      </span>
    </a>

    <a 
      class="nav-item"
      :class="{ 'nav-active': $route.path === '/books' }"
      @click.prevent="$emit('open-tab', '/books', 'library_books', '书籍')"
    >
      <span class="icon">
        <q-icon name="library_books" size="sm"/>
      </span>
    </a>

    <a 
      class="nav-item settings-link"
      :class="{ 'nav-active': $route.path === '/settings' }"
      @click.prevent="$emit('open-tab', '/settings', 'settings', '设置')"
    >
      <span class="icon">
        <q-icon name="settings" size="sm"/>
      </span>
    </a>
  </nav>
</template>

<script setup lang="ts">
// 定义组件的事件
defineEmits<{
  (e: 'open-tab', path: string, icon: string, label: string): void
}>();
</script>

<style scoped>
.side-nav {
  padding-top: 20px;
  padding-bottom: 20px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  height: 100vh;
}

.side-nav > :not(.settings-link) {
  flex-shrink: 0;
}

.nav-item {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 10px;
  margin: 0 8px;
  border-radius: 8px;
  color: #666;
  text-decoration: none;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.nav-item:hover {
  background: rgba(0, 0, 0, 0.06);
  color: #333;
  transform: translateY(-1px);
}

.nav-active {
  background: rgba(60, 60, 60, 0.08);
  color: #000;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 6px;
}

.label {
  font-size: 14px;
  font-weight: 500;
}

.settings-link {
  margin-top: auto;
  margin-bottom: 0;
}
</style> 