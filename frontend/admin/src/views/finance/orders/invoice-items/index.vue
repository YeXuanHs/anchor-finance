<template>
  <div class="invoice-items-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>账单明细</span>
          <el-button type="primary" @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="账单ID">
          <el-input v-model="searchForm.invoice_id" placeholder="请输入账单ID" clearable />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="searchForm.description" placeholder="请输入描述关键词" clearable />
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
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border show-summary :summary-method="getSummary">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="invoice_id" label="账单ID" width="100" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewInvoice(row)">
              {{ row.invoice_id }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="amount" label="金额" width="120" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="数量" width="80" align="center" />
        <el-table-column label="小计" width="120" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount((row.amount || 0) * (row.quantity || 1)) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="120" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewInvoice(row)">查看账单</el-button>
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { TableColumnCtx } from 'element-plus'
import request from '@/utils/http'

interface SummaryMethodProps<T = any> {
  columns: TableColumnCtx<T>[]
  data: T[]
}

const loading = ref(false)

const searchForm = reactive({
  invoice_id: '',
  description: '',
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

const getSummary = (param: SummaryMethodProps) => {
  const { columns, data } = param
  const sums: string[] = []
  columns.forEach((column, index) => {
    if (index === 0) {
      sums[index] = '合计'
      return
    }
    if (column.property === 'amount') {
      const total = data.reduce((sum, row) => sum + (Number(row.amount) || 0), 0)
      sums[index] = `¥${formatAmount(total)}`
    } else if (index === columns.length - 3) {
      const total = data.reduce((sum, row) => sum + (Number(row.amount) || 0) * (Number(row.quantity) || 1), 0)
      sums[index] = `¥${formatAmount(total)}`
    } else {
      sums[index] = ''
    }
  })
  return sums
}

const fetchList = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      invoice_id: searchForm.invoice_id || undefined,
      description: searchForm.description || undefined
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({ url: '/api/admin/invoice-items', params })
    tableData.value = data.list || data || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取账单明细失败:', error)
    ElMessage.error('获取账单明细失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

const handleReset = () => {
  searchForm.invoice_id = ''
  searchForm.description = ''
  searchForm.date_range = []
  handleSearch()
}

const handleViewInvoice = (row: any) => {
  ElMessage.info(`查看账单 #${row.invoice_id}`)
}

const handleExport = () => {
  ElMessage.info('导出功能开发中...')
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchList()
}

const handlePageChange = () => {
  fetchList()
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.invoice-items-page {
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
