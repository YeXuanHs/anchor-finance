<template>
  <div class="email-page page-container">
    <div class="art-card">
      <h3>邮件配置</h3>
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
        v-loading="loading"
        style="max-width: 700px"
      >
        <el-form-item label="SMTP主机" prop="smtp_host">
          <el-input v-model="formData.smtp_host" placeholder="如 smtp.qq.com" />
        </el-form-item>
        <el-form-item label="SMTP端口" prop="smtp_port">
          <el-input-number v-model="formData.smtp_port" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="邮箱账号" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="formData.password" type="password" show-password placeholder="邮箱密码或授权码" />
        </el-form-item>
        <el-form-item label="发件人地址" prop="from_address">
          <el-input v-model="formData.from_address" placeholder="noreply@example.com" />
        </el-form-item>
        <el-form-item label="发件人名称" prop="from_name">
          <el-input v-model="formData.from_name" placeholder="如：智简魔方" />
        </el-form-item>
        <el-form-item label="加密方式" prop="encryption">
          <el-radio-group v-model="formData.encryption">
            <el-radio value="none">无</el-radio>
            <el-radio value="ssl">SSL</el-radio>
            <el-radio value="tls">TLS</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="启用邮件" prop="enabled">
          <el-switch v-model="formData.enabled" />
        </el-form-item>

        <el-divider />

        <el-form-item label="测试收件人" prop="test_email">
          <el-input v-model="testEmail" placeholder="输入邮箱地址发送测试邮件" style="width: 300px" />
          <el-button type="success" :loading="testLoading" @click="handleTest" style="margin-left: 12px">
            发送测试
          </el-button>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="saveLoading" @click="handleSave">保存配置</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const saveLoading = ref(false)
const testLoading = ref(false)
const formRef = ref<FormInstance>()
const testEmail = ref('')

const formData = reactive({
  smtp_host: '',
  smtp_port: 465,
  username: '',
  password: '',
  from_address: '',
  from_name: '',
  encryption: 'ssl' as string,
  enabled: false
})

const formRules: FormRules = {
  smtp_host: [{ required: true, message: '请输入SMTP主机', trigger: 'blur' }],
  smtp_port: [{ required: true, message: '请输入SMTP端口', trigger: 'blur' }],
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  from_address: [{ required: true, message: '请输入发件人地址', trigger: 'blur' }, { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }]
}

const fetchConfig = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/admin/settings/email')
    if (data) {
      Object.assign(formData, {
        smtp_host: data.smtp_host || '',
        smtp_port: data.smtp_port || 465,
        username: data.username || '',
        password: data.password || '',
        from_address: data.from_address || '',
        from_name: data.from_name || '',
        encryption: data.encryption || 'ssl',
        enabled: data.enabled ?? false
      })
    }
  } catch {} finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  saveLoading.value = true
  try {
    await request.put('/api/admin/settings/email', formData)
    ElMessage.success('保存成功')
  } catch {} finally {
    saveLoading.value = false
  }
}

const handleTest = async () => {
  if (!testEmail.value) {
    ElMessage.warning('请输入测试收件人邮箱')
    return
  }
  testLoading.value = true
  try {
    await request.post('/api/admin/settings/email/test', { to: testEmail.value })
    ElMessage.success('测试邮件已发送，请检查收件箱')
  } catch {} finally {
    testLoading.value = false
  }
}

onMounted(fetchConfig)
</script>

<style scoped lang="scss">
.email-page {
  .art-card {
    h3 {
      font-size: 18px;
      font-weight: 600;
      margin: 0 0 24px;
      padding-bottom: 16px;
      border-bottom: 1px solid var(--el-border-color-lighter);
    }
  }
}
</style>
