<template>
  <div class="finance-dashboard">
    <!-- 全局搜索 -->
    <el-card shadow="hover" class="search-card mb-4">
      <el-input
        v-model="globalSearch"
        placeholder="全局搜索：输入客户名、邮箱、产品、订单号、工单号..."
        size="large"
        clearable
        @keyup.enter="handleGlobalSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
        <template #append>
          <el-button @click="handleGlobalSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
        </template>
      </el-input>
      <div class="search-tags">
        <el-tag
          v-for="tag in searchTags"
          :key="tag.label"
          :type="tag.type"
          class="search-tag"
          @click="quickSearch(tag.query)"
        >
          {{ tag.label }}
        </el-tag>
      </div>
    </el-card>

    <!-- 快捷操作 -->
    <el-card shadow="hover" class="mb-4">
      <template #header>
        <span>快捷操作</span>
      </template>
      <el-row :gutter="16">
        <el-col :span="4" v-for="action in quickActions" :key="action.label">
          <el-button
            class="quick-action-btn"
            :icon="action.icon"
            @click="$router.push(action.route)"
          >
            {{ action.label }}
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <el-row :gutter="20" class="mb-4">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>今日收入</span>
              <el-icon><Money /></el-icon>
            </div>
          </template>
          <div class="stat-value">¥{{ stats.today_income?.toLocaleString() || '0' }}</div>
          <div class="stat-footer">
            <span :class="stats.income_change >= 0 ? 'text-green' : 'text-red'">
              {{ stats.income_change >= 0 ? '+' : '' }}{{ stats.income_change?.toFixed(1) || '0' }}%
            </span>
            较昨日
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>新增客户</span>
              <el-icon><User /></el-icon>
            </div>
          </template>
          <div class="stat-value">{{ stats.new_clients || 0 }}</div>
          <div class="stat-footer">
            <span :class="stats.client_change >= 0 ? 'text-green' : 'text-red'">
              {{ stats.client_change >= 0 ? '+' : '' }}{{ stats.client_change?.toFixed(1) || '0' }}%
            </span>
            较昨日
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>待处理工单</span>
              <el-icon><ChatDotRound /></el-icon>
            </div>
          </template>
          <div class="stat-value">{{ stats.pending_tickets || 0 }}</div>
          <div class="stat-footer">
            <span class="text-orange">{{ stats.urgent_tickets || 0 }} 个紧急</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>待处理订单</span>
              <el-icon><Document /></el-icon>
            </div>
          </template>
          <div class="stat-value">{{ stats.pending_orders || 0 }}</div>
          <div class="stat-footer">
            <span class="text-blue">待审核: {{ stats.review_orders || 0 }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mb-4">
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

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>最近订单</span>
              <el-button type="primary" link @click="$router.push('/finance/orders/list')">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentOrders" style="width: 100%" v-loading="ordersLoading">
            <el-table-column prop="order_no" label="订单号" width="150" />
            <el-table-column prop="client_name" label="客户" width="100" />
            <el-table-column prop="product_name" label="产品" />
            <el-table-column prop="amount" label="金额" width="100">
              <template #default="{ row }">
                ¥{{ row.amount?.toLocaleString() || '0' }}
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
              <el-button type="primary" link @click="$router.push('/finance/tickets/list')">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentTickets" style="width: 100%" v-loading="ticketsLoading">
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Money, User, ChatDotRound, Document, Search, Plus, Tickets, ShoppingCart, UserFilled, Setting } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import request from '@/utils/http'

const router = useRouter()
const chartPeriod = ref('week')

// 全局搜索
const globalSearch = ref('')
const searchTags: { label: string; type: any; query: string }[] = [
  { label: '客户管理', type: 'primary', query: 'clients' },
  { label: '订单列表', type: 'success', query: 'orders' },
  { label: '工单列表', type: 'warning', query: 'tickets' },
  { label: '产品管理', type: 'info', query: 'products' },
]

// 快捷操作
const quickActions = [
  { label: '新增客户', icon: Plus, route: '/finance/clients/list' },
  { label: '创建订单', icon: ShoppingCart, route: '/finance/orders/list' },
  { label: '处理工单', icon: Tickets, route: '/finance/tickets/list' },
  { label: '客户列表', icon: UserFilled, route: '/finance/clients/list' },
  { label: '产品管理', icon: Setting, route: '/finance/products/list' },
  { label: '系统设置', icon: Setting, route: '/finance/system/general' },
]

const handleGlobalSearch = () => {
  if (!globalSearch.value.trim()) return
  router.push({
    path: '/finance/clients/list',
    query: { search: globalSearch.value }
  })
}

const quickSearch = (query: string) => {
  const routeMap: Record<string, string> = {
    clients: '/finance/clients/list',
    orders: '/finance/orders/list',
    tickets: '/finance/tickets/list',
    products: '/finance/products/list',
  }
  router.push(routeMap[query] || '/finance/dashboard')
}

// 加载状态
const incomeLoading = ref(false)
const productLoading = ref(false)
const ordersLoading = ref(false)
const ticketsLoading = ref(false)

// 统计数据
const stats = ref({
  today_income: 0,
  income_change: 0,
  new_clients: 0,
  client_change: 0,
  pending_tickets: 0,
  urgent_tickets: 0,
  pending_orders: 0,
  review_orders: 0
})

// 最近订单
const recentOrders = ref([])

// 最近工单
const recentTickets = ref([])

// 图表引用
const incomeChartRef = ref<HTMLElement>()
const productChartRef = ref<HTMLElement>()

// 图表实例
let incomeChart: echarts.ECharts | null = null
let productChart: echarts.ECharts | null = null

// 获取统计数据
const fetchStats = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/dashboard/stats'
    })
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
    const data = await request.get({
      url: '/api/admin/dashboard/product-distribution'
    })
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

// 订单状态文本
const getOrderStatusText = (status: number) => {
  const map: Record<number, string> = {
    0: '待付款',
    1: '待开通',
    2: '开通中',
    3: '已完成',
    4: '已取消',
    5: '已退款'
  }
  return map[status] || '未知'
}

// 订单状态类型
const getOrderStatusType = (status: number) => {
  const map: Record<number, any> = {
    0: 'warning',
    1: 'primary',
    2: 'primary',
    3: 'success',
    4: 'info',
    5: 'danger'
  }
  return map[status] || 'info'
}

// 优先级文本
const getPriorityText = (priority: number) => {
  const map: Record<number, string> = {
    1: '低',
    2: '普通',
    3: '高',
    4: '紧急'
  }
  return map[priority] || '未知'
}

// 优先级类型
const getPriorityType = (priority: number) => {
  const map: Record<number, any> = {
    1: 'info',
    2: 'primary',
    3: 'warning',
    4: 'danger'
  }
  return map[priority] || 'info'
}

// 工单状态文本
const getTicketStatusText = (status: number) => {
  const map: Record<number, string> = {
    0: '待处理',
    1: '已回复',
    2: '处理中',
    3: '已关闭'
  }
  return map[status] || '未知'
}

// 工单状态类型
const getTicketStatusType = (status: number) => {
  const map: Record<number, any> = {
    0: 'warning',
    1: 'success',
    2: 'primary',
    3: 'info'
  }
  return map[status] || 'info'
}

// 初始化收入趋势图表
const initIncomeChart = (chartData: { dates: string[]; values: number[] }) => {
  if (!incomeChartRef.value) return

  if (!incomeChart) {
    incomeChart = echarts.init(incomeChartRef.value)
  }

  const option = {
    tooltip: {
      trigger: 'axis',
      formatter: '{b}<br/>收入: ¥{c}'
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: chartData.dates,
      axisLine: {
        lineStyle: {
          color: '#E5E6EB'
        }
      },
      axisLabel: {
        color: '#86909C'
      }
    },
    yAxis: {
      type: 'value',
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: '#86909C',
        formatter: (value: number) => `¥${(value / 10000).toFixed(1)}万`
      },
      splitLine: {
        lineStyle: {
          type: 'dashed',
          color: '#E5E6EB'
        }
      }
    },
    series: [
      {
        name: '收入',
        type: 'line',
        smooth: true,
        data: chartData.values,
        itemStyle: {
          color: '#4080FF'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(64,128,255,0.3)' },
            { offset: 1, color: 'rgba(64,128,255,0.05)' }
          ])
        }
      }
    ]
  }

  incomeChart.setOption(option, true)
}

// 初始化产品分布图表
const initProductChart = (chartData: { name: string; value: number }[]) => {
  if (!productChartRef.value) return

  if (!productChart) {
    productChart = echarts.init(productChartRef.value)
  }

  const colors = ['#4080FF', '#36D391', '#F59E0B', '#EF4444', '#8B5CF6', '#EC4899', '#14B8A6']

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center',
      textStyle: {
        color: '#86909C'
      }
    },
    series: [
      {
        name: '产品分布',
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['40%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 14,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: chartData.map((item, index) => ({
          ...item,
          itemStyle: { color: colors[index % colors.length] }
        }))
      }
    ]
  }

  productChart.setOption(option, true)
}

// 窗口大小变化时重绘图表
const handleResize = () => {
  incomeChart?.resize()
  productChart?.resize()
}

onMounted(() => {
  // 加载所有数据
  fetchStats()
  fetchIncomeData()
  fetchProductDistribution()
  fetchRecentOrders()
  fetchRecentTickets()

  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  incomeChart?.dispose()
  productChart?.dispose()
})
</script>

<style scoped lang="scss">
.finance-dashboard {
  padding: 20px;
}

.search-card {
  .search-tags {
    margin-top: 12px;
    display: flex;
    gap: 8px;

    .search-tag {
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        transform: translateY(-2px);
        opacity: 0.8;
      }
    }
  }
}

.quick-action-btn {
  width: 100%;
  height: 60px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 4px;
  font-size: 13px;
  border-radius: 8px;
  transition: all 0.2s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
}

.stat-card {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    color: #86909C;
  }

  .stat-value {
    font-size: 28px;
    font-weight: 600;
    color: #1D2129;
    margin: 10px 0;
  }

  .stat-footer {
    font-size: 13px;
    color: #86909C;
  }
}

.text-green {
  color: #36D391;
}

.text-red {
  color: #EF4444;
}

.text-orange {
  color: #F59E0B;
}

.text-blue {
  color: #4080FF;
}

.chart-container {
  height: 300px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
