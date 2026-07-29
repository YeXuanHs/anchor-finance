<template>
  <div class="contract-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>合同管理</span>
        </div>
      </template>

      <el-table :data="contracts" style="width: 100%">
        <el-table-column prop="contract_no" label="合同编号" />
        <el-table-column prop="product_name" label="产品名称" />
        <el-table-column prop="amount" label="合同金额">
          <template #default="{ row }">
            <span class="amount">¥{{ row.amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="start_date" label="开始日期" />
        <el-table-column prop="end_date" label="结束日期" />
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button size="small" @click="viewContract(row)">查看</el-button>
            <el-button size="small" type="primary" @click="downloadContract(row)">下载</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import request from '@/utils/request'

const router = useRouter()

const contracts = ref([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/contracts')
    contracts.value = data.data?.list || data.list || data.data || []
  } catch (e) { console.error(e) } finally { loading.value = false }
})

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    active: 'success',
    expired: 'info',
    pending: 'warning'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '生效中',
    expired: '已过期',
    pending: '待签署'
  }
  return map[status] || status
}

const viewContract = (contract: any) => {
  router.push(`/user/contracts/${contract.id}`)
}

const downloadContract = (contract: any) => {
  window.open(`/api/v1/contracts/${contract.id}/download`)
}
</script>

<style scoped lang="scss">
.contract-page {
  .amount {
    color: #f56c6c;
    font-weight: bold;
  }
}
</style>
