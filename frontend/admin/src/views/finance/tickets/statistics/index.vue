<template>
  <div class="ticket-statistics-page">
    <!-- 筛选条件 -->
    <el-card shadow="never" class="filter-card">
      <el-form inline>
        <el-form-item :label="$t('ticketStatistic.timeRange')">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            :range-separator="$t('common.to')"
            :start-placeholder="$t('common.startDate')"
            :end-placeholder="$t('common.endDate')"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item :label="$t('ticketStatistic.department')">
          <el-select v-model="selectedDept" :placeholder="$t('common.all')" clearable style="width: 150px">
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">
            <el-icon><Search /></el-icon>
            {{ $t('common.search') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="16" class="stat-row">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-label">{{ $t('ticketStatistic.totalTickets') }}</div>
          <div class="stat-value">{{ stats.total || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-label">{{ $t('ticketStatistic.pending') }}</div>
          <div class="stat-value text-orange">{{ stats.open || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-label">{{ $t('ticketStatistic.avgResponseTime') }}</div>
          <div class="stat-value">{{ stats.avg_response_time || '-' }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-label">{{ $t('ticketStatistic.satisfaction') }}</div>
          <div class="stat-value text-green">{{ stats.satisfaction_rate || '-' }}%</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表 -->
    <el-row :gutter="16" class="chart-row">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header><span>{{ $t('ticketStatistic.ticketTrend') }}</span></template>
          <div ref="trendChartRef" class="chart-container" v-loading="chartLoading"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <template #header><span>{{ $t('ticketStatistic.deptDistribution') }}</span></template>
          <div ref="deptChartRef" class="chart-container" v-loading="chartLoading"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 管理员绩效 -->
    <el-card shadow="never" class="table-card">
      <template #header><span>{{ $t('ticketStatistic.adminPerformance') }}</span></template>
      <el-table :data="adminPerformance" border stripe>
        <el-table-column prop="admin_name" :label="$t('ticketStatistic.adminName')" width="120" />
        <el-table-column prop="total_tickets" :label="$t('ticketStatistic.processedTickets')" width="120" align="center" />
        <el-table-column prop="avg_response_time" :label="$t('ticketStatistic.avgResponseTime')" width="150" align="center" />
        <el-table-column prop="avg_resolve_time" :label="$t('ticketStatistic.avgResolveTime')" width="150" align="center" />
        <el-table-column prop="satisfaction_rate" :label="$t('ticketStatistic.satisfaction')" width="100" align="center">
          <template #default="{ row }">
            <span class="text-green">{{ row.satisfaction_rate }}%</span>
          </template>
        </el-table-column>
        <el-table-column prop="open_tickets" :label="$t('ticketStatistic.pending')" width="100" align="center" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { $t } from '@/locales'
import request from '@/utils/http'

const dateRange = ref<[Date, Date] | null>(null)
const selectedDept = ref<number | null>(null)
const departments = ref<any[]>([])
const chartLoading = ref(false)
const stats = ref<any>({})
const adminPerformance = ref<any[]>([])

const trendChartRef = ref<HTMLElement>()
const deptChartRef = ref<HTMLElement>()
let trendChart: echarts.ECharts | null = null
let deptChart: echarts.ECharts | null = null

const fetchData = async () => {
  chartLoading.value = true
  try {
    const params: any = {}
    if (dateRange.value) {
      params.start_date = dateRange.value[0].toISOString().split('T')[0]
      params.end_date = dateRange.value[1].toISOString().split('T')[0]
    }
    if (selectedDept.value) params.department_id = selectedDept.value

    const data = await request.get({ url: '/api/admin/tickets/statistics', params })
    stats.value = data?.stats || {}
    adminPerformance.value = data?.admin_performance || []
    initTrendChart(data?.trend || [])
    initDeptChart(data?.dept_distribution || [])
  } catch (error) {
    console.error($t('ticketStatistic.fetchFailed'), error)
  } finally {
    chartLoading.value = false
  }
}

const fetchDepartments = async () => {
  try {
    const data = await request.get({ url: '/api/admin/ticket-departments' })
    departments.value = data || []
  } catch (error) {
    console.error($t('ticketStatistic.fetchDeptFailed'), error)
  }
}

const initTrendChart = (data: any[]) => {
  if (!trendChartRef.value) return
  if (!trendChart) trendChart = echarts.init(trendChartRef.value)
  trendChart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: data.map((d) => d.date) },
    yAxis: { type: 'value' },
    series: [
      { name: $t('ticketStatistic.newTickets'), type: 'line', smooth: true, data: data.map((d) => d.new_tickets), itemStyle: { color: '#4080FF' } },
      { name: $t('ticketStatistic.closedTickets'), type: 'line', smooth: true, data: data.map((d) => d.closed_tickets), itemStyle: { color: '#36D391' } }
    ]
  }, true)
}

const initDeptChart = (data: any[]) => {
  if (!deptChartRef.value) return
  if (!deptChart) deptChart = echarts.init(deptChartRef.value)
  const colors = ['#4080FF', '#36D391', '#F59E0B', '#EF4444', '#8B5CF6']
  deptChart.setOption({
    tooltip: { trigger: 'item', formatter: '{a} <br/>{b}: {c} ({d}%)' },
    legend: { orient: 'vertical', right: 10, top: 'center' },
    series: [{
      name: $t('ticketStatistic.deptDistribution'),
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['40%', '50%'],
      data: data.map((d, i) => ({ ...d, itemStyle: { color: colors[i % colors.length] } }))
    }]
  }, true)
}

const handleResize = () => {
  trendChart?.resize()
  deptChart?.resize()
}

onMounted(() => {
  fetchData()
  fetchDepartments()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  deptChart?.dispose()
})
</script>

<style scoped lang="scss">
.ticket-statistics-page {
  padding: 16px;
}

.filter-card, .table-card {
  margin-bottom: 16px;
}

.stat-row {
  margin-bottom: 16px;
}

.stat-card {
  text-align: center;
  padding: 20px;
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

.chart-row {
  margin-bottom: 16px;
}

.chart-container {
  height: 300px;
}

.text-orange { color: #F59E0B; }
.text-green { color: #36D391; }
</style>
