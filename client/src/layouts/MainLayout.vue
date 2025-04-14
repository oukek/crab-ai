<template>
  <q-layout view="lHh Lpr lFf">
    <!-- 添加 header 作为标签页导航 -->
    <q-header class="header-nav">
      <div class="tabs-container">
        <div
          v-for="tab in tabs"
          :key="tab.path"
          class="tab"
          :class="{ active: currentTab === tab.path }"
          @click="switchTab(tab.path)"
        >
          <div class="tab-content">
            <q-icon :name="tab.icon" size="20px" />
            <span class="tab-label">{{ tab.label }}</span>
            <button 
              class="close-button"
              @click.stop="closeTab(tab.path)"
              v-show="tabs.length > 1"
            >
              <q-icon name="close" size="16px" />
            </button>
          </div>
        </div>
      </div>
    </q-header>

    <q-drawer
      show-if-above
      bordered
      mini
    >
      <SideNav @open-tab="openNewTab"/>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import SideNav from 'components/SideNav.vue';
import { ref, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';

const router = useRouter();
const route = useRoute();

interface Tab {
  path: string;
  icon: string;
  label: string;
}

const currentTab = ref('/');
const tabs = ref<Tab[]>([
  { path: '/', icon: 'home', label: '首页' }
]);

// 切换标签页
const switchTab = (path: string) => {
  currentTab.value = path;
  router.push(path);
};

// 关闭标签页
const closeTab = (path: string) => {
  const index = tabs.value.findIndex(tab => tab.path === path);
  if (index === -1) return;
  
  // 如果关闭的是当前标签，需要切换到其他标签
  if (currentTab.value === path) {
    const newIndex = index === 0 ? 1 : index - 1;
    switchTab(tabs.value[newIndex].path);
  }
  
  tabs.value.splice(index, 1);
};

// 打开新标签页
const openNewTab = (path: string, icon: string, label: string) => {
  // 如果标签页已存在，直接切换
  const existingTab = tabs.value.find(tab => tab.path === path);
  if (existingTab) {
    switchTab(path);
    return;
  }

  // 添加新标签页
  tabs.value.push({ path, icon, label });
  switchTab(path);
};

// 监听路由变化，自动添加标签页
watch(() => route.path, (newPath) => {
  const existingTab = tabs.value.find(tab => tab.path === newPath);
  if (!existingTab) {
    let icon = 'article';
    let label = '新标签页';
    
    // 根据路径设置图标和标签
    if (newPath === '/') {
      icon = 'home';
      label = '首页';
    } else if (newPath === '/books') {
      icon = 'library_books';
      label = '书籍';
    } else if (newPath === '/settings') {
      icon = 'settings';
      label = '设置';
    }
    
    openNewTab(newPath, icon, label);
  }
  currentTab.value = newPath;
}, { immediate: true });
</script>

<style scoped>
.header-nav {
  background: rgba(255, 255, 255, 0.8) !important;
  border-bottom: 1px solid #e0e0e0;
}

.tabs-container {
  display: flex;
  height: 40px;
  padding: 0 16px;
  gap: 8px;
}

.tab {
  position: relative;
  min-width: 140px;
  height: 32px;
  margin-top: 8px;
  border-radius: 6px 6px 0 0;
  background: rgba(240, 240, 240, 0.6);
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab:hover {
  background: rgba(255, 255, 255, 0.9);
}

.tab.active {
  background: white;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.tab-content {
  display: flex;
  align-items: center;
  padding: 0 12px;
  height: 100%;
  color: #666;
  gap: 8px;
}

.tab.active .tab-content {
  color: #333;
}

.tab-label {
  font-size: 13px;
  font-weight: 500;
}

.close-button {
  opacity: 0;
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  border: none;
  background: transparent;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #666;
  padding: 0;
  transition: all 0.2s ease;
}

.tab:hover .close-button {
  opacity: 1;
}

.close-button:hover {
  background: rgba(0, 0, 0, 0.1);
}
</style>
