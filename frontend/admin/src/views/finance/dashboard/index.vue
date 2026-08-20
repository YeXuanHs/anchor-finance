<template>
  <div class="home-page">
    <!-- 顶部问候语 + 版本 -->
    <div class="welcome-section">
      <div class="welcome-text">
        <h2>{{ $t('dashboard.greeting.hello') }}，{{ adminName }}！{{ $t('dashboard.greeting.todayIs') }}{{ fullDate }}</h2>
      </div>
      <div class="version-info" v-if="systemInfo.version">
        <span>{{ $t('dashboard.currentVersion') }}：{{ systemInfo.version }}</span>
        <span v-if="systemInfo.latest_version"> （{{ $t('dashboard.latestVersion') }}：{{ systemInfo.latest_version }}）</span>
      </div>
    </div>

    <!-- 第一行：4个统计卡片 -->
    <el-row :gutter="16" class="stat-row">
      <!-- 订单概览 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="never" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>{{ $t('dashboard.orderOverview') }}</span>
            </div>
          </template>
          <div class="stat-grid">
            <div class="stat-item" @click="router.push('/order-list')">
              <div class="stat-label">{{ $t('dashboard.todayPaidOrders') }}</div>
              <div class="stat-value">{{ stats.today_paid_orders || 0 }}</div>
              <div class="stat-sub">{{ $t('dashboard.totalOrders') }}: {{ stats.today_total_orders || 0 }}{{ $t('dashboard.unitPiece') }}</div>
            </div>
            <div class="stat-item" @click="router.push('/order-list')">
              <div class="stat-label">{{ $t('dashboard.monthPaidOrders') }}</div>
              <div class="stat-value">{{ stats.month_paid_orders || 0 }}</div>
              <div class="stat-sub">{{ $t('dashboard.totalOrders') }}: {{ stats.month_total_orders || 0 }}{{ $t('dashboard.unitPiece') }}</div>
            </div>
            <div class="stat-item" @click="router.push('/order-list')">
              <div class="stat-label">{{ $t('dashboard.weekPaidOrders') }}</div>
              <div class="stat-value">{{ stats.week_paid_orders || 0 }}</div>
              <div class="stat-sub">{{ $t('dashboard.totalOrders') }}: {{ stats.week_total_orders || 0 }}{{ $t('dashboard.unitPiece') }}</div>
            </div>
            <div class="stat-item" @click="router.push('/order-list')">
              <div class="stat-label">{{ $t('dashboard.weekPendingOrders') }}</div>
              <div class="stat-value">{{ stats.week_pending_orders || 0 }}</div>
              <div class="stat-sub">{{ $t('dashboard.totalOrders') }}: {{ stats.week_total_orders || 0 }}{{ $t('dashboard.unitPiece') }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 代办事项 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="never" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>{{ $t('dashboard.todoItems') }}</span>
            </div>
          </template>
          <div class="todo-list">
            <router-link to="/support-ticket" class="todo-item">
              <el-badge :value="stats.pending_tickets || 0" :max="999" type="warning">
                <el-icon><ChatDotRound /></el-icon>
              </el-badge>
              <span class="todo-label">{{ $t('dashboard.pendingTickets') }}</span>
            </router-link>
            <router-link to="/order-list?status=Pending" class="todo-item">
              <el-badge :value="stats.pending_verification_orders || 0" :max="999" type="primary">
                <el-icon><Document /></el-icon>
              </el-badge>
              <span class="todo-label">{{ $t('dashboard.pendingVerificationOrders') }}</span>
            </router-link>
            <router-link to="/customer-authentication" class="todo-item">
              <el-badge :value="stats.pending_certification || 0" :max="999" type="info">
                <el-icon><User /></el-icon>
              </el-badge>
              <span class="todo-label">{{ $t('dashboard.realNameAuth') }}</span>
            </router-link>
            <router-link to="/customer-withdrawal" class="todo-item">
              <el-badge :value="stats.pending_withdrawals || 0" :max="999" type="danger">
                <el-icon><Money /></el-icon>
              </el-badge>
              <span class="todo-label">{{ $t('dashboard.withdrawalApplication') }}</span>
            </router-link>
          </div>
        </el-card>
      </el-col>

      <!-- 收入概览 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="never" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>{{ $t('dashboard.incomeOverview') }}</span>
            </div>
          </template>
          <div class="income-grid">
            <div class="income-item">
              <div class="income-value">¥{{ formatMoney(stats.today_income) }}</div>
              <div class="income-label">{{ $t('dashboard.todayIncome') }}</div>
            </div>
            <div class="income-item">
              <div class="income-value">¥{{ formatMoney(stats.month_income) }}</div>
              <div class="income-label">{{ $t('dashboard.monthIncome') }}</div>
            </div>
            <div class="income-item">
              <div class="income-value">¥{{ formatMoney(stats.year_income) }}</div>
              <div class="income-label">{{ $t('dashboard.yearIncome') }}</div>
            </div>
            <div class="income-item">
              <div class="income-value">¥{{ formatMoney(stats.total_income) }}</div>
              <div class="income-label">{{ $t('dashboard.totalIncome') }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 交易统计 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="never" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>{{ $t('dashboard.transactionStats') }}</span>
            </div>
          </template>
          <div ref="transactionChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 第二行 -->
    <el-row :gutter="16" class="info-row">
      <!-- 在线管理员 -->
      <el-col :xs="24" :sm="12" :lg="8">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('dashboard.onlineAdmins') }}</span>
            </div>
          </template>
          <div class="admin-count">{{ $t('dashboard.onlineCount') }}: {{ onlineAdmins.length }}</div>
          <div class="admin-list">
            <div v-for="admin in onlineAdmins" :key="admin.id" class="admin-item">
              <el-avatar :size="32" class="admin-avatar">
                {{ admin.username?.charAt(0)?.toUpperCase() }}
              </el-avatar>
              <div class="admin-info">
                <div class="admin-name">{{ admin.username }}</div>
                <div class="admin-time">{{ $t('dashboard.recently') }}：{{ admin.last_active || $t('dashboard.justNow') }}</div>
              </div>
            </div>
            <el-empty v-if="onlineAdmins.length === 0" :description="$t('dashboard.noOnlineAdmins')" :image-size="48" />
          </div>
        </el-card>
      </el-col>

      <!-- 即将到期 -->
      <el-col :xs="24" :sm="12" :lg="8">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('dashboard.expiringSoon') }}</span>
            </div>
          </template>
          <div class="expire-grid">
            <router-link
              v-for="item in expiringProducts"
              :key="item.type"
              :to="'/customer-product?from=home-page&productType=' + item.type"
              class="expire-item"
            >
              <span class="expire-count">{{ item.count }}</span>
              <span class="expire-type">{{ item.name }}</span>
            </router-link>
            <el-empty v-if="expiringProducts.length === 0" :description="$t('dashboard.noExpiring')" :image-size="48" />
          </div>
        </el-card>
      </el-col>

      <!-- 本月销量排行 -->
      <el-col :xs="24" :sm="24" :lg="8">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('dashboard.monthSalesRanking') }}</span>
            </div>
          </template>
          <div ref="salesChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 第三行 -->
    <el-row :gutter="16" class="info-row">
      <!-- 客户概况 -->
      <el-col :xs="24" :sm="24" :lg="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('dashboard.clientOverview') }}</span>
            </div>
          </template>
          <div class="client-stats">
            <span>{{ stats.new_clients_today || 0 }}{{ $t('dashboard.todayNew') }}</span>
            <span class="stats-sep">{{ stats.online_clients || 0 }}{{ $t('dashboard.currentOnline') }}</span>
            <span>{{ stats.total_clients || 0 }}{{ $t('dashboard.totalSystem') }}</span>
          </div>
          <div class="client-list">
            <div v-for="client in onlineClients" :key="client.id" class="client-item">
              <div class="client-main">
                <router-link :to="'/customer-view/abstract?id=' + client.id" class="client-name">
                  {{ client.username }}
                </router-link>
                <span class="client-ip">{{ client.ip }}</span>
              </div>
              <div class="client-time">{{ client.last_active }}</div>
            </div>
            <el-empty v-if="onlineClients.length === 0" :description="$t('dashboard.noOnlineClients')" :image-size="48" />
          </div>
        </el-card>
      </el-col>

      <!-- 待处理工单 -->
      <el-col :xs="24" :sm="24" :lg="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('dashboard.pendingTicketList') }}</span>
            </div>
          </template>
          <div class="ticket-list">
            <div v-for="ticket in recentTickets" :key="ticket.id" class="ticket-item">
              <div class="ticket-main">
                <el-tag size="small" type="success">{{ $t('dashboard.opening') }}</el-tag>
                <router-link :to="'/support-ticket-detail?id=' + ticket.id + '&tid=' + ticket.ticket_no" class="ticket-subject">
                  #{{ ticket.ticket_no }}-{{ ticket.subject }}
                </router-link>
                <span class="ticket-user">{{ ticket.client_name }}</span>
              </div>
              <div class="ticket-time">{{ ticket.last_reply_time || ticket.created_at }}</div>
            </div>
            <el-empty v-if="recentTickets.length === 0" :description="$t('dashboard.noPendingTickets')" :image-size="48" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 底部 -->
    <div class="footer-text">{{ $t('dashboard.footer') }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { ChatDotRound, Document, User, Money } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import * as echarts from 'echarts'
import request from '@/utils/http'

const router = useRouter()
const { t } = useI18n()

const version = ref('1.0.0')
const systemInfo = ref<any>({})

const adminName = computed(() => {
  return localStorage.getItem('admin_name') || 'admin'
})

const fullDate = computed(() => {
  const now = new Date()
  const weekdays = [t('dashboard.weekday.sunday'), t('dashboard.weekday.monday'), t('dashboard.weekday.tuesday'), t('dashboard.weekday.wednesday'), t('dashboard.weekday.thursday'), t('dashboard.weekday.friday'), t('dashboard.weekday.saturday')]
  const year = now.getFullYear()
  const month = now.getMonth() + 1
  const day = now.getDate()
  const weekday = weekdays[now.getDay()]
  return `${year}年${month}月${day}日，${weekday}`
})

const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return amount.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const stats = ref<any>({
  today_paid_orders: 0,
  today_total_orders: 0,
  month_paid_orders: 0,
  month_total_orders: 0,
  week_paid_orders: 0,
  week_total_orders: 0,
  week_pending_orders: 0,
  pending_tickets: 0,
  pending_verification_orders: 0,
  pending_certification: 0,
  pending_withdrawals: 0,
  today_income: 0,
  month_income: 0,
  year_income: 0,
  total_income: 0,
  new_clients_today: 0,
  online_clients: 0,
  total_clients: 0,
})

const onlineAdmins = ref<any[]>([])
const onlineClients = ref<any[]>([])
const expiringProducts = ref<any[]>([])
const recentTickets = ref<any[]>([])

const transactionChartRef = ref<HTMLElement>()
const salesChartRef = ref<HTMLElement>()
let transactionChart: echarts.ECharts | null = null
let salesChart: echarts.ECharts | null = null

const fetchDashboardData = async () => {
  try {
    const data = await request.get({ url: '/api/admin/dashboard/stats' })
    if (data) {
      stats.value = { ...stats.value, ...data }
      setTimeout(() => {
        initTransactionChart()
        initSalesChart()
      }, 300)
    }
  } catch (error) {
    console.error('获取仪表盘数据失败:', error)
  }
}

const fetchSystemInfo = async () => {
  try {
    const data = await request.get({ url: '/api/admin/system/info' })
    systemInfo.value = data || {}
  } catch (error) {
    console.error('获取系统信息失败:', error)
  }
}

const fetchOnlineAdmins = async () => {
  try {
    const data = await request.get({ url: '/api/admin/online-admins' })
    onlineAdmins.value = data || []
  } catch (error) {
    console.error('获取在线管理员失败:', error)
  }
}

const fetchOnlineClients = async () => {
  try {
    const data = await request.get({ url: '/api/admin/dashboard/online-clients' })
    onlineClients.value = data || []
  } catch (error) {
    console.error('获取在线客户失败:', error)
  }
}

const fetchExpiringProducts = async () => {
  try {
    const data = await request.get({ url: '/api/admin/dashboard/expiring-products' })
    expiringProducts.value = data || []
  } catch (error) {
    console.error('获取即将到期产品失败:', error)
  }
}

const fetchRecentTickets = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/tickets',
      params: { page: 1, page_size: 5, status: '0' }
    })
    recentTickets.value = data?.list || data || []
  } catch (error) {
    console.error('获取最近工单失败:', error)
  }
}

const initTransactionChart = () => {
  if (!transactionChartRef.value) return
  if (!transactionChart) transactionChart = echarts.init(transactionChartRef.value)

  const option = {
    tooltip: { trigger: 'item' },
    legend: { bottom: 0, textStyle: { fontSize: 12, color: '#86909C' } },
    series: [{
      type: 'pie',
      radius: ['40%', '65%'],
      center: ['50%', '45%'],
      avoidLabelOverlap: false,
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: false },
      emphasis: { label: { show: true, fontSize: 14, fontWeight: 'bold' } },
      data: [
        { value: stats.value.today_income || 0, name: t('dashboard.todayIncome'), itemStyle: { color: 'var(--el-color-primary)' } },
        { value: stats.value.month_income || 0, name: t('dashboard.monthIncome'), itemStyle: { color: '#36D391' } },
        { value: stats.value.year_income || 0, name: t('dashboard.yearIncome'), itemStyle: { color: '#F59E0B' } },
      ]
    }]
  }

  transactionChart.setOption(option, true)
}

const initSalesChart = () => {
  if (!salesChartRef.value) return
  if (!salesChart) salesChart = echarts.init(salesChartRef.value)

  const option = {
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'value',
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: '#86909C' },
      splitLine: { lineStyle: { type: 'dashed', color: '#E5E6EB' } }
    },
    yAxis: {
      type: 'category',
      data: (stats.value.sales_ranking || []).map((item: any) => item.name).reverse(),
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: '#4E5969' }
    },
    series: [{
      type: 'bar',
      data: (stats.value.sales_ranking || []).map((item: any) => item.value).reverse(),
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
          { offset: 0, color: 'var(--el-color-primary)' },
          { offset: 1, color: 'var(--el-color-primary-light-3)' }
        ]),
        borderRadius: [0, 4, 4, 0]
      },
      barWidth: 20
    }]
  }

  salesChart.setOption(option, true)
}

const handleResize = () => {
  transactionChart?.resize()
  salesChart?.resize()
}

onMounted(() => {
  fetchDashboardData()
  fetchSystemInfo()
  fetchOnlineAdmins()
  fetchOnlineClients()
  fetchExpiringProducts()
  fetchRecentTickets()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  transactionChart?.dispose()
  salesChart?.dispose()
})
</script>

<style scoped lang="scss">
.home-page {
  padding: 16px;
}

.welcome-section {
  margin-bottom: 16px;
  padding: 16px 20px;
  background: linear-gradient(135deg, var(--el-color-primary) 0%, var(--el-color-primary-light-3) 100%);
  border-radius: 8px;
  color: #fff;

  .welcome-text {
    h2 {
      margin: 0;
      font-size: 18px;
      font-weight: 500;
    }
  }

  .version-info {
    margin-top: 6px;
    font-size: 13px;
    opacity: 0.85;
  }
}

.stat-row {
  margin-bottom: 16px;

  .el-card {
    :deep(.el-card__header) {
      padding: 12px 16px;
      border-bottom: 1px solid #f0f0f0;
    }
  }
}

.info-row {
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  font-weight: 500;
  color: #1d2129;
}

.stat-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;

  .stat-item {
    padding: 8px;
    border-radius: 6px;
    cursor: pointer;
    transition: background 0.2s;

    &:hover {
      background: #f5f5f5;
    }
  }

  .stat-label {
    font-size: 12px;
    color: #86909c;
    margin-bottom: 4px;
  }

  .stat-value {
    font-size: 20px;
    font-weight: 600;
    color: #1d2129;
    margin-bottom: 2px;
  }

  .stat-sub {
    font-size: 11px;
    color: #86909c;
  }
}

.todo-list {
  display: flex;
  flex-direction: column;
  gap: 8px;

  .todo-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px 12px;
    border-radius: 6px;
    text-decoration: none;
    color: inherit;
    transition: background 0.2s;

    &:hover {
      background: #f5f5f5;
    }

    .el-icon {
      font-size: 20px;
      color: var(--el-color-primary);
    }

    .todo-label {
      font-size: 13px;
      color: #4e5969;
    }
  }
}

.income-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;

  .income-item {
    text-align: center;
    padding: 8px;

    .income-value {
      font-size: 18px;
      font-weight: 600;
      color: var(--el-color-primary);
      margin-bottom: 4px;
    }

    .income-label {
      font-size: 12px;
      color: #86909c;
    }
  }
}

.chart-container {
  height: 200px;
}

.admin-count {
  font-size: 13px;
  color: #86909c;
  margin-bottom: 12px;
}

.admin-list {
  max-height: 280px;
  overflow-y: auto;
}

.admin-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 0;
  border-bottom: 1px solid #f2f3f5;

  &:last-child {
    border-bottom: none;
  }

  .admin-avatar {
    background: var(--el-color-primary);
    color: #fff;
    font-size: 14px;
    font-weight: 500;
  }

  .admin-info {
    flex: 1;
  }

  .admin-name {
    font-size: 13px;
    color: #1d2129;
  }

  .admin-time {
    font-size: 11px;
    color: #86909c;
  }
}

.expire-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;

  .expire-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    border-radius: 6px;
    text-decoration: none;
    color: inherit;
    transition: background 0.2s;

    &:hover {
      background: #f5f5f5;
    }

    .expire-count {
      font-size: 16px;
      font-weight: 600;
      color: var(--el-color-primary);
    }

    .expire-type {
      font-size: 13px;
      color: #4e5969;
    }
  }
}

.client-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
  font-size: 13px;
  color: #4e5969;

  .stats-sep {
    &::before {
      content: '';
      display: inline-block;
      width: 1px;
      height: 12px;
      background: #e5e6eb;
      margin-right: 16px;
      vertical-align: middle;
    }
    &::after {
      content: '';
      display: inline-block;
      width: 1px;
      height: 12px;
      background: #e5e6eb;
      margin-left: 16px;
      vertical-align: middle;
    }
  }
}

.client-list {
  max-height: 300px;
  overflow-y: auto;
}

.client-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f2f3f5;

  &:last-child {
    border-bottom: none;
  }

  .client-main {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 1;
  }

  .client-name {
    font-size: 13px;
    color: var(--el-color-primary);
    text-decoration: none;
    font-weight: 500;

    &:hover {
      text-decoration: underline;
    }
  }

  .client-ip {
    font-size: 12px;
    color: #86909c;
  }

  .client-time {
    font-size: 12px;
    color: #86909c;
    white-space: nowrap;
  }
}

.ticket-list {
  max-height: 300px;
  overflow-y: auto;
}

.ticket-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f2f3f5;

  &:last-child {
    border-bottom: none;
  }

  .ticket-main {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
    min-width: 0;
  }

  .ticket-subject {
    font-size: 13px;
    color: var(--el-color-primary);
    text-decoration: none;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;

    &:hover {
      text-decoration: underline;
    }
  }

  .ticket-user {
    font-size: 12px;
    color: #86909c;
    white-space: nowrap;
  }

  .ticket-time {
    font-size: 12px;
    color: #86909c;
    white-space: nowrap;
    margin-left: 8px;
  }
}

.footer-text {
  text-align: center;
  padding: 16px 0;
  font-size: 12px;
  color: #c9cdd4;
}
</style>
