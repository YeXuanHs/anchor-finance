<template>
  <div class="credit-bill-detail-page">
    <div class="page-header">
      <el-button text @click="router.back()">
        <el-icon><ArrowLeft /></el-icon> 返回
      </el-button>
      <h1 class="page-title">账单详情</h1>
    </div>

    <div v-loading="loading" class="detail-content">
      <!-- 账单概况 -->
      <el-card shadow="never" class="summary-card">
        <div class="summary-top">
          <div class="bill-id">
            <span class="label">账单编号：</span>
            <span class="mono">{{ billInfo.id }}</span>
          </div>
          <el-tag :type="statusTagType(billInfo.status)" size="large" effect="light">
            {{ statusText(billInfo.status) }}
          </el-tag>
        </div>

        <el-row :gutter="20" class="summary-grid">
          <el-col :span="6">
            <div class="summary-item">
              <span class="item-label">账单金额</span>
              <span class="item-value">¥{{ billInfo.amount?.toFixed(2) || '0.00' }}</span>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="summary-item">
              <span class="item-label">已还金额</span>
              <span class="item-value paid">¥{{ billInfo.paid_amount?.toFixed(2) || '0.00' }}</span>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="summary-item">
              <span class="item-label">待还金额</span>
              <span class="item-value pending">¥{{ billInfo.remaining?.toFixed(2) || '0.00' }}</span>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="summary-item">
              <span class="item-label">逾期费用</span>
              <span class="item-value overdue">¥{{ billInfo.overdue_fee?.toFixed(2) || '0.00' }}</span>
            </div>
          </el-col>
        </el-row>
      </el-card>

      <!-- 账单信息 -->
      <el-card shadow="never" class="info-card">
        <template #header>
          <span class="card-title">账单信息</span>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="账单类型">{{ billTypeText(billInfo.type) }}</el-descriptions-item>
          <el-descriptions-item label="关联产品">{{ billInfo.product_name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="关联订单">
            <el-button v-if="billInfo.order_no" type="primary" link @click="goOrder(billInfo.order_no)">
              {{ billInfo.order_no }}
            </el-button>
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item label="账单周期">{{ billInfo.period || '-' }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ billInfo.created_at || '-' }}</el-descriptions-item>
          <el-descriptions-item label="到期时间">{{ billInfo.due_date || '-' }}</el-descriptions-item>
          <el-descriptions-item label="支付时间">{{ billInfo.paid_at || '-' }}</el-descriptions-item>
          <el-descriptions-item label="支付方式">{{ billInfo.payment_method || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 逾期信息 -->
      <el-card v-if="billInfo.status === 'overdue'" shadow="never" class="overdue-card">
        <template #header>
          <div class="card-header-row">
            <span class="card-title">逾期信息</span>
            <el-tag type="danger" effect="dark">已逾期</el-tag>
          </div>
        </template>
        <el-alert
          :title="`该账单已逾期 ${billInfo.overdue_days || 0} 天`"
          type="error"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />
        <el-descriptions :column="2" border>
          <el-descriptions-item label="逾期天数">{{ billInfo.overdue_days || 0 }} 天</el-descriptions-item>
          <el-descriptions-item label="逾期费率">{{ billInfo.overdue_rate || '0.05' }}%/天</el-descriptions-item>
          <el-descriptions-item label="逾期费用">¥{{ billInfo.overdue_fee?.toFixed(2) || '0.00' }}</el-descriptions-item>
          <el-descriptions-item label="应还总额">¥{{ billInfo.total_due?.toFixed(2) || '0.00' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 还款记录 -->
      <el-card shadow="never" class="records-card">
        <template #header>
          <div class="card-header-row">
            <span class="card-title">还款记录</span>
            <el-button
              v-if="billInfo.status !== 'paid'"
              type="primary"
              size="small"
              @click="repayDialogVisible = true"
            >
              立即还款
            </el-button>
          </div>
        </template>
        <el-table :data="repayRecords" stripe empty-text="暂无还款记录">
          <el-table-column prop="id" label="还款编号" width="160">
            <template #default="{ row }">
              <span class="mono">{{ row.id }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="amount" label="还款金额" width="130">
            <template #default="{ row }">
              <span class="text-success">¥{{ row.amount?.toFixed(2) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="method" label="还款方式" width="120" />
          <el-table-column prop="created_at" label="还款时间" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'success' ? 'success' : 'warning'" size="small">
                {{ row.status === 'success' ? '成功' : '处理中' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    <!-- 还款弹窗 -->
    <el-dialog v-model="repayDialogVisible" title="信用额度还款" width="450px">
      <el-form label-width="80px">
        <el-form-item label="待还金额">
          <span class="current-amount text-danger">¥{{ billInfo.remaining?.toFixed(2) || '0.00' }}</span>
        </el-form-item>
        <el-form-item label="还款金额" required>
          <el-input-number
            v-model="repayAmount"
            :min="1"
            :max="(billInfo.remaining || 0) + (billInfo.overdue_fee || 0)"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="repayDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleRepay">确认还款</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import request from '@/utils/http'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const submitting = ref(false)
const repayDialogVisible = ref(false)
const repayAmount = ref(0)

const billInfo = ref<any>({
  id: '',
  amount: 0,
  paid_amount: 0,
  remaining: 0,
  overdue_fee: 0,
  overdue_days: 0,
  overdue_rate: '0.05',
  total_due: 0,
  status: 'unpaid',
  type: '',
  product_name: '',
  order_no: '',
  period: '',
  created_at: '',
  due_date: '',
  paid_at: '',
  payment_method: ''
})

const repayRecords = ref<any[]>([])

const statusTagType = (status: string) => {
  const map: Record<string, string> = { paid: 'success', unpaid: 'warning', overdue: 'danger', partial: 'primary' }
  return map[status] || 'info'
}

const statusText = (status: string) => {
  const map: Record<string, string> = { paid: '已支付', unpaid: '待支付', overdue: '已逾期', partial: '部分还款' }
  return map[status] || status
}

const billTypeText = (type: string) => {
  const map: Record<string, string> = { credit_use: '额度使用', overdue: '逾期还款', fee: '手续费' }
  return map[type] || type || '-'
}

function goOrder(orderNo: string) {
  router.push(`/user/orders/${orderNo}`)
}

async function fetchBillDetail() {
  const id = route.params.id as string
  if (!id) return
  loading.value = true
  try {
    const res = await request.get(`/v1/credit/bills/${id}`)
    billInfo.value = res.data?.data || res.data || billInfo.value
    repayAmount.value = billInfo.value.remaining || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function fetchRepayRecords() {
  const id = route.params.id as string
  if (!id) return
  try {
    const res = await request.get(`/v1/credit/bills/${id}/repayments`)
    repayRecords.value = res.data?.data?.list || res.data?.list || []
  } catch {
    repayRecords.value = []
  }
}

async function handleRepay() {
  if (!repayAmount.value || repayAmount.value <= 0) {
    ElMessage.warning('请输入有效的还款金额')
    return
  }
  submitting.value = true
  try {
    await request.post(`/v1/credit/bills/${route.params.id}/repay`, { amount: repayAmount.value })
    ElMessage.success('还款成功')
    repayDialogVisible.value = false
    fetchBillDetail()
    fetchRepayRecords()
  } catch (e: any) {
    ElMessage.error(e.message || '还款失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchBillDetail()
  fetchRepayRecords()
})
</script>

<style scoped lang="scss">
.credit-bill-detail-page {
  .page-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 20px;

    .page-title {
      font-size: 20px;
      font-weight: 700;
      color: #303133;
      margin: 0;
    }
  }

  .detail-content {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .summary-card {
    border-radius: 12px;
    border: 1px solid #e8ecf1;

    .summary-top {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 24px;

      .bill-id {
        font-size: 14px;
        color: #606266;

        .label { color: #909399; }
        .mono { font-family: 'Monaco', 'Menlo', monospace; color: #303133; font-weight: 600; }
      }
    }

    .summary-grid {
      .summary-item {
        text-align: center;
        padding: 16px;
        background: #f5f7fa;
        border-radius: 8px;

        .item-label {
          display: block;
          font-size: 13px;
          color: #909399;
          margin-bottom: 8px;
        }

        .item-value {
          font-size: 22px;
          font-weight: 700;
          color: #303133;

          &.paid { color: #67c23a; }
          &.pending { color: #e6a23c; }
          &.overdue { color: #f56c6c; }
        }
      }
    }
  }

  .info-card, .overdue-card, .records-card {
    border-radius: 12px;
    border: 1px solid #e8ecf1;
  }

  .overdue-card {
    border-color: #fde2e2;
  }

  .card-title {
    font-size: 15px;
    font-weight: 600;
    color: #303133;
  }

  .card-header-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .mono {
    font-family: 'Monaco', 'Menlo', monospace;
    font-size: 13px;
    color: #606266;
  }

  .text-success {
    color: #67c23a;
    font-weight: 600;
  }

  .text-danger {
    color: #f56c6c;
    font-weight: 600;
  }

  .current-amount {
    font-size: 18px;
    font-weight: bold;
  }
}
</style>
