<template>
  <div class="dashboard-page" v-loading="loading">
    <!-- 欢迎信息 -->
    <div class="welcome-section">
      <div class="welcome-text">
        <h2>欢迎回来，{{ username }}</h2>
        <p>这是您的控制台概览</p>
      </div>
      <div class="welcome-actions">
        <el-button type="primary" @click="$router.push('/products')">
          <el-icon style="margin-right: 6px;"><Plus /></el-icon>
          订购产品
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(26, 115, 232, 0.1); color: #1a73e8;">
          <el-icon :size="24"><Wallet /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">账户余额</div>
          <div class="stat-value">¥{{ stats.balance?.toFixed(2) || '0.00' }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(52, 199, 89, 0.1); color: #34c759;">
          <el-icon :size="24"><Box /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">活跃服务</div>
          <div class="stat-value">{{ stats.active_services || 0 }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(255, 149, 0, 0.1); color: #ff9500;">
          <el-icon :size="24"><Document /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">待付订单</div>
          <div class="stat-value">{{ stats.pending_orders || 0 }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(255, 59, 48, 0.1); color: #ff3b30;">
          <el-icon :size="24"><Tickets /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">待处理工单</div>
          <div class="stat-value">{{ stats.open_tickets || 0 }}</div>
        </div>
      </div>
    </div>

    <!-- 云服务器 -->
    <div class="section-card">
      <div class="section-header">
        <h3>我的云服务器</h3>
        <el-button link type="primary" @click="$router.push('/user/products')">查看全部</el-button>
      </div>
      <el-table :data="servers" style="width: 100%">
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="ip" label="IP地址" />
        <el-table-column prop="config" label="配置" />
        <el-table-column prop="expire_date" label="到期时间" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button link type="primary" size="small">管理</el-button>
            <el-button link type="primary" size="small" v-if="row.status === 'active'">续费</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 最近订单 -->
    <div class="section-card">
      <div class="section-header">
        <h3>最近订单</h3>
        <el-button link type="primary" @click="$router.push('/user/orders')">查看全部</el-button>
      </div>
      <el-table :data="recentOrders" style="width: 100%">
        <el-table-column prop="order_no" label="订单号" width="180" />
        <el-table-column prop="product_name" label="产品" />
        <el-table-column prop="amount" label="金额">
          <template #default="{ row }">
            <span style="color: #1a73e8; font-weight: 600;">¥{{ row.amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="getOrderStatusType(row.status)" size="small">
              {{ getOrderStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="180" />
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="payOrder(row)" v-if="row.status === 'pending'">
              去支付
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
import { Plus, Wallet, Box, Document, Tickets } from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()
const loading = ref(false)
const username = ref('用户')

const stats = ref({
  balance: 0,
  active_services: 0,
  pending_orders: 0,
  open_tickets: 0
})

const servers = ref([])
const recentOrders = ref([])

// 获取数据
const fetchDashboard = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/user/dashboard')
    if (data?.data) {
      stats.value = data.data.stats || {}
      servers.value = data.data.servers || []
      recentOrders.value = data.data.recent_orders || []
      username.value = data.data.username || '用户'
    }
  } catch (error) {
    console.error('获取仪表盘数据失败:', error)
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
    active: '运行中',
    expired: '已过期',
    suspended: '已暂停',
    pending: '待开通'
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
    pending: '待支付',
    paid: '已支付',
    cancelled: '已取消',
    refunded: '已退款'
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
