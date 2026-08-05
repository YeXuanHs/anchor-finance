<template>
  <div class="api-log-page">
    <div class="page-header">
      <h1 class="page-title">API 调用日志</h1>
    </div>

    <el-card shadow="never" class="filter-card">
      <el-form :model="filters" inline>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="filters.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 300px"
          />
        </el-form-item>
        <el-form-item label="接口路径">
          <el-input v-model="filters.endpoint" placeholder="如 /api/v1/users" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="状态码">
          <el-select v-model="filters.statusCode" placeholder="全部" clearable style="width: 140px">
            <el-option label="2xx 成功" value="2xx" />
            <el-option label="4xx 客户端错误" value="4xx" />
            <el-option label="5xx 服务端错误" value="5xx" />
          </el-select>
        </el-form-item>
        <el-form-item label="请求方式">
          <el-select v-model="filters.method" placeholder="全部" clearable style="width: 130px">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table :data="logs" style="width: 100%" v-loading="loading" stripe>
        <el-table-column prop="time" label="请求时间" width="170" sortable />
        <el-table-column prop="method" label="方式" width="90">
          <template #default="{ row }">
            <el-tag
              :type="methodTagType(row.method)"
              size="small"
              effect="dark"
              round
            >
              {{ row.method }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="endpoint" label="接口路径" min-width="250" show-overflow-tooltip>
          <template #default="{ row }">
            <code class="endpoint-text">{{ row.endpoint }}</code>
          </template>
        </el-table-column>
        <el-table-column prop="statusCode" label="状态码" width="100">
          <template #default="{ row }">
            <el-tag
              :type="statusTagType(row.statusCode)"
              size="small"
              effect="light"
              round
            >
              {{ row.statusCode }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="耗时" width="100" sortable>
          <template #default="{ row }">
            <span :class="{ 'slow-request': row.duration > 1000 }">
              {{ row.duration }}ms
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP 地址" width="150" />
        <el-table-column prop="apiKeyName" label="密钥名称" width="140" />
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" size="small" @click="showDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="showDetailDialog" title="请求详情" width="680px">
      <template v-if="currentLog">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="请求时间">{{ currentLog.time }}</el-descriptions-item>
          <el-descriptions-item label="请求方式">
            <el-tag :type="methodTagType(currentLog.method)" size="small" effect="dark" round>
              {{ currentLog.method }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="接口路径" :span="2">
            <code>{{ currentLog.endpoint }}</code>
          </el-descriptions-item>
          <el-descriptions-item label="状态码">
            <el-tag :type="statusTagType(currentLog.statusCode)" size="small" effect="light" round>
              {{ currentLog.statusCode }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="耗时">{{ currentLog.duration }}ms</el-descriptions-item>
          <el-descriptions-item label="IP 地址">{{ currentLog.ip }}</el-descriptions-item>
          <el-descriptions-item label="密钥名称">{{ currentLog.apiKeyName }}</el-descriptions-item>
        </el-descriptions>

        <div class="detail-section">
          <h4>请求头</h4>
          <pre class="code-block">{{ currentLog.requestHeaders }}</pre>
        </div>

        <div class="detail-section">
          <h4>请求体</h4>
          <pre class="code-block">{{ currentLog.requestBody || '(无)' }}</pre>
        </div>

        <div class="detail-section">
          <h4>响应体</h4>
          <pre class="code-block">{{ currentLog.responseBody }}</pre>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import request from '@/utils/request'

interface ApiLog {
  id: string
  time: string
  method: string
  endpoint: string
  statusCode: number
  duration: number
  ip: string
  apiKeyName: string
  requestHeaders: string
  requestBody: string
  responseBody: string
}

const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showDetailDialog = ref(false)
const currentLog = ref<ApiLog | null>(null)

const filters = reactive({
  dateRange: null as any,
  endpoint: '',
  statusCode: '',
  method: ''
})

const logs = ref<ApiLog[]>([])

async function loadData() {
  loading.value = true
  try {
    const res = await request.get('/api/v1/api-logs', {
      params: { page: currentPage.value, page_size: pageSize.value, method: filters.method, endpoint: filters.endpoint, status_code: filters.statusCode }
    })
    logs.value = res.data?.data || []
    total.value = res.data?.total || 0
  } catch { /* ignore */ }
  loading.value = false
}

function methodTagType(method: string) {
  const map: Record<string, string> = { GET: 'success', POST: 'primary', PUT: 'warning', DELETE: 'danger' }
  return map[method] || 'info'
}

function statusTagType(code: number) {
  if (code >= 200 && code < 300) return 'success'
  if (code >= 400 && code < 500) return 'warning'
  return 'danger'
}

function handleSearch() {
  currentPage.value = 1
  loadData()
}

function handleReset() {
  filters.dateRange = null
  filters.endpoint = ''
  filters.statusCode = ''
  filters.method = ''
}

function showDetail(row: ApiLog) {
  currentLog.value = row
  showDetailDialog.value = true
}

onMounted(() => { loadData() })
</script>

<style scoped lang="scss">
.api-log-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.filter-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;

  :deep(.el-form-item) {
    margin-bottom: 0;
  }
}

.table-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.endpoint-text {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: #303133;
}

.slow-request {
  color: #f56c6c;
  font-weight: 600;
}

.detail-section {
  margin-top: 16px;

  h4 {
    font-size: 14px;
    font-weight: 600;
    color: #303133;
    margin: 0 0 8px 0;
  }
}

.code-block {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 12px 16px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: #606266;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
}
</style>
