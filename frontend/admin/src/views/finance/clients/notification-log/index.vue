<template>
  <div class="notification-log-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsNotificationLog.title') }}</span>
          <el-button type="primary" @click="handleSendNotification">
            <el-icon><Promotion /></el-icon>
            {{ $t('clientsNotificationLog.sendNotification') }}
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('common.type')">
          <el-select v-model="searchForm.type" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('clientsNotificationLog.email')" value="email" />
            <el-option :label="$t('clientsNotificationLog.sms')" value="sms" />
            <el-option :label="$t('clientsNotificationLog.siteMessage')" value="site" />
            <el-option :label="$t('clientsNotificationLog.wechat')" value="wechat" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('clientsNotificationLog.sendSuccess')" :value="1" />
            <el-option :label="$t('clientsNotificationLog.sendFailed')" :value="0" />
          </el-select>
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
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="type" :label="$t('clientsNotificationLog.notifyType')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="title" :label="$t('clientsNotificationLog.titleCol')" min-width="180" show-overflow-tooltip />
        <el-table-column prop="content" :label="$t('clientsNotificationLog.content')" min-width="250" show-overflow-tooltip />
        <el-table-column prop="status" :label="$t('common.status')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('common.success') : $t('common.failed') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="error_msg" :label="$t('clientsNotificationLog.errorMsg')" width="150" show-overflow-tooltip />
        <el-table-column prop="operator" :label="$t('common.operator')" width="100" />
        <el-table-column prop="created_at" :label="$t('clientsNotificationLog.sentTime')" width="170" />
        <el-table-column :label="$t('common.action')" width="100" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">{{ $t('common.detail') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <el-dialog v-model="detailVisible" :title="$t('clientsNotificationLog.detailTitle')" width="650px">
      <el-descriptions :column="2" border>
        <el-descriptions-item :label="$t('clientsNotificationLog.notifyType')">{{ getTypeText(detailData.type) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('common.status')">
          <el-tag :type="detailData.status === 1 ? 'success' : 'danger'" size="small">{{ detailData.status === 1 ? $t('common.success') : $t('common.failed') }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('clientsNotificationLog.titleCol')" :span="2">{{ detailData.title }}</el-descriptions-item>
        <el-descriptions-item :label="$t('clientsNotificationLog.content')" :span="2"><div class="content-text">{{ detailData.content }}</div></el-descriptions-item>
        <el-descriptions-item :label="$t('clientsNotificationLog.receiver')" :span="2">{{ detailData.receiver || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('common.operator')">{{ detailData.operator }}</el-descriptions-item>
        <el-descriptions-item :label="$t('clientsNotificationLog.sentTime')">{{ detailData.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">{{ $t('common.close') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="sendVisible" :title="$t('clientsNotificationLog.sendNotification')" width="600px">
      <el-form :model="sendForm" :rules="sendRules" ref="sendFormRef" label-width="100px">
        <el-form-item :label="$t('clientsNotificationLog.notifyType')" prop="type">
          <el-select v-model="sendForm.type" :placeholder="$t('common.select')">
            <el-option :label="$t('clientsNotificationLog.email')" value="email" />
            <el-option :label="$t('clientsNotificationLog.sms')" value="sms" />
            <el-option :label="$t('clientsNotificationLog.siteMessage')" value="site" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('clientsNotificationLog.titleCol')" prop="title">
          <el-input v-model="sendForm.title" :placeholder="$t('clientsNotificationLog.enterTitle')" />
        </el-form-item>
        <el-form-item :label="$t('clientsNotificationLog.content')" prop="content">
          <el-input v-model="sendForm.content" type="textarea" :rows="5" :placeholder="$t('clientsNotificationLog.enterContent')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="sendVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSendSubmit" :loading="sendLoading">{{ $t('clientsNotificationLog.send') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Promotion } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const route = useRoute()
const clientId = route.params.id as string

const loading = ref(false)
const sendLoading = ref(false)
const searchForm = reactive({ type: '', status: undefined as number | undefined, date_range: [] as string[] })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const detailVisible = ref(false)
const detailData = ref<any>({})
const sendVisible = ref(false)
const sendFormRef = ref<FormInstance>()
const sendForm = reactive({ type: 'email', title: '', content: '' })

const sendRules: FormRules = {
  type: [{ required: true, message: $t('clientsNotificationLog.selectNotifyType'), trigger: 'change' }],
  title: [{ required: true, message: $t('clientsNotificationLog.enterTitle'), trigger: 'blur' }],
  content: [{ required: true, message: $t('clientsNotificationLog.enterContent'), trigger: 'blur' }]
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = { email: $t('clientsNotificationLog.email'), sms: $t('clientsNotificationLog.sms'), site: $t('clientsNotificationLog.siteMessage'), wechat: $t('clientsNotificationLog.wechat') }
  return map[type] || type
}
const getTypeTag = (type: string) => {
  const map: Record<string, string> = { email: 'primary', sms: 'success', site: 'warning', wechat: 'success' }
  return (map[type] || 'info') as any
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size, client_id: clientId }
    if (searchForm.type) params.type = searchForm.type
    if (searchForm.status !== undefined) params.status = searchForm.status
    if (searchForm.date_range?.length === 2) { params.start_date = searchForm.date_range[0]; params.end_date = searchForm.date_range[1] }
    const data = await request.get({ url: '/api/admin/notifications/logs', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (e) { ElMessage.error($t('common.fetchFailed')) } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.type = ''; searchForm.status = undefined; searchForm.date_range = []; handleSearch() }
const handleViewDetail = (row: any) => { detailData.value = { ...row }; detailVisible.value = true }
const handleSendNotification = () => { sendForm.type = 'email'; sendForm.title = ''; sendForm.content = ''; sendVisible.value = true }

const handleSendSubmit = async () => {
  if (!sendFormRef.value) return
  await sendFormRef.value.validate(async (valid) => {
    if (!valid) return
    sendLoading.value = true
    try {
      await request.post({ url: `/api/admin/clients/${clientId}/notifications`, params: sendForm })
      ElMessage.success($t('common.sendSuccess'))
      sendVisible.value = false
      fetchData()
    } catch (e) { ElMessage.error($t('common.sendFailed')) } finally { sendLoading.value = false }
  })
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }
onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.notification-log-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.content-text { white-space: pre-wrap; word-break: break-all; }
</style>
