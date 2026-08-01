<template>
  <div class="email-view-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>客户邮件查看</span>
          <el-button type="primary" @click="handleSendEmail">
            <el-icon><Promotion /></el-icon>
            发送邮件
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="方向">
          <el-select v-model="searchForm.direction" placeholder="全部" clearable>
            <el-option label="收到" value="in" />
            <el-option label="发出" value="out" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="已读" :value="1" />
            <el-option label="未读" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="主题/内容" clearable />
        </el-form-item>
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border @row-click="handleRowClick">
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="direction" label="方向" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.direction === 'in' ? 'success' : 'primary'" size="small">
              {{ row.direction === 'in' ? '收到' : '发出' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="主题" min-width="250" show-overflow-tooltip>
          <template #default="{ row }">
            <span :class="{ 'unread-text': row.status === 0 }">{{ row.subject }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="from" label="发件人" width="180" show-overflow-tooltip />
        <el-table-column prop="to" label="收件人" width="180" show-overflow-tooltip />
        <el-table-column prop="has_attachment" label="附件" width="70" align="center">
          <template #default="{ row }">
            <el-icon v-if="row.has_attachment" color="#409eff"><Paperclip /></el-icon>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'info' : 'danger'" size="small">
              {{ row.status === 1 ? '已读' : '未读' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="170" />
        <el-table-column label="操作" width="120" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click.stop="handleView(row)">查看</el-button>
            <el-popconfirm title="确定删除该邮件吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link @click.stop>删除</el-button>
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

    <!-- 邮件详情对话框 -->
    <el-dialog v-model="detailVisible" title="邮件详情" width="750px" destroy-on-close>
      <div class="email-detail">
        <div class="email-header">
          <h3 class="email-subject">{{ detailData.subject }}</h3>
          <div class="email-meta">
            <div class="meta-item">
              <span class="label">发件人:</span>
              <span>{{ detailData.from }}</span>
            </div>
            <div class="meta-item">
              <span class="label">收件人:</span>
              <span>{{ detailData.to }}</span>
            </div>
            <div class="meta-item">
              <span class="label">时间:</span>
              <span>{{ detailData.created_at }}</span>
            </div>
          </div>
        </div>
        <el-divider />
        <div class="email-body" v-html="detailData.content"></div>
        <div class="email-attachments" v-if="detailData.attachments?.length">
          <el-divider content-position="left">附件</el-divider>
          <div class="attachment-list">
            <div v-for="att in detailData.attachments" :key="att.id" class="attachment-item">
              <el-icon><Document /></el-icon>
              <span>{{ att.name }}</span>
              <el-button type="primary" link @click="handleDownloadAttachment(att)">下载</el-button>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="handleReply">回复</el-button>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 发送邮件对话框 -->
    <el-dialog v-model="sendVisible" title="发送邮件" width="700px">
      <el-form :model="sendForm" :rules="sendRules" ref="sendFormRef" label-width="80px">
        <el-form-item label="收件人" prop="to">
          <el-input v-model="sendForm.to" placeholder="请输入收件人邮箱" />
        </el-form-item>
        <el-form-item label="主题" prop="subject">
          <el-input v-model="sendForm.subject" placeholder="请输入邮件主题" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="sendForm.content" type="textarea" :rows="10" placeholder="请输入邮件内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="sendVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSendSubmit" :loading="sendLoading">发送</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Promotion, Paperclip, Document } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const route = useRoute()
const clientId = route.params.id as string

const loading = ref(false)
const sendLoading = ref(false)

const searchForm = reactive({
  direction: '',
  status: undefined as number | undefined,
  keyword: '',
  date_range: [] as string[]
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])
const detailVisible = ref(false)
const detailData = ref<any>({})
const sendVisible = ref(false)
const sendFormRef = ref<FormInstance>()

const sendForm = reactive({
  to: '',
  subject: '',
  content: ''
})

const sendRules: FormRules = {
  to: [{ required: true, message: '请输入收件人', trigger: 'blur' }],
  subject: [{ required: true, message: '请输入主题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size, client_id: clientId }
    if (searchForm.direction) params.direction = searchForm.direction
    if (searchForm.status !== undefined) params.status = searchForm.status
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({ url: '/api/admin/emails', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取数据失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => {
  searchForm.direction = ''
  searchForm.status = undefined
  searchForm.keyword = ''
  searchForm.date_range = []
  handleSearch()
}

const handleRowClick = (row: any) => {
  if (row.status === 0) {
    row.status = 1
  }
}

const handleView = async (row: any) => {
  try {
    const data = await request.get({ url: `/api/admin/emails/${row.id}` })
    detailData.value = data
    if (row.status === 0) {
      row.status = 1
    }
    detailVisible.value = true
  } catch (error) {
    console.error('获取邮件详情失败:', error)
    ElMessage.error('获取邮件详情失败')
  }
}

const handleReply = () => {
  sendForm.to = detailData.value.from
  sendForm.subject = `Re: ${detailData.value.subject}`
  sendForm.content = `\n\n-------- 原始邮件 --------\n${detailData.value.content}`
  detailVisible.value = false
  sendVisible.value = true
}

const handleSendEmail = () => {
  sendForm.to = ''
  sendForm.subject = ''
  sendForm.content = ''
  sendVisible.value = true
}

const handleSendSubmit = async () => {
  if (!sendFormRef.value) return
  await sendFormRef.value.validate(async (valid) => {
    if (!valid) return
    sendLoading.value = true
    try {
      await request.post({ url: `/api/admin/clients/${clientId}/emails`, params: sendForm })
      ElMessage.success('发送成功')
      sendVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error('发送失败')
    } finally {
      sendLoading.value = false
    }
  })
}

const handleDownloadAttachment = (att: any) => {
  window.open(`/api/admin/email-attachments/${att.id}/download`, '_blank')
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/emails/${row.id}` })
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.email-view-page {
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

.unread-text {
  font-weight: 600;
  color: #303133;
}

.email-detail {
  .email-header {
    .email-subject {
      font-size: 18px;
      font-weight: 600;
      color: #303133;
      margin: 0 0 16px;
    }

    .email-meta {
      .meta-item {
        margin-bottom: 8px;
        font-size: 14px;
        color: #606266;

        .label {
          color: #909399;
          margin-right: 8px;
        }
      }
    }
  }

  .email-body {
    padding: 16px 0;
    line-height: 1.8;
    white-space: pre-wrap;
    word-break: break-all;
  }

  .email-attachments {
    .attachment-list {
      display: flex;
      flex-direction: column;
      gap: 8px;

      .attachment-item {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 8px 12px;
        background: #f5f7fa;
        border-radius: 4px;
      }
    }
  }
}
</style>
