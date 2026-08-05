<template>
  <div class="aff-buy-record-page">
    <div class="page-header">
      <h1 class="page-title">推广购买记录</h1>
      <div class="header-actions">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          size="default"
          style="width: 280px;"
        />
        <el-input
          v-model="searchKey"
          placeholder="搜索订单号/产品名称"
          clearable
          class="search-input"
          :prefix-icon="Search"
        />
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="stats-row">
      <el-card shadow="never" class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: linear-gradient(135deg, #0056FF, #4080FF);">
            <el-icon :size="24"><ShoppingCart /></el-icon>
          </div>
          <div class="stat-info">
            <span class="stat-value">{{ totalOrders }}</span>
            <span class="stat-label">推广订单总数</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: linear-gradient(135deg, #fa8c16, #ffc53d);">
            <el-icon :size="24"><Wallet /></el-icon>
          </div>
          <div class="stat-info">
            <span class="stat-value">¥{{ totalAmount }}</span>
            <span class="stat-label">推广消费总额</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: linear-gradient(135deg, #52c41a, #73d13d);">
            <el-icon :size="24"><Coin /></el-icon>
          </div>
          <div class="stat-info">
            <span class="stat-value">¥{{ totalCommission }}</span>
            <span class="stat-label">累计返利</span>
          </div>
        </div>
      </el-card>
    </div>

    <!-- Status Tabs -->
    <el-tabs v-model="activeTab" class="filter-tabs">
      <el-tab-pane v-for="tab in statusTabs" :key="tab.value" :name="tab.value">
        <template #label>
          <span class="tab-label">
            {{ tab.label }}
            <el-badge
              v-if="tab.count > 0"
              :value="tab.count"
              :max="99"
              class="tab-badge"
            />
          </span>
        </template>
      </el-tab-pane>
    </el-tabs>

    <!-- Data Table -->
    <el-card shadow="never" class="table-card">
      <el-table :data="paginatedRecords" stripe style="width: 100%" v-loading="loading" empty-text="暂无推广购买记录">
        <el-table-column prop="orderId" label="订单号" min-width="170" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="mono-text">{{ row.orderId }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="productName" label="产品名称" min-width="160" show-overflow-tooltip />
        <el-table-column prop="buyer" label="购买用户" min-width="120">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="24" class="record-avatar">{{ row.buyer.charAt(0) }}</el-avatar>
              <span>{{ row.buyer }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="消费金额" width="120">
          <template #default="{ row }">
            <span class="amount-text">¥{{ row.amount }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="commission" label="返利金额" width="120">
          <template #default="{ row }">
            <span class="commission-text">¥{{ row.commission }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small" effect="light" round>
              {{ row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="下单时间" width="170" sortable />
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="filteredRecords.length"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          background
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Search, ShoppingCart, Wallet, Coin } from '@element-plus/icons-vue'
import request from '@/utils/request'

const searchKey = ref('')
const activeTab = ref('all')
const currentPage = ref(1)
const pageSize = ref(10)
const dateRange = ref<[Date, Date] | null>(null)

interface BuyRecord {
  orderId: string
  productName: string
  buyer: string
  amount: string
  commission: string
  status: string
  statusText: string
  createdAt: string
}

const records = ref<BuyRecord[]>([])

const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/affiliate/records')
    if (data?.data) {
      records.value = data.data.list || data.data || []
    }
  } catch (e) {
    console.error('Failed to fetch buy records:', e)
  } finally {
    loading.value = false
  }
})

const totalOrders = computed(() => records.value.length)
const totalAmount = computed(() => {
  const sum = records.value.reduce((acc, r) => acc + parseFloat(r.amount.replace(/,/g, '')), 0)
  return sum.toLocaleString('zh-CN', { minimumFractionDigits: 2 })
})
const totalCommission = computed(() => {
  const sum = records.value.reduce((acc, r) => acc + parseFloat(r.commission.replace(/,/g, '')), 0)
  return sum.toLocaleString('zh-CN', { minimumFractionDigits: 2 })
})

const statusTabs = computed(() => [
  { label: '全部', value: 'all', count: records.value.length },
  { label: '已结算', value: 'settled', count: records.value.filter(r => r.status === 'settled').length },
  { label: '待结算', value: 'pending', count: records.value.filter(r => r.status === 'pending').length },
  { label: '已退款', value: 'cancelled', count: records.value.filter(r => r.status === 'cancelled').length }
])

const filteredRecords = computed(() => {
  let result = records.value
  if (activeTab.value !== 'all') {
    result = result.filter(r => r.status === activeTab.value)
  }
  if (searchKey.value) {
    const key = searchKey.value.toLowerCase()
    result = result.filter(r =>
      r.orderId.toLowerCase().includes(key) || r.productName.toLowerCase().includes(key)
    )
  }
  return result
})

const paginatedRecords = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredRecords.value.slice(start, start + pageSize.value)
})

function getStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    settled: 'success',
    pending: 'warning',
    cancelled: 'danger'
  }
  return map[status] || 'info'
}
</script>

<style scoped>
.aff-buy-record-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-input {
  width: 240px;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.stat-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.stat-card :deep(.el-card__body) {
  padding: 20px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
}

.stat-label {
  font-size: 13px;
  color: #909399;
}

.filter-tabs {
  background: #fff;
  border-radius: 12px;
  padding: 4px 16px 0;
  border: 1px solid #e8ecf1;
}

.filter-tabs :deep(.el-tabs__header) {
  margin: 0;
}

.tab-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.tab-badge :deep(.el-badge__content) {
  font-size: 10px;
}

.table-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.table-card :deep(.el-card__body) {
  padding: 0;
}

.table-card :deep(.el-table th.el-table__cell) {
  background: #fafbfc;
  color: #606266;
  font-weight: 600;
}

.mono-text {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 13px;
  color: #606266;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.record-avatar {
  background: linear-gradient(135deg, #0056FF, #4080FF);
  color: #fff;
  font-size: 11px;
  font-weight: 600;
}

.amount-text {
  font-weight: 600;
  color: #303133;
}

.commission-text {
  font-weight: 600;
  color: #fa8c16;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px 20px;
  border-top: 1px solid #e8ecf1;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .header-actions {
    width: 100%;
    flex-direction: column;
  }

  .search-input {
    width: 100%;
  }

  .stats-row {
    grid-template-columns: 1fr;
  }
}
</style>
