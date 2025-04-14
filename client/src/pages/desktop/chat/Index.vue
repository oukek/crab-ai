<template>
  <div class="chat-container">
    <!-- 顶部工具栏 -->
    <div class="toolbar">
      <button class="new-chat-btn">新的会话</button>
      <div class="toolbar-right">
        <button class="history-btn">历史记录</button>
      </div>
    </div>

    <!-- 中间内容区域 -->
    <div class="content-area">
      <ChatMessages :messages="messages" />
    </div>

    <!-- 底部输入区域 -->
    <div class="input-section">
      <div class="context-area">context区</div>
      <div class="input-container">
        <textarea 
          v-model="inputMessage" 
          placeholder="请输入内容..."
          @keydown.enter.prevent="sendMessage"
        ></textarea>
      </div>
      <div class="bottom-controls">
        <div class="left-controls">
          <q-select
            class="mode-select"
            outlined
            dense
            v-model="selectedMode"
            :options="modeOptions"
            label="模式选择"
          />
          <q-select
            class="model-select"
            outlined
            dense
            v-model="selectedModel"
            :options="modelOptions"
            label="模型选择"
          />
        </div>
        <button class="send-btn" @click="sendMessage">发送</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import ChatMessages from './components/ChatMessages.vue';

const selectedMode = ref(null);
const selectedModel = ref(null);
const inputMessage = ref('');

const modeOptions = ref([
  '模式1',
  '模式2',
  '模式3'
]);

const modelOptions = ref([
  '模型1',
  '模型2',
  '模型3'
]);

interface Message {
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
}

const messages = ref<Message[]>([
  {
    role: 'assistant',
    content: '你好！我是AI助手，有什么我可以帮你的吗？',
    timestamp: new Date()
  }
]);

const sendMessage = () => {
  if (!inputMessage.value.trim()) return;
  
  // 添加用户消息
  messages.value.push({
    role: 'user',
    content: inputMessage.value,
    timestamp: new Date()
  });

  // 模拟AI回复
  setTimeout(() => {
    messages.value.push({
      role: 'assistant',
      content: '这是一个模拟的AI回复消息。',
      timestamp: new Date()
    });
  }, 1000);

  // 清空输入
  inputMessage.value = '';
};
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  padding: 1rem;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem;
  border-bottom: 1px solid #e5e5e5;
}

.toolbar-right {
  display: flex;
  gap: 1rem;
}

.content-area {
  flex: 1;
  overflow-y: auto;
  padding: 1rem 0;
}

.input-section {
  border-top: 1px solid #e5e5e5;
  padding-top: 1rem;
}

.context-area {
  padding: 0.5rem;
  margin-bottom: 1rem;
  border: 1px solid #e5e5e5;
  border-radius: 4px;
}

.input-container {
  margin-bottom: 1rem;
}

textarea {
  width: 100%;
  min-height: 100px;
  padding: 0.5rem;
  border: 1px solid #e5e5e5;
  border-radius: 4px;
  resize: vertical;
}

.bottom-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.left-controls {
  display: flex;
  gap: 1rem;
}

.mode-select,
.model-select {
  width: 150px;
}

button {
  padding: 0.5rem 1rem;
  border: 1px solid #e5e5e5;
  border-radius: 4px;
  background: white;
  cursor: pointer;
}

button:hover {
  background: #f5f5f5;
}

.send-btn {
  background: #1a73e8;
  color: white;
  border: none;
}

.send-btn:hover {
  background: #1557b0;
}
</style>
