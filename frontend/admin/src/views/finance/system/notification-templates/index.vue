<template>
  <div class="notification-templates-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>通知模板管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加模板
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="模板名称/事件" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="邮件" value="email" />
            <el-option label="短信" value="sms" />
            <el-option label="站内信" value="internal" />
            <el-option label="微信" value="wechat" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="模板名称" width="180" />
        <el-table-column prop="event" label="触发事件" width="200" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="typeTagMap[row.type]" size="small">{{ typeTextMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="标题/主题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="warning" link @click="handleResetEvent(row)">重置事件</el-button>
            <el-popconfirm title="确定删除该模板吗？" @confirm="handleDelete(row)">
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

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="800px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="模板名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入模板名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="类型" prop="type">
              <el-select v-model="formData.type" placeholder="请选择类型" style="width: 100%">
                <el-option label="邮件" value="email" />
                <el-option label="短信" value="sms" />
                <el-option label="站内信" value="internal" />
                <el-option label="微信" value="wechat" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="触发事件" prop="event">
          <el-input v-model="formData.event" placeholder="如: order.created, user.registered" />
        </el-form-item>
        <el-form-item label="标题/主题" prop="subject">
          <el-input v-model="formData.subject" placeholder="请输入标题或主题" />
        </el-form-item>
        <el-form-item label="模板内容" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="10" placeholder="支持变量: {username}, {order_id}, {amount} 等" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="状态">
              <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序">
              <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'NotificationTemplatesManage' })

const typeTextMap: Record<string, string> = { email: '邮件', sms: '短信', internal: '站内信', wechat: '微信' }
const typeTagMap: Record<string, string> = { email: 'primary', sms: 'success', internal: 'info', wechat: 'warning' }

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('添加模板')
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: '',
  type: '',
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  type: 'email',
  event: '',
  subject: '',
  content: '',
  status: 1,
  sort: 0
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  event: [{ required: true, message: '请输入触发事件', trigger: 'blur' }],
  content: [{ required: true, message: '请输入模板内容', trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/notifications/templates',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取通知模板列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { Object.assign(searchForm, { keyword: '', type: '', status: undefined }); handleSearch() }

const resetForm = () => {
  formData.id = undefined
  formData.name = ''
  formData.type = 'email'
  formData.event = ''
  formData.subject = ''
  formData.content = ''
  formData.status = 1
  formData.sort = 0
}

const handleAdd = () => {
  dialogTitle.value = '添加模板'
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑模板'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleResetEvent = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定重置模板「${row.name}」的事件绑定吗？`, '重置确认')
    await request.put({ url: `/api/admin/notifications/templates/${row.id}/reset-event` })
    ElMessage.success('事件重置成功')
    fetchData()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error('重置失败')
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/notifications/templates/${row.id}` })
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
      if (formData.id) {
        await request.put({ url: `/api/admin/notifications/templates/${formData.id}`, params: { ...formData } })
      } else {
        await request.post({ url: '/api/admin/notifications/templates', params: { ...formData } })
      }
      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error('操作失败')
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
.notification-templates-page {
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
