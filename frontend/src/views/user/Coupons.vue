<template>
  <div class="coupons-page">
    <div class="page-header">
      <h1 class="page-title">{{ $t('menu.coupons') }}</h1>
      <div class="redeem-area">
        <el-input v-model="redeemCode" :placeholder="$t('coupon.enterCode')" clearable style="width: 200px;" />
        <el-button type="primary" @click="handleRedeem">{{ $t('coupon.redeem') }}</el-button>
      </div>
    </div>

    <el-tabs v-model="activeTab" class="coupon-tabs">
      <el-tab-pane :label="$t('coupon.available')" name="available" />
      <el-tab-pane :label="$t('coupon.used')" name="used" />
      <el-tab-pane :label="$t('coupon.expired')" name="expired" />
    </el-tabs>

    <div class="coupons-grid">
      <div
        v-for="coupon in filteredCoupons"
        :key="coupon.id"
        class="coupon-card"
        :class="{ disabled: activeTab !== 'available' }"
      >
        <div class="coupon-left" :style="{ background: coupon.bg || '#409eff' }">
          <span class="coupon-value">
            <span class="coupon-symbol">¥</span>{{ coupon.value || coupon.amount }}
          </span>
          <span class="coupon-condition">{{ $t('coupon.minSpend', { amount: coupon.threshold || coupon.min_amount || 0 }) }}</span>
        </div>
        <div class="coupon-right">
          <h4 class="coupon-name">{{ coupon.name || coupon.code }}</h4>
          <p class="coupon-desc">{{ coupon.description }}</p>
          <p class="coupon-expire">{{ $t('coupon.expireAt') }} {{ coupon.expireDate || coupon.expire_at }}</p>
          <el-button
            v-if="activeTab === 'available'"
            type="primary"
            size="small"
            @click="handleUse(coupon)"
          >{{ $t('coupon.useNow') }}</el-button>
          <span v-else-if="activeTab === 'used'" class="coupon-used-tag">{{ $t('coupon.used') }}</span>
          <span v-else class="coupon-expired-tag">{{ $t('coupon.expired') }}</span>
        </div>
      </div>
    </div>

    <el-empty v-if="filteredCoupons.length === 0" :description="$t('coupon.empty')" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import request from '@/utils/request'

const { t } = useI18n()
const activeTab = ref('available')
const redeemCode = ref('')
const loading = ref(false)

interface Coupon {
  id: number
  code: string
  name: string
  description: string
  value: number
  amount: number
  threshold: number
  min_amount: number
  expireDate: string
  expire_at: string
  status: string
  bg: string
}

const coupons = ref<Coupon[]>([])

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v2/promo-codes')
    coupons.value = data.data?.list || data.list || data.data || []
  } catch (e) { console.error(e) } finally { loading.value = false }
})

const filteredCoupons = computed(() => {
  return coupons.value.filter(c => c.status === activeTab.value)
})

async function handleRedeem() {
  if (!redeemCode.value) { ElMessage.warning(t('coupon.enterCodeHint')); return }
  try {
    await request.post('/api/v2/promo-codes/validate', { code: redeemCode.value })
    ElMessage.success(t('coupon.redeemSuccess'))
    redeemCode.value = ''
    const { data } = await request.get('/api/v2/promo-codes')
    coupons.value = data.data?.list || data.list || data.data || []
  } catch (e: any) { ElMessage.error(e?.message || t('coupon.redeemFailed')) }
}

function handleUse(coupon: Coupon) { ElMessage.info(t('coupon.useCoupon', { name: coupon.name || coupon.code })) }
</script>

<style scoped>
.coupons-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.redeem-area { display: flex; gap: 8px; }

.coupon-tabs {
  background: #fff;
  border-radius: 12px;
  padding: 4px 16px 0;
  border: 1px solid #e8ecf1;
}

.coupon-tabs :deep(.el-tabs__header) { margin: 0; }

.coupons-grid { display: flex; flex-direction: column; gap: 12px; }

.coupon-card {
  display: flex;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e8ecf1;
  background: #fff;
  transition: all 0.2s;
}

.coupon-card:hover:not(.disabled) {
  box-shadow: 0 2px 12px rgba(0,0,0,0.08);
}

.coupon-card.disabled { opacity: 0.6; }

.coupon-left {
  width: 140px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 20px 16px;
  flex-shrink: 0;
}

.coupon-value { color: #fff; font-size: 32px; font-weight: 700; line-height: 1; }
.coupon-symbol { font-size: 16px; font-weight: 600; }
.coupon-condition { font-size: 12px; color: rgba(255,255,255,0.8); }

.coupon-right {
  flex: 1;
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.coupon-name { font-size: 16px; font-weight: 600; color: #303133; margin: 0; }
.coupon-desc { font-size: 13px; color: #909399; margin: 0; }
.coupon-expire { font-size: 12px; color: #c0c4cc; margin: 0; }
.coupon-used-tag, .coupon-expired-tag { font-size: 13px; color: #c0c4cc; font-weight: 500; }

@media (max-width: 768px) {
  .page-header { flex-direction: column; gap: 12px; align-items: flex-start; }
  .coupon-card { flex-direction: column; }
  .coupon-left { width: auto; padding: 16px; }
}
</style>
