<template>
  <div class="rbac-pages-page">
    <el-row :gutter="20">
      <!-- 左侧：页面权限树 -->
      <el-col :span="10">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>页面权限树</span>
              <div>
                <el-button size="small" @click="handleExpandAll">展开全部</el-button>
                <el-button size="small" @click="handleCollapseAll">折叠全部</el-button>
              </div>
            </div>
          </template>
          <div class="tree-container">
            <el-tree
              ref="treeRef"
              :data="permissionTree"
              :props="{ label: 'name', children: 'children' }"
              node-key="id"
              :default-expand-all="isExpandAll"
              highlight-current
              @node-click="handleNodeClick"
            >
              <template #default="{ node, data }">
                <div class="tree-node">
                  <span>{{ data.name }}</span>
                  <span class="tree-node-meta">
                    <el-tag v-if="data.type" size="small" type="info">{{ data.type }}</el-tag>
                  </span>
                </div>
              </template>
            </el-tree>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧：权限详情/编辑 -->
      <el-col :span="14">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ isEditing ? '编辑权限' : '权限详情' }}</span>
              <el-button type="primary" size="small" @click="handleAdd">
                <el-icon><Plus /></el-icon>
                添加权限
              </el-button>
            </div>
          </template>

          <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px" v-loading="formLoading">
            <el-form-item label="权限名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入权限名称" />
            </el-form-item>
            <el-form-item label="权限标识" prop="code">
              <el-input v-model="formData.code" placeholder="如: user:list, order:create" />
            </el-form-item>
            <el-form-item label="类型" prop="type">
              <el-select v-model="formData.type" placeholder="请选择类型" style="width: 100%">
                <el-option label="目录" value="directory" />
                <el-option label="菜单" value="menu" />
                <el-option label="按钮" value="button" />
                <el-option label="API" value="api" />
              </el-select>
            </el-form-item>
            <el-form-item label="父级权限">
              <el-tree-select
                v-model="formData.parent_id"
                :data="permissionTree"
                :props="{ label: 'name', children: 'children', value: 'id' }"
                placeholder="无（顶级权限）"
                clearable
                check-strictly
                style="width: 100%"
              />
            </el-form-item>
            <el-form-item label="路由路径">
              <el-input v-model="formData.path" placeholder="如: /finance/users" />
            </el-form-item>
            <el-form-item label="图标">
              <el-input v-model="formData.icon" placeholder="如: ep:user" />
            </el-form-item>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="排序">
                  <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="状态">
                  <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="描述">
              <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入权限描述" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
                {{ formData.id ? '更新' : '添加' }}
              </el-button>
              <el-button v-if="formData.id" type="danger" @click="handleDelete">删除</el-button>
              <el-button @click="handleResetForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'RbacPagesManage' })

interface PermissionNode {
  id: number
  name: string
  code: string
  type: string
  children?: PermissionNode[]
}

const treeRef = ref()
const formRef = ref<FormInstance>()
const isExpandAll = ref(true)
const isEditing = ref(false)
const formLoading = ref(false)
const submitLoading = ref(false)

const permissionTree = ref<PermissionNode[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  code: '',
  type: 'menu',
  parent_id: undefined as number | undefined,
  path: '',
  icon: '',
  sort: 0,
  status: 1,
  description: ''
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入权限名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入权限标识', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }]
}

const fetchTree = async () => {
  try {
    const data = await request.get({ url: '/api/admin/rbac-pages' })
    permissionTree.value = data || []
  } catch (error) {
    ElMessage.error('获取权限树失败')
  }
}

const handleNodeClick = async (data: PermissionNode) => {
  isEditing.value = true
  formLoading.value = true
  try {
    const detail = await request.get({ url: `/api/admin/rbac-pages/${data.id}` })
    Object.assign(formData, detail)
  } catch (error) {
    Object.assign(formData, data)
  } finally {
    formLoading.value = false
  }
}

const handleAdd = () => {
  handleResetForm()
  isEditing.value = false
}

const handleResetForm = () => {
  formData.id = undefined
  formData.name = ''
  formData.code = ''
  formData.type = 'menu'
  formData.parent_id = undefined
  formData.path = ''
  formData.icon = ''
  formData.sort = 0
  formData.status = 1
  formData.description = ''
  isEditing.value = false
}

const handleExpandAll = () => {
  isExpandAll.value = true
  // Trigger re-render
  const temp = [...permissionTree.value]
  permissionTree.value = []
  setTimeout(() => { permissionTree.value = temp }, 0)
}

const handleCollapseAll = () => {
  isExpandAll.value = false
  const temp = [...permissionTree.value]
  permissionTree.value = []
  setTimeout(() => { permissionTree.value = temp }, 0)
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/rbac-pages/${formData.id}`, params: { ...formData } })
      } else {
        await request.post({ url: '/api/admin/rbac-pages', params: { ...formData } })
      }
      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      handleResetForm()
      fetchTree()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

const handleDelete = async () => {
  if (!formData.id) return
  try {
    await ElMessageBox.confirm('确定删除该权限吗？子权限将一并删除。', '删除确认', { type: 'warning' })
    await request.del({ url: `/api/admin/rbac-pages/${formData.id}` })
    ElMessage.success('删除成功')
    handleResetForm()
    fetchTree()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error('删除失败')
  }
}

onMounted(() => { fetchTree() })
</script>

<style scoped lang="scss">
.rbac-pages-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tree-container {
  max-height: 600px;
  overflow-y: auto;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 4px;
  padding: 8px;
}

.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding-right: 8px;

  .tree-node-meta {
    margin-left: 8px;
  }
}
</style>
