<template>
  <div class="ticket-prereply-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>工单预回复管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加预回复
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="标题/内容" clearable />
        </el-form-item>
        <el-form-item label="部门">
          <el-select v-model="searchForm.department_id" placeholder="全部" clearable>
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="150" show-overflow-tooltip />
        <el-table-column prop="department_name" label="所属部门" width="120" />
        <el-table-column prop="content" label="回复内容" min-width="300" show-overflow-tooltip />
        <el-table-column prop="sort_order" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="use_count" label="使用次数" width="100" align="center" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该预回复吗？" @confirm="handleDelete(row)">
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入预回复标题" />
        </el-form-item>
        <el-form-item label="所属部门" prop="department_id">
          <el-select v-model="formData.department_id" placeholder="请选择部门">
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="回复内容" prop="content">
          <el-input
            v-model="formData.content"
            type="textarea"
            :rows="6"
            placeholder="请输入预回复内容"
          />
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="999" />
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

interface PreReply {
  id: number
  title: string
  content: string
  department_id: number
  department_name: string
  sort_order: number
  status: number
  use_count: number
  created_at: string
}

interface Department {
  id: number
  name: string
}

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('添加预回复')
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: '',
  department_id: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<PreReply[]>([])
const departments = ref<Department[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  title: '',
  content: '',
  department_id: undefined as number | undefined,
  sort_order: 0,
  status: 1
})

const formRules: FormRules = {
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' },
    { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  department_id: [
    { required: true, message: '请选择所属部门', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请输入回复内容', trigger: 'blur' }
  ]
}

const fetchPreReplies = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/ticket-prereplies',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取预回复列表失败:', error)
    ElMessage.error('获取预回复列表失败')
  } finally {
    loading.value = false
  }
}

const fetchDepartments = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/ticket-departments'
    })
    departments.value = data || []
  } catch (error) {
    console.error('获取部门列表失败:', error)
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchPreReplies()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.department_id = undefined
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = '添加预回复'
  formData.id = undefined
  formData.title = ''
  formData.content = ''
  formData.department_id = undefined
  formData.sort_order = 0
  formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: PreReply) => {
  dialogTitle.value = '编辑预回复'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: PreReply) => {
  try {
    await request.del({
      url: `/api/admin/ticket-prereplies/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchPreReplies()
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
        await request.put({
          url: `/api/admin/ticket-prereplies/${formData.id}`,
          params: formData
        })
      } else {
        await request.post({
          url: '/api/admin/ticket-prereplies',
          params: formData
        })
      }

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchPreReplies()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchPreReplies()
}

const handlePageChange = () => {
  fetchPreReplies()
}

onMounted(() => {
  fetchPreReplies()
  fetchDepartments()
})
</script>

<style scoped lang="scss">
.ticket-prereply-page {
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
</style>
