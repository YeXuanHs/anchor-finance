<template>
  <div class="resource-pool-settings-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsResourcePool.settings') }}</span>
          <el-button type="primary" @click="handleSave" :loading="submitLoading">
            <el-icon><Check /></el-icon>
            {{ $t('clientsResourcePool.saveConfig') }}
          </el-button>
        </div>
      </template>

      <el-form :model="settingsForm" :rules="settingsRules" ref="settingsFormRef" label-width="140px" v-loading="loading">
        <el-divider content-position="left">{{ $t('clientsResourcePool.basicConfig') }}</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('clientsResourcePool.poolName')" prop="pool_name">
              <el-input v-model="settingsForm.pool_name" :placeholder="$t('clientsResourcePool.enterPoolName')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('clientsResourcePool.poolCapacity')" prop="capacity">
              <el-input-number v-model="settingsForm.capacity" :min="1" :max="10000" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('clientsResourcePool.poolStatus')" prop="status">
              <el-switch v-model="settingsForm.status" :active-value="1" :inactive-value="0" :active-text="$t('common.enabled')" :inactive-text="$t('common.disabled')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('clientsResourcePool.poolType')" prop="pool_type">
              <el-select v-model="settingsForm.pool_type" :placeholder="$t('clientsResourcePool.selectPoolType')" style="width: 100%">
                <el-option :label="$t('clientsResourcePool.sharedType')" value="shared" />
                <el-option :label="$t('clientsResourcePool.exclusiveType')" value="exclusive" />
                <el-option :label="$t('clientsResourcePool.hybridType')" value="hybrid" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('clientsResourcePool.poolDescription')" prop="description">
          <el-input v-model="settingsForm.description" type="textarea" :rows="3" :placeholder="$t('clientsResourcePool.enterDescription')" />
        </el-form-item>

        <el-divider content-position="left">{{ $t('clientsResourcePool.allocationRules') }}</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('clientsResourcePool.allocationStrategy')" prop="allocation_strategy">
              <el-select v-model="settingsForm.allocation_strategy" :placeholder="$t('clientsResourcePool.selectStrategy')" style="width: 100%">
                <el-option :label="$t('clientsResourcePool.roundRobin')" value="round_robin" />
                <el-option :label="$t('clientsResourcePool.weighted')" value="weighted" />
                <el-option :label="$t('clientsResourcePool.priority')" value="priority" />
                <el-option :label="$t('clientsResourcePool.random')" value="random" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('clientsResourcePool.maxAllocation')" prop="max_allocation">
              <el-input-number v-model="settingsForm.max_allocation" :min="1" :max="1000" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('clientsResourcePool.allocationThreshold')" prop="allocation_threshold">
              <el-input-number v-model="settingsForm.allocation_threshold" :min="0" :max="100" style="width: 100%" />
              <div class="form-tip">{{ $t('clientsResourcePool.allocationThresholdTip') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('clientsResourcePool.cooldownSeconds')" prop="cooldown_seconds">
              <el-input-number v-model="settingsForm.cooldown_seconds" :min="0" :max="3600" style="width: 100%" />
              <div class="form-tip">{{ $t('clientsResourcePool.cooldownTip') }}</div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">{{ $t('clientsResourcePool.autoConfig') }}</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('clientsResourcePool.autoAssign')" prop="auto_assign">
              <el-switch v-model="settingsForm.auto_assign" :active-text="$t('clientsResourcePool.enable')" :inactive-text="$t('clientsResourcePool.disable')" />
              <div class="form-tip">{{ $t('clientsResourcePool.autoAssignTip') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('clientsResourcePool.autoReclaim')" prop="auto_reclaim">
              <el-switch v-model="settingsForm.auto_reclaim" :active-text="$t('clientsResourcePool.enable')" :inactive-text="$t('clientsResourcePool.disable')" />
              <div class="form-tip">{{ $t('clientsResourcePool.autoReclaimTip') }}</div>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('clientsResourcePool.idleTimeout')" prop="idle_timeout">
              <el-input-number v-model="settingsForm.idle_timeout" :min="1" :max="1440" style="width: 100%" />
              <div class="form-tip">{{ $t('clientsResourcePool.idleTimeoutTip') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('clientsResourcePool.notifyThreshold')" prop="notify_threshold">
              <el-input-number v-model="settingsForm.notify_threshold" :min="0" :max="100" style="width: 100%" />
              <div class="form-tip">{{ $t('clientsResourcePool.notifyThresholdTip') }}</div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">{{ $t('clientsResourcePool.notifyConfig') }}</el-divider>
        <el-form-item :label="$t('clientsResourcePool.notifyMethod')" prop="notify_methods">
          <el-checkbox-group v-model="settingsForm.notify_methods">
            <el-checkbox label="email">{{ $t('clientsResourcePool.emailNotify') }}</el-checkbox>
            <el-checkbox label="sms">{{ $t('clientsResourcePool.smsNotify') }}</el-checkbox>
            <el-checkbox label="webhook">{{ $t('clientsResourcePool.webhookNotify') }}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item :label="$t('clientsResourcePool.webhookUrl')" prop="webhook_url" v-if="settingsForm.notify_methods.includes('webhook')">
          <el-input v-model="settingsForm.webhook_url" :placeholder="$t('clientsResourcePool.enterWebhookUrl')" />
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Check } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'ResourcePoolSettings' })

const loading = ref(false)
const submitLoading = ref(false)
const settingsFormRef = ref<FormInstance>()

const settingsForm = reactive({
  pool_name: '', capacity: 100, status: 1, pool_type: 'shared', description: '',
  allocation_strategy: 'round_robin', max_allocation: 50, allocation_threshold: 80, cooldown_seconds: 60,
  auto_assign: true, auto_reclaim: true, idle_timeout: 30, notify_threshold: 90,
  notify_methods: ['email'] as string[], webhook_url: ''
})

const settingsRules: FormRules = {
  pool_name: [
    { required: true, message: $t('clientsResourcePool.enterPoolName'), trigger: 'blur' },
    { min: 2, max: 50, message: $t('clientsResourcePool.nameLength'), trigger: 'blur' }
  ],
  capacity: [{ required: true, message: $t('clientsResourcePool.enterCapacity'), trigger: 'blur' }],
  pool_type: [{ required: true, message: $t('clientsResourcePool.selectPoolType'), trigger: 'change' }],
  allocation_strategy: [{ required: true, message: $t('clientsResourcePool.selectStrategy'), trigger: 'change' }]
}

const fetchSettings = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/resource-pool/settings' })
    Object.assign(settingsForm, data)
  } catch (error) {
    ElMessage.error($t('clientsResourcePool.fetchSettingsFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!settingsFormRef.value) return
  await settingsFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      await request.put({ url: '/api/admin/resource-pool/settings', params: settingsForm, showSuccessMessage: true })
      ElMessage.success($t('common.saveSuccess'))
    } catch (error) {
      ElMessage.error($t('common.saveFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

onMounted(() => { fetchSettings() })
</script>

<style scoped lang="scss">
.resource-pool-settings-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.form-tip { font-size: 12px; color: #909399; margin-top: 4px; }
:deep(.el-divider__text) { font-weight: 600; color: #303133; }
</style>
