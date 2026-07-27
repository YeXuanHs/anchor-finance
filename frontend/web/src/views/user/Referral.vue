<template>
  <div class="referral-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">推荐返利</h1>
        <p class="page-desc">邀请好友注册，赚取丰厚佣金奖励</p>
      </div>
      <div class="header-illustration">
        <svg viewBox="0 0 200 120" fill="none" width="200" height="120">
          <rect x="20" y="30" width="160" height="80" rx="8" fill="#fff" fill-opacity="0.2" />
          <circle cx="70" cy="55" r="12" fill="#fff" fill-opacity="0.4" />
          <circle cx="130" cy="55" r="12" fill="#fff" fill-opacity="0.4" />
          <path d="M82 55h36" stroke="#fff" stroke-width="2" stroke-dasharray="4 4" stroke-opacity="0.6" />
          <rect x="40" y="80" width="120" height="6" rx="3" fill="#fff" fill-opacity="0.3" />
          <path d="M95 45l5 5 10-10" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill-opacity="0.6" />
        </svg>
      </div>
    </div>

    <!-- Referral Link -->
    <n-card class="link-card" :bordered="false">
      <div class="link-section">
        <div class="link-info">
          <div class="link-label">您的专属推荐链接</div>
          <div class="link-url">
            <n-input
              v-model:value="referralLink"
              readonly
              size="large"
              round
              class="link-input"
            >
              <template #prefix>
                <n-icon :component="LinkOutline" />
              </template>
            </n-input>
            <n-button type="primary" size="large" round @click="copyLink">
              <template #icon>
                <n-icon :component="CopyOutline" />
              </template>
              复制链接
            </n-button>
          </div>
          <div class="link-tip">分享此链接给好友，好友注册并完成首单后您将获得佣金奖励</div>
        </div>
        <div class="share-buttons">
          <n-button round>
            <template #icon>
              <n-icon :component="LogoWechat" />
            </template>
            微信分享
          </n-button>
          <n-button round>
            <template #icon>
              <n-icon :component="ShareSocialOutline" />
            </template>
            更多分享
          </n-button>
        </div>
      </div>
    </n-card>

    <!-- Statistics -->
    <div class="stats-grid">
      <div class="stat-card" style="--accent: #1890ff">
        <div class="stat-icon">
          <n-icon :size="28" :component="PeopleOutline" />
        </div>
        <div class="stat-info">
          <n-statistic :value="stats.totalReferrals" />
          <span class="stat-label">推荐人数</span>
        </div>
      </div>
      <div class="stat-card" style="--accent: #52c41a">
        <div class="stat-icon">
          <n-icon :size="28" :component="CashOutline" />
        </div>
        <div class="stat-info">
          <n-statistic :value="stats.totalCommission" prefix="¥" />
          <span class="stat-label">累计佣金</span>
        </div>
      </div>
      <div class="stat-card" style="--accent: #fa8c16">
        <div class="stat-icon">
          <n-icon :size="28" :component="WalletOutline" />
        </div>
        <div class="stat-info">
          <n-statistic :value="stats.withdrawable" prefix="¥" />
          <span class="stat-label">可提现金额</span>
        </div>
      </div>
      <div class="stat-card" style="--accent: #722ed1">
        <div class="stat-icon">
          <n-icon :size="28" :component="TrophyOutline" />
        </div>
        <div class="stat-info">
          <n-statistic :value="stats.pendingCommission" prefix="¥" />
          <span class="stat-label">待结算佣金</span>
        </div>
      </div>
    </div>

    <!-- Withdraw & Table -->
    <div class="content-grid">
      <!-- Withdraw Card -->
      <n-card class="withdraw-card" title="佣金提现" :bordered="false">
        <template #header-extra>
          <n-tag type="success" size="small" round>实时到账</n-tag>
        </template>
        <div class="withdraw-form">
          <div class="withdraw-balance">
            <span class="balance-label">可提现余额</span>
            <span class="balance-value">¥{{ stats.withdrawable.toFixed(2) }}</span>
          </div>
          <n-input-number
            v-model:value="withdrawAmount"
            :max="stats.withdrawable"
            :min="1"
            placeholder="请输入提现金额"
            size="large"
            round
          >
            <template #prefix>¥</template>
          </n-input-number>
          <n-button
            type="primary"
            block
            round
            size="large"
            :disabled="!withdrawAmount || withdrawAmount <= 0"
            :loading="withdrawing"
            @click="handleWithdraw"
          >
            申请提现
          </n-button>
          <div class="withdraw-tips">
            <div class="tip-item">
              <n-icon :size="14" :component="InformationCircleOutline" />
              <span>最低提现金额 ¥1.00</span>
            </div>
            <div class="tip-item">
              <n-icon :size="14" :component="InformationCircleOutline" />
              <span>提现将在 1-3 个工作日内到账</span>
            </div>
          </div>
        </div>
      </n-card>

      <!-- Referral Rules -->
      <n-card class="rules-card" title="推荐规则" :bordered="false">
        <div class="rules-list">
          <div class="rule-item">
            <div class="rule-step">1</div>
            <div class="rule-content">
              <div class="rule-title">分享链接</div>
              <div class="rule-desc">将您的专属推荐链接分享给好友</div>
            </div>
          </div>
          <div class="rule-item">
            <div class="rule-step">2</div>
            <div class="rule-content">
              <div class="rule-title">好友注册</div>
              <div class="rule-desc">好友通过您的链接完成注册</div>
            </div>
          </div>
          <div class="rule-item">
            <div class="rule-step">3</div>
            <div class="rule-content">
              <div class="rule-title">获得佣金</div>
              <div class="rule-desc">好友完成首单后，您获得订单金额 10% 的佣金</div>
            </div>
          </div>
          <div class="rule-item">
            <div class="rule-step">4</div>
            <div class="rule-content">
              <div class="rule-title">申请提现</div>
              <div class="rule-desc">佣金可随时申请提现至账户余额</div>
            </div>
          </div>
        </div>
      </n-card>
    </div>

    <!-- Referral Records -->
    <n-card class="records-card" title="推荐记录" :bordered="false">
      <template #header-extra>
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索被推荐人"
          size="small"
          round
          clearable
          style="width: 200px"
        >
          <template #prefix>
            <n-icon :component="SearchOutline" />
          </template>
        </n-input>
      </template>
      <n-data-table
        :columns="recordColumns"
        :data="filteredRecords"
        :bordered="false"
        :single-line="false"
        :pagination="pagination"
        size="small"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import { NTag, NButton, useMessage } from 'naive-ui'
import {
  LinkOutline,
  CopyOutline,
  ShareSocialOutline,
  PeopleOutline,
  CashOutline,
  WalletOutline,
  TrophyOutline,
  InformationCircleOutline,
  SearchOutline,
  LogoWechat
} from '@vicons/ionicons5'

const message = useMessage()

const referralLink = ref('https://example.com/ref/USER2026ABC')
const withdrawAmount = ref<number | null>(null)
const withdrawing = ref(false)
const searchKeyword = ref('')

// Statistics
const stats = ref({
  totalReferrals: 28,
  totalCommission: 3680.00,
  withdrawable: 2150.00,
  pendingCommission: 530.00
})

// Referral records
interface ReferralRecord {
  id: string
  username: string
  email: string
  registerDate: string
  firstOrderDate: string
  orderAmount: number
  commission: number
  status: 'settled' | 'pending' | 'cancelled'
  statusText: string
}

const records = ref<ReferralRecord[]>([
  {
    id: 'REF001',
    username: '张三',
    email: 'zhang***@example.com',
    registerDate: '2026-07-20',
    firstOrderDate: '2026-07-21',
    orderAmount: 2999.00,
    commission: 299.90,
    status: 'settled',
    statusText: '已结算'
  },
  {
    id: 'REF002',
    username: '李四',
    email: 'li***@example.com',
    registerDate: '2026-07-18',
    firstOrderDate: '2026-07-19',
    orderAmount: 1599.00,
    commission: 159.90,
    status: 'settled',
    statusText: '已结算'
  },
  {
    id: 'REF003',
    username: '王五',
    email: 'wang***@example.com',
    registerDate: '2026-07-15',
    firstOrderDate: '2026-07-16',
    orderAmount: 899.00,
    commission: 89.90,
    status: 'settled',
    statusText: '已结算'
  },
  {
    id: 'REF004',
    username: '赵六',
    email: 'zhao***@example.com',
    registerDate: '2026-07-25',
    firstOrderDate: '2026-07-26',
    orderAmount: 3599.00,
    commission: 359.90,
    status: 'pending',
    statusText: '待结算'
  },
  {
    id: 'REF005',
    username: '孙七',
    email: 'sun***@example.com',
    registerDate: '2026-07-22',
    firstOrderDate: '2026-07-23',
    orderAmount: 1999.00,
    commission: 199.90,
    status: 'pending',
    statusText: '待结算'
  },
  {
    id: 'REF006',
    username: '周八',
    email: 'zhou***@example.com',
    registerDate: '2026-07-10',
    firstOrderDate: '2026-07-12',
    orderAmount: 599.00,
    commission: 59.90,
    status: 'settled',
    statusText: '已结算'
  },
  {
    id: 'REF007',
    username: '吴九',
    email: 'wu***@example.com',
    registerDate: '2026-07-08',
    firstOrderDate: '-',
    orderAmount: 0,
    commission: 0,
    status: 'cancelled',
    statusText: '未下单'
  },
  {
    id: 'REF008',
    username: '郑十',
    email: 'zheng***@example.com',
    registerDate: '2026-06-28',
    firstOrderDate: '2026-06-30',
    orderAmount: 4299.00,
    commission: 429.90,
    status: 'settled',
    statusText: '已结算'
  }
])

const filteredRecords = computed(() => {
  if (!searchKeyword.value) return records.value
  return records.value.filter(
    (r) =>
      r.username.includes(searchKeyword.value) ||
      r.email.includes(searchKeyword.value)
  )
})

const pagination: PaginationProps = {
  pageSize: 10
}

// Table columns
const recordColumns: DataTableColumns<ReferralRecord> = [
  { title: 'ID', key: 'id', width: 100 },
  { title: '被推荐人', key: 'username', width: 100 },
  { title: '邮箱', key: 'email', ellipsis: { tooltip: true } },
  { title: '注册时间', key: 'registerDate', width: 120 },
  { title: '首单时间', key: 'firstOrderDate', width: 120 },
  {
    title: '订单金额',
    key: 'orderAmount',
    width: 120,
    render: (row) =>
      h('span', { style: 'font-weight: 600; color: #262626' }, row.orderAmount > 0 ? `¥${row.orderAmount.toFixed(2)}` : '-')
  },
  {
    title: '佣金',
    key: 'commission',
    width: 120,
    render: (row) =>
      h('span', { style: 'font-weight: 600; color: #52c41a' }, row.commission > 0 ? `¥${row.commission.toFixed(2)}` : '-')
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row) => {
      const typeMap: Record<string, 'success' | 'warning' | 'default'> = {
        settled: 'success',
        pending: 'warning',
        cancelled: 'default'
      }
      return h(
        NTag,
        { type: typeMap[row.status] || 'default', size: 'small', round: true, bordered: false },
        { default: () => row.statusText }
      )
    }
  }
]

function copyLink() {
  navigator.clipboard.writeText(referralLink.value).then(() => {
    message.success('推荐链接已复制到剪贴板')
  }).catch(() => {
    message.warning('复制失败，请手动复制')
  })
}

function handleWithdraw() {
  if (!withdrawAmount.value || withdrawAmount.value <= 0) {
    message.warning('请输入有效的提现金额')
    return
  }
  withdrawing.value = true
  setTimeout(() => {
    withdrawing.value = false
    stats.value.withdrawable -= withdrawAmount.value!
    message.success(`提现申请已提交，金额 ¥${withdrawAmount.value!.toFixed(2)}`)
    withdrawAmount.value = null
  }, 1500)
}
</script>

<style scoped>
.referral-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* ==================== Page Header ==================== */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 50%, #0050b3 100%);
  border-radius: 12px;
  padding: 32px;
  position: relative;
  overflow: hidden;
}

.page-header::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -10%;
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
  border-radius: 50%;
}

.header-content {
  position: relative;
  z-index: 1;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px 0;
}

.page-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
}

.header-illustration {
  position: relative;
  z-index: 1;
}

/* ==================== Link Card ==================== */
.link-card {
  border-radius: 12px;
  border: 1px solid #f0f0f0;
  background: linear-gradient(135deg, #f0f5ff 0%, #e6f0ff 100%);
}

.link-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.link-label {
  font-size: 14px;
  font-weight: 600;
  color: #262626;
}

.link-url {
  display: flex;
  gap: 12px;
}

.link-input {
  flex: 1;
}

.link-tip {
  font-size: 12px;
  color: #8c8c8c;
}

.share-buttons {
  display: flex;
  gap: 12px;
}

/* ==================== Statistics ==================== */
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

.stat-label {
  font-size: 14px;
  color: #8c8c8c;
}

/* ==================== Content Grid ==================== */
.content-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.withdraw-card,
.rules-card {
  border-radius: 12px;
  border: 1px solid #f0f0f0;
}

/* ==================== Withdraw Form ==================== */
.withdraw-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.withdraw-balance {
  display: flex;
  align-items: baseline;
  gap: 8px;
  padding: 16px;
  background: linear-gradient(135deg, #f0f5ff 0%, #e6f0ff 100%);
  border-radius: 8px;
}

.balance-label {
  font-size: 14px;
  color: #8c8c8c;
}

.balance-value {
  font-size: 28px;
  font-weight: 700;
  color: #1890ff;
}

.withdraw-tips {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tip-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #8c8c8c;
}

/* ==================== Rules ==================== */
.rules-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.rule-item {
  display: flex;
  gap: 16px;
  align-items: flex-start;
}

.rule-step {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
}

.rule-content {
  flex: 1;
}

.rule-title {
  font-size: 14px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 4px;
}

.rule-desc {
  font-size: 13px;
  color: #8c8c8c;
}

/* ==================== Records Card ==================== */
.records-card {
  border-radius: 12px;
  border: 1px solid #f0f0f0;
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
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .header-illustration {
    display: none;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .link-url {
    flex-direction: column;
  }

  .share-buttons {
    flex-direction: column;
  }
}
</style>
