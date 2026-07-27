<template>
  <div>
    <n-card :bordered="false" rounded>
      <template #header>
        <n-space align="center" justify="space-between">
          <span style="font-size: 18px; font-weight: 600">优惠券管理</span>
          <n-space>
            <n-input
              v-model:value="searchKeyword"
              placeholder="搜索优惠券名称"
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
            <n-button type="primary" @click="openCouponModal()">
              <template #icon><n-icon><AddIcon /></n-icon></template>
              创建优惠券
            </n-button>
          </n-space>
        </n-space>
      </template>

      <n-data-table
        :columns="couponColumns"
        :data="filteredCoupons"
        :loading="loading"
        :bordered="false"
        :row-key="(row: any) => row.id"
        size="small"
      />
    </n-card>

    <!-- Coupon Edit Modal -->
    <n-modal
      v-model:show="couponModalVisible"
      preset="card"
      :title="editingCoupon ? '编辑优惠券' : '创建优惠券'"
      style="width: 600px"
      :bordered="false"
      :segmented="{ content: true, footer: true }"
    >
      <n-form ref="couponFormRef" :model="couponForm" :rules="couponRules" label-placement="left" label-width="100">
        <n-form-item label="优惠券名称" path="name">
          <n-input v-model:value="couponForm.name" placeholder="请输入优惠券名称" />
        </n-form-item>
        <n-form-item label="优惠券类型" path="type">
          <n-select v-model:value="couponForm.type" :options="typeOptions" placeholder="选择优惠券类型" />
        </n-form-item>
        <n-form-item label="面值" path="value">
          <n-input-number v-model:value="couponForm.value" :min="0" :precision="2" style="width: 100%">
            <template #prefix>{{ couponForm.type === 'discount' ? '折' : '¥' }}</template>
          </n-input-number>
        </n-form-item>
        <n-form-item label="使用条件" path="minAmount">
          <n-input-number v-model:value="couponForm.minAmount" :min="0" :precision="2" style="width: 100%">
            <template #prefix>满¥</template>
            <template #suffix>可用</template>
          </n-input-number>
        </n-form-item>
        <n-form-item label="有效期" path="validity">
          <n-date-picker
            v-model:value="couponForm.validity"
            type="daterange"
            style="width: 100%"
            :shortcuts="dateShortcuts"
          />
        </n-form-item>
        <n-form-item label="发放数量" path="total">
          <n-input-number v-model:value="couponForm.total" :min="1" style="width: 100%" />
        </n-form-item>
        <n-form-item label="状态" path="status">
          <n-switch v-model:value="couponForm.status" checked-value="active" unchecked-value="inactive">
            <template #checked>启用</template>
            <template #unchecked>禁用</template>
          </n-switch>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="couponModalVisible = false">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleCouponSubmit">确定</n-button>
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
} from '@vicons/ionicons5'

const message = useMessage()
const loading = ref(false)
const submitting = ref(false)
const couponModalVisible = ref(false)
const searchKeyword = ref('')
const filterType = ref<string | null>(null)
const filterStatus = ref<string | null>(null)
const couponFormRef = ref<FormInst | null>(null)
const editingCoupon = ref<any>(null)

// ---- Options ----
const typeOptions = [
  { label: '满减券', value: 'full_reduction' },
  { label: '折扣券', value: 'discount' },
  { label: '固定金额券', value: 'fixed' },
]

const typeFilterOptions = [
  { label: '满减券', value: 'full_reduction' },
  { label: '折扣券', value: 'discount' },
  { label: '固定金额券', value: 'fixed' },
]

const statusFilterOptions = [
  { label: '启用', value: 'active' },
  { label: '禁用', value: 'inactive' },
]

const typeNameMap: Record<string, string> = {
  full_reduction: '满减券',
  discount: '折扣券',
  fixed: '固定金额券',
}

const statusNameMap: Record<string, string> = {
  active: '启用',
  inactive: '禁用',
}

const dateShortcuts = {
  '最近7天': () => {
    const end = new Date()
    const start = new Date()
    start.setTime(start.getTime() - 7 * 24 * 3600 * 1000)
    return [start.getTime(), end.getTime()]
  },
  '最近30天': () => {
    const end = new Date()
    const start = new Date()
    start.setTime(start.getTime() - 30 * 24 * 3600 * 1000)
    return [start.getTime(), end.getTime()]
  },
  '最近90天': () => {
    const end = new Date()
    const start = new Date()
    start.setTime(start.getTime() - 90 * 24 * 3600 * 1000)
    return [start.getTime(), end.getTime()]
  },
}

// ---- Coupons ----
const coupons = ref([
  {
    id: 1,
    name: '新用户专享券',
    type: 'full_reduction',
    value: 50,
    minAmount: 200,
    startDate: '2026-01-01',
    endDate: '2026-12-31',
    total: 1000,
    used: 256,
    status: 'active',
  },
  {
    id: 2,
    name: 'VIP折扣券',
    type: 'discount',
    value: 8.5,
    minAmount: 100,
    startDate: '2026-01-01',
    endDate: '2026-12-31',
    total: 500,
    used: 89,
    status: 'active',
  },
  {
    id: 3,
    name: '新年特惠券',
    type: 'fixed',
    value: 30,
    minAmount: 0,
    startDate: '2026-01-01',
    endDate: '2026-02-28',
    total: 2000,
    used: 2000,
    status: 'inactive',
  },
  {
    id: 4,
    name: '企业客户专享',
    type: 'full_reduction',
    value: 200,
    minAmount: 1000,
    startDate: '2026-03-01',
    endDate: '2026-12-31',
    total: 200,
    used: 45,
    status: 'active',
  },
  {
    id: 5,
    name: '老用户回馈券',
    type: 'discount',
    value: 9,
    minAmount: 50,
    startDate: '2026-01-01',
    endDate: '2026-12-31',
    total: 800,
    used: 234,
    status: 'active',
  },
])

const filteredCoupons = computed(() => {
  let result = coupons.value

  if (searchKeyword.value.trim()) {
    const kw = searchKeyword.value.trim().toLowerCase()
    result = result.filter((c) => c.name.toLowerCase().includes(kw))
  }

  if (filterType.value) {
    result = result.filter((c) => c.type === filterType.value)
  }

  if (filterStatus.value) {
    result = result.filter((c) => c.status === filterStatus.value)
  }

  return result
})

// ---- Coupon Form ----
const couponForm = reactive({
  name: '',
  type: null as string | null,
  value: 0,
  minAmount: 0,
  validity: null as [number, number] | null,
  total: 100,
  status: 'active',
})

const couponRules: FormRules = {
  name: { required: true, message: '请输入优惠券名称', trigger: 'blur' },
  type: { required: true, message: '请选择优惠券类型', trigger: 'change' },
  value: { required: true, type: 'number', message: '请输入面值', trigger: 'blur' },
  total: { required: true, type: 'number', message: '请输入发放数量', trigger: 'blur' },
}

// ---- Coupon Table Columns ----
const couponColumns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 60, sorter: (a, b) => a.id - b.id },
  { title: '名称', key: 'name', ellipsis: { tooltip: true } },
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
          type: row.type === 'full_reduction' ? 'info' : row.type === 'discount' ? 'warning' : 'success',
        },
        { default: () => typeNameMap[row.type] || row.type }
      ),
  },
  {
    title: '面值',
    key: 'value',
    width: 100,
    render: (row) =>
      h('span', { style: 'font-weight:600;color:#1890ff' }, row.type === 'discount' ? `${row.value}折` : `¥${row.value}`),
  },
  {
    title: '使用条件',
    key: 'minAmount',
    width: 120,
    render: (row) => (row.minAmount > 0 ? `满¥${row.minAmount}可用` : '无门槛'),
  },
  {
    title: '有效期',
    key: 'validity',
    width: 180,
    render: (row) => `${row.startDate} ~ ${row.endDate}`,
  },
  {
    title: '已使用',
    key: 'used',
    width: 100,
    render: (row) => `${row.used} / ${row.total}`,
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) =>
      h(NSwitch, {
        value: row.status === 'active',
        size: 'small',
        onUpdateValue: () => handleToggleCoupon(row),
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
              h(
                NButton,
                { size: 'small', quaternary: true, type: 'primary', onClick: () => openCouponModal(row) },
                { icon: () => h(NIcon, null, { default: () => h(EditIcon) }) }
              ),
            default: () => '编辑',
          }),
          h(NPopconfirm, { onPositiveClick: () => handleToggleCoupon(row) }, {
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
            default: () => `确认${row.status === 'active' ? '禁用' : '启用'}该优惠券？`,
          }),
          h(NPopconfirm, { onPositiveClick: () => handleDeleteCoupon(row.id) }, {
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
            default: () => `确认删除优惠券「${row.name}」？`,
          }),
        ],
      }),
  },
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
function openCouponModal(coupon?: any) {
  editingCoupon.value = coupon || null
  if (coupon) {
    Object.assign(couponForm, {
      name: coupon.name,
      type: coupon.type,
      value: coupon.value,
      minAmount: coupon.minAmount,
      validity: [new Date(coupon.startDate).getTime(), new Date(coupon.endDate).getTime()],
      total: coupon.total,
      status: coupon.status,
    })
  } else {
    Object.assign(couponForm, {
      name: '',
      type: null,
      value: 0,
      minAmount: 0,
      validity: null,
      total: 100,
      status: 'active',
    })
  }
  couponModalVisible.value = true
}

async function handleCouponSubmit() {
  try {
    await couponFormRef.value?.validate()
  } catch {
    return
  }
  submitting.value = true
  try {
    // TODO: API call
    message.success(editingCoupon.value ? '优惠券更新成功' : '优惠券创建成功')
    couponModalVisible.value = false
  } catch (err: any) {
    message.error(err.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

function handleToggleCoupon(coupon: any) {
  coupon.status = coupon.status === 'active' ? 'inactive' : 'active'
  message.success(`优惠券「${coupon.name}」已${coupon.status === 'active' ? '启用' : '禁用'}`)
}

function handleDeleteCoupon(id: number) {
  coupons.value = coupons.value.filter((c) => c.id !== id)
  message.success('优惠券已删除')
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
