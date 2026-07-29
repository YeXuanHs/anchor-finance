<template>
  <div class="logs-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="操作内容/用户" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option v-for="t in logTypes" :key="t.value" :label="t.label" :value="t.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="级别">
          <el-select v-model="searchForm.level" placeholder="全部" clearable>
            <el-option label="信息" value="info" />
            <el-option label="警告" value="warning" />
            <el-option label="错误" value="error" />
            <el-option label="危险" value="danger" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
          <el-date-picker v-model="searchForm.date_range" type="daterange" start-placeholder="开始日期" end-placeholder="结束日期" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>搜索
          </el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>系统日志</h3>
        <div>
          <el-button @click="handleRefresh">
            <el-icon><Refresh /></el-icon>刷新
          </el-button>
          <el-button type="danger" @click="handleCleanup">
            <el-icon><Delete /></el-icon>清理日志
          </el-button>
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>导出
          </el-button>
        </div>
      </div>

      <el-table :data="logs" style="width: 100%" v-loading="loading" @row-click="handleViewDetail">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="type" label="类型" width="110">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="level" label="级别" width="80">
          <template #default="{ row }">
            <el-tag :type="getLevelType(row.level)" size="small">{{ getLevelLabel(row.level) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="operator" label="操作人" width="120" />
        <el-table-column prop="ip" label="IP地址" width="130" />
        <el-table-column prop="module" label="模块" width="100" />
        <el-table-column prop="action" label="操作" min-width="200" show-overflow-tooltip />
        <el-table-column prop="target" label="操作对象" min-width="150" show-overflow-tooltip />
        <el-table-column prop="created_at" label="时间" width="170" />
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click.stop="handleViewDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[20, 50, 100, 200]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 日志详情 -->
    <el-dialog title="日志详情" v-model="detailVisible" width="700px">
      <div class="log-detail" v-if="currentLog">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">{{ currentLog.id }}</el-descriptions-item>
          <el-descriptions-item label="时间">{{ currentLog.created_at }}</el-descriptions-item>
          <el-descriptions-item label="类型">
            <el-tag :type="getTypeTag(currentLog.type)" size="small">{{ getTypeLabel(currentLog.type) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="级别">
            <el-tag :type="getLevelType(currentLog.level)" size="small">{{ getLevelLabel(currentLog.level) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="操作人">{{ currentLog.operator || '-' }}</el-descriptions-item>
          <el-descriptions-item label="IP地址">{{ currentLog.ip || '-' }}</el-descriptions-item>
          <el-descriptions-item label="模块">{{ currentLog.module || '-' }}</el-descriptions-item>
          <el-descriptions-item label="UA" :span="2">
            <span class="ua-text">{{ currentLog.user_agent || '-' }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="操作" :span="2">{{ currentLog.action }}</el-descriptions-item>
          <el-descriptions-item label="操作对象" :span="2">{{ currentLog.target || '-' }}</el-descriptions-item>
        </el-descriptions>

        <div class="detail-section" v-if="currentLog.request_data">
          <h4>请求数据</h4>
          <pre class="json-block">{{ formatJson(currentLog.request_data) }}</pre>
        </div>

        <div class="detail-section" v-if="currentLog.response_data">
          <h4>响应数据</h4>
          <pre class="json-block">{{ formatJson(currentLog.response_data) }}</pre>
        </div>

        <div class="detail-section" v-if="currentLog.extra">
          <h4>附加信息</h4>
          <pre class="json-block">{{ formatJson(currentLog.extra) }}</pre>
        </div>
      </div>
    </el-dialog>

    <!-- 清理对话框 -->
    <el-dialog title="清理日志" v-model="cleanupVisible" width="450px">
      <el-form :model="cleanupForm" label-width="100px">
        <el-form-item label="保留天数">
          <el-input-number v-model="cleanupForm.days" :min="1" :max="365" />
          <span style="margin-left: 8px; color: var(--text-secondary);">天之前的日志将被清理</span>
        </el-form-item>
        <el-form-item label="日志类型">
          <el-select v-model="cleanupForm.type" placeholder="全部类型" clearable>
            <el-option v-for="t in logTypes" :key="t.value" :label="t.label" :value="t.value" />
          </el-select>
        </el-form-item>
        <el-alert type="warning" :closable="false">
          清理操作不可恢复，请确认后再执行。
        </el-alert>
      </el-form>
      <template #footer>
        <el-button @click="cleanupVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmCleanup" :loading="cleanupLoading">确认清理</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Delete, Download } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const logs = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const detailVisible = ref(false)
const cleanupVisible = ref(false)
const cleanupLoading = ref(false)
const currentLog = ref<any>(null)

const logTypes = [
  { label: '登录日志', value: 'login' },
  { label: '操作日志', value: 'operation' },
  { label: '系统日志', value: 'system' },
  { label: '安全日志', value: 'security' },
  { label: 'API日志', value: 'api' },
  { label: '错误日志', value: 'error' }
]

const searchForm = ref({ keyword: '', type: '', level: '', date_range: null as any })
const cleanupForm = ref({ days: 30, type: '' })

const getTypeLabel = (val: string) => logTypes.find(t => t.value === val)?.label || val
const getTypeTag = (type: string) => {
  const map: Record<string, string> = { login: '', operation: 'success', system: 'info', security: 'danger', api: 'warning', error: 'danger' }
  return map[type] || 'info'
}
const getLevelType = (level: string) => {
  const map: Record<string, string> = { info: '', warning: 'warning', error: 'danger', danger: 'danger' }
  return map[level] || 'info'
}
const getLevelLabel = (level: string) => {
  const map: Record<string, string> = { info: '信息', warning: '警告', error: '错误', danger: '危险' }
  return map[level] || level
}

const formatJson = (data: any) => {
  if (!data) return ''
  try {
    return typeof data === 'string' ? JSON.stringify(JSON.parse(data), null, 2) : JSON.stringify(data, null, 2)
  } catch {
    return String(data)
  }
}

const fetchLogs = async () => {
  loading.value = true
  try {
    const params: any = { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    delete params.date_range
    if (searchForm.value.date_range?.length === 2) {
      params.start_date = searchForm.value.date_range[0]
      params.end_date = searchForm.value.date_range[1]
    }
    const { data } = await request.get('/admin/api/v1/logs', { params })
    logs.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取日志失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { currentPage.value = 1; fetchLogs() }
const resetSearch = () => { searchForm.value = { keyword: '', type: '', level: '', date_range: null }; handleSearch() }
const handleSizeChange = (val: number) => { pageSize.value = val; fetchLogs() }
const handlePageChange = (val: number) => { currentPage.value = val; fetchLogs() }
const handleRefresh = () => { fetchLogs() }

const handleViewDetail = (row: any) => {
  currentLog.value = row
  detailVisible.value = true
}

const handleCleanup = () => {
  cleanupForm.value = { days: 30, type: '' }
  cleanupVisible.value = true
}

const confirmCleanup = async () => {
  cleanupLoading.value = true
  try {
    const { data } = await request.post('/admin/api/v1/logs/cleanup', cleanupForm.value)
    ElMessage.success(`已清理 ${data.data?.deleted || 0} 条日志`)
    cleanupVisible.value = false
    fetchLogs()
  } catch {
    ElMessage.error('清理失败')
  } finally {
    cleanupLoading.value = false
  }
}

const handleExport = async () => {
  try {
    const params: any = { ...searchForm.value }
    delete params.date_range
    if (searchForm.value.date_range?.length === 2) {
      params.start_date = searchForm.value.date_range[0]
      params.end_date = searchForm.value.date_range[1]
    }
    const { data } = await request.get('/admin/api/v1/logs/export', { params, responseType: 'blob' })
    const url = URL.createObjectURL(data)
    const a = document.createElement('a')
    a.href = url; a.download = 'system-logs.xlsx'; a.click()
    URL.revokeObjectURL(url)
  } catch {
    ElMessage.error('导出失败')
  }
}

onMounted(() => { fetchLogs() })
</script>

<style scoped lang="scss">
.logs-page {
  .table-header {
    display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
  .log-detail {
    .ua-text { font-size: 12px; color: var(--text-secondary); word-break: break-all; }
    .detail-section {
      margin-top: 16px;
      h4 { margin: 0 0 8px; font-size: 14px; font-weight: 600; }
    }
    .json-block {
      background: var(--bg-code, #f5f7fa);
      border: 1px solid var(--border-color-lighter, #e4e7ed);
      border-radius: 8px;
      padding: 12px 16px;
      font-size: 12px;
      font-family: 'Courier New', monospace;
      max-height: 300px;
      overflow: auto;
      white-space: pre-wrap;
      word-break: break-all;
    }
  }
}
</style>
