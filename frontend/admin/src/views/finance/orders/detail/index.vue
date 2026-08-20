<template>
  <div class="order-detail-page">
    <!-- 订单信息卡片 -->
    <el-card shadow="never" class="order-info-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('orderDetail.title') }}</span>
          <div class="header-actions">
            <el-button
              v-if="order.status === 'pending_payment'"
              type="success"
              @click="handleConfirmPayment"
            >
              {{ $t('orderDetail.confirmPayment') }}
            </el-button>
            <el-button
              v-if="order.status === 'pending_activation'"
              type="primary"
              @click="handleActivate"
            >
              {{ $t('orderDetail.activateService') }}
            </el-button>
            <el-button
              v-if="order.status !== 'cancelled' && order.status !== 'refunded'"
              type="danger"
              @click="handleCancel"
            >
              {{ $t('orderDetail.cancelOrder') }}
            </el-button>
            <el-button @click="$router.back()">{{ $t('common.back') }}</el-button>
          </div>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-descriptions :column="2" border>
            <el-descriptions-item :label="$t('orderDetail.orderNo')">{{ order.order_no }}</el-descriptions-item>
            <el-descriptions-item :label="$t('orderDetail.orderType')">
              <el-tag size="small">{{ getTypeText(order.type) }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('orderDetail.client')">
              <el-button type="primary" link @click="$router.push(`/customer-view/${order.client_id}`)">
                {{ order.client_name }}
              </el-button>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('common.status')">
              <el-tag :type="getStatusType(order.status)" size="small">
                {{ getStatusText(order.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('orderDetail.createdAt')">{{ order.created_at }}</el-descriptions-item>
            <el-descriptions-item :label="$t('orderDetail.paidAt')">{{ order.paid_at || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-col>
        <el-col :span="12">
          <el-descriptions :column="2" border>
            <el-descriptions-item :label="$t('orderDetail.orderAmount')">¥{{ formatMoney(order.amount) }}</el-descriptions-item>
            <el-descriptions-item :label="$t('orderDetail.paymentMethod')">{{ order.payment_method || '-' }}</el-descriptions-item>
            <el-descriptions-item :label="$t('orderDetail.tradeNo')">{{ order.trade_no || '-' }}</el-descriptions-item>
            <el-descriptions-item :label="$t('orderDetail.paymentStatus')">
              <el-tag :type="order.paid ? 'success' : 'warning'" size="small">
                {{ order.paid ? $t('orderDetail.paid') : $t('orderDetail.unpaid') }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('common.remark')" :span="2">{{ order.notes || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-col>
      </el-row>
    </el-card>

    <!-- 产品信息 -->
    <el-card shadow="never" class="product-card">
      <template #header>
        <span>{{ $t('orderDetail.productInfo') }}</span>
      </template>
      <el-table :data="order.items || []" border stripe>
        <el-table-column prop="product_name" :label="$t('orderDetail.productName')" min-width="200" />
        <el-table-column prop="billing_cycle" :label="$t('orderDetail.billingCycle')" width="120" />
        <el-table-column prop="quantity" :label="$t('common.quantity')" width="80" align="center" />
        <el-table-column prop="unit_price" :label="$t('orderDetail.unitPrice')" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.unit_price) }}</template>
        </el-table-column>
        <el-table-column prop="amount" :label="$t('orderDetail.subtotal')" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.amount) }}</template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('common.status')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getItemStatusType(row.status)" size="small">
              {{ getItemStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <div class="order-total">
        <span>{{ $t('orderDetail.orderTotal') }}：</span>
        <span class="total-amount">¥{{ formatMoney(order.amount) }}</span>
      </div>
    </el-card>

    <!-- 服务详情 -->
    <el-card shadow="never" class="service-card" v-if="order.service">
      <template #header>
        <span>{{ $t('orderDetail.serviceDetail') }}</span>
      </template>
      <el-descriptions :column="3" border>
        <el-descriptions-item :label="$t('orderDetail.serviceId')">{{ order.service.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('orderDetail.productName')">{{ order.service.product_name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('orderDetail.domain')">{{ order.service.domain || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('orderDetail.ipAddress')">{{ order.service.ip || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('orderDetail.activatedAt')">{{ order.service.activated_at || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('orderDetail.expiredAt')">{{ order.service.expired_at || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('common.status')">
          <el-tag :type="getServiceStatusType(order.service.status)" size="small">
            {{ getServiceStatusText(order.service.status) }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 操作日志 -->
    <el-card shadow="never" class="log-card">
      <template #header>
        <span>{{ $t('orderDetail.operationLog') }}</span>
      </template>
      <el-timeline>
        <el-timeline-item
          v-for="log in order.logs || []"
          :key="log.id"
          :timestamp="log.created_at"
          placement="top"
        >
          <div class="log-content">
            <span class="log-action">{{ log.action }}</span>
            <span class="log-detail">{{ log.detail }}</span>
          </div>
          <div class="log-operator">{{ $t('orderDetail.operator') }}: {{ log.operator_name }}</div>
        </el-timeline-item>
      </el-timeline>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const route = useRoute()
const router = useRouter()
const orderId = route.params.id

// 订单信息
const order = ref<any>({})

// 格式化金额
const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// 订单类型
const getTypeText = (type: string) => {
  const map: Record<string, () => string> = {
    new: () => $t('orderDetail.typeNew'),
    renewal: () => $t('orderDetail.typeRenewal'),
    upgrade: () => $t('orderDetail.typeUpgrade'),
    refund: () => $t('orderDetail.typeRefund')
  }
  return map[type]?.() || $t('common.unknown')
}

// 订单状态
const getStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    pending_payment: 'warning',
    pending_activation: 'primary',
    active: 'success',
    completed: 'success',
    cancelled: 'info',
    refunded: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, () => string> = {
    pending_payment: () => $t('orderDetail.statusPendingPayment'),
    pending_activation: () => $t('orderDetail.statusPendingActivation'),
    active: () => $t('orderDetail.statusActive'),
    completed: () => $t('orderDetail.statusCompleted'),
    cancelled: () => $t('orderDetail.statusCancelled'),
    refunded: () => $t('orderDetail.statusRefunded')
  }
  return map[status]?.() || $t('common.unknown')
}

// 项目状态
const getItemStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { pending: 'warning', active: 'success', suspended: 'danger', terminated: 'info' }
  return map[status] || 'info'
}

const getItemStatusText = (status: string) => {
  const map: Record<string, () => string> = {
    pending: () => $t('orderDetail.itemStatusPending'),
    active: () => $t('orderDetail.itemStatusActive'),
    suspended: () => $t('orderDetail.itemStatusSuspended'),
    terminated: () => $t('orderDetail.itemStatusTerminated')
  }
  return map[status]?.() || $t('common.unknown')
}

// 服务状态
const getServiceStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { active: 'success', suspended: 'danger', pending: 'warning', terminated: 'info' }
  return map[status] || 'info'
}

const getServiceStatusText = (status: string) => {
  const map: Record<string, () => string> = {
    active: () => $t('orderDetail.serviceStatusActive'),
    suspended: () => $t('orderDetail.serviceStatusSuspended'),
    pending: () => $t('orderDetail.serviceStatusPending'),
    terminated: () => $t('orderDetail.serviceStatusTerminated')
  }
  return map[status]?.() || $t('common.unknown')
}

// 获取订单详情
const fetchOrder = async () => {
  try {
    const data = await request.get({ url: `/api/admin/orders/${orderId}` })
    order.value = data
  } catch (error) {
    console.error('fetch order detail failed:', error)
  }
}

// 确认付款
const handleConfirmPayment = async () => {
  try {
    await ElMessageBox.confirm($t('orderDetail.confirmPaymentMsg'), $t('orderDetail.confirmPaymentTitle'), { type: 'warning' })
    await request.post({ url: `/api/admin/orders/${orderId}/confirm-payment` })
    ElMessage.success($t('orderDetail.paymentConfirmed'))
    fetchOrder()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('confirm payment failed:', error)
    }
  }
}

// 开通服务
const handleActivate = async () => {
  try {
    await ElMessageBox.confirm($t('orderDetail.confirmActivateMsg'), $t('orderDetail.confirmActivateTitle'), { type: 'warning' })
    await request.post({ url: `/api/admin/orders/${orderId}/activate` })
    ElMessage.success($t('orderDetail.serviceActivated'))
    fetchOrder()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('activate service failed:', error)
    }
  }
}

// 取消订单
const handleCancel = async () => {
  try {
    await ElMessageBox.confirm($t('orderDetail.confirmCancelMsg'), $t('orderDetail.confirmCancelTitle'), { type: 'warning' })
    await request.post({ url: `/api/admin/orders/${orderId}/cancel` })
    ElMessage.success($t('orderDetail.orderCancelled'))
    fetchOrder()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('cancel order failed:', error)
    }
  }
}

onMounted(() => {
  fetchOrder()
})
</script>

<style scoped lang="scss">
.order-detail-page {
  padding: 16px;
}

.order-info-card,
.product-card,
.service-card,
.log-card {
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.order-total {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding: 16px;
  font-size: 16px;
}

.total-amount {
  font-size: 24px;
  font-weight: 600;
  color: #F59E0B;
  margin-left: 8px;
}

.log-content {
  display: flex;
  gap: 8px;
}

.log-action {
  font-weight: 500;
}

.log-detail {
  color: #86909C;
}

.log-operator {
  margin-top: 8px;
  font-size: 12px;
  color: #86909C;
}
</style>
