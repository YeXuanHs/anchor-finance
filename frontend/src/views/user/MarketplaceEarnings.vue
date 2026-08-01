<template>
  <div class="marketplace-earnings-page">
    <div class="page-header">
      <h1 class="page-title">收益管理</h1>
      <el-button type="primary" :disabled="withdrawable <= 0" @click="withdrawDialogVisible = true">
        <el-icon><Download /></el-icon> 申请提现
      </el-button>
    </div>

    <!-- 收益概览 -->
    <div class="summary-grid">
      <el-card shadow="never" class="summary-card total-card">
        <div class="summary-inner">
          <el-icon :size="32" color="#fff"><Coin /></el-icon>
          <div class="summary-info">
            <span class="summary-value white">¥{{ earnings.total_income?.toFixed(2) || '0.00' }}</span>
            <span class="summary-label white">总收入</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#fa8c16"><TrendCharts /></el-icon>
          <div class="summary-info">
            <span class="summary-value">¥{{ earnings.month_income?.toFixed(2) || '0.00' }}</span>
            <span class="summary-label">本月收入</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#52c41a"><Wallet /></el-icon>
          <div class="summary-info">
            <span class="summary-value">¥{{ earnings.withdrawable?.toFixed(2) || '0.00' }}</span>
            <span class="summary-label">可提现金额</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#1890ff"><Document /></el-icon>
          <div class="summary-info">
            <span class="summary-value">{{ earnings.order_count || 0 }}</span>
            <span class="summary-label">成交订单数</span>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 提现记录 -->
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header-row">
          <span class="card-title">提现记录</span>
          <el-select v-model="withdrawStatusFilter" placeholder="状态" clearable style="width: 120px" size="small">
            <el-option label="全部" value="" />
            <el-option label="处理中" value="pending" />
            <el-option label="已到账" value="completed" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </div>
      </template>
      <el-table :data="withdrawRecords" stripe v-loading="loading" empty-text="暂无提现记录">
        <el-table-column prop="id" label="提现编号" width="160">
          <template #default="{ row }">
            <span class="mono">{{ row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="提现金额" width="130">
          <template #default="{ row }">
            <span class="amount-text">¥{{ row.amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="fee" label="手续费" width="110">
          <template #default="{ row }">
            ¥{{ row.fee?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="actual_amount" label="实际到账" width="130">
          <template #default="{ row }">
            <span class="text-success">¥{{ row.actual_amount?.toFixed(2) || row.amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="method" label="提现方式" width="120" />
        <el-table-column prop="account" label="提现账号" min-width="160" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="withdrawStatusType(row.status)" size="small">
              {{ withdrawStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="170" />
        <el-table-column prop="completed_at" label="到账时间" width="170">
          <template #default="{ row }">
            {{ row.completed_at || '-' }}
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchWithdrawRecords"
          @current-change="fetchWithdrawRecords"
        />
      </div>
    </el-card>

    <!-- 提现弹窗 -->
    <el-dialog v-model="withdrawDialogVisible" title="申请提现" width="480px">
      <el-form :model="withdrawForm" label-width="90px">
        <el-form-item label="可提现金额">
          <span class="current-amount text-success">¥{{ withdrawable?.toFixed(2) || '0.00' }}</span>
        </el-form-item>
        <el-form-item label="提现金额" required>
          <el-input-number
            v-model="withdrawForm.amount"
            :min="1"
            :max="withdrawable"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="提现方式" required>
          <el-select v-model="withdrawForm.method" style="width: 100%">
            <el-option label="银行卡" value="bank" />
            <el-option label="支付宝" value="alipay" />
            <el-option label="微信" value="wechat" />
          </el-select>
        </el-form-item>
        <el-form-item label="提现账号" required>
          <el-input v-model="withdrawForm.account" placeholder="请输入收款账号" />
        </el-form-item>
        <el-form-item label="真实姓名" required>
          <el-input v-model="withdrawForm.real_name" placeholder="请输入真实姓名" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="withdrawDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleWithdraw">确认提现</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Coin, TrendCharts, Wallet, Document, Download } from '@element-plus/icons-vue'
import request from '@/utils/http'

const loading = ref(false)
const submitting = ref(false)
const withdrawDialogVisible = ref(false)
const withdrawStatusFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const earnings = ref({
  total_income: 0,
  month_income: 0,
  withdrawable: 0,
  order_count: 0
})

const withdrawable = computed(() => earnings.value.withdrawable || 0)

const withdrawForm = ref({
  amount: 0,
  method: 'bank',
  account: '',
  real_name: ''
})

const withdrawRecords = ref<any[]>([])

const withdrawStatusType = (status: string) => {
  const map: Record<string, string> = { pending: 'warning', completed: 'success', rejected: 'danger' }
  return map[status] || 'info'
}

const withdrawStatusText = (status: string) => {
  const map: Record<string, string> = { pending: '处理中', completed: '已到账', rejected: '已拒绝' }
  return map[status] || status
}

async function fetchEarnings() {
  try {
    const res = await request.get('/v1/marketplace/earnings')
    earnings.value = res.data?.data || earnings.value
  } catch {}
}

async function fetchWithdrawRecords() {
  loading.value = true
  try {
    const params: any = { page: currentPage.value, page_size: pageSize.value }
    if (withdrawStatusFilter.value) params.status = withdrawStatusFilter.value
    const res = await request.get('/v1/marketplace/withdrawals', { params })
    withdrawRecords.value = res.data?.data?.list || res.data?.list || []
    total.value = res.data?.data?.total || 0
  } catch {
    withdrawRecords.value = []
  } finally {
    loading.value = false
  }
}

async function handleWithdraw() {
  if (!withdrawForm.value.amount || withdrawForm.value.amount <= 0) {
    ElMessage.warning('请输入有效的提现金额')
    return
  }
  if (!withdrawForm.value.account) {
    ElMessage.warning('请输入提现账号')
    return
  }
  submitting.value = true
  try {
    await request.post('/v1/marketplace/withdrawals', withdrawForm.value)
    ElMessage.success('提现申请已提交')
    withdrawDialogVisible.value = false
    withdrawForm.value = { amount: 0, method: 'bank', account: '', real_name: '' }
    fetchEarnings()
    fetchWithdrawRecords()
  } catch (e: any) {
    ElMessage.error(e.message || '提现申请失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchEarnings()
  fetchWithdrawRecords()
})
</script>

<style scoped lang="scss">
.marketplace-earnings-page {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    .page-title {
      font-size: 20px;
      font-weight: 700;
      color: #303133;
      margin: 0;
    }
  }

  .summary-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 16px;
    margin-bottom: 20px;
  }

  .summary-card {
    border-radius: 12px;
    border: 1px solid #e8ecf1;

    &.total-card {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      border: none;
    }
  }

  .summary-card :deep(.el-card__body) { padding: 20px; }
  .summary-inner { display: flex; align-items: center; gap: 16px; }
  .summary-info { display: flex; flex-direction: column; gap: 4px; }
  .summary-value { font-size: 24px; font-weight: 700; color: #303133; &.white { color: #fff; } }
  .summary-label { font-size: 13px; color: #909399; &.white { color: rgba(255,255,255,0.8); } }

  .table-card {
    border-radius: 12px;
    border: 1px solid #e8ecf1;
  }

  .table-card :deep(.el-card__body) { padding: 0; }
  .table-card :deep(.el-table th.el-table__cell) { background: #fafbfc; font-weight: 600; }

  .card-header-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .card-title {
    font-size: 15px;
    font-weight: 600;
    color: #303133;
  }

  .mono {
    font-family: 'Monaco', 'Menlo', monospace;
    font-size: 13px;
    color: #606266;
  }

  .amount-text { font-weight: 600; color: #303133; }
  .text-success { color: #67c23a; font-weight: 600; }
  .current-amount { font-size: 18px; font-weight: bold; }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    padding: 16px 20px;
    border-top: 1px solid #e8ecf1;
  }

  @media (max-width: 768px) {
    .summary-grid { grid-template-columns: repeat(2, 1fr); }
  }
}
</style>
