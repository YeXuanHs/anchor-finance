<template>
  <div class="traffic-orders-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('trafficOrders.title') }}</span>
          <el-button type="primary" @click="handleExport">
            <el-icon><Download /></el-icon>
            {{ $t('common.export') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('trafficOrders.orderNo')">
          <el-input v-model="searchForm.order_no" :placeholder="$t('trafficOrders.inputOrderNo')" clearable />
        </el-form-item>
        <el-form-item :label="$t('trafficOrders.client')">
          <el-input v-model="searchForm.client_name" :placeholder="$t('trafficOrders.inputClientName')" clearable />
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('trafficOrders.pendingPayment')" :value="0" />
            <el-option :label="$t('common.paid')" :value="1" />
            <el-option :label="$t('trafficOrders.active')" :value="2" />
            <el-option :label="$t('trafficOrders.expired')" :value="3" />
            <el-option :label="$t('common.cancelled')" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('trafficOrders.dateRange')">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            :range-separator="$t('common.to')"
            :start-placeholder="$t('common.startDate')"
            :end-placeholder="$t('common.endDate')"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="order_no" :label="$t('trafficOrders.orderNo')" width="170">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">{{ row.order_no }}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="client_name" :label="$t('trafficOrders.client')" width="120" />
        <el-table-column prop="traffic_package" :label="$t('trafficOrders.trafficPackage')" min-width="150" />
        <el-table-column prop="traffic_amount" :label="$t('trafficOrders.trafficAmount')" width="120" align="center">
          <template #default="{ row }">
            {{ formatTraffic(row.traffic_amount) }}
          </template>
        </el-table-column>
        <el-table-column prop="amount" :label="$t('common.amount')" width="110" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('common.status')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="expire_at" :label="$t('trafficOrders.expireAt')" width="170" />
        <el-table-column prop="created_at" :label="$t('trafficOrders.orderTime')" width="170" />
        <el-table-column :label="$t('common.action')" width="150" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">{{ $t('common.detail') }}</el-button>
            <el-popconfirm v-if="row.status === 1" :title="$t('trafficOrders.confirmRenew')" @confirm="handleRenew(row)">
              <template #reference>
                <el-button type="success" link>{{ $t('trafficOrders.renew') }}</el-button>
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
    <el-dialog v-model="detailVisible" :title="$t('trafficOrders.detailTitle')" width="650px" destroy-on-close>
      <el-descriptions :column="2" border>
        <el-descriptions-item :label="$t('trafficOrders.orderNo')">{{ detailData.order_no }}</el-descriptions-item>
        <el-descriptions-item :label="$t('common.status')">
          <el-tag :type="getStatusType(detailData.status)" size="small">
            {{ getStatusText(detailData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('trafficOrders.client')">{{ detailData.client_name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('trafficOrders.trafficPackage')">{{ detailData.traffic_package }}</el-descriptions-item>
        <el-descriptions-item :label="$t('trafficOrders.trafficAmount')">{{ formatTraffic(detailData.traffic_amount) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('trafficOrders.used')">{{ formatTraffic(detailData.used_amount) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('common.amount')">¥{{ formatAmount(detailData.amount) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('trafficOrders.payMethod')">{{ detailData.pay_method || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('trafficOrders.orderTime')">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('trafficOrders.expireAt')">{{ detailData.expire_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">{{ $t('common.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'
import { exportToCSV } from '@/utils/export'

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

const STATUS_MAP: Record<number, { text: () => string; type: string }> = {
  0: { text: () => $t('trafficOrders.pendingPayment'), type: 'warning' },
  1: { text: () => $t('common.paid'), type: 'primary' },
  2: { text: () => $t('trafficOrders.active'), type: 'success' },
  3: { text: () => $t('trafficOrders.expired'), type: 'info' },
  4: { text: () => $t('common.cancelled'), type: 'info' }
}

const getStatusText = (status: number) => STATUS_MAP[status]?.text() || $t('common.unknown')
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
    console.error('fetch data failed:', error)
    ElMessage.error($t('trafficOrders.fetchFailed'))
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
    ElMessage.success($t('trafficOrders.renewSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('trafficOrders.renewFailed'))
  }
}

const handleExport = async () => {
  try {
    const data = await request.get({ url: '/api/admin/orders', params: { page: 1, page_size: 9999, type: 'traffic' } })
    const list = data.list || data || []
    exportToCSV(list, [
      { key: 'id', title: 'ID' },
      { key: 'order_no', title: $t('trafficOrders.orderNo') },
      { key: 'product_name', title: $t('trafficOrders.exportProduct') },
      { key: 'amount', title: $t('common.amount') },
      { key: 'status', title: $t('common.status') },
      { key: 'created_at', title: $t('common.createdAt') }
    ], $t('trafficOrders.exportFileName'))
    ElMessage.success($t('trafficOrders.exportSuccess'))
  } catch { ElMessage.error($t('trafficOrders.exportFailed')) }
}

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
