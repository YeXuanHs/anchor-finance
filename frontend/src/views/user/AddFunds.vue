<template>
  <div class="add-funds-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('addFunds.title') }}</span>
          <span class="balance">{{ $t('addFunds.currentBalance') }}<em>¥{{ balance.toFixed(2) }}</em></span>
        </div>
      </template>

      <el-form label-position="top">
        <el-form-item :label="$t('addFunds.amount')">
          <div class="amount-presets">
            <el-button
              v-for="preset in presetAmounts"
              :key="preset"
              :type="amount === preset ? 'primary' : 'default'"
              @click="amount = preset"
            >
              ¥{{ preset }}
            </el-button>
          </div>
          <el-input-number
            v-model="amount"
            :min="1"
            :max="50000"
            :precision="2"
            :placeholder="$t('addFunds.customAmount')"
            style="width: 100%; margin-top: 12px"
          />
        </el-form-item>

        <el-form-item :label="$t('addFunds.paymentMethod')">
          <div class="payment-methods" v-loading="loadingGateways">
            <div
              v-for="method in paymentMethods"
              :key="method.id"
              class="payment-method"
              :class="{ active: paymentMethod === method.name }"
              @click="paymentMethod = method.name"
            >
              <img :src="method.icon || getIconUrl(method.code)" :alt="method.title" class="method-icon" />
              <span class="method-label">{{ method.title }}</span>
              <el-icon v-if="paymentMethod === method.name" class="check-icon"><Check /></el-icon>
            </div>
            <el-empty v-if="!loadingGateways && paymentMethods.length === 0" :description="$t('addFunds.noPaymentMethod')" />
          </div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" size="large" :loading="submitting" :disabled="!paymentMethod" @click="handleRecharge">
            {{ $t('addFunds.confirmRecharge') }} ¥{{ amount?.toFixed(2) || '0.00' }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('addFunds.rechargeRecord') }}</span>
        </div>
      </template>

      <el-table :data="records" style="width: 100%">
        <el-table-column prop="created_at" :label="$t('addFunds.time')" width="180" />
        <el-table-column prop="amount" :label="$t('addFunds.amountLabel')">
          <template #default="{ row }">
            <span class="text-success">+¥{{ row.amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="payment_method" :label="$t('addFunds.paymentMethodLabel')" />
        <el-table-column prop="status" :label="$t('addFunds.status')">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="transaction_no" :label="$t('addFunds.transactionNo')" />
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadRecords"
          @current-change="loadRecords"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Check } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import request from '@/utils/request'

const { t } = useI18n()

const balance = ref(0)
const amount = ref<number | undefined>(100)
const paymentMethod = ref('')
const submitting = ref(false)
const loadingGateways = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const presetAmounts = [50, 100, 200, 500, 1000, 2000]

interface PaymentMethod {
  id: number
  name: string
  title: string
  code: string
  icon: string
}

const paymentMethods = ref<PaymentMethod[]>([])

const getIconUrl = (code: string) => {
  const icons: Record<string, string> = {
    alipay: '/assets/payment/alipay.png',
    wechat: '/assets/payment/wechat.png',
    qqpay: '/assets/payment/qqpay.png',
    usdt: '/assets/payment/usdt.png',
    bank: '/assets/payment/bank.png',
    balance: '/assets/payment/balance.png'
  }
  return icons[code] || '/assets/payment/default.png'
}

interface RechargeRecord {
  id: number
  created_at: string
  amount: number
  payment_method: string
  status: string
  transaction_no: string
}

const records = ref<RechargeRecord[]>([])

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    success: 'success',
    pending: 'warning',
    failed: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    success: t('addFunds.statusSuccess'),
    pending: t('addFunds.statusPending'),
    failed: t('addFunds.statusFailed')
  }
  return map[status] || status
}

const loadPaymentMethods = async () => {
  loadingGateways.value = true
  try {
    const { data } = await request.get('/api/v1/payment-methods')
    paymentMethods.value = data?.data || []
    if (paymentMethods.value.length > 0 && !paymentMethod.value) {
      paymentMethod.value = paymentMethods.value[0].name
    }
  } catch {
    paymentMethods.value = []
  } finally {
    loadingGateways.value = false
  }
}

const handleRecharge = async () => {
  if (!amount.value || amount.value <= 0) {
    ElMessage.warning(t('addFunds.enterValidAmount'))
    return
  }
  if (!paymentMethod.value) {
    ElMessage.warning(t('addFunds.selectPaymentMethod'))
    return
  }
  submitting.value = true
  try {
    const { data } = await request.post('/api/v2/balance/recharge', {
      amount: amount.value,
      gateway: paymentMethod.value
    })
    if (data?.data?.pay_url) {
      window.location.href = data.data.pay_url
    } else if (data?.data?.type === 'bank_transfer') {
      ElMessage.info(t('addFunds.bankTransferTip'))
    } else {
      ElMessage.success(t('addFunds.rechargeSubmitted'))
      loadRecords()
    }
  } catch (e: any) {
    ElMessage.error(e.response?.data?.message || t('addFunds.rechargeFailed'))
  } finally {
    submitting.value = false
  }
}

const loadRecords = async () => {
  try {
    const { data } = await request.get('/api/v2/balance/logs', {
      params: { page: currentPage.value, page_size: pageSize.value }
    })
    records.value = data?.data?.list || data?.data?.items || []
    total.value = data?.data?.total || 0
  } catch {
    records.value = []
    total.value = 0
  }
}

onMounted(async () => {
  try {
    const { data } = await request.get('/api/v2/balance')
    balance.value = data?.data?.balance || 0
  } catch {}
  loadPaymentMethods()
  loadRecords()
})
</script>

<style scoped lang="scss">
.add-funds-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-card {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .balance {
      font-size: 14px;
      color: #606266;

      em {
        font-style: normal;
        font-size: 20px;
        font-weight: bold;
        color: #f56c6c;
      }
    }
  }
}

.amount-presets {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.payment-methods {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  min-height: 80px;
}

.payment-method {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 24px;
  border: 2px solid #ebeef5;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;

  &:hover {
    border-color: #c0c4cc;
  }

  &.active {
    border-color: #409eff;
    background: #ecf5ff;
  }

  .method-icon {
    width: 32px;
    height: 32px;
  }

  .method-label {
    font-weight: 500;
  }

  .check-icon {
    color: #409eff;
    margin-left: auto;
  }
}

.text-success {
  color: #67c23a;
  font-weight: bold;
}

.pagination-wrap {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
