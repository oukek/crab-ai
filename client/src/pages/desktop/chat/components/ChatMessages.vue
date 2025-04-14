<template>
  <div class="chat-messages">
    <div v-for="(message, index) in messages" 
         :key="index" 
         :class="['message-wrapper', message.role === 'user' ? 'user-message' : 'ai-message']">
      <div class="message">
        <div class="message-content">
          {{ message.content }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Message {
  role: 'user' | 'assistant';
  content: string;
  timestamp?: Date;
}

defineProps<{
  messages: Message[]
}>();

const formatTime = (date: Date) => {
  return new Intl.DateTimeFormat('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  }).format(date);
};
</script>

<style scoped>
.chat-messages {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  height: 100%;
  overflow-y: auto;
}

.message-wrapper {
  display: flex;
  width: 100%;
}

.message {
  max-width: 80%;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.message-content {
  padding: 0.75rem 1rem;
  border-radius: 1rem;
  font-size: 0.95rem;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
}

.message-time {
  font-size: 0.75rem;
  color: #8e8e93;
  margin: 0 0.5rem;
}

/* User message styles */
.user-message {
  justify-content: flex-end;
}

.user-message .message-content {
  background-color: #007AFF;
  color: white;
  border-bottom-right-radius: 0.25rem;
}

/* AI message styles */
.ai-message .message-content {
  background-color: #F2F2F7;
  color: #1C1C1E;
  border-bottom-left-radius: 0.25rem;
}

/* Modern scrollbar styles */
.chat-messages::-webkit-scrollbar {
  width: 8px;
}

.chat-messages::-webkit-scrollbar-track {
  background: transparent;
}

.chat-messages::-webkit-scrollbar-thumb {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 4px;
}

.chat-messages::-webkit-scrollbar-thumb:hover {
  background-color: rgba(0, 0, 0, 0.3);
}
</style> 