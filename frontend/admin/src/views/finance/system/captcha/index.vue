<template>
  <div class="captcha-settings-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>验证码设置</span>
          <el-button type="primary" :loading="saving" @click="handleSave">保存设置</el-button>
        </div>
      </template>

      <el-form :model="formData" label-width="150px">
        <el-divider content-position="left">验证码类型</el-divider>
        <el-form-item label="验证码类型">
          <el-radio-group v-model="formData.type">
            <el-radio value="none">关闭</el-radio>
            <el-radio value="image">图片验证码</el-radio>
            <el-radio value="slide">滑动验证</el-radio>
            <el-radio value="recaptcha">reCAPTCHA</el-radio>
            <el-radio value="turnstile">Turnstile</el-radio>
          </el-radio-group>
        </el-form-item>

        <template v-if="formData.type === 'recaptcha' || formData.type === 'turnstile'">
          <el-divider content-position="left">密钥配置</el-divider>
          <el-form-item label="Site Key">
            <el-input v-model="formData.site_key" placeholder="请输入Site Key" style="width: 400px" />
          </el-form-item>
          <el-form-item label="Secret Key">
            <el-input v-model="formData.secret_key" type="password" placeholder="请输入Secret Key" show-password style="width: 400px" />
          </el-form-item>
        </template>

        <el-divider content-position="left">使用场景</el-divider>
        <el-form-item label="登录页面">
          <el-switch v-model="formData.login_enabled" />
        </el-form-item>
        <el-form-item label="注册页面">
          <el-switch v-model="formData.register_enabled" />
        </el-form-item>
        <el-form-item label="找回密码">
          <el-switch v-model="formData.reset_password_enabled" />
        </el-form-item>
        <el-form-item label="联系我们">
          <el-switch v-model="formData.contact_enabled" />
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const saving = ref(false)
const formData = reactive({
  type: 'none',
  site_key: '',
  secret_key: '',
  login_enabled: false,
  register_enabled: false,
  reset_password_enabled: false,
  contact_enabled: false
})

const fetchSettings = async () => {
  try {
    const data = await request.get({ url: '/api/admin/settings/captcha' })
    Object.assign(formData, data)
  } catch (error) {
    console.error('获取验证码设置失败:', error)
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await request.post({ url: '/api/admin/settings/captcha', data: formData })
    ElMessage.success('保存成功')
  } catch (error) {
    console.error('保存失败:', error)
  } finally {
    saving.value = false
  }
}

onMounted(() => { fetchSettings() })
</script>

<style scoped lang="scss">
.captcha-settings-page { padding: 16px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
:deep(.el-divider__text) { font-weight: 600; color: #1D2129; }
</style>
