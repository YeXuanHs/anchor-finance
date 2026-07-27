<template>
  <div class="orders-page">
    <div class="page-header">
      <h1 class="page-title">我的订单</h1>
      <n-input
        v-model:value="searchKey"
        placeholder="搜索订单号或产品名称"
        clearable
        class="search-input"
      >
        <template #prefix>
          <n-icon :component="SearchOutline" color="#bfbfbf" />
        </template>
      </n-input>
    </div>

    <n-tabs v-model:value="activeTab" type="line" animated class="filter-tabs">
      <n-tab v-for="tab in statusTabs" :key="tab.value" :name="tab.value">
        {{ tab.label }}
        <n-badge
          v-if="tab.count > 0"
          :value="tab.count"
          :max="99"
          class="tab-badge"
        />
      </n-tab>
    </n-tabs>

    <div class="order-list">
      <div v-for="order in filteredOrders" :key="order.id" class="order-card">
        <div class="order-card-header">
          <div class="header-left">
            <span class="order-number">
              <span class="label">订单号：</span>
              <span class="value">{{ order.id }}</span>
            </span>
            <n-divider vertical />
            <span class="order-time">
              <n-icon :component="TimeOutline" size="14" />
              {{ order.createdAt }}
            </span>
          </div>
          <n-tag :type="getStatusType(order.status)" size="small" round>
            {{ order.statusText }}
          </n-tag>
        </div>

        <div class="order-card-body">
          <div class="order-product">
            <div class="product-icon">
              <n-icon :size="24" :component="CubeOutline" color="#1890ff" />
            </div>
            <div class="product-info">
              <span class="product-name">{{ order.product }}</span>
              <span class="product-spec">{{ order.spec }}</span>
              <span class="product-cycle">{{ order.cycle }}</span>
            </div>
          </div>

          <div class="order-amount">
            <span class="amount-label">订单金额</span>
            <span class="amount-value">¥{{ order.amount }}</span>
          </div>
        </div>

        <div class="order-card-footer">
          <div class="order-actions">
            <n-button
              size="small"
              quaternary
              @click="handleDetail(order)"
            >
              查看详情
            </n-button>
            <n-button
              v-if="order.status === 'pending'"
              type="primary"
              size="small"
              @click="handlePay(order)"
            >
              去支付
            </n-button>
            <n-button
              v-if="order.status === 'active'"
              type="primary"
              size="small"
              @click="handleRenew(order)"
            >
              续费
            </n-button>
          </div>
        </div>
      </div>

      <div v-if="filteredOrders.length === 0" class="empty-state">
        <n-icon :size="64" :component="ReceiptOutline" color="#d9d9d9" />
        <p>暂无订单</p>
        <n-button type="primary" @click="$router.push('/products')">购买产品</n-button>
      </div>
    </div>

    <div v-if="totalPages > 1" class="pagination-wrapper">
      <n-pagination
        v-model:page="currentPage"
        :page-count="totalPages"
        :page-slot="7"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import {
  SearchOutline,
  CubeOutline,
  ReceiptOutline,
  TimeOutline
} from '@vicons/ionicons5'

const message = useMessage()
const searchKey = ref('')
const activeTab = ref('all')
const currentPage = ref(1)

const statusTabs = computed(() => [
  { label: '全部', value: 'all', count: orders.value.length },
  { label: '待支付', value: 'pending', count: orders.value.filter(o => o.status === 'pending').length },
  { label: '已开通', value: 'active', count: orders.value.filter(o => o.status === 'active').length },
  { label: '已完成', value: 'completed', count: orders.value.filter(o => o.status === 'completed').length },
  { label: '已取消', value: 'cancelled', count: orders.value.filter(o => o.status === 'cancelled').length }
])

interface Order {
  id: string
  product: string
  spec: string
  cycle: string
  amount: string
  status: string
  statusText: string
  createdAt: string
}

const orders = ref<Order[]>([
  { id: 'ORD20260725001', product: '香港云服务器', spec: '2核4G / 50G SSD / 5Mbps', cycle: '月度订阅', amount: '49.00', status: 'active', statusText: '已开通', createdAt: '2026-07-25 10:30:00' },
  { id: 'ORD20260724002', product: '美国独立服务器', spec: 'E5-2680v4 / 64G / 1T SSD', cycle: '季度订阅', amount: '2,397.00', status: 'pending', statusText: '待支付', createdAt: '2026-07-24 14:20:00' },
  { id: 'ORD20260720003', product: 'OV SSL证书', spec: '单域名 / 企业验证', cycle: '年度订阅', amount: '199.00', status: 'active', statusText: '已开通', createdAt: '2026-07-20 09:15:00' },
  { id: 'ORD20260715004', product: '香港 VPS', spec: '1核2G / 30G NVMe / 1Gbps', cycle: '月度订阅', amount: '19.00', status: 'completed', statusText: '已完成', createdAt: '2026-07-15 16:45:00' },
  { id: 'ORD20260710005', product: '域名注册', spec: '.com / 首年', cycle: '一次性', amount: '9.00', status: 'completed', statusText: '已完成', createdAt: '2026-07-10 11:00:00' },
  { id: 'ORD20260705006', product: '新加坡 VPS', spec: '2核4G / 60G NVMe', cycle: '月度订阅', amount: '35.00', status: 'cancelled', statusText: '已取消', createdAt: '2026-07-05 09:00:00' }
])

const filteredOrders = computed(() => {
  let result = orders.value
  if (activeTab.value !== 'all') {
    result = result.filter(o => o.status === activeTab.value)
  }
  if (searchKey.value) {
    const key = searchKey.value.toLowerCase()
    result = result.filter(o =>
      o.id.toLowerCase().includes(key) ||
      o.product.toLowerCase().includes(key)
    )
  }
  return result
})

const totalPages = computed(() => Math.ceil(filteredOrders.value.length / 10))

function getStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'error'> = {
    active: 'success',
    pending: 'warning',
    completed: 'info',
    cancelled: 'error'
  }
  return map[status] || 'info'
}

function handlePay(order: Order) {
  message.info(`正在跳转支付页面：${order.id}`)
}

function handleRenew(order: Order) {
  message.info(`正在处理续费：${order.id}`)
}

function handleDetail(order: Order) {
  message.info(`查看订单详情：${order.id}`)
}
</script>

<style scoped>
.orders-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 22px;
  font-weight: 700;
  color: #262626;
  margin: 0;
}

.search-input {
  width: 280px;
}

.filter-tabs {
  background: #fff;
  border-radius: 12px;
  padding: 4px 16px;
  border: 1px solid #f0f0f0;
}

.tab-badge {
  margin-left: 6px;
}

.order-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.order-card {
  background: #fff;
  border-radius: 12px;
  border: 1px solid #f0f0f0;
  overflow: hidden;
  transition: all 0.3s ease;
}

.order-card:hover {
  border-color: #1890ff;
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.1);
}

.order-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: #fafafa;
  border-bottom: 1px solid #f0f0f0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.order-number .label {
  font-size: 13px;
  color: #8c8c8c;
}

.order-number .value {
  font-size: 13px;
  color: #262626;
  font-weight: 500;
  font-family: 'Monaco', 'Menlo', monospace;
}

.order-time {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #8c8c8c;
}

.order-card-body {
  display: flex;
  align-items: center;
  padding: 20px;
  gap: 40px;
}

.order-product {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.product-icon {
  width: 48px;
  height: 48px;
  background: #f0f5ff;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.product-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.product-name {
  font-size: 15px;
  font-weight: 600;
  color: #262626;
}

.product-spec {
  font-size: 12px;
  color: #595959;
  font-family: 'Monaco', 'Menlo', monospace;
}

.product-cycle {
  font-size: 12px;
  color: #8c8c8c;
}

.order-amount {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 120px;
  align-items: flex-end;
}

.amount-label {
  font-size: 12px;
  color: #8c8c8c;
}

.amount-value {
  font-size: 20px;
  font-weight: 700;
  color: #ff4d4f;
}

.order-card-footer {
  display: flex;
  justify-content: flex-end;
  padding: 12px 20px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
}

.order-actions {
  display: flex;
  gap: 8px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 80px 0;
  background: #fff;
  border-radius: 12px;
  border: 1px solid #f0f0f0;
}

.empty-state p {
  margin: 0;
  color: #8c8c8c;
  font-size: 14px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }

  .search-input {
    width: 100%;
  }

  .order-card-body {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }

  .order-amount {
    min-width: auto;
    align-items: flex-start;
  }

  .header-left {
    flex-wrap: wrap;
    gap: 8px;
  }
}
</style>
