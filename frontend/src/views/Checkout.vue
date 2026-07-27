<template>
  <div class="checkout-page">
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
          <el-breadcrumb-item :to="{ path: '/cart' }">购物车</el-breadcrumb-item>
          <el-breadcrumb-item>结算</el-breadcrumb-item>
        </el-breadcrumb>

        <!-- Page Title -->
        <div class="page-title-bar">
          <h1 class="page-title">
            <el-icon :size="22" color="#0056FF"><CreditCard /></el-icon>
            确认订单
          </h1>
        </div>

        <div class="checkout-layout">
          <!-- Left: Order Details -->
          <div class="order-panel">
            <!-- Order Items -->
            <div class="section-card">
              <h3 class="section-title">
                <el-icon :size="18" color="#0056FF"><ShoppingCart /></el-icon>
                商品信息
              </h3>
              <el-descriptions :column="2" border size="small" :label-style="{ width: '120px', fontWeight: 600 }">
                <template v-for="item in orderItems" :key="item.id">
                  <el-descriptions-item :label="item.name">
                    {{ item.cpu }} / {{ item.memory }} / {{ item.disk }} / {{ item.bandwidth }}
                  </el-descriptions-item>
                  <el-descriptions-item label="周期/数量">
                    {{ item.cycle }} x {{ item.quantity }}
                  </el-descriptions-item>
                  <el-descriptions-item label="小计">
                    <span class="item-price">¥{{ (item.unitPrice * item.quantity).toFixed(2) }}</span>
                  </el-descriptions-item>
                </template>
              </el-descriptions>
            </div>

            <!-- Payment Method -->
            <div class="section-card">
              <h3 class="section-title">
                <el-icon :size="18" color="#0056FF"><Wallet /></el-icon>
                选择支付方式
              </h3>
              <el-radio-group v-model="selectedPayment" class="payment-radio-group">
                <div
                  v-for="method in paymentMethods"
                  :key="method.value"
                  class="payment-method-card"
                  :class="{ active: selectedPayment === method.value }"
                  @click="selectedPayment = method.value"
                >
                  <el-radio :value="method.value" class="payment-radio">
                    <div class="method-content">
                      <div class="method-icon" :style="{ background: method.bgColor }">
                        <img v-if="method.icon" :src="method.icon" class="method-icon-img" :alt="method.label" />
                        <el-icon v-else :size="24" color="#f59e0b"><Wallet /></el-icon>
                      </div>
                      <span class="method-name">{{ method.label }}</span>
                    </div>
                  </el-radio>
                </div>
              </el-radio-group>
            </div>

            <!-- Payment Countdown -->
            <div class="countdown-bar" v-if="countdownSeconds > 0">
              <el-icon :size="16" color="#ff6b35"><Timer /></el-icon>
              <span>请在 <strong class="countdown-num">{{ formatTime(countdownSeconds) }}</strong> 内完成支付，超时订单将自动关闭</span>
            </div>
          </div>

          <!-- Right: Price Summary -->
          <div class="price-panel">
            <div class="price-card">
              <h3 class="price-card-title">价格明细</h3>

              <el-descriptions :column="1" size="small" border class="price-descriptions">
                <el-descriptions-item label="商品小计">
                  ¥{{ subtotal.toFixed(2) }}
                </el-descriptions-item>
                <el-descriptions-item v-if="appliedCoupon" label="优惠码">
                  <span class="discount-text">-¥{{ discount.toFixed(2) }}</span>
                  <span style="font-size: 12px; color: #999; margin-left: 8px">({{ appliedCoupon }})</span>
                </el-descriptions-item>
              </el-descriptions>

              <!-- Coupon Input -->
              <div class="coupon-block">
                <div v-if="appliedCoupon" class="applied-coupon">
                  <div class="coupon-info">
                    <el-icon :size="14" color="#00b42a"><CircleCheckFilled /></el-icon>
                    <span>优惠码: {{ appliedCoupon }}</span>
                  </div>
                  <el-button text type="danger" size="small" @click="removeCoupon">移除</el-button>
                </div>
                <div v-else class="coupon-input-row">
                  <el-input
                    v-model="couponCode"
                    placeholder="请输入优惠码"
                    size="default"
                    :disabled="!!appliedCoupon"
                  />
                  <el-button class="btn-gradient" size="default" @click="applyCoupon" :disabled="!couponCode">
                    应用
                  </el-button>
                </div>
                <p v-if="couponMessage" class="coupon-msg" :class="couponSuccess ? 'success' : 'error'">
                  {{ couponMessage }}
                </p>
              </div>

              <el-divider />

              <div class="total-section">
                <span class="total-label">应付金额</span>
                <div class="total-price">
                  <span class="total-currency">¥</span>
                  <span class="total-amount">{{ finalTotal.toFixed(2) }}</span>
                </div>
              </div>

              <el-button
                class="btn-gradient submit-btn"
                size="large"
                round
                :loading="submitLoading"
                @click="handleSubmit"
              >
                提交订单
              </el-button>

              <p class="agreement">
                提交订单即表示您同意
                <a href="#" class="link">《服务协议》</a>
                和
                <a href="#" class="link">《隐私政策》</a>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Monitor,
  Cpu,
  CreditCard,
  ShoppingCart,
  Wallet,
  Timer,
  CircleCheckFilled
} from '@element-plus/icons-vue'

const router = useRouter()
const submitLoading = ref(false)

interface OrderItem {
  id: number
  name: string
  cpu: string
  memory: string
  disk: string
  bandwidth: string
  cycle: string
  quantity: number
  unitPrice: number
}

const orderItems = ref<OrderItem[]>([
  {
    id: 10001,
    name: '云服务器 ECS',
    cpu: '2核',
    memory: '4G',
    disk: '100G SSD',
    bandwidth: '5Mbps',
    cycle: '月付',
    quantity: 1,
    unitPrice: 134
  },
  {
    id: 10002,
    name: '云服务器 ECS',
    cpu: '4核',
    memory: '8G',
    disk: '200G SSD',
    bandwidth: '10Mbps',
    cycle: '年付',
    quantity: 2,
    unitPrice: 2890
  }
])

const paymentMethods = [
  { value: 'alipay', label: '支付宝', icon: '/assets/payment/alipay.svg', bgColor: '#e8f4fd' },
  { value: 'wechat', label: '微信支付', icon: '/assets/payment/wechat.svg', bgColor: '#e8f8ea' },
  { value: 'qqpay', label: 'QQ钱包', icon: '/assets/payment/qqpay.svg', bgColor: '#eff6ff' },
  { value: 'usdt', label: 'USDT', icon: '/assets/payment/usdt.svg', bgColor: '#e8faf3' },
  { value: 'balance', label: '余额支付', icon: '', bgColor: '#fff7ed' }
]
const selectedPayment = ref('alipay')

// Coupon
const couponCode = ref('')
const appliedCoupon = ref('')
const couponMessage = ref('')
const couponSuccess = ref(false)
const discountRate = ref(0)

const subtotal = computed(() =>
  orderItems.value.reduce((sum, item) => sum + item.unitPrice * item.quantity, 0)
)

const discount = computed(() => subtotal.value * discountRate.value)

const finalTotal = computed(() => subtotal.value - discount.value)

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

// Countdown
const countdownSeconds = ref(900) // 15 minutes
let countdownTimer: ReturnType<typeof setInterval> | null = null

function formatTime(seconds: number): string {
  const m = Math.floor(seconds / 60)
  const s = seconds % 60
  return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
}

function handleSubmit() {
  if (!selectedPayment.value) {
    ElMessage.warning('请选择支付方式')
    return
  }
  submitLoading.value = true
  setTimeout(() => {
    submitLoading.value = false
    const orderNo = 'AF' + Date.now()
    ElMessage.success('订单提交成功')
    router.push({
      path: '/payment-result',
      query: {
        status: 'success',
        orderNo,
        amount: finalTotal.value.toFixed(2),
        method: paymentMethods.find(m => m.value === selectedPayment.value)?.label || '支付宝'
      }
    })
  }, 1500)
}

onMounted(() => {
  countdownTimer = setInterval(() => {
    if (countdownSeconds.value > 0) {
      countdownSeconds.value--
    } else {
      ElMessage.warning('订单已超时关闭')
      router.push('/cart')
    }
  }, 1000)
})

onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
  }
})
</script>

<style scoped>
.checkout-page {
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

/* Checkout Layout */
.checkout-layout {
  display: flex;
  gap: 24px;
  align-items: flex-start;
}

/* Order Panel */
.order-panel {
  flex: 1;
  min-width: 0;
}

.section-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: #1a3a5c;
  margin-bottom: 18px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.item-price {
  font-weight: 700;
  color: #0056FF;
  font-size: 15px;
}

/* Payment Methods */
.payment-radio-group {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  width: 100%;
}

.payment-method-card {
  border: 2px solid #f0f0f0;
  border-radius: 12px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s;
}

.payment-method-card:hover {
  border-color: #b3d4ff;
}

.payment-method-card.active {
  border-color: #0056FF;
  background: #EBF3FD;
}

.payment-radio {
  width: 100%;
  height: auto;
}

:deep(.el-radio__input) {
  display: none;
}

:deep(.el-radio__label) {
  padding-left: 0;
}

.method-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.method-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.method-icon-img {
  width: 24px;
  height: 24px;
  object-fit: contain;
}

.method-name {
  font-size: 14px;
  font-weight: 500;
  color: #1a3a5c;
}

/* Countdown */
.countdown-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 18px;
  background: #fff7ed;
  border: 1px solid #ffd8b8;
  border-radius: 10px;
  font-size: 13px;
  color: #666;
}

.countdown-num {
  color: #ff6b35;
  font-size: 16px;
  font-weight: 700;
  font-family: 'Courier New', monospace;
}

/* Price Panel */
.price-panel {
  width: 380px;
  flex-shrink: 0;
}

.price-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 84px;
}

.price-card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a3a5c;
  margin-bottom: 18px;
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

.discount-text {
  color: #00b42a;
  font-weight: 600;
}

.coupon-block {
  margin-bottom: 12px;
}

.applied-coupon {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 14px;
  background: #e8ffea;
  border: 1px solid #aff0b5;
  border-radius: 8px;
}

.coupon-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #00b42a;
}

.coupon-input-row {
  display: flex;
  gap: 8px;
}

.coupon-msg {
  font-size: 12px;
  margin-top: 6px;
}

.coupon-msg.success {
  color: #00b42a;
}

.coupon-msg.error {
  color: #f53f3f;
}

.total-section {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 20px;
}

.total-label {
  font-size: 15px;
  font-weight: 600;
  color: #1a3a5c;
}

.total-price {
  display: flex;
  align-items: baseline;
}

.total-currency {
  font-size: 16px;
  font-weight: 600;
  color: #0056FF;
}

.total-amount {
  font-size: 30px;
  font-weight: 700;
  color: #0056FF;
  line-height: 1;
}

.submit-btn {
  width: 100%;
  height: 46px;
  font-size: 16px;
  font-weight: 600;
}

.agreement {
  text-align: center;
  font-size: 12px;
  color: #ccc;
  margin-top: 14px;
}

.link {
  color: #0056FF;
  text-decoration: none;
}

.link:hover {
  text-decoration: underline;
}

/* Responsive */
@media (max-width: 1024px) {
  .checkout-layout {
    flex-direction: column-reverse;
  }
  .price-panel {
    width: 100%;
  }
  .price-card {
    position: static;
  }
}

@media (max-width: 768px) {
  .nav-links {
    display: none;
  }
  .payment-radio-group {
    grid-template-columns: 1fr;
  }
}
</style>
