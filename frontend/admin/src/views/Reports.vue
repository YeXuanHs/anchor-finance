<template>
  <div class="admin-page">
    <el-row :gutter="16">
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">本月收入</span>
              <span class="stat-value">¥356,800</span>
            </div>
            <div class="stat-icon blue"><el-icon :size="28"><Wallet /></el-icon></div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">总用户数</span>
              <span class="stat-value">12,580</span>
            </div>
            <div class="stat-icon green"><el-icon :size="28"><User /></el-icon></div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">总订单数</span>
              <span class="stat-value">8,920</span>
            </div>
            <div class="stat-icon cyan"><el-icon :size="28"><ShoppingCart /></el-icon></div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">活跃产品</span>
              <span class="stat-value">36</span>
            </div>
            <div class="stat-icon orange"><el-icon :size="28"><Goods /></el-icon></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :lg="16">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>月度收入趋势</span>
              <el-radio-group v-model="chartPeriod" size="small">
                <el-radio-button value="week">近7天</el-radio-button>
                <el-radio-button value="month">近30天</el-radio-button>
                <el-radio-button value="year">近12月</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <v-chart :option="revenueChartOption" autoresize style="height: 400px" />
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="8">
        <el-card shadow="hover">
          <template #header><span>产品销售排行</span></template>
          <div class="rank-list">
            <div v-for="(item, idx) in productRank" :key="idx" class="rank-item">
              <span class="rank-num" :class="{ top: idx < 3 }">{{ idx + 1 }}</span>
              <span class="rank-name">{{ item.name }}</span>
              <span class="rank-value">{{ item.count }} 单</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :lg="12">
        <el-card shadow="hover">
          <template #header><span>用户增长趋势</span></template>
          <v-chart :option="userGrowthOption" autoresize style="height: 300px" />
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="12">
        <el-card shadow="hover">
          <template #header><span>支付方式占比</span></template>
          <v-chart :option="paymentPieOption" autoresize style="height: 300px" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, BarChart, PieChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components'
import { Wallet, User, ShoppingCart, Goods } from '@element-plus/icons-vue'

use([CanvasRenderer, LineChart, BarChart, PieChart, TitleComponent, TooltipComponent, LegendComponent, GridComponent])

const chartPeriod = ref('month')

const productRank = [
  { name: '基础版主机', count: 356 },
  { name: '4核8G云服务器', count: 289 },
  { name: '高级版主机', count: 245 },
  { name: '.com域名注册', count: 198 },
  { name: '1核2G云服务器', count: 167 },
  { name: '企业版主机', count: 134 },
  { name: 'DV SSL证书', count: 89 },
]

function getLast12Months() {
  const months = []
  const now = new Date()
  for (let i = 11; i >= 0; i--) {
    const d = new Date(now.getFullYear(), now.getMonth() - i, 1)
    months.push(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`)
  }
  return months
}

const revenueChartOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: getLast12Months(), axisLabel: { color: '#8c8c8c' } },
  yAxis: { type: 'value', axisLabel: { color: '#8c8c8c', formatter: (v: number) => `¥${(v / 10000).toFixed(0)}万` }, splitLine: { lineStyle: { color: '#f0f0f0' } } },
  series: [{
    type: 'bar', data: [280000, 310000, 295000, 340000, 325000, 356000, 380000, 365000, 390000, 410000, 395000, 420000],
    itemStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: '#0056FF' }, { offset: 1, color: '#4080FF' }] }, borderRadius: [4, 4, 0, 0] },
    barWidth: '40%',
  }],
}))

const userGrowthOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: getLast12Months(), axisLabel: { color: '#8c8c8c' } },
  yAxis: { type: 'value', axisLabel: { color: '#8c8c8c' }, splitLine: { lineStyle: { color: '#f0f0f0' } } },
  series: [{
    type: 'line', smooth: true, data: [820, 932, 901, 1034, 1290, 1330, 1520, 1430, 1680, 1820, 1950, 2100],
    lineStyle: { color: '#52c41a', width: 3 }, itemStyle: { color: '#52c41a' },
    areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(82,196,26,0.2)' }, { offset: 1, color: 'rgba(82,196,26,0.01)' }] } },
  }],
}))

const paymentPieOption = computed(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
  legend: { bottom: '2%', left: 'center', itemWidth: 10, itemHeight: 10 },
  series: [{
    type: 'pie', radius: ['40%', '65%'], center: ['50%', '45%'],
    itemStyle: { borderRadius: 8, borderColor: '#fff', borderWidth: 3 },
    label: { show: false }, emphasis: { label: { show: true, fontSize: 14, fontWeight: 'bold' } },
    data: [
      { value: 1256, name: '支付宝', itemStyle: { color: '#0056FF' } },
      { value: 1589, name: '微信支付', itemStyle: { color: '#52c41a' } },
      { value: 234, name: 'PayPal', itemStyle: { color: '#fa8c16' } },
      { value: 456, name: 'Stripe', itemStyle: { color: '#f5222d' } },
      { value: 567, name: '余额', itemStyle: { color: '#722ed1' } },
    ],
  }],
}))
</script>

<style scoped>
.stat-card { border-radius: 12px; margin-bottom: 16px; }
.stat-content { display: flex; align-items: center; justify-content: space-between; }
.stat-info { display: flex; flex-direction: column; }
.stat-label { font-size: 13px; color: #8c8c8c; margin-bottom: 8px; }
.stat-value { font-size: 28px; font-weight: 700; color: #1a1a2e; }
.stat-icon { width: 56px; height: 56px; border-radius: 14px; display: flex; align-items: center; justify-content: center; }
.stat-icon.blue { background: rgba(0, 86, 255, 0.1); color: #0056FF; }
.stat-icon.green { background: rgba(82, 196, 26, 0.1); color: #52c41a; }
.stat-icon.cyan { background: rgba(19, 194, 194, 0.1); color: #13c2c2; }
.stat-icon.orange { background: rgba(250, 140, 22, 0.1); color: #fa8c16; }
.card-header { display: flex; align-items: center; justify-content: space-between; }
.rank-list { display: flex; flex-direction: column; gap: 12px; }
.rank-item { display: flex; align-items: center; gap: 12px; }
.rank-num { width: 24px; height: 24px; border-radius: 50%; background: #f0f2f5; display: flex; align-items: center; justify-content: center; font-size: 12px; font-weight: 600; color: #666; }
.rank-num.top { background: #0056FF; color: #fff; }
.rank-name { flex: 1; font-size: 14px; }
.rank-value { font-size: 13px; color: #8c8c8c; }
</style>
