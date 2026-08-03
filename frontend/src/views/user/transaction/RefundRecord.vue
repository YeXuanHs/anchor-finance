<template>
  <div class="refund-record">
    <!-- 搜索筛选 -->
    <div class="filter-bar">
      <el-date-picker
        v-model="dateRange"
        type="daterange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        value-format="YYYY-MM-DD"
        style="width: 300px"
      />
      <el-select v-model="filterStatus" placeholder="状态" clearable style="width: 140px">
        <el-option label="全部" value="" />
        <el-option label="成功" value="success" />
        <el-option label="审核中" value="reviewing" />
        <el-option label="已拒绝" value="rejected" />
      </el-select>
      <el-select v-model="filterType" placeholder="退款类型" clearable style="width: 140px">
        <el-option label="全部" value="" />
        <el-option label="订单退款" value="order" />
        <el-option label="充值退款" value="recharge" />
        <el-option label="其他" value="other" />
      </el-select>
      <el-button type="primary" @click="handleSearch">
        <el-icon><Search /></el-icon>
        搜索
      </el-button>
      <el-button @click="handleReset">重置</el-button>
    </div>

    <!-- 数据表格 -->
    <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
      <el-table-column prop="id" label="退款单号" width="160" show-overflow-tooltip />
      <el-table-column prop="order_id" label="原订单号" width="160" show-overflow-tooltip />
      <el-table-column prop="created_at" label="退款时间" width="180" />
      <el-table-column prop="amount" label="退款金额" width="140">
        <template #default="{ row }">
          <span class="amount-text refund">+¥{{ row.amount.toFixed(2) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="type" label="退款类型" width="120">
        <template #default="{ row }">
          <el-tag size="small">{{ typeLabel(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)" size="small">
            {{ statusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="reason" label="退款原因" min-width="160" show-overflow-tooltip />
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="handleDetail(row)">详情</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="[10, 20, 50, 100]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSearch"
      @current-change="handleSearch"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()

interface RefundItem {
  id: string
  order_id: string
  created_at: string
  amount: number
  type: string
  status: string
  reason: string
}

const loading = ref(false)
const dateRange = ref<string[]>([])
const filterStatus = ref('')
const filterType = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const tableData = ref<RefundItem[]>([])

const typeLabel = (type: string) => {
  const map: Record<string, string> = { order: '订单退款', recharge: '充值退款', other: '其他' }
  return map[type] || type
}

const statusLabel = (status: string) => {
  const map: Record<string, string> = { success: '成功', reviewing: '审核中', rejected: '已拒绝' }
  return map[status] || status
}

const statusTagType = (status: string) => {
  const map: Record<string, string> = { success: 'success', reviewing: 'warning', rejected: 'danger' }
  return map[status] || 'info'
}

const handleSearch = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v2/balance/logs', { params: { page: currentPage.value, page_size: pageSize.value, type: 'refund' } })
    tableData.value = res.data.data.list
    total.value = res.data.data.total
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  dateRange.value = []
  filterStatus.value = ''
  filterType.value = ''
  currentPage.value = 1
  handleSearch()
}

const handleDetail = (row: RefundItem) => {
  router.push(`/user/orders/${row.order_id}`)
}

onMounted(() => {
  handleSearch()
})

defineExpose({ handleSearch })
</script>

<style scoped lang="scss">
.refund-record {
  .filter-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
    flex-wrap: wrap;
  }

  .amount-text {
    font-weight: 600;

    &.refund {
      color: #52c41a;
    }
  }

  .el-pagination {
    margin-top: 20px;
    justify-content: flex-end;
  }
}
</style>
