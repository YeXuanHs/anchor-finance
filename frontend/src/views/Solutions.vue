<template>
  <div class="solutions-page">
    <SiteHeader />
    
    <!-- Hero Section -->
    <section class="hero-section">
      <div class="container">
        <h1 class="hero-title">解决方案</h1>
        <p class="hero-desc">针对不同行业需求，提供专业的云计算解决方案</p>
      </div>
    </section>
    
    <!-- Solutions Grid -->
    <section class="section solutions-section">
      <div class="container">
        <div class="solutions-grid">
          <div class="solution-card" v-for="(item, index) in solutions" :key="index" @click="goToSolution(item.type)">
            <div class="solution-icon" :style="{ background: item.gradient }">
              <el-icon :size="36"><component :is="item.icon" /></el-icon>
            </div>
            <h3>{{ item.title }}</h3>
            <p>{{ item.desc }}</p>
            <ul class="solution-features">
              <li v-for="(feature, i) in item.features" :key="i">
                <el-icon><Check /></el-icon>
                {{ feature }}
              </li>
            </ul>
            <el-button type="primary" link>
              了解更多 <el-icon><ArrowRight /></el-icon>
            </el-button>
          </div>
        </div>
      </div>
    </section>
    
    <!-- CTA Section -->
    <section class="section cta-section">
      <div class="container">
        <div class="cta-content">
          <h2>没有找到合适的方案？</h2>
          <p>我们的技术顾问可以为您提供定制化的解决方案</p>
          <div class="cta-actions">
            <el-button type="primary" size="large" round @click="$router.push('/contact')">联系我们</el-button>
            <el-button size="large" round @click="$router.push('/user/tickets')">提交工单</el-button>
          </div>
        </div>
      </div>
    </section>
    
    <SiteFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import request from '@/utils/request'
import { 
  VideoCamera, Monitor, ShoppingBag, Trophy, Shield, OfficeBuilding,
  Connection, Cpu,
  Check, ArrowRight
} from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'

const router = useRouter()

const solutions = ref([
  {
    type: 'game',
    icon: 'Trophy',
    title: '游戏加速',
    desc: '为游戏行业提供专用服务器和加速方案，打造极致游戏体验',
    gradient: 'linear-gradient(135deg, #a55eea, #8854d0)',
    features: ['DDoS防护', '低延迟网络', '全球节点', '7x24运维']
  },
  {
    type: 'video',
    icon: 'VideoCamera',
    title: '视频直播',
    desc: '为视频直播平台提供高性能、低延迟的云服务器',
    gradient: 'linear-gradient(135deg, #ff6b6b, #ee5a24)',
    features: ['大带宽支持', '低延迟传输', '全球CDN加速', '弹性扩容']
  },
  {
    type: 'edu',
    icon: 'Monitor',
    title: '在线教育',
    desc: '为在线教育平台提供流畅的音视频服务',
    gradient: 'linear-gradient(135deg, #45aaf2, #2d98da)',
    features: ['音视频优化', '互动白板', '万人课堂', '数据加密']
  },
  {
    type: 'ecommerce',
    icon: 'ShoppingBag',
    title: '电商平台',
    desc: '为电商平台提供稳定可靠的基础设施',
    gradient: 'linear-gradient(135deg, #26de81, #20bf6b)',
    features: ['弹性扩缩架构', '业务中立', '架构开放', '数据安全']
  },
  {
    type: 'security',
    icon: 'Shield',
    title: '安全防护',
    desc: '为企业提供全方位的安全防护解决方案',
    gradient: 'linear-gradient(135deg, #fd9644, #e67e22)',
    features: ['定制安全方案', '业务数据中立', '专业安全服务', '优质生态体系']
  },
  {
    type: 'caredisaster',
    icon: 'OfficeBuilding',
    title: '容灾备份',
    desc: '在云端构建容灾系统，确保数据安全和业务连续性',
    gradient: 'linear-gradient(135deg, #2bcbba, #0fb9b1)',
    features: ['节约成本', '快速演练', '平台支持丰富', '场景覆盖全面']
  },
  {
    type: 'mixedcloud',
    icon: 'Connection',
    title: '混合云',
    desc: '满足客户特定的安全和合规要求，提供全方位混合云能力',
    gradient: 'linear-gradient(135deg, #4b7bec, #3867d6)',
    features: ['全方位混合云架构', '动态扩展', '云网一体', '多云统一管理']
  },
  {
    type: 'highcalculation',
    icon: 'Cpu',
    title: '高算力',
    desc: '为工业设计、海量数据处理等场景提供卓越的计算服务',
    gradient: 'linear-gradient(135deg, #6c5ce7, #a55eea)',
    features: ['卓越性能体验', '多方面安全防护', '可扩展性强', '快速可获得']
  },
  {
    type: 'website',
    icon: 'Position',
    title: '网站建设',
    desc: '帮助企业快速建站，轻松迈入互联网+时代',
    gradient: 'linear-gradient(135deg, #45aaf2, #2d98da)',
    features: ['弹性扩展', '全景备份', '安全防御', '贴心服务']
  },
  {
    type: 'hosting',
    icon: 'DataLine',
    title: '虚拟主机',
    desc: '即开即用，简单易用，高性价比的网站托管方案',
    gradient: 'linear-gradient(135deg, #fc5c65, #eb3b5a)',
    features: ['即开即用', '简单易用', '安全稳定', '高性价比']
  }
])

const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/solutions')
    if (data?.data?.list?.length) {
      solutions.value = data.data.list
    }
  } catch (e) {
    console.error('Failed to fetch solutions:', e)
  } finally {
    loading.value = false
  }
})

const goToSolution = (type: string) => {
  router.push(`/solutions/${type}`)
}
</script>

<style scoped lang="scss">
.solutions-page {
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
  }
  
  .solutions-section {
    background: #f5f7fa;
    
    .solutions-grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 30px;
    }
    
    .solution-card {
      background: #fff;
      border-radius: 16px;
      padding: 40px 32px;
      cursor: pointer;
      transition: all 0.3s;
      
      &:hover {
        transform: translateY(-8px);
        box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
      }
      
      .solution-icon {
        width: 72px;
        height: 72px;
        border-radius: 20px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;
        margin-bottom: 24px;
      }
      
      h3 {
        font-size: 22px;
        font-weight: 700;
        color: #1a2332;
        margin: 0 0 12px;
      }
      
      p {
        font-size: 14px;
        color: #666;
        margin: 0 0 20px;
        line-height: 1.6;
      }
      
      .solution-features {
        list-style: none;
        padding: 0;
        margin: 0 0 24px;
        
        li {
          display: flex;
          align-items: center;
          gap: 8px;
          font-size: 14px;
          color: #555;
          margin-bottom: 8px;
          
          .el-icon {
            color: #67c23a;
          }
        }
      }
    }
  }
  
  .cta-section {
    background: #fff;
    
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
  .solutions-page .solutions-grid {
    grid-template-columns: 1fr;
  }
}
</style>
