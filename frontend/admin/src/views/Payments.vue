<template>
  <div>
    <n-card :bordered="false" rounded>
      <n-tabs v-model:value="activeTab" type="line" animated>
        <!-- Payment Methods Tab -->
        <n-tab-pane name="methods" tab="支付方式">
          <template #header-extra>
            <n-space style="padding: 12px 0 4px">
              <n-input
                v-model:value="searchKeyword"
                placeholder="搜索支付名称"
                clearable
                style="width: 240px"
                @clear="handleSearch"
                @keydown.enter="handleSearch"
              >
                <template #prefix>
                  <n-icon><SearchIcon /></n-icon>
                </template>
              </n-input>
              <n-select
                v-model:value="filterType"
                placeholder="支付类型"
                clearable
                style="width: 150px"
                :options="typeOptions"
              />
              <n-select
                v-model:value="filterStatus"
                placeholder="状态"
                clearable
                style="width: 120px"
                :options="statusOptions"
              />
              <n-button type="primary" @click="openPaymentModal()">
                <template #icon><n-icon><AddIcon /></n-icon></template>
                添加支付方式
              </n-button>
            </n-space>
          </template>

          <n-data-table
            :columns="paymentColumns"
            :data="filteredPayments"
            :loading="loading"
            :bordered="false"
            :row-key="(row: any) => row.id"
            size="small"
          />
        </n-tab-pane>

        <!-- Payment Records Tab -->
        <n-tab-pane name="records" tab="支付记录">
          <template #header-extra>
            <n-space style="padding: 12px 0 4px">
              <n-input
                v-model:value="recordSearchKeyword"
                placeholder="搜索订单号/用户名"
                clearable
                style="width: 240px"
                @clear="handleRecordSearch"
                @keydown.enter="handleRecordSearch"
              >
                <template #prefix>
                  <n-icon><SearchIcon /></n-icon>
                </template>
              </n-input>
              <n-date-picker
                v-model:value="recordDateRange"
                type="daterange"
                clearable
                style="width: 300px"
              />
            </n-space>
          </template>

          <n-data-table
            :columns="recordColumns"
            :data="filteredRecords"
            :loading="recordsLoading"
            :bordered="false"
            :row-key="(row: any) => row.id"
            size="small"
          />
        </n-tab-pane>
      </n-tabs>
    </n-card>

    <!-- Payment Edit Modal -->
    <n-modal
      v-model:show="paymentModalVisible"
      preset="card"
      :title="editingPayment ? '编辑支付方式' : '添加支付方式'"
      style="width: 640px"
      :bordered="false"
      :segmented="{ content: true, footer: true }"
    >
      <n-form ref="paymentFormRef" :model="paymentForm" :rules="paymentRules" label-placement="left" label-width="100">
        <n-form-item label="支付名称" path="name">
          <n-input v-model:value="paymentForm.name" placeholder="请输入支付名称" />
        </n-form-item>
        <n-form-item label="支付类型" path="type">
          <n-select v-model:value="paymentForm.type" :options="typeOptions" placeholder="选择支付类型" />
        </n-form-item>
        <n-form-item label="手续费率(%)" path="feeRate">
          <n-input-number v-model:value="paymentForm.feeRate" :min="0" :max="100" :precision="2" style="width: 100%">
            <template #suffix>%</template>
          </n-input-number>
        </n-form-item>
        <n-form-item label="App ID" path="appId">
          <n-input v-model:value="paymentForm.appId" placeholder="支付网关 App ID" />
        </n-form-item>
        <n-form-item label="App Secret" path="appSecret">
          <n-input v-model:value="paymentForm.appSecret" type="password" show-password-on="click" placeholder="支付网关 App Secret" />
        </n-form-item>
        <n-form-item label="回调URL" path="notifyUrl">
          <n-input v-model:value="paymentForm.notifyUrl" placeholder="支付回调通知URL" />
        </n-form-item>
        <n-form-item label="排序" path="sort">
          <n-input-number v-model:value="paymentForm.sort" :min="0" style="width: 100%" />
        </n-form-item>
        <n-form-item label="状态" path="status">
          <n-switch v-model:value="paymentForm.status" checked-value="active" unchecked-value="inactive">
            <template #checked>启用</template>
            <template #unchecked>禁用</template>
          </n-switch>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="paymentModalVisible = false">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handlePaymentSubmit">确定</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { h, ref, reactive, computed, defineComponent } from 'vue'
import {
  useMessage,
  NTag,
  NButton,
  NSwitch,
  NSpace,
  NPopconfirm,
  NTooltip,
  NIcon,
  type DataTableColumns,
  type FormInst,
  type FormRules,
} from 'naive-ui'
import {
  SearchOutline as SearchIcon,
  AddOutline as AddIcon,
  CreateOutline as EditIcon,
  TrashOutline as DeleteIcon,
} from '@vicons/ionicons5'

const message = useMessage()
const activeTab = ref('methods')
const loading = ref(false)
const recordsLoading = ref(false)
const submitting = ref(false)
const paymentModalVisible = ref(false)
const searchKeyword = ref('')
const filterType = ref<string | null>(null)
const filterStatus = ref<string | null>(null)
const paymentFormRef = ref<FormInst | null>(null)
const editingPayment = ref<any>(null)
const recordSearchKeyword = ref('')
const recordDateRange = ref<[number, number] | null>(null)

// ---- Options ----
const typeOptions = [
  { label: '支付宝', value: 'alipay' },
  { label: '微信支付', value: 'wechat' },
  { label: 'PayPal', value: 'paypal' },
  { label: 'Stripe', value: 'stripe' },
  { label: '虎皮椒', value: 'hupijiao' },
  { label: '易支付', value: 'epay' },
  { label: 'USDT', value: 'usdt' },
  { label: '余额', value: 'balance' },
]

const statusOptions = [
  { label: '启用', value: 'active' },
  { label: '禁用', value: 'inactive' },
]

const typeNameMap: Record<string, string> = {
  alipay: '支付宝',
  wechat: '微信支付',
  paypal: 'PayPal',
  stripe: 'Stripe',
  hupijiao: '虎皮椒',
  epay: '易支付',
  usdt: 'USDT',
  balance: '余额',
}

const statusNameMap: Record<string, string> = {
  active: '启用',
  inactive: '禁用',
}

// ---- Payment Methods ----
const payments = ref([
  { id: 1, name: '支付宝', type: 'alipay', feeRate: 0.6, appId: '2021000000000001', status: 'active', sort: 1, totalAmount: 125680.50, totalCount: 1256 },
  { id: 2, name: '微信支付', type: 'wechat', feeRate: 0.6, appId: 'wx1234567890abcdef', status: 'active', sort: 2, totalAmount: 98450.00, totalCount: 1589 },
  { id: 3, name: 'PayPal', type: 'paypal', feeRate: 3.5, appId: 'AYSq3RDGrs99999999999999', status: 'active', sort: 3, totalAmount: 45200.75, totalCount: 234 },
  { id: 4, name: 'Stripe', type: 'stripe', feeRate: 2.9, appId: 'pk_live_1234567890', status: 'active', sort: 4, totalAmount: 78900.00, totalCount: 456 },
  { id: 5, name: '虎皮椒', type: 'hupijiao', feeRate: 1.0, appId: '100001', status: 'inactive', sort: 5, totalAmount: 12300.00, totalCount: 89 },
  { id: 6, name: '易支付', type: 'epay', feeRate: 0.5, appId: '1001', status: 'active', sort: 6, totalAmount: 56700.25, totalCount: 678 },
  { id: 7, name: 'USDT(TRC20)', type: 'usdt', feeRate: 0, appId: 'TRX1234567890', status: 'active', sort: 7, totalAmount: 234500.00, totalCount: 123 },
  { id: 8, name: '余额支付', type: 'balance', feeRate: 0, appId: '-', status: 'active', sort: 8, totalAmount: 45600.00, totalCount: 567 },
])

const filteredPayments = computed(() => {
  return payments.value.filter((p) => {
    if (searchKeyword.value.trim()) {
      const kw = searchKeyword.value.trim().toLowerCase()
      if (!p.name.toLowerCase().includes(kw)) return false
    }
    if (filterType.value && p.type !== filterType.value) return false
    if (filterStatus.value && p.status !== filterStatus.value) return false
    return true
  })
})

// ---- Payment Records ----
const records = ref([
  { id: 1, orderNo: 'ORD20260727001', username: 'user001', paymentName: '支付宝', amount: 299.00, fee: 1.79, status: 'success', createdAt: '2026-07-27 10:15:32' },
  { id: 2, orderNo: 'ORD20260727002', username: 'user002', paymentName: '微信支付', amount: 599.00, fee: 3.59, status: 'success', createdAt: '2026-07-27 10:20:45' },
  { id: 3, orderNo: 'ORD20260727003', username: 'user003', paymentName: 'PayPal', amount: 1299.00, fee: 45.47, status: 'success', createdAt: '2026-07-27 10:25:18' },
  { id: 4, orderNo: 'ORD20260727004', username: 'user004', paymentName: 'Stripe', amount: 89.00, fee: 2.58, status: 'pending', createdAt: '2026-07-27 10:30:56' },
  { id: 5, orderNo: 'ORD20260727005', username: 'user005', paymentName: '支付宝', amount: 399.00, fee: 2.39, status: 'failed', createdAt: '2026-07-27 10:35:22' },
  { id: 6, orderNo: 'ORD20260727006', username: 'user006', paymentName: 'USDT(TRC20)', amount: 199.00, fee: 0, status: 'success', createdAt: '2026-07-27 10:40:11' },
  { id: 7, orderNo: 'ORD20260727007', username: 'user007', paymentName: '余额支付', amount: 69.00, fee: 0, status: 'success', createdAt: '2026-07-27 10:45:38' },
])

const filteredRecords = computed(() => {
  return records.value.filter((r) => {
    if (recordSearchKeyword.value.trim()) {
      const kw = recordSearchKeyword.value.trim().toLowerCase()
      if (!r.orderNo.toLowerCase().includes(kw) && !r.username.toLowerCase().includes(kw)) return false
    }
    // Date filtering would be implemented here with recordDateRange
    return true
  })
})

// ---- Payment Form ----
const paymentForm = reactive({
  name: '',
  type: null as string | null,
  feeRate: 0,
  appId: '',
  appSecret: '',
  notifyUrl: '',
  sort: 0,
  status: 'active',
})

const paymentRules: FormRules = {
  name: { required: true, message: '请输入支付名称', trigger: 'blur' },
  type: { required: true, message: '请选择支付类型', trigger: 'change' },
  feeRate: { required: true, type: 'number', message: '请输入手续费率', trigger: 'blur' },
  appId: { required: true, message: '请输入App ID', trigger: 'blur' },
}

// ---- Payment Table Columns ----
const paymentColumns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 60, sorter: (a, b) => a.id - b.id },
  { title: '名称', key: 'name', ellipsis: { tooltip: true } },
  {
    title: '类型',
    key: 'type',
    width: 100,
    render: (row) => h(NTag, { size: 'small', round: true, bordered: false, type: 'info' }, { default: () => typeNameMap[row.type] || row.type }),
  },
  {
    title: '手续费率',
    key: 'feeRate',
    width: 100,
    sorter: (a, b) => a.feeRate - b.feeRate,
    render: (row) => h('span', { style: 'font-weight:600;color:#1890ff' }, `${row.feeRate}%`),
  },
  {
    title: '累计金额',
    key: 'totalAmount',
    width: 120,
    sorter: (a, b) => a.totalAmount - b.totalAmount,
    render: (row) => h('span', { style: 'font-weight:600;color:#18a058' }, `¥${row.totalAmount.toLocaleString()}`),
  },
  {
    title: '交易笔数',
    key: 'totalCount',
    width: 100,
    sorter: (a, b) => a.totalCount - b.totalCount,
  },
  {
    title: '排序',
    key: 'sort',
    width: 70,
    sorter: (a, b) => a.sort - b.sort,
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) =>
      h(NSwitch, {
        value: row.status === 'active',
        size: 'small',
        onUpdateValue: () => handleTogglePayment(row),
      }),
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row) =>
      h(NSpace, { size: 4 }, {
        default: () => [
          h(NTooltip, {}, {
            trigger: () =>
              h(NButton, { size: 'small', quaternary: true, type: 'primary', onClick: () => openPaymentModal(row) }, {
                icon: () => h(NIcon, null, { default: () => h(EditIcon) }),
              }),
            default: () => '配置',
          }),
          h(NPopconfirm, { onPositiveClick: () => handleTogglePayment(row) }, {
            trigger: () =>
              h(NTooltip, {}, {
                trigger: () =>
                  h(NButton, { size: 'small', quaternary: true, type: 'warning' }, {
                    icon: () => h(NIcon, null, { default: () => h(StatusIcon) }),
                  }),
                default: () => row.status === 'active' ? '禁用' : '启用',
              }),
            default: () => `确认${row.status === 'active' ? '禁用' : '启用'}该支付方式？`,
          }),
          h(NPopconfirm, { onPositiveClick: () => handleDeletePayment(row.id) }, {
            trigger: () =>
              h(NTooltip, {}, {
                trigger: () =>
                  h(NButton, { size: 'small', quaternary: true, type: 'error' }, {
                    icon: () => h(NIcon, null, { default: () => h(DeleteIcon) }),
                  }),
                default: () => '删除',
              }),
            default: () => `确认删除支付方式「${row.name}」？`,
          }),
        ],
      }),
  },
]

// Status icon component
const StatusIcon = defineComponent({
  render: () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 512 512', fill: 'currentColor' }, [
    h('path', { d: 'M256 48C141.13 48 48 141.13 48 256s93.13 208 208 208 208-93.13 208-208S370.87 48 256 48zm-24 312h48v48h-48v-48zm0-80h48v80h-48v-80zm0-160h48v120h-48V120z' }),
  ]),
})

// ---- Record Table Columns ----
const recordStatusMap: Record<string, { label: string; type: 'success' | 'warning' | 'error' | 'info' }> = {
  success: { label: '成功', type: 'success' },
  pending: { label: '待处理', type: 'warning' },
  failed: { label: '失败', type: 'error' },
  refunded: { label: '已退款', type: 'info' },
}

const recordColumns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 60, sorter: (a, b) => a.id - b.id },
  { title: '订单号', key: 'orderNo', width: 180, ellipsis: { tooltip: true } },
  { title: '用户名', key: 'username', width: 100 },
  { title: '支付方式', key: 'paymentName', width: 120 },
  {
    title: '金额',
    key: 'amount',
    width: 120,
    sorter: (a, b) => a.amount - b.amount,
    render: (row) => h('span', { style: 'font-weight:600;color:#18a058' }, `¥${row.amount.toFixed(2)}`),
  },
  {
    title: '手续费',
    key: 'fee',
    width: 100,
    render: (row) => h('span', { style: 'color:#f5a623' }, `¥${row.fee.toFixed(2)}`),
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) => {
      const status = recordStatusMap[row.status]
      return h(NTag, { size: 'small', round: true, bordered: false, type: status?.type || 'default' }, { default: () => status?.label || row.status })
    },
  },
  { title: '时间', key: 'createdAt', width: 180, sorter: (a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime() },
]

// ---- Actions ----
function openPaymentModal(payment?: any) {
  editingPayment.value = payment || null
  if (payment) {
    Object.assign(paymentForm, {
      name: payment.name,
      type: payment.type,
      feeRate: payment.feeRate,
      appId: payment.appId,
      appSecret: '',
      notifyUrl: '',
      sort: payment.sort,
      status: payment.status,
    })
  } else {
    Object.assign(paymentForm, { name: '', type: null, feeRate: 0, appId: '', appSecret: '', notifyUrl: '', sort: 0, status: 'active' })
  }
  paymentModalVisible.value = true
}

async function handlePaymentSubmit() {
  try { await paymentFormRef.value?.validate() } catch { return }
  submitting.value = true
  try {
    // TODO: API call
    message.success(editingPayment.value ? '支付方式更新成功' : '支付方式添加成功')
    paymentModalVisible.value = false
  } catch (err: any) {
    message.error(err.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

function handleTogglePayment(payment: any) {
  payment.status = payment.status === 'active' ? 'inactive' : 'active'
  message.success(`支付方式「${payment.name}」已${payment.status === 'active' ? '启用' : '禁用'}`)
}

function handleDeletePayment(id: number) {
  payments.value = payments.value.filter((p) => p.id !== id)
  message.success('支付方式已删除')
}

function handleSearch() {
  // filter is reactive via computed
}

function handleRecordSearch() {
  // filter is reactive via computed
}
</script>

<style scoped>
:deep(.n-card) {
  background-color: #1e1e2e;
  color: #cdd6f4;
}

:deep(.n-data-table) {
  --n-th-color: #181825;
  --n-td-color: #1e1e2e;
  --n-border-color: #313244;
  --n-th-text-color: #cdd6f4;
  --n-td-text-color: #cdd6f4;
}
</style>
