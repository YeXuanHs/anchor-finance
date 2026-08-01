<template>
  <div class="notification-log-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>客户通知日志</span>
          <el-button type="primary" @click="handleSendNotification">
            <el-icon><Promotion /></el-icon>
            发送通知
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="邮件" value="email" />
            <el-option label="短信" value="sms" />
            <el-option label="站内信" value="site" />
            <el-option label="微信" value="wechat" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="发送成功" :value="1" />
            <el-option label="发送失败" :value="0" />
          </el-select>
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
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="type" label="通知类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="180" show-overflow-tooltip />
        <el-table-column prop="content" label="内容" min-width="250" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="error_msg" label="错误信息" width="150" show-overflow-tooltip />
        <el-table-column prop="operator" label="操作人" width="100" />
        <el-table-column prop="created_at" label="发送时间" width="170" />
        <el-table-column label="操作" width="100" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
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

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="通知详情" width="650px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="通知类型">{{ getTypeText(detailData.type) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="detailData.status === 1 ? 'success' : 'danger'" size="small">
            {{ detailData.status === 1 ? '成功' : '失败' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="标题" :span="2">{{ detailData.title }}</el-descriptions-item>
        <el-descriptions-item label="内容" :span="2">
          <div class="content-text">{{ detailData.content }}</div>
        </el-descriptions-item>
        <el-descriptions-item label="接收地址" :span="2">{{ detailData.receiver || '-' }}</el-descriptions-item>
        <el-descriptions-item label="操作人">{{ detailData.operator }}</el-descriptions-item>
        <el-descriptions-item label="发送时间">{{ detailData.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 发送通知对话框 -->
    <el-dialog v-model="sendVisible" title="发送通知" width="600px">
      <el-form :model="sendForm" :rules="sendRules" ref="sendFormRef" label-width="100px">
        <el-form-item label="通知类型" prop="type">
          <el-select v-model="sendForm.type" placeholder="请选择">
            <el-option label="邮件" value="email" />
            <el-option label="短信" value="sms" />
            <el-option label="站内信" value="site" />
          </el-select>
        </el-form-item>
        <el-form-item label="标题" prop="title">
          <el-input v-model="sendForm.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="sendForm.content" type="textarea" :rows="5" placeholder="请输入通知内容" />
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
import { Promotion } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const route = useRoute()
const clientId = route.params.id as string

const loading = ref(false)
const sendLoading = ref(false)

const searchForm = reactive({
  type: '',
  status: undefined as number | undefined,
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
  type: 'email',
  title: '',
  content: ''
})

const sendRules: FormRules = {
  type: [{ required: true, message: '请选择通知类型', trigger: 'change' }],
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = { email: '邮件', sms: '短信', site: '站内信', wechat: '微信' }
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
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({ url: '/api/admin/notification-logs', params })
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
  searchForm.type = ''
  searchForm.status = undefined
  searchForm.date_range = []
  handleSearch()
}

const handleViewDetail = (row: any) => {
  detailData.value = { ...row }
  detailVisible.value = true
}

const handleSendNotification = () => {
  sendForm.type = 'email'
  sendForm.title = ''
  sendForm.content = ''
  sendVisible.value = true
}

const handleSendSubmit = async () => {
  if (!sendFormRef.value) return
  await sendFormRef.value.validate(async (valid) => {
    if (!valid) return
    sendLoading.value = true
    try {
      await request.post({ url: `/api/admin/clients/${clientId}/notifications`, params: sendForm })
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

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.notification-log-page {
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

.content-text {
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
