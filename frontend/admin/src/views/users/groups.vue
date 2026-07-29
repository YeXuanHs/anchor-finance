<template>
  <div class="user-groups-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="分组名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>用户分组</h3>
        <el-button type="primary" @click="openAddDialog">
          <el-icon><Plus /></el-icon>
          添加分组
        </el-button>
      </div>

      <el-table :data="groups" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="分组名称" width="180" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="user_count" label="用户数" width="100" align="center" />
        <el-table-column prop="discount" label="折扣" width="100" align="center">
          <template #default="{ row }">
            <span>{{ row.discount ? `${row.discount}%` : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="is_default" label="默认分组" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_default ? 'success' : 'info'" size="small">
              {{ row.is_default ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editGroup(row)">编辑</el-button>
            <el-button type="primary" link @click="openPermissionDialog(row)">权限</el-button>
            <el-button type="primary" link @click="viewMembers(row)">成员</el-button>
            <el-button type="danger" link @click="deleteGroup(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchGroups"
          @current-change="fetchGroups"
        />
      </div>
    </div>

    <el-dialog v-model="showFormDialog" :title="editingItem ? '编辑分组' : '添加分组'" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="分组名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入分组名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="折扣">
          <el-input-number v-model="formData.discount" :min="0" :max="100" />
          <span style="margin-left: 8px;">%</span>
        </el-form-item>
        <el-form-item label="默认分组">
          <el-switch v-model="formData.is_default" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFormDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showPermissionDialog" :title="`分组权限 - ${currentGroup?.name}`" width="600px">
      <el-tree
        ref="permTreeRef"
        :data="permissionTree"
        show-checkbox
        node-key="id"
        :default-checked-keys="checkedPermissions"
        :props="{ label: 'label', children: 'children' }"
      />
      <template #footer>
        <el-button @click="showPermissionDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSavePermissions">保存</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="showMembersDrawer" :title="`分组成员 - ${currentGroup?.name}`" size="500px">
      <el-table :data="members" v-loading="membersLoading">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="created_at" label="注册时间" width="160" />
      </el-table>
      <div class="pagination">
        <el-pagination
          v-model:current-page="memberPage"
          v-model:page-size="memberPageSize"
          :page-sizes="[10, 20, 50]"
          :total="memberTotal"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchMembers"
          @current-change="fetchMembers"
        />
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const membersLoading = ref(false)
const groups = ref<any[]>([])
const members = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const memberPage = ref(1)
const memberPageSize = ref(20)
const memberTotal = ref(0)
const showFormDialog = ref(false)
const showPermissionDialog = ref(false)
const showMembersDrawer = ref(false)
const editingItem = ref<any>(null)
const currentGroup = ref<any>(null)
const formRef = ref<FormInstance>()
const permTreeRef = ref()
const permissionTree = ref<any[]>([])
const checkedPermissions = ref<number[]>([])

const searchForm = ref({ keyword: '' })
const formData = ref({ name: '', description: '', discount: 0, is_default: false })
const formRules = { name: [{ required: true, message: '请输入分组名称', trigger: 'blur' }] }

const fetchGroups = async () => {
  loading.value = true
  try {
    const params = { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    const { data } = await request.get('/admin/api/v1/users/groups', { params })
    groups.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取分组列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { currentPage.value = 1; fetchGroups() }
const resetSearch = () => { searchForm.value = { keyword: '' }; handleSearch() }

const openAddDialog = () => {
  editingItem.value = null
  formData.value = { name: '', description: '', discount: 0, is_default: false }
  showFormDialog.value = true
}

const editGroup = (group: any) => {
  editingItem.value = group
  formData.value = { name: group.name, description: group.description || '', discount: group.discount || 0, is_default: group.is_default || false }
  showFormDialog.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (editingItem.value) {
      await request.put(`/admin/api/v1/users/groups/${editingItem.value.id}`, formData.value)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/users/groups', formData.value)
      ElMessage.success('添加成功')
    }
    showFormDialog.value = false
    fetchGroups()
  } catch {
    ElMessage.error(editingItem.value ? '更新失败' : '添加失败')
  } finally {
    submitLoading.value = false
  }
}

const deleteGroup = async (group: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除分组「${group.name}」吗？该分组下的用户将被移至默认分组。`, '提示', { type: 'warning' })
    await request.delete(`/admin/api/v1/users/groups/${group.id}`)
    ElMessage.success('删除成功')
    fetchGroups()
  } catch {}
}

const openPermissionDialog = async (group: any) => {
  currentGroup.value = group
  try {
    const { data: treeData } = await request.get('/admin/api/v1/permissions/tree')
    permissionTree.value = treeData.data || []
    const { data: groupPerms } = await request.get(`/admin/api/v1/users/groups/${group.id}/permissions`)
    checkedPermissions.value = (groupPerms.data || []).map((p: any) => p.id)
  } catch {
    ElMessage.error('获取权限数据失败')
  }
  showPermissionDialog.value = true
}

const handleSavePermissions = async () => {
  if (!permTreeRef.value) return
  const checkedKeys = permTreeRef.value.getCheckedKeys(false)
  const halfCheckedKeys = permTreeRef.value.getHalfCheckedKeys()
  submitLoading.value = true
  try {
    await request.put(`/admin/api/v1/users/groups/${currentGroup.value.id}/permissions`, {
      permission_ids: [...checkedKeys, ...halfCheckedKeys]
    })
    ElMessage.success('权限保存成功')
    showPermissionDialog.value = false
  } catch {
    ElMessage.error('保存权限失败')
  } finally {
    submitLoading.value = false
  }
}

const viewMembers = async (group: any) => {
  currentGroup.value = group
  memberPage.value = 1
  showMembersDrawer.value = true
  fetchMembers()
}

const fetchMembers = async () => {
  membersLoading.value = true
  try {
    const { data } = await request.get(`/admin/api/v1/users/groups/${currentGroup.value.id}/members`, {
      params: { page: memberPage.value, page_size: memberPageSize.value }
    })
    members.value = data.data?.list || []
    memberTotal.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取成员列表失败')
  } finally {
    membersLoading.value = false
  }
}

onMounted(fetchGroups)
</script>

<style scoped lang="scss">
.user-groups-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
}
</style>
