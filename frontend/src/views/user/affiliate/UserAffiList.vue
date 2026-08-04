<template>
  <div class="user-affi-list-page">
    <div class="page-header">
      <h1 class="page-title">推荐用户列表</h1>
      <div class="header-actions">
        <el-select v-model="filterLevel" placeholder="推荐层级" clearable style="width: 140px;">
          <el-option label="全部层级" value="" />
          <el-option label="一级推荐" value="1" />
          <el-option label="二级推荐" value="2" />
        </el-select>
        <el-input
          v-model="searchKey"
          placeholder="搜索用户名/邮箱"
          clearable
          class="search-input"
          :prefix-icon="Search"
        />
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="stats-row">
      <el-card shadow="never" class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: linear-gradient(135deg, #0056FF, #4080FF);">
            <el-icon :size="24"><User /></el-icon>
          </div>
          <div class="stat-info">
            <span class="stat-value">{{ totalUsers }}</span>
            <span class="stat-label">推荐用户总数</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: linear-gradient(135deg, #52c41a, #73d13d);">
            <el-icon :size="24"><CircleCheck /></el-icon>
          </div>
          <div class="stat-info">
            <span class="stat-value">{{ activeUsers }}</span>
            <span class="stat-label">活跃用户</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: linear-gradient(135deg, #fa8c16, #ffc53d);">
            <el-icon :size="24"><Wallet /></el-icon>
          </div>
          <div class="stat-info">
            <span class="stat-value">¥{{ totalSpent }}</span>
            <span class="stat-label">用户总消费</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: linear-gradient(135deg, #722ed1, #b37feb);">
            <el-icon :size="24"><Coin /></el-icon>
          </div>
          <div class="stat-info">
            <span class="stat-value">¥{{ totalCommission }}</span>
            <span class="stat-label">累计返利</span>
          </div>
        </div>
      </el-card>
    </div>

    <!-- Data Table -->
    <el-card shadow="never" class="table-card">
      <el-table :data="paginatedUsers" stripe style="width: 100%" v-loading="loading" empty-text="暂无推荐用户">
        <el-table-column prop="username" label="用户" min-width="150">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="32" class="record-avatar">{{ row.avatar }}</el-avatar>
              <div class="user-info">
                <span class="user-name">{{ row.username }}</span>
                <span class="user-email">{{ row.email }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="level" label="推荐层级" width="110">
          <template #default="{ row }">
            <el-tag :type="row.level === 1 ? 'primary' : 'info'" size="small" effect="light">
              {{ row.level === 1 ? '一级推荐' : '二级推荐' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="orderCount" label="订单数" width="90" align="center" />
        <el-table-column prop="totalSpent" label="累计消费" width="120">
          <template #default="{ row }">
            <span class="amount-text">¥{{ row.totalSpent }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="commission" label="贡献返利" width="120">
          <template #default="{ row }">
            <span class="commission-text">¥{{ row.commission }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small" effect="light" round>
              {{ row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="registeredAt" label="注册时间" width="170" sortable />
        <el-table-column prop="lastActiveAt" label="最后活跃" width="170" />
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="filteredUsers.length"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          background
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Search, User, CircleCheck, Wallet, Coin } from '@element-plus/icons-vue'
import request from '@/utils/request'

const searchKey = ref('')
const filterLevel = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

interface AffiUser {
  username: string
  email: string
  avatar: string
  level: number
  orderCount: number
  totalSpent: string
  commission: string
  status: string
  statusText: string
  registeredAt: string
  lastActiveAt: string
}

const users = ref<AffiUser[]>([])

const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v2/affiliate/info')
    if (data?.data) {
      users.value = data.data.list || data.data || []
    }
  } catch (e) {
    console.error('Failed to fetch affiliate data:', e)
  } finally {
    loading.value = false
  }
})

const totalUsers = computed(() => users.value.length)
const activeUsers = computed(() => users.value.filter(u => u.status === 'active').length)
const totalSpent = computed(() => {
  const sum = users.value.reduce((acc, u) => acc + parseFloat(u.totalSpent.replace(/,/g, '')), 0)
  return sum.toLocaleString('zh-CN', { minimumFractionDigits: 2 })
})
const totalCommission = computed(() => {
  const sum = users.value.reduce((acc, u) => acc + parseFloat(u.commission.replace(/,/g, '')), 0)
  return sum.toLocaleString('zh-CN', { minimumFractionDigits: 2 })
})

const filteredUsers = computed(() => {
  let result = users.value
  if (filterLevel.value) {
    result = result.filter(u => u.level === Number(filterLevel.value))
  }
  if (searchKey.value) {
    const key = searchKey.value.toLowerCase()
    result = result.filter(u =>
      u.username.toLowerCase().includes(key) || u.email.toLowerCase().includes(key)
    )
  }
  return result
})

const paginatedUsers = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredUsers.value.slice(start, start + pageSize.value)
})
</script>

<style scoped>
.user-affi-list-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-input {
  width: 240px;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.stat-card :deep(.el-card__body) {
  padding: 20px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
}

.stat-label {
  font-size: 13px;
  color: #909399;
}

.table-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.table-card :deep(.el-card__body) {
  padding: 0;
}

.table-card :deep(.el-table th.el-table__cell) {
  background: #fafbfc;
  color: #606266;
  font-weight: 600;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.record-avatar {
  background: linear-gradient(135deg, #0056FF, #4080FF);
  color: #fff;
  font-size: 13px;
  font-weight: 600;
  flex-shrink: 0;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.user-name {
  font-weight: 600;
  color: #303133;
  font-size: 14px;
}

.user-email {
  font-size: 12px;
  color: #909399;
}

.amount-text {
  font-weight: 600;
  color: #303133;
}

.commission-text {
  font-weight: 600;
  color: #fa8c16;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px 20px;
  border-top: 1px solid #e8ecf1;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .header-actions {
    width: 100%;
    flex-direction: column;
  }

  .search-input {
    width: 100%;
  }

  .stats-row {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
