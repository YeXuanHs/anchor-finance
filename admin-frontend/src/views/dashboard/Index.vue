<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <a-row :gutter="24" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic
            title="总客户数"
            :value="stats.totalCustomers"
            :value-style="{ color: '#165DFF' }"
          >
            <template #prefix><icon-user /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic
            title="本月收入"
            :value="stats.monthlyIncome"
            :precision="2"
            :value-style="{ color: '#00B42A' }"
          >
            <template #prefix>¥</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic
            title="待处理工单"
            :value="stats.pendingTickets"
            :value-style="{ color: '#FF7D00' }"
          >
            <template #prefix><icon-customer-service /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic
            title="活跃服务"
            :value="stats.activeServices"
            :value-style="{ color: '#722ED1' }"
          >
            <template #prefix><icon-desktop /></template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 图表区域 -->
    <a-row :gutter="24" class="chart-row">
      <a-col :span="16">
        <a-card title="收入趋势" class="chart-card">
          <div ref="incomeChartRef" class="chart-container"></div>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card title="最近订单" class="chart-card">
          <a-list :data="recentOrders" :bordered="false">
            <template #item="{ item }">
              <a-list-item>
                <a-list-item-meta
                  :title="item.product"
                  :description="`¥${item.amount}`"
                />
                <template #extra>
                  <a-tag :color="item.status === 'active' ? 'green' : 'orange'">
                    {{ item.status === 'active' ? '已完成' : '待处理' }}
                  </a-tag>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import * as echarts from 'echarts'

const incomeChartRef = ref<HTMLElement>()

const stats = reactive({
  totalCustomers: 1286,
  monthlyIncome: 156800.50,
  pendingTickets: 23,
  activeServices: 856
})

const recentOrders = ref([
  { product: '云服务器 2核4G', amount: '99.00', status: 'active' },
  { product: '虚拟主机 基础版', amount: '29.00', status: 'pending' },
  { product: 'SSL证书', amount: '199.00', status: 'active' },
  { product: 'CDN加速 100GB', amount: '49.00', status: 'active' },
  { product: '云服务器 4核8G', amount: '199.00', status: 'pending' }
])

onMounted(() => {
  initIncomeChart()
})

const initIncomeChart = () => {
  if (!incomeChartRef.value) return

  const chart = echarts.init(incomeChartRef.value)

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月']
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: '¥{value}'
      }
    },
    series: [
      {
        name: '收入',
        type: 'line',
        smooth: true,
        data: [12000, 15000, 18000, 16000, 21000, 19000, 23000, 25000, 22000, 28000, 26000, 30000],
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(22, 93, 255, 0.3)' },
            { offset: 1, color: 'rgba(22, 93, 255, 0.05)' }
          ])
        },
        itemStyle: {
          color: '#165DFF'
        }
      }
    ]
  }

  chart.setOption(option)

  window.addEventListener('resize', () => {
    chart.resize()
  })
}
</script>

<style scoped lang="scss">
.dashboard {
  .stats-row {
    margin-bottom: 24px;
  }

  .stat-card {
    text-align: center;

    :deep(.arco-statistic) {
      .arco-statistic-title {
        font-size: 14px;
        color: #86909c;
      }

      .arco-statistic-value {
        font-size: 28px;
        font-weight: 600;
      }
    }
  }

  .chart-row {
    .chart-card {
      height: 400px;
    }

    .chart-container {
      width: 100%;
      height: 300px;
    }
  }
}
</style>
