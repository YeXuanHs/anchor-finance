<template>
  <div class="hosting-page">
    <!-- Hero Banner -->
    <section class="hero-section">
      <div class="hero-content">
        <div class="hero-badge">网站托管服务</div>
        <h1 class="hero-title">虚拟主机</h1>
        <p class="hero-desc">一键建站、操作简单、性价比高，适合个人和中小企业网站</p>
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
        <p class="section-subtitle text-center">为什么选择我们的虚拟主机</p>
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
        <h2 class="section-title text-center">虚拟主机套餐</h2>
        <p class="section-subtitle text-center">多种配置可选，满足不同建站需求</p>
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
                  <span class="price-unit">/年</span>
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

    <!-- 支持环境 -->
    <section class="section environments-section">
      <div class="container">
        <h2 class="section-title text-center">支持环境</h2>
        <p class="section-subtitle text-center">主流建站环境全面支持</p>
        <div class="env-grid">
          <div v-for="env in environments" :key="env.name" class="env-card">
            <h3>{{ env.name }}</h3>
            <p>{{ env.description }}</p>
            <div class="env-tags">
              <el-tag v-for="tag in env.tags" :key="tag" size="small">{{ tag }}</el-tag>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 常见问题 -->
    <section class="section faq-section">
      <div class="container">
        <h2 class="section-title text-center">常见问题</h2>
        <p class="section-subtitle text-center">关于虚拟主机的常见疑问</p>
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
          <h2>轻松建站，从这里开始</h2>
          <p>选择适合您的虚拟主机套餐，快速搭建您的网站</p>
          <div class="cta-actions">
            <el-button type="primary" size="large" round @click="scrollToProducts">立即选购</el-button>
            <el-button size="large" round class="cta-ghost-btn" @click="$router.push('/contact')">联系客服</el-button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  CircleCheck, Monitor, Lightning, Shield, Timer,
  Connection, Setting, Headset
} from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(true)

const heroFeatures = [
  { icon: 'Lightning', label: '一键部署' },
  { icon: 'Shield', label: '安全防护' },
  { icon: 'Timer', label: '99.9%可用性' },
  { icon: 'Headset', label: '7×24支持' }
]

const advantages = [
  { title: '一键建站', description: '预装WordPress、Discuz等主流程序，5分钟快速建站', icon: 'Lightning', color: '#10b981' },
  { title: '操作简单', description: '可视化控制面板，无需技术基础即可管理网站', icon: 'Monitor', color: '#3b82f6' },
  { title: '安全可靠', description: 'Web应用防火墙、防篡改、自动备份', icon: 'Shield', color: '#f59e0b' },
  { title: '高性价比', description: '资源共享，成本更低，适合中小型网站', icon: 'Document', color: '#8b5cf6' },
  { title: '稳定高速', description: 'SSD存储、CDN加速，网站访问更快', icon: 'Connection', color: '#ef4444' },
  { title: '专业支持', description: '7×24小时技术支持，协助解决建站问题', icon: 'Headset', color: '#0ea5e9' }
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
  { id: 301, name: '体验型', description: '适合个人博客', price: 98, featured: false, features: ['1GB空间', '10GB月流量', '1个站点', 'MySQL数据库', '免费SSL'] },
  { id: 302, name: '基础型', description: '适合个人/小型企业', price: 198, featured: false, features: ['5GB空间', '50GB月流量', '3个站点', 'MySQL数据库', '免费SSL', 'CDN加速'] },
  { id: 303, name: '标准型', description: '适合中小企业官网', price: 398, featured: true, features: ['20GB空间', '200GB月流量', '10个站点', 'MySQL/MSSQL', '免费SSL', 'CDN加速', '独立IP'] },
  { id: 304, name: '企业型', description: '适合中大型企业', price: 798, featured: false, features: ['50GB空间', '不限流量', '不限站点', 'MySQL/MSSQL', '免费SSL', 'CDN加速', '独立IP', '专属客服'] }
])

const environments = [
  { name: 'PHP环境', description: '支持PHP 5.6-8.2多版本切换', tags: ['PHP', 'Laravel', 'ThinkPHP', 'WordPress'] },
  { name: 'Java环境', description: '支持Java/Tomcat运行环境', tags: ['Java', 'Tomcat', 'Spring', 'Maven'] },
  { name: 'Node.js环境', description: '支持Node.js应用部署', tags: ['Node.js', 'Express', 'Koa', 'Nuxt.js'] },
  { name: 'ASP.NET环境', description: '支持ASP.NET应用运行', tags: ['.NET', 'C#', 'MVC', 'SQL Server'] }
]

const faqs = [
  { question: '虚拟主机和云服务器有什么区别？', answer: '虚拟主机是共享服务器资源，操作简单、成本低，适合中小型网站；云服务器是独享资源，可自由配置，适合需要定制化环境的用户。如果您是新手或网站流量不大，建议选择虚拟主机。' },
  { question: '虚拟主机支持哪些程序？', answer: '支持WordPress、Discuz、DedeCMS、帝国CMS、Joomla、Drupal等主流CMS程序。我们提供一键安装功能，让您快速部署网站。' },
  { question: '如何备份网站数据？', answer: '我们提供自动备份功能，每天自动备份网站数据。您也可以在控制面板中手动备份和恢复数据。建议定期下载备份到本地保存。' },
  { question: '虚拟主机支持SSL证书吗？', answer: '支持。我们提供免费SSL证书，您可以在控制面板中一键开启HTTPS。也可以上传自己的SSL证书。' }
]

const scrollToProducts = () => {
  document.getElementById('products')?.scrollIntoView({ behavior: 'smooth' })
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/products', { params: { group: 'hosting' } })
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
.hosting-page {
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

.environments-section {
  background: #fff;
}

.env-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.env-card {
  background: #f8fafc;
  border-radius: 16px;
  padding: 28px;
  border: 1px solid #e2e8f0;
  transition: all 0.3s;

  &:hover {
    border-color: #10b981;
    box-shadow: 0 8px 24px rgba(16, 185, 129, 0.08);
  }

  h3 {
    font-size: 18px;
    font-weight: 600;
    color: #0f172a;
    margin: 0 0 8px;
  }

  p {
    font-size: 14px;
    color: #64748b;
    margin: 0 0 16px;
  }

  .env-tags {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;

    .el-tag {
      background: #ecfdf5;
      color: #10b981;
      border-color: #a7f3d0;
    }
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
