<template>
  <div class="maintenance-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>维护模式</span>
        </div>
      </template>

      <el-form :model="configForm" ref="formRef" label-width="120px" class="config-form">
        <el-form-item label="启用维护模式">
          <el-switch v-model="configForm.enabled" active-text="开启" inactive-text="关闭" />
          <div class="form-tip">开启后，前台将显示维护页面，仅允许白名单IP访问</div>
        </el-form-item>

        <el-form-item label="维护标题" prop="title" :rules="[{ required: true, message: '请输入维护标题', trigger: 'blur' }]">
          <el-input v-model="configForm.title" placeholder="如：系统维护中" />
        </el-form-item>

        <el-form-item label="维护提示" prop="message">
          <el-input v-model="configForm.message" type="textarea" :rows="4" placeholder="维护提示信息，支持HTML" />
        </el-form-item>

        <el-form-item label="预计完成时间">
          <el-date-picker
            v-model="configForm.estimated_end_time"
            type="datetime"
            placeholder="选择预计完成时间"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>

        <el-form-item label="IP白名单">
          <el-input
            v-model="configForm.allowed_ips"
            type="textarea"
            :rows="3"
            placeholder="每行一个IP地址，如：&#10;192.168.1.1&#10;10.0.0.0/24"
          />
          <div class="form-tip">白名单内的IP不受维护模式影响，可正常访问</div>
        </el-form-item>

        <el-form-item label="显示倒计时">
          <el-switch v-model="configForm.show_countdown" />
          <div class="form-tip">维护页面是否显示倒计时</div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="saveLoading">保存配置</el-button>
          <el-button @click="handlePreview">预览维护页面</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const saveLoading = ref(false)

const configForm = reactive({
  enabled: false,
  title: '系统维护中',
  message: '尊敬的用户，系统正在进行维护升级，预计很快恢复。给您带来的不便，敬请谅解。',
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
    ElMessage.success('保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
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
