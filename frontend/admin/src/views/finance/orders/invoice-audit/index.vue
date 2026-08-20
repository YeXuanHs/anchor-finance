<template>
  <div class="invoice-audit-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('invoiceAudit.title') }}</span>
          <div>
            <el-button type="warning" @click="handleBatchAudit" :disabled="selectedIds.length === 0">
              {{ $t('invoiceAudit.batchAudit') }} ({{ selectedIds.length }})
            </el-button>
            <el-button type="success" @click="handleExport">
              <el-icon><Download /></el-icon>
              {{ $t('invoiceAudit.export') }}
            </el-button>
          </div>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('invoiceAudit.invoiceNo')">
          <el-input v-model="searchForm.invoice_no" :placeholder="$t('invoiceAudit.inputInvoiceNo')" clearable />
        </el-form-item>
        <el-form-item :label="$t('invoiceAudit.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('invoiceAudit.all')" clearable>
            <el-option :label="$t('invoiceAudit.pending')" :value="0" />
            <el-option :label="$t('invoiceAudit.approved')" :value="1" />
            <el-option :label="$t('invoiceAudit.rejected')" :value="2" />
            <el-option :label="$t('invoiceAudit.issued')" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('invoiceAudit.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('invoiceAudit.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" :selectable="(row: any) => row.status === 0" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="invoice_no" :label="$t('invoiceAudit.invoiceNo')" width="180" />
        <el-table-column prop="client_name" :label="$t('invoiceAudit.client')" width="120" />
        <el-table-column prop="title" :label="$t('invoiceAudit.invoiceTitle')" min-width="200" />
        <el-table-column prop="tax_no" :label="$t('invoiceAudit.taxNo')" width="160" />
        <el-table-column prop="amount" :label="$t('invoiceAudit.amount')" width="120" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('invoiceAudit.status')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="statusTypeMap[row.status]" size="small">
              {{ statusLabelMap[row.status]() }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('invoiceAudit.createdAt')" width="180" />
        <el-table-column :label="$t('invoiceAudit.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">{{ $t('invoiceAudit.detail') }}</el-button>
            <el-button v-if="row.status === 0" type="success" link @click="handleAudit(row, 1)">{{ $t('invoiceAudit.pass') }}</el-button>
            <el-button v-if="row.status === 0" type="danger" link @click="handleAudit(row, 2)">{{ $t('invoiceAudit.reject') }}</el-button>
            <el-button v-if="row.status === 1" type="warning" link @click="handleMarkIssued(row)">{{ $t('invoiceAudit.markIssued') }}</el-button>
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
    <el-dialog v-model="detailVisible" :title="$t('invoiceAudit.invoiceDetail')" width="650px">
      <el-descriptions :column="2" border v-if="detailData">
        <el-descriptions-item :label="$t('invoiceAudit.invoiceNo')">{{ detailData.invoice_no }}</el-descriptions-item>
        <el-descriptions-item :label="$t('invoiceAudit.client')">{{ detailData.client_name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('invoiceAudit.invoiceTitle')" :span="2">{{ detailData.title }}</el-descriptions-item>
        <el-descriptions-item :label="$t('invoiceAudit.taxNo')">{{ detailData.tax_no || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('invoiceAudit.amount')">
          <span class="amount-text">¥{{ formatAmount(detailData.amount) }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('invoiceAudit.bankName')">{{ detailData.bank_name || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('invoiceAudit.bankAccount')">{{ detailData.bank_account || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('invoiceAudit.address')" :span="2">{{ detailData.address || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('invoiceAudit.phone')">{{ detailData.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('invoiceAudit.email')">{{ detailData.email || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('invoiceAudit.createdAt')" :span="2">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('invoiceAudit.remark')" :span="2">{{ detailData.remark || $t('invoiceAudit.none') }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">{{ $t('invoiceAudit.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download } from '@element-plus/icons-vue'
import { $t } from '@/locales'
import request from '@/utils/http'

const statusTypeMap: Record<number, any> = {
  0: 'warning',
  1: 'success',
  2: 'danger',
  3: 'info'
}

const statusLabelMap: Record<number, () => string> = {
  0: () => $t('invoiceAudit.pending'),
  1: () => $t('invoiceAudit.approved'),
  2: () => $t('invoiceAudit.rejected'),
  3: () => $t('invoiceAudit.issued')
}

const loading = ref(false)
const detailVisible = ref(false)
const detailData = ref<any>(null)
const selectedIds = ref<number[]>([])

const searchForm = reactive({
  invoice_no: '',
  status: undefined as number | undefined
})

const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/invoices/audit',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error($t('invoiceAudit.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.invoice_no = ''; searchForm.status = undefined; handleSearch() }
const handleView = (row: any) => { detailData.value = row; detailVisible.value = true }

const handleAudit = async (row: any, status: number) => {
  const action = status === 1 ? $t('invoiceAudit.passAction') : $t('invoiceAudit.rejectAction')
  try {
    await ElMessageBox.confirm($t('invoiceAudit.confirmAudit', { action }), $t('invoiceAudit.auditConfirm'))
    await request.post({ url: `/api/admin/invoices/${row.id}/audit`, params: { status } })
    ElMessage.success($t('invoiceAudit.auditSuccess', { action }))
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') ElMessage.error($t('invoiceAudit.operationFailed'))
  }
}

const handleMarkIssued = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/invoices/${row.id}/issued` })
    ElMessage.success($t('invoiceAudit.markedIssued'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('invoiceAudit.operationFailed'))
  }
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleSelectionChange = (rows: any[]) => {
  selectedIds.value = rows.map((r: any) => r.id)
}

const handleBatchAudit = async () => {
  if (selectedIds.value.length === 0) return
  try {
    await ElMessageBox.confirm($t('invoiceAudit.batchConfirmMessage', { count: selectedIds.value.length }), $t('invoiceAudit.batchConfirm'))
    await request.post({ url: '/api/admin/invoices/batch-audit', params: { ids: selectedIds.value, status: 1 } })
    ElMessage.success($t('invoiceAudit.batchSuccess'))
    selectedIds.value = []
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') ElMessage.error($t('invoiceAudit.batchFailed'))
  }
}

const handleExport = async () => {
  try {
    const params: any = { ...searchForm }
    const res = await request.get({ url: '/api/admin/invoices/export', params, responseType: 'blob' as any })
    const blob = new Blob([res], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${$t('invoiceAudit.exportFileName')}_${new Date().toISOString().slice(0, 10)}.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success($t('invoiceAudit.exportSuccess'))
  } catch (error) {
    ElMessage.error($t('invoiceAudit.exportFailed'))
  }
}

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.invoice-audit-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.amount-text { font-weight: 600; color: var(--el-color-primary); }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
