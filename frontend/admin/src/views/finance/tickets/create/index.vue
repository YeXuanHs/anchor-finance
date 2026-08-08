<template>
  <div class="create-ticket-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>新建工单</span>
          <el-button @click="$router.back()">返回</el-button>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px" size="default">
        <!-- 客户选择 -->
        <el-divider content-position="left">客户信息</el-divider>

        <el-form-item label="客户" prop="client_id">
          <el-select
            v-model="formData.client_id"
            filterable
            remote
            :remote-method="searchClients"
            :loading="clientSearching"
            placeholder="请输入客户名搜索"
            style="width: 400px"
          >
            <el-option
              v-for="client in clientOptions"
              :key="client.id"
              :label="`${client.username} (${client.email})`"
              :value="client.id"
            />
          </el-select>
        </el-form-item>

        <!-- 工单信息 -->
        <el-divider content-position="left">工单信息</el-divider>

        <el-form-item label="部门" prop="department_id">
          <el-select v-model="formData.department_id" placeholder="请选择部门" style="width: 400px">
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="优先级" prop="priority">
          <el-radio-group v-model="formData.priority">
            <el-radio :value="1">低</el-radio>
            <el-radio :value="2">普通</el-radio>
            <el-radio :value="3">高</el-radio>
            <el-radio :value="4">紧急</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="主题" prop="subject">
          <el-input v-model="formData.subject" placeholder="请输入工单主题" style="width: 600px" />
        </el-form-item>

        <el-form-item label="内容" prop="content">
          <el-input
            v-model="formData.content"
            type="textarea"
            :rows="8"
            placeholder="请输入工单内容"
            style="width: 600px"
          />
        </el-form-item>

        <el-form-item label="附件">
          <el-upload
            action="/api/admin/upload"
            :on-success="handleUploadSuccess"
            :on-remove="handleRemoveFile"
            :file-list="fileList"
            multiple
          >
            <el-button type="primary" plain>
              <el-icon><Upload /></el-icon>
              上传附件
            </el-button>
          </el-upload>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">提交工单</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import type { FormInstance, FormRules, UploadFile } from 'element-plus'
import request from '@/utils/http'

const router = useRouter()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const clientSearching = ref(false)
const clientOptions = ref<any[]>([])
const departments = ref<any[]>([])
const fileList = ref<UploadFile[]>([])

// 表单数据
const formData = reactive({
  client_id: null as number | null,
  department_id: null as number | null,
  priority: 2,
  subject: '',
  content: '',
  attachments: [] as string[]
})

// 表单验证规则
const rules: FormRules = {
  client_id: [{ required: true, message: '请选择客户', trigger: 'change' }],
  department_id: [{ required: true, message: '请选择部门', trigger: 'change' }],
  priority: [{ required: true, message: '请选择优先级', trigger: 'change' }],
  subject: [{ required: true, message: '请输入工单主题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入工单内容', trigger: 'blur' }]
}

// 搜索客户
const searchClients = async (query: string) => {
  if (!query) return
  clientSearching.value = true
  try {
    const data = await request.get({ url: '/api/admin/clients', params: { keyword: query, page_size: 20 } })
    clientOptions.value = data?.list || []
  } catch (error) {
    console.error('搜索客户失败:', error)
  } finally {
    clientSearching.value = false
  }
}

// 获取部门列表
const fetchDepartments = async () => {
  try {
    const data = await request.get({ url: '/api/admin/ticket-departments' })
    departments.value = data || []
  } catch (error) {
    console.error('获取部门列表失败:', error)
  }
}

// 上传成功
const handleUploadSuccess = (response: any, file: UploadFile) => {
  if (response?.url) {
    formData.attachments.push(response.url)
  }
}

// 移除文件
const handleRemoveFile = (file: UploadFile) => {
  const response = file.response as any
  const index = formData.attachments.findIndex((url) => url === response?.url)
  if (index > -1) {
    formData.attachments.splice(index, 1)
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    submitting.value = true
    await request.post({ url: '/api/admin/tickets', data: formData })
    ElMessage.success('工单创建成功')
    router.push('/support-ticket')
  } catch (error) {
    console.error('创建工单失败:', error)
  } finally {
    submitting.value = false
  }
}

// 重置表单
const handleReset = () => {
  formRef.value?.resetFields()
  fileList.value = []
  formData.attachments = []
}

onMounted(() => {
  fetchDepartments()
})
</script>

<style scoped lang="scss">
.create-ticket-page {
  padding: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

:deep(.el-divider__text) {
  font-weight: 600;
  color: #1D2129;
}
</style>
