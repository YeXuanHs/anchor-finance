<template>
  <div class="orders-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="订单号">
          <el-input v-model="searchForm.order_no" placeholder="订单号" clearable />
        </el-form-item>
        <el-form-item label="用户">
          <el-input v-model="searchForm.username" placeholder="用户名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待支付" value="pending" />
            <el-option label="已支付" value="paid" />
            <el-option label="处理中" value="processing" />
            <el-option label="已完成" value="completed" />
            <el-option label="已取消" value="cancelled" />
            <el-option label="已退款" value="refunded" />
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
        <h3>订单列表</h3>
        <el-button type="success" @click="handleExport">
          <el-icon><Download /></el-icon>导出
        </el-button>
      </div>
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="order_no" label="订单号" width="180" />
        <el-table-column prop="product_name" label="产品" min-width="140" />
        <el-table-column prop="username" label="用户" width="120" />
        <el-table-column label="金额" width="120">
          <template #default="{ row }">
            <span class="amount">¥{{ row.amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status]" size="small">{{ statusMap[row.status] || row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="下单时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewDetail(row)">详情</el-button>
            <el-button link type="success" v-if="row.status === 'paid'" @click="completeOrder(row.id)">完成</el-button>
            <el-button link type="warning" v-if="['pending', 'paid'].includes(row.status)" @click="cancelOrder(row.id)">取消</el-button>
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

    <el-drawer v-model="showDetail" title="订单详情" size="600px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="订单号">{{ detail.order_no }}</el-descriptions-item>
        <el-descriptions-item label="产品">{{ detail.product_name }}</el-descriptions-item>
        <el-descriptions-item label="用户">{{ detail.username }} (ID: {{ detail.user_id }})</el-descriptions-item>
        <el-descriptions-item label="金额">¥{{ detail.amount?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusType[detail.status]" size="small">{{ statusMap[detail.status] || detail.status }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="支付方式">{{ detail.payment_method || '-' }}</el-descriptions-item>
        <el-descriptions-item label="交易号">{{ detail.trade_no || '-' }}</el-descriptions-item>
        <el-descriptions-item label="下单时间">{{ detail.created_at }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ detail.updated_at }}</el-descriptions-item>
        <el-descriptions-item label="备注">{{ detail.remark || '-' }}</el-descriptions-item>
      </el-descriptions>
      <el-divider content-position="left">订单日志</el-divider>
      <el-timeline>
        <el-timeline-item v-for="log in orderLogs" :key="log.id" :timestamp="log.created_at" placement="top">
          {{ log.content }}
        </el-timeline-item>
      </el-timeline>
    </el-drawer>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const statusMap: Record<string, string> = { pending: '待支付', paid: '已支付', processing: '处理中', completed: '已完成', cancelled: '已取消', refunded: '已退款' }
const statusType: Record<string, string> = { pending: 'warning', paid: 'primary', processing: 'primary', completed: 'success', cancelled: 'info', refunded: 'danger' }

const loading = ref(false)
const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const showDetail = ref(false)
const detail = reactive<any>({})
const orderLogs = ref<any[]>([])
const searchForm = ref({ order_no: '', username: '', status: '', date_range: null as string[] | null })

const resetSearch = () => { searchForm.value = { order_no: '', username: '', status: '', date_range: null }; fetchData() }

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: page.value, page_size: pageSize.value, order_no: searchForm.value.order_no, username: searchForm.value.username, status: searchForm.value.status }
    if (searchForm.value.date_range) { params.start_date = searchForm.value.date_range[0]; params.end_date = searchForm.value.date_range[1] }
    const { data } = await request.get('/admin/api/v1/orders', { params })
    list.value = data.data || []; total.value = data.total || 0
  } catch {} finally { loading.value = false }
}

const viewDetail = async (row: any) => {
  try {
    const { data } = await request.get(`/admin/api/v1/orders/${row.id}`)
    Object.assign(detail, data.data || data)
    orderLogs.value = data.logs || []
  } catch { Object.assign(detail, row); orderLogs.value = [] }
  showDetail.value = true
}

const completeOrder = async (id: number) => {
  await ElMessageBox.confirm('确定将该订单标记为完成？', '确认')
  try { await request.put(`/admin/api/v1/orders/${id}/complete`); ElMessage.success('操作成功'); fetchData() } catch { ElMessage.error('操作失败') }
}

const cancelOrder = async (id: number) => {
  await ElMessageBox.confirm('确定取消该订单？', '确认')
  try { await request.put(`/admin/api/v1/orders/${id}/cancel`); ElMessage.success('操作成功'); fetchData() } catch { ElMessage.error('操作失败') }
}

const handleExport = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/orders/export', { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([data]))
    const a = document.createElement('a'); a.href = url; a.download = `orders_${Date.now()}.xlsx`; a.click()
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
