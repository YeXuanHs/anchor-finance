<template>
  <div class="prereply-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="标题/内容" clearable />
        </el-form-item>
        <el-form-item label="部门">
          <el-select v-model="searchForm.department_id" placeholder="全部" clearable>
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchForm.category" placeholder="全部" clearable>
            <el-option v-for="cat in categories" :key="cat" :label="cat" :value="cat" />
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
        <h3>预设回复</h3>
        <div>
          <el-button type="primary" @click="openCategoryDialog">管理分类</el-button>
          <el-button type="primary" @click="openAddDialog">
            <el-icon><Plus /></el-icon>
            添加预设回复
          </el-button>
        </div>
      </div>

      <el-table :data="prereplies" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="180" />
        <el-table-column prop="department.name" label="部门" width="120" />
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.category }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="内容" min-width="250" show-overflow-tooltip />
        <el-table-column prop="use_count" label="使用次数" width="100" align="center" />
        <el-table-column prop="sort_order" label="排序" width="80" align="center" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editPrereply(row)">编辑</el-button>
            <el-button type="danger" link @click="deletePrereply(row)">删除</el-button>
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
          @size-change="fetchPrereplies"
          @current-change="fetchPrereplies"
        />
      </div>
    </div>

    <el-dialog v-model="showFormDialog" :title="editingItem ? '编辑预设回复' : '添加预设回复'" width="650px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="部门" prop="department_id">
          <el-select v-model="formData.department_id" placeholder="请选择部门">
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="formData.category" placeholder="选择分类" filterable allow-create>
            <el-option v-for="cat in categories" :key="cat" :label="cat" :value="cat" />
          </el-select>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="8" placeholder="请输入预设回复内容，支持变量如 {username}" />
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

    <el-dialog v-model="showCategoryDialog" title="管理分类" width="500px">
      <div class="category-list">
        <div v-for="(cat, idx) in categories" :key="idx" class="category-item">
          <el-input v-model="categories[idx]" size="small" />
          <el-button type="danger" link @click="removeCategory(idx)">
            <el-icon><Delete /></el-icon>
          </el-button>
        </div>
        <el-button type="primary" link @click="addCategory" style="margin-top: 8px;">
          <el-icon><Plus /></el-icon> 添加分类
        </el-button>
      </div>
      <template #footer>
        <el-button @click="showCategoryDialog = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus, Delete } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const prereplies = ref<any[]>([])
const departments = ref<any[]>([])
const categories = ref<string[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showFormDialog = ref(false)
const showCategoryDialog = ref(false)
const editingItem = ref<any>(null)
const formRef = ref<FormInstance>()

const searchForm = ref({ keyword: '', department_id: '', category: '' })
const formData = ref({ title: '', department_id: '', category: '', content: '', sort_order: 0 })
const formRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  department_id: [{ required: true, message: '请选择部门', trigger: 'change' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const fetchPrereplies = async () => {
  loading.value = true
  try {
    const params = { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    const { data } = await request.get('/admin/api/v1/tickets/prereplies', { params })
    prereplies.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取预设回复失败')
  } finally {
    loading.value = false
  }
}

const fetchDepartments = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/tickets/departments', { params: { page_size: 100 } })
    departments.value = data.data?.list || []
  } catch {}
}

const fetchCategories = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/tickets/prereplies/categories')
    categories.value = data.data || []
  } catch {}
}

const handleSearch = () => { currentPage.value = 1; fetchPrereplies() }
const resetSearch = () => { searchForm.value = { keyword: '', department_id: '', category: '' }; handleSearch() }

const openAddDialog = () => {
  editingItem.value = null
  formData.value = { title: '', department_id: '', category: '', content: '', sort_order: 0 }
  showFormDialog.value = true
}

const editPrereply = (item: any) => {
  editingItem.value = item
  formData.value = { title: item.title, department_id: item.department_id, category: item.category, content: item.content, sort_order: item.sort_order }
  showFormDialog.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (editingItem.value) {
      await request.put(`/admin/api/v1/tickets/prereplies/${editingItem.value.id}`, formData.value)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/tickets/prereplies', formData.value)
      ElMessage.success('添加成功')
    }
    showFormDialog.value = false
    fetchPrereplies()
  } catch {
    ElMessage.error(editingItem.value ? '更新失败' : '添加失败')
  } finally {
    submitLoading.value = false
  }
}

const deletePrereply = async (item: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除预设回复「${item.title}」吗？`, '提示', { type: 'warning' })
    await request.delete(`/admin/api/v1/tickets/prereplies/${item.id}`)
    ElMessage.success('删除成功')
    fetchPrereplies()
  } catch {}
}

const openCategoryDialog = () => { showCategoryDialog.value = true }
const addCategory = () => { categories.value.push('') }
const removeCategory = (idx: number) => { categories.value.splice(idx, 1) }

onMounted(() => {
  fetchPrereplies()
  fetchDepartments()
  fetchCategories()
})
</script>

<style scoped lang="scss">
.prereply-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
  .category-item {
    display: flex;
    gap: 8px;
    align-items: center;
    margin-bottom: 8px;
  }
}
</style>
