<template>
  <div class="reports-page">
    <!-- 顶部统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-icon income">
              <el-icon><Money /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">本月收入</div>
              <div class="stat-value">¥{{ formatAmount(summary.month_income) }}</div>
              <div class="stat-trend" :class="summary.income_trend >= 0 ? 'up' : 'down'">
                {{ summary.income_trend >= 0 ? '+' : '' }}{{ summary.income_trend }}%
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-icon clients">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">新增客户</div>
              <div class="stat-value">{{ summary.month_clients || 0 }}</div>
              <div class="stat-trend" :class="summary.clients_trend >= 0 ? 'up' : 'down'">
                {{ summary.clients_trend >= 0 ? '+' : '' }}{{ summary.clients_trend }}%
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-icon orders">
              <el-icon><ShoppingCart /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">本月订单</div>
              <div class="stat-value">{{ summary.month_orders || 0 }}</div>
              <div class="stat-trend" :class="summary.orders_trend >= 0 ? 'up' : 'down'">
                {{ summary.orders_trend >= 0 ? '+' : '' }}{{ summary.orders_trend }}%
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-icon tickets">
              <el-icon><Tickets /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">待处理工单</div>
              <div class="stat-value">{{ summary.pending_tickets || 0 }}</div>
              <div class="stat-desc">较昨日 {{ summary.tickets_diff >= 0 ? '+' : '' }}{{ summary.tickets_diff || 0 }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 收入报表 -->
    <el-card shadow="never" class="chart-card">
      <template #header>
        <div class="card-header">
          <span>收入报表</span>
          <div>
            <el-radio-group v-model="incomeRange" size="small" @change="fetchIncomeReport">
              <el-radio-button label="week">近7天</el-radio-button>
              <el-radio-button label="month">近30天</el-radio-button>
              <el-radio-button label="year">近12月</el-radio-button>
            </el-radio-group>
          </div>
        </div>
      </template>
      <div ref="incomeChartRef" class="chart-container" v-loading="incomeLoading"></div>
    </el-card>

    <!-- 客户增长 & 产品销售 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="12">
        <el-card shadow="never" class="chart-card">
          <template #header>
            <div class="card-header">
              <span>客户增长报表</span>
              <el-radio-group v-model="clientRange" size="small" @change="fetchClientReport">
                <el-radio-button label="month">月</el-radio-button>
                <el-radio-button label="year">年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="clientChartRef" class="chart-container" v-loading="clientLoading"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never" class="chart-card">
          <template #header>
            <div class="card-header">
              <span>产品销售报表</span>
              <el-radio-group v-model="productRange" size="small" @change="fetchProductReport">
                <el-radio-button label="month">月</el-radio-button>
                <el-radio-button label="year">年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="productChartRef" class="chart-container" v-loading="productLoading"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 工单统计 -->
    <el-card shadow="never" class="chart-card">
      <template #header>
        <div class="card-header">
          <span>工单统计</span>
          <el-radio-group v-model="ticketRange" size="small" @change="fetchTicketReport">
            <el-radio-button label="week">近7天</el-radio-button>
            <el-radio-button label="month">近30天</el-radio-button>
          </el-radio-group>
        </div>
      </template>
      <div ref="ticketChartRef" class="chart-container" v-loading="ticketLoading"></div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { Money, User, ShoppingCart, Tickets } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'
import * as echarts from 'echarts'

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

// 概览数据
const summary = reactive({
  month_income: 0,
  income_trend: 0,
  month_clients: 0,
  clients_trend: 0,
  month_orders: 0,
  orders_trend: 0,
  pending_tickets: 0,
  tickets_diff: 0
})

// 时间范围
const incomeRange = ref('month')
const clientRange = ref('month')
const productRange = ref('month')
const ticketRange = ref('week')

// 加载状态
const incomeLoading = ref(false)
const clientLoading = ref(false)
const productLoading = ref(false)
const ticketLoading = ref(false)

// 图表DOM引用
const incomeChartRef = ref<HTMLDivElement>()
const clientChartRef = ref<HTMLDivElement>()
const productChartRef = ref<HTMLDivElement>()
const ticketChartRef = ref<HTMLDivElement>()

// 图表实例
let incomeChart: echarts.ECharts | null = null
let clientChart: echarts.ECharts | null = null
let productChart: echarts.ECharts | null = null
let ticketChart: echarts.ECharts | null = null

const initChart = (el: HTMLDivElement | undefined): echarts.ECharts | null => {
  if (!el) return null
  return echarts.init(el)
}

// 收入报表
const fetchIncomeReport = async () => {
  incomeLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/reports/revenue',
      params: { range: incomeRange.value }
    })
    if (!incomeChart) incomeChart = initChart(incomeChartRef.value)
    if (!incomeChart) return

    const labels = data.labels || []
    const values = data.values || []

    incomeChart.setOption({
      tooltip: {
        trigger: 'axis',
        axisPointer: { type: 'shadow' },
        formatter: (params: any) => {
          const p = params[0]
          return `${p.name}<br/>收入: ¥${Number(p.value).toLocaleString()}`
        }
      },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: {
        type: 'category',
        data: labels,
        axisLabel: { color: '#666' }
      },
      yAxis: {
        type: 'value',
        axisLabel: {
          color: '#666',
          formatter: (v: number) => v >= 10000 ? (v / 10000) + 'w' : v + ''
        }
      },
      series: [
        {
          name: '收入',
          type: 'bar',
          data: values,
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#409eff' },
              { offset: 1, color: '#79bbff' }
            ]),
            borderRadius: [4, 4, 0, 0]
          },
          barMaxWidth: 40
        }
      ]
    })
  } catch (error) {
    console.error('获取收入报表失败:', error)
  } finally {
    incomeLoading.value = false
  }
}

// 客户增长报表
const fetchClientReport = async () => {
  clientLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/reports/users',
      params: { range: clientRange.value }
    })
    if (!clientChart) clientChart = initChart(clientChartRef.value)
    if (!clientChart) return

    const labels = data.labels || []
    const newClients = data.new_clients || []
    const totalClients = data.total_clients || []

    clientChart.setOption({
      tooltip: { trigger: 'axis' },
      legend: { data: ['新增客户', '客户总量'], bottom: 0 },
      grid: { left: '3%', right: '4%', bottom: '12%', containLabel: true },
      xAxis: {
        type: 'category',
        data: labels,
        axisLabel: { color: '#666' }
      },
      yAxis: [
        {
          type: 'value',
          name: '新增',
          axisLabel: { color: '#666' }
        },
        {
          type: 'value',
          name: '总量',
          axisLabel: { color: '#666' }
        }
      ],
      series: [
        {
          name: '新增客户',
          type: 'bar',
          data: newClients,
          itemStyle: {
            color: '#67c23a',
            borderRadius: [4, 4, 0, 0]
          },
          barMaxWidth: 30
        },
        {
          name: '客户总量',
          type: 'line',
          yAxisIndex: 1,
          data: totalClients,
          smooth: true,
          lineStyle: { color: '#e6a23c', width: 2 },
          itemStyle: { color: '#e6a23c' }
        }
      ]
    })
  } catch (error) {
    console.error('获取客户增长报表失败:', error)
  } finally {
    clientLoading.value = false
  }
}

// 产品销售报表
const fetchProductReport = async () => {
  productLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/reports/product-income',
      params: { range: productRange.value }
    })
    if (!productChart) productChart = initChart(productChartRef.value)
    if (!productChart) return

    const items = data.items || []
    const colors = ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#909399', '#9b59b6', '#3498db']

    productChart.setOption({
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        right: '5%',
        top: 'center',
        type: 'scroll'
      },
      series: [
        {
          name: '产品销售',
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['40%', '50%'],
          avoidLabelOverlap: false,
          itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
          label: { show: false },
          emphasis: {
            label: { show: true, fontSize: 14, fontWeight: 'bold' }
          },
          data: items.map((item: any, idx: number) => ({
            ...item,
            itemStyle: { color: colors[idx % colors.length] }
          }))
        }
      ]
    })
  } catch (error) {
    console.error('获取产品销售报表失败:', error)
  } finally {
    productLoading.value = false
  }
}

// 工单统计报表
const fetchTicketReport = async () => {
  ticketLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/reports/tickets',
      params: { range: ticketRange.value }
    })
    if (!ticketChart) ticketChart = initChart(ticketChartRef.value)
    if (!ticketChart) return

    const labels = data.labels || []
    const opened = data.opened || []
    const closed = data.closed || []

    ticketChart.setOption({
      tooltip: { trigger: 'axis' },
      legend: { data: ['新建工单', '已关闭工单'], bottom: 0 },
      grid: { left: '3%', right: '4%', bottom: '12%', containLabel: true },
      xAxis: {
        type: 'category',
        data: labels,
        axisLabel: { color: '#666' }
      },
      yAxis: {
        type: 'value',
        axisLabel: { color: '#666' }
      },
      series: [
        {
          name: '新建工单',
          type: 'line',
          data: opened,
          smooth: true,
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(245,108,108,0.3)' },
              { offset: 1, color: 'rgba(245,108,108,0.05)' }
            ])
          },
          lineStyle: { color: '#f56c6c', width: 2 },
          itemStyle: { color: '#f56c6c' }
        },
        {
          name: '已关闭工单',
          type: 'line',
          data: closed,
          smooth: true,
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(103,194,58,0.3)' },
              { offset: 1, color: 'rgba(103,194,58,0.05)' }
            ])
          },
          lineStyle: { color: '#67c23a', width: 2 },
          itemStyle: { color: '#67c23a' }
        }
      ]
    })
  } catch (error) {
    console.error('获取工单统计失败:', error)
  } finally {
    ticketLoading.value = false
  }
}

// 获取概览数据
const fetchSummary = async () => {
  try {
    const data = await request.get({ url: '/api/admin/reports/dashboard' })
    Object.assign(summary, data)
  } catch (error) {
    console.error('获取概览数据失败:', error)
  }
}

// 窗口resize
const handleResize = () => {
  incomeChart?.resize()
  clientChart?.resize()
  productChart?.resize()
  ticketChart?.resize()
}

onMounted(async () => {
  window.addEventListener('resize', handleResize)
  fetchSummary()
  await nextTick()
  fetchIncomeReport()
  fetchClientReport()
  fetchProductReport()
  fetchTicketReport()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  incomeChart?.dispose()
  clientChart?.dispose()
  productChart?.dispose()
  ticketChart?.dispose()
})
</script>

<style scoped lang="scss">
.reports-page {
  padding: 20px;
}

.stat-cards {
  margin-bottom: 20px;
}

.stat-card {
  .stat-item {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .stat-icon {
    width: 56px;
    height: 56px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    color: #fff;

    &.income { background: linear-gradient(135deg, #409eff, #79bbff); }
    &.clients { background: linear-gradient(135deg, #67c23a, #95d475); }
    &.orders { background: linear-gradient(135deg, #e6a23c, #eebe77); }
    &.tickets { background: linear-gradient(135deg, #f56c6c, #fab6b6); }
  }

  .stat-info {
    flex: 1;
  }

  .stat-label {
    font-size: 13px;
    color: #909399;
    margin-bottom: 4px;
  }

  .stat-value {
    font-size: 22px;
    font-weight: 700;
    color: #303133;
    line-height: 1.2;
  }

  .stat-trend {
    font-size: 12px;
    margin-top: 4px;

    &.up { color: #67c23a; }
    &.down { color: #f56c6c; }
  }

  .stat-desc {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
  }
}

.chart-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  width: 100%;
  height: 380px;
}

.chart-row {
  margin-bottom: 0;

  .chart-card {
    height: 100%;
  }
}
</style>
