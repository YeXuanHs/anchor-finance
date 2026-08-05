<template>
  <div class="rule-middle-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>规则中间件管理</span>
          <div>
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              添加菜单
            </el-button>
            <el-button type="success" @click="handleSaveAll" :loading="saveLoading">
              <el-icon><Check /></el-icon>
              保存配置
            </el-button>
          </div>
        </div>
      </template>

      <!-- 导航规则选择 -->
      <el-alert title="配置说明：设置各菜单的增删改查权限规则，控制不同角色对菜单的操作权限。" type="info" :closable="false" style="margin-bottom: 16px" />

      <el-table :data="menuTree" v-loading="loading" row-key="id" :tree-props="{ children: 'son' }" style="width: 100%" border>
        <el-table-column prop="name" label="菜单名称" min-width="200" />
        <el-table-column label="添加权限" width="200" align="center">
          <template #default="{ row }">
            <el-select v-model="row.add_role" size="small" style="width: 120px" @change="handleRoleChange(row)">
              <el-option label="全部" :value="1" />
              <el-option label="仅本人" :value="2" />
              <el-option label="禁止" :value="0" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column label="查看权限" width="200" align="center">
          <template #default="{ row }">
            <el-select v-model="row.cat_role" size="small" style="width: 120px" @change="handleRoleChange(row)">
              <el-option label="全部" :value="1" />
              <el-option label="仅本人" :value="2" />
              <el-option label="禁止" :value="0" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column label="删除权限" width="200" align="center">
          <template #default="{ row }">
            <el-select v-model="row.del_role" size="small" style="width: 120px" @change="handleRoleChange(row)">
              <el-option label="全部" :value="1" />
              <el-option label="仅本人" :value="2" />
              <el-option label="禁止" :value="0" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column label="编辑权限" width="200" align="center">
          <template #default="{ row }">
            <el-select v-model="row.edit_role" size="small" style="width: 120px" @change="handleRoleChange(row)">
              <el-option label="全部" :value="1" />
              <el-option label="仅本人" :value="2" />
              <el-option label="禁止" :value="0" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="info" link @click="handleSelectMenus(row)">关联菜单</el-button>
            <el-popconfirm title="确定删除该菜单配置吗？" @confirm="handleDelete(row)">
              <template #reference><el-button type="danger" link>删除</el-button></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="菜单名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入菜单名称" />
        </el-form-item>
        <el-form-item label="上级菜单">
          <el-tree-select v-model="formData.pid" :data="menuOptions" :props="{ label: 'name', value: 'id', children: 'son' }" placeholder="无（顶级菜单）" clearable check-strictly style="width: 100%" />
        </el-form-item>
        <el-form-item label="添加权限">
          <el-select v-model="formData.add_role" style="width: 100%">
            <el-option label="全部" :value="1" />
            <el-option label="仅本人" :value="2" />
            <el-option label="禁止" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="查看权限">
          <el-select v-model="formData.cat_role" style="width: 100%">
            <el-option label="全部" :value="1" />
            <el-option label="仅本人" :value="2" />
            <el-option label="禁止" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="删除权限">
          <el-select v-model="formData.del_role" style="width: 100%">
            <el-option label="全部" :value="1" />
            <el-option label="仅本人" :value="2" />
            <el-option label="禁止" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="编辑权限">
          <el-select v-model="formData.edit_role" style="width: 100%">
            <el-option label="全部" :value="1" />
            <el-option label="仅本人" :value="2" />
            <el-option label="禁止" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="formData.order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 关联菜单对话框 -->
    <el-dialog v-model="menuSelectVisible" title="关联权限菜单" width="500px">
      <el-alert title="选择该菜单关联的权限菜单（用于权限控制）" type="info" :closable="false" style="margin-bottom: 16px" />
      <el-tree
        ref="navTreeRef"
        :data="navTree"
        :props="{ label: 'title', children: 'son' }"
        show-checkbox
        node-key="id"
        :default-checked-keys="selectedMenuIds"
        v-loading="navLoading"
      />
      <template #footer>
        <el-button @click="menuSelectVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveMenuRelation" :loading="menuRelationLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Plus, Check } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const submitLoading = ref(false)
const saveLoading = ref(false)
const navLoading = ref(false)
const menuRelationLoading = ref(false)

const menuTree = ref<any[]>([])
const navTree = ref<any[]>([])
const selectedMenuIds = ref<number[]>([])
const currentMenuId = ref(0)

const dialogVisible = ref(false)
const dialogTitle = ref('添加菜单')
const formRef = ref<FormInstance>()
const formData = reactive({
  id: null as number | null,
  name: '',
  pid: 0,
  add_role: 1,
  cat_role: 1,
  del_role: 1,
  edit_role: 1,
  add_role_menu: [] as string[],
  cat_role_menu: [] as string[],
  del_role_menu: [] as string[],
  edit_role_menu: [] as string[],
  order: 1
})

const menuSelectVisible = ref(false)
const navTreeRef = ref<any>(null)

const formRules: FormRules = {
  name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }]
}

const menuOptions = computed(() => {
  const buildOptions = (tree: any[]): any[] => {
    return tree.map(item => ({
      id: item.id,
      name: item.name,
      son: item.son ? buildOptions(item.son) : []
    }))
  }
  return buildOptions(menuTree.value)
})

const fetchMenuList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/rule-middle/menus' })
    menuTree.value = data || []
  } catch { ElMessage.error('获取菜单列表失败') } finally { loading.value = false }
}

const fetchNavTree = async () => {
  navLoading.value = true
  try {
    const data = await request.get({ url: '/api/admin/rule-middle/nav' })
    navTree.value = data || []
  } catch { ElMessage.error('获取导航树失败') } finally { navLoading.value = false }
}

const handleAdd = () => {
  dialogTitle.value = '添加菜单'
  Object.assign(formData, { id: null, name: '', pid: 0, add_role: 1, cat_role: 1, del_role: 1, edit_role: 1, add_role_menu: [], cat_role_menu: [], del_role_menu: [], edit_role_menu: [], order: 1 })
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑菜单'
  Object.assign(formData, {
    id: row.id, name: row.name, pid: row.pid || 0,
    add_role: row.add_role ?? 1, cat_role: row.cat_role ?? 1, del_role: row.del_role ?? 1, edit_role: row.edit_role ?? 1,
    add_role_menu: Array.isArray(row.add_role_menu) ? row.add_role_menu : (row.add_role_menu ? String(row.add_role_menu).split(',') : []),
    cat_role_menu: Array.isArray(row.cat_role_menu) ? row.cat_role_menu : (row.cat_role_menu ? String(row.cat_role_menu).split(',') : []),
    del_role_menu: Array.isArray(row.del_role_menu) ? row.del_role_menu : (row.del_role_menu ? String(row.del_role_menu).split(',') : []),
    edit_role_menu: Array.isArray(row.edit_role_menu) ? row.edit_role_menu : (row.edit_role_menu ? String(row.edit_role_menu).split(',') : []),
    order: row.order || 1
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/rule-middle/menus/${formData.id}`, data: formData, showSuccessMessage: true })
      } else {
        await request.post({ url: '/api/admin/rule-middle/menus', data: formData, showSuccessMessage: true })
      }
      dialogVisible.value = false; fetchMenuList()
    } catch { ElMessage.error('操作失败') } finally { submitLoading.value = false }
  })
}

const handleDelete = async (row: any) => {
  try { await request.del({ url: `/api/admin/rule-middle/menus/${row.id}` }); ElMessage.success('删除成功'); fetchMenuList() }
  catch { ElMessage.error('删除失败') }
}

const handleRoleChange = (row: any) => {
  // 实时标记变更，等待保存
}

const handleSaveAll = async () => {
  saveLoading.value = true
  try {
    const flattenTree = (tree: any[]): any[] => {
      let result: any[] = []
      tree.forEach(item => {
        result.push({ id: item.id, name: item.name, pid: item.pid, order: item.order, add_role: item.add_role, cat_role: item.cat_role, del_role: item.del_role, edit_role: item.edit_role, add_role_menu: item.add_role_menu, cat_role_menu: item.cat_role_menu, del_role_menu: item.del_role_menu, edit_role_menu: item.edit_role_menu })
        if (item.son?.length) result = result.concat(flattenTree(item.son))
      })
      return result
    }
    const list = flattenTree(menuTree.value)
    await request.post({ url: '/api/admin/rule-middle/menus/save-list', data: { list } })
    ElMessage.success('保存成功')
  } catch { ElMessage.error('保存失败') } finally { saveLoading.value = false }
}

const handleSelectMenus = async (row: any) => {
  currentMenuId.value = row.id
  selectedMenuIds.value = row.add_role_menu ? (Array.isArray(row.add_role_menu) ? row.add_role_menu.map(Number) : String(row.add_role_menu).split(',').map(Number).filter(Boolean)) : []
  menuSelectVisible.value = true
  if (!navTree.value.length) await fetchNavTree()
}

const handleSaveMenuRelation = async () => {
  if (!navTreeRef.value) return
  const checkedKeys = navTreeRef.value.getCheckedKeys()
  menuRelationLoading.value = true
  try {
    await request.put({ url: `/api/admin/rule-middle/menus/${currentMenuId.value}`, data: { add_role_menu: checkedKeys.join(',') } })
    ElMessage.success('关联成功'); menuSelectVisible.value = false; fetchMenuList()
  } catch { ElMessage.error('关联失败') } finally { menuRelationLoading.value = false }
}

onMounted(() => { fetchMenuList() })
</script>

<style scoped lang="scss">
.rule-middle-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
