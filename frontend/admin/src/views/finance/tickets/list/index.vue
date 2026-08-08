<template>
  <div class="ticket-list-page">
    <!-- 标签页切换 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane name="all">
        <template #label>
          <span>全部工单</span>
        </template>
      </el-tab-pane>
      <el-tab-pane name="open">
        <template #label>
          <span :class="{ 'has-unread': unreadCounts.open > 0 }">
            待处理
            <el-badge v-if="unreadCounts.open > 0" :value="unreadCounts.open" class="unread-badge" />
          </span>
        </template>
      </el-tab-pane>
      <el-tab-pane name="in_progress">
        <template #label>
          <span :class="{ 'has-unread': unreadCounts.in_progress > 0 }">
            处理中
            <el-badge v-if="unreadCounts.in_progress > 0" :value="unreadCounts.in_progress" class="unread-badge" />
          </span>
        </template>
      </el-tab-pane>
      <el-tab-pane name="replied">
        <template #label>
          <span :class="{ 'has-unread': unreadCounts.replied > 0 }">
            已回复
            <el-badge v-if="unreadCounts.replied > 0" :value="unreadCounts.replied" class="unread-badge" />
          </span>
        </template>
      </el-tab-pane>
      <el-tab-pane name="closed">
        <template #label>
          <span>已关闭</span>
        </template>
      </el-tab-pane>
    </el-tabs>

    <!-- 搜索筛选区域 -->
    <el-card shadow="never" class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="工单号/主题/客户名"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="优先级">
          <el-select v-model="searchForm.priority" placeholder="全部" clearable style="width: 100px">
            <el-option label="低" :value="1" />
            <el-option label="普通" :value="2" />
            <el-option label="高" :value="3" />
            <el-option label="紧急" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="部门">
          <el-select v-model="searchForm.department_id" placeholder="全部" clearable style="width: 120px">
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="创建时间">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作栏 -->
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <div class="action-left">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新建工单
          </el-button>
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
        <div class="action-right">
          <el-button circle @click="fetchList">
            <el-icon><Refresh /></el-icon>
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        :row-class-name="getRowClassName"
        @sort-change="handleSortChange"
      >
        <el-table-column prop="id" label="ID" width="80" sortable="custom" align="center" />
        <el-table-column prop="ticket_no" label="工单号" width="130">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)" :class="{ 'font-bold': row.is_unread }">
              {{ row.ticket_no }}
            </el-button>
            <span v-if="row.is_unread" class="new-badge">new</span>
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="主题" min-width="200">
          <template #default="{ row }">
            <span :class="{ 'font-bold': row.is_unread }">{{ row.subject }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="client_name" label="客户" width="120">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewClient(row)">
              {{ row.client_name }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="department_name" label="部门" width="100" />
        <el-table-column prop="priority" label="优先级" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="getPriorityType(row.priority)" size="small">
              {{ getPriorityText(row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_reply_at" label="最后回复" width="170" sortable="custom" />
        <el-table-column prop="created_at" label="创建时间" width="170" sortable="custom" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleView(row)">
              查看
            </el-button>
            <el-button
              v-if="row.status !== 'closed'"
              type="warning"
              link
              size="small"
              @click="handleClose(row)"
            >
              关闭
            </el-button>
            <el-button
              v-if="row.status === 'closed'"
              type="success"
              link
              size="small"
              @click="handleReopen(row)"
            >
              重新打开
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, Download } from '@element-plus/icons-vue'
import request from '@/utils/http'

const router = useRouter()
const loading = ref(false)
const tableData = ref([])
const activeTab = ref('all')
const departments = ref<{ id: number; name: string }[]>([])

// 未读计数
const unreadCounts = reactive({
  open: 0,
  in_progress: 0,
  replied: 0
})

// 搜索表单
const searchForm = reactive({
  keyword: '',
  priority: null as number | null,
  department_id: null as number | null,
  date_range: null as [Date, Date] | null
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 排序
const sortParams = reactive({
  sort: 'created_at',
  order: 'desc'
})

// 提示音
let notificationSound: HTMLAudioElement | null = null
let pollTimer: ReturnType<typeof setInterval> | null = null

// 初始化提示音
const initNotificationSound = () => {
  notificationSound = new Audio('/notification.mp3')
  notificationSound.volume = 0.5
}

// 播放提示音
const playNotificationSound = () => {
  if (notificationSound) {
    notificationSound.currentTime = 0
    notificationSound.play().catch(() => {})
  }
}

// 行样式
const getRowClassName = ({ row }: { row: any }) => {
  return row.is_unread ? 'unread-row' : ''
}

// 优先级类型
const getPriorityType = (priority: number): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<number, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { 1: 'info', 2: 'primary', 3: 'warning', 4: 'danger' }
  return map[priority] || 'info'
}

// 优先级文本
const getPriorityText = (priority: number) => {
  const map: Record<number, string> = { 1: '低', 2: '普通', 3: '高', 4: '紧急' }
  return map[priority] || '未知'
}

// 状态类型
const getStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    open: 'warning',
    in_progress: 'primary',
    replied: 'success',
    closed: 'info'
  }
  return map[status] || 'info'
}

// 状态文本
const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    open: '待处理',
    in_progress: '处理中',
    replied: '已回复',
    closed: '已关闭'
  }
  return map[status] || '未知'
}

// 标签页切换
const handleTabChange = (tab: string | number) => {
  pagination.page = 1
  fetchList()
}

// 获取列表数据
const fetchList = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      sort: sortParams.sort,
      order: sortParams.order
    }

    if (activeTab.value !== 'all') params.status = activeTab.value
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.priority) params.priority = searchForm.priority
    if (searchForm.department_id) params.department_id = searchForm.department_id
    if (searchForm.date_range) {
      params.start_date = searchForm.date_range[0].toISOString().split('T')[0]
      params.end_date = searchForm.date_range[1].toISOString().split('T')[0]
    }

    const data = await request.get({ url: '/api/admin/tickets', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取工单列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取未读计数
const fetchUnreadCounts = async () => {
  try {
    const data = await request.get({ url: '/api/admin/tickets/unread-counts' })
    const oldTotal = unreadCounts.open + unreadCounts.in_progress + unreadCounts.replied
    unreadCounts.open = data.open || 0
    unreadCounts.in_progress = data.in_progress || 0
    unreadCounts.replied = data.replied || 0
    const newTotal = unreadCounts.open + unreadCounts.in_progress + unreadCounts.replied
    if (newTotal > oldTotal && oldTotal > 0) {
      playNotificationSound()
    }
  } catch (error) {
    console.error('获取未读计数失败:', error)
  }
}

// 获取部门
const fetchDepartments = async () => {
  try {
    const data = await request.get({ url: '/api/admin/ticket-departments' })
    departments.value = data || []
  } catch (error) {
    console.error('获取部门失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.priority = null
  searchForm.department_id = null
  searchForm.date_range = null
  pagination.page = 1
  fetchList()
}

// 排序变化
const handleSortChange = ({ prop, order }: any) => {
  sortParams.sort = prop || 'created_at'
  sortParams.order = order === 'ascending' ? 'asc' : 'desc'
  fetchList()
}

// 分页大小变化
const handleSizeChange = (size: number) => {
  pagination.page_size = size
  pagination.page = 1
  fetchList()
}

// 页码变化
const handlePageChange = (page: number) => {
  pagination.page = page
  fetchList()
}

// 新建工单
const handleAdd = () => {
  router.push('/add-support-ticket')
}

// 查看工单
const handleView = (row: any) => {
  router.push(`/ticket-detail/${row.id}`)
}

// 查看客户
const handleViewClient = (row: any) => {
  router.push(`/customer-view/${row.client_id}`)
}

// 关闭工单
const handleClose = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要关闭工单 "${row.ticket_no}" 吗？`, '确认关闭', {
      type: 'warning'
    })
    await request.post({ url: `/api/admin/tickets/${row.id}/close` })
    ElMessage.success('工单已关闭')
    fetchList()
    fetchUnreadCounts()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('关闭工单失败:', error)
    }
  }
}

// 重新打开工单
const handleReopen = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要重新打开工单 "${row.ticket_no}" 吗？`, '确认重新打开', {
      type: 'warning'
    })
    await request.post({ url: `/api/admin/tickets/${row.id}/reopen` })
    ElMessage.success('工单已重新打开')
    fetchList()
    fetchUnreadCounts()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重新打开工单失败:', error)
    }
  }
}

// 导出
const handleExport = async () => {
  try {
    const params: any = {}
    if (activeTab.value !== 'all') params.status = activeTab.value
    if (searchForm.keyword) params.keyword = searchForm.keyword

    const response = await request.get({
      url: '/api/admin/tickets/export',
      params,
      responseType: 'blob'
    })

    const url = window.URL.createObjectURL(new Blob([response]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `工单列表_${new Date().toISOString().split('T')[0]}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('导出失败:', error)
  }
}

onMounted(() => {
  initNotificationSound()
  fetchList()
  fetchDepartments()
  fetchUnreadCounts()
  // 每30秒轮询未读计数
  pollTimer = setInterval(fetchUnreadCounts, 30000)
})

onUnmounted(() => {
  if (pollTimer) {
    clearInterval(pollTimer)
  }
})
</script>

<style scoped lang="scss">
.ticket-list-page {
  padding: 16px;
}

.search-card {
  margin-bottom: 16px;

  :deep(.el-card__body) {
    padding-bottom: 0;
  }
}

.action-card {
  margin-bottom: 16px;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.action-left {
  display: flex;
  gap: 8px;
}

.table-card {
  :deep(.el-card__body) {
    padding: 0;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px;
}

.has-unread {
  font-weight: bold;
  color: var(--el-color-primary);
}

.unread-badge {
  margin-left: 4px;

  :deep(.el-badge__content) {
    font-size: 10px;
  }
}

.new-badge {
  display: inline-block;
  padding: 0 4px;
  margin-left: 4px;
  font-size: 10px;
  color: #fff;
  background: #F59E0B;
  border-radius: 2px;
  vertical-align: middle;
}

.font-bold {
  font-weight: bold;
}

:deep(.unread-row) {
  background-color: #FFF7ED !important;

  &:hover > td {
    background-color: #FFEDD5 !important;
  }
}
</style>
