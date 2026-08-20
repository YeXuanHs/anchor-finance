<template>
  <div class="refund-detail-page">
    <!-- 退款申请信息 -->
    <el-card shadow="never" class="detail-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('refundDetail.title') }}</span>
          <el-space>
            <el-tag :type="getStatusTag(refundData.status)" size="large">
              {{ getStatusText(refundData.status) }}
            </el-tag>
            <template v-if="refundData.status === 'pending'">
              <el-button type="success" @click="handleApprove">{{ $t('refundDetail.approveRefund') }}</el-button>
              <el-button type="danger" @click="handleReject">{{ $t('refundDetail.rejectRefund') }}</el-button>
            </template>
          </el-space>
        </div>
      </template>

      <el-descriptions :column="3" border>
        <el-descriptions-item :label="$t('refundDetail.refundNo')">{{ refundData.refund_no }}</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.applyTime')">{{ refundData.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.refundType')">
          <el-tag :type="getRefundTypeTag(refundData.refund_type)" size="small">
            {{ getRefundTypeText(refundData.refund_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.refundAmount')">
          <span class="text-danger amount-text">¥{{ refundData.amount?.toFixed(2) || '0.00' }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.refundReason')" :span="2">{{ refundData.reason || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.refundDescription')" :span="3">{{ refundData.description || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.applicant')">{{ refundData.username }} (ID: {{ refundData.user_id }})</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.contact')">{{ refundData.contact || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.processor')">{{ refundData.processor || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.processedAt')">{{ refundData.processed_at || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.processRemark')" :span="2">{{ refundData.process_remark || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 原订单信息 -->
    <el-card shadow="never" class="detail-card">
      <template #header>
        <span>{{ $t('refundDetail.originalOrder') }}</span>
      </template>

      <el-descriptions :column="3" border>
        <el-descriptions-item :label="$t('refundDetail.orderNo')">
          <el-button type="primary" link @click="handleViewOrder">{{ orderData.order_no }}</el-button>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.createdAt')">{{ orderData.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.orderStatus')">
          <el-tag :type="getOrderStatusTag(orderData.status)" size="small">
            {{ getOrderStatusText(orderData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.productName')">{{ orderData.product_name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.productType')">{{ orderData.product_type || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.billingCycle')">{{ orderData.billing_cycle || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.totalAmount')">
          <span class="text-primary">¥{{ orderData.total_amount?.toFixed(2) || '0.00' }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.paidAmount')">
          <span class="text-primary">¥{{ orderData.paid_amount?.toFixed(2) || '0.00' }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.discountAmount')">¥{{ orderData.discount_amount?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.paymentMethod')">{{ orderData.payment_method || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.paidAt')">{{ orderData.paid_at || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('refundDetail.transactionNo')">{{ orderData.transaction_no || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 退款处理流程 -->
    <el-card shadow="never" class="detail-card">
      <template #header>
        <span>{{ $t('refundDetail.refundProcess') }}</span>
      </template>

      <el-timeline>
        <el-timeline-item
          v-for="(step, index) in refundSteps"
          :key="index"
          :timestamp="step.time"
          :type="step.type"
          :hollow="step.hollow"
        >
          <div class="step-content">
            <div class="step-title">{{ step.title }}</div>
            <div class="step-desc" v-if="step.description">{{ step.description }}</div>
            <div class="step-operator" v-if="step.operator">{{ $t('refundDetail.operator') }}: {{ step.operator }}</div>
          </div>
        </el-timeline-item>
      </el-timeline>
    </el-card>

    <!-- 审批对话框 -->
    <el-dialog v-model="approveDialogVisible" :title="approveForm.action === 'approve' ? $t('refundDetail.approveRefund') : $t('refundDetail.rejectRefund')" width="500px">
      <el-alert
        v-if="approveForm.action === 'approve'"
        :title="$t('refundDetail.refundAlert')"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 20px;"
      />
      <el-form :model="approveForm" :rules="approveRules" ref="approveFormRef" label-width="100px">
        <el-form-item :label="$t('refundDetail.refundNo')">
          <el-input :value="refundData.refund_no" disabled />
        </el-form-item>
        <el-form-item :label="$t('refundDetail.refundAmount')">
          <el-input :value="'¥' + (refundData.amount?.toFixed(2) || '0.00')" disabled />
        </el-form-item>
        <el-form-item :label="$t('refundDetail.approveRemarkLabel')" prop="remark">
          <el-input v-model="approveForm.remark" type="textarea" :rows="3" :placeholder="approveForm.action === 'approve' ? $t('refundDetail.approveRemarkPlaceholder') : $t('refundDetail.rejectReasonPlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="approveDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button
          :type="approveForm.action === 'approve' ? 'success' : 'danger'"
          @click="handleSubmitApprove"
          :loading="submitLoading"
        >
          {{ approveForm.action === 'approve' ? $t('refundDetail.confirmApprove') : $t('refundDetail.confirmReject') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'RefundDetail' })

const route = useRoute()
const router = useRouter()

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)

// 退款数据
const refundData = ref<any>({})

// 订单数据
const orderData = ref<any>({})

// 退款步骤
const refundSteps = ref<any[]>([])

// 审批对话框
const approveDialogVisible = ref(false)
const approveFormRef = ref<FormInstance>()
const approveForm = reactive({
  action: 'approve' as 'approve' | 'reject',
  remark: ''
})

// 审批表单验证规则
const approveRules: FormRules = {
  remark: [
    {
      validator: (rule: any, value: string, callback: any) => {
        if (approveForm.action === 'reject' && !value) {
          callback(new Error($t('refundDetail.enterRejectReason')))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 获取退款状态标签
const getStatusTag = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    pending: 'warning',
    approved: 'primary',
    processing: 'primary',
    completed: 'success',
    rejected: 'danger',
    cancelled: 'info'
  }
  return map[status] || 'info'
}

// 获取退款状态文本
const getStatusText = (status: string) => {
  const map: Record<string, () => string> = {
    pending: () => $t('refundDetail.statusPending'),
    approved: () => $t('refundDetail.statusApproved'),
    processing: () => $t('refundDetail.statusProcessing'),
    completed: () => $t('refundDetail.statusCompleted'),
    rejected: () => $t('refundDetail.statusRejected'),
    cancelled: () => $t('refundDetail.statusCancelled')
  }
  return map[status]?.() || $t('common.unknown')
}

// 获取退款类型标签
const getRefundTypeTag = (type: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    full: 'danger',
    partial: 'warning',
    cancellation: 'info'
  }
  return map[type] || 'info'
}

// 获取退款类型文本
const getRefundTypeText = (type: string) => {
  const map: Record<string, () => string> = {
    full: () => $t('refundDetail.refundTypeFull'),
    partial: () => $t('refundDetail.refundTypePartial'),
    cancellation: () => $t('refundDetail.refundTypeCancellation')
  }
  return map[type]?.() || $t('common.unknown')
}

// 获取订单状态标签
const getOrderStatusTag = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    pending: 'warning',
    paid: 'success',
    cancelled: 'info',
    refunded: 'danger'
  }
  return map[status] || 'info'
}

// 获取订单状态文本
const getOrderStatusText = (status: string) => {
  const map: Record<string, () => string> = {
    pending: () => $t('refundDetail.orderStatusPending'),
    paid: () => $t('refundDetail.orderStatusPaid'),
    cancelled: () => $t('refundDetail.orderStatusCancelled'),
    refunded: () => $t('refundDetail.orderStatusRefunded')
  }
  return map[status]?.() || $t('common.unknown')
}

// 获取退款详情
const fetchRefundDetail = async () => {
  const id = route.params.id
  if (!id) return

  loading.value = true
  try {
    const data = await request.get({
      url: `/api/admin/orders/${id}/refund`
    })
    refundData.value = data.refund || {}
    orderData.value = data.order || {}
    refundSteps.value = data.steps || []
  } catch (error) {
    console.error('fetch refund detail failed:', error)
    ElMessage.error($t('refundDetail.fetchFailed'))
  } finally {
    loading.value = false
  }
}

// 查看原订单
const handleViewOrder = () => {
  if (orderData.value.id) {
    router.push(`/order-detail/${orderData.value.id}`)
  }
}

// 通过退款
const handleApprove = () => {
  approveForm.action = 'approve'
  approveForm.remark = ''
  approveDialogVisible.value = true
}

// 拒绝退款
const handleReject = () => {
  approveForm.action = 'reject'
  approveForm.remark = ''
  approveDialogVisible.value = true
}

// 提交审批
const handleSubmitApprove = async () => {
  if (!approveFormRef.value) return

  await approveFormRef.value.validate(async (valid) => {
    if (!valid) return

    if (approveForm.action === 'reject') {
      await ElMessageBox.confirm($t('refundDetail.confirmRejectMsg'), $t('refundDetail.confirmAction'), {
        confirmButtonText: $t('common.confirm'),
        cancelButtonText: $t('common.cancel'),
        type: 'warning'
      })
    }

    submitLoading.value = true
    try {
      await request.post({
        url: `/api/admin/orders/${route.params.id}/refund/process`,
        params: {
          action: approveForm.action,
          remark: approveForm.remark
        },
        showSuccessMessage: true
      })
      approveDialogVisible.value = false
      fetchRefundDetail()
    } catch (error) {
      ElMessage.error($t('refundDetail.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

onMounted(() => {
  fetchRefundDetail()
})
</script>

<style scoped lang="scss">
.refund-detail-page {
  padding: 20px;
}

.detail-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.text-primary {
  color: #409eff;
  font-weight: 600;
}

.text-danger {
  color: #f56c6c;
}

.amount-text {
  font-size: 18px;
}

.step-content {
  .step-title {
    font-weight: 500;
    color: #303133;
    margin-bottom: 4px;
  }

  .step-desc {
    color: #606266;
    font-size: 13px;
  }

  .step-operator {
    color: #909399;
    font-size: 12px;
    margin-top: 4px;
  }
}
</style>
