<template>
  <div class="task-queue-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>任务队列管理</span>
          <div class="header-actions">
            <el-button type="danger" size="small" @click="handleClearFailed">清除失败任务</el-button>
            <el-button type="primary" size="small" @click="handleRetryAll">重试所有失败</el-button>
          </div>
        </div>
      </template>

      <!-- 统计卡片 -->
      <el-row :gutter="20" class="stat-cards">
        <el-col :span="4">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.pending }}</div>
            <div class="stat-label">待处理</div>
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card shadow="hover" class="stat-card processing">
            <div class="stat-value">{{ stats.processing }}</div>
            <div class="stat-label">处理中</div>
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card shadow="hover" class="stat-card success">
            <div class="stat-value">{{ stats.completed }}</div>
            <div class="stat-label">已完成</div>
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card shadow="hover" class="stat-card danger">
            <div class="stat-value">{{ stats.failed }}</div>
            <div class="stat-label">失败</div>
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.delayed }}</div>
            <div class="stat-label">延迟中</div>
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.total }}</div>
            <div class="stat-label">总计</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="任务类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="邮件发送" value="email" />
            <el-option label="短信发送" value="sms" />
            <el-option label="订单处理" value="order" />
            <el-option label="数据同步" value="sync" />
            <el-option label="报表生成" value="report" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待处理" value="pending" />
            <el-option label="处理中" value="processing" />
            <el-option label="已完成" value="completed" />
            <el-option label="失败" value="failed" />
            <el-option label="延迟" value="delayed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="type" label="任务类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ typeMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="任务名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTypeMap[row.status]" size="small">
              {{ statusLabelMap[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="attempts" label="重试次数" width="80" align="center" />
        <el-table-column prop="max_attempts" label="最大重试" width="80" align="center" />
        <el-table-column prop="run_at" label="计划执行时间" width="180" />
        <el-table-column prop="completed_at" label="完成时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status === 'failed'" type="success" link @click="handleRetry(row)">重试</el-button>
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
            <el-popconfirm v-if="row.status !== 'processing'" title="确定删除吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
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
    <el-dialog v-model="detailVisible" title="任务详情" width="650px">
      <el-descriptions :column="2" border v-if="detailData">
        <el-descriptions-item label="任务ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item label="任务类型">{{ typeMap[detailData.type] || detailData.type }}</el-descriptions-item>
        <el-descriptions-item label="任务名称" :span="2">{{ detailData.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusTypeMap[detailData.status]" size="small">
            {{ statusLabelMap[detailData.status] }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="重试次数">{{ detailData.attempts }} / {{ detailData.max_attempts }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="计划时间">{{ detailData.run_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="完成时间" :span="2">{{ detailData.completed_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="错误信息" :span="2" v-if="detailData.error">
          <pre class="error-text">{{ detailData.error }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="任务数据" :span="2">
          <pre class="data-text">{{ JSON.stringify(detailData.payload, null, 2) }}</pre>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const typeMap: Record<string, string> = {
  email: '邮件发送',
  sms: '短信发送',
  order: '订单处理',
  sync: '数据同步',
  report: '报表生成'
}

const statusTypeMap: Record<string, string> = {
  pending: 'info',
  processing: 'warning',
  completed: 'success',
  failed: 'danger',
  delayed: 'info'
}

const statusLabelMap: Record<string, string> = {
  pending: '待处理',
  processing: '处理中',
  completed: '已完成',
  failed: '失败',
  delayed: '延迟'
}

const loading = ref(false)
const detailVisible = ref(false)
const detailData = ref<any>(null)

const stats = reactive({
  pending: 0,
  processing: 0,
  completed: 0,
  failed: 0,
  delayed: 0,
  total: 0
})

const searchForm = reactive({
  type: '',
  status: ''
})

const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/system/task-queue',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取任务列表失败')
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/system/task-queue/stats' })
    if (data) Object.assign(stats, data)
  } catch (error) {
    console.error('获取统计失败:', error)
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.type = ''; searchForm.status = ''; handleSearch() }
const handleViewDetail = (row: any) => { detailData.value = row; detailVisible.value = true }

const handleRetry = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/system/task-queue/${row.id}/retry` })
    ElMessage.success('已重新加入队列')
    fetchData()
    fetchStats()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/system/task-queue/${row.id}` })
    ElMessage.success('删除成功')
    fetchData()
    fetchStats()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const handleClearFailed = async () => {
  try {
    await request.post({ url: '/api/admin/system/task-queue/clear-failed' })
    ElMessage.success('已清除失败任务')
    fetchData()
    fetchStats()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleRetryAll = async () => {
  try {
    await request.post({ url: '/api/admin/system/task-queue/retry-all' })
    ElMessage.success('已重试所有失败任务')
    fetchData()
    fetchStats()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => {
  fetchData()
  fetchStats()
})
</script>

<style scoped lang="scss">
.task-queue-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.header-actions { display: flex; gap: 8px; }
.search-form { margin-top: 20px; margin-bottom: 20px; }
.stat-cards { margin-bottom: 20px; }
.stat-card {
  text-align: center;
  padding: 8px 0;
  .stat-value { font-size: 24px; font-weight: 600; color: var(--el-color-primary); margin-bottom: 4px; }
  .stat-label { color: var(--el-text-color-secondary); font-size: 14px; }
  &.processing .stat-value { color: var(--el-color-warning); }
  &.success .stat-value { color: var(--el-color-success); }
  &.danger .stat-value { color: var(--el-color-danger); }
}
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.error-text, .data-text {
  margin: 0;
  font-size: 12px;
  max-height: 200px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
