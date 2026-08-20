<template>
  <div class="email-view-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsEmailView.title') }}</span>
          <el-button type="primary" @click="handleSendEmail">
            <el-icon><Promotion /></el-icon>
            {{ $t('clientsEmailView.sendEmail') }}
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('clientsEmailView.direction')">
          <el-select v-model="searchForm.direction" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('clientsEmailView.inbound')" value="inbound" />
            <el-option :label="$t('clientsEmailView.outbound')" value="outbound" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('clientsEmailView.read')" :value="1" />
            <el-option :label="$t('clientsEmailView.unread')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('clientsEmailView.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('clientsEmailView.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('common.dateRange')">
          <el-date-picker v-model="searchForm.date_range" type="daterange" :range-separator="$t('common.to')" :start-placeholder="$t('common.startDate')" :end-placeholder="$t('common.endDate')" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="subject" :label="$t('clientsEmailView.subject')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="from_address" :label="$t('clientsEmailView.sender')" width="180" show-overflow-tooltip />
        <el-table-column prop="to_address" :label="$t('clientsEmailView.receiver')" width="180" show-overflow-tooltip />
        <el-table-column prop="direction" :label="$t('clientsEmailView.direction')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.direction === 'inbound' ? 'success' : 'primary'" size="small">{{ row.direction === 'inbound' ? $t('clientsEmailView.inbound') : $t('clientsEmailView.outbound') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('common.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">{{ row.status === 1 ? $t('clientsEmailView.read') : $t('clientsEmailView.unread') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="has_attachment" :label="$t('clientsEmailView.attachment')" width="80" align="center">
          <template #default="{ row }"><el-icon v-if="row.has_attachment" color="#409EFF"><Paperclip /></el-icon></template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.time')" width="170" />
        <el-table-column :label="$t('common.action')" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">{{ $t('common.view') }}</el-button>
            <el-button type="success" link @click="handleDownload(row)" v-if="row.has_attachment">{{ $t('common.download') }}</el-button>
            <el-popconfirm :title="$t('clientsEmailView.confirmDelete')" @confirm="handleDelete(row)">
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

    <el-dialog v-model="detailVisible" :title="$t('clientsEmailView.emailDetail')" width="700px">
      <el-descriptions :column="1" border>
        <el-descriptions-item :label="$t('clientsEmailView.sender')">{{ detailData.from_address }}</el-descriptions-item>
        <el-descriptions-item :label="$t('clientsEmailView.receiver')">{{ detailData.to_address }}</el-descriptions-item>
        <el-descriptions-item :label="$t('clientsEmailView.subject')">{{ detailData.subject }}</el-descriptions-item>
        <el-descriptions-item :label="$t('common.time')">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('clientsEmailView.content')">
          <div class="email-content" v-html="detailData.content"></div>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button type="primary" @click="handleReply(detailData)">{{ $t('clientsEmailView.reply') }}</el-button>
        <el-button @click="detailVisible = false">{{ $t('common.close') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="sendVisible" :title="$t('clientsEmailView.sendEmail')" width="600px">
      <el-form :model="sendForm" :rules="sendRules" ref="sendFormRef" label-width="80px">
        <el-form-item :label="$t('clientsEmailView.receiver')" prop="to_address">
          <el-input v-model="sendForm.to_address" :placeholder="$t('clientsEmailView.enterReceiver')" />
        </el-form-item>
        <el-form-item :label="$t('clientsEmailView.subject')" prop="subject">
          <el-input v-model="sendForm.subject" :placeholder="$t('clientsEmailView.enterSubject')" />
        </el-form-item>
        <el-form-item :label="$t('clientsEmailView.content')" prop="content">
          <el-input v-model="sendForm.content" type="textarea" :rows="6" :placeholder="$t('clientsEmailView.enterEmailContent')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="sendVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSendSubmit" :loading="sendLoading">{{ $t('clientsEmailView.send') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Promotion, Paperclip } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const route = useRoute()
const clientId = route.params.id as string

const loading = ref(false)
const sendLoading = ref(false)
const searchForm = reactive({ direction: '', status: undefined as number | undefined, keyword: '', date_range: [] as string[] })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const detailVisible = ref(false)
const detailData = ref<any>({})
const sendVisible = ref(false)
const sendFormRef = ref<FormInstance>()
const sendForm = reactive({ to_address: '', subject: '', content: '' })

const sendRules: FormRules = {
  to_address: [{ required: true, message: $t('clientsEmailView.enterReceiver'), trigger: 'blur' }],
  subject: [{ required: true, message: $t('clientsEmailView.enterSubject'), trigger: 'blur' }],
  content: [{ required: true, message: $t('clientsEmailView.enterEmailContent'), trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size, client_id: clientId }
    if (searchForm.direction) params.direction = searchForm.direction
    if (searchForm.status !== undefined) params.status = searchForm.status
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.date_range?.length === 2) { params.start_date = searchForm.date_range[0]; params.end_date = searchForm.date_range[1] }
    const data = await request.get({ url: '/api/admin/email-logs', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (e) { ElMessage.error($t('common.fetchFailed')) } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.direction = ''; searchForm.status = undefined; searchForm.keyword = ''; searchForm.date_range = []; handleSearch() }

const handleView = (row: any) => { detailData.value = { ...row }; detailVisible.value = true }

const handleReply = (row: any) => {
  sendForm.to_address = row.from_address || ''
  sendForm.subject = `Re: ${row.subject || ''}`
  sendForm.content = ''
  detailVisible.value = false
  sendVisible.value = true
}

const handleSendEmail = () => { sendForm.to_address = ''; sendForm.subject = ''; sendForm.content = ''; sendVisible.value = true }

const handleSendSubmit = async () => {
  if (!sendFormRef.value) return
  await sendFormRef.value.validate(async (valid) => {
    if (!valid) return
    sendLoading.value = true
    try {
      await request.post({ url: `/api/admin/clients/${clientId}/emails`, params: sendForm })
      ElMessage.success($t('common.sendSuccess'))
      sendVisible.value = false
      fetchData()
    } catch (e) { ElMessage.error($t('common.sendFailed')) } finally { sendLoading.value = false }
  })
}

const handleDownload = async (row: any) => {
  try {
    const res = await request.get({ url: `/api/admin/email-logs/${row.id}/attachments`, responseType: 'blob' as any })
    const url = window.URL.createObjectURL(new Blob([res]))
    const link = document.createElement('a')
    link.href = url
    link.download = row.attachment_name || 'attachment'
    link.click()
    window.URL.revokeObjectURL(url)
  } catch (e) { ElMessage.error($t('common.downloadFailed')) }
}

const handleDelete = async (row: any) => {
  try { await request.del({ url: `/api/admin/email-logs/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchData() } catch (e) { ElMessage.error($t('common.deleteFailed')) }
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }
onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.email-view-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.email-content { max-height: 400px; overflow-y: auto; line-height: 1.6; }
</style>
