<template>
  <div class="invoice-detail-page" v-loading="loading">
    <div class="page-header">
      <el-button text @click="router.back()">
        <el-icon><ArrowLeft /></el-icon> {{ $t('common.back') }}
      </el-button>
      <h1 class="page-title">{{ $t('invoice.detail') }} #{{ invoice.id }}</h1>
      <div class="header-right">
        <el-tag :type="statusType" size="large" effect="light" round>{{ invoice.statusText }}</el-tag>
      </div>
    </div>

    <div class="detail-grid">
      <el-card shadow="never" class="info-card">
        <template #header><span>{{ $t('invoice.billInfo') }}</span></template>
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('invoice.billNo')">{{ invoice.id }}</el-descriptions-item>
          <el-descriptions-item :label="$t('common.status')">
            <el-tag :type="statusType" size="small">{{ invoice.statusText }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('common.description')">{{ invoice.description }}</el-descriptions-item>
          <el-descriptions-item :label="$t('service.product')">{{ invoice.product }}</el-descriptions-item>
          <el-descriptions-item :label="$t('invoice.amount')">
            <span class="amount-highlight">¥{{ invoice.amount }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('invoice.overdueTime')">{{ invoice.dueDate }}</el-descriptions-item>
          <el-descriptions-item :label="$t('invoice.createTime')">{{ invoice.createTime }}</el-descriptions-item>
          <el-descriptions-item :label="$t('invoice.payTime')">{{ invoice.payTime || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <el-card shadow="never" class="info-card">
        <template #header><span>{{ $t('invoice.paymentInfo') }}</span></template>
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('invoice.payMethod')">{{ invoice.payMethod || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('invoice.transactionNo')">{{ invoice.transactionNo || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('invoice.discount')">{{ invoice.discount ? `¥${invoice.discount}` : '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('invoice.actualAmount')">
            <span class="amount-highlight">¥{{ invoice.actualAmount || invoice.amount }}</span>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
    </div>

    <el-card v-if="invoice.items && invoice.items.length" shadow="never" class="items-card">
      <template #header><span>{{ $t('invoice.items') }}</span></template>
      <el-table :data="invoice.items" stripe>
        <el-table-column prop="name" :label="$t('common.name')" min-width="200" />
        <el-table-column prop="description" :label="$t('common.description')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="quantity" :label="$t('invoice.quantity')" width="100" align="center" />
        <el-table-column prop="unitPrice" :label="$t('invoice.unitPrice')" width="120" align="right">
          <template #default="{ row }">¥{{ row.unitPrice }}</template>
        </el-table-column>
        <el-table-column prop="amount" :label="$t('invoice.amount')" width="120" align="right">
          <template #default="{ row }"><span class="amount-text">¥{{ row.amount }}</span></template>
        </el-table-column>
      </el-table>
    </el-card>

    <div v-if="invoice.status !== 'paid'" class="action-bar">
      <el-button type="primary" size="large" @click="handlePay">{{ $t('invoice.pay') }}</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ArrowLeft } from '@element-plus/icons-vue'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()
const loading = ref(false)

interface InvoiceItem {
  name: string
  description: string
  quantity: number
  unitPrice: string
  amount: string
}

interface InvoiceDetail {
  id: string
  description: string
  product: string
  amount: string
  dueDate: string
  status: string
  statusText: string
  createTime: string
  payTime: string
  payMethod: string
  transactionNo: string
  discount: string
  actualAmount: string
  items: InvoiceItem[]
}

const invoice = ref<InvoiceDetail>({
  id: '', description: '', product: '', amount: '0.00', dueDate: '',
  status: '', statusText: '', createTime: '', payTime: '',
  payMethod: '', transactionNo: '', discount: '', actualAmount: '', items: []
})

const statusType = computed(() => {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    paid: 'success', unpaid: 'warning', overdue: 'danger'
  }
  return map[invoice.value.status] || 'info'
})

onMounted(async () => {
  const id = route.params.id as string
  if (!id) return
  loading.value = true
  try {
    const res = await request.get(`/api/v1/invoices/${id}`)
    const data = res?.data || res
    if (data) Object.assign(invoice.value, data)
  } catch (e) { console.error(e) } finally { loading.value = false }
})

function handlePay() {
  router.push(`/user/invoices/${invoice.value.id}/pay`)
}
</script>

<style scoped>
.invoice-detail-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
  flex: 1;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.info-card, .items-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.amount-highlight {
  font-size: 18px;
  font-weight: 700;
  color: #303133;
}

.amount-text {
  font-weight: 600;
  color: #303133;
}

.action-bar {
  display: flex;
  justify-content: flex-end;
  padding: 20px 0;
}

@media (max-width: 768px) {
  .detail-grid { grid-template-columns: 1fr; }
  .page-header { flex-wrap: wrap; }
}
</style>
