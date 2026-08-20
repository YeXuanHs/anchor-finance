<template>
  <div class="order-detail-page">
    <h2>{{ $t('orderDetail.title') }}</h2>
    
    <div class="content-wrapper" v-loading="loading">
      <!-- 订单信息 -->
      <div class="info-section">
        <div class="info-grid">
          <!-- 左侧信息 -->
          <div class="info-column">
            <div class="info-item">
              <span class="label">{{ $t('orderDetail.client') }}</span>
              <router-link :to="`/customer-view/abstract?id=${order.client_id}`" class="value link">
                {{ order.client_name }}
              </router-link>
            </div>
            <div class="info-item">
              <span class="label">{{ $t('orderDetail.orderNo') }}</span>
              <span class="value">{{ order.id }}</span>
            </div>
            <div class="info-item">
              <span class="label">{{ $t('orderDetail.time') }}</span>
              <span class="value">{{ order.created_at }}</span>
            </div>
            <div class="info-item">
              <span class="label">{{ $t('orderDetail.couponCode') }}</span>
              <span class="value">{{ order.coupon_code || $t('orderDetail.notFilled') }}</span>
            </div>
            <div class="info-item">
              <span class="label">{{ $t('orderDetail.ipAddress') }}</span>
              <span class="value">{{ order.ip || $t('orderDetail.notFilled') }}</span>
            </div>
          </div>
          
          <!-- 右侧信息 -->
          <div class="info-column">
            <div class="info-item">
              <span class="label">{{ $t('orderDetail.billInfo') }}</span>
              <router-link :to="`/bill-detail?id=${order.bill_id}&uid=${order.client_id}`" class="value link">
                {{ order.bill_no }}
              </router-link>
            </div>
            <div class="info-item">
              <span class="label">{{ $t('orderDetail.payMethod') }}</span>
              <span class="value">{{ order.pay_method || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">{{ $t('orderDetail.amount') }}</span>
              <span class="value amount">❖{{ order.amount?.toFixed(2) || '0.00' }}</span>
            </div>
            <div class="info-item">
              <span class="label">{{ $t('orderDetail.status') }}</span>
              <el-select v-model="order.status" :placeholder="$t('orderDetail.selectStatus')" @change="handleStatusChange" style="width: 120px">
                <el-option :label="$t('orderDetail.pending')" :value="0" />
                <el-option :label="$t('orderDetail.activated')" :value="1" />
                <el-option :label="$t('orderDetail.cancelled')" :value="2" />
              </el-select>
            </div>
            <div class="info-item">
              <span class="label">{{ $t('orderDetail.clientRemark') }}</span>
              <span class="value">{{ order.client_remark || $t('orderDetail.notFilled') }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 订单项目 -->
      <div class="section">
        <h3>{{ $t('orderDetail.orderItems') }}</h3>
        <el-table :data="orderItems" style="width: 100%">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="item_name" :label="$t('orderDetail.item')" min-width="150" />
          <el-table-column prop="description" :label="$t('orderDetail.description')" min-width="200" />
          <el-table-column prop="billing_cycle" :label="$t('orderDetail.billingCycle')" width="120" />
          <el-table-column prop="amount" :label="$t('orderDetail.amount')" width="120">
            <template #default="{ row }">
              <span class="amount">❖{{ row.amount?.toFixed(2) || '0.00' }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="status" :label="$t('orderDetail.status')" width="100">
            <template #default="{ row }">
              <span :class="getStatusClass(row.status)">{{ getStatusText(row.status) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="pay_status" :label="$t('orderDetail.payStatus')" width="100">
            <template #default="{ row }">
              <span :class="getPayStatusClass(row.pay_status)">{{ getPayStatusText(row.pay_status) }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="$t('orderDetail.operations')" width="120">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleViewItem(row)">{{ $t('orderDetail.view') }}</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="total">
          {{ $t('orderDetail.total') }}: ❖{{ orderTotal.toFixed(2) }}
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <el-button type="success" @click="handleVerify" :loading="actionLoading">
          <el-icon><Check /></el-icon>
          {{ $t('orderDetail.verify') }}
        </el-button>
        <el-button type="warning" @click="handleCancel" :loading="actionLoading">{{ $t('orderDetail.cancelOrder') }}</el-button>
        <el-button type="danger" @click="handleDelete" :loading="actionLoading">{{ $t('orderDetail.deleteOrder') }}</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Check } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const actionLoading = ref(false)
const order = ref<any>({})
const orderItems = ref<any[]>([])

const STATUS_CLASS_MAP: Record<number, string> = {
  0: 'status-pending',
  1: 'status-active',
  2: 'status-cancelled'
}

const PAY_STATUS_CLASS_MAP: Record<number, string> = {
  0: 'pay-unpaid',
  1: 'pay-paid',
  2: 'pay-refunded'
}

const statusKeyMap: Record<number, string> = {
  0: 'orderDetail.pending',
  1: 'orderDetail.activated',
  2: 'orderDetail.cancelled'
}

const payStatusKeyMap: Record<number, string> = {
  0: 'orderDetail.unpaid',
  1: 'orderDetail.paid',
  2: 'orderDetail.refunded'
}

const getStatusText = (status: number) => {
  return statusKeyMap[status] ? $t(statusKeyMap[status]) : $t('orderDetail.unknown')
}

const getStatusClass = (status: number) => {
  return STATUS_CLASS_MAP[status] || ''
}

const getPayStatusText = (status: number) => {
  return payStatusKeyMap[status] ? $t(payStatusKeyMap[status]) : $t('orderDetail.unknown')
}

const getPayStatusClass = (status: number) => {
  return PAY_STATUS_CLASS_MAP[status] || ''
}

const orderTotal = computed(() => {
  return orderItems.value.reduce((sum, item) => sum + Number(item.amount || 0), 0)
})

const fetchOrder = async () => {
  const id = route.query.id || route.params.id
  if (!id) return
  
  loading.value = true
  try {
    const data = await request.get({ url: `/api/admin/orders/${id}` })
    order.value = data.order || data
    orderItems.value = data.items || []
  } catch (error) {
    console.error($t('orderDetail.fetchFailed') + ':', error)
    ElMessage.error($t('orderDetail.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleStatusChange = async (val: number) => {
  try {
    await request.put({
      url: `/api/admin/orders/${order.value.id}/status`,
      data: { status: val }
    })
    ElMessage.success($t('orderDetail.statusUpdateSuccess'))
  } catch (error) {
    ElMessage.error($t('orderDetail.statusUpdateFailed'))
    fetchOrder()
  }
}

const handleVerify = async () => {
  try {
    await ElMessageBox.confirm($t('orderDetail.verifyConfirmMsg'), $t('orderDetail.verifyConfirmTitle'), {
      confirmButtonText: $t('common.confirm'),
      cancelButtonText: $t('common.cancel'),
      type: 'warning'
    })
    actionLoading.value = true
    await request.post({
      url: `/api/admin/orders/${order.value.id}/verify`
    })
    ElMessage.success($t('orderDetail.verifySuccess'))
    fetchOrder()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error($t('orderDetail.verifyFailed'))
    }
  } finally {
    actionLoading.value = false
  }
}

const handleCancel = async () => {
  try {
    await ElMessageBox.confirm($t('orderDetail.cancelConfirmMsg'), $t('orderDetail.cancelConfirmTitle'), {
      confirmButtonText: $t('common.confirm'),
      cancelButtonText: $t('common.cancel'),
      type: 'warning'
    })
    actionLoading.value = true
    await request.post({
      url: `/api/admin/orders/${order.value.id}/cancel`
    })
    ElMessage.success($t('orderDetail.cancelSuccess'))
    fetchOrder()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error($t('orderDetail.cancelFailed'))
    }
  } finally {
    actionLoading.value = false
  }
}

const handleDelete = async () => {
  try {
    await ElMessageBox.confirm($t('orderDetail.deleteConfirmMsg'), $t('orderDetail.deleteConfirmTitle'), {
      confirmButtonText: $t('common.confirm'),
      cancelButtonText: $t('common.cancel'),
      type: 'warning'
    })
    actionLoading.value = true
    await request.del({
      url: `/api/admin/orders/${order.value.id}`
    })
    ElMessage.success($t('orderDetail.deleteSuccess'))
    router.back()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error($t('orderDetail.deleteFailed'))
    }
  } finally {
    actionLoading.value = false
  }
}

const handleViewItem = (row: any) => {
  router.push(`/customer-view/${row.client_id}?tab=products`)
}

onMounted(() => {
  fetchOrder()
})
</script>

<style scoped lang="scss">
.order-detail-page {
  padding: 20px;
}

h2 {
  margin: 0 0 20px;
  font-size: 18px;
}

.content-wrapper {
  background: #fff;
}

.info-section {
  margin-bottom: 24px;
}

.info-grid {
  display: flex;
  gap: 40px;
}

.info-column {
  flex: 1;
}

.info-item {
  display: flex;
  margin-bottom: 12px;
  
  .label {
    width: 80px;
    color: var(--el-text-color-secondary);
    font-size: 14px;
    flex-shrink: 0;
  }
  
  .value {
    flex: 1;
    font-size: 14px;
    color: var(--el-text-color-primary);
    
    &.link {
      color: var(--el-color-primary);
      text-decoration: none;
      
      &:hover {
        text-decoration: underline;
      }
    }
    
    &.amount {
      font-weight: 600;
      color: #f56c6c;
    }
  }
}

.section {
  margin-bottom: 24px;
  
  h3 {
    font-size: 16px;
    margin: 0 0 16px;
  }
}

.amount {
  font-weight: 600;
  color: #f56c6c;
}

.status-pending {
  color: #e6a23c;
}

.status-active {
  color: #67c23a;
}

.status-cancelled {
  color: #909399;
}

.pay-unpaid {
  color: #e6a23c;
}

.pay-paid {
  color: #67c23a;
}

.pay-refunded {
  color: #f56c6c;
}

.total {
  padding: 12px 0;
  font-size: 14px;
  color: #606266;
  text-align: right;
}

.action-buttons {
  display: flex;
  gap: 12px;
  padding: 16px 0;
  border-top: 1px solid var(--el-border-color-lighter);
}
</style>
