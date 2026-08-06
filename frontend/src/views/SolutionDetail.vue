<template>
  <div class="solution-detail-page">
    <SiteHeader />
    
    <!-- Hero Section -->
    <section class="hero-section" :style="{ background: solution.gradient }">
      <div class="container">
        <div class="hero-content">
          <div class="hero-icon">
            <el-icon :size="48"><component :is="solution.icon" /></el-icon>
          </div>
          <h1 class="hero-title">{{ solution.title }}</h1>
          <p class="hero-desc">{{ solution.desc }}</p>
          <div class="hero-actions">
            <el-button type="primary" size="large" round @click="$router.push('/products')">立即选购</el-button>
            <el-button size="large" round class="ghost-btn">联系我们</el-button>
          </div>
        </div>
      </div>
    </section>
    
    <!-- Features Section -->
    <section class="section features-section">
      <div class="container">
        <div class="section-header">
          <h2>方案优势</h2>
          <p>为什么选择我们的解决方案</p>
        </div>
        
        <div class="features-grid">
          <div class="feature-card" v-for="(feature, index) in solution.features" :key="index">
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
          <p>为您推荐适合的产品配置</p>
        </div>
        
        <div class="products-grid">
          <div class="product-card" v-for="(product, index) in solution.products" :key="index">
            <div class="product-header">
              <h3>{{ product.name }}</h3>
              <div class="product-price">
                <span class="price">¥{{ product.price }}</span>
                <span class="period">/月</span>
              </div>
            </div>
            <ul class="product-specs">
              <li v-for="(spec, i) in product.specs" :key="i">
                <el-icon><Check /></el-icon>
                {{ spec }}
              </li>
            </ul>
            <el-button type="primary" @click="$router.push(`/products/${product.id}`)">立即购买</el-button>
          </div>
        </div>
      </div>
    </section>
    
    <!-- Use Cases Section -->
    <section class="section usecases-section">
      <div class="container">
        <div class="section-header">
          <h2>应用场景</h2>
          <p>了解解决方案的典型应用</p>
        </div>
        
        <div class="usecases-grid">
          <div class="usecase-card" v-for="(usecase, index) in solution.usecases" :key="index">
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
          <h2>准备好开始了吗？</h2>
          <p>立即选购适合您的产品，或联系我们的技术顾问获取定制方案</p>
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
import { Check } from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const route = useRoute()
const loading = ref(true)

const solution = ref({
  title: '解决方案',
  desc: '专业的云计算解决方案',
  icon: 'OfficeBuilding',
  gradient: 'linear-gradient(135deg, #1a237e, #0d47a1)',
  features: [] as any[],
  products: [] as any[],
  usecases: [] as any[]
})

onMounted(async () => {
  try {
    const id = route.params.id
    const res = await request.get(`/api/v1/solutions/${id}`)
    if (res.data?.data) {
      solution.value = { ...solution.value, ...res.data.data }
    }
  } catch (error) {
    console.error('获取方案详情失败:', error)
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
        max-width: 600px;
        margin-left: auto;
        margin-right: auto;
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
      
      .feature-icon {
        width: 56px;
        height: 56px;
        margin: 0 auto 16px;
        background: linear-gradient(135deg, #409eff, #66b1ff);
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
            color: #409eff;
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
      background: linear-gradient(135deg, #1a237e, #0d47a1);
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
