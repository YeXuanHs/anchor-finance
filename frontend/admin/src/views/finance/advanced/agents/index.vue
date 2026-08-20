<template>
  <div class="agent-list-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('page.agents.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('page.agents.addAgent') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('page.agents.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('page.agents.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('page.agents.agentLevel')">
          <el-select v-model="searchForm.level" :placeholder="$t('page.agents.all')" clearable>
            <el-option :label="$t('page.agents.level1')" :value="1" />
            <el-option :label="$t('page.agents.level2')" :value="2" />
            <el-option :label="$t('page.agents.level3')" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('page.agents.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('page.agents.all')" clearable>
            <el-option :label="$t('page.agents.statusActive')" :value="1" />
            <el-option :label="$t('page.agents.statusDisabled')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('page.agents.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('page.agents.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="username" :label="$t('page.agents.username')" width="120" />
        <el-table-column prop="email" :label="$t('page.agents.email')" min-width="180" show-overflow-tooltip />
        <el-table-column prop="phone" :label="$t('page.agents.phone')" width="130" />
        <el-table-column prop="level" :label="$t('page.agents.agentLevel')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getAgentLevelType(row.level)" size="small">
              {{ getAgentLevelText(row.level) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="commission_rate" :label="$t('page.agents.commissionRate')" width="100" align="center">
          <template #default="{ row }">
            <span class="commission-text">{{ row.commission_rate }}%</span>
          </template>
        </el-table-column>
        <el-table-column prop="balance" :label="$t('page.agents.balance')" width="110" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.balance) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="sub_agent_count" :label="$t('page.agents.subAgentCount')" width="90" align="center" />
        <el-table-column prop="client_count" :label="$t('page.agents.clientCount')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('page.agents.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('page.agents.statusActive') : $t('page.agents.statusDisabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('page.agents.createTime')" width="170" />
        <el-table-column :label="$t('page.agents.actions')" width="250" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">{{ $t('page.agents.detail') }}</el-button>
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('page.agents.edit') }}</el-button>
            <el-button type="warning" link @click="handleAdjustCommission(row)">{{ $t('page.agents.adjustCommission') }}</el-button>
            <el-popconfirm
              :title="row.status === 1 ? $t('page.agents.confirmDisable') : $t('page.agents.confirmEnable')"
              @confirm="handleToggleStatus(row)"
            >
              <template #reference>
                <el-button :type="row.status === 1 ? 'danger' : 'success'" link>
                  {{ row.status === 1 ? $t('page.agents.disable') : $t('page.agents.enable') }}
                </el-button>
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

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('page.agents.username')" prop="username">
              <el-input v-model="formData.username" :placeholder="$t('page.agents.usernamePlaceholder')" :disabled="!!formData.id" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('page.agents.email')" prop="email">
              <el-input v-model="formData.email" :placeholder="$t('page.agents.emailPlaceholder')" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('page.agents.phone')" prop="phone">
              <el-input v-model="formData.phone" :placeholder="$t('page.agents.phonePlaceholder')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('page.agents.password')" prop="password" v-if="!formData.id">
              <el-input v-model="formData.password" type="password" :placeholder="$t('page.agents.passwordPlaceholder')" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('page.agents.agentLevel')" prop="level">
              <el-select v-model="formData.level" :placeholder="$t('page.agents.selectLevel')" style="width: 100%">
                <el-option :label="$t('page.agents.level1')" :value="1" />
                <el-option :label="$t('page.agents.level2')" :value="2" />
                <el-option :label="$t('page.agents.level3')" :value="3" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('page.agents.parentAgent')" prop="parent_id">
              <el-select v-model="formData.parent_id" :placeholder="$t('page.agents.topAgent')" clearable filterable style="width: 100%">
                <el-option v-for="agent in agentList" :key="agent.id" :label="agent.username" :value="agent.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('page.agents.commissionRate')" prop="commission_rate">
              <el-input-number v-model="formData.commission_rate" :min="0" :max="100" :precision="1" :step="0.5" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('page.agents.status')" prop="status">
              <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('page.agents.contactName')" prop="contact_name">
          <el-input v-model="formData.contact_name" :placeholder="$t('page.agents.contactNamePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('page.agents.company')" prop="company">
          <el-input v-model="formData.company" :placeholder="$t('page.agents.companyPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('page.agents.address')" prop="address">
          <el-input v-model="formData.address" :placeholder="$t('page.agents.addressPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('page.agents.notes')" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="3" :placeholder="$t('page.agents.notesPlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('page.agents.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('page.agents.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" :title="$t('page.agents.detailTitle')" width="750px" destroy-on-close>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.username')">{{ detailData.username }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.email')">{{ detailData.email || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.phone')">{{ detailData.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.agentLevel')">
          <el-tag :type="getAgentLevelType(detailData.level)" size="small">
            {{ getAgentLevelText(detailData.level) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.commissionRate')">
          <span class="commission-text">{{ detailData.commission_rate }}%</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.balance')">
          <span class="amount-text">¥{{ formatAmount(detailData.balance) }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.totalCommission')">
          <span class="amount-text">¥{{ formatAmount(detailData.total_commission) }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.subAgentCount')">{{ detailData.sub_agent_count || 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.clientCount')">{{ detailData.client_count || 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.contactName')">{{ detailData.contact_name || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.company')">{{ detailData.company || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.address')" :span="2">{{ detailData.address || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.status')">
          <el-tag :type="detailData.status === 1 ? 'success' : 'danger'" size="small">
            {{ detailData.status === 1 ? $t('page.agents.statusActive') : $t('page.agents.statusDisabled') }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.createTime')">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.agents.notes')" :span="2">{{ detailData.remark || $t('page.agents.none') }}</el-descriptions-item>
      </el-descriptions>

      <!-- 下级代理列表 -->
      <div v-if="detailData.sub_agents?.length" class="sub-agents-section">
        <el-divider content-position="left">{{ $t('page.agents.subAgents') }}</el-divider>
        <el-table :data="detailData.sub_agents" size="small" border>
          <el-table-column prop="username" :label="$t('page.agents.username')" />
          <el-table-column prop="level" :label="$t('page.agents.level')" width="80" align="center">
            <template #default="{ row }">
              {{ getAgentLevelText(row.level) }}
            </template>
          </el-table-column>
          <el-table-column prop="commission_rate" :label="$t('page.agents.commissionRate')" width="90" align="center">
            <template #default="{ row }">
              {{ row.commission_rate }}%
            </template>
          </el-table-column>
          <el-table-column prop="client_count" :label="$t('page.agents.clientCount')" width="80" align="center" />
          <el-table-column prop="status" :label="$t('page.agents.status')" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                {{ row.status === 1 ? $t('page.agents.statusActive') : $t('page.agents.statusDisabled') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" :label="$t('page.agents.createTime')" width="170" />
        </el-table>
      </div>

      <template #footer>
        <el-button @click="detailVisible = false">{{ $t('page.agents.close') }}</el-button>
      </template>
    </el-dialog>

    <!-- 调整佣金对话框 -->
    <el-dialog v-model="commissionDialogVisible" :title="$t('page.agents.adjustCommissionTitle')" width="500px" destroy-on-close>
      <el-form :model="commissionForm" :rules="commissionFormRules" ref="commissionFormRef" label-width="100px">
        <el-form-item :label="$t('page.agents.agentUser')">
          <el-input :model-value="commissionForm.username" disabled />
        </el-form-item>
        <el-form-item :label="$t('page.agents.currentRate')">
          <el-tag type="info">{{ commissionForm.current_rate }}%</el-tag>
        </el-form-item>
        <el-form-item :label="$t('page.agents.newRate')" prop="commission_rate">
          <el-input-number
            v-model="commissionForm.commission_rate"
            :min="0"
            :max="100"
            :precision="1"
            :step="0.5"
            style="width: 200px"
          />
          <span style="margin-left: 8px; color: #909399">%</span>
        </el-form-item>
        <el-form-item :label="$t('page.agents.adjustReason')" prop="reason">
          <el-input
            v-model="commissionForm.reason"
            type="textarea"
            :rows="3"
            :placeholder="$t('page.agents.adjustReasonPlaceholder')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="commissionDialogVisible = false">{{ $t('page.agents.cancel') }}</el-button>
        <el-button type="primary" @click="handleCommissionSubmit" :loading="submitLoading">{{ $t('page.agents.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  keyword: '',
  level: undefined as number | undefined,
  status: undefined as number | undefined
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref<any[]>([])

// 代理列表（用于下拉选择上级代理）
const agentList = ref<any[]>([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref($t('page.agents.addAgent'))
const formRef = ref<FormInstance>()

// 详情对话框
const detailVisible = ref(false)
const detailData = ref<any>({})

// 佣金调整对话框
const commissionDialogVisible = ref(false)
const commissionFormRef = ref<FormInstance>()
const commissionForm = reactive({
  id: 0,
  username: '',
  current_rate: 0,
  commission_rate: 0,
  reason: ''
})
const commissionFormRules: FormRules = {
  commission_rate: [
    { required: true, message: $t('page.agents.msgEnterCommissionRate'), trigger: 'blur' }
  ]
}

// 表单数据
const formData = reactive({
  id: undefined as number | undefined,
  username: '',
  email: '',
  phone: '',
  password: '',
  level: 1,
  parent_id: undefined as number | undefined,
  commission_rate: 10,
  status: 1,
  contact_name: '',
  company: '',
  address: '',
  remark: ''
})

// 表单验证规则
const formRules: FormRules = {
  username: [
    { required: true, message: $t('page.agents.msgEnterUsername'), trigger: 'blur' },
    { min: 3, max: 50, message: $t('page.agents.msgUsernameLength'), trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: $t('page.agents.msgEmailFormat'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: $t('page.agents.msgEnterPassword'), trigger: 'blur' },
    { min: 6, message: $t('page.agents.msgPasswordLength'), trigger: 'blur' }
  ],
  level: [
    { required: true, message: $t('page.agents.msgSelectLevel'), trigger: 'change' }
  ],
  commission_rate: [
    { required: true, message: $t('page.agents.msgEnterCommissionRate'), trigger: 'blur' }
  ]
}

// 代理等级映射
const getAgentLevelText = (level: number) => {
  const map: Record<number, string> = {
    1: $t('page.agents.level1'),
    2: $t('page.agents.level2'),
    3: $t('page.agents.level3')
  }
  return map[level] || $t('page.agents.unknown')
}

const getAgentLevelType = (level: number) => {
  const map: Record<number, string> = {
    1: 'danger',
    2: 'warning',
    3: 'primary'
  }
  return (map[level] || 'info') as any
}

// 格式化金额
const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

// 获取代理列表
const fetchAgents = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/agents',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        keyword: searchForm.keyword || undefined,
        level: searchForm.level,
        status: searchForm.status
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取代理列表失败:', error)
    ElMessage.error($t('page.agents.msgFetchFailed'))
  } finally {
    loading.value = false
  }
}

// 获取代理列表（用于下拉选择上级）
const fetchAgentList = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/agents',
      params: { page: 1, page_size: 9999, status: 1 }
    })
    agentList.value = data.list || []
  } catch (error) {
    console.error('获取代理列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchAgents()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.level = undefined
  searchForm.status = undefined
  handleSearch()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = $t('page.agents.addAgent')
  formData.id = undefined
  formData.username = ''
  formData.email = ''
  formData.phone = ''
  formData.password = ''
  formData.level = 1
  formData.parent_id = undefined
  formData.commission_rate = 10
  formData.status = 1
  formData.contact_name = ''
  formData.company = ''
  formData.address = ''
  formData.remark = ''
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = $t('page.agents.editAgent')
  Object.assign(formData, row)
  formData.password = ''
  dialogVisible.value = true
}

// 查看详情
const handleViewDetail = async (row: any) => {
  try {
    const data = await request.get({
      url: `/api/admin/agents/${row.id}`
    })
    detailData.value = data
  } catch (error) {
    detailData.value = { ...row }
  }
  detailVisible.value = true
}

// 调整佣金
const handleAdjustCommission = (row: any) => {
  commissionForm.id = row.id
  commissionForm.username = row.username
  commissionForm.current_rate = row.commission_rate
  commissionForm.commission_rate = row.commission_rate
  commissionForm.reason = ''
  commissionDialogVisible.value = true
}

// 提交佣金调整
const handleCommissionSubmit = async () => {
  if (!commissionFormRef.value) return

  await commissionFormRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      await request.put({
        url: `/api/admin/agents/${commissionForm.id}`,
        params: {
          commission_rate: commissionForm.commission_rate,
          reason: commissionForm.reason || undefined
        }
      })
      ElMessage.success($t('page.agents.msgCommissionAdjustSuccess'))
      commissionDialogVisible.value = false
      fetchAgents()
    } catch (error) {
      ElMessage.error($t('page.agents.msgAdjustFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

// 切换状态
const handleToggleStatus = async (row: any) => {
  try {
    await request.put({
      url: `/api/admin/agents/${row.id}`,
      params: { status: row.status === 1 ? 0 : 1 }
    })
    ElMessage.success(row.status === 1 ? $t('page.agents.msgDisabled') : $t('page.agents.msgEnabled'))
    fetchAgents()
  } catch (error) {
    ElMessage.error($t('page.agents.msgOperationFailed'))
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/agents/${formData.id}` : '/api/admin/agents'

      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }

      ElMessage.success(formData.id ? $t('page.agents.msgUpdateSuccess') : $t('page.agents.msgAddSuccess'))
      dialogVisible.value = false
      fetchAgents()
      fetchAgentList()
    } catch (error) {
      ElMessage.error($t('page.agents.msgOperationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchAgents()
}

// 页码变化
const handlePageChange = () => {
  fetchAgents()
}

onMounted(() => {
  fetchAgents()
  fetchAgentList()
})
</script>

<style scoped lang="scss">
.agent-list-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.amount-text {
  font-weight: 600;
  color: #f56c6c;
}

.commission-text {
  font-weight: 600;
  color: #409eff;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.sub-agents-section {
  margin-top: 16px;
}
</style>
