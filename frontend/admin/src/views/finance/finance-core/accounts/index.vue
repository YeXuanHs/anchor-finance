<template>
  <div class="transaction-list-page">
    <!-- 标签页 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane :label="$t('transaction.all')" name="all" />
      <el-tab-pane :label="$t('transaction.recharge')" name="recharge" />
      <el-tab-pane :label="$t('transaction.product')" name="product" />
      <el-tab-pane :label="$t('transaction.renewal')" name="renewal" />
    </el-tabs>

    <!-- 说明文字 -->
    <h5 class="page-title">{{ $t('transaction.description') }}</h5>

    <!-- 操作按钮 -->
    <div class="action-bar">
      <div class="action-left">
        <el-button type="primary" @click="handleAdd">{{ $t('transaction.addTransaction') }}</el-button>
        <el-button @click="showAdvancedSearch = !showAdvancedSearch">{{ $t('transaction.advancedSearch') }}</el-button>
      </div>
    </div>

    <!-- 高级搜索 -->
    <el-card v-if="showAdvancedSearch" shadow="never" class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item :label="$t('transaction.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('transaction.keywordPlaceholder')" clearable style="width: 200px" @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item :label="$t('transaction.timeRange')">
          <el-date-picker v-model="searchForm.dateRange" type="daterange" :range-separator="$t('common.to')" :start-placeholder="$t('common.startDate')" :end-placeholder="$t('common.endDate')" style="width: 240px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-table v-loading="loading" :data="tableData" border stripe>
      <el-table-column prop="client_name" :label="$t('transaction.client')" width="120">
        <template #default="{ row }">
          <el-button type="primary" link @click="$router.push('/customer-view/abstract?id=' + row.client_id)">{{ row.client_name }}</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" :label="$t('transaction.time')" width="170" sortable="custom" />
      <el-table-column prop="payment_method" :label="$t('transaction.paymentMethod')" width="120">
        <template #default="{ row }">{{ row.payment_method || '-' }}</template>
      </el-table-column>
      <el-table-column prop="description" :label="$t('transaction.description')" min-width="200" show-overflow-tooltip />
      <el-table-column prop="amount" :label="$t('transaction.amount')" width="120" align="right" sortable="custom">
        <template #default="{ row }">
          <span :class="row.amount >= 0 ? 'text-green' : 'text-red'">¥{{ formatMoney(row.amount) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="sales_name" :label="$t('transaction.sales')" width="80">
        <template #default="{ row }">{{ row.sales_name || $t('transaction.none') }}</template>
      </el-table-column>
      <el-table-column prop="transaction_no" :label="$t('transaction.transactionNo')" width="200" show-overflow-tooltip />
      <el-table-column prop="type" :label="$t('transaction.type')" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="getTypeTag(row.type)" size="small">{{ getTypeText(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('common.operation')" width="120" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
          <el-button type="danger" link size="small" @click="handleDelete(row)">{{ $t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const tableData = ref([])
const activeTab = ref('all')
const showAdvancedSearch = ref(false)
const searchForm = reactive({ keyword: '', dateRange: null as any })
const pagination = reactive({ page: 1, page_size: 100, total: 0 })

const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const getTypeTag = (type: string): 'success' | 'danger' | 'warning' | 'info' | 'primary' => {
  const map: Record<string, 'success' | 'danger' | 'warning' | 'info' | 'primary'> = { income: 'success', expense: 'danger', refund: 'warning', recharge: 'success', product: 'primary', renewal: 'info' }
  return map[type] || 'info'
}

const getTypeText = (type: string) => {
  const map: Record<string, () => string> = {
    income: () => $t('transaction.income'),
    expense: () => $t('transaction.expense'),
    refund: () => $t('transaction.refund'),
    recharge: () => $t('transaction.recharge'),
    product: () => $t('transaction.product'),
    renewal: () => $t('transaction.renewal')
  }
  return map[type]?.() || type
}

const handleTabChange = () => { pagination.page = 1; fetchList() }
const handleSearch = () => { pagination.page = 1; fetchList() }
const handleReset = () => { searchForm.keyword = ''; searchForm.dateRange = null; pagination.page = 1; fetchList() }
const handleSizeChange = (size: number) => { pagination.page_size = size; pagination.page = 1; fetchList() }
const handlePageChange = (page: number) => { pagination.page = page; fetchList() }
const handleAdd = () => { ElMessage.info($t('transaction.addComingSoon')) }
const handleEdit = (row: any) => { ElMessage.info($t('transaction.editComingSoon')) }

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('transaction.confirmDelete'), $t('common.confirm'), { type: 'warning' })
    await request.del({ url: `/api/admin/transactions/${row.id}` })
    ElMessage.success($t('common.deleteSuccess'))
    fetchList()
  } catch (e) { if (e !== 'cancel') console.error(e) }
}

const fetchList = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (activeTab.value !== 'all') params.type = activeTab.value
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.dateRange) {
      params.start_date = searchForm.dateRange[0]?.toISOString?.()?.split('T')[0] || ''
      params.end_date = searchForm.dateRange[1]?.toISOString?.()?.split('T')[0] || ''
    }
    const data = await request.get({ url: '/api/admin/transactions', params })
    tableData.value = data?.list || data || []
    pagination.total = data?.total || 0
  } catch (error) {
    console.error('fetch transaction list failed:', error)
  } finally { loading.value = false }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.transaction-list-page { padding: 16px; }
.page-title { margin: 0 0 12px; font-size: 14px; font-weight: 500; color: #4e5969; }
.search-card { margin-bottom: 12px; }
.action-bar { display: flex; justify-content: space-between; margin-bottom: 12px; }
.action-left { display: flex; gap: 8px; }
.pagination-wrapper { display: flex; justify-content: flex-end; margin-top: 12px; }
.text-green { color: #36d391; }
.text-red { color: #ef4444; }
</style>
