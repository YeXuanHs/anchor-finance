<template>
  <div class="reports-page page-container">
    <div class="report-tabs">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="收入报表" name="revenue" />
        <el-tab-pane label="用户报表" name="users" />
        <el-tab-pane label="订单报表" name="orders" />
        <el-tab-pane label="产品报表" name="products" />
      </el-tabs>
    </div>

    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="时间范围">
          <el-date-picker v-model="searchForm.date_range" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="exportReport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="stats-grid">
      <div class="stat-card" v-for="stat in stats" :key="stat.label">
        <div class="stat-label">{{ stat.label }}</div>
        <div class="stat-value">{{ stat.value }}</div>
        <div class="stat-change" :class="stat.changeType">{{ stat.change }}</div>
      </div>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>{{ getTabTitle() }}</h3>
      </div>
      <el-table :data="reportData" style="width: 100%" v-loading="loading" show-summary>
        <el-table-column v-for="col in columns" :key="col.prop" :prop="col.prop" :label="col.label" :width="col.width" />
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Download } from '@element-plus/icons-vue'
import request from '@/utils/request'

const activeTab = ref('revenue')
const loading = ref(false)
const reportData = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const searchForm = ref({ date_range: [] })

const stats = ref([
  { label: '总收入', value: '¥0.00', change: '+0%', changeType: 'up' },
  { label: '订单数', value: '0', change: '+0%', changeType: 'up' },
  { label: '新增用户', value: '0', change: '+0%', changeType: 'up' },
  { label: '活跃用户', value: '0', change: '+0%', changeType: 'up' }
])

const columnsMap: Record<string, any[]> = {
  revenue: [
    { prop: 'date', label: '日期', width: 120 },
    { prop: 'income', label: '收入' },
    { prop: 'refund', label: '退款' },
    { prop: 'net', label: '净收入' },
    { prop: 'order_count', label: '订单数' }
  ],
  users: [
    { prop: 'date', label: '日期', width: 120 },
    { prop: 'new_users', label: '新增用户' },
    { prop: 'active_users', label: '活跃用户' },
    { prop: 'total_users', label: '总用户数' }
  ],
  orders: [
    { prop: 'date', label: '日期', width: 120 },
    { prop: 'total_orders', label: '总订单数' },
    { prop: 'completed', label: '已完成' },
    { prop: 'cancelled', label: '已取消' },
    { prop: 'amount', label: '订单金额' }
  ],
  products: [
    { prop: 'name', label: '产品名称' },
    { prop: 'sales_count', label: '销量' },
    { prop: 'revenue', label: '收入' },
    { prop: 'refund_rate', label: '退款率' }
  ]
}

const columns = computed(() => columnsMap[activeTab.value] || [])

const getTabTitle = () => {
  const map: Record<string, string> = { revenue: '收入明细', users: '用户统计', orders: '订单统计', products: '产品统计' }
  return map[activeTab.value] || ''
}

const buildDateParams = () => {
  const params: any = { page: currentPage.value, page_size: pageSize.value }
  const range = searchForm.value.date_range
  if (range?.length === 2) {
    params.start_date = range[0]
    params.end_date = range[1]
  }
  return params
}

const fetchDashboard = async () => {
  try {
    const { data } = await request.get('/admin/reports/dashboard')
    if (data?.data) {
      const d = data.data
      stats.value = [
        { label: '总收入', value: `¥${d.total_revenue?.toLocaleString() || '0.00'}`, change: `${d.revenue_change || '+0%'}`, changeType: d.revenue_change?.startsWith('-') ? 'down' : 'up' },
        { label: '订单数', value: `${d.total_orders || 0}`, change: `${d.orders_change || '+0%'}`, changeType: d.orders_change?.startsWith('-') ? 'down' : 'up' },
        { label: '新增用户', value: `${d.new_users || 0}`, change: `${d.users_change || '+0%'}`, changeType: d.users_change?.startsWith('-') ? 'down' : 'up' },
        { label: '活跃用户', value: `${d.active_users || 0}`, change: `${d.active_change || '+0%'}`, changeType: d.active_change?.startsWith('-') ? 'down' : 'up' }
      ]
    }
  } catch {}
}

const fetchReportData = async () => {
  loading.value = true
  try {
    const endpoint = `/admin/reports/${activeTab.value}`
    const { data } = await request.get(endpoint, { params: buildDateParams() })
    if (data?.data) {
      reportData.value = data.data.items || data.data
      total.value = data.data.total || 0
    }
  } catch {
    ElMessage.error('获取报表数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchReportData()
}

const exportReport = async () => {
  try {
    const params = { ...buildDateParams(), type: activeTab.value }
    const { data } = await request.get('/admin/reports/export', { params, responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([data]))
    const link = document.createElement('a')
    link.href = url
    link.download = `报表_${activeTab.value}_${new Date().toISOString().slice(0, 10)}.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch {
    ElMessage.error('导出失败')
  }
}

watch(activeTab, () => {
  currentPage.value = 1
  fetchReportData()
})

fetchDashboard()
fetchReportData()
</script>

<style scoped lang="scss">
.reports-page {
  .report-tabs { margin-bottom: 20px; }
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
    margin-bottom: 20px;
  }
  .stat-card {
    background: var(--bg-card);
    border-radius: var(--border-radius);
    padding: 20px;
    box-shadow: var(--shadow-sm);
    .stat-label { font-size: 13px; color: var(--text-secondary); margin-bottom: 8px; }
    .stat-value { font-size: 24px; font-weight: 600; color: var(--text-primary); }
    .stat-change { font-size: 12px; margin-top: 4px; &.up { color: var(--success-color); } &.down { color: var(--danger-color); } }
  }
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
}
</style>
