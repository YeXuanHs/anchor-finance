<template>
  <div class="dashboard">
    <!-- Welcome Banner -->
    <div class="welcome-banner">
      <div class="welcome-content">
        <h1 class="welcome-title">欢迎回来，{{ username }}</h1>
        <p class="welcome-desc">今天是 {{ currentDate }}，祝您工作愉快！</p>
      </div>
      <div class="welcome-illustration">
        <svg viewBox="0 0 200 120" fill="none" width="200" height="120">
          <rect x="20" y="30" width="160" height="80" rx="8" fill="#fff" fill-opacity="0.2" />
          <rect x="30" y="40" width="60" height="8" rx="4" fill="#fff" fill-opacity="0.5" />
          <rect x="30" y="56" width="140" height="4" rx="2" fill="#fff" fill-opacity="0.3" />
          <rect x="30" y="66" width="120" height="4" rx="2" fill="#fff" fill-opacity="0.3" />
          <rect x="30" y="80" width="40" height="20" rx="6" fill="#fff" fill-opacity="0.4" />
          <circle cx="160" cy="50" r="15" fill="#fff" fill-opacity="0.2" />
          <path d="M155 50l5 5 10-10" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill-opacity="0.6" />
        </svg>
      </div>
    </div>

    <!-- Stat Cards -->
    <div class="stats-grid">
      <div v-for="stat in stats" :key="stat.title" class="stat-card" :style="{ '--accent': stat.color }">
        <div class="stat-icon">
          <n-icon :size="28" :component="stat.icon" />
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stat.value }}</span>
          <span class="stat-label">{{ stat.title }}</span>
        </div>
      </div>
    </div>

    <!-- Content Grid -->
    <div class="content-grid">
      <!-- Main Section -->
      <div class="main-section">
        <!-- Recent Orders -->
        <n-card class="section-card" title="最近订单">
          <template #header-extra>
            <router-link to="/user/orders" class="view-all-link">
              查看全部
              <n-icon :size="14" :component="ArrowForwardOutline" />
            </router-link>
          </template>
          <n-data-table
            :columns="orderColumns"
            :data="recentOrders"
            :bordered="false"
            :single-line="false"
            size="small"
          />
        </n-card>

        <!-- Recent Invoices -->
        <n-card class="section-card" title="最近账单">
          <template #header-extra>
            <router-link to="/user/invoices" class="view-all-link">
              查看全部
              <n-icon :size="14" :component="ArrowForwardOutline" />
            </router-link>
          </template>
          <n-data-table
            :columns="invoiceColumns"
            :data="recentInvoices"
            :bordered="false"
            :single-line="false"
            size="small"
          />
        </n-card>
      </div>

      <!-- Side Section -->
      <div class="side-section">
        <!-- Announcements -->
        <n-card class="section-card" title="系统公告">
          <div class="announcement-list">
            <div
              v-for="item in announcements"
              :key="item.id"
              class="announcement-item"
              @click="handleAnnouncement(item)"
            >
              <n-tag :type="item.type" size="small" round>{{ item.tag }}</n-tag>
              <span class="announcement-title">{{ item.title }}</span>
              <span class="announcement-date">{{ item.date }}</span>
            </div>
          </div>
        </n-card>

        <!-- Quick Actions -->
        <n-card class="section-card" title="快捷操作">
          <div class="quick-actions">
            <div
              v-for="action in quickActions"
              :key="action.title"
              class="action-item"
              @click="$router.push(action.path)"
            >
              <div class="action-icon" :style="{ background: action.bg }">
                <n-icon :size="22" :component="action.icon" color="#fff" />
              </div>
              <span class="action-label">{{ action.title }}</span>
            </div>
          </div>
        </n-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import type { DataTableColumns } from 'naive-ui'
import { NTag, NButton } from 'naive-ui'
import {
  CubeOutline,
  WalletOutline,
  ChatbubblesOutline,
  CardOutline,
  ArrowForwardOutline,
  CartOutline,
  ChatboxOutline,
  PersonOutline
} from '@vicons/ionicons5'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const username = computed(() => userStore.username || '用户')

const currentDate = computed(() => {
  const now = new Date()
  const options: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  }
  return now.toLocaleDateString('zh-CN', options)
})

// Stats
const stats = ref([
  {
    title: '我的产品',
    value: '3',
    icon: CubeOutline,
    color: '#1890ff'
  },
  {
    title: '待付账单',
    value: '2',
    icon: WalletOutline,
    color: '#fa8c16'
  },
  {
    title: '活跃工单',
    value: '1',
    icon: ChatbubblesOutline,
    color: '#52c41a'
  },
  {
    title: '账户余额',
    value: '¥1,280',
    icon: CardOutline,
    color: '#722ed1'
  }
])

// Order table columns
const orderColumns: DataTableColumns<any> = [
  { title: '订单号', key: 'id', ellipsis: { tooltip: true }, width: 160 },
  { title: '产品', key: 'product', ellipsis: { tooltip: true } },
  {
    title: '金额',
    key: 'amount',
    width: 100,
    render: (row) => h('span', { style: 'font-weight: 600; color: #262626' }, `¥${row.amount}`)
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: (row) =>
      h(
        NTag,
        { type: getStatusType(row.status), size: 'small', round: true, bordered: false },
        { default: () => row.statusText }
      )
  },
  { title: '时间', key: 'time', width: 110 },
  {
    title: '操作',
    key: 'action',
    width: 100,
    render: (row) =>
      row.status === 'pending'
        ? h(
            NButton,
            { type: 'primary', size: 'tiny', round: true },
            { default: () => '支付' }
          )
        : h(
            NButton,
            { quaternary: true, type: 'primary', size: 'tiny', round: true },
            { default: () => '详情' }
          )
  }
]

const recentOrders = ref([
  { id: 'ORD20260726001', product: '基础财务套餐', amount: '299.00', status: 'active', statusText: '已开通', time: '2026-07-26' },
  { id: 'ORD20260725002', product: '专业税务筹划', amount: '899.00', status: 'pending', statusText: '待支付', time: '2026-07-25' },
  { id: 'ORD20260724003', product: '财务数据分析', amount: '699.00', status: 'completed', statusText: '已完成', time: '2026-07-24' },
  { id: 'ORD20260720004', product: '基础记账服务', amount: '199.00', status: 'cancelled', statusText: '已取消', time: '2026-07-20' },
  { id: 'ORD20260718005', product: '企业财务顾问', amount: '2,999.00', status: 'active', statusText: '已开通', time: '2026-07-18' }
])

// Invoice table columns
const invoiceColumns: DataTableColumns<any> = [
  { title: '账单号', key: 'id', width: 110 },
  { title: '描述', key: 'description', ellipsis: { tooltip: true } },
  {
    title: '金额',
    key: 'amount',
    width: 100,
    render: (row) => h('span', { style: 'font-weight: 600; color: #262626' }, `¥${row.amount}`)
  },
  { title: '到期日', key: 'dueDate', width: 110 },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: (row) =>
      h(
        NTag,
        { type: row.status === 'paid' ? 'success' : 'warning', size: 'small', round: true, bordered: false },
        { default: () => row.statusText }
      )
  },
  {
    title: '操作',
    key: 'action',
    width: 80,
    render: (row) =>
      row.status !== 'paid'
        ? h(
            NButton,
            { type: 'warning', size: 'tiny', round: true },
            { default: () => '支付' }
          )
        : h('span', { style: 'color: #bfbfbf; font-size: 12px' }, '—')
  }
]

const recentInvoices = ref([
  { id: 'INV2026072601', description: '基础财务套餐-月度续费', amount: '299.00', status: 'unpaid', statusText: '待支付', dueDate: '2026-08-15' },
  { id: 'INV2026072502', description: '专业税务筹划-季度服务', amount: '899.00', status: 'unpaid', statusText: '待支付', dueDate: '2026-08-01' },
  { id: 'INV2026072003', description: '财务数据分析-月度服务', amount: '699.00', status: 'paid', statusText: '已支付', dueDate: '2026-07-20' }
])

// Announcements
const announcements = ref([
  { id: 1, title: '系统升级通知：2026年8月1日凌晨维护', tag: '系统', type: 'warning' as const, date: '07-26' },
  { id: 2, title: '新功能上线：支持微信支付和支付宝', tag: '功能', type: 'success' as const, date: '07-24' },
  { id: 3, title: '财务报表导出功能优化完成', tag: '更新', type: 'info' as const, date: '07-20' },
  { id: 4, title: '安全提醒：请定期修改密码', tag: '安全', type: 'error' as const, date: '07-18' }
])

// Quick actions
const quickActions = ref([
  { title: '购买产品', icon: CartOutline, path: '/products', bg: 'linear-gradient(135deg, #1890ff, #096dd9)' },
  { title: '提交工单', icon: ChatboxOutline, path: '/user/tickets', bg: 'linear-gradient(135deg, #52c41a, #389e0d)' },
  { title: '充值余额', icon: WalletOutline, path: '/user/wallet', bg: 'linear-gradient(135deg, #fa8c16, #d46b08)' },
  { title: '个人资料', icon: PersonOutline, path: '/user/profile', bg: 'linear-gradient(135deg, #8c8c8c, #595959)' }
])

function getStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'error'> = {
    active: 'success',
    pending: 'warning',
    completed: 'info',
    cancelled: 'error'
  }
  return map[status] || 'info'
}

function handleAnnouncement(item: any) {
  // TODO: open announcement detail
}
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* ==================== Welcome Banner ==================== */
.welcome-banner {
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 50%, #0050b3 100%);
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
  top: -50%;
  right: -10%;
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
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
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
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
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  border: 1px solid #f0f0f0;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: var(--accent);
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  border-color: transparent;
}

.stat-icon {
  width: 52px;
  height: 52px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: color-mix(in srgb, var(--accent) 10%, transparent);
  color: var(--accent);
  flex-shrink: 0;
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #262626;
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: #8c8c8c;
}

/* ==================== Content Grid ==================== */
.content-grid {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 24px;
}

.main-section {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.side-section {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.section-card {
  border-radius: 12px;
  border: 1px solid #f0f0f0;
}

.section-card :deep(.n-card-header) {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.section-card :deep(.n-card__content) {
  padding: 16px 20px;
}

.view-all-link {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #1890ff;
  text-decoration: none;
  font-size: 13px;
  transition: color 0.2s;
}

.view-all-link:hover {
  color: #40a9ff;
}

/* ==================== Announcements ==================== */
.announcement-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.announcement-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.announcement-item:hover {
  background: #f0f5ff;
}

.announcement-title {
  flex: 1;
  font-size: 13px;
  color: #262626;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.announcement-date {
  font-size: 12px;
  color: #bfbfbf;
  flex-shrink: 0;
}

/* ==================== Quick Actions ==================== */
.quick-actions {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px 8px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  background: #fafafa;
}

.action-item:hover {
  background: #f0f5ff;
  transform: translateY(-2px);
}

.action-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-label {
  font-size: 13px;
  color: #595959;
  font-weight: 500;
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
    grid-template-columns: 1fr;
  }

  .welcome-illustration {
    display: none;
  }

  .quick-actions {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
