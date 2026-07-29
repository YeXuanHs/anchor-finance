<template>
  <div class="invoices-page">
    <div class="page-header">
      <h1 class="page-title">账单管理</h1>
      <div class="header-right">
        <el-select v-model="statusFilter" placeholder="账单状态" clearable style="width: 140px;">
          <el-option label="全部" value="all" />
          <el-option label="待支付" value="unpaid" />
          <el-option label="已支付" value="paid" />
          <el-option label="已逾期" value="overdue" />
        </el-select>
      </div>
    </div>

    <div class="summary-grid">
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#fa8c16"><Wallet /></el-icon>
          <div class="summary-info">
            <span class="summary-value">¥{{ summaryPending }}</span>
            <span class="summary-label">待支付金额</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#52c41a"><CircleCheck /></el-icon>
          <div class="summary-info">
            <span class="summary-value">¥{{ summaryPaid }}</span>
            <span class="summary-label">已支付金额</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#0056FF"><Document /></el-icon>
          <div class="summary-info">
            <span class="summary-value">{{ totalInvoices }}</span>
            <span class="summary-label">账单总数</span>
          </div>
        </div>
      </el-card>
    </div>

    <el-card shadow="never" class="table-card">
      <el-table :data="filteredInvoices" stripe style="width: 100%" empty-text="暂无账单">
        <el-table-column prop="id" label="账单号" width="140">
          <template #default="{ row }">
            <span class="mono-text">{{ row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="账单描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="product" label="关联产品" min-width="140" show-overflow-tooltip />
        <el-table-column prop="amount" label="金额" width="110">
          <template #default="{ row }">
            <span class="amount-text">¥{{ row.amount }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="dueDate" label="到期日" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getInvoiceStatusType(row.status)" size="small" effect="light" round>
              {{ row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status !== 'paid'" type="primary" size="small" @click="handlePay(row)">支付</el-button>
            <el-button type="primary" size="small" link>详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          :total="totalInvoices"
          :page-size="10"
          layout="total, prev, pager, next"
          background
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Wallet, CircleCheck, Document } from '@element-plus/icons-vue'
import request from '@/utils/request'

const statusFilter = ref('all')
const currentPage = ref(1)
const loading = ref(false)
const summaryPending = ref('0.00')
const summaryPaid = ref('0.00')
const totalInvoices = ref(0)

interface Invoice {
  id: string
  description: string
  product: string
  amount: string
  dueDate: string
  status: string
  statusText: string
}

const invoices = ref<Invoice[]>([])

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/invoices', { params: { page: currentPage.value } })
    invoices.value = data.data?.list || data.list || []
    summaryPending.value = data.data?.summaryPending || '0.00'
    summaryPaid.value = data.data?.summaryPaid || '0.00'
    totalInvoices.value = data.data?.total || invoices.value.length
  } catch (e) { console.error(e) } finally { loading.value = false }
})

const filteredInvoices = computed(() => {
  if (statusFilter.value === 'all') return invoices.value
  return invoices.value.filter(i => i.status === statusFilter.value)
})

function getInvoiceStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    paid: 'success', unpaid: 'warning', overdue: 'danger'
  }
  return map[status] || 'info'
}

function handlePay(invoice: Invoice) {
  ElMessage.info(`正在跳转支付：${invoice.id}`)
}
</script>

<style scoped>
.invoices-page {
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

.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
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
  .summary-grid { grid-template-columns: 1fr; }
}
</style>
