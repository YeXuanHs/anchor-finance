<template>
  <div class="domain-page">
    <!-- Hero Banner -->
    <section class="hero-section">
      <div class="hero-content">
        <div class="hero-badge">域名注册服务</div>
        <h1 class="hero-title">域名注册</h1>
        <p class="hero-desc">海量域名后缀可选，快速注册，安全可靠</p>
        <div class="search-box">
          <el-input
            v-model="domainName"
            placeholder="输入您想要的域名..."
            size="large"
            class="domain-input"
          >
            <template #append>
              <el-button type="primary" @click="searchDomain">
                <el-icon><Search /></el-icon>
                查询域名
              </el-button>
            </template>
          </el-input>
        </div>
      </div>
    </section>

    <!-- 热门后缀 -->
    <section class="section suffixes-section">
      <div class="container">
        <h2 class="section-title text-center">热门域名后缀</h2>
        <p class="section-subtitle text-center">选择适合您业务的域名后缀</p>
        <el-skeleton :loading="loading" animated :rows="4">
          <template #default>
            <div class="suffixes-grid">
              <div
                v-for="suffix in suffixes"
                :key="suffix.ext"
                class="suffix-card"
                @click="domainName = `example${suffix.ext}`"
              >
                <div class="suffix-ext">{{ suffix.ext }}</div>
                <div class="suffix-price">
                  <span class="price-value">¥{{ suffix.price }}</span>
                  <span class="price-unit">/年</span>
                </div>
                <p class="suffix-desc">{{ suffix.description }}</p>
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
        <p class="section-subtitle text-center">为什么选择我们注册域名</p>
        <div class="advantages-grid">
          <div v-for="adv in advantages" :key="adv.title" class="advantage-card">
            <el-icon :size="32" :color="adv.color"><component :is="adv.icon" /></el-icon>
            <h3>{{ adv.title }}</h3>
            <p>{{ adv.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 域名服务 -->
    <section class="section services-section">
      <div class="container">
        <h2 class="section-title text-center">域名服务</h2>
        <p class="section-subtitle text-center">全方位域名管理服务</p>
        <div class="services-grid">
          <div v-for="service in services" :key="service.title" class="service-card">
            <el-icon :size="28" color="#8b5cf6"><component :is="service.icon" /></el-icon>
            <h3>{{ service.title }}</h3>
            <p>{{ service.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 常见问题 -->
    <section class="section faq-section">
      <div class="container">
        <h2 class="section-title text-center">常见问题</h2>
        <p class="section-subtitle text-center">关于域名注册的常见疑问</p>
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
          <h2>注册您的专属域名</h2>
          <p>保护您的品牌，开启互联网之旅</p>
          <div class="cta-actions">
            <el-button type="primary" size="large" round @click="scrollToTop">查询域名</el-button>
            <el-button size="large" round class="cta-ghost-btn" @click="$router.push('/contact')">联系我们</el-button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  Search, Refresh, Shield, Transfer, Setting, Connection, Timer, Document
} from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(true)
const domainName = ref('')

interface DomainSuffix {
  ext: string
  price: number
  description: string
}

const suffixes = ref<DomainSuffix[]>([
  { ext: '.com', price: 69, description: '全球最流行的域名后缀' },
  { ext: '.cn', price: 29, description: '中国国家顶级域名' },
  { ext: '.net', price: 79, description: '网络服务相关域名' },
  { ext: '.org', price: 89, description: '组织机构域名' },
  { ext: '.io', price: 299, description: '科技创业公司首选' },
  { ext: '.co', price: 199, description: '公司和商业域名' },
  { ext: '.vip', price: 49, description: '尊贵身份标识' },
  { ext: '.shop', price: 129, description: '电商平台专属域名' }
])

const advantages = [
  { title: '快速注册', description: '域名注册即时生效，无需等待', icon: 'Timer', color: '#8b5cf6' },
  { title: '安全可靠', description: '域名锁定保护，防止被盗', icon: 'Shield', color: '#10b981' },
  { title: '免费DNS', description: '提供免费DNS解析服务', icon: 'Connection', color: '#f59e0b' },
  { title: '智能管理', description: '一站式域名管理平台', icon: 'Setting', color: '#3b82f6' }
]

const services = [
  { title: '域名解析', description: '智能DNS解析，支持多种记录类型', icon: 'Connection' },
  { title: '域名转入', description: '快速转入域名，享受优惠续费', icon: 'Refresh' },
  { title: '域名转移', description: '支持域名转移给其他账户', icon: 'Transfer' },
  { title: '隐私保护', description: 'WHOIS隐私保护，隐藏个人信息', icon: 'Shield' }
]

const faqs = [
  { question: '如何注册域名？', answer: '在搜索框中输入您想要的域名，点击查询。如果域名可用，选择注册年限并完成支付即可。注册完成后即可在控制面板中管理您的域名。' },
  { question: '域名注册后多久生效？', answer: '域名注册成功后即时生效。但由于DNS传播需要时间，全球各地可能需要24-48小时才能完全生效。' },
  { question: '域名到期后怎么办？', answer: '域名到期前我们会发送续费提醒。到期后有30天的续费宽限期，宽限期内仍可正常续费。超过宽限期域名将被释放，可能被他人注册。' },
  { question: '支持域名转入吗？', answer: '支持。您可以在原注册商获取转移密码，然后在我们的平台提交转入申请。转入过程通常需要5-7个工作日，转入成功后会自动续费一年。' }
]

const searchDomain = () => {
  if (!domainName.value) return
  // 跳转到域名查询结果页
}

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/domain/suffixes')
    if (res.data?.data?.length) {
      suffixes.value = res.data.data
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
.domain-page {
  min-height: 100vh;
  background: #f8fafc;
}

.hero-section {
  background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 50%, #6d28d9 100%);
  padding: 120px 20px 80px;
  text-align: center;

  .hero-content {
    max-width: 600px;
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

  .search-box {
    max-width: 500px;
    margin: 0 auto;
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

.suffixes-section {
  background: #fff;
}

.suffixes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 20px;
}

.suffix-card {
  background: #f8fafc;
  border-radius: 16px;
  padding: 24px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;

  &:hover {
    border-color: #8b5cf6;
    box-shadow: 0 8px 24px rgba(139, 92, 246, 0.1);
  }

  .suffix-ext {
    font-size: 24px;
    font-weight: 700;
    color: #0f172a;
    margin-bottom: 12px;
  }

  .suffix-price {
    margin-bottom: 8px;

    .price-value {
      font-size: 28px;
      font-weight: 700;
      color: #8b5cf6;
    }

    .price-unit {
      font-size: 14px;
      color: #94a3b8;
    }
  }

  .suffix-desc {
    font-size: 13px;
    color: #64748b;
    margin: 0;
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

.services-section {
  background: #fff;
}

.services-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;

  @media (max-width: 768px) {
    grid-template-columns: repeat(2, 1fr);
  }
}

.service-card {
  background: #f8fafc;
  border-radius: 16px;
  padding: 28px;
  text-align: center;
  transition: all 0.3s;
  border: 1px solid #e2e8f0;

  &:hover {
    border-color: #8b5cf6;
    box-shadow: 0 8px 24px rgba(139, 92, 246, 0.08);
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
  background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 50%, #6d28d9 100%);

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
