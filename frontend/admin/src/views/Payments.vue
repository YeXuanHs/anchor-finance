<template>
  <div>
    <n-card :bordered="false" rounded>
      <template #header>
        <n-space align="center" justify="space-between">
          <span style="font-size: 18px; font-weight: 600">支付方式管理</span>
          <n-space>
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
              :options="typeFilterOptions"
              placeholder="筛选类型"
              clearable
              style="width: 140px"
            />
            <n-select
              v-model:value="filterStatus"
              :options="statusFilterOptions"
              placeholder="筛选状态"
              clearable
              style="width: 120px"
            />
            <n-button type="primary" @click="openPaymentModal()">
              <template #icon><n-icon><AddIcon /></n-icon></template>
              添加支付方式
            </n-button>
          </n-space>
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
    </n-card>

    <!-- Payment Config Modal -->
    <n-modal
      v-model:show="paymentModalVisible"
      preset="card"
      :title="editingPayment ? '编辑支付方式' : '添加支付方式'"
      style="width: 600px"
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
        <n-form-item label="AppID" path="appId">
          <n-input v-model:value="paymentForm.appId" placeholder="请输入AppID" />
        </n-form-item>
        <n-form-item label="AppSecret" path="appSecret">
          <n-input v-model:value="paymentForm.appSecret" type="password" show-password-on="click" placeholder="请输入AppSecret" />
        </n-form-item>
        <n-form-item label="回调URL" path="notifyUrl">
          <n-input v-model:value="paymentForm.notifyUrl" placeholder="请输入回调URL" />
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

    <!-- Payment Records Modal -->
    <n-modal
      v-model:show="recordsModalVisible"
      preset="card"
      title="支付记录"
      style="width: 900px"
      :bordered="false"
      :segmented="{ content: true, footer: true }"
    >
      <n-data-table
        :columns="recordColumns"
        :data="paymentRecords"
        :bordered="false"
        :row-key="(row: any) => row.id"
        size="small"
        :max-height="400"
      />
      <template #footer>
        <n-space justify="end">
          <n-button @click="recordsModalVisible = false">关闭</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { h, ref, reactive, computed } from 'vue'
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
  ListOutline as ListIcon,
} from '@vicons/ionicons5'

const message = useMessage()
const loading = ref(false)
const submitting = ref(false)
const paymentModalVisible = ref(false)
const recordsModalVisible = ref(false)
const searchKeyword = ref('')
const filterType = ref<string | null>(null)
const filterStatus = ref<string | null>(null)
const paymentFormRef = ref<FormInst | null>(null)
const editingPayment = ref<any>(null)
const selectedPayment = ref<any>(null)

// ---- Options ----
const typeOptions = [
  { label: '支付宝', value: 'alipay' },
  { label: '微信支付', value: 'wechat' },
  { label: 'PayPal', value: 'paypal' },
  { label: 'Stripe', value: 'stripe' },
  { label: '虎皮椒', value: 'hupijiao' },
  { label: '易支付', value: 'epay' },
  { label: 'USDT', value: 'usdt' },
  { label: '余额支付', value: 'balance' },
]

const typeFilterOptions = [
  { label: '支付宝', value: 'alipay' },
  { label: '微信支付', value: 'wechat' },
  { label: 'PayPal', value: 'paypal' },
  { label: 'Stripe', value: 'stripe' },
  { label: '虎皮椒', value: 'hupijiao' },
  { label: '易支付', value: 'epay' },
  { label: 'USDT', value: 'usdt' },
  { label: '余额支付', value: 'balance' },
]

const statusFilterOptions = [
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
  balance: '余额支付',
}

const typeTagMap: Record<string, string> = {
  alipay: 'info',
  wechat: 'success',
  paypal: 'warning',
  stripe: 'primary',
  hupijiao: 'error',
  epay: 'info',
  usdt: 'success',
  balance: 'default',
}

// ---- Payments ----
const payments = ref([
  {
    id: 1,
    name: '支付宝官方',
    type: 'alipay',
    feeRate: 0.6,
    appId: '202100110066xxxx',
    appSecret: '********',
    notifyUrl: 'https://api.example.com/notify/alipay',
    status: 'active',
    todayAmount: 12580.5,
    todayCount: 45,
  },
  {
    id: 2,
    name: '微信支付',
    type: 'wechat',
    feeRate: 0.6,
    appId: 'wx1234567890abcdef',
    appSecret: '********',
    notifyUrl: 'https://api.example.com/notify/wechat',
    status: 'active',
    todayAmount: 8960.0,
    todayCount: 32,
  },
  {
    id: 3,
    name: 'PayPal国际',
    type: 'paypal',
    feeRate: 3.5,
    appId: 'AYSxxxxxxxxxxxxxxx',
    appSecret: '********',
    notifyUrl: 'https://api.example.com/notify/paypal',
    status: 'active',
    todayAmount: 2500.0,
    todayCount: 8,
  },
  {
    id: 4,
    name: 'Stripe',
    type: 'stripe',
    feeRate: 2.9,
    appId: 'pk_live_xxxxxxxxxxxx',
    appSecret: 'sk_live_xxxxxxxxxxxx',
    notifyUrl: 'https://api.example.com/notify/stripe',
    status: 'inactive',
    todayAmount: 0,
    todayCount: 0,
  },
  {
    id: 5,
    name: '虎皮椒',
    type: 'hupijiao',
    feeRate: 1.0,
    appId: 'hupijiao_xxxxx',
    appSecret: '********',
    notifyUrl: 'https://api.example.com/notify/hupijiao',
    status: 'active',
    todayAmount: 3200.0,
    todayCount: 15,
  },
  {
    id: 6,
    name: '易支付',
    type: 'epay',
    feeRate: 1.5,
    appId: 'epay_xxxxx',
    appSecret: '********',
    notifyUrl: 'https://api.example.com/notify/epay',
    status: 'active',
    todayAmount: 1800.0,
    todayCount: 12,
  },
  {
    id: 7,
    name: 'USDT(TRC20)',
    type: 'usdt',
    feeRate: 0,
    appId: 'TRX_xxxxxxxxxxxx',
    appSecret: '********',
    notifyUrl: 'https://api.example.com/notify/usdt',
    status: 'active',
    todayAmount: 5000.0,
    todayCount: 5,
  },
  {
    id: 8,
    name: '余额支付',
    type: 'balance',
    feeRate: 0,
    appId: '-',
    appSecret: '-',
    notifyUrl: '-',
    status: 'active',
    todayAmount: 2100.0,
    todayCount: 28,
  },
])

// ---- Payment Records ----
const paymentRecords = ref([
  { id: 'PAY20260727001', amount: 299.0, method: '支付宝', status: 'success', user: 'user001', time: '2026-07-27 09:15:32' },
  { id: 'PAY20260727002', amount: 599.0, method: '微信支付', status: 'success', user: 'user002', time: '2026-07-27 09:20:15' },
  { id: 'PAY20260727003', amount: 1299.0, method: 'PayPal', status: 'success', user: 'user003', time: '2026-07-27 09:25:48' },
  { id: 'PAY20260727004', amount: 89.0, method: '余额支付', status: 'success', user: 'user004', time: '2026-07-27 09:30:22' },
  { id: 'PAY20260727005', amount: 399.0, method: '虎皮椒', status: 'pending', user: 'user005', time: '2026-07-27 09:35:10' },
  { id: 'PAY20260727006', amount: 69.0, method: 'USDT', status: 'success', user: 'user006', time: '2026-07-27 09:40:55' },
  { id: 'PAY20260727007', amount: 199.0, method: '易支付', status: 'failed', user: 'user007', time: '2026-07-27 09:45:30' },
  { id: 'PAY20260727008', amount: 1299.0, method: '支付宝', status: 'success', user: 'user008', time: '2026-07-27 09:50:18' },
])

const filteredPayments = computed(() => {
  let result = payments.value

  if (searchKeyword.value.trim()) {
    const kw = searchKeyword.value.trim().toLowerCase()
    result = result.filter((p) => p.name.toLowerCase().includes(kw))
  }

  if (filterType.value) {
    result = result.filter((p) => p.type === filterType.value)
  }

  if (filterStatus.value) {
    result = result.filter((p) => p.status === filterStatus.value)
  }

  return result
})

// ---- Payment Form ----
const paymentForm = reactive({
  name: '',
  type: null as string | null,
  feeRate: 0,
  appId: '',
  appSecret: '',
  notifyUrl: '',
  status: 'active',
})

const paymentRules: FormRules = {
  name: { required: true, message: '请输入支付名称', trigger: 'blur' },
  type: { required: true, message: '请选择支付类型', trigger: 'change' },
  feeRate: { required: true, type: 'number', message: '请输入手续费率', trigger: 'blur' },
}

// ---- Payment Table Columns ----
const paymentColumns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 60, sorter: (a, b) => a.id - b.id },
  { title: '支付名称', key: 'name', ellipsis: { tooltip: true } },
  {
    title: '类型',
    key: 'type',
    width: 100,
    render: (row) =>
      h(
        NTag,
        {
          size: 'small',
          round: true,
          bordered: false,
          type: typeTagMap[row.type] as any || 'default',
        },
        { default: () => typeNameMap[row.type] || row.type }
      ),
  },
  {
    title: '手续费率',
    key: 'feeRate',
    width: 100,
    render: (row) => h('span', { style: 'font-weight:600;color:#1890ff' }, row.feeRate > 0 ? `${row.feeRate}%` : '免费'),
  },
  {
    title: '今日收款',
    key: 'todayAmount',
    width: 120,
    sorter: (a, b) => a.todayAmount - b.todayAmount,
    render: (row) => h('span', { style: 'font-weight:600;color:#52c41a' }, `¥${row.todayAmount.toFixed(2)}`),
  },
  {
    title: '今日笔数',
    key: 'todayCount',
    width: 100,
    sorter: (a, b) => a.todayCount - b.todayCount,
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
    width: 180,
    render: (row) =>
      h(NSpace, { size: 4 }, {
        default: () => [
          h(NTooltip, {}, {
            trigger: () =>
              h(
                NButton,
                { size: 'small', quaternary: true, type: 'primary', onClick: () => openPaymentModal(row) },
                { icon: () => h(NIcon, null, { default: () => h(EditIcon) }) }
              ),
            default: () => '编辑',
          }),
          h(NTooltip, {}, {
            trigger: () =>
              h(
                NButton,
                { size: 'small', quaternary: true, type: 'info', onClick: () => openRecordsModal(row) },
                { icon: () => h(NIcon, null, { default: () => h(ListIcon) }) }
              ),
            default: () => '支付记录',
          }),
          h(NPopconfirm, { onPositiveClick: () => handleTogglePayment(row) }, {
            trigger: () =>
              h(NTooltip, {}, {
                trigger: () =>
                  h(
                    NButton,
                    { size: 'small', quaternary: true, type: 'warning' },
                    { icon: () => h(NIcon, null, { default: () => h(DownIcon) }) }
                  ),
                default: () => (row.status === 'active' ? '禁用' : '启用'),
              }),
            default: () => `确认${row.status === 'active' ? '禁用' : '启用'}该支付方式？`,
          }),
          h(NPopconfirm, { onPositiveClick: () => handleDeletePayment(row.id) }, {
            trigger: () =>
              h(NTooltip, {}, {
                trigger: () =>
                  h(
                    NButton,
                    { size: 'small', quaternary: true, type: 'error' },
                    { icon: () => h(NIcon, null, { default: () => h(DeleteIcon) }) }
                  ),
                default: () => '删除',
              }),
            default: () => `确认删除支付方式「${row.name}」？`,
          }),
        ],
      }),
  },
]

// ---- Record Table Columns ----
const recordColumns: DataTableColumns<any> = [
  { title: '订单号', key: 'id', width: 160 },
  {
    title: '金额',
    key: 'amount',
    width: 100,
    render: (row) => h('span', { style: 'font-weight:600;color:#1890ff' }, `¥${row.amount.toFixed(2)}`),
  },
  { title: '支付方式', key: 'method', width: 100 },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) =>
      h(
        NTag,
        {
          size: 'small',
          round: true,
          bordered: false,
          type: row.status === 'success' ? 'success' : row.status === 'pending' ? 'warning' : 'error',
        },
        { default: () => (row.status === 'success' ? '成功' : row.status === 'pending' ? '待支付' : '失败') }
      ),
  },
  { title: '用户', key: 'user', width: 100 },
  { title: '时间', key: 'time', width: 160 },
]

// DownArrow icon for "禁用"
const DownIcon = defineComponent({
  render: () =>
    h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 512 512', fill: 'currentColor' }, [
      h('path', {
        d: 'M256 464l128-128H320V256h-32v80H128l128 128zm0-400v80h-32V144H128l128-128 128 128H320V64h-64z',
      }),
    ]),
})

// ---- Actions ----
function openPaymentModal(payment?: any) {
  editingPayment.value = payment || null
  if (payment) {
    Object.assign(paymentForm, {
      name: payment.name,
      type: payment.type,
      feeRate: payment.feeRate,
      appId: payment.appId,
      appSecret: payment.appSecret,
      notifyUrl: payment.notifyUrl,
      status: payment.status,
    })
  } else {
    Object.assign(paymentForm, {
      name: '',
      type: null,
      feeRate: 0,
      appId: '',
      appSecret: '',
      notifyUrl: '',
      status: 'active',
    })
  }
  paymentModalVisible.value = true
}

function openRecordsModal(payment: any) {
  selectedPayment.value = payment
  // In real app, fetch records for this payment method
  recordsModalVisible.value = true
}

async function handlePaymentSubmit() {
  try {
    await paymentFormRef.value?.validate()
  } catch {
    return
  }
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
</script>

<style scoped>
.n-card {
  margin: 16px;
}
</style>
