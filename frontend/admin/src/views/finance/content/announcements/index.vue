<template>
  <div class="announcements-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('announcement.title') }}</span>
          <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('announcement.addAnnouncement') }}</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('announcement.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('announcement.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('announcement.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('announcement.draft')" value="draft" />
            <el-option :label="$t('announcement.published')" value="published" />
            <el-option :label="$t('announcement.revoked')" value="revoked" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" :label="$t('announcement.announcementTitle')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="author" :label="$t('announcement.author')" width="100" />
        <el-table-column prop="status" :label="$t('announcement.status')" width="100">
          <template #default="{ row }"><el-tag :type="statusTagType(row.status)" size="small">{{ statusText(row.status) }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="publish_at" :label="$t('announcement.publishAt')" width="180" />
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180" />
        <el-table-column :label="$t('announcement.operations')" width="280" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button :type="row.status === 'published' ? 'warning' : 'success'" link @click="handleTogglePublish(row)">{{ row.status === 'published' ? $t('announcement.revoke') : $t('announcement.publish') }}</el-button>
            <el-popconfirm :title="$t('announcement.confirmDelete')" @confirm="handleDelete(row)"><template #reference><el-button type="danger" link>{{ $t('common.delete') }}</el-button></template></el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('announcement.announcementTitle')" prop="title">
          <el-input v-model="formData.title" :placeholder="$t('announcement.enterTitle')" />
        </el-form-item>
        <el-form-item :label="$t('announcement.summary')" prop="summary">
          <el-input v-model="formData.summary" type="textarea" :rows="2" :placeholder="$t('announcement.enterSummary')" />
        </el-form-item>
        <el-form-item :label="$t('announcement.content')" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="8" :placeholder="$t('announcement.enterContent')" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('announcement.author')" prop="author">
              <el-input v-model="formData.author" :placeholder="$t('announcement.enterAuthor')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('announcement.sort')" prop="sort">
              <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('announcement.isTop')" prop="is_top">
          <el-switch v-model="formData.is_top" :active-value="1" :inactive-value="0" :active-text="$t('common.yes')" :inactive-text="$t('common.no')" />
        </el-form-item>
        <el-form-item :label="$t('announcement.status')" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio value="draft">{{ $t('announcement.draft') }}</el-radio>
            <el-radio value="published">{{ $t('announcement.published') }}</el-radio>
          </el-radio-group>
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
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'AnnouncementsManage' })

const loading = ref(false)
const submitLoading = ref(false)
const searchForm = reactive({ keyword: '', status: '' as string })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref($t('announcement.addAnnouncement'))
const formRef = ref<FormInstance>()

const formData = reactive({ id: undefined as number | undefined, title: '', summary: '', content: '', author: '', sort: 0, is_top: 0 as number, status: 'draft' as string })

const formRules: FormRules = {
  title: [{ required: true, message: () => $t('announcement.enterTitle'), trigger: 'blur' }],
  content: [{ required: true, message: () => $t('announcement.enterContent'), trigger: 'blur' }]
}

const statusText = (status: string) => { const map: Record<string, () => string> = { draft: () => $t('announcement.draft'), published: () => $t('announcement.published'), revoked: () => $t('announcement.revoked') }; return map[status]?.() || status }
const statusTagType = (status: string) => { const map: Record<string, any> = { draft: 'info', published: 'success', revoked: 'warning' }; return map[status] || 'info' }

const fetchData = async () => {
  loading.value = true
  try { const data = await request.get({ url: '/api/admin/announcements', params: { page: pagination.page, page_size: pagination.page_size, ...searchForm } }); tableData.value = data.list || []; pagination.total = data.total || 0 } catch { ElMessage.error($t('announcement.fetchFailed')) } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.status = ''; handleSearch() }
const resetForm = () => { formData.id = undefined; formData.title = ''; formData.summary = ''; formData.content = ''; formData.author = ''; formData.sort = 0; formData.is_top = 0; formData.status = 'draft' }
const handleAdd = () => { dialogTitle.value = $t('announcement.addAnnouncement'); resetForm(); dialogVisible.value = true }
const handleEdit = (row: any) => { dialogTitle.value = $t('announcement.editAnnouncement'); Object.assign(formData, row); dialogVisible.value = true }

const handleTogglePublish = async (row: any) => {
  try { const newStatus = row.status === 'published' ? 'revoked' : 'published'; await request.put({ url: `/api/admin/announcements/${row.id}`, params: { status: newStatus } }); ElMessage.success(newStatus === 'published' ? $t('announcement.publishSuccess') : $t('announcement.revokeSuccess')); fetchData() } catch { ElMessage.error($t('common.operateFailed')) }
}

const handleDelete = async (row: any) => { try { await request.del({ url: `/api/admin/announcements/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchData() } catch { ElMessage.error($t('common.deleteFailed')) } }

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/announcements/${formData.id}` : '/api/admin/announcements'
      if (formData.id) await request.put({ url, params: formData })
      else await request.post({ url, params: formData })
      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess')); dialogVisible.value = false; fetchData()
    } catch { ElMessage.error($t('common.operateFailed')) } finally { submitLoading.value = false }
  })
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.announcements-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
