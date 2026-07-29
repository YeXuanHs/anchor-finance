<template>
  <div class="tickets-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="工单号">
          <el-input v-model="searchForm.ticket_no" placeholder="工单号" clearable />
        </el-form-item>
        <el-form-item label="标题">
          <el-input v-model="searchForm.title" placeholder="工单标题" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待处理" value="open" />
            <el-option label="处理中" value="processing" />
            <el-option label="已回复" value="replied" />
            <el-option label="已关闭" value="closed" />
          </el-select>
        </el-form-item>
        <el-form-item label="优先级">
          <el-select v-model="searchForm.priority" placeholder="全部" clearable>
            <el-option label="紧急" value="urgent" />
            <el-option label="高" value="high" />
            <el-option label="中" value="medium" />
            <el-option label="低" value="low" />
          </el-select>
        </el-form-item>
        <el-form-item label="部门">
          <el-select v-model="searchForm.department_id" placeholder="全部" clearable>
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>工单列表</h3>
        <div>
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
      </div>

      <el-table :data="tickets" style="width: 100%" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="ticket_no" label="工单号" width="150" />
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="user.username" label="用户" width="120" />
        <el-table-column prop="department.name" label="部门" width="120" />
        <el-table-column prop="priority" label="优先级" width="100">
          <template #default="{ row }">
            <el-tag :type="getPriorityType(row.priority)" size="small">
              {{ getPriorityText(row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="replies_count" label="回复数" width="80" align="center" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column prop="updated_at" label="最后更新" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewTicket(row)">详情</el-button>
            <el-button type="primary" link @click="replyTicket(row)" v-if="row.status !== 'closed'">回复</el-button>
            <el-button type="danger" link @click="closeTicket(row)" v-if="row.status !== 'closed'">关闭</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchTickets"
          @current-change="fetchTickets"
        />
      </div>
    </div>

    <el-dialog v-model="showReplyDialog" title="回复工单" width="600px">
      <el-form :model="replyForm" label-width="80px">
        <el-form-item label="工单号">
          <el-input :model-value="currentTicket?.ticket_no" disabled />
        </el-form-item>
        <el-form-item label="标题">
          <el-input :model-value="currentTicket?.title" disabled />
        </el-form-item>
        <el-form-item label="回复内容">
          <el-input v-model="replyForm.content" type="textarea" :rows="6" placeholder="请输入回复内容" />
        </el-form-item>
        <el-form-item label="更改状态">
          <el-select v-model="replyForm.status" placeholder="保持当前状态" clearable>
            <el-option label="处理中" value="processing" />
            <el-option label="已回复" value="replied" />
            <el-option label="已关闭" value="closed" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showReplyDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleReply">发送回复</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download } from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()
const loading = ref(false)
const submitLoading = ref(false)
const tickets = ref<any[]>([])
const departments = ref<any[]>([])
const selectedRows = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showReplyDialog = ref(false)
const currentTicket = ref<any>(null)

const searchForm = ref({
  ticket_no: '',
  title: '',
  status: '',
  priority: '',
  department_id: ''
})

const replyForm = ref({
  content: '',
  status: ''
})

const getStatusText = (status: string) => {
  const map: Record<string, string> = { open: '待处理', processing: '处理中', replied: '已回复', closed: '已关闭' }
  return map[status] || status
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = { open: 'warning', processing: 'primary', replied: 'success', closed: 'info' }
  return map[status] || 'info'
}

const getPriorityText = (priority: string) => {
  const map: Record<string, string> = { urgent: '紧急', high: '高', medium: '中', low: '低' }
  return map[priority] || priority
}

const getPriorityType = (priority: string) => {
  const map: Record<string, string> = { urgent: 'danger', high: 'warning', medium: '', low: 'info' }
  return map[priority] || 'info'
}

const fetchTickets = async () => {
  loading.value = true
  try {
    const params = { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    const { data } = await request.get('/admin/api/v1/tickets', { params })
    tickets.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取工单列表失败')
  } finally {
    loading.value = false
  }
}

const fetchDepartments = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/tickets/departments', { params: { page_size: 100 } })
    departments.value = data.data?.list || []
  } catch {}
}

const handleSearch = () => {
  currentPage.value = 1
  fetchTickets()
}

const resetSearch = () => {
  searchForm.value = { ticket_no: '', title: '', status: '', priority: '', department_id: '' }
  handleSearch()
}

const handleSelectionChange = (rows: any[]) => {
  selectedRows.value = rows
}

const viewTicket = (ticket: any) => {
  router.push(`/tickets/detail/${ticket.id}`)
}

const replyTicket = (ticket: any) => {
  currentTicket.value = ticket
  replyForm.value = { content: '', status: '' }
  showReplyDialog.value = true
}

const handleReply = async () => {
  if (!replyForm.value.content.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }
  submitLoading.value = true
  try {
    await request.post(`/admin/api/v1/tickets/${currentTicket.value.id}/reply`, replyForm.value)
    ElMessage.success('回复成功')
    showReplyDialog.value = false
    fetchTickets()
  } catch {
    ElMessage.error('回复失败')
  } finally {
    submitLoading.value = false
  }
}

const closeTicket = async (ticket: any) => {
  try {
    await ElMessageBox.confirm('确定要关闭该工单吗？', '提示', { type: 'warning' })
    await request.put(`/admin/api/v1/tickets/${ticket.id}/close`)
    ElMessage.success('工单已关闭')
    fetchTickets()
  } catch {}
}

const handleExport = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/tickets/export', { params: searchForm.value, responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([data]))
    const a = document.createElement('a')
    a.href = url
    a.download = `tickets_${Date.now()}.xlsx`
    a.click()
    window.URL.revokeObjectURL(url)
  } catch {
    ElMessage.error('导出失败')
  }
}

onMounted(() => {
  fetchTickets()
  fetchDepartments()
})
</script>

<style scoped lang="scss">
.tickets-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
}
</style>
