<template>
  <div class="page-container">
    <art-table :data="tableData" :loading="loading">
      <template #header>
        <el-button type="primary" @click="handleAdd">新增规则</el-button>
      </template>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="规则名称" />
      <el-table-column prop="department" label="目标部门" />
      <el-table-column prop="priority" label="优先级" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.is_enabled ? 'success' : 'info'">{{ row.is_enabled ? '启用' : '禁用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
          <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </art-table>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import request from '@/utils/http'
const loading = ref(false)
const tableData = ref([])
const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get({ url: '/api/admin/ticket-deliver/rules' })
    tableData.value = data || []
  } finally { loading.value = false }
}
const handleAdd = () => {}
const handleEdit = (row: any) => {}
const handleDelete = (row: any) => {}
fetchData()
</script>
