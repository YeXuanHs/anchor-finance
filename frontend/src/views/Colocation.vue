<template>
  <div class="colocation-page">
    <!-- Hero Banner -->
    <section class="hero-section">
      <div class="hero-content">
        <div class="hero-badge">服务器托管服务</div>
        <h1 class="hero-title">服务器托管</h1>
        <p class="hero-desc">专业数据中心、稳定电力供应、恒温恒湿环境，为您的服务器提供最佳运行环境</p>
        <div class="hero-actions">
          <el-button type="primary" size="large" round @click="$router.push('/products?group=colocation')">立即咨询</el-button>
          <el-button size="large" round class="hero-ghost-btn" @click="$router.push('/contact')">联系我们</el-button>
        </div>
      </div>
    </section>

    <!-- 核心优势 -->
    <section class="section advantages-section">
      <div class="container">
        <h2 class="section-title text-center">核心优势</h2>
        <p class="section-subtitle text-center">为什么选择我们的托管服务</p>
        <el-skeleton :loading="loading" animated :rows="4">
          <template #default>
            <div class="advantages-grid">
              <div v-for="adv in advantages" :key="adv.title" class="advantage-card">
                <el-icon :size="32" :color="adv.color"><component :is="adv.icon" /></el-icon>
                <h3>{{ adv.title }}</h3>
                <p>{{ adv.description }}</p>
              </div>
            </div>
          </template>
        </el-skeleton>
      </div>
    </section>

    <!-- 数据中心 -->
    <section class="section datacenters-section">
      <div class="container">
        <h2 class="section-title text-center">数据中心</h2>
        <p class="section-subtitle text-center">全球多个数据中心可选</p>
        <el-skeleton :loading="loading" animated :rows="6">
          <template #default>
            <div class="datacenters-grid">
              <div v-for="dc in datacenters" :key="dc.id" class="dc-card">
                <div class="dc-header">
                  <h3>{{ dc.name }}</h3>
                  <el-tag :type="dc.status === 'available' ? 'success' : 'warning'" size="small">
                    {{ dc.status === 'available' ? '可选' : '即将上线' }}
                  </el-tag>
                </div>
                <div class="dc-info">
                  <div class="info-item">
                    <el-icon><Location /></el-icon>
                    <span>{{ dc.location }}</span>
                  </div>
                  <div class="info-item">
                    <el-icon><Connection /></el-icon>
                    <span>{{ dc.bandwidth }}</span>
                  </div>
                  <div class="info-item">
                    <el-icon><Lightning /></el-icon>
                    <span>{{ dc.power }}</span>
                  </div>
                </div>
                <div class="dc-features">
                  <el-tag v-for="feature in dc.features" :key="feature" size="small" type="info">{{ feature }}</el-tag>
                </div>
                <el-button type="primary" link @click="$router.push('/contact')">
                  了解更多 <el-icon><ArrowRight /></el-icon>
                </el-button>
              </div>
            </div>
          </template>
        </el-skeleton>
      </div>
    </section>

    <!-- 服务保障 -->
    <section class="section guarantee-section">
      <div class="container">
        <h2 class="section-title text-center">服务保障</h2>
        <p class="section-subtitle text-center">全方位保障您的服务器安全运行</p>
        <div class="guarantee-grid">
          <div v-for="item in guarantees" :key="item.title" class="guarantee-card">
            <h3>{{ item.value }}</h3>
            <p>{{ item.title }}</p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  Location, Connection, Lightning, ArrowRight,
  OfficeBuilding, Shield, Headset, Timer
} from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(true)

const advantages = ref([
  { title: 'Tier 3+ 数据中心', description: '国际标准Tier 3+级别数据中心', icon: 'OfficeBuilding', color: '#1a56db' },
  { title: '99.99% 可用性', description: '双路供电、UPS、柴油发电机多重保障', icon: 'Shield', color: '#10b981' },
  { title: '7×24 运维', description: '专业运维团队全天候值守', icon: 'Headset', color: '#f59e0b' },
  { title: '快速上架', description: '签合同后24小时内完成上架', icon: 'Timer', color: '#8b5cf6' }
])

const datacenters = ref([
  {
    id: 1,
    name: '华东数据中心',
    location: '上海',
    bandwidth: '100Gbps+ 总带宽',
    power: '双路市电 + UPS',
    status: 'available',
    features: ['BGP多线', '恒温恒湿', '消防系统']
  },
  {
    id: 2,
    name: '华南数据中心',
    location: '广州',
    bandwidth: '80Gbps+ 总带宽',
    power: '双路市电 + UPS',
    status: 'available',
    features: ['BGP多线', '恒温恒湿', '24小时安保']
  },
  {
    id: 3,
    name: '华北数据中心',
    location: '北京',
    bandwidth: '120Gbps+ 总带宽',
    power: '双路市电 + 柴油发电机',
    status: 'available',
    features: ['BGP多线', 'Tier 3+', '金融级安全']
  },
  {
    id: 4,
    name: '香港数据中心',
    location: '香港',
    bandwidth: '50Gbps+ 总带宽',
    power: '双路市电 + UPS',
    status: 'available',
    features: ['国际线路', '低延迟', '免备案']
  }
])

const guarantees = [
  { title: '网络可用性', value: '99.99%' },
  { title: '电力可用性', value: '99.99%' },
  { title: '故障响应', value: '<15分钟' },
  { title: '上架时间', value: '<24小时' }
]

const fetchData = async () => {
  loading.value = true
  try {
    const [advRes, dcRes] = await Promise.allSettled([
      request.get('/api/v1/colocation/advantages'),
      request.get('/api/v1/colocation/datacenters')
    ])

    if (advRes.status === 'fulfilled' && advRes.value.data?.data) {
      advantages.value = advRes.value.data.data
    }
    if (dcRes.status === 'fulfilled' && dcRes.value.data?.data) {
      datacenters.value = dcRes.value.data.data
    }
  } catch (error) {
    console.error('获取数据失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>

<style scoped lang="scss">
.colocation-page {
  min-height: 100vh;
  background: #f8fafc;
}

.hero-section {
  background: linear-gradient(135deg, #1e40af 0%, #1d4ed8 50%, #2563eb 100%);
  padding: 120px 20px 80px;
  text-align: center;

  .hero-content {
    max-width: 700px;
    margin: 0 auto;
  }

  .hero-badge {
    display: inline-block;
    padding: 6px 16px;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 20px;
    color: #fff;
    font-size: 14px;
    margin-bottom: 20px;
  }

  .hero-title {
    font-size: 48px;
    font-weight: 700;
    color: #fff;
    margin: 0 0 16px;
  }

  .hero-desc {
    font-size: 18px;
    color: rgba(255, 255, 255, 0.85);
    margin: 0 0 32px;
  }

  .hero-actions {
    display: flex;
    gap: 16px;
    justify-content: center;
  }

  .hero-ghost-btn {
    background: rgba(255, 255, 255, 0.2) !important;
    border-color: rgba(255, 255, 255, 0.4) !important;
    color: #fff !important;
  }
}

.section {
  padding: 80px 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.section-title {
  font-size: 32px;
  font-weight: 600;
  color: #0f172a;
  margin: 0 0 12px;

  &.text-center {
    text-align: center;
  }
}

.section-subtitle {
  font-size: 16px;
  color: #64748b;
  margin: 0 0 48px;

  &.text-center {
    text-align: center;
  }
}

// 优势
.advantages-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;

  @media (max-width: 768px) {
    grid-template-columns: repeat(2, 1fr);
  }
}

.advantage-card {
  background: #fff;
  border-radius: 16px;
  padding: 28px;
  text-align: center;
  transition: all 0.3s;
  border: 1px solid #e2e8f0;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.08);
  }

  h3 {
    font-size: 18px;
    font-weight: 600;
    color: #0f172a;
    margin: 16px 0 8px;
  }

  p {
    font-size: 14px;
    color: #64748b;
    margin: 0;
  }
}

// 数据中心
.datacenters-section {
  background: #fff;
}

.datacenters-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.dc-card {
  background: #f8fafc;
  border-radius: 16px;
  padding: 28px;
  border: 1px solid #e2e8f0;
  transition: all 0.3s;

  &:hover {
    border-color: #1a56db;
    box-shadow: 0 8px 24px rgba(26, 86, 219, 0.08);
  }
}

.dc-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  h3 {
    font-size: 18px;
    font-weight: 600;
    color: #0f172a;
    margin: 0;
  }
}

.dc-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #475569;

  .el-icon {
    color: #94a3b8;
  }
}

.dc-features {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

// 服务保障
.guarantee-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;

  @media (max-width: 768px) {
    grid-template-columns: repeat(2, 1fr);
  }
}

.guarantee-card {
  background: linear-gradient(135deg, #1a56db, #3b82f6);
  border-radius: 16px;
  padding: 28px;
  text-align: center;
  color: #fff;

  h3 {
    font-size: 32px;
    font-weight: 700;
    margin: 0 0 8px;
  }

  p {
    font-size: 14px;
    opacity: 0.9;
    margin: 0;
  }
}
</style>
