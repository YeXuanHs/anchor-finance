<template>
  <div class="rbac-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="权限名称/标识" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="菜单" value="menu" />
            <el-option label="按钮" value="button" />
            <el-option label="API" value="api" />
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
        <h3>权限管理</h3>
        <el-button type="primary" @click="openAddDialog(null)">
          <el-icon><Plus /></el-icon>
          添加权限
        </el-button>
      </div>

      <el-table
        :data="filteredTree"
        style="width: 100%"
        v-loading="loading"
        row-key="id"
        :tree-props="{ children: 'children' }"
        default-expand-all
      >
        <el-table-column prop="name" label="权限名称" min-width="200" />
        <el-table-column prop="code" label="权限标识" width="200">
          <template #default="{ row }">
            <el-tag type="info" size="small" effect="plain">{{ row.code }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="typeTagMap[row.type]" size="small">{{ typeTextMap[row.type] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="160" show-overflow-tooltip />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openAddDialog(row)">添加子权限</el-button>
            <el-button type="primary" link @click="openEditDialog(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑权限' : '添加权限'" width="550px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="上级权限">
          <el-tree-select
            v-model="formData.parent_id"
            :data="permissionTreeOptions"
            placeholder="无（顶级权限）"
            clearable
            check-strictly
            :props="{ label: 'name', value: 'id', children: 'children' }"
          />
        </el-form-item>
        <el-form-item label="权限名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入权限名称" />
        </el-form-item>
        <el-form-item label="权限标识" prop="code">
          <el-input v-model="formData.code" placeholder="如: users.index, users.create" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型">
            <el-option label="菜单" value="menu" />
            <el-option label="按钮" value="button" />
            <el-option label="API" value="api" />
          </el-select>
        </el-form-item>
        <el-form-item label="路由路径" prop="path" v-if="formData.type === 'menu'">
          <el-input v-model="formData.path" placeholder="如: /users" />
        </el-form-item>
        <el-form-item label="图标" prop="icon" v-if="formData.type === 'menu'">
          <el-input v-model="formData.icon" placeholder="Element Plus 图标名" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const tableData = ref<any[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const typeTagMap: Record<string, string> = { menu: 'primary', button: 'success', api: 'warning' }
const typeTextMap: Record<string, string> = { menu: '菜单', button: '按钮', api: 'API' }

const searchForm = ref({ keyword: '', type: '' })

const formData = reactive({
  parent_id: null as number | null,
  name: '',
  code: '',
  type: 'menu' as string,
  path: '',
  icon: '',
  sort: 0,
  description: ''
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入权限名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入权限标识', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }]
}

const buildTree = (list: any[], parentId: number | null = null): any[] => {
  return list
    .filter(item => item.parent_id === parentId)
    .sort((a, b) => a.sort - b.sort)
    .map(item => ({
      ...item,
      children: buildTree(list, item.id)
    }))
}

const permissionTree = computed(() => buildTree(tableData.value))

const filteredTree = computed(() => {
  if (!searchForm.value.keyword && !searchForm.value.type) return permissionTree.value
  const filterNode = (nodes: any[]): any[] => {
    return nodes.reduce((acc: any[], node: any) => {
      const matchKeyword = !searchForm.value.keyword ||
        node.name.includes(searchForm.value.keyword) ||
        node.code.includes(searchForm.value.keyword)
      const matchType = !searchForm.value.type || node.type === searchForm.value.type
      const filteredChildren = filterNode(node.children || [])
      if ((matchKeyword && matchType) || filteredChildren.length > 0) {
        acc.push({ ...node, children: filteredChildren })
      }
      return acc
    }, [])
  }
  return filterNode(permissionTree.value)
})

const permissionTreeOptions = computed(() => {
  return [{ id: 0, name: '顶级权限', children: permissionTree.value }]
})

const handleSearch = () => fetchData()

const resetSearch = () => {
  searchForm.value = { keyword: '', type: '' }
  fetchData()
}

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/permissions')
    tableData.value = data.data || []
  } catch {} finally {
    loading.value = false
  }
}

const resetForm = () => {
  Object.assign(formData, {
    parent_id: null, name: '', code: '', type: 'menu',
    path: '', icon: '', sort: 0, description: ''
  })
}

const openAddDialog = (parent: any) => {
  isEdit.value = false
  editId.value = null
  resetForm()
  if (parent) formData.parent_id = parent.id
  dialogVisible.value = true
}

const openEditDialog = (row: any) => {
  isEdit.value = true
  editId.value = row.id
  Object.assign(formData, {
    parent_id: row.parent_id || null, name: row.name, code: row.code,
    type: row.type, path: row.path || '', icon: row.icon || '',
    sort: row.sort || 0, description: row.description || ''
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (isEdit.value) {
      await request.put(`/admin/api/v1/permissions/${editId.value}`, formData)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/permissions', formData)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch {} finally {
    submitLoading.value = false
  }
}

const handleDelete = async (row: any) => {
  if (row.children && row.children.length > 0) {
    ElMessage.warning('请先删除子权限')
    return
  }
  await ElMessageBox.confirm(`确定删除权限「${row.name}」吗？`, '提示', { type: 'warning' })
  try {
    await request.delete(`/admin/api/v1/permissions/${row.id}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch {}
}

onMounted(fetchData)
</script>

<style scoped lang="scss">
.rbac-page {
  .search-bar { margin-bottom: 16px; }
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 18px; font-weight: 600; }
  }
}
</style>
