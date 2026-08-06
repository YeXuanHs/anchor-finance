<template>
  <div class="solution-detail-page">
    <SiteHeader />

    <!-- Hero Section -->
    <section class="hero-section" style="background: linear-gradient(135deg, #45aaf2, #2d98da);">
      <div class="container">
        <div class="hero-content">
          <div class="hero-icon">
            <el-icon :size="48"><Monitor /></el-icon>
          </div>
          <h1 class="hero-title">网站建设解决方案</h1>
          <p class="hero-desc">为企业以及开发者用户实现灵活弹性自动化的基础IT设施建设、按需付费的服务模式以及0成本的运维IT服务体系，帮助企业客户快速建站，轻松迈入互联网+时代。</p>
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
          <p>为企业提供高效便捷的网站建设解决方案</p>
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
          <p>为您推荐适合的网站建设产品配置</p>
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
          <p>覆盖各类网站建设场景需求</p>
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
          <h2>快速开启互联网之旅</h2>
          <p>立即选购适合您的建站产品，或联系我们的技术顾问获取定制方案</p>
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
import { Monitor, Check } from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const route = useRoute()
const loading = ref(true)

const features = ref([
  { icon: 'Position', title: '弹性扩展', desc: '高效便捷、节约成本，根据网站访问量自动弹性扩展资源' },
  { icon: 'Shield', title: '全景备份', desc: '保证数据安全、降低备份成本，支持自动定时备份和快速恢复' },
  { icon: 'Lock', title: '安全防御', desc: '解决网络安全问题，网站免受攻击威胁，提供全方位安全防护' },
  { icon: 'Service', title: '贴心服务', desc: '7×24小时售后服务，专业技术团队随时解决用户疑难问题' }
])

const products = ref([
  { name: '建站入门版', price: 99, featured: false, specs: ['2核CPU', '4GB内存', '100GB SSD', '5Mbps带宽', '免费SSL证书'] },
  { name: '建站专业版', price: 299, featured: true, specs: ['4核CPU', '8GB内存', '500GB SSD', '20Mbps带宽', 'CDN加速'] },
  { name: '建站企业版', price: 699, featured: false, specs: ['8核CPU', '16GB内存', '1TB SSD', '50Mbps带宽', '全方位安全防护'] }
])

const usecases = ref([
  { title: '小型网站', desc: '适用于起步阶段的小型网站，日均PV低于5万，提供稳定可靠的托管服务' },
  { title: '中型网站', desc: '适用于成长阶段的中小型网站，日均PV低于30万，支持弹性扩容' },
  { title: '大型企业网站', desc: '适用于稳定运营的大型网站，日均PV百万以上，提供高性能保障' },
  { title: '电商平台', desc: '支持电商网站高并发访问，保障促销活动期间网站稳定运行' }
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
        background: linear-gradient(135deg, #45aaf2, #2d98da);
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
        border: 2px solid #45aaf2;

        .product-badge {
          position: absolute;
          top: -12px;
          right: 20px;
          background: linear-gradient(135deg, #45aaf2, #2d98da);
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
            color: #45aaf2;
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
        background: linear-gradient(135deg, rgba(69, 170, 242, 0.05), rgba(45, 152, 218, 0.05));
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
      background: linear-gradient(135deg, #45aaf2, #2d98da);
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
