<template>
  <div class="orders-page">
    <div class="page-header">
      <h1 class="page-title">订单管理</h1>
      <el-input
        v-model="searchKey"
        placeholder="搜索订单号或产品名称"
        clearable
        class="search-input"
        :prefix-icon="Search"
      />
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

    <!-- Orders Table -->
    <el-card shadow="never" class="table-card">
      <el-table :data="paginatedOrders" stripe style="width: 100%" empty-text="暂无订单">
        <el-table-column prop="id" label="订单号" min-width="170" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="mono-text">{{ row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="product" label="产品名称" min-width="160" show-overflow-tooltip />
        <el-table-column prop="spec" label="配置规格" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="spec-text">{{ row.spec }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="cycle" label="计费周期" width="110" />
        <el-table-column prop="amount" label="金额" width="110">
          <template #default="{ row }">
            <span class="amount-text">¥{{ row.amount }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small" effect="light" round>
              {{ row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="下单时间" width="170" />
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click="handleDetail(row)">详情</el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="primary"
              size="small"
              @click="handlePay(row)"
            >去支付</el-button>
            <el-button
              v-if="row.status === 'active'"
              type="primary"
              size="small"
              plain
              @click="handleRenew(row)"
            >续费</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="filteredOrders.length"
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
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import request from '@/utils/request'

const searchKey = ref('')
const activeTab = ref('all')
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)

interface Order {
  id: string
  product: string
  spec: string
  cycle: string
  amount: string
  status: string
  statusText: string
  createdAt: string
}

const orders = ref<Order[]>([])

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v2/orders')
    orders.value = data.data?.list || data.list || data.data || []
  } catch (e) { console.error(e) } finally { loading.value = false }
})

const statusTabs = computed(() => [
  { label: '全部', value: 'all', count: orders.value.length },
  { label: '待支付', value: 'pending', count: orders.value.filter(o => o.status === 'pending').length },
  { label: '已开通', value: 'active', count: orders.value.filter(o => o.status === 'active').length },
  { label: '已完成', value: 'completed', count: orders.value.filter(o => o.status === 'completed').length },
  { label: '已取消', value: 'cancelled', count: orders.value.filter(o => o.status === 'cancelled').length }
])

const filteredOrders = computed(() => {
  let result = orders.value
  if (activeTab.value !== 'all') {
    result = result.filter(o => o.status === activeTab.value)
  }
  if (searchKey.value) {
    const key = searchKey.value.toLowerCase()
    result = result.filter(o =>
      o.id.toLowerCase().includes(key) || o.product.toLowerCase().includes(key)
    )
  }
  return result
})

const paginatedOrders = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredOrders.value.slice(start, start + pageSize.value)
})

function getStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    active: 'success',
    pending: 'warning',
    completed: 'info',
    cancelled: 'danger'
  }
  return map[status] || 'info'
}

function handlePay(order: Order) {
  ElMessage.info(`正在跳转支付页面：${order.id}`)
}

function handleRenew(order: Order) {
  ElMessage.info(`正在处理续费：${order.id}`)
}

function handleDetail(order: Order) {
  ElMessage.info(`查看订单详情：${order.id}`)
}
</script>

<style scoped>
.orders-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
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

.search-input {
  width: 280px;
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

.table-card :deep(.el-table) {
  --el-table-border-color: #e8ecf1;
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

.spec-text {
  font-size: 12px;
  color: #909399;
  font-family: 'Monaco', 'Menlo', monospace;
}

.amount-text {
  font-weight: 600;
  color: #303133;
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
    gap: 12px;
    align-items: flex-start;
  }

  .search-input {
    width: 100%;
  }
}
</style>
