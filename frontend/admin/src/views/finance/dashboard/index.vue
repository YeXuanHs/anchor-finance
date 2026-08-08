<template>
  <div class="home-page">
    <!-- 顶部问候语 -->
    <div class="welcome-section">
      <div class="welcome-text">
        <h2>{{ greeting }}，{{ adminName }}！</h2>
        <p class="version-info">锚点财务 v{{ version }} | 兼容智简魔方(zjmf)</p>
      </div>
      <div class="welcome-date">
        <el-icon><Calendar /></el-icon>
        <span>{{ currentDate }}</span>
      </div>
    </div>

    <!-- 第一行：4个统计卡片 -->
    <el-row :gutter="16" class="stat-row">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card" @click="router.push('/order-list')">
          <div class="stat-card-content">
            <div class="stat-info">
              <div class="stat-label">待处理订单</div>
              <div class="stat-value">{{ stats.pending_orders || 0 }}</div>
              <div class="stat-footer">
                <span class="stat-link">查看详情 →</span>
              </div>
            </div>
            <div class="stat-icon order-icon">
              <el-icon :size="32"><Document /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card" @click="router.push('/support-ticket')">
          <div class="stat-card-content">
            <div class="stat-info">
              <div class="stat-label">待处理工单</div>
              <div class="stat-value">{{ stats.pending_tickets || 0 }}</div>
              <div class="stat-footer">
                <span class="stat-link">查看详情 →</span>
              </div>
            </div>
            <div class="stat-icon ticket-icon">
              <el-icon :size="32"><ChatDotRound /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card" @click="router.push('/business-statement')">
          <div class="stat-card-content">
            <div class="stat-info">
              <div class="stat-label">今日收入</div>
              <div class="stat-value">¥{{ formatMoney(stats.today_income) }}</div>
              <div class="stat-footer">
                <span :class="stats.income_change >= 0 ? 'text-green' : 'text-red'">
                  {{ stats.income_change >= 0 ? '+' : '' }}{{ stats.income_change?.toFixed(1) || '0' }}%
                </span>
                较昨日
              </div>
            </div>
            <div class="stat-icon income-icon">
              <el-icon :size="32"><Money /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card" @click="router.push('/customer-list')">
          <div class="stat-card-content">
            <div class="stat-info">
              <div class="stat-label">客户总数</div>
              <div class="stat-value">{{ stats.total_clients || 0 }}</div>
              <div class="stat-footer">
                <span class="text-green">+{{ stats.new_clients_today || 0 }}</span>
                今日新增
              </div>
            </div>
            <div class="stat-icon client-icon">
              <el-icon :size="32"><User /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 第二行：收入趋势 + 产品分布 -->
    <el-row :gutter="16" class="chart-row">
      <el-col :span="16">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>收入趋势</span>
              <el-radio-group v-model="chartPeriod" size="small" @change="fetchIncomeData">
                <el-radio-button label="week">近7天</el-radio-button>
                <el-radio-button label="month">近30天</el-radio-button>
                <el-radio-button label="year">近1年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="incomeChartRef" class="chart-container" v-loading="incomeLoading"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <span>产品分布</span>
          </template>
          <div ref="productChartRef" class="chart-container" v-loading="productLoading"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 第三行：最近订单 + 最近工单 -->
    <el-row :gutter="16" class="list-row">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>最近订单</span>
              <el-button type="primary" link @click="router.push('/order-list')">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentOrders" style="width: 100%" v-loading="ordersLoading" size="small">
            <el-table-column prop="order_no" label="订单号" width="150" />
            <el-table-column prop="client_name" label="客户" width="100" />
            <el-table-column prop="product_name" label="产品" />
            <el-table-column prop="amount" label="金额" width="100">
              <template #default="{ row }">
                ¥{{ formatMoney(row.amount) }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="getOrderStatusType(row.status)" size="small">
                  {{ getOrderStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>最近工单</span>
              <el-button type="primary" link @click="router.push('/support-ticket')">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentTickets" style="width: 100%" v-loading="ticketsLoading" size="small">
            <el-table-column prop="ticket_no" label="工单号" width="120" />
            <el-table-column prop="subject" label="主题" />
            <el-table-column prop="client_name" label="客户" width="100" />
            <el-table-column prop="priority" label="优先级" width="80">
              <template #default="{ row }">
                <el-tag :type="getPriorityType(row.priority)" size="small">
                  {{ getPriorityText(row.priority) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="getTicketStatusType(row.status)" size="small">
                  {{ getTicketStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 第四行：系统日志 + 在线管理员 + 即将到期 -->
    <el-row :gutter="16" class="info-row">
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>系统日志</span>
              <el-button type="primary" link @click="router.push('/system-log')">查看全部</el-button>
            </div>
          </template>
          <div class="log-list">
            <div v-for="log in recentLogs" :key="log.id" class="log-item">
              <span class="log-time">{{ log.created_at }}</span>
              <span class="log-content">{{ log.content }}</span>
            </div>
            <el-empty v-if="recentLogs.length === 0" description="暂无日志" :image-size="60" />
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <span>在线管理员</span>
          </template>
          <div class="admin-list">
            <div v-for="admin in onlineAdmins" :key="admin.id" class="admin-item">
              <el-avatar :size="32" :src="admin.avatar">{{ admin.username?.charAt(0) }}</el-avatar>
              <div class="admin-info">
                <div class="admin-name">{{ admin.username }}</div>
                <div class="admin-time">{{ admin.last_active_at }}</div>
              </div>
            </div>
            <el-empty v-if="onlineAdmins.length === 0" description="暂无在线管理员" :image-size="60" />
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>即将到期产品</span>
              <el-button type="primary" link @click="router.push('/customer-product')">查看全部</el-button>
            </div>
          </template>
          <div class="expire-list">
            <div v-for="item in expiringProducts" :key="item.id" class="expire-item">
              <div class="expire-info">
                <div class="expire-name">{{ item.product_name }}</div>
                <div class="expire-client">{{ item.client_name }}</div>
              </div>
              <div class="expire-date">
                <el-tag :type="item.days_left <= 3 ? 'danger' : 'warning'" size="small">
                  {{ item.days_left }}天后到期
                </el-tag>
              </div>
            </div>
            <el-empty v-if="expiringProducts.length === 0" description="暂无即将到期产品" :image-size="60" />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { Calendar, Document, ChatDotRound, Money, User } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import request from '@/utils/http'

const router = useRouter()
const chartPeriod = ref('week')
const version = ref('1.0.0')

// 问候语
const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '凌晨好'
  if (hour < 9) return '早上好'
  if (hour < 12) return '上午好'
  if (hour < 14) return '中午好'
  if (hour < 17) return '下午好'
  if (hour < 19) return '傍晚好'
  return '晚上好'
})

// 管理员名称
const adminName = computed(() => {
  return localStorage.getItem('admin_name') || '管理员'
})

// 当前日期
const currentDate = computed(() => {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  const weekDays = ['日', '一', '二', '三', '四', '五', '六']
  const weekDay = weekDays[now.getDay()]
  return `${year}年${month}月${day}日 星期${weekDay}`
})

// 格式化金额
const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return amount.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// 加载状态
const incomeLoading = ref(false)
const productLoading = ref(false)
const ordersLoading = ref(false)
const ticketsLoading = ref(false)

// 统计数据
const stats = ref({
  pending_orders: 0,
  pending_tickets: 0,
  today_income: 0,
  income_change: 0,
  total_clients: 0,
  new_clients_today: 0
})

// 最近订单
const recentOrders = ref([])

// 最近工单
const recentTickets = ref([])

// 最近日志
const recentLogs = ref([])

// 在线管理员
const onlineAdmins = ref([])

// 即将到期产品
const expiringProducts = ref([])

// 图表引用
const incomeChartRef = ref<HTMLElement>()
const productChartRef = ref<HTMLElement>()

// 图表实例
let incomeChart: echarts.ECharts | null = null
let productChart: echarts.ECharts | null = null

// 获取统计数据
const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/dashboard/stats' })
    stats.value = data
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 获取收入趋势数据
const fetchIncomeData = async () => {
  incomeLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/dashboard/income-trend',
      params: { period: chartPeriod.value }
    })
    initIncomeChart(data)
  } catch (error) {
    console.error('获取收入趋势失败:', error)
  } finally {
    incomeLoading.value = false
  }
}

// 获取产品分布数据
const fetchProductDistribution = async () => {
  productLoading.value = true
  try {
    const data = await request.get({ url: '/api/admin/dashboard/product-distribution' })
    initProductChart(data)
  } catch (error) {
    console.error('获取产品分布失败:', error)
  } finally {
    productLoading.value = false
  }
}

// 获取最近订单
const fetchRecentOrders = async () => {
  ordersLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/orders',
      params: { page: 1, page_size: 5, sort: 'created_at', order: 'desc' }
    })
    recentOrders.value = data.list || []
  } catch (error) {
    console.error('获取最近订单失败:', error)
  } finally {
    ordersLoading.value = false
  }
}

// 获取最近工单
const fetchRecentTickets = async () => {
  ticketsLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/tickets',
      params: { page: 1, page_size: 5, sort: 'created_at', order: 'desc' }
    })
    recentTickets.value = data.list || []
  } catch (error) {
    console.error('获取最近工单失败:', error)
  } finally {
    ticketsLoading.value = false
  }
}

// 获取最近日志
const fetchRecentLogs = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/system-logs',
      params: { page: 1, page_size: 5 }
    })
    recentLogs.value = data.list || []
  } catch (error) {
    console.error('获取最近日志失败:', error)
  }
}

// 获取在线管理员
const fetchOnlineAdmins = async () => {
  try {
    const data = await request.get({ url: '/api/admin/online-admins' })
    onlineAdmins.value = data || []
  } catch (error) {
    console.error('获取在线管理员失败:', error)
  }
}

// 获取即将到期产品
const fetchExpiringProducts = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/products/expiring',
      params: { days: 7, limit: 5 }
    })
    expiringProducts.value = data || []
  } catch (error) {
    console.error('获取即将到期产品失败:', error)
  }
}

// 订单状态文本
const getOrderStatusText = (status: number) => {
  const map: Record<number, string> = { 0: '待付款', 1: '待开通', 2: '开通中', 3: '已完成', 4: '已取消', 5: '已退款' }
  return map[status] || '未知'
}

// 订单状态类型
const getOrderStatusType = (status: number) => {
  const map: Record<number, string> = { 0: 'warning', 1: 'primary', 2: 'primary', 3: 'success', 4: 'info', 5: 'danger' }
  return map[status] || 'info'
}

// 优先级文本
const getPriorityText = (priority: number) => {
  const map: Record<number, string> = { 1: '低', 2: '普通', 3: '高', 4: '紧急' }
  return map[priority] || '未知'
}

// 优先级类型
const getPriorityType = (priority: number) => {
  const map: Record<number, string> = { 1: 'info', 2: 'primary', 3: 'warning', 4: 'danger' }
  return map[priority] || 'info'
}

// 工单状态文本
const getTicketStatusText = (status: number) => {
  const map: Record<number, string> = { 0: '待处理', 1: '已回复', 2: '处理中', 3: '已关闭' }
  return map[status] || '未知'
}

// 工单状态类型
const getTicketStatusType = (status: number) => {
  const map: Record<number, string> = { 0: 'warning', 1: 'success', 2: 'primary', 3: 'info' }
  return map[status] || 'info'
}

// 初始化收入趋势图表
const initIncomeChart = (chartData: { dates: string[]; values: number[] }) => {
  if (!incomeChartRef.value) return
  if (!incomeChart) incomeChart = echarts.init(incomeChartRef.value)

  const option = {
    tooltip: { trigger: 'axis', formatter: '{b}<br/>收入: ¥{c}' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'category',
      data: chartData.dates,
      axisLine: { lineStyle: { color: '#E5E6EB' } },
      axisLabel: { color: '#86909C' }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: '#86909C', formatter: (value: number) => `¥${(value / 10000).toFixed(1)}万` },
      splitLine: { lineStyle: { type: 'dashed', color: '#E5E6EB' } }
    },
    series: [{
      name: '收入',
      type: 'line',
      smooth: true,
      data: chartData.values,
      itemStyle: { color: 'var(--el-color-primary)' },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(var(--el-color-primary-rgb), 0.3)' },
          { offset: 1, color: 'rgba(var(--el-color-primary-rgb), 0.05)' }
        ])
      }
    }]
  }

  incomeChart.setOption(option, true)
}

// 初始化产品分布图表
const initProductChart = (chartData: { name: string; value: number }[]) => {
  if (!productChartRef.value) return
  if (!productChart) productChart = echarts.init(productChartRef.value)

  const colors = ['#4080FF', '#36D391', '#F59E0B', '#EF4444', '#8B5CF6', '#EC4899', '#14B8A6']

  const option = {
    tooltip: { trigger: 'item', formatter: '{a} <br/>{b}: {c} ({d}%)' },
    legend: { orient: 'vertical', right: 10, top: 'center', textStyle: { color: '#86909C' } },
    series: [{
      name: '产品分布',
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['40%', '50%'],
      avoidLabelOverlap: false,
      itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
      label: { show: false, position: 'center' },
      emphasis: { label: { show: true, fontSize: 14, fontWeight: 'bold' } },
      labelLine: { show: false },
      data: chartData.map((item, index) => ({
        ...item,
        itemStyle: { color: colors[index % colors.length] }
      }))
    }]
  }

  productChart.setOption(option, true)
}

// 窗口大小变化时重绘图表
const handleResize = () => {
  incomeChart?.resize()
  productChart?.resize()
}

onMounted(() => {
  fetchStats()
  fetchIncomeData()
  fetchProductDistribution()
  fetchRecentOrders()
  fetchRecentTickets()
  fetchRecentLogs()
  fetchOnlineAdmins()
  fetchExpiringProducts()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  incomeChart?.dispose()
  productChart?.dispose()
})
</script>

<style scoped lang="scss">
.home-page {
  padding: 20px;
}

.welcome-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 20px;
  background: linear-gradient(135deg, var(--el-color-primary) 0%, var(--el-color-primary-light-3) 100%);
  border-radius: 8px;
  color: #fff;

  .welcome-text {
    h2 {
      margin: 0 0 8px 0;
      font-size: 24px;
      font-weight: 600;
    }

    .version-info {
      margin: 0;
      opacity: 0.8;
      font-size: 14px;
    }
  }

  .welcome-date {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    opacity: 0.9;
  }
}

.stat-row {
  margin-bottom: 16px;
}

.stat-card {
  cursor: pointer;
  transition: all 0.3s;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  .stat-card-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .stat-info {
    flex: 1;
  }

  .stat-label {
    font-size: 14px;
    color: #86909C;
    margin-bottom: 8px;
  }

  .stat-value {
    font-size: 28px;
    font-weight: 600;
    color: #1D2129;
    margin-bottom: 8px;
  }

  .stat-footer {
    font-size: 13px;
    color: #86909C;
  }

  .stat-link {
    color: var(--el-color-primary);
    cursor: pointer;
  }

  .stat-icon {
    width: 64px;
    height: 64px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
  }

  .order-icon { background: linear-gradient(135deg, #4080FF, #6EA8FF); }
  .ticket-icon { background: linear-gradient(135deg, #F59E0B, #FBBF24); }
  .income-icon { background: linear-gradient(135deg, #36D391, #6EE7B7); }
  .client-icon { background: linear-gradient(135deg, #8B5CF6, #A78BFA); }
}

.chart-row, .list-row, .info-row {
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  height: 300px;
}

.log-list {
  max-height: 300px;
  overflow-y: auto;
}

.log-item {
  padding: 8px 0;
  border-bottom: 1px solid #F2F3F5;
  display: flex;
  gap: 12px;

  &:last-child {
    border-bottom: none;
  }

  .log-time {
    font-size: 12px;
    color: #86909C;
    white-space: nowrap;
  }

  .log-content {
    font-size: 13px;
    color: #4E5969;
  }
}

.admin-list {
  max-height: 300px;
  overflow-y: auto;
}

.admin-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #F2F3F5;

  &:last-child {
    border-bottom: none;
  }

  .admin-info {
    flex: 1;
  }

  .admin-name {
    font-size: 14px;
    color: #1D2129;
  }

  .admin-time {
    font-size: 12px;
    color: #86909C;
  }
}

.expire-list {
  max-height: 300px;
  overflow-y: auto;
}

.expire-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #F2F3F5;

  &:last-child {
    border-bottom: none;
  }

  .expire-info {
    flex: 1;
  }

  .expire-name {
    font-size: 14px;
    color: #1D2129;
  }

  .expire-client {
    font-size: 12px;
    color: #86909C;
  }
}

.text-green { color: #36D391; }
.text-red { color: #EF4444; }
</style>
