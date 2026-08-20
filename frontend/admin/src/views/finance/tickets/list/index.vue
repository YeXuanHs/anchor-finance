<template>
  <div class="ticket-list-page">
    <!-- 页面说明 -->
    <h5 class="page-title">{{ $t('ticketList.pageDescription') }}</h5>

    <!-- 未读工单提醒 -->
    <el-alert
      v-if="hasUnread"
      :title="$t('ticketList.unreadAlert')"
      type="warning"
      show-icon
      :closable="false"
      class="unread-alert"
    />

    <!-- 自动刷新 -->
    <div class="refresh-bar">
      <span>{{ $t('ticketList.autoRefresh') }}</span>
      <el-select v-model="refreshInterval" size="small" style="width: 90px" @change="handleRefreshChange">
        <el-option :label="$t('ticketList.off')" :value="0" />
        <el-option label="30s" :value="30" />
        <el-option label="60s" :value="60" />
        <el-option label="120s" :value="120" />
      </el-select>
    </div>

    <!-- 标签页切换 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane :label="$t('ticketList.pending')" name="open" />
      <el-tab-pane :label="$t('ticketList.replied')" name="replied" />
      <el-tab-pane :label="$t('ticketList.waiting')" name="in_progress" />
      <el-tab-pane :label="$t('ticketList.closed')" name="closed" />
      <el-tab-pane :label="$t('ticketList.all')" name="all" />
    </el-tabs>

    <!-- 操作按钮 -->
    <div class="action-bar">
      <div class="action-left">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>{{ $t('ticketList.createTicket') }}
        </el-button>
        <el-button @click="showAdvancedSearch = !showAdvancedSearch">
          {{ $t('ticketList.advancedSearch') }}
        </el-button>
      </div>
    </div>

    <!-- 高级搜索 -->
    <el-card v-if="showAdvancedSearch" shadow="never" class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item :label="$t('ticketList.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('ticketList.keywordPlaceholder')" clearable style="width: 200px" @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item :label="$t('ticketList.priority')">
          <el-select v-model="searchForm.priority" :placeholder="$t('common.all')" clearable style="width: 100px">
            <el-option :label="$t('ticketList.low')" :value="1" />
            <el-option :label="$t('ticketList.normal')" :value="2" />
            <el-option :label="$t('ticketList.high')" :value="3" />
            <el-option :label="$t('ticketList.urgent')" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ticketList.department')">
          <el-select v-model="searchForm.department_id" :placeholder="$t('common.all')" clearable style="width: 120px">
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-table
      v-loading="loading"
      :data="tableData"
      border
      stripe
      :row-class-name="getRowClassName"
      @selection-change="handleSelectionChange"
      @sort-change="handleSortChange"
    >
      <el-table-column type="selection" width="50" />
      <el-table-column prop="id" label="ID" width="80" sortable="custom" align="center">
        <template #default="{ row }">
          <el-link type="primary" @click="handleView(row)">{{ row.id }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="subject" :label="$t('ticketList.ticketTitle')" min-width="250">
        <template #default="{ row }">
          <div class="subject-cell">
            <el-link type="primary" @click="handleView(row)" :class="{ 'font-bold': row.is_unread }">
              {{ row.subject }}
            </el-link>
            <span v-if="row.is_unread" class="new-badge">new</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="client_name" :label="$t('ticketList.submitter')" width="120">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleViewClient(row)">{{ row.client_name }}</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="status" :label="$t('ticketList.status')" width="100" align="center" sortable="custom">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="admin_name" :label="$t('ticketList.handler')" width="100">
        <template #default="{ row }">{{ row.admin_name || '-' }}</template>
      </el-table-column>
      <el-table-column prop="department_name" :label="$t('ticketList.department')" width="120" />
      <el-table-column prop="created_at" :label="$t('ticketList.submitTime')" width="170" sortable="custom" />
      <el-table-column prop="last_reply_at" :label="$t('ticketList.lastReply')" width="120">
        <template #default="{ row }">
          <span v-if="row.last_reply_at">{{ formatRelativeTime(row.last_reply_at) }}</span>
          <span v-else>-</span>
        </template>
      </el-table-column>
    </el-table>

    <!-- 批量操作栏 -->
    <div v-if="selectedRows.length > 0" class="batch-bar">
      <span>{{ $t('ticketList.selectedItems') }}: {{ selectedRows.length }}</span>
      <el-button size="small" @click="handleBatchClose">{{ $t('ticketList.close') }}</el-button>
      <el-button size="small" type="danger" @click="handleBatchDelete">{{ $t('common.delete') }}</el-button>
    </div>

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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/http'
import { $t } from '@/locales'

const router = useRouter()
const loading = ref(false)
const tableData = ref([])
const activeTab = ref('open')
const departments = ref<{ id: number; name: string }[]>([])
const showAdvancedSearch = ref(false)
const selectedRows = ref<any[]>([])
const refreshInterval = ref(0)

const searchForm = reactive({ keyword: '', priority: null as number | null, department_id: null as number | null })
const pagination = reactive({ page: 1, page_size: 100, total: 0 })
const sortParams = reactive({ sort: 'last_reply_at', order: 'desc' })

const hasUnread = computed(() => tableData.value.some((row: any) => row.is_unread))

let pollTimer: ReturnType<typeof setInterval> | null = null
let notificationSound: HTMLAudioElement | null = null

const initNotificationSound = () => {
  notificationSound = new Audio('/ticket_music.3882b501.mp3')
  notificationSound.volume = 0.5
}

const playNotificationSound = () => {
  if (!notificationSound) return
  let count = 0
  const play = () => {
    if (count >= 3) return
    notificationSound!.currentTime = 0
    notificationSound!.play().catch(() => {})
    count++
  }
  play()
  const interval = setInterval(() => {
    if (count >= 3) { clearInterval(interval); return }
    play()
  }, 1500)
}

const getRowClassName = ({ row }: { row: any }) => row.is_unread ? 'unread-row' : ''

const getStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { open: 'warning', in_progress: 'primary', replied: 'success', closed: 'info' }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, () => string> = {
    open: () => $t('ticketList.pending'),
    in_progress: () => $t('ticketList.waiting'),
    replied: () => $t('ticketList.replied'),
    closed: () => $t('ticketList.closed')
  }
  return map[status]?.() || $t('common.unknown')
}

const formatRelativeTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const now = new Date()
  const diff = Math.floor((now.getTime() - date.getTime()) / 1000)
  if (diff < 60) return $t('ticketList.justNow')
  if (diff < 3600) return Math.floor(diff / 60) + $t('ticketList.minutesAgo')
  if (diff < 86400) return Math.floor(diff / 3600) + $t('ticketList.hoursAgo')
  if (diff < 2592000) return Math.floor(diff / 86400) + $t('ticketList.daysAgo')
  return dateStr
}

const handleTabChange = () => { pagination.page = 1; fetchList() }

const fetchList = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size, sort: sortParams.sort, order: sortParams.order }
    if (activeTab.value !== 'all') params.status = activeTab.value
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.priority) params.priority = searchForm.priority
    if (searchForm.department_id) params.department_id = searchForm.department_id
    const data = await request.get({ url: '/api/admin/tickets', params })
    tableData.value = data?.list || data || []
    pagination.total = data?.total || 0
  } catch (error) {
    console.error('fetch ticket list failed:', error)
  } finally { loading.value = false }
}

const fetchDepartments = async () => {
  try {
    const data = await request.get({ url: '/api/admin/ticket-departments' })
    departments.value = data || []
  } catch {}
}

const handleRefreshChange = (val: number) => {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
  if (val > 0) {
    pollTimer = setInterval(() => { fetchList() }, val * 1000)
  }
}

const handleSelectionChange = (rows: any[]) => { selectedRows.value = rows }
const handleSearch = () => { pagination.page = 1; fetchList() }
const handleReset = () => { searchForm.keyword = ''; searchForm.priority = null; searchForm.department_id = null; pagination.page = 1; fetchList() }
const handleSortChange = ({ prop, order }: any) => { sortParams.sort = prop || 'last_reply_at'; sortParams.order = order === 'ascending' ? 'asc' : 'desc'; fetchList() }
const handleSizeChange = (size: number) => { pagination.page_size = size; pagination.page = 1; fetchList() }
const handlePageChange = (page: number) => { pagination.page = page; fetchList() }
const handleAdd = () => { router.push('/add-support-ticket') }
const handleView = (row: any) => { router.push(`/support-ticket-detail?id=${row.id}&tid=${row.ticket_no}`) }
const handleViewClient = (row: any) => { router.push(`/customer-view/abstract?id=${row.client_id}`) }

const handleBatchClose = async () => {
  try {
    await ElMessageBox.confirm($t('ticketList.confirmBatchClose'), $t('common.confirm'), { type: 'warning' })
    for (const row of selectedRows.value) {
      await request.post({ url: `/api/admin/tickets/${row.id}/close` })
    }
    ElMessage.success($t('ticketList.closedSuccess'))
    fetchList()
  } catch (e) { if (e !== 'cancel') console.error(e) }
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm($t('ticketList.confirmBatchDelete'), $t('common.confirm'), { type: 'warning' })
    for (const row of selectedRows.value) {
      await request.del({ url: `/api/admin/tickets/${row.id}` })
    }
    ElMessage.success($t('common.deleteSuccess'))
    fetchList()
  } catch (e) { if (e !== 'cancel') console.error(e) }
}

onMounted(() => {
  initNotificationSound()
  fetchList()
  fetchDepartments()
})
onUnmounted(() => { if (pollTimer) clearInterval(pollTimer) })
</script>

<style scoped lang="scss">
.ticket-list-page { padding: 16px; }
.page-title { margin: 0 0 12px; font-size: 14px; font-weight: 500; color: #4e5969; }
.unread-alert { margin-bottom: 12px; }
.refresh-bar { display: flex; align-items: center; gap: 8px; margin-bottom: 12px; font-size: 13px; color: #86909c; }
.action-bar { display: flex; justify-content: space-between; margin-bottom: 12px; }
.action-left { display: flex; gap: 8px; }
.search-card { margin-bottom: 12px; }
.subject-cell { display: flex; align-items: center; gap: 6px; }
.font-bold { font-weight: bold; }
.new-badge { padding: 1px 6px; background: #ff7d00; color: white; border-radius: 10px; font-size: 11px; flex-shrink: 0; }
.batch-bar { display: flex; align-items: center; gap: 12px; padding: 8px 16px; margin-top: 12px; background: #f5f5f5; border-radius: 4px; font-size: 13px; }
.pagination-wrapper { display: flex; justify-content: flex-end; margin-top: 12px; }
:deep(.el-table .unread-row) { background: #fffbe6 !important; }
</style>
