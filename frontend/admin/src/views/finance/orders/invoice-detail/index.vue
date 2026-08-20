<template>
  <div class="invoice-detail-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('invoiceDetail.title') }}</span>
          <div class="header-actions">
            <el-button @click="handleBack">
              <el-icon><Back /></el-icon>
              {{ $t('common.back') }}
            </el-button>
          </div>
        </div>
      </template>

      <div v-loading="loading" class="loading-container">
        <div v-if="bill" class="bill-content">
          <!-- 账单基本信息 -->
          <el-descriptions :column="2" border class="bill-info">
            <el-descriptions-item :label="$t('invoiceDetail.billNo')">{{ bill.bill_no }}</el-descriptions-item>
            <el-descriptions-item :label="$t('common.status')">
              <el-tag :type="statusTypeMap[bill.status]" size="default">
                {{ statusLabelMap[bill.status] }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('invoiceDetail.relatedClient')">
              <el-button type="primary" link @click="handleViewClient">
                {{ bill.client_username }}
              </el-button>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('invoiceDetail.relatedProduct')">{{ bill.product_name || '-' }}</el-descriptions-item>
            <el-descriptions-item :label="$t('invoiceDetail.billAmount')">
              <span class="amount-text">¥{{ formatAmount(bill.amount) }}</span>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('invoiceDetail.paidAmount')">
              <span class="amount-text">¥{{ formatAmount(bill.paid_amount) }}</span>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('invoiceDetail.unpaidAmount')">
              <span class="amount-text" :class="{ 'text-danger': bill.amount - bill.paid_amount > 0 }">
                ¥{{ formatAmount(bill.amount - bill.paid_amount) }}
              </span>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('invoiceDetail.discountAmount')">¥{{ formatAmount(bill.discount_amount || 0) }}</el-descriptions-item>
            <el-descriptions-item :label="$t('invoiceDetail.createdAt')">{{ bill.created_at }}</el-descriptions-item>
            <el-descriptions-item :label="$t('invoiceDetail.dueDate')">{{ bill.due_date }}</el-descriptions-item>
            <el-descriptions-item :label="$t('invoiceDetail.paidAt')">{{ bill.paid_at || '-' }}</el-descriptions-item>
            <el-descriptions-item :label="$t('common.remark')" :span="2">{{ bill.remark || $t('common.none') }}</el-descriptions-item>
          </el-descriptions>

          <!-- 账单明细列表 -->
          <div class="section">
            <h3>{{ $t('invoiceDetail.billItems') }}</h3>
            <el-table :data="bill.items || []" style="width: 100%" border>
              <el-table-column prop="id" label="ID" width="80" align="center" />
              <el-table-column prop="description" :label="$t('invoiceDetail.description')" min-width="200" />
              <el-table-column prop="quantity" :label="$t('common.quantity')" width="100" align="center" />
              <el-table-column prop="unit_price" :label="$t('orderDetail.unitPrice')" width="120" align="right">
                <template #default="{ row }">
                  ¥{{ formatAmount(row.unit_price) }}
                </template>
              </el-table-column>
              <el-table-column prop="discount" :label="$t('invoiceDetail.discount')" width="100" align="right">
                <template #default="{ row }">
                  {{ row.discount || '0%' }}
                </template>
              </el-table-column>
              <el-table-column prop="amount" :label="$t('invoiceDetail.amount')" width="120" align="right">
                <template #default="{ row }">
                  <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
                </template>
              </el-table-column>
            </el-table>
            <div class="bill-total">
              <span>{{ $t('invoiceDetail.total') }}: </span>
              <span class="amount-text">¥{{ formatAmount(bill.amount) }}</span>
            </div>
          </div>

          <!-- 支付记录 -->
          <div class="section">
            <h3>{{ $t('invoiceDetail.paymentRecords') }}</h3>
            <el-table :data="payments" v-loading="paymentsLoading" style="width: 100%" border>
              <el-table-column prop="id" label="ID" width="80" align="center" />
              <el-table-column prop="payment_no" :label="$t('invoiceDetail.paymentNo')" width="200" />
              <el-table-column prop="gateway" :label="$t('invoiceDetail.paymentGateway')" width="120" />
              <el-table-column prop="amount" :label="$t('invoiceDetail.paymentAmount')" width="120" align="right">
                <template #default="{ row }">
                  <span class="text-success">¥{{ formatAmount(row.amount) }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="status" :label="$t('common.status')" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                    {{ row.status === 1 ? $t('invoiceDetail.success') : $t('invoiceDetail.failed') }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="created_at" :label="$t('invoiceDetail.paymentTime')" width="180" />
              <el-table-column prop="remark" :label="$t('common.remark')" min-width="150" show-overflow-tooltip />
            </el-table>
            <el-empty v-if="!payments.length && !paymentsLoading" :description="$t('invoiceDetail.noPayments')" />
          </div>
        </div>
      </div>
    </el-card>

    <!-- 底部操作栏 -->
    <div class="action-bar" v-if="bill">
      <el-button
        v-if="bill.status === 0 || bill.status === 1"
        type="primary"
        @click="handleSend"
        :loading="actionLoading"
      >
        <el-icon><Promotion /></el-icon>
        {{ $t('invoiceDetail.sendInvoice') }}
      </el-button>
      <el-button
        v-if="bill.status === 2 || bill.status === 3"
        type="success"
        @click="handleRefund"
        :loading="actionLoading"
      >
        <el-icon><RefreshRight /></el-icon>
        {{ $t('invoiceDetail.refund') }}
      </el-button>
      <el-popconfirm
        v-if="bill.status !== 4 && bill.status !== 5"
        :title="$t('invoiceDetail.confirmCancel')"
        @confirm="handleCancel"
      >
        <template #reference>
          <el-button type="danger" :loading="actionLoading">
            <el-icon><CircleClose /></el-icon>
            {{ $t('invoiceDetail.cancelInvoice') }}
          </el-button>
        </template>
      </el-popconfirm>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Back, Promotion, RefreshRight, CircleClose } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const paymentsLoading = ref(false)
const actionLoading = ref(false)
const bill = ref<any>({})
const payments = ref<any[]>([])

const statusTypeMap: Record<number, any> = {
  0: 'info',
  1: 'warning',
  2: 'success',
  3: 'success',
  4: 'danger'
}

const statusLabelMap: Record<number, () => string> = {
  0: () => $t('invoiceDetail.statusPending'),
  1: () => $t('invoiceDetail.statusSent'),
  2: () => $t('invoiceDetail.statusPaid'),
  3: () => $t('invoiceDetail.statusPartialPaid'),
  4: () => $t('invoiceDetail.statusCancelled')
}

const fetchBill = async () => {
  const id = route.params.id
  if (!id) return

  loading.value = true
  try {
    const data = await request.get({ url: `/api/admin/invoices/${id}` })
    bill.value = data.bill || data
  } catch (error) {
    console.error('fetch invoice detail failed:', error)
    ElMessage.error($t('invoiceDetail.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const fetchPayments = async () => {
  const id = route.params.id
  if (!id) return

  paymentsLoading.value = true
  try {
    const data = await request.get({ url: `/api/admin/invoices/${id}/payments` })
    payments.value = data || []
  } catch (error) {
    console.error('fetch payments failed:', error)
  } finally {
    paymentsLoading.value = false
  }
}

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const handleBack = () => {
  router.back()
}

const handleViewClient = () => {
  if (bill.value?.client_id) {
    router.push(`/customer-view/${bill.value.client_id}`)
  }
}

const handleSend = async () => {
  const id = route.params.id
  if (!id) return

  actionLoading.value = true
  try {
    await request.post({ url: `/api/admin/invoices/${id}/email` })
    ElMessage.success($t('invoiceDetail.sendSuccess'))
    fetchBill()
  } catch (error) {
    ElMessage.error($t('invoiceDetail.sendFailed'))
  } finally {
    actionLoading.value = false
  }
}

const handleRefund = async () => {
  const id = route.params.id
  if (!id) return

  try {
    await ElMessageBox.confirm($t('invoiceDetail.confirmRefundMsg'), $t('invoiceDetail.confirmRefundTitle'), {
      type: 'warning'
    })

    actionLoading.value = true
    await request.post({ url: `/api/admin/invoices/${id}/refund` })
    ElMessage.success($t('invoiceDetail.refundSuccess'))
    fetchBill()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error($t('invoiceDetail.refundFailed'))
    }
  } finally {
    actionLoading.value = false
  }
}

const handleCancel = async () => {
  const id = route.params.id
  if (!id) return

  actionLoading.value = true
  try {
    await request.post({ url: `/api/admin/invoices/${id}/cancel` })
    ElMessage.success($t('invoiceDetail.invoiceCancelled'))
    fetchBill()
  } catch (error) {
    ElMessage.error($t('invoiceDetail.cancelFailed'))
  } finally {
    actionLoading.value = false
  }
}

onMounted(() => {
  fetchBill()
  fetchPayments()
})
</script>

<style scoped lang="scss">
.invoice-detail-page {
  padding: 20px;
  padding-bottom: 80px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.loading-container {
  min-height: 400px;
}

.bill-info {
  margin-bottom: 24px;
}

.amount-text {
  font-weight: 600;
  color: var(--el-color-primary);
}

.text-success {
  color: var(--el-color-success);
}

.text-danger {
  color: var(--el-color-danger);
}

.section {
  margin-top: 24px;

  h3 {
    margin: 0 0 16px;
    font-size: 16px;
    font-weight: 600;
  }

  .bill-total {
    text-align: right;
    margin-top: 16px;
    padding: 12px 16px;
    background: var(--el-fill-color-lighter);
    border-radius: 4px;
    font-size: 16px;
  }
}

.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 16px 20px;
  background: var(--el-bg-color);
  border-top: 1px solid var(--el-border-color-lighter);
  display: flex;
  gap: 12px;
  z-index: 100;
}
</style>
