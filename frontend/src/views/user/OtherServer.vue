<template>
  <div class="other-server-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>其他服务器</span>
          <el-button type="primary" @click="$router.push('/products')">
            <el-icon><Plus /></el-icon>
            订购新服务
          </el-button>
        </div>
      </template>

      <el-table :data="servers" style="width: 100%" v-loading="loading">
        <el-table-column prop="product_name" label="产品名称" min-width="160" />
        <el-table-column prop="ip" label="IP地址" width="150" />
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag>{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="billing_cycle" label="计费周期" width="100" />
        <el-table-column prop="amount" label="金额" width="100">
          <template #default="{ row }">
            <span class="amount">¥{{ row.amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="next_due_date" label="到期时间" width="120" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewDetail(row)">详情</el-button>
            <el-button size="small" type="success" @click="handlePowerOn(row)" v-if="row.status === 'suspended' || row.status === 'stopped'">开机</el-button>
            <el-button size="small" type="warning" @click="handlePowerOff(row)" v-if="row.status === 'active'">关机</el-button>
            <el-button size="small" type="info" @click="handleReboot(row)" v-if="row.status === 'active'">重启</el-button>
            <el-button size="small" type="primary" @click="renewHost(row)" v-if="row.status === 'active'">续费</el-button>
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
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()

const loading = ref(false)
const servers = ref([])

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    active: 'success',
    suspended: 'danger',
    stopped: 'warning',
    pending: 'info',
    terminated: 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '运行中',
    suspended: '已暂停',
    stopped: '已关机',
    pending: '待开通',
    terminated: '已终止'
  }
  return map[status] || status
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    storage: '存储服务器',
    backup: '备份服务器',
    gpu: 'GPU服务器',
    database: '数据库服务器',
    mail: '邮件服务器',
    other: '其他'
  }
  return map[type] || type || '其他'
}

const viewDetail = (server: any) => {
  router.push(`/user/products/${server.id}`)
}

const renewHost = (server: any) => {
  router.push({ path: '/user/batch-renew', query: { host_id: server.id } })
}

const handlePowerOn = async (server: any) => {
  try {
    await ElMessageBox.confirm('确定要启动该服务器吗？', '确认开机', { type: 'info' })
    await request.post(`/api/v1/user/products/${server.id}/power-on`)
    ElMessage.success('开机指令已发送')
    loadServers()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e?.message || '操作失败')
  }
}

const handlePowerOff = async (server: any) => {
  try {
    await ElMessageBox.confirm('确定要关闭该服务器吗？此操作可能导致服务中断。', '确认关机', { type: 'warning' })
    await request.post(`/api/v1/user/products/${server.id}/power-off`)
    ElMessage.success('关机指令已发送')
    loadServers()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e?.message || '操作失败')
  }
}

const handleReboot = async (server: any) => {
  try {
    await ElMessageBox.confirm('确定要重启该服务器吗？', '确认重启', { type: 'warning' })
    await request.post(`/api/v1/user/products/${server.id}/reboot`)
    ElMessage.success('重启指令已发送')
    loadServers()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e?.message || '操作失败')
  }
}

const loadServers = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/user/products', { params: { group: 'other' } })
    servers.value = data?.data?.list || data?.data?.items || data?.data || []
  } catch {
    servers.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadServers()
})
</script>

<style scoped lang="scss">
.other-server-page {
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
