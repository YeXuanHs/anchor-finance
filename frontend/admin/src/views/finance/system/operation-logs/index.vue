<template>
  <div class="page-container">
    <art-table :data="tableData" :loading="loading">
      <template #header>
        <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" style="width: 240px" />
        <el-button type="primary" @click="fetchData" style="margin-left: 10px">搜索</el-button>
      </template>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="user_id" label="用户ID" width="100" />
      <el-table-column prop="action" label="操作" />
      <el-table-column prop="module" label="模块" width="120" />
      <el-table-column prop="ip" label="IP地址" width="140" />
      <el-table-column prop="created_at" label="时间" width="180" />
      <el-table-column label="详情" width="100">
        <template #default="{ row }">
          <el-button type="primary" link>查看</el-button>
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
const dateRange = ref([])
const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get({ url: '/api/admin/log-records' })
    tableData.value = data || []
  } finally { loading.value = false }
}
fetchData()
</script>
