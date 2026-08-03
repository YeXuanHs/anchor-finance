<template>
  <div class="page-container">
    <art-card title="运行映射" shadow="never">
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="host_id" label="主机ID" width="100" />
        <el-table-column prop="type" label="任务类型" width="120">
          <template #default="{ row }">
            <el-tag>{{ row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="action" label="操作" min-width="150" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : row.status === 2 ? 'warning' : 'info'">
              {{ row.status === 0 ? '待执行' : row.status === 1 ? '已完成' : '执行中' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column prop="completed_at" label="完成时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleRepeat(row)" :disabled="row.status !== 1">重试</el-button>
            <el-button size="small" @click="handleDetail(row)">详情</el-button>
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
    const { data } = await request.get('/api/admin/run-map')
    tableData.value = data?.data || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleRepeat = async (row: any) => {
  try {
    await request.post(`/api/admin/run-map/${row.id}/repeat`)
    ElMessage.success('重试已触发')
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleDetail = (row: any) => {
  // TODO: 显示详情
}

onMounted(() => fetchData())
</script>

<style scoped lang="scss">
.page-container {
  padding: 20px;
}
</style>
