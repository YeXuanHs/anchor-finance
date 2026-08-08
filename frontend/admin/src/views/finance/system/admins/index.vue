<template>
  <div class="admin-list-page">
    <!-- 操作栏 -->
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <div class="action-left">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加管理员
          </el-button>
        </div>
        <div class="action-right">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索管理员"
            clearable
            style="width: 200px"
            @keyup.enter="fetchList"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button circle @click="fetchList">
            <el-icon><Refresh /></el-icon>
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="email" label="邮箱" min-width="200" />
        <el-table-column prop="role_name" label="角色" width="120">
          <template #default="{ row }">
            <el-tag type="primary" size="small">{{ row.role_name || '未分配' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ row.status === 'active' ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_login_at" label="最后登录" width="170" />
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button
              :type="row.status === 'active' ? 'warning' : 'success'"
              link
              size="small"
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 'active' ? '禁用' : '启用' }}
            </el-button>
            <el-button type="info" link size="small" @click="handleResetPassword(row)">
              重置密码
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      @close="handleDialogClose"
    >
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="请输入用户名" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="formData.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirm_password" v-if="!isEdit">
          <el-input v-model="formData.confirm_password" type="password" placeholder="请再次输入密码" show-password />
        </el-form-item>
        <el-form-item label="角色" prop="role_id">
          <el-select v-model="formData.role_id" placeholder="请选择角色" style="width: 100%">
            <el-option v-for="role in roles" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio value="active">正常</el-radio>
            <el-radio value="disabled">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const tableData = ref([])
const searchKeyword = ref('')
const roles = ref<{ id: number; name: string }[]>([])

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 弹窗
const dialogVisible = ref(false)
const dialogTitle = ref('添加管理员')
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

// 表单数据
const formData = reactive({
  username: '',
  email: '',
  password: '',
  confirm_password: '',
  role_id: null as number | null,
  status: 'active'
})

// 表单验证规则
const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (value !== formData.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  role_id: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

// 获取列表数据
const fetchList = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size
    }
    if (searchKeyword.value) params.keyword = searchKeyword.value

    const data = await request.get({ url: '/api/admin/admins', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取管理员列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取角色列表
const fetchRoles = async () => {
  try {
    const data = await request.get({ url: '/api/admin/roles' })
    roles.value = data || []
  } catch (error) {
    console.error('获取角色列表失败:', error)
  }
}

// 分页大小变化
const handleSizeChange = (size: number) => {
  pagination.page_size = size
  pagination.page = 1
  fetchList()
}

// 页码变化
const handlePageChange = (page: number) => {
  pagination.page = page
  fetchList()
}

// 添加管理员
const handleAdd = () => {
  isEdit.value = false
  dialogTitle.value = '添加管理员'
  editingId.value = null
  resetForm()
  dialogVisible.value = true
}

// 编辑管理员
const handleEdit = (row: any) => {
  isEdit.value = true
  dialogTitle.value = '编辑管理员'
  editingId.value = row.id
  Object.assign(formData, {
    username: row.username,
    email: row.email,
    password: '',
    confirm_password: '',
    role_id: row.role_id,
    status: row.status
  })
  dialogVisible.value = true
}

// 切换状态
const handleToggleStatus = async (row: any) => {
  const newStatus = row.status === 'active' ? 'disabled' : 'active'
  const statusText = newStatus === 'active' ? '启用' : '禁用'

  try {
    await ElMessageBox.confirm(`确定要${statusText}管理员 "${row.username}" 吗？`, '确认操作', {
      type: 'warning'
    })
    await request.post({ url: `/api/admin/admins/${row.id}/status`, data: { status: newStatus } })
    ElMessage.success(`${statusText}成功`)
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(`${statusText}失败:`, error)
    }
  }
}

// 重置密码
const handleResetPassword = async (row: any) => {
  try {
    const { value: newPassword } = await ElMessageBox.prompt(
      `请输入管理员 "${row.username}" 的新密码`,
      '重置密码',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputType: 'password',
        inputValidator: (value) => {
          if (!value || value.length < 6) {
            return '密码长度不能少于 6 个字符'
          }
          return true
        }
      }
    )

    await request.post({ url: `/api/admin/admins/${row.id}/reset-password`, data: { password: newPassword } })
    ElMessage.success('密码重置成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重置密码失败:', error)
    }
  }
}

// 删除管理员
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除管理员 "${row.username}" 吗？此操作不可恢复。`, '确认删除', {
      type: 'warning'
    })
    await request.delete({ url: `/api/admin/admins/${row.id}` })
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    if (isEdit.value && editingId.value) {
      await request.put({ url: `/api/admin/admins/${editingId.value}`, data: formData })
      ElMessage.success('更新成功')
    } else {
      await request.post({ url: '/api/admin/admins', data: formData })
      ElMessage.success('添加成功')
    }

    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    submitting.value = false
  }
}

// 重置表单
const resetForm = () => {
  formData.username = ''
  formData.email = ''
  formData.password = ''
  formData.confirm_password = ''
  formData.role_id = null
  formData.status = 'active'
}

// 弹窗关闭
const handleDialogClose = () => {
  formRef.value?.resetFields()
}

onMounted(() => {
  fetchList()
  fetchRoles()
})
</script>

<style scoped lang="scss">
.admin-list-page {
  padding: 16px;
}

.action-card {
  margin-bottom: 16px;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.action-right {
  display: flex;
  gap: 8px;
}

.table-card {
  :deep(.el-card__body) {
    padding: 0;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px;
}
</style>
