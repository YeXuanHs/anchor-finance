<template>
  <div class="revenue-ranking-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>收入排行</span>
          <div class="header-filters">
            <el-select v-model="selectedPeriod" placeholder="时间段" @change="handlePeriodChange">
              <el-option label="本月" value="month" />
              <el-option label="本季度" value="quarter" />
              <el-option label="本年度" value="year" />
              <el-option label="全部" value="all" />
            </el-select>
          </div>
        </div>
      </template>

      <el-row :gutter="20">
        <!-- 客户收入排行 -->
        <el-col :span="12">
          <div class="ranking-section">
            <h3>客户收入排行 (前10名)</h3>
            <el-table :data="clientRanking" v-loading="clientLoading" style="width: 100%" border>
              <el-table-column prop="rank" label="排名" width="70" align="center">
                <template #default="{ row, $index }">
                  <el-tag v-if="$index < 3" :type="(['danger', 'warning', ''] as any)[$index]" size="small" round>
                    {{ $index + 1 }}
                  </el-tag>
                  <span v-else>{{ $index + 1 }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="client_name" label="客户" min-width="120">
                <template #default="{ row }">
                  <el-button type="primary" link @click="handleViewClient(row)">
                    {{ row.client_name }}
                  </el-button>
                </template>
              </el-table-column>
              <el-table-column prop="revenue" label="收入" width="130" align="right">
                <template #default="{ row }">
                  <span class="amount-text">¥{{ formatAmount(row.revenue) }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="expense" label="支出" width="130" align="right">
                <template #default="{ row }">
                  ¥{{ formatAmount(row.expense) }}
                </template>
              </el-table-column>
              <el-table-column prop="net" label="净收入" width="130" align="right">
                <template #default="{ row }">
                  <span :class="row.net >= 0 ? 'text-success' : 'text-danger'">
                    ¥{{ formatAmount(row.net) }}
                  </span>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-col>

        <!-- 收入/支出对比 -->
        <el-col :span="12">
          <div class="ranking-section">
            <h3>收入/支出对比</h3>
            <div class="comparison-chart" v-loading="comparisonLoading">
              <div class="chart-bars">
                <div
                  v-for="(item, index) in comparisonData"
                  :key="index"
                  class="bar-group"
                >
                  <div class="bar-label">{{ item.label }}</div>
                  <div class="bar-container">
                    <div
                      class="bar income-bar"
                      :style="{ width: getBarWidth(item.income) + '%' }"
                      :title="'收入: ¥' + formatAmount(item.income)"
                    />
                    <div
                      class="bar expense-bar"
                      :style="{ width: getBarWidth(item.expense) + '%' }"
                      :title="'支出: ¥' + formatAmount(item.expense)"
                    />
                  </div>
                  <div class="bar-values">
                    <span class="income-value">¥{{ formatAmount(item.income) }}</span>
                    <span class="expense-value">¥{{ formatAmount(item.expense) }}</span>
                  </div>
                </div>
              </div>
              <div class="chart-legend">
                <span class="legend-item">
                  <span class="dot income" />收入
                </span>
                <span class="legend-item">
                  <span class="dot expense" />支出
                </span>
              </div>
            </div>

            <!-- 汇总数据 -->
            <el-card shadow="never" style="margin-top: 20px;">
              <template #header>
                <span>汇总数据</span>
              </template>
              <el-descriptions :column="2" border size="small">
                <el-descriptions-item label="总收入">
                  <span class="text-success">¥{{ formatAmount(totalStats.total_income) }}</span>
                </el-descriptions-item>
                <el-descriptions-item label="总支出">
                  <span class="text-danger">¥{{ formatAmount(totalStats.total_expense) }}</span>
                </el-descriptions-item>
                <el-descriptions-item label="净收入">
                  <span :class="totalStats.total_income - totalStats.total_expense >= 0 ? 'text-success' : 'text-danger'">
                    ¥{{ formatAmount(totalStats.total_income - totalStats.total_expense) }}
                  </span>
                </el-descriptions-item>
                <el-descriptions-item label="利润率">
                  {{ totalStats.total_income > 0 ? (((totalStats.total_income - totalStats.total_expense) / totalStats.total_income) * 100).toFixed(1) : 0 }}%
                </el-descriptions-item>
              </el-descriptions>
            </el-card>
          </div>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const router = useRouter()

const clientLoading = ref(false)
const comparisonLoading = ref(false)

const selectedPeriod = ref('year')

const clientRanking = ref<any[]>([])
const comparisonData = ref<any[]>([])

const totalStats = reactive({
  total_income: 0,
  total_expense: 0
})

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const getBarWidth = (value: number) => {
  const maxIncome = Math.max(...comparisonData.value.map(d => d.income), 1)
  return (value / maxIncome) * 100
}

const fetchClientRanking = async () => {
  clientLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/reports/revenue-ranking',
      params: { period: selectedPeriod.value }
    })
    clientRanking.value = data || []
  } catch (error) {
    console.error('获取客户收入排行失败:', error)
    ElMessage.error('获取客户收入排行失败')
  } finally {
    clientLoading.value = false
  }
}

const fetchComparison = async () => {
  comparisonLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/reports/revenue-ranking/comparison',
      params: { period: selectedPeriod.value }
    })
    comparisonData.value = data.items || []
    totalStats.total_income = data.total_income || 0
    totalStats.total_expense = data.total_expense || 0
  } catch (error) {
    console.error('获取收入支出对比失败:', error)
    ElMessage.error('获取收入支出对比失败')
  } finally {
    comparisonLoading.value = false
  }
}

const handlePeriodChange = () => {
  fetchClientRanking()
  fetchComparison()
}

const handleViewClient = (row: any) => {
  router.push(`/customer-view/${row.client_id}`)
}

onMounted(() => {
  fetchClientRanking()
  fetchComparison()
})
</script>

<style scoped lang="scss">
.revenue-ranking-page {
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

.ranking-section {
  h3 {
    margin: 0 0 16px;
    font-size: 16px;
    font-weight: 600;
  }
}

.amount-text {
  font-weight: 600;
  color: var(--el-color-primary);
}

.text-success {
  color: var(--el-color-success);
}

.text-danger {
  color: var(--el-color-danger);
}

.comparison-chart {
  padding: 20px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;

  .chart-bars {
    .bar-group {
      margin-bottom: 16px;

      &:last-child {
        margin-bottom: 0;
      }

      .bar-label {
        font-weight: 500;
        margin-bottom: 6px;
        font-size: 14px;
      }

      .bar-container {
        display: flex;
        flex-direction: column;
        gap: 4px;

        .bar {
          height: 20px;
          border-radius: 4px;
          min-width: 4px;
          transition: width 0.3s ease;

          &.income-bar {
            background: var(--el-color-success);
          }

          &.expense-bar {
            background: var(--el-color-danger);
          }
        }
      }

      .bar-values {
        display: flex;
        gap: 24px;
        margin-top: 4px;
        font-size: 13px;

        .income-value {
          color: var(--el-color-success);
        }

        .expense-value {
          color: var(--el-color-danger);
        }
      }
    }
  }

  .chart-legend {
    display: flex;
    gap: 24px;
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid var(--el-border-color-lighter);

    .legend-item {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 13px;

      .dot {
        width: 12px;
        height: 12px;
        border-radius: 2px;

        &.income {
          background: var(--el-color-success);
        }

        &.expense {
          background: var(--el-color-danger);
        }
      }
    }
  }
}
</style>