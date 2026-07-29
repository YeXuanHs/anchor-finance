<template>
  <div class="pages-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="页面标题/路径" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="已发布" value="published" />
            <el-option label="草稿" value="draft" />
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
        <h3>自定义页面</h3>
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          添加页面
        </el-button>
      </div>

      <el-table :data="list" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="title" label="标题" min-width="160" show-overflow-tooltip />
        <el-table-column prop="slug" label="URL路径" width="180">
          <template #default="{ row }">
            <el-text type="info" size="small">/{{ row.slug }}</el-text>
          </template>
        </el-table-column>
        <el-table-column prop="view_count" label="浏览数" width="90" />
        <el-table-column prop="is_home" label="首页" width="70">
          <template #default="{ row }">
            <el-tag v-if="row.is_home" type="warning" size="small">首页</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">
              {{ row.status === 'published' ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="170" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="success" link @click="handlePreview(row)">预览</el-button>
            <el-button
              :type="row.status === 'published' ? 'warning' : 'success'"
              link
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 'published' ? '下线' : '发布' }}
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
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
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑页面' : '添加页面'"
      width="900px"
      destroy-on-close
      top="3vh"
    >
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本内容" name="content">
          <el-form :model="formData" :rules="rules" ref="formRef" label-width="100px">
            <el-row :gutter="16">
              <el-col :span="16">
                <el-form-item label="标题" prop="title">
                  <el-input v-model="formData.title" placeholder="页面标题" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="URL路径" prop="slug">
                  <el-input v-model="formData.slug" placeholder="about-us" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="页面内容" prop="content">
              <el-input
                v-model="formData.content"
                type="textarea"
                :rows="15"
                placeholder="支持HTML富文本内容"
              />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="SEO设置" name="seo">
          <el-form label-width="100px">
            <el-form-item label="SEO标题">
              <el-input v-model="formData.seo_title" placeholder="SEO标题（留空则使用页面标题）" />
            </el-form-item>
            <el-form-item label="SEO描述">
              <el-input v-model="formData.seo_description" type="textarea" :rows="3" placeholder="搜索引擎描述" />
            </el-form-item>
            <el-form-item label="SEO关键词">
              <el-input v-model="formData.seo_keywords" placeholder="关键词，用逗号分隔" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="高级设置" name="advanced">
          <el-form label-width="100px">
            <el-form-item label="状态">
              <el-radio-group v-model="formData.status">
                <el-radio value="published">发布</el-radio>
                <el-radio value="draft">草稿</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="设为首页">
              <el-switch v-model="formData.is_home" />
              <span style="margin-left: 8px; color: #909399; font-size: 12px">开启后此页面将作为站点首页展示</span>
            </el-form-item>
            <el-form-item label="自定义样式">
              <el-input v-model="formData.custom_css" type="textarea" :rows="5" placeholder="自定义CSS样式" />
            </el-form-item>
            <el-form-item label="自定义脚本">
              <el-input v-model="formData.custom_js" type="textarea" :rows="5" placeholder="自定义JavaScript代码" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const list = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const searchForm = ref({ keyword: '', status: '' })

const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)
const activeTab = ref('content')

const formData = reactive({
  title: '',
  slug: '',
  content: '',
  seo_title: '',
  seo_description: '',
  seo_keywords: '',
  status: 'draft',
  is_home: false,
  custom_css: '',
  custom_js: ''
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入页面标题', trigger: 'blur' }],
  slug: [
    { required: true, message: '请输入URL路径', trigger: 'blur' },
    { pattern: /^[a-z0-9-]+$/, message: '仅支持小写字母、数字和连字符', trigger: 'blur' }
  ],
  content: [{ required: true, message: '请输入页面内容', trigger: 'blur' }]
}

const fetchList = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/content/pages', {
      params: { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    })
    list.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

const handleSearch = () => { currentPage.value = 1; fetchList() }

const resetSearch = () => {
  searchForm.value = { keyword: '', status: '' }
  handleSearch()
}

const resetForm = () => {
  Object.assign(formData, {
    title: '', slug: '', content: '', seo_title: '', seo_description: '',
    seo_keywords: '', status: 'draft', is_home: false, custom_css: '', custom_js: ''
  })
  editingId.value = null
  activeTab.value = 'content'
}

const handleAdd = () => { isEdit.value = false; resetForm(); dialogVisible.value = true }

const handleEdit = (row: any) => {
  isEdit.value = true
  editingId.value = row.id
  Object.assign(formData, {
    title: row.title, slug: row.slug, content: row.content || '',
    seo_title: row.seo_title || '', seo_description: row.seo_description || '',
    seo_keywords: row.seo_keywords || '', status: row.status,
    is_home: row.is_home || false, custom_css: row.custom_css || '', custom_js: row.custom_js || ''
  })
  activeTab.value = 'content'
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (isEdit.value && editingId.value) {
      await request.put(`/admin/api/v1/content/pages/${editingId.value}`, formData)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/content/pages', formData)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    fetchList()
  } catch { /* handled */ } finally {
    submitLoading.value = false
  }
}

const handlePreview = (row: any) => {
  window.open(`/page/${row.slug}`, '_blank')
}

const handleToggleStatus = async (row: any) => {
  const newStatus = row.status === 'published' ? 'draft' : 'published'
  const action = newStatus === 'published' ? '发布' : '下线'
  await ElMessageBox.confirm(`确定要${action}页面「${row.title}」吗？`, '提示', { type: 'warning' })
  try {
    await request.put(`/admin/api/v1/content/pages/${row.id}`, { status: newStatus })
    ElMessage.success(`${action}成功`)
    fetchList()
  } catch { /* handled */ }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm(`确定要删除页面「${row.title}」吗？此操作不可恢复。`, '警告', { type: 'error' })
  try {
    await request.delete(`/admin/api/v1/content/pages/${row.id}`)
    ElMessage.success('删除成功')
    fetchList()
  } catch { /* handled */ }
}

onMounted(fetchList)
</script>

<style scoped lang="scss">
.pages-page {
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
