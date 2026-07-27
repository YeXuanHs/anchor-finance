<template>
  <div class="dashboard">
    <!-- Skeleton Loading -->
    <template v-if="loading">
      <!-- Welcome Skeleton -->
      <el-skeleton :rows="2" animated class="skeleton-banner">
        <template #template>
          <div class="banner-skeleton">
            <el-skeleton-item variant="h1" style="width: 40%; height: 28px;" />
            <el-skeleton-item variant="text" style="width: 55%; height: 16px; margin-top: 12px;" />
          </div>
        </template>
      </el-skeleton>

      <!-- Stat Cards Skeleton -->
      <div class="stats-grid">
        <el-card v-for="i in 4" :key="i" class="stat-card-skeleton" shadow="never">
          <el-skeleton :rows="1" animated>
            <template #template>
              <div style="display: flex; align-items: center; gap: 16px;">
                <el-skeleton-item variant="circle" style="width: 48px; height: 48px; flex-shrink: 0;" />
                <div style="flex: 1;">
                  <el-skeleton-item variant="h3" style="width: 70px; height: 28px;" />
                  <el-skeleton-item variant="text" style="width: 90px; height: 14px; margin-top: 6px;" />
                </div>
              </div>
            </template>
          </el-skeleton>
        </el-card>
      </div>

      <!-- Content Skeleton -->
      <div class="content-grid">
        <el-card shadow="never"><el-skeleton :rows="6" animated /></el-card>
        <el-card shadow="never"><el-skeleton :rows="6" animated /></el-card>
      </div>
    </template>

    <!-- Real Content -->
    <template v-else>
      <!-- Welcome Banner -->
      <div class="welcome-banner">
        <div class="welcome-content">
          <h1 class="welcome-title">欢迎回来，{{ username }}</h1>
          <p class="welcome-desc">今天是 {{ currentDate }}，祝您工作愉快！</p>
          <div class="welcome-actions">
            <el-button type="primary" class="welcome-btn" @click="$router.push('/products')">
              <el-icon><ShoppingCart /></el-icon>订购产品
            </el-button>
            <el-button plain class="welcome-btn-outline" @click="$router.push('/user/tickets')">
              <el-icon><ChatLineRound /></el-icon>提交工单
            </el-button>
          </div>
        </div>
        <div class="welcome-illustration">
          <svg viewBox="0 0 200 120" fill="none" width="180" height="110">
            <rect x="20" y="30" width="160" height="80" rx="10" fill="#fff" fill-opacity="0.1" />
            <rect x="30" y="40" width="60" height="8" rx="4" fill="#fff" fill-opacity="0.3" />
            <rect x="30" y="56" width="140" height="4" rx="2" fill="#fff" fill-opacity="0.18" />
            <rect x="30" y="66" width="120" height="4" rx="2" fill="#fff" fill-opacity="0.18" />
            <rect x="30" y="80" width="40" height="20" rx="6" fill="#fff" fill-opacity="0.25" />
            <circle cx="160" cy="50" r="15" fill="#fff" fill-opacity="0.1" />
            <path d="M155 50l5 5 10-10" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill-opacity="0.4" />
          </svg>
        </div>
      </div>

      <!-- Stat Cards -->
      <div class="stats-grid">
        <el-card
          v-for="stat in stats"
          :key="stat.title"
          class="stat-card"
          shadow="never"
          :style="{ '--accent': stat.color }"
        >
          <div class="stat-card-inner">
            <div class="stat-icon" :style="{ background: stat.bg }">
              <el-icon :size="24" :color="stat.color"><component :is="stat.icon" /></el-icon>
            </div>
            <div class="stat-info">
              <span class="stat-value">{{ stat.value }}</span>
              <span class="stat-label">{{ stat.title }}</span>
            </div>
          </div>
        </el-card>
      </div>

      <!-- Content Grid: Cloud Servers + Recent Orders -->
      <div class="content-grid">
        <!-- Cloud Server Overview -->
        <el-card shadow="never" class="section-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">云服务器概览</span>
              <router-link to="/user/products" class="view-all">
                查看全部 <el-icon :size="14"><ArrowRight /></el-icon>
              </router-link>
            </div>
          </template>
          <div class="server-list">
            <div v-for="server in servers" :key="server.id" class="server-item">
              <div class="server-main">
                <div class="server-icon-wrap">
                  <el-icon :size="20" color="#0056FF"><Monitor /></el-icon>
                </div>
                <div class="server-info">
                  <span class="server-name">{{ server.name }}</span>
                  <span class="server-ip">{{ server.ip }}</span>
                </div>
              </div>
              <div class="server-meta">
                <div class="server-spec">{{ server.spec }}</div>
                <el-tag :type="server.statusType" size="small" effect="light" round>{{ server.statusText }}</el-tag>
              </div>
            </div>
          </div>
          <el-empty v-if="servers.length === 0" description="暂无云服务器" :image-size="60">
            <el-button type="primary" size="small" @click="$router.push('/products')">立即购买</el-button>
          </el-empty>
        </el-card>

        <!-- Recent Orders -->
        <el-card shadow="never" class="section-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">最近订单</span>
              <router-link to="/user/orders" class="view-all">
                查看全部 <el-icon :size="14"><ArrowRight /></el-icon>
              </router-link>
            </div>
          </template>
          <el-table :data="recentOrders" stripe style="width: 100%" size="small">
            <el-table-column prop="id" label="订单号" min-width="160" show-overflow-tooltip>
              <template #default="{ row }">
                <span class="mono-text">{{ row.id }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="product" label="产品" min-width="130" show-overflow-tooltip />
            <el-table-column prop="amount" label="金额" width="100">
              <template #default="{ row }">
                <span class="amount-text">¥{{ row.amount }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="90">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small" effect="light" round>
                  {{ row.statusText }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="time" label="时间" width="100" />
          </el-table>
        </el-card>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import {
  Wallet, Box, Tickets, CreditCard, ShoppingCart, ChatLineRound, ArrowRight, Monitor
} from '@element-plus/icons-vue'

const userStore = useUserStore()
const username = computed(() => userStore.username || '用户')
const loading = ref(true)

const currentDate = computed(() => {
  const now = new Date()
  const options: Intl.DateTimeFormatOptions = {
    year: 'numeric', month: 'long', day: 'numeric', weekday: 'long'
  }
  return now.toLocaleDateString('zh-CN', options)
})

interface StatItem {
  title: string
  value: string
  icon: any
  color: string
  bg: string
}

const stats = ref<StatItem[]>([
  { title: '账户余额', value: '¥1,280', icon: CreditCard, color: '#0056FF', bg: 'rgba(0, 86, 255, 0.08)' },
  { title: '运行中服务', value: '3', icon: Box, color: '#52c41a', bg: 'rgba(82, 196, 26, 0.08)' },
  { title: '订单数量', value: '12', icon: Wallet, color: '#fa8c16', bg: 'rgba(250, 140, 22, 0.08)' },
  { title: '待处理工单', value: '1', icon: Tickets, color: '#f5222d', bg: 'rgba(245, 34, 45, 0.08)' }
])

interface Server {
  id: number
  name: string
  ip: string
  spec: string
  statusType: string
  statusText: string
}

const servers = ref<Server[]>([
  { id: 1, name: '香港云服务器', ip: '103.24.xx.1', spec: '2核4G / 50G SSD / 5Mbps', statusType: 'success', statusText: '运行中' },
  { id: 2, name: '美国独立服务器', ip: '198.55.xx.10', spec: 'E5-2680v4 / 64G / 1T', statusType: 'warning', statusText: '初始化' },
  { id: 3, name: '新加坡 VPS', ip: '161.117.xx.5', spec: '2核4G / 60G NVMe', statusType: 'danger', statusText: '已暂停' }
])

const recentOrders = ref([
  { id: 'ORD20260726001', product: '基础财务套餐', amount: '299.00', status: 'active', statusText: '已开通', time: '2026-07-26' },
  { id: 'ORD20260725002', product: '专业税务筹划', amount: '899.00', status: 'pending', statusText: '待支付', time: '2026-07-25' },
  { id: 'ORD20260724003', product: '财务数据分析', amount: '699.00', status: 'completed', statusText: '已完成', time: '2026-07-24' },
  { id: 'ORD20260720004', product: '基础记账服务', amount: '199.00', status: 'cancelled', statusText: '已取消', time: '2026-07-20' }
])

function getStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    active: 'success',
    pending: 'warning',
    completed: 'info',
    cancelled: 'danger'
  }
  return map[status] || 'info'
}

onMounted(() => {
  setTimeout(() => {
    loading.value = false
  }, 600)
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* ==================== Skeleton ==================== */
.skeleton-banner {
  padding: 32px;
  background: linear-gradient(135deg, #1a3a5c 0%, #0056FF 100%);
  border-radius: 12px;
}

.banner-skeleton {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.banner-skeleton :deep(.el-skeleton__item) {
  background: rgba(255, 255, 255, 0.15);
}

.stat-card-skeleton {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.stat-card-skeleton :deep(.el-card__body) {
  padding: 20px;
}

/* ==================== Welcome Banner ==================== */
.welcome-banner {
  background: linear-gradient(135deg, #1a3a5c 0%, #0056FF 50%, #4080FF 100%);
  border-radius: 12px;
  padding: 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  overflow: hidden;
  position: relative;
}

.welcome-banner::before {
  content: '';
  position: absolute;
  top: -60%;
  right: -10%;
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(255,255,255,0.06) 0%, transparent 70%);
  border-radius: 50%;
}

.welcome-content {
  position: relative;
  z-index: 1;
}

.welcome-title {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px 0;
}

.welcome-desc {
  font-size: 14px;
  color: rgba(255,255,255,0.7);
  margin: 0 0 20px 0;
}

.welcome-actions {
  display: flex;
  gap: 12px;
}

.welcome-btn {
  background: rgba(255,255,255,0.15) !important;
  border: 1px solid rgba(255,255,255,0.3) !important;
  color: #fff !important;
  border-radius: 6px;
}

.welcome-btn:hover {
  background: rgba(255,255,255,0.25) !important;
}

.welcome-btn-outline {
  background: transparent !important;
  border: 1px solid rgba(255,255,255,0.3) !important;
  color: #fff !important;
  border-radius: 6px;
}

.welcome-btn-outline:hover {
  background: rgba(255,255,255,0.1) !important;
}

.welcome-illustration {
  position: relative;
  z-index: 1;
}

/* ==================== Stat Cards ==================== */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  transition: all 0.3s ease;
  background: #fff;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 12px rgba(0,0,0,0.08);
}

.stat-card :deep(.el-card__body) {
  padding: 20px;
}

.stat-card-inner {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
}

.stat-label {
  font-size: 13px;
  color: #909399;
}

/* ==================== Content Grid ==================== */
.content-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.section-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  background: #fff;
}

.section-card :deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #e8ecf1;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.view-all {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #0056FF;
  text-decoration: none;
  font-size: 13px;
  transition: color 0.2s;
}

.view-all:hover {
  color: #4080FF;
}

/* ==================== Server List ==================== */
.server-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.server-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 0;
  border-bottom: 1px solid #f2f3f5;
}

.server-item:last-child {
  border-bottom: none;
}

.server-main {
  display: flex;
  align-items: center;
  gap: 12px;
}

.server-icon-wrap {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: rgba(0,86,255,0.06);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.server-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.server-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.server-ip {
  font-size: 12px;
  color: #909399;
  font-family: 'Monaco', 'Menlo', monospace;
}

.server-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.server-spec {
  font-size: 12px;
  color: #909399;
  font-family: 'Monaco', 'Menlo', monospace;
}

/* ==================== Table Styles ==================== */
.mono-text {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 13px;
  color: #606266;
}

.amount-text {
  font-weight: 600;
  color: #303133;
}

/* ==================== Responsive ==================== */
@media (max-width: 1200px) {
  .content-grid {
    grid-template-columns: 1fr;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr 1fr;
  }

  .welcome-illustration {
    display: none;
  }

  .welcome-banner {
    padding: 24px;
  }

  .welcome-title {
    font-size: 20px;
  }

  .server-meta {
    flex-direction: column;
    align-items: flex-end;
    gap: 4px;
  }
}

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
