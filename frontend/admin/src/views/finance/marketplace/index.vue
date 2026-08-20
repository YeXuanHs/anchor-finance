<template>
  <div class="marketplace-config">
    <el-card class="config-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.marketplace.title') }}</span>
          <el-button type="primary" :loading="saving" @click="saveConfig">
            {{ $t('finance.marketplace.saveConfig') }}
          </el-button>
        </div>
      </template>

      <el-form :model="config" label-width="160px">
        <el-form-item :label="$t('finance.marketplace.enableMarketplace')">
          <el-switch v-model="config.enabled" />
          <span class="form-tip">{{ $t('finance.marketplace.enableMarketplaceTip') }}</span>
        </el-form-item>

        <el-divider>{{ $t('finance.marketplace.tradeSettings') }}</el-divider>

        <el-form-item :label="$t('finance.marketplace.feeRate')">
          <el-input-number
            v-model="config.fee_rate"
            :min="0"
            :max="100"
            :precision="2"
            :step="0.5"
          />
          <span class="form-tip">%</span>
        </el-form-item>

        <el-form-item :label="$t('finance.marketplace.minFee')">
          <el-input-number
            v-model="config.min_fee"
            :min="0"
            :precision="2"
            :step="1"
          />
          <span class="form-tip">{{ $t('finance.marketplace.yuan') }}</span>
        </el-form-item>

        <el-form-item :label="$t('finance.marketplace.allowFeeOnly')">
          <el-switch v-model="config.allow_fee_only" />
          <span class="form-tip">{{ $t('finance.marketplace.allowFeeOnlyTip') }}</span>
        </el-form-item>

        <el-divider>{{ $t('finance.marketplace.listingLimits') }}</el-divider>

        <el-form-item :label="$t('finance.marketplace.maxListingDays')">
          <el-input-number
            v-model="config.max_listing_days"
            :min="1"
            :max="365"
          />
          <span class="form-tip">{{ $t('finance.marketplace.maxListingDaysTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('finance.marketplace.minHoldDays')">
          <el-input-number
            v-model="config.min_hold_days"
            :min="0"
            :max="365"
          />
          <span class="form-tip">{{ $t('finance.marketplace.minHoldDaysTip') }}</span>
        </el-form-item>

        <el-divider>{{ $t('finance.marketplace.buyerRequirements') }}</el-divider>

        <el-form-item :label="$t('finance.marketplace.requireRealName')">
          <el-switch v-model="config.require_real_name" />
          <span class="form-tip">{{ $t('finance.marketplace.requireRealNameTip') }}</span>
        </el-form-item>

        <el-divider>{{ $t('finance.marketplace.autoOperations') }}</el-divider>

        <el-form-item :label="$t('finance.marketplace.autoTransfer')">
          <el-switch v-model="config.auto_transfer" />
          <span class="form-tip">{{ $t('finance.marketplace.autoTransferTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('finance.marketplace.emailNotify')">
          <el-switch v-model="config.notify_email" />
          <span class="form-tip">{{ $t('finance.marketplace.emailNotifyTip') }}</span>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 挂售列表 -->
    <el-card class="listings-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.marketplace.listingManagement') }}</span>
          <el-radio-group v-model="listingFilter" @change="fetchListings">
            <el-radio-button label="">{{ $t('finance.marketplace.all') }}</el-radio-button>
            <el-radio-button label="1">{{ $t('finance.marketplace.onSale') }}</el-radio-button>
            <el-radio-button label="2">{{ $t('finance.marketplace.sold') }}</el-radio-button>
            <el-radio-button label="3">{{ $t('finance.marketplace.delisted') }}</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <el-table :data="listings" v-loading="loadingListings" style="width: 100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="product_name" :label="$t('finance.marketplace.host')" min-width="150" />
        <el-table-column :label="$t('finance.marketplace.seller')" width="100">
          <template #default="{ row }">
            {{ row.user?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.marketplace.price')" width="100">
          <template #default="{ row }">
            <span class="price">¥{{ row.sell_price?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.marketplace.status')" width="80">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="view_count" :label="$t('finance.marketplace.views')" width="70" />
        <el-table-column :label="$t('finance.marketplace.listingTime')" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper" v-if="listingTotal > 20">
        <el-pagination
          v-model:current-page="listingPage"
          :page-size="20"
          :total="listingTotal"
          layout="total, prev, pager, next"
          @current-change="fetchListings"
        />
      </div>
    </el-card>

    <!-- 订单列表 -->
    <el-card class="orders-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.marketplace.orderManagement') }}</span>
          <el-radio-group v-model="orderFilter" @change="fetchOrders">
            <el-radio-button label="">{{ $t('finance.marketplace.all') }}</el-radio-button>
            <el-radio-button label="0">{{ $t('finance.marketplace.pendingPayment') }}</el-radio-button>
            <el-radio-button label="1">{{ $t('finance.marketplace.paid') }}</el-radio-button>
            <el-radio-button label="2">{{ $t('finance.marketplace.transferred') }}</el-radio-button>
            <el-radio-button label="3">{{ $t('finance.marketplace.completed') }}</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <el-table :data="orders" v-loading="loadingOrders" style="width: 100%">
        <el-table-column prop="order_no" :label="$t('finance.marketplace.orderNo')" width="180" />
        <el-table-column :label="$t('finance.marketplace.buyer')" width="100">
          <template #default="{ row }">
            {{ row.buyer?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.marketplace.seller')" width="100">
          <template #default="{ row }">
            {{ row.seller?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.marketplace.amount')" width="100">
          <template #default="{ row }">
            <span class="price">¥{{ row.total_amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.marketplace.fee')" width="80">
          <template #default="{ row }">
            ¥{{ row.fee?.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.marketplace.method')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.payment_method === 'full' ? 'success' : 'warning'" size="small">
              {{ row.payment_method === 'full' ? $t('finance.marketplace.fullPayment') : $t('finance.marketplace.feeOnly') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.marketplace.status')" width="80">
          <template #default="{ row }">
            <el-tag :type="getOrderStatusType(row.status)" size="small">
              {{ getOrderStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.marketplace.transferStatus')" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.transfer_status > 0" :type="getTransferStatusType(row.transfer_status)" size="small">
              {{ getTransferStatusText(row.transfer_status) }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.marketplace.orderTime')" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper" v-if="orderTotal > 20">
        <el-pagination
          v-model:current-page="orderPage"
          :page-size="20"
          :total="orderTotal"
          layout="total, prev, pager, next"
          @current-change="fetchOrders"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const config = ref({
  enabled: true,
  fee_rate: 5,
  min_fee: 1,
  max_listing_days: 30,
  min_hold_days: 7,
  require_real_name: false,
  allow_fee_only: true,
  auto_transfer: true,
  notify_email: true
})
const saving = ref(false)

const listings = ref<any[]>([])
const listingTotal = ref(0)
const listingPage = ref(1)
const listingFilter = ref('')
const loadingListings = ref(false)

const orders = ref<any[]>([])
const orderTotal = ref(0)
const orderPage = ref(1)
const orderFilter = ref('')
const loadingOrders = ref(false)

onMounted(() => {
  fetchConfig()
  fetchListings()
  fetchOrders()
})

async function fetchConfig() {
  try {
    const res = await request.get({ url: '/api/admin/marketplace/config' })
    if (res) {
      config.value = res
    }
  } catch (e) {
    console.error(e)
  }
}

async function saveConfig() {
  saving.value = true
  try {
    await request.put({ url: '/api/admin/marketplace/config', params: config.value })
    ElMessage.success($t('finance.marketplace.msgSaveSuccess'))
  } catch (e: any) {
    ElMessage.error(e.message || $t('finance.marketplace.msgSaveFailed'))
  } finally {
    saving.value = false
  }
}

async function fetchListings() {
  loadingListings.value = true
  try {
    const res = await request.get({
      url: '/api/admin/marketplace/listings',
      params: {
        page: listingPage.value,
        page_size: 20,
        status: listingFilter.value
      }
    })
    listings.value = res?.list || []
    listingTotal.value = res?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loadingListings.value = false
  }
}

async function fetchOrders() {
  loadingOrders.value = true
  try {
    const res = await request.get({
      url: '/api/admin/marketplace/orders',
      params: {
        page: orderPage.value,
        page_size: 20,
        status: orderFilter.value
      }
    })
    orders.value = res?.list || []
    orderTotal.value = res?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loadingOrders.value = false
  }
}

function getStatusType(status: number) {
  const map: Record<number, any> = { 1: 'success', 2: 'info', 3: 'warning' }
  return map[status] || 'info'
}

function getStatusText(status: number) {
  const map: Record<number, string> = { 1: $t('finance.marketplace.onSale'), 2: $t('finance.marketplace.sold'), 3: $t('finance.marketplace.delisted') }
  return map[status] || $t('finance.marketplace.unknown')
}

function getOrderStatusType(status: number) {
  const map: Record<number, any> = { 0: 'warning', 1: 'primary', 2: 'success', 3: 'success', 4: 'info', 5: 'danger' }
  return map[status] || 'info'
}

function getOrderStatusText(status: number) {
  const map: Record<number, string> = { 0: $t('finance.marketplace.pendingPayment'), 1: $t('finance.marketplace.paid'), 2: $t('finance.marketplace.transferred'), 3: $t('finance.marketplace.completed'), 4: $t('finance.marketplace.cancelled'), 5: $t('finance.marketplace.refunded') }
  return map[status] || $t('finance.marketplace.unknown')
}

function getTransferStatusType(status: number) {
  const map: Record<number, any> = { 1: 'warning', 2: 'success', 3: 'danger' }
  return map[status] || 'info'
}

function getTransferStatusText(status: number) {
  const map: Record<number, string> = { 1: $t('finance.marketplace.transferring'), 2: $t('finance.marketplace.transferSuccess'), 3: $t('finance.marketplace.transferFailed') }
  return map[status] || $t('finance.marketplace.unknown')
}

function formatDate(date: string): string {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}
</script>

<style scoped lang="scss">
.marketplace-config {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-tip {
  margin-left: 12px;
  font-size: 12px;
  color: #999;
}

.price {
  color: #ff4757;
  font-weight: 600;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.el-divider {
  margin: 24px 0;
}
</style>
