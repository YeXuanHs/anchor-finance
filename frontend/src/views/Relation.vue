<template>
  <div class="relation-page">
    <SiteHeader />
    <section class="hero-section">
      <div class="container">
        <h1>生态合作</h1>
        <p class="hero-desc">携手共建，互利共赢</p>
      </div>
    </section>
    <section class="content-section">
      <div class="container">
        <div class="section-title"><h2>合作模式</h2></div>
        <div class="models-grid">
          <div class="model-card" v-for="m in models" :key="m.title">
            <h3>{{ m.title }}</h3>
            <p>{{ m.desc }}</p>
            <ul><li v-for="b in m.benefits" :key="b">{{ b }}</li></ul>
          </div>
        </div>
        <div class="apply-section">
          <h2>申请合作</h2>
          <p>如果您对我们的合作模式感兴趣，请联系我们：</p>
          <router-link to="/contact" class="btn-primary">立即申请</router-link>
        </div>
      </div>
    </section>
    <SiteFooter />
  </div>
</template>
<script setup lang="ts">
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
const models = [
  { title: '代理合作', desc: '成为我们的代理商，享受优惠价格和专属支持', benefits: ['优惠的代理价格', '专属客户经理', '营销素材支持', '技术培训支持'] },
  { title: '技术合作', desc: '技术互补，共同为客户提供更完善的解决方案', benefits: ['API对接支持', '技术文档共享', '联合解决方案', '优先技术支持'] },
  { title: '渠道合作', desc: '资源共享，扩大市场覆盖范围', benefits: ['品牌联合推广', '客户资源共享', '佣金分成', '市场活动合作'] }
]
</script>
<style scoped lang="scss">
.relation-page { min-height: 100vh; }
.hero-section { background: linear-gradient(135deg, #1a365d 0%, #2563eb 100%); padding: 80px 0 60px; text-align: center; color: #fff;
  h1 { font-size: 36px; font-weight: 700; margin-bottom: 12px; }
  .hero-desc { font-size: 16px; opacity: 0.8; }
}
.container { max-width: 1200px; margin: 0 auto; padding: 0 20px; }
.section-title { text-align: center; margin-bottom: 40px; h2 { font-size: 28px; color: #1a365d; } }
.content-section { padding: 80px 0; }
.models-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 24px; margin-bottom: 60px; }
.model-card { background: #fff; border-radius: 12px; padding: 32px; box-shadow: 0 2px 12px rgba(0,0,0,0.06);
  h3 { font-size: 20px; color: #1a365d; margin-bottom: 12px; }
  p { color: #6b7280; margin-bottom: 16px; }
  ul { padding-left: 20px; li { color: #4b5563; line-height: 2; } }
}
.apply-section { text-align: center; background: #f0f9ff; border-radius: 12px; padding: 48px;
  h2 { font-size: 24px; color: #1a365d; margin-bottom: 12px; }
  p { color: #6b7280; margin-bottom: 24px; }
  .btn-primary { display: inline-block; padding: 12px 32px; background: #2563eb; color: #fff; border-radius: 8px; text-decoration: none; font-weight: 600; }
}
</style>
