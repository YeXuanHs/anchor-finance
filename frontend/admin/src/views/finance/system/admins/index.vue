<template>
  <div class="admin-list-page">
    <!-- 操作栏 -->
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <div class="action-left">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('admin.addAdmin') }}
          </el-button>
        </div>
        <div class="action-right">
          <el-input
            v-model="searchKeyword"
            :placeholder="$t('admin.searchPlaceholder')"
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
        <el-table-column prop="username" :label="$t('admin.username')" width="120" />
        <el-table-column prop="email" :label="$t('admin.email')" min-width="200" />
        <el-table-column prop="role_name" :label="$t('admin.role')" width="120">
          <template #default="{ row }">
            <el-tag type="primary" size="small">{{ row.role_name || $t('admin.noRole') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('admin.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ row.status === 'active' ? $t('admin.normal') : $t('admin.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_login_at" :label="$t('admin.lastLogin')" width="170" />
        <el-table-column prop="created_at" :label="$t('admin.createdAt')" width="170" />
        <el-table-column :label="$t('admin.operations')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              {{ $t('admin.edit') }}
            </el-button>
            <el-button
              :type="row.status === 'active' ? 'warning' : 'success'"
              link
              size="small"
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 'active' ? $t('admin.disable') : $t('admin.enable') }}
            </el-button>
            <el-button type="info" link size="small" @click="handleResetPassword(row)">
              {{ $t('admin.resetPassword') }}
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">
              {{ $t('admin.delete') }}
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
        <el-form-item :label="$t('admin.username')" prop="username">
          <el-input v-model="formData.username" :placeholder="$t('admin.enterUsername')" :disabled="isEdit" />
        </el-form-item>
        <el-form-item :label="$t('admin.email')" prop="email">
          <el-input v-model="formData.email" :placeholder="$t('admin.enterEmail')" />
        </el-form-item>
        <el-form-item :label="$t('admin.password')" prop="password" v-if="!isEdit">
          <el-input v-model="formData.password" type="password" :placeholder="$t('admin.enterPassword')" show-password />
        </el-form-item>
        <el-form-item :label="$t('admin.confirmPassword')" prop="confirm_password" v-if="!isEdit">
          <el-input v-model="formData.confirm_password" type="password" :placeholder="$t('admin.enterConfirmPassword')" show-password />
        </el-form-item>
        <el-form-item :label="$t('admin.role')" prop="role_id">
          <el-select v-model="formData.role_id" :placeholder="$t('admin.selectRole')" style="width: 100%">
            <el-option v-for="role in roles" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('admin.status')" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio value="active">{{ $t('admin.normal') }}</el-radio>
            <el-radio value="disabled">{{ $t('admin.disabled') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">{{ $t('common.confirm') }}</el-button>
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
import { $t } from '@/locales'

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
const dialogTitle = ref($t('admin.addAdmin'))
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
    { required: true, message: () => $t('admin.enterUsername'), trigger: 'blur' },
    { min: 3, max: 20, message: () => $t('admin.usernameLength'), trigger: 'blur' }
  ],
  email: [
    { required: true, message: () => $t('admin.enterEmail'), trigger: 'blur' },
    { type: 'email', message: () => $t('admin.invalidEmail'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: () => $t('admin.enterPassword'), trigger: 'blur' },
    { min: 6, message: () => $t('admin.passwordLength'), trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: () => $t('admin.enterConfirmPassword'), trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (value !== formData.password) {
          callback(new Error($t('admin.passwordMismatch')))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  role_id: [
    { required: true, message: () => $t('admin.selectRole'), trigger: 'change' }
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
    console.error('Failed to fetch admin list:', error)
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
    console.error('Failed to fetch role list:', error)
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
  dialogTitle.value = $t('admin.addAdmin')
  editingId.value = null
  resetForm()
  dialogVisible.value = true
}

// 编辑管理员
const handleEdit = (row: any) => {
  isEdit.value = true
  dialogTitle.value = $t('admin.editAdmin')
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
  const statusText = newStatus === 'active' ? $t('admin.enable') : $t('admin.disable')

  try {
    await ElMessageBox.confirm(
      $t('admin.confirmToggle', { status: statusText, username: row.username }),
      $t('common.confirmAction'),
      { type: 'warning' }
    )
    await request.post({ url: `/api/admin/admins/${row.id}/status`, data: { status: newStatus } })
    ElMessage.success($t('common.operateSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Toggle status failed:', error)
    }
  }
}

// 重置密码
const handleResetPassword = async (row: any) => {
  try {
    const { value: newPassword } = await ElMessageBox.prompt(
      $t('admin.confirmResetPassword', { username: row.username }),
      $t('admin.resetPassword'),
      {
        confirmButtonText: $t('common.confirm'),
        cancelButtonText: $t('common.cancel'),
        inputType: 'password',
        inputValidator: (value) => {
          if (!value || value.length < 6) {
            return $t('admin.passwordLength')
          }
          return true
        }
      }
    )

    await request.post({ url: `/api/admin/admins/${row.id}/reset-password`, data: { password: newPassword } })
    ElMessage.success($t('admin.resetSuccess'))
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Reset password failed:', error)
    }
  }
}

// 删除管理员
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      $t('admin.confirmDelete', { username: row.username }),
      $t('common.confirmAction'),
      { type: 'warning' }
    )
    await request.del({ url: `/api/admin/admins/${row.id}` })
    ElMessage.success($t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete failed:', error)
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
      ElMessage.success($t('common.updateSuccess'))
    } else {
      await request.post({ url: '/api/admin/admins', data: formData })
      ElMessage.success($t('common.addSuccess'))
    }

    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('Submit failed:', error)
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
