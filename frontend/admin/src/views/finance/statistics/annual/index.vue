<template>
  <div class="annual-statistics-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>年度统计</span>
          <div class="header-filters">
            <el-select v-model="selectedYear" placeholder="选择年份" @change="handleYearChange">
              <el-option
                v-for="year in yearOptions"
                :key="year"
                :label="year + '年'"
                :value="year"
              />
            </el-select>
            <el-select v-model="selectedCurrency" placeholder="货币" @change="handleCurrencyChange">
              <el-option label="人民币 (CNY)" value="CNY" />
              <el-option label="美元 (USD)" value="USD" />
            </el-select>
          </div>
        </div>
      </template>

      <!-- 年度概览卡片 -->
      <el-row :gutter="20" class="overview-cards">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card income">
            <div class="stat-icon">
              <el-icon :size="40"><TrendCharts /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">年度收入</div>
              <div class="stat-value">¥{{ formatAmount(stats.total_income) }}</div>
              <div class="stat-compare" :class="stats.income_growth >= 0 ? 'up' : 'down'">
                {{ stats.income_growth >= 0 ? '+' : '' }}{{ stats.income_growth }}% 同比
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card expense">
            <div class="stat-icon">
              <el-icon :size="40"><Bottom /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">年度支出</div>
              <div class="stat-value">¥{{ formatAmount(stats.total_expense) }}</div>
              <div class="stat-compare" :class="stats.expense_growth <= 0 ? 'up' : 'down'">
                {{ stats.expense_growth >= 0 ? '+' : '' }}{{ stats.expense_growth }}% 同比
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card balance">
            <div class="stat-icon">
              <el-icon :size="40"><Coin /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">年度结余</div>
              <div class="stat-value">¥{{ formatAmount(stats.total_income - stats.total_expense) }}</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card orders">
            <div class="stat-icon">
              <el-icon :size="40"><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">年度订单数</div>
              <div class="stat-value">{{ stats.total_orders }}</div>
              <div class="stat-compare" :class="stats.order_growth >= 0 ? 'up' : 'down'">
                {{ stats.order_growth >= 0 ? '+' : '' }}{{ stats.order_growth }}% 同比
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 月度收入趋势图 -->
      <div class="chart-section">
        <h3>月度收入趋势</h3>
        <div class="chart-container" v-loading="chartLoading">
          <div class="monthly-chart">
            <div class="chart-header">
              <span class="legend income">收入</span>
              <span class="legend expense">支出</span>
              <span class="legend balance">结余</span>
            </div>
            <div class="chart-body">
              <div
                v-for="(item, index) in monthlyData"
                :key="index"
                class="chart-column"
              >
                <div class="chart-bars">
                  <div
                    class="bar income-bar"
                    :style="{ height: getBarHeight(item.income) + '%' }"
                    :title="'收入: ¥' + formatAmount(item.income)"
                  />
                  <div
                    class="bar expense-bar"
                    :style="{ height: getBarHeight(item.expense) + '%' }"
                    :title="'支出: ¥' + formatAmount(item.expense)"
                  />
                </div>
                <div class="chart-label">{{ item.month }}月</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 月度明细表格 -->
      <div class="table-section">
        <h3>月度明细</h3>
        <el-table :data="monthlyData" v-loading="tableLoading" style="width: 100%" border>
          <el-table-column prop="month" label="月份" width="100" align="center">
            <template #default="{ row }">
              {{ row.month }}月
            </template>
          </el-table-column>
          <el-table-column prop="income" label="收入" width="140" align="right">
            <template #default="{ row }">
              <span class="text-success">¥{{ formatAmount(row.income) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="expense" label="支出" width="140" align="right">
            <template #default="{ row }">
              <span class="text-danger">¥{{ formatAmount(row.expense) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="结余" width="140" align="right">
            <template #default="{ row }">
              <span :class="row.income - row.expense >= 0 ? 'text-success' : 'text-danger'">
                ¥{{ formatAmount(row.income - row.expense) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="order_count" label="订单数" width="100" align="center" />
          <el-table-column prop="new_clients" label="新客户" width="100" align="center" />
          <el-table-column label="月环比" width="120" align="center">
            <template #default="{ row, $index }">
              <span
                v-if="$index > 0"
                :class="monthlyData[$index].income - monthlyData[$index - 1].income >= 0 ? 'text-success' : 'text-danger'"
              >
                {{ monthlyData[$index].income - monthlyData[$index - 1].income >= 0 ? '+' : '' }}
                {{ ((monthlyData[$index].income - monthlyData[$index - 1].income) / (monthlyData[$index - 1].income || 1) * 100).toFixed(1) }}%
              </span>
              <span v-else>-</span>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { TrendCharts, Bottom, Coin, Document } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const chartLoading = ref(false)
const tableLoading = ref(false)

const currentYear = new Date().getFullYear()
const selectedYear = ref(currentYear)
const selectedCurrency = ref('CNY')

const yearOptions = computed(() => {
  const years = []
  for (let y = currentYear; y >= currentYear - 5; y--) {
    years.push(y)
  }
  return years
})

const stats = reactive({
  total_income: 0,
  total_expense: 0,
  total_orders: 0,
  income_growth: 0,
  expense_growth: 0,
  order_growth: 0
})

const monthlyData = ref<any[]>([])

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const getBarHeight = (value: number) => {
  if (!monthlyData.value.length) return 0
  const max = Math.max(...monthlyData.value.map(d => Math.max(d.income, d.expense)))
  return max > 0 ? (value / max) * 100 : 0
}

const fetchStatistics = async () => {
  chartLoading.value = true
  tableLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/reports/year-income-statistics',
      params: {
        year: selectedYear.value,
        currency: selectedCurrency.value
      }
    })
    stats.total_income = data.total_income || 0
    stats.total_expense = data.total_expense || 0
    stats.total_orders = data.total_orders || 0
    stats.income_growth = data.income_growth || 0
    stats.expense_growth = data.expense_growth || 0
    stats.order_growth = data.order_growth || 0
    monthlyData.value = data.monthly || []
  } catch (error) {
    console.error('获取年度统计失败:', error)
    ElMessage.error('获取年度统计失败')
  } finally {
    chartLoading.value = false
    tableLoading.value = false
  }
}

const handleYearChange = () => {
  fetchStatistics()
}

const handleCurrencyChange = () => {
  fetchStatistics()
}

onMounted(() => {
  fetchStatistics()
})
</script>

<style scoped lang="scss">
.annual-statistics-page {
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
  margin-bottom: 24px;

  .stat-card {
    .stat-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 60px;
      height: 60px;
      border-radius: 12px;
      margin-bottom: 12px;

      &.income {
        background: var(--el-color-success-light-9);
        color: var(--el-color-success);
      }

      &.expense {
        background: var(--el-color-danger-light-9);
        color: var(--el-color-danger);
      }

      &.balance {
        background: var(--el-color-warning-light-9);
        color: var(--el-color-warning);
      }

      &.orders {
        background: var(--el-color-primary-light-9);
        color: var(--el-color-primary);
      }
    }

    .stat-info {
      .stat-label {
        color: var(--el-text-color-secondary);
        font-size: 14px;
        margin-bottom: 4px;
      }

      .stat-value {
        font-size: 22px;
        font-weight: 600;
        color: var(--el-text-color-primary);
        margin-bottom: 4px;
      }

      .stat-compare {
        font-size: 13px;

        &.up {
          color: var(--el-color-success);
        }

        &.down {
          color: var(--el-color-danger);
        }
      }
    }
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

.monthly-chart {
  .chart-header {
    display: flex;
    gap: 24px;
    margin-bottom: 16px;
    justify-content: flex-end;

    .legend {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 13px;
      color: var(--el-text-color-secondary);

      &::before {
        content: '';
        display: inline-block;
        width: 12px;
        height: 12px;
        border-radius: 2px;
      }

      &.income::before {
        background: var(--el-color-success);
      }

      &.expense::before {
        background: var(--el-color-danger);
      }

      &.balance::before {
        background: var(--el-color-warning);
      }
    }
  }

  .chart-body {
    display: flex;
    align-items: flex-end;
    height: 240px;
    gap: 8px;
    padding-bottom: 30px;
    position: relative;

    .chart-column {
      flex: 1;
      display: flex;
      flex-direction: column;
      align-items: center;
      height: 100%;

      .chart-bars {
        flex: 1;
        display: flex;
        gap: 4px;
        align-items: flex-end;
        width: 100%;

        .bar {
          flex: 1;
          border-radius: 4px 4px 0 0;
          min-height: 4px;
          transition: height 0.3s ease;

          &.income-bar {
            background: var(--el-color-success);
          }

          &.expense-bar {
            background: var(--el-color-danger);
          }
        }
      }

      .chart-label {
        margin-top: 8px;
        font-size: 12px;
        color: var(--el-text-color-secondary);
      }
    }
  }
}

.text-success {
  color: var(--el-color-success);
}

.text-danger {
  color: var(--el-color-danger);
}
</style>