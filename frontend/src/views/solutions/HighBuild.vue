<template>
  <div class="solution-detail-page">
    <SiteHeader />

    <!-- Hero Section -->
    <section class="hero-section" style="background: linear-gradient(135deg, #e74c3c, #e67e22);">
      <div class="container">
        <div class="hero-content">
          <div class="hero-icon">
            <el-icon :size="48"><Warning /></el-icon>
          </div>
          <h1 class="hero-title">高防解决方案</h1>
          <p class="hero-desc">为企业提供全方位的网络安全防护，专业级DDoS防护、Web应用防火墙、入侵检测与防御，保障业务稳定运行，抵御各类网络攻击。</p>
          <div class="hero-actions">
            <el-button type="primary" size="large" round @click="$router.push('/products')">立即选购</el-button>
            <el-button size="large" round class="ghost-btn" @click="$router.push('/contact')">联系我们</el-button>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="section features-section">
      <div class="container">
        <div class="section-header">
          <h2>方案优势</h2>
          <p>全方位多层次的安全防护体系</p>
        </div>
        <div class="features-grid">
          <div class="feature-card" v-for="(feature, index) in features" :key="index">
            <div class="feature-icon">
              <el-icon :size="24"><component :is="feature.icon" /></el-icon>
            </div>
            <h3>{{ feature.title }}</h3>
            <p>{{ feature.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Products Section -->
    <section class="section products-section">
      <div class="container">
        <div class="section-header">
          <h2>推荐产品</h2>
          <p>为您推荐适合的高防产品配置</p>
        </div>
        <div class="products-grid">
          <div class="product-card" v-for="(product, index) in products" :key="index" :class="{ 'product-featured': product.featured }">
            <div class="product-badge" v-if="product.featured">推荐</div>
            <div class="product-header">
              <h3>{{ product.name }}</h3>
              <div class="product-price">
                <span class="price">¥{{ product.price }}</span>
                <span class="period">/月起</span>
              </div>
            </div>
            <ul class="product-specs">
              <li v-for="(spec, i) in product.specs" :key="i">
                <el-icon><Check /></el-icon>
                {{ spec }}
              </li>
            </ul>
            <el-button type="primary" @click="$router.push('/products')">立即购买</el-button>
          </div>
        </div>
      </div>
    </section>

    <!-- Use Cases Section -->
    <section class="section usecases-section">
      <div class="container">
        <div class="section-header">
          <h2>应用场景</h2>
          <p>覆盖企业全场景安全防护需求</p>
        </div>
        <div class="usecases-grid">
          <div class="usecase-card" v-for="(usecase, index) in usecases" :key="index">
            <h3>{{ usecase.title }}</h3>
            <p>{{ usecase.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="section cta-section">
      <div class="container">
        <div class="cta-content">
          <h2>守护业务安全</h2>
          <p>立即选购适合您的高防产品，或联系我们的安全顾问获取定制方案</p>
          <div class="cta-actions">
            <el-button type="primary" size="large" round @click="$router.push('/products')">立即选购</el-button>
            <el-button size="large" round @click="$router.push('/contact')">联系我们</el-button>
          </div>
        </div>
      </div>
    </section>

    <SiteFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Warning, Lock, Monitor, CircleCheck, Check } from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const route = useRoute()
const loading = ref(true)

const features = ref([
  { icon: 'Shield', title: 'DDoS防护', desc: '专业级DDoS防护能力，支持SYN Flood、UDP Flood、ICMP Flood等各类攻击防护，清洗能力高达T级' },
  { icon: 'Lock', title: 'WAF防火墙', desc: 'Web应用防火墙，有效防御SQL注入、XSS跨站、CC攻击等Web层面的安全威胁' },
  { icon: 'Monitor', title: '实时监控', desc: '7x24小时实时监控与告警，异常流量秒级发现，快速响应各类安全事件' },
  { icon: 'Warning', title: '入侵检测', desc: '智能入侵检测与防御系统，基于AI算法识别未知威胁，主动防御高级持续性威胁' }
])

const products = ref([
  { name: '高防基础版', price: 599, featured: false, specs: ['10Gbps防护', 'Web应用防火墙', '基础CC防护', '实时监控告警', '7x24技术支持'] },
  { name: '高防专业版', price: 1999, featured: true, specs: ['50Gbps防护', '高级WAF防护', 'CC智能防御', '入侵检测防御', '专属安全顾问'] },
  { name: '高防企业版', price: 4999, featured: false, specs: ['100Gbps+防护', '全栈安全防护', 'AI智能防御', '威胁情报分析', '定制安全方案'] }
])

const usecases = ref([
  { title: '游戏行业', desc: '防护游戏服务器DDoS攻击，保障玩家体验，防止游戏资产损失和用户流失' },
  { title: '金融行业', desc: '满足金融合规要求，保护交易系统安全，防止数据泄露和业务中断' },
  { title: '电商平台', desc: '保障大促期间业务稳定，防护CC攻击和恶意刷单，确保交易顺畅' },
  { title: '政务系统', desc: '满足等保合规要求，保护政务数据安全，防止黑客攻击和信息泄露' }
])

onMounted(async () => {
  try {
    const slug = route.path.split('/').pop()
    const res = await request.get(`/api/v1/solutions/${slug}`)
    if (res.data?.data) {
      const data = res.data.data
      if (data.features) features.value = data.features
      if (data.products) products.value = data.products
      if (data.usecases) usecases.value = data.usecases
    }
  } catch (error) {
    console.error('获取方案数据失败:', error)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped lang="scss">
.solution-detail-page {
  padding-top: 64px;

  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
  }

  .hero-section {
    color: #fff;
    padding: 100px 0;

    .hero-content {
      text-align: center;

      .hero-icon {
        width: 80px;
        height: 80px;
        margin: 0 auto 24px;
        background: rgba(255, 255, 255, 0.2);
        border-radius: 20px;
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .hero-title {
        font-size: 42px;
        font-weight: 700;
        margin: 0 0 16px;
      }

      .hero-desc {
        font-size: 18px;
        opacity: 0.9;
        margin: 0 0 30px;
        max-width: 700px;
        margin-left: auto;
        margin-right: auto;
        line-height: 1.6;
      }

      .hero-actions {
        display: flex;
        gap: 16px;
        justify-content: center;

        .ghost-btn {
          background: transparent;
          border: 2px solid #fff;
          color: #fff;

          &:hover {
            background: rgba(255, 255, 255, 0.1);
          }
        }
      }
    }
  }

  .section {
    padding: 80px 0;

    .section-header {
      text-align: center;
      margin-bottom: 50px;

      h2 {
        font-size: 32px;
        font-weight: 700;
        color: #1a2332;
        margin: 0 0 12px;
      }

      p {
        font-size: 16px;
        color: #666;
        margin: 0;
      }
    }
  }

  .features-section {
    background: #fff;

    .features-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 30px;
    }

    .feature-card {
      text-align: center;
      padding: 32px 24px;
      background: #f5f7fa;
      border-radius: 12px;
      transition: all 0.3s;

      &:hover {
        transform: translateY(-4px);
        box-shadow: 0 12px 24px rgba(0, 0, 0, 0.08);
      }

      .feature-icon {
        width: 56px;
        height: 56px;
        margin: 0 auto 16px;
        background: linear-gradient(135deg, #e74c3c, #e67e22);
        border-radius: 14px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;
      }

      h3 {
        font-size: 18px;
        font-weight: 600;
        color: #1a2332;
        margin: 0 0 8px;
      }

      p {
        font-size: 14px;
        color: #666;
        margin: 0;
        line-height: 1.6;
      }
    }
  }

  .products-section {
    background: #f5f7fa;

    .products-grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 30px;
    }

    .product-card {
      background: #fff;
      border-radius: 12px;
      padding: 32px;
      text-align: center;
      position: relative;
      transition: all 0.3s;

      &:hover {
        transform: translateY(-4px);
        box-shadow: 0 12px 24px rgba(0, 0, 0, 0.08);
      }

      &.product-featured {
        border: 2px solid #e74c3c;

        .product-badge {
          position: absolute;
          top: -12px;
          right: 20px;
          background: linear-gradient(135deg, #e74c3c, #e67e22);
          color: #fff;
          padding: 4px 16px;
          border-radius: 12px;
          font-size: 12px;
          font-weight: 600;
        }
      }

      .product-header {
        margin-bottom: 24px;

        h3 {
          font-size: 20px;
          font-weight: 600;
          color: #1a2332;
          margin: 0 0 16px;
        }

        .product-price {
          .price {
            font-size: 36px;
            font-weight: 700;
            color: #e74c3c;
          }

          .period {
            font-size: 14px;
            color: #909399;
          }
        }
      }

      .product-specs {
        list-style: none;
        padding: 0;
        margin: 0 0 24px;

        li {
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 8px;
          font-size: 14px;
          color: #555;
          padding: 8px 0;
          border-bottom: 1px solid #f0f0f0;

          &:last-child {
            border-bottom: none;
          }

          .el-icon {
            color: #67c23a;
          }
        }
      }
    }
  }

  .usecases-section {
    background: #fff;

    .usecases-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 24px;
    }

    .usecase-card {
      padding: 24px;
      background: #f5f7fa;
      border-radius: 12px;
      transition: all 0.3s;

      &:hover {
        background: linear-gradient(135deg, rgba(231, 76, 60, 0.05), rgba(230, 126, 34, 0.05));
      }

      h3 {
        font-size: 16px;
        font-weight: 600;
        color: #1a2332;
        margin: 0 0 8px;
      }

      p {
        font-size: 14px;
        color: #666;
        margin: 0;
        line-height: 1.6;
      }
    }
  }

  .cta-section {
    background: #f5f7fa;

    .cta-content {
      text-align: center;
      padding: 60px;
      background: linear-gradient(135deg, #e74c3c, #e67e22);
      border-radius: 20px;
      color: #fff;

      h2 {
        font-size: 32px;
        font-weight: 700;
        margin: 0 0 16px;
      }

      p {
        font-size: 16px;
        opacity: 0.8;
        margin: 0 0 30px;
      }

      .cta-actions {
        display: flex;
        gap: 16px;
        justify-content: center;
      }
    }
  }
}

@media (max-width: 768px) {
  .solution-detail-page {
    .features-grid,
    .products-grid,
    .usecases-grid {
      grid-template-columns: 1fr;
    }
  }
}
</style>
