<template>
  <div class="page-container">
    <art-card :title="$t('finance.community.title')" shadow="never">
      <template #header-extra>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          {{ $t('finance.community.publishPost') }}
        </el-button>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('finance.community.titleLabel')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('finance.community.enterKeyword')" clearable />
        </el-form-item>
        <el-form-item :label="$t('finance.community.category')">
          <el-select v-model="searchForm.category" :placeholder="$t('finance.community.all')" clearable>
            <el-option :label="$t('finance.community.notice')" value="notice" />
            <el-option :label="$t('finance.community.discussion')" value="discussion" />
            <el-option :label="$t('finance.community.share')" value="share" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('finance.community.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('finance.community.all')" clearable>
            <el-option :label="$t('finance.community.published')" :value="1" />
            <el-option :label="$t('finance.community.unpublished')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('finance.community.search') }}</el-button>
          <el-button @click="handleSearchReset">{{ $t('finance.community.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" :label="$t('finance.community.titleLabel')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="author" :label="$t('finance.community.author')" width="120" />
        <el-table-column prop="category" :label="$t('finance.community.category')" width="120" />
        <el-table-column prop="view_count" :label="$t('finance.community.views')" width="100" />
        <el-table-column prop="comment_count" :label="$t('finance.community.comments')" width="100" />
        <el-table-column prop="status" :label="$t('finance.community.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? $t('finance.community.published') : $t('finance.community.unpublished') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('finance.community.publishTime')" width="180" />
        <el-table-column :label="$t('finance.community.actions')" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">{{ $t('finance.community.edit') }}</el-button>
            <el-button size="small" @click="handleComments(row)">{{ $t('finance.community.comments') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">{{ $t('finance.community.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </art-card>

    <!-- 创建/编辑帖子对话框 -->
    <el-dialog v-model="postDialogVisible" :title="postDialogTitle" width="700px" destroy-on-close>
      <el-form :model="postForm" :rules="postFormRules" ref="postFormRef" label-width="100px">
        <el-form-item :label="$t('finance.community.titleLabel')" prop="title">
          <el-input v-model="postForm.title" :placeholder="$t('finance.community.enterPostTitle')" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('finance.community.category')" prop="category">
              <el-select v-model="postForm.category" :placeholder="$t('finance.community.selectCategory')" style="width: 100%">
                <el-option :label="$t('finance.community.notice')" value="notice" />
                <el-option :label="$t('finance.community.discussion')" value="discussion" />
                <el-option :label="$t('finance.community.share')" value="share" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('finance.community.status')" prop="status">
              <el-select v-model="postForm.status" style="width: 100%">
                <el-option :label="$t('finance.community.published')" :value="1" />
                <el-option :label="$t('finance.community.unpublished')" :value="0" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('finance.community.content')" prop="content">
          <el-input v-model="postForm.content" type="textarea" :rows="10" :placeholder="$t('finance.community.enterPostContent')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="postDialogVisible = false">{{ $t('finance.community.cancel') }}</el-button>
        <el-button type="primary" @click="handlePostSubmit" :loading="postSubmitLoading">{{ $t('finance.community.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 评论列表对话框 -->
    <el-dialog v-model="commentsDialogVisible" :title="$t('finance.community.commentManagement')" width="700px">
      <el-table :data="commentsList" v-loading="commentsLoading" stripe>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="author" :label="$t('finance.community.commenter')" width="100" />
        <el-table-column prop="content" :label="$t('finance.community.content')" min-width="300" show-overflow-tooltip />
        <el-table-column prop="created_at" :label="$t('finance.community.time')" width="170" />
        <el-table-column :label="$t('finance.community.actions')" width="80">
          <template #default="{ row }">
            <el-button type="danger" link @click="handleDeleteComment(row)">{{ $t('finance.community.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="commentsDialogVisible = false">{{ $t('finance.community.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const tableData = ref([])

const searchForm = reactive({
  keyword: '',
  category: '',
  status: undefined as number | undefined
})

const handleSearch = () => fetchData()
const handleSearchReset = () => {
  searchForm.keyword = ''
  searchForm.category = ''
  searchForm.status = undefined
  fetchData()
}

// 帖子对话框
const postDialogVisible = ref(false)
const postDialogTitle = ref($t('finance.community.publishPost'))
const postSubmitLoading = ref(false)
const postFormRef = ref<FormInstance>()
const postForm = reactive({
  id: undefined as number | undefined,
  title: '',
  category: '',
  status: 1 as number,
  content: ''
})
const postFormRules: FormRules = {
  title: [{ required: true, message: $t('finance.community.enterPostTitle'), trigger: 'blur' }],
  category: [{ required: true, message: $t('finance.community.selectCategory'), trigger: 'change' }],
  content: [{ required: true, message: $t('finance.community.enterPostContent'), trigger: 'blur' }]
}

const resetPostForm = () => {
  postForm.id = undefined
  postForm.title = ''
  postForm.category = ''
  postForm.status = 1
  postForm.content = ''
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {}
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.category) params.category = searchForm.category
    if (searchForm.status !== undefined) params.status = searchForm.status
    const res = await request.get({ url: '/api/admin/community/posts', params })
    tableData.value = res || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  postDialogTitle.value = $t('finance.community.publishPost')
  resetPostForm()
  postDialogVisible.value = true
}

const handleEdit = (row: any) => {
  postDialogTitle.value = $t('finance.community.editPost')
  postForm.id = row.id
  postForm.title = row.title
  postForm.category = row.category
  postForm.status = row.status
  postForm.content = row.content || ''
  postDialogVisible.value = true
}

const handlePostSubmit = async () => {
  if (!postFormRef.value) return
  await postFormRef.value.validate(async (valid) => {
    if (!valid) return
    postSubmitLoading.value = true
    try {
      if (postForm.id) {
        await request.put({ url: `/api/admin/community/posts/${postForm.id}`, params: { ...postForm } })
        ElMessage.success($t('finance.community.updateSuccess'))
      } else {
        await request.post({ url: '/api/admin/community/posts', params: { ...postForm } })
        ElMessage.success($t('finance.community.publishSuccess'))
      }
      postDialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('finance.community.operationFailed'))
    } finally {
      postSubmitLoading.value = false
    }
  })
}

// 评论对话框
const commentsDialogVisible = ref(false)
const commentsLoading = ref(false)
const commentsList = ref<any[]>([])
const currentPostId = ref<number>(0)

const handleComments = async (row: any) => {
  currentPostId.value = row.id
  commentsDialogVisible.value = true
  commentsLoading.value = true
  try {
    const res = await request.get({ url: `/api/admin/community/posts/${row.id}/comments` })
    commentsList.value = res || []
  } catch (error) {
    ElMessage.error($t('finance.community.fetchCommentsFailed'))
  } finally {
    commentsLoading.value = false
  }
}

const handleDeleteComment = async (row: any) => {
  await ElMessageBox.confirm($t('finance.community.confirmDeleteComment'), $t('finance.community.tip'), { type: 'warning' })
  try {
    await request.del({ url: `/api/admin/community/comments/${row.id}` })
    ElMessage.success($t('finance.community.deleteSuccess'))
    handleComments({ id: currentPostId.value })
  } catch (error) {
    console.error(error)
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm($t('finance.community.confirmDeletePost'), $t('finance.community.tip'), { type: 'warning' })
  try {
    await request.del({ url: `/api/admin/community/posts/${row.id}` })
    ElMessage.success($t('finance.community.deleteSuccess'))
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => fetchData())
</script>

<style scoped lang="scss">
.page-container {
  padding: 20px;
}
.search-form {
  margin-bottom: 20px;
}
</style>
