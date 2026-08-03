<template>
  <div class="marketplace-transactions-page">
    <div class="page-header">
      <h1 class="page-title">交易记录</h1>
      <div class="header-right">
        <el-select v-model="typeFilter" placeholder="交易类型" clearable style="width: 130px">
          <el-option label="全部" value="" />
          <el-option label="收入" value="income" />
          <el-option label="支出" value="expense" />
          <el-option label="退款" value="refund" />
        </el-select>
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          style="width: 280px"
        />
        <el-button type="primary" @click="fetchList">
          <el-icon><Search /></el-icon> 搜索
        </el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>
    </div>

    <!-- 概览 -->
    <div class="summary-grid">
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#52c41a"><Top /></el-icon>
          <div class="summary-info">
            <span class="summary-value">¥{{ summary.total_income?.toFixed(2) || '0.00' }}</span>
            <span class="summary-label">总收入</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#ff4d4f"><Bottom /></el-icon>
          <div class="summary-info">
            <span class="summary-value">¥{{ summary.total_expense?.toFixed(2) || '0.00' }}</span>
            <span class="summary-label">总支出</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#1890ff"><Document /></el-icon>
          <div class="summary-info">
            <span class="summary-value">{{ summary.total_count || 0 }}</span>
            <span class="summary-label">交易总数</span>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 交易表格 -->
    <el-card shadow="never" class="table-card">
      <el-table :data="tableData" stripe v-loading="loading" empty-text="暂无交易记录">
        <el-table-column prop="id" label="交易编号" width="180">
          <template #default="{ row }">
            <span class="mono">{{ row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="交易时间" width="170" />
        <el-table-column prop="type" label="交易类型" width="100">
          <template #default="{ row }">
            <el-tag :type="typeTagType(row.type)" size="small">
              {{ typeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="130">
          <template #default="{ row }">
            <span :class="row.type === 'income' ? 'text-success' : 'text-danger'">
              {{ row.type === 'income' ? '+' : '-' }}¥{{ row.amount?.toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="buyer" label="买家/卖家" min-width="140" show-overflow-tooltip />
        <el-table-column prop="product_name" label="商品名称" min-width="180" show-overflow-tooltip />
        <el-table-column prop="order_no" label="关联订单" width="160">
          <template #default="{ row }">
            <el-button v-if="row.order_no" type="primary" link @click="goOrder(row.order_no)">
              {{ row.order_no }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" size="small">
              {{ statusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="140" show-overflow-tooltip />
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, Top, Bottom, Document } from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()
const loading = ref(false)
const typeFilter = ref('')
const dateRange = ref<string[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const summary = ref({
  total_income: 0,
  total_expense: 0,
  total_count: 0
})

const tableData = ref<any[]>([])

const typeTagType = (type: string) => {
  const map: Record<string, string> = { income: 'success', expense: 'danger', refund: 'warning' }
  return map[type] || 'info'
}

const typeText = (type: string) => {
  const map: Record<string, string> = { income: '收入', expense: '支出', refund: '退款' }
  return map[type] || type
}

const statusTagType = (status: string) => {
  const map: Record<string, string> = { completed: 'success', pending: 'warning', failed: 'danger', cancelled: 'info' }
  return map[status] || 'info'
}

const statusText = (status: string) => {
  const map: Record<string, string> = { completed: '已完成', pending: '处理中', failed: '失败', cancelled: '已取消' }
  return map[status] || status
}

function goOrder(orderNo: string) {
  router.push(`/user/orders/${orderNo}`)
}

function handleReset() {
  typeFilter.value = ''
  dateRange.value = []
  currentPage.value = 1
  fetchList()
}

async function fetchList() {
  loading.value = true
  try {
    const params: any = { page: currentPage.value, page_size: pageSize.value }
    if (typeFilter.value) params.type = typeFilter.value
    if (dateRange.value?.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    const res = await request.get('/api/v1/marketplace/transactions', { params })
    tableData.value = res.data?.data?.list || res.data?.list || []
    total.value = res.data?.data?.total || 0
  } catch {
    tableData.value = []
  } finally {
    loading.value = false
  }
}

async function fetchSummary() {
  try {
    const res = await request.get('/api/v1/marketplace/transactions/summary')
    summary.value = res.data?.data || summary.value
  } catch {}
}

onMounted(() => {
  fetchList()
  fetchSummary()
})
</script>

<style scoped lang="scss">
.marketplace-transactions-page {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-wrap: wrap;
    gap: 12px;

    .page-title {
      font-size: 20px;
      font-weight: 700;
      color: #303133;
      margin: 0;
    }

    .header-right {
      display: flex;
      gap: 12px;
      align-items: center;
    }
  }

  .summary-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
    margin-bottom: 20px;
  }

  .summary-card {
    border-radius: 12px;
    border: 1px solid #e8ecf1;
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

  .mono {
    font-family: 'Monaco', 'Menlo', monospace;
    font-size: 13px;
    color: #606266;
  }

  .text-success { color: #67c23a; font-weight: 600; }
  .text-danger { color: #f56c6c; font-weight: 600; }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    padding: 16px 20px;
    border-top: 1px solid #e8ecf1;
  }

  @media (max-width: 768px) {
    .page-header { flex-direction: column; align-items: flex-start; }
    .header-right { width: 100%; flex-wrap: wrap; }
    .summary-grid { grid-template-columns: 1fr; }
  }
}
</style>
