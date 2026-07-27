<template>
  <n-card title="账单管理">
    <template #header-extra>
      <n-space>
        <n-select v-model:value="filters.status" :options="statusOptions" placeholder="账单状态" clearable style="width: 140px" />
        <n-date-picker v-model:value="filters.dateRange" type="daterange" clearable style="width: 260px" />
        <n-input v-model:value="filters.keyword" placeholder="搜索账单号/用户名" clearable style="width: 200px" @keydown.enter="handleSearch" />
        <n-button type="primary" @click="handleSearch">搜索</n-button>
      </n-space>
    </template>

    <n-data-table :columns="columns" :data="invoices" :loading="loading" :bordered="false" :pagination="pagination" @update:page="handlePageChange" @update:page-size="handlePageSizeChange" />

    <!-- Invoice Detail Modal -->
    <n-modal v-model:show="detailVisible" preset="card" :title="`账单详情 - ${currentInvoice?.invoiceNo || ''}`" style="width: 520px">
      <n-descriptions bordered :column="1" label-placement="left">
        <n-descriptions-item label="账单号">{{ currentInvoice?.invoiceNo }}</n-descriptions-item>
        <n-descriptions-item label="关联订单">{{ currentInvoice?.orderNo }}</n-descriptions-item>
        <n-descriptions-item label="用户">{{ currentInvoice?.user }}</n-descriptions-item>
        <n-descriptions-item label="金额">¥{{ currentInvoice?.amount?.toFixed(2) }}</n-descriptions-item>
        <n-descriptions-item label="状态">
          <n-tag :type="statusMap[currentInvoice?.status || '']?.type" size="small" round>
            {{ statusMap[currentInvoice?.status || '']?.label }}
          </n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="创建时间">{{ currentInvoice?.createdAt }}</n-descriptions-item>
        <n-descriptions-item label="到期时间">{{ currentInvoice?.dueDate }}</n-descriptions-item>
        <n-descriptions-item v-if="currentInvoice?.paidAt" label="支付时间">{{ currentInvoice?.paidAt }}</n-descriptions-item>
      </n-descriptions>
      <template #footer>
        <n-space justify="end">
          <n-button v-if="currentInvoice?.status === 'unpaid'" type="success" @click="markPaid">标记已支付</n-button>
          <n-button @click="detailVisible = false">关闭</n-button>
        </n-space>
      </template>
    </n-modal>
  </n-card>
</template>

<script setup lang="ts">
import { h, ref, reactive } from 'vue'
import { useMessage, NTag, NButton, NSpace, NPopconfirm, type DataTableColumns, type PaginationProps } from 'naive-ui'

const message = useMessage()
const loading = ref(false)
const detailVisible = ref(false)
const currentInvoice = ref<any>(null)

const filters = reactive({
  status: null as string | null,
  dateRange: null as [number, number] | null,
  keyword: '',
})

const statusOptions = [
  { label: '已支付', value: 'paid' },
  { label: '未支付', value: 'unpaid' },
  { label: '已过期', value: 'expired' },
  { label: '已退款', value: 'refunded' },
]

const statusMap: Record<string, { label: string; type: 'success' | 'info' | 'warning' | 'error' }> = {
  paid: { label: '已支付', type: 'success' },
  unpaid: { label: '未支付', type: 'warning' },
  expired: { label: '已过期', type: 'error' },
  refunded: { label: '已退款', type: 'info' },
}

const pagination = reactive<PaginationProps>({
  page: 1,
  pageSize: 10,
  itemCount: 30,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
})

const invoices = ref([
  { id: 1, invoiceNo: 'INV-2024001', orderNo: 'ORD-2024001', user: '张三', amount: 299.00, status: 'paid', createdAt: '2024-03-15 14:30', dueDate: '2024-03-22', paidAt: '2024-03-15 14:35' },
  { id: 2, invoiceNo: 'INV-2024002', orderNo: 'ORD-2024002', user: '李四', amount: 599.00, status: 'unpaid', createdAt: '2024-03-15 13:22', dueDate: '2024-03-22', paidAt: null },
  { id: 3, invoiceNo: 'INV-2024003', orderNo: 'ORD-2024003', user: '王五', amount: 1299.00, status: 'unpaid', createdAt: '2024-03-15 11:45', dueDate: '2024-03-22', paidAt: null },
  { id: 4, invoiceNo: 'INV-2024004', orderNo: 'ORD-2024004', user: '赵六', amount: 89.00, status: 'paid', createdAt: '2024-03-15 10:18', dueDate: '2024-03-22', paidAt: '2024-03-15 11:00' },
  { id: 5, invoiceNo: 'INV-2024005', orderNo: 'ORD-2024006', user: '周八', amount: 69.00, status: 'expired', createdAt: '2024-03-01 09:30', dueDate: '2024-03-08', paidAt: null },
  { id: 6, invoiceNo: 'INV-2024006', orderNo: 'ORD-2024007', user: '吴九', amount: 299.00, status: 'unpaid', createdAt: '2024-03-13 22:10', dueDate: '2024-03-20', paidAt: null },
])

const columns: DataTableColumns<any> = [
  { title: '账单号', key: 'invoiceNo', width: 150 },
  { title: '关联订单', key: 'orderNo', width: 150 },
  { title: '用户', key: 'user', width: 100 },
  { title: '金额', key: 'amount', width: 120, render: (row) => `¥${row.amount.toFixed(2)}` },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row) => {
      const s = statusMap[row.status]
      return h(NTag, { type: s.type, size: 'small', round: true }, { default: () => s.label })
    },
  },
  { title: '创建时间', key: 'createdAt', width: 160 },
  { title: '到期时间', key: 'dueDate', width: 120 },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row) =>
      h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, { size: 'small', onClick: () => openDetail(row) }, { default: () => '详情' }),
          h(NPopconfirm, { onPositiveClick: () => handleDelete(row.id) }, {
            trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '删除' }),
            default: () => '确认删除该账单？',
          }),
        ],
      }),
  },
]

function openDetail(invoice: any) {
  currentInvoice.value = invoice
  detailVisible.value = true
}

function markPaid() {
  if (currentInvoice.value) {
    currentInvoice.value.status = 'paid'
    currentInvoice.value.paidAt = new Date().toLocaleString()
    message.success('账单已标记为已支付')
    detailVisible.value = false
  }
}

function handleDelete(id: number) {
  invoices.value = invoices.value.filter((i) => i.id !== id)
  message.success('账单已删除')
}

function handleSearch() {
  pagination.page = 1
}

function handlePageChange(page: number) {
  pagination.page = page
}

function handlePageSizeChange(pageSize: number) {
  pagination.pageSize = pageSize
  pagination.page = 1
}
</script>
