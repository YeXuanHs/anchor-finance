<template>
  <div class="invoice-list-page">
    <div class="page-header">
      <h1 class="page-title">发票管理</h1>
      <el-button type="primary" @click="$router.push('/user/invoice/apply')">
        <el-icon><Plus /></el-icon>申请发票
      </el-button>
    </div>

    <!-- 搜索筛选 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="发票状态">
          <el-select v-model="filterForm.status" placeholder="全部状态" clearable style="width: 140px;">
            <el-option label="全部" value="" />
            <el-option label="待审核" value="pending" />
            <el-option label="已开具" value="issued" />
            <el-option label="已邮寄" value="shipped" />
            <el-option label="已驳回" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item label="发票类型">
          <el-select v-model="filterForm.type" placeholder="全部类型" clearable style="width: 140px;">
            <el-option label="全部" value="" />
            <el-option label="增值税普通发票" value="normal" />
            <el-option label="增值税专用发票" value="special" />
          </el-select>
        </el-form-item>
        <el-form-item label="申请时间">
          <el-date-picker
            v-model="filterForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 260px;"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 统计卡片 -->
    <div class="summary-grid">
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#fa8c16"><Document /></el-icon>
          <div class="summary-info">
            <span class="summary-value">{{ stats.total }}</span>
            <span class="summary-label">发票总数</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#0056FF"><Timer /></el-icon>
          <div class="summary-info">
            <span class="summary-value">{{ stats.pending }}</span>
            <span class="summary-label">待审核</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#52c41a"><CircleCheck /></el-icon>
          <div class="summary-info">
            <span class="summary-value">{{ stats.issued }}</span>
            <span class="summary-label">已开具</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#722ed1"><Wallet /></el-icon>
          <div class="summary-info">
            <span class="summary-value">¥{{ stats.totalAmount }}</span>
            <span class="summary-label">开票总额</span>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table :data="filteredInvoices" stripe style="width: 100%" v-loading="loading" empty-text="暂无发票记录">
        <el-table-column prop="invoiceNo" label="发票号" width="160">
          <template #default="{ row }">
            <span class="mono-text">{{ row.invoiceNo }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="发票类型" width="160">
          <template #default="{ row }">
            <el-tag :type="row.type === 'special' ? 'warning' : 'info'" size="small" effect="light">
              {{ row.type === 'special' ? '增值税专用发票' : '增值税普通发票' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="发票抬头" min-width="180" show-overflow-tooltip />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            <span class="amount-text">¥{{ row.amount }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="applyDate" label="申请时间" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small" effect="light" round>
              {{ row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click="handleViewDetail(row)">详情</el-button>
            <el-button v-if="row.status === 'issued'" type="primary" size="small" link @click="handleDownload(row)">下载</el-button>
            <el-button v-if="row.status === 'pending'" type="danger" size="small" link @click="handleCancel(row)">撤回</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="发票详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="发票号">{{ currentInvoice?.invoiceNo }}</el-descriptions-item>
        <el-descriptions-item label="发票类型">{{ currentInvoice?.type === 'special' ? '增值税专用发票' : '增值税普通发票' }}</el-descriptions-item>
        <el-descriptions-item label="发票抬头">{{ currentInvoice?.title }}</el-descriptions-item>
        <el-descriptions-item label="纳税人识别号">{{ currentInvoice?.taxNo }}</el-descriptions-item>
        <el-descriptions-item label="开票金额">¥{{ currentInvoice?.amount }}</el-descriptions-item>
        <el-descriptions-item label="申请时间">{{ currentInvoice?.applyDate }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentInvoice?.status || '')" size="small" effect="light" round>
            {{ currentInvoice?.statusText }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="快递单号">{{ currentInvoice?.trackingNo || '-' }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ currentInvoice?.remark || '-' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Document, Timer, CircleCheck, Wallet } from '@element-plus/icons-vue'
import request from '@/utils/request'

interface Invoice {
  id: number
  invoiceNo: string
  type: 'normal' | 'special'
  title: string
  taxNo: string
  amount: string
  applyDate: string
  status: string
  statusText: string
  trackingNo?: string
  remark?: string
}

const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const detailVisible = ref(false)
const currentInvoice = ref<Invoice | null>(null)

const filterForm = reactive({
  status: '',
  type: '',
  dateRange: null as [Date, Date] | null
})

const stats = reactive({
  total: 0,
  pending: 0,
  issued: 0,
  totalAmount: '0.00'
})

const invoices = ref<Invoice[]>([])

const filteredInvoices = computed(() => {
  let result = invoices.value
  if (filterForm.status) {
    result = result.filter(i => i.status === filterForm.status)
  }
  if (filterForm.type) {
    result = result.filter(i => i.type === filterForm.type)
  }
  total.value = result.length
  return result.slice((currentPage.value - 1) * pageSize.value, currentPage.value * pageSize.value)
})

function getStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    pending: 'warning', issued: 'success', shipped: 'info', rejected: 'danger'
  }
  return map[status] || 'info'
}

function handleSearch() {
  currentPage.value = 1
  fetchInvoices()
}

function handleReset() {
  filterForm.status = ''
  filterForm.type = ''
  filterForm.dateRange = null
  currentPage.value = 1
}

function handleSizeChange(val: number) {
  pageSize.value = val
  currentPage.value = 1
  fetchInvoices()
}

function handlePageChange(val: number) {
  currentPage.value = val
  fetchInvoices()
}

function handleViewDetail(row: Invoice) {
  currentInvoice.value = row
  detailVisible.value = true
}

function handleDownload(row: Invoice) {
  ElMessage.success(`正在下载发票：${row.invoiceNo}`)
}

function handleCancel(row: Invoice) {
  ElMessageBox.confirm('确定要撤回该发票申请吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    row.status = 'cancelled'
    row.statusText = '已撤回'
    ElMessage.success('已撤回申请')
  }).catch(() => {})
}

onMounted(() => {
  fetchInvoices()
})

async function fetchInvoices() {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/invoices', {
      params: {
        page: currentPage.value,
        limit: pageSize.value,
        status: filterForm.status || undefined,
        type: filterForm.type || undefined
      }
    })
    if (data?.data) {
      invoices.value = data.data.list || data.data || []
      total.value = data.data.total || invoices.value.length
      if (data.data.stats) {
        stats.total = data.data.stats.total || 0
        stats.pending = data.data.stats.pending || 0
        stats.issued = data.data.stats.issued || 0
        stats.totalAmount = data.data.stats.totalAmount || '0.00'
      }
    }
  } catch (e) {
    console.error('Failed to fetch invoices:', e)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.invoice-list-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.filter-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.filter-card :deep(.el-card__body) {
  padding: 16px 20px 0;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
  gap: 0;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.summary-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  background: #fff;
}

.summary-card :deep(.el-card__body) { padding: 20px; }
.summary-inner { display: flex; align-items: center; gap: 16px; }
.summary-info { display: flex; flex-direction: column; gap: 4px; }
.summary-value { font-size: 24px; font-weight: 700; color: #303133; }
.summary-label { font-size: 13px; color: #909399; }

.table-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.table-card :deep(.el-card__body) { padding: 0; }
.table-card :deep(.el-table th.el-table__cell) { background: #fafbfc; font-weight: 600; }
.mono-text { font-family: 'Monaco', 'Menlo', monospace; font-size: 13px; color: #606266; }
.amount-text { font-weight: 600; color: #303133; }

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px 20px;
  border-top: 1px solid #e8ecf1;
}

@media (max-width: 768px) {
  .page-header { flex-direction: column; gap: 12px; align-items: flex-start; }
  .summary-grid { grid-template-columns: repeat(2, 1fr); }
  .filter-form :deep(.el-form-item) { margin-bottom: 12px; }
}
</style>
