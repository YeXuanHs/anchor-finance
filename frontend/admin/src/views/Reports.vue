<template>
  <div class="reports">
    <!-- Income Stats -->
    <n-grid :cols="4" :x-gap="16" :y-gap="16" responsive="screen" :item-responsive="true">
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card" :bordered="false" rounded>
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">今日收入</span>
              <span class="stat-value">¥{{ incomeStats.today.toLocaleString() }}</span>
              <span class="stat-change" :class="incomeStats.todayChange >= 0 ? 'up' : 'down'">
                <n-icon size="14">
                  <TrendUpIcon v-if="incomeStats.todayChange >= 0" />
                  <TrendDownIcon v-else />
                </n-icon>
                {{ Math.abs(incomeStats.todayChange) }}%
              </span>
            </div>
            <div class="stat-icon blue">
              <n-icon size="28"><CashIcon /></n-icon>
            </div>
          </div>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card" :bordered="false" rounded>
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">本月收入</span>
              <span class="stat-value">¥{{ incomeStats.month.toLocaleString() }}</span>
              <span class="stat-change" :class="incomeStats.monthChange >= 0 ? 'up' : 'down'">
                <n-icon size="14">
                  <TrendUpIcon v-if="incomeStats.monthChange >= 0" />
                  <TrendDownIcon v-else />
                </n-icon>
                {{ Math.abs(incomeStats.monthChange) }}%
              </span>
            </div>
            <div class="stat-icon cyan">
              <n-icon size="28"><WalletIcon /></n-icon>
            </div>
          </div>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card" :bordered="false" rounded>
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">总收入</span>
              <span class="stat-value">¥{{ incomeStats.total.toLocaleString() }}</span>
            </div>
            <div class="stat-icon green">
              <n-icon size="28"><TrendingUpIcon /></n-icon>
            </div>
          </div>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card" :bordered="false" rounded>
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">待收金额</span>
              <span class="stat-value text-orange">¥{{ incomeStats.pending.toLocaleString() }}</span>
            </div>
            <div class="stat-icon orange">
              <n-icon size="28"><TimeIcon /></n-icon>
            </div>
          </div>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- User & Order Stats -->
    <n-grid :cols="7" :x-gap="16" :y-gap="16" style="margin-top: 16px" responsive="screen" :item-responsive="true">
      <n-gi span="7 m:4 l:3">
        <n-card title="用户统计" :bordered="false" rounded class="section-card">
          <n-grid :cols="3" :x-gap="12">
            <n-gi>
              <div class="mini-stat">
                <n-statistic label="总用户" :value="userStats.total">
                  <template #suffix>
                    <n-icon size="16" color="#1890ff"><PeopleIcon /></n-icon>
                  </template>
                </n-statistic>
              </div>
            </n-gi>
            <n-gi>
              <div class="mini-stat">
                <n-statistic label="新增用户" :value="userStats.newToday">
                  <template #prefix>
                    <n-icon size="16" color="#52c41a"><PersonAddIcon /></n-icon>
                  </template>
                </n-statistic>
              </div>
            </n-gi>
            <n-gi>
              <div class="mini-stat">
                <n-statistic label="活跃用户" :value="userStats.active">
                  <template #prefix>
                    <n-icon size="16" color="#fa8c16"><FlashIcon /></n-icon>
                  </template>
                </n-statistic>
              </div>
            </n-gi>
          </n-grid>
        </n-card>
      </n-gi>
      <n-gi span="7 m:3 l:4">
        <n-card title="订单统计" :bordered="false" rounded class="section-card">
          <n-grid :cols="4" :x-gap="12">
            <n-gi>
              <div class="mini-stat">
                <n-statistic label="总订单" :value="orderStats.total" />
              </div>
            </n-gi>
            <n-gi>
              <div class="mini-stat">
                <n-statistic label="待付款">
                  <template #default>
                    <span style="color: #fa8c16">{{ orderStats.pending }}</span>
                  </template>
                </n-statistic>
              </div>
            </n-gi>
            <n-gi>
              <div class="mini-stat">
                <n-statistic label="已完成">
                  <template #default>
                    <span style="color: #52c41a">{{ orderStats.completed }}</span>
                  </template>
                </n-statistic>
              </div>
            </n-gi>
            <n-gi>
              <div class="mini-stat">
                <n-statistic label="已取消">
                  <template #default>
                    <span style="color: #d9d9d9">{{ orderStats.cancelled }}</span>
                  </template>
                </n-statistic>
              </div>
            </n-gi>
          </n-grid>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- Charts Row -->
    <n-grid :cols="3" :x-gap="16" style="margin-top: 16px" responsive="screen" :item-responsive="true">
      <n-gi span="3 l:2">
        <n-card title="收入趋势" :bordered="false" rounded class="chart-card">
          <template #header-extra>
            <n-tabs v-model:value="chartPeriod" size="small" type="segment" animated>
              <n-tab-pane name="week" tab="近7天" />
              <n-tab-pane name="month" tab="近30天" />
              <n-tab-pane name="quarter" tab="近90天" />
            </n-tabs>
          </template>
          <div class="bar-chart">
            <div class="bar-chart-inner">
              <div v-for="(item, index) in revenueData" :key="index" class="bar-item">
                <div class="bar-wrapper">
                  <div
                    class="bar"
                    :style="{ height: (item.value / maxRevenue) * 100 + '%' }"
                  >
                    <span class="bar-tooltip">¥{{ item.value.toLocaleString() }}</span>
                  </div>
                </div>
                <span class="bar-label">{{ item.label }}</span>
              </div>
            </div>
            <div class="bar-y-axis">
              <span>¥{{ maxRevenue.toLocaleString() }}</span>
              <span>¥{{ Math.round(maxRevenue * 0.75).toLocaleString() }}</span>
              <span>¥{{ Math.round(maxRevenue * 0.5).toLocaleString() }}</span>
              <span>¥{{ Math.round(maxRevenue * 0.25).toLocaleString() }}</span>
              <span>¥0</span>
            </div>
          </div>
        </n-card>
      </n-gi>
      <n-gi span="3 l:1">
        <n-card title="产品销量排行" :bordered="false" rounded class="chart-card">
          <div class="rank-list">
            <div v-for="(item, index) in productRank" :key="item.name" class="rank-item">
              <div class="rank-badge" :class="{ 'top-three': index < 3 }">{{ index + 1 }}</div>
              <div class="rank-info">
                <span class="rank-name">{{ item.name }}</span>
                <div class="rank-bar-bg">
                  <div
                    class="rank-bar-fill"
                    :style="{ width: (item.sales / productRank[0].sales) * 100 + '%' }"
                  />
                </div>
              </div>
              <span class="rank-count">{{ item.sales }}单</span>
            </div>
          </div>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- Recent Transactions -->
    <n-grid :cols="1" :x-gap="16" style="margin-top: 16px">
      <n-gi>
        <n-card title="最近交易记录" :bordered="false" rounded>
          <template #header-extra>
            <n-date-picker v-model:value="dateRange" type="daterange" size="small" style="margin-right: 12px" />
            <n-button text type="primary">导出报表</n-button>
          </template>
          <n-data-table
            :columns="transactionColumns"
            :data="recentTransactions"
            :bordered="false"
            :pagination="pagination"
            size="small"
          />
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { h, reactive, ref, computed } from 'vue'
import {
  TrendOutline as TrendUpIcon,
  TrendDownOutline as TrendDownIcon,
  CashOutline as CashIcon,
  WalletOutline as WalletIcon,
  TrendingUpOutline as TrendingUpIcon,
  TimeOutline as TimeIcon,
  PeopleOutline as PeopleIcon,
  PersonAddOutline as PersonAddIcon,
  FlashOutline as FlashIcon,
} from '@vicons/ionicons5'
import { NTag, NButton, type DataTableColumns } from 'naive-ui'

// ---- Income Stats ----
const incomeStats = reactive({
  today: 15680,
  todayChange: 12.5,
  month: 368520,
  monthChange: 8.3,
  total: 4256800,
  pending: 42600,
})

// ---- User Stats ----
const userStats = reactive({
  total: 12580,
  newToday: 34,
  active: 8920,
})

// ---- Order Stats ----
const orderStats = reactive({
  total: 28650,
  pending: 340,
  completed: 27130,
  cancelled: 1180,
})

// ---- Revenue Bar Chart ----
const chartPeriod = ref('week')

function generateBarData(count: number): { label: string; value: number }[] {
  const data: { label: string; value: number }[] = []
  const now = new Date()
  for (let i = count - 1; i >= 0; i--) {
    const d = new Date(now)
    d.setDate(d.getDate() - i)
    data.push({
      label: `${d.getMonth() + 1}/${d.getDate()}`,
      value: Math.floor(Math.random() * 15000 + 3000),
    })
  }
  return data
}

const revenueData = computed(() => {
  const count = chartPeriod.value === 'week' ? 7 : chartPeriod.value === 'month' ? 30 : 90
  return generateBarData(count)
})

const maxRevenue = computed(() => {
  const max = Math.max(...revenueData.value.map((d) => d.value))
  return Math.ceil(max / 1000) * 1000
})

// ---- Product Rank ----
const productRank = ref([
  { name: '基础版主机', sales: 1280 },
  { name: '4核8G云服务器', sales: 960 },
  { name: '高级版主机', sales: 856 },
  { name: '1核2G云服务器', sales: 720 },
  { name: '企业版主机', sales: 580 },
  { name: '域名注册', sales: 420 },
  { name: 'SSL证书', sales: 310 },
])

// ---- Recent Transactions ----
const dateRange = ref<[number, number] | null>(null)

const statusMap: Record<string, { label: string; color: string }> = {
  success: { label: '成功', color: 'success' },
  pending: { label: '处理中', color: 'warning' },
  failed: { label: '失败', color: 'error' },
  refund: { label: '已退款', color: 'info' },
}

const recentTransactions = ref([
  { id: 'TXN20260727001', user: '张三', type: '订单支付', product: '基础版主机', amount: 299, status: 'success', time: '2026-07-27 09:30' },
  { id: 'TXN20260727002', user: '李四', type: '续费', product: '高级版主机', amount: 599, status: 'success', time: '2026-07-27 09:15' },
  { id: 'TXN20260727003', user: '王五', type: '新购', product: '企业版主机', amount: 1299, status: 'pending', time: '2026-07-27 08:45' },
  { id: 'TXN20260727004', user: '赵六', type: '订单支付', product: '1核2G云服务器', amount: 89, status: 'success', time: '2026-07-26 22:18' },
  { id: 'TXN20260727005', user: '孙七', type: '退款', product: '4核8G云服务器', amount: -399, status: 'refund', time: '2026-07-26 20:55' },
  { id: 'TXN20260727006', user: '周八', type: '新购', product: '域名注册', amount: 68, status: 'failed', time: '2026-07-26 18:30' },
  { id: 'TXN20260727007', user: '吴九', type: '续费', product: 'SSL证书', amount: 199, status: 'success', time: '2026-07-26 16:42' },
  { id: 'TXN20260727008', user: '郑十', type: '订单支付', product: '基础版主机', amount: 299, status: 'success', time: '2026-07-26 14:10' },
])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => { pagination.page = page },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
  },
})

const transactionColumns: DataTableColumns<any> = [
  { title: '交易号', key: 'id', width: 160, ellipsis: { tooltip: true } },
  { title: '用户', key: 'user', width: 80 },
  { title: '类型', key: 'type', width: 100 },
  { title: '产品', key: 'product', ellipsis: { tooltip: true } },
  {
    title: '金额',
    key: 'amount',
    width: 110,
    render: (row) => {
      const color = row.amount >= 0 ? '#1890ff' : '#f5222d'
      const prefix = row.amount >= 0 ? '+' : ''
      return h('span', { style: `font-weight:600;color:${color}` }, `¥${prefix}${row.amount}`)
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: (row) => {
      const s = statusMap[row.status]
      return h(NTag, { type: s.color as any, size: 'small', round: true, bordered: false }, { default: () => s.label })
    },
  },
  { title: '时间', key: 'time', width: 155 },
  {
    title: '操作',
    key: 'action',
    width: 80,
    render: () =>
      h(NButton, { text: true, type: 'primary', size: 'small' }, { default: () => '详情' }),
  },
]
</script>

<style scoped>
.reports {
  padding: 0;
}

.stat-card {
  border-radius: 12px;
  transition: box-shadow 0.3s, transform 0.2s;
}

.stat-card:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.stat-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-label {
  font-size: 13px;
  color: #8c8c8c;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a2e;
  line-height: 1.2;
}

.stat-value.text-orange {
  color: #fa8c16;
}

.stat-change {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  font-size: 13px;
  margin-top: 6px;
}

.stat-change.up {
  color: #52c41a;
}

.stat-change.down {
  color: #f5222d;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon.blue {
  background: rgba(24, 144, 255, 0.1);
  color: #1890ff;
}

.stat-icon.cyan {
  background: rgba(19, 194, 194, 0.1);
  color: #13c2c2;
}

.stat-icon.green {
  background: rgba(82, 196, 26, 0.1);
  color: #52c41a;
}

.stat-icon.orange {
  background: rgba(250, 140, 22, 0.1);
  color: #fa8c16;
}

.section-card {
  border-radius: 12px;
}

.mini-stat {
  text-align: center;
  padding: 8px 0;
}

.chart-card {
  border-radius: 12px;
}

/* Bar Chart */
.bar-chart {
  position: relative;
  display: flex;
  height: 300px;
  padding-bottom: 28px;
}

.bar-chart-inner {
  flex: 1;
  display: flex;
  align-items: flex-end;
  gap: 4px;
  padding: 0 8px;
  border-bottom: 1px solid #f0f0f0;
  border-left: 1px solid #f0f0f0;
}

.bar-y-axis {
  position: absolute;
  left: -8px;
  top: 0;
  bottom: 28px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  font-size: 11px;
  color: #bfbfbf;
  transform: translateX(-100%);
  padding-right: 8px;
  white-space: nowrap;
}

.bar-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  max-width: 60px;
}

.bar-wrapper {
  flex: 1;
  width: 100%;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.bar {
  width: 70%;
  min-width: 8px;
  max-width: 40px;
  background: linear-gradient(180deg, #1890ff 0%, #40a9ff 100%);
  border-radius: 4px 4px 0 0;
  transition: height 0.5s ease;
  position: relative;
  cursor: pointer;
}

.bar:hover {
  background: linear-gradient(180deg, #096dd9 0%, #1890ff 100%);
}

.bar:hover .bar-tooltip {
  opacity: 1;
  transform: translateX(-50%) translateY(-4px);
}

.bar-tooltip {
  position: absolute;
  top: -28px;
  left: 50%;
  transform: translateX(-50%) translateY(0);
  background: rgba(0, 0, 0, 0.75);
  color: #fff;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
  white-space: nowrap;
  opacity: 0;
  transition: all 0.2s ease;
  pointer-events: none;
}

.bar-label {
  font-size: 11px;
  color: #8c8c8c;
  margin-top: 6px;
  white-space: nowrap;
}

/* Rank List */
.rank-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.rank-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.rank-badge {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: #8c8c8c;
  background: #f5f5f5;
  flex-shrink: 0;
}

.rank-badge.top-three {
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  color: #fff;
}

.rank-info {
  flex: 1;
  min-width: 0;
}

.rank-name {
  font-size: 13px;
  color: #1a1a2e;
  display: block;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rank-bar-bg {
  height: 6px;
  background: #f5f5f5;
  border-radius: 3px;
  overflow: hidden;
}

.rank-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #1890ff, #40a9ff);
  border-radius: 3px;
  transition: width 0.6s ease;
}

.rank-count {
  font-size: 13px;
  color: #8c8c8c;
  flex-shrink: 0;
  width: 48px;
  text-align: right;
}
</style>
