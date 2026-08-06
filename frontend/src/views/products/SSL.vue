<template>
  <div class="ssl-page">
    <!-- Hero Banner -->
    <section class="hero-section">
      <div class="hero-content">
        <div class="hero-badge">SSL证书服务</div>
        <h1 class="hero-title">SSL证书</h1>
        <p class="hero-desc">保护网站数据安全，提升用户信任度，HTTPS加密必备</p>
        <div class="hero-actions">
          <el-button type="primary" size="large" round @click="scrollToProducts">立即选购</el-button>
          <el-button size="large" round class="hero-ghost-btn" @click="$router.push('/contact')">联系我们</el-button>
        </div>
      </div>
    </section>

    <!-- 证书类型 -->
    <section id="products" class="section types-section">
      <div class="container">
        <h2 class="section-title text-center">证书类型</h2>
        <p class="section-subtitle text-center">根据您的需求选择合适的证书</p>
        <el-skeleton :loading="loading" animated :rows="6">
          <template #default>
            <div class="types-grid">
              <div
                v-for="cert in certTypes"
                :key="cert.id"
                class="type-card"
                :class="{ featured: cert.featured }"
              >
                <div class="type-badge" v-if="cert.featured">推荐</div>
                <div class="type-icon" :style="{ background: cert.gradient }">
                  <el-icon :size="32" color="#fff"><Lock /></el-icon>
                </div>
                <h3>{{ cert.name }}</h3>
                <p class="type-desc">{{ cert.description }}</p>
                <div class="type-price">
                  <span class="price-value">¥{{ cert.price }}</span>
                  <span class="price-unit">/年</span>
                </div>
                <ul class="type-features">
                  <li v-for="feature in cert.features" :key="feature">
                    <el-icon><CircleCheck /></el-icon>
                    {{ feature }}
                  </li>
                </ul>
                <el-button
                  :type="cert.featured ? 'primary' : 'default'"
                  round
                  @click="$router.push(`/products/${cert.id}`)"
                >
                  立即选购
                </el-button>
              </div>
            </div>
          </template>
        </el-skeleton>
      </div>
    </section>

    <!-- 核心优势 -->
    <section class="section advantages-section">
      <div class="container">
        <h2 class="section-title text-center">核心优势</h2>
        <p class="section-subtitle text-center">为什么选择我们的SSL证书</p>
        <div class="advantages-grid">
          <div v-for="adv in advantages" :key="adv.title" class="advantage-card">
            <el-icon :size="32" :color="adv.color"><component :is="adv.icon" /></el-icon>
            <h3>{{ adv.title }}</h3>
            <p>{{ adv.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 常见问题 -->
    <section class="section faq-section">
      <div class="container">
        <h2 class="section-title text-center">常见问题</h2>
        <p class="section-subtitle text-center">关于SSL证书的常见疑问</p>
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
          <h2>保护您的网站安全</h2>
          <p>选择适合您的SSL证书，立即开启HTTPS加密</p>
          <div class="cta-actions">
            <el-button type="primary" size="large" round @click="scrollToProducts">立即选购</el-button>
            <el-button size="large" round class="cta-ghost-btn" @click="$router.push('/contact')">联系我们</el-button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Lock, CircleCheck, Warning, Timer, Setting, Service } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(true)

interface CertType {
  id: number
  name: string
  description: string
  price: number
  gradient: string
  featured: boolean
  features: string[]
}

const certTypes = ref<CertType[]>([
  {
    id: 501,
    name: 'DV证书（免费）',
    description: '域名验证证书，适合个人网站',
    price: 0,
    gradient: 'linear-gradient(135deg, #10b981, #34d399)',
    featured: false,
    features: ['域名验证', '快速签发', '浏览器信任', '基础加密']
  },
  {
    id: 502,
    name: 'OV证书',
    description: '组织验证证书，适合企业网站',
    price: 299,
    gradient: 'linear-gradient(135deg, #3b82f6, #60a5fa)',
    featured: true,
    features: ['组织验证', '企业信息展示', '更高信任度', '256位加密', '支持通配符']
  },
  {
    id: 503,
    name: 'EV证书',
    description: '扩展验证证书，适合金融电商',
    price: 999,
    gradient: 'linear-gradient(135deg, #8b5cf6, #a78bfa)',
    featured: false,
    features: ['扩展验证', '绿色地址栏', '最高信任度', '256位加密', '支持多域名', '赔付保障']
  }
])

const advantages = [
  { title: '快速签发', description: 'DV证书分钟级签发，OV证书1-3个工作日', icon: 'Timer', color: '#3b82f6' },
  { title: '一键部署', description: '支持一键部署到CDN和服务器', icon: 'Setting', color: '#10b981' },
  { title: '安全可靠', description: '256位加密，保护数据传输安全', icon: 'Shield', color: '#f59e0b' },
  { title: '专业支持', description: '7×24小时技术支持，协助安装配置', icon: 'Service', color: '#8b5cf6' }
]

const faqs = [
  { question: '什么是SSL证书？', answer: 'SSL证书是一种数字证书，用于在用户浏览器和服务器之间建立加密连接，保护数据传输安全。安装SSL证书后，网站地址会显示HTTPS和安全锁标志。' },
  { question: 'DV、OV、EV证书有什么区别？', answer: 'DV证书只验证域名所有权，签发最快；OV证书需要验证组织信息，会在证书中显示企业信息；EV证书需要最严格的验证，会在浏览器地址栏显示绿色企业名称。' },
  { question: 'SSL证书多久能签发？', answer: 'DV证书通常在几分钟内签发；OV证书需要1-3个工作日；EV证书需要3-7个工作日。签发时间取决于验证速度。' },
  { question: '证书过期后怎么办？', answer: '证书过期前我们会发送续费提醒。您可以提前续费，也可以在过期后重新购买。过期后网站会显示不安全警告，建议及时续费。' }
]

const scrollToProducts = () => {
  document.getElementById('products')?.scrollIntoView({ behavior: 'smooth' })
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/ssl/certificates')
    if (res.data?.data?.length) {
      certTypes.value = res.data.data
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
.ssl-page {
  min-height: 100vh;
  background: #f8fafc;
}

.hero-section {
  background: linear-gradient(135deg, #14b8a6 0%, #0d9488 50%, #0f766e 100%);
  padding: 120px 20px 80px;
  text-align: center;

  .hero-content {
    max-width: 700px;
    margin: 0 auto;
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

.types-section {
  background: #fff;
}

.types-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  align-items: start;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.type-card {
  background: #f8fafc;
  border-radius: 20px;
  padding: 32px;
  text-align: center;
  position: relative;
  transition: all 0.3s;
  border: 2px solid #e2e8f0;

  &.featured {
    border-color: #14b8a6;
    background: #fff;
    transform: scale(1.05);
    box-shadow: 0 16px 40px rgba(20, 184, 166, 0.15);
  }

  .type-badge {
    position: absolute;
    top: -12px;
    left: 50%;
    transform: translateX(-50%);
    padding: 4px 20px;
    background: linear-gradient(135deg, #14b8a6, #2dd4bf);
    color: #fff;
    border-radius: 12px;
    font-size: 13px;
    font-weight: 600;
  }

  .type-icon {
    width: 64px;
    height: 64px;
    border-radius: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 8px auto 16px;
  }

  h3 {
    font-size: 22px;
    font-weight: 600;
    color: #0f172a;
    margin: 0 0 8px;
  }

  .type-desc {
    font-size: 14px;
    color: #64748b;
    margin: 0 0 20px;
  }
}

.type-price {
  margin-bottom: 24px;

  .price-value {
    font-size: 40px;
    font-weight: 700;
    color: #14b8a6;
  }

  .price-unit {
    font-size: 14px;
    color: #94a3b8;
  }
}

.type-features {
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
      color: #14b8a6;
    }
  }
}

.advantages-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;

  @media (max-width: 768px) {
    grid-template-columns: repeat(2, 1fr);
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
  background: linear-gradient(135deg, #14b8a6 0%, #0d9488 50%, #0f766e 100%);

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
