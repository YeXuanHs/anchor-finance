<template>
  <div class="send-message-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.sendMessage.pageTitle') }}</span>
        </div>
      </template>

      <!-- 发送表单 -->
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px" class="send-form">
        <el-form-item :label="$t('finance.sendMessage.messageType')" prop="type">
          <el-radio-group v-model="formData.type">
            <el-radio-button label="sms">{{ $t('finance.sendMessage.sms') }}</el-radio-button>
            <el-radio-button label="email">{{ $t('finance.sendMessage.email') }}</el-radio-button>
            <el-radio-button label="notice">{{ $t('finance.sendMessage.notice') }}</el-radio-button>
          </el-radio-group>
        </el-form-item>

        <el-form-item :label="$t('finance.sendMessage.sendTarget')" prop="target">
          <el-radio-group v-model="formData.target" @change="onTargetChange">
            <el-radio label="all">{{ $t('finance.sendMessage.allClients') }}</el-radio>
            <el-radio label="group">{{ $t('finance.sendMessage.byGroup') }}</el-radio>
            <el-radio label="custom">{{ $t('finance.sendMessage.specifiedUsers') }}</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item :label="$t('finance.sendMessage.selectGroup')" prop="group_ids" v-if="formData.target === 'group'">
          <el-select
            v-model="formData.group_ids"
            multiple
            :placeholder="$t('finance.sendMessage.selectGroupPlaceholder')"
            style="width: 100%"
          >
            <el-option
              v-for="group in groupOptions"
              :key="group.id"
              :label="group.name"
              :value="group.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item :label="$t('finance.sendMessage.specifiedUsers')" prop="client_ids" v-if="formData.target === 'custom'">
          <el-select
            v-model="formData.client_ids"
            multiple
            filterable
            remote
            :remote-method="searchClients"
            :loading="clientSearchLoading"
            :placeholder="$t('finance.sendMessage.searchClientPlaceholder')"
            style="width: 100%"
          >
            <el-option
              v-for="client in clientOptions"
              :key="client.id"
              :label="`${client.username} (${client.email || client.phone || ''})`"
              :value="client.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item :label="$t('finance.sendMessage.emailSubject')" prop="subject" v-if="formData.type === 'email'">
          <el-input v-model="formData.subject" :placeholder="$t('finance.sendMessage.enterEmailSubject')" />
        </el-form-item>

        <el-form-item :label="$t('finance.sendMessage.messageContent')" prop="content">
          <el-input
            v-model="formData.content"
            type="textarea"
            :rows="6"
            :placeholder="contentPlaceholder"
          />
        </el-form-item>

        <el-form-item :label="$t('finance.sendMessage.scheduledSend')" v-if="formData.type !== 'sms'">
          <el-switch v-model="formData.schedule_enabled" />
          <el-date-picker
            v-if="formData.schedule_enabled"
            v-model="formData.schedule_time"
            type="datetime"
            :placeholder="$t('finance.sendMessage.selectSendTime')"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="margin-left: 12px"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSend" :loading="sendLoading">
            {{ formData.schedule_enabled ? $t('finance.sendMessage.scheduledSendBtn') : $t('finance.sendMessage.sendNow') }}
          </el-button>
          <el-button @click="handleResetForm">{{ $t('finance.sendMessage.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 发送记录 -->
    <el-card shadow="never" style="margin-top: 20px">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.sendMessage.sendRecords') }}</span>
          <el-button type="primary" link @click="fetchRecords">
            <el-icon><Refresh /></el-icon>
            {{ $t('finance.sendMessage.refresh') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索 -->
      <el-form :inline="true" :model="recordSearch" class="search-form">
        <el-form-item :label="$t('finance.sendMessage.type')">
          <el-select v-model="recordSearch.type" :placeholder="$t('finance.sendMessage.all')" clearable>
            <el-option :label="$t('finance.sendMessage.sms')" value="sms" />
            <el-option :label="$t('finance.sendMessage.email')" value="email" />
            <el-option :label="$t('finance.sendMessage.notice')" value="notice" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('finance.sendMessage.status')">
          <el-select v-model="recordSearch.status" :placeholder="$t('finance.sendMessage.all')" clearable>
            <el-option :label="$t('finance.sendMessage.statusSending')" value="sending" />
            <el-option :label="$t('finance.sendMessage.statusCompleted')" value="completed" />
            <el-option :label="$t('finance.sendMessage.statusPartial')" value="partial" />
            <el-option :label="$t('finance.sendMessage.statusFailed')" value="failed" />
            <el-option :label="$t('finance.sendMessage.statusPending')" value="pending" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleRecordSearch">{{ $t('finance.sendMessage.search') }}</el-button>
          <el-button @click="handleRecordReset">{{ $t('finance.sendMessage.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="records" v-loading="recordsLoading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="type" :label="$t('finance.sendMessage.type')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="typeTagMap[row.type]" size="small">
              {{ typeLabelMap[row.type] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="subject" :label="$t('finance.sendMessage.subject')" min-width="160" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.subject || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="content" :label="$t('finance.sendMessage.content')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="target_desc" :label="$t('finance.sendMessage.sendTarget')" width="150" show-overflow-tooltip />
        <el-table-column prop="total_count" :label="$t('finance.sendMessage.total')" width="80" align="center" />
        <el-table-column prop="success_count" :label="$t('finance.sendMessage.success')" width="80" align="center">
          <template #default="{ row }">
            <span style="color: #67c23a">{{ row.success_count }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="fail_count" :label="$t('finance.sendMessage.failed')" width="80" align="center">
          <template #default="{ row }">
            <span :style="{ color: row.fail_count > 0 ? '#f56c6c' : '#909399' }">{{ row.fail_count }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('finance.sendMessage.status')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="statusTagMap[row.status]" size="small">
              {{ statusLabelMap[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('finance.sendMessage.createdAt')" width="170" />
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="recordPagination.page"
          v-model:page-size="recordPagination.page_size"
          :page-sizes="[10, 20, 50]"
          :total="recordPagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleRecordSizeChange"
          @current-change="handleRecordPageChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const typeLabelMap: Record<string, string> = {
  sms: $t('finance.sendMessage.sms'),
  email: $t('finance.sendMessage.email'),
  notice: $t('finance.sendMessage.notice')
}

const typeTagMap: Record<string, any> = {
  sms: 'warning',
  email: '',
  notice: 'success'
}

const statusLabelMap: Record<string, string> = {
  sending: $t('finance.sendMessage.statusSending'),
  completed: $t('finance.sendMessage.statusCompleted'),
  partial: $t('finance.sendMessage.statusPartial'),
  failed: $t('finance.sendMessage.statusFailed'),
  pending: $t('finance.sendMessage.statusPending')
}

const statusTagMap: Record<string, any> = {
  sending: 'warning',
  completed: 'success',
  partial: 'warning',
  failed: 'danger',
  pending: 'info'
}

const sendLoading = ref(false)
const recordsLoading = ref(false)
const clientSearchLoading = ref(false)
const formRef = ref<FormInstance>()

const groupOptions = ref<any[]>([])
const clientOptions = ref<any[]>([])

const formData = reactive({
  type: 'notice',
  target: 'all',
  group_ids: [] as number[],
  client_ids: [] as number[],
  subject: '',
  content: '',
  schedule_enabled: false,
  schedule_time: ''
})

const contentPlaceholder = computed(() => {
  const map: Record<string, string> = {
    sms: $t('finance.smsPlaceholder'),
    email: $t('finance.emailPlaceholder'),
    notice: $t('finance.noticePlaceholder')
  }
  return map[formData.type] || $t('finance.sendMessage.enterMessageContent')
})

const formRules: FormRules = {
  type: [{ required: true, message: $t('finance.sendMessage.selectMessageType'), trigger: 'change' }],
  target: [{ required: true, message: $t('finance.sendMessage.selectSendTarget'), trigger: 'change' }],
  group_ids: [
    { required: true, type: 'array', min: 1, message: $t('finance.sendMessage.selectAtLeastOneGroup'), trigger: 'change' }
  ],
  client_ids: [
    { required: true, type: 'array', min: 1, message: $t('finance.sendMessage.selectAtLeastOneUser'), trigger: 'change' }
  ],
  subject: [
    { required: true, message: $t('finance.sendMessage.enterEmailSubject'), trigger: 'blur' }
  ],
  content: [
    { required: true, message: $t('finance.sendMessage.enterMessageContent'), trigger: 'blur' }
  ]
}

const onTargetChange = () => {
  formData.group_ids = []
  formData.client_ids = []
}

// 搜索客户
const searchClients = async (query: string) => {
  if (!query) {
    clientOptions.value = []
    return
  }
  clientSearchLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/users',
      params: { keyword: query, page_size: 30 }
    })
    clientOptions.value = data.list || data || []
  } catch (error) {
    console.error('搜索客户失败:', error)
  } finally {
    clientSearchLoading.value = false
  }
}

// 获取分组列表
const fetchGroups = async () => {
  try {
    const data = await request.get({ url: '/api/admin/client-groups' })
    groupOptions.value = data || []
  } catch (error) {
    console.error('获取分组列表失败:', error)
  }
}

// 发送消息
const handleSend = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    sendLoading.value = true
    try {
      const params: any = {
        type: formData.type,
        target: formData.target,
        content: formData.content
      }

      if (formData.type === 'email' && formData.subject) {
        params.subject = formData.subject
      }
      if (formData.target === 'group') {
        params.group_ids = formData.group_ids
      } else if (formData.target === 'custom') {
        params.client_ids = formData.client_ids
      }
      if (formData.schedule_enabled && formData.schedule_time) {
        params.schedule_time = formData.schedule_time
      }

      await request.post({ url: '/api/admin/messages/send', params })
      ElMessage.success(formData.schedule_enabled ? $t('finance.sendMessage.scheduleSet') : $t('finance.sendMessage.sendSuccess'))
      handleResetForm()
      fetchRecords()
    } catch (error) {
      ElMessage.error($t('finance.sendMessage.sendFailed'))
    } finally {
      sendLoading.value = false
    }
  })
}

// 重置表单
const handleResetForm = () => {
  formData.type = 'notice'
  formData.target = 'all'
  formData.group_ids = []
  formData.client_ids = []
  formData.subject = ''
  formData.content = ''
  formData.schedule_enabled = false
  formData.schedule_time = ''
}

// 发送记录
const recordSearch = reactive({
  type: '',
  status: ''
})

const recordPagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const records = ref<any[]>([])

const fetchRecords = async () => {
  recordsLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/messages/batch/records',
      params: {
        page: recordPagination.page,
        page_size: recordPagination.page_size,
        type: recordSearch.type || undefined,
        status: recordSearch.status || undefined
      }
    })
    records.value = data.list || []
    recordPagination.total = data.total || 0
  } catch (error) {
    console.error('获取发送记录失败:', error)
    ElMessage.error($t('finance.sendMessage.fetchRecordsFailed'))
  } finally {
    recordsLoading.value = false
  }
}

const handleRecordSearch = () => {
  recordPagination.page = 1
  fetchRecords()
}

const handleRecordReset = () => {
  recordSearch.type = ''
  recordSearch.status = ''
  handleRecordSearch()
}

const handleRecordSizeChange = () => {
  recordPagination.page = 1
  fetchRecords()
}

const handleRecordPageChange = () => {
  fetchRecords()
}

onMounted(() => {
  fetchGroups()
  fetchRecords()
})
</script>

<style scoped lang="scss">
.send-message-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.send-form {
  max-width: 700px;
}

.search-form {
  margin-bottom: 16px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
