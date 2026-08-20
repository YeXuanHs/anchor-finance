<template>
  <div class="kb-page">
    <div class="page-header">
      <div class="header-left">
        <h2>{{ $t('knowledgeBase.title') }}</h2>
        <span class="subtitle">{{ $t('knowledgeBase.subtitle') }}</span>
      </div>
    </div>

    <el-tabs v-model="activeTab">
      <!-- 分类管理 -->
      <el-tab-pane :label="$t('knowledgeBase.categoryManage')" name="categories">
        <div class="tab-header">
          <el-button type="primary" @click="showAddCategory">{{ $t('knowledgeBase.addCategory') }}</el-button>
        </div>
        <el-table :data="categories" v-loading="loadingCategories" stripe>
          <el-table-column prop="name" :label="$t('knowledgeBase.categoryName')" min-width="200" />
          <el-table-column prop="description" :label="$t('common.description')" min-width="300" />
          <el-table-column prop="sort_order" :label="$t('common.sort')" width="80" align="center" />
          <el-table-column prop="is_active" :label="$t('common.status')" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.is_active ? 'success' : 'info'" size="small">
                {{ row.is_active ? $t('common.enabled') : $t('common.disabled') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('common.action')" width="150" align="center">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="showEditCategory(row)">{{ $t('common.edit') }}</el-button>
              <el-popconfirm :title="$t('common.confirmDelete')" @confirm="deleteCategory(row.id)">
                <template #reference>
                  <el-button type="danger" link size="small">{{ $t('common.delete') }}</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>

        <el-dialog v-model="catDialogVisible" :title="isEditCat ? $t('knowledgeBase.editCategory') : $t('knowledgeBase.addCategory')" width="500px">
          <el-form :model="catForm" label-width="80px">
            <el-form-item :label="$t('common.name')" required>
              <el-input v-model="catForm.name" :placeholder="$t('knowledgeBase.enterCategoryName')" />
            </el-form-item>
            <el-form-item :label="$t('common.description')">
              <el-input v-model="catForm.description" type="textarea" :rows="3" />
            </el-form-item>
            <el-form-item :label="$t('common.sort')">
              <el-input-number v-model="catForm.sort_order" :min="0" />
            </el-form-item>
            <el-form-item :label="$t('common.enabled')">
              <el-switch v-model="catForm.is_active" />
            </el-form-item>
          </el-form>
          <template #footer>
            <el-button @click="catDialogVisible = false">{{ $t('common.cancel') }}</el-button>
            <el-button type="primary" @click="submitCategory">{{ $t('common.confirm') }}</el-button>
          </template>
        </el-dialog>
      </el-tab-pane>

      <!-- 文章管理 -->
      <el-tab-pane :label="$t('knowledgeBase.articleManage')" name="articles">
        <div class="tab-header">
          <el-select v-model="filterCategory" :placeholder="$t('knowledgeBase.filterCategory')" clearable style="width: 200px" @change="loadArticles">
            <el-option :label="$t('knowledgeBase.allCategories')" :value="0" />
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
          <el-input v-model="searchKeyword" :placeholder="$t('knowledgeBase.searchArticles')" clearable style="width: 300px" @clear="loadArticles">
            <template #append>
              <el-button @click="loadArticles">{{ $t('common.search') }}</el-button>
            </template>
          </el-input>
          <el-button type="primary" @click="showAddArticle">{{ $t('knowledgeBase.addArticle') }}</el-button>
        </div>

        <el-table :data="articles" v-loading="loadingArticles" stripe>
          <el-table-column prop="title" :label="$t('knowledgeBase.articleTitle')" min-width="250" />
          <el-table-column prop="category_id" :label="$t('knowledgeBase.category')" width="120">
            <template #default="{ row }">
              {{ getCategoryName(row.category_id) }}
            </template>
          </el-table-column>
          <el-table-column prop="is_faq" label="FAQ" width="80" align="center">
            <template #default="{ row }">
              <el-tag v-if="row.is_faq" type="warning" size="small">FAQ</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="view_count" :label="$t('knowledgeBase.views')" width="80" align="center" />
          <el-table-column prop="help_count" :label="$t('knowledgeBase.helpful')" width="80" align="center" />
          <el-table-column prop="is_active" :label="$t('common.status')" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.is_active ? 'success' : 'info'" size="small">
                {{ row.is_active ? $t('common.enabled') : $t('common.disabled') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('common.action')" width="150" align="center">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="showEditArticle(row)">{{ $t('common.edit') }}</el-button>
              <el-popconfirm :title="$t('common.confirmDelete')" @confirm="deleteArticle(row.id)">
                <template #reference>
                  <el-button type="danger" link size="small">{{ $t('common.delete') }}</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>

        <el-pagination
          v-if="articleTotal > 0"
          :current-page="currentPage"
          :page-size="20"
          :total="articleTotal"
          layout="total, prev, pager, next"
          @current-change="handlePageChange"
          style="margin-top: 16px; justify-content: flex-end;"
        />

        <el-dialog v-model="articleDialogVisible" :title="isEditArticle ? $t('knowledgeBase.editArticle') : $t('knowledgeBase.addArticle')" width="800px" top="5vh">
          <el-form :model="articleForm" label-width="100px">
            <el-form-item :label="$t('knowledgeBase.articleTitle')" required>
              <el-input v-model="articleForm.title" :placeholder="$t('knowledgeBase.enterArticleTitle')" />
            </el-form-item>
            <el-form-item :label="$t('knowledgeBase.category')" required>
              <el-select v-model="articleForm.category_id" :placeholder="$t('knowledgeBase.selectCategory')">
                <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('knowledgeBase.summary')">
              <el-input v-model="articleForm.summary" type="textarea" :rows="2" :placeholder="$t('knowledgeBase.enterSummary')" />
            </el-form-item>
            <el-form-item :label="$t('knowledgeBase.content')" required>
              <el-input v-model="articleForm.content" type="textarea" :rows="10" :placeholder="$t('knowledgeBase.enterContent')" />
            </el-form-item>
            <el-form-item :label="$t('knowledgeBase.keywords')">
              <el-input v-model="articleForm.keywords" :placeholder="$t('knowledgeBase.enterKeywords')" />
            </el-form-item>
            <el-form-item :label="$t('knowledgeBase.tags')">
              <el-input v-model="articleForm.tags" :placeholder="$t('knowledgeBase.enterTags')" />
            </el-form-item>
            <el-row :gutter="20">
              <el-col :span="8">
                <el-form-item label="FAQ">
                  <el-switch v-model="articleForm.is_faq" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item :label="$t('common.enabled')">
                  <el-switch v-model="articleForm.is_active" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item :label="$t('common.sort')">
                  <el-input-number v-model="articleForm.sort_order" :min="0" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
          <template #footer>
            <el-button @click="articleDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="submitArticle">确定</el-button>
          </template>
        </el-dialog>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

const activeTab = ref('categories')
const loadingCategories = ref(false)
const loadingArticles = ref(false)
const categories = ref<any[]>([])
const articles = ref<any[]>([])
const articleTotal = ref(0)
const currentPage = ref(1)
const filterCategory = ref(0)
const searchKeyword = ref('')

// 分类表单
const catDialogVisible = ref(false)
const isEditCat = ref(false)
const editCatId = ref(0)
const catForm = ref({ name: '', description: '', sort_order: 0, is_active: true })

// 文章表单
const articleDialogVisible = ref(false)
const isEditArticle = ref(false)
const editArticleId = ref(0)
const articleForm = ref({
  title: '', category_id: null as number | null, summary: '', content: '',
  keywords: '', tags: '', is_faq: false, is_active: true, sort_order: 0
})

const loadCategories = async () => {
  loadingCategories.value = true
  try {
    const res = await request.get({ url: '/api/admin/knowledge/categories', params: { show_inactive: true } })
    categories.value = res || []
  } catch (e) { ElMessage.error('加载分类失败') }
  finally { loadingCategories.value = false }
}

const loadArticles = async () => {
  loadingArticles.value = true
  try {
    const params: any = { page: currentPage.value, page_size: 20 }
    if (filterCategory.value > 0) params.category_id = filterCategory.value
    if (searchKeyword.value) params.keyword = searchKeyword.value
    const res = await request.get({ url: '/api/admin/knowledge/articles', params })
    articles.value = res?.items || []
    articleTotal.value = res?.total || 0
  } catch (e) { ElMessage.error('加载文章失败') }
  finally { loadingArticles.value = false }
}

const getCategoryName = (id: number) => categories.value.find(c => c.id === id)?.name || '-'

const showAddCategory = () => {
  isEditCat.value = false
  catForm.value = { name: '', description: '', sort_order: 0, is_active: true }
  catDialogVisible.value = true
}

const showEditCategory = (row: any) => {
  isEditCat.value = true
  editCatId.value = row.id
  catForm.value = { ...row }
  catDialogVisible.value = true
}

const submitCategory = async () => {
  if (!catForm.value.name) { ElMessage.warning('请输入分类名称'); return }
  try {
    if (isEditCat.value) {
      await request.put({ url: `/api/admin/knowledge/categories/${editCatId.value}`, params: catForm.value })
    } else {
      await request.post({ url: '/api/admin/knowledge/categories', params: catForm.value })
    }
    ElMessage.success('操作成功')
    catDialogVisible.value = false
    loadCategories()
  } catch (e) { ElMessage.error('操作失败') }
}

const deleteCategory = async (id: number) => {
  try {
    await request.del({ url: `/api/admin/knowledge/categories/${id}` })
    ElMessage.success('删除成功')
    loadCategories()
  } catch (e: any) { ElMessage.error('删除失败') }
}

const showAddArticle = () => {
  isEditArticle.value = false
  articleForm.value = { title: '', category_id: null, summary: '', content: '', keywords: '', tags: '', is_faq: false, is_active: true, sort_order: 0 }
  articleDialogVisible.value = true
}

const showEditArticle = (row: any) => {
  isEditArticle.value = true
  editArticleId.value = row.id
  articleForm.value = { ...row }
  articleDialogVisible.value = true
}

const submitArticle = async () => {
  if (!articleForm.value.title || !articleForm.value.category_id) { ElMessage.warning('请填写标题和分类'); return }
  try {
    if (isEditArticle.value) {
      await request.put({ url: `/api/admin/knowledge/articles/${editArticleId.value}`, params: articleForm.value })
    } else {
      await request.post({ url: '/api/admin/knowledge/articles', params: articleForm.value })
    }
    ElMessage.success('操作成功')
    articleDialogVisible.value = false
    loadArticles()
  } catch (e) { ElMessage.error('操作失败') }
}

const deleteArticle = async (id: number) => {
  try {
    await request.del({ url: `/api/admin/knowledge/articles/${id}` })
    ElMessage.success('删除成功')
    loadArticles()
  } catch (e) { ElMessage.error('删除失败') }
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  loadArticles()
}

onMounted(() => {
  loadCategories()
  loadArticles()
})
</script>

<style scoped>
.kb-page { padding: 20px; }
.page-header { margin-bottom: 20px; }
.page-header h2 { margin: 0 0 4px 0; font-size: 20px; }
.subtitle { color: #909399; font-size: 14px; }
.tab-header { display: flex; gap: 12px; align-items: center; margin-bottom: 16px; }
</style>
