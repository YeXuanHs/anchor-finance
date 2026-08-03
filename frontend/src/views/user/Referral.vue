<template>
  <div class="referral-page">
    <div class="page-header">
      <h1 class="page-title">{{ $t('menu.referral') }}</h1>
    </div>

    <!-- Banner Card -->
    <el-card shadow="never" class="banner-card">
      <div class="banner-content">
        <div class="banner-text">
          <h2>{{ $t('referral.registerThroughPromotion') }}</h2>
          <p>{{ $t('referral.successfulPromotion') }}</p>
        </div>
        <div class="banner-stats">
          <div class="banner-stat">
            <span class="stat-num">{{ totalReferrals }}</span>
            <span class="stat-label">{{ $t('affiliate.registeredUsers') }}</span>
          </div>
          <div class="banner-stat">
            <span class="stat-num">¥{{ totalEarnings }}</span>
            <span class="stat-label">{{ $t('affiliate.commission') }}</span>
          </div>
          <div class="banner-stat">
            <span class="stat-num">¥{{ availableBalance }}</span>
            <span class="stat-label">{{ $t('affiliate.availableAmount') }}</span>
          </div>
        </div>
      </div>
    </el-card>

    <!-- Quick Navigation -->
    <div class="quick-nav">
      <el-card shadow="never" class="nav-card" @click="$router.push('/user/affiliate/buy-record')">
        <div class="nav-content">
          <div class="nav-icon" style="background: linear-gradient(135deg, #0056FF, #4080FF);">
            <el-icon :size="28"><ShoppingCart /></el-icon>
          </div>
          <div class="nav-info">
            <span class="nav-title">{{ $t('affiliate.orderTime') }}</span>
            <span class="nav-desc">{{ $t('referral.promotionRecord') }}</span>
          </div>
          <el-icon class="nav-arrow"><ArrowRight /></el-icon>
        </div>
      </el-card>
      <el-card shadow="never" class="nav-card" @click="$router.push('/user/affiliate/user-list')">
        <div class="nav-content">
          <div class="nav-icon" style="background: linear-gradient(135deg, #52c41a, #73d13d);">
            <el-icon :size="28"><User /></el-icon>
          </div>
          <div class="nav-info">
            <span class="nav-title">{{ $t('affiliate.registeredUsers') }}</span>
            <span class="nav-desc">{{ $t('affiliate.sharePromotionTip') }}</span>
          </div>
          <el-icon class="nav-arrow"><ArrowRight /></el-icon>
        </div>
      </el-card>
      <el-card shadow="never" class="nav-card" @click="$router.push('/user/affiliate/withdraw')">
        <div class="nav-content">
          <div class="nav-icon" style="background: linear-gradient(135deg, #fa8c16, #ffc53d);">
            <el-icon :size="28"><Wallet /></el-icon>
          </div>
          <div class="nav-info">
            <span class="nav-title">{{ $t('affiliate.withdrawalAmount') }}</span>
            <span class="nav-desc">{{ $t('affiliate.immediateWithdrawal') }}</span>
          </div>
          <el-icon class="nav-arrow"><ArrowRight /></el-icon>
        </div>
      </el-card>
    </div>

    <!-- Referral Code Card -->
    <el-card shadow="never" class="code-card">
      <template #header>
        <span class="card-title">{{ $t('affiliate.offerLink') }}</span>
      </template>
      <div class="code-area">
        <el-input v-model="referralCode" readonly size="large" class="code-input">
          <template #append>
            <el-button @click="handleCopy">
              <el-icon><CopyDocument /></el-icon>{{ $t('common.copy') }}
            </el-button>
          </template>
        </el-input>
        <p class="code-hint">{{ $t('referral.registerThroughPromotion') }}</p>
      </div>
      <div class="share-links">
        <span class="share-label">{{ $t('affiliate.recommendedLinks') }}：</span>
        <el-input v-model="shareLink" readonly size="small" style="flex: 1;">
          <template #append>
            <el-button @click="handleCopyLink">{{ $t('affiliate.sharePromotionLinks') }}</el-button>
          </template>
        </el-input>
      </div>
    </el-card>

    <!-- Recent Records -->
    <el-card shadow="never" class="records-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">{{ $t('referral.promotionRecord') }}</span>
          <el-button type="primary" link @click="$router.push('/user/affiliate/buy-record')">{{ $t('home.viewMore') }}</el-button>
        </div>
      </template>
      <el-table :data="referralRecords" stripe style="width: 100%" :empty-text="$t('common.noData')">
        <el-table-column prop="username" :label="$t('tableSearch.clientUsername')" min-width="120">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="28" class="record-avatar">{{ row.username.charAt(0) }}</el-avatar>
              <span>{{ row.username }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="registerDate" :label="$t('transaction.registrationTime')" width="130" />
        <el-table-column prop="spent" :label="$t('affiliate.orderAmount')" width="120">
          <template #default="{ row }">
            <span class="amount-text">¥{{ row.spent }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="commission" :label="$t('affiliate.commission')" width="120">
          <template #default="{ row }">
            <span class="commission-text">¥{{ row.commission }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('common.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'settled' ? 'success' : 'warning'" size="small" effect="light" round>
              {{ row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Rules Card -->
    <el-card shadow="never" class="rules-card">
      <template #header>
        <span class="card-title">{{ $t('coupon.redeem') }}</span>
      </template>
      <el-timeline>
        <el-timeline-item v-for="(rule, i) in rules" :key="i" :type="i === 0 ? 'primary' : ''">
          {{ rule }}
        </el-timeline-item>
      </el-timeline>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { CopyDocument, ShoppingCart, User, Wallet, ArrowRight } from '@element-plus/icons-vue'
import request from '@/utils/request'

const { t } = useI18n()

const loading = ref(false)
const referralCode = ref('')
const shareLink = ref('')
const totalReferrals = ref(0)
const totalEarnings = ref('0.00')
const availableBalance = ref('0.00')

const referralRecords = ref<any[]>([])

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/affiliate/info')
    const info = data.data || {}
    referralCode.value = info.referralCode || ''
    shareLink.value = info.shareLink || ''
    totalReferrals.value = info.totalReferrals || 0
    totalEarnings.value = info.totalEarnings || '0.00'
    availableBalance.value = info.availableBalance || '0.00'
    referralRecords.value = info.records || []
  } catch (e) { console.error(e) } finally { loading.value = false }
})

const rules = [
  t('referral.registerThroughPromotion'),
  t('referral.successfulPromotion'),
  t('affiliate.getRebateTip'),
  t('affiliate.sharePromotionTip'),
  t('affiliate.registerThroughPromotion')
]

function handleCopy() {
  navigator.clipboard?.writeText(referralCode.value)
  ElMessage.success(t('affiliate.copySuccess'))
}

function handleCopyLink() {
  navigator.clipboard?.writeText(shareLink.value)
  ElMessage.success(t('affiliate.copySuccess'))
}
</script>

<style scoped>
.referral-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.banner-card {
  border-radius: 12px;
  border: none;
  background: linear-gradient(135deg, #1a3a5c 0%, #0056FF 100%);
}

.banner-card :deep(.el-card__body) { padding: 32px; }

.banner-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.banner-text h2 { font-size: 22px; font-weight: 700; color: #fff; margin: 0 0 8px 0; }
.banner-text p { font-size: 14px; color: rgba(255,255,255,0.75); margin: 0; }
.banner-text strong { color: #ffd54f; }

.banner-stats { display: flex; gap: 32px; }

.banner-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-num { font-size: 28px; font-weight: 700; color: #fff; }
.stat-label { font-size: 13px; color: rgba(255,255,255,0.65); }

/* Quick Navigation */
.quick-nav {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.nav-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  cursor: pointer;
  transition: all 0.2s;
}

.nav-card:hover {
  box-shadow: 0 4px 12px rgba(0, 86, 255, 0.12);
  transform: translateY(-2px);
}

.nav-card :deep(.el-card__body) {
  padding: 20px;
}

.nav-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.nav-icon {
  width: 52px;
  height: 52px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.nav-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.nav-desc {
  font-size: 13px;
  color: #909399;
}

.nav-arrow {
  color: #c0c4cc;
  font-size: 18px;
}

.code-card, .records-card, .rules-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  background: #fff;
}

.card-title { font-size: 15px; font-weight: 600; color: #303133; }
.card-header { display: flex; align-items: center; justify-content: space-between; }

.code-area { display: flex; flex-direction: column; gap: 12px; }

.code-input :deep(.el-input__inner) {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 18px;
  font-weight: 600;
  color: #0056FF;
  letter-spacing: 2px;
}

.code-hint { font-size: 13px; color: #909399; margin: 0; }

.share-links {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e8ecf1;
}

.share-label { font-size: 14px; color: #606266; white-space: nowrap; }

.user-cell { display: flex; align-items: center; gap: 8px; }

.record-avatar {
  background: linear-gradient(135deg, #0056FF, #4080FF);
  color: #fff;
  font-size: 12px;
  font-weight: 600;
}

.amount-text { font-weight: 600; color: #303133; }
.commission-text { font-weight: 600; color: #fa8c16; }

.records-card :deep(.el-table th.el-table__cell) { background: #fafbfc; font-weight: 600; }
.rules-card :deep(.el-timeline) { padding-left: 0; }

@media (max-width: 768px) {
  .banner-content { flex-direction: column; gap: 20px; text-align: center; }
  .banner-stats { justify-content: center; }
  .quick-nav { grid-template-columns: 1fr; }
}
</style>
