<template>
  <div class="departments-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="部门名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>工单部门</h3>
        <el-button type="primary" @click="openAddDialog">
          <el-icon><Plus /></el-icon>
          添加部门
        </el-button>
      </div>

      <el-table :data="departments" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="部门名称" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="admin_count" label="管理员数" width="100" align="center" />
        <el-table-column prop="ticket_count" label="工单数" width="100" align="center" />
        <el-table-column prop="auto_assign" label="自动分配" width="100">
          <template #default="{ row }">
            <el-tag :type="row.auto_assign ? 'success' : 'info'" size="small">
              {{ row.auto_assign ? '开启' : '关闭' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="80" align="center" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editDepartment(row)">编辑</el-button>
            <el-button type="primary" link @click="openAdminDialog(row)">分配客服</el-button>
            <el-button type="danger" link @click="deleteDepartment(row)">删除</el-button>
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
          @size-change="fetchDepartments"
          @current-change="fetchDepartments"
        />
      </div>
    </div>

    <el-dialog v-model="showFormDialog" :title="editingItem ? '编辑部门' : '添加部门'" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="部门名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入部门名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="自动分配">
          <el-switch v-model="formData.auto_assign" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="formData.sort_order" :min="0" :max="9999" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFormDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showAdminDialog" :title="`分配客服 - ${currentDept?.name}`" width="600px">
      <div class="admin-transfer">
        <el-transfer
          v-model="selectedAdminIds"
          :data="allAdmins"
          :titles="['可选客服', '已分配客服']"
          :props="{ key: 'id', label: 'username' }"
        />
      </div>
      <template #footer>
        <el-button @click="showAdminDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleAssignAdmins">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const departments = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showFormDialog = ref(false)
const showAdminDialog = ref(false)
const editingItem = ref<any>(null)
const currentDept = ref<any>(null)
const formRef = ref<FormInstance>()
const allAdmins = ref<any[]>([])
const selectedAdminIds = ref<number[]>([])

const searchForm = ref({ keyword: '' })
const formData = ref({ name: '', description: '', auto_assign: false, sort_order: 0 })
const formRules = { name: [{ required: true, message: '请输入部门名称', trigger: 'blur' }] }

const fetchDepartments = async () => {
  loading.value = true
  try {
    const params = { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    const { data } = await request.get('/admin/api/v1/tickets/departments', { params })
    departments.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取部门列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { currentPage.value = 1; fetchDepartments() }
const resetSearch = () => { searchForm.value = { keyword: '' }; handleSearch() }

const openAddDialog = () => {
  editingItem.value = null
  formData.value = { name: '', description: '', auto_assign: false, sort_order: 0 }
  showFormDialog.value = true
}

const editDepartment = (dept: any) => {
  editingItem.value = dept
  formData.value = { name: dept.name, description: dept.description, auto_assign: dept.auto_assign, sort_order: dept.sort_order }
  showFormDialog.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (editingItem.value) {
      await request.put(`/admin/api/v1/tickets/departments/${editingItem.value.id}`, formData.value)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/tickets/departments', formData.value)
      ElMessage.success('添加成功')
    }
    showFormDialog.value = false
    fetchDepartments()
  } catch {
    ElMessage.error(editingItem.value ? '更新失败' : '添加失败')
  } finally {
    submitLoading.value = false
  }
}

const deleteDepartment = async (dept: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除部门「${dept.name}」吗？`, '提示', { type: 'warning' })
    await request.delete(`/admin/api/v1/tickets/departments/${dept.id}`)
    ElMessage.success('删除成功')
    fetchDepartments()
  } catch {}
}

const openAdminDialog = async (dept: any) => {
  currentDept.value = dept
  try {
    const { data: adminsData } = await request.get('/admin/api/v1/admins', { params: { page_size: 200 } })
    allAdmins.value = adminsData.data?.list || []
    const { data: deptAdmins } = await request.get(`/admin/api/v1/tickets/departments/${dept.id}/admins`)
    selectedAdminIds.value = (deptAdmins.data || []).map((a: any) => a.id)
  } catch {
    ElMessage.error('获取客服数据失败')
  }
  showAdminDialog.value = true
}

const handleAssignAdmins = async () => {
  submitLoading.value = true
  try {
    await request.put(`/admin/api/v1/tickets/departments/${currentDept.value.id}/admins`, { admin_ids: selectedAdminIds.value })
    ElMessage.success('分配成功')
    showAdminDialog.value = false
    fetchDepartments()
  } catch {
    ElMessage.error('分配失败')
  } finally {
    submitLoading.value = false
  }
}

onMounted(fetchDepartments)
</script>

<style scoped lang="scss">
.departments-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
  .admin-transfer { display: flex; justify-content: center; }
}
</style>
