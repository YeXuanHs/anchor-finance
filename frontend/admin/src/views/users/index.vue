<template>
  <div class="users-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="用户名/邮箱/手机" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="正常" value="active" />
            <el-option label="禁用" value="disabled" />
            <el-option label="待验证" value="pending" />
          </el-select>
        </el-form-item>
        <el-form-item label="分组">
          <el-select v-model="searchForm.group_id" placeholder="全部" clearable>
            <el-option v-for="g in groups" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>用户列表</h3>
        <div>
          <el-button type="danger" :disabled="selectedRows.length === 0" @click="handleBatchDisable">批量禁用</el-button>
          <el-button type="success" :disabled="selectedRows.length === 0" @click="handleBatchEnable">批量启用</el-button>
          <el-button type="primary" @click="openAddDialog">
            <el-icon><Plus /></el-icon>
            添加用户
          </el-button>
        </div>
      </div>

      <el-table :data="users" style="width: 100%" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="140" />
        <el-table-column prop="email" label="邮箱" min-width="180" show-overflow-tooltip />
        <el-table-column prop="phone" label="手机" width="130" />
        <el-table-column prop="group.name" label="分组" width="120" />
        <el-table-column prop="level.name" label="等级" width="100" />
        <el-table-column prop="balance" label="余额" width="120" align="right">
          <template #default="{ row }">
            <span class="amount">¥{{ row.balance?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : row.status === 'disabled' ? 'danger' : 'warning'" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewUser(row)">详情</el-button>
            <el-button type="primary" link @click="editUser(row)">编辑</el-button>
            <el-button type="danger" link @click="deleteUser(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchUsers"
          @current-change="fetchUsers"
        />
      </div>
    </div>

    <el-dialog v-model="showFormDialog" :title="editingItem ? '编辑用户' : '添加用户'" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="请输入用户名" :disabled="!!editingItem" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机">
          <el-input v-model="formData.phone" placeholder="请输入手机" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!editingItem">
          <el-input v-model="formData.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="分组">
          <el-select v-model="formData.group_id" placeholder="请选择分组" clearable>
            <el-option v-for="g in groups" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="余额">
          <el-input-number v-model="formData.balance" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status">
            <el-radio value="active">正常</el-radio>
            <el-radio value="disabled">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFormDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()
const loading = ref(false)
const submitLoading = ref(false)
const users = ref<any[]>([])
const groups = ref<any[]>([])
const selectedRows = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showFormDialog = ref(false)
const editingItem = ref<any>(null)
const formRef = ref<FormInstance>()

const searchForm = ref({ keyword: '', status: '', group_id: '' })
const formData = ref({ username: '', email: '', phone: '', password: '', group_id: '', balance: 0, status: 'active' })
const formRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }, { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = { active: '正常', disabled: '禁用', pending: '待验证' }
  return map[status] || status
}

const fetchUsers = async () => {
  loading.value = true
  try {
    const params = { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    const { data } = await request.get('/admin/api/v1/users', { params })
    users.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

const fetchGroups = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/users/groups', { params: { page_size: 100 } })
    groups.value = data.data?.list || []
  } catch {}
}

const handleSearch = () => { currentPage.value = 1; fetchUsers() }
const resetSearch = () => { searchForm.value = { keyword: '', status: '', group_id: '' }; handleSearch() }
const handleSelectionChange = (rows: any[]) => { selectedRows.value = rows }

const openAddDialog = () => {
  editingItem.value = null
  formData.value = { username: '', email: '', phone: '', password: '', group_id: '', balance: 0, status: 'active' }
  showFormDialog.value = true
}

const editUser = (user: any) => {
  editingItem.value = user
  formData.value = { username: user.username, email: user.email, phone: user.phone || '', password: '', group_id: user.group_id || '', balance: user.balance || 0, status: user.status }
  showFormDialog.value = true
}

const viewUser = (user: any) => { router.push(`/users/detail/${user.id}`) }

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (editingItem.value) {
      const payload = { ...formData.value }
      delete (payload as any).password
      await request.put(`/admin/api/v1/users/${editingItem.value.id}`, payload)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/users', formData.value)
      ElMessage.success('添加成功')
    }
    showFormDialog.value = false
    fetchUsers()
  } catch {
    ElMessage.error(editingItem.value ? '更新失败' : '添加失败')
  } finally {
    submitLoading.value = false
  }
}

const deleteUser = async (user: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除用户「${user.username}」吗？此操作不可恢复。`, '提示', { type: 'warning' })
    await request.delete(`/admin/api/v1/users/${user.id}`)
    ElMessage.success('删除成功')
    fetchUsers()
  } catch {}
}

const handleBatchDisable = async () => {
  try {
    await ElMessageBox.confirm(`确定要禁用选中的 ${selectedRows.value.length} 个用户吗？`, '提示', { type: 'warning' })
    const ids = selectedRows.value.map(r => r.id)
    await request.post('/admin/api/v1/users/batch-disable', { ids })
    ElMessage.success('批量禁用成功')
    fetchUsers()
  } catch {}
}

const handleBatchEnable = async () => {
  try {
    const ids = selectedRows.value.map(r => r.id)
    await request.post('/admin/api/v1/users/batch-enable', { ids })
    ElMessage.success('批量启用成功')
    fetchUsers()
  } catch {}
}

onMounted(() => {
  fetchUsers()
  fetchGroups()
})
</script>

<style scoped lang="scss">
.users-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
}
</style>
