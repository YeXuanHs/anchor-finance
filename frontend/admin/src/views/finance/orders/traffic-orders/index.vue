<template>
  <div class="traffic-orders-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>流量包订单</span>
          <el-button type="primary" @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="订单号">
          <el-input v-model="searchForm.order_no" placeholder="请输入订单号" clearable />
        </el-form-item>
        <el-form-item label="客户">
          <el-input v-model="searchForm.client_name" placeholder="请输入客户名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待付款" :value="0" />
            <el-option label="已付款" :value="1" />
            <el-option label="已生效" :value="2" />
            <el-option label="已过期" :value="3" />
            <el-option label="已取消" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
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
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="order_no" label="订单号" width="170">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">{{ row.order_no }}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="client_name" label="客户" width="120" />
        <el-table-column prop="traffic_package" label="流量包" min-width="150" />
        <el-table-column prop="traffic_amount" label="流量额度" width="120" align="center">
          <template #default="{ row }">
            {{ formatTraffic(row.traffic_amount) }}
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="110" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="expire_at" label="到期时间" width="170" />
        <el-table-column prop="created_at" label="下单时间" width="170" />
        <el-table-column label="操作" width="150" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
            <el-popconfirm v-if="row.status === 1" title="确定要续期该流量包吗？" @confirm="handleRenew(row)">
              <template #reference>
                <el-button type="success" link>续期</el-button>
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
    <el-dialog v-model="detailVisible" title="订单详情" width="650px" destroy-on-close>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="订单号">{{ detailData.order_no }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(detailData.status)" size="small">
            {{ getStatusText(detailData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="客户">{{ detailData.client_name }}</el-descriptions-item>
        <el-descriptions-item label="流量包">{{ detailData.traffic_package }}</el-descriptions-item>
        <el-descriptions-item label="流量额度">{{ formatTraffic(detailData.traffic_amount) }}</el-descriptions-item>
        <el-descriptions-item label="已使用">{{ formatTraffic(detailData.used_amount) }}</el-descriptions-item>
        <el-descriptions-item label="金额">¥{{ formatAmount(detailData.amount) }}</el-descriptions-item>
        <el-descriptions-item label="支付方式">{{ detailData.pay_method || '-' }}</el-descriptions-item>
        <el-descriptions-item label="下单时间">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="到期时间">{{ detailData.expire_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)

const searchForm = reactive({
  order_no: '',
  client_name: '',
  status: undefined as number | undefined,
  date_range: [] as string[]
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])
const detailVisible = ref(false)
const detailData = ref<any>({})

const STATUS_MAP: Record<number, { text: string; type: string }> = {
  0: { text: '待付款', type: 'warning' },
  1: { text: '已付款', type: 'primary' },
  2: { text: '已生效', type: 'success' },
  3: { text: '已过期', type: 'info' },
  4: { text: '已取消', type: 'info' }
}

const getStatusText = (status: number) => STATUS_MAP[status]?.text || '未知'
const getStatusType = (status: number) => (STATUS_MAP[status]?.type || 'info') as any

const formatAmount = (amount: number | undefined) =>
  amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'

const formatTraffic = (mb: number | undefined) => {
  if (!mb) return '0 MB'
  if (mb >= 1024) return `${(mb / 1024).toFixed(1)} GB`
  return `${mb} MB`
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (searchForm.order_no) params.order_no = searchForm.order_no
    if (searchForm.client_name) params.client_name = searchForm.client_name
    if (searchForm.status !== undefined) params.status = searchForm.status
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({ url: '/api/admin/traffic-orders', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取数据失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => {
  searchForm.order_no = ''
  searchForm.client_name = ''
  searchForm.status = undefined
  searchForm.date_range = []
  handleSearch()
}

const handleViewDetail = (row: any) => {
  detailData.value = { ...row }
  detailVisible.value = true
}

const handleRenew = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/traffic-orders/${row.id}/renew` })
    ElMessage.success('续期成功')
    fetchData()
  } catch (error) {
    ElMessage.error('续期失败')
  }
}

const handleExport = () => { ElMessage.info('导出功能开发中...') }

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.traffic-orders-page {
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
  color: #f56c6c;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
