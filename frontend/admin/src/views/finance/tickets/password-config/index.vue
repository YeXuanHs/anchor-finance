<template>
  <div class="ticket-password-config-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('passwordConfig.title') }}</span>
        </div>
      </template>

      <!-- 配置表单 -->
      <el-form :model="configForm" :rules="configRules" ref="configFormRef" label-width="120px" class="config-form">
        <el-form-item :label="$t('passwordConfig.enablePassword')" prop="enable_password">
          <el-switch v-model="configForm.enable_password" :active-text="$t('common.enable')" :inactive-text="$t('common.disable')" />
          <div class="form-tip">{{ $t('passwordConfig.enablePasswordTip') }}</div>
        </el-form-item>

        <el-form-item :label="$t('passwordConfig.passwordRule')" prop="password_rule" v-if="configForm.enable_password">
          <el-radio-group v-model="configForm.password_rule">
            <el-radio value="fixed">{{ $t('passwordConfig.ruleFixed') }}</el-radio>
            <el-radio value="random">{{ $t('passwordConfig.ruleRandom') }}</el-radio>
            <el-radio value="custom">{{ $t('passwordConfig.ruleCustom') }}</el-radio>
          </el-radio-group>
          <div class="form-tip">
            <template v-if="configForm.password_rule === 'fixed'">{{ $t('passwordConfig.ruleFixedTip') }}</template>
            <template v-else-if="configForm.password_rule === 'random'">{{ $t('passwordConfig.ruleRandomTip') }}</template>
            <template v-else>{{ $t('passwordConfig.ruleCustomTip') }}</template>
          </div>
        </el-form-item>

        <el-form-item :label="$t('passwordConfig.fixedPassword')" prop="fixed_password" v-if="configForm.enable_password && configForm.password_rule === 'fixed'">
          <el-input v-model="configForm.fixed_password" :placeholder="$t('passwordConfig.enterFixedPassword')" show-password />
        </el-form-item>

        <el-form-item :label="$t('passwordConfig.passwordLength')" prop="password_length" v-if="configForm.enable_password && configForm.password_rule === 'random'">
          <el-input-number v-model="configForm.password_length" :min="4" :max="20" />
          <div class="form-tip">{{ $t('passwordConfig.randomLengthTip') }}</div>
        </el-form-item>

        <el-form-item :label="$t('passwordConfig.passwordChars')" prop="password_chars" v-if="configForm.enable_password && configForm.password_rule === 'custom'">
          <el-input v-model="configForm.password_chars" :placeholder="$t('passwordConfig.enterPasswordChars')" />
          <div class="form-tip">{{ $t('passwordConfig.passwordCharsTip') }}</div>
        </el-form-item>

        <el-form-item :label="$t('passwordConfig.passwordLength')" prop="password_length" v-if="configForm.enable_password && configForm.password_rule === 'custom'">
          <el-input-number v-model="configForm.password_length" :min="4" :max="20" />
          <div class="form-tip">{{ $t('passwordConfig.customLengthTip') }}</div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="saveLoading">{{ $t('passwordConfig.saveConfig') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

interface PasswordConfig {
  enable_password: boolean
  password_rule: 'fixed' | 'random' | 'custom'
  fixed_password: string
  password_length: number
  password_chars: string
}

const configFormRef = ref<FormInstance>()
const saveLoading = ref(false)

const configForm = reactive<PasswordConfig>({
  enable_password: false,
  password_rule: 'fixed',
  fixed_password: '',
  password_length: 6,
  password_chars: '0123456789'
})

const configRules: FormRules = {
  fixed_password: [
    { required: true, message: () => $t('passwordConfig.enterFixedPassword'), trigger: 'blur' }
  ],
  password_length: [
    { required: true, message: () => $t('passwordConfig.enterPasswordLength'), trigger: 'blur' }
  ],
  password_chars: [
    { required: true, message: () => $t('passwordConfig.enterPasswordChars'), trigger: 'blur' }
  ]
}

const fetchConfig = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/ticket-prereply/password-config'
    })
    if (data) {
      Object.assign(configForm, data)
    }
  } catch (error) {
    console.error($t('passwordConfig.fetchFailed'), error)
  }
}

const handleSave = async () => {
  if (!configFormRef.value) return

  await configFormRef.value.validate(async (valid) => {
    if (!valid) return

    saveLoading.value = true
    try {
      await request.post({
        url: '/api/admin/ticket-prereply/password-config',
        params: { ...configForm }
      })
      ElMessage.success($t('common.saveSuccess'))
    } catch (error) {
      ElMessage.error($t('common.saveFailed'))
    } finally {
      saveLoading.value = false
    }
  })
}

const handleReset = () => {
  configForm.enable_password = false
  configForm.password_rule = 'fixed'
  configForm.fixed_password = ''
  configForm.password_length = 6
  configForm.password_chars = '0123456789'
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped lang="scss">
.ticket-password-config-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.config-form {
  max-width: 600px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}
</style>
