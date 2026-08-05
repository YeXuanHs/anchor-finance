<template>
  <div class="upgrades-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>升级管理</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="订单号">
          <el-input v-model="searchForm.order_id" placeholder="订单号" clearable />
        </el-form-item>
        <el-form-item label="客户">
          <el-input v-model="searchForm.username" placeholder="客户用户名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待处理" value="pending" />
            <el-option label="待支付" value="awaiting_payment" />
            <el-option label="已完成" value="completed" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
        </el-form-item>
        <el-form-item label="升级类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="产品升级" value="product" />
            <el-option label="配置升降级" value="configoption" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 统计卡片 -->
      <el-row :gutter="20" class="stat-section">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.total }}</div>
            <div class="stat-label">升级总数</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.pending }}</div>
            <div class="stat-label">待处理</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.awaiting_payment }}</div>
            <div class="stat-label">待支付</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.completed }}</div>
            <div class="stat-label">已完成</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 升级列表 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="order_id" label="订单号" width="120">
          <template #default="{ row }">
            <el-link type="primary" @click="goToOrder(row.order_id)">{{ row.order_id }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="客户" width="120" />
        <el-table-column prop="type" label="升级类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'product' ? 'primary' : 'warning'" size="small">
              {{ row.type === 'product' ? '产品升级' : '配置升降级' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="product_name" label="产品名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="original_value" label="原始配置" min-width="120" show-overflow-tooltip />
        <el-table-column prop="upgrade_value" label="升级配置" min-width="120" show-overflow-tooltip />
        <el-table-column prop="amount" label="升级费用" width="120" align="right">
          <template #default="{ row }">
            <span class="amount">{{ row.currency_prefix }}{{ row.amount }}{{ row.currency_suffix }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="170" />
        <el-table-column prop="completed_at" label="完成时间" width="170">
          <template #default="{ row }">{{ row.completed_at || '-' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleDetail(row)">详情</el-button>
            <el-button v-if="row.status === 'pending'" type="success" link @click="handleApprove(row)">批准</el-button>
            <el-button v-if="row.status === 'pending'" type="warning" link @click="handleCancel(row)">取消</el-button>
            <el-button v-if="row.status === 'awaiting_payment'" type="info" link @click="handleMarkPaid(row)">标记已付</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="升级详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="升级ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item label="订单号">{{ detailData.order_id }}</el-descriptions-item>
        <el-descriptions-item label="客户">{{ detailData.username }}</el-descriptions-item>
        <el-descriptions-item label="升级类型">{{ detailData.type === 'product' ? '产品升级' : '配置升降级' }}</el-descriptions-item>
        <el-descriptions-item label="产品名称" :span="2">{{ detailData.product_name }}</el-descriptions-item>
        <el-descriptions-item label="原始配置">{{ detailData.original_value }}</el-descriptions-item>
        <el-descriptions-item label="升级配置">{{ detailData.upgrade_value }}</el-descriptions-item>
        <el-descriptions-item label="升级费用">
          {{ detailData.currency_prefix }}{{ detailData.amount }}{{ detailData.currency_suffix }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(detailData.status)">{{ getStatusLabel(detailData.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="申请时间">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="完成时间">{{ detailData.completed_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ detailData.remark || '-' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const router = useRouter()

const statusMap: Record<string, { label: string; type: string }> = {
  pending: { label: '待处理', type: 'warning' },
  awaiting_payment: { label: '待支付', type: 'primary' },
  completed: { label: '已完成', type: 'success' },
  cancelled: { label: '已取消', type: 'info' }
}

const getStatusType = (status: string) => (statusMap[status]?.type || 'info') as any
const getStatusLabel = (status: string) => statusMap[status]?.label || status

const loading = ref(false)
const detailDialogVisible = ref(false)

const searchForm = reactive({ order_id: '', username: '', status: '', type: '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const stats = reactive({ total: 0, pending: 0, awaiting_payment: 0, completed: 0 })
const detailData = ref<any>({})

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (searchForm.order_id) params.order_id = searchForm.order_id
    if (searchForm.username) params.username = searchForm.username
    if (searchForm.status) params.status = searchForm.status
    if (searchForm.type) params.type = searchForm.type
    const res = await request.get({ url: '/api/admin/upgrades', params })
    tableData.value = res?.list || res?.data || res || []
    pagination.total = res?.total || 0
  } catch { ElMessage.error('获取升级列表失败') } finally { loading.value = false }
}

const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/upgrades/stats' })
    if (data) Object.assign(stats, data)
  } catch { /* ignore */ }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.order_id = ''; searchForm.username = ''; searchForm.status = ''; searchForm.type = ''; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const goToOrder = (id: string | number) => {
  router.push({ path: `/finance/orders/order-detail/${id}` })
}

const handleDetail = (row: any) => {
  detailData.value = { ...row }
  detailDialogVisible.value = true
}

const handleApprove = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定批准升级请求 #${row.id} 吗？`, '批准升级')
    await request.post({ url: `/api/admin/upgrades/${row.id}/approve`, showSuccessMessage: true })
    fetchData(); fetchStats()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error('操作失败') }
}

const handleCancel = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定取消升级请求 #${row.id} 吗？`, '取消升级')
    await request.post({ url: `/api/admin/upgrades/${row.id}/cancel`, showSuccessMessage: true })
    fetchData(); fetchStats()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error('操作失败') }
}

const handleMarkPaid = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定标记升级请求 #${row.id} 为已支付吗？`, '标记已付')
    await request.post({ url: `/api/admin/upgrades/${row.id}/mark-paid`, showSuccessMessage: true })
    fetchData(); fetchStats()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error('操作失败') }
}

onMounted(() => { fetchData(); fetchStats() })
</script>

<style scoped lang="scss">
.upgrades-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 16px; .el-form-item { margin-bottom: 0; } }
.stat-section { margin-bottom: 24px; }
.stat-card { text-align: center; padding: 8px 0; }
.stat-value { font-size: 20px; font-weight: 600; color: var(--el-color-primary); margin-bottom: 4px; }
.stat-label { color: var(--el-text-color-secondary); font-size: 14px; }
.amount { color: var(--el-color-danger); font-weight: 600; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
