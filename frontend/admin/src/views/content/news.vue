<template>
  <div class="news-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="标题/内容" clearable />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchForm.category" placeholder="全部" clearable>
            <el-option v-for="item in categoryOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="已发布" value="published" />
            <el-option label="草稿" value="draft" />
            <el-option label="已归档" value="archived" />
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
        <h3>新闻公告</h3>
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          新增新闻
        </el-button>
      </div>

      <el-table :data="list" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="title" label="标题" min-width="180" show-overflow-tooltip />
        <el-table-column prop="category" label="分类" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ categoryLabel(row.category) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="author" label="作者" width="100" />
        <el-table-column prop="views" label="浏览量" width="90" />
        <el-table-column prop="is_top" label="置顶" width="70">
          <template #default="{ row }">
            <el-tag v-if="row.is_top" type="warning" size="small">置顶</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'published' ? 'success' : row.status === 'draft' ? 'info' : 'warning'" size="small">
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="published_at" label="发布时间" width="170" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button
              :type="row.is_top ? 'warning' : 'success'"
              link
              @click="handleToggleTop(row)"
            >
              {{ row.is_top ? '取消置顶' : '置顶' }}
            </el-button>
            <el-button
              v-if="row.status === 'draft'"
              type="success"
              link
              @click="handlePublish(row)"
            >
              发布
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
      :title="isEdit ? '编辑新闻' : '新增新闻'"
      width="750px"
      destroy-on-close
    >
      <el-form :model="formData" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入新闻标题" />
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="分类" prop="category">
              <el-select v-model="formData.category" placeholder="请选择分类" style="width: 100%">
                <el-option v-for="item in categoryOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="作者">
              <el-input v-model="formData.author" placeholder="作者名称" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="摘要">
          <el-input v-model="formData.summary" type="textarea" :rows="2" placeholder="新闻摘要" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="8" placeholder="请输入新闻内容" />
        </el-form-item>
        <el-form-item label="封面图">
          <el-input v-model="formData.cover" placeholder="封面图片URL" />
        </el-form-item>
        <el-form-item label="标签">
          <el-select
            v-model="formData.tags"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="输入标签后回车"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio value="draft">草稿</el-radio>
            <el-radio value="published">已发布</el-radio>
            <el-radio value="archived">已归档</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
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

const searchForm = ref({ keyword: '', category: '', status: '' })

const categoryOptions = [
  { label: '公司新闻', value: 'company' },
  { label: '产品动态', value: 'product' },
  { label: '行业资讯', value: 'industry' },
  { label: '活动公告', value: 'event' },
  { label: '系统公告', value: 'system' }
]

const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

const formData = reactive({
  title: '',
  category: '',
  author: '',
  content: '',
  cover: '',
  summary: '',
  tags: [] as string[],
  status: 'draft'
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

const categoryLabel = (val: string) => categoryOptions.find((c) => c.value === val)?.label || val

const statusLabel = (status: string) => {
  const map: Record<string, string> = { draft: '草稿', published: '已发布', archived: '已归档' }
  return map[status] || status
}

const fetchList = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/content/news', {
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
  searchForm.value = { keyword: '', category: '', status: '' }
  handleSearch()
}

const resetForm = () => {
  Object.assign(formData, { title: '', category: '', author: '', content: '', cover: '', summary: '', tags: [], status: 'draft' })
  editingId.value = null
}

const handleAdd = () => { isEdit.value = false; resetForm(); dialogVisible.value = true }

const handleEdit = (row: any) => {
  isEdit.value = true
  editingId.value = row.id
  Object.assign(formData, {
    title: row.title, category: row.category, author: row.author || '',
    content: row.content || '', cover: row.cover || '', summary: row.summary || '',
    tags: row.tags || [], status: row.status
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (isEdit.value && editingId.value) {
      await request.put(`/admin/api/v1/content/news/${editingId.value}`, formData)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/content/news', formData)
      ElMessage.success('新增成功')
    }
    dialogVisible.value = false
    fetchList()
  } catch { /* handled */ } finally {
    submitLoading.value = false
  }
}

const handleToggleTop = async (row: any) => {
  try {
    await request.put(`/admin/api/v1/content/news/${row.id}/top`, { is_top: !row.is_top })
    ElMessage.success(row.is_top ? '已取消置顶' : '已置顶')
    fetchList()
  } catch { /* handled */ }
}

const handlePublish = async (row: any) => {
  try {
    await request.put(`/admin/api/v1/content/news/${row.id}`, { status: 'published' })
    ElMessage.success('发布成功')
    fetchList()
  } catch { /* handled */ }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm(`确定要删除新闻「${row.title}」吗？`, '警告', { type: 'error' })
  try {
    await request.delete(`/admin/api/v1/content/news/${row.id}`)
    ElMessage.success('删除成功')
    fetchList()
  } catch { /* handled */ }
}

onMounted(fetchList)
</script>

<style scoped lang="scss">
.news-page {
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
