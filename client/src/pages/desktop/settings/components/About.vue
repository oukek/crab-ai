<template>
  <div class="about-container q-pa-md">
    <!-- Logo和版本信息区域 -->
    <div class="app-info q-mb-xl">
      <q-img
        src="/app-logo.png"
        spinner-color="primary"
        style="height: 80px; width: 80px"
        class="app-logo q-mb-md"
      />
      <div class="text-h5 text-weight-bold q-mb-sm">Crab AI</div>
      <div class="version-tag">Version 1.0.0</div>
    </div>

    <!-- 信息列表区域 -->
    <div class="info-list">
      <div class="info-item" v-for="(item, index) in infoItems" :key="index">
        <div class="info-label">{{ item.label }}</div>
        <div class="info-content">
          <a 
            v-if="item.isLink" 
            :href="item.content" 
            target="_blank" 
            class="link-text"
          >
            {{ item.content }}
          </a>
          <span v-else>{{ item.content }}</span>
        </div>
      </div>
    </div>

    <!-- 版权信息 -->
    <div class="copyright text-grey-7 q-mt-xl">
      © 2024 Crab AI. All rights reserved.
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const infoItems = ref([
  {
    label: '官方网站',
    content: 'https://crab-ai.com',
    isLink: true
  },
  {
    label: '技术支持',
    content: 'support@crab-ai.com',
    isLink: false
  },
  {
    label: '联系我们',
    content: 'contact@crab-ai.com',
    isLink: false
  },
  {
    label: '开源协议',
    content: 'MIT License',
    isLink: false
  }
])
</script>

<style scoped>
.about-container {
  max-width: 600px;
  margin: 0 auto;
  padding-top: 48px;
}

.app-info {
  text-align: center;
}

.app-logo {
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.version-tag {
  display: inline-block;
  padding: 4px 12px;
  background: #F5F5F5;
  border-radius: 20px;
  font-size: 14px;
  color: #666;
}

.info-list {
  background: #FAFAFA;
  border-radius: 12px;
  padding: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #EEEEEE;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  width: 100px;
  color: #666;
  font-size: 14px;
}

.info-content {
  flex: 1;
  font-size: 14px;
}

.link-text {
  color: #1976D2;
  text-decoration: none;
  transition: color 0.3s;
}

.link-text:hover {
  color: #1565C0;
}

.copyright {
  text-align: center;
  font-size: 12px;
}
</style> 