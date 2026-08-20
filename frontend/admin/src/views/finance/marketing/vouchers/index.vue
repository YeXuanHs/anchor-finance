<template>
  <div class="vouchers-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.vouchers.title') }}</span>
          <div>
            <el-button type="primary" @click="showExpressDialog()">
              <el-icon><Van /></el-icon>
              {{ $t('finance.vouchers.expressManagement') }}
            </el-button>
            <el-button type="primary" @click="showRateDialog">
              <el-icon><Setting /></el-icon>
              {{ $t('finance.vouchers.rateConfig') }}
            </el-button>
          </div>
        </div>
      </template>

      <!-- 工具栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('finance.vouchers.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('finance.vouchers.allStatus')" clearable>
            <el-option :label="$t('finance.vouchers.statusPending')" value="Pending" />
            <el-option :label="$t('finance.vouchers.statusCancelled')" value="Cancelled" />
            <el-option :label="$t('finance.vouchers.statusReject')" value="Reject" />
            <el-option :label="$t('finance.vouchers.statusUnpaid')" value="Unpaid" />
            <el-option :label="$t('finance.vouchers.statusSend')" value="Send" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('finance.vouchers.search')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('finance.vouchers.searchPlaceholder')" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('finance.vouchers.searchBtn') }}</el-button>
          <el-button @click="handleReset">{{ $t('finance.vouchers.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border stripe>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="username" :label="$t('finance.vouchers.username')" width="120" show-overflow-tooltip />
        <el-table-column prop="create_time" :label="$t('finance.vouchers.applyTime')" width="180" />
        <el-table-column prop="title" :label="$t('finance.vouchers.invoiceTitle')" min-width="150" show-overflow-tooltip />
        <el-table-column :label="$t('finance.vouchers.invoiceType')" width="160">
          <template #default="{ row }">
            <div>{{ getTypeLabel(row.issue_type) }}</div>
            <div class="text-secondary">{{ getVoucherTypeLabel(row.voucher_type) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="amount" :label="$t('finance.vouchers.amount')" width="100" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('finance.vouchers.status')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.vouchers.shippingAddress')" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            {{ formatAddress(row) }}
          </template>
        </el-table-column>
        <el-table-column prop="express_name" :label="$t('finance.vouchers.expressMethod')" width="120" />
        <el-table-column prop="notes" :label="$t('finance.vouchers.notes')" width="150" show-overflow-tooltip />
        <el-table-column prop="check_time" :label="$t('finance.vouchers.auditTime')" width="180" />
        <el-table-column :label="$t('finance.vouchers.actions')" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="showDetail(row)">{{ $t('finance.vouchers.detail') }}</el-button>
            <el-button v-if="row.status === 'Pending' || row.status === 'Unpaid'" type="warning" link @click="showAuditDialog(row)">{{ $t('finance.vouchers.audit') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <!-- 发票详情对话框 -->
    <el-dialog v-model="detailDialogVisible" :title="$t('finance.vouchers.invoiceDetail')" width="800px" destroy-on-close>
      <div v-loading="detailLoading" class="detail-container">
        <template v-if="detailData">
          <el-card shadow="never" class="detail-section">
            <template #header><span class="section-title">{{ $t('finance.vouchers.basicInfo') }}</span></template>
            <el-descriptions :column="2" border>
              <el-descriptions-item :label="$t('finance.vouchers.applyId')">{{ detailData.voucher?.id }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.userId')">{{ detailData.uid }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.applyTime')">{{ detailData.voucher?.create_time }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.status')">
                <el-tag :type="getStatusType(detailData.voucher?.status)">{{ getStatusLabel(detailData.voucher?.status) }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.auditTime')">{{ detailData.voucher?.check_time || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.notes')">{{ detailData.voucher?.notes || '-' }}</el-descriptions-item>
            </el-descriptions>
          </el-card>
          <el-card shadow="never" class="detail-section">
            <template #header><span class="section-title">{{ $t('finance.vouchers.invoiceTitle') }}</span></template>
            <el-descriptions :column="2" border>
              <el-descriptions-item :label="$t('finance.vouchers.titleName')">{{ detailData.voucher?.title }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.type')">{{ getTypeLabel(detailData.voucher?.issue_type) }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.invoiceType')">{{ getVoucherTypeLabel(detailData.voucher?.voucher_type) }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.taxId')">{{ detailData.voucher?.tax_id || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.bank')">{{ detailData.voucher?.bank || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.account')">{{ detailData.voucher?.account || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.address')">{{ detailData.voucher?.address || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.phone')">{{ detailData.voucher?.phone || '-' }}</el-descriptions-item>
            </el-descriptions>
          </el-card>
          <el-card shadow="never" class="detail-section">
            <template #header><span class="section-title">{{ $t('finance.vouchers.shippingAddress') }}</span></template>
            <el-descriptions :column="2" border>
              <el-descriptions-item :label="$t('finance.vouchers.province')">{{ detailData.voucher?.province || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.city')">{{ detailData.voucher?.city || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.region')">{{ detailData.voucher?.region || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.detailAddress')" :span="2">{{ detailData.voucher?.detail || '-' }}</el-descriptions-item>
            </el-descriptions>
          </el-card>
          <el-card shadow="never" class="detail-section">
            <template #header><span class="section-title">{{ $t('finance.vouchers.expressInfo') }}</span></template>
            <el-descriptions :column="2" border>
              <el-descriptions-item :label="$t('finance.vouchers.expressName')">{{ detailData.voucher?.express_name || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="$t('finance.vouchers.expressFee')">¥{{ formatAmount(detailData.voucher?.express_price) }}</el-descriptions-item>
            </el-descriptions>
          </el-card>
          <el-card shadow="never" class="detail-section">
            <template #header><span class="section-title">{{ $t('finance.vouchers.relatedBills') }}</span></template>
            <el-table :data="detailData.invoices" border size="small">
              <el-table-column prop="id" :label="$t('finance.vouchers.billId')" width="80" />
              <el-table-column prop="subtotal" :label="$t('finance.vouchers.billAmount')" width="120" align="right">
                <template #default="{ row }">¥{{ formatAmount(row.subtotal) }}</template>
              </el-table-column>
              <el-table-column prop="taxed" :label="$t('finance.vouchers.taxRate')" width="100" />
              <el-table-column prop="taxed_amount" :label="$t('finance.vouchers.taxAmount')" width="120" align="right">
                <template #default="{ row }">¥{{ formatAmount(row.taxed_amount) }}</template>
              </el-table-column>
              <el-table-column :label="$t('finance.vouchers.billDetail')" min-width="200">
                <template #default="{ row }">
                  <div v-for="item in row.items" :key="item.id" class="invoice-item">{{ item.description }}</div>
                </template>
              </el-table-column>
            </el-table>
            <div class="total-amount">
              <span>{{ $t('finance.vouchers.invoiceTotal') }}：</span>
              <span class="amount-text">¥{{ formatAmount(detailData.voucher_amount) }}</span>
            </div>
          </el-card>
        </template>
      </div>
    </el-dialog>

    <!-- 审核对话框 -->
    <el-dialog v-model="auditDialogVisible" :title="$t('finance.vouchers.auditInvoice')" width="500px" destroy-on-close>
      <el-form :model="auditForm" :rules="auditRules" ref="auditFormRef" label-width="80px">
        <el-form-item :label="$t('finance.vouchers.invoiceId')"><el-input :model-value="auditForm.id" disabled /></el-form-item>
        <el-form-item :label="$t('finance.vouchers.username')"><el-input :model-value="auditForm.username" disabled /></el-form-item>
        <el-form-item :label="$t('finance.vouchers.auditStatus')" prop="status">
          <el-select v-model="auditForm.status" :placeholder="$t('finance.vouchers.selectAuditStatus')" style="width: 100%">
            <el-option :label="$t('finance.vouchers.optionReject')" value="Reject" />
            <el-option :label="$t('finance.vouchers.optionSend')" value="Send" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('finance.vouchers.notes')" prop="notes">
          <el-input v-model="auditForm.notes" type="textarea" :rows="3" :placeholder="$t('finance.vouchers.notesPlaceholder')" maxlength="500" show-word-limit />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="auditDialogVisible = false">{{ $t('finance.vouchers.cancel') }}</el-button>
        <el-button type="primary" @click="handleAuditSubmit" :loading="auditLoading">{{ $t('finance.vouchers.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 费率配置对话框 -->
    <el-dialog v-model="rateDialogVisible" :title="$t('finance.vouchers.rateConfig')" width="400px" destroy-on-close>
      <el-form :model="rateForm" :rules="rateRules" ref="rateFormRef" label-width="100px">
        <el-form-item :label="$t('finance.vouchers.enableInvoiceMgmt')" prop="voucher_manager">
          <el-switch v-model="rateForm.voucher_manager" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item :label="$t('finance.vouchers.taxRatePercent')" prop="rate">
          <el-input-number v-model="rateForm.rate" :min="0" :max="100" :precision="2" :step="0.1" controls-position="right" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rateDialogVisible = false">{{ $t('finance.vouchers.cancel') }}</el-button>
        <el-button type="primary" @click="handleRateSubmit" :loading="rateLoading">{{ $t('finance.vouchers.save') }}</el-button>
      </template>
    </el-dialog>

    <!-- 快递管理对话框 -->
    <el-dialog v-model="expressDialogVisible" :title="$t('finance.vouchers.expressManagement')" width="700px" destroy-on-close>
      <div class="express-header">
        <el-button type="primary" size="small" @click="showExpressForm()">{{ $t('finance.vouchers.addExpress') }}</el-button>
      </div>
      <el-table :data="expressList" v-loading="expressLoading" border size="small">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" :label="$t('finance.vouchers.expressName')" min-width="150" />
        <el-table-column prop="price" :label="$t('finance.vouchers.fee')" width="120" align="right">
          <template #default="{ row }">¥{{ formatAmount(row.price) }}</template>
        </el-table-column>
        <el-table-column prop="create_time" :label="$t('finance.vouchers.createTime')" width="180" />
        <el-table-column :label="$t('finance.vouchers.actions')" width="150" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="showExpressForm(row)">{{ $t('finance.vouchers.edit') }}</el-button>
            <el-popconfirm :title="$t('finance.vouchers.confirmDeleteExpress')" @confirm="handleDeleteExpress(row)">
              <template #reference><el-button type="danger" link>{{ $t('finance.vouchers.delete') }}</el-button></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <!-- 快递表单对话框 -->
    <el-dialog v-model="expressFormVisible" :title="expressForm.id ? $t('finance.vouchers.editExpress') : $t('finance.vouchers.addExpress')" width="400px" destroy-on-close>
      <el-form :model="expressForm" :rules="expressRules" ref="expressFormRef" label-width="80px">
        <el-form-item :label="$t('finance.vouchers.name')" prop="name">
          <el-input v-model="expressForm.name" :placeholder="$t('finance.vouchers.expressNamePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('finance.vouchers.fee')" prop="price">
          <el-input-number v-model="expressForm.price" :min="0" :precision="2" :step="1" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="expressFormVisible = false">{{ $t('finance.vouchers.cancel') }}</el-button>
        <el-button type="primary" @click="handleExpressSubmit" :loading="expressSubmitLoading">{{ $t('finance.vouchers.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Setting, Van } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const statusMap: Record<string, string> = { Pending: $t('finance.vouchers.statusPending'), Cancelled: $t('finance.vouchers.statusCancelled'), Reject: $t('finance.vouchers.statusReject'), Unpaid: $t('finance.vouchers.statusUnpaid'), Send: $t('finance.vouchers.statusSend') }
const statusTypeMap: Record<string, string> = { Pending: 'warning', Cancelled: 'info', Reject: 'danger', Unpaid: 'info', Send: 'success' }
const typeMap: Record<string, string> = { person: $t('finance.vouchers.personal'), company: $t('finance.vouchers.company') }
const voucherTypeMap: Record<string, string> = { common: $t('finance.vouchers.vatNormal'), dedicated: $t('finance.vouchers.vatDedicated') }

const loading = ref(false)
const detailLoading = ref(false)
const auditLoading = ref(false)
const rateLoading = ref(false)
const expressLoading = ref(false)
const expressSubmitLoading = ref(false)

const searchForm = reactive({ status: '' as string, keyword: '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const detailDialogVisible = ref(false)
const detailData = ref<any>(null)

const auditDialogVisible = ref(false)
const auditFormRef = ref<FormInstance>()
const auditForm = reactive({ id: 0, username: '', status: '' as string, notes: '' })

const rateDialogVisible = ref(false)
const rateFormRef = ref<FormInstance>()
const rateForm = reactive({ voucher_manager: 0, rate: 0 })

const expressDialogVisible = ref(false)
const expressFormVisible = ref(false)
const expressFormRef = ref<FormInstance>()
const expressList = ref<any[]>([])
const expressForm = reactive({ id: null as number | null, name: '', price: 0 })

const auditRules: FormRules = { status: [{ required: true, message: $t('finance.vouchers.msgSelectAuditStatus'), trigger: 'change' }] }
const rateRules: FormRules = { rate: [{ required: true, message: $t('finance.vouchers.msgEnterTaxRate'), trigger: 'blur' }] }
const expressRules: FormRules = { name: [{ required: true, message: $t('finance.vouchers.msgEnterExpressName'), trigger: 'blur' }] }

const formatAmount = (amount: number | undefined) => amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
const formatAddress = (row: any) => { const parts = [row.province, row.city, row.region].filter(Boolean); return parts.join('') || '-' }
const getStatusLabel = (status: string) => statusMap[status] || status
const getStatusType = (status: string) => (statusTypeMap[status] || 'info') as any
const getTypeLabel = (type: string) => typeMap[type] || type
const getVoucherTypeLabel = (type: string) => voucherTypeMap[type] || type

const fetchVouchers = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size, status: searchForm.status || undefined, keyword: searchForm.keyword || undefined }
    const data = await request.get({ url: '/api/admin/voucher-list', params })
    tableData.value = data.voucher || data.list || []
    pagination.total = data.total || 0
  } catch { ElMessage.error($t('finance.vouchers.msgFetchListFailed')) } finally { loading.value = false }
}

const fetchDetail = async (id: number) => {
  detailLoading.value = true
  try { detailData.value = await request.get({ url: `/api/admin/voucher-detail/${id}` }) }
  catch { ElMessage.error($t('finance.vouchers.msgFetchDetailFailed')) } finally { detailLoading.value = false }
}

const fetchRateConfig = async () => {
  try { const data = await request.get({ url: '/api/admin/voucher-rate' }); rateForm.voucher_manager = data.voucher_manager || 0; rateForm.rate = data.rate || 0 }
  catch { ElMessage.error($t('finance.vouchers.msgFetchRateFailed')) }
}

const fetchExpressList = async () => {
  expressLoading.value = true
  try { const data = await request.get({ url: '/api/admin/voucher/express' }); expressList.value = data.express || data || [] }
  catch { ElMessage.error($t('finance.vouchers.msgFetchExpressFailed')) } finally { expressLoading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchVouchers() }
const handleReset = () => { searchForm.status = ''; searchForm.keyword = ''; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchVouchers() }
const handlePageChange = () => { fetchVouchers() }

const showDetail = async (row: any) => { detailData.value = null; detailDialogVisible.value = true; await fetchDetail(row.id) }
const showAuditDialog = (row: any) => { auditForm.id = row.id; auditForm.username = row.username; auditForm.status = ''; auditForm.notes = ''; auditDialogVisible.value = true }

const handleAuditSubmit = async () => {
  if (!auditFormRef.value) return
  await auditFormRef.value.validate(async (valid) => {
    if (!valid) return
    auditLoading.value = true
    try { await request.post({ url: '/api/admin/voucher-status', data: { id: auditForm.id, status: auditForm.status, notes: auditForm.notes } }); ElMessage.success($t('finance.vouchers.msgAuditSuccess')); auditDialogVisible.value = false; fetchVouchers() }
    catch { ElMessage.error($t('finance.vouchers.msgAuditFailed')) } finally { auditLoading.value = false }
  })
}

const showRateDialog = async () => { await fetchRateConfig(); rateDialogVisible.value = true }

const handleRateSubmit = async () => {
  if (!rateFormRef.value) return
  await rateFormRef.value.validate(async (valid) => {
    if (!valid) return
    rateLoading.value = true
    try { await request.post({ url: '/api/admin/voucher-rate', data: { voucher_manager: rateForm.voucher_manager, rate: rateForm.rate } }); ElMessage.success($t('finance.vouchers.msgRateSaveSuccess')); rateDialogVisible.value = false }
    catch { ElMessage.error($t('finance.vouchers.msgRateSaveFailed')) } finally { rateLoading.value = false }
  })
}

const showExpressDialog = () => { fetchExpressList(); expressDialogVisible.value = true }

const showExpressForm = (row?: any) => {
  if (row) { expressForm.id = row.id; expressForm.name = row.name; expressForm.price = row.price }
  else { expressForm.id = null; expressForm.name = ''; expressForm.price = 0 }
  expressFormVisible.value = true
}

const handleExpressSubmit = async () => {
  if (!expressFormRef.value) return
  await expressFormRef.value.validate(async (valid) => {
    if (!valid) return
    expressSubmitLoading.value = true
    try {
      if (expressForm.id) {
        await request.post({ url: '/api/admin/voucher/express', data: { id: expressForm.id, name: expressForm.name, price: expressForm.price } })
      } else {
        await request.post({ url: '/api/admin/voucher/express', data: { name: expressForm.name, price: expressForm.price } })
      }
      ElMessage.success($t('finance.vouchers.msgOperationSuccess')); expressFormVisible.value = false; fetchExpressList()
    } catch { ElMessage.error($t('finance.vouchers.msgOperationFailed')) } finally { expressSubmitLoading.value = false }
  })
}

const handleDeleteExpress = async (row: any) => {
  try { await request.del({ url: '/api/admin/voucher/express', data: { id: row.id } }); ElMessage.success($t('finance.vouchers.msgDeleteSuccess')); fetchExpressList() }
  catch { ElMessage.error($t('finance.vouchers.msgDeleteFailed')) }
}

onMounted(() => { fetchVouchers() })
</script>

<style scoped lang="scss">
.vouchers-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.amount-text { font-weight: 600; color: #f56c6c; }
.text-secondary { font-size: 12px; color: #909399; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.detail-container { min-height: 200px; }
.detail-section { margin-bottom: 20px; &:last-child { margin-bottom: 0; } }
.section-title { font-weight: 600; }
.invoice-item { padding: 2px 0; font-size: 12px; color: #606266; }
.total-amount { margin-top: 16px; text-align: right; font-size: 14px; font-weight: 600; }
.express-header { margin-bottom: 16px; display: flex; justify-content: flex-end; }
</style>
