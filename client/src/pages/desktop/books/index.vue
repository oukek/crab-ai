<template>
  <div class="notes-container">
    <!-- 左侧树形列表 -->
    <div class="sidebar">
      <div class="sidebar-header">
        <div class="title">笔记</div>
        <q-btn flat round dense icon="add" size="sm" class="add-btn" />
      </div>
      
      <q-tree
        :nodes="treeData"
        node-key="id"
        selected-color="grey-3"
        v-model:selected="selectedNote"
        default-expand-all
      >
        <template v-slot:default-header="prop">
          <div class="custom-node">
            <q-icon :name="prop.node.icon || 'description'" size="xs" class="q-mr-sm" />
            <div class="node-label">{{ prop.node.label }}</div>
          </div>
        </template>
      </q-tree>
    </div>

    <!-- 右侧内容区域 -->
    <div class="content">
      <div class="content-header">
        <div class="breadcrumb">
          <q-breadcrumbs separator="/" class="text-grey-7">
            <q-breadcrumbs-el label="我的笔记" icon="folder" />
            <q-breadcrumbs-el label="技术文档" />
          </q-breadcrumbs>
        </div>
        <div class="actions">
          <q-btn flat round dense icon="search" size="sm" class="action-btn" />
          <q-btn flat round dense icon="more_horiz" size="sm" class="action-btn" />
        </div>
      </div>
      
      <div class="editor-container">
        <!-- 这里可以集成富文本编辑器 -->
        <div class="editor-placeholder">
          在这里开始编写...
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

// 树形数据结构
const treeData = ref([
  {
    id: 1,
    label: '我的笔记',
    icon: 'folder',
    children: [
      {
        id: 2,
        label: '技术文档',
        icon: 'folder',
        children: [
          { id: 3, label: 'Vue 学习笔记', icon: 'description' },
          { id: 4, label: 'TypeScript 教程', icon: 'description' }
        ]
      },
      {
        id: 5,
        label: '项目计划',
        icon: 'folder',
        children: [
          { id: 6, label: '2024年规划', icon: 'description' }
        ]
      }
    ]
  }
]);

const selectedNote = ref(null);
</script>

<style scoped>
.notes-container {
  display: flex;
  height: 100vh;
  background-color: #ffffff;
}

.sidebar {
  width: 260px;
  border-right: 1px solid rgba(0, 0, 0, 0.08);
  background-color: #f9f9f9;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}

.title {
  font-size: 16px;
  font-weight: 500;
  color: #1d1d1f;
}

.add-btn {
  color: #666;
}

.custom-node {
  display: flex;
  align-items: center;
  padding: 4px 0;
}

.node-label {
  font-size: 14px;
  color: #1d1d1f;
}

.content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #ffffff;
}

.content-header {
  padding: 16px 24px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.breadcrumb {
  font-size: 14px;
}

.actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  color: #666;
}

.editor-container {
  flex: 1;
  padding: 24px;
  background-color: #ffffff;
}

.editor-placeholder {
  color: #999;
  font-size: 14px;
}

/* Quasar overrides for Mac/Figma style */
:deep(.q-tree__node--selected) {
  background: rgba(0, 0, 0, 0.04) !important;
}

:deep(.q-tree__node:hover:not(.q-tree__node--selected)) {
  background: rgba(0, 0, 0, 0.02);
}

:deep(.q-tree__node-header) {
  padding: 4px 8px;
  border-radius: 6px;
  margin: 2px 8px;
}

:deep(.q-tree__children) {
  padding-left: 20px;
}

:deep(.q-btn) {
  border-radius: 6px;
}

:deep(.q-breadcrumbs) {
  font-size: 13px;
}
</style> 