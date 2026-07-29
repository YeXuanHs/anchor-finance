<template>
  <div class="transfer-page">
    <SiteHeader />
    <section class="hero-section">
      <div class="container">
        <h1>上云迁移</h1>
        <p class="hero-desc">专业迁移服务，零停机平滑上云</p>
      </div>
    </section>
    <section class="content-section">
      <div class="container">
        <div class="section-title"><h2>迁移流程</h2></div>
        <div class="steps">
          <div class="step" v-for="(s, i) in steps" :key="i">
            <div class="step-num">{{ i + 1 }}</div>
            <h3>{{ s.title }}</h3>
            <p>{{ s.desc }}</p>
          </div>
        </div>
        <div class="section-title" style="margin-top:60px"><h2>支持的迁移场景</h2></div>
        <div class="scenarios-grid">
          <div class="scenario-card" v-for="s in scenarios" :key="s.title">
            <h3>{{ s.title }}</h3>
            <p>{{ s.desc }}</p>
          </div>
        </div>
        <div class="cta-box">
          <h2>免费迁移评估</h2>
          <p>提交您的需求，我们的工程师将在24小时内为您提供迁移方案和报价</p>
          <router-link to="/contact" class="btn-primary">申请迁移</router-link>
        </div>
      </div>
    </section>
    <SiteFooter />
  </div>
</template>
<script setup lang="ts">
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
const steps = [
  { title: '需求评估', desc: '分析现有环境，评估迁移复杂度和工作量' },
  { title: '方案设计', desc: '制定详细迁移方案，包括时间计划和回退方案' },
  { title: '环境准备', desc: '搭建目标环境，配置网络和安全策略' },
  { title: '数据迁移', desc: '分批迁移数据，确保数据完整性和一致性' },
  { title: '测试验证', desc: '功能测试、性能测试、安全测试' },
  { title: '切换上线', desc: 'DNS切换、流量导入、监控观察' }
]
const scenarios = [
  { title: '传统IDC上云', desc: '从传统物理服务器迁移到云服务器' },
  { title: '跨云迁移', desc: '从其他云平台迁移到我们的云平台' },
  { title: '混合云部署', desc: '实现私有云和公有云的混合架构' },
  { title: '数据库迁移', desc: 'MySQL、PostgreSQL、MongoDB等数据库迁移' }
]
</script>
<style scoped lang="scss">
.transfer-page { min-height: 100vh; }
.hero-section { background: linear-gradient(135deg, #1a365d 0%, #2563eb 100%); padding: 80px 0 60px; text-align: center; color: #fff;
  h1 { font-size: 36px; font-weight: 700; margin-bottom: 12px; }
  .hero-desc { font-size: 16px; opacity: 0.8; }
}
.container { max-width: 1200px; margin: 0 auto; padding: 0 20px; }
.section-title { text-align: center; margin-bottom: 40px; h2 { font-size: 28px; color: #1a365d; } }
.content-section { padding: 80px 0; }
.steps { display: flex; flex-wrap: wrap; gap: 24px; justify-content: center; }
.step { flex: 1; min-width: 180px; max-width: 200px; text-align: center;
  .step-num { width: 48px; height: 48px; border-radius: 50%; background: #2563eb; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 20px; font-weight: 700; margin: 0 auto 12px; }
  h3 { font-size: 16px; color: #1a365d; margin-bottom: 8px; }
  p { font-size: 14px; color: #6b7280; }
}
.scenarios-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin-bottom: 60px; }
.scenario-card { background: #fff; border-radius: 12px; padding: 24px; box-shadow: 0 2px 12px rgba(0,0,0,0.06); text-align: center;
  h3 { font-size: 16px; color: #1a365d; margin-bottom: 8px; }
  p { font-size: 14px; color: #6b7280; }
}
.cta-box { background: linear-gradient(135deg, #1a365d 0%, #2563eb 100%); border-radius: 12px; padding: 48px; text-align: center; color: #fff;
  h2 { font-size: 24px; margin-bottom: 12px; }
  p { opacity: 0.8; margin-bottom: 24px; }
  .btn-primary { display: inline-block; padding: 12px 32px; background: #fff; color: #2563eb; border-radius: 8px; text-decoration: none; font-weight: 600; }
}
</style>
