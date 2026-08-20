<template>
  <div class="rbac-pages-page">
    <el-row :gutter="20">
      <el-col :span="10">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('rbacPages.permissionTree') }}</span>
              <div>
                <el-button size="small" @click="handleExpandAll">{{ $t('rbacPages.expandAll') }}</el-button>
                <el-button size="small" @click="handleCollapseAll">{{ $t('rbacPages.collapseAll') }}</el-button>
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

      <el-col :span="14">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ isEditing ? $t('rbacPages.editPermission') : $t('rbacPages.permissionDetail') }}</span>
              <el-button type="primary" size="small" @click="handleAdd">
                <el-icon><Plus /></el-icon>
                {{ $t('rbacPages.addPermission') }}
              </el-button>
            </div>
          </template>

          <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px" v-loading="formLoading">
            <el-form-item :label="$t('rbacPages.permissionName')" prop="name">
              <el-input v-model="formData.name" :placeholder="$t('rbacPages.enterPermissionName')" />
            </el-form-item>
            <el-form-item :label="$t('rbacPages.permissionCode')" prop="code">
              <el-input v-model="formData.code" :placeholder="$t('rbacPages.enterCodePlaceholder')" />
            </el-form-item>
            <el-form-item :label="$t('rbacPages.type')" prop="type">
              <el-select v-model="formData.type" :placeholder="$t('rbacPages.selectType')" style="width: 100%">
                <el-option :label="$t('rbacPages.directory')" value="directory" />
                <el-option :label="$t('rbacPages.menu')" value="menu" />
                <el-option :label="$t('rbacPages.button')" value="button" />
                <el-option label="API" value="api" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('rbacPages.parentPermission')">
              <el-tree-select
                v-model="formData.parent_id"
                :data="permissionTree"
                :props="{ label: 'name', children: 'children', value: 'id' } as any"
                :placeholder="$t('rbacPages.noParentPermission')"
                clearable
                check-strictly
                style="width: 100%"
              />
            </el-form-item>
            <el-form-item :label="$t('rbacPages.routePath')">
              <el-input v-model="formData.path" :placeholder="$t('rbacPages.routePathPlaceholder')" />
            </el-form-item>
            <el-form-item :label="$t('rbacPages.icon')">
              <el-input v-model="formData.icon" :placeholder="$t('rbacPages.iconPlaceholder')" />
            </el-form-item>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('rbacPages.sort')">
                  <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('rbacPages.status')">
                  <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" :active-text="$t('rbacPages.enable')" :inactive-text="$t('rbacPages.disable')" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item :label="$t('rbacPages.description')">
              <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('rbacPages.enterPermissionDesc')" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
                {{ formData.id ? $t('common.edit') : $t('common.add') }}
              </el-button>
              <el-button v-if="formData.id" type="danger" @click="handleDelete">{{ $t('common.delete') }}</el-button>
              <el-button @click="handleResetForm">{{ $t('common.reset') }}</el-button>
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
import { $t } from '@/locales'

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
  name: [{ required: true, message: $t('rbacPages.enterPermissionName'), trigger: 'blur' }],
  code: [{ required: true, message: $t('rbacPages.enterPermissionCode'), trigger: 'blur' }],
  type: [{ required: true, message: $t('rbacPages.selectType'), trigger: 'change' }]
}

const fetchTree = async () => {
  try {
    const data = await request.get({ url: '/api/admin/rbac-pages' })
    permissionTree.value = data || []
  } catch (error) {
    ElMessage.error($t('rbacPages.fetchTreeFailed'))
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
      ElMessage.success(formData.id ? $t('rbacPages.updateSuccess') : $t('rbacPages.addSuccess'))
      handleResetForm()
      fetchTree()
    } catch (error) {
      ElMessage.error($t('rbacPages.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleDelete = async () => {
  if (!formData.id) return
  try {
    await ElMessageBox.confirm($t('rbacPages.confirmDelete'), $t('common.delete'), { type: 'warning' })
    await request.del({ url: `/api/admin/rbac-pages/${formData.id}` })
    ElMessage.success($t('rbacPages.deleteSuccess'))
    handleResetForm()
    fetchTree()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error($t('rbacPages.deleteFailed'))
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
