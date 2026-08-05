<template>
  <div class="recharge-record">
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
        <el-option label="待处理" value="pending" />
        <el-option label="失败" value="failed" />
      </el-select>
      <el-select v-model="filterPayMethod" placeholder="支付方式" clearable style="width: 140px">
        <el-option label="全部" value="" />
        <el-option label="支付宝" value="alipay" />
        <el-option label="微信" value="wechat" />
        <el-option label="银行卡" value="bank" />
      </el-select>
      <el-button type="primary" @click="handleSearch">
        <el-icon><Search /></el-icon>
        搜索
      </el-button>
      <el-button @click="handleReset">重置</el-button>
    </div>

    <!-- 数据表格 -->
    <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
      <el-table-column prop="id" label="订单号" width="160" show-overflow-tooltip />
      <el-table-column prop="created_at" label="充值时间" width="180" />
      <el-table-column prop="amount" label="充值金额" width="140">
        <template #default="{ row }">
          <span class="amount-text income">+¥{{ row.amount.toFixed(2) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="pay_method" label="支付方式" width="120">
        <template #default="{ row }">
          <el-tag :type="payMethodTagType(row.pay_method)" size="small">
            {{ payMethodLabel(row.pay_method) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)" size="small">
            {{ statusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" min-width="160" show-overflow-tooltip />
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

interface RechargeItem {
  id: string
  created_at: string
  amount: number
  pay_method: string
  status: string
  remark: string
}

const loading = ref(false)
const dateRange = ref<string[]>([])
const filterStatus = ref('')
const filterPayMethod = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const tableData = ref<RechargeItem[]>([])

const payMethodLabel = (method: string) => {
  const map: Record<string, string> = { alipay: '支付宝', wechat: '微信', bank: '银行卡' }
  return map[method] || method
}

const payMethodTagType = (method: string) => {
  const map: Record<string, string> = { alipay: 'primary', wechat: 'success', bank: 'warning' }
  return map[method] || 'info'
}

const statusLabel = (status: string) => {
  const map: Record<string, string> = { success: '成功', pending: '待处理', failed: '失败' }
  return map[status] || status
}

const statusTagType = (status: string) => {
  const map: Record<string, string> = { success: 'success', pending: 'warning', failed: 'danger' }
  return map[status] || 'info'
}

const handleSearch = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/balance/logs', { params: { page: currentPage.value, page_size: pageSize.value, type: 'recharge' } })
    tableData.value = res.data.data.list
    total.value = res.data.data.total
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  dateRange.value = []
  filterStatus.value = ''
  filterPayMethod.value = ''
  currentPage.value = 1
  handleSearch()
}

const handleDetail = (row: RechargeItem) => {
  router.push(`/user/orders/${row.order_id}`)
}

onMounted(() => {
  handleSearch()
})

defineExpose({ handleSearch })
</script>

<style scoped lang="scss">
.recharge-record {
  .filter-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
    flex-wrap: wrap;
  }

  .amount-text {
    font-weight: 600;

    &.income {
      color: #52c41a;
    }
  }

  .el-pagination {
    margin-top: 20px;
    justify-content: flex-end;
  }
}
</style>
