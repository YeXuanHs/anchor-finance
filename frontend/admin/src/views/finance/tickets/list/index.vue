<template>
  <div class="ticket-list-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>工单列表</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="工单标题/内容/用户名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待处理" value="open" />
            <el-option label="处理中" value="in_progress" />
            <el-option label="已回复" value="replied" />
            <el-option label="已关闭" value="closed" />
          </el-select>
        </el-form-item>
        <el-form-item label="部门">
          <el-select v-model="searchForm.department_id" placeholder="全部" clearable>
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="优先级">
          <el-select v-model="searchForm.priority" placeholder="全部" clearable>
            <el-option label="低" value="low" />
            <el-option label="中" value="medium" />
            <el-option label="高" value="high" />
            <el-option label="紧急" value="urgent" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="subject" label="工单标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="client_username" label="客户" width="120" />
        <el-table-column prop="department_name" label="部门" width="120" />
        <el-table-column prop="priority" label="优先级" width="80">
          <template #default="{ row }">
            <el-tag :type="priorityTypeMap[row.priority]" size="small">
              {{ priorityLabelMap[row.priority] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTypeMap[row.status]" size="small">
              {{ statusLabelMap[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reply_count" label="回复数" width="80" align="center" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column prop="updated_at" label="最后更新" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">查看</el-button>
            <el-button type="primary" link @click="handleReply(row)">回复</el-button>
            <el-popconfirm
              v-if="row.status !== 'closed'"
              title="确定关闭该工单吗？"
              @confirm="handleClose(row)"
            >
              <template #reference>
                <el-button type="warning" link>关闭</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 工单详情/回复对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="800px" top="5vh">
      <div class="ticket-detail" v-if="currentTicket">
        <!-- 工单信息 -->
        <div class="ticket-info">
          <h3>{{ currentTicket.subject }}</h3>
          <div class="ticket-meta">
            <span>客户: {{ currentTicket.client_username }}</span>
            <span>部门: {{ currentTicket.department_name }}</span>
            <span>优先级: {{ priorityLabelMap[currentTicket.priority] }}</span>
            <span>状态: {{ statusLabelMap[currentTicket.status] }}</span>
          </div>
          <div class="ticket-content">{{ currentTicket.content }}</div>
          <!-- 工单附件 -->
          <div class="ticket-attachments" v-if="currentTicket.attachments && currentTicket.attachments.length">
            <div class="attachment-title">附件：</div>
            <div class="attachment-list">
              <div v-for="(file, index) in currentTicket.attachments" :key="index" class="attachment-item">
                <el-link type="primary" :href="file.url" target="_blank">
                  <el-icon><Document /></el-icon>
                  {{ file.file_name }}
                  <span class="file-size">({{ formatFileSize(file.file_size) }})</span>
                </el-link>
              </div>
            </div>
          </div>
        </div>

        <!-- 回复列表 -->
        <div class="reply-list" v-if="currentTicket.replies && currentTicket.replies.length">
          <div
            class="reply-item"
            v-for="reply in currentTicket.replies"
            :key="reply.id"
            :class="{ 'is-admin': reply.is_admin }"
          >
            <div class="reply-header">
              <span class="reply-user">{{ reply.is_admin ? '管理员' : '客户' }}: {{ reply.username }}</span>
              <span class="reply-time">{{ reply.created_at }}</span>
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

        <!-- 回复表单 -->
        <div class="reply-form" v-if="showReplyForm">
          <el-divider>回复</el-divider>
          <el-form :model="replyForm" ref="replyFormRef">
            <el-form-item prop="content" :rules="[{ required: true, message: '请输入回复内容', trigger: 'blur' }]">
              <el-input
                v-model="replyForm.content"
                type="textarea"
                :rows="4"
                placeholder="请输入回复内容"
              />
            </el-form-item>
            <!-- 附件上传 -->
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
                    支持上传图片、文档、压缩包等文件，单个文件不超过 10MB，最多 5 个
                  </div>
                </template>
              </el-upload>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSubmitReply" :loading="submitLoading">提交回复</el-button>
              <el-button @click="showReplyForm = false">取消</el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Document, Upload } from '@element-plus/icons-vue'
import type { FormInstance, UploadFile, UploadProps } from 'element-plus'
import request from '@/utils/http'
import { useUserStore } from '@/store/modules/user'

interface Ticket {
  id: number
  subject: string
  content: string
  client_username: string
  department_id: number
  department_name: string
  priority: string
  status: string
  reply_count: number
  created_at: string
  updated_at: string
  attachments?: Attachment[]
  replies?: TicketReply[]
}

interface TicketReply {
  id: number
  content: string
  username: string
  is_admin: boolean
  created_at: string
  attachments?: Attachment[]
}

interface Attachment {
  id: number
  file_name: string
  file_path: string
  file_size: number
  url: string
}

interface Department {
  id: number
  name: string
}

const priorityTypeMap: Record<string, any> = {
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

const statusTypeMap: Record<string, any> = {
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

const loading = ref(false)
const submitLoading = ref(false)
const showReplyForm = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('工单详情')
const replyFormRef = ref<FormInstance>()

const userStore = useUserStore()

// 上传请求头
const uploadHeaders = computed(() => ({
  Authorization: userStore.accessToken
}))

const searchForm = reactive({
  keyword: '',
  status: '',
  department_id: undefined as number | undefined,
  priority: ''
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<Ticket[]>([])
const departments = ref<Department[]>([])
const currentTicket = ref<Ticket | null>(null)

const replyForm = reactive({
  content: '',
  attachments: [] as UploadFile[]
})

// 格式化文件大小
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 上传前校验
const beforeUpload: UploadProps['beforeUpload'] = (file) => {
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB!')
    return false
  }
  return true
}

// 上传成功
const handleUploadSuccess = (response: any, file: UploadFile) => {
  if (response?.data) {
    // 将上传成功的信息保存到文件对象
    file.url = response.data.url
    ;(file as any).attachment_id = response.data.id
  } else {
    ElMessage.error(response?.msg || '上传失败')
  }
}

// 上传失败
const handleUploadError = () => {
  ElMessage.error('文件上传失败')
}

// 超出限制
const handleExceed = () => {
  ElMessage.warning('最多只能上传 5 个文件')
}

// 获取上传接口地址（根据当前工单ID动态生成）
const getUploadUrl = computed(() => {
  if (currentTicket.value) {
    return `/api/admin/tickets/${currentTicket.value.id}/attachments`
  }
  return ''
})

const fetchTickets = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/tickets',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取工单列表失败:', error)
    ElMessage.error('获取工单列表失败')
  } finally {
    loading.value = false
  }
}

const fetchDepartments = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/ticket-depts'
    })
    departments.value = data || []
  } catch (error) {
    console.error('获取部门列表失败:', error)
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchTickets()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  searchForm.department_id = undefined
  searchForm.priority = ''
  handleSearch()
}

const handleView = async (row: Ticket) => {
  dialogTitle.value = '工单详情'
  showReplyForm.value = false
  try {
    const data = await request.get({
      url: `/api/admin/tickets/${row.id}`
    })
    // 合并工单信息和附件
    currentTicket.value = {
      ...data.ticket,
      attachments: data.attachments || [],
      replies: data.replies || []
    }
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取工单详情失败')
  }
}

const handleReply = async (row: Ticket) => {
  dialogTitle.value = '回复工单'
  showReplyForm.value = true
  replyForm.content = ''
  replyForm.attachments = []
  try {
    const data = await request.get({
      url: `/api/admin/tickets/${row.id}`
    })
    // 合并工单信息和附件
    currentTicket.value = {
      ...data.ticket,
      attachments: data.attachments || [],
      replies: data.replies || []
    }
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取工单详情失败')
  }
}

const handleSubmitReply = async () => {
  if (!replyFormRef.value) return

  await replyFormRef.value.validate(async (valid) => {
    if (!valid || !currentTicket.value) return

    submitLoading.value = true
    try {
      // 收集附件ID
      const attachmentIds = replyForm.attachments
        .filter(file => file.attachment_id)
        .map(file => file.attachment_id)

      await request.post({
        url: `/api/admin/tickets/${currentTicket.value.id}/reply`,
        params: {
          content: replyForm.content,
          attachment_ids: attachmentIds
        }
      })
      ElMessage.success('回复成功')
      showReplyForm.value = false
      handleView(currentTicket.value)
      fetchTickets()
    } catch (error) {
      ElMessage.error('回复失败')
    } finally {
      submitLoading.value = false
    }
  })
}

const handleClose = async (row: Ticket) => {
  try {
    await request.post({
      url: `/api/admin/tickets/${row.id}/close`
    })
    ElMessage.success('工单已关闭')
    fetchTickets()
  } catch (error) {
    ElMessage.error('关闭失败')
  }
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchTickets()
}

const handlePageChange = () => {
  fetchTickets()
}

onMounted(() => {
  fetchTickets()
  fetchDepartments()
})
</script>

<style scoped lang="scss">
.ticket-list-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.ticket-detail {
  .ticket-info {
    h3 {
      margin: 0 0 12px;
      font-size: 18px;
    }

    .ticket-meta {
      display: flex;
      gap: 16px;
      margin-bottom: 16px;
      color: var(--el-text-color-secondary);
      font-size: 14px;
    }

    .ticket-content {
      padding: 16px;
      background: var(--el-fill-color-lighter);
      border-radius: 4px;
      line-height: 1.6;
      white-space: pre-wrap;
    }

    .ticket-attachments {
      margin-top: 12px;
      padding: 12px;
      background: var(--el-fill-color-lighter);
      border-radius: 4px;

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

  .reply-list {
    margin-top: 24px;

    .reply-item {
      padding: 16px;
      margin-bottom: 12px;
      border: 1px solid var(--el-border-color-lighter);
      border-radius: 4px;

      &.is-admin {
        background: var(--el-color-primary-light-9);
        border-color: var(--el-color-primary-light-7);
      }

      .reply-header {
        display: flex;
        justify-content: space-between;
        margin-bottom: 8px;

        .reply-user {
          font-weight: 500;
        }

        .reply-time {
          color: var(--el-text-color-secondary);
          font-size: 13px;
        }
      }

      .reply-content {
        line-height: 1.6;
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

        .attachment-item {
          .file-size {
            color: var(--el-text-color-secondary);
            font-size: 12px;
          }
        }
      }
    }
  }

  .reply-form {
    margin-top: 16px;
  }
}
</style>
