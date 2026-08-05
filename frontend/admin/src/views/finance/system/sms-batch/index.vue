<template>
  <div class="sms-batch-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>短信批量发送管理</span>
          <el-button type="primary" @click="handleCreateBatch">
            <el-icon><Plus /></el-icon>
            创建批次
          </el-button>
        </div>
      </template>

      <!-- 统计卡片 -->
      <el-row :gutter="16" class="stats-row">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">总批次</div>
            <div class="stat-value">{{ stats.total || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">已发送</div>
            <div class="stat-value success">{{ stats.sent || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">发送成功率</div>
            <div class="stat-value primary">{{ stats.success_rate || '0%' }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">今日发送</div>
            <div class="stat-value warning">{{ stats.today_sent || 0 }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="批次名称">
          <el-input v-model="searchForm.keyword" placeholder="批次名称" clearable />
        </el-form-item>
        <el-form-item label="发送状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待发送" value="pending" />
            <el-option label="发送中" value="sending" />
            <el-option label="已完成" value="completed" />
            <el-option label="已失败" value="failed" />
          </el-select>
        </el-form-item>
        <el-form-item label="短信类型">
          <el-select v-model="searchForm.sms_type" placeholder="全部" clearable>
            <el-option label="营销短信" value="marketing" />
            <el-option label="通知短信" value="notification" />
            <el-option label="验证码" value="verification" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="批次名称" min-width="150" />
        <el-table-column prop="sms_type" label="短信类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="smsTypeTagMap[row.sms_type]" size="small">
              {{ smsTypeTextMap[row.sms_type] || row.sms_type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="template_name" label="使用模板" width="120" />
        <el-table-column prop="total_count" label="目标数量" width="90" align="center" />
        <el-table-column prop="sent_count" label="已发送" width="90" align="center" />
        <el-table-column prop="success_count" label="成功数" width="90" align="center">
          <template #default="{ row }">
            <span class="success-text">{{ row.success_count || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="failed_count" label="失败数" width="90" align="center">
          <template #default="{ row }">
            <span :class="{ 'failed-text': row.failed_count > 0 }">{{ row.failed_count || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="statusTagMap[row.status]" size="small">
              {{ statusTextMap[row.status] || row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
            <el-button v-if="row.status === 'pending'" type="success" link @click="handleSend(row)">发送</el-button>
            <el-button v-if="row.status === 'sending'" type="warning" link @click="handlePause(row)">暂停</el-button>
            <el-popconfirm v-if="row.status !== 'sending'" title="确定删除该批次吗？" @confirm="handleDelete(row)">
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

    <!-- 创建批次对话框 -->
    <el-dialog v-model="createDialogVisible" title="创建短信批次" width="700px" destroy-on-close>
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
        <el-form-item label="批次名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入批次名称" />
        </el-form-item>
        <el-form-item label="短信类型" prop="sms_type">
          <el-select v-model="createForm.sms_type" placeholder="请选择短信类型" style="width: 100%">
            <el-option label="营销短信" value="marketing" />
            <el-option label="通知短信" value="notification" />
            <el-option label="验证码" value="verification" />
          </el-select>
        </el-form-item>
        <el-form-item label="短信模板" prop="template_id">
          <el-select v-model="createForm.template_id" placeholder="请选择短信模板" style="width: 100%">
            <el-option v-for="t in templateList" :key="t.id" :label="t.title" :value="t.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标用户" prop="target_type">
          <el-radio-group v-model="createForm.target_type">
            <el-radio value="all">全部用户</el-radio>
            <el-radio value="group">指定分组</el-radio>
            <el-radio value="custom">自定义号码</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="createForm.target_type === 'group'" label="用户分组">
          <el-select v-model="createForm.group_id" placeholder="请选择用户分组" style="width: 100%">
            <el-option v-for="g in groupList" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="createForm.target_type === 'custom'" label="手机号码">
          <el-input v-model="createForm.phones" type="textarea" :rows="4" placeholder="每行一个手机号码" />
        </el-form-item>
        <el-form-item label="定时发送">
          <el-date-picker v-model="createForm.scheduled_at" type="datetime" placeholder="选择发送时间（留空立即发送）" style="width: 100%" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="createForm.remark" type="textarea" :rows="2" placeholder="批次备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitBatch" :loading="createLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 批次详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="批次详情" width="800px" destroy-on-close>
      <div v-loading="detailLoading">
        <el-descriptions :column="2" border class="detail-info">
          <el-descriptions-item label="批次名称">{{ currentBatch.name }}</el-descriptions-item>
          <el-descriptions-item label="短信类型">{{ smsTypeTextMap[currentBatch.sms_type] }}</el-descriptions-item>
          <el-descriptions-item label="目标数量">{{ currentBatch.total_count }}</el-descriptions-item>
          <el-descriptions-item label="已发送">{{ currentBatch.sent_count }}</el-descriptions-item>
          <el-descriptions-item label="成功数"><span class="success-text">{{ currentBatch.success_count }}</span></el-descriptions-item>
          <el-descriptions-item label="失败数"><span :class="{ 'failed-text': currentBatch.failed_count > 0 }">{{ currentBatch.failed_count }}</span></el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="statusTagMap[currentBatch.status]">{{ statusTextMap[currentBatch.status] }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ currentBatch.created_at }}</el-descriptions-item>
        </el-descriptions>

        <!-- 运营商检测 -->
        <div class="operator-section">
          <h4>运营商检测</h4>
          <el-row :gutter="16">
            <el-col :span="8">
              <el-card shadow="hover" class="operator-card">
                <div class="operator-label">中国移动</div>
                <div class="operator-value">{{ currentBatch.cmcc_count || 0 }}</div>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card shadow="hover" class="operator-card">
                <div class="operator-label">中国联通</div>
                <div class="operator-value">{{ currentBatch.cucc_count || 0 }}</div>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card shadow="hover" class="operator-card">
                <div class="operator-label">中国电信</div>
                <div class="operator-value">{{ currentBatch.ctcc_count || 0 }}</div>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <!-- 发送记录 -->
        <div class="send-logs-section">
          <h4>发送记录</h4>
          <el-table :data="sendLogs" border size="small" max-height="300">
            <el-table-column prop="phone" label="手机号" width="130" />
            <el-table-column prop="status" label="状态" width="80" align="center">
              <template #default="{ row }">
                <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
                  {{ row.status === 'success' ? '成功' : '失败' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="error_msg" label="错误信息" min-width="200" show-overflow-tooltip />
            <el-table-column prop="sent_at" label="发送时间" width="180" />
          </el-table>
        </div>
      </div>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'SmsBatchManage' })

const smsTypeTextMap: Record<string, string> = { marketing: '营销短信', notification: '通知短信', verification: '验证码' }
const smsTypeTagMap: Record<string, string> = { marketing: 'warning', notification: 'primary', verification: 'success' }
const statusTextMap: Record<string, string> = { pending: '待发送', sending: '发送中', completed: '已完成', failed: '已失败', paused: '已暂停' }
const statusTagMap: Record<string, string> = { pending: 'info', sending: 'warning', completed: 'success', failed: 'danger', paused: 'info' }

const loading = ref(false)
const createLoading = ref(false)
const detailLoading = ref(false)
const createDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const createFormRef = ref<FormInstance>()

const stats = reactive({ total: 0, sent: 0, success_rate: '0%', today_sent: 0 })
const searchForm = reactive({ keyword: '', status: '', sms_type: '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const templateList = ref<any[]>([])
const groupList = ref<any[]>([])
const sendLogs = ref<any[]>([])
const currentBatch = reactive<any>({})

const createForm = reactive({
  name: '',
  sms_type: 'marketing',
  template_id: '',
  target_type: 'all',
  group_id: '',
  phones: '',
  scheduled_at: '',
  remark: ''
})

const createRules: FormRules = {
  name: [{ required: true, message: '请输入批次名称', trigger: 'blur' }],
  sms_type: [{ required: true, message: '请选择短信类型', trigger: 'change' }],
  template_id: [{ required: true, message: '请选择短信模板', trigger: 'change' }],
  target_type: [{ required: true, message: '请选择目标用户类型', trigger: 'change' }]
}

const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/sms/batches/stats' })
    Object.assign(stats, data)
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/sms/batches',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取批次列表失败')
  } finally {
    loading.value = false
  }
}

const fetchTemplates = async () => {
  try {
    const data = await request.get({ url: '/api/admin/sms/templates', params: { page_size: 100 } })
    templateList.value = data.templates || data.list || []
  } catch (error) {
    console.error('获取模板列表失败:', error)
  }
}

const fetchGroups = async () => {
  try {
    const data = await request.get({ url: '/api/admin/client-groups', params: { page_size: 100 } })
    groupList.value = data.list || data || []
  } catch (error) {
    console.error('获取分组列表失败:', error)
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { Object.assign(searchForm, { keyword: '', status: '', sms_type: '' }); handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleCreateBatch = () => {
  Object.assign(createForm, { name: '', sms_type: 'marketing', template_id: '', target_type: 'all', group_id: '', phones: '', scheduled_at: '', remark: '' })
  createDialogVisible.value = true
}

const handleSubmitBatch = async () => {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (!valid) return
    createLoading.value = true
    try {
      await request.post({ url: '/api/admin/sms/batches', params: { ...createForm } })
      ElMessage.success('批次创建成功')
      createDialogVisible.value = false
      fetchData()
      fetchStats()
    } catch (error) {
      ElMessage.error('创建失败')
    } finally {
      createLoading.value = false
    }
  })
}

const handleViewDetail = async (row: any) => {
  Object.assign(currentBatch, row)
  detailDialogVisible.value = true
  detailLoading.value = true
  try {
    const data = await request.get({ url: `/api/admin/sms/batches/${row.id}` })
    Object.assign(currentBatch, data)
    sendLogs.value = data.send_logs || []
  } catch (error) {
    ElMessage.error('获取详情失败')
  } finally {
    detailLoading.value = false
  }
}

const handleSend = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定发送批次 "${row.name}" 吗？`, '确认发送', { type: 'warning' })
    await request.post({ url: `/api/admin/sms/batches/${row.id}/send` })
    ElMessage.success('发送任务已提交')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('发送失败')
  }
}

const handlePause = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/sms/batches/${row.id}/pause` })
    ElMessage.success('已暂停')
    fetchData()
  } catch (error) {
    ElMessage.error('暂停失败')
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/sms/batches/${row.id}` })
    ElMessage.success('删除成功')
    fetchData()
    fetchStats()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

onMounted(() => { fetchStats(); fetchData(); fetchTemplates(); fetchGroups() })
</script>

<style scoped lang="scss">
.sms-batch-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.stats-row { margin-bottom: 20px; }
.stat-card { text-align: center; }
.stat-label { font-size: 13px; color: var(--el-text-color-secondary); margin-bottom: 8px; }
.stat-value {
  font-size: 28px; font-weight: 600; color: var(--el-text-color-primary);
  &.success { color: var(--el-color-success); }
  &.primary { color: var(--el-color-primary); }
  &.warning { color: var(--el-color-warning); }
}
.success-text { color: var(--el-color-success); font-weight: 600; }
.failed-text { color: var(--el-color-danger); font-weight: 600; }
.detail-info { margin-bottom: 20px; }
.operator-section { margin-top: 20px; h4 { margin: 0 0 12px; font-size: 15px; font-weight: 600; } }
.operator-card { text-align: center; }
.operator-label { font-size: 13px; color: var(--el-text-color-secondary); margin-bottom: 8px; }
.operator-value { font-size: 22px; font-weight: 600; color: var(--el-text-color-primary); }
.send-logs-section { margin-top: 20px; h4 { margin: 0 0 12px; font-size: 15px; font-weight: 600; } }
</style>
