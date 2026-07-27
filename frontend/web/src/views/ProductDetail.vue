<template>
  <div class="detail-page">
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

    <!-- Breadcrumb -->
    <div class="breadcrumb-bar">
      <div class="breadcrumb-inner">
        <n-breadcrumb>
          <n-breadcrumb-item @click="$router.push('/')">首页</n-breadcrumb-item>
          <n-breadcrumb-item @click="$router.push('/products')">产品</n-breadcrumb-item>
          <n-breadcrumb-item>{{ product.name }}</n-breadcrumb-item>
        </n-breadcrumb>
      </div>
    </div>

    <!-- Main -->
    <div class="main-content">
      <div class="content-inner">
        <!-- Left: Config Form -->
        <div class="config-panel">
          <!-- Product Header -->
          <div class="product-header">
            <div class="product-icon-box">
              <n-icon size="36" color="#1890ff"><ServerOutline /></n-icon>
            </div>
            <div>
              <h1 class="product-title">{{ product.name }}</h1>
              <p class="product-desc">{{ product.description }}</p>
            </div>
          </div>

          <!-- Billing Cycle -->
          <div class="config-section">
            <h3 class="config-label">付费周期</h3>
            <div class="cycle-tabs">
              <div
                v-for="cycle in billingCycles"
                :key="cycle.value"
                class="cycle-tab"
                :class="{ active: selectedCycle === cycle.value }"
                @click="selectedCycle = cycle.value"
              >
                <span class="cycle-name">{{ cycle.label }}</span>
                <span class="cycle-price">¥{{ cycle.price }}</span>
                <span class="cycle-savings" v-if="cycle.savings">省{{ cycle.savings }}%</span>
              </div>
            </div>
          </div>

          <!-- Region -->
          <div class="config-section">
            <h3 class="config-label">区域选择</h3>
            <div class="option-grid region-grid">
              <div
                v-for="region in regions"
                :key="region.value"
                class="option-item region-item"
                :class="{ active: selectedRegion === region.value }"
                @click="selectedRegion = region.value"
              >
                <span class="region-flag">{{ region.flag }}</span>
                <span class="option-main">{{ region.label }}</span>
              </div>
            </div>
          </div>

          <!-- CPU -->
          <div class="config-section">
            <h3 class="config-label">CPU</h3>
            <div class="option-grid">
              <div
                v-for="cpu in cpuOptions"
                :key="cpu.value"
                class="option-item"
                :class="{ active: selectedCpu === cpu.value }"
                @click="selectedCpu = cpu.value"
              >
                <span class="option-main">{{ cpu.label }}</span>
                <span class="option-diff" v-if="cpu.diff">{{ cpu.diff }}</span>
              </div>
            </div>
          </div>

          <!-- Memory -->
          <div class="config-section">
            <h3 class="config-label">内存</h3>
            <div class="option-grid">
              <div
                v-for="mem in memoryOptions"
                :key="mem.value"
                class="option-item"
                :class="{ active: selectedMemory === mem.value }"
                @click="selectedMemory = mem.value"
              >
                <span class="option-main">{{ mem.label }}</span>
                <span class="option-diff" v-if="mem.diff">{{ mem.diff }}</span>
              </div>
            </div>
          </div>

          <!-- Disk -->
          <div class="config-section">
            <h3 class="config-label">系统盘</h3>
            <div class="option-grid">
              <div
                v-for="disk in diskOptions"
                :key="disk.value"
                class="option-item"
                :class="{ active: selectedDisk === disk.value }"
                @click="selectedDisk = disk.value"
              >
                <span class="option-main">{{ disk.label }}</span>
                <span class="option-diff" v-if="disk.diff">{{ disk.diff }}</span>
              </div>
            </div>
          </div>

          <!-- Bandwidth -->
          <div class="config-section">
            <h3 class="config-label">带宽</h3>
            <div class="option-grid">
              <div
                v-for="bw in bandwidthOptions"
                :key="bw.value"
                class="option-item"
                :class="{ active: selectedBandwidth === bw.value }"
                @click="selectedBandwidth = bw.value"
              >
                <span class="option-main">{{ bw.label }}</span>
                <span class="option-diff" v-if="bw.diff">{{ bw.diff }}</span>
              </div>
            </div>
          </div>

          <!-- OS -->
          <div class="config-section">
            <h3 class="config-label">操作系统</h3>
            <div class="os-tabs">
              <n-radio-group v-model:value="osCategory" size="small" class="os-category">
                <n-radio-button value="linux">Linux</n-radio-button>
                <n-radio-button value="windows">Windows</n-radio-button>
              </n-radio-group>
            </div>
            <div class="option-grid os-grid">
              <div
                v-for="os in filteredOsOptions"
                :key="os.value"
                class="option-item os-item"
                :class="{ active: selectedOs === os.value }"
                @click="selectedOs = os.value"
              >
                <span class="os-icon">{{ os.icon }}</span>
                <span class="option-main">{{ os.label }}</span>
                <span class="option-diff" v-if="os.diff">{{ os.diff }}</span>
              </div>
            </div>
          </div>

          <!-- Quantity -->
          <div class="config-section">
            <h3 class="config-label">数量</h3>
            <n-input-number
              v-model:value="quantity"
              :min="1"
              :max="20"
              size="large"
              style="width: 160px;"
            >
              <template #minus-icon>
                <n-icon><RemoveOutline /></n-icon>
              </template>
              <template #add-icon>
                <n-icon><AddOutline /></n-icon>
              </template>
            </n-input-number>
          </div>
        </div>

        <!-- Right: Price Summary (Sticky) -->
        <div class="price-panel">
          <div class="price-card">
            <h3 class="price-card-title">已选配置</h3>
            <div class="price-detail-list">
              <div class="price-detail-item">
                <span class="detail-label">付费周期</span>
                <span class="detail-value">{{ currentCycleLabel }}</span>
              </div>
              <div class="price-detail-item">
                <span class="detail-label">区域</span>
                <span class="detail-value">{{ currentRegionLabel }}</span>
              </div>
              <div class="price-detail-item">
                <span class="detail-label">CPU</span>
                <span class="detail-value">{{ selectedCpu }}</span>
              </div>
              <div class="price-detail-item">
                <span class="detail-label">内存</span>
                <span class="detail-value">{{ selectedMemory }}</span>
              </div>
              <div class="price-detail-item">
                <span class="detail-label">系统盘</span>
                <span class="detail-value">{{ selectedDisk }}</span>
              </div>
              <div class="price-detail-item">
                <span class="detail-label">带宽</span>
                <span class="detail-value">{{ selectedBandwidth }}</span>
              </div>
              <div class="price-detail-item">
                <span class="detail-label">操作系统</span>
                <span class="detail-value">{{ currentOsLabel }}</span>
              </div>
              <div class="price-detail-item">
                <span class="detail-label">数量</span>
                <span class="detail-value">{{ quantity }} 台</span>
              </div>
            </div>

            <n-divider />

            <div class="price-row">
              <span class="price-row-label">产品价格</span>
              <span class="price-row-value">¥{{ baseCyclePrice.toFixed(2) }}</span>
            </div>

            <div class="price-row">
              <span class="price-row-label">配置加价</span>
              <span class="price-row-value">¥{{ configExtraPrice.toFixed(2) }}</span>
            </div>

            <div class="coupon-section">
              <div class="coupon-input-row">
                <n-input
                  v-model:value="couponCode"
                  placeholder="请输入优惠码"
                  size="small"
                />
                <n-button type="primary" size="small" @click="applyCoupon" :disabled="!couponCode">
                  验证
                </n-button>
              </div>
              <p v-if="couponMessage" class="coupon-msg" :class="couponSuccess ? 'success' : 'error'">
                {{ couponMessage }}
              </p>
            </div>

            <div v-if="couponDiscount > 0" class="price-row discount-row">
              <span class="price-row-label">优惠码</span>
              <span class="price-row-value discount-value">-¥{{ couponDiscountAmount.toFixed(2) }}</span>
            </div>

            <n-divider />

            <div class="total-section">
              <span class="total-label">总计</span>
              <div class="total-price">
                <span class="total-currency">¥</span>
                <span class="total-amount">{{ totalPrice.toFixed(2) }}</span>
              </div>
            </div>

            <n-button type="primary" size="large" block round class="buy-btn" @click="handleBuy">
              立即购买
            </n-button>
            <n-button size="large" block round class="cart-btn" @click="handleAddToCart">
              加入购物车
            </n-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import {
  AnchorOutline,
  ServerOutline,
  RemoveOutline,
  AddOutline
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()
const message = useMessage()

const productId = Number(route.params.id)

const product = ref({
  id: productId,
  name: '香港云服务器',
  description: 'BGP国际多线，延迟低至10ms，弹性扩展，分钟级创建，适用于Web应用、游戏加速、外贸电商等多种业务场景'
})

// Billing Cycle
const billingCycles = [
  { value: 'monthly', label: '月付', price: 99, savings: 0 },
  { value: 'quarterly', label: '季付', price: 267, savings: 10 },
  { value: 'semiannual', label: '半年付', price: 499, savings: 16 },
  { value: 'annual', label: '年付', price: 899, savings: 25 }
]
const selectedCycle = ref('monthly')

// Region
const regions = [
  { value: 'hongkong', label: '香港', flag: '🇭🇰' },
  { value: 'usa', label: '美国', flag: '🇺🇸' },
  { value: 'japan', label: '日本', flag: '🇯🇵' },
  { value: 'singapore', label: '新加坡', flag: '🇸🇬' }
]
const selectedRegion = ref('hongkong')

// CPU
const cpuOptions = [
  { value: '1核', label: '1核', diff: '' },
  { value: '2核', label: '2核', diff: '+¥20' },
  { value: '4核', label: '4核', diff: '+¥60' },
  { value: '8核', label: '8核', diff: '+¥140' },
  { value: '16核', label: '16核', diff: '+¥300' }
]
const selectedCpu = ref('1核')

// Memory
const memoryOptions = [
  { value: '1G', label: '1G', diff: '' },
  { value: '2G', label: '2G', diff: '+¥15' },
  { value: '4G', label: '4G', diff: '+¥45' },
  { value: '8G', label: '8G', diff: '+¥105' },
  { value: '16G', label: '16G', diff: '+¥225' },
  { value: '32G', label: '32G', diff: '+¥465' }
]
const selectedMemory = ref('1G')

// Disk
const diskOptions = [
  { value: '50G', label: '50G SSD', diff: '' },
  { value: '100G', label: '100G SSD', diff: '+¥10' },
  { value: '200G', label: '200G SSD', diff: '+¥30' }
]
const selectedDisk = ref('50G')

// Bandwidth
const bandwidthOptions = [
  { value: '1M', label: '1Mbps', diff: '' },
  { value: '5M', label: '5Mbps', diff: '+¥20' },
  { value: '10M', label: '10Mbps', diff: '+¥50' },
  { value: '20M', label: '20Mbps', diff: '+¥110' },
  { value: '50M', label: '50Mbps', diff: '+¥280' }
]
const selectedBandwidth = ref('1M')

// OS
const osCategory = ref('linux')
const osOptions = [
  { value: 'centos7', label: 'CentOS 7', diff: '', category: 'linux', icon: '🐧' },
  { value: 'centos8', label: 'CentOS 8', diff: '', category: 'linux', icon: '🐧' },
  { value: 'ubuntu20', label: 'Ubuntu 20.04', diff: '', category: 'linux', icon: '🟠' },
  { value: 'ubuntu22', label: 'Ubuntu 22.04', diff: '', category: 'linux', icon: '🟠' },
  { value: 'debian11', label: 'Debian 11', diff: '', category: 'linux', icon: '🌀' },
  { value: 'debian12', label: 'Debian 12', diff: '', category: 'linux', icon: '🌀' },
  { value: 'win2019', label: 'Windows 2019', diff: '+¥30', category: 'windows', icon: '🪟' },
  { value: 'win2022', label: 'Windows 2022', diff: '+¥30', category: 'windows', icon: '🪟' }
]
const selectedOs = ref('centos7')

const filteredOsOptions = computed(() =>
  osOptions.filter(os => os.category === osCategory.value)
)

// Quantity
const quantity = ref(1)

// Coupon
const couponCode = ref('')
const couponMessage = ref('')
const couponSuccess = ref(false)
const couponDiscount = ref(0)

function applyCoupon() {
  if (!couponCode.value) {
    couponMessage.value = '请输入优惠码'
    couponSuccess.value = false
    return
  }
  if (couponCode.value.toUpperCase() === 'ANCHOR2024') {
    couponMessage.value = '优惠码有效，享受9折优惠'
    couponSuccess.value = true
    couponDiscount.value = 0.1
  } else if (couponCode.value.toUpperCase() === 'NEW10') {
    couponMessage.value = '优惠码有效，立减10%'
    couponSuccess.value = true
    couponDiscount.value = 0.1
  } else {
    couponMessage.value = '优惠码无效'
    couponSuccess.value = false
    couponDiscount.value = 0
  }
}

// Price calculation
function getExtraPrice(options: { value: string; diff: string }[], selected: string): number {
  const opt = options.find(o => o.value === selected)
  if (!opt || !opt.diff) return 0
  const match = opt.diff.match(/\+¥(\d+)/)
  return match ? Number(match[1]) : 0
}

const baseCyclePrice = computed(() => {
  const cycle = billingCycles.find(c => c.value === selectedCycle.value)
  return cycle ? cycle.price : 99
})

const configExtraPrice = computed(() => {
  let extra = 0
  extra += getExtraPrice(cpuOptions, selectedCpu.value)
  extra += getExtraPrice(memoryOptions, selectedMemory.value)
  extra += getExtraPrice(diskOptions, selectedDisk.value)
  extra += getExtraPrice(bandwidthOptions, selectedBandwidth.value)
  extra += getExtraPrice(osOptions as any, selectedOs.value)
  return extra
})

const couponDiscountAmount = computed(() => {
  const subtotal = (baseCyclePrice.value + configExtraPrice.value) * quantity.value
  return subtotal * couponDiscount.value
})

const totalPrice = computed(() => {
  const subtotal = (baseCyclePrice.value + configExtraPrice.value) * quantity.value
  return subtotal - couponDiscountAmount.value
})

// Labels
const currentCycleLabel = computed(() => {
  const cycle = billingCycles.find(c => c.value === selectedCycle.value)
  return cycle ? `${cycle.label} ¥${cycle.price}` : ''
})

const currentRegionLabel = computed(() =>
  regions.find(r => r.value === selectedRegion.value)?.label || ''
)

const currentOsLabel = computed(() =>
  osOptions.find(o => o.value === selectedOs.value)?.label || ''
)

function handleBuy() {
  const token = localStorage.getItem('token')
  if (!token) {
    message.warning('请先登录')
    router.push('/login')
    return
  }
  message.success('订单已创建，正在跳转...')
  router.push('/cart')
}

function handleAddToCart() {
  message.success('已加入购物车')
}
</script>

<style scoped>
.detail-page {
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

/* Breadcrumb */
.breadcrumb-bar {
  background: #fff;
  border-bottom: 1px solid #f0f1f5;
  margin-top: 64px;
}

.breadcrumb-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 14px 24px;
}

/* Main Content */
.main-content {
  padding: 24px 0 40px;
}

.content-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  gap: 24px;
}

/* Config Panel */
.config-panel {
  flex: 1;
  min-width: 0;
}

.product-header {
  display: flex;
  align-items: center;
  gap: 20px;
  background: #fff;
  border-radius: 12px;
  padding: 28px;
  margin-bottom: 20px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.product-icon-box {
  width: 72px;
  height: 72px;
  background: linear-gradient(135deg, #e6f7ff, #bae7ff);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.product-title {
  font-size: 24px;
  font-weight: 700;
  color: #1d2129;
  margin-bottom: 6px;
}

.product-desc {
  font-size: 14px;
  color: #86909c;
  line-height: 1.6;
}

.config-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px 28px;
  margin-bottom: 16px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.config-label {
  font-size: 15px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 16px;
}

/* Cycle Tabs */
.cycle-tabs {
  display: flex;
  gap: 12px;
}

.cycle-tab {
  flex: 1;
  text-align: center;
  padding: 14px 12px;
  border: 2px solid #f0f1f5;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.cycle-tab:hover {
  border-color: #bae7ff;
}

.cycle-tab.active {
  border-color: #1890ff;
  background: linear-gradient(135deg, #e6f7ff, #f0faff);
}

.cycle-name {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #86909c;
  margin-bottom: 4px;
}

.cycle-tab.active .cycle-name {
  color: #1890ff;
  font-weight: 600;
}

.cycle-price {
  display: block;
  font-size: 20px;
  font-weight: 700;
  color: #1d2129;
}

.cycle-tab.active .cycle-price {
  color: #1890ff;
}

.cycle-savings {
  display: inline-block;
  font-size: 11px;
  color: #ff7a45;
  background: #fff7e6;
  border-radius: 4px;
  padding: 1px 6px;
  margin-top: 6px;
}

.cycle-tab.active .cycle-savings {
  background: #fff1f0;
}

/* Option Grid */
.option-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 18px;
  border: 2px solid #f0f1f5;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
  color: #1d2129;
}

.option-item:hover {
  border-color: #bae7ff;
}

.option-item.active {
  border-color: #1890ff;
  background: #e6f7ff;
  color: #1890ff;
}

.option-main {
  font-weight: 500;
}

.option-diff {
  font-size: 12px;
  color: #ff7a45;
}

.option-item.active .option-diff {
  color: #fa541c;
}

/* Region */
.region-item {
  min-width: 120px;
  justify-content: center;
}

.region-flag {
  font-size: 20px;
  line-height: 1;
}

/* OS */
.os-tabs {
  margin-bottom: 14px;
}

.os-grid {
  margin-top: 4px;
}

.os-item {
  min-width: 140px;
}

.os-icon {
  font-size: 18px;
  line-height: 1;
}

/* Price Panel */
.price-panel {
  width: 320px;
  flex-shrink: 0;
}

.price-card {
  background: #fff;
  border-radius: 12px;
  padding: 28px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
  position: sticky;
  top: 88px;
}

.price-card-title {
  font-size: 17px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 20px;
}

.price-detail-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.price-detail-item {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.detail-label {
  color: #86909c;
}

.detail-value {
  color: #1d2129;
  font-weight: 500;
}

.price-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.price-row-label {
  font-size: 14px;
  color: #4e5969;
}

.price-row-value {
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
}

.discount-value {
  color: #52c41a;
}

.coupon-section {
  margin-bottom: 12px;
}

.coupon-input-row {
  display: flex;
  gap: 8px;
}

.coupon-input-row .n-input {
  flex: 1;
}

.coupon-msg {
  font-size: 12px;
  margin-top: 6px;
}

.coupon-msg.success {
  color: #52c41a;
}

.coupon-msg.error {
  color: #ff4d4f;
}

.total-section {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 20px;
}

.total-label {
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
}

.total-price {
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

.buy-btn {
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border: none;
}

.buy-btn:hover {
  background: linear-gradient(135deg, #40a9ff, #1890ff);
}

.cart-btn {
  margin-top: 10px;
  height: 44px;
  font-size: 15px;
  font-weight: 500;
  border-color: #1890ff;
  color: #1890ff;
}

.cart-btn:hover {
  background: #e6f7ff;
  border-color: #40a9ff;
  color: #40a9ff;
}

/* Responsive */
@media (max-width: 1024px) {
  .content-inner {
    flex-direction: column;
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
  .cycle-tabs {
    flex-wrap: wrap;
  }
  .cycle-tab {
    min-width: calc(50% - 6px);
  }
}
</style>
