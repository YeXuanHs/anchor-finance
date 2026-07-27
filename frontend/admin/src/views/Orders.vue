<template>
  <n-card title="订单管理">
    <template #header-extra>
      <n-space>
        <n-select v-model:value="filters.status" :options="statusOptions" placeholder="订单状态" clearable style="width: 140px" />
        <n-date-picker v-model:value="filters.dateRange" type="daterange" clearable style="width: 260px" />
        <n-input v-model:value="filters.keyword" placeholder="搜索订单号/用户名" clearable style="width: 200px" @keydown.enter="handleSearch" />
        <n-button type="primary" @click="handleSearch">搜索</n-button>
      </n-space>
    </template>

    <n-data-table :columns="columns" :data="orders" :loading="loading" :bordered="false" :pagination="pagination" @update:page="handlePageChange" @update:page-size="handlePageSizeChange" />

    <!-- Order Detail Drawer -->
    <n-drawer v-model:show="drawerVisible" :width="500">
      <n-drawer-content :title="`订单详情 - ${currentOrder?.id || ''}`">
        <n-descriptions bordered :column="1" label-placement="left">
          <n-descriptions-item label="订单号">{{ currentOrder?.id }}</n-descriptions-item>
          <n-descriptions-item label="用户">{{ currentOrder?.user }}</n-descriptions-item>
          <n-descriptions-item label="产品">{{ currentOrder?.product }}</n-descriptions-item>
          <n-descriptions-item label="金额">¥{{ currentOrder?.amount }}</n-descriptions-item>
          <n-descriptions-item label="状态">
            <n-tag :type="statusMap[currentOrder?.status || '']?.type" size="small" round>
              {{ statusMap[currentOrder?.status || '']?.label }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="下单时间">{{ currentOrder?.createdAt }}</n-descriptions-item>
          <n-descriptions-item label="备注">{{ currentOrder?.remark || '无' }}</n-descriptions-item>
        </n-descriptions>

        <n-divider>更新状态</n-divider>
        <n-space>
          <n-button v-if="currentOrder?.status === 'pending'" type="info" size="small" @click="updateStatus('processing')">开始处理</n-button>
          <n-button v-if="currentOrder?.status === 'processing'" type="success" size="small" @click="updateStatus('completed')">标记完成</n-button>
          <n-button v-if="currentOrder?.status !== 'cancelled' && currentOrder?.status !== 'completed'" type="error" size="small" @click="updateStatus('cancelled')">取消订单</n-button>
        </n-space>
      </n-drawer-content>
    </n-drawer>
  </n-card>
</template>

<script setup lang="ts">
import { h, ref, reactive } from 'vue'
import { useMessage, NTag, NButton, NSpace, NPopconfirm, type DataTableColumns, type PaginationProps } from 'naive-ui'

const message = useMessage()
const loading = ref(false)
const drawerVisible = ref(false)
const currentOrder = ref<any>(null)

const filters = reactive({
  status: null as string | null,
  dateRange: null as [number, number] | null,
  keyword: '',
})

const statusOptions = [
  { label: '待支付', value: 'pending' },
  { label: '处理中', value: 'processing' },
  { label: '已完成', value: 'completed' },
  { label: '已取消', value: 'cancelled' },
]

const statusMap: Record<string, { label: string; type: 'success' | 'info' | 'warning' | 'error' }> = {
  pending: { label: '待支付', type: 'warning' },
  processing: { label: '处理中', type: 'info' },
  completed: { label: '已完成', type: 'success' },
  cancelled: { label: '已取消', type: 'error' },
  refunded: { label: '已退款', type: 'error' },
}

const pagination = reactive<PaginationProps>({
  page: 1,
  pageSize: 10,
  itemCount: 50,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
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

const columns: DataTableColumns<any> = [
  { title: '订单号', key: 'id', width: 150 },
  { title: '用户', key: 'user', width: 100 },
  { title: '产品', key: 'product' },
  { title: '金额', key: 'amount', width: 100, render: (row) => `¥${row.amount}` },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row) => {
      const s = statusMap[row.status]
      return h(NTag, { type: s.type, size: 'small', round: true }, { default: () => s.label })
    },
  },
  { title: '下单时间', key: 'createdAt', width: 180 },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row) =>
      h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, { size: 'small', onClick: () => openDrawer(row) }, { default: () => '详情' }),
          h(NPopconfirm, { onPositiveClick: () => handleDelete(row.id) }, {
            trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '删除' }),
            default: () => '确认删除该订单？',
          }),
        ],
      }),
  },
]

function openDrawer(order: any) {
  currentOrder.value = order
  drawerVisible.value = true
}

function updateStatus(status: string) {
  if (currentOrder.value) {
    currentOrder.value.status = status
    message.success(`订单状态已更新为: ${statusMap[status]?.label}`)
    drawerVisible.value = false
  }
}

function handleDelete(id: string) {
  orders.value = orders.value.filter((o) => o.id !== id)
  message.success('订单已删除')
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
