<template>
  <div class="database-info-page">
    <art-card title="数据库信息" shadow="never">
      <template #header>
        <div class="card-header">
          <span>数据库信息</span>
          <div>
            <el-button type="primary" @click="fetchDatabaseInfo" :loading="loading">刷新</el-button>
            <el-button @click="handleOptimize" :loading="optimizeLoading">优化全部表</el-button>
            <el-button type="warning" @click="handleBackup">备份数据库</el-button>
          </div>
        </div>
      </template>

      <!-- 概览 -->
      <el-row :gutter="16" class="stat-cards">
        <el-col :span="8">
          <el-card shadow="hover">
            <el-statistic title="数据表总数" :value="dbInfo.total_count" />
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <el-statistic title="总行数" :value="dbInfo.total_rows" />
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <el-statistic title="总大小">
              <template #default>{{ dbInfo.total_size }}</template>
            </el-statistic>
          </el-card>
        </el-col>
      </el-row>

      <!-- 表详情 -->
      <el-table :data="dbInfo.report_array" v-loading="loading" stripe border style="margin-top: 20px" max-height="600">
        <el-table-column prop="name" label="表名" sortable />
        <el-table-column prop="rows" label="行数" width="120" sortable />
        <el-table-column prop="size" label="大小" width="120" sortable />
      </el-table>
    </art-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const optimizeLoading = ref(false)

const dbInfo = reactive({
  total_count: 0,
  total_rows: 0,
  total_size: '',
  report_array: [] as Array<{ name: string; rows: number; size: string }>
})

const fetchDatabaseInfo = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/system/database-info' })
    if (res) Object.assign(dbInfo, res)
  } catch (error) {
    ElMessage.error('获取数据库信息失败')
  } finally {
    loading.value = false
  }
}

const handleOptimize = async () => {
  try {
    await ElMessageBox.confirm('确定要优化所有数据库表吗？', '提示')
    optimizeLoading.value = true
    await request.post({ url: '/api/admin/system/optimize-tables', showSuccessMessage: true })
    fetchDatabaseInfo()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('优化失败')
  } finally {
    optimizeLoading.value = false
  }
}

const handleBackup = async () => {
  try {
    await ElMessageBox.confirm('确定要备份数据库吗？', '提示')
    await request.post({ url: '/api/admin/system/backup-database', showSuccessMessage: true })
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('备份失败')
  }
}

onMounted(() => fetchDatabaseInfo())
</script>

<style scoped lang="scss">
.database-info-page {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.stat-cards {
  margin-bottom: 16px;
}
</style>
