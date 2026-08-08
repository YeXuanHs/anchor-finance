<template>
  <div class="invoice-list-page">
    <!-- 标签页 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="全部" name="all" />
      <el-tab-pane label="未付" name="unpaid" />
      <el-tab-pane label="已付" name="paid" />
      <el-tab-pane label="已取消" name="cancelled" />
      <el-tab-pane label="已退款" name="refunded" />
    </el-tabs>

    <!-- 搜索筛选 -->
    <el-card shadow="never" class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="发票号/客户名" clearable style="width: 200px" @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch"><el-icon><Search /></el-icon>搜索</el-button>
          <el-button @click="handleReset"><el-icon><Refresh /></el-icon>重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 汇总 -->
    <el-card shadow="never" class="summary-card">
      <el-row :gutter="20">
        <el-col :span="6"><div class="summary-item"><div class="label">总金额</div><div class="value">¥{{ formatMoney(summary.total) }}</div></div></el-col>
        <el-col :span="6"><div class="summary-item"><div class="label">已付</div><div class="value text-green">¥{{ formatMoney(summary.paid) }}</div></div></el-col>
        <el-col :span="6"><div class="summary-item"><div class="label">未付</div><div class="value text-orange">¥{{ formatMoney(summary.unpaid) }}</div></div></el-col>
        <el-col :span="6"><div class="summary-item"><div class="label">已退款</div><div class="value text-red">¥{{ formatMoney(summary.refunded) }}</div></div></el-col>
      </el-row>
    </el-card>

    <!-- 操作栏 -->
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <div class="action-left">
          <el-button type="danger" :disabled="selectedIds.length === 0" @click="handleBatchCancel">批量取消</el-button>
          <el-button type="success" :disabled="selectedIds.length === 0" @click="handleBatchPay">批量确认付款</el-button>
          <el-button @click="handleExport"><el-icon><Download /></el-icon>导出</el-button>
        </div>
        <div class="action-right">
          <el-button circle @click="fetchList"><el-icon><Refresh /></el-icon></el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="invoice_no" label="发票号" width="150">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">{{ row.invoice_no }}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="client_name" label="客户" width="120">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push(`/customer-view/${row.client_id}`)">{{ row.client_name }}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.amount) }}</template>
        </el-table-column>
        <el-table-column prop="paid_amount" label="已付" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.paid_amount) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="due_date" label="到期日期" width="120" />
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleView(row)">查看</el-button>
            <el-button v-if="row.status === 'unpaid'" type="success" link size="small" @click="handleMarkPaid(row)">标记已付</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination-wrapper">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10,20,50,100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Download } from '@element-plus/icons-vue'
import request from '@/utils/http'

const router = useRouter()
const loading = ref(false)
const tableData = ref([])
const activeTab = ref('all')
const selectedIds = ref<number[]>([])
const searchForm = reactive({ keyword: '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const summary = ref<any>({ total: 0, paid: 0, unpaid: 0, refunded: 0 })

const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const getStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { unpaid: 'warning', paid: 'success', cancelled: 'info', refunded: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = { unpaid: '未付', paid: '已付', cancelled: '已取消', refunded: '已退款' }
  return map[status] || '未知'
}

const handleTabChange = () => { pagination.page = 1; fetchList() }
const handleSearch = () => { pagination.page = 1; fetchList() }
const handleReset = () => { searchForm.keyword = ''; pagination.page = 1; fetchList() }
const handleSizeChange = (size: number) => { pagination.page_size = size; pagination.page = 1; fetchList() }
const handlePageChange = (page: number) => { pagination.page = page; fetchList() }
const handleSelectionChange = (rows: any[]) => { selectedIds.value = rows.map((r) => r.id) }

const fetchList = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (activeTab.value !== 'all') params.status = activeTab.value
    if (searchForm.keyword) params.keyword = searchForm.keyword
    const data = await request.get({ url: '/api/admin/invoices', params })
    tableData.value = data?.list || []
    pagination.total = data?.total || 0
    if (data?.summary) summary.value = data.summary
  } catch (error) {
    console.error('获取发票列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleView = (row: any) => { router.push(`/invoice-detail/${row.id}`) }

const handleMarkPaid = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要标记发票 ${row.invoice_no} 为已付吗？`, '确认操作', { type: 'warning' })
    await request.post({ url: `/api/admin/invoices/${row.id}/mark-paid` })
    ElMessage.success('操作成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') console.error('操作失败:', error)
  }
}

const handleBatchCancel = async () => {
  try {
    await ElMessageBox.confirm(`确定要批量取消选中的 ${selectedIds.value.length} 张发票吗？`, '确认操作', { type: 'warning' })
    await request.post({ url: '/api/admin/invoices/batch-cancel', data: { ids: selectedIds.value } })
    ElMessage.success('操作成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') console.error('操作失败:', error)
  }
}

const handleBatchPay = async () => {
  try {
    await ElMessageBox.confirm(`确定要批量确认选中的 ${selectedIds.value.length} 张发票已付款吗？`, '确认操作', { type: 'warning' })
    await request.post({ url: '/api/admin/invoices/batch-pay', data: { ids: selectedIds.value } })
    ElMessage.success('操作成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') console.error('操作失败:', error)
  }
}

const handleExport = async () => {
  try {
    const response = await request.get({ url: '/api/admin/invoices/export', responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `发票列表_${new Date().toISOString().split('T')[0]}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('导出失败:', error)
  }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.invoice-list-page { padding: 16px; }
.search-card { margin-bottom: 16px; :deep(.el-card__body) { padding-bottom: 0; } }
.summary-card { margin-bottom: 16px; }
.summary-item { text-align: center; padding: 10px 0; }
.label { font-size: 14px; color: #86909C; margin-bottom: 8px; }
.value { font-size: 24px; font-weight: 600; }
.text-green { color: #36D391; }
.text-orange { color: #F59E0B; }
.text-red { color: #EF4444; }
.action-card { margin-bottom: 16px; }
.action-bar { display: flex; justify-content: space-between; align-items: center; }
.table-card { :deep(.el-card__body) { padding: 0; } }
.pagination-wrapper { display: flex; justify-content: flex-end; padding: 16px; }
</style>
