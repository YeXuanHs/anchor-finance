<template>
  <div class="statistics-overview-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>统计概览</span>
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

      <!-- 概览卡片 -->
      <el-row :gutter="20" class="overview-cards">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card income">
            <div class="stat-icon">
              <el-icon :size="40"><TrendCharts /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">总收入</div>
              <div class="stat-value">¥{{ formatAmount(stats.total_income) }}</div>
              <div class="stat-compare" :class="stats.income_growth >= 0 ? 'up' : 'down'">
                {{ stats.income_growth >= 0 ? '+' : '' }}{{ stats.income_growth }}%
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card orders">
            <div class="stat-icon">
              <el-icon :size="40"><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">订单数</div>
              <div class="stat-value">{{ stats.total_orders }}</div>
              <div class="stat-compare" :class="stats.order_growth >= 0 ? 'up' : 'down'">
                {{ stats.order_growth >= 0 ? '+' : '' }}{{ stats.order_growth }}%
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card clients">
            <div class="stat-icon">
              <el-icon :size="40"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">客户数</div>
              <div class="stat-value">{{ stats.total_clients }}</div>
              <div class="stat-compare" :class="stats.client_growth >= 0 ? 'up' : 'down'">
                {{ stats.client_growth >= 0 ? '+' : '' }}{{ stats.client_growth }}%
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card tickets">
            <div class="stat-icon">
              <el-icon :size="40"><ChatDotRound /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">待处理工单</div>
              <div class="stat-value">{{ stats.pending_tickets }}</div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 快捷导航 -->
      <div class="quick-links">
        <h3>详细统计</h3>
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card shadow="hover" class="link-card" @click="router.push('/finance/statistics/revenue-ranking')">
              <el-icon :size="32"><Trophy /></el-icon>
              <div class="link-text">收入排行</div>
              <div class="link-desc">客户和产品收入排行</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="link-card" @click="router.push('/finance/statistics/annual')">
              <el-icon :size="32"><Calendar /></el-icon>
              <div class="link-text">年度统计</div>
              <div class="link-desc">年度收入支出分析</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="link-card" @click="router.push('/finance/statistics/product-revenue')">
              <el-icon :size="32"><Goods /></el-icon>
              <div class="link-text">产品收入</div>
              <div class="link-desc">产品收入趋势分析</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="link-card" @click="router.push('/finance/statistics/new-customers')">
              <el-icon :size="32"><UserFilled /></el-icon>
              <div class="link-text">新客户统计</div>
              <div class="link-desc">新增客户趋势分析</div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 近期数据 -->
      <div class="recent-section">
        <h3>近期订单</h3>
        <el-table :data="recentOrders" v-loading="ordersLoading" style="width: 100%" border>
          <el-table-column prop="order_no" label="订单号" width="170" />
          <el-table-column prop="client_name" label="客户" width="120" />
          <el-table-column prop="product_name" label="产品" min-width="150" />
          <el-table-column prop="amount" label="金额" width="120" align="right">
            <template #default="{ row }">
              <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="时间" width="180" />
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { TrendCharts, Document, User, ChatDotRound, Trophy, Calendar, Goods, UserFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const router = useRouter()
const ordersLoading = ref(false)
const dateRange = ref<[string, string] | null>(null)

const stats = reactive({
  total_income: 0,
  total_orders: 0,
  total_clients: 0,
  pending_tickets: 0,
  income_growth: 0,
  order_growth: 0,
  client_growth: 0
})

const recentOrders = ref<any[]>([])

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const statusTypeMap: Record<number, string> = {
  0: 'warning', 1: 'primary', 2: 'success', 3: 'success', 4: 'info', 5: 'info', 6: 'danger'
}
const statusTextMap: Record<number, string> = {
  0: '待付款', 1: '待审核', 2: '审核通过', 3: '已开通', 4: '已完成', 5: '已取消', 6: '已退款'
}
const getStatusType = (s: number) => (statusTypeMap[s] || 'info') as any
const getStatusText = (s: number) => statusTextMap[s] || '未知'

const fetchStats = async () => {
  try {
    const params: any = {}
    if (dateRange.value) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    const data = await request.get({ url: '/api/admin/reports/dashboard', params })
    if (data) Object.assign(stats, data)
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const fetchRecentOrders = async () => {
  ordersLoading.value = true
  try {
    const data = await request.get({ url: '/api/admin/orders', params: { page: 1, page_size: 5 } })
    recentOrders.value = data.list || []
  } catch (error) {
    console.error('获取近期订单失败:', error)
  } finally {
    ordersLoading.value = false
  }
}

const handleDateChange = () => {
  fetchStats()
}

onMounted(() => {
  fetchStats()
  fetchRecentOrders()
})
</script>

<style scoped lang="scss">
.statistics-overview-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.header-filters { display: flex; gap: 12px; }
.overview-cards { margin-bottom: 24px; }
.stat-card {
  .stat-icon {
    display: flex; align-items: center; justify-content: center;
    width: 60px; height: 60px; border-radius: 12px; margin-bottom: 12px;
    &.income { background: var(--el-color-success-light-9); color: var(--el-color-success); }
    &.orders { background: var(--el-color-primary-light-9); color: var(--el-color-primary); }
    &.clients { background: var(--el-color-warning-light-9); color: var(--el-color-warning); }
    &.tickets { background: var(--el-color-danger-light-9); color: var(--el-color-danger); }
  }
  .stat-info {
    .stat-label { color: var(--el-text-color-secondary); font-size: 14px; margin-bottom: 4px; }
    .stat-value { font-size: 22px; font-weight: 600; margin-bottom: 4px; }
    .stat-compare { font-size: 13px; &.up { color: var(--el-color-success); } &.down { color: var(--el-color-danger); } }
  }
}
.quick-links {
  margin-bottom: 24px;
  h3 { margin: 0 0 16px; font-size: 16px; font-weight: 600; }
}
.link-card {
  text-align: center; padding: 20px; cursor: pointer;
  &:hover { border-color: var(--el-color-primary); }
  .link-text { margin-top: 8px; font-size: 16px; font-weight: 600; }
  .link-desc { margin-top: 4px; color: var(--el-text-color-secondary); font-size: 13px; }
}
.recent-section {
  h3 { margin: 0 0 16px; font-size: 16px; font-weight: 600; }
}
.amount-text { font-weight: 600; color: var(--el-color-primary); }
</style>
