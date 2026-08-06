<template>
  <div class="email-templates-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>邮件模板</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="模板名称/标识" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
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
        <el-table-column prop="name" label="模板名称" min-width="150" />
        <el-table-column prop="code" label="模板标识" width="160" />
        <el-table-column prop="subject" label="邮件主题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="180" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
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

    <!-- 编辑对话框 -->
    <el-dialog v-model="dialogVisible" title="编辑邮件模板" width="900px" top="5vh">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入模板名称" />
        </el-form-item>
        <el-form-item label="模板标识" prop="code">
          <el-input v-model="formData.code" disabled />
        </el-form-item>
        <el-form-item label="邮件主题" prop="subject">
          <el-input v-model="formData.subject" placeholder="请输入邮件主题" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="模板内容" prop="content">
          <div class="editor-wrapper">
            <ArtWangEditor
              ref="editorRef"
              v-model="formData.content"
              height="400px"
              placeholder="请输入邮件模板内容..."
            />
          </div>
        </el-form-item>
        <el-form-item label="可用变量">
          <div class="variables-info">
            <el-tag
              v-for="variable in formData.variables"
              :key="variable"
              class="variable-tag"
              type="info"
              effect="plain"
            >
              {{ '{' + variable + '}' }}
            </el-tag>
            <span v-if="!formData.variables?.length" class="empty-text">暂无可用变量</span>
          </div>
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
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

interface EmailTemplate {
  id: number
  name: string
  code: string
  subject: string
  content: string
  status: number
  variables: string[]
  updated_at: string
}

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const editorRef = ref()

const searchForm = reactive({
  keyword: '',
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<EmailTemplate[]>([])

const formData = reactive({
  id: 0,
  name: '',
  code: '',
  subject: '',
  content: '',
  status: 1,
  variables: [] as string[]
})

const formRules: FormRules = {
  name: [
    { required: true, message: '请输入模板名称', trigger: 'blur' }
  ],
  subject: [
    { required: true, message: '请输入邮件主题', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入模板内容', trigger: 'blur' }
  ]
}

// 获取模板列表
const fetchTemplates = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/email-templates',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取邮件模板列表失败:', error)
    ElMessage.error('获取邮件模板列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchTemplates()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = undefined
  handleSearch()
}

// 编辑
const handleEdit = (row: EmailTemplate) => {
  Object.assign(formData, {
    id: row.id,
    name: row.name,
    code: row.code,
    subject: row.subject,
    content: row.content,
    status: row.status,
    variables: row.variables || []
  })
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      await request.put({
        url: `/api/admin/email-templates/${formData.id}`,
        params: {
          name: formData.name,
          subject: formData.subject,
          content: formData.content,
          status: formData.status
        },
        showSuccessMessage: true
      })
      ElMessage.success('更新成功')
      dialogVisible.value = false
      fetchTemplates()
    } catch (error) {
      ElMessage.error('更新失败')
    } finally {
      submitLoading.value = false
    }
  })
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchTemplates()
}

// 页码变化
const handlePageChange = () => {
  fetchTemplates()
}

onMounted(() => {
  fetchTemplates()
})
</script>

<style scoped lang="scss">
.email-templates-page {
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

.editor-wrapper {
  width: 100%;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  overflow: hidden;
}

.variables-info {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.variable-tag {
  cursor: pointer;
}

.empty-text {
  color: var(--el-text-color-placeholder);
  font-size: 13px;
}
</style>
