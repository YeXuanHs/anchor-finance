<template>
  <div class="roles-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>角色权限管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加角色
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="角色名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="admin_count" label="管理员数" width="100" align="center" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handlePermission(row)">权限</el-button>
            <el-popconfirm title="确定删除该角色吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑角色对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入角色描述" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 权限分配对话框 -->
    <el-dialog v-model="permissionDialogVisible" title="分配权限" width="600px">
      <div v-loading="permissionLoading">
        <div class="permission-header">
          <span>角色：<strong>{{ currentRole.name }}</strong></span>
          <div>
            <el-button size="small" @click="handleCheckAll">全选</el-button>
            <el-button size="small" @click="handleUncheckAll">取消全选</el-button>
            <el-button size="small" @click="handleExpandAll">展开全部</el-button>
            <el-button size="small" @click="handleCollapseAll">折叠全部</el-button>
          </div>
        </div>
        <el-tree
          ref="treeRef"
          :data="permissionTree"
          :props="{ label: 'label', children: 'children' }"
          show-checkbox
          node-key="id"
          :default-expand-all="isExpandAll"
          class="permission-tree"
        />
      </div>
      <template #footer>
        <el-button @click="permissionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSavePermission" :loading="permissionSubmitLoading">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import type ElTree from 'element-plus/es/components/tree'
import request from '@/utils/http'

interface Role {
  id: number
  name: string
  description: string
  admin_count: number
  status: number
  created_at: string
}

interface PermissionNode {
  id: number
  label: string
  children?: PermissionNode[]
}

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('添加角色')
const formRef = ref<FormInstance>()
const treeRef = ref<InstanceType<typeof ElTree>>()
const isExpandAll = ref(true)

const searchForm = reactive({
  keyword: ''
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<Role[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  description: '',
  status: 1
})

const formRules: FormRules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ]
}

// 权限相关
const permissionDialogVisible = ref(false)
const permissionLoading = ref(false)
const permissionSubmitLoading = ref(false)
const currentRole = reactive<Role>({
  id: 0,
  name: '',
  description: '',
  admin_count: 0,
  status: 1,
  created_at: ''
})
const permissionTree = ref<PermissionNode[]>([])

// 获取角色列表
const fetchRoles = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/rbac/roles',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取角色列表失败:', error)
    ElMessage.error('获取角色列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchRoles()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  handleSearch()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = '添加角色'
  formData.id = undefined
  formData.name = ''
  formData.description = ''
  formData.status = 1
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: Role) => {
  dialogTitle.value = '编辑角色'
  formData.id = row.id
  formData.name = row.name
  formData.description = row.description
  formData.status = row.status
  dialogVisible.value = true
}

// 删除
const handleDelete = async (row: Role) => {
  try {
    await request.del({
      url: `/api/admin/rbac/roles/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchRoles()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({
          url: `/api/admin/rbac/roles/${formData.id}`,
          params: {
            name: formData.name,
            description: formData.description,
            status: formData.status
          },
          showSuccessMessage: true
        })
      } else {
        await request.post({
          url: '/api/admin/rbac/roles',
          params: {
            name: formData.name,
            description: formData.description,
            status: formData.status
          },
          showSuccessMessage: true
        })
      }

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchRoles()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

// 打开权限分配
const handlePermission = async (row: Role) => {
  Object.assign(currentRole, row)
  permissionDialogVisible.value = true
  permissionLoading.value = true

  try {
    const [treeData, checkedData] = await Promise.all([
      request.get({ url: '/api/admin/rbac/permissions' }),
      request.get({ url: `/api/admin/rbac/roles/${row.id}/permissions` })
    ])
    permissionTree.value = treeData || []
    await nextTick()
    if (treeRef.value && checkedData) {
      treeRef.value.setCheckedKeys(checkedData, false)
    }
  } catch (error) {
    console.error('获取权限数据失败:', error)
    ElMessage.error('获取权限数据失败')
  } finally {
    permissionLoading.value = false
  }
}

// 全选
const handleCheckAll = () => {
  if (!treeRef.value) return
  const allKeys = getAllKeys(permissionTree.value)
  treeRef.value.setCheckedKeys(allKeys, false)
}

// 取消全选
const handleUncheckAll = () => {
  if (!treeRef.value) return
  treeRef.value.setCheckedKeys([], false)
}

// 展开全部
const handleExpandAll = () => {
  isExpandAll.value = true
  setExpandAll(permissionTree.value, true)
}

// 折叠全部
const handleCollapseAll = () => {
  isExpandAll.value = false
  setExpandAll(permissionTree.value, false)
}

// 递归获取所有节点 key
const getAllKeys = (nodes: PermissionNode[]): number[] => {
  const keys: number[] = []
  const traverse = (list: PermissionNode[]) => {
    list.forEach((node) => {
      keys.push(node.id)
      if (node.children?.length) {
        traverse(node.children)
      }
    })
  }
  traverse(nodes)
  return keys
}

// 设置展开/折叠
const setExpandAll = (nodes: PermissionNode[], expand: boolean) => {
  nodes.forEach((node) => {
    if (node.children?.length) {
      const treeNode = treeRef.value?.store.nodesMap[node.id] as any
      if (treeNode) {
        treeNode.expanded = expand
      }
      setExpandAll(node.children, expand)
    }
  })
}

// 保存权限
const handleSavePermission = async () => {
  if (!treeRef.value) return

  permissionSubmitLoading.value = true
  try {
    const checkedKeys = treeRef.value.getCheckedKeys(false) as number[]
    const halfCheckedKeys = treeRef.value.getHalfCheckedKeys() as number[]
    const permissionIds = [...checkedKeys, ...halfCheckedKeys]

    await request.put({
      url: `/api/admin/rbac/roles/${currentRole.id}/permissions`,
      params: { permission_ids: permissionIds },
      showSuccessMessage: true
    })
    ElMessage.success('权限保存成功')
    permissionDialogVisible.value = false
  } catch (error) {
    ElMessage.error('权限保存失败')
  } finally {
    permissionSubmitLoading.value = false
  }
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchRoles()
}

// 页码变化
const handlePageChange = () => {
  fetchRoles()
}

onMounted(() => {
  fetchRoles()
})
</script>

<style scoped lang="scss">
.roles-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.permission-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.permission-tree {
  max-height: 400px;
  overflow-y: auto;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 4px;
  padding: 8px;
}
</style>
