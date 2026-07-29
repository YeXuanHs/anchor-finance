<template>
  <div class="roles-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="角色名称/标识" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>角色管理</h3>
        <el-button type="primary" @click="openAddDialog">
          <el-icon><Plus /></el-icon>
          添加角色
        </el-button>
      </div>

      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" min-width="140" />
        <el-table-column prop="code" label="角色标识" width="160">
          <template #default="{ row }">
            <el-tag type="info" size="small" effect="plain">{{ row.code }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="admin_count" label="管理员数" width="100" />
        <el-table-column prop="is_system" label="系统角色" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_system ? 'danger' : 'info'" size="small">
              {{ row.is_system ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openEditDialog(row)">编辑</el-button>
            <el-button type="primary" link @click="openPermissionDialog(row)">权限</el-button>
            <el-button type="danger" link @click="handleDelete(row)" :disabled="row.is_system">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchData"
          @size-change="fetchData"
        />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑角色' : '添加角色'" width="550px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色标识" prop="code">
          <el-input v-model="formData.code" placeholder="如: admin, editor, viewer" :disabled="isEdit && formData.is_system" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入角色描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="permDialogVisible" title="分配权限" width="550px">
      <div style="margin-bottom: 12px;">
        <el-checkbox v-model="expandAll" @change="toggleExpand">展开/折叠</el-checkbox>
        <el-checkbox v-model="checkAll" @change="toggleCheckAll">全选/取消全选</el-checkbox>
      </div>
      <el-tree
        ref="permTreeRef"
        :data="permissionTree"
        show-checkbox
        node-key="id"
        :default-expand-all="expandAll"
        :default-checked-keys="checkedPermissions"
        :props="{ label: 'name', children: 'children' }"
        style="max-height: 400px; overflow-y: auto;"
      />
      <template #footer>
        <el-button @click="permDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="permSubmitLoading" @click="savePermissions">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const formRef = ref<FormInstance>()
const permDialogVisible = ref(false)
const permSubmitLoading = ref(false)
const permissionTree = ref<any[]>([])
const checkedPermissions = ref<number[]>([])
const permTreeRef = ref()
const expandAll = ref(true)
const checkAll = ref(false)
const currentRoleId = ref<number | null>(null)

const searchForm = ref({ keyword: '' })

const formData = reactive({
  name: '',
  code: '',
  description: '',
  is_system: false
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入角色标识', trigger: 'blur' }]
}

const handleSearch = () => { page.value = 1; fetchData() }

const resetSearch = () => {
  searchForm.value = { keyword: '' }
  handleSearch()
}

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/roles', {
      params: { page: page.value, page_size: pageSize.value, ...searchForm.value }
    })
    tableData.value = data.data || []
    total.value = data.total || 0
  } catch {} finally {
    loading.value = false
  }
}

const fetchPermissions = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/permissions')
    const list = data.data || []
    const buildTree = (items: any[], parentId: number | null = null): any[] => {
      return items
        .filter((i: any) => i.parent_id === parentId)
        .map((item: any) => ({ ...item, children: buildTree(items, item.id) }))
    }
    permissionTree.value = buildTree(list)
  } catch {}
}

const resetForm = () => {
  Object.assign(formData, { name: '', code: '', description: '', is_system: false })
}

const openAddDialog = () => {
  isEdit.value = false
  editId.value = null
  resetForm()
  dialogVisible.value = true
}

const openEditDialog = (row: any) => {
  isEdit.value = true
  editId.value = row.id
  Object.assign(formData, {
    name: row.name, code: row.code, description: row.description || '', is_system: row.is_system
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (isEdit.value) {
      await request.put(`/admin/api/v1/roles/${editId.value}`, formData)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/roles', formData)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch {} finally {
    submitLoading.value = false
  }
}

const handleDelete = async (row: any) => {
  if (row.is_system) {
    ElMessage.warning('系统角色不可删除')
    return
  }
  await ElMessageBox.confirm(`确定删除角色「${row.name}」吗？`, '提示', { type: 'warning' })
  try {
    await request.delete(`/admin/api/v1/roles/${row.id}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch {}
}

const openPermissionDialog = async (row: any) => {
  currentRoleId.value = row.id
  checkAll.value = false
  try {
    const { data } = await request.get(`/admin/api/v1/roles/${row.id}/permissions`)
    checkedPermissions.value = data.permission_ids || []
  } catch {
    checkedPermissions.value = []
  }
  permDialogVisible.value = true
}

const toggleExpand = (val: boolean) => {
  // Tree re-renders with default-expand-all
}

const toggleCheckAll = (val: boolean) => {
  if (permTreeRef.value) {
    if (val) {
      const allIds = getAllNodeIds(permissionTree.value)
      permTreeRef.value.setCheckedKeys(allIds)
    } else {
      permTreeRef.value.setCheckedKeys([])
    }
  }
}

const getAllNodeIds = (nodes: any[]): number[] => {
  const ids: number[] = []
  const walk = (list: any[]) => {
    list.forEach(n => {
      ids.push(n.id)
      if (n.children?.length) walk(n.children)
    })
  }
  walk(nodes)
  return ids
}

const savePermissions = async () => {
  if (!permTreeRef.value || !currentRoleId.value) return
  permSubmitLoading.value = true
  const checkedKeys = permTreeRef.value.getCheckedKeys(false)
  const halfCheckedKeys = permTreeRef.value.getHalfCheckedKeys()
  try {
    await request.put(`/admin/api/v1/roles/${currentRoleId.value}/permissions`, {
      permission_ids: [...checkedKeys, ...halfCheckedKeys]
    })
    ElMessage.success('权限已更新')
    permDialogVisible.value = false
  } catch {} finally {
    permSubmitLoading.value = false
  }
}

onMounted(() => {
  fetchData()
  fetchPermissions()
})
</script>

<style scoped lang="scss">
.roles-page {
  .search-bar { margin-bottom: 16px; }
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 18px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
}
</style>
