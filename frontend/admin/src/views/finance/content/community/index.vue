<template>
  <div class="page-container">
    <art-card title="社区管理" shadow="never">
      <template #header-extra>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          发布帖子
        </el-button>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="标题">
          <el-input v-model="searchForm.keyword" placeholder="请输入关键词" clearable />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchForm.category" placeholder="全部" clearable>
            <el-option label="公告" value="notice" />
            <el-option label="讨论" value="discussion" />
            <el-option label="分享" value="share" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="已发布" :value="1" />
            <el-option label="已下架" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleSearchReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="author" label="作者" width="120" />
        <el-table-column prop="category" label="分类" width="120" />
        <el-table-column prop="view_count" label="浏览" width="100" />
        <el-table-column prop="comment_count" label="评论" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '已发布' : '已下架' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="发布时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleComments(row)">评论</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </art-card>

    <!-- 创建/编辑帖子对话框 -->
    <el-dialog v-model="postDialogVisible" :title="postDialogTitle" width="700px" destroy-on-close>
      <el-form :model="postForm" :rules="postFormRules" ref="postFormRef" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="postForm.title" placeholder="请输入帖子标题" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="分类" prop="category">
              <el-select v-model="postForm.category" placeholder="请选择分类" style="width: 100%">
                <el-option label="公告" value="notice" />
                <el-option label="讨论" value="discussion" />
                <el-option label="分享" value="share" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="postForm.status" style="width: 100%">
                <el-option label="已发布" :value="1" />
                <el-option label="已下架" :value="0" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="内容" prop="content">
          <el-input v-model="postForm.content" type="textarea" :rows="10" placeholder="请输入帖子内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="postDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handlePostSubmit" :loading="postSubmitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 评论列表对话框 -->
    <el-dialog v-model="commentsDialogVisible" title="评论管理" width="700px">
      <el-table :data="commentsList" v-loading="commentsLoading" stripe>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="author" label="评论者" width="100" />
        <el-table-column prop="content" label="内容" min-width="300" show-overflow-tooltip />
        <el-table-column prop="created_at" label="时间" width="170" />
        <el-table-column label="操作" width="80">
          <template #default="{ row }">
            <el-button type="danger" link @click="handleDeleteComment(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="commentsDialogVisible = false">关闭</el-button>
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
const postDialogTitle = ref('发布帖子')
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
  title: [{ required: true, message: '请输入帖子标题', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  content: [{ required: true, message: '请输入帖子内容', trigger: 'blur' }]
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
    const { data } = await request.get('/api/admin/community/posts', { params })
    tableData.value = data?.data || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  postDialogTitle.value = '发布帖子'
  resetPostForm()
  postDialogVisible.value = true
}

const handleEdit = (row: any) => {
  postDialogTitle.value = '编辑帖子'
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
        ElMessage.success('更新成功')
      } else {
        await request.post({ url: '/api/admin/community/posts', params: { ...postForm } })
        ElMessage.success('发布成功')
      }
      postDialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error('操作失败')
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
    const { data } = await request.get(`/api/admin/community/posts/${row.id}/comments`)
    commentsList.value = data?.data || []
  } catch (error) {
    ElMessage.error('获取评论失败')
  } finally {
    commentsLoading.value = false
  }
}

const handleDeleteComment = async (row: any) => {
  await ElMessageBox.confirm('确定删除该评论？', '提示', { type: 'warning' })
  try {
    await request.delete(`/api/admin/community/comments/${row.id}`)
    ElMessage.success('删除成功')
    handleComments({ id: currentPostId.value })
  } catch (error) {
    console.error(error)
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm('确定删除该帖子？', '提示', { type: 'warning' })
  try {
    await request.delete(`/api/admin/community/posts/${row.id}`)
    ElMessage.success('删除成功')
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
