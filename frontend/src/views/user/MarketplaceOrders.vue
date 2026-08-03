<template>
  <div class="marketplace-orders">
    <div class="page-header">
      <el-button @click="$router.back()" text>
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>我的订单</h2>
    </div>

    <el-tabs v-model="activeTab" @tab-change="fetchOrders">
      <el-tab-pane label="我买到的" name="buyer" />
      <el-tab-pane label="我卖出的" name="seller" />
    </el-tabs>

    <div class="orders-list" v-loading="loading">
      <div v-if="orders.length === 0 && !loading" class="empty-tip">
        <el-empty description="暂无订单" />
      </div>

      <div v-for="order in orders" :key="order.id" class="order-card">
        <div class="order-header">
          <div class="order-info">
            <span class="order-no">订单号: {{ order.order_no }}</span>
            <el-tag :type="getStatusType(order.status)" size="small">
              {{ getStatusText(order.status) }}
            </el-tag>
          </div>
          <span class="order-time">{{ formatDate(order.created_at) }}</span>
        </div>

        <div class="order-body">
          <div class="product-info">
            <h3>{{ order.listing?.product_name || '主机' }}</h3>
            <div class="product-meta">
              <span v-if="activeTab === 'buyer'">
                卖家: {{ order.seller?.username || '-' }}
              </span>
              <span v-else>
                买家: {{ order.buyer?.username || '-' }}
              </span>
            </div>
          </div>

          <div class="price-info">
            <div class="price-detail">
              <div v-if="order.amount > 0" class="price-row">
                <span>商品价格</span>
                <span>¥{{ order.amount.toFixed(2) }}</span>
              </div>
              <div class="price-row">
                <span>手续费</span>
                <span>¥{{ order.fee.toFixed(2) }}</span>
              </div>
              <div class="price-row total">
                <span>总计</span>
                <span>¥{{ order.total_amount.toFixed(2) }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="order-footer">
          <div class="payment-type">
            <el-tag v-if="order.payment_method === 'full'" type="success" size="small">
              全额购买
            </el-tag>
            <el-tag v-else type="warning" size="small">
              仅付手续费
            </el-tag>
          </div>

          <div class="order-actions">
            <!-- 买家操作 -->
            <template v-if="activeTab === 'buyer'">
              <el-button
                v-if="order.status === 0"
                type="primary"
                size="small"
                @click="payOrder(order)"
              >
                立即支付
              </el-button>
              <el-button
                v-if="order.status === 0"
                size="small"
                @click="cancelOrder(order)"
              >
                取消订单
              </el-button>
            </template>

            <!-- 卖家操作 -->
            <template v-if="activeTab === 'seller'">
              <el-button
                v-if="order.status === 2"
                type="success"
                size="small"
                @click="completeOrder(order)"
              >
                确认交易完成
              </el-button>
              <el-button
                v-if="order.payment_method === 'fee_only' && order.status >= 1"
                size="small"
                @click="viewBuyerContact(order)"
              >
                查看买家信息
              </el-button>
            </template>

            <!-- 聊天 -->
            <el-button
              size="small"
              @click="openChat(order)"
            >
              <el-icon><ChatDotRound /></el-icon>
              联系对方
            </el-button>
          </div>
        </div>

        <!-- 转移状态 -->
        <div v-if="order.transfer_status > 0" class="transfer-status">
          <el-divider />
          <div class="transfer-info">
            <span class="transfer-label">转移状态:</span>
            <el-tag :type="getTransferStatusType(order.transfer_status)" size="small">
              {{ getTransferStatusText(order.transfer_status) }}
            </el-tag>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination-wrapper" v-if="total > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        @current-change="fetchOrders"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, ChatDotRound } from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()

const activeTab = ref<'buyer' | 'seller'>('buyer')
const loading = ref(false)
const orders = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

onMounted(() => {
  fetchOrders()
})

async function fetchOrders() {
  loading.value = true
  try {
    const endpoint = activeTab.value === 'buyer'
      ? '/v1/marketplace/orders/buyer'
      : '/v1/marketplace/orders/seller'

    const res = await request.get(endpoint, {
      params: {
        page: currentPage.value,
        page_size: pageSize.value
      }
    })
    orders.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function payOrder(order: any) {
  try {
    await ElMessageBox.confirm('确定要支付此订单吗？', '确认支付', { type: 'warning' })
    await request.post(`/v1/marketplace/orders/${order.id}/pay`)
    ElMessage.success('支付成功')
    fetchOrders()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '支付失败')
    }
  }
}

async function cancelOrder(order: any) {
  try {
    await ElMessageBox.confirm('确定要取消此订单吗？款项将退回余额。', '取消订单', { type: 'warning' })
    await request.post(`/v1/marketplace/orders/${order.id}/cancel`)
    ElMessage.success('订单已取消')
    fetchOrders()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '取消失败')
    }
  }
}

async function completeOrder(order: any) {
  try {
    await ElMessageBox.confirm('确认交易已经完成？确认后服务器将自动转移给买家。', '确认完成', { type: 'success' })
    await request.post(`/v1/marketplace/orders/${order.id}/complete`)
    ElMessage.success('交易已完成')
    fetchOrders()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '操作失败')
    }
  }
}

function viewBuyerContact(order: any) {
  if (order.buyer) {
    ElMessageBox.alert(
      `<div>
        <p><b>用户名:</b> ${order.buyer.username}</p>
        <p><b>邮箱:</b> ${order.buyer.email || '-'}</p>
        <p><b>手机:</b> ${order.buyer.phone || '-'}</p>
        <p><b>QQ:</b> ${order.buyer.qq || '-'}</p>
      </div>`,
      '买家信息',
      { dangerouslyUseHTMLString: true }
    )
  }
}

function openChat(order: any) {
  const otherId = activeTab.value === 'buyer' ? order.seller_id : order.buyer_id
  router.push(`/user/marketplace/chat/${order.listing_id}/${otherId}`)
}

function getStatusType(status: number) {
  const map: Record<number, string> = {
    0: 'warning',
    1: 'primary',
    2: 'success',
    3: 'success',
    4: 'info',
    5: 'danger'
  }
  return map[status] || 'info'
}

function getStatusText(status: number) {
  const map: Record<number, string> = {
    0: '待支付',
    1: '已支付',
    2: '已转移',
    3: '已完成',
    4: '已取消',
    5: '已退款'
  }
  return map[status] || '未知'
}

function getTransferStatusType(status: number) {
  const map: Record<number, string> = {
    1: 'warning',
    2: 'success',
    3: 'danger'
  }
  return map[status] || 'info'
}

function getTransferStatusText(status: number) {
  const map: Record<number, string> = {
    1: '转移中',
    2: '转移成功',
    3: '转移失败'
  }
  return map[status] || '未知'
}

function formatDate(date: string): string {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}
</script>

<style scoped lang="scss">
.marketplace-orders {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;

  h2 {
    margin: 0;
  }
}

.orders-list {
  min-height: 200px;
}

.order-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.order-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.order-no {
  font-size: 13px;
  color: #666;
}

.order-time {
  font-size: 12px;
  color: #999;
}

.order-body {
  display: flex;
  justify-content: space-between;
  padding: 16px 0;
}

.product-info {
  h3 {
    margin: 0 0 8px 0;
    font-size: 16px;
  }

  .product-meta {
    font-size: 13px;
    color: #666;
  }
}

.price-detail {
  text-align: right;
}

.price-row {
  display: flex;
  justify-content: space-between;
  gap: 24px;
  font-size: 13px;
  margin-bottom: 4px;

  &.total {
    font-weight: bold;
    font-size: 15px;
    color: #ff4757;
    border-top: 1px solid #eee;
    padding-top: 8px;
    margin-top: 8px;
  }
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.order-actions {
  display: flex;
  gap: 8px;
}

.transfer-status {
  margin-top: 12px;
}

.transfer-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.transfer-label {
  font-size: 13px;
  color: #666;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

.empty-tip {
  padding: 60px 0;
}
</style>
