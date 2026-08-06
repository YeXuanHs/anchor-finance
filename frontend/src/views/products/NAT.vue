<template>
  <div class="nat-page">
    <!-- Hero Banner -->
    <section class="hero-section">
      <div class="hero-content">
        <div class="hero-badge">网络地址转换</div>
        <h1 class="hero-title">NAT网关</h1>
        <p class="hero-desc">安全高效的网络地址转换服务，轻松实现内网服务器访问互联网</p>
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
        <p class="section-subtitle text-center">为什么选择NAT网关</p>
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
        <h2 class="section-title text-center">NAT网关套餐</h2>
        <p class="section-subtitle text-center">灵活的规格配置，满足不同规模业务需求</p>
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

    <!-- 适用场景 -->
    <section class="section scenarios-section">
      <div class="container">
        <h2 class="section-title text-center">适用场景</h2>
        <p class="section-subtitle text-center">NAT网关适用于多种业务场景</p>
        <div class="scenarios-grid">
          <div v-for="scenario in scenarios" :key="scenario.title" class="scenario-card">
            <el-icon :size="28" color="#ef4444"><component :is="scenario.icon" /></el-icon>
            <h3>{{ scenario.title }}</h3>
            <p>{{ scenario.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 常见问题 -->
    <section class="section faq-section">
      <div class="container">
        <h2 class="section-title text-center">常见问题</h2>
        <p class="section-subtitle text-center">关于NAT网关的常见疑问</p>
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
          <h2>需要安全的网络出口？</h2>
          <p>选择适合您的NAT网关，轻松实现内网服务器安全访问互联网</p>
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
  CircleCheck
} from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(true)

const heroFeatures = [
  { icon: 'Connection', label: '安全隔离' },
  { icon: 'Lightning', label: '高性能' },
  { icon: 'Shield', label: 'DDoS防护' },
  { icon: 'DataLine', label: '实时监控' }
]

const advantages = [
  { title: '安全隔离', description: '内网服务器无需公网IP，通过NAT网关安全访问互联网', icon: 'Shield', color: '#ef4444' },
  { title: '高性能转发', description: '高吞吐量、低延迟的网络地址转换能力', icon: 'Lightning', color: '#10b981' },
  { title: '灵活配置', description: '支持端口转发、SNAT、DNAT等多种配置', icon: 'Setting', color: '#3b82f6' },
  { title: '实时监控', description: '流量、连接数、带宽使用率实时监控', icon: 'DataLine', color: '#f59e0b' },
  { title: '高可用性', description: '多节点冗余，自动故障切换，保障业务连续性', icon: 'Connection', color: '#8b5cf6' },
  { title: '专业支持', description: '7×24小时技术支持，协助配置和故障排查', icon: 'Headset', color: '#0ea5e9' }
]

const scenarios = [
  { title: 'Web服务', description: '内网Web服务器对外提供服务', icon: 'Monitor' },
  { title: '数据库访问', description: '内网数据库安全访问互联网', icon: 'Lock' },
  { title: 'API服务', description: '内网API服务对外发布', icon: 'Connection' },
  { title: '邮件服务', description: '内网邮件服务器收发邮件', icon: 'DataLine' }
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
  { id: 601, name: '小型', description: '适合小型业务', price: 99, featured: false, features: ['100Mbps带宽', '10000连接数', '5条转发规则', '基础监控'] },
  { id: 602, name: '中型', description: '适合中型业务', price: 299, featured: true, features: ['500Mbps带宽', '50000连接数', '20条转发规则', '高级监控', 'DDoS防护'] },
  { id: 603, name: '大型', description: '适合大型业务', price: 599, featured: false, features: ['1Gbps带宽', '100000连接数', '50条转发规则', '高级监控', 'DDoS防护', '专属IP'] },
  { id: 604, name: '企业级', description: '适合企业级应用', price: 999, featured: false, features: ['10Gbps带宽', '不限连接数', '不限规则', '高级监控', 'DDoS防护', '多专属IP', '专属客服'] }
])

const faqs = [
  { question: '什么是NAT网关？', answer: 'NAT网关（Network Address Translation Gateway）是一种网络地址转换服务，它允许多个内网服务器共享一个或多个公网IP地址访问互联网，同时也可以将互联网流量转发到内网服务器。' },
  { question: '为什么需要NAT网关？', answer: '使用NAT网关可以提高网络安全性（内网服务器无需暴露公网IP）、节省公网IP资源、灵活管理网络访问策略、以及提供统一的互联网出口。' },
  { question: 'NAT网关支持哪些协议？', answer: '支持TCP、UDP、ICMP等常见协议。可以根据端口号进行精确的流量转发控制。' },
  { question: 'NAT网关的性能如何？', answer: '我们提供不同规格的NAT网关，从100Mbps到10Gbps带宽可选，连接数从1万到无限制。您可以根据业务规模选择合适的规格。' }
]

const scrollToProducts = () => {
  document.getElementById('products')?.scrollIntoView({ behavior: 'smooth' })
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/products', { params: { group: 'nat' } })
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
.nat-page {
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

.scenarios-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;

  @media (max-width: 768px) {
    grid-template-columns: repeat(2, 1fr);
  }
}

.scenario-card {
  background: #fff;
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
