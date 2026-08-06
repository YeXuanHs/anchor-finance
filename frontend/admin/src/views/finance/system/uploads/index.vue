<template>
  <div class="uploads-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>文件管理</span>
          <el-button type="primary" size="small" @click="handleUpload">上传文件</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="文件类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="图片" value="image" />
            <el-option label="文档" value="document" />
            <el-option label="附件" value="attachment" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="文件名" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="文件名" min-width="200" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ typeMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="100">
          <template #default="{ row }">{{ formatSize(row.size) }}</template>
        </el-table-column>
        <el-table-column prop="mime" label="MIME" width="150" />
        <el-table-column prop="uploader" label="上传者" width="120" />
        <el-table-column prop="created_at" label="上传时间" width="180" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handlePreview(row)">预览</el-button>
            <el-button type="success" link @click="handleDownload(row)">下载</el-button>
            <el-popconfirm title="确定删除吗？" @confirm="handleDelete(row)">
              <template #reference><el-button type="danger" link>删除</el-button></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size"
          :page-sizes="[20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <el-dialog v-model="uploadVisible" title="上传文件" width="500px">
      <el-upload drag :action="uploadUrl" :headers="uploadHeaders" :on-success="handleUploadSuccess" :on-error="handleUploadError" multiple>
        <el-icon :size="48"><UploadFilled /></el-icon>
        <div>拖拽文件到此处或 <em>点击上传</em></div>
      </el-upload>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import request from '@/utils/http'

const typeMap: Record<string, string> = { image: '图片', document: '文档', attachment: '附件', other: '其他' }
const loading = ref(false)
const uploadVisible = ref(false)
const searchForm = reactive({ type: '', keyword: '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const uploadUrl = '/api/admin/uploads'
const uploadHeaders = computed(() => ({ Authorization: `Bearer ${localStorage.getItem('token') || ''}` }))

const formatSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/uploads', params: { page: pagination.page, page_size: pagination.page_size, ...searchForm } })
    tableData.value = data.list || []; pagination.total = data.total || 0
  } catch { ElMessage.error('获取文件列表失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.type = ''; searchForm.keyword = ''; handleSearch() }
const handleUpload = () => { uploadVisible.value = true }
const handlePreview = (row: any) => { window.open(row.url, '_blank') }
const handleDownload = (row: any) => { window.open(`/api/admin/uploads/${row.id}/download`, '_blank') }
const handleDelete = async (row: any) => {
  try { await request.del({ url: `/api/admin/uploads/${row.id}` }); ElMessage.success('删除成功'); fetchData() } catch { ElMessage.error('删除失败') }
}
const handleUploadSuccess = () => { ElMessage.success('上传成功'); uploadVisible.value = false; fetchData() }
const handleUploadError = () => { ElMessage.error('上传失败') }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.uploads-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
