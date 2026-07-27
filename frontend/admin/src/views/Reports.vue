<template>
  <div class="reports">
    <!-- Date Range Picker -->
    <n-card :bordered="false" rounded class="filter-card">
      <n-space align="center" justify="space-between">
        <n-space align="center">
          <span class="filter-label">统计周期</span>
          <n-tabs v-model:value="activeTab" type="segment" animated>
            <n-tab-pane name="today" tab="今日" />
            <n-tab-pane name="week" tab="本周" />
            <n-tab-pane name="month" tab="本月" />
            <n-tab-pane name="year" tab="本年" />
          </n-tabs>
        </n-space>
        <n-date-picker v-model:value="dateRange" type="daterange" clearable />
      </n-space>
    </n-card>

    <!-- Income Stats -->
    <div class="section-title">收入统计</div>
    <n-grid :cols="4" :x-gap="16" :y-gap="16" responsive="screen" :item-responsive="true">
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card income" :bordered="false" rounded>
          <n-statistic label="今日收入" :value="incomeStats.todayRevenue">
            <template #prefix>¥</template>
            <template #suffix>
              <span class="stat-trend up">+12.5%</span>
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card income" :bordered="false" rounded>
          <n-statistic label="本月收入" :value="incomeStats.monthRevenue">
            <template #prefix>¥</template>
            <template #suffix>
              <span class="stat-trend up">+8.3%</span>
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card income" :bordered="false" rounded>
          <n-statistic label="总收入" :value="incomeStats.totalRevenue">
            <template #prefix>¥</template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card warning" :bordered="false" rounded>
          <n-statistic label="待收金额" :value="incomeStats.pendingAmount">
            <template #prefix>¥</template>
            <template #suffix>
              <span class="stat-trend warning">需跟进</span>
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- User Stats -->
    <div class="section-title">用户统计</div>
    <n-grid :cols="3" :x-gap="16" :y-gap="16" responsive="screen" :item-responsive="true">
      <n-gi span="3 m:1">
        <n-card class="stat-card user" :bordered="false" rounded>
          <n-statistic label="总用户数" :value="userStats.totalUsers">
            <template #suffix>
              <span class="stat-trend up">+156 本月</span>
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi span="3 m:1">
        <n-card class="stat-card user" :bordered="false" rounded>
          <n-statistic label="新增用户" :value="userStats.newUsers">
            <template #suffix>
              <span class="stat-trend up">+15.2%</span>
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi span="3 m:1">
        <n-card class="stat-card user" :bordered="false" rounded>
          <n-statistic label="活跃用户" :value="userStats.activeUsers">
            <template #suffix>
              <span class="stat-trend neutral">近7天</span>
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- Order Stats -->
    <div class="section-title">订单统计</div>
    <n-grid :cols="4" :x-gap="16" :y-gap="16" responsive="screen" :item-responsive="true">
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card order" :bordered="false" rounded>
          <n-statistic label="总订单" :value="orderStats.totalOrders" />
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card order-pending" :bordered="false" rounded>
          <n-statistic label="待付款" :value="orderStats.pendingPayment">
            <template #suffix>
              <span class="stat-trend warning">需催付</span>
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card order-success" :bordered="false" rounded>
          <n-statistic label="已完成" :value="orderStats.completed" />
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card order-cancel" :bordered="false" rounded>
          <n-statistic label="已取消" :value="orderStats.cancelled" />
        </n-card>
      </n-gi>
    </n-grid>

    <!-- Charts & Rankings -->
    <n-grid :cols="3" :x-gap="16" :y-gap="16" style="margin-top: 24px" responsive="screen" :item-responsive="true">
      <!-- Income Trend Chart -->
      <n-gi span="3 l:2">
        <n-card title="收入趋势" :bordered="false" rounded class="chart-card">
          <div class="bar-chart">
            <div class="chart-y-axis">
              <span>¥50k</span>
              <span>¥40k</span>
              <span>¥30k</span>
              <span>¥20k</span>
              <span>¥10k</span>
              <span>0</span>
            </div>
            <div class="chart-bars">
              <div v-for="(item, index) in chartData" :key="index" class="bar-item">
                <div class="bar-tooltip">¥{{ item.value.toLocaleString() }}</div>
                <div class="bar" :style="{ height: (item.value / maxValue) * 100 + '%' }">
                  <div class="bar-fill" />
                </div>
                <span class="bar-label">{{ item.label }}</span>
              </div>
            </div>
          </div>
        </n-card>
      </n-gi>

      <!-- Product Sales Ranking -->
      <n-gi span="3 l:1">
        <n-card title="产品销量排行" :bordered="false" rounded class="ranking-card">
          <div class="ranking-list">
            <div v-for="(item, index) in productRanking" :key="index" class="ranking-item">
              <span class="rank" :class="'rank-' + (index + 1)">{{ index + 1 }}</span>
              <span class="product-name">{{ item.name }}</span>
              <span class="sales-count">{{ item.sales }} 单</span>
              <n-progress
                type="line"
                :percentage="(item.sales / productRanking[0].sales) * 100"
                :show-indicator="false"
                :height="6"
                :border-radius="3"
                :fill-border-radius="3"
              />
            </div>
          </div>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- Recent Transactions -->
    <n-card title="最近交易记录" :bordered="false" rounded class="table-card" style="margin-top: 24px">
      <template #header-extra>
        <n-button text type="primary">查看全部</n-button>
      </template>
      <n-data-table
        :columns="transactionColumns"
        :data="recentTransactions"
        :bordered="false"
        :pagination="{ pageSize: 8 }"
        size="small"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { h, ref, computed } from 'vue'
import { NTag, NProgress, type DataTableColumns } from 'naive-ui'

// Active tab
const activeTab = ref('month')
const dateRange = ref<[number, number] | null>(null)

// Income Stats
const incomeStats = reactive({
  todayRevenue: 12680,
  monthRevenue: 289560,
  totalRevenue: 1586420,
  pendingAmount: 45230,
})

// User Stats
const userStats = reactive({
  totalUsers: 12856,
  newUsers: 342,
  activeUsers: 2156,
})

// Order Stats
const orderStats = reactive({
  totalOrders: 8654,
  pendingPayment: 156,
  completed: 7892,
  cancelled: 606,
})

// Chart Data (last 12 months)
const chartData = ref([
  { label: '1月', value: 32000 },
  { label: '2月', value: 28000 },
  { label: '3月', value: 35000 },
  { label: '4月', value: 42000 },
  { label: '5月', value: 38000 },
  { label: '6月', value: 45000 },
  { label: '7月', value: 48000 },
  { label: '8月', value: 41000 },
  { label: '9月', value: 50000 },
  { label: '10月', value: 46000 },
  { label: '11月', value: 52000 },
  { label: '12月', value: 55000 },
])

const maxValue = computed(() => Math.max(...chartData.value.map((d) => d.value)))

// Product Sales Ranking
const productRanking = ref([
  { name: '基础版主机', sales: 1256 },
  { name: '高级版主机', sales: 986 },
  { name: '4核8G云服务器', sales: 756 },
  { name: '企业版主机', sales: 623 },
  { name: '1核2G云服务器', sales: 589 },
  { name: 'CDN加速服务', sales: 456 },
  { name: 'SSL证书', sales: 321 },
])

// Transaction status
const statusMap: Record<string, { label: string; type: string }> = {
  success: { label: '成功', type: 'success' },
  pending: { label: '处理中', type: 'warning' },
  failed: { label: '失败', type: 'error' },
  refunded: { label: '已退款', type: 'info' },
}

// Recent Transactions
const recentTransactions = ref([
  { id: 'TX20260727001', user: '张三', type: '购买产品', product: '基础版主机', amount: 299, status: 'success', time: '2026-07-27 10:15' },
  { id: 'TX20260727002', user: '李四', type: '续费', product: '高级版主机', amount: 599, status: 'success', time: '2026-07-27 09:42' },
  { id: 'TX20260727003', user: '王五', type: '购买产品', product: '企业版主机', amount: 1299, status: 'pending', time: '2026-07-27 09:18' },
  { id: 'TX20260726004', user: '赵六', type: '退款', product: '1核2G云服务器', amount: -89, status: 'refunded', time: '2026-07-26 16:35' },
  { id: 'TX20260726005', user: '孙七', type: '购买产品', product: 'CDN加速服务', amount: 199, status: 'success', time: '2026-07-26 14:22' },
  { id: 'TX20260726006', user: '周八', type: '续费', product: '4核8G云服务器', amount: 399, status: 'failed', time: '2026-07-26 11:08' },
  { id: 'TX20260725007', user: '吴九', type: '购买产品', product: 'SSL证书', amount: 129, status: 'success', time: '2026-07-25 18:45' },
  { id: 'TX20260725008', user: '郑十', type: '升级', product: '基础版→高级版', amount: 300, status: 'success', time: '2026-07-25 15:30' },
])

// Table Columns
const transactionColumns: DataTableColumns<any> = [
  { title: '交易号', key: 'id', width: 150, ellipsis: { tooltip: true } },
  { title: '用户', key: 'user', width: 80 },
  { title: '类型', key: 'type', width: 90 },
  { title: '产品', key: 'product', ellipsis: { tooltip: true } },
  {
    title: '金额',
    key: 'amount',
    width: 100,
    render: (row) =>
      h(
        'span',
        { style: `font-weight:600;color:${row.amount >= 0 ? '#18a058' : '#d03050'}` },
        `${row.amount >= 0 ? '+' : ''}¥${Math.abs(row.amount).toLocaleString()}`
      ),
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: (row) => {
      const s = statusMap[row.status]
      return h(NTag, { type: s.type as any, size: 'small', round: true, bordered: false }, { default: () => s.label })
    },
  },
  { title: '时间', key: 'time', width: 150 },
]
</script>

<style scoped>
.reports {
  padding: 0;
}

.filter-card {
  margin-bottom: 24px;
  border-radius: 12px;
}

.filter-label {
  font-size: 14px;
  color: #8c8c8c;
  margin-right: 8px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #ffffff;
  margin: 24px 0 12px;
  padding-left: 12px;
  border-left: 3px solid #18a058;
}

.stat-card {
  border-radius: 12px;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
}

.stat-card.income {
  background: linear-gradient(135deg, rgba(24, 160, 88, 0.15) 0%, rgba(24, 160, 88, 0.05) 100%);
  border: 1px solid rgba(24, 160, 88, 0.2);
}

.stat-card.warning {
  background: linear-gradient(135deg, rgba(240, 160, 32, 0.15) 0%, rgba(240, 160, 32, 0.05) 100%);
  border: 1px solid rgba(240, 160, 32, 0.2);
}

.stat-card.user {
  background: linear-gradient(135deg, rgba(32, 128, 240, 0.15) 0%, rgba(32, 128, 240, 0.05) 100%);
  border: 1px solid rgba(32, 128, 240, 0.2);
}

.stat-card.order {
  background: linear-gradient(135deg, rgba(128, 128, 240, 0.15) 0%, rgba(128, 128, 240, 0.05) 100%);
  border: 1px solid rgba(128, 128, 240, 0.2);
}

.stat-card.order-pending {
  background: linear-gradient(135deg, rgba(240, 160, 32, 0.15) 0%, rgba(240, 160, 32, 0.05) 100%);
  border: 1px solid rgba(240, 160, 32, 0.2);
}

.stat-card.order-success {
  background: linear-gradient(135deg, rgba(24, 160, 88, 0.15) 0%, rgba(24, 160, 88, 0.05) 100%);
  border: 1px solid rgba(24, 160, 88, 0.2);
}

.stat-card.order-cancel {
  background: linear-gradient(135deg, rgba(160, 160, 160, 0.15) 0%, rgba(160, 160, 160, 0.05) 100%);
  border: 1px solid rgba(160, 160, 160, 0.2);
}

.stat-trend {
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
  margin-left: 8px;
}

.stat-trend.up {
  color: #18a058;
  background: rgba(24, 160, 88, 0.1);
}

.stat-trend.warning {
  color: #f0a020;
  background: rgba(240, 160, 32, 0.1);
}

.stat-trend.neutral {
  color: #8c8c8c;
  background: rgba(140, 140, 140, 0.1);
}

.chart-card {
  border-radius: 12px;
}

.bar-chart {
  height: 300px;
  display: flex;
  gap: 16px;
  padding: 20px 0;
}

.chart-y-axis {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  font-size: 12px;
  color: #8c8c8c;
  width: 50px;
  text-align: right;
}

.chart-bars {
  flex: 1;
  display: flex;
  align-items: flex-end;
  gap: 12px;
  border-left: 1px solid rgba(255, 255, 255, 0.1);
  padding-left: 12px;
}

.bar-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
  position: relative;
}

.bar-tooltip {
  font-size: 11px;
  color: #ffffff;
  background: rgba(0, 0, 0, 0.8);
  padding: 4px 8px;
  border-radius: 4px;
  margin-bottom: 8px;
  opacity: 0;
  transition: opacity 0.2s;
  white-space: nowrap;
}

.bar-item:hover .bar-tooltip {
  opacity: 1;
}

.bar {
  width: 100%;
  max-width: 40px;
  border-radius: 6px 6px 0 0;
  overflow: hidden;
  transition: height 0.5s ease;
  flex: 1;
  display: flex;
  align-items: flex-end;
}

.bar-fill {
  width: 100%;
  height: 100%;
  background: linear-gradient(180deg, #18a058 0%, #0c6b3a 100%);
  border-radius: 6px 6px 0 0;
}

.bar-item:hover .bar-fill {
  background: linear-gradient(180deg, #36ad6a 0%, #18a058 100%);
}

.bar-label {
  font-size: 12px;
  color: #8c8c8c;
  margin-top: 8px;
}

.ranking-card {
  border-radius: 12px;
}

.ranking-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.ranking-item {
  display: grid;
  grid-template-columns: 32px 1fr 60px;
  gap: 12px;
  align-items: center;
}

.rank {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  background: rgba(140, 140, 140, 0.2);
  color: #8c8c8c;
}

.rank-1 {
  background: linear-gradient(135deg, #ffd700 0%, #ffaa00 100%);
  color: #000;
}

.rank-2 {
  background: linear-gradient(135deg, #c0c0c0 0%, #a0a0a0 100%);
  color: #000;
}

.rank-3 {
  background: linear-gradient(135deg, #cd7f32 0%, #a0652a 100%);
  color: #fff;
}

.product-name {
  font-size: 14px;
  color: #ffffff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.sales-count {
  font-size: 13px;
  color: #8c8c8c;
  text-align: right;
}

.table-card {
  border-radius: 12px;
}
</style>
