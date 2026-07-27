<template>
  <div class="dashboard">
    <!-- Stat Cards -->
    <el-row :gutter="16">
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">今日收入</span>
              <span class="stat-value">¥{{ stats.todayRevenue.toLocaleString() }}</span>
              <span class="stat-change" :class="stats.revenueChange >= 0 ? 'up' : 'down'">
                <el-icon :size="14"><Top v-if="stats.revenueChange >= 0" /><Bottom v-else /></el-icon>
                {{ Math.abs(stats.revenueChange) }}%
              </span>
            </div>
            <div class="stat-icon blue">
              <el-icon :size="28"><Wallet /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">新增用户</span>
              <span class="stat-value">{{ stats.newUsers }}</span>
              <span class="stat-change" :class="stats.usersChange >= 0 ? 'up' : 'down'">
                <el-icon :size="14"><Top v-if="stats.usersChange >= 0" /><Bottom v-else /></el-icon>
                {{ Math.abs(stats.usersChange) }}%
              </span>
            </div>
            <div class="stat-icon green">
              <el-icon :size="28"><User /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">待处理工单</span>
              <span class="stat-value" :class="{ 'text-red': stats.openTickets > 0 }">{{ stats.openTickets }}</span>
              <span class="stat-change neutral" v-if="stats.openTickets > 0">需要处理</span>
            </div>
            <div class="stat-icon" :class="stats.openTickets > 0 ? 'red' : 'green'">
              <el-icon :size="28"><ChatDotSquare /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">服务器数量</span>
              <span class="stat-value">{{ stats.serverCount }}</span>
              <span class="stat-change up">
                <el-icon :size="14"><Top /></el-icon>
                运行正常
              </span>
            </div>
            <div class="stat-icon cyan">
              <el-icon :size="28"><Monitor /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Charts Row -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :lg="16">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>收入趋势（近30天）</span>
            </div>
          </template>
          <v-chart :option="revenueChartOption" autoresize style="height: 360px" />
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="8">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>订单状态分布</span>
            </div>
          </template>
          <v-chart :option="orderPieOption" autoresize style="height: 360px" />
        </el-card>
      </el-col>
    </el-row>

    <!-- Recent Orders -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :lg="16">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>最近订单</span>
              <el-button text type="primary" @click="$router.push('/admin/orders')">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentOrders" size="small" stripe>
            <el-table-column prop="id" label="订单号" width="160" show-overflow-tooltip />
            <el-table-column prop="user" label="用户" width="90" />
            <el-table-column prop="product" label="产品" show-overflow-tooltip />
            <el-table-column prop="amount" label="金额" width="100">
              <template #default="{ row }">
                <span style="font-weight: 600; color: #0056FF">¥{{ row.amount }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="statusMap[row.status]?.type" size="small" round>
                  {{ statusMap[row.status]?.label }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="time" label="时间" width="160" />
          </el-table>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="8">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>最近工单</span>
              <el-button text type="primary" @click="$router.push('/admin/tickets')">查看全部</el-button>
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
  </div>
</template>

<script setup lang="ts">
import { reactive, computed, ref } from 'vue'
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
  Bottom,
  Wallet,
  User,
  ChatDotSquare,
  Monitor,
} from '@element-plus/icons-vue'

use([CanvasRenderer, LineChart, PieChart, TitleComponent, TooltipComponent, LegendComponent, GridComponent])

const stats = reactive({
  todayRevenue: 12680,
  revenueChange: 12.5,
  newUsers: 34,
  usersChange: -2.1,
  openTickets: 7,
  serverCount: 24,
})

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
    formatter: (params: any) => {
      const p = params[0]
      return `<div style="font-weight:600">${p.axisValue}</div><div>收入: <b style="color:#0056FF">¥${p.value.toLocaleString()}</b></div>`
    },
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
      lineStyle: { color: '#0056FF', width: 3 },
      itemStyle: { color: '#0056FF', borderWidth: 2, borderColor: '#fff' },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(0, 86, 255, 0.25)' },
            { offset: 0.5, color: 'rgba(0, 86, 255, 0.08)' },
            { offset: 1, color: 'rgba(0, 86, 255, 0.01)' },
          ],
        },
      },
      emphasis: {
        focus: 'series',
        itemStyle: { borderWidth: 3, shadowBlur: 10, shadowColor: 'rgba(0,86,255,0.3)' },
      },
      data: generateRevenueData(),
    },
  ],
}))

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
        itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0,0,0,0.15)' },
      },
      labelLine: { show: false },
      data: [
        { value: 340, name: '待支付', itemStyle: { color: '#0056FF' } },
        { value: 580, name: '已支付', itemStyle: { color: '#52c41a' } },
        { value: 1220, name: '已开通', itemStyle: { color: '#fa8c16' } },
        { value: 180, name: '已取消', itemStyle: { color: '#d9d9d9' } },
      ],
    },
  ],
}))

const statusMap: Record<string, { label: string; type: string }> = {
  pending: { label: '待支付', type: 'warning' },
  paid: { label: '已支付', type: 'info' },
  active: { label: '已开通', type: 'success' },
  cancelled: { label: '已取消', type: 'info' },
}

const recentOrders = ref([
  { id: 'AF20260726001', user: '张三', product: '基础版主机', amount: 299, status: 'active', time: '2026-07-26 14:30' },
  { id: 'AF20260726002', user: '李四', product: '高级版主机', amount: 599, status: 'paid', time: '2026-07-26 13:22' },
  { id: 'AF20260726003', user: '王五', product: '企业版主机', amount: 1299, status: 'pending', time: '2026-07-26 11:45' },
  { id: 'AF20260726004', user: '赵六', product: '1核2G云服务器', amount: 89, status: 'active', time: '2026-07-26 10:18' },
  { id: 'AF20260726005', user: '孙七', product: '4核8G云服务器', amount: 399, status: 'cancelled', time: '2026-07-25 16:55' },
])

const recentTickets = ref([
  { id: 1, subject: '服务器无法连接SSH', user: '张三', priority: '紧急', time: '10分钟前' },
  { id: 2, subject: '域名解析未生效', user: '李四', priority: '高', time: '30分钟前' },
  { id: 3, subject: '申请退款基础版主机', user: '王五', priority: '中', time: '1小时前' },
  { id: 4, subject: '控制面板登录异常', user: '赵六', priority: '低', time: '2小时前' },
  { id: 5, subject: 'SSL证书安装咨询', user: '孙七', priority: '低', time: '3小时前' },
])

function priorityType(priority: string) {
  const map: Record<string, string> = {
    '紧急': 'danger',
    '高': 'warning',
    '中': 'info',
    '低': 'info',
  }
  return (map[priority] || 'info') as any
}
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.stat-card {
  border-radius: 12px;
  margin-bottom: 16px;
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

.stat-value.text-red {
  color: #f5222d;
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

.stat-change.neutral {
  color: #fa8c16;
  font-size: 12px;
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
  background: rgba(0, 86, 255, 0.1);
  color: #0056FF;
}

.stat-icon.cyan {
  background: rgba(19, 194, 194, 0.1);
  color: #13c2c2;
}

.stat-icon.green {
  background: rgba(82, 196, 26, 0.1);
  color: #52c41a;
}

.stat-icon.red {
  background: rgba(245, 34, 45, 0.1);
  color: #f5222d;
}

.chart-card {
  border-radius: 12px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}

.ticket-list {
  display: flex;
  flex-direction: column;
}

.ticket-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid #f0f2f5;
}

.ticket-item:last-child {
  border-bottom: none;
}

.ticket-info {
  flex: 1;
  min-width: 0;
}

.ticket-subject {
  font-size: 14px;
  color: #1a1a2e;
  font-weight: 500;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ticket-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #8c8c8c;
  margin-top: 4px;
}
</style>
