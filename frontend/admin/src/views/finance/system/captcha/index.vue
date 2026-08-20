<template>
  <div class="captcha-settings-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('captcha.title') }}</span>
          <el-button type="primary" :loading="saving" @click="handleSave">{{ $t('captcha.saveSettings') }}</el-button>
        </div>
      </template>

      <el-form :model="formData" label-width="150px">
        <el-divider content-position="left">{{ $t('captcha.captchaType') }}</el-divider>
        <el-form-item :label="$t('captcha.captchaType')">
          <el-radio-group v-model="formData.type">
            <el-radio value="none">{{ $t('captcha.close') }}</el-radio>
            <el-radio value="image">{{ $t('captcha.imageCaptcha') }}</el-radio>
            <el-radio value="slide">{{ $t('captcha.slideCaptcha') }}</el-radio>
            <el-radio value="recaptcha">reCAPTCHA</el-radio>
            <el-radio value="turnstile">Turnstile</el-radio>
          </el-radio-group>
        </el-form-item>

        <template v-if="formData.type === 'recaptcha' || formData.type === 'turnstile'">
          <el-divider content-position="left">{{ $t('captcha.secretConfig') }}</el-divider>
          <el-form-item label="Site Key">
            <el-input v-model="formData.site_key" :placeholder="$t('captcha.enterSiteKey')" style="width: 400px" />
          </el-form-item>
          <el-form-item label="Secret Key">
            <el-input v-model="formData.secret_key" type="password" :placeholder="$t('captcha.enterSecretKey')" show-password style="width: 400px" />
          </el-form-item>
        </template>

        <el-divider content-position="left">{{ $t('captcha.useScenarios') }}</el-divider>
        <el-form-item :label="$t('captcha.loginPage')">
          <el-switch v-model="formData.login_enabled" />
        </el-form-item>
        <el-form-item :label="$t('captcha.registerPage')">
          <el-switch v-model="formData.register_enabled" />
        </el-form-item>
        <el-form-item :label="$t('captcha.resetPassword')">
          <el-switch v-model="formData.reset_password_enabled" />
        </el-form-item>
        <el-form-item :label="$t('captcha.contactUs')">
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
import { $t } from '@/locales'

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
    const data = await request.get({ url: '/api/admin/captcha-config' })
    Object.assign(formData, data)
  } catch (error) {
    console.error('Failed to fetch captcha settings:', error)
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await request.put({ url: '/api/admin/captcha-config/basic', data: formData })
    ElMessage.success($t('captcha.saveSuccess'))
  } catch (error) {
    console.error('Save failed:', error)
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
