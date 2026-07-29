<template>
  <div class="host-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>我的主机</span>
          <el-button type="primary" @click="$router.push('/products')">
            <el-icon><Plus /></el-icon>
            订购新主机
          </el-button>
        </div>
      </template>

      <el-table :data="hosts" style="width: 100%" v-loading="loading">
        <el-table-column prop="product_name" label="产品名称" />
        <el-table-column prop="domain" label="域名/IP" />
        <el-table-column prop="billing_cycle" label="计费周期" />
        <el-table-column prop="amount" label="金额">
          <template #default="{ row }">
            <span class="amount">¥{{ row.amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="next_due_date" label="下次续费" />
        <el-table-column label="操作" width="250">
          <template #default="{ row }">
            <el-button size="small" @click="viewDetail(row)">详情</el-button>
            <el-button size="small" type="success" @click="renewHost(row)" v-if="row.status === 'active'">续费</el-button>
            <el-button size="small" type="warning" @click="upgradeHost(row)" v-if="row.status === 'active'">升降级</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()

const loading = ref(false)
const hosts = ref([])

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    active: 'success',
    suspended: 'danger',
    pending: 'warning',
    terminated: 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '使用中',
    suspended: '已暂停',
    pending: '待开通',
    terminated: '已终止'
  }
  return map[status] || status
}

const viewDetail = (host: any) => {
  router.push(`/user/products/${host.id}`)
}

const renewHost = (host: any) => {
  router.push({ path: '/user/batch-renew', query: { host_id: host.id } })
}

const upgradeHost = (host: any) => {
  router.push({ path: '/user/upgrade', query: { host_id: host.id } })
}

const loadHosts = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/user/products')
    hosts.value = data?.data?.list || data?.data?.items || data?.data || []
  } catch {
    hosts.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadHosts()
})
</script>

<style scoped lang="scss">
.host-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .amount {
    color: #f56c6c;
    font-weight: bold;
  }
}
</style>
