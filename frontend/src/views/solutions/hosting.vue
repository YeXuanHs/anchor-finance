<template>
  <div class="solution-detail-page">
    <SiteHeader />

    <!-- Hero Section -->
    <section class="hero-section" style="background: linear-gradient(135deg, #fc5c65, #eb3b5a);">
      <div class="container">
        <div class="hero-content">
          <div class="hero-icon">
            <el-icon :size="48"><OfficeBuilding /></el-icon>
          </div>
          <h1 class="hero-title">虚拟主机解决方案</h1>
          <p class="hero-desc">为核心数据库、高性能计算业务提供云端专用的高性能、安全隔离的物理集群，实现配置随意扩展，灵活开展业务，即开即用，简单易用。</p>
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
          <p>为用户提供简单易用的虚拟主机服务</p>
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
          <p>为您推荐适合的虚拟主机套餐</p>
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
          <p>覆盖虚拟主机常见使用场景</p>
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
          <h2>轻松开启网站之旅</h2>
          <p>立即选购适合您的虚拟主机产品，或联系我们的技术顾问获取定制方案</p>
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
import { OfficeBuilding, Check } from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const route = useRoute()
const loading = ref(true)

const features = ref([
  { icon: 'Position', title: '即开即用', desc: '一键开通虚拟主机，无需复杂配置，快速部署您的网站应用' },
  { icon: 'Cpu', title: '简单易用', desc: '可视化控制面板，轻松管理网站文件、数据库、域名等' },
  { icon: 'Lock', title: '安全稳定', desc: '独立资源隔离，DDoS防护，数据自动备份，保障网站安全稳定运行' },
  { icon: 'DataLine', title: '高性价比', desc: '按需选择配置，灵活升级，低成本快速建站，适合个人和中小企业' }
])

const products = ref([
  { name: '基础型', price: 29, featured: false, specs: ['1GB空间', '10GB月流量', '1个网站', '1个数据库', '基础防护'] },
  { name: '普惠型', price: 59, featured: true, specs: ['5GB空间', '50GB月流量', '3个网站', '3个数据库', 'SSL证书'] },
  { name: '入门型', price: 99, featured: false, specs: ['10GB空间', '100GB月流量', '5个网站', '5个数据库', 'CDN加速'] },
  { name: '旗舰型', price: 199, featured: false, specs: ['50GB空间', '无限流量', '无限网站', '无限数据库', '全方位防护'] }
])

const usecases = ref([
  { title: '个人博客', desc: '为个人博客提供稳定可靠的托管服务，轻松搭建个人网站，分享生活点滴' },
  { title: '企业官网', desc: '为企业官网提供高性能托管，展示企业形象，提升品牌影响力' },
  { title: '论坛社区', desc: '支持论坛社区类网站部署，提供稳定的运行环境和数据存储' },
  { title: '小型电商', desc: '适合小型电商网站托管，支持在线商城搭建，开启电商之旅' }
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
        background: linear-gradient(135deg, #fc5c65, #eb3b5a);
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
      grid-template-columns: repeat(4, 1fr);
      gap: 24px;
    }

    .product-card {
      background: #fff;
      border-radius: 12px;
      padding: 28px 20px;
      text-align: center;
      position: relative;
      transition: all 0.3s;

      &:hover {
        transform: translateY(-4px);
        box-shadow: 0 12px 24px rgba(0, 0, 0, 0.08);
      }

      &.product-featured {
        border: 2px solid #fc5c65;

        .product-badge {
          position: absolute;
          top: -12px;
          right: 16px;
          background: linear-gradient(135deg, #fc5c65, #eb3b5a);
          color: #fff;
          padding: 4px 12px;
          border-radius: 12px;
          font-size: 12px;
          font-weight: 600;
        }
      }

      .product-header {
        margin-bottom: 20px;

        h3 {
          font-size: 18px;
          font-weight: 600;
          color: #1a2332;
          margin: 0 0 12px;
        }

        .product-price {
          .price {
            font-size: 32px;
            font-weight: 700;
            color: #fc5c65;
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
        margin: 0 0 20px;

        li {
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 6px;
          font-size: 13px;
          color: #555;
          padding: 6px 0;
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
        background: linear-gradient(135deg, rgba(252, 92, 101, 0.05), rgba(235, 59, 90, 0.05));
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
      background: linear-gradient(135deg, #fc5c65, #eb3b5a);
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
    .usecases-grid {
      grid-template-columns: 1fr;
    }

    .products-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }
}
</style>
