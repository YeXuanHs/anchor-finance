<template>
  <div class="product-revenue-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>产品收入统计</span>
          <div class="header-filters">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              @change="handleDateChange"
            />
          </div>
        </div>
      </template>

      <!-- 产品收入排行 -->
      <div class="ranking-section">
        <h3>产品收入排行</h3>
        <el-table :data="productRanking" v-loading="rankingLoading" style="width: 100%" border>
          <el-table-column prop="rank" label="排名" width="80" align="center">
            <template #default="{ row, $index }">
              <el-tag v-if="$index < 3" :type="(['danger', 'warning', ''] as any)[$index]" size="small" round>
                {{ $index + 1 }}
              </el-tag>
              <span v-else>{{ $index + 1 }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="product_name" label="产品名称" min-width="200" />
          <el-table-column prop="product_type" label="类型" width="120" />
          <el-table-column prop="revenue" label="收入" width="140" align="right">
            <template #default="{ row }">
              <span class="amount-text">¥{{ formatAmount(row.revenue) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="order_count" label="订单数" width="100" align="center" />
          <el-table-column prop="client_count" label="客户数" width="100" align="center" />
          <el-table-column prop="avg_revenue" label="平均收入" width="120" align="right">
            <template #default="{ row }">
              ¥{{ formatAmount(row.avg_revenue) }}
            </template>
          </el-table-column>
          <el-table-column prop="revenue_share" label="占比" width="120" align="center">
            <template #default="{ row }">
              {{ row.revenue_share }}%
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 收入趋势图 -->
      <div class="chart-section">
        <h3>产品收入趋势</h3>
        <div class="chart-container" v-loading="trendLoading">
          <div class="trend-list">
            <div
              v-for="(item, index) in productTrend"
              :key="index"
              class="trend-item"
            >
              <div class="trend-header">
                <span class="trend-name">{{ item.product_name }}</span>
                <span class="trend-amount">¥{{ formatAmount(item.revenue) }}</span>
              </div>
              <el-progress
                :percentage="item.revenue_percent"
                :stroke-width="10"
                :color="getProgressColor(index)"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- 月度对比 -->
      <div class="table-section">
        <h3>月度产品收入对比</h3>
        <el-table :data="monthlyComparison" v-loading="comparisonLoading" style="width: 100%" border>
          <el-table-column prop="product_name" label="产品" min-width="150" />
          <el-table-column
            v-for="month in 12"
            :key="month"
            :label="month + '月'"
            width="100"
            align="right"
          >
            <template #default="{ row }">
              ¥{{ formatAmount(row.monthly[month] || 0) }}
            </template>
          </el-table-column>
          <el-table-column prop="total" label="合计" width="140" align="right" fixed="right">
            <template #default="{ row }">
              <span class="amount-text">¥{{ formatAmount(row.total) }}</span>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const rankingLoading = ref(false)
const trendLoading = ref(false)
const comparisonLoading = ref(false)

const dateRange = ref<[string, string] | null>(null)

const productRanking = ref<any[]>([])
const productTrend = ref<any[]>([])
const monthlyComparison = ref<any[]>([])

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const getProgressColor = (index: number) => {
  const colors = [
    '#5470c6', '#91cc75', '#fac858', '#ee6666', '#73c0de',
    '#3ba272', '#fc8452', '#9a60b4', '#ea7ccc', '#48b8d0'
  ]
  return colors[index % colors.length]
}

const fetchRanking = async () => {
  rankingLoading.value = true
  try {
    const params: any = {}
    if (dateRange.value) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    const data = await request.get({
      url: '/api/admin/reports/product-income',
      params
    })
    productRanking.value = data || []
  } catch (error) {
    console.error('获取产品收入排行失败:', error)
    ElMessage.error('获取产品收入排行失败')
  } finally {
    rankingLoading.value = false
  }
}

const fetchTrend = async () => {
  trendLoading.value = true
  try {
    const params: any = {}
    if (dateRange.value) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    const data = await request.get({
      url: '/api/admin/reports/product-income/trend',
      params
    })
    productTrend.value = data || []
  } catch (error) {
    console.error('获取产品收入趋势失败:', error)
  } finally {
    trendLoading.value = false
  }
}

const fetchComparison = async () => {
  comparisonLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/reports/product-income/comparison',
      params: { year: new Date().getFullYear() }
    })
    monthlyComparison.value = data || []
  } catch (error) {
    console.error('获取月度对比失败:', error)
  } finally {
    comparisonLoading.value = false
  }
}

const handleDateChange = () => {
  fetchRanking()
  fetchTrend()
}

onMounted(() => {
  fetchRanking()
  fetchTrend()
  fetchComparison()
})
</script>

<style scoped lang="scss">
.product-revenue-page {
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

.ranking-section, .chart-section, .table-section {
  margin-top: 24px;

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

.chart-container {
  padding: 20px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;
}

.trend-list {
  .trend-item {
    margin-bottom: 16px;

    &:last-child {
      margin-bottom: 0;
    }

    .trend-header {
      display: flex;
      justify-content: space-between;
      margin-bottom: 6px;

      .trend-name {
        font-weight: 500;
      }

      .trend-amount {
        color: var(--el-color-primary);
        font-weight: 600;
      }
    }
  }
}
</style>