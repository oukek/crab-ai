<template>
  <q-page class="settings-page q-pa-md">
    <div class="settings-container">
      <!-- 左侧模型服务列表 -->
      <div class="model-services">
        <q-list padding>
          <q-item-label header>模型服务</q-item-label>
          <q-item clickable v-ripple :active="currentSection === 'model-service'" @click="currentSection = 'model-service'">
            <q-item-section avatar>
              <q-icon name="cloud" />
            </q-item-section>
            <q-item-section>模型服务</q-item-section>
          </q-item>
          
          <q-item clickable v-ripple :active="currentSection === 'default-model'" @click="currentSection = 'default-model'">
            <q-item-section avatar>
              <q-icon name="widgets" />
            </q-item-section>
            <q-item-section>默认模型</q-item-section>
          </q-item>

          <q-item clickable v-ripple :active="currentSection === 'web-search'" @click="currentSection = 'web-search'">
            <q-item-section avatar>
              <q-icon name="search" />
            </q-item-section>
            <q-item-section>网络搜索</q-item-section>
          </q-item>

          <q-item clickable v-ripple :active="currentSection === 'mcp-server'" @click="currentSection = 'mcp-server'">
            <q-item-section avatar>
              <q-icon name="computer" />
            </q-item-section>
            <q-item-section>MCP 服务器</q-item-section>
          </q-item>

          <q-separator spaced />

          <q-item clickable v-ripple :active="currentSection === 'about'" @click="currentSection = 'about'">
            <q-item-section avatar>
              <q-icon name="info" />
            </q-item-section>
            <q-item-section>关于</q-item-section>
          </q-item>
        </q-list>
      </div>

      <!-- 右侧设置内容区域 -->
      <div class="settings-content">
        <ModelService v-if="currentSection === 'model-service'" />
        <About v-if="currentSection === 'about'" />
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, defineComponent } from 'vue'
import ModelService from './components/ModelService.vue'
import About from './components/About.vue'

defineComponent({
  name: 'SettingsPage'
})

// 状态管理
const modelEnabled = ref<boolean>(true)
const apiKey = ref<string>('')
const apiUrl = ref<string>('https://api.siliconflow.cn')
const currentSection = ref<string>('model-service')
</script>

<style scoped>
.settings-page {
  height: 100vh;
}

.settings-container {
  display: grid;
  grid-template-columns: 240px 1fr;
  gap: 20px;
  height: 100%;
}

.model-services {
  background: #fff;
  border-radius: 8px;
  border: 1px solid #eee;
  height: 100%;
}

.settings-content {
  background: #fff;
  border-radius: 8px;
  border: 1px solid #eee;
  height: 100%;
}

.api-section {
  max-width: 600px;
}
</style> 