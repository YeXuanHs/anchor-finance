<template>
  <div class="role-list-page">
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <div class="action-left">
          <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('role.addRole') }}</el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" :label="$t('role.roleName')" min-width="150" />
        <el-table-column prop="description" :label="$t('common.description')" min-width="200" />
        <el-table-column prop="admins_count" :label="$t('role.adminCount')" width="100" align="center" />
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="170" />
        <el-table-column :label="$t('role.operations')" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="warning" link size="small" @click="handlePermission(row)">{{ $t('role.permission') }}</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)" :disabled="row.is_system">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" @close="handleDialogClose">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item :label="$t('role.roleName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('role.enterRoleName')" />
        </el-form-item>
        <el-form-item :label="$t('common.description')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('role.enterDescription')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="permissionDialogVisible" :title="$t('role.permissionSetting')" width="600px">
      <div class="permission-header">
        <span>{{ $t('role.role') }}: {{ currentRole?.name }}</span>
        <el-button type="primary" size="small" @click="handleSelectAll">{{ $t('role.selectAll') }}</el-button>
        <el-button size="small" @click="handleDeselectAll">{{ $t('role.deselectAll') }}</el-button>
      </div>
      <el-tree ref="permissionTreeRef" :data="permissionTree" show-checkbox node-key="id" :default-checked-keys="checkedPermissions" :props="{ label: 'title', children: 'children' }" />
      <template #footer>
        <el-button @click="permissionDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="savingPermission" @click="handleSavePermission">{{ $t('common.save') }}</el-button>
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
import { $t } from '@/locales'

const loading = ref(false)
const tableData = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref($t('role.addRole'))
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)
const permissionDialogVisible = ref(false)
const currentRole = ref<any>(null)
const permissionTree = ref<any[]>([])
const checkedPermissions = ref<number[]>([])
const permissionTreeRef = ref<any>()
const savingPermission = ref(false)

const formData = reactive({ name: '', description: '' })
const rules = { name: [{ required: true, message: () => $t('role.enterRoleName'), trigger: 'blur' }] }

const fetchList = async () => {
  loading.value = true
  try { const data = await request.get({ url: '/api/admin/roles' }); tableData.value = data || [] } catch (error) { console.error('fetch roles failed:', error) } finally { loading.value = false }
}

const handleAdd = () => { isEdit.value = false; dialogTitle.value = $t('role.addRole'); editingId.value = null; formData.name = ''; formData.description = ''; dialogVisible.value = true }
const handleEdit = (row: any) => { isEdit.value = true; dialogTitle.value = $t('role.editRole'); editingId.value = row.id; formData.name = row.name; formData.description = row.description || ''; dialogVisible.value = true }

const handleDelete = async (row: any) => {
  try { await ElMessageBox.confirm($t('role.confirmDelete', { name: row.name }), $t('common.tips'), { type: 'warning' }); await request.del({ url: `/api/admin/roles/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchList() } catch (error) { if (error !== 'cancel') console.error('delete failed:', error) }
}

const handlePermission = async (row: any) => {
  currentRole.value = row
  try {
    const treeData = await request.get({ url: '/api/admin/permissions/tree' }); permissionTree.value = treeData || []
    const roleData = await request.get({ url: `/api/admin/roles/${row.id}/permissions` }); checkedPermissions.value = roleData?.permission_ids || []
  } catch (error) { console.error('fetch permissions failed:', error) }
  permissionDialogVisible.value = true
}

const handleSelectAll = () => { permissionTreeRef.value?.setCheckedKeys(getAllKeys(permissionTree.value)) }
const handleDeselectAll = () => { permissionTreeRef.value?.setCheckedKeys([]) }
const getAllKeys = (tree: any[]): number[] => { const keys: number[] = []; const traverse = (nodes: any[]) => { for (const node of nodes) { keys.push(node.id); if (node.children?.length) traverse(node.children) } }; traverse(tree); return keys }

const handleSavePermission = async () => {
  if (!currentRole.value) return
  savingPermission.value = true
  try {
    const checkedKeys = permissionTreeRef.value?.getCheckedKeys() || []
    const halfCheckedKeys = permissionTreeRef.value?.getHalfCheckedKeys() || []
    await request.post({ url: `/api/admin/roles/${currentRole.value.id}/permissions`, data: { permission_ids: [...checkedKeys, ...halfCheckedKeys] } })
    ElMessage.success($t('role.permissionSaved')); permissionDialogVisible.value = false
  } catch (error) { console.error('save permission failed:', error) } finally { savingPermission.value = false }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate(); submitting.value = true
    if (isEdit.value && editingId.value) { await request.put({ url: `/api/admin/roles/${editingId.value}`, data: formData }); ElMessage.success($t('common.updateSuccess')) }
    else { await request.post({ url: '/api/admin/roles', data: formData }); ElMessage.success($t('common.addSuccess')) }
    dialogVisible.value = false; fetchList()
  } catch (error) { if (error !== false) console.error('submit failed:', error) } finally { submitting.value = false }
}

const handleDialogClose = () => { formRef.value?.resetFields() }
onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.role-list-page { padding: 16px; }
.action-card { margin-bottom: 16px; }
.permission-header { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; }
</style>
