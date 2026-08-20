<template>
  <div class="upgrades-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('upgrades.pageTitle') }}</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('upgrades.orderId')">
          <el-input v-model="searchForm.order_id" :placeholder="$t('upgrades.orderId')" clearable />
        </el-form-item>
        <el-form-item :label="$t('upgrades.customer')">
          <el-input v-model="searchForm.username" :placeholder="$t('upgrades.customerUsername')" clearable />
        </el-form-item>
        <el-form-item :label="$t('upgrades.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('upgrades.all')" clearable>
            <el-option :label="$t('upgrades.pending')" value="pending" />
            <el-option :label="$t('upgrades.awaitingPayment')" value="awaiting_payment" />
            <el-option :label="$t('upgrades.completed')" value="completed" />
            <el-option :label="$t('upgrades.cancelled')" value="cancelled" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('upgrades.upgradeType')">
          <el-select v-model="searchForm.type" :placeholder="$t('upgrades.all')" clearable>
            <el-option :label="$t('upgrades.productUpgrade')" value="product" />
            <el-option :label="$t('upgrades.configUpgrade')" value="configoption" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('upgrades.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('upgrades.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 统计卡片 -->
      <el-row :gutter="20" class="stat-section">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.total }}</div>
            <div class="stat-label">{{ $t('upgrades.totalUpgrades') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.pending }}</div>
            <div class="stat-label">{{ $t('upgrades.pending') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.awaiting_payment }}</div>
            <div class="stat-label">{{ $t('upgrades.awaitingPayment') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.completed }}</div>
            <div class="stat-label">{{ $t('upgrades.completed') }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 升级列表 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="order_id" :label="$t('upgrades.orderId')" width="120">
          <template #default="{ row }">
            <el-link type="primary" @click="goToOrder(row.order_id)">{{ row.order_id }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="username" :label="$t('upgrades.customer')" width="120" />
        <el-table-column prop="type" :label="$t('upgrades.upgradeType')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'product' ? 'primary' : 'warning'" size="small">
              {{ row.type === 'product' ? $t('upgrades.productUpgrade') : $t('upgrades.configUpgrade') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="product_name" :label="$t('upgrades.productName')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="original_value" :label="$t('upgrades.originalConfig')" min-width="120" show-overflow-tooltip />
        <el-table-column prop="upgrade_value" :label="$t('upgrades.upgradeConfig')" min-width="120" show-overflow-tooltip />
        <el-table-column prop="amount" :label="$t('upgrades.upgradeFee')" width="120" align="right">
          <template #default="{ row }">
            <span class="amount">{{ row.currency_prefix }}{{ row.amount }}{{ row.currency_suffix }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('upgrades.status')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('upgrades.createdAt')" width="170" />
        <el-table-column prop="completed_at" :label="$t('upgrades.completedAt')" width="170">
          <template #default="{ row }">{{ row.completed_at || '-' }}</template>
        </el-table-column>
        <el-table-column :label="$t('upgrades.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleDetail(row)">{{ $t('upgrades.detail') }}</el-button>
            <el-button v-if="row.status === 'pending'" type="success" link @click="handleApprove(row)">{{ $t('upgrades.approve') }}</el-button>
            <el-button v-if="row.status === 'pending'" type="warning" link @click="handleCancel(row)">{{ $t('upgrades.cancel') }}</el-button>
            <el-button v-if="row.status === 'awaiting_payment'" type="info" link @click="handleMarkPaid(row)">{{ $t('upgrades.markPaid') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialogVisible" :title="$t('upgrades.upgradeDetail')" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item :label="$t('upgrades.upgradeId')">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('upgrades.orderId')">{{ detailData.order_id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('upgrades.customer')">{{ detailData.username }}</el-descriptions-item>
        <el-descriptions-item :label="$t('upgrades.upgradeType')">{{ detailData.type === 'product' ? $t('upgrades.productUpgrade') : $t('upgrades.configUpgrade') }}</el-descriptions-item>
        <el-descriptions-item :label="$t('upgrades.productName')" :span="2">{{ detailData.product_name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('upgrades.originalConfig')">{{ detailData.original_value }}</el-descriptions-item>
        <el-descriptions-item :label="$t('upgrades.upgradeConfig')">{{ detailData.upgrade_value }}</el-descriptions-item>
        <el-descriptions-item :label="$t('upgrades.upgradeFee')">
          {{ detailData.currency_prefix }}{{ detailData.amount }}{{ detailData.currency_suffix }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('upgrades.status')">
          <el-tag :type="getStatusType(detailData.status)">{{ getStatusLabel(detailData.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('upgrades.createdAt')">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('upgrades.completedAt')">{{ detailData.completed_at || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('upgrades.remark')" :span="2">{{ detailData.remark || '-' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">{{ $t('upgrades.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const router = useRouter()

const statusMap: Record<string, { label: () => string; type: string }> = {
  pending: { label: () => $t('upgrades.pending'), type: 'warning' },
  awaiting_payment: { label: () => $t('upgrades.awaitingPayment'), type: 'primary' },
  completed: { label: () => $t('upgrades.completed'), type: 'success' },
  cancelled: { label: () => $t('upgrades.cancelled'), type: 'info' }
}

const getStatusType = (status: string) => (statusMap[status]?.type || 'info') as any
const getStatusLabel = (status: string) => statusMap[status]?.label() || status

const loading = ref(false)
const detailDialogVisible = ref(false)

const searchForm = reactive({ order_id: '', username: '', status: '', type: '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const stats = reactive({ total: 0, pending: 0, awaiting_payment: 0, completed: 0 })
const detailData = ref<any>({})

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (searchForm.order_id) params.order_id = searchForm.order_id
    if (searchForm.username) params.username = searchForm.username
    if (searchForm.status) params.status = searchForm.status
    if (searchForm.type) params.type = searchForm.type
    const res = await request.get({ url: '/api/admin/upgrades', params })
    tableData.value = res?.list || res?.data || res || []
    pagination.total = res?.total || 0
  } catch { ElMessage.error($t('upgrades.fetchListFailed')) } finally { loading.value = false }
}

const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/upgrades/stats' })
    if (data) Object.assign(stats, data)
  } catch { /* ignore */ }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.order_id = ''; searchForm.username = ''; searchForm.status = ''; searchForm.type = ''; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const goToOrder = (id: string | number) => {
  router.push({ path: `/order-detail/${id}` })
}

const handleDetail = (row: any) => {
  detailData.value = { ...row }
  detailDialogVisible.value = true
}

const handleApprove = async (row: any) => {
  try {
    await ElMessageBox.confirm(`${$t('upgrades.confirmApprove')} #${row.id}?`, $t('upgrades.approveUpgrade'))
    await request.post({ url: `/api/admin/upgrades/${row.id}/approve`, showSuccessMessage: true })
    fetchData(); fetchStats()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error($t('upgrades.operationFailed')) }
}

const handleCancel = async (row: any) => {
  try {
    await ElMessageBox.confirm(`${$t('upgrades.confirmCancel')} #${row.id}?`, $t('upgrades.cancelUpgrade'))
    await request.post({ url: `/api/admin/upgrades/${row.id}/cancel`, showSuccessMessage: true })
    fetchData(); fetchStats()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error($t('upgrades.operationFailed')) }
}

const handleMarkPaid = async (row: any) => {
  try {
    await ElMessageBox.confirm(`${$t('upgrades.confirmMarkPaid')} #${row.id} ${$t('upgrades.asPaid')}?`, $t('upgrades.markPaid'))
    await request.post({ url: `/api/admin/upgrades/${row.id}/mark-paid`, showSuccessMessage: true })
    fetchData(); fetchStats()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error($t('upgrades.operationFailed')) }
}

onMounted(() => { fetchData(); fetchStats() })
</script>

<style scoped lang="scss">
.upgrades-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 16px; .el-form-item { margin-bottom: 0; } }
.stat-section { margin-bottom: 24px; }
.stat-card { text-align: center; padding: 8px 0; }
.stat-value { font-size: 20px; font-weight: 600; color: var(--el-color-primary); margin-bottom: 4px; }
.stat-label { color: var(--el-text-color-secondary); font-size: 14px; }
.amount { color: var(--el-color-danger); font-weight: 600; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
