<template>
  <div class="solution-detail-page">
    <SiteHeader />

    <!-- Hero Section -->
    <section class="hero-section" style="background: linear-gradient(135deg, #4b7bec, #3867d6);">
      <div class="container">
        <div class="hero-content">
          <div class="hero-icon">
            <el-icon :size="48"><Connection /></el-icon>
          </div>
          <h1 class="hero-title">混合云解决方案</h1>
          <p class="hero-desc">提供客户在本地数据中心使用云服务的能力，满足客户特定的安全和合规要求，通过持续迭代演进提供满足业务要求的云服务，同时解决部分业务场景低时延的限定要求。</p>
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
          <p>为企业提供全方位的混合云解决方案</p>
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
          <p>为您推荐适合的混合云产品配置</p>
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
          <p>覆盖企业常见混合云场景需求</p>
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
          <h2>构建灵活混合云架构</h2>
          <p>立即选购适合您的混合云产品，或联系我们的技术顾问获取定制方案</p>
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
import { Connection, Check } from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const route = useRoute()
const loading = ref(true)

const features = ref([
  { icon: 'Position', title: '全方位混合云架构', desc: '裸机、虚拟机全方位计算服务，公有云、私有云、托管云全方位混合云能力' },
  { icon: 'Cpu', title: '动态扩展', desc: '在本地数据中心及公有云实现资源混合部署，充分利用公有云资源优势实现动态扩展' },
  { icon: 'DataLine', title: '云网一体', desc: '通过高速专线服务实现云平台间VPC级别的高质量、安全和高可靠网络互通' },
  { icon: 'Lock', title: '多云统一管理', desc: '对接纳管本地私有云和VMware虚拟化环境，实现对企业IaaS云计算环境资源的统一管理' }
])

const products = ref([
  { name: '混合云基础版', price: 599, featured: false, specs: ['基础云主机', '专线接入', '统一管理面板', '基础安全防护'] },
  { name: '混合云专业版', price: 1299, featured: true, specs: ['高性能云主机', '高速专线', '多云管理', '高级安全防护'] },
  { name: '混合云企业版', price: 2999, featured: false, specs: ['专属集群', '多中心架构', '智能调度', '全方位安全防护'] }
])

const usecases = ref([
  { title: '业务波峰波谷', desc: '将平稳业务部署于私有云，使用公有云作为灵活扩展储备，访问量升高时分担业务压力' },
  { title: '多中心灾备', desc: '设计可动态扩容的基础架构，主系统故障时快速扩容灾备节点接替主节点工作' },
  { title: '合规可控', desc: '将敏感业务和数据部署在私有云满足监管要求，将用户访问业务部署在公有云' },
  { title: '数据迁移', desc: '帮助企业平滑迁移上云，确保业务连续性，降低迁移风险和成本' }
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
        background: linear-gradient(135deg, #4b7bec, #3867d6);
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
        border: 2px solid #4b7bec;

        .product-badge {
          position: absolute;
          top: -12px;
          right: 20px;
          background: linear-gradient(135deg, #4b7bec, #3867d6);
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
            color: #4b7bec;
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
        background: linear-gradient(135deg, rgba(75, 123, 236, 0.05), rgba(56, 103, 214, 0.05));
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
      background: linear-gradient(135deg, #4b7bec, #3867d6);
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
