<template>
  <div class="dashboard-page">
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card" v-for="stat in stats" :key="stat.label">
        <div class="stat-icon" :class="stat.type">
          <el-icon :size="24"><component :is="stat.icon" /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-value">{{ stat.value }}</div>
          <div class="stat-change" :class="stat.changeType">
            {{ stat.change }}
          </div>
        </div>
      </div>
    </div>
    
    <!-- 图表区域 -->
    <div class="charts-grid">
      <div class="chart-card">
        <div class="chart-header">
          <h3>收入趋势</h3>
          <el-radio-group v-model="chartPeriod" size="small">
            <el-radio-button label="week">本周</el-radio-button>
            <el-radio-button label="month">本月</el-radio-button>
            <el-radio-button label="year">本年</el-radio-button>
          </el-radio-group>
        </div>
        <div class="chart-content">
          <div class="chart-placeholder">
            <el-icon :size="48" color="#e5e5ea"><TrendCharts /></el-icon>
            <p>收入趋势图表</p>
          </div>
        </div>
      </div>
      
      <div class="chart-card">
        <div class="chart-header">
          <h3>产品分布</h3>
        </div>
        <div class="chart-content">
          <div class="chart-placeholder">
            <el-icon :size="48" color="#e5e5ea"><PieChart /></el-icon>
            <p>产品分布图表</p>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 最近订单 -->
    <div class="recent-section">
      <div class="section-card">
        <div class="section-header">
          <h3>最近订单</h3>
          <el-button type="primary" link @click="$router.push('/orders')">
            查看全部
          </el-button>
        </div>
        <el-table :data="recentOrders" style="width: 100%">
          <el-table-column prop="order_no" label="订单号" width="180" />
          <el-table-column prop="product" label="产品" />
          <el-table-column prop="user" label="用户" />
          <el-table-column prop="amount" label="金额">
            <template #default="{ row }">
              <span class="amount">¥{{ row.amount?.toFixed(2) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态">
            <template #default="{ row }">
              <span class="status-tag" :class="row.status">
                {{ getStatusText(row.status) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="时间" width="180" />
        </el-table>
      </div>
    </div>
    
    <!-- 最近工单 -->
    <div class="recent-section">
      <div class="section-card">
        <div class="section-header">
          <h3>待处理工单</h3>
          <el-button type="primary" link @click="$router.push('/tickets')">
            查看全部
          </el-button>
        </div>
        <el-table :data="recentTickets" style="width: 100%">
          <el-table-column prop="ticket_no" label="工单号" width="150" />
          <el-table-column prop="title" label="标题" />
          <el-table-column prop="user" label="用户" />
          <el-table-column prop="priority" label="优先级">
            <template #default="{ row }">
              <el-tag :type="getPriorityType(row.priority)" size="small">
                {{ row.priority }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="时间" width="180" />
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button type="primary" link @click="viewTicket(row)">
                处理
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const chartPeriod = ref('month')

const stats = [
  { label: '今日收入', value: '¥12,580', change: '+12.5%', changeType: 'up', icon: 'Wallet', type: 'primary' },
  { label: '新增用户', value: '156', change: '+8.2%', changeType: 'up', icon: 'User', type: 'success' },
  { label: '待处理工单', value: '23', change: '-5', changeType: 'down', icon: 'Tickets', type: 'warning' },
  { label: '服务器总数', value: '1,234', change: '+12', changeType: 'up', icon: 'Monitor', type: 'info' }
]

const recentOrders = [
  { order_no: 'ORD20260727001', product: '云服务器 - 2核4G', user: '张三', amount: 299.00, status: 'active', created_at: '2026-07-27 18:00' },
  { order_no: 'ORD20260727002', product: '独立服务器 - E5-2680', user: '李四', amount: 1299.00, status: 'pending', created_at: '2026-07-27 17:30' },
  { order_no: 'ORD20260727003', product: '虚拟主机 - 企业版', user: '王五', amount: 99.00, status: 'active', created_at: '2026-07-27 17:00' }
]

const recentTickets = [
  { ticket_no: 'TK20260727001', title: '服务器无法连接', user: '张三', priority: '高', created_at: '2026-07-27 18:00' },
  { ticket_no: 'TK20260727002', title: '如何升级配置', user: '李四', priority: '中', created_at: '2026-07-27 17:30' },
  { ticket_no: 'TK20260727003', title: '退款申请', user: '王五', priority: '高', created_at: '2026-07-27 17:00' }
]

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '已完成',
    pending: '待处理',
    cancelled: '已取消'
  }
  return map[status] || status
}

const getPriorityType = (priority: string) => {
  const map: Record<string, string> = {
    '高': 'danger',
    '中': 'warning',
    '低': 'info'
  }
  return map[priority] || 'info'
}

const viewTicket = (ticket: any) => {
  router.push(`/tickets/${ticket.ticket_no}`)
}
</script>

<style scoped lang="scss">
.dashboard-page {
  animation: fadeIn 0.3s ease;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 20px;
  
  @media (max-width: 1200px) {
    grid-template-columns: repeat(2, 1fr);
  }
}

.stat-card {
  background: var(--bg-card);
  border-radius: var(--border-radius);
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-md);
  }
  
  .stat-icon {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    
    &.primary {
      background: var(--primary-bg);
      color: var(--primary-color);
    }
    
    &.success {
      background: rgba(52, 199, 89, 0.1);
      color: var(--success-color);
    }
    
    &.warning {
      background: rgba(255, 149, 0, 0.1);
      color: var(--warning-color);
    }
    
    &.info {
      background: rgba(142, 142, 147, 0.1);
      color: var(--info-color);
    }
  }
  
  .stat-info {
    flex: 1;
    
    .stat-label {
      font-size: 13px;
      color: var(--text-secondary);
      margin-bottom: 4px;
    }
    
    .stat-value {
      font-size: 24px;
      font-weight: 600;
      color: var(--text-primary);
    }
    
    .stat-change {
      font-size: 12px;
      margin-top: 4px;
      
      &.up {
        color: var(--success-color);
      }
      
      &.down {
        color: var(--danger-color);
      }
    }
  }
}

.charts-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
  
  @media (max-width: 1200px) {
    grid-template-columns: 1fr;
  }
}

.chart-card {
  background: var(--bg-card);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  
  .chart-header {
    padding: 16px 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid var(--border-color);
    
    h3 {
      font-size: 16px;
      font-weight: 600;
      margin: 0;
    }
  }
  
  .chart-content {
    padding: 20px;
    height: 300px;
    
    .chart-placeholder {
      height: 100%;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      color: var(--text-secondary);
      
      p {
        margin-top: 12px;
      }
    }
  }
}

.recent-section {
  margin-bottom: 20px;
}

.section-card {
  background: var(--bg-card);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  
  .section-header {
    padding: 16px 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid var(--border-color);
    
    h3 {
      font-size: 16px;
      font-weight: 600;
      margin: 0;
    }
  }
}

.amount {
  color: var(--danger-color);
  font-weight: 600;
}

.status-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  
  &::before {
    content: '';
    width: 6px;
    height: 6px;
    border-radius: 50%;
  }
  
  &.active {
    background: rgba(52, 199, 89, 0.1);
    color: var(--success-color);
    
    &::before {
      background: var(--success-color);
    }
  }
  
  &.pending {
    background: rgba(255, 149, 0, 0.1);
    color: var(--warning-color);
    
    &::before {
      background: var(--warning-color);
    }
  }
}
</style>
