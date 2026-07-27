<template>
  <div class="coupons-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">我的优惠券</h1>
        <p class="page-desc">管理您的优惠券，查看使用记录</p>
      </div>
      <n-button type="primary" round @click="showClaimModal = true">
        <template #icon>
          <n-icon :component="AddOutline" />
        </template>
        领取优惠券
      </n-button>
    </div>

    <!-- Statistics -->
    <div class="stats-row">
      <div class="stat-card" style="--accent: #1890ff">
        <div class="stat-icon">
          <n-icon :size="28" :component="TicketOutline" />
        </div>
        <div class="stat-info">
          <n-statistic :value="stats.available" />
          <span class="stat-label">可用优惠券</span>
        </div>
      </div>
      <div class="stat-card" style="--accent: #52c41a">
        <div class="stat-icon">
          <n-icon :size="28" :component="CheckmarkCircleOutline" />
        </div>
        <div class="stat-info">
          <n-statistic :value="stats.used" />
          <span class="stat-label">已使用</span>
        </div>
      </div>
      <div class="stat-card" style="--accent: #fa8c16">
        <div class="stat-icon">
          <n-icon :size="28" :component="TimeOutline" />
        </div>
        <div class="stat-info">
          <n-statistic :value="stats.expired" />
          <span class="stat-label">已过期</span>
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <n-card class="content-card" :bordered="false">
      <n-tabs v-model:value="activeTab" type="line" animated>
        <n-tab-pane name="available" tab="可用优惠券">
          <div v-if="availableCoupons.length" class="coupon-grid">
            <n-grid :cols="3" :x-gap="16" :y-gap="16" responsive="screen">
              <n-gi v-for="coupon in availableCoupons" :key="coupon.id">
                <div class="coupon-card" :class="[coupon.type]">
                  <div class="coupon-left">
                    <div class="coupon-amount">
                      <span v-if="coupon.discountType === 'fixed'" class="amount-symbol">¥</span>
                      <span class="amount-value">{{ coupon.value }}</span>
                      <span v-if="coupon.discountType === 'percent'" class="amount-symbol">%</span>
                    </div>
                    <div class="coupon-condition">
                      满{{ coupon.minAmount }}可用
                    </div>
                  </div>
                  <div class="coupon-divider">
                    <div class="divider-circle top"></div>
                    <div class="divider-line"></div>
                    <div class="divider-circle bottom"></div>
                  </div>
                  <div class="coupon-right">
                    <div class="coupon-name">{{ coupon.name }}</div>
                    <div class="coupon-desc">{{ coupon.description }}</div>
                    <div class="coupon-meta">
                      <n-tag size="small" :type="getCouponTagType(coupon.category)" round>
                        {{ coupon.category }}
                      </n-tag>
                      <span class="coupon-expire">{{ coupon.expireDate }} 到期</span>
                    </div>
                    <div class="coupon-actions">
                      <n-button size="tiny" type="primary" round @click="copyCouponCode(coupon.code)">
                        复制券码
                      </n-button>
                      <n-button size="tiny" quaternary type="primary" round @click="useCoupon(coupon)">
                        立即使用
                      </n-button>
                    </div>
                  </div>
                </div>
              </n-gi>
            </n-grid>
          </div>
          <n-empty v-else description="暂无可用优惠券" />
        </n-tab-pane>

        <n-tab-pane name="used" tab="使用记录">
          <n-data-table
            :columns="usedColumns"
            :data="usedCoupons"
            :bordered="false"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <n-tab-pane name="expired" tab="已过期">
          <n-data-table
            :columns="expiredColumns"
            :data="expiredCoupons"
            :bordered="false"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>
      </n-tabs>
    </n-card>

    <!-- Claim Modal -->
    <n-modal v-model:show="showClaimModal" preset="card" title="领取优惠券" style="width: 480px">
      <div class="claim-form">
        <n-input
          v-model:value="claimCode"
          placeholder="请输入优惠券码"
          size="large"
          round
        />
        <n-button type="primary" block round :loading="claiming" @click="handleClaim">
          立即领取
        </n-button>
      </div>
      <div class="available-claim-list">
        <div class="claim-section-title">可领取的优惠券</div>
        <div
          v-for="item in claimableList"
          :key="item.id"
          class="claim-item"
          @click="handleQuickClaim(item)"
        >
          <div class="claim-amount">
            <span v-if="item.discountType === 'fixed'">¥</span>{{ item.value }}
            <span v-if="item.discountType === 'percent'">%</span>
          </div>
          <div class="claim-info">
            <div class="claim-name">{{ item.name }}</div>
            <div class="claim-condition">满{{ item.minAmount }}可用</div>
          </div>
          <n-button size="tiny" type="primary" round>领取</n-button>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import type { DataTableColumns } from 'naive-ui'
import { NTag, NButton, useMessage } from 'naive-ui'
import {
  TicketOutline,
  CheckmarkCircleOutline,
  TimeOutline,
  AddOutline,
  CopyOutline
} from '@vicons/ionicons5'

const message = useMessage()

const activeTab = ref('available')
const showClaimModal = ref(false)
const claimCode = ref('')
const claiming = ref(false)

// Statistics
const stats = ref({
  available: 5,
  used: 12,
  expired: 3
})

// Coupon type
interface Coupon {
  id: string
  name: string
  code: string
  type: string
  discountType: 'fixed' | 'percent'
  value: number
  minAmount: number
  description: string
  category: string
  expireDate: string
  status: 'available' | 'used' | 'expired'
  usedDate?: string
}

// Available coupons
const availableCoupons = ref<Coupon[]>([
  {
    id: '1',
    name: '新用户专享券',
    code: 'NEW2026',
    type: 'primary',
    discountType: 'fixed',
    value: 50,
    minAmount: 200,
    description: '新用户首单专享优惠',
    category: '通用',
    expireDate: '2026-08-31',
    status: 'available'
  },
  {
    id: '2',
    name: '财务服务折扣券',
    code: 'FINANCE20',
    type: 'success',
    discountType: 'percent',
    value: 20,
    minAmount: 500,
    description: '财务分析服务8折优惠',
    category: '财务',
    expireDate: '2026-09-15',
    status: 'available'
  },
  {
    id: '3',
    name: '满减优惠券',
    code: 'SAVE100',
    type: 'warning',
    discountType: 'fixed',
    value: 100,
    minAmount: 800,
    description: '全场满800减100',
    category: '通用',
    expireDate: '2026-10-01',
    status: 'available'
  },
  {
    id: '4',
    name: '税务筹划专享',
    code: 'TAX50',
    type: 'info',
    discountType: 'fixed',
    value: 50,
    minAmount: 300,
    description: '税务筹划服务专享',
    category: '税务',
    expireDate: '2026-08-20',
    status: 'available'
  },
  {
    id: '5',
    name: '会员日特惠',
    code: 'VIP30',
    type: 'error',
    discountType: 'percent',
    value: 30,
    minAmount: 1000,
    description: '会员日全场7折',
    category: '限时',
    expireDate: '2026-08-05',
    status: 'available'
  }
])

// Used coupons
const usedCoupons = ref<Coupon[]>([
  {
    id: '101',
    name: '首次注册奖励',
    code: 'WELCOME',
    type: 'primary',
    discountType: 'fixed',
    value: 30,
    minAmount: 100,
    description: '新用户注册奖励',
    category: '通用',
    expireDate: '2026-07-15',
    status: 'used',
    usedDate: '2026-07-10'
  },
  {
    id: '102',
    name: '618大促券',
    code: 'PROMO618',
    type: 'warning',
    discountType: 'fixed',
    value: 200,
    minAmount: 1500,
    description: '618大促满减',
    category: '限时',
    expireDate: '2026-06-30',
    status: 'used',
    usedDate: '2026-06-18'
  }
])

// Expired coupons
const expiredCoupons = ref<Coupon[]>([
  {
    id: '201',
    name: '春节特惠券',
    code: 'CNY2026',
    type: 'error',
    discountType: 'fixed',
    value: 88,
    minAmount: 500,
    description: '春节限时优惠',
    category: '限时',
    expireDate: '2026-02-28',
    status: 'expired'
  },
  {
    id: '202',
    name: '体验金',
    code: 'TRIAL',
    type: 'info',
    discountType: 'fixed',
    value: 20,
    minAmount: 50,
    description: '新服务体验金',
    category: '通用',
    expireDate: '2026-03-31',
    status: 'expired'
  }
])

// Claimable list
const claimableList = ref([
  { id: 'c1', name: '夏日清凉券', value: 30, discountType: 'fixed' as const, minAmount: 150 },
  { id: 'c2', name: '推荐返利券', value: 15, discountType: 'percent' as const, minAmount: 300 },
  { id: 'c3', name: '续费专享', value: 50, discountType: 'fixed' as const, minAmount: 400 }
])

// Used table columns
const usedColumns: DataTableColumns<Coupon> = [
  { title: '券码', key: 'code', width: 120 },
  { title: '名称', key: 'name', ellipsis: { tooltip: true } },
  {
    title: '面值',
    key: 'value',
    width: 100,
    render: (row) =>
      h('span', { style: 'font-weight: 600; color: #262626' },
        row.discountType === 'fixed' ? `¥${row.value}` : `${row.value}%`
      )
  },
  { title: '使用条件', key: 'minAmount', width: 120, render: (row) => h('span', {}, `满${row.minAmount}可用`) },
  { title: '使用日期', key: 'usedDate', width: 120 },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: () => h(NTag, { type: 'success', size: 'small', round: true, bordered: false }, { default: () => '已使用' })
  }
]

// Expired table columns
const expiredColumns: DataTableColumns<Coupon> = [
  { title: '券码', key: 'code', width: 120 },
  { title: '名称', key: 'name', ellipsis: { tooltip: true } },
  {
    title: '面值',
    key: 'value',
    width: 100,
    render: (row) =>
      h('span', { style: 'font-weight: 600; color: #bfbfbf' },
        row.discountType === 'fixed' ? `¥${row.value}` : `${row.value}%`
      )
  },
  { title: '过期日期', key: 'expireDate', width: 120 },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: () => h(NTag, { type: 'default', size: 'small', round: true, bordered: false }, { default: () => '已过期' })
  }
]

function getCouponTagType(category: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'error' | 'default'> = {
    '通用': 'info',
    '财务': 'success',
    '税务': 'warning',
    '限时': 'error'
  }
  return map[category] || 'default'
}

function copyCouponCode(code: string) {
  navigator.clipboard.writeText(code).then(() => {
    message.success('优惠券码已复制')
  }).catch(() => {
    message.warning('复制失败，请手动复制')
  })
}

function useCoupon(coupon: Coupon) {
  message.info(`跳转到使用优惠券: ${coupon.name}`)
}

function handleClaim() {
  if (!claimCode.value.trim()) {
    message.warning('请输入优惠券码')
    return
  }
  claiming.value = true
  setTimeout(() => {
    claiming.value = false
    message.success('优惠券领取成功')
    claimCode.value = ''
    showClaimModal.value = false
  }, 1000)
}

function handleQuickClaim(item: { id: string; name: string }) {
  message.success(`已领取: ${item.name}`)
}
</script>

<style scoped>
.coupons-page {
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

/* ==================== Statistics ==================== */
.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
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

/* ==================== Content Card ==================== */
.content-card {
  border-radius: 12px;
  border: 1px solid #f0f0f0;
}

/* ==================== Coupon Grid ==================== */
.coupon-card {
  display: flex;
  border-radius: 12px;
  background: #fff;
  border: 1px solid #f0f0f0;
  overflow: hidden;
  transition: all 0.3s ease;
  min-height: 140px;
}

.coupon-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  border-color: transparent;
}

.coupon-card.primary {
  border-left: 4px solid #1890ff;
}

.coupon-card.success {
  border-left: 4px solid #52c41a;
}

.coupon-card.warning {
  border-left: 4px solid #fa8c16;
}

.coupon-card.info {
  border-left: 4px solid #722ed1;
}

.coupon-card.error {
  border-left: 4px solid #f5222d;
}

.coupon-left {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 16px;
  min-width: 100px;
  background: linear-gradient(135deg, #f0f5ff 0%, #e6f0ff 100%);
}

.coupon-amount {
  display: flex;
  align-items: baseline;
  color: #1890ff;
}

.amount-symbol {
  font-size: 16px;
  font-weight: 600;
}

.amount-value {
  font-size: 32px;
  font-weight: 700;
  line-height: 1;
}

.coupon-condition {
  font-size: 12px;
  color: #8c8c8c;
  margin-top: 4px;
}

.coupon-divider {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 1px;
  position: relative;
}

.divider-line {
  flex: 1;
  width: 1px;
  background: repeating-linear-gradient(
    to bottom,
    #f0f0f0 0px,
    #f0f0f0 4px,
    transparent 4px,
    transparent 8px
  );
}

.divider-circle {
  width: 16px;
  height: 8px;
  background: #f5f5f5;
}

.divider-circle.top {
  border-radius: 0 0 8px 8px;
}

.divider-circle.bottom {
  border-radius: 8px 8px 0 0;
}

.coupon-right {
  flex: 1;
  padding: 16px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 8px;
}

.coupon-name {
  font-size: 14px;
  font-weight: 600;
  color: #262626;
}

.coupon-desc {
  font-size: 12px;
  color: #8c8c8c;
}

.coupon-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.coupon-expire {
  font-size: 12px;
  color: #bfbfbf;
}

.coupon-actions {
  display: flex;
  gap: 8px;
  margin-top: 4px;
}

/* ==================== Claim Modal ==================== */
.claim-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.available-claim-list {
  border-top: 1px solid #f0f0f0;
  padding-top: 16px;
}

.claim-section-title {
  font-size: 14px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 12px;
}

.claim-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.claim-item:hover {
  background: #f0f5ff;
}

.claim-amount {
  font-size: 24px;
  font-weight: 700;
  color: #1890ff;
  min-width: 80px;
  text-align: center;
}

.claim-info {
  flex: 1;
}

.claim-name {
  font-size: 14px;
  font-weight: 500;
  color: #262626;
}

.claim-condition {
  font-size: 12px;
  color: #8c8c8c;
}

/* ==================== Responsive ==================== */
@media (max-width: 1200px) {
  .stats-row {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .stats-row {
    grid-template-columns: 1fr;
  }
}
</style>
