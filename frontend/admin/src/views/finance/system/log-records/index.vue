<template>
  <div class="log-records-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('page.logRecords.title') }}</span>
          <el-button type="danger" @click="handleClearLogs">
            <el-icon><Delete /></el-icon>
            {{ $t('page.logRecords.clearLogs') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('page.logRecords.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('page.logRecords.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('page.logRecords.logType')">
          <el-select v-model="searchForm.type" :placeholder="$t('page.logRecords.all')" clearable>
            <el-option :label="$t('page.logRecords.loginLog')" value="login" />
            <el-option :label="$t('page.logRecords.operationLog')" value="operation" />
            <el-option :label="$t('page.logRecords.errorLog')" value="error" />
            <el-option :label="$t('page.logRecords.systemLog')" value="system" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('page.logRecords.dateRange')">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            :range-separator="$t('page.logRecords.rangeSeparator')"
            :start-placeholder="$t('page.logRecords.startDate')"
            :end-placeholder="$t('page.logRecords.endDate')"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('page.logRecords.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('page.logRecords.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="type" :label="$t('page.logRecords.type')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user" :label="$t('page.logRecords.user')" width="120" />
        <el-table-column prop="content" :label="$t('page.logRecords.logContent')" min-width="300" show-overflow-tooltip />
        <el-table-column prop="ip" :label="$t('page.logRecords.ipAddress')" width="140" />
        <el-table-column prop="created_at" :label="$t('page.logRecords.time')" width="170" />
        <el-table-column :label="$t('page.logRecords.actions')" width="80" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">{{ $t('page.logRecords.detail') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" :title="$t('page.logRecords.logDetail')" width="600px">
      <el-descriptions :column="1" border v-if="currentLog">
        <el-descriptions-item label="ID">{{ currentLog.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.logRecords.type')">
          <el-tag :type="getTypeTag(currentLog.type)" size="small">{{ getTypeText(currentLog.type) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('page.logRecords.user')">{{ currentLog.user }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.logRecords.logContent')">{{ currentLog.content }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.logRecords.ipAddress')">{{ currentLog.ip }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.logRecords.time')">{{ currentLog.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.logRecords.detailInfo')">
          <pre style="white-space: pre-wrap; word-break: break-all;">{{ currentLog.detail || $t('page.logRecords.none') }}</pre>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">{{ $t('page.logRecords.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

const loading = ref(false)

const searchForm = reactive({
  keyword: '',
  type: '',
  date_range: [] as string[]
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])
const detailVisible = ref(false)
const currentLog = ref<any>(null)

const getTypeTag = (type: string) => {
  const map: Record<string, string> = {
    login: 'success',
    operation: 'primary',
    error: 'danger',
    system: 'info'
  }
  return (map[type] || 'info') as any
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    login: $t('page.logRecords.loginLog'),
    operation: $t('page.logRecords.operationLog'),
    error: $t('page.logRecords.errorLog'),
    system: $t('page.logRecords.systemLog')
  }
  return map[type] || type
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      keyword: searchForm.keyword || undefined,
      type: searchForm.type || undefined
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({ url: '/api/admin/log-records', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取日志记录失败:', error)
    ElMessage.error($t('page.logRecords.fetchDataFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.type = ''
  searchForm.date_range = []
  handleSearch()
}

const handleViewDetail = (row: any) => {
  currentLog.value = row
  detailVisible.value = true
}

const handleClearLogs = async () => {
  await ElMessageBox.confirm($t('page.logRecords.confirmClearLogs'), $t('page.logRecords.warning'), {
    type: 'warning'
  })
  try {
    await request.post({ url: '/api/admin/log-records/export' })
    ElMessage.success($t('page.logRecords.clearSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('page.logRecords.clearFailed'))
  }
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchData()
}

const handlePageChange = () => {
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.log-records-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
