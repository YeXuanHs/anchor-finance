<template>
  <div class="page-container">
    <art-card title="批量同步" shadow="never">
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="任务名称" min-width="150" />
        <el-table-column prop="type" label="同步类型" width="120">
          <template #default="{ row }">
            <el-tag>{{ row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : row.status === 2 ? 'warning' : 'info'">
              {{ row.status === 0 ? '待执行' : row.status === 1 ? '已完成' : '执行中' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total" label="总数" width="100" />
        <el-table-column prop="synced" label="已同步" width="100" />
        <el-table-column prop="failed" label="失败" width="100" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleExecute(row)" :disabled="row.status === 2">执行</el-button>
            <el-button size="small" @click="handleLogs(row)">日志</el-button>
          </template>
        </el-table-column>
      </el-table>
    </art-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const tableData = ref([])

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/batch-sync')
    tableData.value = data?.data || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleExecute = async (row: any) => {
  try {
    await request.post(`/admin/batch-sync/${row.id}/execute`)
    ElMessage.success('开始执行')
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleLogs = (row: any) => {
  // TODO: 显示日志对话框
}

onMounted(() => fetchData())
</script>

<style scoped lang="scss">
.page-container {
  padding: 20px;
}
</style>
