<template>
  <div class="attachments-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsAttachments.title') }}</span>
          <el-button type="primary" @click="handleUpload">
            <el-icon><Upload /></el-icon>
            {{ $t('clientsAttachments.upload') }}
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('clientsAttachments.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('clientsAttachments.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('common.type')">
          <el-select v-model="searchForm.type" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('clientsAttachments.image')" value="image" />
            <el-option :label="$t('clientsAttachments.document')" value="document" />
            <el-option :label="$t('clientsAttachments.archive')" value="archive" />
            <el-option :label="$t('clientsAttachments.otherType')" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="filename" :label="$t('clientsAttachments.filename')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="size" :label="$t('clientsAttachments.size')" width="120" align="center">
          <template #default="{ row }">{{ formatSize(row.size) }}</template>
        </el-table-column>
        <el-table-column prop="remark" :label="$t('clientsAttachments.remark')" width="150" show-overflow-tooltip />
        <el-table-column prop="uploader" :label="$t('clientsAttachments.uploader')" width="120" />
        <el-table-column prop="created_at" :label="$t('clientsAttachments.uploadTime')" width="170" />
        <el-table-column :label="$t('common.action')" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleDownload(row)">{{ $t('common.download') }}</el-button>
            <el-button type="primary" link @click="handleEditRemark(row)">{{ $t('common.edit') }}</el-button>
            <el-popconfirm :title="$t('clientsAttachments.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('common.delete') }}</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <el-dialog v-model="uploadVisible" :title="$t('clientsAttachments.upload')" width="500px">
      <el-upload drag :auto-upload="false" :on-change="handleFileChange" :file-list="uploadFileList" :limit="1">
        <el-icon class="el-icon--upload"><Upload /></el-icon>
        <div class="el-upload__text">{{ $t('clientsAttachments.dragOrClick') }}</div>
      </el-upload>
      <el-form label-width="80px" style="margin-top: 16px">
        <el-form-item :label="$t('clientsAttachments.remark')">
          <el-input v-model="uploadRemark" :placeholder="$t('clientsAttachments.enterRemark')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="uploadVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleUploadSubmit" :loading="uploadLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="remarkVisible" :title="$t('clientsAttachments.editRemark')" width="400px">
      <el-form label-width="80px">
        <el-form-item :label="$t('clientsAttachments.remark')">
          <el-input v-model="editRemark" :placeholder="$t('clientsAttachments.enterRemark')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="remarkVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleRemarkSubmit" :loading="remarkLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Upload } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { UploadFile } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const uploadLoading = ref(false)
const remarkLoading = ref(false)
const searchForm = reactive({ keyword: '', type: '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const uploadVisible = ref(false)
const uploadFileList = ref<UploadFile[]>([])
const uploadRemark = ref('')
const selectedFile = ref<File | null>(null)

const remarkVisible = ref(false)
const editRemarkId = ref<number>()
const editRemark = ref('')

const formatSize = (bytes: number) => {
  if (!bytes) return '-'
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0
  let size = bytes
  while (size >= 1024 && i < units.length - 1) { size /= 1024; i++ }
  return `${size.toFixed(2)} ${units[i]}`
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/attachments', params: { page: pagination.page, page_size: pagination.page_size, ...searchForm } })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (e) { ElMessage.error($t('common.fetchFailed')) } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.type = ''; handleSearch() }

const handleUpload = () => { uploadFileList.value = []; uploadRemark.value = ''; selectedFile.value = null; uploadVisible.value = true }

const handleFileChange = (file: UploadFile) => { selectedFile.value = file.raw || null }

const handleUploadSubmit = async () => {
  if (!selectedFile.value) { ElMessage.warning($t('clientsAttachments.selectFile')); return }
  uploadLoading.value = true
  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)
    if (uploadRemark.value) formData.append('remark', uploadRemark.value)
    await request.post({ url: '/api/admin/attachments/upload', data: formData, headers: { 'Content-Type': 'multipart/form-data' } })
    ElMessage.success($t('common.uploadSuccess'))
    uploadVisible.value = false
    fetchData()
  } catch (e) { ElMessage.error($t('common.uploadFailed')) } finally { uploadLoading.value = false }
}

const handleDownload = async (row: any) => {
  try {
    const res = await request.get({ url: `/api/admin/attachments/${row.id}/download`, responseType: 'blob' as any })
    const url = window.URL.createObjectURL(new Blob([res]))
    const link = document.createElement('a')
    link.href = url
    link.download = row.filename || 'file'
    link.click()
    window.URL.revokeObjectURL(url)
  } catch (e) { ElMessage.error($t('common.downloadFailed')) }
}

const handleEditRemark = (row: any) => { editRemarkId.value = row.id; editRemark.value = row.remark || ''; remarkVisible.value = true }

const handleRemarkSubmit = async () => {
  remarkLoading.value = true
  try {
    await request.put({ url: `/api/admin/attachments/${editRemarkId.value}`, params: { remark: editRemark.value } })
    ElMessage.success($t('common.updateSuccess'))
    remarkVisible.value = false
    fetchData()
  } catch (e) { ElMessage.error($t('common.updateFailed')) } finally { remarkLoading.value = false }
}

const handleDelete = async (row: any) => {
  try { await request.del({ url: `/api/admin/attachments/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchData() } catch (e) { ElMessage.error($t('common.deleteFailed')) }
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }
onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.attachments-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
