<template>
  <div class="sales-statistics-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsSalesStatistics.title') }}</span>
          <el-button type="success" @click="handleExport">
            <el-icon><Download /></el-icon>
            {{ $t('clientsSalesStatistics.exportReport') }}
          </el-button>
        </div>
      </template>

      <el-row :gutter="16" class="stats-row">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">&yen;{{ formatAmount(summary.today_sales) }}</div>
            <div class="stat-label">{{ $t('clientsSalesStatistics.todaySales') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">&yen;{{ formatAmount(summary.month_sales) }}</div>
            <div class="stat-label">{{ $t('clientsSalesStatistics.monthSales') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ summary.today_orders || 0 }}</div>
            <div class="stat-label">{{ $t('clientsSalesStatistics.todayOrders') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ summary.month_orders || 0 }}</div>
            <div class="stat-label">{{ $t('clientsSalesStatistics.monthOrders') }}</div>
          </el-card>
        </el-col>
      </el-row>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('common.dateRange')">
          <el-date-picker v-model="searchForm.date_range" type="daterange" :range-separator="$t('common.to')" :start-placeholder="$t('common.startDate')" :end-placeholder="$t('common.endDate')" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="date" :label="$t('clientsSalesStatistics.date')" width="130" align="center" />
        <el-table-column prop="order_count" :label="$t('clientsSalesStatistics.orderCount')" width="110" align="center" />
        <el-table-column prop="sales_amount" :label="$t('clientsSalesStatistics.salesAmount')" width="140" align="right">
          <template #default="{ row }">&yen;{{ formatAmount(row.sales_amount) }}</template>
        </el-table-column>
        <el-table-column prop="new_clients" :label="$t('clientsSalesStatistics.newClients')" width="110" align="center" />
        <el-table-column prop="refund_amount" :label="$t('clientsSalesStatistics.refundAmount')" width="130" align="right">
          <template #default="{ row }">&yen;{{ formatAmount(row.refund_amount) }}</template>
        </el-table-column>
        <el-table-column prop="net_income" :label="$t('clientsSalesStatistics.netIncome')" width="140" align="right">
          <template #default="{ row }"><span class="text-green">&yen;{{ formatAmount(row.net_income) }}</span></template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const searchForm = reactive({ date_range: [] as string[] })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const summary = ref<any>({})

const formatAmount = (amount: number | undefined) => amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (searchForm.date_range?.length === 2) { params.start_date = searchForm.date_range[0]; params.end_date = searchForm.date_range[1] }
    const data = await request.get({ url: '/api/admin/sales/statistics', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
    if (data.summary) summary.value = data.summary
  } catch (e) { ElMessage.error($t('common.fetchFailed')) } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.date_range = []; handleSearch() }
const handleExport = async () => {
  try {
    await request.get({ url: '/api/admin/sales/export', params: { ...searchForm } })
    ElMessage.success($t('common.exportSuccess'))
  } catch (e) { ElMessage.error($t('common.exportFailed')) }
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }
onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.sales-statistics-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.stats-row { margin-bottom: 20px; }
.stat-card { text-align: center; .stat-value { font-size: 24px; font-weight: 700; color: var(--el-color-primary); margin-bottom: 8px; } .stat-label { font-size: 13px; color: var(--el-text-color-secondary); } }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.text-green { color: var(--el-color-success); }
</style>
