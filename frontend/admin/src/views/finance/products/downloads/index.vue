<template>
  <div class="downloads-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('download.title') }}</span>
          <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('download.addDownload') }}</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('download.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('download.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('download.linkedProduct')">
          <el-select v-model="searchForm.product_id" :placeholder="$t('common.all')" clearable filterable>
            <el-option v-for="item in productOptions" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" :label="$t('download.fileName')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="product_name" :label="$t('download.linkedProduct')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="file_size" :label="$t('download.fileSize')" width="120" align="center"><template #default="{ row }">{{ formatFileSize(row.file_size) }}</template></el-table-column>
        <el-table-column prop="download_count" :label="$t('download.downloadCount')" width="100" align="center" />
        <el-table-column prop="version" :label="$t('download.version')" width="100" align="center" />
        <el-table-column prop="sort" :label="$t('download.sort')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('download.status')" width="80" align="center">
          <template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">{{ row.status === 1 ? $t('common.enable') : $t('common.disable') }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180" />
        <el-table-column :label="$t('download.operations')" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-popconfirm :title="$t('download.confirmDelete')" @confirm="handleDelete(row)"><template #reference><el-button type="danger" link>{{ $t('common.delete') }}</el-button></template></el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('download.fileName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('download.enterFileName')" />
        </el-form-item>
        <el-form-item :label="$t('download.linkedProduct')" prop="product_id">
          <el-select v-model="formData.product_id" :placeholder="$t('download.selectProduct')" filterable style="width: 100%">
            <el-option v-for="item in productOptions" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('download.version')" prop="version">
          <el-input v-model="formData.version" :placeholder="$t('download.enterVersion')" />
        </el-form-item>
        <el-form-item :label="$t('download.fileUrl')" prop="file_url">
          <el-input v-model="formData.file_url" :placeholder="$t('download.enterFileUrl')">
            <template #append><el-upload :show-file-list="false" :before-upload="handleBeforeUpload" accept="*"><el-button>{{ $t('download.upload') }}</el-button></el-upload></template>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('common.description')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('download.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('download.sort')" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('download.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" :active-text="$t('common.enable')" :inactive-text="$t('common.disable')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules, UploadRawFile } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'ProductDownloads' })

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref($t('download.addDownload'))
const formRef = ref<FormInstance>()
const productOptions = ref<any[]>([])

const searchForm = reactive({ keyword: '', product_id: undefined as number | undefined })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const formData = reactive({ id: undefined as number | undefined, name: '', product_id: undefined as number | undefined, version: '', file_url: '', description: '', sort: 0, status: 1 })

const formRules: FormRules = {
  name: [{ required: true, message: () => $t('download.enterFileName'), trigger: 'blur' }],
  product_id: [{ required: true, message: () => $t('download.selectProduct'), trigger: 'change' }]
}

const formatFileSize = (bytes?: number) => {
  if (!bytes) return '-'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

const handleBeforeUpload = (file: UploadRawFile) => { formData.file_url = URL.createObjectURL(file); formData.name = formData.name || file.name; return false }

const fetchData = async () => {
  loading.value = true
  try { const data = await request.get({ url: '/api/admin/product-downloads', params: { page: pagination.page, page_size: pagination.page_size, ...searchForm } }); tableData.value = data.list || data || []; pagination.total = data.total || 0 } catch { ElMessage.error($t('download.fetchFailed')) } finally { loading.value = false }
}

const fetchProducts = async () => { try { const data = await request.get({ url: '/api/admin/products', params: { page_size: 999 } }); productOptions.value = data.list || data || [] } catch {} }

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.product_id = undefined; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleAdd = () => { dialogTitle.value = $t('download.addDownload'); formData.id = undefined; formData.name = ''; formData.product_id = undefined; formData.version = ''; formData.file_url = ''; formData.description = ''; formData.sort = 0; formData.status = 1; dialogVisible.value = true }
const handleEdit = (row: any) => { dialogTitle.value = $t('download.editDownload'); Object.assign(formData, row); dialogVisible.value = true }

const handleDelete = async (row: any) => { try { await request.del({ url: `/api/admin/product-downloads/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchData() } catch { ElMessage.error($t('common.deleteFailed')) } }

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) await request.put({ url: `/api/admin/product-downloads/${formData.id}`, params: formData })
      else await request.post({ url: '/api/admin/product-downloads', params: formData })
      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess')); dialogVisible.value = false; fetchData()
    } catch { ElMessage.error($t('common.operateFailed')) } finally { submitLoading.value = false }
  })
}

onMounted(() => { fetchData(); fetchProducts() })
</script>

<style scoped lang="scss">
.downloads-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
