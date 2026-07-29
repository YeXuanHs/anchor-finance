<template>
  <div class="ha-page">
    <SiteHeader />
    <section class="hero-section">
      <div class="container">
        <h1>高可用架构解决方案</h1>
        <p class="hero-desc">构建永不中断的业务系统</p>
      </div>
    </section>
    <section class="content-section">
      <div class="container">
        <div class="section-title"><h2>方案特点</h2></div>
        <div class="features-grid">
          <div class="feature-card" v-for="f in features" :key="f.title">
            <h3>{{ f.title }}</h3>
            <p>{{ f.desc }}</p>
          </div>
        </div>
        <div class="section-title"><h2>适用场景</h2></div>
        <div class="scenarios">
          <div class="scenario" v-for="s in scenarios" :key="s">{{ s }}</div>
        </div>
        <div class="cta-box">
          <h2>定制高可用方案</h2>
          <p>联系技术顾问，获取专属架构设计</p>
          <router-link to="/contact" class="btn-primary">咨询方案</router-link>
        </div>
      </div>
    </section>
    <SiteFooter />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const route = useRoute()
const loading = ref(true)

const features = ref([
  { title: '多活部署', desc: '业务部署在多个可用区，单点故障不影响整体服务' },
  { title: '自动故障转移', desc: '秒级检测故障，自动切换到健康节点' },
  { title: '负载均衡', desc: '智能分发流量，充分利用每一台服务器' },
  { title: '数据同步', desc: '主从实时同步，保证数据一致性' },
  { title: '弹性伸缩', desc: '根据负载自动扩缩容，应对流量高峰' },
  { title: '全链路监控', desc: '从网络到应用的全链路监控和告警' }
])
const scenarios = ref(['金融交易系统', '电商大促', '在线游戏', 'SaaS平台', '政务系统', '医疗信息系统'])

onMounted(async () => {
  try {
    const slug = route.path.split('/').pop()
    const res = await request.get(`/api/v1/solutions/${slug}`)
    if (res.data?.data) {
      const data = res.data.data
      if (data.features) features.value = data.features
      if (data.scenarios) scenarios.value = data.scenarios
    }
  } catch (error) {
    console.error('获取方案数据失败:', error)
  } finally {
    loading.value = false
  }
})
</script>
<style scoped lang="scss">
.ha-page { min-height: 100vh; }
.hero-section { background: linear-gradient(135deg, #1a365d 0%, #2563eb 100%); padding: 80px 0 60px; text-align: center; color: #fff;
  h1 { font-size: 36px; font-weight: 700; margin-bottom: 12px; }
  .hero-desc { font-size: 16px; opacity: 0.8; }
}
.container { max-width: 1200px; margin: 0 auto; padding: 0 20px; }
.section-title { text-align: center; margin-bottom: 40px; h2 { font-size: 28px; color: #1a365d; } }
.content-section { padding: 80px 0; }
.features-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 24px; margin-bottom: 60px; }
.feature-card { background: #fff; border-radius: 12px; padding: 28px; box-shadow: 0 2px 12px rgba(0,0,0,0.06);
  h3 { font-size: 18px; color: #1a365d; margin-bottom: 8px; }
  p { font-size: 14px; color: #6b7280; line-height: 1.6; }
}
.scenarios { display: flex; flex-wrap: wrap; gap: 12px; justify-content: center; margin-bottom: 60px; }
.scenario { padding: 10px 20px; background: #eff6ff; color: #2563eb; border-radius: 20px; font-size: 14px; }
.cta-box { background: linear-gradient(135deg, #1a365d 0%, #2563eb 100%); border-radius: 12px; padding: 48px; text-align: center; color: #fff;
  h2 { font-size: 24px; margin-bottom: 12px; }
  p { opacity: 0.8; margin-bottom: 24px; }
  .btn-primary { display: inline-block; padding: 12px 32px; background: #fff; color: #2563eb; border-radius: 8px; text-decoration: none; font-weight: 600; }
}
</style>
