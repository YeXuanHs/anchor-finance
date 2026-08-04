<template>
  <div class="sales-statistics-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>销售统计</span>
          <el-button type="primary" @click="handleExport">
            <el-icon><Download /></el-icon>
            导出报表
          </el-button>
        </div>
      </template>

      <!-- 统计卡片 -->
      <el-row :gutter="20" class="stat-cards">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-item">
              <div class="stat-label">今日销售额</div>
              <div class="stat-value">¥{{ formatAmount(stats.today_amount) }}</div>
              <div class="stat-change" :class="{ up: stats.today_change >= 0, down: stats.today_change < 0 }">
                {{ stats.today_change >= 0 ? '+' : '' }}{{ stats.today_change }}%
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-item">
              <div class="stat-label">本月销售额</div>
              <div class="stat-value">¥{{ formatAmount(stats.month_amount) }}</div>
              <div class="stat-change" :class="{ up: stats.month_change >= 0, down: stats.month_change < 0 }">
                {{ stats.month_change >= 0 ? '+' : '' }}{{ stats.month_change }}%
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-item">
              <div class="stat-label">今日订单数</div>
              <div class="stat-value">{{ stats.today_orders || 0 }}</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-item">
              <div class="stat-label">本月订单数</div>
              <div class="stat-value">{{ stats.month_orders || 0 }}</div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
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
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="date" label="日期" width="120" align="center" />
        <el-table-column prop="order_count" label="订单数" width="100" align="center" />
        <el-table-column prop="amount" label="销售额" width="140" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="new_clients" label="新客户" width="100" align="center" />
        <el-table-column prop="refund_amount" label="退款额" width="140" align="right">
          <template #default="{ row }">
            <span class="refund-text">¥{{ formatAmount(row.refund_amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="net_amount" label="净收入" width="140" align="right">
          <template #default="{ row }">
            <span class="net-text">¥{{ formatAmount(row.net_amount) }}</span>
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
import request from '@/utils/http'
import { exportToCSV } from '@/utils/export'

const loading = ref(false)

const searchForm = reactive({
  date_range: [] as string[]
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])
const stats = ref({
  today_amount: 0,
  today_change: 0,
  month_amount: 0,
  month_change: 0,
  today_orders: 0,
  month_orders: 0
})

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({ url: '/api/admin/sales/statistics', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
    if (data.stats) stats.value = data.stats
  } catch (error) {
    console.error('获取数据失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.date_range = []; handleSearch() }

const handleExport = async () => {
  try {
    const data = await request.get({ url: '/api/admin/statistics/sales', params: { page: 1, page_size: 9999 } })
    const list = data.list || data || []
    exportToCSV(list, [
      { key: 'client_name', title: '客户' },
      { key: 'order_count', title: '订单数' },
      { key: 'total_amount', title: '总金额' },
      { key: 'last_order_at', title: '最近下单' }
    ], '销售统计')
    ElMessage.success('导出成功')
  } catch { ElMessage.error('导出失败') }
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.sales-statistics-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-cards {
  margin-bottom: 24px;

  .stat-card {
    .stat-item {
      text-align: center;

      .stat-label {
        font-size: 14px;
        color: #909399;
        margin-bottom: 8px;
      }

      .stat-value {
        font-size: 24px;
        font-weight: 700;
        color: #303133;
        margin-bottom: 4px;
      }

      .stat-change {
        font-size: 13px;

        &.up {
          color: #67c23a;
        }

        &.down {
          color: #f56c6c;
        }
      }
    }
  }
}

.search-form {
  margin-bottom: 20px;
}

.amount-text {
  font-weight: 600;
  color: #f56c6c;
}

.refund-text {
  color: #e6a23c;
}

.net-text {
  font-weight: 600;
  color: #67c23a;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
