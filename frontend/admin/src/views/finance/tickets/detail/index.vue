<template>
  <div class="ticket-detail-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>工单详情</span>
          <div class="header-actions">
            <el-button @click="handleBack">
              <el-icon><Back /></el-icon>
              返回
            </el-button>
          </div>
        </div>
      </template>

      <div v-loading="loading" class="loading-container">
        <div v-if="ticket" class="ticket-content">
          <!-- 工单基本信息 -->
          <div class="ticket-header">
            <div class="ticket-title">
              <h2>{{ ticket.subject }}</h2>
              <div class="ticket-badges">
                <el-tag :type="statusTypeMap[ticket.status]" size="default">
                  {{ statusLabelMap[ticket.status] }}
                </el-tag>
                <el-tag :type="priorityTypeMap[ticket.priority]" size="default">
                  {{ priorityLabelMap[ticket.priority] }}
                </el-tag>
              </div>
            </div>
            <div class="ticket-meta">
              <span>ID: {{ ticket.id }}</span>
              <span>创建时间: {{ ticket.created_at }}</span>
              <span>最后更新: {{ ticket.updated_at }}</span>
              <span>回复数: {{ ticket.reply_count }}</span>
            </div>
          </div>

          <el-divider />

          <!-- 关联信息 -->
          <el-row :gutter="20" class="related-info">
            <el-col :span="12">
              <el-card shadow="never" class="info-card">
                <template #header>
                  <span>关联客户</span>
                </template>
                <el-descriptions :column="1" size="small">
                  <el-descriptions-item label="客户ID">{{ ticket.client_id }}</el-descriptions-item>
                  <el-descriptions-item label="用户名">{{ ticket.client_username }}</el-descriptions-item>
                  <el-descriptions-item label="邮箱">{{ ticket.client_email || '-' }}</el-descriptions-item>
                  <el-descriptions-item>
                    <el-button type="primary" link @click="handleViewClient">
                      查看客户详情
                    </el-button>
                  </el-descriptions-item>
                </el-descriptions>
              </el-card>
            </el-col>
            <el-col :span="12">
              <el-card shadow="never" class="info-card">
                <template #header>
                  <span>工单信息</span>
                </template>
                <el-descriptions :column="1" size="small">
                  <el-descriptions-item label="部门">{{ ticket.department_name }}</el-descriptions-item>
                  <el-descriptions-item label="关联产品">{{ ticket.product_name || '无' }}</el-descriptions-item>
                  <el-descriptions-item label="关联服务">{{ ticket.service_name || '无' }}</el-descriptions-item>
                </el-descriptions>
              </el-card>
            </el-col>
          </el-row>

          <el-divider />

          <!-- 工单内容 -->
          <div class="ticket-body">
            <h3>工单内容</h3>
            <div class="ticket-content-text">{{ ticket.content }}</div>
            <!-- 附件 -->
            <div class="ticket-attachments" v-if="ticket.attachments && ticket.attachments.length">
              <div class="attachment-title">附件：</div>
              <div class="attachment-list">
                <div v-for="(file, index) in ticket.attachments" :key="index" class="attachment-item">
                  <el-link type="primary" :href="file.url" target="_blank">
                    <el-icon><Document /></el-icon>
                    {{ file.file_name }}
                    <span class="file-size">({{ formatFileSize(file.file_size) }})</span>
                  </el-link>
                </div>
              </div>
            </div>
          </div>

          <el-divider />

          <!-- 回复/对话列表 -->
          <div class="reply-section">
            <h3>对话记录</h3>
            <div class="reply-list" v-if="ticket.replies && ticket.replies.length">
              <div
                class="reply-item"
                v-for="reply in ticket.replies"
                :key="reply.id"
                :class="{ 'is-admin': reply.is_admin }"
              >
                <div class="reply-header">
                  <div class="reply-user-info">
                    <el-avatar :size="32">{{ reply.username?.charAt(0) }}</el-avatar>
                    <div class="reply-meta">
                      <span class="reply-user">
                        {{ reply.is_admin ? '管理员' : '客户' }}: {{ reply.username }}
                      </span>
                      <span class="reply-time">{{ reply.created_at }}</span>
                    </div>
                  </div>
                </div>
                <div class="reply-content">{{ reply.content }}</div>
                <!-- 回复附件 -->
                <div class="reply-attachments" v-if="reply.attachments && reply.attachments.length">
                  <div class="attachment-list">
                    <div v-for="(file, index) in reply.attachments" :key="index" class="attachment-item">
                      <el-link type="primary" :href="file.url" target="_blank">
                        <el-icon><Document /></el-icon>
                        {{ file.file_name }}
                        <span class="file-size">({{ formatFileSize(file.file_size) }})</span>
                      </el-link>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <el-empty v-else description="暂无回复" />
          </div>

          <el-divider />

          <!-- 内部备注 -->
          <div class="internal-notes" v-if="ticket.internal_notes && ticket.internal_notes.length">
            <h3>内部备注</h3>
            <div class="notes-list">
              <div class="note-item" v-for="note in ticket.internal_notes" :key="note.id">
                <div class="note-header">
                  <span class="note-user">{{ note.username }}</span>
                  <span class="note-time">{{ note.created_at }}</span>
                </div>
                <div class="note-content">{{ note.content }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 回复对话框 -->
    <el-dialog v-model="replyDialogVisible" title="回复工单" width="600px">
      <el-form :model="replyForm" ref="replyFormRef" :rules="replyRules" label-width="80px">
        <el-form-item label="回复内容" prop="content">
          <el-input
            v-model="replyForm.content"
            type="textarea"
            :rows="5"
            placeholder="请输入回复内容"
          />
        </el-form-item>
        <el-form-item label="附件">
          <el-upload
            v-model:file-list="replyForm.attachments"
            :action="getUploadUrl"
            :headers="uploadHeaders"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
            multiple
            :limit="5"
            :on-exceed="handleExceed"
          >
            <el-button type="primary" plain>
              <el-icon><Upload /></el-icon>
              选择文件
            </el-button>
            <template #tip>
              <div class="el-upload__tip">
                支持上传图片、文档等文件，单个文件不超过 10MB，最多 5 个
              </div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="replyDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitReply" :loading="replyLoading">提交回复</el-button>
      </template>
    </el-dialog>

    <!-- 转部门对话框 -->
    <el-dialog v-model="transferDialogVisible" title="转部门" width="400px">
      <el-form :model="transferForm" label-width="80px">
        <el-form-item label="目标部门">
          <el-select v-model="transferForm.department_id" placeholder="请选择部门" style="width: 100%">
            <el-option
              v-for="dept in departments"
              :key="dept.id"
              :label="dept.name"
              :value="dept.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="transferForm.remark" type="textarea" :rows="2" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="transferDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleTransfer" :loading="transferLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 底部操作栏 -->
    <div class="action-bar" v-if="ticket && ticket.status !== 'closed'">
      <el-button type="primary" @click="handleOpenReplyDialog">
        <el-icon><ChatLineRound /></el-icon>
        回复
      </el-button>
      <el-button type="warning" @click="handleOpenTransferDialog">
        <el-icon><Switch /></el-icon>
        转部门
      </el-button>
      <el-popconfirm title="确定关闭该工单吗？" @confirm="handleClose">
        <template #reference>
          <el-button type="danger">
            <el-icon><CircleClose /></el-icon>
            关闭工单
          </el-button>
        </template>
      </el-popconfirm>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Back, Document, Upload, ChatLineRound, Switch, CircleClose } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules, UploadFile, UploadProps } from 'element-plus'
import request from '@/utils/http'
import { useUserStore } from '@/store/modules/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const replyLoading = ref(false)
const transferLoading = ref(false)
const ticket = ref<any>({})
const departments = ref<any[]>([])

const replyDialogVisible = ref(false)
const transferDialogVisible = ref(false)
const replyFormRef = ref<FormInstance>()

const uploadHeaders = computed(() => ({
  Authorization: userStore.accessToken
}))

const getUploadUrl = computed(() => {
  if (ticket.value?.id) {
    return `/api/admin/tickets/${ticket.value.id}/attachments`
  }
  return ''
})

const replyForm = reactive({
  content: '',
  attachments: [] as UploadFile[]
})

const replyRules: FormRules = {
  content: [
    { required: true, message: '请输入回复内容', trigger: 'blur' }
  ]
}

const transferForm = reactive({
  department_id: undefined as number | undefined,
  remark: ''
})

const priorityTypeMap: Record<string, string> = {
  low: 'info',
  medium: '',
  high: 'warning',
  urgent: 'danger'
}

const priorityLabelMap: Record<string, string> = {
  low: '低',
  medium: '中',
  high: '高',
  urgent: '紧急'
}

const statusTypeMap: Record<string, string> = {
  open: 'warning',
  in_progress: '',
  replied: 'success',
  closed: 'info'
}

const statusLabelMap: Record<string, string> = {
  open: '待处理',
  in_progress: '处理中',
  replied: '已回复',
  closed: '已关闭'
}

const fetchTicket = async () => {
  const id = route.params.id
  if (!id) return

  loading.value = true
  try {
    const data = await request.get({ url: `/api/admin/tickets/${id}` })
    ticket.value = {
      ...data.ticket,
      attachments: data.attachments || [],
      replies: data.replies || [],
      internal_notes: data.internal_notes || []
    }
  } catch (error) {
    console.error('获取工单详情失败:', error)
    ElMessage.error('获取工单详情失败')
  } finally {
    loading.value = false
  }
}

const fetchDepartments = async () => {
  try {
    const data = await request.get({ url: '/api/admin/ticket-depts' })
    departments.value = data || []
  } catch (error) {
    console.error('获取部门列表失败:', error)
  }
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB!')
    return false
  }
  return true
}

const handleUploadSuccess = (response: any, file: UploadFile) => {
  if (response?.data) {
    file.url = response.data.url
    ;(file as any).attachment_id = response.data.id
  } else {
    ElMessage.error(response?.msg || '上传失败')
  }
}

const handleUploadError = () => {
  ElMessage.error('文件上传失败')
}

const handleExceed = () => {
  ElMessage.warning('最多只能上传 5 个文件')
}

const handleBack = () => {
  router.back()
}

const handleViewClient = () => {
  if (ticket.value?.client_id) {
    router.push(`/finance/clients/detail/${ticket.value.client_id}`)
  }
}

const handleOpenReplyDialog = () => {
  replyForm.content = ''
  replyForm.attachments = []
  replyDialogVisible.value = true
}

const handleSubmitReply = async () => {
  if (!replyFormRef.value || !ticket.value?.id) return

  await replyFormRef.value.validate(async (valid) => {
    if (!valid) return

    replyLoading.value = true
    try {
      const attachmentIds = replyForm.attachments
        .filter((file: any) => file.attachment_id)
        .map((file: any) => file.attachment_id)

      await request.post({
        url: `/api/admin/tickets/${ticket.value.id}/reply`,
        params: {
          content: replyForm.content,
          attachment_ids: attachmentIds
        }
      })
      ElMessage.success('回复成功')
      replyDialogVisible.value = false
      fetchTicket()
    } catch (error) {
      ElMessage.error('回复失败')
    } finally {
      replyLoading.value = false
    }
  })
}

const handleOpenTransferDialog = () => {
  transferForm.department_id = undefined
  transferForm.remark = ''
  transferDialogVisible.value = true
}

const handleTransfer = async () => {
  if (!transferForm.department_id) {
    ElMessage.error('请选择目标部门')
    return
  }
  if (!ticket.value?.id) return

  transferLoading.value = true
  try {
    await request.post({
      url: `/api/admin/tickets/${ticket.value.id}/transfer`,
      params: transferForm
    })
    ElMessage.success('转部门成功')
    transferDialogVisible.value = false
    fetchTicket()
  } catch (error) {
    ElMessage.error('转部门失败')
  } finally {
    transferLoading.value = false
  }
}

const handleClose = async () => {
  if (!ticket.value?.id) return

  try {
    await request.post({ url: `/api/admin/tickets/${ticket.value.id}/close` })
    ElMessage.success('工单已关闭')
    fetchTicket()
  } catch (error) {
    ElMessage.error('关闭失败')
  }
}

onMounted(() => {
  fetchTicket()
  fetchDepartments()
})
</script>

<style scoped lang="scss">
.ticket-detail-page {
  padding: 20px;
  padding-bottom: 80px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.loading-container {
  min-height: 400px;
}

.ticket-header {
  .ticket-title {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;

    h2 {
      margin: 0;
      font-size: 22px;
    }

    .ticket-badges {
      display: flex;
      gap: 8px;
    }
  }

  .ticket-meta {
    display: flex;
    gap: 24px;
    color: var(--el-text-color-secondary);
    font-size: 14px;
  }
}

.related-info {
  margin: 16px 0;

  .info-card {
    height: 100%;
  }
}

.ticket-body {
  h3 {
    margin: 0 0 12px;
    font-size: 16px;
  }

  .ticket-content-text {
    padding: 16px;
    background: var(--el-fill-color-lighter);
    border-radius: 4px;
    line-height: 1.8;
    white-space: pre-wrap;
  }

  .ticket-attachments {
    margin-top: 16px;

    .attachment-title {
      font-weight: 500;
      margin-bottom: 8px;
    }

    .attachment-list {
      display: flex;
      flex-wrap: wrap;
      gap: 12px;
    }

    .attachment-item {
      .file-size {
        color: var(--el-text-color-secondary);
        font-size: 12px;
      }
    }
  }
}

.reply-section {
  h3 {
    margin: 0 0 16px;
    font-size: 16px;
  }
}

.reply-list {
  .reply-item {
    padding: 16px;
    margin-bottom: 12px;
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 8px;

    &.is-admin {
      background: var(--el-color-primary-light-9);
      border-color: var(--el-color-primary-light-7);
    }

    .reply-header {
      margin-bottom: 12px;

      .reply-user-info {
        display: flex;
        align-items: center;
        gap: 12px;

        .reply-meta {
          display: flex;
          flex-direction: column;

          .reply-user {
            font-weight: 500;
            font-size: 14px;
          }

          .reply-time {
            color: var(--el-text-color-secondary);
            font-size: 12px;
          }
        }
      }
    }

    .reply-content {
      line-height: 1.8;
      white-space: pre-wrap;
    }

    .reply-attachments {
      margin-top: 12px;
      padding-top: 12px;
      border-top: 1px dashed var(--el-border-color-lighter);

      .attachment-list {
        display: flex;
        flex-wrap: wrap;
        gap: 12px;
      }
    }
  }
}

.internal-notes {
  h3 {
    margin: 0 0 12px;
    font-size: 16px;
  }

  .notes-list {
    .note-item {
      padding: 12px;
      margin-bottom: 8px;
      background: var(--el-fill-color-lighter);
      border-radius: 4px;

      .note-header {
        display: flex;
        justify-content: space-between;
        margin-bottom: 8px;

        .note-user {
          font-weight: 500;
        }

        .note-time {
          color: var(--el-text-color-secondary);
          font-size: 13px;
        }
      }

      .note-content {
        line-height: 1.6;
        white-space: pre-wrap;
      }
    }
  }
}

.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 16px 20px;
  background: var(--el-bg-color);
  border-top: 1px solid var(--el-border-color-lighter);
  display: flex;
  gap: 12px;
  z-index: 100;
}
</style>