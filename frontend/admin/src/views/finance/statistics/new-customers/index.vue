<template>
  <div class="new-customers-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('page.newCustomers.title') }}</span>
          <div class="header-filters">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="/"
              :start-placeholder="$t('page.newCustomers.startDate')"
              :end-placeholder="$t('page.newCustomers.endDate')"
              value-format="YYYY-MM-DD"
              @change="handleDateChange"
            />
          </div>
        </div>
      </template>

      <!-- 概览卡片 -->
      <el-row :gutter="20" class="overview-cards">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.new_clients }}</div>
            <div class="stat-label">{{ $t('page.newCustomers.newClients') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.new_orders }}</div>
            <div class="stat-label">{{ $t('page.newCustomers.newOrders') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.completed_orders }}</div>
            <div class="stat-label">{{ $t('page.newCustomers.completedOrders') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.new_tickets }}</div>
            <div class="stat-label">{{ $t('page.newCustomers.newTickets') }}</div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px;">
        <el-col :span="12">
          <el-card shadow="hover" class="stat-card accent">
            <div class="stat-value">{{ stats.replied_tickets }}</div>
            <div class="stat-label">{{ $t('page.newCustomers.repliedTickets') }}</div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card shadow="hover" class="stat-card accent">
            <div class="stat-value">{{ stats.conversion_rate }}%</div>
            <div class="stat-label">{{ $t('page.newCustomers.conversionRate') }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 趋势图表 -->
      <div class="chart-section">
        <h3>{{ $t('page.newCustomers.dailyTrend') }}</h3>
        <div class="chart-container" v-loading="chartLoading">
          <div class="trend-chart">
            <div class="chart-body">
              <div
                v-for="(item, index) in dailyData"
                :key="index"
                class="chart-bar-wrapper"
                :title="item.date + ': ' + item.count + $t('page.newCustomers.personCount')"
              >
                <div
                  class="chart-bar"
                  :style="{ height: getBarHeight(item.count) + '%' }"
                />
                <div class="chart-label" v-if="index % Math.ceil(dailyData.length / 10) === 0">
                  {{ item.date.split('-').slice(1).join('/') }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 明细表格 -->
      <div class="table-section">
        <h3>{{ $t('page.newCustomers.dailyDetails') }}</h3>
        <el-table :data="dailyData" v-loading="tableLoading" style="width: 100%" border>
          <el-table-column prop="date" :label="$t('page.newCustomers.date')" width="120" />
          <el-table-column prop="new_clients" :label="$t('page.newCustomers.newClients')" width="120" align="center" />
          <el-table-column prop="new_orders" :label="$t('page.newCustomers.newOrder')" width="120" align="center" />
          <el-table-column prop="completed_orders" :label="$t('page.newCustomers.completedOrder')" width="120" align="center" />
          <el-table-column prop="new_tickets" :label="$t('page.newCustomers.newTicket')" width="120" align="center" />
          <el-table-column prop="replied_tickets" :label="$t('page.newCustomers.repliedTicket')" width="120" align="center" />
          <el-table-column :label="$t('page.newCustomers.conversionRate')" width="120" align="center">
            <template #default="{ row }">
              {{ row.new_clients > 0 ? ((row.new_orders / row.new_clients) * 100).toFixed(1) : '0.0' }}%
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const chartLoading = ref(false)
const tableLoading = ref(false)

const dateRange = ref<[string, string] | null>(null)

const stats = reactive({
  new_clients: 0,
  new_orders: 0,
  completed_orders: 0,
  new_tickets: 0,
  replied_tickets: 0,
  conversion_rate: 0
})

const dailyData = ref<any[]>([])

const getBarHeight = (value: number) => {
  if (!dailyData.value.length) return 0
  const max = Math.max(...dailyData.value.map(d => d.count || d.new_clients))
  return max > 0 ? (value / max) * 100 : 0
}

const fetchStatistics = async () => {
  chartLoading.value = true
  tableLoading.value = true
  try {
    const params: any = {}
    if (dateRange.value) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    const data = await request.get({
      url: '/api/admin/reports/new-client-statistics',
      params
    })
    stats.new_clients = data.new_clients || 0
    stats.new_orders = data.new_orders || 0
    stats.completed_orders = data.completed_orders || 0
    stats.new_tickets = data.new_tickets || 0
    stats.replied_tickets = data.replied_tickets || 0
    stats.conversion_rate = data.conversion_rate || 0
    dailyData.value = (data.daily || []).map((item: any) => ({
      ...item,
      count: item.new_clients
    }))
  } catch (error) {
    console.error('获取新客户统计失败:', error)
    ElMessage.error($t('page.newCustomers.fetchFailed'))
  } finally {
    chartLoading.value = false
    tableLoading.value = false
  }
}

const handleDateChange = () => {
  fetchStatistics()
}

onMounted(() => {
  fetchStatistics()
})
</script>

<style scoped lang="scss">
.new-customers-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-filters {
  display: flex;
  gap: 12px;
}

.overview-cards {
  margin-bottom: 12px;
}

.stat-card {
  text-align: center;
  padding: 8px 0;

  .stat-value {
    font-size: 28px;
    font-weight: 600;
    color: var(--el-color-primary);
    margin-bottom: 4px;
  }

  .stat-label {
    color: var(--el-text-color-secondary);
    font-size: 14px;
  }

  &.accent .stat-value {
    color: var(--el-color-success);
  }
}

.chart-section, .table-section {
  margin-top: 24px;

  h3 {
    margin: 0 0 16px;
    font-size: 16px;
    font-weight: 600;
  }
}

.chart-container {
  padding: 20px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;
}

.trend-chart {
  .chart-body {
    display: flex;
    align-items: flex-end;
    height: 200px;
    gap: 2px;

    .chart-bar-wrapper {
      flex: 1;
      display: flex;
      flex-direction: column;
      align-items: center;
      height: 100%;
      cursor: pointer;

      .chart-bar {
        width: 100%;
        background: var(--el-color-primary);
        border-radius: 4px 4px 0 0;
        min-height: 4px;
        transition: height 0.3s ease, background 0.2s;

        &:hover {
          background: var(--el-color-primary-dark-2);
        }
      }

      .chart-label {
        margin-top: 8px;
        font-size: 11px;
        color: var(--el-text-color-secondary);
        white-space: nowrap;
      }
    }
  }
}
</style>