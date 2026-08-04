<template>
  <div class="ticket-create">
    <div class="page-header">
      <el-page-header @back="$router.push('/user/tickets')">
        <template #content>
          <span class="page-title">提交工单</span>
        </template>
      </el-page-header>
    </div>
    
    <el-card class="form-card">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="工单类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择工单类型">
            <el-option label="技术支持" value="technical" />
            <el-option label="财务问题" value="billing" />
            <el-option label="销售咨询" value="sales" />
            <el-option label="投诉建议" value="complaint" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="优先级" prop="priority">
          <el-radio-group v-model="form.priority">
            <el-radio value="low">低</el-radio>
            <el-radio value="medium">中</el-radio>
            <el-radio value="high">高</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="关联产品">
          <el-select v-model="form.product_id" placeholder="请选择关联产品（可选）" clearable>
            <el-option v-for="p in products" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="主题" prop="subject">
          <el-input v-model="form.subject" placeholder="请输入工单主题" />
        </el-form-item>
        
        <el-form-item label="内容" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="8" placeholder="请详细描述您的问题" />
        </el-form-item>
        
        <el-form-item label="附件">
          <el-upload
            action="/api/v1/tickets/upload"
            :on-success="handleUploadSuccess"
            :file-list="form.attachments"
          >
            <el-button>上传附件</el-button>
            <template #tip>
              <div class="el-upload__tip">支持 jpg/png/pdf/zip 格式，单个文件不超过 10MB</div>
            </template>
          </el-upload>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="submitTicket" :loading="submitting">提交工单</el-button>
          <el-button @click="$router.push('/user/tickets')">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const formRef = ref()
const submitting = ref(false)
const products = ref<any[]>([])

const form = ref({
  type: '',
  priority: 'medium',
  product_id: '',
  subject: '',
  content: '',
  attachments: [] as any[]
})

const rules = {
  type: [{ required: true, message: '请选择工单类型', trigger: 'change' }],
  priority: [{ required: true, message: '请选择优先级', trigger: 'change' }],
  subject: [{ required: true, message: '请输入工单主题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入工单内容', trigger: 'blur' }]
}

const handleUploadSuccess = (response: any) => {
  form.value.attachments.push({
    name: response.data.name,
    url: response.data.url
  })
}

const fetchProducts = async () => {
  try {
    const { data } = await request.get('/api/v2/user/products')
    if (data?.data) {
      products.value = data.data
    }
  } catch (error) {
    // 忽略
  }
}

const submitTicket = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true
    
    await request.post('/api/v2/tickets', form.value)
    ElMessage.success('工单已提交')
    router.push('/user/tickets')
  } catch (error) {
    if (error !== false) {
      ElMessage.error('提交失败')
    }
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchProducts()
})
</script>

<style scoped lang="scss">
.ticket-create {
  .page-header {
    margin-bottom: 24px;
  }
  
  .page-title {
    font-size: 18px;
    font-weight: 600;
  }
  
  .form-card {
    max-width: 800px;
  }
}
</style>
