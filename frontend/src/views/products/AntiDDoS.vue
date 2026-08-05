<template>
  <div class="antiddos-page">
    <!-- Hero Banner -->
    <section class="hero-section">
      <div class="hero-content">
        <div class="hero-badge">安全防护服务</div>
        <h1 class="hero-title">Anti-DDoS 高防</h1>
        <p class="hero-desc">全方位DDoS防护能力，保障业务稳定运行，抵御大流量攻击</p>
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
        <p class="section-subtitle text-center">为什么选择我们的DDoS防护</p>
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
        <h2 class="section-title text-center">高防套餐</h2>
        <p class="section-subtitle text-center">灵活防护配置，满足不同规模业务需求</p>
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
        <p class="section-subtitle text-center">全面的DDoS防护解决方案</p>
        <div class="features-grid">
          <div v-for="feature in productFeatures" :key="feature.title" class="feature-card">
            <el-icon :size="28" color="#ef4444"><component :is="feature.icon" /></el-icon>
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
        <p class="section-subtitle text-center">关于DDoS防护的常见疑问</p>
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
          <h2>保护您的业务安全</h2>
          <p>选择适合的高防套餐，让业务远离DDoS攻击威胁</p>
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
  CircleCheck, Shield, Lightning, Monitor, Warning,
  Lock, DataLine, TrendCharts, Headset, View
} from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(true)

const heroFeatures = [
  { icon: 'Shield', label: 'T级防护' },
  { icon: 'Lightning', label: '秒级切换' },
  { icon: 'View', label: '实时监控' },
  { icon: 'Lock', label: '零误封' }
]

const advantages = [
  { title: '超大防护容量', description: '单节点T级防护能力，分布式清洗集群，轻松抵御大流量攻击', icon: 'Shield', color: '#ef4444' },
  { title: '智能识别清洗', description: 'AI驱动的流量分析引擎，精准识别攻击流量，毫秒级响应', icon: 'View', color: '#3b82f6' },
  { title: '零业务影响', description: '透明化防护模式，正常业务流量零损耗，访问延迟无感知', icon: 'Lightning', color: '#10b981' },
  { title: '多协议支持', description: '支持TCP/UDP/HTTP/HTTPS等协议防护，覆盖各类业务场景', icon: 'Monitor', color: '#f59e0b' },
  { title: '实时监控报告', description: '攻击流量实时可视化，详细攻击报告自动生成，防护状态一目了然', icon: 'DataLine', color: '#8b5cf6' },
  { title: '7×24专家值守', description: '资深安全专家团队全天候值守，攻击发生时即时响应处理', icon: 'Headset', color: '#0ea5e9' }
]

const productFeatures = [
  { title: 'CC防护', description: '智能CC攻击识别与防护，保障Web业务稳定', icon: 'Warning' },
  { title: '流量清洗', description: '多层级流量清洗机制，精准过滤攻击流量', icon: 'Shield' },
  { title: '黑名单管理', description: '灵活的IP黑白名单策略，精细化访问控制', icon: 'Lock' },
  { title: '弹性扩展', description: '防护能力按需升级，灵活应对突发攻击', icon: 'TrendCharts' },
  { title: '攻击告警', description: '攻击事件实时推送，多渠道告警通知', icon: 'Monitor' },
  { title: '数据报表', description: '详细的防护数据报表，攻击趋势分析', icon: 'DataLine' }
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
  { id: 201, name: '基础防护', description: '适合个人站点、小型应用', price: 299, featured: false, features: ['10Gbps防护', 'TCP/UDP防护', 'CC基础防护', '实时监控', '5个域名'] },
  { id: 202, name: '高级防护', description: '适合企业官网、电商平台', price: 799, featured: false, features: ['50Gbps防护', '全协议防护', 'CC高级防护', '攻击报告', '20个域名'] },
  { id: 203, name: '旗舰防护', description: '适合游戏、金融等高安全需求', price: 1999, featured: true, features: ['200Gbps防护', '全协议防护', 'CC无限防护', '专属安全顾问', '无限域名', '秒级切换'] },
  { id: 204, name: '定制防护', description: '适合大型企业、特殊业务场景', price: 4999, featured: false, features: ['T级防护', '全协议防护', '定制清洗策略', '7×24专家值守', '无限域名', '专属集群'] }
])

const faqs = [
  { question: '什么是DDoS攻击？', answer: 'DDoS（分布式拒绝服务攻击）是通过大量恶意流量占用目标服务器资源，导致正常用户无法访问的攻击方式。攻击者通常利用僵尸网络发起大规模流量冲击。' },
  { question: '高防服务器如何防护DDoS攻击？', answer: '我们通过分布式清洗集群对流量进行实时分析和过滤，利用AI算法精准识别攻击流量并将其清洗，只将正常流量回源到您的服务器，确保业务不受影响。' },
  { question: '防护是否会增加访问延迟？', answer: '我们的高防节点采用BGP多线接入，清洗延迟在毫秒级别，对正常用户访问几乎无感知。同时支持智能回源，优化访问路径。' },
  { question: '被攻击时如何处理？', answer: '当检测到攻击时，系统会自动切换到高防模式进行流量清洗。同时我们会在第一时间通知您，并提供详细的攻击报告。如有需要，安全专家会协助您调整防护策略。' }
]

const scrollToProducts = () => {
  document.getElementById('products')?.scrollIntoView({ behavior: 'smooth' })
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/products', { params: { group: 'antiddos' } })
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
.antiddos-page {
  min-height: 100vh;
  background: #f8fafc;
}

.hero-section {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 50%, #b91c1c 100%);
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
    border-color: #ef4444;
    background: #fff;
    transform: scale(1.05);
    box-shadow: 0 16px 40px rgba(239, 68, 68, 0.15);
  }

  .product-badge {
    position: absolute;
    top: -12px;
    left: 50%;
    transform: translateX(-50%);
    padding: 4px 20px;
    background: linear-gradient(135deg, #ef4444, #f87171);
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
    color: #ef4444;
  }

  .price-value {
    font-size: 40px;
    font-weight: 700;
    color: #ef4444;
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
    border-color: #ef4444;
    box-shadow: 0 8px 24px rgba(239, 68, 68, 0.08);
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
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 50%, #b91c1c 100%);

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
