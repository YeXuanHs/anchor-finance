<template>
  <div class="transaction-list-page">
    <!-- 标签页 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="全部" name="all" />
      <el-tab-pane label="收入" name="income" />
      <el-tab-pane label="支出" name="expense" />
      <el-tab-pane label="退款" name="refund" />
    </el-tabs>

    <!-- 搜索筛选 -->
    <el-card shadow="never" class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="交易号/客户名" clearable style="width: 200px" @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker v-model="searchForm.dateRange" type="daterange" range-separator="至" start-placeholder="开始" end-placeholder="结束" style="width: 240px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch"><el-icon><Search /></el-icon>搜索</el-button>
          <el-button @click="handleReset"><el-icon><Refresh /></el-icon>重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 汇总卡片 -->
    <el-card shadow="never" class="summary-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <div class="summary-item income">
            <div class="summary-label">总收入</div>
            <div class="summary-value">¥{{ formatMoney(summary.total_income) }}</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="summary-item expense">
            <div class="summary-label">总支出</div>
            <div class="summary-value">¥{{ formatMoney(summary.total_expense) }}</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="summary-item refund">
            <div class="summary-label">总退款</div>
            <div class="summary-value">¥{{ formatMoney(summary.total_refund) }}</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="summary-item net">
            <div class="summary-label">净收入</div>
            <div class="summary-value">¥{{ formatMoney(summary.net_income) }}</div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 操作栏 -->
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <div class="action-left">
          <el-button @click="handleExport"><el-icon><Download /></el-icon>导出</el-button>
        </div>
        <div class="action-right">
          <el-button circle @click="fetchList"><el-icon><Refresh /></el-icon></el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe show-summary :summary-method="getSummaries">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="transaction_no" label="交易号" width="150" />
        <el-table-column prop="client_name" label="客户" width="120">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push(`/customer-view/${row.client_id}`)">{{ row.client_name }}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="120" align="right">
          <template #default="{ row }">
            <span :class="row.type === 'income' ? 'text-green' : 'text-red'">
              {{ row.type === 'income' ? '+' : '-' }}¥{{ formatMoney(row.amount) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="balance_after" label="余额" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.balance_after) }}</template>
        </el-table-column>
        <el-table-column prop="gateway" label="支付方式" width="100" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="created_at" label="时间" width="170" />
      </el-table>
      <div class="pagination-wrapper">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10,20,50,100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search, Refresh, Download } from '@element-plus/icons-vue'
import request from '@/utils/http'

const loading = ref(false)
const tableData = ref([])
const activeTab = ref('all')
const searchForm = reactive({ keyword: '', dateRange: null })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const summary = ref<any>({ total_income: 0, total_expense: 0, total_refund: 0, net_income: 0 })

const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const getTypeTag = (type: string) => {
  const map: Record<string, string> = { income: 'success', expense: 'danger', refund: 'warning' }
  return map[type] || 'info'
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = { income: '收入', expense: '支出', refund: '退款' }
  return map[type] || '未知'
}

const handleTabChange = () => { pagination.page = 1; fetchList() }
const handleSearch = () => { pagination.page = 1; fetchList() }
const handleReset = () => { searchForm.keyword = ''; searchForm.dateRange = null; pagination.page = 1; fetchList() }
const handleSizeChange = (size: number) => { pagination.page_size = size; pagination.page = 1; fetchList() }
const handlePageChange = (page: number) => { pagination.page = page; fetchList() }

const fetchList = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (activeTab.value !== 'all') params.type = activeTab.value
    if (searchForm.keyword) params.keyword = searchForm.keyword
    const data = await request.get({ url: '/api/admin/transactions', params })
    tableData.value = data?.list || []
    pagination.total = data?.total || 0
    if (data?.summary) summary.value = data.summary
  } catch (error) {
    console.error('获取交易列表失败:', error)
  } finally {
    loading.value = false
  }
}

const getSummaries = (param: any) => {
  const { columns, data } = param
  const sums: string[] = []
  columns.forEach((column: any, index: number) => {
    if (index === 0) { sums[index] = '合计'; return }
    if (column.property === 'amount') {
      const values = data.map((item: any) => Number(item[column.property] || 0))
      sums[index] = `¥${formatMoney(values.reduce((a: number, b: number) => a + b, 0))}`
    } else { sums[index] = '' }
  })
  return sums
}

const handleExport = async () => {
  try {
    const response = await request.get({ url: '/api/admin/transactions/export', responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `交易记录_${new Date().toISOString().split('T')[0]}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('导出失败:', error)
  }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.transaction-list-page { padding: 16px; }
.search-card { margin-bottom: 16px; :deep(.el-card__body) { padding-bottom: 0; } }
.summary-card { margin-bottom: 16px; }
.summary-item { text-align: center; padding: 10px 0; }
.summary-label { font-size: 14px; color: #86909C; margin-bottom: 8px; }
.summary-value { font-size: 24px; font-weight: 600; }
.income .summary-value { color: #36D391; }
.expense .summary-value { color: #EF4444; }
.refund .summary-value { color: #F59E0B; }
.net .summary-value { color: #4080FF; }
.action-card { margin-bottom: 16px; }
.action-bar { display: flex; justify-content: space-between; align-items: center; }
.table-card { :deep(.el-card__body) { padding: 0; } }
.pagination-wrapper { display: flex; justify-content: flex-end; padding: 16px; }
.text-green { color: #36D391; }
.text-red { color: #EF4444; }
</style>
