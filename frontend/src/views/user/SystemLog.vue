<template>
  <div class="system-log-page">
    <div class="page-header">
      <h1 class="page-title">系统操作日志</h1>
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
        <el-form-item label="操作类型">
          <el-select v-model="filters.type" placeholder="全部" clearable style="width: 150px">
            <el-option label="产品操作" value="product" />
            <el-option label="订单操作" value="order" />
            <el-option label="账户操作" value="account" />
            <el-option label="财务操作" value="finance" />
            <el-option label="系统设置" value="system" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="filters.keyword" placeholder="搜索操作描述..." clearable style="width: 200px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">重置</el-button>
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table :data="logs" style="width: 100%" v-loading="loading" stripe>
        <el-table-column prop="time" label="操作时间" width="170" sortable />
        <el-table-column prop="type" label="操作类型" width="110">
          <template #default="{ row }">
            <el-tag :type="typeTagColor(row.type)" size="small" effect="light">
              {{ typeLabel(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="action" label="操作行为" width="130" />
        <el-table-column prop="description" label="操作描述" min-width="280" show-overflow-tooltip />
        <el-table-column prop="operator" label="操作人" width="120" />
        <el-table-column prop="ip" label="IP 地址" width="150" />
        <el-table-column prop="result" label="结果" width="100">
          <template #default="{ row }">
            <el-tag :type="row.result === 'success' ? 'success' : 'danger'" size="small" effect="light" round>
              {{ row.result === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
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
    <el-dialog v-model="showDetailDialog" title="操作详情" width="600px">
      <template v-if="currentLog">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="操作时间">{{ currentLog.time }}</el-descriptions-item>
          <el-descriptions-item label="操作类型">
            <el-tag :type="typeTagColor(currentLog.type)" size="small" effect="light">
              {{ typeLabel(currentLog.type) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="操作行为">{{ currentLog.action }}</el-descriptions-item>
          <el-descriptions-item label="操作人">{{ currentLog.operator }}</el-descriptions-item>
          <el-descriptions-item label="IP 地址">{{ currentLog.ip }}</el-descriptions-item>
          <el-descriptions-item label="结果">
            <el-tag :type="currentLog.result === 'success' ? 'success' : 'danger'" size="small" effect="light" round>
              {{ currentLog.result === 'success' ? '成功' : '失败' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="操作描述" :span="2">{{ currentLog.description }}</el-descriptions-item>
        </el-descriptions>

        <div class="detail-section" v-if="currentLog.detail">
          <h4>详细信息</h4>
          <pre class="code-block">{{ currentLog.detail }}</pre>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Download } from '@element-plus/icons-vue'
import request from '@/utils/request'

interface SystemLog {
  id: string
  time: string
  type: string
  action: string
  description: string
  operator: string
  ip: string
  result: 'success' | 'failed'
  detail: string
}

const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showDetailDialog = ref(false)
const currentLog = ref<SystemLog | null>(null)

const filters = reactive({
  dateRange: null as any,
  type: '',
  keyword: ''
})

const logs = ref<SystemLog[]>([])

function typeLabel(type: string) {
  const map: Record<string, string> = {
    product: '产品操作',
    order: '订单操作',
    account: '账户操作',
    finance: '财务操作',
    system: '系统设置'
  }
  return map[type] || type
}

function typeTagColor(type: string) {
  const map: Record<string, string> = {
    product: '',
    order: 'success',
    account: 'warning',
    finance: 'danger',
    system: 'info'
  }
  return map[type] || ''
}

async function loadData() {
  loading.value = true
  try {
    const res = await request.get('/api/v2/system-logs', {
      params: { page: currentPage.value, page_size: pageSize.value, level: filters.type, module: filters.keyword, date_range: filters.dateRange }
    })
    logs.value = res.data?.data || []
    total.value = res.data?.total || 0
  } catch { /* ignore */ }
  loading.value = false
}

function handleSearch() {
  currentPage.value = 1
  loadData()
}

function handleReset() {
  filters.dateRange = null
  filters.type = ''
  filters.keyword = ''
  currentPage.value = 1
  loadData()
}

async function handleExport() {
  try {
    const res = await request.get('/api/v2/system-logs/export', { params: { level: filters.type, module: filters.keyword }, responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([res.data]))
    const a = document.createElement('a'); a.href = url; a.download = 'system_logs.csv'; a.click()
    window.URL.revokeObjectURL(url)
  } catch { ElMessage.error('导出失败') }
}

function showDetail(row: SystemLog) {
  currentLog.value = row
  showDetailDialog.value = true
}

onMounted(() => { loadData() })

watch([currentPage, pageSize], () => { loadData() })
</script>

<style scoped lang="scss">
.system-log-page {
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
