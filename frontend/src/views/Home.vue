<template>
  <div class="home">
    <!-- Fixed Header -->
    <header class="header" :class="{ 'header-scrolled': scrolled }">
      <div class="header-inner">
        <div class="logo" @click="$router.push('/')">
          <img :src="configStore.getLogo('home') || '/logo.png'" :alt="$t('landing.brandName')" class="logo-img" />
          <span class="logo-text">{{ $t('landing.brandName') }}</span>
        </div>
        <nav class="nav-links">
          <!-- 动态导航 - 从数据库读取 -->
          <template v-for="item in topNavs" :key="item.id">
            <!-- 有子菜单的 -->
            <el-dropdown v-if="item.children && item.children.length > 0" trigger="hover" @command="(cmd: string) => $router.push(cmd)">
              <span class="nav-link">
                {{ item.name }}<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-for="child in item.children" :key="child.id" :command="child.url">
                    {{ child.name }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <!-- 没有子菜单的 -->
            <router-link v-else :to="item.url" class="nav-link">{{ item.name }}</router-link>
          </template>
        </nav>
        <div class="header-actions">
          <LanguageSwitch />
          <el-button text class="login-btn" @click="$router.push('/login')">
            <el-icon :size="16" style="margin-right: 4px;"><User /></el-icon>
            {{ $t('landing.login') }}
          </el-button>
          <el-button type="primary" round size="default" @click="$router.push('/register')">
            {{ $t('landing.freeRegister') }}
          </el-button>
        </div>
      </div>
    </header>

    <!-- Hero Carousel - 带左右切换按钮 -->
    <section id="hero" class="hero-section">
      <div class="carousel-wrapper">
        <el-carousel
          ref="carouselRef"
          :interval="5000"
          :autoplay="true"
          indicator-position="none"
          height="640px"
          class="hero-carousel"
          motion-blur
        >
          <el-carousel-item v-for="banner in banners" :key="banner.id">
            <div class="carousel-slide" :style="{ background: banner.bg_color || 'linear-gradient(135deg, #1a237e, #0d47a1)' }">
              <video v-if="banner.video" class="slide-video" autoplay muted loop playsinline>
                <source :src="banner.video" type="video/webm" />
              </video>
              <div class="slide-overlay"></div>
              <div class="slide-content">
                <div class="slide-badge" v-if="banner.badge">
                  <el-icon :size="14"><Lightning /></el-icon>
                  {{ banner.badge }}
                </div>
                <h1 class="slide-title">{{ banner.title }}</h1>
                <p class="slide-desc">{{ banner.description }}</p>
                <div class="slide-actions">
                  <el-button type="primary" size="large" round @click="$router.push(banner.link || '/products')">
                    <el-icon style="margin-right: 6px;"><Position /></el-icon>
                    {{ banner.btn_text || $t('landing.buyNow') }}
                  </el-button>
                  <el-button size="large" round class="slide-ghost-btn" @click="$router.push('/products')">
                    {{ $t('landing.learnMore') }}
                  </el-button>
                </div>
              </div>
            </div>
          </el-carousel-item>
        </el-carousel>
        
        <!-- 左右切换按钮 -->
        <button class="carousel-btn carousel-btn-prev" @click="carouselRef?.prev()">
          <el-icon :size="24"><ArrowLeft /></el-icon>
        </button>
        <button class="carousel-btn carousel-btn-next" @click="carouselRef?.next()">
          <el-icon :size="24"><ArrowRight /></el-icon>
        </button>
        
        <!-- 底部指示器 -->
        <div class="carousel-indicators">
          <button
            v-for="(banner, index) in banners"
            :key="banner.id"
            class="indicator"
            :class="{ active: currentSlide === index }"
            @click="goToSlide(index)"
          ></button>
        </div>
      </div>
      
      <!-- 轮播图下方快捷入口 -->
      <div class="quick-entry">
        <div class="entry-item" v-for="entry in quickEntries" :key="entry.id" @click="$router.push(entry.link)">
          <el-icon :size="28"><component :is="entry.icon" /></el-icon>
          <span>{{ entry.title }}</span>
        </div>
      </div>
    </section>

    <!-- 热门产品推荐 - 销量前4 -->
    <section id="hot-products" class="section hot-products-section">
      <div class="container">
        <div class="section-header">
          <h2 class="section-title">{{ $t('landing.hotProducts') }}</h2>
          <p class="section-subtitle">{{ $t('landing.hotProductsDesc') }}</p>
        </div>
        <div class="products-grid">
          <div
            v-for="product in hotProducts"
            :key="product.id"
            class="product-card"
            @click="$router.push(`/products/${product.id}`)"
          >
            <div class="product-badge" v-if="product.badge">{{ product.badge }}</div>
            <div class="product-icon" :style="{ background: product.gradient || 'linear-gradient(135deg, #1a73e8, #4a90e2)' }">
              <el-icon :size="32"><Monitor /></el-icon>
            </div>
            <h3 class="product-name">{{ product.name }}</h3>
            <p class="product-desc">{{ product.description }}</p>
            <div class="product-price">
              <span class="price-symbol">¥</span>
              <span class="price-value">{{ product.price }}</span>
              <span class="price-unit">{{ $t('landing.priceFrom') }}</span>
            </div>
            <el-button type="primary" round class="product-btn">{{ $t('landing.buyNow') }}</el-button>
          </div>
        </div>
        <div class="section-more">
          <el-button @click="$router.push('/products')">{{ $t('landing.viewMoreProducts') }} <el-icon><ArrowRight /></el-icon></el-button>
        </div>
      </div>
    </section>

    <!-- 解决方案 -->
    <section id="solutions" class="section solutions-section">
      <div class="container">
        <div class="section-header">
          <h2 class="section-title">{{ $t('landing.solutions') }}</h2>
          <p class="section-subtitle">{{ $t('landing.solutionsDesc') }}</p>
        </div>
        <div class="solutions-grid">
          <div v-for="solution in solutions" :key="solution.id" class="solution-card">
            <div class="solution-icon">
              <el-icon :size="36"><component :is="solution.icon" /></el-icon>
            </div>
            <h3>{{ solution.title }}</h3>
            <p>{{ solution.description }}</p>
            <ul class="solution-features">
              <li v-for="feature in solution.features" :key="feature">{{ feature }}</li>
            </ul>
            <el-button link type="primary">{{ $t('landing.understandDetail') }} <el-icon><ArrowRight /></el-icon></el-button>
          </div>
        </div>
      </div>
    </section>

    <!-- 核心优势 -->
    <section id="features" class="section features-section">
      <div class="container">
        <div class="section-header">
          <h2 class="section-title">{{ $t('landing.whyChooseUs') }}</h2>
          <p class="section-subtitle">{{ $t('landing.whyChooseUsDesc') }}</p>
        </div>
        <div class="features-grid">
          <div v-for="feature in features" :key="feature.id" class="feature-card">
            <div class="feature-icon">
              <el-icon :size="32"><component :is="feature.icon" /></el-icon>
            </div>
            <h3>{{ feature.title }}</h3>
            <p>{{ feature.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 数据统计 -->
    <section class="section stats-section">
      <div class="container">
        <div class="stats-grid">
          <div v-for="stat in stats" :key="stat.id" class="stat-item">
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- 新闻公告 -->
    <section id="announcements" class="section announcements-section">
      <div class="container">
        <div class="section-header">
          <h2 class="section-title">{{ $t('landing.newsAnnouncements') }}</h2>
          <p class="section-subtitle">{{ $t('landing.newsDesc') }}</p>
        </div>
        <div class="news-grid">
          <div v-for="news in announcements" :key="news.id" class="news-card" @click="$router.push(`/news/${news.id}`)">
            <div class="news-date">
              <span class="day">{{ new Date(news.created_at).getDate() }}</span>
              <span class="month">{{ new Date(news.created_at).toLocaleDateString('zh', { month: 'short' }) }}</span>
            </div>
            <div class="news-content">
              <h3>{{ news.title }}</h3>
              <p>{{ news.summary }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 合作伙伴 -->
    <section id="partners" class="section partners-section">
      <div class="container">
        <div class="section-header">
          <h2 class="section-title">{{ $t('landing.partners') }}</h2>
          <p class="section-subtitle">{{ $t('landing.partnersDesc') }}</p>
        </div>
        <div class="partners-grid">
          <div v-for="partner in partners" :key="partner.id" class="partner-card">
            <img :src="partner.logo" :alt="partner.name" class="partner-logo" />
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer id="footer" class="footer">
      <div class="container">
        <div class="footer-grid">
          <div class="footer-col">
            <div class="footer-logo">
              <img :src="configStore.getLogo('home') || '/logo.png'" :alt="$t('landing.brandName')" />
              <span>{{ $t('landing.brandName') }}</span>
            </div>
            <p class="footer-desc">{{ siteSettings.site_description || $t('landing.defaultDesc') }}</p>
            <div class="footer-contact">
              <p v-if="siteSettings.contact_phone"><el-icon><Phone /></el-icon> {{ siteSettings.contact_phone }}</p>
              <p v-if="siteSettings.contact_email"><el-icon><Message /></el-icon> {{ siteSettings.contact_email }}</p>
              <p v-if="siteSettings.contact_address"><el-icon><Location /></el-icon> {{ siteSettings.contact_address }}</p>
            </div>
          </div>
          <div class="footer-col">
            <h4>{{ $t('landing.productServices') }}</h4>
            <ul>
              <li v-for="group in productGroups" :key="group.id">
                <router-link :to="`/products?group=${group.id}`">{{ group.name }}</router-link>
              </li>
            </ul>
          </div>
          <!-- 动态底部导航 -->
          <div v-for="nav in bottomNavs" :key="nav.id" class="footer-col">
            <h4>{{ nav.name }}</h4>
            <ul>
              <li v-for="child in nav.children" :key="child.id">
                <router-link :to="child.url">{{ child.name }}</router-link>
              </li>
            </ul>
          </div>
        </div>
        <div class="footer-bottom">
          <p>{{ siteSettings.copyright || `© ${new Date().getFullYear()} ${$t('landing.brandName')} ${$t('landing.allRightsReserved')}` }}</p>
          <p v-if="siteSettings.icp">{{ siteSettings.icp }}</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  ArrowDown, ArrowLeft, ArrowRight, User, Lightning, Position, Monitor,
  Phone, Message, Location
} from '@element-plus/icons-vue'
import request from '@/utils/request'
import { useConfigStore } from '@/stores/config'
import LanguageSwitch from '@/components/LanguageSwitch.vue'

const { t } = useI18n()
const carouselRef = ref()
const scrolled = ref(false)
const currentSlide = ref(0)

// 导航数据
const topNavs = ref<any[]>([])
const bottomNavs = ref<any[]>([])

// 从API获取的数据（默认数据作为后备，优先从 API 加载）
const banners = ref<any[]>([
  { id: 1, title: t('landing.highPerfCloud'), description: t('landing.highPerfDesc'), badge: t('landing.hotRecommend'), video: '/carousel/2.webm', btn_text: t('landing.buyNow'), link: '/products' },
  { id: 2, title: t('landing.globalNodes'), description: t('landing.globalNodesDesc'), badge: t('landing.globalLayout'), video: '/carousel/3.webm', btn_text: t('landing.buyNow'), link: '/products' },
  { id: 3, title: t('landing.proTechSupport'), description: t('landing.proTechSupportDesc'), badge: t('landing.proService'), video: '/carousel/4.webm', btn_text: t('landing.contactUs'), link: '/tickets/create' }
])

const productGroups = ref<any[]>([])
const hotProducts = ref<any[]>([])
const announcements = ref<any[]>([])
const partners = ref<any[]>([])
const siteSettings = ref<Record<string, any>>({})
const configStore = useConfigStore()

// 快捷入口
const quickEntries = [
  { id: 1, title: t('landing.cloudServer'), icon: 'Monitor', link: '/products?group=cloud' },
  { id: 2, title: t('landing.dedicatedServer'), icon: 'Cpu', link: '/products?group=dedicated' },
  { id: 3, title: t('landing.domainRegistration'), icon: 'Connection', link: '/products?group=domain' },
  { id: 4, title: t('landing.sslCertificate'), icon: 'Shield', link: '/products?group=ssl' }
]

// 解决方案
const solutions = [
  {
    id: 1,
    title: t('landing.enterpriseCloud'),
    icon: 'OfficeBuilding',
    description: t('landing.enterpriseCloudDesc'),
    features: [t('landing.featureElastic'), t('landing.featureHA'), t('landing.featureSecure'), t('landing.featureCost')]
  },
  {
    id: 2,
    title: t('landing.ecommerceSolution'),
    icon: 'ShoppingBag',
    description: t('landing.ecommerceSolutionDesc'),
    features: [t('landing.featureHighConcurrency'), t('landing.featureCDN'), t('landing.featureDataSecurity'), t('landing.featureStable')]
  },
  {
    id: 3,
    title: t('landing.gameAcceleration'),
    icon: 'TrendCharts',
    description: t('landing.gameAccelerationDesc'),
    features: [t('landing.featureGlobalNodes'), t('landing.featureDDoS'), t('landing.featureLowLatency'), t('landing.featureElasticScale')]
  },
  {
    id: 4,
    title: t('landing.bigDataAnalytics'),
    icon: 'DataLine',
    description: t('landing.bigDataAnalyticsDesc'),
    features: [t('landing.featureMassStorage'), t('landing.featureRealtimeAnalysis'), t('landing.featureDataMining'), t('landing.featureVisualization')]
  }
]

// 核心优势
const features = [
  { id: 1, title: t('landing.availability99'), icon: 'Shield', description: t('landing.availability99Desc') },
  { id: 2, title: t('landing.elasticScaling'), icon: 'TrendCharts', description: t('landing.elasticScalingDesc') },
  { id: 3, title: t('landing.support24x7'), icon: 'Headset', description: t('landing.support24x7Desc') },
  { id: 4, title: t('landing.globalNodesFeature'), icon: 'Connection', description: t('landing.globalNodesFeatureDesc') },
  { id: 5, title: t('landing.securityReliable'), icon: 'Shield', description: t('landing.securityReliableDesc') },
  { id: 6, title: t('landing.costEffective'), icon: 'TrendCharts', description: t('landing.costEffectiveDesc') }
]

// 数据统计
const stats = [
  { id: 1, value: '10,000+', label: t('landing.statClients') },
  { id: 2, value: '99.9%', label: t('landing.statAvailability') },
  { id: 3, value: '30+', label: t('landing.statDataCenters') },
  { id: 4, value: '7×24', label: t('landing.statTechSupport') }
]

// 获取数据
const fetchData = async () => {
  try {
    // 获取顶部导航
    const navRes = await request.get('/api/v1/nav/top')
    if (navRes.data?.data) {
      topNavs.value = navRes.data.data
    }
    
    // 获取底部导航
    const bottomNavRes = await request.get('/api/v1/nav/bottom')
    if (bottomNavRes.data?.data) {
      bottomNavs.value = bottomNavRes.data.data
    }
    
    // 获取轮播图（从站点设置中获取 banners，或使用独立 API）
    try {
      const bannerRes = await request.get('/api/v1/banners')
      if (bannerRes.data?.data?.length) {
        banners.value = bannerRes.data.data
      }
    } catch {
      // banners API 不可用时保留默认数据
    }
    
    // 获取产品分组
    const groupRes = await request.get('/api/v1/product-groups')
    if (groupRes.data?.data) {
      productGroups.value = groupRes.data.data
    }
    
    // 获取热门产品（销量前4）
    const hotRes = await request.get('/api/v1/products/hot', { params: { limit: 4 } })
    if (hotRes.data?.data) {
      hotProducts.value = hotRes.data.data
    }
    
    // 获取新闻公告
    const newsRes = await request.get('/api/v1/news', { params: { limit: 3 } })
    if (newsRes.data?.data) {
      announcements.value = newsRes.data.data
    }
    
    // 获取合作伙伴（后端暂无此 API，保留空数组）
    // const partnerRes = await request.get('/api/v1/partners')
    // if (partnerRes.data?.data) {
    //   partners.value = partnerRes.data.data
    // }
    
    // 获取站点设置
    const settingRes = await request.get('/api/v1/system/settings')
    if (settingRes.data?.data) {
      siteSettings.value = settingRes.data.data
    }
  } catch (error) {
    console.error('获取数据失败:', error)
  }
}

// 轮播图控制
const goToSlide = (index: number) => {
  carouselRef.value?.setActiveItem(index)
  currentSlide.value = index
}

// 监听滚动
const handleScroll = () => {
  scrolled.value = window.scrollY > 50
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
  fetchData()
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped lang="scss">
.home {
  min-height: 100vh;
}

// Header
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  transition: all 0.3s ease;
  background: transparent;
  
  &.header-scrolled {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
    
    .nav-link, .logo-text, .login-btn {
      color: #333 !important;
    }
  }
  
  .header-inner {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
    height: 70px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  
  .logo {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    
    .logo-img {
      width: 36px;
      height: 36px;
    }
    
    .logo-text {
      font-size: 20px;
      font-weight: 600;
      color: #fff;
    }
  }
  
  .nav-links {
    display: flex;
    align-items: center;
    gap: 30px;
    
    .nav-link {
      color: rgba(255, 255, 255, 0.9);
      text-decoration: none;
      font-size: 15px;
      cursor: pointer;
      transition: color 0.3s;
      
      &:hover {
        color: #fff;
      }
    }
  }
  
  .header-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    
    .login-btn {
      color: rgba(255, 255, 255, 0.9);
    }
  }
}

// Carousel
.hero-section {
  position: relative;
  
  .carousel-wrapper {
    position: relative;
  }
  
  .carousel-btn {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    z-index: 10;
    width: 48px;
    height: 48px;
    border: none;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(10px);
    color: #fff;
    cursor: pointer;
    transition: all 0.3s;
    display: flex;
    align-items: center;
    justify-content: center;
    
    &:hover {
      background: rgba(255, 255, 255, 0.4);
    }
    
    &-prev {
      left: 20px;
    }
    
    &-next {
      right: 20px;
    }
  }
  
  .carousel-indicators {
    position: absolute;
    bottom: 80px;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    gap: 8px;
    z-index: 10;
    
    .indicator {
      width: 32px;
      height: 4px;
      border: none;
      border-radius: 2px;
      background: rgba(255, 255, 255, 0.4);
      cursor: pointer;
      transition: all 0.3s;
      
      &.active {
        width: 48px;
        background: #fff;
      }
    }
  }
}

.carousel-slide {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  
  .slide-video {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  .slide-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.4);
  }
  
  .slide-content {
    position: relative;
    z-index: 1;
    text-align: center;
    color: #fff;
    padding: 0 20px;
    
    .slide-badge {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 6px 16px;
      background: rgba(255, 255, 255, 0.2);
      border-radius: 20px;
      font-size: 14px;
      margin-bottom: 20px;
    }
    
    .slide-title {
      font-size: 48px;
      font-weight: 700;
      margin: 0 0 16px;
      text-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
    }
    
    .slide-desc {
      font-size: 18px;
      margin: 0 0 32px;
      opacity: 0.9;
    }
    
    .slide-actions {
      display: flex;
      gap: 16px;
      justify-content: center;
      
      .slide-ghost-btn {
        background: transparent;
        border-color: rgba(255, 255, 255, 0.6);
        color: #fff;
        
        &:hover {
          background: rgba(255, 255, 255, 0.1);
          border-color: #fff;
        }
      }
    }
  }
}

// 快捷入口
.quick-entry {
  position: absolute;
  bottom: -40px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 10;
  display: flex;
  gap: 20px;
  background: #fff;
  padding: 24px 40px;
  border-radius: 16px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  
  .entry-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 12px 24px;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.3s;
    color: #666;
    
    &:hover {
      background: #f5f7fa;
      color: #1a73e8;
    }
    
    span {
      font-size: 14px;
    }
  }
}

// Sections
.section {
  padding: 80px 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.section-header {
  text-align: center;
  margin-bottom: 48px;
  
  .section-title {
    font-size: 32px;
    font-weight: 600;
    color: #1d1d1f;
    margin: 0 0 12px;
  }
  
  .section-subtitle {
    font-size: 16px;
    color: #86868b;
    margin: 0;
  }
}

// Hot Products
.hot-products-section {
  background: #f5f7fa;
  padding-top: 120px;
  
  .products-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 24px;
    
    @media (max-width: 1200px) {
      grid-template-columns: repeat(2, 1fr);
    }
  }
  
  .product-card {
    background: #fff;
    border-radius: 16px;
    padding: 32px 24px;
    text-align: center;
    cursor: pointer;
    transition: all 0.3s;
    position: relative;
    
    &:hover {
      transform: translateY(-8px);
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
    }
    
    .product-badge {
      position: absolute;
      top: 16px;
      right: 16px;
      padding: 4px 12px;
      background: #ff3b30;
      color: #fff;
      border-radius: 12px;
      font-size: 12px;
    }
    
    .product-icon {
      width: 72px;
      height: 72px;
      border-radius: 20px;
      display: flex;
      align-items: center;
      justify-content: center;
      margin: 0 auto 20px;
      color: #fff;
    }
    
    .product-name {
      font-size: 18px;
      font-weight: 600;
      margin: 0 0 8px;
    }
    
    .product-desc {
      font-size: 14px;
      color: #86868b;
      margin: 0 0 20px;
      line-height: 1.5;
    }
    
    .product-price {
      margin-bottom: 20px;
      
      .price-symbol {
        font-size: 16px;
        color: #1a73e8;
      }
      
      .price-value {
        font-size: 32px;
        font-weight: 700;
        color: #1a73e8;
      }
      
      .price-unit {
        font-size: 14px;
        color: #86868b;
      }
    }
  }
  
  .section-more {
    text-align: center;
    margin-top: 40px;
  }
}

// Solutions
.solutions-section {
  .solutions-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 24px;
    
    @media (max-width: 1200px) {
      grid-template-columns: repeat(2, 1fr);
    }
  }
  
  .solution-card {
    background: #fff;
    border-radius: 16px;
    padding: 32px;
    border: 1px solid #e5e5ea;
    transition: all 0.3s;
    
    &:hover {
      border-color: #1a73e8;
      box-shadow: 0 10px 30px rgba(26, 115, 232, 0.1);
    }
    
    .solution-icon {
      width: 64px;
      height: 64px;
      background: linear-gradient(135deg, #1a73e8, #4a90e2);
      border-radius: 16px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
      margin-bottom: 20px;
    }
    
    h3 {
      font-size: 20px;
      font-weight: 600;
      margin: 0 0 12px;
    }
    
    p {
      font-size: 14px;
      color: #86868b;
      margin: 0 0 16px;
    }
    
    .solution-features {
      list-style: none;
      padding: 0;
      margin: 0 0 20px;
      
      li {
        padding: 6px 0;
        font-size: 14px;
        color: #606266;
        
        &::before {
          content: '✓';
          color: #1a73e8;
          margin-right: 8px;
        }
      }
    }
  }
}

// Features
.features-section {
  background: #f5f7fa;
  
  .features-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 24px;
    
    @media (max-width: 768px) {
      grid-template-columns: repeat(2, 1fr);
    }
  }
  
  .feature-card {
    background: #fff;
    border-radius: 16px;
    padding: 32px;
    text-align: center;
    transition: all 0.3s;
    
    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
    }
    
    .feature-icon {
      width: 64px;
      height: 64px;
      background: linear-gradient(135deg, #1a73e8, #4a90e2);
      border-radius: 16px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
      margin: 0 auto 20px;
    }
    
    h3 {
      font-size: 18px;
      font-weight: 600;
      margin: 0 0 12px;
    }
    
    p {
      font-size: 14px;
      color: #86868b;
      margin: 0;
    }
  }
}

// Stats
.stats-section {
  background: linear-gradient(135deg, #1a237e, #0d47a1);
  
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 24px;
    text-align: center;
    color: #fff;
    
    @media (max-width: 768px) {
      grid-template-columns: repeat(2, 1fr);
    }
  }
  
  .stat-item {
    .stat-value {
      font-size: 40px;
      font-weight: 700;
      margin-bottom: 8px;
    }
    
    .stat-label {
      font-size: 16px;
      opacity: 0.8;
    }
  }
}

// Announcements
.announcements-section {
  .news-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 24px;
    
    @media (max-width: 768px) {
      grid-template-columns: 1fr;
    }
  }
  
  .news-card {
    background: #fff;
    border-radius: 16px;
    padding: 24px;
    display: flex;
    gap: 20px;
    cursor: pointer;
    transition: all 0.3s;
    border: 1px solid #e5e5ea;
    
    &:hover {
      border-color: #1a73e8;
      box-shadow: 0 10px 30px rgba(26, 115, 232, 0.1);
    }
    
    .news-date {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      width: 60px;
      height: 60px;
      background: linear-gradient(135deg, #1a73e8, #4a90e2);
      border-radius: 12px;
      color: #fff;
      flex-shrink: 0;
      
      .day {
        font-size: 24px;
        font-weight: 700;
        line-height: 1;
      }
      
      .month {
        font-size: 12px;
      }
    }
    
    .news-content {
      flex: 1;
      
      h3 {
        font-size: 16px;
        font-weight: 600;
        margin: 0 0 8px;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }
      
      p {
        font-size: 14px;
        color: #86868b;
        margin: 0;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }
    }
  }
}

// Partners
.partners-section {
  background: #f5f7fa;
  
  .partners-grid {
    display: grid;
    grid-template-columns: repeat(6, 1fr);
    gap: 24px;
    
    @media (max-width: 768px) {
      grid-template-columns: repeat(3, 1fr);
    }
  }
  
  .partner-card {
    background: #fff;
    border-radius: 12px;
    padding: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s;
    border: 1px solid #e5e5ea;
    
    &:hover {
      border-color: #1a73e8;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    }
    
    .partner-logo {
      max-width: 100%;
      height: 40px;
      object-fit: contain;
      filter: grayscale(100%);
      opacity: 0.6;
      transition: all 0.3s;
    }
    
    &:hover .partner-logo {
      filter: grayscale(0);
      opacity: 1;
    }
  }
}

// Footer
.footer {
  background: #1d1d1f;
  color: #fff;
  padding: 60px 0 30px;
  
  .footer-grid {
    display: grid;
    grid-template-columns: 2fr 1fr 1fr 1fr;
    gap: 40px;
    margin-bottom: 40px;
    
    @media (max-width: 768px) {
      grid-template-columns: 1fr;
    }
  }
  
  .footer-col {
    .footer-logo {
      display: flex;
      align-items: center;
      gap: 10px;
      margin-bottom: 16px;
      
      img {
        width: 32px;
        height: 32px;
      }
      
      span {
        font-size: 18px;
        font-weight: 600;
      }
    }
    
    .footer-desc {
      font-size: 14px;
      color: rgba(255, 255, 255, 0.6);
      margin: 0 0 20px;
      line-height: 1.6;
    }
    
    .footer-contact {
      p {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 14px;
        color: rgba(255, 255, 255, 0.6);
        margin: 8px 0;
      }
    }
    
    h4 {
      font-size: 16px;
      font-weight: 600;
      margin: 0 0 20px;
    }
    
    ul {
      list-style: none;
      padding: 0;
      margin: 0;
      
      li {
        margin-bottom: 12px;
        
        a {
          color: rgba(255, 255, 255, 0.6);
          text-decoration: none;
          font-size: 14px;
          transition: color 0.3s;
          
          &:hover {
            color: #fff;
          }
        }
      }
    }
  }
  
  .footer-bottom {
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    padding-top: 24px;
    text-align: center;
    
    p {
      font-size: 13px;
      color: rgba(255, 255, 255, 0.4);
      margin: 4px 0;
    }
  }
}
</style>
