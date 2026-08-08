<template>
  <div class="role-list-page">
    <!-- 操作栏 -->
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <div class="action-left">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加角色
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="角色名称" min-width="150" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="admins_count" label="管理员数" width="100" align="center" />
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="warning" link size="small" @click="handlePermission(row)">
              权限
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)" :disabled="row.is_system">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      @close="handleDialogClose"
    >
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入角色描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 权限设置弹窗 -->
    <el-dialog
      v-model="permissionDialogVisible"
      title="权限设置"
      width="600px"
    >
      <div class="permission-header">
        <span>角色: {{ currentRole?.name }}</span>
        <el-button type="primary" size="small" @click="handleSelectAll">全选</el-button>
        <el-button size="small" @click="handleDeselectAll">取消全选</el-button>
      </div>
      <el-tree
        ref="permissionTreeRef"
        :data="permissionTree"
        show-checkbox
        node-key="id"
        :default-checked-keys="checkedPermissions"
        :props="{ label: 'title', children: 'children' }"
      />
      <template #footer>
        <el-button @click="permissionDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="savingPermission" @click="handleSavePermission">保存</el-button>
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
const dialogTitle = ref('添加角色')
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

// 权限弹窗
const permissionDialogVisible = ref(false)
const currentRole = ref<any>(null)
const permissionTree = ref<any[]>([])
const checkedPermissions = ref<number[]>([])
const permissionTreeRef = ref<any>()
const savingPermission = ref(false)

// 表单数据
const formData = reactive({
  name: '',
  description: ''
})

// 表单验证规则
const rules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }]
}

// 获取列表数据
const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/roles' })
    tableData.value = data || []
  } catch (error) {
    console.error('获取角色列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 添加角色
const handleAdd = () => {
  isEdit.value = false
  dialogTitle.value = '添加角色'
  editingId.value = null
  formData.name = ''
  formData.description = ''
  dialogVisible.value = true
}

// 编辑角色
const handleEdit = (row: any) => {
  isEdit.value = true
  dialogTitle.value = '编辑角色'
  editingId.value = row.id
  formData.name = row.name
  formData.description = row.description || ''
  dialogVisible.value = true
}

// 删除角色
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除角色 "${row.name}" 吗？`, '确认删除', {
      type: 'warning'
    })
    await request.delete({ url: `/api/admin/roles/${row.id}` })
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 权限设置
const handlePermission = async (row: any) => {
  currentRole.value = row

  try {
    // 获取权限树
    const treeData = await request.get({ url: '/api/admin/permissions/tree' })
    permissionTree.value = treeData || []

    // 获取角色已有的权限
    const roleData = await request.get({ url: `/api/admin/roles/${row.id}/permissions` })
    checkedPermissions.value = roleData?.permission_ids || []
  } catch (error) {
    console.error('获取权限数据失败:', error)
  }

  permissionDialogVisible.value = true
}

// 全选
const handleSelectAll = () => {
  const allKeys = getAllKeys(permissionTree.value)
  permissionTreeRef.value?.setCheckedKeys(allKeys)
}

// 取消全选
const handleDeselectAll = () => {
  permissionTreeRef.value?.setCheckedKeys([])
}

// 获取所有节点的 key
const getAllKeys = (tree: any[]): number[] => {
  const keys: number[] = []
  const traverse = (nodes: any[]) => {
    for (const node of nodes) {
      keys.push(node.id)
      if (node.children?.length) {
        traverse(node.children)
      }
    }
  }
  traverse(tree)
  return keys
}

// 保存权限
const handleSavePermission = async () => {
  if (!currentRole.value) return

  savingPermission.value = true
  try {
    const checkedKeys = permissionTreeRef.value?.getCheckedKeys() || []
    const halfCheckedKeys = permissionTreeRef.value?.getHalfCheckedKeys() || []
    const allKeys = [...checkedKeys, ...halfCheckedKeys]

    await request.post({
      url: `/api/admin/roles/${currentRole.value.id}/permissions`,
      data: { permission_ids: allKeys }
    })
    ElMessage.success('权限保存成功')
    permissionDialogVisible.value = false
  } catch (error) {
    console.error('保存权限失败:', error)
  } finally {
    savingPermission.value = false
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    if (isEdit.value && editingId.value) {
      await request.put({ url: `/api/admin/roles/${editingId.value}`, data: formData })
      ElMessage.success('更新成功')
    } else {
      await request.post({ url: '/api/admin/roles', data: formData })
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
.role-list-page {
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

.permission-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #EBEEF5;

  span {
    font-weight: 500;
  }
}
</style>
