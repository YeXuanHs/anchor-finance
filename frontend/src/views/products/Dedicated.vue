<template>
  <div class="dedicated-page">
    <!-- Hero Banner -->
    <section class="hero-section">
      <div class="hero-content">
        <div class="hero-badge">物理服务器租用</div>
        <h1 class="hero-title">独立服务器</h1>
        <p class="hero-desc">独享物理资源、极致性能、完全控制权，满足高负载业务需求</p>
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
        <p class="section-subtitle text-center">为什么选择独立服务器</p>
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
        <h2 class="section-title text-center">独立服务器套餐</h2>
        <p class="section-subtitle text-center">高性能物理服务器，独享资源</p>
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
        <p class="section-subtitle text-center">独立服务器适用于多种业务场景</p>
        <div class="scenarios-grid">
          <div v-for="scenario in scenarios" :key="scenario.title" class="scenario-card">
            <el-icon :size="28" color="#f59e0b"><component :is="scenario.icon" /></el-icon>
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
        <p class="section-subtitle text-center">关于独立服务器的常见疑问</p>
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
          <h2>需要定制化配置？</h2>
          <p>联系我们获取专属服务器方案，满足您的特殊需求</p>
          <div class="cta-actions">
            <el-button type="primary" size="large" round @click="scrollToProducts">查看套餐</el-button>
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
  CircleCheck, Cpu, Monitor, Lightning, Shield, Timer,
  Connection, Headset, OfficeBuilding, ShoppingCart, VideoCamera, DataLine
} from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(true)

const heroFeatures = [
  { icon: 'Cpu', label: '独享资源' },
  { icon: 'Shield', label: '物理隔离' },
  { icon: 'Lightning', label: '极致性能' },
  { icon: 'Setting', label: '完全控制' }
]

const advantages = [
  { title: '独享资源', description: 'CPU、内存、硬盘完全独享，不受其他用户影响', icon: 'Cpu', color: '#f59e0b' },
  { title: '极致性能', description: '企业级硬件配置，满足高并发、大计算需求', icon: 'Lightning', color: '#10b981' },
  { title: '完全控制', description: 'root/admin权限，自由安装软件和配置环境', icon: 'Setting', color: '#3b82f6' },
  { title: '物理隔离', description: '物理服务器独立运行，数据安全更有保障', icon: 'Shield', color: '#8b5cf6' },
  { title: '大带宽', description: '支持独享大带宽，适合高流量业务', icon: 'Connection', color: '#ef4444' },
  { title: '专属支持', description: '专属技术经理一对一服务，7×24小时响应', icon: 'Headset', color: '#0ea5e9' }
]

const scenarios = [
  { title: '大型网站', description: '高流量门户网站、社区论坛', icon: 'Monitor' },
  { title: '电商平台', description: '大型电商、跨境电商平台', icon: 'ShoppingCart' },
  { title: '游戏服务', description: '游戏服务器、游戏加速', icon: 'VideoCamera' },
  { title: '大数据', description: '数据采集、存储和分析', icon: 'DataLine' }
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
  { id: 201, name: '入门级', description: '适合中小型业务', price: 599, featured: false, features: ['E3-1230v6', '16GB DDR4', '1TB SSD', '10Mbps独享', '5个IP'] },
  { id: 202, name: '标准级', description: '适合中大型业务', price: 999, featured: true, features: ['E5-2680v4', '32GB DDR4', '2TB SSD', '30Mbps独享', '10个IP'] },
  { id: 203, name: '企业级', description: '适合大型企业应用', price: 1999, featured: false, features: ['双路E5-2690v4', '128GB DDR4', '4TB SSD', '100Mbps独享', '20个IP'] },
  { id: 204, name: '旗舰级', description: '适合超大规模业务', price: 3999, featured: false, features: ['双路Platinum 8280', '256GB DDR4', '8TB SSD', '1Gbps独享', '30个IP'] }
])

const faqs = [
  { question: '独立服务器和云服务器有什么区别？', answer: '独立服务器是物理服务器，资源完全独享，性能更稳定；云服务器是虚拟化的，资源共享但弹性更好。如果您的业务需要稳定高性能和完全控制权，建议选择独立服务器。' },
  { question: '独立服务器支持哪些操作系统？', answer: '支持所有主流操作系统，包括CentOS、Ubuntu、Debian、Windows Server等。您也可以提供自定义镜像进行安装。' },
  { question: '服务器开通需要多长时间？', answer: '标准配置通常在24小时内开通完成。定制化配置可能需要1-3个工作日。具体时间请咨询客服确认。' },
  { question: '是否支持硬件升级？', answer: '支持。您可以在租用期间随时升级CPU、内存、硬盘等硬件配置。升级费用按剩余租期计算。' }
]

const scrollToProducts = () => {
  document.getElementById('products')?.scrollIntoView({ behavior: 'smooth' })
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v2/products', { params: { group: 'dedicated' } })
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
.dedicated-page {
  min-height: 100vh;
  background: #f8fafc;
}

.hero-section {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 50%, #b45309 100%);
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
    border-color: #f59e0b;
    background: #fff;
    transform: scale(1.05);
    box-shadow: 0 16px 40px rgba(245, 158, 11, 0.15);
  }

  .product-badge {
    position: absolute;
    top: -12px;
    left: 50%;
    transform: translateX(-50%);
    padding: 4px 20px;
    background: linear-gradient(135deg, #f59e0b, #fbbf24);
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
    color: #f59e0b;
  }

  .price-value {
    font-size: 40px;
    font-weight: 700;
    color: #f59e0b;
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
    border-color: #f59e0b;
    box-shadow: 0 8px 24px rgba(245, 158, 11, 0.08);
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
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 50%, #b45309 100%);

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
