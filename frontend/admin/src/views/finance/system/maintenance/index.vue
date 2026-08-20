<template>
  <div class="maintenance-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('page.maintenance.title') }}</span>
        </div>
      </template>

      <el-form :model="configForm" ref="formRef" label-width="120px" class="config-form">
        <el-form-item :label="$t('page.maintenance.enableMaintenance')">
          <el-switch v-model="configForm.enabled" :active-text="$t('page.maintenance.on')" :inactive-text="$t('page.maintenance.off')" />
          <div class="form-tip">{{ $t('page.maintenance.enableMaintenanceTip') }}</div>
        </el-form-item>

        <el-form-item :label="$t('page.maintenance.maintenanceTitle')" prop="title" :rules="[{ required: true, message: $t('page.maintenance.enterTitle'), trigger: 'blur' }]">
          <el-input v-model="configForm.title" :placeholder="$t('page.maintenance.titlePlaceholder')" />
        </el-form-item>

        <el-form-item :label="$t('page.maintenance.maintenanceMessage')" prop="message">
          <el-input v-model="configForm.message" type="textarea" :rows="4" :placeholder="$t('page.maintenance.messagePlaceholder')" />
        </el-form-item>

        <el-form-item :label="$t('page.maintenance.estimatedEndTime')">
          <el-date-picker
            v-model="configForm.estimated_end_time"
            type="datetime"
            :placeholder="$t('page.maintenance.selectEndTime')"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>

        <el-form-item :label="$t('page.maintenance.ipWhitelist')">
          <el-input
            v-model="configForm.allowed_ips"
            type="textarea"
            :rows="3"
            :placeholder="$t('page.maintenance.ipWhitelistPlaceholder')"
          />
          <div class="form-tip">{{ $t('page.maintenance.ipWhitelistTip') }}</div>
        </el-form-item>

        <el-form-item :label="$t('page.maintenance.showCountdown')">
          <el-switch v-model="configForm.show_countdown" />
          <div class="form-tip">{{ $t('page.maintenance.showCountdownTip') }}</div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="saveLoading">{{ $t('page.maintenance.saveConfig') }}</el-button>
          <el-button @click="handlePreview">{{ $t('page.maintenance.preview') }}</el-button>
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

const saveLoading = ref(false)

const configForm = reactive({
  enabled: false,
  title: $t('page.maintenance.defaultTitle'),
  message: $t('page.maintenance.defaultMessage'),
  estimated_end_time: '',
  allowed_ips: '',
  show_countdown: true
})

const fetchConfig = async () => {
  try {
    const data = await request.get({ url: '/api/admin/maintenance/status' })
    if (data) {
      Object.assign(configForm, data)
    }
  } catch (error) {
    console.error('获取维护配置失败:', error)
  }
}

const handleSave = async () => {
  saveLoading.value = true
  try {
    await request.put({ url: '/api/admin/maintenance/settings', params: { ...configForm } })
    ElMessage.success($t('page.maintenance.saveSuccess'))
  } catch (error) {
    ElMessage.error($t('page.maintenance.saveFailed'))
  } finally {
    saveLoading.value = false
  }
}

const handlePreview = () => {
  window.open('/maintenance-preview', '_blank')
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped lang="scss">
.maintenance-page {
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
