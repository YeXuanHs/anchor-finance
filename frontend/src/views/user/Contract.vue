<template>
  <div class="contract-page">
    <el-card>
      <template #header><span>我的合同</span></template>
      <el-table :data="contracts" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="合同标题" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'signed' ? 'success' : 'warning'">{{ row.status === 'signed' ? '已签署' : '待签署' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" link>查看</el-button>
            <el-button v-if="row.status !== 'signed'" type="success" link>签署</el-button>
            <el-button type="info" link>下载PDF</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import request from '@/utils/http'
const loading = ref(false)
const contracts = ref([])
const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/contracts')
    contracts.value = data || []
  } finally { loading.value = false }
}
fetchData()
</script>
