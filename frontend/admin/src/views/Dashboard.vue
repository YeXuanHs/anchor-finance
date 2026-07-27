<template>
  <div class="dashboard">
    <!-- Stat Cards -->
    <n-grid :cols="4" :x-gap="16" :y-gap="16" responsive="screen" :item-responsive="true">
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card" :bordered="false" rounded>
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">今日收入</span>
              <span class="stat-value">¥{{ stats.todayRevenue.toLocaleString() }}</span>
              <span class="stat-change" :class="stats.revenueChange >= 0 ? 'up' : 'down'">
                <n-icon size="14">
                  <TrendUpIcon v-if="stats.revenueChange >= 0" />
                  <TrendDownIcon v-else />
                </n-icon>
                {{ Math.abs(stats.revenueChange) }}%
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
              <span class="stat-label">今日订单</span>
              <span class="stat-value">{{ stats.todayOrders }}</span>
              <span class="stat-change" :class="stats.ordersChange >= 0 ? 'up' : 'down'">
                <n-icon size="14">
                  <TrendUpIcon v-if="stats.ordersChange >= 0" />
                  <TrendDownIcon v-else />
                </n-icon>
                {{ Math.abs(stats.ordersChange) }}%
              </span>
            </div>
            <div class="stat-icon cyan">
              <n-icon size="28"><CartIcon /></n-icon>
            </div>
          </div>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card" :bordered="false" rounded>
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">新增用户</span>
              <span class="stat-value">{{ stats.newUsers }}</span>
              <span class="stat-change" :class="stats.usersChange >= 0 ? 'up' : 'down'">
                <n-icon size="14">
                  <TrendUpIcon v-if="stats.usersChange >= 0" />
                  <TrendDownIcon v-else />
                </n-icon>
                {{ Math.abs(stats.usersChange) }}%
              </span>
            </div>
            <div class="stat-icon green">
              <n-icon size="28"><PeopleIcon /></n-icon>
            </div>
          </div>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card" :bordered="false" rounded>
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">待处理工单</span>
              <span class="stat-value" :class="{ 'text-red': stats.openTickets > 0 }">{{ stats.openTickets }}</span>
              <span class="stat-change neutral" v-if="stats.openTickets > 0">
                需要处理
              </span>
            </div>
            <div class="stat-icon" :class="stats.openTickets > 0 ? 'red' : 'green'">
              <n-icon size="28"><TicketIcon /></n-icon>
            </div>
          </div>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- Charts Row -->
    <n-grid :cols="3" :x-gap="16" style="margin-top: 16px" responsive="screen" :item-responsive="true">
      <n-gi span="3 l:2">
        <n-card title="收入趋势（近30天）" :bordered="false" rounded class="chart-card">
          <v-chart :option="revenueChartOption" autoresize style="height: 360px" />
        </n-card>
      </n-gi>
      <n-gi span="3 l:1">
        <n-card title="订单状态分布" :bordered="false" rounded class="chart-card">
          <v-chart :option="orderPieOption" autoresize style="height: 360px" />
        </n-card>
      </n-gi>
    </n-grid>

    <!-- Recent Orders & Tickets -->
    <n-grid :cols="3" :x-gap="16" style="margin-top: 16px" responsive="screen" :item-responsive="true">
      <n-gi span="3 l:2">
        <n-card title="最近订单" :bordered="false" rounded>
          <template #header-extra>
            <n-button text type="primary" @click="$router.push('/admin/orders')">查看全部</n-button>
          </template>
          <n-data-table
            :columns="orderColumns"
            :data="recentOrders"
            :bordered="false"
            :pagination="false"
            size="small"
          />
        </n-card>
      </n-gi>
      <n-gi span="3 l:1">
        <n-card title="最近工单" :bordered="false" rounded>
          <template #header-extra>
            <n-button text type="primary" @click="$router.push('/admin/tickets')">查看全部</n-button>
          </template>
          <n-list hoverable clickable>
            <n-list-item v-for="ticket in recentTickets" :key="ticket.id">
              <n-thing>
                <template #header>
                  <span class="ticket-subject">{{ ticket.subject }}</span>
                </template>
                <template #header-extra>
                  <n-tag :type="priorityType(ticket.priority)" size="small" round>
                    {{ ticket.priority }}
                  </n-tag>
                </template>
                <template #description>
                  <div class="ticket-meta">
                    <span>{{ ticket.user }}</span>
                    <span>{{ ticket.time }}</span>
                  </div>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { h, reactive, computed } from 'vue'
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
  TrendUpOutline as TrendUpIcon,
  TrendDownOutline as TrendDownIcon,
  CashOutline as CashIcon,
  CartOutline as CartIcon,
  PeopleOutline as PeopleIcon,
  ChatbubblesOutline as TicketIcon,
} from '@vicons/ionicons5'
import { NTag, NButton, type DataTableColumns } from 'naive-ui'

use([CanvasRenderer, LineChart, PieChart, TitleComponent, TooltipComponent, LegendComponent, GridComponent])

// ---- Stats ----
const stats = reactive({
  todayRevenue: 12680,
  revenueChange: 12.5,
  todayOrders: 86,
  ordersChange: 8.3,
  newUsers: 34,
  usersChange: -2.1,
  openTickets: 7,
})

// ---- Revenue Chart (Last 30 days) ----
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
      return `<div style="font-weight:600">${p.axisValue}</div><div>收入: <b style="color:#1890ff">¥${p.value.toLocaleString()}</b></div>`
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
      lineStyle: { color: '#1890ff', width: 3 },
      itemStyle: { color: '#1890ff', borderWidth: 2, borderColor: '#fff' },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(24, 144, 255, 0.25)' },
            { offset: 0.5, color: 'rgba(24, 144, 255, 0.08)' },
            { offset: 1, color: 'rgba(24, 144, 255, 0.01)' },
          ],
        },
      },
      emphasis: {
        focus: 'series',
        itemStyle: { borderWidth: 3, shadowBlur: 10, shadowColor: 'rgba(24,144,255,0.3)' },
      },
      data: generateRevenueData(),
    },
  ],
}))

// ---- Order Pie Chart ----
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
        { value: 340, name: '待支付', itemStyle: { color: '#1890ff' } },
        { value: 580, name: '已支付', itemStyle: { color: '#52c41a' } },
        { value: 1220, name: '已开通', itemStyle: { color: '#fa8c16' } },
        { value: 180, name: '已取消', itemStyle: { color: '#d9d9d9' } },
      ],
    },
  ],
}))

// ---- Recent Orders Table ----
const statusMap: Record<string, { label: string; color: string }> = {
  pending: { label: '待支付', color: 'warning' },
  paid: { label: '已支付', color: 'info' },
  active: { label: '已开通', color: 'success' },
  cancelled: { label: '已取消', color: 'default' },
}

const recentOrders = ref([
  { id: 'AF20260726001', user: '张三', product: '基础版主机', amount: 299, status: 'active', time: '2026-07-26 14:30' },
  { id: 'AF20260726002', user: '李四', product: '高级版主机', amount: 599, status: 'paid', time: '2026-07-26 13:22' },
  { id: 'AF20260726003', user: '王五', product: '企业版主机', amount: 1299, status: 'pending', time: '2026-07-26 11:45' },
  { id: 'AF20260726004', user: '赵六', product: '1核2G云服务器', amount: 89, status: 'active', time: '2026-07-26 10:18' },
  { id: 'AF20260726005', user: '孙七', product: '4核8G云服务器', amount: 399, status: 'cancelled', time: '2026-07-25 16:55' },
])

const orderColumns: DataTableColumns<any> = [
  { title: '订单号', key: 'id', width: 150, ellipsis: { tooltip: true } },
  { title: '用户', key: 'user', width: 80 },
  { title: '产品', key: 'product', ellipsis: { tooltip: true } },
  {
    title: '金额',
    key: 'amount',
    width: 100,
    render: (row) => h('span', { style: 'font-weight:600;color:#1890ff' }, `¥${row.amount}`),
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
  { title: '时间', key: 'time', width: 150 },
]

// ---- Recent Tickets ----
const recentTickets = ref([
  { id: 1, subject: '服务器无法连接SSH', user: '张三', priority: '紧急', time: '10分钟前' },
  { id: 2, subject: '域名解析未生效', user: '李四', priority: '高', time: '30分钟前' },
  { id: 3, subject: '申请退款基础版主机', user: '王五', priority: '中', time: '1小时前' },
  { id: 4, subject: '控制面板登录异常', user: '赵六', priority: '低', time: '2小时前' },
  { id: 5, subject: 'SSL证书安装咨询', user: '孙七', priority: '低', time: '3小时前' },
])

function priorityType(priority: string): 'error' | 'warning' | 'info' | 'success' | 'default' {
  const map: Record<string, 'error' | 'warning' | 'info' | 'success' | 'default'> = {
    '紧急': 'error',
    '高': 'warning',
    '中': 'info',
    '低': 'default',
  }
  return map[priority] || 'default'
}
</script>

<style scoped>
.dashboard {
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

.stat-icon.red {
  background: rgba(245, 34, 45, 0.1);
  color: #f5222d;
}

.chart-card {
  border-radius: 12px;
}

.ticket-subject {
  font-size: 14px;
  color: #1a1a2e;
}

.ticket-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #8c8c8c;
}
</style>
