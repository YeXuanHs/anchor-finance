<template>
  <div class="withdraw-page">
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane :label="$t('withdraw.all')" name="all" />
      <el-tab-pane :label="$t('withdraw.pending')" name="pending" />
      <el-tab-pane :label="$t('withdraw.approved')" name="approved" />
      <el-tab-pane :label="$t('withdraw.rejected')" name="rejected" />
      <el-tab-pane :label="$t('withdraw.completed')" name="completed" />
    </el-tabs>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="withdraw_no" :label="$t('withdraw.withdrawNo')" width="150" />
        <el-table-column prop="client_name" :label="$t('withdraw.client')" width="120">
          <template #default="{ row }"><el-button type="primary" link @click="$router.push(`/customer-view/${row.client_id}`)">{{ row.client_name }}</el-button></template>
        </el-table-column>
        <el-table-column prop="amount" :label="$t('withdraw.amount')" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.amount) }}</template>
        </el-table-column>
        <el-table-column prop="fee" :label="$t('withdraw.fee')" width="100" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.fee) }}</template>
        </el-table-column>
        <el-table-column prop="actual_amount" :label="$t('withdraw.actualAmount')" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.actual_amount) }}</template>
        </el-table-column>
        <el-table-column prop="withdraw_method" :label="$t('withdraw.withdrawMethod')" width="100" />
        <el-table-column prop="account_info" :label="$t('withdraw.accountInfo')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" :label="$t('withdraw.status')" width="100" align="center">
          <template #default="{ row }"><el-tag :type="getStatusType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('withdraw.applyTime')" width="170" />
        <el-table-column :label="$t('withdraw.operations')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status === 'pending'" type="success" link size="small" @click="handleApprove(row)">{{ $t('withdraw.approve') }}</el-button>
            <el-button v-if="row.status === 'pending'" type="danger" link size="small" @click="handleReject(row)">{{ $t('withdraw.reject') }}</el-button>
            <el-button v-if="row.status === 'approved'" type="primary" link size="small" @click="handleComplete(row)">{{ $t('withdraw.markComplete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>
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
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const getStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { pending: 'warning', approved: 'primary', rejected: 'danger', completed: 'success' }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, () => string> = { pending: () => $t('withdraw.pending'), approved: () => $t('withdraw.approved'), rejected: () => $t('withdraw.rejected'), completed: () => $t('withdraw.completed') }
  return map[status]?.() || $t('common.unknown')
}

const handleTabChange = () => { pagination.page = 1; fetchList() }

const fetchList = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (activeTab.value !== 'all') params.status = activeTab.value
    const data = await request.get({ url: '/api/admin/withdrawals', params })
    tableData.value = data?.list || []; pagination.total = data?.total || 0
  } catch (error) { console.error('fetch withdraw list failed:', error) } finally { loading.value = false }
}

const handleSizeChange = (size: number) => { pagination.page_size = size; pagination.page = 1; fetchList() }
const handlePageChange = (page: number) => { pagination.page = page; fetchList() }

const handleApprove = async (row: any) => {
  try { await ElMessageBox.confirm($t('withdraw.confirmApprove'), $t('withdraw.confirmApproveTitle'), { type: 'warning' }); await request.post({ url: `/api/admin/withdrawals/${row.id}/approve` }); ElMessage.success($t('withdraw.approvedSuccess')); fetchList() } catch (error) { if (error !== 'cancel') console.error('approve failed:', error) }
}

const handleReject = async (row: any) => {
  try {
    const { value: reason } = await ElMessageBox.prompt($t('withdraw.enterRejectReason'), $t('withdraw.rejectTitle'), { confirmButtonText: $t('common.confirm'), cancelButtonText: $t('common.cancel'), inputValidator: (v) => !!v || $t('withdraw.enterRejectReason') })
    await request.post({ url: `/api/admin/withdrawals/${row.id}/reject`, data: { reason } }); ElMessage.success($t('withdraw.rejectedSuccess')); fetchList()
  } catch (error) { if (error !== 'cancel') console.error('reject failed:', error) }
}

const handleComplete = async (row: any) => {
  try { await ElMessageBox.confirm($t('withdraw.confirmComplete'), $t('withdraw.confirmCompleteTitle'), { type: 'warning' }); await request.post({ url: `/api/admin/withdrawals/${row.id}/complete` }); ElMessage.success($t('withdraw.completedSuccess')); fetchList() } catch (error) { if (error !== 'cancel') console.error('complete failed:', error) }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.withdraw-page { padding: 16px; }
.table-card { :deep(.el-card__body) { padding: 0; } }
.pagination-wrapper { display: flex; justify-content: flex-end; padding: 16px; }
</style>
