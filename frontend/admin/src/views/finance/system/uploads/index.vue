<template>
  <div class="uploads-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.uploads.title') }}</span>
          <el-button type="primary" size="small" @click="handleUpload">{{ $t('finance.uploads.uploadFile') }}</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('finance.uploads.fileType')">
          <el-select v-model="searchForm.type" :placeholder="$t('finance.uploads.all')" clearable>
            <el-option :label="$t('finance.uploads.image')" value="image" />
            <el-option :label="$t('finance.uploads.document')" value="document" />
            <el-option :label="$t('finance.uploads.attachment')" value="attachment" />
            <el-option :label="$t('finance.uploads.other')" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('finance.uploads.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('finance.uploads.fileName')" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('finance.uploads.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('finance.uploads.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="$t('finance.uploads.fileName')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="type" :label="$t('finance.uploads.type')" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ typeMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="size" :label="$t('finance.uploads.size')" width="100">
          <template #default="{ row }">{{ formatSize(row.size) }}</template>
        </el-table-column>
        <el-table-column prop="mime" label="MIME" width="150" />
        <el-table-column prop="uploader" :label="$t('finance.uploads.uploader')" width="120" />
        <el-table-column prop="created_at" :label="$t('finance.uploads.uploadTime')" width="180" />
        <el-table-column :label="$t('finance.uploads.actions')" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handlePreview(row)">{{ $t('finance.uploads.preview') }}</el-button>
            <el-button type="success" link @click="handleDownload(row)">{{ $t('finance.uploads.download') }}</el-button>
            <el-popconfirm :title="$t('finance.uploads.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference><el-button type="danger" link>{{ $t('finance.uploads.delete') }}</el-button></template>
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

    <el-dialog v-model="uploadVisible" :title="$t('finance.uploads.uploadFile')" width="500px">
      <el-upload drag :action="uploadUrl" :headers="uploadHeaders" :on-success="handleUploadSuccess" :on-error="handleUploadError" multiple>
        <el-icon :size="48"><UploadFilled /></el-icon>
        <div>{{ $t('finance.uploads.dragHere') }} <em>{{ $t('finance.uploads.clickUpload') }}</em></div>
      </el-upload>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import { $t } from '@/locales'
import request from '@/utils/http'

const typeMap: Record<string, string> = {
  image: $t('finance.uploads.image'),
  document: $t('finance.uploads.document'),
  attachment: $t('finance.uploads.attachment'),
  other: $t('finance.uploads.other')
}
const loading = ref(false)
const uploadVisible = ref(false)
const searchForm = reactive({ type: '', keyword: '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const uploadUrl = '/api/admin/upload'
const uploadHeaders = computed(() => ({ Authorization: `Bearer ${localStorage.getItem('token') || ''}` }))

const formatSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/upload', params: { page: pagination.page, page_size: pagination.page_size, ...searchForm } })
    tableData.value = data.list || []; pagination.total = data.total || 0
  } catch { ElMessage.error($t('finance.uploads.fetchFailed')) } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.type = ''; searchForm.keyword = ''; handleSearch() }
const handleUpload = () => { uploadVisible.value = true }
const handlePreview = (row: any) => { window.open(row.url, '_blank') }
const handleDownload = (row: any) => { window.open(`/api/admin/upload/${row.id}/download`, '_blank') }
const handleDelete = async (row: any) => {
  try { await request.del({ url: `/api/admin/upload/${row.id}` }); ElMessage.success($t('finance.uploads.deleteSuccess')); fetchData() } catch { ElMessage.error($t('finance.uploads.deleteFailed')) }
}
const handleUploadSuccess = () => { ElMessage.success($t('finance.uploads.uploadSuccess')); uploadVisible.value = false; fetchData() }
const handleUploadError = () => { ElMessage.error($t('finance.uploads.uploadFailed')) }
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
