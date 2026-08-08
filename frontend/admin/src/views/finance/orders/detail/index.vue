<template>
  <div class="order-detail-page">
    <!-- 订单信息卡片 -->
    <el-card shadow="never" class="order-info-card">
      <template #header>
        <div class="card-header">
          <span>订单信息</span>
          <div class="header-actions">
            <el-button
              v-if="order.status === 'pending_payment'"
              type="success"
              @click="handleConfirmPayment"
            >
              确认付款
            </el-button>
            <el-button
              v-if="order.status === 'pending_activation'"
              type="primary"
              @click="handleActivate"
            >
              开通服务
            </el-button>
            <el-button
              v-if="order.status !== 'cancelled' && order.status !== 'refunded'"
              type="danger"
              @click="handleCancel"
            >
              取消订单
            </el-button>
            <el-button @click="$router.back()">返回</el-button>
          </div>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="订单号">{{ order.order_no }}</el-descriptions-item>
            <el-descriptions-item label="订单类型">
              <el-tag size="small">{{ getTypeText(order.type) }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="客户">
              <el-button type="primary" link @click="$router.push(`/customer-view/${order.client_id}`)">
                {{ order.client_name }}
              </el-button>
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(order.status)" size="small">
                {{ getStatusText(order.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="下单时间">{{ order.created_at }}</el-descriptions-item>
            <el-descriptions-item label="支付时间">{{ order.paid_at || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-col>
        <el-col :span="12">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="订单金额">¥{{ formatMoney(order.amount) }}</el-descriptions-item>
            <el-descriptions-item label="支付方式">{{ order.payment_method || '-' }}</el-descriptions-item>
            <el-descriptions-item label="交易号">{{ order.trade_no || '-' }}</el-descriptions-item>
            <el-descriptions-item label="支付状态">
              <el-tag :type="order.paid ? 'success' : 'warning'" size="small">
                {{ order.paid ? '已支付' : '未支付' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="备注" :span="2">{{ order.notes || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-col>
      </el-row>
    </el-card>

    <!-- 产品信息 -->
    <el-card shadow="never" class="product-card">
      <template #header>
        <span>产品信息</span>
      </template>
      <el-table :data="order.items || []" border stripe>
        <el-table-column prop="product_name" label="产品名称" min-width="200" />
        <el-table-column prop="billing_cycle" label="计费周期" width="120" />
        <el-table-column prop="quantity" label="数量" width="80" align="center" />
        <el-table-column prop="unit_price" label="单价" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.unit_price) }}</template>
        </el-table-column>
        <el-table-column prop="amount" label="小计" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.amount) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getItemStatusType(row.status)" size="small">
              {{ getItemStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <div class="order-total">
        <span>订单总计：</span>
        <span class="total-amount">¥{{ formatMoney(order.amount) }}</span>
      </div>
    </el-card>

    <!-- 服务详情 -->
    <el-card shadow="never" class="service-card" v-if="order.service">
      <template #header>
        <span>服务详情</span>
      </template>
      <el-descriptions :column="3" border>
        <el-descriptions-item label="服务ID">{{ order.service.id }}</el-descriptions-item>
        <el-descriptions-item label="产品名称">{{ order.service.product_name }}</el-descriptions-item>
        <el-descriptions-item label="域名/标识">{{ order.service.domain || '-' }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ order.service.ip || '-' }}</el-descriptions-item>
        <el-descriptions-item label="开通时间">{{ order.service.activated_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="到期时间">{{ order.service.expired_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getServiceStatusType(order.service.status)" size="small">
            {{ getServiceStatusText(order.service.status) }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 操作日志 -->
    <el-card shadow="never" class="log-card">
      <template #header>
        <span>操作日志</span>
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
          <div class="log-operator">操作人: {{ log.operator_name }}</div>
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
  const map: Record<string, string> = { new: '新购', renewal: '续费', upgrade: '升级', refund: '退款' }
  return map[type] || '未知'
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
  const map: Record<string, string> = {
    pending_payment: '待付款',
    pending_activation: '待开通',
    active: '进行中',
    completed: '已完成',
    cancelled: '已取消',
    refunded: '已退款'
  }
  return map[status] || '未知'
}

// 项目状态
const getItemStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { pending: 'warning', active: 'success', suspended: 'danger', terminated: 'info' }
  return map[status] || 'info'
}

const getItemStatusText = (status: string) => {
  const map: Record<string, string> = { pending: '待开通', active: '正常', suspended: '暂停', terminated: '已终止' }
  return map[status] || '未知'
}

// 服务状态
const getServiceStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { active: 'success', suspended: 'danger', pending: 'warning', terminated: 'info' }
  return map[status] || 'info'
}

const getServiceStatusText = (status: string) => {
  const map: Record<string, string> = { active: '正常', suspended: '暂停', pending: '待开通', terminated: '已终止' }
  return map[status] || '未知'
}

// 获取订单详情
const fetchOrder = async () => {
  try {
    const data = await request.get({ url: `/api/admin/orders/${orderId}` })
    order.value = data
  } catch (error) {
    console.error('获取订单详情失败:', error)
  }
}

// 确认付款
const handleConfirmPayment = async () => {
  try {
    await ElMessageBox.confirm('确定要确认此订单已付款吗？', '确认付款', { type: 'warning' })
    await request.post({ url: `/api/admin/orders/${orderId}/confirm-payment` })
    ElMessage.success('确认付款成功')
    fetchOrder()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('确认付款失败:', error)
    }
  }
}

// 开通服务
const handleActivate = async () => {
  try {
    await ElMessageBox.confirm('确定要开通此订单的服务吗？', '开通服务', { type: 'warning' })
    await request.post({ url: `/api/admin/orders/${orderId}/activate` })
    ElMessage.success('服务开通成功')
    fetchOrder()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('开通服务失败:', error)
    }
  }
}

// 取消订单
const handleCancel = async () => {
  try {
    await ElMessageBox.confirm('确定要取消此订单吗？', '取消订单', { type: 'warning' })
    await request.post({ url: `/api/admin/orders/${orderId}/cancel` })
    ElMessage.success('订单已取消')
    fetchOrder()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('取消订单失败:', error)
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
