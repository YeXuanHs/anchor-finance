<template>
  <div class="credit-used-detail-page">
    <div class="page-header">
      <h1 class="page-title">已用额度明细</h1>
      <div class="header-right">
        <el-select v-model="statusFilter" placeholder="状态筛选" clearable style="width: 140px">
          <el-option label="全部" value="" />
          <el-option label="使用中" value="active" />
          <el-option label="已还款" value="repaid" />
          <el-option label="已逾期" value="overdue" />
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
      </div>
    </div>

    <!-- 概览卡片 -->
    <div class="summary-grid">
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#fa8c16"><Wallet /></el-icon>
          <div class="summary-info">
            <span class="summary-value">¥{{ summary.total_used?.toFixed(2) || '0.00' }}</span>
            <span class="summary-label">已用总额</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#52c41a"><CircleCheck /></el-icon>
          <div class="summary-info">
            <span class="summary-value">¥{{ summary.total_repaid?.toFixed(2) || '0.00' }}</span>
            <span class="summary-label">已还总额</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#ff4d4f"><Warning /></el-icon>
          <div class="summary-info">
            <span class="summary-value">¥{{ summary.total_overdue?.toFixed(2) || '0.00' }}</span>
            <span class="summary-label">逾期金额</span>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 明细表格 -->
    <el-card shadow="never" class="table-card">
      <el-table :data="tableData" stripe v-loading="loading" empty-text="暂无额度使用记录">
        <el-table-column prop="id" label="记录编号" width="160">
          <template #default="{ row }">
            <span class="mono">{{ row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="使用时间" width="170" />
        <el-table-column prop="amount" label="使用金额" width="130">
          <template #default="{ row }">
            <span class="text-danger">-¥{{ row.amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="bill_id" label="关联账单" width="160">
          <template #default="{ row }">
            <el-button v-if="row.bill_id" type="primary" link @click="goBill(row.bill_id)">
              {{ row.bill_id }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="order_no" label="关联订单" width="160">
          <template #default="{ row }">
            <el-button v-if="row.order_no" type="primary" link @click="goOrder(row.order_no)">
              {{ row.order_no }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="用途说明" min-width="180" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" size="small">
              {{ statusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="repaid_at" label="还款时间" width="170">
          <template #default="{ row }">
            {{ row.repaid_at || '-' }}
          </template>
        </el-table-column>
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
import { Search, Wallet, CircleCheck, Warning } from '@element-plus/icons-vue'
import request from '@/utils/http'

const router = useRouter()
const loading = ref(false)
const statusFilter = ref('')
const dateRange = ref<string[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const summary = ref({
  total_used: 0,
  total_repaid: 0,
  total_overdue: 0
})

const tableData = ref<any[]>([])

const statusTagType = (status: string) => {
  const map: Record<string, string> = { active: 'warning', repaid: 'success', overdue: 'danger' }
  return map[status] || 'info'
}

const statusText = (status: string) => {
  const map: Record<string, string> = { active: '使用中', repaid: '已还款', overdue: '已逾期' }
  return map[status] || status
}

function goBill(billId: string) {
  router.push(`/user/credit-bill/${billId}`)
}

function goOrder(orderNo: string) {
  router.push(`/user/orders/${orderNo}`)
}

async function fetchList() {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (statusFilter.value) params.status = statusFilter.value
    if (dateRange.value?.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    const res = await request.get('/v1/credit/used-detail', { params })
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
    const res = await request.get('/v1/credit/used-summary')
    summary.value = res.data?.data || summary.value
  } catch {}
}

onMounted(() => {
  fetchList()
  fetchSummary()
})
</script>

<style scoped lang="scss">
.credit-used-detail-page {
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

  .text-danger {
    color: #f56c6c;
    font-weight: 600;
  }

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
