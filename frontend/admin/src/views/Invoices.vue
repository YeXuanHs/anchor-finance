<template>
  <div class="admin-page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">账单管理</span>
          <div class="card-actions">
            <el-select v-model="filters.status" placeholder="账单状态" clearable style="width: 130px">
              <el-option v-for="o in statusOptions" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
            <el-date-picker v-model="filters.dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" style="width: 260px" clearable />
            <el-input v-model="filters.keyword" placeholder="搜索账单号/用户名" clearable style="width: 200px" @keydown.enter="handleSearch">
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
          </div>
        </div>
      </template>

      <el-table :data="invoices" v-loading="loading" stripe size="small">
        <el-table-column prop="invoiceNo" label="账单号" width="150" />
        <el-table-column prop="orderNo" label="关联订单" width="150" />
        <el-table-column prop="user" label="用户" width="90" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            <span style="font-weight: 600; color: #0056FF">¥{{ row.amount.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusMap[row.status]?.type as any" size="small" round>{{ statusMap[row.status]?.label }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="160" />
        <el-table-column prop="dueDate" label="到期时间" width="120" />
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="openDetail(row)">详情</el-button>
            <el-popconfirm title="确认删除该账单？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button text type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="30"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @current-change="(p: number) => pagination.page = p"
          @size-change="(s: number) => { pagination.pageSize = s; pagination.page = 1 }"
        />
      </div>
    </el-card>

    <!-- Invoice Detail Dialog -->
    <el-dialog v-model="detailVisible" :title="`账单详情 - ${currentInvoice?.invoiceNo || ''}`" width="520px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="账单号">{{ currentInvoice?.invoiceNo }}</el-descriptions-item>
        <el-descriptions-item label="关联订单">{{ currentInvoice?.orderNo }}</el-descriptions-item>
        <el-descriptions-item label="用户">{{ currentInvoice?.user }}</el-descriptions-item>
        <el-descriptions-item label="金额">¥{{ currentInvoice?.amount?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusMap[currentInvoice?.status || '']?.type as any" size="small" round>{{ statusMap[currentInvoice?.status || '']?.label }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentInvoice?.createdAt }}</el-descriptions-item>
        <el-descriptions-item label="到期时间">{{ currentInvoice?.dueDate }}</el-descriptions-item>
        <el-descriptions-item v-if="currentInvoice?.paidAt" label="支付时间">{{ currentInvoice?.paidAt }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button v-if="currentInvoice?.status === 'unpaid'" type="success" @click="markPaid">标记已支付</el-button>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'

const loading = ref(false)
const detailVisible = ref(false)
const currentInvoice = ref<any>(null)

const filters = reactive({ status: null as string | null, dateRange: null as any, keyword: '' })
const pagination = reactive({ page: 1, pageSize: 10 })

const statusOptions = [
  { label: '已支付', value: 'paid' }, { label: '未支付', value: 'unpaid' },
  { label: '已过期', value: 'expired' }, { label: '已退款', value: 'refunded' },
]
const statusMap: Record<string, { label: string; type: string }> = {
  paid: { label: '已支付', type: 'success' }, unpaid: { label: '未支付', type: 'warning' },
  expired: { label: '已过期', type: 'danger' }, refunded: { label: '已退款', type: 'info' },
}

const invoices = ref([
  { id: 1, invoiceNo: 'INV-2024001', orderNo: 'ORD-2024001', user: '张三', amount: 299.00, status: 'paid', createdAt: '2024-03-15 14:30', dueDate: '2024-03-22', paidAt: '2024-03-15 14:35' },
  { id: 2, invoiceNo: 'INV-2024002', orderNo: 'ORD-2024002', user: '李四', amount: 599.00, status: 'unpaid', createdAt: '2024-03-15 13:22', dueDate: '2024-03-22', paidAt: null },
  { id: 3, invoiceNo: 'INV-2024003', orderNo: 'ORD-2024003', user: '王五', amount: 1299.00, status: 'unpaid', createdAt: '2024-03-15 11:45', dueDate: '2024-03-22', paidAt: null },
  { id: 4, invoiceNo: 'INV-2024004', orderNo: 'ORD-2024004', user: '赵六', amount: 89.00, status: 'paid', createdAt: '2024-03-15 10:18', dueDate: '2024-03-22', paidAt: '2024-03-15 11:00' },
  { id: 5, invoiceNo: 'INV-2024005', orderNo: 'ORD-2024006', user: '周八', amount: 69.00, status: 'expired', createdAt: '2024-03-01 09:30', dueDate: '2024-03-08', paidAt: null },
  { id: 6, invoiceNo: 'INV-2024006', orderNo: 'ORD-2024007', user: '吴九', amount: 299.00, status: 'unpaid', createdAt: '2024-03-13 22:10', dueDate: '2024-03-20', paidAt: null },
])

function openDetail(invoice: any) { currentInvoice.value = invoice; detailVisible.value = true }

function markPaid() {
  if (currentInvoice.value) {
    currentInvoice.value.status = 'paid'
    currentInvoice.value.paidAt = new Date().toLocaleString()
    ElMessage.success('账单已标记为已支付')
    detailVisible.value = false
  }
}

function handleDelete(id: number) {
  invoices.value = invoices.value.filter((i) => i.id !== id)
  ElMessage.success('账单已删除')
}

function handleSearch() { pagination.page = 1 }
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 12px; }
.card-title { font-size: 16px; font-weight: 600; }
.card-actions { display: flex; align-items: center; gap: 12px; }
</style>
