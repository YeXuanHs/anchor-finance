<template>
  <div class="send-message-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>群发消息</span>
        </div>
      </template>

      <!-- 发送表单 -->
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px" class="send-form">
        <el-form-item label="消息类型" prop="type">
          <el-radio-group v-model="formData.type">
            <el-radio-button label="sms">短信</el-radio-button>
            <el-radio-button label="email">邮件</el-radio-button>
            <el-radio-button label="notice">站内信</el-radio-button>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="发送对象" prop="target">
          <el-radio-group v-model="formData.target" @change="onTargetChange">
            <el-radio label="all">全部客户</el-radio>
            <el-radio label="group">按分组</el-radio>
            <el-radio label="custom">指定用户</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="选择分组" prop="group_ids" v-if="formData.target === 'group'">
          <el-select
            v-model="formData.group_ids"
            multiple
            placeholder="请选择客户分组"
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

        <el-form-item label="指定用户" prop="client_ids" v-if="formData.target === 'custom'">
          <el-select
            v-model="formData.client_ids"
            multiple
            filterable
            remote
            :remote-method="searchClients"
            :loading="clientSearchLoading"
            placeholder="输入搜索客户，可多选"
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

        <el-form-item label="邮件标题" prop="subject" v-if="formData.type === 'email'">
          <el-input v-model="formData.subject" placeholder="请输入邮件标题" />
        </el-form-item>

        <el-form-item label="消息内容" prop="content">
          <el-input
            v-model="formData.content"
            type="textarea"
            :rows="6"
            :placeholder="contentPlaceholder"
          />
        </el-form-item>

        <el-form-item label="定时发送" v-if="formData.type !== 'sms'">
          <el-switch v-model="formData.schedule_enabled" />
          <el-date-picker
            v-if="formData.schedule_enabled"
            v-model="formData.schedule_time"
            type="datetime"
            placeholder="选择发送时间"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="margin-left: 12px"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSend" :loading="sendLoading">
            {{ formData.schedule_enabled ? '定时发送' : '立即发送' }}
          </el-button>
          <el-button @click="handleResetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 发送记录 -->
    <el-card shadow="never" style="margin-top: 20px">
      <template #header>
        <div class="card-header">
          <span>发送记录</span>
          <el-button type="primary" link @click="fetchRecords">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <!-- 搜索 -->
      <el-form :inline="true" :model="recordSearch" class="search-form">
        <el-form-item label="类型">
          <el-select v-model="recordSearch.type" placeholder="全部" clearable>
            <el-option label="短信" value="sms" />
            <el-option label="邮件" value="email" />
            <el-option label="站内信" value="notice" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="recordSearch.status" placeholder="全部" clearable>
            <el-option label="发送中" value="sending" />
            <el-option label="已完成" value="completed" />
            <el-option label="部分失败" value="partial" />
            <el-option label="已失败" value="failed" />
            <el-option label="待发送" value="pending" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleRecordSearch">搜索</el-button>
          <el-button @click="handleRecordReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="records" v-loading="recordsLoading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="type" label="类型" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="typeTagMap[row.type]" size="small">
              {{ typeLabelMap[row.type] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="标题" min-width="160" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.subject || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="content" label="内容" min-width="200" show-overflow-tooltip />
        <el-table-column prop="target_desc" label="发送对象" width="150" show-overflow-tooltip />
        <el-table-column prop="total_count" label="总数" width="80" align="center" />
        <el-table-column prop="success_count" label="成功" width="80" align="center">
          <template #default="{ row }">
            <span style="color: #67c23a">{{ row.success_count }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="fail_count" label="失败" width="80" align="center">
          <template #default="{ row }">
            <span :style="{ color: row.fail_count > 0 ? '#f56c6c' : '#909399' }">{{ row.fail_count }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="statusTagMap[row.status]" size="small">
              {{ statusLabelMap[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170" />
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

const typeLabelMap: Record<string, string> = {
  sms: '短信',
  email: '邮件',
  notice: '站内信'
}

const typeTagMap: Record<string, string> = {
  sms: 'warning',
  email: '',
  notice: 'success'
}

const statusLabelMap: Record<string, string> = {
  sending: '发送中',
  completed: '已完成',
  partial: '部分失败',
  failed: '已失败',
  pending: '待发送'
}

const statusTagMap: Record<string, string> = {
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
    sms: '请输入短信内容（70字以内）',
    email: '请输入邮件正文（支持HTML）',
    notice: '请输入站内信内容'
  }
  return map[formData.type] || '请输入消息内容'
})

const formRules: FormRules = {
  type: [{ required: true, message: '请选择消息类型', trigger: 'change' }],
  target: [{ required: true, message: '请选择发送对象', trigger: 'change' }],
  group_ids: [
    { required: true, type: 'array', min: 1, message: '请至少选择一个分组', trigger: 'change' }
  ],
  client_ids: [
    { required: true, type: 'array', min: 1, message: '请至少选择一个用户', trigger: 'change' }
  ],
  subject: [
    { required: true, message: '请输入邮件标题', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入消息内容', trigger: 'blur' }
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
      ElMessage.success(formData.schedule_enabled ? '已设置定时发送' : '发送成功')
      handleResetForm()
      fetchRecords()
    } catch (error) {
      ElMessage.error('发送失败')
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
      url: '/api/admin/messages/records',
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
    ElMessage.error('获取发送记录失败')
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
