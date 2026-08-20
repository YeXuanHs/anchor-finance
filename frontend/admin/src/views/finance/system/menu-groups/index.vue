<template>
  <div class="menu-groups-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('menuGroups.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('menuGroups.addGroup') }}
          </el-button>
        </div>
      </template>

      <!-- 统计卡片 -->
      <el-row :gutter="16" class="stats-row">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">{{ $t('menuGroups.totalGroups') }}</div>
            <div class="stat-value">{{ stats.total || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">{{ $t('menuGroups.activated') }}</div>
            <div class="stat-value success">{{ stats.active || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">{{ $t('menuGroups.memberNav') }}</div>
            <div class="stat-value primary">{{ stats.member || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">{{ $t('menuGroups.websiteNav') }}</div>
            <div class="stat-value warning">{{ stats.website || 0 }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('menuGroups.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('menuGroups.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('menuGroups.groupType')">
          <el-select v-model="searchForm.type" :placeholder="$t('menuGroups.all')" clearable>
            <el-option :label="$t('menuGroups.memberCenterNav')" value="member" />
            <el-option :label="$t('menuGroups.websiteTopNav')" value="header" />
            <el-option :label="$t('menuGroups.websiteBottomNav')" value="footer" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('menuGroups.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('menuGroups.all')" clearable>
            <el-option :label="$t('menuGroups.activatedLabel')" :value="1" />
            <el-option :label="$t('menuGroups.notActivated')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('menuGroups.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('menuGroups.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border row-key="id">
        <el-table-column prop="id" :label="$t('menuGroups.id')" width="80" align="center" />
        <el-table-column prop="name" :label="$t('menuGroups.groupName')" min-width="150" />
        <el-table-column prop="slug" :label="$t('menuGroups.slug')" width="150">
          <template #default="{ row }">
            <el-tag size="small">{{ row.slug }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" :label="$t('menuGroups.groupTypeColumn')" width="130" align="center">
          <template #default="{ row }">
            <el-tag :type="typeTagMap[row.type]" size="small">
              {{ typeTextMap[row.type as keyof typeof typeTextMap] || row.type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" :label="$t('menuGroups.description')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="menu_count" :label="$t('menuGroups.menuCount')" width="90" align="center">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.menu_count || 0 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" :label="$t('menuGroups.sort')" width="80" align="center" />
        <el-table-column prop="is_active" :label="$t('menuGroups.statusColumn')" width="80" align="center">
          <template #default="{ row }">
            <el-switch v-model="row.is_active" :active-value="1" :inactive-value="0" @change="handleToggleActive(row)" />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('menuGroups.createTime')" width="180" />
        <el-table-column :label="$t('menuGroups.operations')" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('menuGroups.edit') }}</el-button>
            <el-button type="success" link @click="handleManageMenus(row)">{{ $t('menuGroups.menuManagement') }}</el-button>
            <el-popconfirm :title="$t('menuGroups.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('menuGroups.delete') }}</el-button>
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

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('menuGroups.groupName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('menuGroups.inputGroupName')" />
        </el-form-item>
        <el-form-item :label="$t('menuGroups.slug')" prop="slug">
          <el-input v-model="formData.slug" :placeholder="$t('menuGroups.slugPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('menuGroups.groupTypeColumn')" prop="type">
          <el-select v-model="formData.type" :placeholder="$t('menuGroups.selectGroupType')" style="width: 100%">
            <el-option :label="$t('menuGroups.memberCenterNav')" value="member" />
            <el-option :label="$t('menuGroups.websiteTopNav')" value="header" />
            <el-option :label="$t('menuGroups.websiteBottomNav')" value="footer" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('menuGroups.description')">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('menuGroups.inputDescription')" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('menuGroups.sortLabel')">
              <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('menuGroups.activeStatus')">
              <el-switch v-model="formData.is_active" :active-value="1" :inactive-value="0" :active-text="$t('menuGroups.activate')" :inactive-text="$t('menuGroups.notActivate')" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('menuGroups.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('menuGroups.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 菜单管理对话框 -->
    <el-dialog v-model="menuDialogVisible" :title="$t('menuGroups.menuManagement') + ' - ' + currentGroup.name" width="800px" destroy-on-close>
      <div v-loading="menuLoading">
        <div class="menu-header">
          <span>{{ $t('menuGroups.currentMenuList') }}</span>
          <el-button type="primary" size="small" @click="handleAddMenu">{{ $t('menuGroups.addMenuItem') }}</el-button>
        </div>
        <el-table :data="groupMenus" style="width: 100%" border size="small" row-key="id">
          <el-table-column prop="id" :label="$t('menuGroups.id')" width="60" />
          <el-table-column prop="name" :label="$t('menuGroups.menuName')" min-width="150" />
          <el-table-column prop="url" :label="$t('menuGroups.link')" min-width="200" show-overflow-tooltip />
          <el-table-column prop="icon" :label="$t('menuGroups.icon')" width="100" />
          <el-table-column prop="sort" :label="$t('menuGroups.sortColumn')" width="80" align="center" />
          <el-table-column prop="is_active" :label="$t('menuGroups.statusMenu')" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.is_active ? 'success' : 'info'" size="small">{{ row.is_active ? $t('menuGroups.show') : $t('menuGroups.hide') }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('menuGroups.operations')" width="150">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="handleEditMenu(row)">{{ $t('menuGroups.editMenu') }}</el-button>
              <el-popconfirm :title="$t('menuGroups.confirmRemoveMenu')" @confirm="handleRemoveMenu(row)">
                <template #reference>
                  <el-button type="danger" link size="small">{{ $t('menuGroups.remove') }}</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <template #footer>
        <el-button @click="menuDialogVisible = false">{{ $t('menuGroups.close') }}</el-button>
      </template>
    </el-dialog>

    <!-- 添加/编辑菜单项对话框 -->
    <el-dialog v-model="menuItemDialogVisible" :title="menuItemDialogTitle" width="500px" destroy-on-close>
      <el-form :model="menuItemForm" :rules="menuItemRules" ref="menuItemFormRef" label-width="100px">
        <el-form-item :label="$t('menuGroups.menuName')" prop="name">
          <el-input v-model="menuItemForm.name" :placeholder="$t('menuGroups.inputMenuName')" />
        </el-form-item>
        <el-form-item :label="$t('menuGroups.link')" prop="url">
          <el-input v-model="menuItemForm.url" :placeholder="$t('menuGroups.inputLink')" />
        </el-form-item>
        <el-form-item :label="$t('menuGroups.icon')">
          <el-input v-model="menuItemForm.icon" :placeholder="$t('menuGroups.iconPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('menuGroups.sortLabel')">
          <el-input-number v-model="menuItemForm.sort" :min="0" />
        </el-form-item>
        <el-form-item :label="$t('menuGroups.statusMenu')">
          <el-switch v-model="menuItemForm.is_active" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="menuItemDialogVisible = false">{{ $t('menuGroups.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveMenuItem" :loading="menuItemSubmitLoading">{{ $t('menuGroups.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

defineOptions({ name: 'MenuGroupsManage' })

const typeTextMap = computed(() => ({
  member: $t('menuGroups.memberCenterNav'),
  header: $t('menuGroups.websiteTopNav'),
  footer: $t('menuGroups.websiteBottomNav')
}))
const typeTagMap: Record<string, any> = { member: 'primary', header: 'success', footer: 'warning' }

const loading = ref(false)
const submitLoading = ref(false)
const menuLoading = ref(false)
const menuItemSubmitLoading = ref(false)
const dialogVisible = ref(false)
const menuDialogVisible = ref(false)
const menuItemDialogVisible = ref(false)
const isEditing = ref(false)
const isEditingMenu = ref(false)
const dialogTitle = computed(() => isEditing.value ? $t('menuGroups.editGroupTitle') : $t('menuGroups.addGroupTitle'))
const menuItemDialogTitle = computed(() => isEditingMenu.value ? $t('menuGroups.editMenuItemTitle') : $t('menuGroups.addMenuItemTitle'))
const formRef = ref<FormInstance>()
const menuItemFormRef = ref<FormInstance>()
const currentGroupId = ref(0)
const editingMenuItemId = ref<number | null>(null)

const stats = reactive({ total: 0, active: 0, member: 0, website: 0 })
const searchForm = reactive({ keyword: '', type: '', status: undefined as number | undefined })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const groupMenus = ref<any[]>([])
const currentGroup = reactive({ id: 0, name: '' })

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  slug: '',
  type: 'member',
  description: '',
  sort: 0,
  is_active: 1
})

const menuItemForm = reactive({
  name: '',
  url: '',
  icon: '',
  sort: 0,
  is_active: 1
})

const formRules = computed<FormRules>(() => ({
  name: [{ required: true, message: $t('menuGroups.inputGroupName'), trigger: 'blur' }],
  slug: [{ required: true, message: $t('menuGroups.inputSlug'), trigger: 'blur' }],
  type: [{ required: true, message: $t('menuGroups.selectGroupType'), trigger: 'change' }]
}))

const menuItemRules = computed<FormRules>(() => ({
  name: [{ required: true, message: $t('menuGroups.inputMenuName'), trigger: 'blur' }],
  url: [{ required: true, message: $t('menuGroups.inputLink'), trigger: 'blur' }]
}))

const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/menu-groups/stats' })
    Object.assign(stats, data)
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/menu-groups',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error($t('menuGroups.fetchListFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { Object.assign(searchForm, { keyword: '', type: '', status: undefined }); handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleAdd = () => {
  isEditing.value = false
  Object.assign(formData, { id: undefined, name: '', slug: '', type: 'member', description: '', sort: 0, is_active: 1 })
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  isEditing.value = true
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleToggleActive = async (row: any) => {
  try {
    await request.put({ url: `/api/admin/menu-groups/${row.id}`, params: { is_active: row.is_active } })
    ElMessage.success(row.is_active ? $t('menuGroups.activatedSuccess') : $t('menuGroups.deactivatedSuccess'))
    fetchStats()
  } catch (error) {
    row.is_active = row.is_active ? 0 : 1
    ElMessage.error($t('menuGroups.operationFailed'))
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/menu-groups/${row.id}` })
    ElMessage.success($t('menuGroups.deleteSuccess'))
    fetchData()
    fetchStats()
  } catch (error) {
    ElMessage.error($t('menuGroups.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/menu-groups/${formData.id}`, params: { ...formData } })
      } else {
        await request.post({ url: '/api/admin/menu-groups', params: { ...formData } })
      }
      ElMessage.success(formData.id ? $t('menuGroups.updateSuccess') : $t('menuGroups.addSuccess'))
      dialogVisible.value = false
      fetchData()
      fetchStats()
    } catch (error) {
      ElMessage.error($t('menuGroups.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleManageMenus = async (row: any) => {
  currentGroup.id = row.id
  currentGroup.name = row.name
  currentGroupId.value = row.id
  menuDialogVisible.value = true
  menuLoading.value = true
  try {
    const data = await request.get({ url: `/api/admin/menu-groups/${row.id}/menus` })
    groupMenus.value = data || []
  } catch (error) {
    ElMessage.error($t('menuGroups.fetchMenuFailed'))
  } finally {
    menuLoading.value = false
  }
}

const handleAddMenu = () => {
  isEditingMenu.value = false
  editingMenuItemId.value = null
  Object.assign(menuItemForm, { name: '', url: '', icon: '', sort: 0, is_active: 1 })
  menuItemDialogVisible.value = true
}

const handleEditMenu = (row: any) => {
  isEditingMenu.value = true
  editingMenuItemId.value = row.id
  Object.assign(menuItemForm, { name: row.name, url: row.url, icon: row.icon, sort: row.sort, is_active: row.is_active })
  menuItemDialogVisible.value = true
}

const handleSaveMenuItem = async () => {
  if (!menuItemFormRef.value) return
  await menuItemFormRef.value.validate(async (valid) => {
    if (!valid) return
    menuItemSubmitLoading.value = true
    try {
      if (editingMenuItemId.value) {
        await request.put({ url: `/api/admin/menu-groups/${currentGroupId.value}/menus/${editingMenuItemId.value}`, params: { ...menuItemForm } })
      } else {
        await request.post({ url: `/api/admin/menu-groups/${currentGroupId.value}/menus`, params: { ...menuItemForm } })
      }
      ElMessage.success(editingMenuItemId.value ? $t('menuGroups.updateSuccess') : $t('menuGroups.addSuccess'))
      menuItemDialogVisible.value = false
      handleManageMenus({ id: currentGroup.id, name: currentGroup.name })
    } catch (error) {
      ElMessage.error($t('menuGroups.operationFailed'))
    } finally {
      menuItemSubmitLoading.value = false
    }
  })
}

const handleRemoveMenu = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/menu-groups/${currentGroupId.value}/menus/${row.id}` })
    ElMessage.success($t('menuGroups.removeSuccess'))
    handleManageMenus({ id: currentGroup.id, name: currentGroup.name })
  } catch (error) {
    ElMessage.error($t('menuGroups.removeFailed'))
  }
}

onMounted(() => { fetchStats(); fetchData() })
</script>

<style scoped lang="scss">
.menu-groups-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.stats-row { margin-bottom: 20px; }
.stat-card { text-align: center; }
.stat-label { font-size: 13px; color: var(--el-text-color-secondary); margin-bottom: 8px; }
.stat-value {
  font-size: 28px; font-weight: 600; color: var(--el-text-color-primary);
  &.success { color: var(--el-color-success); }
  &.primary { color: var(--el-color-primary); }
  &.warning { color: var(--el-color-warning); }
}
.menu-header {
  display: flex; justify-content: space-between; align-items: center;
  margin-bottom: 16px; padding-bottom: 12px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}
</style>
