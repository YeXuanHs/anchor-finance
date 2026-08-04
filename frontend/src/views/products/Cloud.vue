<template>
  <div class="cloud-page">
    <!-- Hero Banner -->
    <section class="hero-section">
      <div class="hero-content">
        <div class="hero-badge">云计算服务</div>
        <h1 class="hero-title">云服务器</h1>
        <p class="hero-desc">高性能、弹性扩展、安全可靠的云服务器，助力您的业务快速成长</p>
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
        <p class="section-subtitle text-center">为什么选择我们的云服务器</p>
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
        <h2 class="section-title text-center">云服务器套餐</h2>
        <p class="section-subtitle text-center">灵活配置，满足不同业务需求</p>
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
        <p class="section-subtitle text-center">全方位满足您的业务需求</p>
        <div class="features-grid">
          <div v-for="feature in productFeatures" :key="feature.title" class="feature-card">
            <el-icon :size="28" color="#3b82f6"><component :is="feature.icon" /></el-icon>
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
        <p class="section-subtitle text-center">关于云服务器的常见疑问</p>
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
          <h2>准备好开始了吗？</h2>
          <p>选择适合您的云服务器套餐，立即开启云端之旅</p>
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
import {
  CircleCheck, Monitor, Cpu, Lightning, Shield, Timer,
  Connection, Setting, DataLine, TrendCharts, Headset
} from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(true)

const heroFeatures = [
  { icon: 'Cpu', label: '高性能CPU' },
  { icon: 'Lightning', label: '秒级开通' },
  { icon: 'Shield', label: 'DDoS防护' },
  { icon: 'Connection', label: '多线BGP' }
]

const advantages = [
  { title: '高性能计算', description: '采用最新一代Intel/AMD处理器，NVMe SSD存储', icon: 'Cpu', color: '#3b82f6' },
  { title: '弹性扩展', description: 'CPU、内存、带宽按需升降，灵活应对业务变化', icon: 'TrendCharts', color: '#10b981' },
  { title: '安全可靠', description: '多层安全防护，免费DDoS防护，数据自动备份', icon: 'Shield', color: '#f59e0b' },
  { title: '快速部署', description: '一键部署常见应用环境，分钟级开通服务器', icon: 'Lightning', color: '#8b5cf6' },
  { title: '全球节点', description: '覆盖全球多个数据中心，就近接入极速体验', icon: 'Connection', color: '#ef4444' },
  { title: '7×24支持', description: '专业技术团队全天候在线，工单/IM快速响应', icon: 'Headset', color: '#0ea5e9' }
]

const productFeatures = [
  { title: '快照备份', description: '支持手动/自动快照，一键回滚数据', icon: 'Timer' },
  { title: '安全组', description: '灵活配置入站/出站规则，精细化访问控制', icon: 'Shield' },
  { title: '负载均衡', description: '自动分发流量，提高应用可用性', icon: 'DataLine' },
  { title: '弹性IP', description: 'IP地址灵活绑定和解绑，无缝迁移', icon: 'Connection' },
  { title: '监控告警', description: '实时监控资源使用，异常自动告警', icon: 'Monitor' },
  { title: 'API接口', description: '丰富的API接口，支持自动化运维', icon: 'Setting' }
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
  { id: 101, name: '入门型', description: '适合个人博客、小型网站', price: 49, featured: false, features: ['1核CPU', '1GB内存', '40GB SSD', '1Mbps带宽', '免费DDoS防护'] },
  { id: 102, name: '基础型', description: '适合中小企业官网、论坛', price: 99, featured: false, features: ['2核CPU', '4GB内存', '80GB SSD', '3Mbps带宽', '免费DDoS防护'] },
  { id: 103, name: '进阶型', description: '适合电商平台、应用服务', price: 199, featured: true, features: ['4核CPU', '8GB内存', '160GB SSD', '5Mbps带宽', '免费DDoS防护', '负载均衡'] },
  { id: 104, name: '专业型', description: '适合大型应用、游戏服务', price: 399, featured: false, features: ['8核CPU', '16GB内存', '320GB SSD', '10Mbps带宽', '高级DDoS防护', '专属客服'] }
])

const faqs = [
  { question: '云服务器和传统服务器有什么区别？', answer: '云服务器基于云计算技术，具有弹性扩展、按需付费、快速部署等优势。相比传统服务器，云服务器可以随时调整配置，无需担心硬件故障，且成本更低。' },
  { question: '如何选择合适的云服务器配置？', answer: '建议根据您的业务类型和访问量选择。个人博客/小型网站选择入门型即可；企业官网/论坛建议基础型；电商平台/应用服务建议进阶型；大型应用建议专业型或更高配置。' },
  { question: '云服务器支持哪些操作系统？', answer: '支持主流的Linux发行版（CentOS、Ubuntu、Debian等）和Windows Server系统。您也可以选择预装环境的镜像，如LAMP、LNMP、Docker等。' },
  { question: '数据安全如何保障？', answer: '我们提供多层次的数据安全保障：RAID磁盘阵列、自动快照备份、免费DDoS防护、安全组规则配置等。建议您同时做好应用层面的安全防护。' }
]

const scrollToProducts = () => {
  document.getElementById('products')?.scrollIntoView({ behavior: 'smooth' })
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v2/products', { params: { group: 'cloud' } })
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
.cloud-page {
  min-height: 100vh;
  background: #f8fafc;
}

.hero-section {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 50%, #1d4ed8 100%);
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
    border-color: #3b82f6;
    background: #fff;
    transform: scale(1.05);
    box-shadow: 0 16px 40px rgba(59, 130, 246, 0.15);
  }

  .product-badge {
    position: absolute;
    top: -12px;
    left: 50%;
    transform: translateX(-50%);
    padding: 4px 20px;
    background: linear-gradient(135deg, #3b82f6, #60a5fa);
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
    color: #3b82f6;
  }

  .price-value {
    font-size: 40px;
    font-weight: 700;
    color: #3b82f6;
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
    border-color: #3b82f6;
    box-shadow: 0 8px 24px rgba(59, 130, 246, 0.08);
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
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 50%, #1d4ed8 100%);

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
