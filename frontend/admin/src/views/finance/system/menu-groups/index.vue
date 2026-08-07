<template>
  <div class="menu-groups-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>菜单分组管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加分组
          </el-button>
        </div>
      </template>

      <!-- 统计卡片 -->
      <el-row :gutter="16" class="stats-row">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">总分组数</div>
            <div class="stat-value">{{ stats.total || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">已激活</div>
            <div class="stat-value success">{{ stats.active || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">会员中心导航</div>
            <div class="stat-value primary">{{ stats.member || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">官网导航</div>
            <div class="stat-value warning">{{ stats.website || 0 }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="分组名称" clearable />
        </el-form-item>
        <el-form-item label="分组类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="会员中心导航" value="member" />
            <el-option label="官网顶部导航" value="header" />
            <el-option label="官网底部导航" value="footer" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="已激活" :value="1" />
            <el-option label="未激活" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border row-key="id">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="分组名称" min-width="150" />
        <el-table-column prop="slug" label="标识" width="150">
          <template #default="{ row }">
            <el-tag size="small">{{ row.slug }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="分组类型" width="130" align="center">
          <template #default="{ row }">
            <el-tag :type="typeTagMap[row.type]" size="small">
              {{ typeTextMap[row.type] || row.type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="menu_count" label="菜单数" width="90" align="center">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.menu_count || 0 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="is_active" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-switch v-model="row.is_active" :active-value="1" :inactive-value="0" @change="handleToggleActive(row)" />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="success" link @click="handleManageMenus(row)">菜单管理</el-button>
            <el-popconfirm title="确定删除该菜单分组吗？" @confirm="handleDelete(row)">
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

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="分组名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入分组名称" />
        </el-form-item>
        <el-form-item label="标识" prop="slug">
          <el-input v-model="formData.slug" placeholder="英文标识，如: member-nav" />
        </el-form-item>
        <el-form-item label="分组类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择分组类型" style="width: 100%">
            <el-option label="会员中心导航" value="member" />
            <el-option label="官网顶部导航" value="header" />
            <el-option label="官网底部导航" value="footer" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入分组描述" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="排序">
              <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="激活状态">
              <el-switch v-model="formData.is_active" :active-value="1" :inactive-value="0" active-text="激活" inactive-text="未激活" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 菜单管理对话框 -->
    <el-dialog v-model="menuDialogVisible" :title="`菜单管理 - ${currentGroup.name}`" width="800px" destroy-on-close>
      <div v-loading="menuLoading">
        <div class="menu-header">
          <span>当前菜单列表</span>
          <el-button type="primary" size="small" @click="handleAddMenu">添加菜单项</el-button>
        </div>
        <el-table :data="groupMenus" style="width: 100%" border size="small" row-key="id">
          <el-table-column prop="id" label="ID" width="60" />
          <el-table-column prop="name" label="菜单名称" min-width="150" />
          <el-table-column prop="url" label="链接" min-width="200" show-overflow-tooltip />
          <el-table-column prop="icon" label="图标" width="100" />
          <el-table-column prop="sort" label="排序" width="80" align="center" />
          <el-table-column prop="is_active" label="状态" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.is_active ? 'success' : 'info'" size="small">{{ row.is_active ? '显示' : '隐藏' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="handleEditMenu(row)">编辑</el-button>
              <el-popconfirm title="确定移除该菜单项吗？" @confirm="handleRemoveMenu(row)">
                <template #reference>
                  <el-button type="danger" link size="small">移除</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <template #footer>
        <el-button @click="menuDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 添加/编辑菜单项对话框 -->
    <el-dialog v-model="menuItemDialogVisible" :title="menuItemDialogTitle" width="500px" destroy-on-close>
      <el-form :model="menuItemForm" :rules="menuItemRules" ref="menuItemFormRef" label-width="100px">
        <el-form-item label="菜单名称" prop="name">
          <el-input v-model="menuItemForm.name" placeholder="请输入菜单名称" />
        </el-form-item>
        <el-form-item label="链接" prop="url">
          <el-input v-model="menuItemForm.url" placeholder="请输入链接地址" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="menuItemForm.icon" placeholder="图标类名" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="menuItemForm.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="menuItemForm.is_active" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="menuItemDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveMenuItem" :loading="menuItemSubmitLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'MenuGroupsManage' })

const typeTextMap: Record<string, string> = { member: '会员中心', header: '官网顶部', footer: '官网底部' }
const typeTagMap: Record<string, any> = { member: 'primary', header: 'success', footer: 'warning' }

const loading = ref(false)
const submitLoading = ref(false)
const menuLoading = ref(false)
const menuItemSubmitLoading = ref(false)
const dialogVisible = ref(false)
const menuDialogVisible = ref(false)
const menuItemDialogVisible = ref(false)
const dialogTitle = ref('添加菜单分组')
const menuItemDialogTitle = ref('添加菜单项')
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

const formRules: FormRules = {
  name: [{ required: true, message: '请输入分组名称', trigger: 'blur' }],
  slug: [{ required: true, message: '请输入标识', trigger: 'blur' }],
  type: [{ required: true, message: '请选择分组类型', trigger: 'change' }]
}

const menuItemRules: FormRules = {
  name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  url: [{ required: true, message: '请输入链接地址', trigger: 'blur' }]
}

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
    ElMessage.error('获取菜单分组列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { Object.assign(searchForm, { keyword: '', type: '', status: undefined }); handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleAdd = () => {
  dialogTitle.value = '添加菜单分组'
  Object.assign(formData, { id: undefined, name: '', slug: '', type: 'member', description: '', sort: 0, is_active: 1 })
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑菜单分组'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleToggleActive = async (row: any) => {
  try {
    await request.put({ url: `/api/admin/menu-groups/${row.id}`, params: { is_active: row.is_active } })
    ElMessage.success(row.is_active ? '已激活' : '已停用')
    fetchStats()
  } catch (error) {
    row.is_active = row.is_active ? 0 : 1
    ElMessage.error('操作失败')
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/menu-groups/${row.id}` })
    ElMessage.success('删除成功')
    fetchData()
    fetchStats()
  } catch (error) {
    ElMessage.error('删除失败')
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
      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchData()
      fetchStats()
    } catch (error) {
      ElMessage.error('操作失败')
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
    ElMessage.error('获取菜单列表失败')
  } finally {
    menuLoading.value = false
  }
}

const handleAddMenu = () => {
  menuItemDialogTitle.value = '添加菜单项'
  editingMenuItemId.value = null
  Object.assign(menuItemForm, { name: '', url: '', icon: '', sort: 0, is_active: 1 })
  menuItemDialogVisible.value = true
}

const handleEditMenu = (row: any) => {
  menuItemDialogTitle.value = '编辑菜单项'
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
      ElMessage.success(editingMenuItemId.value ? '更新成功' : '添加成功')
      menuItemDialogVisible.value = false
      handleManageMenus({ id: currentGroup.id, name: currentGroup.name })
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      menuItemSubmitLoading.value = false
    }
  })
}

const handleRemoveMenu = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/menu-groups/${currentGroupId.value}/menus/${row.id}` })
    ElMessage.success('移除成功')
    handleManageMenus({ id: currentGroup.id, name: currentGroup.name })
  } catch (error) {
    ElMessage.error('移除失败')
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
