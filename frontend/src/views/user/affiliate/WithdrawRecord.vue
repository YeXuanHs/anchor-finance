<template>
  <div class="withdraw-record-page">
    <div class="page-header">
      <h1 class="page-title">推广提现记录</h1>
      <div class="header-actions">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          size="default"
          style="width: 280px;"
        />
        <el-button type="primary" @click="handleWithdraw">
          <el-icon><Plus /></el-icon>申请提现
        </el-button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="stats-row">
      <el-card shadow="never" class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: linear-gradient(135deg, #0056FF, #4080FF);">
            <el-icon :size="24"><Wallet /></el-icon>
          </div>
          <div class="stat-info">
            <span class="stat-value">¥{{ availableBalance }}</span>
            <span class="stat-label">可提现余额</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: linear-gradient(135deg, #52c41a, #73d13d);">
            <el-icon :size="24"><CircleCheck /></el-icon>
          </div>
          <div class="stat-info">
            <span class="stat-value">¥{{ totalWithdrawn }}</span>
            <span class="stat-label">累计已提现</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: linear-gradient(135deg, #fa8c16, #ffc53d);">
            <el-icon :size="24"><Clock /></el-icon>
          </div>
          <div class="stat-info">
            <span class="stat-value">¥{{ pendingAmount }}</span>
            <span class="stat-label">待审核金额</span>
          </div>
        </div>
      </el-card>
    </div>

    <!-- Status Tabs -->
    <el-tabs v-model="activeTab" class="filter-tabs">
      <el-tab-pane v-for="tab in statusTabs" :key="tab.value" :name="tab.value">
        <template #label>
          <span class="tab-label">
            {{ tab.label }}
            <el-badge
              v-if="tab.count > 0"
              :value="tab.count"
              :max="99"
              class="tab-badge"
            />
          </span>
        </template>
      </el-tab-pane>
    </el-tabs>

    <!-- Data Table -->
    <el-card shadow="never" class="table-card">
      <el-table :data="paginatedRecords" stripe style="width: 100%" v-loading="loading" empty-text="暂无提现记录">
        <el-table-column prop="id" label="提现单号" min-width="170" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="mono-text">{{ row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="提现金额" width="130">
          <template #default="{ row }">
            <span class="amount-text">¥{{ row.amount }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="method" label="提现方式" width="120">
          <template #default="{ row }">
            <div class="method-cell">
              <el-icon :size="16" class="method-icon"><component :is="getMethodIcon(row.method)" /></el-icon>
              <span>{{ row.methodText }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="account" label="收款账户" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="account-text">{{ row.account }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small" effect="light" round>
              {{ row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="申请时间" width="170" sortable />
        <el-table-column prop="completedAt" label="完成时间" width="170" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'pending'"
              type="danger"
              size="small"
              link
              @click="handleCancel(row)"
            >撤回</el-button>
            <el-button
              type="primary"
              size="small"
              link
              @click="handleDetail(row)"
            >详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="filteredRecords.length"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          background
        />
      </div>
    </el-card>

    <!-- Withdraw Dialog -->
    <el-dialog v-model="withdrawDialogVisible" title="申请提现" width="480px" :close-on-click-modal="false">
      <el-form :model="withdrawForm" label-width="100px" class="withdraw-form">
        <el-form-item label="可提现余额">
          <span class="balance-text">¥{{ availableBalance }}</span>
        </el-form-item>
        <el-form-item label="提现金额" required>
          <el-input-number
            v-model="withdrawForm.amount"
            :min="1"
            :max="parseFloat(availableBalance.replace(/,/g, ''))"
            :precision="2"
            :step="100"
            style="width: 100%;"
          />
        </el-form-item>
        <el-form-item label="提现方式" required>
          <el-radio-group v-model="withdrawForm.method">
            <el-radio value="alipay">支付宝</el-radio>
            <el-radio value="wechat">微信</el-radio>
            <el-radio value="bank">银行卡</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="收款账户" required>
          <el-input v-model="withdrawForm.account" placeholder="请输入收款账户" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="withdrawForm.remark" type="textarea" :rows="2" placeholder="选填" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="withdrawDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitWithdraw">提交申请</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Wallet, CircleCheck, Clock, CreditCard, ChatDotRound, Memo } from '@element-plus/icons-vue'
import request from '@/utils/request'

const searchKey = ref('')
const activeTab = ref('all')
const currentPage = ref(1)
const pageSize = ref(10)
const dateRange = ref<[Date, Date] | null>(null)
const withdrawDialogVisible = ref(false)

const withdrawForm = ref({
  amount: 100,
  method: 'alipay',
  account: '',
  remark: ''
})

const availableBalance = ref('2,580.00')

interface WithdrawRecord {
  id: string
  amount: string
  method: string
  methodText: string
  account: string
  status: string
  statusText: string
  createdAt: string
  completedAt: string
}

const records = ref<WithdrawRecord[]>([])

const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v2/affiliate/withdraws')
    if (data?.data) {
      records.value = data.data.list || data.data || []
      if (data.data.balance !== undefined) {
        availableBalance.value = data.data.balance
      }
    }
  } catch (e) {
    console.error('Failed to fetch withdraw records:', e)
  } finally {
    loading.value = false
  }
})

const totalWithdrawn = computed(() => {
  const sum = records.value
    .filter(r => r.status === 'completed')
    .reduce((acc, r) => acc + parseFloat(r.amount.replace(/,/g, '')), 0)
  return sum.toLocaleString('zh-CN', { minimumFractionDigits: 2 })
})

const pendingAmount = computed(() => {
  const sum = records.value
    .filter(r => r.status === 'pending')
    .reduce((acc, r) => acc + parseFloat(r.amount.replace(/,/g, '')), 0)
  return sum.toLocaleString('zh-CN', { minimumFractionDigits: 2 })
})

const statusTabs = computed(() => [
  { label: '全部', value: 'all', count: records.value.length },
  { label: '审核中', value: 'pending', count: records.value.filter(r => r.status === 'pending').length },
  { label: '已到账', value: 'completed', count: records.value.filter(r => r.status === 'completed').length },
  { label: '已拒绝', value: 'rejected', count: records.value.filter(r => r.status === 'rejected').length },
  { label: '已撤回', value: 'cancelled', count: records.value.filter(r => r.status === 'cancelled').length }
])

const filteredRecords = computed(() => {
  let result = records.value
  if (activeTab.value !== 'all') {
    result = result.filter(r => r.status === activeTab.value)
  }
  if (searchKey.value) {
    const key = searchKey.value.toLowerCase()
    result = result.filter(r =>
      r.id.toLowerCase().includes(key) || r.account.toLowerCase().includes(key)
    )
  }
  return result
})

const paginatedRecords = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredRecords.value.slice(start, start + pageSize.value)
})

function getStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    completed: 'success',
    pending: 'warning',
    rejected: 'danger',
    cancelled: 'info'
  }
  return map[status] || 'info'
}

function getMethodIcon(method: string) {
  const map: Record<string, any> = {
    alipay: CreditCard,
    wechat: ChatDotRound,
    bank: Memo
  }
  return map[method] || CreditCard
}

function handleWithdraw() {
  withdrawDialogVisible.value = true
}

function handleSubmitWithdraw() {
  if (!withdrawForm.value.amount || withdrawForm.value.amount <= 0) {
    ElMessage.warning('请输入有效的提现金额')
    return
  }
  if (!withdrawForm.value.account) {
    ElMessage.warning('请输入收款账户')
    return
  }
  ElMessage.success('提现申请已提交')
  withdrawDialogVisible.value = false
  withdrawForm.value = { amount: 100, method: 'alipay', account: '', remark: '' }
}

function handleCancel(row: WithdrawRecord) {
  ElMessageBox.confirm('确定要撤回该提现申请吗？', '撤回确认', {
    confirmButtonText: '确定撤回',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    row.status = 'cancelled'
    row.statusText = '已撤回'
    ElMessage.success('已撤回提现申请')
  }).catch(() => {})
}

function handleDetail(row: WithdrawRecord) {
  ElMessage.info(`查看提现详情：${row.id}`)
}
</script>

<style scoped>
.withdraw-record-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.stat-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.stat-card :deep(.el-card__body) {
  padding: 20px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
}

.stat-label {
  font-size: 13px;
  color: #909399;
}

.filter-tabs {
  background: #fff;
  border-radius: 12px;
  padding: 4px 16px 0;
  border: 1px solid #e8ecf1;
}

.filter-tabs :deep(.el-tabs__header) {
  margin: 0;
}

.tab-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.tab-badge :deep(.el-badge__content) {
  font-size: 10px;
}

.table-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.table-card :deep(.el-card__body) {
  padding: 0;
}

.table-card :deep(.el-table th.el-table__cell) {
  background: #fafbfc;
  color: #606266;
  font-weight: 600;
}

.mono-text {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 13px;
  color: #606266;
}

.method-cell {
  display: flex;
  align-items: center;
  gap: 6px;
}

.method-icon {
  color: #909399;
}

.account-text {
  font-size: 13px;
  color: #606266;
}

.amount-text {
  font-weight: 600;
  color: #303133;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px 20px;
  border-top: 1px solid #e8ecf1;
}

.withdraw-form {
  padding: 0 20px;
}

.balance-text {
  font-size: 20px;
  font-weight: 700;
  color: #0056FF;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .header-actions {
    width: 100%;
    flex-direction: column;
  }

  .stats-row {
    grid-template-columns: 1fr;
  }
}
</style>
