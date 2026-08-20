<template>
  <div class="security-settings-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('security.title') }}</span>
          <el-button type="primary" :loading="saving" @click="handleSave">{{ $t('security.saveSettings') }}</el-button>
        </div>
      </template>

      <el-form :model="formData" label-width="150px">
        <el-divider content-position="left">{{ $t('security.passwordPolicy') }}</el-divider>
        <el-form-item :label="$t('security.minPasswordLength')">
          <el-input-number v-model="formData.min_password_length" :min="6" :max="32" />
        </el-form-item>
        <el-form-item :label="$t('security.passwordComplexity')">
          <el-checkbox-group v-model="formData.password_complexity">
            <el-checkbox value="uppercase">{{ $t('security.uppercase') }}</el-checkbox>
            <el-checkbox value="lowercase">{{ $t('security.lowercase') }}</el-checkbox>
            <el-checkbox value="number">{{ $t('security.number') }}</el-checkbox>
            <el-checkbox value="special">{{ $t('security.special') }}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item :label="$t('security.passwordExpiry')">
          <el-input-number v-model="formData.password_expiry_days" :min="0" :max="365" />
          <span class="form-tip">{{ $t('security.daysNeverExpire') }}</span>
        </el-form-item>

        <el-divider content-position="left">{{ $t('security.loginSecurity') }}</el-divider>
        <el-form-item :label="$t('security.maxLoginAttempts')">
          <el-input-number v-model="formData.max_login_attempts" :min="3" :max="20" />
          <span class="form-tip">{{ $t('security.times') }}</span>
        </el-form-item>
        <el-form-item :label="$t('security.lockoutDuration')">
          <el-input-number v-model="formData.lockout_duration" :min="5" :max="1440" />
          <span class="form-tip">{{ $t('security.minutes') }}</span>
        </el-form-item>
        <el-form-item :label="$t('security.sessionTimeout')">
          <el-input-number v-model="formData.session_timeout" :min="5" :max="1440" />
          <span class="form-tip">{{ $t('security.minutes') }}</span>
        </el-form-item>
        <el-form-item :label="$t('security.singleSession')">
          <el-switch v-model="formData.single_session" />
          <span class="form-tip">{{ $t('security.singleSessionTip') }}</span>
        </el-form-item>

        <el-divider content-position="left">{{ $t('security.ipSecurity') }}</el-divider>
        <el-form-item :label="$t('security.ipWhitelist')">
          <el-input v-model="formData.ip_whitelist" type="textarea" :rows="3" :placeholder="$t('security.ipWhitelistPlaceholder')" style="width: 400px" />
        </el-form-item>
        <el-form-item :label="$t('security.ipBlacklist')">
          <el-input v-model="formData.ip_blacklist" type="textarea" :rows="3" :placeholder="$t('security.ipBlacklistPlaceholder')" style="width: 400px" />
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
  min_password_length: 8,
  password_complexity: ['lowercase', 'number'],
  password_expiry_days: 0,
  max_login_attempts: 5,
  lockout_duration: 30,
  session_timeout: 120,
  single_session: false,
  ip_whitelist: '',
  ip_blacklist: ''
})

const fetchSettings = async () => {
  try {
    const data = await request.get({ url: '/api/admin/config/security' })
    Object.assign(formData, data)
  } catch (error) {
    console.error('Failed to fetch security settings:', error)
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await request.put({ url: '/api/admin/config/security', data: formData })
    ElMessage.success($t('security.saveSuccess'))
  } catch (error) {
    console.error('Save failed:', error)
  } finally {
    saving.value = false
  }
}

onMounted(() => { fetchSettings() })
</script>

<style scoped lang="scss">
.security-settings-page { padding: 16px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.form-tip { margin-left: 10px; font-size: 12px; color: #86909C; }
:deep(.el-divider__text) { font-weight: 600; color: #1D2129; }
</style>
