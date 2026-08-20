<template>
  <div class="affiliate-records-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.affiliateRecords.title') }}</span>
          <el-button type="primary" @click="handleExport">
            <el-icon><Download /></el-icon>
            {{ $t('finance.affiliateRecords.export') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('finance.affiliateRecords.affiliateUser')">
          <el-input v-model="searchForm.affiliate_username" :placeholder="$t('finance.affiliateRecords.affiliateUserPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('finance.affiliateRecords.sourceUser')">
          <el-input v-model="searchForm.client_username" :placeholder="$t('finance.affiliateRecords.sourceUserPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('finance.affiliateRecords.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('finance.affiliateRecords.all')" clearable>
            <el-option :label="$t('finance.affiliateRecords.statusPending')" :value="0" />
            <el-option :label="$t('finance.affiliateRecords.statusConfirmed')" :value="1" />
            <el-option :label="$t('finance.affiliateRecords.statusRejected')" :value="2" />
            <el-option :label="$t('finance.affiliateRecords.statusWithdrawn')" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('finance.affiliateRecords.dateRange')">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            :range-separator="$t('finance.affiliateRecords.rangeSeparator')"
            :start-placeholder="$t('finance.affiliateRecords.startDate')"
            :end-placeholder="$t('finance.affiliateRecords.endDate')"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('finance.affiliateRecords.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('finance.affiliateRecords.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 汇总统计 -->
      <el-row :gutter="20" class="summary-row">
        <el-col :span="6">
          <div class="summary-card">
            <div class="summary-label">{{ $t('finance.affiliateRecords.totalCommission') }}</div>
            <div class="summary-value total">¥{{ formatAmount(summary.total_commission) }}</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="summary-card">
            <div class="summary-label">{{ $t('finance.affiliateRecords.pendingCommission') }}</div>
            <div class="summary-value pending">¥{{ formatAmount(summary.pending_commission) }}</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="summary-card">
            <div class="summary-label">{{ $t('finance.affiliateRecords.confirmedCommission') }}</div>
            <div class="summary-value confirmed">¥{{ formatAmount(summary.confirmed_commission) }}</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="summary-card">
            <div class="summary-label">{{ $t('finance.affiliateRecords.withdrawnCommission') }}</div>
            <div class="summary-value withdrawn">¥{{ formatAmount(summary.withdrawn_commission) }}</div>
          </div>
        </el-col>
      </el-row>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="affiliate_username" :label="$t('finance.affiliateRecords.affiliateUser')" width="120" />
        <el-table-column prop="client_username" :label="$t('finance.affiliateRecords.sourceUser')" width="120" />
        <el-table-column prop="order_no" :label="$t('finance.affiliateRecords.relatedOrder')" width="170">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewOrder(row)">
              {{ row.order_no }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="order_amount" :label="$t('finance.affiliateRecords.orderAmount')" width="110" align="right">
          <template #default="{ row }">
            ¥{{ formatAmount(row.order_amount) }}
          </template>
        </el-table-column>
        <el-table-column prop="commission_rate" :label="$t('finance.affiliateRecords.commissionRate')" width="100" align="center">
          <template #default="{ row }">
            <span class="commission-text">{{ row.commission_rate }}%</span>
          </template>
        </el-table-column>
        <el-table-column prop="commission_amount" :label="$t('finance.affiliateRecords.commissionAmount')" width="110" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.commission_amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('finance.affiliateRecords.status')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('finance.affiliateRecords.createdAt')" width="170" />
        <el-table-column :label="$t('finance.affiliateRecords.actions')" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">{{ $t('finance.affiliateRecords.detail') }}</el-button>
            <el-button
              v-if="row.status === 0"
              type="success"
              link
              @click="handleConfirm(row)"
            >
              {{ $t('finance.affiliateRecords.confirm') }}
            </el-button>
            <el-button
              v-if="row.status === 0"
              type="danger"
              link
              @click="handleReject(row)"
            >
              {{ $t('finance.affiliateRecords.reject') }}
            </el-button>
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
    <el-dialog v-model="detailVisible" :title="$t('finance.affiliateRecords.detailTitle')" width="650px" destroy-on-close>
      <el-descriptions :column="2" border>
        <el-descriptions-item :label="$t('finance.affiliateRecords.recordId')">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('finance.affiliateRecords.status')">
          <el-tag :type="getStatusType(detailData.status)" size="small">
            {{ getStatusText(detailData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('finance.affiliateRecords.affiliateUser')">{{ detailData.affiliate_username }}</el-descriptions-item>
        <el-descriptions-item :label="$t('finance.affiliateRecords.sourceUser')">{{ detailData.client_username }}</el-descriptions-item>
        <el-descriptions-item :label="$t('finance.affiliateRecords.relatedOrder')">{{ detailData.order_no }}</el-descriptions-item>
        <el-descriptions-item :label="$t('finance.affiliateRecords.orderAmount')">
          ¥{{ formatAmount(detailData.order_amount) }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('finance.affiliateRecords.commissionRate')">
          <span class="commission-text">{{ detailData.commission_rate }}%</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('finance.affiliateRecords.commissionAmount')">
          <span class="amount-text">¥{{ formatAmount(detailData.commission_amount) }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('finance.affiliateRecords.createdAt')">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('finance.affiliateRecords.confirmedAt')">{{ detailData.confirmed_at || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('finance.affiliateRecords.notes')" :span="2">{{ detailData.remark || $t('finance.affiliateRecords.none') }}</el-descriptions-item>
      </el-descriptions>

      <!-- 确认/拒绝操作 -->
      <div v-if="detailData.status === 0" class="action-section">
        <el-divider content-position="left">{{ $t('finance.affiliateRecords.auditAction') }}</el-divider>
        <el-input
          v-model="confirmRemark"
          type="textarea"
          :rows="3"
          :placeholder="$t('finance.affiliateRecords.auditRemarkPlaceholder')"
        />
        <div class="action-buttons">
          <el-button type="success" @click="handleConfirmFromDetail" :loading="actionLoading">
            {{ $t('finance.affiliateRecords.confirmCommission') }}
          </el-button>
          <el-button type="danger" @click="handleRejectFromDetail" :loading="actionLoading">
            {{ $t('finance.affiliateRecords.rejectCommission') }}
          </el-button>
        </div>
      </div>

      <template #footer>
        <el-button @click="detailVisible = false">{{ $t('finance.affiliateRecords.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { exportToCSV } from '@/utils/export'
import { $t } from '@/locales'

// 加载状态
const loading = ref(false)
const actionLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  affiliate_username: '',
  client_username: '',
  status: undefined as number | undefined,
  date_range: [] as string[]
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref<any[]>([])

// 汇总统计
const summary = reactive({
  total_commission: 0,
  pending_commission: 0,
  confirmed_commission: 0,
  withdrawn_commission: 0
})

// 详情对话框
const detailVisible = ref(false)
const detailData = ref<any>({})
const confirmRemark = ref('')

// 状态映射
const STATUS_MAP: Record<number, { text: string; type: string }> = {
  0: { text: $t('finance.affiliateRecords.statusPending'), type: 'warning' },
  1: { text: $t('finance.affiliateRecords.statusConfirmed'), type: 'success' },
  2: { text: $t('finance.affiliateRecords.statusRejected'), type: 'danger' },
  3: { text: $t('finance.affiliateRecords.statusWithdrawn'), type: 'info' }
}

// 获取状态文本
const getStatusText = (status: number) => {
  return STATUS_MAP[status]?.text || $t('finance.affiliateRecords.unknown')
}

// 获取状态类型
const getStatusType = (status: number) => {
  return (STATUS_MAP[status]?.type || 'info') as any
}

// 格式化金额
const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

// 获取佣金记录列表
const fetchRecords = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      affiliate_username: searchForm.affiliate_username || undefined,
      client_username: searchForm.client_username || undefined,
      status: searchForm.status
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({
      url: '/api/admin/affiliate/user-affi-record',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
    if (data.summary) {
      Object.assign(summary, data.summary)
    }
  } catch (error) {
    console.error('获取佣金记录失败:', error)
    ElMessage.error($t('finance.affiliateRecords.msgFetchFailed'))
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchRecords()
}

// 重置
const handleReset = () => {
  searchForm.affiliate_username = ''
  searchForm.client_username = ''
  searchForm.status = undefined
  searchForm.date_range = []
  handleSearch()
}

// 查看详情
const handleViewDetail = (row: any) => {
  detailData.value = { ...row }
  confirmRemark.value = ''
  detailVisible.value = true
}

// 查看关联订单
const handleViewOrder = (row: any) => {
  window.open(`/finance/orders/list?order_no=${row.order_no}`, '_blank')
}

// 确认佣金
const handleConfirm = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('finance.affiliateRecords.msgConfirmCommission'), $t('finance.affiliateRecords.confirmCommissionTitle'), {
      confirmButtonText: $t('finance.affiliateRecords.confirm'),
      cancelButtonText: $t('finance.affiliateRecords.close'),
      type: 'warning'
    })
    await request.post({
      url: `/api/admin/affiliate/records/${row.id}/confirm`,
      params: { status: 1 }
    })
    ElMessage.success($t('finance.affiliateRecords.msgConfirmSuccess'))
    fetchRecords()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error($t('finance.affiliateRecords.msgOperationFailed'))
    }
  }
}

// 拒绝佣金
const handleReject = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('finance.affiliateRecords.msgRejectCommission'), $t('finance.affiliateRecords.rejectCommissionTitle'), {
      confirmButtonText: $t('finance.affiliateRecords.confirm'),
      cancelButtonText: $t('finance.affiliateRecords.close'),
      type: 'warning'
    })
    await request.post({
      url: `/api/admin/affiliate/records/${row.id}/confirm`,
      params: { status: 2 }
    })
    ElMessage.success($t('finance.affiliateRecords.msgRejectSuccess'))
    fetchRecords()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error($t('finance.affiliateRecords.msgOperationFailed'))
    }
  }
}

// 从详情页确认
const handleConfirmFromDetail = async () => {
  actionLoading.value = true
  try {
    await request.post({
      url: `/api/admin/affiliate/records/${detailData.value.id}/confirm`,
      params: {
        status: 1,
        remark: confirmRemark.value || undefined
      }
    })
    ElMessage.success($t('finance.affiliateRecords.msgConfirmSuccess'))
    detailVisible.value = false
    fetchRecords()
  } catch (error) {
    ElMessage.error($t('finance.affiliateRecords.msgOperationFailed'))
  } finally {
    actionLoading.value = false
  }
}

// 从详情页拒绝
const handleRejectFromDetail = async () => {
  actionLoading.value = true
  try {
    await request.post({
      url: `/api/admin/affiliate/records/${detailData.value.id}/confirm`,
      params: {
        status: 2,
        remark: confirmRemark.value || undefined
      }
    })
    ElMessage.success($t('finance.affiliateRecords.msgRejectSuccess'))
    detailVisible.value = false
    fetchRecords()
  } catch (error) {
    ElMessage.error($t('finance.affiliateRecords.msgOperationFailed'))
  } finally {
    actionLoading.value = false
  }
}

// 导出
const handleExport = async () => {
  try {
    const data = await request.get({ url: '/api/admin/affiliate/user-affi-record', params: { page: 1, page_size: 9999 } })
    const list = data.list || data || []
    exportToCSV(list, [
      { key: 'id', title: 'ID' },
      { key: 'user_id', title: '用户ID' },
      { key: 'amount', title: '金额' },
      { key: 'status', title: '状态' },
      { key: 'created_at', title: '时间' }
    ], '推广记录')
    ElMessage.success($t('finance.affiliateRecords.msgExportSuccess'))
  } catch { ElMessage.error($t('finance.affiliateRecords.msgExportFailed')) }
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchRecords()
}

// 页码变化
const handlePageChange = () => {
  fetchRecords()
}

onMounted(() => {
  fetchRecords()
})
</script>

<style scoped lang="scss">
.affiliate-records-page {
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

.summary-row {
  margin-bottom: 20px;
}

.summary-card {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 16px;
  text-align: center;

  .summary-label {
    font-size: 13px;
    color: #909399;
    margin-bottom: 8px;
  }

  .summary-value {
    font-size: 22px;
    font-weight: 600;

    &.total {
      color: #303133;
    }

    &.pending {
      color: #e6a23c;
    }

    &.confirmed {
      color: #67c23a;
    }

    &.withdrawn {
      color: #409eff;
    }
  }
}

.amount-text {
  font-weight: 600;
  color: #f56c6c;
}

.commission-text {
  font-weight: 600;
  color: #409eff;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.action-section {
  margin-top: 16px;

  .action-buttons {
    margin-top: 12px;
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }
}
</style>
