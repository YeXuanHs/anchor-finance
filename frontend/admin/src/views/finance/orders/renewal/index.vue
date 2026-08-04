<template>
  <div class="renewal-orders-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>续费订单管理</span>
          <div class="header-actions">
            <el-button type="primary" @click="handleBatchRenewal">批量续费</el-button>
          </div>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="订单号/客户名/产品名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待续费" :value="0" />
            <el-option label="已续费" :value="1" />
            <el-option label="已过期" :value="2" />
            <el-option label="已取消" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="到期范围">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="order_no" label="原订单号" width="170" />
        <el-table-column prop="client_name" label="客户" width="120" />
        <el-table-column prop="product_name" label="产品/服务" min-width="160" show-overflow-tooltip />
        <el-table-column prop="amount" label="续费金额" width="120" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="billing_cycle" label="周期" width="100" />
        <el-table-column prop="due_date" label="到期时间" width="170" />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 0"
              type="success"
              link
              @click="handleRenew(row)"
            >
              续费
            </el-button>
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
            <el-popconfirm
              v-if="row.status === 0"
              title="确定取消该续费吗？"
              @confirm="handleCancel(row)"
            >
              <template #reference>
                <el-button type="danger" link>取消</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
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

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="续费订单详情" width="650px">
      <el-descriptions :column="2" border v-if="detailData">
        <el-descriptions-item label="原订单号">{{ detailData.order_no }}</el-descriptions-item>
        <el-descriptions-item label="客户">{{ detailData.client_name }}</el-descriptions-item>
        <el-descriptions-item label="产品" :span="2">{{ detailData.product_name }}</el-descriptions-item>
        <el-descriptions-item label="续费金额">
          <span class="amount-text">¥{{ formatAmount(detailData.amount) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="计费周期">{{ detailData.billing_cycle }}</el-descriptions-item>
        <el-descriptions-item label="到期时间">{{ detailData.due_date }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(detailData.status)" size="small">
            {{ getStatusText(detailData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ detailData.remark || '无' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const detailVisible = ref(false)
const detailData = ref<any>(null)
const selectedRows = ref<any[]>([])

const searchForm = reactive({
  keyword: '',
  status: undefined as number | undefined,
  date_range: [] as string[]
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const statusMap: Record<number, { text: string; type: string }> = {
  0: { text: '待续费', type: 'warning' },
  1: { text: '已续费', type: 'success' },
  2: { text: '已过期', type: 'danger' },
  3: { text: '已取消', type: 'info' }
}

const getStatusText = (status: number) => statusMap[status]?.text || '未知'
const getStatusType = (status: number) => (statusMap[status]?.type || 'info') as any

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      keyword: searchForm.keyword || undefined,
      status: searchForm.status
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({
      url: '/api/admin/renewal/list',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取续费订单失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = undefined
  searchForm.date_range = []
  handleSearch()
}

const handleSelectionChange = (rows: any[]) => {
  selectedRows.value = rows
}

const handleViewDetail = (row: any) => {
  detailData.value = row
  detailVisible.value = true
}

const handleRenew = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定续费订单 ${row.order_no} 吗？金额: ¥${formatAmount(row.amount)}`, '续费确认')
    await request.post({
      url: `/api/admin/orders/${row.id}/renew`
    })
    ElMessage.success('续费成功')
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('续费失败')
    }
  }
}

const handleCancel = async (row: any) => {
  try {
    await request.post({
      url: `/api/admin/orders/${row.id}/cancel-renewal`
    })
    ElMessage.success('已取消')
    fetchData()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleBatchRenewal = async () => {
  if (!selectedRows.value.length) {
    ElMessage.warning('请先选择要续费的订单')
    return
  }
  try {
    await ElMessageBox.confirm(`确定批量续费 ${selectedRows.value.length} 个订单吗？`, '批量续费确认')
    const ids = selectedRows.value.map(r => r.id)
    await request.post({
      url: '/api/admin/renewal/multi-renew',
      params: { ids }
    })
    ElMessage.success('批量续费成功')
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('批量续费失败')
    }
  }
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchData()
}

const handlePageChange = () => {
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.renewal-orders-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.amount-text {
  font-weight: 600;
  color: var(--el-color-primary);
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
