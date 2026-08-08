<template>
  <div class="order-list-page">
    <!-- 标签页切换 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="全部订单" name="all" />
      <el-tab-pane label="待付款" name="pending_payment" />
      <el-tab-pane label="待开通" name="pending_activation" />
      <el-tab-pane label="进行中" name="active" />
      <el-tab-pane label="已完成" name="completed" />
      <el-tab-pane label="已取消" name="cancelled" />
      <el-tab-pane label="已退款" name="refunded" />
    </el-tabs>

    <!-- 搜索筛选区域 -->
    <el-card shadow="never" class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="订单号/客户名/产品名"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="订单类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable style="width: 120px">
            <el-option label="新购" value="new" />
            <el-option label="续费" value="renewal" />
            <el-option label="升级" value="upgrade" />
            <el-option label="退款" value="refund" />
          </el-select>
        </el-form-item>
        <el-form-item label="支付方式">
          <el-select v-model="searchForm.payment_method" placeholder="全部" clearable style="width: 120px">
            <el-option v-for="method in paymentMethods" :key="method.id" :label="method.name" :value="method.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="下单时间">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作栏 -->
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <div class="action-left">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加订单
          </el-button>
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
        <div class="action-right">
          <el-button circle @click="fetchList">
            <el-icon><Refresh /></el-icon>
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        @sort-change="handleSortChange"
      >
        <el-table-column prop="id" label="ID" width="80" sortable="custom" align="center" />
        <el-table-column prop="order_no" label="订单号" width="150">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">
              {{ row.order_no }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="client_name" label="客户" min-width="120">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewClient(row)">
              {{ row.client_name }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="product_name" label="产品" min-width="150" />
        <el-table-column prop="type" label="类型" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)" size="small">
              {{ getTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="100" align="right">
          <template #default="{ row }">
            ¥{{ formatMoney(row.amount) }}
          </template>
        </el-table-column>
        <el-table-column prop="payment_method" label="支付方式" width="100" />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="下单时间" width="170" sortable="custom" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleView(row)">
              查看
            </el-button>
            <el-button
              v-if="row.status === 'pending_payment'"
              type="success"
              link
              size="small"
              @click="handleConfirmPayment(row)"
            >
              确认付款
            </el-button>
            <el-button
              v-if="row.status === 'pending_activation'"
              type="warning"
              link
              size="small"
              @click="handleActivate(row)"
            >
              开通
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, Download } from '@element-plus/icons-vue'
import request from '@/utils/http'

const router = useRouter()
const loading = ref(false)
const tableData = ref([])
const activeTab = ref('all')
const paymentMethods = ref<{ id: number; name: string }[]>([])

// 搜索表单
const searchForm = reactive({
  keyword: '',
  type: '',
  payment_method: '',
  date_range: null as [Date, Date] | null
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 排序
const sortParams = reactive({
  sort: 'created_at',
  order: 'desc'
})

// 格式化金额
const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// 类型标签
const getTypeTagType = (type: string) => {
  const map: Record<string, string> = {
    new: 'primary',
    renewal: 'success',
    upgrade: 'warning',
    refund: 'danger'
  }
  return map[type] || 'info'
}

// 类型文本
const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    new: '新购',
    renewal: '续费',
    upgrade: '升级',
    refund: '退款'
  }
  return map[type] || '未知'
}

// 状态类型
const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending_payment: 'warning',
    pending_activation: 'primary',
    active: 'success',
    completed: 'success',
    cancelled: 'info',
    refunded: 'danger'
  }
  return map[status] || 'info'
}

// 状态文本
const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending_payment: '待付款',
    pending_activation: '待开通',
    active: '进行中',
    completed: '已完成',
    cancelled: '已取消',
    refunded: '已退款'
  }
  return map[status] || '未知'
}

// 标签页切换
const handleTabChange = (tab: string) => {
  searchForm.type = ''
  pagination.page = 1
  fetchList()
}

// 获取列表数据
const fetchList = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      sort: sortParams.sort,
      order: sortParams.order
    }

    if (activeTab.value !== 'all') params.status = activeTab.value
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.type) params.type = searchForm.type
    if (searchForm.payment_method) params.payment_method = searchForm.payment_method
    if (searchForm.date_range) {
      params.start_date = searchForm.date_range[0].toISOString().split('T')[0]
      params.end_date = searchForm.date_range[1].toISOString().split('T')[0]
    }

    const data = await request.get({ url: '/api/admin/orders', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取订单列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.type = ''
  searchForm.payment_method = ''
  searchForm.date_range = null
  pagination.page = 1
  fetchList()
}

// 排序变化
const handleSortChange = ({ prop, order }: any) => {
  sortParams.sort = prop || 'created_at'
  sortParams.order = order === 'ascending' ? 'asc' : 'desc'
  fetchList()
}

// 分页大小变化
const handleSizeChange = (size: number) => {
  pagination.page_size = size
  pagination.page = 1
  fetchList()
}

// 页码变化
const handlePageChange = (page: number) => {
  pagination.page = page
  fetchList()
}

// 添加订单
const handleAdd = () => {
  router.push('/add-order')
}

// 查看订单
const handleView = (row: any) => {
  router.push(`/order-detail/${row.id}`)
}

// 查看客户
const handleViewClient = (row: any) => {
  router.push(`/customer-view/${row.client_id}`)
}

// 确认付款
const handleConfirmPayment = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要确认订单 "${row.order_no}" 已付款吗？`, '确认付款', {
      type: 'warning'
    })
    await request.post({ url: `/api/admin/orders/${row.id}/confirm-payment` })
    ElMessage.success('确认付款成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('确认付款失败:', error)
    }
  }
}

// 开通订单
const handleActivate = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要开通订单 "${row.order_no}" 吗？`, '确认开通', {
      type: 'warning'
    })
    await request.post({ url: `/api/admin/orders/${row.id}/activate` })
    ElMessage.success('开通成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('开通失败:', error)
    }
  }
}

// 导出
const handleExport = async () => {
  try {
    const params: any = {}
    if (activeTab.value !== 'all') params.status = activeTab.value
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.type) params.type = searchForm.type

    const response = await request.get({
      url: '/api/admin/orders/export',
      params,
      responseType: 'blob'
    })

    const url = window.URL.createObjectURL(new Blob([response]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `订单列表_${new Date().toISOString().split('T')[0]}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('导出失败:', error)
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.order-list-page {
  padding: 16px;
}

.search-card {
  margin-bottom: 16px;

  :deep(.el-card__body) {
    padding-bottom: 0;
  }
}

.action-card {
  margin-bottom: 16px;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.action-left {
  display: flex;
  gap: 8px;
}

.table-card {
  :deep(.el-card__body) {
    padding: 0;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px;
}
</style>
