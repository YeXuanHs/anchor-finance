<template>
  <div class="system-messages-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>系统消息管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            发送消息
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="标题/内容" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="系统通知" value="system" />
            <el-option label="活动通知" value="activity" />
            <el-option label="安全警告" value="security" />
            <el-option label="订单通知" value="order" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="未读" value="unread" />
            <el-option label="已读" value="read" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
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
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="typeTagMap[row.type]" size="small">{{ typeTextMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="target" label="接收对象" width="120" />
        <el-table-column prop="is_read" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_read ? 'info' : 'danger'" size="small">
              {{ row.is_read ? '已读' : '未读' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="发送时间" width="180" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleDetail(row)">详情</el-button>
            <el-button v-if="!row.is_read" type="success" link @click="handleMarkRead(row)">标记已读</el-button>
            <el-popconfirm title="确定删除该消息吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
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

    <!-- 发送消息对话框 -->
    <el-dialog v-model="dialogVisible" title="发送系统消息" width="700px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入消息标题" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型" style="width: 100%">
            <el-option label="系统通知" value="system" />
            <el-option label="活动通知" value="activity" />
            <el-option label="安全警告" value="security" />
            <el-option label="订单通知" value="order" />
          </el-select>
        </el-form-item>
        <el-form-item label="接收对象">
          <el-select v-model="formData.target" placeholder="全部用户" clearable style="width: 100%">
            <el-option label="全部用户" value="all" />
            <el-option label="指定用户" value="specific" />
            <el-option label="指定分组" value="group" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="formData.target === 'specific'" label="用户ID">
          <el-input v-model="formData.user_ids" placeholder="多个用户ID用逗号分隔" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="8" placeholder="请输入消息内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">发送</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="消息详情" width="600px">
      <el-descriptions :column="1" border v-if="detailData">
        <el-descriptions-item label="ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item label="标题">{{ detailData.title }}</el-descriptions-item>
        <el-descriptions-item label="类型">
          <el-tag :type="typeTagMap[detailData.type]" size="small">{{ typeTextMap[detailData.type] || detailData.type }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="接收对象">{{ detailData.target }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="detailData.is_read ? 'info' : 'danger'" size="small">
            {{ detailData.is_read ? '已读' : '未读' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="内容">
          <div style="white-space: pre-wrap;">{{ detailData.content }}</div>
        </el-descriptions-item>
        <el-descriptions-item label="发送时间">{{ detailData.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
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

defineOptions({ name: 'SystemMessagesManage' })

const typeTextMap: Record<string, string> = { system: '系统通知', activity: '活动通知', security: '安全警告', order: '订单通知' }
const typeTagMap: Record<string, any> = { system: 'primary', activity: 'success', security: 'danger', order: 'warning' }

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const detailVisible = ref(false)
const detailData = ref<any>(null)
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: '',
  type: '',
  status: '',
  date_range: [] as string[]
})

const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const formData = reactive({
  title: '',
  type: 'system',
  target: 'all',
  user_ids: '',
  content: ''
})

const formRules: FormRules = {
  title: [{ required: true, message: '请输入消息标题', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  content: [{ required: true, message: '请输入消息内容', trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      keyword: searchForm.keyword || undefined,
      type: searchForm.type || undefined,
      status: searchForm.status || undefined
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({ url: '/api/admin/system-messages', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取系统消息失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { Object.assign(searchForm, { keyword: '', type: '', status: '', date_range: [] }); handleSearch() }

const handleAdd = () => {
  Object.assign(formData, { title: '', type: 'system', target: 'all', user_ids: '', content: '' })
  dialogVisible.value = true
}

const handleDetail = (row: any) => { detailData.value = row; detailVisible.value = true }

const handleMarkRead = async (row: any) => {
  try {
    await request.put({ url: `/api/admin/system-messages/${row.id}/read` })
    ElMessage.success('标记成功')
    fetchData()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/system-messages/${row.id}` })
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      await request.post({ url: '/api/admin/system-messages', params: { ...formData } })
      ElMessage.success('发送成功')
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error('发送失败')
    } finally {
      submitLoading.value = false
    }
  })
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.system-messages-page {
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
</style>
