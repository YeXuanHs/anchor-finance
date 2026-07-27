<template>
  <div class="cart-page">
    <!-- Header -->
    <header class="page-header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <div class="logo-icon">
            <el-icon :size="22" color="#fff"><Monitor /></el-icon>
          </div>
          <span class="logo-text">智简魔方</span>
        </router-link>
        <nav class="nav-links">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/products" class="nav-link">产品</router-link>
          <router-link to="/" class="nav-link">公告</router-link>
          <router-link to="/" class="nav-link">帮助</router-link>
        </nav>
        <div class="header-actions">
          <el-button text @click="$router.push('/login')">登录</el-button>
          <el-button class="btn-gradient" round size="small" @click="$router.push('/register')">免费注册</el-button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <div class="main-content">
      <div class="content-inner">
        <!-- Breadcrumb -->
        <el-breadcrumb separator="/" class="breadcrumb">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item>购物车</el-breadcrumb-item>
        </el-breadcrumb>

        <!-- Page Title -->
        <div class="page-title-bar">
          <h1 class="page-title">
            <el-icon :size="22" color="#0056FF"><ShoppingCart /></el-icon>
            购物车
            <span class="item-count" v-if="cartItems.length">({{ cartItems.length }})</span>
          </h1>
        </div>

        <!-- Cart Content -->
        <div v-if="cartItems.length > 0" class="cart-layout">
          <!-- Cart Items Table -->
          <div class="cart-items-panel">
            <el-table :data="cartItems" style="width: 100%" :header-cell-style="headerCellStyle">
              <el-table-column label="商品信息" min-width="220">
                <template #default="{ row }">
                  <div class="item-info-cell">
                    <div class="item-icon-wrap">
                      <el-icon :size="20" color="#0056FF"><Cpu /></el-icon>
                    </div>
                    <div>
                      <div class="item-name">{{ row.name }}</div>
                      <div class="item-id">ID: {{ row.id }}</div>
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="配置摘要" min-width="200">
                <template #default="{ row }">
                  <div class="config-tags">
                    <el-tag size="small" type="info" effect="plain" round>{{ row.cpu }}</el-tag>
                    <el-tag size="small" type="info" effect="plain" round>{{ row.memory }}</el-tag>
                    <el-tag size="small" type="info" effect="plain" round>{{ row.disk }}</el-tag>
                    <el-tag size="small" type="info" effect="plain" round>{{ row.bandwidth }}</el-tag>
                  </div>
                  <div class="config-os">{{ row.os }}</div>
                </template>
              </el-table-column>
              <el-table-column label="周期" width="120">
                <template #default="{ row }">
                  <el-select
                    v-model="row.cycle"
                    size="small"
                    style="width: 90px"
                    @change="updateItemPrice(row)"
                  >
                    <el-option label="月付" value="月付" />
                    <el-option label="季付" value="季付" />
                    <el-option label="半年付" value="半年付" />
                    <el-option label="年付" value="年付" />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="单价" width="110">
                <template #default="{ row }">
                  <span class="price-text">¥{{ row.unitPrice.toFixed(2) }}</span>
                </template>
              </el-table-column>
              <el-table-column label="数量" width="130">
                <template #default="{ row }">
                  <el-input-number
                    v-model="row.quantity"
                    :min="1"
                    :max="20"
                    size="small"
                    style="width: 100px"
                  />
                </template>
              </el-table-column>
              <el-table-column label="小计" width="120">
                <template #default="{ row }">
                  <span class="total-text">¥{{ (row.unitPrice * row.quantity).toFixed(2) }}</span>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="80" align="center">
                <template #default="{ row }">
                  <el-button text type="danger" size="small" @click="removeItem(row.id)">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>

          <!-- Summary Panel -->
          <div class="summary-panel">
            <div class="summary-card">
              <h3 class="summary-title">订单摘要</h3>

              <!-- Price Breakdown with el-descriptions -->
              <el-descriptions :column="1" size="small" border class="price-descriptions">
                <el-descriptions-item label="商品小计">
                  ¥{{ subtotal.toFixed(2) }}
                </el-descriptions-item>
                <el-descriptions-item v-if="discount > 0" label="优惠折扣">
                  <span class="discount-value">-¥{{ discount.toFixed(2) }}</span>
                </el-descriptions-item>
                <el-descriptions-item label="应付总价">
                  <span class="final-total">¥{{ finalTotal.toFixed(2) }}</span>
                </el-descriptions-item>
              </el-descriptions>

              <!-- Applied Coupon -->
              <div v-if="appliedCoupon" class="applied-coupon">
                <div class="coupon-info">
                  <el-icon :size="16" color="#00b42a"><CircleCheckFilled /></el-icon>
                  <span>优惠码: {{ appliedCoupon }}</span>
                </div>
                <el-button text type="danger" size="small" @click="removeCoupon">移除</el-button>
              </div>

              <!-- Coupon Input -->
              <div v-else class="coupon-row">
                <el-input
                  v-model="couponCode"
                  placeholder="请输入优惠码"
                  size="default"
                />
                <el-button class="btn-gradient" size="default" @click="applyCoupon">应用</el-button>
              </div>
              <p v-if="couponMessage" class="coupon-msg" :class="couponSuccess ? 'success' : 'error'">
                {{ couponMessage }}
              </p>

              <el-button
                class="btn-gradient checkout-btn"
                size="large"
                round
                :loading="checkoutLoading"
                @click="handleCheckout"
              >
                去结算
              </el-button>

              <div class="summary-badges">
                <div class="badge-item">
                  <el-icon :size="14" color="#00b42a"><CircleCheckFilled /></el-icon>
                  <span>7天无理由退款</span>
                </div>
                <div class="badge-item">
                  <el-icon :size="14" color="#00b42a"><CircleCheckFilled /></el-icon>
                  <span>安全支付</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty Cart -->
        <div v-else class="empty-cart">
          <el-empty description="购物车是空的">
            <el-button class="btn-gradient" round @click="$router.push('/products')">去选购产品</el-button>
          </el-empty>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Monitor,
  Cpu,
  ShoppingCart,
  Delete,
  CircleCheckFilled
} from '@element-plus/icons-vue'

const router = useRouter()
const checkoutLoading = ref(false)

const headerCellStyle = {
  background: '#EBF3FD',
  color: '#1a3a5c',
  fontWeight: '600',
  fontSize: '13px'
}

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

const cycleMultiplier: Record<string, number> = {
  '月付': 1,
  '季付': 2.8,
  '半年付': 5.2,
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
    basePrice: 289
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
  ElMessage.success('已移除')
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
  ElMessage.success('优惠码已移除')
}

function handleCheckout() {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  checkoutLoading.value = true
  setTimeout(() => {
    checkoutLoading.value = false
    router.push('/checkout')
  }, 800)
}
</script>

<style scoped>
.cart-page {
  min-height: 100vh;
  background: #f5f7fa;
}

/* Header */
.page-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: #fff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
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
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #0056FF 0%, #4080FF 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  font-size: 17px;
  font-weight: 700;
  color: #1a3a5c;
}

.nav-links {
  display: flex;
  gap: 28px;
}

.nav-link {
  color: #666;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: color 0.2s;
  padding: 4px 0;
}

.nav-link:hover,
.nav-link.active {
  color: #0056FF;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

/* Gradient Button */
.btn-gradient {
  background: linear-gradient(135deg, #0056FF 0%, #4080FF 100%) !important;
  border: none !important;
  color: #fff !important;
  font-weight: 500;
}

.btn-gradient:hover {
  opacity: 0.9;
}

/* Main Content */
.main-content {
  padding-top: 84px;
  padding-bottom: 40px;
}

.content-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
}

/* Breadcrumb */
.breadcrumb {
  margin-bottom: 20px;
}

/* Page Title */
.page-title-bar {
  margin-bottom: 20px;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #1a3a5c;
  display: flex;
  align-items: center;
  gap: 10px;
}

.item-count {
  font-size: 15px;
  font-weight: 400;
  color: #999;
}

/* Cart Layout */
.cart-layout {
  display: flex;
  gap: 24px;
  align-items: flex-start;
}

/* Cart Items Panel */
.cart-items-panel {
  flex: 1;
  min-width: 0;
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

:deep(.el-table) {
  --el-table-border-color: #f0f0f0;
}

:deep(.el-table th.el-table__cell) {
  font-weight: 600;
}

:deep(.el-table .el-table__row:hover > td.el-table__cell) {
  background-color: #fafbfc;
}

.item-info-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.item-icon-wrap {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #EBF3FD 0%, #d6e6ff 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.item-name {
  font-size: 14px;
  font-weight: 600;
  color: #1a3a5c;
}

.item-id {
  font-size: 11px;
  color: #ccc;
  margin-top: 2px;
}

.config-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 4px;
}

.config-os {
  font-size: 11px;
  color: #999;
}

.price-text,
.total-text {
  font-size: 13px;
  font-weight: 600;
  color: #1a3a5c;
}

.total-text {
  color: #0056FF;
}

/* Summary Panel */
.summary-panel {
  width: 360px;
  flex-shrink: 0;
}

.summary-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 84px;
}

.summary-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a3a5c;
  margin-bottom: 20px;
}

.price-descriptions {
  margin-bottom: 16px;
}

:deep(.el-descriptions__label) {
  font-weight: 500;
  color: #666;
}

:deep(.el-descriptions__content) {
  font-weight: 600;
  color: #1a3a5c;
}

.discount-value {
  color: #00b42a;
  font-weight: 600;
}

.final-total {
  color: #0056FF;
  font-size: 18px;
  font-weight: 700;
}

.applied-coupon {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 14px;
  background: #e8ffea;
  border: 1px solid #aff0b5;
  border-radius: 8px;
  margin-bottom: 16px;
}

.coupon-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #00b42a;
}

.coupon-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.coupon-msg {
  font-size: 12px;
  margin-bottom: 16px;
}

.coupon-msg.success {
  color: #00b42a;
}

.coupon-msg.error {
  color: #f53f3f;
}

.checkout-btn {
  width: 100%;
  height: 44px;
  font-size: 15px;
  font-weight: 600;
  margin-top: 16px;
}

.summary-badges {
  margin-top: 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.badge-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #999;
}

/* Empty Cart */
.empty-cart {
  padding: 80px 0;
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
}
</style>
