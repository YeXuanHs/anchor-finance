<template>
  <div class="home">
    <!-- Fixed Header -->
    <header class="header" :class="{ 'header-scrolled': scrolled }">
      <div class="header-inner">
        <div class="logo" @click="$router.push('/')">
          <div class="logo-icon">
            <n-icon size="26" color="#fff">
              <AnchorOutline />
            </n-icon>
          </div>
          <span class="logo-text">锚点财务</span>
        </div>
        <nav class="nav-links">
          <router-link to="/" class="nav-link" :class="{ active: activeNav === 'home' }" @click="activeNav = 'home'">首页</router-link>
          <router-link to="/products" class="nav-link" :class="{ active: activeNav === 'products' }" @click="activeNav = 'products'">产品</router-link>
          <a href="#announcements" class="nav-link" :class="{ active: activeNav === 'news' }" @click="activeNav = 'news'">新闻</a>
          <a href="#footer" class="nav-link" :class="{ active: activeNav === 'help' }" @click="activeNav = 'help'">帮助</a>
        </nav>
        <div class="header-actions">
          <n-button text class="login-btn" @click="$router.push('/login')">
            <n-icon size="16" style="margin-right: 4px;"><PersonOutline /></n-icon>
            登录
          </n-button>
          <n-button type="primary" round size="small" strong @click="$router.push('/register')">
            免费注册
          </n-button>
        </div>
      </div>
    </header>

    <!-- Hero Carousel -->
    <section id="hero" class="hero-section">
      <n-carousel
        :autoplay="true"
        :interval="5000"
        :show-dots="true"
        dot-type="line"
        effect="fade"
        class="hero-carousel"
      >
        <!-- Slide 1 -->
        <div class="carousel-slide slide-1">
          <div class="slide-shapes">
            <div class="shape shape-1"></div>
            <div class="shape shape-2"></div>
            <div class="shape shape-3"></div>
          </div>
          <div class="slide-content">
            <div class="slide-badge">
              <n-icon size="14"><FlashOutline /></n-icon>
              热门推荐
            </div>
            <h1 class="slide-title">高性能云服务器</h1>
            <p class="slide-desc">企业级云服务器，99.9% SLA 保障，弹性配置分钟级部署，助力业务高速发展</p>
            <div class="slide-actions">
              <n-button type="primary" size="large" round strong @click="$router.push('/products')">
                <template #icon><n-icon><RocketOutline /></n-icon></template>
                立即选购
              </n-button>
              <n-button size="large" round ghost class="slide-ghost-btn">
                了解更多
              </n-button>
            </div>
          </div>
        </div>
        <!-- Slide 2 -->
        <div class="carousel-slide slide-2">
          <div class="slide-shapes">
            <div class="shape shape-1"></div>
            <div class="shape shape-2"></div>
            <div class="shape shape-3"></div>
          </div>
          <div class="slide-content">
            <div class="slide-badge">
              <n-icon size="14"><GlobeOutline /></n-icon>
              全球布局
            </div>
            <h1 class="slide-title">全球节点覆盖</h1>
            <p class="slide-desc">遍布全球 30+ 数据中心节点，就近接入极速体验，让您的业务触达世界每个角落</p>
            <div class="slide-actions">
              <n-button type="primary" size="large" round strong @click="$router.push('/products')">
                <template #icon><n-icon><RocketOutline /></n-icon></template>
                立即选购
              </n-button>
              <n-button size="large" round ghost class="slide-ghost-btn">
                了解更多
              </n-button>
            </div>
          </div>
        </div>
        <!-- Slide 3 -->
        <div class="carousel-slide slide-3">
          <div class="slide-shapes">
            <div class="shape shape-1"></div>
            <div class="shape shape-2"></div>
            <div class="shape shape-3"></div>
          </div>
          <div class="slide-content">
            <div class="slide-badge">
              <n-icon size="14"><HeadsetOutline /></n-icon>
              专业服务
            </div>
            <h1 class="slide-title">专业技术支持</h1>
            <p class="slide-desc">7×24 小时专业工程师团队全天候在线，1 分钟极速响应，为您的业务保驾护航</p>
            <div class="slide-actions">
              <n-button type="primary" size="large" round strong @click="$router.push('/products')">
                <template #icon><n-icon><RocketOutline /></n-icon></template>
                立即选购
              </n-button>
              <n-button size="large" round ghost class="slide-ghost-btn">
                了解更多
              </n-button>
            </div>
          </div>
        </div>
        <!-- Slide 4 -->
        <div class="carousel-slide slide-4">
          <div class="slide-shapes">
            <div class="shape shape-1"></div>
            <div class="shape shape-2"></div>
            <div class="shape shape-3"></div>
          </div>
          <div class="slide-content">
            <div class="slide-badge">
              <n-icon size="14"><ShieldCheckmarkOutline /></n-icon>
              安全保障
            </div>
            <h1 class="slide-title">安全可靠</h1>
            <p class="slide-desc">多副本冗余存储，DDoS 防护，自动备份容灾切换，全方位守护您的数据资产</p>
            <div class="slide-actions">
              <n-button type="primary" size="large" round strong @click="$router.push('/products')">
                <template #icon><n-icon><RocketOutline /></n-icon></template>
                立即选购
              </n-button>
              <n-button size="large" round ghost class="slide-ghost-btn">
                了解更多
              </n-button>
            </div>
          </div>
        </div>
      </n-carousel>
    </section>

    <!-- Product Showcase -->
    <section id="products" class="products-section">
      <div class="section-container">
        <div class="section-header">
          <span class="section-tag">产品服务</span>
          <h2 class="section-title">热门产品</h2>
          <p class="section-subtitle">为您提供全方位的云计算解决方案，满足各类业务场景需求</p>
        </div>
        <n-grid :x-gap="24" :y-gap="24" :cols="4" responsive="screen" item-responsive>
          <n-gi span="4 m:2 l:1" v-for="product in hotProducts" :key="product.id">
            <div class="product-card" @click="$router.push(`/products/${product.id}`)">
              <div class="product-card-top" :style="{ background: product.gradient }">
                <div class="product-icon-wrap">
                  <n-icon size="36" color="#fff">
                    <component :is="product.icon" />
                  </n-icon>
                </div>
              </div>
              <div class="product-card-body">
                <h3 class="product-name">{{ product.name }}</h3>
                <p class="product-desc">{{ product.description }}</p>
                <div class="product-price">
                  <span class="price-currency">¥</span>
                  <span class="price-amount">{{ product.price }}</span>
                  <span class="price-unit">/月起</span>
                </div>
                <n-button type="primary" block round strong class="product-btn">
                  立即选购
                </n-button>
              </div>
            </div>
          </n-gi>
        </n-grid>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="features-section">
      <div class="section-container">
        <div class="section-header">
          <span class="section-tag">核心优势</span>
          <h2 class="section-title">为什么选择我们</h2>
          <p class="section-subtitle">六大核心优势，为您的业务保驾护航</p>
        </div>
        <n-grid :x-gap="24" :y-gap="24" :cols="3" responsive="screen" item-responsive>
          <n-gi span="3 m:1" v-for="feature in features" :key="feature.title">
            <div class="feature-card">
              <div class="feature-icon-circle" :style="{ background: feature.gradient }">
                <n-icon size="28" color="#fff">
                  <component :is="feature.icon" />
                </n-icon>
              </div>
              <h3 class="feature-title">{{ feature.title }}</h3>
              <p class="feature-desc">{{ feature.description }}</p>
            </div>
          </n-gi>
        </n-grid>
      </div>
    </section>

    <!-- Stats Section -->
    <section class="stats-section">
      <div class="stats-bg-pattern"></div>
      <div class="section-container">
        <n-grid :x-gap="32" :cols="4" responsive="screen" item-responsive>
          <n-gi span="4 m:2 l:1" v-for="stat in stats" :key="stat.label">
            <div class="stat-item">
              <div class="stat-icon-ring">
                <n-icon size="24" color="#fff">
                  <component :is="stat.icon" />
                </n-icon>
              </div>
              <div class="stat-number">
                <span class="stat-value">{{ stat.value }}</span>
                <span class="stat-suffix">{{ stat.suffix }}</span>
              </div>
              <div class="stat-label">{{ stat.label }}</div>
            </div>
          </n-gi>
        </n-grid>
      </div>
    </section>

    <!-- Announcements Section -->
    <section id="announcements" class="announcements-section">
      <div class="section-container">
        <div class="section-header">
          <span class="section-tag">最新动态</span>
          <h2 class="section-title">新闻公告</h2>
          <p class="section-subtitle">了解锚点财务最新产品动态与行业资讯</p>
        </div>
        <n-grid :x-gap="24" :y-gap="24" :cols="3" responsive="screen" item-responsive>
          <n-gi span="3 m:1" v-for="item in announcements" :key="item.id">
            <div class="announcement-card">
              <div class="announcement-tag" :style="{ background: item.tagColor }">{{ item.tag }}</div>
              <h3 class="announcement-title">{{ item.title }}</h3>
              <p class="announcement-summary">{{ item.summary }}</p>
              <div class="announcement-meta">
                <span class="announcement-date">
                  <n-icon size="14"><CalendarOutline /></n-icon>
                  {{ item.date }}
                </span>
                <a href="#" class="announcement-link">
                  阅读全文
                  <n-icon size="14"><ArrowForwardOutline /></n-icon>
                </a>
              </div>
            </div>
          </n-gi>
        </n-grid>
      </div>
    </section>

    <!-- Footer -->
    <footer id="footer" class="footer">
      <div class="section-container">
        <div class="footer-grid">
          <div class="footer-brand">
            <div class="logo">
              <div class="logo-icon">
                <n-icon size="24" color="#fff">
                  <AnchorOutline />
                </n-icon>
              </div>
              <span class="logo-text">锚点财务</span>
            </div>
            <p class="footer-desc">专业的云计算服务提供商，致力于为企业提供高性能、高可用的云基础设施服务。我们以技术驱动创新，助力千行百业数字化转型。</p>
            <div class="footer-socials">
              <a href="#" class="social-link" title="微信"><n-icon size="18"><LogoWechat /></n-icon></a>
              <a href="#" class="social-link" title="微博"><n-icon size="18"><LogoGithub /></n-icon></a>
              <a href="#" class="social-link" title="邮箱"><n-icon size="18"><MailOutline /></n-icon></a>
            </div>
          </div>
          <div class="footer-column">
            <h4 class="footer-col-title">关于我们</h4>
            <a href="#" class="footer-link">公司介绍</a>
            <a href="#" class="footer-link">发展历程</a>
            <a href="#" class="footer-link">加入我们</a>
            <a href="#" class="footer-link">合作伙伴</a>
          </div>
          <div class="footer-column">
            <h4 class="footer-col-title">产品服务</h4>
            <router-link to="/products" class="footer-link">云服务器</router-link>
            <a href="#" class="footer-link">对象存储</a>
            <a href="#" class="footer-link">CDN 加速</a>
            <a href="#" class="footer-link">云数据库</a>
          </div>
          <div class="footer-column">
            <h4 class="footer-col-title">帮助支持</h4>
            <a href="#" class="footer-link">帮助文档</a>
            <a href="#" class="footer-link">API 文档</a>
            <a href="#" class="footer-link">提交工单</a>
            <a href="#" class="footer-link">常见问题</a>
          </div>
          <div class="footer-column">
            <h4 class="footer-col-title">联系方式</h4>
            <p class="footer-contact">
              <n-icon size="16"><CallOutline /></n-icon>
              400-888-8888
            </p>
            <p class="footer-contact">
              <n-icon size="16"><MailOutline /></n-icon>
              support@anchorfinance.com
            </p>
            <p class="footer-contact">
              <n-icon size="16"><LocationOutline /></n-icon>
              北京市朝阳区科技园区
            </p>
          </div>
        </div>
        <div class="footer-bottom">
          <p>&copy; 2024-2026 锚点财务 版权所有 | 京ICP备XXXXXXXX号</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import {
  AnchorOutline,
  ServerOutline,
  GlobeOutline,
  HeadsetOutline,
  FlashOutline,
  ShieldCheckmarkOutline,
  RocketOutline,
  PeopleOutline,
  CloudOutline,
  TimeOutline,
  CallOutline,
  MailOutline,
  LocationOutline,
  PersonOutline,
  CalendarOutline,
  ArrowForwardOutline,
  LogoWechat,
  LogoGithub,
  HardwareChipOutline,
  WifiOutline
} from '@vicons/ionicons5'

const scrolled = ref(false)
const activeNav = ref('home')

function handleScroll() {
  scrolled.value = window.scrollY > 20
}

onMounted(() => window.addEventListener('scroll', handleScroll))
onUnmounted(() => window.removeEventListener('scroll', handleScroll))

const hotProducts = [
  {
    id: 1,
    name: '云服务器 ECS',
    description: '弹性计算服务，分钟级创建，按需付费，适用于各类业务场景',
    price: 49,
    icon: ServerOutline,
    gradient: 'linear-gradient(135deg, #1890ff, #096dd9)'
  },
  {
    id: 2,
    name: '对象存储 OSS',
    description: '海量安全存储，支持多协议接入，数据可靠性达 99.999999999%',
    price: 9,
    icon: CloudOutline,
    gradient: 'linear-gradient(135deg, #36cfc9, #13c2c2)'
  },
  {
    id: 3,
    name: 'CDN 加速',
    description: '全球 2800+ 节点，智能调度，毫秒级响应，提升用户访问体验',
    price: 19,
    icon: GlobeOutline,
    gradient: 'linear-gradient(135deg, #597ef7, #2f54eb)'
  },
  {
    id: 4,
    name: '云数据库 RDS',
    description: '高性能关系型数据库，自动备份，容灾切换，省心省力',
    price: 99,
    icon: HardwareChipOutline,
    gradient: 'linear-gradient(135deg, #ffc53d, #faad14)'
  }
]

const features = [
  {
    title: '99.9% SLA 保障',
    description: '企业级服务等级协议，确保您的业务稳定运行，故障自动赔偿',
    icon: ShieldCheckmarkOutline,
    gradient: 'linear-gradient(135deg, #1890ff, #096dd9)'
  },
  {
    title: '7×24 技术支持',
    description: '专业工程师团队全天候在线，1 分钟响应，快速解决技术问题',
    icon: HeadsetOutline,
    gradient: 'linear-gradient(135deg, #36cfc9, #13c2c2)'
  },
  {
    title: '全球节点覆盖',
    description: '遍布全球 30+ 区域数据中心，就近接入，为全球用户提供极速体验',
    icon: GlobeOutline,
    gradient: 'linear-gradient(135deg, #597ef7, #2f54eb)'
  },
  {
    title: '弹性扩展',
    description: '支持分钟级资源升降配，按需使用，灵活应对业务高峰',
    icon: FlashOutline,
    gradient: 'linear-gradient(135deg, #ffc53d, #faad14)'
  },
  {
    title: '数据安全',
    description: '多副本冗余存储，自动备份，DDoS 防护，全方位保障数据安全',
    icon: ShieldCheckmarkOutline,
    gradient: 'linear-gradient(135deg, #ff7a45, #fa541c)'
  },
  {
    title: '快速部署',
    description: '一键部署常用应用环境，分钟级完成服务器创建与配置',
    icon: RocketOutline,
    gradient: 'linear-gradient(135deg, #b37feb, #722ed1)'
  }
]

const stats = [
  { label: '服务客户', value: '10,000', suffix: '+', icon: PeopleOutline },
  { label: '产品实例', value: '50,000', suffix: '+', icon: ServerOutline },
  { label: '数据中心', value: '20', suffix: '+', icon: GlobeOutline },
  { label: '运行时间', value: '99.99', suffix: '%', icon: TimeOutline }
]

const announcements = [
  {
    id: 1,
    tag: '产品更新',
    tagColor: 'linear-gradient(135deg, #1890ff, #096dd9)',
    title: '云服务器 ECS 全面升级至第七代实例',
    summary: '全新第七代云服务器实例正式上线，搭载最新一代处理器，计算性能提升 40%，内存带宽提升 30%，为您的业务提供更强劲的算力支持。',
    date: '2026-07-25'
  },
  {
    id: 2,
    tag: '优惠活动',
    tagColor: 'linear-gradient(135deg, #ff7a45, #fa541c)',
    title: '新用户专享：首月 1 折起，限时特惠',
    summary: '即日起至 2026 年 8 月 31 日，新注册用户可享受云服务器首月 1 折优惠，对象存储、CDN 等产品同步参与活动，最高可省 2000 元。',
    date: '2026-07-20'
  },
  {
    id: 3,
    tag: '技术博客',
    tagColor: 'linear-gradient(135deg, #36cfc9, #13c2c2)',
    title: '深度解析：如何构建高可用架构方案',
    summary: '本文详细介绍了基于锚点财务云服务构建高可用架构的最佳实践，涵盖多可用区部署、负载均衡配置、自动容灾切换等核心内容。',
    date: '2026-07-15'
  }
]
</script>

<style scoped>
/* ==================== Header ==================== */
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  z-index: 1000;
  transition: all 0.35s cubic-bezier(0.4, 0, 0.2, 1);
  background: transparent;
}

.header-scrolled {
  background: rgba(255, 255, 255, 0.96);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  box-shadow: 0 1px 12px rgba(0, 0, 0, 0.08);
}

.header-inner {
  max-width: 1200px;
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
  cursor: pointer;
  user-select: none;
}

.logo-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.35);
  transition: transform 0.3s;
}

.logo:hover .logo-icon {
  transform: rotate(-12deg) scale(1.05);
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 1px;
  transition: color 0.35s;
}

.header-scrolled .logo-text {
  color: #1d2129;
}

.nav-links {
  display: flex;
  gap: 36px;
}

.nav-link {
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  font-size: 15px;
  font-weight: 500;
  transition: color 0.25s;
  position: relative;
  padding: 4px 0;
}

.nav-link:hover,
.nav-link.active {
  color: #fff;
}

.header-scrolled .nav-link {
  color: #4e5969;
}

.header-scrolled .nav-link:hover,
.header-scrolled .nav-link.active {
  color: #1890ff;
}

.nav-link.active::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 50%;
  transform: translateX(-50%);
  width: 20px;
  height: 2px;
  background: #1890ff;
  border-radius: 1px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.login-btn {
  color: rgba(255, 255, 255, 0.85) !important;
  font-weight: 500;
}

.header-scrolled .login-btn {
  color: #4e5969 !important;
}

/* ==================== Hero Carousel ==================== */
.hero-section {
  margin-top: -64px;
}

.hero-carousel {
  height: 640px;
}

.carousel-slide {
  height: 640px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.slide-shapes {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
}

.shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.08;
  background: #fff;
}

.shape-1 {
  width: 400px;
  height: 400px;
  top: -100px;
  right: -80px;
  animation: float 8s ease-in-out infinite;
}

.shape-2 {
  width: 200px;
  height: 200px;
  bottom: -40px;
  left: 10%;
  animation: float 6s ease-in-out infinite reverse;
}

.shape-3 {
  width: 120px;
  height: 120px;
  top: 30%;
  left: 5%;
  animation: float 10s ease-in-out infinite 1s;
}

@keyframes float {
  0%, 100% { transform: translateY(0) scale(1); }
  50% { transform: translateY(-20px) scale(1.05); }
}

.slide-1 {
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 40%, #0050b3 100%);
}

.slide-2 {
  background: linear-gradient(135deg, #13c2c2 0%, #08979c 40%, #006d75 100%);
}

.slide-3 {
  background: linear-gradient(135deg, #597ef7 0%, #2f54eb 40%, #1d39c4 100%);
}

.slide-4 {
  background: linear-gradient(135deg, #1890ff 0%, #2f54eb 40%, #1d39c4 100%);
}

.slide-content {
  text-align: center;
  color: #fff;
  z-index: 1;
  padding: 0 24px;
}

.slide-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 16px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(8px);
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 24px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.slide-title {
  font-size: 52px;
  font-weight: 800;
  margin-bottom: 16px;
  text-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
  letter-spacing: 2px;
}

.slide-desc {
  font-size: 18px;
  opacity: 0.9;
  margin-bottom: 40px;
  max-width: 540px;
  margin-left: auto;
  margin-right: auto;
  line-height: 1.7;
}

.slide-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.slide-ghost-btn {
  border-color: rgba(255, 255, 255, 0.5) !important;
  color: #fff !important;
  transition: all 0.25s !important;
}

.slide-ghost-btn:hover {
  border-color: #fff !important;
  background: rgba(255, 255, 255, 0.12) !important;
}

/* ==================== Section Common ==================== */
.section-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

.section-header {
  text-align: center;
  margin-bottom: 48px;
}

.section-tag {
  display: inline-block;
  padding: 4px 14px;
  background: linear-gradient(135deg, #e6f7ff, #bae7ff);
  color: #1890ff;
  font-size: 13px;
  font-weight: 600;
  border-radius: 20px;
  margin-bottom: 16px;
  letter-spacing: 1px;
}

.section-title {
  font-size: 32px;
  font-weight: 700;
  color: #1d2129;
  margin-bottom: 12px;
  letter-spacing: 1px;
}

.section-subtitle {
  font-size: 16px;
  color: #86909c;
  line-height: 1.6;
}

/* ==================== Products Section ==================== */
.products-section {
  padding: 80px 0;
  background: #f5f7fa;
}

.product-card {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.35s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  border: 1px solid #e8e8e8;
}

.product-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 16px 40px rgba(24, 144, 255, 0.15);
  border-color: transparent;
}

.product-card-top {
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.product-icon-wrap {
  width: 64px;
  height: 64px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(8px);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.35s;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.product-card:hover .product-icon-wrap {
  transform: scale(1.1) rotate(-3deg);
}

.product-card-body {
  padding: 24px;
  text-align: center;
}

.product-name {
  font-size: 18px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 8px;
}

.product-desc {
  font-size: 13px;
  color: #86909c;
  line-height: 1.6;
  margin-bottom: 20px;
  min-height: 42px;
}

.product-price {
  margin-bottom: 20px;
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 2px;
}

.price-currency {
  font-size: 16px;
  font-weight: 600;
  color: #1890ff;
}

.price-amount {
  font-size: 32px;
  font-weight: 700;
  color: #1890ff;
  line-height: 1;
}

.price-unit {
  font-size: 14px;
  color: #86909c;
}

.product-btn {
  width: 140px;
}

/* ==================== Features Section ==================== */
.features-section {
  padding: 80px 0;
  background: #fff;
}

.feature-card {
  text-align: center;
  padding: 40px 28px;
  border-radius: 12px;
  transition: all 0.35s cubic-bezier(0.4, 0, 0.2, 1);
  background: #f5f7fa;
  border: 1px solid transparent;
}

.feature-card:hover {
  background: #fff;
  box-shadow: 0 12px 32px rgba(24, 144, 255, 0.1);
  transform: translateY(-4px);
  border-color: #e8e8e8;
}

.feature-icon-circle {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  box-shadow: 0 6px 20px rgba(24, 144, 255, 0.25);
  transition: transform 0.35s;
}

.feature-card:hover .feature-icon-circle {
  transform: scale(1.1);
}

.feature-title {
  font-size: 18px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 10px;
}

.feature-desc {
  font-size: 14px;
  color: #86909c;
  line-height: 1.7;
}

/* ==================== Stats Section ==================== */
.stats-section {
  padding: 72px 0;
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
  position: relative;
  overflow: hidden;
}

.stats-bg-pattern {
  position: absolute;
  inset: 0;
  background-image:
    radial-gradient(circle at 20% 50%, rgba(255,255,255,0.06) 0%, transparent 50%),
    radial-gradient(circle at 80% 50%, rgba(255,255,255,0.04) 0%, transparent 50%);
  pointer-events: none;
}

.stat-item {
  text-align: center;
  position: relative;
  z-index: 1;
}

.stat-icon-ring {
  width: 48px;
  height: 48px;
  border: 2px solid rgba(255, 255, 255, 0.25);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.stat-number {
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 2px;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 40px;
  font-weight: 800;
  color: #fff;
  line-height: 1;
  letter-spacing: 1px;
}

.stat-suffix {
  font-size: 20px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.85);
}

.stat-label {
  font-size: 15px;
  color: rgba(255, 255, 255, 0.7);
  font-weight: 500;
}

/* ==================== Announcements Section ==================== */
.announcements-section {
  padding: 80px 0;
  background: #f5f7fa;
}

.announcement-card {
  background: #fff;
  border-radius: 12px;
  padding: 28px;
  transition: all 0.35s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.announcement-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.08);
  border-color: transparent;
}

.announcement-tag {
  display: inline-block;
  padding: 3px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 16px;
  align-self: flex-start;
}

.announcement-title {
  font-size: 18px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 12px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.announcement-summary {
  font-size: 14px;
  color: #86909c;
  line-height: 1.7;
  margin-bottom: 20px;
  flex: 1;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.announcement-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.announcement-date {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #c9cdd4;
}

.announcement-link {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #1890ff;
  text-decoration: none;
  font-weight: 500;
  transition: gap 0.25s;
}

.announcement-link:hover {
  gap: 8px;
}

/* ==================== Footer ==================== */
.footer {
  background: #1d2129;
  padding: 60px 0 0;
  color: #fff;
}

.footer-grid {
  display: grid;
  grid-template-columns: 1.5fr 1fr 1fr 1fr 1.2fr;
  gap: 40px;
  padding-bottom: 48px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.footer-brand .logo {
  margin-bottom: 16px;
}

.footer-brand .logo-text {
  color: #fff;
}

.footer-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.45);
  line-height: 1.7;
  margin-bottom: 20px;
}

.footer-socials {
  display: flex;
  gap: 12px;
}

.social-link {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.06);
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.5);
  text-decoration: none;
  transition: all 0.25s;
}

.social-link:hover {
  background: #1890ff;
  color: #fff;
}

.footer-col-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 20px;
  color: #fff;
}

.footer-link {
  display: block;
  color: rgba(255, 255, 255, 0.45);
  text-decoration: none;
  font-size: 14px;
  margin-bottom: 12px;
  transition: color 0.25s, padding-left 0.25s;
}

.footer-link:hover {
  color: #1890ff;
  padding-left: 4px;
}

.footer-contact {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.45);
  margin-bottom: 12px;
}

.footer-bottom {
  padding: 20px 0;
  text-align: center;
  color: rgba(255, 255, 255, 0.3);
  font-size: 13px;
}

/* ==================== Responsive ==================== */
@media (max-width: 992px) {
  .footer-grid {
    grid-template-columns: 1fr 1fr;
  }
  .slide-title {
    font-size: 36px;
  }
  .hero-carousel,
  .carousel-slide {
    height: 560px;
  }
}

@media (max-width: 768px) {
  .nav-links {
    display: none;
  }
  .hero-carousel,
  .carousel-slide {
    height: 500px;
  }
  .slide-title {
    font-size: 28px;
    letter-spacing: 0;
  }
  .slide-desc {
    font-size: 15px;
  }
  .slide-badge {
    font-size: 12px;
    padding: 4px 12px;
  }
  .section-title {
    font-size: 26px;
  }
  .footer-grid {
    grid-template-columns: 1fr;
  }
  .stat-value {
    font-size: 32px;
  }
}
</style>
