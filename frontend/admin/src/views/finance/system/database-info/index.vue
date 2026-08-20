<template>
  <div class="database-info-page">
    <art-card :title="$t('page.databaseInfo.title')" shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('page.databaseInfo.title') }}</span>
          <div>
            <el-button type="primary" @click="fetchDatabaseInfo" :loading="loading">{{ $t('page.databaseInfo.refresh') }}</el-button>
            <el-button @click="handleOptimize" :loading="optimizeLoading">{{ $t('page.databaseInfo.optimizeAll') }}</el-button>
            <el-button type="warning" @click="handleBackup">{{ $t('page.databaseInfo.backup') }}</el-button>
          </div>
        </div>
      </template>

      <!-- 概览 -->
      <el-row :gutter="16" class="stat-cards">
        <el-col :span="8">
          <el-card shadow="hover">
            <el-statistic :title="$t('page.databaseInfo.totalTables')" :value="dbInfo.total_count" />
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <el-statistic :title="$t('page.databaseInfo.totalRows')" :value="dbInfo.total_rows" />
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <el-statistic :title="$t('page.databaseInfo.totalSize')">
              <template #default>{{ dbInfo.total_size }}</template>
            </el-statistic>
          </el-card>
        </el-col>
      </el-row>

      <!-- 表详情 -->
      <el-table :data="dbInfo.report_array" v-loading="loading" stripe border style="margin-top: 20px" max-height="600">
        <el-table-column prop="name" :label="$t('page.databaseInfo.tableName')" sortable />
        <el-table-column prop="rows" :label="$t('page.databaseInfo.rows')" width="120" sortable />
        <el-table-column prop="size" :label="$t('page.databaseInfo.size')" width="120" sortable />
      </el-table>
    </art-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

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
    const res = await request.get({ url: '/api/admin/system/database' })
    if (res) Object.assign(dbInfo, res)
  } catch (error) {
    ElMessage.error($t('page.databaseInfo.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleOptimize = async () => {
  try {
    await ElMessageBox.confirm($t('page.databaseInfo.optimizeConfirm'), $t('page.databaseInfo.tips'))
    optimizeLoading.value = true
    await request.post({ url: '/api/admin/system/database/optimize', showSuccessMessage: true })
    fetchDatabaseInfo()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error($t('page.databaseInfo.optimizeFailed'))
  } finally {
    optimizeLoading.value = false
  }
}

const handleBackup = async () => {
  try {
    await ElMessageBox.confirm($t('page.databaseInfo.backupConfirm'), $t('page.databaseInfo.tips'))
    await request.post({ url: '/api/admin/system/database/backup', showSuccessMessage: true })
  } catch (error) {
    if (error !== 'cancel') ElMessage.error($t('page.databaseInfo.backupFailed'))
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
