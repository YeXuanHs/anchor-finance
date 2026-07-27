<template>
  <div class="orders-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="订单号">
          <el-input v-model="searchForm.order_no" placeholder="订单号" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待支付" value="pending" />
            <el-option label="已完成" value="completed" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <div class="art-card">
      <div class="table-header">
        <h3>订单列表</h3>
        <el-button @click="exportOrders">
          <el-icon><Download /></el-icon>
          导出
        </el-button>
      </div>
      
      <el-table :data="orders" style="width: 100%" v-loading="loading">
        <el-table-column prop="order_no" label="订单号" width="180" />
        <el-table-column prop="product" label="产品" />
        <el-table-column prop="user" label="用户" />
        <el-table-column prop="amount" label="金额">
          <template #default="{ row }">
            <span class="amount">¥{{ row.amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <span class="status-tag" :class="row.status">
              {{ getStatusText(row.status) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="下单时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewOrder(row)">详情</el-button>
            <el-button type="danger" link @click="cancelOrder(row)" v-if="row.status === 'pending'">
              取消
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const loading = ref(false)
const orders = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const searchForm = ref({
  order_no: '',
  status: ''
})

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待支付',
    completed: '已完成',
    cancelled: '已取消',
    refunded: '已退款'
  }
  return map[status] || status
}

const handleSearch = () => {
  currentPage.value = 1
}

const resetSearch = () => {
  searchForm.value = { order_no: '', status: '' }
  handleSearch()
}

const viewOrder = (order: any) => {
  // TODO: 查看订单详情
}

const cancelOrder = async (order: any) => {
  // TODO: 取消订单
}

const exportOrders = () => {
  // TODO: 导出订单
}
</script>

<style scoped lang="scss">
.orders-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    
    h3 {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
    }
  }
  
  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
