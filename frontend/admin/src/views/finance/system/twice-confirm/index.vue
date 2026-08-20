<template>
  <div class="twice-confirm-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('page.twiceConfirm.title') }}</span>
        </div>
      </template>

      <el-form :model="configForm" ref="configFormRef" label-width="160px" class="config-form">
        <el-divider content-position="left">{{ $t('page.twiceConfirm.sensitiveOps') }}</el-divider>

        <el-form-item :label="$t('page.twiceConfirm.deleteClient')">
          <el-switch v-model="configForm.delete_client" />
          <span class="form-tip">{{ $t('page.twiceConfirm.deleteClientTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('page.twiceConfirm.deleteOrder')">
          <el-switch v-model="configForm.delete_order" />
          <span class="form-tip">{{ $t('page.twiceConfirm.deleteOrderTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('page.twiceConfirm.batchOperation')">
          <el-switch v-model="configForm.batch_operation" />
          <span class="form-tip">{{ $t('page.twiceConfirm.batchOperationTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('page.twiceConfirm.fundOperation')">
          <el-switch v-model="configForm.fund_operation" />
          <span class="form-tip">{{ $t('page.twiceConfirm.fundOperationTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('page.twiceConfirm.changePassword')">
          <el-switch v-model="configForm.change_password" />
          <span class="form-tip">{{ $t('page.twiceConfirm.changePasswordTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('page.twiceConfirm.systemSettings')">
          <el-switch v-model="configForm.system_settings" />
          <span class="form-tip">{{ $t('page.twiceConfirm.systemSettingsTip') }}</span>
        </el-form-item>

        <el-divider content-position="left">{{ $t('page.twiceConfirm.confirmMethod') }}</el-divider>

        <el-form-item :label="$t('page.twiceConfirm.confirmMethod')">
          <el-checkbox-group v-model="configForm.methods">
            <el-checkbox value="password">{{ $t('page.twiceConfirm.passwordConfirm') }}</el-checkbox>
            <el-checkbox value="email_code">{{ $t('page.twiceConfirm.emailCode') }}</el-checkbox>
            <el-checkbox value="sms_code">{{ $t('page.twiceConfirm.smsCode') }}</el-checkbox>
            <el-checkbox value="totp">{{ $t('page.twiceConfirm.totp') }}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item :label="$t('page.twiceConfirm.codeExpire')">
          <el-input-number v-model="configForm.code_expire" :min="30" :max="300" :step="30" />
          {{ $t('page.twiceConfirm.seconds') }}
        </el-form-item>

        <el-form-item :label="$t('page.twiceConfirm.maxRetries')">
          <el-input-number v-model="configForm.max_retries" :min="1" :max="10" />
          <span class="form-tip">{{ $t('page.twiceConfirm.maxRetriesTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('page.twiceConfirm.lockDuration')">
          <el-input-number v-model="configForm.lock_duration" :min="60" :max="3600" :step="60" />
          {{ $t('page.twiceConfirm.seconds') }}
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="saveLoading">{{ $t('page.twiceConfirm.saveConfig') }}</el-button>
          <el-button @click="handleReset">{{ $t('page.twiceConfirm.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const configFormRef = ref<FormInstance>()
const saveLoading = ref(false)

const configForm = reactive({
  delete_client: true,
  delete_order: true,
  batch_operation: true,
  fund_operation: true,
  change_password: true,
  system_settings: false,
  methods: ['password'] as string[],
  code_expire: 60,
  max_retries: 3,
  lock_duration: 300
})

const fetchConfig = async () => {
  try {
    const data = await request.get({ url: '/api/admin/config/second-verify' })
    if (data) {
      Object.assign(configForm, data)
    }
  } catch (error) {
    console.error('获取配置失败:', error)
  }
}

const handleSave = async () => {
  saveLoading.value = true
  try {
    await request.put({
      url: '/api/admin/config/second-verify',
      params: { ...configForm }
    })
    ElMessage.success($t('page.twiceConfirm.saveSuccess'))
  } catch (error) {
    ElMessage.error($t('page.twiceConfirm.saveFailed'))
  } finally {
    saveLoading.value = false
  }
}

const handleReset = () => {
  fetchConfig()
}

onMounted(() => { fetchConfig() })
</script>

<style scoped lang="scss">
.twice-confirm-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.config-form { max-width: 700px; }
.form-tip { margin-left: 12px; color: var(--el-text-color-secondary); font-size: 13px; }
</style>
