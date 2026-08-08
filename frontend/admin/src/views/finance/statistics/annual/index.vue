<template>
  <div class="annual-statistics-page">
    <!-- 筛选条件 -->
    <el-card shadow="never" class="filter-card">
      <el-form inline>
        <el-form-item label="年份">
          <el-select v-model="selectedYear" @change="fetchData">
            <el-option v-for="year in yearOptions" :key="year" :label="`${year}年`" :value="year" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 年度概览 -->
    <el-card shadow="never" class="overview-card">
      <template #header>
        <span>{{ selectedYear }}年 度概览</span>
      </template>
      <el-row :gutter="20">
        <el-col :span="6">
          <div class="stat-item">
            <div class="stat-label">年度总收入</div>
            <div class="stat-value text-green">¥{{ formatMoney(overview.total_income) }}</div>
            <div class="stat-compare" :class="overview.income_growth >= 0 ? 'text-green' : 'text-red'">
              {{ overview.income_growth >= 0 ? '+' : '' }}{{ overview.income_growth?.toFixed(1) || 0 }}% 较上年
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item">
            <div class="stat-label">年度总订单</div>
            <div class="stat-value">{{ overview.total_orders || 0 }}</div>
            <div class="stat-compare" :class="overview.orders_growth >= 0 ? 'text-green' : 'text-red'">
              {{ overview.orders_growth >= 0 ? '+' : '' }}{{ overview.orders_growth?.toFixed(1) || 0 }}% 较上年
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item">
            <div class="stat-label">新增客户</div>
            <div class="stat-value">{{ overview.new_clients || 0 }}</div>
            <div class="stat-compare" :class="overview.clients_growth >= 0 ? 'text-green' : 'text-red'">
              {{ overview.clients_growth >= 0 ? '+' : '' }}{{ overview.clients_growth?.toFixed(1) || 0 }}% 较上年
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item">
            <div class="stat-label">新增工单</div>
            <div class="stat-value">{{ overview.total_tickets || 0 }}</div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 月度趋势图表 -->
    <el-card shadow="never" class="chart-card">
      <template #header>
        <span>月度趋势</span>
      </template>
      <div ref="chartRef" class="chart-container" v-loading="chartLoading"></div>
    </el-card>

    <!-- 月度明细表格 -->
    <el-card shadow="never" class="table-card">
      <template #header>
        <span>月度明细</span>
      </template>
      <el-table :data="monthlyData" border stripe show-summary :summary-method="getSummaries">
        <el-table-column prop="month" label="月份" width="100" align="center" />
        <el-table-column prop="income" label="收入" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.income) }}</template>
        </el-table-column>
        <el-table-column prop="expense" label="支出" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.expense) }}</template>
        </el-table-column>
        <el-table-column prop="refund" label="退款" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.refund) }}</template>
        </el-table-column>
        <el-table-column prop="orders" label="订单数" width="100" align="center" />
        <el-table-column prop="new_clients" label="新增客户" width="100" align="center" />
        <el-table-column prop="tickets" label="工单数" width="100" align="center" />
        <el-table-column prop="net_income" label="净收入" width="120" align="right">
          <template #default="{ row }">
            <span :class="row.net_income >= 0 ? 'text-green' : 'text-red'">
              ¥{{ formatMoney(row.net_income) }}
            </span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import request from '@/utils/http'

const selectedYear = ref(new Date().getFullYear())
const chartLoading = ref(false)
const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

// 年份选项
const yearOptions = Array.from({ length: 5 }, (_, i) => new Date().getFullYear() - i)

// 概览数据
const overview = ref<any>({})

// 月度数据
const monthlyData = ref<any[]>([])

// 格式化金额
const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// 获取数据
const fetchData = async () => {
  chartLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/statistics/annual',
      params: { year: selectedYear.value }
    })
    overview.value = data?.overview || {}
    monthlyData.value = data?.monthly || []
    initChart(data?.monthly || [])
  } catch (error) {
    console.error('获取年度统计失败:', error)
  } finally {
    chartLoading.value = false
  }
}

// 初始化图表
const initChart = (data: any[]) => {
  if (!chartRef.value) return
  if (!chart) chart = echarts.init(chartRef.value)

  const months = data.map((d) => d.month)
  const incomeData = data.map((d) => d.income || 0)
  const expenseData = data.map((d) => d.expense || 0)
  const ordersData = data.map((d) => d.orders || 0)

  chart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['收入', '支出', '订单数'] },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: months },
    yAxis: [
      { type: 'value', name: '金额(元)', axisLabel: { formatter: '¥{value}' } },
      { type: 'value', name: '订单数', position: 'right' }
    ],
    series: [
      { name: '收入', type: 'bar', data: incomeData, itemStyle: { color: '#36D391' } },
      { name: '支出', type: 'bar', data: expenseData, itemStyle: { color: '#EF4444' } },
      { name: '订单数', type: 'line', yAxisIndex: 1, data: ordersData, itemStyle: { color: '#4080FF' } }
    ]
  }, true)
}

// 合计方法
const getSummaries = (param: any) => {
  const { columns, data } = param
  const sums: string[] = []
  columns.forEach((column: any, index: number) => {
    if (index === 0) { sums[index] = '合计'; return }
    if (['income', 'expense', 'refund', 'net_income'].includes(column.property)) {
      const values = data.map((item: any) => Number(item[column.property] || 0))
      sums[index] = `¥${formatMoney(values.reduce((a: number, b: number) => a + b, 0))}`
    } else if (['orders', 'new_clients', 'tickets'].includes(column.property)) {
      const values = data.map((item: any) => Number(item[column.property] || 0))
      sums[index] = String(values.reduce((a: number, b: number) => a + b, 0))
    } else {
      sums[index] = ''
    }
  })
  return sums
}

const handleResize = () => chart?.resize()

onMounted(() => {
  fetchData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chart?.dispose()
})
</script>

<style scoped lang="scss">
.annual-statistics-page {
  padding: 16px;
}

.filter-card, .overview-card, .chart-card, .table-card {
  margin-bottom: 16px;
}

.stat-item {
  text-align: center;
  padding: 10px 0;
}

.stat-label {
  font-size: 14px;
  color: #86909C;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
}

.stat-compare {
  font-size: 13px;
  margin-top: 4px;
}

.chart-container {
  height: 400px;
}

.text-green { color: #36D391; }
.text-red { color: #EF4444; }
</style>
