<template>
  <div class="resource-pool-tasks-page">
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="任务类型">
          <el-select v-model="searchForm.task_type" placeholder="全部" clearable>
            <el-option label="资源分配" value="allocate" />
            <el-option label="资源回收" value="reclaim" />
            <el-option label="资源迁移" value="migrate" />
            <el-option label="资源扩容" value="expand" />
            <el-option label="健康检查" value="health_check" />
          </el-select>
        </el-form-item>
        <el-form-item label="任务状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待执行" value="pending" />
            <el-option label="执行中" value="running" />
            <el-option label="已完成" value="completed" />
            <el-option label="失败" value="failed" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
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

    <el-card shadow="never" class="art-table-card">
      <template #header>
        <div class="card-header">
          <span>任务队列</span>
          <el-space>
            <el-tag type="warning">待执行: {{ stats.pending || 0 }}</el-tag>
            <el-tag type="primary">执行中: {{ stats.running || 0 }}</el-tag>
            <el-tag type="danger">失败: {{ stats.failed || 0 }}</el-tag>
          </el-space>
        </div>
      </template>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="task_no" label="任务编号" width="150" />
        <el-table-column prop="task_type" label="任务类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getTaskTypeTag(row.task_type)" size="small">
              {{ getTaskTypeText(row.task_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="target_resource" label="目标资源" min-width="180" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTag(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="progress" label="进度" width="120">
          <template #default="{ row }">
            <el-progress
              :percentage="row.progress || 0"
              :status="row.status === 'completed' ? 'success' : row.status === 'failed' ? 'exception' : undefined"
              :stroke-width="8"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column prop="completed_at" label="完成时间" width="180">
          <template #default="{ row }">
            {{ row.completed_at || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="operator" label="操作人" width="100" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
            <template v-if="row.status === 'failed'">
              <el-popconfirm title="确定重试该任务吗？" @confirm="handleRetry(row)">
                <template #reference>
                  <el-button type="warning" link>重试</el-button>
                </template>
              </el-popconfirm>
            </template>
            <template v-if="row.status === 'pending' || row.status === 'running'">
              <el-popconfirm title="确定取消该任务吗？" @confirm="handleCancel(row)">
                <template #reference>
                  <el-button type="danger" link>取消</el-button>
                </template>
              </el-popconfirm>
            </template>
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

    <!-- 任务详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="任务详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="任务编号">{{ detailData.task_no }}</el-descriptions-item>
        <el-descriptions-item label="任务类型">{{ getTaskTypeText(detailData.task_type) }}</el-descriptions-item>
        <el-descriptions-item label="目标资源" :span="2">{{ detailData.target_resource }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusTag(detailData.status)">{{ getStatusText(detailData.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="进度">{{ detailData.progress || 0 }}%</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="完成时间">{{ detailData.completed_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="操作人">{{ detailData.operator || '-' }}</el-descriptions-item>
        <el-descriptions-item label="耗时">{{ detailData.duration || '-' }}</el-descriptions-item>
        <el-descriptions-item label="参数" :span="2">
          <pre class="detail-pre">{{ detailData.params ? JSON.stringify(detailData.params, null, 2) : '-' }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="执行结果" :span="2">
          <pre class="detail-pre">{{ detailData.result ? JSON.stringify(detailData.result, null, 2) : '-' }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="错误信息" :span="2" v-if="detailData.error_message">
          <span class="text-danger">{{ detailData.error_message }}</span>
        </el-descriptions-item>
      </el-descriptions>
      <el-divider v-if="detailData.logs && detailData.logs.length" />
      <div v-if="detailData.logs && detailData.logs.length" class="task-logs">
        <h4>执行日志</h4>
        <el-timeline>
          <el-timeline-item
            v-for="(log, index) in detailData.logs"
            :key="index"
            :timestamp="log.time"
            :type="log.level === 'error' ? 'danger' : log.level === 'warn' ? 'warning' : 'primary'"
          >
            {{ log.message }}
          </el-timeline-item>
        </el-timeline>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'ResourcePoolTasks' })

// 加载状态
const loading = ref(false)

// 统计数据
const stats = ref({
  pending: 0,
  running: 0,
  failed: 0
})

// 搜索表单
const searchForm = reactive({
  task_type: undefined as string | undefined,
  status: undefined as string | undefined,
  date_range: [] as string[]
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref([])

// 详情对话框
const detailDialogVisible = ref(false)
const detailData = ref<any>({})

// 获取任务类型标签
const getTaskTypeTag = (type: string) => {
  const map: Record<string, string> = {
    allocate: 'primary',
    reclaim: 'warning',
    migrate: 'success',
    expand: 'info',
    health_check: ''
  }
  return map[type] || 'info'
}

// 获取任务类型文本
const getTaskTypeText = (type: string) => {
  const map: Record<string, string> = {
    allocate: '资源分配',
    reclaim: '资源回收',
    migrate: '资源迁移',
    expand: '资源扩容',
    health_check: '健康检查'
  }
  return map[type] || '未知'
}

// 获取状态标签
const getStatusTag = (status: string) => {
  const map: Record<string, string> = {
    pending: 'info',
    running: 'primary',
    completed: 'success',
    failed: 'danger',
    cancelled: 'warning'
  }
  return map[status] || 'info'
}

// 获取状态文本
const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待执行',
    running: '执行中',
    completed: '已完成',
    failed: '失败',
    cancelled: '已取消'
  }
  return map[status] || '未知'
}

// 获取任务列表
const fetchTasks = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size
    }
    if (searchForm.task_type) params.task_type = searchForm.task_type
    if (searchForm.status) params.status = searchForm.status
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }

    const data = await request.get({
      url: '/api/admin/resource-pool/tasks',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
    stats.value = {
      pending: data.pending || 0,
      running: data.running || 0,
      failed: data.failed || 0
    }
  } catch (error) {
    console.error('获取任务列表失败:', error)
    ElMessage.error('获取任务列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchTasks()
}

// 重置
const handleReset = () => {
  searchForm.task_type = undefined
  searchForm.status = undefined
  searchForm.date_range = []
  handleSearch()
}

// 查看详情
const handleViewDetail = async (row: any) => {
  try {
    const data = await request.get({
      url: `/api/admin/resource-pool/tasks/${row.id}`
    })
    detailData.value = data || row
    detailDialogVisible.value = true
  } catch (error) {
    detailData.value = row
    detailDialogVisible.value = true
  }
}

// 重试任务
const handleRetry = async (row: any) => {
  try {
    await request.post({
      url: `/api/admin/resource-pool/tasks/${row.id}/retry`,
      showSuccessMessage: true
    })
    ElMessage.success('任务已重新提交')
    fetchTasks()
  } catch (error) {
    ElMessage.error('重试失败')
  }
}

// 取消任务
const handleCancel = async (row: any) => {
  try {
    await request.post({
      url: `/api/admin/resource-pool/tasks/${row.id}/cancel`,
      showSuccessMessage: true
    })
    ElMessage.success('任务已取消')
    fetchTasks()
  } catch (error) {
    ElMessage.error('取消失败')
  }
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchTasks()
}

// 页码变化
const handlePageChange = () => {
  fetchTasks()
}

onMounted(() => {
  fetchTasks()
})
</script>

<style scoped lang="scss">
.resource-pool-tasks-page {
  padding: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  .el-form-item {
    margin-bottom: 0;
  }
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.text-danger {
  color: #f56c6c;
}

.detail-pre {
  margin: 0;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 200px;
  overflow-y: auto;
}

.task-logs {
  h4 {
    margin: 0 0 16px 0;
    font-size: 14px;
    color: #303133;
  }
}
</style>
