<template>
  <div class="order-detail" v-loading="loading">
    <div class="page-header">
      <el-page-header @back="$router.push('/user/orders')">
        <template #content>
          <span class="page-title">订单详情</span>
        </template>
      </el-page-header>
    </div>
    
    <div class="content-wrapper" v-if="order">
      <!-- 订单状态 -->
      <div class="status-card">
        <div class="status-info">
          <el-tag :type="getStatusType(order.status)" size="large">
            {{ getStatusText(order.status) }}
          </el-tag>
          <span class="order-no">订单号：{{ order.order_no }}</span>
        </div>
        <div class="status-actions">
          <el-button v-if="order.status === 'pending'" type="primary" @click="payOrder">立即支付</el-button>
          <el-button v-if="order.status === 'pending'" @click="cancelOrder">取消订单</el-button>
        </div>
      </div>
      
      <!-- 订单信息 -->
      <el-card class="info-card">
        <template #header>
          <span>订单信息</span>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="产品名称">{{ order.product_name }}</el-descriptions-item>
          <el-descriptions-item label="计费周期">{{ order.cycle }}</el-descriptions-item>
          <el-descriptions-item label="数量">{{ order.quantity }}</el-descriptions-item>
          <el-descriptions-item label="订单金额">
            <span class="price">¥{{ order.amount?.toFixed(2) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="下单时间">{{ order.created_at }}</el-descriptions-item>
          <el-descriptions-item label="支付方式">{{ order.payment_method || '-' }}</el-descriptions-item>
          <el-descriptions-item label="备注" :span="2">{{ order.remark || '无' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>
      
      <!-- 产品配置 -->
      <el-card class="info-card" v-if="order.config">
        <template #header>
          <span>产品配置</span>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item v-for="(value, key) in order.config" :key="key" :label="key">
            {{ value }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const order = ref<any>(null)

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    paid: 'success',
    cancelled: 'info',
    completed: 'success'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    cancelled: '已取消',
    completed: '已完成'
  }
  return map[status] || status
}

const fetchOrder = async () => {
  loading.value = true
  try {
    const { data } = await request.get(`/api/v1/orders/${route.params.id}`)
    if (data?.data) {
      order.value = data.data
    }
  } catch (error) {
    ElMessage.error('获取订单信息失败')
  } finally {
    loading.value = false
  }
}

const payOrder = () => {
  router.push(`/checkout?order=${order.value.id}`)
}

const cancelOrder = async () => {
  try {
    await ElMessageBox.confirm('确定要取消这个订单吗？', '提示', { type: 'warning' })
    await request.post(`/api/v1/orders/${order.value.id}/cancel`)
    ElMessage.success('订单已取消')
    fetchOrder()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('取消失败')
    }
  }
}

onMounted(() => {
  fetchOrder()
})
</script>

<style scoped lang="scss">
.order-detail {
  .page-header {
    margin-bottom: 24px;
  }
  
  .page-title {
    font-size: 18px;
    font-weight: 600;
  }
  
  .status-card {
    background: #fff;
    border-radius: 12px;
    padding: 24px;
    margin-bottom: 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    .status-info {
      display: flex;
      align-items: center;
      gap: 16px;
      
      .order-no {
        color: #909399;
        font-size: 14px;
      }
    }
  }
  
  .info-card {
    margin-bottom: 20px;
    
    .price {
      color: #f56c6c;
      font-weight: 600;
      font-size: 16px;
    }
  }
}
</style>
