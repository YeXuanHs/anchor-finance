<template>
  <div class="knowledge-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.knowledge.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('finance.knowledge.addArticle') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('finance.knowledge.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('finance.knowledge.titleOrContent')" clearable />
        </el-form-item>
        <el-form-item :label="$t('finance.knowledge.category')">
          <el-select v-model="searchForm.category" :placeholder="$t('finance.knowledge.all')" clearable>
            <el-option :label="$t('finance.knowledge.faq')" value="faq" />
            <el-option :label="$t('finance.knowledge.tutorial')" value="tutorial" />
            <el-option :label="$t('finance.knowledge.docs')" value="docs" />
            <el-option :label="$t('finance.knowledge.api')" value="api" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('finance.knowledge.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('finance.knowledge.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" :label="$t('finance.knowledge.titleLabel')" min-width="250" show-overflow-tooltip />
        <el-table-column prop="category" :label="$t('finance.knowledge.category')" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ categoryText(row.category) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" :label="$t('finance.knowledge.sort')" width="80" />
        <el-table-column prop="views" :label="$t('finance.knowledge.views')" width="80" />
        <el-table-column prop="is_published" :label="$t('finance.knowledge.status')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_published ? 'success' : 'info'" size="small">
              {{ row.is_published ? $t('finance.knowledge.published') : $t('finance.knowledge.draft') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('finance.knowledge.createdAt')" width="180" />
        <el-table-column prop="updated_at" :label="$t('finance.knowledge.updatedAt')" width="180" />
        <el-table-column :label="$t('finance.knowledge.actions')" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('finance.knowledge.edit') }}</el-button>
            <el-popconfirm :title="$t('finance.knowledge.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('finance.knowledge.delete') }}</el-button>
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="800px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="16">
            <el-form-item :label="$t('finance.knowledge.titleLabel')" prop="title">
              <el-input v-model="formData.title" :placeholder="$t('finance.knowledge.enterTitle')" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('finance.knowledge.category')" prop="category">
              <el-select v-model="formData.category" :placeholder="$t('finance.knowledge.selectCategory')" style="width: 100%">
                <el-option :label="$t('finance.knowledge.faq')" value="faq" />
                <el-option :label="$t('finance.knowledge.tutorial')" value="tutorial" />
                <el-option :label="$t('finance.knowledge.docs')" value="docs" />
                <el-option :label="$t('finance.knowledge.api')" value="api" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('finance.knowledge.summary')" prop="summary">
          <el-input v-model="formData.summary" type="textarea" :rows="2" :placeholder="$t('finance.knowledge.enterSummary')" />
        </el-form-item>
        <el-form-item :label="$t('finance.knowledge.content')" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="12" :placeholder="$t('finance.knowledge.enterContent')" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="$t('finance.knowledge.sort')" prop="sort">
              <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('finance.knowledge.isTop')" prop="is_top">
              <el-switch v-model="formData.is_top" :active-value="1" :inactive-value="0" :active-text="$t('finance.knowledge.yes')" :inactive-text="$t('finance.knowledge.no')" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('finance.knowledge.isPublished')" prop="is_published">
              <el-switch v-model="formData.is_published" :active-value="1" :inactive-value="0" :active-text="$t('finance.knowledge.yes')" :inactive-text="$t('finance.knowledge.no')" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('finance.knowledge.tags')" prop="tags">
          <el-input v-model="formData.tags" :placeholder="$t('finance.knowledge.tagsPlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('finance.knowledge.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('finance.knowledge.confirm') }}</el-button>
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
import { $t } from '@/locales'

defineOptions({ name: 'KnowledgeManage' })

const loading = ref(false)
const submitLoading = ref(false)

const searchForm = reactive({
  keyword: '',
  category: '' as string
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref([])

const dialogVisible = ref(false)
const dialogTitle = ref($t('finance.knowledge.addArticle'))
const formRef = ref<FormInstance>()

const formData = reactive({
  id: undefined as number | undefined,
  title: '',
  category: '' as string,
  summary: '',
  content: '',
  sort: 0,
  is_top: 0 as number,
  is_published: 0 as number,
  tags: ''
})

const formRules: FormRules = {
  title: [
    { required: true, message: $t('finance.knowledge.enterTitle'), trigger: 'blur' },
    { min: 2, max: 200, message: $t('finance.knowledge.lengthLimit'), trigger: 'blur' }
  ],
  category: [
    { required: true, message: $t('finance.knowledge.selectCategory'), trigger: 'change' }
  ],
  content: [
    { required: true, message: $t('finance.knowledge.enterContent'), trigger: 'blur' }
  ]
}

const categoryText = (category: string) => {
  const map: Record<string, string> = {
    faq: $t('finance.knowledge.faq'),
    tutorial: $t('finance.knowledge.tutorial'),
    docs: $t('finance.knowledge.docs'),
    api: $t('finance.knowledge.api')
  }
  return map[category] || category
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/knowledge/articles',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取知识库列表失败:', error)
    ElMessage.error($t('finance.knowledge.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.category = ''
  handleSearch()
}

const resetForm = () => {
  formData.id = undefined
  formData.title = ''
  formData.category = ''
  formData.summary = ''
  formData.content = ''
  formData.sort = 0
  formData.is_top = 0
  formData.is_published = 0
  formData.tags = ''
}

const handleAdd = () => {
  dialogTitle.value = $t('finance.knowledge.addArticle')
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('finance.knowledge.editArticle')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/knowledge/articles/${row.id}` })
    ElMessage.success($t('finance.knowledge.deleteSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('finance.knowledge.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/knowledge/articles/${formData.id}` : '/api/admin/knowledge/articles'
      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }
      ElMessage.success(formData.id ? $t('finance.knowledge.updateSuccess') : $t('finance.knowledge.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('finance.knowledge.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchData()
}

const handlePageChange = () => {
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.knowledge-page {
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
