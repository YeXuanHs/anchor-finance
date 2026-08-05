<template>
  <div class="products-page">
    <!-- Header -->
    <header class="page-header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <img src="/logo.png" :alt="siteName" class="logo-img" />
          <span class="logo-text">{{ siteName }}</span>
        </router-link>
        <nav class="nav-links">
          <router-link to="/" class="nav-link">首页</router-link>
          <el-dropdown trigger="hover" @command="(cmd: string) => $router.push(`/products?group=${cmd}`)">
            <span class="nav-link active">
              产品<el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item v-for="group in productGroups" :key="group.id" :command="group.id">
                  {{ group.name }}
                </el-dropdown-item>
                <el-dropdown-item divided command="">全部产品</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <router-link to="/" class="nav-link">公告</router-link>
          <router-link to="/" class="nav-link">帮助</router-link>
        </nav>
        <div class="header-actions">
          <el-button text @click="$router.push('/login')">登录</el-button>
          <el-button type="primary" round size="small" @click="$router.push('/register')">免费注册</el-button>
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
              <el-option v-for="g in productGroups" :key="g.id" :label="g.name" :value="g.id" />
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
            <el-select v-model="sortBy" placeholder="默认排序" clearable style="width: 160px">
              <el-option label="默认排序" value="" />
              <el-option label="价格从低到高" value="price_asc" />
              <el-option label="价格从高到低" value="price_desc" />
              <el-option label="销量优先" value="sales_desc" />
            </el-select>
          </div>

          <div class="filter-group">
            <el-input
              v-model="keyword"
              placeholder="搜索产品..."
              clearable
              style="width: 200px"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </div>
      </div>
    </div>

    <!-- Products Grid -->
    <div class="products-content">
      <div class="products-grid" v-loading="loading">
        <div
          v-for="product in filteredProducts"
          :key="product.id"
          class="product-card"
          @click="$router.push(`/products/${product.id}`)"
        >
          <div class="product-badge" v-if="product.badge">{{ product.badge }}</div>
          <div class="product-header">
            <div class="product-icon" :style="{ background: product.gradient || 'linear-gradient(135deg, #1a73e8, #4a90e2)' }">
              <el-icon :size="28" color="#fff"><Monitor /></el-icon>
            </div>
            <h3 class="product-name">{{ product.name }}</h3>
            <p class="product-desc">{{ product.description }}</p>
          </div>
          <div class="product-specs">
            <div class="spec-item" v-for="spec in product.specs" :key="spec.label">
              <span class="spec-value">{{ spec.value }}</span>
              <span class="spec-label">{{ spec.label }}</span>
            </div>
          </div>
          <div class="product-footer">
            <div class="product-price">
              <span class="price-symbol">¥</span>
              <span class="price-value">{{ product.price }}</span>
              <span class="price-unit">/月</span>
            </div>
            <el-button type="primary" round size="small">立即选购</el-button>
          </div>
        </div>
      </div>

      <div class="empty-state" v-if="!loading && filteredProducts.length === 0">
        <el-icon :size="64" color="#e5e5ea"><Box /></el-icon>
        <p>暂无符合条件的产品</p>
        <el-button @click="resetFilters">重置筛选</el-button>
      </div>
    </div>

    <!-- Footer -->
    <footer class="page-footer">
      <div class="footer-inner">
        <p>&copy; {{ new Date().getFullYear() }} {{ siteName }} All Rights Reserved</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowDown, Search, Monitor, Box } from '@element-plus/icons-vue'
import request from '@/utils/request'

const route = useRoute()
const loading = ref(false)

const productGroups = ref([])
const products = ref([])
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

const selectedGroup = ref(route.query.group as string || '')
const priceRange = ref('')
const sortBy = ref('')
const keyword = ref('')

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    // 获取产品分组
    const groupRes = await request.get('/api/v1/product-groups')
    if (groupRes.data?.data) {
      productGroups.value = groupRes.data.data
    }
    
    // 获取产品列表
    const productRes = await request.get('/api/v1/products', {
      params: {
        group: selectedGroup.value,
        sort: sortBy.value,
        keyword: keyword.value
      }
    })
    if (productRes.data?.data) {
      products.value = productRes.data.data
    }
  } catch (error) {
    console.error('获取产品数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 筛选后的产品
const filteredProducts = computed(() => {
  let result = [...products.value]
  
  // 分组筛选
  if (selectedGroup.value) {
    result = result.filter(p => p.group_id === selectedGroup.value)
  }
  
  // 价格筛选
  if (priceRange.value) {
    const [min, max] = priceRange.value.split('-').map(Number)
    result = result.filter(p => {
      if (max) {
        return p.price >= min && p.price <= max
      }
      return p.price >= min
    })
  }
  
  // 关键词搜索
  if (keyword.value) {
    const kw = keyword.value.toLowerCase()
    result = result.filter(p => 
      p.name.toLowerCase().includes(kw) || 
      p.description?.toLowerCase().includes(kw)
    )
  }
  
  // 排序
  if (sortBy.value === 'price_asc') {
    result.sort((a, b) => a.price - b.price)
  } else if (sortBy.value === 'price_desc') {
    result.sort((a, b) => b.price - a.price)
  } else if (sortBy.value === 'sales_desc') {
    result.sort((a, b) => (b.sales || 0) - (a.sales || 0))
  }
  
  return result
})

const resetFilters = () => {
  selectedGroup.value = ''
  priceRange.value = ''
  sortBy.value = ''
  keyword.value = ''
}

onMounted(() => {
  fetchData()
  fetchSiteName()
})
</script>

<style scoped lang="scss">
.products-page {
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
      cursor: pointer;
      
      &.active {
        color: #1a73e8;
        font-weight: 500;
      }
      
      &:hover {
        color: #1a73e8;
      }
    }
  }
}

.filter-section {
  background: #fff;
  border-bottom: 1px solid #e5e5ea;
  padding: 16px 0;
  
  .filter-inner {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
  }
  
  .filter-bar {
    display: flex;
    align-items: center;
    gap: 20px;
    margin-top: 16px;
    flex-wrap: wrap;
    
    .filter-group {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .filter-label {
        font-size: 14px;
        color: #909399;
        white-space: nowrap;
      }
    }
  }
}

.products-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px 20px;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  
  @media (max-width: 992px) {
    grid-template-columns: repeat(2, 1fr);
  }
  
  @media (max-width: 576px) {
    grid-template-columns: 1fr;
  }
}

.product-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
  border: 1px solid #e5e5ea;
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.08);
    border-color: #1a73e8;
  }
  
  .product-badge {
    position: absolute;
    top: 12px;
    right: 12px;
    padding: 4px 10px;
    background: #ff3b30;
    color: #fff;
    border-radius: 8px;
    font-size: 12px;
  }
  
  .product-header {
    margin-bottom: 20px;
    
    .product-icon {
      width: 56px;
      height: 56px;
      border-radius: 14px;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-bottom: 16px;
    }
    
    .product-name {
      font-size: 18px;
      font-weight: 600;
      margin: 0 0 8px;
      color: #1d1d1f;
    }
    
    .product-desc {
      font-size: 14px;
      color: #86868b;
      margin: 0;
      line-height: 1.5;
    }
  }
  
  .product-specs {
    display: flex;
    gap: 16px;
    padding: 16px 0;
    border-top: 1px solid #f0f0f0;
    border-bottom: 1px solid #f0f0f0;
    margin-bottom: 16px;
    
    .spec-item {
      flex: 1;
      text-align: center;
      
      .spec-value {
        display: block;
        font-size: 14px;
        font-weight: 600;
        color: #1d1d1f;
      }
      
      .spec-label {
        font-size: 12px;
        color: #909399;
      }
    }
  }
  
  .product-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    
    .product-price {
      .price-symbol {
        font-size: 14px;
        color: #1a73e8;
      }
      
      .price-value {
        font-size: 28px;
        font-weight: 700;
        color: #1a73e8;
      }
      
      .price-unit {
        font-size: 14px;
        color: #909399;
      }
    }
  }
}

.empty-state {
  text-align: center;
  padding: 80px 0;
  
  p {
    margin: 16px 0;
    color: #909399;
  }
}

.page-footer {
  background: #1d1d1f;
  color: rgba(255, 255, 255, 0.6);
  padding: 24px 0;
  
  .footer-inner {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
    text-align: center;
    
    p {
      margin: 0;
      font-size: 14px;
    }
  }
}
</style>
