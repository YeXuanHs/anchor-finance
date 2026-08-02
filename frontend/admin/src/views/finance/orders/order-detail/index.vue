<template>
  <div class="order-detail-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>订单详情</span>
          <div class="header-actions">
            <el-button @click="handleBack">
              <el-icon><Back /></el-icon>
              返回
            </el-button>
          </div>
        </div>
      </template>

      <div v-loading="loading" class="loading-container">
        <div v-if="order" class="order-content">
          <!-- 订单基本信息 -->
          <el-descriptions :column="2" border class="order-info">
            <el-descriptions-item label="订单编号">{{ order.order_no }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getOrderStatusType(order.status)" size="default">
                {{ getOrderStatusText(order.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="关联客户">
              <el-button type="primary" link @click="handleViewClient">
                {{ order.client_name }}
              </el-button>
            </el-descriptions-item>
            <el-descriptions-item label="客户邮箱">{{ order.client_email || '-' }}</el-descriptions-item>
            <el-descriptions-item label="产品/服务" :span="2">{{ order.product_name }}</el-descriptions-item>
            <el-descriptions-item label="订单金额">
              <span class="amount-text">¥{{ formatAmount(order.amount) }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="实付金额">
              <span class="amount-text">¥{{ formatAmount(order.paid_amount || order.amount) }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="支付方式">{{ order.pay_method || '-' }}</el-descriptions-item>
            <el-descriptions-item label="支付流水号">{{ order.payment_no || '-' }}</el-descriptions-item>
            <el-descriptions-item label="下单时间">{{ order.created_at }}</el-descriptions-item>
            <el-descriptions-item label="更新时间">{{ order.updated_at || '-' }}</el-descriptions-item>
            <el-descriptions-item label="备注" :span="2">{{ order.remark || '无' }}</el-descriptions-item>
          </el-descriptions>

          <!-- 关联产品信息 -->
          <div class="section">
            <h3>产品信息</h3>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="产品名称">{{ order.product_name }}</el-descriptions-item>
              <el-descriptions-item label="产品ID">{{ order.product_id }}</el-descriptions-item>
              <el-descriptions-item label="产品类型">{{ order.product_type || '-' }}</el-descriptions-item>
              <el-descriptions-item label="计费周期">{{ order.billing_cycle || '-' }}</el-descriptions-item>
              <el-descriptions-item label="域名">{{ order.domain || '-' }}</el-descriptions-item>
              <el-descriptions-item label="IP地址">{{ order.ip || '-' }}</el-descriptions-item>
              <el-descriptions-item label="规格配置">{{ order.config || '-' }}</el-descriptions-item>
              <el-descriptions-item label="到期时间">{{ order.due_date || '-' }}</el-descriptions-item>
            </el-descriptions>
          </div>

          <!-- 订单操作日志 -->
          <div class="section">
            <h3>操作日志</h3>
            <el-timeline>
              <el-timeline-item
                v-for="log in logs"
                :key="log.id"
                :timestamp="log.created_at"
                :type="log.type || 'primary'"
                placement="top"
              >
                <div class="log-content">
                  <span class="log-user">{{ log.username }}</span>
                  <span class="log-action">{{ log.action }}</span>
                  <div class="log-detail" v-if="log.description">{{ log.description }}</div>
                </div>
              </el-timeline-item>
            </el-timeline>
            <el-empty v-if="!logs.length" description="暂无操作日志" />
          </div>
        </div>
      </div>
    </el-card>

    <!-- 底部操作栏 -->
    <div class="action-bar" v-if="order">
      <el-button
        v-if="order.status === 2"
        type="primary"
        @click="handleActivate"
        :loading="actionLoading"
      >
        <el-icon><Check /></el-icon>
        开通服务
      </el-button>
      <el-button
        v-if="order.status === 3"
        type="warning"
        @click="handleSuspend"
        :loading="actionLoading"
      >
        <el-icon><VideoPause /></el-icon>
        暂停服务
      </el-button>
      <el-button
        v-if="order.status === 3"
        type="danger"
        @click="handleTerminate"
        :loading="actionLoading"
      >
        <el-icon><CircleClose /></el-icon>
        终止服务
      </el-button>
      <el-button
        v-if="order.status === 3"
        type="success"
        @click="handleRenew"
        :loading="actionLoading"
      >
        <el-icon><Refresh /></el-icon>
        续费
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Back, Check, VideoPause, CircleClose, Refresh } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const actionLoading = ref(false)
const order = ref<any>({})
const logs = ref<any[]>([])

const ORDER_STATUS_MAP: Record<number, { text: string; type: string }> = {
  0: { text: '待付款', type: 'warning' },
  1: { text: '待审核', type: 'primary' },
  2: { text: '审核通过', type: 'success' },
  3: { text: '已开通', type: 'success' },
  4: { text: '已完成', type: 'info' },
  5: { text: '已取消', type: 'info' },
  6: { text: '已退款', type: 'danger' }
}

const getOrderStatusText = (status: number) => {
  return ORDER_STATUS_MAP[status]?.text || '未知'
}

const getOrderStatusType = (status: number) => {
  return (ORDER_STATUS_MAP[status]?.type || 'info') as any
}

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const fetchOrder = async () => {
  const id = route.params.id
  if (!id) return

  loading.value = true
  try {
    const data = await request.get({ url: `/api/admin/orders/${id}` })
    order.value = data.order || data
  } catch (error) {
    console.error('获取订单详情失败:', error)
    ElMessage.error('获取订单详情失败')
  } finally {
    loading.value = false
  }
}

const fetchLogs = async () => {
  const id = route.params.id
  if (!id) return

  try {
    const data = await request.get({ url: `/api/admin/orders/${id}/logs` })
    logs.value = data || []
  } catch (error) {
    console.error('获取操作日志失败:', error)
  }
}

const handleBack = () => {
  router.back()
}

const handleViewClient = () => {
  if (order.value?.client_id) {
    router.push(`/finance/clients/detail/${order.value.client_id}`)
  }
}

const handleActivate = async () => {
  const id = route.params.id
  if (!id) return

  try {
    await ElMessageBox.confirm('确定要开通该订单的服务吗？', '开通确认', { type: 'warning' })
    actionLoading.value = true
    await request.post({ url: `/api/admin/orders/${id}/activate` })
    ElMessage.success('服务开通成功')
    fetchOrder()
    fetchLogs()
  } catch (error: any) {
    if (error !== 'cancel') ElMessage.error('开通失败')
  } finally {
    actionLoading.value = false
  }
}

const handleSuspend = async () => {
  const id = route.params.id
  if (!id) return

  try {
    await ElMessageBox.confirm('确定要暂停该服务吗？', '暂停确认', { type: 'warning' })
    actionLoading.value = true
    await request.post({ url: `/api/admin/orders/${id}/suspend` })
    ElMessage.success('服务已暂停')
    fetchOrder()
    fetchLogs()
  } catch (error: any) {
    if (error !== 'cancel') ElMessage.error('暂停失败')
  } finally {
    actionLoading.value = false
  }
}

const handleTerminate = async () => {
  const id = route.params.id
  if (!id) return

  try {
    await ElMessageBox.confirm('确定要终止该服务吗？此操作不可恢复。', '终止确认', { type: 'danger' })
    actionLoading.value = true
    await request.post({ url: `/api/admin/orders/${id}/terminate` })
    ElMessage.success('服务已终止')
    fetchOrder()
    fetchLogs()
  } catch (error: any) {
    if (error !== 'cancel') ElMessage.error('终止失败')
  } finally {
    actionLoading.value = false
  }
}

const handleRenew = async () => {
  const id = route.params.id
  if (!id) return

  try {
    await ElMessageBox.confirm('确定要为该服务续费吗？', '续费确认', { type: 'warning' })
    actionLoading.value = true
    await request.post({ url: `/api/admin/orders/${id}/renew` })
    ElMessage.success('续费成功')
    fetchOrder()
    fetchLogs()
  } catch (error: any) {
    if (error !== 'cancel') ElMessage.error('续费失败')
  } finally {
    actionLoading.value = false
  }
}

onMounted(() => {
  fetchOrder()
  fetchLogs()
})
</script>

<style scoped lang="scss">
.order-detail-page {
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

.order-info {
  margin-bottom: 24px;
}

.amount-text {
  font-weight: 600;
  color: #f56c6c;
}

.section {
  margin-top: 24px;

  h3 {
    margin: 0 0 16px;
    font-size: 16px;
    font-weight: 600;
  }
}

.log-content {
  .log-user {
    font-weight: 500;
    margin-right: 8px;
  }

  .log-action {
    color: var(--el-text-color-secondary);
  }

  .log-detail {
    margin-top: 4px;
    color: var(--el-text-color-regular);
    font-size: 13px;
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