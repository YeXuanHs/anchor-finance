<template>
  <div class="products-page">
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
          <router-link to="/products" class="nav-link active">产品</router-link>
          <router-link to="/" class="nav-link">公告</router-link>
          <router-link to="/" class="nav-link">帮助</router-link>
        </nav>
        <div class="header-actions">
          <el-button text @click="$router.push('/login')">登录</el-button>
          <el-button class="btn-gradient" round size="small" @click="$router.push('/register')">免费注册</el-button>
        </div>
      </div>
    </header>

    <!-- Filter Section -->
    <div class="filter-section">
      <div class="filter-inner">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item>产品中心</el-breadcrumb-item>
        </el-breadcrumb>

        <div class="filter-bar">
          <div class="filter-group">
            <span class="filter-label">分类</span>
            <el-select v-model="selectedGroup" placeholder="全部分类" clearable style="width: 160px">
              <el-option v-for="g in productGroups" :key="g.key" :label="g.label" :value="g.key" />
            </el-select>
          </div>

          <div class="filter-group">
            <span class="filter-label">价格</span>
            <el-select v-model="priceRange" placeholder="价格区间" clearable style="width: 160px">
              <el-option label="不限" value="" />
              <el-option label="0 - 50元" value="0-50" />
              <el-option label="50 - 200元" value="50-200" />
              <el-option label="200 - 500元" value="200-500" />
              <el-option label="500元以上" value="500+" />
            </el-select>
          </div>

          <div class="filter-group">
            <span class="filter-label">排序</span>
            <el-select v-model="sortBy" placeholder="排序方式" style="width: 160px">
              <el-option label="默认排序" value="default" />
              <el-option label="价格升序" value="price-asc" />
              <el-option label="价格降序" value="price-desc" />
              <el-option label="最新上架" value="newest" />
            </el-select>
          </div>

          <el-button text type="primary" @click="resetFilters" class="reset-btn">
            <el-icon style="margin-right: 4px"><RefreshRight /></el-icon>
            重置
          </el-button>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="main-content">
      <div class="content-inner">
        <!-- Skeleton Loading -->
        <template v-if="loading">
          <el-row :gutter="20">
            <el-col v-for="i in 9" :key="i" :xs="24" :sm="12" :lg="8">
              <div class="skeleton-card">
                <el-skeleton animated>
                  <template #template>
                    <div style="padding: 24px">
                      <div style="display: flex; align-items: center; gap: 14px; margin-bottom: 16px">
                        <el-skeleton-item variant="rect" style="width: 52px; height: 52px; border-radius: 12px" />
                        <div style="flex: 1">
                          <el-skeleton-item variant="h3" style="width: 60%; height: 20px; margin-bottom: 8px" />
                          <el-skeleton-item variant="text" style="width: 100%; height: 14px" />
                        </div>
                      </div>
                      <div style="display: flex; gap: 6px; margin-bottom: 16px">
                        <el-skeleton-item variant="rect" style="width: 64px; height: 24px; border-radius: 4px" />
                        <el-skeleton-item variant="rect" style="width: 64px; height: 24px; border-radius: 4px" />
                        <el-skeleton-item variant="rect" style="width: 64px; height: 24px; border-radius: 4px" />
                      </div>
                      <el-skeleton-item variant="rect" style="width: 100%; height: 1px; margin-bottom: 16px" />
                      <div style="display: flex; justify-content: space-between; align-items: center">
                        <el-skeleton-item variant="h3" style="width: 120px; height: 28px" />
                        <el-skeleton-item variant="button" style="width: 88px; height: 34px; border-radius: 6px" />
                      </div>
                    </div>
                  </template>
                </el-skeleton>
              </div>
            </el-col>
          </el-row>
        </template>

        <!-- Products Grid -->
        <template v-else>
          <el-row :gutter="20">
            <el-col
              v-for="product in paginatedProducts"
              :key="product.id"
              :xs="24"
              :sm="12"
              :lg="8"
            >
              <div class="product-card" @click="$router.push(`/products/${product.id}`)">
                <div class="card-badge" v-if="product.badge">{{ product.badge }}</div>
                <div class="card-header">
                  <div class="card-icon-wrap">
                    <el-icon :size="26" color="#0056FF"><component :is="product.icon" /></el-icon>
                  </div>
                  <div class="card-header-text">
                    <h3 class="card-title">{{ product.name }}</h3>
                    <p class="card-desc">{{ product.description }}</p>
                  </div>
                </div>
                <div class="card-specs">
                  <span class="spec-item" v-for="(tag, idx) in product.tags" :key="idx">
                    <el-icon :size="12" color="#0056FF"><Check /></el-icon>
                    {{ tag }}
                  </span>
                </div>
                <div class="card-footer">
                  <div class="card-price">
                    <span class="price-symbol">¥</span>
                    <span class="price-amount">{{ product.price }}</span>
                    <span class="price-unit">/月起</span>
                  </div>
                  <el-button class="btn-gradient btn-sm" round @click.stop="$router.push(`/products/${product.id}`)">
                    立即购买
                  </el-button>
                </div>
              </div>
            </el-col>
          </el-row>

          <!-- Empty State -->
          <div v-if="filteredProducts.length === 0" class="empty-state">
            <el-empty description="暂无匹配的产品">
              <el-button class="btn-gradient" round @click="resetFilters">重置筛选</el-button>
            </el-empty>
          </div>

          <!-- Pagination -->
          <div v-if="totalPages > 1" class="pagination-wrap">
            <el-pagination
              v-model:current-page="currentPage"
              :page-count="totalPages"
              :page-size="pageSize"
              :total="filteredProducts.length"
              layout="prev, pager, next, jumper"
              background
            />
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  Monitor,
  Cpu,
  Monitor as DesktopIcon,
  Connection,
  Position,
  Setting,
  Check,
  RefreshRight
} from '@element-plus/icons-vue'

const loading = ref(true)
const selectedGroup = ref('')
const priceRange = ref('')
const sortBy = ref('default')
const currentPage = ref(1)
const pageSize = 9

interface ProductGroup {
  key: string
  label: string
}

const productGroups: ProductGroup[] = [
  { key: 'ecs', label: '云服务器' },
  { key: 'vps', label: 'VPS' },
  { key: 'dedicated', label: '独立服务器' },
  { key: 'domain', label: '域名' },
  { key: 'ssl', label: 'SSL证书' }
]

interface Product {
  id: number
  name: string
  description: string
  price: number
  group: string
  icon: any
  tags: string[]
  badge?: string
}

const products = ref<Product[]>([
  {
    id: 1,
    name: '香港云服务器',
    description: 'BGP国际多线，延迟低至10ms，适合外贸、游戏加速等业务场景',
    price: 49,
    group: 'ecs',
    icon: DesktopIcon,
    tags: ['BGP多线', 'SSD存储', '弹性扩展'],
    badge: '热卖'
  },
  {
    id: 2,
    name: '美国云服务器',
    description: '接入CN2 GIA优质线路，大带宽不限流量，适合全球化业务部署',
    price: 69,
    group: 'ecs',
    icon: DesktopIcon,
    tags: ['CN2 GIA', '大带宽', '不限流量']
  },
  {
    id: 3,
    name: '日本云服务器',
    description: '东京机房，NTT线路直连中国大陆，低延迟高稳定性',
    price: 59,
    group: 'ecs',
    icon: DesktopIcon,
    tags: ['NTT线路', '低延迟', '东京机房'],
    badge: '新品'
  },
  {
    id: 4,
    name: '香港 VPS',
    description: 'KVM虚拟化架构，独立IP，root权限，适合个人建站和开发测试',
    price: 19,
    group: 'vps',
    icon: Cpu,
    tags: ['KVM架构', '独立IP', 'root权限']
  },
  {
    id: 5,
    name: '美国 VPS',
    description: '高性能NVMe SSD存储，1Gbps大带宽，适合中小型项目部署',
    price: 29,
    group: 'vps',
    icon: Cpu,
    tags: ['NVMe SSD', '1Gbps带宽', '高性价比']
  },
  {
    id: 6,
    name: '新加坡 VPS',
    description: '东南亚优质节点，覆盖亚太地区用户，适合跨境电商和出海业务',
    price: 35,
    group: 'vps',
    icon: Cpu,
    tags: ['东南亚节点', '跨境电商', '出海首选']
  },
  {
    id: 7,
    name: '香港独立服务器',
    description: 'E5高性能CPU，64G内存起步，独享带宽，适合大型企业和高并发业务',
    price: 599,
    group: 'dedicated',
    icon: Setting,
    tags: ['独享带宽', '高性能CPU', '大内存'],
    badge: '推荐'
  },
  {
    id: 8,
    name: '美国独立服务器',
    description: '双路E5处理器，128G内存，100Mbps独享带宽，适合数据密集型业务',
    price: 799,
    group: 'dedicated',
    icon: Setting,
    tags: ['双路E5', '100Mbps独享', '大存储']
  },
  {
    id: 9,
    name: '域名注册',
    description: '支持.com/.cn/.net等主流后缀，首年注册优惠，DNS解析免费',
    price: 9,
    group: 'domain',
    icon: Connection,
    tags: ['.com/.cn/.net', '首年优惠', '免费DNS']
  },
  {
    id: 10,
    name: '域名转入',
    description: '域名转入即享一年续费，免费提供Whois隐私保护和DNS解析服务',
    price: 55,
    group: 'domain',
    icon: Connection,
    tags: ['免费续费1年', '隐私保护', '免费DNS']
  },
  {
    id: 11,
    name: 'DV SSL证书',
    description: '域名验证型证书，10分钟快速签发，支持单域名和通配符，适合个人站点',
    price: 0,
    group: 'ssl',
    icon: Position,
    tags: ['免费签发', '快速验证', '单域名'],
    badge: '免费'
  },
  {
    id: 12,
    name: 'OV SSL证书',
    description: '企业验证型证书，地址栏安全锁标识，提升用户信任度，适合企业官网',
    price: 199,
    group: 'ssl',
    icon: Position,
    tags: ['企业验证', '安全锁', '信任标识']
  },
  {
    id: 13,
    name: 'EV SSL证书',
    description: '扩展验证型证书，绿色地址栏显示企业名称，最高级别安全认证',
    price: 599,
    group: 'ssl',
    icon: Position,
    tags: ['扩展验证', '绿色地址栏', '最高安全']
  }
])

const filteredProducts = computed(() => {
  let list = [...products.value]

  if (selectedGroup.value) {
    list = list.filter(p => p.group === selectedGroup.value)
  }

  if (priceRange.value) {
    if (priceRange.value === '500+') {
      list = list.filter(p => p.price >= 500)
    } else {
      const [min, max] = priceRange.value.split('-').map(Number)
      list = list.filter(p => p.price >= min && p.price <= max)
    }
  }

  if (sortBy.value === 'price-asc') {
    list.sort((a, b) => a.price - b.price)
  } else if (sortBy.value === 'price-desc') {
    list.sort((a, b) => b.price - a.price)
  } else if (sortBy.value === 'newest') {
    list.sort((a, b) => b.id - a.id)
  }

  return list
})

const totalPages = computed(() => Math.ceil(filteredProducts.value.length / pageSize))

const paginatedProducts = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredProducts.value.slice(start, start + pageSize)
})

function resetFilters() {
  selectedGroup.value = ''
  priceRange.value = ''
  sortBy.value = 'default'
  currentPage.value = 1
}

onMounted(() => {
  setTimeout(() => {
    loading.value = false
  }, 800)
})
</script>

<style scoped>
.products-page {
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

.btn-sm {
  height: 32px;
  font-size: 13px;
  padding: 0 16px;
}

/* Filter Section */
.filter-section {
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  margin-top: 60px;
  padding: 20px 0;
}

.filter-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-top: 16px;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-label {
  font-size: 13px;
  color: #999;
  font-weight: 500;
  white-space: nowrap;
}

.reset-btn {
  margin-left: auto;
}

/* Main Content */
.main-content {
  padding: 24px 0 40px;
}

.content-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
}

/* Skeleton Card */
.skeleton-card {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  margin-bottom: 20px;
}

/* Product Card */
.product-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  margin-bottom: 20px;
}

.product-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 86, 255, 0.15);
}

.card-badge {
  position: absolute;
  top: 12px;
  right: -28px;
  background: linear-gradient(135deg, #0056FF 0%, #4080FF 100%);
  color: #fff;
  font-size: 11px;
  font-weight: 600;
  padding: 2px 32px;
  transform: rotate(45deg);
}

.card-header {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  margin-bottom: 16px;
}

.card-icon-wrap {
  width: 52px;
  height: 52px;
  background: linear-gradient(135deg, #EBF3FD 0%, #d6e6ff 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.card-header-text {
  flex: 1;
  min-width: 0;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a3a5c;
  margin-bottom: 6px;
}

.card-desc {
  font-size: 13px;
  color: #999;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-specs {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.spec-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #666;
  background: #EBF3FD;
  padding: 3px 10px;
  border-radius: 4px;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-price {
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.price-symbol {
  font-size: 14px;
  font-weight: 600;
  color: #0056FF;
}

.price-amount {
  font-size: 26px;
  font-weight: 700;
  color: #0056FF;
  line-height: 1;
}

.price-unit {
  font-size: 12px;
  color: #999;
  margin-left: 2px;
}

/* Empty State */
.empty-state {
  padding: 60px 0;
}

/* Pagination */
.pagination-wrap {
  display: flex;
  justify-content: center;
  margin-top: 32px;
  padding: 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

:deep(.el-pagination.is-background .el-pager li:not(.is-disabled).is-active) {
  background: linear-gradient(135deg, #0056FF 0%, #4080FF 100%) !important;
}

/* Responsive */
@media (max-width: 768px) {
  .nav-links {
    display: none;
  }
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }
  .filter-group {
    width: 100%;
  }
  .filter-group .el-select {
    flex: 1;
  }
  .reset-btn {
    margin-left: 0;
  }
}
</style>
