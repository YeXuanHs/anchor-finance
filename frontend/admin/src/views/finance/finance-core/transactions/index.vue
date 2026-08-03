<template>
  <div class="transactions-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>交易流水</span>
          <el-button type="primary" @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="交易号/用户名" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="充值" value="deposit" />
            <el-option label="支付" value="payment" />
            <el-option label="退款" value="refund" />
            <el-option label="扣款" value="deduction" />
            <el-option label="佣金" value="commission" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="成功" :value="1" />
            <el-option label="失败" :value="0" />
            <el-option label="处理中" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="网关">
          <el-select v-model="searchForm.gateway" placeholder="全部" clearable>
            <el-option v-for="gw in gateways" :key="gw" :label="gw" :value="gw" />
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

      <!-- 汇总信息 -->
      <div class="summary-cards">
        <el-row :gutter="20">
          <el-col :span="6">
            <div class="summary-item">
              <div class="summary-label">总收入</div>
              <div class="summary-value text-success">¥{{ formatAmount(summary.total_income) }}</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="summary-item">
              <div class="summary-label">总支出</div>
              <div class="summary-value text-danger">¥{{ formatAmount(summary.total_expense) }}</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="summary-item">
              <div class="summary-label">总退款</div>
              <div class="summary-value text-warning">¥{{ formatAmount(summary.total_refund) }}</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="summary-item">
              <div class="summary-label">交易笔数</div>
              <div class="summary-value">{{ summary.total_count }}</div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="transaction_no" label="交易号" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">
              {{ row.transaction_no }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="client_username" label="用户" width="120">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewClient(row)">
              {{ row.client_username }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)" size="small">
              {{ getTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="130" align="right">
          <template #default="{ row }">
            <span :class="getAmountClass(row.type, row.amount)">
              {{ row.amount >= 0 ? '+' : '' }}¥{{ formatAmount(row.amount) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="balance" label="余额" width="120" align="right">
          <template #default="{ row }">
            ¥{{ formatAmount(row.balance) }}
          </template>
        </el-table-column>
        <el-table-column prop="gateway" label="网关" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="180" show-overflow-tooltip />
        <el-table-column prop="created_at" label="时间" width="180" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
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
    <el-dialog v-model="detailVisible" title="交易详情" width="600px">
      <el-descriptions :column="2" border v-if="currentTransaction">
        <el-descriptions-item label="交易ID">{{ currentTransaction.id }}</el-descriptions-item>
        <el-descriptions-item label="交易号">{{ currentTransaction.transaction_no }}</el-descriptions-item>
        <el-descriptions-item label="用户">{{ currentTransaction.client_username }}</el-descriptions-item>
        <el-descriptions-item label="类型">
          <el-tag :type="getTypeTagType(currentTransaction.type)" size="small">
            {{ getTypeText(currentTransaction.type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="金额">
          <span :class="getAmountClass(currentTransaction.type, currentTransaction.amount)" style="font-size: 18px; font-weight: 600;">
            {{ currentTransaction.amount >= 0 ? '+' : '' }}¥{{ formatAmount(currentTransaction.amount) }}
          </span>
        </el-descriptions-item>
        <el-descriptions-item label="余额">¥{{ formatAmount(currentTransaction.balance) }}</el-descriptions-item>
        <el-descriptions-item label="网关">{{ currentTransaction.gateway }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentTransaction.status)" size="small">
            {{ getStatusText(currentTransaction.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentTransaction.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="关联订单">{{ currentTransaction.order_no || '-' }}</el-descriptions-item>
        <el-descriptions-item label="关联账单">{{ currentTransaction.bill_no || '-' }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ currentTransaction.ip || '-' }}</el-descriptions-item>
        <el-descriptions-item label="时间">{{ currentTransaction.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'
import { exportToCSV } from '@/utils/export'

const router = useRouter()

const loading = ref(false)
const detailVisible = ref(false)
const currentTransaction = ref<any>({})

const gateways = ref<string[]>([])

const searchForm = reactive({
  keyword: '',
  type: '',
  status: undefined as number | undefined,
  gateway: '',
  date_range: [] as string[]
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])

const summary = reactive({
  total_income: 0,
  total_expense: 0,
  total_refund: 0,
  total_count: 0
})

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const getTypeTagType = (type: string) => {
  const map: Record<string, string> = {
    deposit: 'success',
    payment: 'warning',
    refund: 'info',
    deduction: 'danger',
    commission: ''
  }
  return (map[type] || 'info') as any
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    deposit: '充值',
    payment: '支付',
    refund: '退款',
    deduction: '扣款',
    commission: '佣金'
  }
  return map[type] || type
}

const getAmountClass = (type: string, amount: number) => {
  if (type === 'deposit' || type === 'commission') return 'text-success'
  if (type === 'payment' || type === 'deduction') return 'text-danger'
  if (type === 'refund') return 'text-warning'
  return ''
}

const getStatusType = (status: number) => {
  const map: Record<number, string> = {
    0: 'danger',
    1: 'success',
    2: 'warning'
  }
  return (map[status] || 'info') as any
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = {
    0: '失败',
    1: '成功',
    2: '处理中'
  }
  return map[status] || '未知'
}

const fetchTransactions = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      keyword: searchForm.keyword || undefined,
      type: searchForm.type || undefined,
      status: searchForm.status,
      gateway: searchForm.gateway || undefined
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({
      url: '/api/admin/accounts',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
    summary.total_income = data.total_income || 0
    summary.total_expense = data.total_expense || 0
    summary.total_refund = data.total_refund || 0
    summary.total_count = data.total_count || 0
  } catch (error) {
    console.error('获取交易流水失败:', error)
    ElMessage.error('获取交易流水失败')
  } finally {
    loading.value = false
  }
}

const fetchGateways = async () => {
  try {
    const data = await request.get({ url: '/api/admin/payment-gateways' })
    gateways.value = data || []
  } catch (error) {
    console.error('获取网关列表失败:', error)
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchTransactions()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.type = ''
  searchForm.status = undefined
  searchForm.gateway = ''
  searchForm.date_range = []
  handleSearch()
}

const handleViewDetail = (row: any) => {
  currentTransaction.value = { ...row }
  detailVisible.value = true
}

const handleViewClient = (row: any) => {
  router.push(`/finance/clients/detail/${row.client_id}`)
}

const handleExport = async () => {
  try {
    const data = await request.get({ url: '/api/admin/accounts', params: { page: 1, page_size: 9999 } })
    const list = data.list || data || []
    exportToCSV(list, [
      { key: 'id', title: 'ID' },
      { key: 'user_id', title: '用户ID' },
      { key: 'amount_in', title: '收入' },
      { key: 'amount_out', title: '支出' },
      { key: 'currency', title: '货币' },
      { key: 'description', title: '描述' },
      { key: 'created_at', title: '时间' }
    ], '交易流水')
    ElMessage.success('导出成功')
  } catch { ElMessage.error('导出失败') }
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchTransactions()
}

const handlePageChange = () => {
  fetchTransactions()
}

onMounted(() => {
  fetchTransactions()
  fetchGateways()
})
</script>

<style scoped lang="scss">
.transactions-page {
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

.summary-cards {
  margin-bottom: 20px;
  padding: 16px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;

  .summary-item {
    text-align: center;

    .summary-label {
      color: var(--el-text-color-secondary);
      font-size: 14px;
      margin-bottom: 4px;
    }

    .summary-value {
      font-size: 20px;
      font-weight: 600;
      color: var(--el-text-color-primary);
    }
  }
}

.text-success {
  color: var(--el-color-success);
}

.text-danger {
  color: var(--el-color-danger);
}

.text-warning {
  color: var(--el-color-warning);
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>