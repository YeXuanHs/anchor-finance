<template>
  <div class="operation-logs-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('operationLog.title') }}</span>
          <el-button type="danger" @click="handleClearLogs"><el-icon><Delete /></el-icon>{{ $t('operationLog.clearLogs') }}</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('operationLog.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('operationLog.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('operationLog.dateRange')">
          <el-date-picker v-model="searchForm.date_range" type="daterange" :range-separator="$t('common.to')" :start-placeholder="$t('common.startDate')" :end-placeholder="$t('common.endDate')" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="admin_name" :label="$t('operationLog.operator')" width="120" />
        <el-table-column prop="description" :label="$t('operationLog.description')" min-width="300" show-overflow-tooltip />
        <el-table-column prop="ip" :label="$t('operationLog.ip')" width="140" />
        <el-table-column prop="created_at" :label="$t('operationLog.operationTime')" width="170" />
        <el-table-column :label="$t('operationLog.operations')" width="80" fixed="right" align="center">
          <template #default="{ row }"><el-button type="primary" link @click="handleViewDetail(row)">{{ $t('operationLog.detail') }}</el-button></template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <el-dialog v-model="detailVisible" :title="$t('operationLog.detailTitle')" width="600px">
      <el-descriptions :column="1" border v-if="currentLog">
        <el-descriptions-item label="ID">{{ currentLog.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('operationLog.operator')">{{ currentLog.admin_name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('operationLog.description')">{{ currentLog.description }}</el-descriptions-item>
        <el-descriptions-item :label="$t('operationLog.ip')">{{ currentLog.ip }}</el-descriptions-item>
        <el-descriptions-item :label="$t('operationLog.operationTime')">{{ currentLog.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('operationLog.detailInfo')"><pre style="white-space: pre-wrap; word-break: break-all;">{{ currentLog.detail || '-' }}</pre></el-descriptions-item>
      </el-descriptions>
      <template #footer><el-button @click="detailVisible = false">{{ $t('common.close') }}</el-button></template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const searchForm = reactive({ keyword: '', date_range: [] as string[] })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const detailVisible = ref(false)
const currentLog = ref<any>(null)

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size, keyword: searchForm.keyword || undefined }
    if (searchForm.date_range?.length === 2) { params.start_date = searchForm.date_range[0]; params.end_date = searchForm.date_range[1] }
    const data = await request.get({ url: '/api/admin/log-records', params }); tableData.value = data.list || []; pagination.total = data.total || 0
  } catch { ElMessage.error($t('operationLog.fetchFailed')) } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.date_range = []; handleSearch() }
const handleViewDetail = (row: any) => { currentLog.value = row; detailVisible.value = true }

const handleClearLogs = async () => {
  await ElMessageBox.confirm($t('operationLog.confirmClear'), $t('common.warning'), { type: 'warning' })
  try { await request.del({ url: '/api/admin/system-logs/clear-level' }); ElMessage.success($t('operationLog.clearSuccess')); fetchData() } catch { ElMessage.error($t('operationLog.clearFailed')) }
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.operation-logs-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
