<template>
  <div class="page-container">
    <art-card title="客户服务" shadow="never">
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="product_name" label="产品" min-width="150" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleSuspend(row)" :disabled="row.status !== 1">暂停</el-button>
            <el-button size="small" @click="handleTerminate(row)" :disabled="row.status === 0">终止</el-button>
            <el-button size="small" @click="handleRenew(row)">续费</el-button>
          </template>
        </el-table-column>
      </el-table>
    </art-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const tableData = ref([])

const getStatusType = (status: number) => {
  const map: Record<number, string> = { 0: 'info', 1: 'success', 2: 'warning', 3: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = { 0: '已删除', 1: '使用中', 2: '已暂停', 3: '已终止' }
  return map[status] || '未知'
}

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/admin/client-services')
    tableData.value = data?.data || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSuspend = async (row: any) => {
  await ElMessageBox.confirm('确定暂停该服务？', '提示', { type: 'warning' })
  try {
    await request.post(`/api/admin/client-services/${row.id}/suspend`)
    ElMessage.success('暂停成功')
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleTerminate = async (row: any) => {
  await ElMessageBox.confirm('确定终止该服务？此操作不可逆！', '警告', { type: 'error' })
  try {
    await request.post(`/api/admin/client-services/${row.id}/terminate`)
    ElMessage.success('终止成功')
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleRenew = async (row: any) => {
  // TODO: 续费对话框
}

onMounted(() => fetchData())
</script>

<style scoped lang="scss">
.page-container {
  padding: 20px;
}
</style>
