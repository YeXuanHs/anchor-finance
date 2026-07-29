<template>
  <div class="notification-page page-container">
    <div class="art-card">
      <div class="table-header">
        <h3>通知设置</h3>
      </div>

      <el-table :data="channels" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="渠道名称" width="160">
          <template #default="{ row }">
            <div class="channel-name">
              <el-icon :size="18" :style="{ color: row.color }">
                <component :is="row.icon" />
              </el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small" effect="plain">{{ row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="说明" show-overflow-tooltip />
        <el-table-column prop="enabled" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small" effect="light">
              {{ row.enabled ? '已启用' : '已禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">配置</el-button>
            <el-button
              :type="row.enabled ? 'danger' : 'success'"
              link
              @click="handleToggle(row)"
            >
              {{ row.enabled ? '禁用' : '启用' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 配置对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="`${editingChannel?.name} 配置`"
      width="580px"
      destroy-on-close
    >
      <el-form
        ref="dialogFormRef"
        :model="dialogForm"
        :rules="dialogRules"
        label-width="120px"
      >
        <!-- 邮件配置 -->
        <template v-if="editingChannel?.id === 'email'">
          <el-form-item label="SMTP服务器" prop="smtp_host">
            <el-input v-model="dialogForm.smtp_host" placeholder="如 smtp.qq.com" clearable />
          </el-form-item>
          <el-form-item label="SMTP端口" prop="smtp_port">
            <el-input-number v-model="dialogForm.smtp_port" :min="1" :max="65535" />
          </el-form-item>
          <el-form-item label="发件人邮箱" prop="from_address">
            <el-input v-model="dialogForm.from_address" placeholder="noreply@example.com" clearable />
          </el-form-item>
          <el-form-item label="邮箱密码" prop="smtp_password">
            <el-input
              v-model="dialogForm.smtp_password"
              type="password"
              placeholder="SMTP授权码"
              show-password
              clearable
            />
          </el-form-item>
          <el-form-item label="加密方式">
            <el-select v-model="dialogForm.smtp_encryption" style="width: 100%">
              <el-option label="SSL" value="ssl" />
              <el-option label="TLS" value="tls" />
              <el-option label="无" value="none" />
            </el-select>
          </el-form-item>
        </template>

        <!-- 短信配置 -->
        <template v-if="editingChannel?.id === 'sms'">
          <el-form-item label="服务商">
            <el-select v-model="dialogForm.sms_provider" style="width: 100%">
              <el-option label="阿里云" value="aliyun" />
              <el-option label="腾讯云" value="tencent" />
            </el-select>
          </el-form-item>
          <el-form-item label="API Key" prop="api_key">
            <el-input v-model="dialogForm.api_key" placeholder="请输入API Key" clearable />
          </el-form-item>
          <el-form-item label="API Secret" prop="api_secret">
            <el-input
              v-model="dialogForm.api_secret"
              type="password"
              placeholder="请输入API Secret"
              show-password
              clearable
            />
          </el-form-item>
        </template>

        <!-- 微信配置 -->
        <template v-if="editingChannel?.id === 'wechat'">
          <el-form-item label="App ID" prop="app_id">
            <el-input v-model="dialogForm.app_id" placeholder="请输入App ID" clearable />
          </el-form-item>
          <el-form-item label="App Secret" prop="app_secret">
            <el-input
              v-model="dialogForm.app_secret"
              type="password"
              placeholder="请输入App Secret"
              show-password
              clearable
            />
          </el-form-item>
          <el-form-item label="模板ID" prop="template_id">
            <el-input v-model="dialogForm.template_id" placeholder="请输入模板消息ID" clearable />
          </el-form-item>
        </template>

        <!-- Webhook配置 -->
        <template v-if="editingChannel?.id === 'webhook'">
          <el-form-item label="Webhook URL" prop="webhook_url">
            <el-input v-model="dialogForm.webhook_url" placeholder="https://..." clearable />
          </el-form-item>
          <el-form-item label="请求方式">
            <el-select v-model="dialogForm.webhook_method" style="width: 100%">
              <el-option label="POST" value="POST" />
              <el-option label="GET" value="GET" />
            </el-select>
          </el-form-item>
          <el-form-item label="请求头">
            <el-input
              v-model="dialogForm.webhook_headers"
              type="textarea"
              :rows="3"
              placeholder='{"Content-Type": "application/json"}'
            />
          </el-form-item>
          <el-form-item label="Secret">
            <el-input
              v-model="dialogForm.webhook_secret"
              type="password"
              placeholder="用于签名验证（可选）"
              show-password
              clearable
            />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="dialogLoading" @click="handleTest">测试</el-button>
        <el-button type="primary" :loading="dialogLoading" @click="handleDialogSubmit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import request from '@/utils/request'

interface NotificationChannel {
  id: string
  name: string
  type: string
  icon: string
  color: string
  description: string
  enabled: boolean
  config: Record<string, any>
}

const loading = ref(false)
const dialogVisible = ref(false)
const dialogLoading = ref(false)
const dialogFormRef = ref<FormInstance>()
const editingChannel = ref<NotificationChannel | null>(null)

const channels = ref<NotificationChannel[]>([])

const dialogForm = reactive<Record<string, any>>({
  smtp_host: '',
  smtp_port: 465,
  from_address: '',
  smtp_password: '',
  smtp_encryption: 'ssl',
  sms_provider: 'aliyun',
  api_key: '',
  api_secret: '',
  app_id: '',
  app_secret: '',
  template_id: '',
  webhook_url: '',
  webhook_method: 'POST',
  webhook_headers: '',
  webhook_secret: ''
})

const dialogRules: FormRules = {
  smtp_host: [{ required: true, message: '请输入SMTP服务器', trigger: 'blur' }],
  from_address: [{ required: true, message: '请输入发件人邮箱', trigger: 'blur' }],
  smtp_password: [{ required: true, message: '请输入邮箱密码', trigger: 'blur' }],
  api_key: [{ required: true, message: '请输入API Key', trigger: 'blur' }],
  api_secret: [{ required: true, message: '请输入API Secret', trigger: 'blur' }],
  app_id: [{ required: true, message: '请输入App ID', trigger: 'blur' }],
  app_secret: [{ required: true, message: '请输入App Secret', trigger: 'blur' }],
  webhook_url: [{ required: true, message: '请输入Webhook URL', trigger: 'blur' }]
}

const defaultChannels: NotificationChannel[] = [
  {
    id: 'email', name: '邮件通知', type: '邮件', icon: 'Message', color: '#409EFF',
    description: '通过SMTP发送邮件通知', enabled: false, config: {}
  },
  {
    id: 'sms', name: '短信通知', type: '短信', icon: 'Iphone', color: '#67C23A',
    description: '通过短信服务商发送通知', enabled: false, config: {}
  },
  {
    id: 'wechat', name: '微信通知', type: '微信', icon: 'ChatDotRound', color: '#07C160',
    description: '通过微信公众号模板消息推送', enabled: false, config: {}
  },
  {
    id: 'webhook', name: 'Webhook', type: 'Webhook', icon: 'Link', color: '#E6A23C',
    description: '通过HTTP回调推送通知', enabled: false, config: {}
  }
]

const fetchChannels = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/admin/system/notification')
    if (data?.data?.length) {
      channels.value = defaultChannels.map(dc => {
        const saved = data.data.find((c: any) => c.id === dc.id)
        return saved ? { ...dc, ...saved, config: saved.config || {} } : dc
      })
    } else {
      channels.value = [...defaultChannels]
    }
  } catch {
    channels.value = [...defaultChannels]
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: NotificationChannel) => {
  editingChannel.value = row
  Object.keys(dialogForm).forEach(key => {
    dialogForm[key] = row.config?.[key] ?? dialogForm[key]
  })
  dialogVisible.value = true
}

const handleToggle = async (row: NotificationChannel) => {
  row.enabled = !row.enabled
  try {
    await request.post(`/api/admin/system/notification/${row.id}/toggle`, { enabled: row.enabled })
    ElMessage.success(`${row.name}已${row.enabled ? '启用' : '禁用'}`)
  } catch {
    row.enabled = !row.enabled
    ElMessage.error('操作失败')
  }
}

const handleTest = async () => {
  if (!editingChannel.value) return
  dialogLoading.value = true
  try {
    await request.post(`/api/admin/system/notification/${editingChannel.value.id}/test`)
    ElMessage.success('测试消息已发送')
  } catch {
    ElMessage.error('测试发送失败，请检查配置')
  } finally {
    dialogLoading.value = false
  }
}

const handleDialogSubmit = async () => {
  if (!editingChannel.value) return

  dialogLoading.value = true
  try {
    await request.post(`/api/admin/system/notification/${editingChannel.value.id}`, { config: { ...dialogForm } })
    editingChannel.value.config = { ...dialogForm }
    dialogVisible.value = false
    ElMessage.success('配置已保存')
  } catch {
    ElMessage.error('保存失败，请重试')
  } finally {
    dialogLoading.value = false
  }
}

onMounted(() => {
  fetchChannels()
})
</script>

<style scoped lang="scss">
.notification-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
    }
  }

  .channel-name {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 500;
  }
}
</style>
