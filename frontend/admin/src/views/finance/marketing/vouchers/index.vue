<template>
  <div class="vouchers-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>发票管理</span>
          <div>
            <el-button type="primary" @click="showExpressDialog()">
              <el-icon><Van /></el-icon>
              快递管理
            </el-button>
            <el-button type="primary" @click="showRateDialog">
              <el-icon><Setting /></el-icon>
              费率配置
            </el-button>
          </div>
        </div>
      </template>

      <!-- 工具栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部状态" clearable>
            <el-option label="待审核" value="Pending" />
            <el-option label="已取消" value="Cancelled" />
            <el-option label="已驳回" value="Reject" />
            <el-option label="待支付" value="Unpaid" />
            <el-option label="已发出" value="Send" />
          </el-select>
        </el-form-item>
        <el-form-item label="搜索">
          <el-input v-model="searchForm.keyword" placeholder="用户名/备注" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border stripe>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="username" label="用户名" width="120" show-overflow-tooltip />
        <el-table-column prop="create_time" label="申请时间" width="180" />
        <el-table-column prop="title" label="发票抬头" min-width="150" show-overflow-tooltip />
        <el-table-column label="发票类型" width="160">
          <template #default="{ row }">
            <div>{{ getTypeLabel(row.issue_type) }}</div>
            <div class="text-secondary">{{ getVoucherTypeLabel(row.voucher_type) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="100" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="收件地址" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            {{ formatAddress(row) }}
          </template>
        </el-table-column>
        <el-table-column prop="express_name" label="快递方式" width="120" />
        <el-table-column prop="notes" label="备注" width="150" show-overflow-tooltip />
        <el-table-column prop="check_time" label="审核时间" width="180" />
        <el-table-column label="操作" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="showDetail(row)">详情</el-button>
            <el-button v-if="row.status === 'Pending' || row.status === 'Unpaid'" type="warning" link @click="showAuditDialog(row)">审核</el-button>
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
    <el-dialog v-model="detailDialogVisible" title="发票详情" width="800px" destroy-on-close>
      <div v-loading="detailLoading" class="detail-container">
        <template v-if="detailData">
          <el-card shadow="never" class="detail-section">
            <template #header><span class="section-title">基本信息</span></template>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="申请ID">{{ detailData.voucher?.id }}</el-descriptions-item>
              <el-descriptions-item label="用户ID">{{ detailData.uid }}</el-descriptions-item>
              <el-descriptions-item label="申请时间">{{ detailData.voucher?.create_time }}</el-descriptions-item>
              <el-descriptions-item label="状态">
                <el-tag :type="getStatusType(detailData.voucher?.status)">{{ getStatusLabel(detailData.voucher?.status) }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="审核时间">{{ detailData.voucher?.check_time || '-' }}</el-descriptions-item>
              <el-descriptions-item label="备注">{{ detailData.voucher?.notes || '-' }}</el-descriptions-item>
            </el-descriptions>
          </el-card>
          <el-card shadow="never" class="detail-section">
            <template #header><span class="section-title">发票抬头</span></template>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="抬头名称">{{ detailData.voucher?.title }}</el-descriptions-item>
              <el-descriptions-item label="类型">{{ getTypeLabel(detailData.voucher?.issue_type) }}</el-descriptions-item>
              <el-descriptions-item label="发票类型">{{ getVoucherTypeLabel(detailData.voucher?.voucher_type) }}</el-descriptions-item>
              <el-descriptions-item label="税号">{{ detailData.voucher?.tax_id || '-' }}</el-descriptions-item>
              <el-descriptions-item label="开户行">{{ detailData.voucher?.bank || '-' }}</el-descriptions-item>
              <el-descriptions-item label="账号">{{ detailData.voucher?.account || '-' }}</el-descriptions-item>
              <el-descriptions-item label="地址">{{ detailData.voucher?.address || '-' }}</el-descriptions-item>
              <el-descriptions-item label="电话">{{ detailData.voucher?.phone || '-' }}</el-descriptions-item>
            </el-descriptions>
          </el-card>
          <el-card shadow="never" class="detail-section">
            <template #header><span class="section-title">收件地址</span></template>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="省份">{{ detailData.voucher?.province || '-' }}</el-descriptions-item>
              <el-descriptions-item label="城市">{{ detailData.voucher?.city || '-' }}</el-descriptions-item>
              <el-descriptions-item label="区县">{{ detailData.voucher?.region || '-' }}</el-descriptions-item>
              <el-descriptions-item label="详细地址" :span="2">{{ detailData.voucher?.detail || '-' }}</el-descriptions-item>
            </el-descriptions>
          </el-card>
          <el-card shadow="never" class="detail-section">
            <template #header><span class="section-title">快递信息</span></template>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="快递名称">{{ detailData.voucher?.express_name || '-' }}</el-descriptions-item>
              <el-descriptions-item label="快递费用">¥{{ formatAmount(detailData.voucher?.express_price) }}</el-descriptions-item>
            </el-descriptions>
          </el-card>
          <el-card shadow="never" class="detail-section">
            <template #header><span class="section-title">关联账单</span></template>
            <el-table :data="detailData.invoices" border size="small">
              <el-table-column prop="id" label="账单ID" width="80" />
              <el-table-column prop="subtotal" label="账单金额" width="120" align="right">
                <template #default="{ row }">¥{{ formatAmount(row.subtotal) }}</template>
              </el-table-column>
              <el-table-column prop="taxed" label="税率" width="100" />
              <el-table-column prop="taxed_amount" label="税额" width="120" align="right">
                <template #default="{ row }">¥{{ formatAmount(row.taxed_amount) }}</template>
              </el-table-column>
              <el-table-column label="账单明细" min-width="200">
                <template #default="{ row }">
                  <div v-for="item in row.items" :key="item.id" class="invoice-item">{{ item.description }}</div>
                </template>
              </el-table-column>
            </el-table>
            <div class="total-amount">
              <span>发票金额合计：</span>
              <span class="amount-text">¥{{ formatAmount(detailData.voucher_amount) }}</span>
            </div>
          </el-card>
        </template>
      </div>
    </el-dialog>

    <!-- 审核对话框 -->
    <el-dialog v-model="auditDialogVisible" title="审核发票" width="500px" destroy-on-close>
      <el-form :model="auditForm" :rules="auditRules" ref="auditFormRef" label-width="80px">
        <el-form-item label="发票ID"><el-input :model-value="auditForm.id" disabled /></el-form-item>
        <el-form-item label="用户名"><el-input :model-value="auditForm.username" disabled /></el-form-item>
        <el-form-item label="审核状态" prop="status">
          <el-select v-model="auditForm.status" placeholder="请选择审核状态" style="width: 100%">
            <el-option label="驳回 (Reject)" value="Reject" />
            <el-option label="已发出 (Send)" value="Send" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注" prop="notes">
          <el-input v-model="auditForm.notes" type="textarea" :rows="3" placeholder="请输入备注（可选，最多500字符）" maxlength="500" show-word-limit />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="auditDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAuditSubmit" :loading="auditLoading">确认</el-button>
      </template>
    </el-dialog>

    <!-- 费率配置对话框 -->
    <el-dialog v-model="rateDialogVisible" title="费率配置" width="400px" destroy-on-close>
      <el-form :model="rateForm" :rules="rateRules" ref="rateFormRef" label-width="100px">
        <el-form-item label="启用发票管理" prop="voucher_manager">
          <el-switch v-model="rateForm.voucher_manager" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="税率(%)" prop="rate">
          <el-input-number v-model="rateForm.rate" :min="0" :max="100" :precision="2" :step="0.1" controls-position="right" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rateDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleRateSubmit" :loading="rateLoading">保存</el-button>
      </template>
    </el-dialog>

    <!-- 快递管理对话框 -->
    <el-dialog v-model="expressDialogVisible" title="快递管理" width="700px" destroy-on-close>
      <div class="express-header">
        <el-button type="primary" size="small" @click="showExpressForm()">添加快递</el-button>
      </div>
      <el-table :data="expressList" v-loading="expressLoading" border size="small">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="快递名称" min-width="150" />
        <el-table-column prop="price" label="费用" width="120" align="right">
          <template #default="{ row }">¥{{ formatAmount(row.price) }}</template>
        </el-table-column>
        <el-table-column prop="create_time" label="创建时间" width="180" />
        <el-table-column label="操作" width="150" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="showExpressForm(row)">编辑</el-button>
            <el-popconfirm title="确定删除该快递方式吗？" @confirm="handleDeleteExpress(row)">
              <template #reference><el-button type="danger" link>删除</el-button></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <!-- 快递表单对话框 -->
    <el-dialog v-model="expressFormVisible" :title="expressForm.id ? '编辑快递' : '添加快递'" width="400px" destroy-on-close>
      <el-form :model="expressForm" :rules="expressRules" ref="expressFormRef" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="expressForm.name" placeholder="请输入快递名称" />
        </el-form-item>
        <el-form-item label="费用" prop="price">
          <el-input-number v-model="expressForm.price" :min="0" :precision="2" :step="1" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="expressFormVisible = false">取消</el-button>
        <el-button type="primary" @click="handleExpressSubmit" :loading="expressSubmitLoading">确定</el-button>
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

const statusMap: Record<string, string> = { Pending: '待审核', Cancelled: '已取消', Reject: '已驳回', Unpaid: '待支付', Send: '已发出' }
const statusTypeMap: Record<string, string> = { Pending: 'warning', Cancelled: 'info', Reject: 'danger', Unpaid: 'info', Send: 'success' }
const typeMap: Record<string, string> = { person: '个人', company: '企业' }
const voucherTypeMap: Record<string, string> = { common: '增值税普通发票', dedicated: '增值税专用发票' }

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

const auditRules: FormRules = { status: [{ required: true, message: '请选择审核状态', trigger: 'change' }] }
const rateRules: FormRules = { rate: [{ required: true, message: '请输入税率', trigger: 'blur' }] }
const expressRules: FormRules = { name: [{ required: true, message: '请输入快递名称', trigger: 'blur' }] }

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
  } catch { ElMessage.error('获取发票列表失败') } finally { loading.value = false }
}

const fetchDetail = async (id: number) => {
  detailLoading.value = true
  try { detailData.value = await request.get({ url: `/api/admin/voucher-detail/${id}` }) }
  catch { ElMessage.error('获取发票详情失败') } finally { detailLoading.value = false }
}

const fetchRateConfig = async () => {
  try { const data = await request.get({ url: '/api/admin/voucher-rate' }); rateForm.voucher_manager = data.voucher_manager || 0; rateForm.rate = data.rate || 0 }
  catch { ElMessage.error('获取费率配置失败') }
}

const fetchExpressList = async () => {
  expressLoading.value = true
  try { const data = await request.get({ url: '/api/admin/voucher/express' }); expressList.value = data.express || data || [] }
  catch { ElMessage.error('获取快递列表失败') } finally { expressLoading.value = false }
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
    try { await request.post({ url: '/api/admin/voucher-status', data: { id: auditForm.id, status: auditForm.status, notes: auditForm.notes } }); ElMessage.success('审核成功'); auditDialogVisible.value = false; fetchVouchers() }
    catch { ElMessage.error('审核失败') } finally { auditLoading.value = false }
  })
}

const showRateDialog = async () => { await fetchRateConfig(); rateDialogVisible.value = true }

const handleRateSubmit = async () => {
  if (!rateFormRef.value) return
  await rateFormRef.value.validate(async (valid) => {
    if (!valid) return
    rateLoading.value = true
    try { await request.post({ url: '/api/admin/voucher-rate', data: { voucher_manager: rateForm.voucher_manager, rate: rateForm.rate } }); ElMessage.success('费率配置保存成功'); rateDialogVisible.value = false }
    catch { ElMessage.error('费率配置保存失败') } finally { rateLoading.value = false }
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
      ElMessage.success('操作成功'); expressFormVisible.value = false; fetchExpressList()
    } catch { ElMessage.error('操作失败') } finally { expressSubmitLoading.value = false }
  })
}

const handleDeleteExpress = async (row: any) => {
  try { await request.del({ url: '/api/admin/voucher/express', data: { id: row.id } }); ElMessage.success('删除成功'); fetchExpressList() }
  catch { ElMessage.error('删除失败') }
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
