<template>
  <div class="renewal-orders-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('renewal.title') }}</span>
          <div class="header-actions">
            <el-button type="primary" @click="handleBatchRenewal">{{ $t('renewal.batchRenewal') }}</el-button>
          </div>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('renewal.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('renewal.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('renewal.pendingRenewal')" :value="0" />
            <el-option :label="$t('renewal.renewed')" :value="1" />
            <el-option :label="$t('renewal.expired')" :value="2" />
            <el-option :label="$t('common.cancelled')" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('renewal.dueDateRange')">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            :range-separator="$t('common.to')"
            :start-placeholder="$t('common.startDate')"
            :end-placeholder="$t('common.endDate')"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="order_no" :label="$t('renewal.originalOrderNo')" width="170" />
        <el-table-column prop="client_name" :label="$t('renewal.client')" width="120" />
        <el-table-column prop="product_name" :label="$t('renewal.product')" min-width="160" show-overflow-tooltip />
        <el-table-column prop="amount" :label="$t('renewal.renewalAmount')" width="120" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="billing_cycle" :label="$t('renewal.cycle')" width="100" />
        <el-table-column prop="due_date" :label="$t('renewal.dueDate')" width="170" />
        <el-table-column prop="status" :label="$t('common.status')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.action')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 0"
              type="success"
              link
              @click="handleRenew(row)"
            >
              {{ $t('renewal.renew') }}
            </el-button>
            <el-button type="primary" link @click="handleViewDetail(row)">{{ $t('common.detail') }}</el-button>
            <el-popconfirm
              v-if="row.status === 0"
              :title="$t('renewal.confirmCancel')"
              @confirm="handleCancel(row)"
            >
              <template #reference>
                <el-button type="danger" link>{{ $t('common.cancel') }}</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" :title="$t('renewal.detailTitle')" width="650px">
      <el-descriptions :column="2" border v-if="detailData">
        <el-descriptions-item :label="$t('renewal.originalOrderNo')">{{ detailData.order_no }}</el-descriptions-item>
        <el-descriptions-item :label="$t('renewal.client')">{{ detailData.client_name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('renewal.productLabel')" :span="2">{{ detailData.product_name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('renewal.renewalAmount')">
          <span class="amount-text">¥{{ formatAmount(detailData.amount) }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('renewal.billingCycle')">{{ detailData.billing_cycle }}</el-descriptions-item>
        <el-descriptions-item :label="$t('renewal.dueDate')">{{ detailData.due_date }}</el-descriptions-item>
        <el-descriptions-item :label="$t('common.status')">
          <el-tag :type="getStatusType(detailData.status)" size="small">
            {{ getStatusText(detailData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('common.createdAt')" :span="2">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('common.remark')" :span="2">{{ detailData.remark || $t('renewal.none') }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">{{ $t('common.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

const loading = ref(false)
const detailVisible = ref(false)
const detailData = ref<any>(null)
const selectedRows = ref<any[]>([])

const searchForm = reactive({
  keyword: '',
  status: undefined as number | undefined,
  date_range: [] as string[]
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const statusMap: Record<number, { text: () => string; type: string }> = {
  0: { text: () => $t('renewal.pendingRenewal'), type: 'warning' },
  1: { text: () => $t('renewal.renewed'), type: 'success' },
  2: { text: () => $t('renewal.expired'), type: 'danger' },
  3: { text: () => $t('common.cancelled'), type: 'info' }
}

const getStatusText = (status: number) => statusMap[status]?.text() || $t('common.unknown')
const getStatusType = (status: number) => (statusMap[status]?.type || 'info') as any

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      keyword: searchForm.keyword || undefined,
      status: searchForm.status
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({
      url: '/api/admin/multi-renew',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error($t('renewal.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = undefined
  searchForm.date_range = []
  handleSearch()
}

const handleSelectionChange = (rows: any[]) => {
  selectedRows.value = rows
}

const handleViewDetail = (row: any) => {
  detailData.value = row
  detailVisible.value = true
}

const handleRenew = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('renewal.confirmRenewOrder', { orderNo: row.order_no, amount: '¥' + formatAmount(row.amount) }), $t('renewal.renewConfirm'))
    await request.post({
      url: `/api/admin/multi-renew/${row.id}/execute`
    })
    ElMessage.success($t('renewal.renewSuccess'))
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error($t('renewal.renewFailed'))
    }
  }
}

const handleCancel = async (row: any) => {
  try {
    await request.post({
      url: `/api/admin/multi-renew/${row.id}/cancel`
    })
    ElMessage.success($t('common.cancelled'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('common.operateFailed'))
  }
}

const handleBatchRenewal = async () => {
  if (!selectedRows.value.length) {
    ElMessage.warning($t('renewal.selectOrderFirst'))
    return
  }
  try {
    await ElMessageBox.confirm($t('renewal.confirmBatchRenew', { count: selectedRows.value.length }), $t('renewal.batchRenewConfirm'))
    const ids = selectedRows.value.map(r => r.id)
    await request.post({
      url: '/api/admin/multi-renew',
      params: { ids }
    })
    ElMessage.success($t('renewal.batchRenewSuccess'))
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error($t('renewal.batchRenewFailed'))
    }
  }
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchData()
}

const handlePageChange = () => {
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.renewal-orders-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.amount-text {
  font-weight: 600;
  color: var(--el-color-primary);
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
