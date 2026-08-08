<template>
  <div class="department-page">
    <!-- 操作栏 -->
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <div class="action-left">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加部门
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="部门名称" min-width="150" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="admins_count" label="管理员数" width="100" align="center" />
        <el-table-column prop="tickets_count" label="工单数" width="100" align="center" />
        <el-table-column prop="sort_order" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="'active'"
              :inactive-value="'disabled'"
              @change="handleToggleStatus(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="warning" link size="small" @click="handleAssignAdmin(row)">分配管理员</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" @close="handleDialogClose">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="部门名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入部门名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio value="active">启用</el-radio>
            <el-radio value="disabled">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 分配管理员弹窗 -->
    <el-dialog v-model="assignDialogVisible" title="分配管理员" width="500px">
      <el-transfer
        v-model="selectedAdminIds"
        :data="allAdmins"
        :titles="['可选管理员', '已分配']"
        :props="{ key: 'id', label: 'username' }"
      />
      <template #footer>
        <el-button @click="assignDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="assigning" @click="handleSaveAssign">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const tableData = ref([])

// 弹窗
const dialogVisible = ref(false)
const dialogTitle = ref('添加部门')
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

// 分配管理员弹窗
const assignDialogVisible = ref(false)
const assigning = ref(false)
const currentDeptId = ref<number | null>(null)
const allAdmins = ref<any[]>([])
const selectedAdminIds = ref<number[]>([])

// 表单数据
const formData = reactive({
  name: '',
  description: '',
  sort_order: 0,
  status: 'active'
})

// 表单验证规则
const rules = {
  name: [{ required: true, message: '请输入部门名称', trigger: 'blur' }]
}

// 获取列表数据
const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/ticket-departments' })
    tableData.value = data || []
  } catch (error) {
    console.error('获取部门列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 添加部门
const handleAdd = () => {
  isEdit.value = false
  dialogTitle.value = '添加部门'
  editingId.value = null
  formData.name = ''
  formData.description = ''
  formData.sort_order = 0
  formData.status = 'active'
  dialogVisible.value = true
}

// 编辑部门
const handleEdit = (row: any) => {
  isEdit.value = true
  dialogTitle.value = '编辑部门'
  editingId.value = row.id
  formData.name = row.name
  formData.description = row.description || ''
  formData.sort_order = row.sort_order || 0
  formData.status = row.status
  dialogVisible.value = true
}

// 切换状态
const handleToggleStatus = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/ticket-departments/${row.id}/status`, data: { status: row.status } })
    ElMessage.success('状态更新成功')
  } catch (error) {
    console.error('更新状态失败:', error)
    fetchList()
  }
}

// 删除部门
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除部门 "${row.name}" 吗？`, '确认删除', { type: 'warning' })
    await request.delete({ url: `/api/admin/ticket-departments/${row.id}` })
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 分配管理员
const handleAssignAdmin = async (row: any) => {
  currentDeptId.value = row.id
  try {
    const [adminsData, deptData] = await Promise.all([
      request.get({ url: '/api/admin/admins' }),
      request.get({ url: `/api/admin/ticket-departments/${row.id}/admins` })
    ])
    allAdmins.value = adminsData?.list || adminsData || []
    selectedAdminIds.value = deptData?.admin_ids || []
  } catch (error) {
    console.error('获取管理员数据失败:', error)
  }
  assignDialogVisible.value = true
}

// 保存分配
const handleSaveAssign = async () => {
  if (!currentDeptId.value) return
  assigning.value = true
  try {
    await request.post({
      url: `/api/admin/ticket-departments/${currentDeptId.value}/admins`,
      data: { admin_ids: selectedAdminIds.value }
    })
    ElMessage.success('分配成功')
    assignDialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('分配失败:', error)
  } finally {
    assigning.value = false
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    submitting.value = true
    if (isEdit.value && editingId.value) {
      await request.put({ url: `/api/admin/ticket-departments/${editingId.value}`, data: formData })
      ElMessage.success('更新成功')
    } else {
      await request.post({ url: '/api/admin/ticket-departments', data: formData })
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

// 弹窗关闭
const handleDialogClose = () => {
  formRef.value?.resetFields()
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.department-page {
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

.table-card {
  :deep(.el-card__body) {
    padding: 0;
  }
}
</style>
