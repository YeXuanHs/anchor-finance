<template>
  <div class="dashboard-page">
    <!-- Stat Cards -->
    <el-row :gutter="16" class="stat-row">
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card blue">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">今日收入</span>
              <span class="stat-value">¥{{ stats.todayRevenue.toLocaleString() }}</span>
              <span class="stat-trend">
                <el-icon><Top /></el-icon>
                {{ stats.revenueChange }}%
              </span>
            </div>
            <div class="stat-icon">
              <el-icon :size="48"><Wallet /></el-icon>
            </div>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card green">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">新增用户</span>
              <span class="stat-value">{{ stats.newUsers }}</span>
              <span class="stat-trend">
                <el-icon><Top /></el-icon>
                {{ stats.usersChange }}%
              </span>
            </div>
            <div class="stat-icon">
              <el-icon :size="48"><User /></el-icon>
            </div>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card orange">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">待处理工单</span>
              <span class="stat-value">{{ stats.openTickets }}</span>
              <span class="stat-hint">需要处理</span>
            </div>
            <div class="stat-icon">
              <el-icon :size="48"><ChatDotRound /></el-icon>
            </div>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card red">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">服务器数量</span>
              <span class="stat-value">{{ stats.serverCount }}</span>
              <span class="stat-hint">运行中</span>
            </div>
            <div class="stat-icon">
              <el-icon :size="48"><Monitor /></el-icon>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- Charts Row -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :lg="16">
        <el-card class="admin-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>收入趋势（近30天）</span>
              <el-radio-group v-model="chartPeriod" size="small">
                <el-radio-button value="week">本周</el-radio-button>
                <el-radio-button value="month">本月</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container">
            <v-chart :option="revenueChartOption" autoresize style="height: 360px" />
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="8">
        <el-card class="admin-card" shadow="never">
          <template #header>
            <span>订单状态分布</span>
          </template>
          <div class="chart-container">
            <v-chart :option="orderPieOption" autoresize style="height: 360px" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Recent Orders & Tickets -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :lg="16">
        <el-card class="admin-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>最近订单</span>
              <el-button type="primary" link @click="$router.push('/admin/orders')">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentOrders" style="width: 100%" size="small">
            <el-table-column prop="id" label="订单号" width="150" show-overflow-tooltip />
            <el-table-column prop="user" label="用户" width="80" />
            <el-table-column prop="product" label="产品" show-overflow-tooltip />
            <el-table-column prop="amount" label="金额" width="100">
              <template #default="{ row }">
                <span class="amount">¥{{ row.amount }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="statusMap[row.status]?.type" size="small" round>
                  {{ statusMap[row.status]?.label }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="time" label="时间" width="150" />
          </el-table>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="8">
        <el-card class="admin-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>最近工单</span>
              <el-button type="primary" link @click="$router.push('/admin/tickets')">查看全部</el-button>
            </div>
          </template>
          <div class="ticket-list">
            <div v-for="ticket in recentTickets" :key="ticket.id" class="ticket-item">
              <div class="ticket-info">
                <span class="ticket-subject">{{ ticket.subject }}</span>
                <div class="ticket-meta">
                  <span>{{ ticket.user }}</span>
                  <span>{{ ticket.time }}</span>
                </div>
              </div>
              <el-tag :type="priorityType(ticket.priority)" size="small" round>
                {{ ticket.priority }}
              </el-tag>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- System Status -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="24">
        <el-card class="admin-card" shadow="never">
          <template #header>
            <span>系统状态</span>
          </template>
          <el-row :gutter="24">
            <el-col :xs="24" :sm="8">
              <div class="status-item">
                <div class="status-header">
                  <span>CPU 使用率</span>
                  <span class="status-value">45%</span>
                </div>
                <el-progress :percentage="45" :stroke-width="8" color="#409eff" />
              </div>
            </el-col>
            <el-col :xs="24" :sm="8">
              <div class="status-item">
                <div class="status-header">
                  <span>内存使用率</span>
                  <span class="status-value">62%</span>
                </div>
                <el-progress :percentage="62" :stroke-width="8" color="#67c23a" />
              </div>
            </el-col>
            <el-col :xs="24" :sm="8">
              <div class="status-item">
                <div class="status-header">
                  <span>磁盘使用率</span>
                  <span class="status-value">78%</span>
                </div>
                <el-progress :percentage="78" :stroke-width="8" color="#e6a23c" />
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, PieChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
} from 'echarts/components'
import {
  Top,
  Wallet,
  User,
  ChatDotRound,
  Monitor,
} from '@element-plus/icons-vue'

use([CanvasRenderer, LineChart, PieChart, TitleComponent, TooltipComponent, LegendComponent, GridComponent])

definePageMeta({
  layout: 'admin',
})

const chartPeriod = ref('month')

// Stats
const stats = reactive({
  todayRevenue: 12680,
  revenueChange: 12.5,
  newUsers: 34,
  usersChange: 8.3,
  openTickets: 7,
  serverCount: 24,
})

// Status map
const statusMap: Record<string, { label: string; type: string }> = {
  pending: { label: '待支付', type: 'warning' },
  paid: { label: '已支付', type: 'info' },
  active: { label: '已开通', type: 'success' },
  cancelled: { label: '已取消', type: 'danger' },
}

// Revenue Chart
function getLast30Days(): string[] {
  const days: string[] = []
  const now = new Date()
  for (let i = 29; i >= 0; i--) {
    const d = new Date(now)
    d.setDate(d.getDate() - i)
    days.push(`${d.getMonth() + 1}/${d.getDate()}`)
  }
  return days
}

function generateRevenueData(): number[] {
  const data: number[] = []
  for (let i = 0; i < 30; i++) {
    data.push(Math.floor(Math.random() * 8000 + 4000))
  }
  return data
}

const revenueChartOption = computed(() => ({
  tooltip: {
    trigger: 'axis',
    backgroundColor: 'rgba(255,255,255,0.95)',
    borderColor: '#eee',
    textStyle: { color: '#333' },
  },
  grid: { left: '3%', right: '4%', bottom: '3%', top: '8%', containLabel: true },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: getLast30Days(),
    axisLine: { lineStyle: { color: '#e8e8e8' } },
    axisLabel: { color: '#8c8c8c', fontSize: 11 },
    axisTick: { show: false },
  },
  yAxis: {
    type: 'value',
    axisLabel: {
      color: '#8c8c8c',
      formatter: (v: number) => v >= 1000 ? `${(v / 1000).toFixed(0)}k` : `${v}`,
    },
    splitLine: { lineStyle: { color: '#f0f0f0', type: 'dashed' } },
    axisLine: { show: false },
    axisTick: { show: false },
  },
  series: [
    {
      name: '收入',
      type: 'line',
      smooth: true,
      symbol: 'circle',
      symbolSize: 6,
      showSymbol: false,
      lineStyle: { color: '#409eff', width: 3 },
      itemStyle: { color: '#409eff', borderWidth: 2, borderColor: '#fff' },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(64, 158, 255, 0.25)' },
            { offset: 0.5, color: 'rgba(64, 158, 255, 0.08)' },
            { offset: 1, color: 'rgba(64, 158, 255, 0.01)' },
          ],
        },
      },
      data: generateRevenueData(),
    },
  ],
}))

// Pie Chart
const orderPieOption = computed(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
  legend: { bottom: '2%', left: 'center', itemWidth: 10, itemHeight: 10, textStyle: { fontSize: 12 } },
  series: [
    {
      type: 'pie',
      radius: ['45%', '72%'],
      center: ['50%', '45%'],
      avoidLabelOverlap: false,
      itemStyle: { borderRadius: 8, borderColor: '#fff', borderWidth: 3 },
      label: { show: false },
      emphasis: {
        label: { show: true, fontSize: 14, fontWeight: 'bold' },
      },
      labelLine: { show: false },
      data: [
        { value: 340, name: '待支付', itemStyle: { color: '#e6a23c' } },
        { value: 580, name: '已支付', itemStyle: { color: '#409eff' } },
        { value: 1220, name: '已开通', itemStyle: { color: '#67c23a' } },
        { value: 180, name: '已取消', itemStyle: { color: '#909399' } },
      ],
    },
  ],
}))

// Recent Orders
const recentOrders = ref([
  { id: 'AF20260726001', user: '张三', product: '基础版主机', amount: 299, status: 'active', time: '2026-07-26 14:30' },
  { id: 'AF20260726002', user: '李四', product: '高级版主机', amount: 599, status: 'paid', time: '2026-07-26 13:22' },
  { id: 'AF20260726003', user: '王五', product: '企业版主机', amount: 1299, status: 'pending', time: '2026-07-26 11:45' },
  { id: 'AF20260726004', user: '赵六', product: '1核2G云服务器', amount: 89, status: 'active', time: '2026-07-26 10:18' },
  { id: 'AF20260726005', user: '孙七', product: '4核8G云服务器', amount: 399, status: 'cancelled', time: '2026-07-25 16:55' },
])

// Recent Tickets
const recentTickets = ref([
  { id: 1, subject: '服务器无法连接SSH', user: '张三', priority: '紧急', time: '10分钟前' },
  { id: 2, subject: '域名解析未生效', user: '李四', priority: '高', time: '30分钟前' },
  { id: 3, subject: '申请退款基础版主机', user: '王五', priority: '中', time: '1小时前' },
  { id: 4, subject: '控制面板登录异常', user: '赵六', priority: '低', time: '2小时前' },
  { id: 5, subject: 'SSL证书安装咨询', user: '孙七', priority: '低', time: '3小时前' },
])

function priorityType(priority: string): '' | 'success' | 'warning' | 'info' | 'danger' {
  const map: Record<string, '' | 'success' | 'warning' | 'info' | 'danger'> = {
    '紧急': 'danger',
    '高': 'warning',
    '中': 'info',
    '低': '',
  }
  return map[priority] || ''
}
</script>

<style scoped>
.dashboard-page {
  padding: 0;
}

.stat-row {
  margin-bottom: 0;
}

.stat-card {
  border-radius: 12px;
  padding: 24px;
  color: #ffffff;
  transition: all 0.3s;
  cursor: pointer;
  margin-bottom: 16px;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
}

.stat-card.blue {
  background: linear-gradient(135deg, #409eff 0%, #337ecc 100%);
}

.stat-card.green {
  background: linear-gradient(135deg, #67c23a 0%, #529b2e 100%);
}

.stat-card.orange {
  background: linear-gradient(135deg, #e6a23c 0%, #cf9236 100%);
}

.stat-card.red {
  background: linear-gradient(135deg, #f56c6c 0%, #dd6161 100%);
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
  font-size: 14px;
  opacity: 0.85;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  line-height: 1.2;
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  margin-top: 8px;
  opacity: 0.9;
}

.stat-hint {
  font-size: 12px;
  margin-top: 8px;
  opacity: 0.8;
}

.stat-icon {
  opacity: 0.3;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.chart-container {
  width: 100%;
}

.amount {
  font-weight: 600;
  color: #409eff;
}

.ticket-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ticket-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background: #fafafa;
  border-radius: 8px;
  transition: background-color 0.3s;
}

.ticket-item:hover {
  background: #f0f0f0;
}

.ticket-info {
  flex: 1;
  min-width: 0;
}

.ticket-subject {
  font-size: 14px;
  color: #333;
  display: block;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.ticket-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.status-item {
  padding: 16px 0;
}

.status-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.status-header span:first-child {
  font-size: 14px;
  color: #666;
}

.status-value {
  font-weight: 600;
  color: #333;
}
</style>
