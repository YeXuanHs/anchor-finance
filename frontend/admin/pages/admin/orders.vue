<template>
  <div class="orders-page">
    <el-card class="admin-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>订单管理</span>
          <div class="header-actions">
            <el-select v-model="filters.status" placeholder="订单状态" clearable style="width: 140px">
              <el-option v-for="opt in statusOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
            </el-select>
            <el-date-picker v-model="filters.dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" style="width: 260px" />
            <el-input v-model="filters.keyword" placeholder="搜索订单号/用户名" clearable style="width: 200px" @keyup.enter="handleSearch" />
            <el-button type="primary" @click="handleSearch">搜索</el-button>
          </div>
        </div>
      </template>

      <el-table :data="orders" style="width: 100%" v-loading="loading" size="default">
        <el-table-column prop="id" label="订单号" width="150" />
        <el-table-column prop="user" label="用户" width="100" />
        <el-table-column prop="product" label="产品" show-overflow-tooltip />
        <el-table-column prop="amount" label="金额" width="100">
          <template #default="{ row }">
            <span class="amount">¥{{ row.amount }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusMap[row.status]?.type" size="small" round>
              {{ statusMap[row.status]?.label }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="下单时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDrawer(row)">详情</el-button>
            <el-popconfirm title="确认删除该订单？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="50"
          layout="total, sizes, prev, pager, next"
        />
      </div>
    </el-card>

    <!-- Order Detail Drawer -->
    <el-drawer v-model="drawerVisible" :title="`订单详情 - ${currentOrder?.id || ''}`" size="500px">
      <template v-if="currentOrder">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="订单号">{{ currentOrder.id }}</el-descriptions-item>
          <el-descriptions-item label="用户">{{ currentOrder.user }}</el-descriptions-item>
          <el-descriptions-item label="产品">{{ currentOrder.product }}</el-descriptions-item>
          <el-descriptions-item label="金额">
            <span class="amount">¥{{ currentOrder.amount }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="statusMap[currentOrder.status]?.type" size="small">
              {{ statusMap[currentOrder.status]?.label }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="下单时间">{{ currentOrder.createdAt }}</el-descriptions-item>
          <el-descriptions-item label="备注">{{ currentOrder.remark || '无' }}</el-descriptions-item>
        </el-descriptions>

        <el-divider>更新状态</el-divider>
        <el-space>
          <el-button v-if="currentOrder.status === 'pending'" type="primary" size="small" @click="updateStatus('processing')">开始处理</el-button>
          <el-button v-if="currentOrder.status === 'processing'" type="success" size="small" @click="updateStatus('completed')">标记完成</el-button>
          <el-button v-if="currentOrder.status !== 'cancelled' && currentOrder.status !== 'completed'" type="danger" size="small" @click="updateStatus('cancelled')">取消订单</el-button>
        </el-space>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'

definePageMeta({
  layout: 'admin',
})

const loading = ref(false)
const drawerVisible = ref(false)
const currentOrder = ref<any>(null)

const filters = reactive({
  status: null as string | null,
  dateRange: null as any,
  keyword: '',
})

const statusOptions = [
  { label: '待支付', value: 'pending' },
  { label: '处理中', value: 'processing' },
  { label: '已完成', value: 'completed' },
  { label: '已取消', value: 'cancelled' },
]

const statusMap: Record<string, { label: string; type: '' | 'success' | 'warning' | 'info' | 'danger' }> = {
  pending: { label: '待支付', type: 'warning' },
  processing: { label: '处理中', type: 'info' },
  completed: { label: '已完成', type: 'success' },
  cancelled: { label: '已取消', type: 'danger' },
  refunded: { label: '已退款', type: 'danger' },
}

const pagination = reactive({
  page: 1,
  pageSize: 10,
})

const orders = ref([
  { id: 'ORD-2024001', user: '张三', product: '基础版主机', amount: 299, status: 'completed', createdAt: '2024-03-15 14:30:00', remark: '' },
  { id: 'ORD-2024002', user: '李四', product: '高级版主机', amount: 599, status: 'processing', createdAt: '2024-03-15 13:22:00', remark: '加急处理' },
  { id: 'ORD-2024003', user: '王五', product: '企业版主机', amount: 1299, status: 'pending', createdAt: '2024-03-15 11:45:00', remark: '' },
  { id: 'ORD-2024004', user: '赵六', product: '1核2G云服务器', amount: 89, status: 'completed', createdAt: '2024-03-15 10:18:00', remark: '' },
  { id: 'ORD-2024005', user: '孙七', product: '4核8G云服务器', amount: 399, status: 'cancelled', createdAt: '2024-03-14 16:55:00', remark: '用户取消' },
  { id: 'ORD-2024006', user: '周八', product: '.com域名注册', amount: 69, status: 'completed', createdAt: '2024-03-14 09:30:00', remark: '' },
  { id: 'ORD-2024007', user: '吴九', product: '基础版主机', amount: 299, status: 'pending', createdAt: '2024-03-13 22:10:00', remark: '' },
])

function openDrawer(order: any) {
  currentOrder.value = order
  drawerVisible.value = true
}

function updateStatus(status: string) {
  if (currentOrder.value) {
    currentOrder.value.status = status
    ElMessage.success(`订单状态已更新为: ${statusMap[status]?.label}`)
    drawerVisible.value = false
  }
}

function handleDelete(id: string) {
  orders.value = orders.value.filter((o) => o.id !== id)
  ElMessage.success('订单已删除')
}

function handleSearch() {
  pagination.page = 1
}
</script>

<style scoped>
.orders-page {
  padding: 0;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.amount {
  font-weight: 600;
  color: #409eff;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
