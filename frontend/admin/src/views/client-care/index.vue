<template>
  <div class="client-care-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="用户名/邮箱" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option v-for="t in careTypes" :key="t.value" :label="t.label" :value="t.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待发送" value="pending" />
            <el-option label="已发送" value="sent" />
            <el-option label="发送失败" value="failed" />
            <el-option label="已取消" value="cancelled" />
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

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card" v-for="stat in stats" :key="stat.label">
        <div class="stat-icon" :class="stat.type">
          <el-icon :size="24"><component :is="stat.icon" /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-value">{{ stat.value }}</div>
        </div>
      </div>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>客户关怀任务</h3>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>创建任务
        </el-button>
      </div>

      <el-table :data="tasks" style="width: 100%" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="username" label="用户" min-width="120" />
        <el-table-column prop="email" label="邮箱" min-width="180" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="channel" label="渠道" width="100">
          <template #default="{ row }">
            <el-tag v-for="ch in (row.channel || '').split(',')" :key="ch" size="small" class="channel-tag">
              {{ getChannelLabel(ch) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="trigger" label="触发条件" width="130" show-overflow-tooltip />
        <el-table-column prop="scheduled_at" label="计划时间" width="170" />
        <el-table-column prop="sent_at" label="发送时间" width="170" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link size="small" @click="handleViewDetail(row)">详情</el-button>
            <el-button type="success" link size="small" @click="handleSendNow(row)" v-if="row.status === 'pending'">
              立即发送
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="table-footer">
        <div class="batch-actions" v-if="selectedTasks.length">
          <span>已选 {{ selectedTasks.length }} 项</span>
          <el-button size="small" type="danger" @click="handleBatchDelete">批量删除</el-button>
        </div>
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="600px" :close-on-click-modal="false">
      <el-form :model="formData" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="用户" prop="user_id">
          <el-select v-model="formData.user_id" filterable remote :remote-method="searchUsers" placeholder="搜索用户名/邮箱" :loading="userSearching" style="width: 100%;">
            <el-option v-for="u in userOptions" :key="u.id" :label="`${u.username} (${u.email})`" :value="u.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型">
            <el-option v-for="t in careTypes" :key="t.value" :label="t.label" :value="t.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="渠道" prop="channels">
          <el-checkbox-group v-model="formData.channels">
            <el-checkbox value="email">邮件</el-checkbox>
            <el-checkbox value="sms">短信</el-checkbox>
            <el-checkbox value="wechat">微信</el-checkbox>
            <el-checkbox value="站内信">站内信</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="触发条件" prop="trigger">
          <el-select v-model="formData.trigger" placeholder="选择触发条件" clearable>
            <el-option label="手动触发" value="manual" />
            <el-option label="用户生日" value="birthday" />
            <el-option label="服务到期前7天" value="expiry_7d" />
            <el-option label="服务到期前30天" value="expiry_30d" />
            <el-option label="30天未登录" value="inactive_30d" />
            <el-option label="消费满额" value="amount_threshold" />
          </el-select>
        </el-form-item>
        <el-form-item label="计划时间">
          <el-date-picker v-model="formData.scheduled_at" type="datetime" placeholder="选择时间" style="width: 100%;" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="5" placeholder="支持变量：{username}、{email}、{expire_date}" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog title="任务详情" v-model="detailVisible" width="550px">
      <el-descriptions :column="1" border v-if="currentTask">
        <el-descriptions-item label="ID">{{ currentTask.id }}</el-descriptions-item>
        <el-descriptions-item label="用户">{{ currentTask.username }} ({{ currentTask.email }})</el-descriptions-item>
        <el-descriptions-item label="类型">{{ getTypeLabel(currentTask.type) }}</el-descriptions-item>
        <el-descriptions-item label="渠道">{{ currentTask.channel }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentTask.status)">{{ getStatusLabel(currentTask.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="触发条件">{{ currentTask.trigger || '手动' }}</el-descriptions-item>
        <el-descriptions-item label="计划时间">{{ currentTask.scheduled_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="发送时间">{{ currentTask.sent_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="内容">{{ currentTask.content }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Bell, UserFilled, Warning, CircleCheck } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const submitting = ref(false)
const tasks = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const dialogVisible = ref(false)
const detailVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref()
const currentTask = ref<any>(null)
const selectedTasks = ref<any[]>([])
const userSearching = ref(false)
const userOptions = ref<any[]>([])

const careTypes = [
  { label: '生日祝福', value: 'birthday' },
  { label: '到期提醒', value: 'expiry' },
  { label: '流失预警', value: 'churn' },
  { label: '回访', value: 'follow_up' },
  { label: '促销通知', value: 'promotion' },
  { label: '满意度调查', value: 'survey' }
]

const searchForm = ref({ keyword: '', type: '', status: '', date_range: null as any })

const stats = ref([
  { label: '待发送', value: '0', icon: 'Bell', type: 'primary' },
  { label: '已发送', value: '0', icon: 'CircleCheck', type: 'success' },
  { label: '失败', value: '0', icon: 'Warning', type: 'warning' },
  { label: '本月回访', value: '0', icon: 'UserFilled', type: 'info' }
])

const defaultForm = () => ({
  id: null,
  user_id: null,
  type: '',
  channels: ['email'],
  trigger: 'manual',
  content: '',
  scheduled_at: ''
})
const formData = reactive(defaultForm())

const rules = {
  user_id: [{ required: true, message: '请选择用户', trigger: 'change' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  channels: [{ required: true, type: 'array' as const, message: '请选择渠道', trigger: 'change' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const getTypeTag = (type: string) => {
  const map: Record<string, string> = { birthday: 'success', expiry: 'warning', churn: 'danger', follow_up: '', promotion: 'warning', survey: 'info' }
  return map[type] || 'info'
}
const getTypeLabel = (val: string) => careTypes.find(t => t.value === val)?.label || val
const getChannelLabel = (ch: string) => ({ email: '邮件', sms: '短信', wechat: '微信', '站内信': '站内信' } as Record<string, string>)[ch] || ch
const getStatusType = (status: string) => ({ pending: 'warning', sent: 'success', failed: 'danger', cancelled: 'info' } as Record<string, string>)[status] || 'info'
const getStatusLabel = (status: string) => ({ pending: '待发送', sent: '已发送', failed: '失败', cancelled: '已取消' } as Record<string, string>)[status] || status

const fetchTasks = async () => {
  loading.value = true
  try {
    const params: any = { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    delete params.date_range
    if (searchForm.value.date_range?.length === 2) {
      params.start_date = searchForm.value.date_range[0]
      params.end_date = searchForm.value.date_range[1]
    }
    const { data } = await request.get('/admin/api/v1/client-care', { params })
    tasks.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取列表失败')
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/client-care/stats')
    if (data.data) {
      stats.value[0].value = String(data.data.pending || 0)
      stats.value[1].value = String(data.data.sent || 0)
      stats.value[2].value = String(data.data.failed || 0)
      stats.value[3].value = String(data.data.follow_up_month || 0)
    }
  } catch {}
}

const searchUsers = async (query: string) => {
  if (!query) return
  userSearching.value = true
  try {
    const { data } = await request.get('/admin/api/v1/users/search', { params: { keyword: query } })
    userOptions.value = data.data?.list || []
  } catch {} finally {
    userSearching.value = false
  }
}

const handleSearch = () => { currentPage.value = 1; fetchTasks() }
const resetSearch = () => { searchForm.value = { keyword: '', type: '', status: '', date_range: null }; handleSearch() }
const handleSizeChange = (val: number) => { pageSize.value = val; fetchTasks() }
const handlePageChange = (val: number) => { currentPage.value = val; fetchTasks() }
const handleSelectionChange = (sel: any[]) => { selectedTasks.value = sel }

const handleCreate = () => {
  Object.assign(formData, defaultForm())
  dialogTitle.value = '创建关怀任务'
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  Object.assign(formData, { ...row, channels: row.channel ? row.channel.split(',') : ['email'] })
  dialogTitle.value = '编辑关怀任务'
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try { await formRef.value?.validate() } catch { return }
  submitting.value = true
  try {
    const payload = { ...formData, channel: formData.channels.join(',') }
    if (formData.id) {
      await request.put(`/admin/api/v1/client-care/${formData.id}`, payload)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/client-care', payload)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchTasks()
    fetchStats()
  } catch {
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
  }
}

const handleViewDetail = (row: any) => { currentTask.value = row; detailVisible.value = true }

const handleSendNow = async (row: any) => {
  try {
    await ElMessageBox.confirm('确认立即发送该关怀消息？', '确认', { type: 'info' })
    await request.post(`/admin/api/v1/client-care/${row.id}/send`)
    ElMessage.success('发送成功')
    fetchTasks()
    fetchStats()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error('发送失败')
  }
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(`确认删除该关怀任务？`, '确认删除', { type: 'warning' }).then(async () => {
    try {
      await request.delete(`/admin/api/v1/client-care/${row.id}`)
      ElMessage.success('删除成功')
      fetchTasks()
      fetchStats()
    } catch {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

const handleBatchDelete = () => {
  ElMessageBox.confirm(`确认删除选中的 ${selectedTasks.value.length} 个任务？`, '批量删除', { type: 'warning' }).then(async () => {
    try {
      await request.post('/admin/api/v1/client-care/batch-delete', { ids: selectedTasks.value.map((t: any) => t.id) })
      ElMessage.success('删除成功')
      fetchTasks()
      fetchStats()
    } catch {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

onMounted(() => { fetchTasks(); fetchStats() })
</script>

<style scoped lang="scss">
.client-care-page {
  .stats-grid {
    display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin-bottom: 20px;
    .stat-card {
      background: var(--bg-card); border-radius: var(--border-radius); padding: 20px;
      display: flex; align-items: center; gap: 16px; box-shadow: var(--shadow-sm);
      .stat-icon {
        width: 48px; height: 48px; border-radius: 12px;
        display: flex; align-items: center; justify-content: center;
        &.primary { background: var(--primary-bg, rgba(64,158,255,0.1)); color: var(--primary-color, #409eff); }
        &.success { background: rgba(52, 199, 89, 0.1); color: var(--success-color, #34c759); }
        &.warning { background: rgba(255, 149, 0, 0.1); color: var(--warning-color, #ff9500); }
        &.info { background: rgba(142, 142, 147, 0.1); color: var(--info-color, #8e8e93); }
      }
      .stat-info {
        .stat-label { font-size: 13px; color: var(--text-secondary); }
        .stat-value { font-size: 24px; font-weight: 600; color: var(--text-primary); }
      }
    }
  }
  .table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .channel-tag { margin-right: 4px; }
  .table-footer { margin-top: 16px; display: flex; justify-content: space-between; align-items: center;
    .batch-actions { display: flex; align-items: center; gap: 12px; font-size: 13px; color: var(--text-secondary); }
  }
}
</style>
