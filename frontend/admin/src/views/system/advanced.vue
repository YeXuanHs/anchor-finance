<template>
  <div class="advanced-config-page page-container">
    <div class="art-card">
      <h3>高级设置</h3>
      <el-form :model="formData" label-width="150px" style="max-width: 600px;">
        <el-form-item label="调试模式">
          <el-switch v-model="formData.debug" />
          <span style="margin-left: 8px; color: var(--text-secondary);">生产环境请勿开启</span>
        </el-form-item>
        <el-form-item label="缓存驱动">
          <el-select v-model="formData.cache_driver" placeholder="请选择">
            <el-option label="文件" value="file" />
            <el-option label="Redis" value="redis" />
            <el-option label="Memcached" value="memcached" />
          </el-select>
        </el-form-item>
        <el-form-item label="会话驱动">
          <el-select v-model="formData.session_driver" placeholder="请选择">
            <el-option label="文件" value="file" />
            <el-option label="Redis" value="redis" />
            <el-option label="数据库" value="database" />
          </el-select>
        </el-form-item>
        <el-form-item label="队列驱动">
          <el-select v-model="formData.queue_driver" placeholder="请选择">
            <el-option label="同步" value="sync" />
            <el-option label="Redis" value="redis" />
            <el-option label="数据库" value="database" />
          </el-select>
        </el-form-item>
        <el-form-item label="API限流">
          <el-input-number v-model="formData.rate_limit" :min="0" />
          <span style="margin-left: 8px;">请求/分钟（0为不限制）</span>
        </el-form-item>
        <el-form-item label="上传文件大小限制">
          <el-input-number v-model="formData.upload_max_size" :min="1" />
          <span style="margin-left: 8px;">MB</span>
        </el-form-item>
        <el-form-item label="允许的文件类型">
          <el-input v-model="formData.allowed_types" placeholder="jpg,png,gif,pdf,doc" />
        </el-form-item>
        <el-form-item label="时区">
          <el-select v-model="formData.timezone" placeholder="请选择">
            <el-option label="Asia/Shanghai" value="Asia/Shanghai" />
            <el-option label="UTC" value="UTC" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="saveConfig">保存配置</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const formData = ref({
  debug: false,
  cache_driver: 'file',
  session_driver: 'file',
  queue_driver: 'sync',
  rate_limit: 60,
  upload_max_size: 10,
  allowed_types: 'jpg,png,gif,pdf,doc,docx,xls,xlsx',
  timezone: 'Asia/Shanghai'
})

const fetchConfig = async () => {
  try {
    const { data } = await request.get('/admin/settings/advanced')
    if (data?.data) {
      Object.assign(formData.value, data.data)
    }
  } catch {
    ElMessage.error('加载配置失败')
  }
}

const saveConfig = async () => {
  try {
    await request.put('/admin/settings/advanced', formData.value)
    ElMessage.success('配置保存成功')
  } catch {
    ElMessage.error('配置保存失败')
  }
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped lang="scss">
.advanced-config-page {
  h3 { margin: 0 0 24px; font-size: 16px; font-weight: 600; }
}
</style>
