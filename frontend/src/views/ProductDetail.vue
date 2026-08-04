<template>
  <div class="product-detail-page" v-loading="loading">
    <!-- Header -->
    <header class="page-header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <img src="/logo.png" :alt="siteName" class="logo-img" />
          <span class="logo-text">{{ siteName }}</span>
        </router-link>
        <nav class="nav-links">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/products" class="nav-link">产品</router-link>
          <span class="nav-link active">{{ product.name || '产品详情' }}</span>
        </nav>
        <div class="header-actions">
          <el-button text @click="$router.push('/login')">登录</el-button>
          <el-button type="primary" round size="small" @click="$router.push('/register')">免费注册</el-button>
        </div>
      </div>
    </header>

    <div class="detail-content">
      <div class="detail-inner">
        <!-- 面包屑 -->
        <el-breadcrumb separator="/" class="breadcrumb">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: '/products' }">产品中心</el-breadcrumb-item>
          <el-breadcrumb-item>{{ product.name }}</el-breadcrumb-item>
        </el-breadcrumb>

        <div class="detail-grid">
          <!-- 左侧配置区 -->
          <div class="config-section">
            <div class="config-card">
              <h2>{{ product.name }}</h2>
              <p class="product-desc">{{ product.description }}</p>

              <!-- 产品特性 -->
              <div class="features-list" v-if="product.features?.length">
                <div class="feature-item" v-for="feature in product.features" :key="feature">
                  <el-icon color="#1a73e8"><CircleCheckFilled /></el-icon>
                  <span>{{ feature }}</span>
                </div>
              </div>

              <el-divider />

              <!-- 数据中心 -->
              <div class="config-group" v-if="regions.length">
                <h3>数据中心</h3>
                <div class="region-grid">
                  <div
                    v-for="region in regions"
                    :key="region.id"
                    class="region-item"
                    :class="{ active: selectedRegion === region.id }"
                    @click="selectedRegion = region.id"
                  >
                    <img :src="`/assets/flags/${region.country_code}.png`" :alt="region.name" class="region-flag" />
                    <span>{{ region.name }}</span>
                  </div>
                </div>
              </div>

              <!-- 操作系统 -->
              <div class="config-group" v-if="osTypes.length">
                <h3>操作系统</h3>
                <el-tabs v-model="selectedOsType">
                  <el-tab-pane v-for="os in osTypes" :key="os.id" :label="os.name" :name="os.id">
                    <div class="os-grid">
                      <div
                        v-for="version in os.versions"
                        :key="version.id"
                        class="os-item"
                        :class="{ active: selectedOs === version.id }"
                        @click="selectedOs = version.id"
                      >
                        <img :src="`/assets/os/${os.name}.svg`" :alt="os.name" class="os-icon" />
                        <span>{{ version.name }}</span>
                      </div>
                    </div>
                  </el-tab-pane>
                </el-tabs>
              </div>

              <!-- CPU -->
              <div class="config-group" v-if="cpuOptions.length">
                <h3>CPU</h3>
                <el-radio-group v-model="selectedCpu">
                  <el-radio-button v-for="cpu in cpuOptions" :key="cpu.id" :label="cpu.id">
                    {{ cpu.name }}
                    <span class="price-add" v-if="cpu.price_add">+¥{{ cpu.price_add }}/月</span>
                  </el-radio-button>
                </el-radio-group>
              </div>

              <!-- 内存 -->
              <div class="config-group" v-if="memoryOptions.length">
                <h3>内存</h3>
                <el-radio-group v-model="selectedMemory">
                  <el-radio-button v-for="mem in memoryOptions" :key="mem.id" :label="mem.id">
                    {{ mem.name }}
                    <span class="price-add" v-if="mem.price_add">+¥{{ mem.price_add }}/月</span>
                  </el-radio-button>
                </el-radio-group>
              </div>

              <!-- 硬盘 -->
              <div class="config-group" v-if="diskOptions.length">
                <h3>硬盘</h3>
                <el-radio-group v-model="selectedDisk">
                  <el-radio-button v-for="disk in diskOptions" :key="disk.id" :label="disk.id">
                    {{ disk.name }}
                    <span class="price-add" v-if="disk.price_add">+¥{{ disk.price_add }}/月</span>
                  </el-radio-button>
                </el-radio-group>
              </div>

              <!-- 带宽 -->
              <div class="config-group" v-if="bandwidthOptions.length">
                <h3>带宽</h3>
                <el-radio-group v-model="selectedBandwidth">
                  <el-radio-button v-for="bw in bandwidthOptions" :key="bw.id" :label="bw.id">
                    {{ bw.name }}
                    <span class="price-add" v-if="bw.price_add">+¥{{ bw.price_add }}/月</span>
                  </el-radio-button>
                </el-radio-group>
              </div>
            </div>
          </div>

          <!-- 右侧价格区 -->
          <div class="price-section">
            <div class="price-card">
              <h3>价格信息</h3>
              
              <!-- 计费周期 -->
              <div class="billing-cycles" v-if="billingCycles.length">
                <h4>计费周期</h4>
                <div
                  v-for="cycle in billingCycles"
                  :key="cycle.id"
                  class="cycle-item"
                  :class="{ active: selectedCycle === cycle.id, discount: cycle.discount }"
                  @click="selectedCycle = cycle.id"
                >
                  <span class="cycle-name">{{ cycle.name }}</span>
                  <span class="cycle-price">¥{{ cycle.price }}</span>
                  <span class="cycle-discount" v-if="cycle.discount">{{ cycle.discount }}</span>
                </div>
              </div>

              <el-divider />

              <!-- 价格明细 -->
              <div class="price-summary">
                <div class="summary-item">
                  <span>基础价格</span>
                  <span>¥{{ basePrice }}</span>
                </div>
                <div class="summary-item" v-if="configPrice > 0">
                  <span>配置加价</span>
                  <span>+¥{{ configPrice }}</span>
                </div>
                <div class="summary-item total">
                  <span>总计</span>
                  <span class="total-price">¥{{ totalPrice }}</span>
                </div>
              </div>

              <!-- 优惠码 -->
              <div class="coupon-input">
                <el-input v-model="couponCode" placeholder="请输入优惠码">
                  <template #append>
                    <el-button @click="applyCoupon">应用</el-button>
                  </template>
                </el-input>
              </div>

              <el-button type="primary" size="large" round class="buy-btn" @click="addToCart">
                <el-icon style="margin-right: 6px;"><ShoppingCart /></el-icon>
                立即购买
              </el-button>

              <div class="price-tips">
                <p><el-icon><CircleCheck /></el-icon> 7天无理由退款</p>
                <p><el-icon><CircleCheck /></el-icon> 99.9%可用性保障</p>
                <p><el-icon><CircleCheck /></el-icon> 7×24技术支持</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { CircleCheckFilled, CircleCheck, ShoppingCart } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const siteName = ref('')

const fetchSiteName = async () => {
  try {
    const res = await request.get('/api/v1/settings/public')
    if (res?.data?.site_name) {
      siteName.value = res.data.site_name
    }
  } catch {
    // Use empty
  }
}

// 产品数据 - 从API获取
const product = ref({
  id: route.params.id,
  name: '',
  description: '',
  features: [],
  price: 0
})

const regions = ref([])
const osTypes = ref([])
const billingCycles = ref([])
const cpuOptions = ref([])
const memoryOptions = ref([])
const diskOptions = ref([])
const bandwidthOptions = ref([])

// 选择的配置
const selectedRegion = ref('')
const selectedOsType = ref('')
const selectedOs = ref('')
const selectedCpu = ref('')
const selectedMemory = ref('')
const selectedDisk = ref('')
const selectedBandwidth = ref('')
const selectedCycle = ref('')
const couponCode = ref('')

// 价格计算
const basePrice = computed(() => {
  const cycle = billingCycles.value.find(c => c.id === selectedCycle.value)
  return cycle?.price || product.value.price || 0
})

const configPrice = computed(() => {
  let price = 0
  const cpu = cpuOptions.value.find(c => c.id === selectedCpu.value)
  const mem = memoryOptions.value.find(m => m.id === selectedMemory.value)
  const disk = diskOptions.value.find(d => d.id === selectedDisk.value)
  const bw = bandwidthOptions.value.find(b => b.id === selectedBandwidth.value)
  
  if (cpu?.price_add) price += cpu.price_add
  if (mem?.price_add) price += mem.price_add
  if (disk?.price_add) price += disk.price_add
  if (bw?.price_add) price += bw.price_add
  
  return price
})

const totalPrice = computed(() => basePrice.value + configPrice.value)

// 获取产品数据
const fetchProduct = async () => {
  loading.value = true
  try {
    const { data } = await request.get(`/api/v2/products/${route.params.id}`)
    if (data?.data) {
      product.value = data.data
      regions.value = data.data.regions || []
      osTypes.value = data.data.os_types || []
      billingCycles.value = data.data.billing_cycles || []
      cpuOptions.value = data.data.cpu_options || []
      memoryOptions.value = data.data.memory_options || []
      diskOptions.value = data.data.disk_options || []
      bandwidthOptions.value = data.data.bandwidth_options || []
      
      // 默认选中第一个
      if (regions.value.length) selectedRegion.value = regions.value[0].id
      if (osTypes.value.length) {
        selectedOsType.value = osTypes.value[0].id
        if (osTypes.value[0].versions?.length) {
          selectedOs.value = osTypes.value[0].versions[0].id
        }
      }
      if (billingCycles.value.length) selectedCycle.value = billingCycles.value[0].id
      if (cpuOptions.value.length) selectedCpu.value = cpuOptions.value[0].id
      if (memoryOptions.value.length) selectedMemory.value = memoryOptions.value[0].id
      if (diskOptions.value.length) selectedDisk.value = diskOptions.value[0].id
      if (bandwidthOptions.value.length) selectedBandwidth.value = bandwidthOptions.value[0].id
    }
  } catch (error) {
    console.error('获取产品数据失败:', error)
    ElMessage.error('获取产品数据失败')
  } finally {
    loading.value = false
  }
}

const applyCoupon = async () => {
  if (!couponCode.value) return
  
  try {
    const { data } = await request.post('/api/v2/promo-codes/validate', {
      code: couponCode.value
    })
    if (data?.ok) {
      ElMessage.success('优惠码已应用')
    } else {
      ElMessage.error(data?.message || '优惠码无效')
    }
  } catch (error) {
    ElMessage.error('验证优惠码失败')
  }
}

const addToCart = async () => {
  try {
    const { data } = await request.post('/api/v2/cart/add', {
      product_id: product.value.id,
      region: selectedRegion.value,
      os: selectedOs.value,
      cpu: selectedCpu.value,
      memory: selectedMemory.value,
      disk: selectedDisk.value,
      bandwidth: selectedBandwidth.value,
      billing_cycle: selectedCycle.value,
      coupon_code: couponCode.value
    })
    if (data?.ok) {
      ElMessage.success('已添加到购物车')
      router.push('/cart')
    } else {
      ElMessage.error(data?.message || '添加失败')
    }
  } catch (error) {
    ElMessage.error('添加购物车失败')
  }
}

onMounted(() => {
  fetchProduct()
  fetchSiteName()
})
</script>

<style scoped lang="scss">
.product-detail-page {
  min-height: 100vh;
  background: #f5f7fa;
}

.page-header {
  background: #fff;
  border-bottom: 1px solid #e5e5ea;
  position: sticky;
  top: 0;
  z-index: 100;
  
  .header-inner {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  
  .logo {
    display: flex;
    align-items: center;
    gap: 10px;
    text-decoration: none;
    
    .logo-img {
      width: 32px;
      height: 32px;
    }
    
    .logo-text {
      font-size: 18px;
      font-weight: 600;
      color: #1d1d1f;
    }
  }
  
  .nav-links {
    display: flex;
    align-items: center;
    gap: 24px;
    
    .nav-link {
      color: #606266;
      text-decoration: none;
      font-size: 15px;
      
      &.active {
        color: #1a73e8;
        font-weight: 500;
      }
    }
  }
}

.detail-content {
  padding: 24px 0;
}

.detail-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.breadcrumb {
  margin-bottom: 24px;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 24px;
  
  @media (max-width: 992px) {
    grid-template-columns: 1fr;
  }
}

.config-section {
  .config-card {
    background: #fff;
    border-radius: 12px;
    padding: 32px;
    border: 1px solid #e5e5ea;
    
    h2 {
      font-size: 24px;
      font-weight: 600;
      margin: 0 0 8px;
    }
    
    .product-desc {
      color: #86868b;
      margin: 0 0 24px;
    }
    
    .features-list {
      display: flex;
      flex-wrap: wrap;
      gap: 16px;
      margin-bottom: 24px;
      
      .feature-item {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 14px;
        color: #606266;
      }
    }
  }
}

.config-group {
  margin-bottom: 24px;
  
  h3 {
    font-size: 16px;
    font-weight: 600;
    margin: 0 0 12px;
  }
}

.region-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  
  .region-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    border: 1px solid #e5e5ea;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    
    &:hover {
      border-color: #1a73e8;
    }
    
    &.active {
      border-color: #1a73e8;
      background: rgba(26, 115, 232, 0.05);
      color: #1a73e8;
    }
    
    .region-flag {
      width: 20px;
      height: 14px;
      object-fit: cover;
    }
  }
}

.os-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  
  .os-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    border: 1px solid #e5e5ea;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    
    &:hover {
      border-color: #1a73e8;
    }
    
    &.active {
      border-color: #1a73e8;
      background: rgba(26, 115, 232, 0.05);
      color: #1a73e8;
    }
    
    .os-icon {
      width: 20px;
      height: 20px;
    }
  }
}

.price-add {
  font-size: 12px;
  color: #f56c6c;
  margin-left: 4px;
}

.price-section {
  .price-card {
    background: #fff;
    border-radius: 12px;
    padding: 32px;
    border: 1px solid #e5e5ea;
    position: sticky;
    top: 88px;
    
    h3 {
      font-size: 18px;
      font-weight: 600;
      margin: 0 0 20px;
    }
    
    h4 {
      font-size: 14px;
      color: #909399;
      margin: 0 0 12px;
    }
  }
}

.billing-cycles {
  .cycle-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    border: 1px solid #e5e5ea;
    border-radius: 8px;
    margin-bottom: 8px;
    cursor: pointer;
    transition: all 0.2s;
    
    &:hover {
      border-color: #1a73e8;
    }
    
    &.active {
      border-color: #1a73e8;
      background: rgba(26, 115, 232, 0.05);
    }
    
    .cycle-name {
      font-weight: 500;
    }
    
    .cycle-price {
      color: #1a73e8;
      font-weight: 600;
    }
    
    .cycle-discount {
      padding: 2px 8px;
      background: #ff3b30;
      color: #fff;
      border-radius: 4px;
      font-size: 12px;
    }
  }
}

.price-summary {
  margin-bottom: 20px;
  
  .summary-item {
    display: flex;
    justify-content: space-between;
    padding: 8px 0;
    font-size: 14px;
    color: #606266;
    
    &.total {
      border-top: 1px solid #e5e5ea;
      padding-top: 16px;
      margin-top: 8px;
      font-weight: 600;
      color: #1d1d1f;
      
      .total-price {
        font-size: 24px;
        color: #1a73e8;
      }
    }
  }
}

.coupon-input {
  margin-bottom: 20px;
}

.buy-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
}

.price-tips {
  margin-top: 20px;
  
  p {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    color: #909399;
    margin: 8px 0;
  }
}
</style>
