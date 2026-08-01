<template>
  <div class="marketplace-config">
    <el-card class="config-card">
      <template #header>
        <div class="card-header">
          <span>交易市场配置</span>
          <el-button type="primary" :loading="saving" @click="saveConfig">
            保存配置
          </el-button>
        </div>
      </template>

      <el-form :model="config" label-width="160px">
        <el-form-item label="启用交易市场">
          <el-switch v-model="config.enabled" />
          <span class="form-tip">开启后用户可以在市场挂售和购买主机</span>
        </el-form-item>

        <el-divider>交易设置</el-divider>

        <el-form-item label="手续费率">
          <el-input-number
            v-model="config.fee_rate"
            :min="0"
            :max="100"
            :precision="2"
            :step="0.5"
          />
          <span class="form-tip">%</span>
        </el-form-item>

        <el-form-item label="最低手续费">
          <el-input-number
            v-model="config.min_fee"
            :min="0"
            :precision="2"
            :step="1"
          />
          <span class="form-tip">元</span>
        </el-form-item>

        <el-form-item label="允许仅付手续费模式">
          <el-switch v-model="config.allow_fee_only" />
          <span class="form-tip">开启后买家可选择仅支付手续费获取卖家联系方式私下交易</span>
        </el-form-item>

        <el-divider>挂售限制</el-divider>

        <el-form-item label="挂售最长天数">
          <el-input-number
            v-model="config.max_listing_days"
            :min="1"
            :max="365"
          />
          <span class="form-tip">天，超过自动下架</span>
        </el-form-item>

        <el-form-item label="最少持有天数">
          <el-input-number
            v-model="config.min_hold_days"
            :min="0"
            :max="365"
          />
          <span class="form-tip">天，购买后需持有此天数才能挂售</span>
        </el-form-item>

        <el-divider>买家要求</el-divider>

        <el-form-item label="买家需要实名认证">
          <el-switch v-model="config.require_real_name" />
          <span class="form-tip">开启后未实名用户无法购买</span>
        </el-form-item>

        <el-divider>自动操作</el-divider>

        <el-form-item label="自动转移">
          <el-switch v-model="config.auto_transfer" />
          <span class="form-tip">付款成功后自动转移主机所有权</span>
        </el-form-item>

        <el-form-item label="邮件通知">
          <el-switch v-model="config.notify_email" />
          <span class="form-tip">交易相关事件发送邮件通知</span>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 挂售列表 -->
    <el-card class="listings-card">
      <template #header>
        <div class="card-header">
          <span>挂售管理</span>
          <el-radio-group v-model="listingFilter" @change="fetchListings">
            <el-radio-button label="">全部</el-radio-button>
            <el-radio-button label="1">在售</el-radio-button>
            <el-radio-button label="2">已售</el-radio-button>
            <el-radio-button label="3">已下架</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <el-table :data="listings" v-loading="loadingListings" style="width: 100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="product_name" label="主机" min-width="150" />
        <el-table-column label="卖家" width="100">
          <template #default="{ row }">
            {{ row.user?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="售价" width="100">
          <template #default="{ row }">
            <span class="price">¥{{ row.sell_price?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="view_count" label="浏览" width="70" />
        <el-table-column label="挂售时间" width="160">
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
          <span>订单管理</span>
          <el-radio-group v-model="orderFilter" @change="fetchOrders">
            <el-radio-button label="">全部</el-radio-button>
            <el-radio-button label="0">待支付</el-radio-button>
            <el-radio-button label="1">已支付</el-radio-button>
            <el-radio-button label="2">已转移</el-radio-button>
            <el-radio-button label="3">已完成</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <el-table :data="orders" v-loading="loadingOrders" style="width: 100%">
        <el-table-column prop="order_no" label="订单号" width="180" />
        <el-table-column label="买家" width="100">
          <template #default="{ row }">
            {{ row.buyer?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="卖家" width="100">
          <template #default="{ row }">
            {{ row.seller?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="金额" width="100">
          <template #default="{ row }">
            <span class="price">¥{{ row.total_amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="手续费" width="80">
          <template #default="{ row }">
            ¥{{ row.fee?.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column label="方式" width="80">
          <template #default="{ row }">
            <el-tag :type="row.payment_method === 'full' ? 'success' : 'warning'" size="small">
              {{ row.payment_method === 'full' ? '全额' : '仅手续费' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="getOrderStatusType(row.status)" size="small">
              {{ getOrderStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="转移状态" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.transfer_status > 0" :type="getTransferStatusType(row.transfer_status)" size="small">
              {{ getTransferStatusText(row.transfer_status) }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="下单时间" width="160">
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
    const res = await request.get('/admin/marketplace/config')
    if (res.data) {
      config.value = res.data
    }
  } catch (e) {
    console.error(e)
  }
}

async function saveConfig() {
  saving.value = true
  try {
    await request.put('/admin/marketplace/config', config.value)
    ElMessage.success('保存成功')
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
  } finally {
    saving.value = false
  }
}

async function fetchListings() {
  loadingListings.value = true
  try {
    const res = await request.get('/admin/marketplace/listings', {
      params: {
        page: listingPage.value,
        page_size: 20,
        status: listingFilter.value
      }
    })
    listings.value = res.data?.list || []
    listingTotal.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loadingListings.value = false
  }
}

async function fetchOrders() {
  loadingOrders.value = true
  try {
    const res = await request.get('/admin/marketplace/orders', {
      params: {
        page: orderPage.value,
        page_size: 20,
        status: orderFilter.value
      }
    })
    orders.value = res.data?.list || []
    orderTotal.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loadingOrders.value = false
  }
}

function getStatusType(status: number) {
  const map: Record<number, string> = { 1: 'success', 2: 'info', 3: 'warning' }
  return map[status] || 'info'
}

function getStatusText(status: number) {
  const map: Record<number, string> = { 1: '在售', 2: '已售', 3: '已下架' }
  return map[status] || '未知'
}

function getOrderStatusType(status: number) {
  const map: Record<number, string> = { 0: 'warning', 1: 'primary', 2: 'success', 3: 'success', 4: 'info', 5: 'danger' }
  return map[status] || 'info'
}

function getOrderStatusText(status: number) {
  const map: Record<number, string> = { 0: '待支付', 1: '已支付', 2: '已转移', 3: '已完成', 4: '已取消', 5: '已退款' }
  return map[status] || '未知'
}

function getTransferStatusType(status: number) {
  const map: Record<number, string> = { 1: 'warning', 2: 'success', 3: 'danger' }
  return map[status] || 'info'
}

function getTransferStatusText(status: number) {
  const map: Record<number, string> = { 1: '转移中', 2: '转移成功', 3: '转移失败' }
  return map[status] || '未知'
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
