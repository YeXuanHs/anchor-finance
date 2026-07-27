<template>
  <div class="cart-page">
    <!-- Header -->
    <header class="header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <div class="logo-icon">
            <n-icon size="24" color="#fff"><AnchorOutline /></n-icon>
          </div>
          <span class="logo-text">锚点财务</span>
        </router-link>
        <nav class="nav-links">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/products" class="nav-link">产品</router-link>
          <a href="#" class="nav-link">公告</a>
          <a href="#" class="nav-link">帮助</a>
        </nav>
        <div class="header-actions">
          <n-button text @click="$router.push('/login')">登录</n-button>
          <n-button type="primary" round size="small" @click="$router.push('/register')">免费注册</n-button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <div class="main-content">
      <div class="content-inner">
        <!-- Breadcrumb -->
        <n-breadcrumb class="breadcrumb">
          <n-breadcrumb-item @click="$router.push('/')">
            <template #icon>
              <n-icon :component="HomeOutline" />
            </template>
            首页
          </n-breadcrumb-item>
          <n-breadcrumb-item>购物车</n-breadcrumb-item>
        </n-breadcrumb>

        <!-- Page Title -->
        <div class="page-title-bar">
          <h1 class="page-title">
            <n-icon size="24" color="#1890ff"><CartOutline /></n-icon>
            购物车
            <span class="item-count" v-if="cartItems.length">({{ cartItems.length }})</span>
          </h1>
        </div>

        <!-- Cart Content -->
        <div v-if="cartItems.length > 0" class="cart-layout">
          <!-- Cart Items -->
          <div class="cart-items-panel">
            <!-- Table Header -->
            <div class="cart-table-header">
              <span class="col-product">商品信息</span>
              <span class="col-config">配置摘要</span>
              <span class="col-cycle">周期</span>
              <span class="col-price">单价</span>
              <span class="col-qty">数量</span>
              <span class="col-total">小计</span>
              <span class="col-action">操作</span>
            </div>

            <!-- Cart Items -->
            <div
              v-for="item in cartItems"
              :key="item.id"
              class="cart-item"
            >
              <div class="col-product">
                <div class="item-icon-wrap">
                  <n-icon size="22" color="#1890ff"><ServerOutline /></n-icon>
                </div>
                <div class="item-info">
                  <h3 class="item-name">{{ item.name }}</h3>
                  <p class="item-id">ID: {{ item.id }}</p>
                </div>
              </div>
              <div class="col-config">
                <div class="config-tags">
                  <n-tag size="small" :bordered="false" type="info">{{ item.cpu }}</n-tag>
                  <n-tag size="small" :bordered="false" type="info">{{ item.memory }}</n-tag>
                  <n-tag size="small" :bordered="false" type="info">{{ item.disk }}</n-tag>
                  <n-tag size="small" :bordered="false" type="info">{{ item.bandwidth }}</n-tag>
                </div>
                <p class="config-os">{{ item.os }}</p>
              </div>
              <div class="col-cycle">
                <n-select
                  v-model:value="item.cycle"
                  :options="cycleOptions"
                  size="small"
                  style="width: 90px;"
                  @update:value="updateItemPrice(item)"
                />
              </div>
              <div class="col-price">
                <span class="price-text">¥{{ item.unitPrice.toFixed(2) }}</span>
              </div>
              <div class="col-qty">
                <n-input-number
                  v-model:value="item.quantity"
                  :min="1"
                  :max="20"
                  size="small"
                  style="width: 80px;"
                />
              </div>
              <div class="col-total">
                <span class="total-text">¥{{ (item.unitPrice * item.quantity).toFixed(2) }}</span>
              </div>
              <div class="col-action">
                <n-button text type="error" size="small" @click="removeItem(item.id)">
                  <template #icon><n-icon><TrashOutline /></n-icon></template>
                </n-button>
              </div>
            </div>
          </div>

          <!-- Summary Panel -->
          <div class="summary-panel">
            <div class="summary-card">
              <h3 class="summary-title">订单摘要</h3>

              <div class="summary-row">
                <span class="summary-label">商品小计</span>
                <span class="summary-value">¥{{ subtotal.toFixed(2) }}</span>
              </div>

              <!-- Applied Coupon -->
              <div v-if="appliedCoupon" class="applied-coupon">
                <div class="coupon-info">
                  <n-icon size="16" color="#52c41a"><CheckmarkCircleOutline /></n-icon>
                  <span>优惠码: {{ appliedCoupon }}</span>
                </div>
                <n-button text type="error" size="small" @click="removeCoupon">移除</n-button>
              </div>

              <!-- Coupon Input -->
              <div v-else class="coupon-row">
                <n-input
                  v-model:value="couponCode"
                  placeholder="请输入优惠码"
                  size="small"
                />
                <n-button type="primary" size="small" ghost @click="applyCoupon">应用</n-button>
              </div>
              <p v-if="couponMessage" class="coupon-msg" :class="couponSuccess ? 'success' : 'error'">
                {{ couponMessage }}
              </p>

              <div class="summary-row discount-row" v-if="discount > 0">
                <span class="summary-label">优惠折扣</span>
                <span class="summary-value discount-value">-¥{{ discount.toFixed(2) }}</span>
              </div>

              <n-divider />

              <div class="summary-row total-row">
                <span class="summary-label">总计</span>
                <div class="summary-total">
                  <span class="total-currency">¥</span>
                  <span class="total-amount">{{ finalTotal.toFixed(2) }}</span>
                </div>
              </div>

              <n-button
                type="primary"
                size="large"
                block
                round
                class="checkout-btn"
                @click="handleCheckout"
              >
                去结算
              </n-button>

              <div class="summary-badges">
                <div class="badge-item">
                  <n-icon size="14" color="#52c41a"><CheckmarkCircleOutline /></n-icon>
                  <span>7天无理由退款</span>
                </div>
                <div class="badge-item">
                  <n-icon size="14" color="#52c41a"><CheckmarkCircleOutline /></n-icon>
                  <span>安全支付</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty Cart -->
        <div v-else class="empty-cart">
          <div class="empty-illustration">
            <svg viewBox="0 0 200 160" width="200" height="160" fill="none">
              <ellipse cx="100" cy="140" rx="80" ry="12" fill="#f0f1f5"/>
              <path d="M60 50h80l20 70H40z" stroke="#d9d9d9" stroke-width="2" fill="#fafafa" stroke-linejoin="round"/>
              <path d="M40 120h120" stroke="#d9d9d9" stroke-width="2"/>
              <circle cx="65" cy="130" r="8" fill="#d9d9d9"/>
              <circle cx="135" cy="130" r="8" fill="#d9d9d9"/>
              <path d="M80 80h40M88 90h24" stroke="#bfbfbf" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </div>
          <h2 class="empty-title">购物车是空的</h2>
          <p class="empty-desc">快去挑选心仪的产品吧</p>
          <n-button type="primary" size="large" round @click="$router.push('/products')">
            去选购产品
          </n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import {
  AnchorOutline,
  ServerOutline,
  CartOutline,
  TrashOutline,
  CheckmarkCircleOutline,
  HomeOutline
} from '@vicons/ionicons5'

const router = useRouter()
const message = useMessage()

interface CartItem {
  id: number
  name: string
  cpu: string
  memory: string
  disk: string
  bandwidth: string
  os: string
  cycle: string
  quantity: number
  unitPrice: number
  basePrice: number
}

const cycleOptions = [
  { label: '月付', value: '月付' },
  { label: '季付', value: '季付' },
  { label: '年付', value: '年付' }
]

const cycleMultiplier: Record<string, number> = {
  '月付': 1,
  '季付': 2.8,
  '年付': 10
}

const cartItems = ref<CartItem[]>([
  {
    id: 10001,
    name: '云服务器 ECS',
    cpu: '2核',
    memory: '4G',
    disk: '100G SSD',
    bandwidth: '5Mbps',
    os: 'CentOS 7',
    cycle: '月付',
    quantity: 1,
    unitPrice: 134,
    basePrice: 134
  },
  {
    id: 10002,
    name: '云服务器 ECS',
    cpu: '4核',
    memory: '8G',
    disk: '200G SSD',
    bandwidth: '10Mbps',
    os: 'Ubuntu 22.04',
    cycle: '年付',
    quantity: 2,
    unitPrice: 2890,
    basePrice: 134
  }
])

const couponCode = ref('')
const appliedCoupon = ref('')
const couponMessage = ref('')
const couponSuccess = ref(false)
const discountRate = ref(0)

const subtotal = computed(() =>
  cartItems.value.reduce((sum, item) => sum + item.unitPrice * item.quantity, 0)
)

const discount = computed(() => subtotal.value * discountRate.value)

const finalTotal = computed(() => subtotal.value - discount.value)

function updateItemPrice(item: CartItem) {
  item.unitPrice = Math.round(item.basePrice * cycleMultiplier[item.cycle])
}

function removeItem(id: number) {
  cartItems.value = cartItems.value.filter(item => item.id !== id)
  message.success('已移除')
}

function applyCoupon() {
  if (!couponCode.value) {
    couponMessage.value = '请输入优惠码'
    couponSuccess.value = false
    return
  }
  if (couponCode.value.toUpperCase() === 'ANCHOR2024') {
    couponMessage.value = '优惠码有效，享受9折优惠'
    couponSuccess.value = true
    discountRate.value = 0.1
    appliedCoupon.value = couponCode.value.toUpperCase()
    couponCode.value = ''
  } else {
    couponMessage.value = '优惠码无效'
    couponSuccess.value = false
    discountRate.value = 0
  }
}

function removeCoupon() {
  appliedCoupon.value = ''
  discountRate.value = 0
  couponMessage.value = ''
  message.success('优惠码已移除')
}

function handleCheckout() {
  const token = localStorage.getItem('token')
  if (!token) {
    message.warning('请先登录')
    router.push('/login')
    return
  }
  message.success('正在跳转到结算页面...')
  router.push('/user/orders')
}
</script>

<style scoped>
.cart-page {
  min-height: 100vh;
  background: #f7f8fa;
}

/* Header */
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: #fff;
  box-shadow: 0 1px 8px rgba(0, 0, 0, 0.06);
  z-index: 100;
}

.header-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
}

.logo-icon {
  width: 34px;
  height: 34px;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  color: #1d2129;
}

.nav-links {
  display: flex;
  gap: 32px;
}

.nav-link {
  color: #4e5969;
  text-decoration: none;
  font-size: 15px;
  font-weight: 500;
  transition: color 0.2s;
}

.nav-link:hover,
.nav-link.active {
  color: #1890ff;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* Main Content */
.main-content {
  padding-top: 88px;
  padding-bottom: 40px;
}

.content-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

/* Breadcrumb */
.breadcrumb {
  margin-bottom: 20px;
}

/* Page Title */
.page-title-bar {
  margin-bottom: 24px;
}

.page-title {
  font-size: 22px;
  font-weight: 700;
  color: #1d2129;
  display: flex;
  align-items: center;
  gap: 10px;
}

.item-count {
  font-size: 16px;
  font-weight: 400;
  color: #86909c;
}

/* Cart Layout */
.cart-layout {
  display: flex;
  gap: 24px;
}

/* Cart Items Panel */
.cart-items-panel {
  flex: 1;
  min-width: 0;
}

.cart-table-header {
  display: grid;
  grid-template-columns: 2fr 2.5fr 0.8fr 0.8fr 0.8fr 0.8fr 0.5fr;
  gap: 12px;
  padding: 14px 20px;
  background: #fff;
  border-radius: 12px 12px 0 0;
  font-size: 13px;
  font-weight: 600;
  color: #86909c;
  border-bottom: 1px solid #f0f1f5;
}

.cart-item {
  display: grid;
  grid-template-columns: 2fr 2.5fr 0.8fr 0.8fr 0.8fr 0.8fr 0.5fr;
  gap: 12px;
  align-items: center;
  padding: 20px;
  background: #fff;
  border-bottom: 1px solid #f7f8fa;
  transition: background 0.2s;
}

.cart-item:hover {
  background: #fafbfc;
}

.cart-item:last-child {
  border-radius: 0 0 12px 12px;
  border-bottom: none;
}

/* Product Column */
.col-product {
  display: flex;
  align-items: center;
  gap: 12px;
}

.item-icon-wrap {
  width: 42px;
  height: 42px;
  background: linear-gradient(135deg, #e6f7ff, #bae7ff);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.item-name {
  font-size: 15px;
  font-weight: 600;
  color: #1d2129;
}

.item-id {
  font-size: 12px;
  color: #c9cdd4;
  margin-top: 2px;
}

/* Config Column */
.config-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 4px;
}

.config-os {
  font-size: 12px;
  color: #86909c;
}

/* Price */
.price-text,
.total-text {
  font-size: 14px;
  font-weight: 600;
  color: #1d2129;
}

.total-text {
  color: #1890ff;
}

/* Summary Panel */
.summary-panel {
  width: 340px;
  flex-shrink: 0;
}

.summary-card {
  background: #fff;
  border-radius: 12px;
  padding: 28px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
  position: sticky;
  top: 88px;
}

.summary-title {
  font-size: 17px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 24px;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.summary-label {
  font-size: 14px;
  color: #4e5969;
}

.summary-value {
  font-size: 15px;
  font-weight: 600;
  color: #1d2129;
}

.discount-value {
  color: #52c41a;
}

.applied-coupon {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: #f6ffed;
  border: 1px solid #b7eb8f;
  border-radius: 8px;
  margin-bottom: 12px;
}

.coupon-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #52c41a;
}

.coupon-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.coupon-msg {
  font-size: 12px;
  margin-bottom: 12px;
}

.coupon-msg.success {
  color: #52c41a;
}

.coupon-msg.error {
  color: #ff4d4f;
}

.total-row {
  margin-bottom: 20px;
}

.summary-total {
  display: flex;
  align-items: baseline;
}

.total-currency {
  font-size: 18px;
  font-weight: 600;
  color: #1890ff;
}

.total-amount {
  font-size: 32px;
  font-weight: 700;
  color: #1890ff;
  line-height: 1;
}

.checkout-btn {
  height: 48px;
  font-size: 16px;
  font-weight: 600;
}

.summary-badges {
  margin-top: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.badge-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #86909c;
}

/* Empty Cart */
.empty-cart {
  text-align: center;
  padding: 80px 0;
}

.empty-illustration {
  margin-bottom: 24px;
}

.empty-title {
  font-size: 20px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 8px;
}

.empty-desc {
  font-size: 14px;
  color: #86909c;
  margin-bottom: 32px;
}

/* Responsive */
@media (max-width: 1024px) {
  .cart-layout {
    flex-direction: column-reverse;
  }
  .summary-panel {
    width: 100%;
  }
  .summary-card {
    position: static;
  }
}

@media (max-width: 768px) {
  .nav-links {
    display: none;
  }
  .cart-table-header {
    display: none;
  }
  .cart-item {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
  }
  .col-product {
    width: 100%;
  }
  .col-config {
    width: 100%;
  }
  .col-cycle,
  .col-qty,
  .col-price,
  .col-total,
  .col-action {
    flex: 1;
    min-width: 0;
  }
}
</style>
