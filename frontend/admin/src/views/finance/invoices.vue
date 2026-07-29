<template>
  <div class="invoices-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="发票号">
          <el-input v-model="searchForm.invoice_no" placeholder="发票号" clearable />
        </el-form-item>
        <el-form-item label="用户">
          <el-input v-model="searchForm.username" placeholder="用户名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待开票" value="pending" />
            <el-option label="已开票" value="issued" />
            <el-option label="已作废" value="void" />
            <el-option label="已冲红" value="reversed" />
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
        <h3>发票管理</h3>
        <el-button type="success" @click="handleExport">
          <el-icon><Download /></el-icon>导出
        </el-button>
      </div>
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="invoice_no" label="发票号" width="200" />
        <el-table-column prop="username" label="用户" width="120" />
        <el-table-column prop="title" label="发票抬头" min-width="160" show-overflow-tooltip />
        <el-table-column prop="tax_no" label="税号" width="160" />
        <el-table-column label="金额" width="120">
          <template #default="{ row }">
            <span class="amount">¥{{ row.amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="税额" width="100">
          <template #default="{ row }">¥{{ row.tax_amount?.toFixed(2) || '0.00' }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status]" size="small">{{ statusMap[row.status] || row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewDetail(row)">详情</el-button>
            <el-button link type="success" v-if="row.status === 'pending'" @click="issueInvoice(row.id)">开票</el-button>
            <el-button link type="danger" v-if="row.status === 'issued'" @click="voidInvoice(row.id)">作废</el-button>
          </template>
        </el-table-column>
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

    <el-drawer v-model="showDetail" title="发票详情" size="550px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="发票号">{{ detail.invoice_no }}</el-descriptions-item>
        <el-descriptions-item label="发票类型">{{ detail.type === 'special' ? '增值税专用发票' : '增值税普通发票' }}</el-descriptions-item>
        <el-descriptions-item label="发票抬头">{{ detail.title }}</el-descriptions-item>
        <el-descriptions-item label="税号">{{ detail.tax_no }}</el-descriptions-item>
        <el-descriptions-item label="开户银行">{{ detail.bank_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="银行账号">{{ detail.bank_account || '-' }}</el-descriptions-item>
        <el-descriptions-item label="地址电话">{{ detail.address || '-' }} {{ detail.phone || '' }}</el-descriptions-item>
        <el-descriptions-item label="金额">¥{{ detail.amount?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="税额">¥{{ detail.tax_amount?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item label="价税合计">¥{{ ((detail.amount || 0) + (detail.tax_amount || 0)).toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusType[detail.status]" size="small">{{ statusMap[detail.status] || detail.status }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="关联订单">{{ detail.order_no || '-' }}</el-descriptions-item>
        <el-descriptions-item label="申请时间">{{ detail.created_at }}</el-descriptions-item>
        <el-descriptions-item label="开票时间">{{ detail.issued_at || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const statusMap: Record<string, string> = { pending: '待开票', issued: '已开票', void: '已作废', reversed: '已冲红' }
const statusType: Record<string, string> = { pending: 'warning', issued: 'success', void: 'info', reversed: 'danger' }

const loading = ref(false)
const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const showDetail = ref(false)
const detail = reactive<any>({})
const searchForm = ref({ invoice_no: '', username: '', status: '', date_range: null as string[] | null })

const resetSearch = () => { searchForm.value = { invoice_no: '', username: '', status: '', date_range: null }; fetchData() }

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: page.value, page_size: pageSize.value, invoice_no: searchForm.value.invoice_no, username: searchForm.value.username, status: searchForm.value.status }
    if (searchForm.value.date_range) { params.start_date = searchForm.value.date_range[0]; params.end_date = searchForm.value.date_range[1] }
    const { data } = await request.get('/admin/api/v1/invoices', { params })
    list.value = data.data || []; total.value = data.total || 0
  } catch {} finally { loading.value = false }
}

const viewDetail = async (row: any) => {
  try {
    const { data } = await request.get(`/admin/api/v1/invoices/${row.id}`)
    Object.assign(detail, data.data || data)
  } catch { Object.assign(detail, row) }
  showDetail.value = true
}

const issueInvoice = async (id: number) => {
  await ElMessageBox.confirm('确认开具该发票？', '确认')
  try { await request.put(`/admin/api/v1/invoices/${id}/issue`); ElMessage.success('开票成功'); fetchData() } catch { ElMessage.error('开票失败') }
}

const voidInvoice = async (id: number) => {
  await ElMessageBox.confirm('确认作废该发票？此操作不可撤销。', '确认')
  try { await request.put(`/admin/api/v1/invoices/${id}/void`); ElMessage.success('已作废'); fetchData() } catch { ElMessage.error('操作失败') }
}

const handleExport = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/invoices/export', { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([data]))
    const a = document.createElement('a'); a.href = url; a.download = `invoices_${Date.now()}.xlsx`; a.click()
    window.URL.revokeObjectURL(url)
  } catch { ElMessage.error('导出失败') }
}

onMounted(fetchData)
</script>
<style scoped lang="scss">
.search-bar { margin-bottom: 16px; }
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; h3 { font-size: 18px; font-weight: 600; margin: 0; } }
.amount { color: var(--danger-color); font-weight: 600; }
</style>
