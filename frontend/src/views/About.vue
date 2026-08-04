<template>
  <div class="about-page">
    <SiteHeader />
    
    <!-- Hero Section -->
    <section class="hero-section">
      <div class="container">
        <h1 class="hero-title">关于{{ siteSettings.site_name || '我们' }}</h1>
        <p class="hero-desc">专业、可靠、创新的财务管理系统解决方案</p>
      </div>
    </section>
    
    <!-- Company Introduction -->
    <section class="section company-intro">
      <div class="container">
        <div class="section-header">
          <h2>公司介绍</h2>
          <p>了解我们的故事和使命</p>
        </div>
        <div class="intro-content">
          <div class="intro-text">
            <p>{{ siteSettings.site_name || '我们' }}是一家专注于财务管理系统研发的科技公司，致力于为企业提供专业、高效、安全的财务管理解决方案。</p>
            <p>我们的团队由资深的财务专家和技术工程师组成，拥有丰富的行业经验和深厚的技术积累。通过不断的技术创新和产品优化，我们已经为数千家企业提供了优质的财务管理系统服务。</p>
            <p>我们的使命是通过技术创新，帮助企业实现财务管理的数字化转型，提升管理效率，降低运营成本，助力企业快速发展。</p>
          </div>
          <div class="intro-stats">
            <div class="stat-item">
              <div class="stat-number">5000+</div>
              <div class="stat-label">服务企业</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">99.9%</div>
              <div class="stat-label">系统可用性</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">50+</div>
              <div class="stat-label">技术专利</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">24/7</div>
              <div class="stat-label">技术支持</div>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- Development History -->
    <section class="section history-section">
      <div class="container">
        <div class="section-header">
          <h2>发展历程</h2>
          <p>我们的成长足迹</p>
        </div>
        <div class="timeline">
          <div class="timeline-item" v-for="(item, index) in history" :key="index">
            <div class="timeline-year">{{ item.year }}</div>
            <div class="timeline-content">
              <h3>{{ item.title }}</h3>
              <p>{{ item.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- Core Values -->
    <section class="section values-section">
      <div class="container">
        <div class="section-header">
          <h2>核心价值观</h2>
          <p>指导我们前行的理念</p>
        </div>
        <div class="values-grid">
          <div class="value-card" v-for="(item, index) in values" :key="index">
            <div class="value-icon">
              <el-icon :size="32"><component :is="item.icon" /></el-icon>
            </div>
            <h3>{{ item.title }}</h3>
            <p>{{ item.desc }}</p>
          </div>
        </div>
      </div>
    </section>
    
    <!-- Team Section -->
    <section class="section team-section">
      <div class="container">
        <div class="section-header">
          <h2>我们的团队</h2>
          <p>专业、敬业、创新的团队</p>
        </div>
        <div class="team-grid">
          <div class="team-card" v-for="(member, index) in team" :key="index">
            <div class="team-avatar">
              <el-avatar :size="80">{{ member.name.charAt(0) }}</el-avatar>
            </div>
            <h3>{{ member.name }}</h3>
            <p class="team-role">{{ member.role }}</p>
            <p class="team-desc">{{ member.desc }}</p>
          </div>
        </div>
      </div>
    </section>
    
    <!-- Contact Section -->
    <section class="section contact-section">
      <div class="container">
        <div class="section-header">
          <h2>联系我们</h2>
          <p>期待与您的合作</p>
        </div>
        <div class="contact-grid">
          <div class="contact-card">
            <el-icon :size="32"><Location /></el-icon>
            <h3>公司地址</h3>
            <p>{{ siteSettings.contact_address }}</p>
          </div>
          <div class="contact-card">
            <el-icon :size="32"><Phone /></el-icon>
            <h3>联系电话</h3>
            <p>{{ siteSettings.contact_phone }}</p>
          </div>
          <div class="contact-card">
            <el-icon :size="32"><Message /></el-icon>
            <h3>电子邮箱</h3>
            <p>{{ siteSettings.contact_email }}</p>
          </div>
          <div class="contact-card">
            <el-icon :size="32"><Clock /></el-icon>
            <h3>工作时间</h3>
            <p>{{ siteSettings.work_time }}</p>
          </div>
        </div>
      </div>
    </section>
    
    <SiteFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Location, Phone, Message, Clock, Aim, Shield, Cpu, Headset } from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const siteSettings = ref({
  contact_address: '',
  contact_phone: '',
  contact_email: '',
  work_time: ''
})

const fetchSiteSettings = async () => {
  try {
    const res = await request.get('/api/v1/settings/public')
    if (res?.data) {
      siteSettings.value = { ...siteSettings.value, ...res.data }
    }
  } catch {
    // Use defaults
  }
}

onMounted(() => {
  fetchSiteSettings()
})

const history = ref([
  { year: '2018', title: '公司成立', desc: '公司正式成立，开始研发财务管理核心系统' },
  { year: '2019', title: '产品发布', desc: '首款产品正式上线，获得市场广泛认可' },
  { year: '2020', title: '快速增长', desc: '服务企业突破1000家，完成A轮融资' },
  { year: '2021', title: '技术突破', desc: '获得多项技术专利，产品全面升级' },
  { year: '2022', title: '行业领先', desc: '成为行业领先的财务管理系统提供商' },
  { year: '2023', title: '持续创新', desc: '推出AI智能财务分析功能，服务企业超5000家' }
])

const values = ref([
  { icon: 'Aim', title: '专业专注', desc: '深耕财务管理领域，提供最专业的解决方案' },
  { icon: 'Shield', title: '安全可靠', desc: '采用银行级安全标准，保障数据安全' },
  { icon: 'Cpu', title: '技术创新', desc: '持续技术创新，引领行业发展' },
  { icon: 'Headset', title: '客户至上', desc: '7x24小时技术支持，客户满意是我们的追求' }
])

const team = ref([
  { name: '张明', role: '首席执行官', desc: '10年财务管理行业经验，前知名财务软件公司高管' },
  { name: '李华', role: '首席技术官', desc: '资深技术专家，主导多个大型系统架构设计' },
  { name: '王芳', role: '产品总监', desc: '深入了解企业财务需求，打造用户喜爱的产品' },
  { name: '赵强', role: '技术总监', desc: '全栈工程师，擅长高并发系统设计与优化' }
])
</script>

<style scoped lang="scss">
.about-page {
  padding-top: 64px;
  
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
  }
  
  .hero-section {
    background: linear-gradient(135deg, #1a237e, #0d47a1);
    color: #fff;
    padding: 100px 0;
    text-align: center;
    
    .hero-title {
      font-size: 42px;
      font-weight: 700;
      margin: 0 0 16px;
    }
    
    .hero-desc {
      font-size: 18px;
      opacity: 0.8;
      margin: 0;
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
  
  .company-intro {
    background: #fff;
    
    .intro-content {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 60px;
      align-items: center;
      
      .intro-text {
        p {
          font-size: 15px;
          line-height: 1.8;
          color: #555;
          margin: 0 0 16px;
        }
      }
      
      .intro-stats {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 30px;
        
        .stat-item {
          text-align: center;
          padding: 30px;
          background: #f5f7fa;
          border-radius: 12px;
          
          .stat-number {
            font-size: 36px;
            font-weight: 700;
            color: #409eff;
            margin-bottom: 8px;
          }
          
          .stat-label {
            font-size: 14px;
            color: #666;
          }
        }
      }
    }
  }
  
  .history-section {
    background: #f5f7fa;
    
    .timeline {
      position: relative;
      
      &::before {
        content: '';
        position: absolute;
        left: 50%;
        top: 0;
        bottom: 0;
        width: 2px;
        background: #ddd;
        transform: translateX(-50%);
      }
      
      .timeline-item {
        display: flex;
        align-items: flex-start;
        margin-bottom: 40px;
        position: relative;
        
        &:nth-child(even) {
          flex-direction: row-reverse;
          
          .timeline-content {
            text-align: right;
            padding-left: 0;
            padding-right: 60px;
          }
        }
        
        .timeline-year {
          position: absolute;
          left: 50%;
          transform: translateX(-50%);
          background: #409eff;
          color: #fff;
          padding: 8px 20px;
          border-radius: 20px;
          font-weight: 600;
          z-index: 1;
        }
        
        .timeline-content {
          width: 45%;
          padding-left: 60px;
          
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
    }
  }
  
  .values-section {
    background: #fff;
    
    .values-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 30px;
      
      .value-card {
        text-align: center;
        padding: 40px 24px;
        background: #f5f7fa;
        border-radius: 12px;
        transition: all 0.3s;
        
        &:hover {
          transform: translateY(-5px);
          box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
        }
        
        .value-icon {
          width: 64px;
          height: 64px;
          margin: 0 auto 20px;
          background: linear-gradient(135deg, #409eff, #66b1ff);
          border-radius: 16px;
          display: flex;
          align-items: center;
          justify-content: center;
          color: #fff;
        }
        
        h3 {
          font-size: 18px;
          font-weight: 600;
          color: #1a2332;
          margin: 0 0 12px;
        }
        
        p {
          font-size: 14px;
          color: #666;
          margin: 0;
          line-height: 1.6;
        }
      }
    }
  }
  
  .team-section {
    background: #f5f7fa;
    
    .team-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 30px;
      
      .team-card {
        text-align: center;
        padding: 40px 24px;
        background: #fff;
        border-radius: 12px;
        
        .team-avatar {
          margin-bottom: 20px;
          
          .el-avatar {
            background: linear-gradient(135deg, #409eff, #66b1ff);
            color: #fff;
            font-size: 28px;
          }
        }
        
        h3 {
          font-size: 18px;
          font-weight: 600;
          color: #1a2332;
          margin: 0 0 8px;
        }
        
        .team-role {
          font-size: 14px;
          color: #409eff;
          margin: 0 0 12px;
        }
        
        .team-desc {
          font-size: 13px;
          color: #666;
          margin: 0;
          line-height: 1.6;
        }
      }
    }
  }
  
  .contact-section {
    background: #fff;
    
    .contact-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 30px;
      
      .contact-card {
        text-align: center;
        padding: 40px 24px;
        background: #f5f7fa;
        border-radius: 12px;
        
        .el-icon {
          color: #409eff;
          margin-bottom: 20px;
        }
        
        h3 {
          font-size: 18px;
          font-weight: 600;
          color: #1a2332;
          margin: 0 0 12px;
        }
        
        p {
          font-size: 14px;
          color: #666;
          margin: 0;
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .about-page {
    .company-intro .intro-content {
      grid-template-columns: 1fr;
      gap: 30px;
    }
    
    .history-section .timeline {
      &::before {
        left: 20px;
      }
      
      .timeline-item,
      .timeline-item:nth-child(even) {
        flex-direction: column;
        padding-left: 50px;
        
        .timeline-year {
          left: 20px;
          transform: none;
        }
        
        .timeline-content {
          width: 100%;
          padding: 0;
          text-align: left;
        }
      }
    }
    
    .values-grid,
    .team-grid,
    .contact-grid {
      grid-template-columns: 1fr;
    }
  }
}
</style>
