<template>
  <div class="order-detail-page page-container">
    <div class="page-header">
      <el-button @click="$router.back()">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h3>订单详情</h3>
      <div class="header-actions">
        <el-button type="success" v-if="order.status === 'pending'" @click="confirmOrder">确认订单</el-button>
        <el-button type="danger" v-if="order.status === 'pending'" @click="cancelOrder">取消订单</el-button>
      </div>
    </div>

    <div class="detail-grid">
      <div class="art-card">
        <h4>基本信息</h4>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="订单号">{{ order.order_no }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(order.status)">{{ getStatusText(order.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="用户">{{ order.username }}</el-descriptions-item>
          <el-descriptions-item label="产品">{{ order.product_name }}</el-descriptions-item>
          <el-descriptions-item label="金额">
            <span class="amount">¥{{ order.amount?.toFixed(2) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="支付方式">{{ order.payment_method }}</el-descriptions-item>
          <el-descriptions-item label="下单时间">{{ order.created_at }}</el-descriptions-item>
          <el-descriptions-item label="支付时间">{{ order.paid_at || '-' }}</el-descriptions-item>
        </el-descriptions>
      </div>

      <div class="art-card">
        <h4>产品配置</h4>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="产品名称">{{ order.product_name }}</el-descriptions-item>
          <el-descriptions-item label="产品周期">{{ order.billing_cycle }}</el-descriptions-item>
          <el-descriptions-item label="配置信息">
            <pre v-if="order.config">{{ order.config }}</pre>
            <span v-else>-</span>
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <div class="art-card">
        <h4>账单信息</h4>
        <el-table :data="order.invoices || []" style="width: 100%">
          <el-table-column prop="invoice_no" label="账单号" />
          <el-table-column prop="amount" label="金额">
            <template #default="{ row }">
              <span class="amount">¥{{ row.amount?.toFixed(2) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态">
            <template #default="{ row }">
              <el-tag :type="row.status === 'paid' ? 'success' : 'warning'" size="small">
                {{ row.status === 'paid' ? '已支付' : '未支付' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" />
        </el-table>
      </div>

      <div class="art-card">
        <h4>操作日志</h4>
        <el-timeline>
          <el-timeline-item v-for="log in order.logs" :key="log.id" :timestamp="log.created_at" placement="top">
            <p>{{ log.content }}</p>
            <p class="log-operator">操作人：{{ log.operator }}</p>
          </el-timeline-item>
        </el-timeline>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const route = useRoute()
const order = ref<any>({})

const getStatusType = (status: string) => {
  const map: Record<string, string> = { pending: 'warning', paid: 'primary', completed: 'success', cancelled: 'info', refunded: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = { pending: '待支付', paid: '已支付', completed: '已完成', cancelled: '已取消', refunded: '已退款' }
  return map[status] || status
}

const fetchOrder = async () => {
  const id = route.params.id
  try {
    const { data } = await request.get(`/admin/api/v1/orders/${id}`)
    order.value = data.data || data
  } catch {}
}

const confirmOrder = async () => {}
const cancelOrder = async () => {}

onMounted(() => { fetchOrder() })
</script>

<style scoped lang="scss">
.order-detail-page {
  .page-header {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;
    h3 { margin: 0; font-size: 18px; font-weight: 600; flex: 1; }
  }
  .detail-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
    .art-card {
      h4 { margin: 0 0 16px; font-size: 15px; font-weight: 600; }
    }
  }
  .amount { color: var(--danger-color); font-weight: 600; }
  .log-operator { color: var(--text-secondary); font-size: 12px; margin-top: 4px; }
  pre { margin: 0; white-space: pre-wrap; word-break: break-all; }
}
</style>
