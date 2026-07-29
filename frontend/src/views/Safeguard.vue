<template>
  <div class="safeguard-page">
    <SiteHeader />
    <section class="hero-section">
      <div class="container">
        <h1>服务保障</h1>
        <p class="hero-desc">多重保障，让您的业务无忧运行</p>
      </div>
    </section>
    <section class="content-section">
      <div class="container">
        <div class="section-title"><h2>我们的承诺</h2></div>
        <div class="guarantees-grid">
          <div class="guarantee-card" v-for="g in guarantees" :key="g.title">
            <div class="icon">{{ g.icon }}</div>
            <h3>{{ g.title }}</h3>
            <p>{{ g.desc }}</p>
          </div>
        </div>
        <div class="sla-section">
          <h2>SLA 服务等级协议</h2>
          <table class="sla-table">
            <thead><tr><th>服务等级</th><th>可用性</th><th>故障响应</th><th>赔偿标准</th></tr></thead>
            <tbody>
              <tr><td>基础版</td><td>99.9%</td><td>30分钟</td><td>月费10%赔偿</td></tr>
              <tr><td>专业版</td><td>99.95%</td><td>15分钟</td><td>月费25%赔偿</td></tr>
              <tr><td>企业版</td><td>99.99%</td><td>5分钟</td><td>月费100%赔偿</td></tr>
            </tbody>
          </table>
        </div>
      </div>
    </section>
    <SiteFooter />
  </div>
</template>
<script setup lang="ts">
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
const guarantees = [
  { icon: '⚡', title: '99.99% 可用性', desc: '基于多活架构和冗余设计，保障业务连续性' },
  { icon: '🛡️', title: 'DDoS 防护', desc: '免费基础DDoS防护，可选高防方案' },
  { icon: '💾', title: '数据备份', desc: '每日自动备份，支持7天内任意时间点恢复' },
  { icon: '🔒', title: '数据安全', desc: '传输加密、存储加密，符合等保要求' },
  { icon: '📞', title: '7x24 支持', desc: '全天候技术支持，电话/工单/在线客服' },
  { icon: '💰', title: '赔偿保障', desc: '未达到SLA承诺，按标准进行赔偿' }
]
</script>
<style scoped lang="scss">
.safeguard-page { min-height: 100vh; }
.hero-section { background: linear-gradient(135deg, #1a365d 0%, #2563eb 100%); padding: 80px 0 60px; text-align: center; color: #fff;
  h1 { font-size: 36px; font-weight: 700; margin-bottom: 12px; }
  .hero-desc { font-size: 16px; opacity: 0.8; }
}
.container { max-width: 1200px; margin: 0 auto; padding: 0 20px; }
.section-title { text-align: center; margin-bottom: 40px; h2 { font-size: 28px; color: #1a365d; } }
.content-section { padding: 80px 0; }
.guarantees-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 24px; margin-bottom: 60px; }
.guarantee-card { background: #fff; border-radius: 12px; padding: 32px; text-align: center; box-shadow: 0 2px 12px rgba(0,0,0,0.06);
  .icon { font-size: 40px; margin-bottom: 16px; }
  h3 { font-size: 18px; color: #1a365d; margin-bottom: 8px; }
  p { color: #6b7280; line-height: 1.6; }
}
.sla-section { background: #fff; border-radius: 12px; padding: 40px; box-shadow: 0 2px 12px rgba(0,0,0,0.06);
  h2 { font-size: 24px; color: #1a365d; margin-bottom: 24px; text-align: center; }
}
.sla-table { width: 100%; border-collapse: collapse;
  th, td { padding: 14px 20px; text-align: center; border-bottom: 1px solid #e5e7eb; }
  th { background: #f8fafc; color: #1a365d; font-weight: 600; }
  td { color: #4b5563; }
  tr:hover td { background: #f0f9ff; }
}
</style>
