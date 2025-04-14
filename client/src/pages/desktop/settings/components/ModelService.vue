<template>
  <div class="model-service-container">
    <!-- 左侧服务商列表 -->
    <div class="providers-list">
      <q-list padding>
        <q-item-label header>服务商</q-item-label>
        <q-item 
          v-for="provider in providers" 
          :key="provider.id"
          clickable 
          v-ripple
          :active="selectedProvider === provider.id"
          @click="selectedProvider = provider.id"
        >
          <q-item-section avatar>
            <q-avatar size="sm">
              <img :src="provider.icon" :alt="provider.name">
            </q-avatar>
          </q-item-section>
          <q-item-section>{{ provider.name }}</q-item-section>
        </q-item>
      </q-list>
    </div>

    <!-- 右侧配置区域 -->
    <div class="provider-config q-pa-md">
      <template v-if="currentProvider">
        <div class="row items-center q-mb-md">
          <div class="text-h6">{{ currentProvider.name }}</div>
          <q-toggle v-model="providerEnabled" class="q-ml-sm" />
          <q-btn flat round icon="open_in_new" size="sm" class="q-ml-auto" />
        </div>

        <!-- API密钥设置 -->
        <div class="api-section q-mb-lg">
          <div class="text-subtitle1 q-mb-sm">API 密钥</div>
          <q-input 
            v-model="apiKey"
            outlined
            type="password"
            :placeholder="`请在这里填写${currentProvider.name}的API密钥`"
          >
            <template v-slot:append>
              <q-btn flat round icon="content_copy" size="sm" />
              <q-btn flat round icon="visibility" size="sm" />
            </template>
          </q-input>
          <div class="text-caption text-grey">
            <q-icon name="info" size="xs" /> 点击这里获取密钥 充值
          </div>
        </div>

        <!-- API地址设置 -->
        <div class="api-section">
          <div class="text-subtitle1 q-mb-sm">API 地址</div>
          <q-input 
            v-model="apiUrl"
            outlined
            :placeholder="currentProvider.defaultApiUrl"
          >
            <template v-slot:append>
              <q-btn flat round icon="refresh" size="sm" />
              <q-btn flat round icon="content_copy" size="sm" />
            </template>
          </q-input>
          <div class="text-caption text-grey q-mt-sm">
            {{ currentProvider.apiEndpoint }}
          </div>
        </div>

        <!-- 可用模型列表 -->
        <div class="models-section q-mt-lg">
          <div class="text-subtitle1 q-mb-md">可用模型</div>
          <q-expansion-item
            v-for="model in currentProvider.models"
            :key="model.id"
            expand-separator
            icon="memory"
            :label="model.name"
            default-opened
          >
            <q-card>
              <q-card-section>
                <div class="row items-center">
                  <q-avatar size="sm" class="q-mr-sm">
                    <img :src="model.icon" :alt="model.name">
                  </q-avatar>
                  <div>{{ model.fullName }}</div>
                  <q-chip v-if="model.isNew" color="warning" text-color="white" class="q-ml-sm">
                    新入
                  </q-chip>
                  <q-btn flat round icon="settings" size="sm" class="q-ml-auto" />
                </div>
              </q-card-section>
            </q-card>
          </q-expansion-item>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

// 模拟数据，实际使用时可以通过props传入或从API获取
const providers = ref([
  {
    id: 'google-ai-studio',
    name: 'Google AI Studio',
    icon: '/path_to_icon',
    defaultApiUrl: 'https://api.google-ai-studio.cn',
    apiEndpoint: 'https://api.google-ai-studio.cn/v1/chat/completions',
    models: [
      {
        id: 'bge-m3',
        name: 'BAAI',
        fullName: 'BAAI/bge-m3',
        icon: '/path_to_baai_icon',
        isNew: true
      }
    ]
  },
  // 可以添加更多服务商
])

const selectedProvider = ref('google-ai-studio')
const providerEnabled = ref(true)
const apiKey = ref('')
const apiUrl = ref('')

const currentProvider = computed(() => 
  providers.value.find(p => p.id === selectedProvider.value)
)
</script>

<style scoped>
.model-service-container {
  display: grid;
  grid-template-columns: 240px 1fr;
  gap: 20px;
  height: 100%;
}

.providers-list {
  background: #fff;
  border-radius: 8px;
  border: 1px solid #eee;
}

.provider-config {
  background: #fff;
  border-radius: 8px;
  border: 1px solid #eee;
}

.api-section {
  max-width: 600px;
}
</style> 