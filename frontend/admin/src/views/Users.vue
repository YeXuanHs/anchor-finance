<template>
  <div>
    <n-card title="用户管理" :bordered="false" rounded>
      <template #header-extra>
        <n-space>
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索用户名 / 邮箱 / 手机"
            clearable
            style="width: 280px"
            @clear="handleSearch"
            @keydown.enter="handleSearch"
          >
            <template #prefix>
              <n-icon><SearchIcon /></n-icon>
            </template>
          </n-input>
          <n-button type="primary" @click="openModal()">
            <template #icon><n-icon><AddIcon /></n-icon></template>
            添加用户
          </n-button>
        </n-space>
      </template>

      <n-data-table
        :columns="columns"
        :data="paginatedUsers"
        :loading="loading"
        :bordered="false"
        :row-key="(row: any) => row.id"
        size="small"
        style="margin-top: 4px"
      />

      <div class="pagination-wrap">
        <n-pagination
          v-model:page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :item-count="filteredUsers.length"
          :page-sizes="pagination.pageSizes"
          show-size-picker
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </n-card>

    <!-- Add/Edit Modal -->
    <n-modal
      v-model:show="modalVisible"
      preset="card"
      :title="editingUser ? '编辑用户' : '添加用户'"
      style="width: 520px"
      :bordered="false"
      :segmented="{ content: true, footer: true }"
    >
      <n-form ref="formRef" :model="formData" :rules="rules" label-placement="left" label-width="80">
        <n-form-item label="用户名" path="username">
          <n-input v-model:value="formData.username" placeholder="请输入用户名" />
        </n-form-item>
        <n-form-item label="邮箱" path="email">
          <n-input v-model:value="formData.email" placeholder="请输入邮箱" />
        </n-form-item>
        <n-form-item label="手机" path="phone">
          <n-input v-model:value="formData.phone" placeholder="请输入手机号" />
        </n-form-item>
        <n-form-item v-if="!editingUser" label="密码" path="password">
          <n-input v-model:value="formData.password" type="password" show-password-on="click" placeholder="请输入密码" />
        </n-form-item>
        <n-form-item label="状态" path="status">
          <n-switch v-model:value="formData.status" checked-value="active" unchecked-value="disabled">
            <template #checked>启用</template>
            <template #unchecked>禁用</template>
          </n-switch>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="modalVisible = false">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleSubmit">确定</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- User Detail Drawer -->
    <n-drawer v-model:show="drawerVisible" :width="480" placement="right">
      <n-drawer-content :title="drawerUser?.username || '用户详情'" closable>
        <template v-if="drawerUser">
          <n-descriptions bordered label-placement="left" :column="1" label-style="width:100px">
            <n-descriptions-item label="ID">{{ drawerUser.id }}</n-descriptions-item>
            <n-descriptions-item label="用户名">{{ drawerUser.username }}</n-descriptions-item>
            <n-descriptions-item label="邮箱">{{ drawerUser.email }}</n-descriptions-item>
            <n-descriptions-item label="手机">{{ drawerUser.phone || '-' }}</n-descriptions-item>
            <n-descriptions-item label="余额">
              <n-text type="info" strong>¥{{ drawerUser.balance.toFixed(2) }}</n-text>
            </n-descriptions-item>
            <n-descriptions-item label="状态">
              <n-tag :type="drawerUser.status === 'active' ? 'success' : 'error'" size="small" round>
                {{ drawerUser.status === 'active' ? '启用' : '禁用' }}
              </n-tag>
            </n-descriptions-item>
            <n-descriptions-item label="注册时间">{{ drawerUser.createdAt }}</n-descriptions-item>
            <n-descriptions-item label="最后登录">{{ drawerUser.lastLogin || '-' }}</n-descriptions-item>
          </n-descriptions>

          <n-divider />

          <n-space vertical>
            <n-text strong>最近订单</n-text>
            <n-data-table
              :columns="drawerOrderColumns"
              :data="drawerOrders"
              :bordered="false"
              :pagination="false"
              size="small"
            />
          </n-space>
        </template>
      </n-drawer-content>
    </n-drawer>
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
  type DataTableColumns,
  type FormInst,
  type FormRules,
} from 'naive-ui'
import {
  SearchOutline as SearchIcon,
  AddOutline as AddIcon,
  EyeOutline as ViewIcon,
  CreateOutline as EditIcon,
  TrashOutline as DeleteIcon,
} from '@vicons/ionicons5'

const message = useMessage()
const loading = ref(false)
const submitting = ref(false)
const modalVisible = ref(false)
const drawerVisible = ref(false)
const searchKeyword = ref('')
const formRef = ref<FormInst | null>(null)
const editingUser = ref<any>(null)
const drawerUser = ref<any>(null)

const formData = reactive({
  username: '',
  email: '',
  phone: '',
  password: '',
  status: 'active' as string,
})

const rules: FormRules = {
  username: { required: true, message: '请输入用户名', trigger: 'blur' },
  email: { required: true, message: '请输入邮箱', trigger: 'blur' },
  password: { required: true, message: '请输入密码', trigger: 'blur' },
}

const pagination = reactive({
  page: 1,
  pageSize: 10,
  pageSizes: [10, 20, 50],
})

// ---- Mock Data ----
const users = ref([
  { id: 1, username: 'zhangsan', email: 'zhangsan@example.com', phone: '13800138001', balance: 1500.00, status: 'active', createdAt: '2026-01-15 09:30', lastLogin: '2026-07-26 08:12' },
  { id: 2, username: 'lisi', email: 'lisi@example.com', phone: '13800138002', balance: 800.50, status: 'active', createdAt: '2026-02-20 14:22', lastLogin: '2026-07-25 16:45' },
  { id: 3, username: 'wangwu', email: 'wangwu@example.com', phone: '13800138003', balance: 0, status: 'disabled', createdAt: '2026-03-10 11:05', lastLogin: '' },
  { id: 4, username: 'zhaoliu', email: 'zhaoliu@example.com', phone: '13800138004', balance: 3200.00, status: 'active', createdAt: '2026-04-05 16:18', lastLogin: '2026-07-24 10:30' },
  { id: 5, username: 'sunqi', email: 'sunqi@example.com', phone: '13800138005', balance: 450.75, status: 'active', createdAt: '2026-05-18 08:44', lastLogin: '2026-07-23 22:10' },
  { id: 6, username: 'zhouba', email: 'zhouba@example.com', phone: '13800138006', balance: 2100.00, status: 'active', createdAt: '2026-05-25 13:55', lastLogin: '2026-07-26 09:00' },
  { id: 7, username: 'wujiu', email: 'wujiu@example.com', phone: '13800138007', balance: 60.00, status: 'disabled', createdAt: '2026-06-02 10:12', lastLogin: '2026-06-15 17:30' },
  { id: 8, username: 'zhengshi', email: 'zhengshi@example.com', phone: '13800138008', balance: 5600.00, status: 'active', createdAt: '2026-06-10 07:28', lastLogin: '2026-07-26 11:20' },
  { id: 9, username: 'qianyi', email: 'qianyi@example.com', phone: '13800138009', balance: 120.00, status: 'active', createdAt: '2026-06-18 15:40', lastLogin: '2026-07-22 14:55' },
  { id: 10, username: 'fenger', email: 'fenger@example.com', phone: '13800138010', balance: 980.00, status: 'active', createdAt: '2026-07-01 09:00', lastLogin: '2026-07-25 20:15' },
  { id: 11, username: 'chensan', email: 'chensan@example.com', phone: '13800138011', balance: 340.00, status: 'active', createdAt: '2026-07-05 11:30', lastLogin: '2026-07-26 07:45' },
  { id: 12, username: 'huangsi', email: 'huangsi@example.com', phone: '13800138012', balance: 0, status: 'disabled', createdAt: '2026-07-10 14:10', lastLogin: '' },
])

// ---- Filter & Pagination ----
const filteredUsers = computed(() => {
  if (!searchKeyword.value.trim()) return users.value
  const kw = searchKeyword.value.trim().toLowerCase()
  return users.value.filter(
    (u) =>
      u.username.toLowerCase().includes(kw) ||
      u.email.toLowerCase().includes(kw) ||
      u.phone.includes(kw)
  )
})

const paginatedUsers = computed(() => {
  const start = (pagination.page - 1) * pagination.pageSize
  return filteredUsers.value.slice(start, start + pagination.pageSize)
})

// ---- Drawer Orders ----
const drawerOrders = ref([
  { id: 'AF20260726001', product: '基础版主机', amount: 299, status: 'active', time: '2026-07-26 14:30' },
  { id: 'AF20260715003', product: '4核8G云服务器', amount: 399, status: 'paid', time: '2026-07-15 09:20' },
])

const drawerOrderStatusMap: Record<string, { label: string; type: string }> = {
  pending: { label: '待支付', type: 'warning' },
  paid: { label: '已支付', type: 'info' },
  active: { label: '已开通', type: 'success' },
  cancelled: { label: '已取消', type: 'default' },
}

const drawerOrderColumns: DataTableColumns<any> = [
  { title: '订单号', key: 'id', width: 140 },
  { title: '产品', key: 'product' },
  { title: '金额', key: 'amount', width: 80, render: (row) => `¥${row.amount}` },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) => {
      const s = drawerOrderStatusMap[row.status]
      return h(NTag, { type: s.type as any, size: 'tiny', round: true, bordered: false }, { default: () => s.label })
    },
  },
  { title: '时间', key: 'time', width: 140 },
]

// ---- Table Columns ----
const columns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 60, sorter: (a, b) => a.id - b.id },
  { title: '用户名', key: 'username', width: 120 },
  { title: '邮箱', key: 'email', ellipsis: { tooltip: true } },
  { title: '手机', key: 'phone', width: 130 },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) =>
      h(NSwitch, {
        value: row.status === 'active',
        size: 'small',
        onUpdateValue: () => handleToggleStatus(row),
      }),
  },
  { title: '注册时间', key: 'createdAt', width: 150, sorter: (a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime() },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    render: (row) =>
      h(NSpace, { size: 4 }, {
        default: () => [
          h(NTooltip, {}, {
            trigger: () =>
              h(NButton, { size: 'small', quaternary: true, type: 'info', onClick: () => openDrawer(row) }, {
                icon: () => h(NIcon, null, { default: () => h(ViewIcon) }),
              }),
            default: () => '查看',
          }),
          h(NTooltip, {}, {
            trigger: () =>
              h(NButton, { size: 'small', quaternary: true, type: 'primary', onClick: () => openModal(row) }, {
                icon: () => h(NIcon, null, { default: () => h(EditIcon) }),
              }),
            default: () => '编辑',
          }),
          h(NPopconfirm, { onPositiveClick: () => handleDelete(row.id) }, {
            trigger: () =>
              h(NTooltip, {}, {
                trigger: () =>
                  h(NButton, { size: 'small', quaternary: true, type: 'error' }, {
                    icon: () => h(NIcon, null, { default: () => h(DeleteIcon) }),
                  }),
                default: () => '删除',
              }),
            default: () => `确认删除用户「${row.username}」？`,
          }),
        ],
      }),
  },
]

// ---- Actions ----
function openModal(user?: any) {
  editingUser.value = user || null
  if (user) {
    Object.assign(formData, { username: user.username, email: user.email, phone: user.phone, status: user.status, password: '' })
  } else {
    Object.assign(formData, { username: '', email: '', phone: '', password: '', status: 'active' })
  }
  modalVisible.value = true
}

function openDrawer(user: any) {
  drawerUser.value = user
  drawerVisible.value = true
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  submitting.value = true
  try {
    // TODO: API call
    message.success(editingUser.value ? '用户更新成功' : '用户添加成功')
    modalVisible.value = false
  } catch (err: any) {
    message.error(err.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

function handleToggleStatus(user: any) {
  user.status = user.status === 'active' ? 'disabled' : 'active'
  message.success(`用户「${user.username}」已${user.status === 'active' ? '启用' : '禁用'}`)
}

function handleDelete(id: number) {
  users.value = users.value.filter((u) => u.id !== id)
  message.success('用户已删除')
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

<style scoped>
.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
