<template>
  <div class="ticket-statistics-page">
    <!-- 筛选条件 -->
    <el-card shadow="never" class="filter-card">
      <el-form inline>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item label="部门">
          <el-select v-model="selectedDept" placeholder="全部" clearable style="width: 150px">
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="16" class="stat-row">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-label">总工单数</div>
          <div class="stat-value">{{ stats.total || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-label">待处理</div>
          <div class="stat-value text-orange">{{ stats.open || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-label">平均响应时间</div>
          <div class="stat-value">{{ stats.avg_response_time || '-' }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-label">满意度</div>
          <div class="stat-value text-green">{{ stats.satisfaction_rate || '-' }}%</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表 -->
    <el-row :gutter="16" class="chart-row">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header><span>工单趋势</span></template>
          <div ref="trendChartRef" class="chart-container" v-loading="chartLoading"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <template #header><span>部门分布</span></template>
          <div ref="deptChartRef" class="chart-container" v-loading="chartLoading"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 管理员绩效 -->
    <el-card shadow="never" class="table-card">
      <template #header><span>管理员绩效</span></template>
      <el-table :data="adminPerformance" border stripe>
        <el-table-column prop="admin_name" label="管理员" width="120" />
        <el-table-column prop="total_tickets" label="处理工单数" width="120" align="center" />
        <el-table-column prop="avg_response_time" label="平均响应时间" width="150" align="center" />
        <el-table-column prop="avg_resolve_time" label="平均解决时间" width="150" align="center" />
        <el-table-column prop="satisfaction_rate" label="满意度" width="100" align="center">
          <template #default="{ row }">
            <span class="text-green">{{ row.satisfaction_rate }}%</span>
          </template>
        </el-table-column>
        <el-table-column prop="open_tickets" label="待处理" width="100" align="center" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
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

// 获取数据
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
    console.error('获取统计数据失败:', error)
  } finally {
    chartLoading.value = false
  }
}

// 获取部门
const fetchDepartments = async () => {
  try {
    const data = await request.get({ url: '/api/admin/ticket-departments' })
    departments.value = data || []
  } catch (error) {
    console.error('获取部门失败:', error)
  }
}

// 初始化趋势图
const initTrendChart = (data: any[]) => {
  if (!trendChartRef.value) return
  if (!trendChart) trendChart = echarts.init(trendChartRef.value)
  trendChart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: data.map((d) => d.date) },
    yAxis: { type: 'value' },
    series: [
      { name: '新建工单', type: 'line', smooth: true, data: data.map((d) => d.new_tickets), itemStyle: { color: '#4080FF' } },
      { name: '已关闭', type: 'line', smooth: true, data: data.map((d) => d.closed_tickets), itemStyle: { color: '#36D391' } }
    ]
  }, true)
}

// 初始化部门分布图
const initDeptChart = (data: any[]) => {
  if (!deptChartRef.value) return
  if (!deptChart) deptChart = echarts.init(deptChartRef.value)
  const colors = ['#4080FF', '#36D391', '#F59E0B', '#EF4444', '#8B5CF6']
  deptChart.setOption({
    tooltip: { trigger: 'item', formatter: '{a} <br/>{b}: {c} ({d}%)' },
    legend: { orient: 'vertical', right: 10, top: 'center' },
    series: [{
      name: '部门分布',
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
