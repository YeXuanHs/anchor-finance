<template>
  <div class="dashboard-page" v-loading="loading">
    <div class="welcome-section">
      <div class="welcome-text">
        <h2>{{ $t('dashboard.welcome') }}，{{ username }}</h2>
        <p>{{ $t('dashboard.overview') }}</p>
      </div>
      <div class="welcome-actions">
        <el-button type="primary" @click="$router.push('/products')">
          <el-icon style="margin-right: 6px;"><Plus /></el-icon>
          {{ $t('dashboard.orderProduct') }}
        </el-button>
      </div>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(26, 115, 232, 0.1); color: #1a73e8;">
          <el-icon :size="24"><Wallet /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">{{ $t('dashboard.balance') }}</div>
          <div class="stat-value">¥{{ stats.balance?.toFixed(2) || '0.00' }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(52, 199, 89, 0.1); color: #34c759;">
          <el-icon :size="24"><Box /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">{{ $t('dashboard.activeServices') }}</div>
          <div class="stat-value">{{ stats.active_services || 0 }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(255, 149, 0, 0.1); color: #ff9500;">
          <el-icon :size="24"><Document /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">{{ $t('dashboard.pendingOrders') }}</div>
          <div class="stat-value">{{ stats.pending_orders || 0 }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(255, 59, 48, 0.1); color: #ff3b30;">
          <el-icon :size="24"><Tickets /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">{{ $t('dashboard.openTickets') }}</div>
          <div class="stat-value">{{ stats.open_tickets || 0 }}</div>
        </div>
      </div>
    </div>

    <div class="section-card">
      <div class="section-header">
        <h3>{{ $t('dashboard.myCloudServers') }}</h3>
        <el-button link type="primary" @click="$router.push('/user/products')">{{ $t('dashboard.viewAll') }}</el-button>
      </div>
      <el-table :data="servers" style="width: 100%">
        <el-table-column prop="name" :label="$t('dashboard.name')" />
        <el-table-column prop="ip" :label="$t('dashboard.ipAddress')" />
        <el-table-column prop="config" :label="$t('dashboard.config')" />
        <el-table-column prop="expire_date" :label="$t('dashboard.dueDate')" />
        <el-table-column prop="status" :label="$t('dashboard.status')">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('dashboard.manage')" width="150">
          <template #default="{ row }">
            <el-button link type="primary" size="small">{{ $t('dashboard.manage') }}</el-button>
            <el-button link type="primary" size="small" v-if="row.status === 'active'">{{ $t('dashboard.renew') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="section-card">
      <div class="section-header">
        <h3>{{ $t('dashboard.recentOrders') }}</h3>
        <el-button link type="primary" @click="$router.push('/user/orders')">{{ $t('dashboard.viewAll') }}</el-button>
      </div>
      <el-table :data="recentOrders" style="width: 100%">
        <el-table-column prop="order_no" :label="$t('dashboard.orderNo')" width="180" />
        <el-table-column prop="product_name" :label="$t('dashboard.product')" />
        <el-table-column prop="amount" :label="$t('dashboard.amount')">
          <template #default="{ row }">
            <span style="color: #1a73e8; font-weight: 600;">¥{{ row.amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('dashboard.status')">
          <template #default="{ row }">
            <el-tag :type="getOrderStatusType(row.status)" size="small">
              {{ getOrderStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.time')" width="180" />
        <el-table-column :label="$t('common.operating')" width="120">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="payOrder(row)" v-if="row.status === 'pending'">
              {{ $t('dashboard.goPay') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Plus, Wallet, Box, Document, Tickets } from '@element-plus/icons-vue'
import request from '@/utils/request'

const { t } = useI18n()
const router = useRouter()
const loading = ref(false)
const username = ref('')

const stats = ref({
  balance: 0,
  active_services: 0,
  pending_orders: 0,
  open_tickets: 0
})

const servers = ref([])
const recentOrders = ref([])

const fetchDashboard = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/user/dashboard')
    if (data?.data) {
      stats.value = data.data.stats || {}
      servers.value = data.data.servers || []
      recentOrders.value = data.data.recent_orders || []
      username.value = data.data.username || ''
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    active: 'success',
    expired: 'danger',
    suspended: 'warning',
    pending: 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: t('dashboard.statusRunning'),
    expired: t('dashboard.statusExpired'),
    suspended: t('dashboard.statusSuspended'),
    pending: t('dashboard.statusPending')
  }
  return map[status] || status
}

const getOrderStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    paid: 'success',
    cancelled: 'info',
    refunded: 'danger'
  }
  return map[status] || 'info'
}

const getOrderStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: t('dashboard.statusUnpaid'),
    paid: t('dashboard.statusPaid'),
    cancelled: t('dashboard.statusCancelled'),
    refunded: t('dashboard.statusRefunded')
  }
  return map[status] || status
}

const payOrder = (order: any) => {
  router.push(`/checkout?order=${order.order_no}`)
}

onMounted(() => {
  fetchDashboard()
})
</script>

<style scoped lang="scss">
.dashboard-page {
  .welcome-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    
    .welcome-text {
      h2 {
        font-size: 24px;
        font-weight: 600;
        margin: 0 0 4px;
      }
      
      p {
        color: #86868b;
        margin: 0;
      }
    }
  }
  
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 16px;
    margin-bottom: 24px;
    
    @media (max-width: 992px) {
      grid-template-columns: repeat(2, 1fr);
    }
  }
  
  .stat-card {
    background: #fff;
    border-radius: 12px;
    padding: 20px;
    display: flex;
    align-items: center;
    gap: 16px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    
    .stat-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    
    .stat-info {
      .stat-label {
        font-size: 13px;
        color: #86868b;
        margin-bottom: 4px;
      }
      
      .stat-value {
        font-size: 24px;
        font-weight: 600;
        color: #1d1d1f;
      }
    }
  }
  
  .section-card {
    background: #fff;
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 16px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    
    .section-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;
      
      h3 {
        font-size: 16px;
        font-weight: 600;
        margin: 0;
      }
    }
  }
}
</style>
