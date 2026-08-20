<template>
  <div class="rule-middle-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('ruleMiddle.title') }}</span>
          <div>
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              {{ $t('ruleMiddle.addMenu') }}
            </el-button>
            <el-button type="success" @click="handleSaveAll" :loading="saveLoading">
              <el-icon><Check /></el-icon>
              {{ $t('ruleMiddle.saveConfig') }}
            </el-button>
          </div>
        </div>
      </template>

      <el-alert :title="$t('ruleMiddle.configDesc')" type="info" :closable="false" style="margin-bottom: 16px" />

      <el-table :data="menuTree" v-loading="loading" row-key="id" :tree-props="{ children: 'son' }" style="width: 100%" border>
        <el-table-column prop="name" :label="$t('ruleMiddle.menuName')" min-width="200" />
        <el-table-column :label="$t('ruleMiddle.addPermission')" width="200" align="center">
          <template #default="{ row }">
            <el-select v-model="row.add_role" size="small" style="width: 120px" @change="handleRoleChange(row)">
              <el-option :label="$t('ruleMiddle.all')" :value="1" />
              <el-option :label="$t('ruleMiddle.ownOnly')" :value="2" />
              <el-option :label="$t('ruleMiddle.forbidden')" :value="0" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column :label="$t('ruleMiddle.viewPermission')" width="200" align="center">
          <template #default="{ row }">
            <el-select v-model="row.cat_role" size="small" style="width: 120px" @change="handleRoleChange(row)">
              <el-option :label="$t('ruleMiddle.all')" :value="1" />
              <el-option :label="$t('ruleMiddle.ownOnly')" :value="2" />
              <el-option :label="$t('ruleMiddle.forbidden')" :value="0" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column :label="$t('ruleMiddle.deletePermission')" width="200" align="center">
          <template #default="{ row }">
            <el-select v-model="row.del_role" size="small" style="width: 120px" @change="handleRoleChange(row)">
              <el-option :label="$t('ruleMiddle.all')" :value="1" />
              <el-option :label="$t('ruleMiddle.ownOnly')" :value="2" />
              <el-option :label="$t('ruleMiddle.forbidden')" :value="0" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column :label="$t('ruleMiddle.editPermission')" width="200" align="center">
          <template #default="{ row }">
            <el-select v-model="row.edit_role" size="small" style="width: 120px" @change="handleRoleChange(row)">
              <el-option :label="$t('ruleMiddle.all')" :value="1" />
              <el-option :label="$t('ruleMiddle.ownOnly')" :value="2" />
              <el-option :label="$t('ruleMiddle.forbidden')" :value="0" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.action')" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="info" link @click="handleSelectMenus(row)">{{ $t('ruleMiddle.relatedMenus') }}</el-button>
            <el-popconfirm :title="$t('ruleMiddle.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference><el-button type="danger" link>{{ $t('common.delete') }}</el-button></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('ruleMiddle.menuName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('ruleMiddle.enterMenuName')" />
        </el-form-item>
        <el-form-item :label="$t('ruleMiddle.parentMenu')">
          <el-tree-select v-model="formData.pid" :data="menuOptions" :props="{ label: 'name', value: 'id', children: 'son' } as any" :placeholder="$t('ruleMiddle.noParentMenu')" clearable check-strictly style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('ruleMiddle.addPermission')">
          <el-select v-model="formData.add_role" style="width: 100%">
            <el-option :label="$t('ruleMiddle.all')" :value="1" />
            <el-option :label="$t('ruleMiddle.ownOnly')" :value="2" />
            <el-option :label="$t('ruleMiddle.forbidden')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ruleMiddle.viewPermission')">
          <el-select v-model="formData.cat_role" style="width: 100%">
            <el-option :label="$t('ruleMiddle.all')" :value="1" />
            <el-option :label="$t('ruleMiddle.ownOnly')" :value="2" />
            <el-option :label="$t('ruleMiddle.forbidden')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ruleMiddle.deletePermission')">
          <el-select v-model="formData.del_role" style="width: 100%">
            <el-option :label="$t('ruleMiddle.all')" :value="1" />
            <el-option :label="$t('ruleMiddle.ownOnly')" :value="2" />
            <el-option :label="$t('ruleMiddle.forbidden')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ruleMiddle.editPermission')">
          <el-select v-model="formData.edit_role" style="width: 100%">
            <el-option :label="$t('ruleMiddle.all')" :value="1" />
            <el-option :label="$t('ruleMiddle.ownOnly')" :value="2" />
            <el-option :label="$t('ruleMiddle.forbidden')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.sort')">
          <el-input-number v-model="formData.order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="menuSelectVisible" :title="$t('ruleMiddle.relatedPermissionMenu')" width="500px">
      <el-alert :title="$t('ruleMiddle.relatedPermissionDesc')" type="info" :closable="false" style="margin-bottom: 16px" />
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
        <el-button @click="menuSelectVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveMenuRelation" :loading="menuRelationLoading">{{ $t('common.confirm') }}</el-button>
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
import { $t } from '@/locales'

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
const dialogTitle = ref($t('ruleMiddle.addMenu'))
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
  name: [{ required: true, message: $t('ruleMiddle.enterMenuName'), trigger: 'blur' }]
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
  } catch { ElMessage.error($t('ruleMiddle.fetchMenuFailed')) } finally { loading.value = false }
}

const fetchNavTree = async () => {
  navLoading.value = true
  try {
    const data = await request.get({ url: '/api/admin/rule-middle/nav' })
    navTree.value = data || []
  } catch { ElMessage.error($t('ruleMiddle.fetchNavFailed')) } finally { navLoading.value = false }
}

const handleAdd = () => {
  dialogTitle.value = $t('ruleMiddle.addMenu')
  Object.assign(formData, { id: null, name: '', pid: 0, add_role: 1, cat_role: 1, del_role: 1, edit_role: 1, add_role_menu: [], cat_role_menu: [], del_role_menu: [], edit_role_menu: [], order: 1 })
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('ruleMiddle.editMenu')
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
    } catch { ElMessage.error($t('ruleMiddle.operationFailed')) } finally { submitLoading.value = false }
  })
}

const handleDelete = async (row: any) => {
  try { await request.del({ url: `/api/admin/rule-middle/menus/${row.id}` }); ElMessage.success($t('ruleMiddle.deleteSuccess')); fetchMenuList() }
  catch { ElMessage.error($t('ruleMiddle.deleteFailed')) }
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
    ElMessage.success($t('ruleMiddle.saveSuccess'))
  } catch { ElMessage.error($t('ruleMiddle.saveFailed')) } finally { saveLoading.value = false }
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
    ElMessage.success($t('ruleMiddle.relateSuccess')); menuSelectVisible.value = false; fetchMenuList()
  } catch { ElMessage.error($t('ruleMiddle.relateFailed')) } finally { menuRelationLoading.value = false }
}

onMounted(() => { fetchMenuList() })
</script>

<style scoped lang="scss">
.rule-middle-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
