<template>
  <div class="ddos-overview-page" v-loading="loading">
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(239, 68, 68, 0.1); color: #ef4444;">
          <el-icon :size="24"><Warning /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">今日攻击次数</div>
          <div class="stat-value">{{ stats.attacks_today || 0 }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(16, 185, 129, 0.1); color: #10b981;">
          <el-icon :size="24"><Shield /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">防护IP数量</div>
          <div class="stat-value">{{ stats.protected_ips || 0 }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(59, 130, 246, 0.1); color: #3b82f6;">
          <el-icon :size="24"><DataLine /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">清洗峰值</div>
          <div class="stat-value">{{ stats.peak_traffic || '0 Gbps' }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(168, 85, 247, 0.1); color: #a855f7;">
          <el-icon :size="24"><TrendCharts /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">累计清洗流量</div>
          <div class="stat-value">{{ stats.total_cleaned || '0 TB' }}</div>
        </div>
      </div>
    </div>

    <!-- 流量监控图 -->
    <div class="section-card">
      <div class="section-header">
        <h3>实时流量监控</h3>
        <div class="header-actions">
          <el-radio-group v-model="timeRange" size="small" @change="fetchTrafficData">
            <el-radio-button label="1h">1小时</el-radio-button>
            <el-radio-button label="6h">6小时</el-radio-button>
            <el-radio-button label="24h">24小时</el-radio-button>
            <el-radio-button label="7d">7天</el-radio-button>
          </el-radio-group>
        </div>
      </div>
      <div class="chart-placeholder">
        <div class="chart-bars">
          <div v-for="(item, index) in trafficData" :key="index" class="chart-bar-group">
            <div class="chart-bar" :style="{ height: `${(item.value / maxTraffic) * 100}%` }">
              <div class="bar-fill" :class="{ attack: item.isAttack }"></div>
            </div>
            <span class="bar-label">{{ item.label }}</span>
          </div>
        </div>
        <div class="chart-legend">
          <div class="legend-item"><span class="legend-dot normal"></span> 正常流量</div>
          <div class="legend-item"><span class="legend-dot attack"></span> 攻击流量</div>
        </div>
      </div>
    </div>

    <!-- 最近攻击事件 & 防护状态 -->
    <div class="content-grid">
      <!-- 最近攻击事件 -->
      <div class="section-card">
        <div class="section-header">
          <h3>最近攻击事件</h3>
          <el-button link type="primary" @click="$router.push('/user/ddos/my-ip')">查看全部</el-button>
        </div>
        <div class="attack-list">
          <div v-for="attack in recentAttacks" :key="attack.id" class="attack-item">
            <div class="attack-icon" :class="attack.level">
              <el-icon><Warning /></el-icon>
            </div>
            <div class="attack-info">
              <div class="attack-title">{{ attack.type }}攻击 - {{ attack.target_ip }}</div>
              <div class="attack-meta">
                <span>{{ attack.peak }}峰值</span>
                <span>持续{{ attack.duration }}</span>
                <span>{{ attack.time }}</span>
              </div>
            </div>
            <el-tag :type="attack.status === 'mitigated' ? 'success' : 'danger'" size="small">
              {{ attack.status === 'mitigated' ? '已清洗' : '进行中' }}
            </el-tag>
          </div>
          <el-empty v-if="!recentAttacks.length" description="暂无攻击记录" :image-size="80" />
        </div>
      </div>

      <!-- 防护状态 -->
      <div class="section-card">
        <div class="section-header">
          <h3>防护状态</h3>
        </div>
        <div class="protection-list">
          <div v-for="ip in protectionStatus" :key="ip.ip" class="protection-item">
            <div class="ip-info">
              <div class="ip-address">{{ ip.ip }}</div>
              <div class="ip-domain">{{ ip.domain || '未绑定域名' }}</div>
            </div>
            <div class="ip-status">
              <el-tag :type="ip.protected ? 'success' : 'danger'" size="small">
                {{ ip.protected ? '防护中' : '未防护' }}
              </el-tag>
              <span class="bandwidth">{{ ip.bandwidth }}</span>
            </div>
          </div>
          <el-empty v-if="!protectionStatus.length" description="暂无防护IP" :image-size="80" />
        </div>
      </div>
    </div>

    <!-- 防护套餐信息 -->
    <div class="section-card">
      <div class="section-header">
        <h3>我的防护套餐</h3>
        <el-button type="primary" @click="$router.push('/products/antiddos')">升级套餐</el-button>
      </div>
      <div class="package-info">
        <div class="package-card" v-if="currentPackage">
          <div class="package-header">
            <h4>{{ currentPackage.name }}</h4>
            <el-tag type="success">生效中</el-tag>
          </div>
          <div class="package-details">
            <div class="detail-item">
              <span class="label">防护带宽</span>
              <span class="value">{{ currentPackage.bandwidth }}</span>
            </div>
            <div class="detail-item">
              <span class="label">IP数量</span>
              <span class="value">{{ currentPackage.ip_count }}个</span>
            </div>
            <div class="detail-item">
              <span class="label">到期时间</span>
              <span class="value">{{ currentPackage.expire_date }}</span>
            </div>
            <div class="detail-item">
              <span class="label">月费用</span>
              <span class="value price">¥{{ currentPackage.price }}/月</span>
            </div>
          </div>
          <el-progress :percentage="currentPackage.usage || 0" :stroke-width="8" />
          <div class="usage-label">本月已使用 {{ currentPackage.usage || 0 }}% 流量</div>
        </div>
        <el-empty v-else description="暂无防护套餐" :image-size="80">
          <el-button type="primary" @click="$router.push('/products/antiddos')">立即订购</el-button>
        </el-empty>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Warning, Shield, DataLine, TrendCharts } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const timeRange = ref('24h')

const stats = ref({
  attacks_today: 0,
  protected_ips: 0,
  peak_traffic: '0 Gbps',
  total_cleaned: '0 TB'
})

const trafficData = ref([
  { label: '00:00', value: 120, isAttack: false },
  { label: '02:00', value: 80, isAttack: false },
  { label: '04:00', value: 60, isAttack: false },
  { label: '06:00', value: 90, isAttack: false },
  { label: '08:00', value: 150, isAttack: false },
  { label: '10:00', value: 200, isAttack: false },
  { label: '12:00', value: 350, isAttack: true },
  { label: '14:00', value: 180, isAttack: false },
  { label: '16:00', value: 160, isAttack: false },
  { label: '18:00', value: 140, isAttack: false },
  { label: '20:00', value: 170, isAttack: false },
  { label: '22:00', value: 130, isAttack: false }
])

const maxTraffic = computed(() => {
  return Math.max(...trafficData.value.map(item => item.value), 1)
})

const recentAttacks = ref([
  {
    id: 1,
    type: 'DDoS',
    target_ip: '192.168.1.100',
    peak: '45.6 Gbps',
    duration: '15分钟',
    time: '2024-01-15 14:30',
    status: 'mitigated',
    level: 'high'
  },
  {
    id: 2,
    type: 'CC',
    target_ip: '192.168.1.101',
    peak: '12.3万QPS',
    duration: '8分钟',
    time: '2024-01-15 10:15',
    status: 'mitigated',
    level: 'medium'
  },
  {
    id: 3,
    type: 'SYN Flood',
    target_ip: '192.168.1.100',
    peak: '28.9 Gbps',
    duration: '22分钟',
    time: '2024-01-14 22:45',
    status: 'mitigated',
    level: 'low'
  }
])

const protectionStatus = ref([
  { ip: '192.168.1.100', domain: 'example.com', protected: true, bandwidth: '100Gbps' },
  { ip: '192.168.1.101', domain: 'api.example.com', protected: true, bandwidth: '50Gbps' },
  { ip: '192.168.1.102', domain: '', protected: false, bandwidth: '10Gbps' }
])

const currentPackage = ref({
  name: '高级防护套餐',
  bandwidth: '200Gbps',
  ip_count: 10,
  expire_date: '2024-12-31',
  price: 799,
  usage: 35
})

const fetchTrafficData = async () => {
  try {
    const { data } = await request.get('/api/v1/user/ddos/traffic', { params: { range: timeRange.value } })
    if (data?.data) {
      trafficData.value = data.data
    }
  } catch {
    // 使用默认数据
  }
}

const fetchOverview = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/user/ddos/overview')
    if (data?.data) {
      stats.value = data.data.stats || stats.value
      recentAttacks.value = data.data.recent_attacks || recentAttacks.value
      protectionStatus.value = data.data.protection_status || protectionStatus.value
      currentPackage.value = data.data.current_package || currentPackage.value
    }
  } catch {
    // 使用默认数据
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchOverview()
  fetchTrafficData()
})
</script>

<style scoped lang="scss">
.ddos-overview-page {
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 16px;
    margin-bottom: 16px;

    @media (max-width: 992px) {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  .stat-card {
    background: #fff;
    border-radius: 12px;
    padding: 20px;
    display: flex;
    align-items: center;
    gap: 16px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);

    .stat-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .stat-info {
      .stat-label {
        font-size: 13px;
        color: #86868b;
        margin-bottom: 4px;
      }

      .stat-value {
        font-size: 24px;
        font-weight: 600;
        color: #1d1d1f;
      }
    }
  }

  .section-card {
    background: #fff;
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 16px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);

    .section-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;

      h3 {
        font-size: 16px;
        font-weight: 600;
        margin: 0;
      }
    }
  }

  .chart-placeholder {
    .chart-bars {
      display: flex;
      align-items: flex-end;
      gap: 8px;
      height: 200px;
      padding: 0 10px;
      border-bottom: 1px solid #e5e7eb;
    }

    .chart-bar-group {
      flex: 1;
      display: flex;
      flex-direction: column;
      align-items: center;
      height: 100%;
      justify-content: flex-end;
    }

    .chart-bar {
      width: 100%;
      max-width: 40px;
      border-radius: 4px 4px 0 0;
      overflow: hidden;
      min-height: 4px;

      .bar-fill {
        width: 100%;
        height: 100%;
        background: linear-gradient(to top, #3b82f6, #60a5fa);
        border-radius: 4px 4px 0 0;

        &.attack {
          background: linear-gradient(to top, #ef4444, #f87171);
        }
      }
    }

    .bar-label {
      font-size: 11px;
      color: #9ca3af;
      margin-top: 8px;
    }

    .chart-legend {
      display: flex;
      gap: 20px;
      justify-content: center;
      margin-top: 16px;

      .legend-item {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 13px;
        color: #6b7280;
      }

      .legend-dot {
        width: 10px;
        height: 10px;
        border-radius: 50%;

        &.normal {
          background: #3b82f6;
        }

        &.attack {
          background: #ef4444;
        }
      }
    }
  }

  .content-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
    margin-bottom: 16px;

    @media (max-width: 992px) {
      grid-template-columns: 1fr;
    }
  }

  .attack-list {
    .attack-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px 0;
      border-bottom: 1px solid #f3f4f6;

      &:last-child {
        border-bottom: none;
      }
    }

    .attack-icon {
      width: 36px;
      height: 36px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;

      &.high {
        background: rgba(239, 68, 68, 0.1);
        color: #ef4444;
      }

      &.medium {
        background: rgba(245, 158, 11, 0.1);
        color: #f59e0b;
      }

      &.low {
        background: rgba(59, 130, 246, 0.1);
        color: #3b82f6;
      }
    }

    .attack-info {
      flex: 1;

      .attack-title {
        font-size: 14px;
        font-weight: 500;
        color: #1f2937;
        margin-bottom: 4px;
      }

      .attack-meta {
        display: flex;
        gap: 12px;
        font-size: 12px;
        color: #9ca3af;
      }
    }
  }

  .protection-list {
    .protection-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 0;
      border-bottom: 1px solid #f3f4f6;

      &:last-child {
        border-bottom: none;
      }
    }

    .ip-info {
      .ip-address {
        font-size: 14px;
        font-weight: 600;
        color: #1f2937;
        font-family: monospace;
      }

      .ip-domain {
        font-size: 12px;
        color: #9ca3af;
        margin-top: 2px;
      }
    }

    .ip-status {
      display: flex;
      align-items: center;
      gap: 8px;

      .bandwidth {
        font-size: 12px;
        color: #6b7280;
      }
    }
  }

  .package-info {
    .package-card {
      background: linear-gradient(135deg, #fef2f2 0%, #fff1f2 100%);
      border: 1px solid #fecaca;
      border-radius: 12px;
      padding: 20px;
    }

    .package-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;

      h4 {
        font-size: 18px;
        font-weight: 600;
        color: #1f2937;
        margin: 0;
      }
    }

    .package-details {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 16px;
      margin-bottom: 16px;

      @media (max-width: 768px) {
        grid-template-columns: repeat(2, 1fr);
      }
    }

    .detail-item {
      .label {
        display: block;
        font-size: 13px;
        color: #6b7280;
        margin-bottom: 4px;
      }

      .value {
        font-size: 16px;
        font-weight: 600;
        color: #1f2937;

        &.price {
          color: #ef4444;
        }
      }
    }

    .usage-label {
      text-align: center;
      font-size: 12px;
      color: #6b7280;
      margin-top: 8px;
    }
  }
}
</style>
