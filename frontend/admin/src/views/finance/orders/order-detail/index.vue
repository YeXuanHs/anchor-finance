<template>
  <div class="order-detail-page">
    <h2>订单详情</h2>
    
    <div class="content-wrapper" v-loading="loading">
      <!-- 订单信息 -->
      <div class="info-section">
        <div class="info-grid">
          <!-- 左侧信息 -->
          <div class="info-column">
            <div class="info-item">
              <span class="label">客户</span>
              <router-link :to="`/customer-view/abstract?id=${order.client_id}`" class="value link">
                {{ order.client_name }}
              </router-link>
            </div>
            <div class="info-item">
              <span class="label">订单号</span>
              <span class="value">{{ order.id }}</span>
            </div>
            <div class="info-item">
              <span class="label">时间</span>
              <span class="value">{{ order.created_at }}</span>
            </div>
            <div class="info-item">
              <span class="label">优惠码</span>
              <span class="value">{{ order.coupon_code || '未填写' }}</span>
            </div>
            <div class="info-item">
              <span class="label">IP地址</span>
              <span class="value">{{ order.ip || '未填写' }}</span>
            </div>
          </div>
          
          <!-- 右侧信息 -->
          <div class="info-column">
            <div class="info-item">
              <span class="label">账单信息</span>
              <router-link :to="`/bill-detail?id=${order.bill_id}&uid=${order.client_id}`" class="value link">
                {{ order.bill_no }}
              </router-link>
            </div>
            <div class="info-item">
              <span class="label">付款方式</span>
              <span class="value">{{ order.pay_method || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">金额</span>
              <span class="value amount">❖{{ order.amount?.toFixed(2) || '0.00' }}</span>
            </div>
            <div class="info-item">
              <span class="label">状态</span>
              <el-select v-model="order.status" placeholder="请选择" @change="handleStatusChange" style="width: 120px">
                <el-option label="待核验" :value="0" />
                <el-option label="已激活" :value="1" />
                <el-option label="已取消" :value="2" />
              </el-select>
            </div>
            <div class="info-item">
              <span class="label">客户备注</span>
              <span class="value">{{ order.client_remark || '未填写' }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 订单项目 -->
      <div class="section">
        <h3>订单项目</h3>
        <el-table :data="orderItems" style="width: 100%">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="item_name" label="条目" min-width="150" />
          <el-table-column prop="description" label="描述" min-width="200" />
          <el-table-column prop="billing_cycle" label="付款周期" width="120" />
          <el-table-column prop="amount" label="金额" width="120">
            <template #default="{ row }">
              <span class="amount">❖{{ row.amount?.toFixed(2) || '0.00' }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <span :class="getStatusClass(row.status)">{{ getStatusText(row.status) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="pay_status" label="付款状态" width="100">
            <template #default="{ row }">
              <span :class="getPayStatusClass(row.pay_status)">{{ getPayStatusText(row.pay_status) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleViewItem(row)">查看</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="total">
          合计: ❖{{ orderTotal.toFixed(2) }}
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <el-button type="success" @click="handleVerify" :loading="actionLoading">
          <el-icon><Check /></el-icon>
          核验通过
        </el-button>
        <el-button type="warning" @click="handleCancel" :loading="actionLoading">取消订单</el-button>
        <el-button type="danger" @click="handleDelete" :loading="actionLoading">删除订单</el-button>
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

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const actionLoading = ref(false)
const order = ref<any>({})
const orderItems = ref<any[]>([])

// 状态映射
const STATUS_MAP: Record<number, { text: string; class: string }> = {
  0: { text: '待核验', class: 'status-pending' },
  1: { text: '已激活', class: 'status-active' },
  2: { text: '已取消', class: 'status-cancelled' }
}

// 付款状态映射
const PAY_STATUS_MAP: Record<number, { text: string; class: string }> = {
  0: { text: '未支付', class: 'pay-unpaid' },
  1: { text: '已支付', class: 'pay-paid' },
  2: { text: '已退款', class: 'pay-refunded' }
}

// 获取状态文本
const getStatusText = (status: number) => {
  return STATUS_MAP[status]?.text || '未知'
}

// 获取状态样式类
const getStatusClass = (status: number) => {
  return STATUS_MAP[status]?.class || ''
}

// 获取付款状态文本
const getPayStatusText = (status: number) => {
  return PAY_STATUS_MAP[status]?.text || '未知'
}

// 获取付款状态样式类
const getPayStatusClass = (status: number) => {
  return PAY_STATUS_MAP[status]?.class || ''
}

// 订单总计
const orderTotal = computed(() => {
  return orderItems.value.reduce((sum, item) => sum + Number(item.amount || 0), 0)
})

// 获取订单详情
const fetchOrder = async () => {
  const id = route.query.id || route.params.id
  if (!id) return
  
  loading.value = true
  try {
    const data = await request.get({ url: `/api/admin/orders/${id}` })
    order.value = data.order || data
    orderItems.value = data.items || []
  } catch (error) {
    console.error('获取订单详情失败:', error)
    ElMessage.error('获取订单详情失败')
  } finally {
    loading.value = false
  }
}

// 状态变化
const handleStatusChange = async (val: number) => {
  try {
    await request.put({
      url: `/api/admin/orders/${order.value.id}/status`,
      data: { status: val }
    })
    ElMessage.success('状态更新成功')
  } catch (error) {
    ElMessage.error('状态更新失败')
    fetchOrder()
  }
}

// 核验通过
const handleVerify = async () => {
  try {
    await ElMessageBox.confirm('确定要核验通过该订单吗？', '核验确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    actionLoading.value = true
    await request.post({
      url: `/api/admin/orders/${order.value.id}/verify`
    })
    ElMessage.success('核验通过')
    fetchOrder()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('核验失败')
    }
  } finally {
    actionLoading.value = false
  }
}

// 取消订单
const handleCancel = async () => {
  try {
    await ElMessageBox.confirm('确定要取消该订单吗？', '取消确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    actionLoading.value = true
    await request.post({
      url: `/api/admin/orders/${order.value.id}/cancel`
    })
    ElMessage.success('订单已取消')
    fetchOrder()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('取消失败')
    }
  } finally {
    actionLoading.value = false
  }
}

// 删除订单
const handleDelete = async () => {
  try {
    await ElMessageBox.confirm('确定要删除该订单吗？此操作不可恢复。', '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    actionLoading.value = true
    await request.del({
      url: `/api/admin/orders/${order.value.id}`
    })
    ElMessage.success('订单已删除')
    router.back()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  } finally {
    actionLoading.value = false
  }
}

// 查看项目详情
const handleViewItem = (row: any) => {
  // TODO: 跳转到项目详情
  ElMessage.info('查看项目详情功能开发中')
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
