<template>
  <div>
    <n-card :bordered="false" rounded>
      <template #header-extra>
        <n-space style="padding: 12px 0 4px">
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
            placeholder="优惠券类型"
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
          <n-button type="primary" @click="openCouponModal()">
            <template #icon><n-icon><AddIcon /></n-icon></template>
            创建优惠券
          </n-button>
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
      style="width: 640px"
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
            <template #prefix>
              <span v-if="couponForm.type === 'fixed'">¥</span>
              <span v-else-if="couponForm.type === 'percent'">%</span>
            </template>
          </n-input-number>
        </n-form-item>
        <n-form-item label="使用条件" path="minAmount">
          <n-input-number v-model:value="couponForm.minAmount" :min="0" :precision="2" style="width: 100%">
            <template #prefix>¥</template>
            <template #suffix>满减门槛</template>
          </n-input-number>
        </n-form-item>
        <n-form-item label="发放总量" path="totalCount">
          <n-input-number v-model:value="couponForm.totalCount" :min="1" style="width: 100%" />
        </n-form-item>
        <n-form-item label="有效期" path="dateRange">
          <n-date-picker
            v-model:value="couponForm.dateRange"
            type="daterange"
            style="width: 100%"
            clearable
          />
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
  { label: '满减', value: 'full' },
  { label: '折扣', value: 'percent' },
  { label: '固定金额', value: 'fixed' },
]

const statusOptions = [
  { label: '启用', value: 'active' },
  { label: '禁用', value: 'inactive' },
  { label: '已过期', value: 'expired' },
]

const typeNameMap: Record<string, string> = {
  full: '满减',
  percent: '折扣',
  fixed: '固定金额',
}

const statusNameMap: Record<string, string> = {
  active: '启用',
  inactive: '禁用',
  expired: '已过期',
}

// ---- Coupons ----
const coupons = ref([
  { id: 1, name: '新用户专享券', type: 'fixed', value: 50, minAmount: 100, totalCount: 1000, usedCount: 356, status: 'active', startDate: '2026-01-01', endDate: '2026-12-31' },
  { id: 2, name: '满500减100', type: 'full', value: 100, minAmount: 500, totalCount: 500, usedCount: 128, status: 'active', startDate: '2026-03-01', endDate: '2026-09-30' },
  { id: 3, name: '8折优惠券', type: 'percent', value: 20, minAmount: 200, totalCount: 800, usedCount: 445, status: 'active', startDate: '2026-01-15', endDate: '2026-06-30' },
  { id: 4, name: '双十一狂欢券', type: 'fixed', value: 100, minAmount: 300, totalCount: 2000, usedCount: 2000, status: 'expired', startDate: '2025-11-01', endDate: '2025-11-30' },
  { id: 5, name: '年中大促折扣', type: 'percent', value: 15, minAmount: 150, totalCount: 600, usedCount: 89, status: 'inactive', startDate: '2026-06-01', endDate: '2026-06-30' },
  { id: 6, name: 'VIP专属券', type: 'fixed', value: 200, minAmount: 1000, totalCount: 100, usedCount: 23, status: 'active', startDate: '2026-01-01', endDate: '2026-12-31' },
])

const filteredCoupons = computed(() => {
  return coupons.value.filter((c) => {
    if (searchKeyword.value.trim()) {
      const kw = searchKeyword.value.trim().toLowerCase()
      if (!c.name.toLowerCase().includes(kw)) return false
    }
    if (filterType.value && c.type !== filterType.value) return false
    if (filterStatus.value && c.status !== filterStatus.value) return false
    return true
  })
})

// ---- Coupon Form ----
const couponForm = reactive({
  name: '',
  type: null as string | null,
  value: 0,
  minAmount: 0,
  totalCount: 100,
  dateRange: null as [number, number] | null,
  status: 'active',
})

const couponRules: FormRules = {
  name: { required: true, message: '请输入优惠券名称', trigger: 'blur' },
  type: { required: true, message: '请选择优惠券类型', trigger: 'change' },
  value: { required: true, type: 'number', message: '请输入面值', trigger: 'blur' },
  totalCount: { required: true, type: 'number', message: '请输入发放总量', trigger: 'blur' },
}

// ---- Coupon Table Columns ----
const couponColumns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 60, sorter: (a, b) => a.id - b.id },
  { title: '名称', key: 'name', ellipsis: { tooltip: true } },
  {
    title: '类型',
    key: 'type',
    width: 90,
    render: (row) => h(NTag, { size: 'small', round: true, bordered: false, type: row.type === 'fixed' ? 'success' : row.type === 'percent' ? 'warning' : 'info' }, { default: () => typeNameMap[row.type] || row.type }),
  },
  {
    title: '面值',
    key: 'value',
    width: 100,
    sorter: (a, b) => a.value - b.value,
    render: (row) => {
      const prefix = row.type === 'percent' ? '%' : '¥'
      return h('span', { style: 'font-weight:600;color:#18a058' }, `${row.type === 'percent' ? row.value : row.value}${prefix}`)
    },
  },
  {
    title: '使用条件',
    key: 'minAmount',
    width: 120,
    render: (row) => h('span', null, row.minAmount > 0 ? `满¥${row.minAmount}可用` : '无门槛'),
  },
  {
    title: '有效期',
    key: 'dateRange',
    width: 200,
    render: (row) => h('span', null, `${row.startDate} ~ ${row.endDate}`),
  },
  {
    title: '已使用',
    key: 'usedCount',
    width: 100,
    sorter: (a, b) => a.usedCount - b.usedCount,
    render: (row) => h('span', null, `${row.usedCount} / ${row.totalCount}`),
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) =>
      h(NSwitch, {
        value: row.status === 'active',
        size: 'small',
        disabled: row.status === 'expired',
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
              h(NButton, { size: 'small', quaternary: true, type: 'primary', onClick: () => openCouponModal(row) }, {
                icon: () => h(NIcon, null, { default: () => h(EditIcon) }),
              }),
            default: () => '编辑',
          }),
          h(NPopconfirm, { onPositiveClick: () => handleToggleCoupon(row) }, {
            trigger: () =>
              h(NTooltip, {}, {
                trigger: () =>
                  h(NButton, { size: 'small', quaternary: true, type: 'warning', disabled: row.status === 'expired' }, {
                    icon: () => h(NIcon, null, { default: () => h(StatusIcon) }),
                  }),
                default: () => row.status === 'active' ? '禁用' : '启用',
              }),
            default: () => `确认${row.status === 'active' ? '禁用' : '启用'}该优惠券？`,
          }),
          h(NPopconfirm, { onPositiveClick: () => handleDeleteCoupon(row.id) }, {
            trigger: () =>
              h(NTooltip, {}, {
                trigger: () =>
                  h(NButton, { size: 'small', quaternary: true, type: 'error' }, {
                    icon: () => h(NIcon, null, { default: () => h(DeleteIcon) }),
                  }),
                default: () => '删除',
              }),
            default: () => `确认删除优惠券「${row.name}」？`,
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

// ---- Actions ----
function openCouponModal(coupon?: any) {
  editingCoupon.value = coupon || null
  if (coupon) {
    Object.assign(couponForm, {
      name: coupon.name,
      type: coupon.type,
      value: coupon.value,
      minAmount: coupon.minAmount,
      totalCount: coupon.totalCount,
      dateRange: null,
      status: coupon.status,
    })
  } else {
    Object.assign(couponForm, { name: '', type: null, value: 0, minAmount: 0, totalCount: 100, dateRange: null, status: 'active' })
  }
  couponModalVisible.value = true
}

async function handleCouponSubmit() {
  try { await couponFormRef.value?.validate() } catch { return }
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
  if (coupon.status === 'expired') return
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
