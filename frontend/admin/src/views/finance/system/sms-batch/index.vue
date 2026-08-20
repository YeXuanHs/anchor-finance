<template>
  <div class="sms-batch-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('smsBatch.title') }}</span>
          <el-button type="primary" @click="handleCreateBatch">
            <el-icon><Plus /></el-icon>
            {{ $t('smsBatch.createBatch') }}
          </el-button>
        </div>
      </template>

      <!-- 统计卡片 -->
      <el-row :gutter="16" class="stats-row">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">{{ $t('smsBatch.totalBatches') }}</div>
            <div class="stat-value">{{ stats.total || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">{{ $t('smsBatch.sent') }}</div>
            <div class="stat-value success">{{ stats.sent || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">{{ $t('smsBatch.successRate') }}</div>
            <div class="stat-value primary">{{ stats.success_rate || '0%' }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">{{ $t('smsBatch.todaySent') }}</div>
            <div class="stat-value warning">{{ stats.today_sent || 0 }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('smsBatch.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('smsBatch.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('smsBatch.sendStatus')">
          <el-select v-model="searchForm.status" :placeholder="$t('smsBatch.all')" clearable>
            <el-option :label="$t('smsBatch.pending')" value="pending" />
            <el-option :label="$t('smsBatch.sending')" value="sending" />
            <el-option :label="$t('smsBatch.completed')" value="completed" />
            <el-option :label="$t('smsBatch.failed')" value="failed" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('smsBatch.smsType')">
          <el-select v-model="searchForm.sms_type" :placeholder="$t('smsBatch.all')" clearable>
            <el-option :label="$t('smsBatch.marketing')" value="marketing" />
            <el-option :label="$t('smsBatch.notification')" value="notification" />
            <el-option :label="$t('smsBatch.verification')" value="verification" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('smsBatch.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('smsBatch.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" :label="$t('smsBatch.id')" width="70" align="center" />
        <el-table-column prop="name" :label="$t('smsBatch.batchName')" min-width="150" />
        <el-table-column prop="sms_type" :label="$t('smsBatch.smsTypeColumn')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="smsTypeTagMap[row.sms_type]" size="small">
              {{ smsTypeTextMap[row.sms_type] || row.sms_type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="template_name" :label="$t('smsBatch.templateUsed')" width="120" />
        <el-table-column prop="total_count" :label="$t('smsBatch.targetCount')" width="90" align="center" />
        <el-table-column prop="sent_count" :label="$t('smsBatch.sentCount')" width="90" align="center" />
        <el-table-column prop="success_count" :label="$t('smsBatch.successCount')" width="90" align="center">
          <template #default="{ row }">
            <span class="success-text">{{ row.success_count || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="failed_count" :label="$t('smsBatch.failedCount')" width="90" align="center">
          <template #default="{ row }">
            <span :class="{ 'failed-text': row.failed_count > 0 }">{{ row.failed_count || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('smsBatch.status')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="statusTagMap[row.status]" size="small">
              {{ statusTextMap[row.status] || row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('smsBatch.createTime')" width="180" />
        <el-table-column :label="$t('smsBatch.operations')" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">{{ $t('smsBatch.detail') }}</el-button>
            <el-button v-if="row.status === 'pending'" type="success" link @click="handleSend(row)">{{ $t('smsBatch.send') }}</el-button>
            <el-button v-if="row.status === 'sending'" type="warning" link @click="handlePause(row)">{{ $t('smsBatch.pause') }}</el-button>
            <el-popconfirm v-if="row.status !== 'sending'" :title="$t('smsBatch.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('smsBatch.delete') }}</el-button>
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
    <el-dialog v-model="createDialogVisible" :title="$t('smsBatch.createSmsBatch')" width="700px" destroy-on-close>
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
        <el-form-item :label="$t('smsBatch.batchName')" prop="name">
          <el-input v-model="createForm.name" :placeholder="$t('smsBatch.inputBatchName')" />
        </el-form-item>
        <el-form-item :label="$t('smsBatch.smsTypeColumn')" prop="sms_type">
          <el-select v-model="createForm.sms_type" :placeholder="$t('smsBatch.selectSmsType')" style="width: 100%">
            <el-option :label="$t('smsBatch.marketing')" value="marketing" />
            <el-option :label="$t('smsBatch.notification')" value="notification" />
            <el-option :label="$t('smsBatch.verification')" value="verification" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('smsBatch.templateUsed')" prop="template_id">
          <el-select v-model="createForm.template_id" :placeholder="$t('smsBatch.selectSmsTemplate')" style="width: 100%">
            <el-option v-for="t in templateList" :key="t.id" :label="t.title" :value="t.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('smsBatch.targetCount')" prop="target_type">
          <el-radio-group v-model="createForm.target_type">
            <el-radio value="all">{{ $t('smsBatch.allUsers') }}</el-radio>
            <el-radio value="group">{{ $t('smsBatch.specificGroup') }}</el-radio>
            <el-radio value="custom">{{ $t('smsBatch.customNumbers') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="createForm.target_type === 'group'" :label="$t('smsBatch.userGroup')">
          <el-select v-model="createForm.group_id" :placeholder="$t('smsBatch.selectUserGroup')" style="width: 100%">
            <el-option v-for="g in groupList" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="createForm.target_type === 'custom'" :label="$t('smsBatch.phoneNumbers')">
          <el-input v-model="createForm.phones" type="textarea" :rows="4" :placeholder="$t('smsBatch.phonePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('smsBatch.scheduledSend')">
          <el-date-picker v-model="createForm.scheduled_at" type="datetime" :placeholder="$t('smsBatch.selectSendTime')" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('smsBatch.remark')">
          <el-input v-model="createForm.remark" type="textarea" :rows="2" :placeholder="$t('smsBatch.batchRemark')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">{{ $t('smsBatch.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmitBatch" :loading="createLoading">{{ $t('smsBatch.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 批次详情对话框 -->
    <el-dialog v-model="detailDialogVisible" :title="$t('smsBatch.batchDetail')" width="800px" destroy-on-close>
      <div v-loading="detailLoading">
        <el-descriptions :column="2" border class="detail-info">
          <el-descriptions-item :label="$t('smsBatch.batchName')">{{ currentBatch.name }}</el-descriptions-item>
          <el-descriptions-item :label="$t('smsBatch.smsTypeColumn')">{{ smsTypeTextMap[currentBatch.sms_type] }}</el-descriptions-item>
          <el-descriptions-item :label="$t('smsBatch.targetCount')">{{ currentBatch.total_count }}</el-descriptions-item>
          <el-descriptions-item :label="$t('smsBatch.sentCount')">{{ currentBatch.sent_count }}</el-descriptions-item>
          <el-descriptions-item :label="$t('smsBatch.successCount')"><span class="success-text">{{ currentBatch.success_count }}</span></el-descriptions-item>
          <el-descriptions-item :label="$t('smsBatch.failedCount')"><span :class="{ 'failed-text': currentBatch.failed_count > 0 }">{{ currentBatch.failed_count }}</span></el-descriptions-item>
          <el-descriptions-item :label="$t('smsBatch.status')">
            <el-tag :type="statusTagMap[currentBatch.status]">{{ statusTextMap[currentBatch.status] }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('smsBatch.createTime')">{{ currentBatch.created_at }}</el-descriptions-item>
        </el-descriptions>

        <!-- 运营商检测 -->
        <div class="operator-section">
          <h4>{{ $t('smsBatch.operatorDetection') }}</h4>
          <el-row :gutter="16">
            <el-col :span="8">
              <el-card shadow="hover" class="operator-card">
                <div class="operator-label">{{ $t('smsBatch.chinaMobile') }}</div>
                <div class="operator-value">{{ currentBatch.cmcc_count || 0 }}</div>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card shadow="hover" class="operator-card">
                <div class="operator-label">{{ $t('smsBatch.chinaUnicom') }}</div>
                <div class="operator-value">{{ currentBatch.cucc_count || 0 }}</div>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card shadow="hover" class="operator-card">
                <div class="operator-label">{{ $t('smsBatch.chinaTelecom') }}</div>
                <div class="operator-value">{{ currentBatch.ctcc_count || 0 }}</div>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <!-- 发送记录 -->
        <div class="send-logs-section">
          <h4>{{ $t('smsBatch.sendLogs') }}</h4>
          <el-table :data="sendLogs" border size="small" max-height="300">
            <el-table-column prop="phone" :label="$t('smsBatch.phone')" width="130" />
            <el-table-column prop="status" :label="$t('smsBatch.status')" width="80" align="center">
              <template #default="{ row }">
                <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
                  {{ row.status === 'success' ? $t('smsBatch.success') : $t('smsBatch.failure') }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="error_msg" :label="$t('smsBatch.errorMsg')" min-width="200" show-overflow-tooltip />
            <el-table-column prop="sent_at" :label="$t('smsBatch.sendTime')" width="180" />
          </el-table>
        </div>
      </div>
      <template #footer>
        <el-button @click="detailDialogVisible = false">{{ $t('smsBatch.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'SmsBatchManage' })

const smsTypeTextMap = computed<Record<string, string>>(() => ({
  marketing: $t('smsBatch.marketing'),
  notification: $t('smsBatch.notification'),
  verification: $t('smsBatch.verification')
}))
const smsTypeTagMap: Record<string, any> = { marketing: 'warning', notification: 'primary', verification: 'success' }
const statusTextMap = computed<Record<string, string>>(() => ({
  pending: $t('smsBatch.pending'),
  sending: $t('smsBatch.sending'),
  completed: $t('smsBatch.completed'),
  failed: $t('smsBatch.failed'),
  paused: $t('smsBatch.paused')
}))
const statusTagMap: Record<string, any> = { pending: 'info', sending: 'warning', completed: 'success', failed: 'danger', paused: 'info' }

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

const createRules = computed<FormRules>(() => ({
  name: [{ required: true, message: $t('smsBatch.inputBatchName'), trigger: 'blur' }],
  sms_type: [{ required: true, message: $t('smsBatch.selectSmsType'), trigger: 'change' }],
  template_id: [{ required: true, message: $t('smsBatch.selectSmsTemplate'), trigger: 'change' }],
  target_type: [{ required: true, message: $t('smsBatch.selectTargetType'), trigger: 'change' }]
}))

const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/sms/batches/stats' })
    Object.assign(stats, data)
  } catch (error) {
    console.error($t('smsBatch.fetchStatsFailed'), error)
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
    ElMessage.error($t('smsBatch.fetchBatchListFailed'))
  } finally {
    loading.value = false
  }
}

const fetchTemplates = async () => {
  try {
    const data = await request.get({ url: '/api/admin/sms/templates', params: { page_size: 100 } })
    templateList.value = data.templates || data.list || []
  } catch (error) {
    console.error($t('smsBatch.fetchTemplateFailed'), error)
  }
}

const fetchGroups = async () => {
  try {
    const data = await request.get({ url: '/api/admin/client-groups', params: { page_size: 100 } })
    groupList.value = data.list || data || []
  } catch (error) {
    console.error($t('smsBatch.fetchGroupFailed'), error)
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
      ElMessage.success($t('smsBatch.batchCreated'))
      createDialogVisible.value = false
      fetchData()
      fetchStats()
    } catch (error) {
      ElMessage.error($t('smsBatch.createFailed'))
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
    ElMessage.error($t('smsBatch.fetchDetailFailed'))
  } finally {
    detailLoading.value = false
  }
}

const handleSend = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('smsBatch.confirmSendBatch', { name: row.name }), $t('smsBatch.confirmSend'), { type: 'warning' })
    await request.post({ url: `/api/admin/sms/batches/${row.id}/send` })
    ElMessage.success($t('smsBatch.sendTaskSubmitted'))
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error($t('smsBatch.sendFailed'))
  }
}

const handlePause = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/sms/batches/${row.id}/pause` })
    ElMessage.success($t('smsBatch.pausedSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('smsBatch.pauseFailed'))
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/sms/batches/${row.id}` })
    ElMessage.success($t('smsBatch.deleteSuccess'))
    fetchData()
    fetchStats()
  } catch (error) {
    ElMessage.error($t('smsBatch.deleteFailed'))
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
