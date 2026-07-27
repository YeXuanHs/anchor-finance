<template>
  <div class="invoices-page">
    <div class="page-header">
      <h1 class="page-title">我的账单</h1>
      <n-input
        v-model:value="searchKey"
        placeholder="搜索账单号或描述"
        clearable
        class="search-input"
      >
        <template #prefix>
          <n-icon :component="SearchOutline" color="#bfbfbf" />
        </template>
      </n-input>
    </div>

    <n-tabs v-model:value="activeTab" type="line" animated class="filter-tabs">
      <n-tab v-for="tab in statusTabs" :key="tab.value" :name="tab.value">
        {{ tab.label }}
        <n-badge
          v-if="tab.count > 0"
          :value="tab.count"
          :max="99"
          class="tab-badge"
        />
      </n-tab>
    </n-tabs>

    <div v-if="selectedIds.length > 0" class="batch-bar">
      <n-checkbox v-model:checked="selectAll" @update:checked="handleSelectAll">
        全选
      </n-checkbox>
      <span class="selected-info">
        已选择 <strong>{{ selectedIds.length }}</strong> 项，
        合计 <strong class="total-amount">¥{{ selectedTotal }}</strong>
      </span>
      <n-button
        type="primary"
        size="small"
        @click="handleBatchPay"
      >
        批量支付
      </n-button>
    </div>

    <div class="invoice-list">
      <div
        v-for="invoice in filteredInvoices"
        :key="invoice.id"
        class="invoice-card"
        :class="{ selected: selectedIds.includes(invoice.id) }"
      >
        <div class="invoice-checkbox" v-if="invoice.status === 'unpaid'">
          <n-checkbox
            :checked="selectedIds.includes(invoice.id)"
            @update:checked="(checked: boolean) => handleSelect(invoice.id, checked)"
          />
        </div>

        <div class="invoice-icon">
          <n-icon :size="28" :component="WalletOutline" :color="getStatusColor(invoice.status)" />
        </div>

        <div class="invoice-info">
          <div class="invoice-top">
            <span class="invoice-id">{{ invoice.id }}</span>
            <n-tag :type="getStatusType(invoice.status)" size="small" round>
              {{ invoice.statusText }}
            </n-tag>
          </div>
          <div class="invoice-desc">{{ invoice.description }}</div>
          <div class="invoice-bottom">
            <span class="invoice-due">
              <n-icon :component="CalendarOutline" size="14" />
              到期日：{{ invoice.dueDate }}
            </span>
          </div>
        </div>

        <div class="invoice-amount-section">
          <span class="invoice-amount">¥{{ invoice.amount }}</span>
          <div class="invoice-actions">
            <n-button
              v-if="invoice.status === 'unpaid'"
              type="primary"
              size="small"
              @click="handlePay(invoice)"
            >
              支付
            </n-button>
            <n-button
              size="small"
              quaternary
              @click="handleDetail(invoice)"
            >
              查看详情
            </n-button>
          </div>
        </div>
      </div>

      <div v-if="filteredInvoices.length === 0" class="empty-state">
        <n-icon :size="64" :component="WalletOutline" color="#d9d9d9" />
        <p>暂无账单记录</p>
      </div>
    </div>

    <div v-if="totalPages > 1" class="pagination-wrapper">
      <n-pagination
        v-model:page="currentPage"
        :page-count="totalPages"
        :page-slot="7"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import {
  SearchOutline,
  WalletOutline,
  CalendarOutline
} from '@vicons/ionicons5'

const message = useMessage()
const searchKey = ref('')
const activeTab = ref('all')
const currentPage = ref(1)
const selectedIds = ref<string[]>([])

const statusTabs = computed(() => [
  { label: '全部', value: 'all', count: invoices.value.length },
  { label: '待支付', value: 'unpaid', count: invoices.value.filter(i => i.status === 'unpaid').length },
  { label: '已支付', value: 'paid', count: invoices.value.filter(i => i.status === 'paid').length },
  { label: '已取消', value: 'cancelled', count: invoices.value.filter(i => i.status === 'cancelled').length }
])

interface Invoice {
  id: string
  orderId: string
  description: string
  amount: string
  status: string
  statusText: string
  dueDate: string
  paidAt: string | null
}

const invoices = ref<Invoice[]>([
  { id: 'INV20260725001', orderId: 'ORD20260725001', description: '香港云服务器-月度', amount: '49.00', status: 'paid', statusText: '已支付', dueDate: '2026-08-25', paidAt: '2026-07-25' },
  { id: 'INV20260724002', orderId: 'ORD20260724002', description: '美国独立服务器-季度', amount: '2,397.00', status: 'unpaid', statusText: '待支付', dueDate: '2026-07-31', paidAt: null },
  { id: 'INV20260720003', orderId: 'ORD20260720003', description: 'OV SSL证书-年度', amount: '199.00', status: 'paid', statusText: '已支付', dueDate: '2027-07-20', paidAt: '2026-07-20' },
  { id: 'INV20260715004', orderId: 'ORD20260715004', description: '香港 VPS-月度', amount: '19.00', status: 'unpaid', statusText: '待支付', dueDate: '2026-07-20', paidAt: null },
  { id: 'INV20260710005', orderId: 'ORD20260710005', description: '域名注册-首年', amount: '9.00', status: 'paid', statusText: '已支付', dueDate: '2027-07-10', paidAt: '2026-07-10' },
  { id: 'INV20260705006', orderId: 'ORD20260705006', description: '新加坡 VPS-月度', amount: '35.00', status: 'cancelled', statusText: '已取消', dueDate: '2026-07-10', paidAt: null }
])

const unpaidInvoices = computed(() =>
  invoices.value.filter(i => i.status === 'unpaid')
)

const selectAll = computed({
  get: () => {
    const unpaidIds = unpaidInvoices.value.map(i => i.id)
    return unpaidIds.length > 0 && unpaidIds.every(id => selectedIds.value.includes(id))
  },
  set: () => {}
})

const selectedTotal = computed(() => {
  return invoices.value
    .filter(i => selectedIds.value.includes(i.id))
    .reduce((sum, i) => sum + parseFloat(i.amount.replace(',', '')), 0)
    .toFixed(2)
})

const filteredInvoices = computed(() => {
  let result = invoices.value
  if (activeTab.value !== 'all') {
    result = result.filter(i => i.status === activeTab.value)
  }
  if (searchKey.value) {
    const key = searchKey.value.toLowerCase()
    result = result.filter(i =>
      i.id.toLowerCase().includes(key) ||
      i.description.toLowerCase().includes(key)
    )
  }
  return result
})

const totalPages = computed(() => Math.ceil(filteredInvoices.value.length / 10))

function getStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'error' | 'default'> = {
    paid: 'success',
    unpaid: 'warning',
    cancelled: 'default'
  }
  return map[status] || 'default'
}

function getStatusColor(status: string) {
  const map: Record<string, string> = {
    paid: '#52c41a',
    unpaid: '#faad14',
    cancelled: '#bfbfbf'
  }
  return map[status] || '#8c8c8c'
}

function handleSelect(id: string, checked: boolean) {
  if (checked) {
    selectedIds.value.push(id)
  } else {
    selectedIds.value = selectedIds.value.filter(i => i !== id)
  }
}

function handleSelectAll(checked: boolean) {
  if (checked) {
    selectedIds.value = unpaidInvoices.value.map(i => i.id)
  } else {
    selectedIds.value = []
  }
}

function handlePay(invoice: Invoice) {
  message.info(`正在跳转支付页面：${invoice.id}`)
}

function handleDetail(invoice: Invoice) {
  message.info(`查看账单详情：${invoice.id}`)
}

function handleBatchPay() {
  message.info(`批量支付 ${selectedIds.value.length} 项账单，合计 ¥${selectedTotal.value}`)
}
</script>

<style scoped>
.invoices-page {
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
  font-size: 22px;
  font-weight: 700;
  color: #262626;
  margin: 0;
}

.search-input {
  width: 280px;
}

.filter-tabs {
  background: #fff;
  border-radius: 12px;
  padding: 4px 16px;
  border: 1px solid #f0f0f0;
}

.tab-badge {
  margin-left: 6px;
}

.batch-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 20px;
  background: #f0f5ff;
  border-radius: 12px;
  border: 1px solid #1890ff;
}

.selected-info {
  font-size: 13px;
  color: #595959;
}

.selected-info strong {
  color: #262626;
}

.total-amount {
  color: #ff4d4f;
  font-size: 16px;
}

.invoice-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.invoice-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: #fff;
  border-radius: 12px;
  border: 1px solid #f0f0f0;
  transition: all 0.3s ease;
}

.invoice-card:hover {
  border-color: #d6e8fa;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.08);
}

.invoice-card.selected {
  border-color: #1890ff;
  background: #f0f5ff;
}

.invoice-checkbox {
  flex-shrink: 0;
}

.invoice-icon {
  width: 52px;
  height: 52px;
  background: #f0f5ff;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.invoice-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.invoice-top {
  display: flex;
  align-items: center;
  gap: 12px;
}

.invoice-id {
  font-size: 13px;
  font-weight: 500;
  color: #262626;
  font-family: 'Monaco', 'Menlo', monospace;
}

.invoice-desc {
  font-size: 15px;
  font-weight: 600;
  color: #262626;
}

.invoice-bottom {
  display: flex;
  align-items: center;
  gap: 20px;
  font-size: 12px;
  color: #8c8c8c;
}

.invoice-due {
  display: flex;
  align-items: center;
  gap: 4px;
}

.invoice-amount-section {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
  flex-shrink: 0;
}

.invoice-amount {
  font-size: 22px;
  font-weight: 700;
  color: #ff4d4f;
}

.invoice-actions {
  display: flex;
  gap: 8px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 80px 0;
  background: #fff;
  border-radius: 12px;
  border: 1px solid #f0f0f0;
}

.empty-state p {
  margin: 0;
  color: #8c8c8c;
  font-size: 14px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }

  .search-input {
    width: 100%;
  }

  .invoice-card {
    flex-wrap: wrap;
  }

  .invoice-checkbox {
    order: -1;
  }

  .invoice-amount-section {
    width: 100%;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    padding-top: 12px;
    border-top: 1px solid #f0f0f0;
  }

  .batch-bar {
    flex-wrap: wrap;
  }
}
</style>
