<template>
  <div class="detail-page">
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

    <!-- Breadcrumb -->
    <div class="breadcrumb-bar">
      <div class="breadcrumb-inner">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: '/products' }">产品</el-breadcrumb-item>
          <el-breadcrumb-item>{{ product.name }}</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
    </div>

    <!-- Main Content -->
    <div class="main-content">
      <div class="content-inner">
        <!-- Left: Config Form -->
        <div class="config-panel">
          <!-- Product Header -->
          <div class="product-header">
            <div class="product-icon-box">
              <el-icon :size="32" color="#0056FF"><Cpu /></el-icon>
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
            <h3 class="config-label">数据中心</h3>
            <div class="option-grid region-grid">
              <div
                v-for="region in regions"
                :key="region.value"
                class="option-item region-item"
                :class="{ active: selectedRegion === region.value }"
                @click="selectedRegion = region.value"
              >
                <img :src="region.flag" class="region-flag-img" :alt="region.label" />
                <span class="option-main">{{ region.label }}</span>
              </div>
            </div>
          </div>

          <!-- OS -->
          <div class="config-section">
            <h3 class="config-label">操作系统</h3>
            <div class="os-tabs">
              <el-radio-group v-model="selectedOsType" size="default">
                <el-radio-button v-for="osType in osTypes" :key="osType.value" :value="osType.value">
                  {{ osType.label }}
                </el-radio-button>
              </el-radio-group>
            </div>
            <div class="os-version-row">
              <span class="os-version-label">版本</span>
              <el-select v-model="selectedOsVersion" style="width: 280px" size="default">
                <el-option v-for="ver in currentOsVersions" :key="ver" :label="ver" :value="ver" />
              </el-select>
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

          <!-- Features -->
          <div class="config-section">
            <h3 class="config-label">产品特性</h3>
            <div class="features-grid">
              <div v-for="(feature, idx) in productFeatures" :key="idx" class="feature-item">
                <el-icon :size="16" color="#0056FF"><component :is="feature.icon" /></el-icon>
                <span>{{ feature.text }}</span>
              </div>
            </div>
          </div>

          <!-- Quantity -->
          <div class="config-section">
            <h3 class="config-label">数量</h3>
            <el-input-number
              v-model="quantity"
              :min="1"
              :max="20"
              size="large"
              style="width: 160px"
            />
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
                <span class="detail-label">数据中心</span>
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
                <span class="detail-value">{{ selectedOsVersion }}</span>
              </div>
              <div class="price-detail-item">
                <span class="detail-label">数量</span>
                <span class="detail-value">{{ quantity }} 台</span>
              </div>
            </div>

            <el-divider />

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
                <el-input
                  v-model="couponCode"
                  placeholder="请输入优惠码"
                  size="small"
                />
                <el-button class="btn-gradient" size="small" @click="applyCoupon" :disabled="!couponCode">
                  验证
                </el-button>
              </div>
              <p v-if="couponMessage" class="coupon-msg" :class="couponSuccess ? 'success' : 'error'">
                {{ couponMessage }}
              </p>
            </div>

            <div v-if="couponDiscount > 0" class="price-row discount-row">
              <span class="price-row-label">优惠码</span>
              <span class="price-row-value discount-value">-¥{{ couponDiscountAmount.toFixed(2) }}</span>
            </div>

            <el-divider />

            <div class="total-section">
              <span class="total-label">总计</span>
              <div class="total-price">
                <span class="total-currency">¥</span>
                <span class="total-amount">{{ totalPrice.toFixed(2) }}</span>
              </div>
            </div>

            <el-button
              class="btn-gradient buy-btn"
              size="large"
              round
              :loading="buyLoading"
              @click="handleBuy"
            >
              立即购买
            </el-button>
            <el-button
              size="large"
              round
              class="cart-btn"
              :loading="cartLoading"
              @click="handleAddToCart"
            >
              加入购物车
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Monitor,
  Cpu,
  Check,
  Lightning,
  Shield,
  Connection,
  Timer,
  TrendCharts
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

const productId = Number(route.params.id)
const buyLoading = ref(false)
const cartLoading = ref(false)

const product = ref({
  id: productId,
  name: '香港云服务器',
  description: 'BGP国际多线，延迟低至10ms，弹性扩展，分钟级创建，适用于Web应用、游戏加速、外贸电商等多种业务场景'
})

const productFeatures = [
  { icon: Lightning, text: '高性能NVMe SSD存储' },
  { icon: Shield, text: 'DDoS防护保障' },
  { icon: Connection, text: 'BGP多线接入' },
  { icon: Timer, text: '99.9% SLA保障' },
  { icon: TrendCharts, text: '弹性伸缩' },
  { icon: Check, text: '7天无理由退款' }
]

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
  { value: 'HK', label: '香港', flag: '/assets/flags/HK.png' },
  { value: 'US', label: '美国', flag: '/assets/flags/US.png' },
  { value: 'JP', label: '日本', flag: '/assets/flags/JP.png' },
  { value: 'SG', label: '新加坡', flag: '/assets/flags/SG.png' }
]
const selectedRegion = ref('HK')

// OS
const osTypes = [
  { value: 'centos', label: 'CentOS' },
  { value: 'ubuntu', label: 'Ubuntu' },
  { value: 'debian', label: 'Debian' },
  { value: 'windows', label: 'Windows' }
]
const selectedOsType = ref('centos')
const selectedOsVersion = ref('CentOS 7.9')

const osVersionsMap: Record<string, string[]> = {
  centos: ['CentOS 7.9', 'CentOS 8.5', 'CentOS Stream 9'],
  ubuntu: ['Ubuntu 20.04', 'Ubuntu 22.04', 'Ubuntu 24.04'],
  debian: ['Debian 10', 'Debian 11', 'Debian 12'],
  windows: ['Windows Server 2019', 'Windows Server 2022']
}

const currentOsVersions = computed(() => osVersionsMap[selectedOsType.value] || [])

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
  { value: '200G', label: '200G SSD', diff: '+¥30' },
  { value: '500G', label: '500G SSD', diff: '+¥80' }
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

function handleBuy() {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  buyLoading.value = true
  setTimeout(() => {
    buyLoading.value = false
    ElMessage.success('订单已创建，正在跳转...')
    router.push('/checkout')
  }, 1000)
}

function handleAddToCart() {
  cartLoading.value = true
  setTimeout(() => {
    cartLoading.value = false
    ElMessage.success('已加入购物车')
  }, 800)
}
</script>

<style scoped>
.detail-page {
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

/* Breadcrumb */
.breadcrumb-bar {
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  margin-top: 60px;
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
  padding: 24px;
  margin-bottom: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.product-icon-box {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #EBF3FD 0%, #d6e6ff 100%);
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.product-title {
  font-size: 22px;
  font-weight: 700;
  color: #1a3a5c;
  margin-bottom: 6px;
}

.product-desc {
  font-size: 13px;
  color: #999;
  line-height: 1.6;
}

.config-section {
  background: #fff;
  border-radius: 12px;
  padding: 20px 24px;
  margin-bottom: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.config-label {
  font-size: 14px;
  font-weight: 600;
  color: #1a3a5c;
  margin-bottom: 14px;
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
  border: 2px solid #f0f0f0;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.cycle-tab:hover {
  border-color: #b3d4ff;
}

.cycle-tab.active {
  border-color: #0056FF;
  background: linear-gradient(135deg, #EBF3FD 0%, #f0f5ff 100%);
}

.cycle-name {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #999;
  margin-bottom: 4px;
}

.cycle-tab.active .cycle-name {
  color: #0056FF;
  font-weight: 600;
}

.cycle-price {
  display: block;
  font-size: 20px;
  font-weight: 700;
  color: #1a3a5c;
}

.cycle-tab.active .cycle-price {
  color: #0056FF;
}

.cycle-savings {
  display: inline-block;
  font-size: 11px;
  color: #ff6b35;
  background: #fff7ed;
  border-radius: 4px;
  padding: 1px 6px;
  margin-top: 6px;
}

.cycle-tab.active .cycle-savings {
  background: #fff1eb;
}

/* OS Tabs */
.os-tabs {
  margin-bottom: 14px;
}

.os-version-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.os-version-label {
  font-size: 13px;
  color: #666;
  white-space: nowrap;
}

:deep(.el-radio-group .el-radio-button__inner) {
  border-color: #0056FF;
  color: #0056FF;
}

:deep(.el-radio-group .el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background: linear-gradient(135deg, #0056FF 0%, #4080FF 100%);
  border-color: #0056FF;
  color: #fff;
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
  border: 2px solid #f0f0f0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 13px;
  color: #1a3a5c;
}

.option-item:hover {
  border-color: #b3d4ff;
}

.option-item.active {
  border-color: #0056FF;
  background: #EBF3FD;
  color: #0056FF;
}

.option-main {
  font-weight: 500;
}

.option-diff {
  font-size: 11px;
  color: #ff6b35;
}

.option-item.active .option-diff {
  color: #ff4757;
}

/* Region */
.region-item {
  min-width: 120px;
  justify-content: center;
}

.region-flag-img {
  width: 24px;
  height: 16px;
  object-fit: cover;
  border-radius: 2px;
}

/* Features */
.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #666;
}

/* Price Panel */
.price-panel {
  width: 320px;
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

.price-detail-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.price-detail-item {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.detail-label {
  color: #999;
}

.detail-value {
  color: #1a3a5c;
  font-weight: 500;
}

.price-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.price-row-label {
  font-size: 13px;
  color: #666;
}

.price-row-value {
  font-size: 15px;
  font-weight: 600;
  color: #1a3a5c;
}

.discount-value {
  color: #00b42a;
}

.coupon-section {
  margin-bottom: 10px;
}

.coupon-input-row {
  display: flex;
  gap: 8px;
}

.coupon-input-row .el-input {
  flex: 1;
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
  margin-bottom: 18px;
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

.buy-btn {
  width: 100%;
  height: 44px;
  font-size: 15px;
  font-weight: 600;
}

.cart-btn {
  width: 100%;
  margin-top: 10px;
  margin-left: 0 !important;
  height: 40px;
  font-size: 14px;
  font-weight: 500;
  border-color: #0056FF;
  color: #0056FF;
}

.cart-btn:hover {
  background: #EBF3FD;
  border-color: #4080FF;
  color: #4080FF;
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
  .features-grid {
    grid-template-columns: 1fr;
  }
}
</style>
