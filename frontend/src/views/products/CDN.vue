<template>
  <div class="cdn-page">
    <!-- Hero Banner -->
    <section class="hero-section">
      <div class="hero-content">
        <div class="hero-badge">加速服务</div>
        <h1 class="hero-title">CDN 加速</h1>
        <p class="hero-desc">全球节点智能分发，极速稳定的内容分发网络，让您的业务触达全球</p>
        <div class="hero-actions">
          <el-button type="primary" size="large" round @click="scrollToProducts">立即选购</el-button>
          <el-button size="large" round class="hero-ghost-btn" @click="$router.push('/contact')">联系我们</el-button>
        </div>
      </div>
      <div class="hero-features">
        <div class="hero-feature" v-for="feat in heroFeatures" :key="feat.label">
          <el-icon :size="24"><component :is="feat.icon" /></el-icon>
          <span>{{ feat.label }}</span>
        </div>
      </div>
    </section>

    <!-- 核心优势 -->
    <section class="section advantages-section">
      <div class="container">
        <h2 class="section-title text-center">核心优势</h2>
        <p class="section-subtitle text-center">为什么选择我们的CDN加速</p>
        <div class="advantages-grid">
          <div v-for="item in advantages" :key="item.title" class="advantage-card">
            <el-icon :size="32" :color="item.color"><component :is="item.icon" /></el-icon>
            <h3>{{ item.title }}</h3>
            <p>{{ item.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 产品列表 -->
    <section id="products" class="section products-section">
      <div class="container">
        <h2 class="section-title text-center">CDN套餐</h2>
        <p class="section-subtitle text-center">灵活计费，按需选择</p>
        <el-skeleton :loading="loading" animated :rows="8">
          <template #default>
            <div class="products-grid">
              <div
                v-for="plan in plans"
                :key="plan.id"
                class="product-card"
                :class="{ featured: plan.featured }"
              >
                <div class="product-badge" v-if="plan.featured">推荐</div>
                <h3>{{ plan.name }}</h3>
                <p class="product-desc">{{ plan.description }}</p>
                <div class="product-price">
                  <span class="price-symbol">¥</span>
                  <span class="price-value">{{ plan.price }}</span>
                  <span class="price-unit">/月</span>
                </div>
                <ul class="product-features">
                  <li v-for="feature in plan.features" :key="feature">
                    <el-icon><CircleCheck /></el-icon>
                    {{ feature }}
                  </li>
                </ul>
                <el-button
                  :type="plan.featured ? 'primary' : 'default'"
                  round
                  @click="$router.push(`/products/${plan.id}`)"
                >
                  立即选购
                </el-button>
              </div>
            </div>
          </template>
        </el-skeleton>
      </div>
    </section>

    <!-- 产品特点 -->
    <section class="section features-section">
      <div class="container">
        <h2 class="section-title text-center">产品特点</h2>
        <p class="section-subtitle text-center">全方位的内容分发加速方案</p>
        <div class="features-grid">
          <div v-for="feature in productFeatures" :key="feature.title" class="feature-card">
            <el-icon :size="28" color="#10b981"><component :is="feature.icon" /></el-icon>
            <h3>{{ feature.title }}</h3>
            <p>{{ feature.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 常见问题 -->
    <section class="section faq-section">
      <div class="container">
        <h2 class="section-title text-center">常见问题</h2>
        <p class="section-subtitle text-center">关于CDN加速的常见疑问</p>
        <el-collapse>
          <el-collapse-item v-for="faq in faqs" :key="faq.question" :title="faq.question">
            <p>{{ faq.answer }}</p>
          </el-collapse-item>
        </el-collapse>
      </div>
    </section>

    <!-- CTA -->
    <section class="section cta-section">
      <div class="container">
        <div class="cta-content">
          <h2>加速您的业务</h2>
          <p>选择适合的CDN套餐，让内容触达全球用户</p>
          <div class="cta-actions">
            <el-button type="primary" size="large" round @click="scrollToProducts">立即选购</el-button>
            <el-button size="large" round class="cta-ghost-btn" @click="$router.push('/contact')">联系销售</el-button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { CircleCheck } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(true)

const heroFeatures = [
  { icon: 'Globe', label: '全球节点' },
  { icon: 'Lightning', label: '毫秒响应' },
  { icon: 'Connection', label: '智能调度' },
  { icon: 'Monitor', label: '实时监控' }
]

const advantages = [
  { title: '全球节点覆盖', description: '遍布全球200+加速节点，覆盖六大洲，就近接入极速体验', icon: 'Globe', color: '#10b981' },
  { title: '智能调度引擎', description: '基于实时网络质量监测，智能选择最优加速路径', icon: 'Connection', color: '#3b82f6' },
  { title: '极速缓存加速', description: '多级缓存架构，静态资源毫秒级响应，大幅提升访问速度', icon: 'Lightning', color: '#f59e0b' },
  { title: '安全防护能力', description: '集成DDoS防护和WAF防火墙，加速的同时保障安全', icon: 'Monitor', color: '#ef4444' },
  { title: '灵活计费模式', description: '按流量/按带宽/按请求次数多种计费方式，灵活选择', icon: 'TrendCharts', color: '#8b5cf6' },
  { title: '7×24技术支持', description: '专业技术团队全天候在线，快速响应解决各类问题', icon: 'Headset', color: '#0ea5e9' }
]

const productFeatures = [
  { title: '静态加速', description: '网站图片、CSS、JS等静态资源极速分发', icon: 'VideoPlay' },
  { title: '动态加速', description: 'API接口、动态页面智能路由优化', icon: 'Lightning' },
  { title: 'HTTPS加速', description: '全链路HTTPS加密，SSL证书一键部署', icon: 'Monitor' },
  { title: '视频加速', description: '大文件下载、视频点播/直播流畅播放', icon: 'VideoPlay' },
  { title: '边缘计算', description: '边缘节点自定义逻辑，灵活处理请求', icon: 'DataLine' },
  { title: '数据分析', description: '实时流量分析、命中率统计、日志下载', icon: 'TrendCharts' }
]

interface ProductPlan {
  id: number
  name: string
  description: string
  price: number
  featured: boolean
  features: string[]
}

const plans = ref<ProductPlan[]>([
  { id: 301, name: '基础版', description: '适合个人站点、小型博客', price: 59, featured: false, features: ['100GB流量/月', '50+节点', 'HTTP/HTTPS', '基础报表', '5个域名'] },
  { id: 302, name: '专业版', description: '适合中小企业官网、商城', price: 199, featured: false, features: ['500GB流量/月', '100+节点', '全协议加速', '详细报表', '20个域名', 'WAF防护'] },
  { id: 303, name: '企业版', description: '适合大型平台、视频网站', price: 599, featured: true, features: ['2TB流量/月', '200+节点', '全协议加速', '实时分析', '无限域名', 'WAF防护', '边缘计算'] },
  { id: 304, name: '旗舰版', description: '适合全球化业务、定制需求', price: 1499, featured: false, features: ['10TB流量/月', '全球节点', '全协议加速', '专属架构师', '无限域名', '全套安全', '定制方案'] }
])

const faqs = [
  { question: '什么是CDN？', answer: 'CDN（内容分发网络）通过在全球部署的边缘节点缓存您的内容，使用户可以从距离最近的节点获取资源，从而大幅提升访问速度和用户体验。' },
  { question: 'CDN加速适用于哪些场景？', answer: 'CDN适用于网站加速、文件下载、视频点播/直播、API加速、游戏更新等多种场景。凡是有大量用户访问的内容分发需求，都可以使用CDN加速。' },
  { question: '流量用完后会怎样？', answer: '当月流量用完后，您可以选择升级套餐或按量付费继续使用。系统会在流量即将耗尽时发送提醒，确保您的业务不受影响。' },
  { question: '如何接入CDN加速？', answer: '接入非常简单：1.添加加速域名；2.配置CNAME解析到CDN提供的域名；3.根据需要配置缓存规则。整个过程通常在10分钟内完成。' }
]

const scrollToProducts = () => {
  document.getElementById('products')?.scrollIntoView({ behavior: 'smooth' })
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/products', { params: { group: 'cdn' } })
    if (res.data?.data?.length) {
      plans.value = res.data.data
    }
  } catch (error) {
    console.error('获取数据失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>

<style scoped lang="scss">
.cdn-page {
  min-height: 100vh;
  background: #f8fafc;
}

.hero-section {
  background: linear-gradient(135deg, #10b981 0%, #059669 50%, #047857 100%);
  padding: 120px 20px 60px;
  text-align: center;
  position: relative;

  .hero-content {
    max-width: 700px;
    margin: 0 auto 40px;
  }

  .hero-badge {
    display: inline-block;
    padding: 6px 16px;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 20px;
    color: #fff;
    font-size: 14px;
    margin-bottom: 20px;
  }

  .hero-title {
    font-size: 48px;
    font-weight: 700;
    color: #fff;
    margin: 0 0 16px;
  }

  .hero-desc {
    font-size: 18px;
    color: rgba(255, 255, 255, 0.85);
    margin: 0 0 32px;
  }

  .hero-actions {
    display: flex;
    gap: 16px;
    justify-content: center;
  }

  .hero-ghost-btn {
    background: rgba(255, 255, 255, 0.2) !important;
    border-color: rgba(255, 255, 255, 0.4) !important;
    color: #fff !important;
  }

  .hero-features {
    display: flex;
    gap: 40px;
    justify-content: center;
    flex-wrap: wrap;
  }

  .hero-feature {
    display: flex;
    align-items: center;
    gap: 8px;
    color: rgba(255, 255, 255, 0.9);
    font-size: 15px;
  }
}

.section {
  padding: 80px 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.section-title {
  font-size: 32px;
  font-weight: 600;
  color: #0f172a;
  margin: 0 0 12px;

  &.text-center {
    text-align: center;
  }
}

.section-subtitle {
  font-size: 16px;
  color: #64748b;
  margin: 0 0 48px;

  &.text-center {
    text-align: center;
  }
}

// 优势
.advantages-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;

  @media (max-width: 768px) {
    grid-template-columns: repeat(2, 1fr);
  }

  @media (max-width: 480px) {
    grid-template-columns: 1fr;
  }
}

.advantage-card {
  background: #fff;
  border-radius: 16px;
  padding: 28px;
  text-align: center;
  transition: all 0.3s;
  border: 1px solid #e2e8f0;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.08);
  }

  h3 {
    font-size: 18px;
    font-weight: 600;
    color: #0f172a;
    margin: 16px 0 8px;
  }

  p {
    font-size: 14px;
    color: #64748b;
    margin: 0;
  }
}

// 产品列表
.products-section {
  background: #fff;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
  align-items: start;

  @media (max-width: 992px) {
    grid-template-columns: repeat(2, 1fr);
  }

  @media (max-width: 576px) {
    grid-template-columns: 1fr;
  }
}

.product-card {
  background: #f8fafc;
  border-radius: 20px;
  padding: 32px;
  text-align: center;
  position: relative;
  transition: all 0.3s;
  border: 2px solid #e2e8f0;

  &.featured {
    border-color: #10b981;
    background: #fff;
    transform: scale(1.05);
    box-shadow: 0 16px 40px rgba(16, 185, 129, 0.15);
  }

  .product-badge {
    position: absolute;
    top: -12px;
    left: 50%;
    transform: translateX(-50%);
    padding: 4px 20px;
    background: linear-gradient(135deg, #10b981, #34d399);
    color: #fff;
    border-radius: 12px;
    font-size: 13px;
    font-weight: 600;
  }

  h3 {
    font-size: 22px;
    font-weight: 600;
    color: #0f172a;
    margin: 8px 0 8px;
  }

  .product-desc {
    font-size: 14px;
    color: #64748b;
    margin: 0 0 20px;
  }
}

.product-price {
  margin-bottom: 24px;

  .price-symbol {
    font-size: 16px;
    color: #10b981;
  }

  .price-value {
    font-size: 40px;
    font-weight: 700;
    color: #10b981;
  }

  .price-unit {
    font-size: 14px;
    color: #94a3b8;
  }
}

.product-features {
  list-style: none;
  padding: 0;
  margin: 0 0 24px;
  text-align: left;

  li {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    color: #334155;
    padding: 6px 0;

    .el-icon {
      color: #10b981;
    }
  }
}

// 产品特点
.features-section {
  background: #fff;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;

  @media (max-width: 768px) {
    grid-template-columns: repeat(2, 1fr);
  }

  @media (max-width: 480px) {
    grid-template-columns: 1fr;
  }
}

.feature-card {
  background: #f8fafc;
  border-radius: 16px;
  padding: 24px;
  text-align: center;
  transition: all 0.3s;
  border: 1px solid #e2e8f0;

  &:hover {
    border-color: #10b981;
    box-shadow: 0 8px 24px rgba(16, 185, 129, 0.08);
  }

  h3 {
    font-size: 16px;
    font-weight: 600;
    color: #0f172a;
    margin: 12px 0 8px;
  }

  p {
    font-size: 13px;
    color: #64748b;
    margin: 0;
  }
}

// FAQ
.faq-section {
  background: #fff;
}

:deep(.el-collapse-item__header) {
  font-size: 16px;
  font-weight: 600;
  color: #0f172a;
}

:deep(.el-collapse-item__content) {
  font-size: 15px;
  color: #475569;
  line-height: 1.6;
}

// CTA
.cta-section {
  background: linear-gradient(135deg, #10b981 0%, #059669 50%, #047857 100%);

  .cta-content {
    text-align: center;
    color: #fff;

    h2 {
      font-size: 36px;
      font-weight: 700;
      margin: 0 0 16px;
    }

    p {
      font-size: 18px;
      opacity: 0.9;
      margin: 0 0 32px;
    }
  }

  .cta-actions {
    display: flex;
    gap: 16px;
    justify-content: center;
  }

  .cta-ghost-btn {
    background: rgba(255, 255, 255, 0.2) !important;
    border-color: rgba(255, 255, 255, 0.4) !important;
    color: #fff !important;
  }
}
</style>
