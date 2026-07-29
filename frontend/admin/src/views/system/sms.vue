<template>
  <div class="sms-page page-container">
    <div class="art-card">
      <div class="table-header">
        <h3>短信配置</h3>
      </div>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="140px"
        style="max-width: 680px"
        v-loading="loading"
      >
        <el-form-item label="短信服务商" prop="provider">
          <el-select v-model="formData.provider" placeholder="请选择短信服务商" style="width: 100%">
            <el-option label="阿里云短信" value="aliyun" />
            <el-option label="腾讯云短信" value="tencent" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>

        <el-form-item label="API Key" prop="api_key">
          <el-input v-model="formData.api_key" placeholder="请输入API Key" clearable />
        </el-form-item>

        <el-form-item label="API Secret" prop="api_secret">
          <el-input
            v-model="formData.api_secret"
            type="password"
            placeholder="请输入API Secret"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item label="短信签名" prop="sign_name">
          <el-input v-model="formData.sign_name" placeholder="请输入短信签名" clearable />
        </el-form-item>

        <el-form-item label="模板ID" prop="template_id">
          <el-input v-model="formData.template_id" placeholder="请输入模板ID" clearable />
        </el-form-item>

        <el-form-item label="启用短信服务">
          <el-switch v-model="formData.enabled" />
          <span class="form-tip">{{ formData.enabled ? '已启用' : '已禁用' }}</span>
        </el-form-item>

        <el-divider />

        <el-form-item label="测试手机号">
          <el-input v-model="testPhone" placeholder="输入测试手机号" clearable style="width: 280px" />
          <el-button
            type="warning"
            :loading="testLoading"
            :disabled="!testPhone"
            style="margin-left: 12px"
            @click="handleTestSend"
          >
            发送测试短信
          </el-button>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="saving" @click="handleSave">保存配置</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import request from '@/utils/request'

const formRef = ref<FormInstance>()
const loading = ref(false)
const saving = ref(false)
const testLoading = ref(false)
const testPhone = ref('')

const formData = reactive({
  provider: 'aliyun',
  api_key: '',
  api_secret: '',
  sign_name: '',
  template_id: '',
  enabled: true
})

const rules: FormRules = {
  provider: [{ required: true, message: '请选择短信服务商', trigger: 'change' }],
  api_key: [{ required: true, message: '请输入API Key', trigger: 'blur' }],
  api_secret: [{ required: true, message: '请输入API Secret', trigger: 'blur' }],
  sign_name: [{ required: true, message: '请输入短信签名', trigger: 'blur' }],
  template_id: [{ required: true, message: '请输入模板ID', trigger: 'blur' }]
}

const fetchConfig = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/admin/system/sms')
    if (data?.data) {
      Object.assign(formData, data.data)
    }
  } catch {
    // 使用默认值
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    await request.post('/api/admin/system/sms', formData)
    ElMessage.success('短信配置保存成功')
  } catch {
    ElMessage.error('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

const handleTestSend = async () => {
  if (!testPhone.value) return

  testLoading.value = true
  try {
    await request.post('/api/admin/system/sms/test', { phone: testPhone.value })
    ElMessage.success('测试短信已发送，请注意查收')
  } catch {
    ElMessage.error('测试发送失败，请检查配置')
  } finally {
    testLoading.value = false
  }
}

const handleReset = () => {
  formRef.value?.resetFields()
  fetchConfig()
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped lang="scss">
.sms-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
    }
  }

  .form-tip {
    margin-left: 12px;
    font-size: 13px;
    color: var(--el-text-color-secondary);
  }
}
</style>
