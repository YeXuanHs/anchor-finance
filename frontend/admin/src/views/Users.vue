<template>
  <div class="admin-page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">用户管理</span>
          <div class="card-actions">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索用户名 / 邮箱 / 手机"
              clearable
              style="width: 260px"
              @clear="handleSearch"
              @keydown.enter="handleSearch"
            >
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
            <el-button type="primary" @click="openModal()">
              <el-icon><Plus /></el-icon>添加用户
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="paginatedUsers" v-loading="loading" stripe size="small">
        <el-table-column prop="id" label="ID" width="60" sortable />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="email" label="邮箱" show-overflow-tooltip />
        <el-table-column prop="phone" label="手机" width="130" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-switch
              :model-value="row.status === 'active'"
              size="small"
              @change="handleToggleStatus(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="注册时间" width="160" sortable />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-tooltip content="查看" placement="top">
              <el-button text type="primary" :icon="View" @click="openDrawer(row)" />
            </el-tooltip>
            <el-tooltip content="编辑" placement="top">
              <el-button text type="primary" :icon="Edit" @click="openModal(row)" />
            </el-tooltip>
            <el-popconfirm title="确认删除该用户？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-tooltip content="删除" placement="top">
                  <el-button text type="danger" :icon="Delete" />
                </el-tooltip>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="filteredUsers.length"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @current-change="handlePageChange"
          @size-change="handlePageSizeChange"
        />
      </div>
    </el-card>

    <!-- Add/Edit Dialog -->
    <el-dialog
      v-model="modalVisible"
      :title="editingUser ? '编辑用户' : '添加用户'"
      width="520px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机" prop="phone">
          <el-input v-model="formData.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item v-if="!editingUser" label="密码" prop="password">
          <el-input v-model="formData.password" type="password" show-password placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="formData.status" active-value="active" inactive-value="disabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="modalVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- User Detail Drawer -->
    <el-drawer v-model="drawerVisible" title="用户详情" size="480px">
      <template v-if="drawerUser">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="ID">{{ drawerUser.id }}</el-descriptions-item>
          <el-descriptions-item label="用户名">{{ drawerUser.username }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ drawerUser.email }}</el-descriptions-item>
          <el-descriptions-item label="手机">{{ drawerUser.phone || '-' }}</el-descriptions-item>
          <el-descriptions-item label="余额">
            <span style="color: #0056FF; font-weight: 600">¥{{ drawerUser.balance.toFixed(2) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="drawerUser.status === 'active' ? 'success' : 'danger'" size="small">
              {{ drawerUser.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="注册时间">{{ drawerUser.createdAt }}</el-descriptions-item>
          <el-descriptions-item label="最后登录">{{ drawerUser.lastLogin || '-' }}</el-descriptions-item>
        </el-descriptions>

        <el-divider />
        <h4 style="margin-bottom: 12px">最近订单</h4>
        <el-table :data="drawerOrders" size="small" stripe>
          <el-table-column prop="id" label="订单号" width="150" />
          <el-table-column prop="product" label="产品" />
          <el-table-column prop="amount" label="金额" width="80">
            <template #default="{ row }">¥{{ row.amount }}</template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="80">
            <template #default="{ row }">
              <el-tag :type="drawerOrderStatusMap[row.status]?.type as any" size="small">
                {{ drawerOrderStatusMap[row.status]?.label }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="time" label="时间" width="150" />
        </el-table>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Plus, View, Edit, Delete } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const modalVisible = ref(false)
const drawerVisible = ref(false)
const searchKeyword = ref('')
const formRef = ref<FormInstance>()
const editingUser = ref<any>(null)
const drawerUser = ref<any>(null)

const formData = reactive({
  username: '',
  email: '',
  phone: '',
  password: '',
  status: 'active',
})

const rules: FormRules = {
  username: { required: true, message: '请输入用户名', trigger: 'blur' },
  email: { required: true, message: '请输入邮箱', trigger: 'blur' },
  password: { required: true, message: '请输入密码', trigger: 'blur' },
}

const pagination = reactive({ page: 1, pageSize: 10 })

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

const filteredUsers = computed(() => {
  if (!searchKeyword.value.trim()) return users.value
  const kw = searchKeyword.value.trim().toLowerCase()
  return users.value.filter(
    (u) => u.username.toLowerCase().includes(kw) || u.email.toLowerCase().includes(kw) || u.phone.includes(kw)
  )
})

const paginatedUsers = computed(() => {
  const start = (pagination.page - 1) * pagination.pageSize
  return filteredUsers.value.slice(start, start + pagination.pageSize)
})

const drawerOrders = ref([
  { id: 'AF20260726001', product: '基础版主机', amount: 299, status: 'active', time: '2026-07-26 14:30' },
  { id: 'AF20260715003', product: '4核8G云服务器', amount: 399, status: 'paid', time: '2026-07-15 09:20' },
])

const drawerOrderStatusMap: Record<string, { label: string; type: string }> = {
  pending: { label: '待支付', type: 'warning' },
  paid: { label: '已支付', type: 'info' },
  active: { label: '已开通', type: 'success' },
  cancelled: { label: '已取消', type: 'info' },
}

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
  if (!formRef.value) return
  try {
    await formRef.value.validate()
  } catch {
    return
  }
  submitting.value = true
  try {
    ElMessage.success(editingUser.value ? '用户更新成功' : '用户添加成功')
    modalVisible.value = false
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

function handleToggleStatus(user: any) {
  user.status = user.status === 'active' ? 'disabled' : 'active'
  ElMessage.success(`用户「${user.username}」已${user.status === 'active' ? '启用' : '禁用'}`)
}

function handleDelete(id: number) {
  users.value = users.value.filter((u) => u.id !== id)
  ElMessage.success('用户已删除')
}

function handleSearch() { pagination.page = 1 }
function handlePageChange(page: number) { pagination.page = page }
function handlePageSizeChange(size: number) { pagination.pageSize = size; pagination.page = 1 }
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
}

.card-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}
</style>
