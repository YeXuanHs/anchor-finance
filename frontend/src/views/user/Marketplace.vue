<template>
  <div class="marketplace-container">
    <!-- 顶部横幅 -->
    <div class="marketplace-banner">
      <h1>主机交易市场</h1>
      <p>闲置主机快速变现，优质资源低价入手</p>
      <div class="banner-stats">
        <div class="stat-item">
          <span class="stat-num">{{ stats.totalListings }}</span>
          <span class="stat-label">在售主机</span>
        </div>
        <div class="stat-item">
          <span class="stat-num">{{ stats.totalSellers }}</span>
          <span class="stat-label">卖家数量</span>
        </div>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <div class="filter-left">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索主机名称..."
          prefix-icon="Search"
          style="width: 250px"
          @keyup.enter="fetchListings"
        />
      </div>
      <div class="filter-right">
        <el-button type="primary" @click="$router.push('/user/marketplace/sell')">
          <el-icon><Plus /></el-icon>
          我要挂售
        </el-button>
        <el-button @click="$router.push('/user/marketplace/orders')">
          我的订单
        </el-button>
        <el-button @click="$router.push('/user/marketplace/chat')">
          <el-icon><ChatDotRound /></el-icon>
          消息
          <el-badge v-if="unreadCount > 0" :value="unreadCount" class="msg-badge" />
        </el-button>
      </div>
    </div>

    <!-- 商品列表 -->
    <div class="listings-grid" v-loading="loading">
      <div v-if="listings.length === 0 && !loading" class="empty-tip">
        <el-empty description="暂无在售主机" />
      </div>
      <div
        v-for="item in listings"
        :key="item.id"
        class="listing-card"
        @click="viewListing(item)"
      >
        <div class="card-header">
          <span class="product-name">{{ item.product_name }}</span>
          <el-tag size="small" type="success">在售</el-tag>
        </div>

        <div class="card-body">
          <div class="specs-row">
            <div class="spec-item">
              <span class="spec-label">CPU</span>
              <span class="spec-value">{{ getHostSpec(item, 'cpu') }}核</span>
            </div>
            <div class="spec-item">
              <span class="spec-label">内存</span>
              <span class="spec-value">{{ getHostSpec(item, 'memory') }}G</span>
            </div>
            <div class="spec-item">
              <span class="spec-label">硬盘</span>
              <span class="spec-value">{{ getHostSpec(item, 'disk') }}G</span>
            </div>
            <div class="spec-item">
              <span class="spec-label">带宽</span>
              <span class="spec-value">{{ getHostSpec(item, 'bandwidth') }}M</span>
            </div>
          </div>

          <div v-if="item.host?.dedicated_ip" class="ip-info">
            <span>IP: {{ item.host.dedicated_ip }}</span>
          </div>

          <div v-if="item.expires_at" class="expire-info">
            到期: {{ formatDate(item.expires_at) }}
            <span class="expire-days">剩余 {{ getDaysLeft(item.expires_at) }} 天</span>
          </div>
        </div>

        <div class="card-footer">
          <div class="price-info">
            <div class="original-price">
              官网月付 ¥{{ item.original_price?.toFixed(2) }}
            </div>
            <div class="sell-price">
              ¥{{ item.sell_price?.toFixed(2) }}
              <span class="price-unit">/月</span>
            </div>
          </div>
          <div class="seller-info">
            <el-avatar :size="20">{{ item.user?.username?.charAt(0) || '?' }}</el-avatar>
            <span class="seller-name">{{ item.user?.username || '未知' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination-wrapper" v-if="total > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        @current-change="fetchListings"
      />
    </div>

    <!-- 详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      :title="currentListing?.product_name"
      width="600px"
    >
      <div v-if="currentListing" class="listing-detail">
        <div class="detail-section">
          <h3>主机配置</h3>
          <div class="detail-specs">
            <div class="spec-row">
              <span class="label">CPU:</span>
              <span>{{ getHostSpec(currentListing, 'cpu') }}核</span>
            </div>
            <div class="spec-row">
              <span class="label">内存:</span>
              <span>{{ getHostSpec(currentListing, 'memory') }}G</span>
            </div>
            <div class="spec-row">
              <span class="label">硬盘:</span>
              <span>{{ getHostSpec(currentListing, 'disk') }}G</span>
            </div>
            <div class="spec-row">
              <span class="label">带宽:</span>
              <span>{{ getHostSpec(currentListing, 'bandwidth') }}M</span>
            </div>
            <div class="spec-row" v-if="currentListing.host?.os">
              <span class="label">系统:</span>
              <span>{{ currentListing.host.os }}</span>
            </div>
            <div class="spec-row" v-if="currentListing.host?.dedicated_ip">
              <span class="label">IP:</span>
              <span>{{ currentListing.host.dedicated_ip }}</span>
            </div>
          </div>
        </div>

        <div class="detail-section" v-if="currentListing.description">
          <h3>卖家描述</h3>
          <p class="description">{{ currentListing.description }}</p>
        </div>

        <div class="detail-section">
          <h3>卖家信息</h3>
          <div class="seller-detail">
            <el-avatar :size="40">{{ currentListing.user?.username?.charAt(0) || '?' }}</el-avatar>
            <div class="seller-meta">
              <span class="name">{{ currentListing.user?.username }}</span>
              <span class="views">浏览 {{ currentListing.view_count }} 次</span>
            </div>
          </div>
        </div>

        <div class="price-section">
          <div class="price-compare">
            <div class="price-item">
              <span class="label">官网月付</span>
              <span class="value original">¥{{ currentListing.original_price?.toFixed(2) }}</span>
            </div>
            <div class="price-item">
              <span class="label">售价</span>
              <span class="value sell">¥{{ currentListing.sell_price?.toFixed(2) }}</span>
            </div>
            <div class="price-item">
              <span class="label">节省</span>
              <span class="value save">¥{{ ((currentListing.original_price || 0) - (currentListing.sell_price || 0)).toFixed(2) }}</span>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="contactSeller(currentListing)">
            <el-icon><ChatDotRound /></el-icon>
            联系卖家
          </el-button>
          <el-button type="warning" @click="showPurchaseDialog(currentListing, 'fee_only')">
            仅付手续费联系卖家
          </el-button>
          <el-button type="primary" @click="showPurchaseDialog(currentListing, 'full')">
            全额购买
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 购买确认弹窗 -->
    <el-dialog
      v-model="purchaseVisible"
      title="确认购买"
      width="500px"
    >
      <div v-if="purchaseListing" class="purchase-confirm">
        <div class="purchase-info">
          <h3>{{ purchaseListing.product_name }}</h3>
          <div class="purchase-price">
            <span v-if="purchaseType === 'full'">
              售价: ¥{{ purchaseListing.sell_price?.toFixed(2) }}
            </span>
            <span v-else>
              仅需支付手续费即可获取卖家联系方式
            </span>
          </div>
        </div>

        <div class="fee-breakdown">
          <div class="fee-row" v-if="purchaseType === 'full'">
            <span>商品价格</span>
            <span>¥{{ purchaseListing.sell_price?.toFixed(2) }}</span>
          </div>
          <div class="fee-row">
            <span>手续费 ({{ config.feeRate }}%)</span>
            <span>¥{{ calculateFee(purchaseListing).toFixed(2) }}</span>
          </div>
          <div class="fee-row total">
            <span>合计</span>
            <span>¥{{ calculateTotal(purchaseListing, purchaseType).toFixed(2) }}</span>
          </div>
        </div>

        <div class="balance-info">
          当前余额: <span class="balance">¥{{ userBalance.toFixed(2) }}</span>
          <span v-if="userBalance < calculateTotal(purchaseListing, purchaseType)" class="insufficient">
            (余额不足)
          </span>
        </div>

        <el-checkbox v-model="agreedTerms">
          我已阅读并同意《交易市场用户条款》
        </el-checkbox>
      </div>

      <template #footer>
        <el-button @click="purchaseVisible = false">取消</el-button>
        <el-button
          type="primary"
          :disabled="!agreedTerms || userBalance < calculateTotal(purchaseListing, purchaseType)"
          :loading="purchasing"
          @click="confirmPurchase"
        >
          确认购买
        </el-button>
      </template>
    </el-dialog>

    <!-- 卖家联系方式弹窗 -->
    <el-dialog
      v-model="sellerContactVisible"
      title="卖家联系方式"
      width="400px"
    >
      <div class="seller-contact">
        <p>请与卖家私下完成交易。</p>
        <p>交易完成后，请联系卖家在个人中心确认交易完成，服务器将自动转移。</p>
        <div class="contact-info" v-if="sellerContact">
          <div class="contact-item">
            <span class="label">用户名:</span>
            <span>{{ sellerContact.username }}</span>
          </div>
          <div class="contact-item" v-if="sellerContact.email">
            <span class="label">邮箱:</span>
            <span>{{ sellerContact.email }}</span>
          </div>
          <div class="contact-item" v-if="sellerContact.phone">
            <span class="label">手机:</span>
            <span>{{ sellerContact.phone }}</span>
          </div>
          <div class="contact-item" v-if="sellerContact.qq">
            <span class="label">QQ:</span>
            <span>{{ sellerContact.qq }}</span>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus, ChatDotRound } from '@element-plus/icons-vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const listings = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchKeyword = ref('')
const unreadCount = ref(0)
const config = ref<any>({
  feeRate: 5,
  minFee: 1
})

// 详情弹窗
const detailVisible = ref(false)
const currentListing = ref<any>(null)

// 购买弹窗
const purchaseVisible = ref(false)
const purchaseListing = ref<any>(null)
const purchaseType = ref<'full' | 'fee_only'>('full')
const agreedTerms = ref(false)
const purchasing = ref(false)
const userBalance = computed(() => userStore.userInfo?.balance || 0)

// 卖家联系方式
const sellerContactVisible = ref(false)
const sellerContact = ref<any>(null)

const stats = computed(() => ({
  totalListings: total.value,
  totalSellers: new Set(listings.value.map(l => l.user_id)).size
}))

onMounted(() => {
  fetchListings()
  fetchConfig()
  fetchUnreadCount()
})

async function fetchListings() {
  loading.value = true
  try {
    const res = await request.get('/v1/marketplace/listings', {
      params: {
        page: currentPage.value,
        page_size: pageSize.value,
        keyword: searchKeyword.value
      }
    })
    listings.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function fetchConfig() {
  try {
    // 公开配置可以通过接口获取
    config.value = {
      feeRate: 5,
      minFee: 1
    }
  } catch (e) {
    console.error(e)
  }
}

async function fetchUnreadCount() {
  try {
    const res = await request.get('/v1/marketplace/unread-count')
    unreadCount.value = res.data?.unread_count || 0
  } catch (e) {
    console.error(e)
  }
}

function viewListing(item: any) {
  currentListing.value = item
  detailVisible.value = true
}

function contactSeller(listing: any) {
  if (!listing?.user_id) return
  router.push(`/user/marketplace/chat/${listing.id}/${listing.user_id}`)
}

function showPurchaseDialog(listing: any, type: 'full' | 'fee_only') {
  purchaseListing.value = listing
  purchaseType.value = type
  agreedTerms.value = false
  purchaseVisible.value = true
}

function calculateFee(listing: any): number {
  if (!listing) return 0
  const fee = listing.sell_price * (config.value.feeRate / 100)
  return Math.max(fee, config.value.minFee)
}

function calculateTotal(listing: any, type: string): number {
  if (!listing) return 0
  if (type === 'full') {
    return listing.sell_price + calculateFee(listing)
  }
  return calculateFee(listing)
}

async function confirmPurchase() {
  if (!purchaseListing.value) return
  purchasing.value = true
  try {
    const res = await request.post('/v1/marketplace/orders', {
      listing_id: purchaseListing.value.id,
      payment_method: purchaseType.value
    })
    ElMessage.success('购买成功')
    purchaseVisible.value = false
    detailVisible.value = false

    if (purchaseType.value === 'fee_only' && res.data?.seller) {
      sellerContact.value = res.data.seller
      sellerContactVisible.value = true
    }

    fetchListings()
    userStore.fetchUserInfo()
  } catch (e: any) {
    ElMessage.error(e.message || '购买失败')
  } finally {
    purchasing.value = false
  }
}

function getHostSpec(listing: any, key: string): string {
  const host = listing?.host
  if (!host) return '-'
  const config = host.config || {}
  return config[key] || host[key] || '-'
}

function formatDate(date: string): string {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

function getDaysLeft(date: string): number {
  if (!date) return 0
  const diff = new Date(date).getTime() - Date.now()
  return Math.max(0, Math.ceil(diff / (1000 * 60 * 60 * 24)))
}
</script>

<style scoped lang="scss">
.marketplace-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.marketplace-banner {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 40px;
  border-radius: 12px;
  text-align: center;
  margin-bottom: 24px;

  h1 {
    font-size: 28px;
    margin-bottom: 8px;
  }

  p {
    font-size: 16px;
    opacity: 0.9;
    margin-bottom: 24px;
  }

  .banner-stats {
    display: flex;
    justify-content: center;
    gap: 60px;
  }

  .stat-item {
    display: flex;
    flex-direction: column;

    .stat-num {
      font-size: 32px;
      font-weight: bold;
    }

    .stat-label {
      font-size: 14px;
      opacity: 0.8;
    }
  }
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);

  .filter-right {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .msg-badge {
    margin-left: 4px;
  }
}

.listings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
  min-height: 200px;
}

.listing-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  cursor: pointer;
  transition: all 0.3s;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;

  .product-name {
    font-size: 16px;
    font-weight: 600;
    color: #333;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
    margin-right: 8px;
  }
}

.card-body {
  margin-bottom: 12px;
}

.specs-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  margin-bottom: 8px;
}

.spec-item {
  text-align: center;
  padding: 6px;
  background: #f5f7fa;
  border-radius: 4px;

  .spec-label {
    display: block;
    font-size: 11px;
    color: #999;
  }

  .spec-value {
    font-size: 13px;
    font-weight: 500;
    color: #333;
  }
}

.ip-info,
.expire-info {
  font-size: 12px;
  color: #666;
  margin-top: 6px;

  .expire-days {
    margin-left: 8px;
    color: #e6a23c;
  }
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.price-info {
  .original-price {
    font-size: 12px;
    color: #999;
    text-decoration: line-through;
  }

  .sell-price {
    font-size: 20px;
    font-weight: bold;
    color: #ff4757;

    .price-unit {
      font-size: 12px;
      font-weight: normal;
    }
  }
}

.seller-info {
  display: flex;
  align-items: center;
  gap: 6px;

  .seller-name {
    font-size: 13px;
    color: #666;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

.listing-detail {
  .detail-section {
    margin-bottom: 20px;

    h3 {
      font-size: 15px;
      font-weight: 600;
      margin-bottom: 12px;
      padding-bottom: 8px;
      border-bottom: 1px solid #f0f0f0;
    }
  }

  .detail-specs {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }

  .spec-row {
    display: flex;
    gap: 8px;

    .label {
      color: #999;
      min-width: 40px;
    }
  }

  .description {
    color: #666;
    line-height: 1.6;
  }

  .seller-detail {
    display: flex;
    align-items: center;
    gap: 12px;

    .seller-meta {
      display: flex;
      flex-direction: column;

      .name {
        font-weight: 500;
      }

      .views {
        font-size: 12px;
        color: #999;
      }
    }
  }

  .price-section {
    margin-top: 20px;
    padding: 16px;
    background: #f5f7fa;
    border-radius: 8px;
  }

  .price-compare {
    display: flex;
    justify-content: space-around;
  }

  .price-item {
    text-align: center;

    .label {
      display: block;
      font-size: 12px;
      color: #999;
      margin-bottom: 4px;
    }

    .value {
      font-size: 18px;
      font-weight: bold;

      &.original {
        color: #999;
      }

      &.sell {
        color: #ff4757;
      }

      &.save {
        color: #2ed573;
      }
    }
  }
}

.purchase-confirm {
  .purchase-info {
    margin-bottom: 20px;

    h3 {
      margin-bottom: 8px;
    }

    .purchase-price {
      font-size: 16px;
      color: #333;
    }
  }

  .fee-breakdown {
    background: #f5f7fa;
    padding: 16px;
    border-radius: 8px;
    margin-bottom: 16px;
  }

  .fee-row {
    display: flex;
    justify-content: space-between;
    padding: 6px 0;

    &.total {
      border-top: 1px solid #ddd;
      margin-top: 8px;
      padding-top: 12px;
      font-weight: bold;
      font-size: 16px;
    }
  }

  .balance-info {
    margin-bottom: 16px;

    .balance {
      font-weight: bold;
      color: #409eff;
    }

    .insufficient {
      color: #ff4757;
      font-size: 13px;
    }
  }
}

.seller-contact {
  p {
    margin-bottom: 12px;
    color: #666;
    line-height: 1.6;
  }

  .contact-info {
    background: #f5f7fa;
    padding: 16px;
    border-radius: 8px;
    margin-top: 16px;
  }

  .contact-item {
    display: flex;
    gap: 8px;
    margin-bottom: 8px;

    .label {
      color: #999;
      min-width: 60px;
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
