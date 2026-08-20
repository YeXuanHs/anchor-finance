<template>
  <div class="order-list-page">
    <!-- 标签页切换 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane :label="$t('orderList.all')" name="all" />
      <el-tab-pane :label="$t('orderList.pendingPayment')" name="pending_payment" />
      <el-tab-pane :label="$t('orderList.active')" name="active" />
      <el-tab-pane :label="$t('orderList.cancelled')" name="cancelled" />
    </el-tabs>

    <!-- 说明文字和操作按钮 -->
    <div class="page-header">
      <div class="page-desc">
        <el-icon><InfoFilled /></el-icon>
        <span>{{ $t('orderList.description') }}</span>
        <el-button type="primary" link>{{ $t('orderList.helpDoc') }}</el-button>
      </div>
      <div class="page-actions">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          {{ $t('orderList.addOrder') }}
        </el-button>
        <el-button @click="showAdvancedSearch = !showAdvancedSearch">
          {{ $t('orderList.advancedSearch') }}
        </el-button>
      </div>
    </div>

    <!-- 高级搜索区域 -->
    <el-card v-if="showAdvancedSearch" shadow="never" class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item :label="$t('orderList.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('orderList.keywordPlaceholder')" clearable style="width: 200px" @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item :label="$t('orderList.orderType')">
          <el-select v-model="searchForm.type" :placeholder="$t('common.all')" clearable style="width: 120px">
            <el-option :label="$t('orderList.newPurchase')" value="new" />
            <el-option :label="$t('orderList.renewal')" value="renewal" />
            <el-option :label="$t('orderList.upgrade')" value="upgrade" />
            <el-option :label="$t('orderList.refund')" value="refund" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('orderList.orderTime')">
          <el-date-picker v-model="searchForm.date_range" type="daterange" :range-separator="$t('common.to')" :start-placeholder="$t('common.startDate')" :end-placeholder="$t('common.endDate')" style="width: 240px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-table v-loading="loading" :data="tableData" border stripe @sort-change="handleSortChange">
      <el-table-column type="selection" width="50" />
      <el-table-column prop="id" label="ID" width="80" sortable="custom" align="center" />
      <el-table-column prop="client_name" :label="$t('orderList.clientName')" min-width="120">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleViewClient(row)">{{ row.client_name }}</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="product_name" :label="$t('orderList.product')" min-width="150" />
      <el-table-column prop="ip" label="IP" width="130" show-overflow-tooltip />
      <el-table-column prop="created_at" :label="$t('orderList.orderTime')" width="170" sortable="custom" />
      <el-table-column prop="amount" :label="$t('orderList.amount')" width="100" align="right">
        <template #default="{ row }">¥{{ formatMoney(row.amount) }}</template>
      </el-table-column>
      <el-table-column :label="$t('orderList.paymentStatusMethod')" width="150">
        <template #default="{ row }">
          <div>{{ getStatusText(row.status) }}</div>
          <div class="text-gray-400 text-xs">{{ row.payment_method || '-' }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="status" :label="$t('orderList.status')" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="client_remark" :label="$t('orderList.clientRemark')" min-width="120" show-overflow-tooltip />
      <el-table-column :label="$t('orderList.commissionSales')" width="100">
        <template #default="{ row }">{{ row.sales_name || '-' }}</template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size"
        :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange" @current-change="handlePageChange" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { InfoFilled, Plus } from '@element-plus/icons-vue'
import request from '@/utils/http'
import { $t } from '@/locales'

const router = useRouter()
const loading = ref(false)
const tableData = ref([])
const activeTab = ref('all')
const showAdvancedSearch = ref(false)

const searchForm = reactive({
  keyword: '',
  type: '',
  date_range: null as [Date, Date] | null
})

const pagination = reactive({ page: 1, page_size: 100, total: 0 })
const sortParams = reactive({ sort: 'id', order: 'desc' })

const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const getStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    pending_payment: 'warning', pending_activation: 'primary', active: 'success', completed: 'success', cancelled: 'info', refunded: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, () => string> = {
    pending_payment: () => $t('orderList.pendingPayment'),
    pending_activation: () => $t('orderList.pendingActivation'),
    active: () => $t('orderList.active'),
    completed: () => $t('orderList.completed'),
    cancelled: () => $t('orderList.cancelled'),
    refunded: () => $t('orderList.refunded')
  }
  return map[status]?.() || $t('common.unknown')
}

const handleTabChange = () => { pagination.page = 1; fetchList() }

const fetchList = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size, sort: sortParams.sort, order: sortParams.order }
    if (activeTab.value !== 'all') params.status = activeTab.value
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.type) params.type = searchForm.type
    if (searchForm.date_range) {
      params.start_date = searchForm.date_range[0].toISOString().split('T')[0]
      params.end_date = searchForm.date_range[1].toISOString().split('T')[0]
    }
    const data = await request.get({ url: '/api/admin/orders', params })
    tableData.value = data?.list || data || []
    pagination.total = data?.total || 0
  } catch (error) {
    console.error('fetch order list failed:', error)
  } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchList() }
const handleReset = () => { searchForm.keyword = ''; searchForm.type = ''; searchForm.date_range = null; handleSearch() }
const handleSortChange = ({ prop, order }: any) => { sortParams.sort = prop || 'id'; sortParams.order = order === 'ascending' ? 'asc' : 'desc'; fetchList() }
const handleSizeChange = (size: number) => { pagination.page_size = size; pagination.page = 1; fetchList() }
const handlePageChange = (page: number) => { pagination.page = page; fetchList() }
const handleAdd = () => { router.push('/order-create') }
const handleViewClient = (row: any) => { router.push(`/customer-view?id=${row.client_id}`) }

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.order-list-page {
  padding: 16px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}
.page-desc {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: var(--el-color-primary-light-9);
  border-radius: 4px;
  font-size: 14px;
  .el-icon { color: var(--el-color-primary); }
}
.page-actions {
  display: flex;
  gap: 8px;
}
.search-card {
  margin-bottom: 16px;
}
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
