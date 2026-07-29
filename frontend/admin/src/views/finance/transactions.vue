<template>
  <div class="transactions-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="用户">
          <el-input v-model="searchForm.username" placeholder="用户名/ID" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="充值" value="recharge" />
            <el-option label="消费" value="payment" />
            <el-option label="退款" value="refund" />
            <el-option label="提现" value="withdraw" />
            <el-option label="转入" value="transfer_in" />
            <el-option label="转出" value="transfer_out" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker v-model="searchForm.date_range" type="daterange" range-separator="至" start-placeholder="开始" end-placeholder="结束" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="art-card">
      <div class="table-header">
        <h3>交易记录</h3>
        <el-button type="success" @click="handleExport">
          <el-icon><Download /></el-icon>导出
        </el-button>
      </div>
      <el-table :data="list" v-loading="loading" show-summary :summary-method="getSummary">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="trade_no" label="交易号" width="200" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="typeStyle[row.type]" size="small">{{ typeMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="金额" width="140">
          <template #default="{ row }">
            <span :style="{ color: row.amount > 0 ? '#67c23a' : '#f56c6c', fontWeight: 600 }">
              {{ row.amount > 0 ? '+' : '' }}¥{{ Math.abs(row.amount).toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="余额" width="120">
          <template #default="{ row }">¥{{ row.balance_after?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="160" show-overflow-tooltip />
        <el-table-column prop="related_order" label="关联订单" width="160" />
        <el-table-column prop="created_at" label="时间" width="180" />
      </el-table>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @change="fetchData"
        style="margin-top: 16px; justify-content: flex-end"
      />
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const typeMap: Record<string, string> = { recharge: '充值', payment: '消费', refund: '退款', withdraw: '提现', transfer_in: '转入', transfer_out: '转出' }
const typeStyle: Record<string, string> = { recharge: 'success', payment: 'danger', refund: 'warning', withdraw: 'info', transfer_in: 'success', transfer_out: 'info' }

const loading = ref(false)
const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const searchForm = ref({ username: '', type: '', date_range: null as string[] | null })

const resetSearch = () => { searchForm.value = { username: '', type: '', date_range: null }; fetchData() }

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: page.value, page_size: pageSize.value, username: searchForm.value.username, type: searchForm.value.type }
    if (searchForm.value.date_range) { params.start_date = searchForm.value.date_range[0]; params.end_date = searchForm.value.date_range[1] }
    const { data } = await request.get('/admin/api/v1/transactions', { params })
    list.value = data.data || []; total.value = data.total || 0
  } catch {} finally { loading.value = false }
}

const getSummary = ({ columns, data }: any) => {
  const sums: string[] = []
  columns.forEach((col: any, index: number) => {
    if (index === 0) { sums[index] = '合计'; return }
    if (col.label === '金额') {
      const totalVal = data.reduce((sum: number, row: any) => sum + (row.amount || 0), 0)
      sums[index] = `${totalVal > 0 ? '+' : ''}¥${Math.abs(totalVal).toFixed(2)}`
    } else { sums[index] = '' }
  })
  return sums
}

const handleExport = async () => {
  try {
    const params: any = { ...searchForm.value }
    delete params.date_range
    if (searchForm.value.date_range) { params.start_date = searchForm.value.date_range[0]; params.end_date = searchForm.value.date_range[1] }
    const { data } = await request.get('/admin/api/v1/transactions/export', { params, responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([data]))
    const a = document.createElement('a'); a.href = url; a.download = `transactions_${Date.now()}.xlsx`; a.click()
    window.URL.revokeObjectURL(url)
  } catch { ElMessage.error('导出失败') }
}

onMounted(fetchData)
</script>
<style scoped lang="scss">
.search-bar { margin-bottom: 16px; }
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; h3 { font-size: 18px; font-weight: 600; margin: 0; } }
</style>
