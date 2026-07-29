<template>
  <div class="management-page">
    <SiteHeader />
    <section class="hero-section">
      <div class="container">
        <h1>管家服务</h1>
        <p class="hero-desc">专业的服务器运维托管，让您专注业务发展</p>
      </div>
    </section>
    <section class="services-section">
      <div class="container">
        <div class="section-title"><h2>服务内容</h2></div>
        <div class="services-grid">
          <div class="service-card" v-for="s in services" :key="s.title">
            <div class="icon">{{ s.icon }}</div>
            <h3>{{ s.title }}</h3>
            <p>{{ s.desc }}</p>
          </div>
        </div>
      </div>
    </section>
    <section class="cta-section">
      <div class="container">
        <h2>需要管家服务？</h2>
        <p>联系我们的技术顾问，为您定制专属运维方案</p>
        <router-link to="/contact" class="btn-primary">联系我们</router-link>
      </div>
    </section>
    <SiteFooter />
  </div>
</template>
<script setup lang="ts">
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
const services = [
  { icon: '🔧', title: '系统运维', desc: '操作系统安装、配置、更新和安全加固' },
  { icon: '📊', title: '性能监控', desc: '7x24小时性能监控，及时发现并解决潜在问题' },
  { icon: '🔒', title: '安全防护', desc: '防火墙配置、入侵检测、安全审计和漏洞修复' },
  { icon: '💾', title: '数据备份', desc: '定期自动备份，支持异地容灾和快速恢复' },
  { icon: '🚀', title: '环境部署', desc: 'Web服务器、数据库、缓存等环境的部署和优化' },
  { icon: '📞', title: '技术支持', desc: '专属技术顾问，电话/工单/远程支持' }
]
</script>
<style scoped lang="scss">
.management-page { min-height: 100vh; }
.hero-section { background: linear-gradient(135deg, #1a365d 0%, #2563eb 100%); padding: 80px 0 60px; text-align: center; color: #fff;
  h1 { font-size: 36px; font-weight: 700; margin-bottom: 12px; }
  .hero-desc { font-size: 16px; opacity: 0.8; }
}
.container { max-width: 1200px; margin: 0 auto; padding: 0 20px; }
.section-title { text-align: center; margin-bottom: 40px;
  h2 { font-size: 28px; color: #1a365d; }
}
.services-section { padding: 80px 0; }
.services-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 24px; }
.service-card { background: #fff; border-radius: 12px; padding: 32px; text-align: center; box-shadow: 0 2px 12px rgba(0,0,0,0.06); transition: transform 0.3s;
  &:hover { transform: translateY(-4px); }
  .icon { font-size: 40px; margin-bottom: 16px; }
  h3 { font-size: 18px; color: #1a365d; margin-bottom: 8px; }
  p { color: #6b7280; line-height: 1.6; }
}
.cta-section { background: linear-gradient(135deg, #1a365d 0%, #2563eb 100%); padding: 60px 0; text-align: center; color: #fff;
  h2 { font-size: 28px; margin-bottom: 12px; }
  p { opacity: 0.8; margin-bottom: 24px; }
  .btn-primary { display: inline-block; padding: 12px 32px; background: #fff; color: #2563eb; border-radius: 8px; text-decoration: none; font-weight: 600; }
}
</style>
