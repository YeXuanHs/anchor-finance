<template>
  <div class="user-detail page-container">
    <div class="page-header">
      <el-button @click="$router.back()">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>用户详情</h2>
    </div>
    
    <div class="detail-grid">
      <div class="info-card art-card">
        <h3>基本信息</h3>
        <div class="info-item">
          <span class="label">用户名</span>
          <span class="value">{{ user.username }}</span>
        </div>
        <div class="info-item">
          <span class="label">邮箱</span>
          <span class="value">{{ user.email }}</span>
        </div>
        <div class="info-item">
          <span class="label">手机</span>
          <span class="value">{{ user.phone || '-' }}</span>
        </div>
        <div class="info-item">
          <span class="label">余额</span>
          <span class="value amount">¥{{ user.balance?.toFixed(2) }}</span>
        </div>
        <div class="info-item">
          <span class="label">状态</span>
          <span class="value">
            <span class="status-tag" :class="user.status">
              {{ user.status === 'active' ? '正常' : '禁用' }}
            </span>
          </span>
        </div>
        <div class="info-item">
          <span class="label">注册时间</span>
          <span class="value">{{ user.created_at }}</span>
        </div>
      </div>
      
      <div class="info-card art-card">
        <h3>产品/服务</h3>
        <el-table :data="user.products" style="width: 100%">
          <el-table-column prop="name" label="产品名称" />
          <el-table-column prop="status" label="状态">
            <template #default="{ row }">
              <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
                {{ row.status === 'active' ? '使用中' : '已过期' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="expire_date" label="到期时间" />
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const route = useRoute()
const userId = route.params.id

const user = ref({
  id: userId,
  username: '',
  email: '',
  phone: '',
  balance: 0,
  status: 'active',
  created_at: '',
  products: [] as any[]
})

const fetchUser = async () => {
  try {
    const { data } = await request.get(`/admin/users/${userId}`)
    if (data?.data) {
      Object.assign(user.value, data.data)
    }
  } catch {
    ElMessage.error('获取用户信息失败')
  }
}

onMounted(() => {
  fetchUser()
})
</script>

<style scoped lang="scss">
.user-detail {
  .page-header {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;
    
    h2 {
      margin: 0;
      font-size: 20px;
    }
  }
  
  .detail-grid {
    display: grid;
    grid-template-columns: 1fr 2fr;
    gap: 20px;
    
    @media (max-width: 1200px) {
      grid-template-columns: 1fr;
    }
  }
  
  .info-card {
    h3 {
      font-size: 16px;
      font-weight: 600;
      margin: 0 0 16px;
      padding-bottom: 12px;
      border-bottom: 1px solid var(--border-color);
    }
    
    .info-item {
      display: flex;
      padding: 12px 0;
      border-bottom: 1px solid #f5f5f5;
      
      &:last-child {
        border-bottom: none;
      }
      
      .label {
        width: 100px;
        color: var(--text-secondary);
      }
      
      .value {
        flex: 1;
        color: var(--text-primary);
        
        &.amount {
          color: var(--danger-color);
          font-weight: 600;
        }
      }
    }
  }
}
</style>
