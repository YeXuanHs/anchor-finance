<template>
  <div class="refund-detail-page">
    <!-- 退款申请信息 -->
    <el-card shadow="never" class="detail-card">
      <template #header>
        <div class="card-header">
          <span>退款申请信息</span>
          <el-space>
            <el-tag :type="getStatusTag(refundData.status)" size="large">
              {{ getStatusText(refundData.status) }}
            </el-tag>
            <template v-if="refundData.status === 'pending'">
              <el-button type="success" @click="handleApprove">通过退款</el-button>
              <el-button type="danger" @click="handleReject">拒绝退款</el-button>
            </template>
          </el-space>
        </div>
      </template>

      <el-descriptions :column="3" border>
        <el-descriptions-item label="退款编号">{{ refundData.refund_no }}</el-descriptions-item>
        <el-descriptions-item label="申请时间">{{ refundData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="退款类型">
          <el-tag :type="getRefundTypeTag(refundData.refund_type)" size="small">
            {{ getRefundTypeText(refundData.refund_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="退款金额">
          <span class="text-danger amount-text">¥{{ refundData.amount?.toFixed(2) || '0.00' }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="退款原因" :span="2">{{ refundData.reason || '-' }}</el-descriptions-item>
        <el-descriptions-item label="退款说明" :span="3">{{ refundData.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="申请用户">{{ refundData.username }} (ID: {{ refundData.user_id }})</el-descriptions-item>
        <el-descriptions-item label="联系方式">{{ refundData.contact || '-' }}</el-descriptions-item>
        <el-descriptions-item label="处理人">{{ refundData.processor || '-' }}</el-descriptions-item>
        <el-descriptions-item label="处理时间">{{ refundData.processed_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="处理意见" :span="2">{{ refundData.process_remark || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 原订单信息 -->
    <el-card shadow="never" class="detail-card">
      <template #header>
        <span>原订单信息</span>
      </template>

      <el-descriptions :column="3" border>
        <el-descriptions-item label="订单编号">
          <el-button type="primary" link @click="handleViewOrder">{{ orderData.order_no }}</el-button>
        </el-descriptions-item>
        <el-descriptions-item label="下单时间">{{ orderData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="订单状态">
          <el-tag :type="getOrderStatusTag(orderData.status)" size="small">
            {{ getOrderStatusText(orderData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="产品名称">{{ orderData.product_name }}</el-descriptions-item>
        <el-descriptions-item label="产品类型">{{ orderData.product_type || '-' }}</el-descriptions-item>
        <el-descriptions-item label="计费周期">{{ orderData.billing_cycle || '-' }}</el-descriptions-item>
        <el-descriptions-item label="订单金额">
          <span class="text-primary">¥{{ orderData.total_amount?.toFixed(2) || '0.00' }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="实付金额">
          <span class="text-primary">¥{{ orderData.paid_amount?.toFixed(2) || '0.00' }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="优惠金额">¥{{ orderData.discount_amount?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item label="支付方式">{{ orderData.payment_method || '-' }}</el-descriptions-item>
        <el-descriptions-item label="支付时间">{{ orderData.paid_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="交易单号">{{ orderData.transaction_no || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 退款处理流程 -->
    <el-card shadow="never" class="detail-card">
      <template #header>
        <span>退款处理流程</span>
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
            <div class="step-operator" v-if="step.operator">操作人: {{ step.operator }}</div>
          </div>
        </el-timeline-item>
      </el-timeline>
    </el-card>

    <!-- 审批对话框 -->
    <el-dialog v-model="approveDialogVisible" :title="approveForm.action === 'approve' ? '通过退款' : '拒绝退款'" width="500px">
      <el-alert
        v-if="approveForm.action === 'approve'"
        title="请确认退款金额将原路返回到用户支付账户"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 20px;"
      />
      <el-form :model="approveForm" :rules="approveRules" ref="approveFormRef" label-width="100px">
        <el-form-item label="退款编号">
          <el-input :value="refundData.refund_no" disabled />
        </el-form-item>
        <el-form-item label="退款金额">
          <el-input :value="'¥' + (refundData.amount?.toFixed(2) || '0.00')" disabled />
        </el-form-item>
        <el-form-item label="审批意见" prop="remark">
          <el-input v-model="approveForm.remark" type="textarea" :rows="3" :placeholder="approveForm.action === 'approve' ? '请输入审批意见（选填）' : '请输入拒绝原因（必填）'" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="approveDialogVisible = false">取消</el-button>
        <el-button
          :type="approveForm.action === 'approve' ? 'success' : 'danger'"
          @click="handleSubmitApprove"
          :loading="submitLoading"
        >
          {{ approveForm.action === 'approve' ? '确认通过' : '确认拒绝' }}
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
          callback(new Error('请输入拒绝原因'))
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
  const map: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    processing: '退款中',
    completed: '已退款',
    rejected: '已拒绝',
    cancelled: '已取消'
  }
  return map[status] || '未知'
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
  const map: Record<string, string> = {
    full: '全额退款',
    partial: '部分退款',
    cancellation: '取消退款'
  }
  return map[type] || '未知'
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
  const map: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    cancelled: '已取消',
    refunded: '已退款'
  }
  return map[status] || '未知'
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
    console.error('获取退款详情失败:', error)
    ElMessage.error('获取退款详情失败')
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
      await ElMessageBox.confirm('确定拒绝该退款申请吗？', '确认操作', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
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
      ElMessage.error('操作失败')
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
