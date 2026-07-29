<template>
  <div class="dediserver-page">
    <SiteHeader />
    <section class="hero-section">
      <div class="container">
        <h1>站群服务器</h1>
        <p class="hero-desc">多IP段独立服务器，SEO站群首选</p>
      </div>
    </section>
    <section class="content-section">
      <div class="container">
        <div class="features-grid">
          <div class="feature-card" v-for="f in features" :key="f.title">
            <div class="icon">{{ f.icon }}</div>
            <h3>{{ f.title }}</h3>
            <p>{{ f.desc }}</p>
          </div>
        </div>
        <div class="section-title"><h2>推荐配置</h2></div>
        <div class="plans-grid" v-loading="loading">
          <div class="plan-card" v-for="p in plans" :key="p.name" :class="{ featured: p.featured }">
            <h3>{{ p.name }}</h3>
            <div class="price">¥{{ p.price }}<span>/月</span></div>
            <ul><li v-for="spec in p.specs" :key="spec">{{ spec }}</li></ul>
            <router-link to="/products" class="btn-primary">立即选购</router-link>
          </div>
        </div>
      </div>
    </section>
    <SiteFooter />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const loading = ref(false)

const features = [
  { icon: '🌐', title: '多IP段', desc: '支持多达253个独立IP，不同C段可选' },
  { icon: '⚡', title: '高性能', desc: '至强处理器，SSD存储，万兆网络' },
  { icon: '🔒', title: '独立资源', desc: '独享CPU、内存、带宽，不受邻居影响' },
  { icon: '🛡️', title: 'DDoS防护', desc: '免费基础防护，可选高防方案' }
]

const plans = ref([])

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/products', { params: { group: 'dedicated' } })
    if (data?.data) {
      plans.value = data.data.list || data.data || plans.value
    }
  } catch (e) {
    console.error('Failed to fetch dedicated server products:', e)
  } finally {
    loading.value = false
  }
})
</script>
<style scoped lang="scss">
.dediserver-page { min-height: 100vh; }
.hero-section { background: linear-gradient(135deg, #1a365d 0%, #2563eb 100%); padding: 80px 0 60px; text-align: center; color: #fff;
  h1 { font-size: 36px; font-weight: 700; margin-bottom: 12px; }
  .hero-desc { font-size: 16px; opacity: 0.8; }
}
.container { max-width: 1200px; margin: 0 auto; padding: 0 20px; }
.section-title { text-align: center; margin-bottom: 40px; h2 { font-size: 28px; color: #1a365d; } }
.content-section { padding: 80px 0; }
.features-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 24px; margin-bottom: 60px; }
.feature-card { text-align: center; padding: 24px;
  .icon { font-size: 40px; margin-bottom: 12px; }
  h3 { font-size: 18px; color: #1a365d; margin-bottom: 8px; }
  p { font-size: 14px; color: #6b7280; }
}
.plans-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 24px; }
.plan-card { background: #fff; border-radius: 12px; padding: 32px; text-align: center; box-shadow: 0 2px 12px rgba(0,0,0,0.06); border: 2px solid transparent;
  &.featured { border-color: #2563eb; transform: scale(1.05); }
  h3 { font-size: 20px; color: #1a365d; margin-bottom: 16px; }
  .price { font-size: 36px; font-weight: 700; color: #2563eb; margin-bottom: 20px; span { font-size: 14px; color: #6b7280; font-weight: 400; } }
  ul { list-style: none; padding: 0; margin-bottom: 24px; li { padding: 8px 0; color: #4b5563; border-bottom: 1px solid #f3f4f6; } }
  .btn-primary { display: inline-block; padding: 10px 28px; background: #2563eb; color: #fff; border-radius: 8px; text-decoration: none; }
}
</style>
