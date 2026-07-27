<template>
  <div class="products-page">
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
          <router-link to="/products" class="nav-link active">产品</router-link>
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
          <n-breadcrumb-item>产品中心</n-breadcrumb-item>
        </n-breadcrumb>
      </div>
    </div>

    <!-- Main Content -->
    <div class="main-content">
      <div class="content-inner">
        <!-- Sidebar -->
        <aside class="sidebar">
          <div class="sidebar-card">
            <h3 class="sidebar-title">
              <n-icon size="18" color="#1890ff"><ListOutline /></n-icon>
              产品分类
            </h3>
            <div class="group-list">
              <div
                v-for="group in productGroups"
                :key="group.key"
                class="group-item"
                :class="{ active: selectedGroup === group.key }"
                @click="selectedGroup = group.key"
              >
                <n-icon size="18" :color="selectedGroup === group.key ? '#1890ff' : '#86909c'">
                  <component :is="group.icon" />
                </n-icon>
                <span class="group-label">{{ group.label }}</span>
                <span class="group-count">{{ group.count }}</span>
              </div>
            </div>
          </div>
        </aside>

        <!-- Product List -->
        <main class="product-list">
          <!-- Toolbar -->
          <div class="toolbar">
            <div class="toolbar-left">
              <span class="result-count">共 <strong>{{ filteredProducts.length }}</strong> 个产品</span>
            </div>
            <div class="toolbar-right">
              <n-select
                v-model:value="sortBy"
                :options="sortOptions"
                size="small"
                style="width: 140px;"
                placeholder="排序方式"
              />
              <n-button-group size="small">
                <n-button :type="viewMode === 'grid' ? 'primary' : 'default'" @click="viewMode = 'grid'">
                  <template #icon><n-icon><GridOutline /></n-icon></template>
                </n-button>
                <n-button :type="viewMode === 'list' ? 'primary' : 'default'" @click="viewMode = 'list'">
                  <template #icon><n-icon><ListOutline /></n-icon></template>
                </n-button>
              </n-button-group>
            </div>
          </div>

          <!-- Grid View -->
          <div v-if="viewMode === 'grid'" class="products-grid">
            <div
              v-for="product in paginatedProducts"
              :key="product.id"
              class="product-card"
              @click="$router.push(`/products/${product.id}`)"
            >
              <div class="card-badge" v-if="product.badge">{{ product.badge }}</div>
              <div class="card-icon-wrap">
                <n-icon size="32" color="#1890ff">
                  <component :is="product.icon" />
                </n-icon>
              </div>
              <h3 class="card-title">{{ product.name }}</h3>
              <p class="card-desc">{{ product.description }}</p>
              <div class="card-tags">
                <n-tag v-for="tag in product.tags" :key="tag" size="small" :bordered="false" type="info">
                  {{ tag }}
                </n-tag>
              </div>
              <div class="card-footer">
                <div class="card-price">
                  <span class="currency">¥</span>
                  <span class="amount">{{ product.price }}</span>
                  <span class="unit">.00/月起</span>
                </div>
                <n-button type="primary" size="small" round @click.stop="$router.push(`/products/${product.id}`)">
                  立即购买
                </n-button>
              </div>
            </div>
          </div>

          <!-- List View -->
          <div v-else class="products-list">
            <div
              v-for="product in paginatedProducts"
              :key="product.id"
              class="product-list-item"
              @click="$router.push(`/products/${product.id}`)"
            >
              <div class="list-icon-wrap">
                <n-icon size="28" color="#1890ff">
                  <component :is="product.icon" />
                </n-icon>
              </div>
              <div class="list-info">
                <div class="list-header">
                  <h3 class="list-title">{{ product.name }}</h3>
                  <n-tag v-if="product.badge" size="small" type="warning" :bordered="false">{{ product.badge }}</n-tag>
                </div>
                <p class="list-desc">{{ product.description }}</p>
                <div class="list-tags">
                  <n-tag v-for="tag in product.tags" :key="tag" size="small" :bordered="false" type="info">
                    {{ tag }}
                  </n-tag>
                </div>
              </div>
              <div class="list-price-area">
                <div class="list-price">
                  <span class="currency">¥</span>
                  <span class="amount">{{ product.price }}</span>
                  <span class="unit">.00/月起</span>
                </div>
                <n-button type="primary" size="small" round @click.stop="$router.push(`/products/${product.id}`)">
                  立即购买
                </n-button>
              </div>
            </div>
          </div>

          <!-- Empty State -->
          <div v-if="filteredProducts.length === 0" class="empty-state">
            <n-icon size="64" color="#c9cdd4"><CloudOfflineOutline /></n-icon>
            <p>暂无匹配的产品</p>
            <n-button type="primary" @click="resetFilters">重置筛选</n-button>
          </div>

          <!-- Pagination -->
          <div v-if="totalPages > 1" class="pagination-wrap">
            <n-pagination
              v-model:page="currentPage"
              :page-count="totalPages"
              :page-slot="7"
              show-quick-jumper
            />
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  AnchorOutline,
  ServerOutline,
  DesktopOutline,
  HardwareChipOutline,
  GlobeOutline,
  ShieldCheckmarkOutline,
  GridOutline,
  ListOutline,
  CloudOfflineOutline
} from '@vicons/ionicons5'

const viewMode = ref<'grid' | 'list'>('grid')
const selectedGroup = ref('all')
const sortBy = ref('default')
const currentPage = ref(1)
const pageSize = 9

const sortOptions = [
  { label: '默认排序', value: 'default' },
  { label: '价格升序', value: 'price-asc' },
  { label: '价格降序', value: 'price-desc' },
  { label: '最新上架', value: 'newest' }
]

interface ProductGroup {
  key: string
  label: string
  icon: any
  count: number
}

const productGroups = computed<ProductGroup[]>(() => [
  { key: 'all', label: '全部产品', icon: GridOutline, count: products.value.length },
  { key: 'ecs', label: '云服务器', icon: ServerOutline, count: products.value.filter(p => p.group === 'ecs').length },
  { key: 'vps', label: 'VPS', icon: HardwareChipOutline, count: products.value.filter(p => p.group === 'vps').length },
  { key: 'dedicated', label: '独立服务器', icon: DesktopOutline, count: products.value.filter(p => p.group === 'dedicated').length },
  { key: 'domain', label: '域名', icon: GlobeOutline, count: products.value.filter(p => p.group === 'domain').length },
  { key: 'ssl', label: 'SSL证书', icon: ShieldCheckmarkOutline, count: products.value.filter(p => p.group === 'ssl').length }
])

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
    icon: ServerOutline,
    tags: ['BGP多线', 'SSD存储', '弹性扩展'],
    badge: '热卖'
  },
  {
    id: 2,
    name: '美国云服务器',
    description: '接入CN2 GIA优质线路，大带宽不限流量，适合全球化业务部署',
    price: 69,
    group: 'ecs',
    icon: ServerOutline,
    tags: ['CN2 GIA', '大带宽', '不限流量']
  },
  {
    id: 3,
    name: '日本云服务器',
    description: '东京机房，NTT线路直连中国大陆，低延迟高稳定性',
    price: 59,
    group: 'ecs',
    icon: ServerOutline,
    tags: ['NTT线路', '低延迟', '东京机房'],
    badge: '新品'
  },
  {
    id: 4,
    name: '香港 VPS',
    description: 'KVM虚拟化架构，独立IP，root权限，适合个人建站和开发测试',
    price: 19,
    group: 'vps',
    icon: HardwareChipOutline,
    tags: ['KVM架构', '独立IP', 'root权限']
  },
  {
    id: 5,
    name: '美国 VPS',
    description: '高性能NVMe SSD存储，1Gbps大带宽，适合中小型项目部署',
    price: 29,
    group: 'vps',
    icon: HardwareChipOutline,
    tags: ['NVMe SSD', '1Gbps带宽', '高性价比']
  },
  {
    id: 6,
    name: '新加坡 VPS',
    description: '东南亚优质节点，覆盖亚太地区用户，适合跨境电商和出海业务',
    price: 35,
    group: 'vps',
    icon: HardwareChipOutline,
    tags: ['东南亚节点', '跨境电商', '出海首选']
  },
  {
    id: 7,
    name: '香港独立服务器',
    description: 'E5高性能CPU，64G内存起步，独享带宽，适合大型企业和高并发业务',
    price: 599,
    group: 'dedicated',
    icon: DesktopOutline,
    tags: ['独享带宽', '高性能CPU', '大内存'],
    badge: '推荐'
  },
  {
    id: 8,
    name: '美国独立服务器',
    description: '双路E5处理器，128G内存，100Mbps独享带宽，适合数据密集型业务',
    price: 799,
    group: 'dedicated',
    icon: DesktopOutline,
    tags: ['双路E5', '100Mbps独享', '大存储']
  },
  {
    id: 9,
    name: '域名注册',
    description: '支持.com/.cn/.net等主流后缀，首年注册优惠，DNS解析免费',
    price: 9,
    group: 'domain',
    icon: GlobeOutline,
    tags: ['.com/.cn/.net', '首年优惠', '免费DNS']
  },
  {
    id: 10,
    name: '域名转入',
    description: '域名转入即享一年续费，免费提供Whois隐私保护和DNS解析服务',
    price: 55,
    group: 'domain',
    icon: GlobeOutline,
    tags: ['免费续费1年', '隐私保护', '免费DNS']
  },
  {
    id: 11,
    name: 'DV SSL证书',
    description: '域名验证型证书，10分钟快速签发，支持单域名和通配符，适合个人站点',
    price: 0,
    group: 'ssl',
    icon: ShieldCheckmarkOutline,
    tags: ['免费签发', '快速验证', '单域名'],
    badge: '免费'
  },
  {
    id: 12,
    name: 'OV SSL证书',
    description: '企业验证型证书，地址栏安全锁标识，提升用户信任度，适合企业官网',
    price: 199,
    group: 'ssl',
    icon: ShieldCheckmarkOutline,
    tags: ['企业验证', '安全锁', '信任标识']
  },
  {
    id: 13,
    name: 'EV SSL证书',
    description: '扩展验证型证书，绿色地址栏显示企业名称，最高级别安全认证',
    price: 599,
    group: 'ssl',
    icon: ShieldCheckmarkOutline,
    tags: ['扩展验证', '绿色地址栏', '最高安全']
  }
])

const filteredProducts = computed(() => {
  let list = [...products.value]

  if (selectedGroup.value !== 'all') {
    list = list.filter(p => p.group === selectedGroup.value)
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
  selectedGroup.value = 'all'
  sortBy.value = 'default'
  currentPage.value = 1
}
</script>

<style scoped>
.products-page {
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

/* Sidebar */
.sidebar {
  width: 240px;
  flex-shrink: 0;
}

.sidebar-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.sidebar-title {
  font-size: 15px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.group-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.group-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
  color: #4e5969;
}

.group-item:hover {
  background: #f2f3f5;
  color: #1d2129;
}

.group-item.active {
  background: #e6f7ff;
  color: #1890ff;
}

.group-label {
  flex: 1;
  font-weight: 500;
}

.group-count {
  font-size: 12px;
  color: #c9cdd4;
  background: #f2f3f5;
  padding: 1px 8px;
  border-radius: 10px;
}

.group-item.active .group-count {
  background: rgba(24, 144, 255, 0.1);
  color: #1890ff;
}

/* Toolbar */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.result-count {
  font-size: 14px;
  color: #86909c;
}

.result-count strong {
  color: #1890ff;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* Grid View */
.products-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.product-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid #f0f1f5;
  position: relative;
  overflow: hidden;
}

.product-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 28px rgba(24, 144, 255, 0.1);
  border-color: rgba(24, 144, 255, 0.2);
}

.card-badge {
  position: absolute;
  top: 12px;
  right: -28px;
  background: linear-gradient(135deg, #ff7a45, #fa541c);
  color: #fff;
  font-size: 12px;
  padding: 2px 32px;
  transform: rotate(45deg);
}

.card-icon-wrap {
  width: 56px;
  height: 56px;
  background: linear-gradient(135deg, #e6f7ff, #bae7ff);
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}

.card-title {
  font-size: 17px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 8px;
}

.card-desc {
  font-size: 13px;
  color: #86909c;
  line-height: 1.6;
  margin-bottom: 14px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  min-height: 42px;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 18px;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid #f2f3f5;
}

.card-price,
.list-price {
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.card-price .currency,
.list-price .currency {
  font-size: 14px;
  font-weight: 600;
  color: #1890ff;
}

.card-price .amount,
.list-price .amount {
  font-size: 26px;
  font-weight: 700;
  color: #1890ff;
  line-height: 1;
}

.card-price .unit,
.list-price .unit {
  font-size: 13px;
  color: #c9cdd4;
  margin-left: 2px;
}

/* List View */
.products-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.product-list-item {
  display: flex;
  align-items: center;
  gap: 20px;
  background: #fff;
  border-radius: 12px;
  padding: 20px 24px;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid #f0f1f5;
}

.product-list-item:hover {
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.08);
  border-color: rgba(24, 144, 255, 0.2);
}

.list-icon-wrap {
  width: 52px;
  height: 52px;
  background: linear-gradient(135deg, #e6f7ff, #bae7ff);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.list-info {
  flex: 1;
  min-width: 0;
}

.list-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.list-title {
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
}

.list-desc {
  font-size: 13px;
  color: #86909c;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.list-tags {
  display: flex;
  gap: 6px;
}

.list-price-area {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 10px;
  flex-shrink: 0;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 80px 0;
  color: #c9cdd4;
}

.empty-state p {
  margin: 16px 0 24px;
  font-size: 15px;
}

/* Pagination */
.pagination-wrap {
  display: flex;
  justify-content: center;
  margin-top: 32px;
  padding: 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

/* Responsive */
@media (max-width: 1200px) {
  .products-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 1024px) {
  .sidebar {
    display: none;
  }
}

@media (max-width: 768px) {
  .nav-links {
    display: none;
  }
  .products-grid {
    grid-template-columns: 1fr;
  }
}
</style>
