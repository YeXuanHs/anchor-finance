<template>
  <div class="accounts-record">
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
      <el-select v-model="filterType" placeholder="变动类型" clearable style="width: 140px">
        <el-option label="全部" value="" />
        <el-option label="充值" value="recharge" />
        <el-option label="消费" value="expense" />
        <el-option label="退款" value="refund" />
        <el-option label="提现" value="withdraw" />
        <el-option label="转入" value="transfer_in" />
        <el-option label="转出" value="transfer_out" />
      </el-select>
      <el-input
        v-model="searchKeyword"
        placeholder="搜索备注..."
        clearable
        style="width: 200px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-button type="primary" @click="handleSearch">
        <el-icon><Search /></el-icon>
        搜索
      </el-button>
      <el-button @click="handleReset">重置</el-button>
    </div>

    <!-- 数据表格 -->
    <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
      <el-table-column prop="id" label="流水号" width="160" show-overflow-tooltip />
      <el-table-column prop="created_at" label="时间" width="180" />
      <el-table-column prop="type" label="类型" width="120">
        <template #default="{ row }">
          <el-tag :type="typeTagType(row.type)" size="small">
            {{ typeLabel(row.type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="金额" width="140">
        <template #default="{ row }">
          <span :class="['amount-text', row.amount > 0 ? 'income' : 'expense']">
            {{ row.amount > 0 ? '+' : '' }}¥{{ Math.abs(row.amount).toFixed(2) }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="before_balance" label="变动前余额" width="140">
        <template #default="{ row }">
          <span>¥{{ row.before_balance.toFixed(2) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="after_balance" label="变动后余额" width="140">
        <template #default="{ row }">
          <span>¥{{ row.after_balance.toFixed(2) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="related_order" label="关联订单" width="160" show-overflow-tooltip />
      <el-table-column prop="remark" label="备注" min-width="160" show-overflow-tooltip />
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
import { Search } from '@element-plus/icons-vue'
import request from '@/utils/request'

interface AccountsItem {
  id: string
  created_at: string
  type: string
  amount: number
  before_balance: number
  after_balance: number
  related_order: string
  remark: string
}

const loading = ref(false)
const dateRange = ref<string[]>([])
const filterType = ref('')
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const tableData = ref<AccountsItem[]>([])

const typeLabel = (type: string) => {
  const map: Record<string, string> = {
    recharge: '充值',
    expense: '消费',
    refund: '退款',
    withdraw: '提现',
    transfer_in: '转入',
    transfer_out: '转出'
  }
  return map[type] || type
}

const typeTagType = (type: string) => {
  const map: Record<string, string> = {
    recharge: 'success',
    expense: 'danger',
    refund: 'info',
    withdraw: 'warning',
    transfer_in: 'success',
    transfer_out: 'danger'
  }
  return map[type] || 'info'
}

const handleSearch = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/balance/logs', { params: { page: currentPage.value, page_size: pageSize.value } })
    tableData.value = res.data.data.list
    total.value = res.data.data.total
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  dateRange.value = []
  filterType.value = ''
  searchKeyword.value = ''
  currentPage.value = 1
  handleSearch()
}

onMounted(() => {
  handleSearch()
})

defineExpose({ handleSearch })
</script>

<style scoped lang="scss">
.accounts-record {
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

    &.expense {
      color: #ff4d4f;
    }
  }

  .el-pagination {
    margin-top: 20px;
    justify-content: flex-end;
  }
}
</style>
