<template>
  <div class="page-container">
    <art-table :data="tableData" :loading="loading">
      <template #header>
        <el-button type="primary" @click="handleAdd">新增服务器</el-button>
      </template>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="服务器名称" />
      <el-table-column prop="ip" label="IP地址" width="140" />
      <el-table-column prop="datacenter" label="数据中心" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'online' ? 'success' : 'danger'">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="300" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link>详情</el-button>
          <el-button type="success" link>开机</el-button>
          <el-button type="warning" link>关机</el-button>
          <el-button type="info" link>重启</el-button>
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
    const { data } = await request.get({ url: '/api/admin/dcim/servers' })
    tableData.value = data || []
  } finally { loading.value = false }
}
const handleAdd = () => {}
fetchData()
</script>
