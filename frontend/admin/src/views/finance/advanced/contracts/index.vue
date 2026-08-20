<template>
  <div class="contract-list-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.contracts.pageTitle') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('finance.contracts.addContract') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('finance.contracts.contractNo')">
          <el-input v-model="searchForm.contract_no" :placeholder="$t('finance.contracts.enterContractNo')" clearable />
        </el-form-item>
        <el-form-item :label="$t('finance.contracts.clientName')">
          <el-input v-model="searchForm.client_name" :placeholder="$t('finance.contracts.enterClientName')" clearable />
        </el-form-item>
        <el-form-item :label="$t('finance.contracts.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('finance.contracts.all')" clearable>
            <el-option :label="$t('finance.contracts.statusDraft')" :value="0" />
            <el-option :label="$t('finance.contracts.statusActive')" :value="1" />
            <el-option :label="$t('finance.contracts.statusExpired')" :value="2" />
            <el-option :label="$t('finance.contracts.statusTerminated')" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('finance.contracts.dateRange')">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            :range-separator="$t('finance.contracts.to')"
            :start-placeholder="$t('finance.contracts.startDate')"
            :end-placeholder="$t('finance.contracts.endDate')"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('finance.contracts.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('finance.contracts.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="contract_no" :label="$t('finance.contracts.contractNo')" width="160">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">
              {{ row.contract_no }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="title" :label="$t('finance.contracts.contractTitle')" min-width="180" show-overflow-tooltip />
        <el-table-column prop="client_name" :label="$t('finance.contracts.relatedClient')" width="120" />
        <el-table-column prop="amount" :label="$t('finance.contracts.contractAmount')" width="120" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('finance.contracts.status')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getContractStatusType(row.status)" size="small">
              {{ getContractStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="start_date" :label="$t('finance.contracts.startDate')" width="120" />
        <el-table-column prop="end_date" :label="$t('finance.contracts.endDate')" width="120" />
        <el-table-column prop="created_at" :label="$t('finance.contracts.createdAt')" width="170" />
        <el-table-column :label="$t('finance.contracts.actions')" width="250" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">{{ $t('finance.contracts.detail') }}</el-button>
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('finance.contracts.edit') }}</el-button>
            <el-button type="success" link @click="handleGeneratePDF(row)">{{ $t('finance.contracts.generatePDF') }}</el-button>
            <el-popconfirm :title="$t('finance.contracts.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('finance.contracts.delete') }}</el-button>
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('finance.contracts.contractNo')" prop="contract_no">
              <el-input v-model="formData.contract_no" :placeholder="$t('finance.contracts.enterContractNo')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('finance.contracts.contractTitle')" prop="title">
              <el-input v-model="formData.title" :placeholder="$t('finance.contracts.enterContractTitle')" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('finance.contracts.relatedClient')" prop="client_id">
              <el-select v-model="formData.client_id" :placeholder="$t('finance.contracts.selectClient')" filterable style="width: 100%">
                <el-option v-for="client in clientList" :key="client.id" :label="client.username" :value="client.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('finance.contracts.contractAmount')" prop="amount">
              <el-input-number v-model="formData.amount" :min="0" :precision="2" :step="100" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('finance.contracts.startDate')" prop="start_date">
              <el-date-picker v-model="formData.start_date" type="date" value-format="YYYY-MM-DD" :placeholder="$t('finance.contracts.selectStartDate')" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('finance.contracts.endDate')" prop="end_date">
              <el-date-picker v-model="formData.end_date" type="date" value-format="YYYY-MM-DD" :placeholder="$t('finance.contracts.selectEndDate')" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('finance.contracts.status')" prop="status">
          <el-select v-model="formData.status" :placeholder="$t('finance.contracts.selectStatus')">
            <el-option :label="$t('finance.contracts.statusDraft')" :value="0" />
            <el-option :label="$t('finance.contracts.statusActive')" :value="1" />
            <el-option :label="$t('finance.contracts.statusExpired')" :value="2" />
            <el-option :label="$t('finance.contracts.statusTerminated')" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('finance.contracts.contractContent')" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="6" :placeholder="$t('finance.contracts.enterContent')" />
        </el-form-item>
        <el-form-item :label="$t('finance.contracts.remark')" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="3" :placeholder="$t('finance.contracts.enterRemark')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('finance.contracts.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('finance.contracts.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" :title="$t('finance.contracts.contractDetail')" width="750px" destroy-on-close>
      <el-descriptions :column="2" border>
        <el-descriptions-item :label="$t('finance.contracts.contractNo')">{{ detailData.contract_no }}</el-descriptions-item>
        <el-descriptions-item :label="$t('finance.contracts.status')">
          <el-tag :type="getContractStatusType(detailData.status)" size="small">
            {{ getContractStatusText(detailData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('finance.contracts.contractTitle')" :span="2">{{ detailData.title }}</el-descriptions-item>
        <el-descriptions-item :label="$t('finance.contracts.relatedClient')">{{ detailData.client_name || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('finance.contracts.contractAmount')">
          <span class="amount-text">¥{{ formatAmount(detailData.amount) }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('finance.contracts.startDate')">{{ detailData.start_date || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('finance.contracts.endDate')">{{ detailData.end_date || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('finance.contracts.createdAt')" :span="2">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('finance.contracts.contractContent')" :span="2">
          <div class="contract-content">{{ detailData.content || $t('finance.contracts.none') }}</div>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('finance.contracts.remark')" :span="2">{{ detailData.remark || $t('finance.contracts.none') }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">{{ $t('finance.contracts.close') }}</el-button>
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
  contract_no: '',
  client_name: '',
  status: undefined as number | undefined,
  date_range: [] as string[]
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref<any[]>([])

// 客户列表（用于下拉选择）
const clientList = ref<any[]>([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref($t('finance.contracts.addContract'))
const formRef = ref<FormInstance>()

// 详情对话框
const detailVisible = ref(false)
const detailData = ref<any>({})

// 表单数据
const formData = reactive({
  id: undefined as number | undefined,
  contract_no: '',
  title: '',
  client_id: undefined as number | undefined,
  amount: 0,
  start_date: '',
  end_date: '',
  status: 0,
  content: '',
  remark: ''
})

// 表单验证规则
const formRules: FormRules = {
  contract_no: [
    { required: true, message: $t('finance.contracts.enterContractNo'), trigger: 'blur' }
  ],
  title: [
    { required: true, message: $t('finance.contracts.enterContractTitle'), trigger: 'blur' }
  ],
  client_id: [
    { required: true, message: $t('finance.contracts.selectClient'), trigger: 'change' }
  ],
  start_date: [
    { required: true, message: $t('finance.contracts.selectStartDate'), trigger: 'change' }
  ],
  end_date: [
    { required: true, message: $t('finance.contracts.selectEndDate'), trigger: 'change' }
  ]
}

// 合同状态映射
const CONTRACT_STATUS_MAP: Record<number, { text: string; type: string }> = {
  0: { text: $t('finance.contracts.statusDraft'), type: 'info' },
  1: { text: $t('finance.contracts.statusActive'), type: 'success' },
  2: { text: $t('finance.contracts.statusExpired'), type: 'warning' },
  3: { text: $t('finance.contracts.statusTerminated'), type: 'danger' }
}

// 获取合同状态文本
const getContractStatusText = (status: number) => {
  return CONTRACT_STATUS_MAP[status]?.text || $t('finance.contracts.unknown')
}

// 获取合同状态类型
const getContractStatusType = (status: number) => {
  return (CONTRACT_STATUS_MAP[status]?.type || 'info') as any
}

// 格式化金额
const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

// 获取合同列表
const fetchContracts = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      contract_no: searchForm.contract_no || undefined,
      client_name: searchForm.client_name || undefined,
      status: searchForm.status
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({
      url: '/api/admin/contracts',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取合同列表失败:', error)
    ElMessage.error($t('finance.contracts.fetchListFailed'))
  } finally {
    loading.value = false
  }
}

// 获取客户列表（用于下拉选择）
const fetchClients = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/users',
      params: { page: 1, page_size: 9999 }
    })
    clientList.value = data.list || []
  } catch (error) {
    console.error('获取客户列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchContracts()
}

// 重置
const handleReset = () => {
  searchForm.contract_no = ''
  searchForm.client_name = ''
  searchForm.status = undefined
  searchForm.date_range = []
  handleSearch()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = $t('finance.contracts.addContract')
  formData.id = undefined
  formData.contract_no = ''
  formData.title = ''
  formData.client_id = undefined
  formData.amount = 0
  formData.start_date = ''
  formData.end_date = ''
  formData.status = 0
  formData.content = ''
  formData.remark = ''
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = $t('finance.contracts.editContract')
  Object.assign(formData, row)
  dialogVisible.value = true
}

// 查看详情
const handleViewDetail = (row: any) => {
  detailData.value = { ...row }
  detailVisible.value = true
}

// 生成PDF
const handleGeneratePDF = async (row: any) => {
  try {
    const data = await request.get({
      url: `/api/admin/contracts/${row.id}/pdf`
    })
    if (data.url) {
      window.open(data.url, '_blank')
      ElMessage.success($t('finance.contracts.pdfSuccess'))
    }
  } catch (error) {
    ElMessage.error($t('finance.contracts.pdfFailed'))
  }
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/contracts/${row.id}`
    })
    ElMessage.success($t('finance.contracts.deleteSuccess'))
    fetchContracts()
  } catch (error) {
    ElMessage.error($t('finance.contracts.deleteFailed'))
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/contracts/${formData.id}` : '/api/admin/contracts'

      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }

      ElMessage.success(formData.id ? $t('finance.contracts.updateSuccess') : $t('finance.contracts.addSuccess'))
      dialogVisible.value = false
      fetchContracts()
    } catch (error) {
      ElMessage.error($t('finance.contracts.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchContracts()
}

// 页码变化
const handlePageChange = () => {
  fetchContracts()
}

onMounted(() => {
  fetchContracts()
  fetchClients()
})
</script>

<style scoped lang="scss">
.contract-list-page {
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

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.contract-content {
  max-height: 200px;
  overflow-y: auto;
  white-space: pre-wrap;
  line-height: 1.6;
}
</style>
