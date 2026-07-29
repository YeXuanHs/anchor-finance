<template>
  <div class="certification-config-page page-container">
    <div class="art-card">
      <h3>实名认证配置</h3>
      <el-form :model="formData" label-width="150px" style="max-width: 600px;">
        <el-form-item label="启用实名认证">
          <el-switch v-model="formData.enabled" />
        </el-form-item>
        <el-form-item label="认证方式">
          <el-checkbox-group v-model="formData.methods">
            <el-checkbox value="manual">人工审核</el-checkbox>
            <el-checkbox value="auto">自动审核</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="支持证件类型">
          <el-checkbox-group v-model="formData.id_types">
            <el-checkbox value="id_card">身份证</el-checkbox>
            <el-checkbox value="passport">护照</el-checkbox>
            <el-checkbox value="driver_license">驾驶证</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="认证有效期">
          <el-input-number v-model="formData.valid_days" :min="0" />
          <span style="margin-left: 8px;">天（0为永久有效）</span>
        </el-form-item>
        <el-form-item label="未认证限制">
          <el-checkbox-group v-model="formData.restrictions">
            <el-checkbox value="order">禁止下单</el-checkbox>
            <el-checkbox value="ticket">禁止提交工单</el-checkbox>
            <el-checkbox value="withdraw">禁止提现</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="认证提示">
          <el-input v-model="formData.notice" type="textarea" :rows="3" placeholder="用户看到的认证提示信息" />
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
  enabled: false,
  methods: ['manual'],
  id_types: ['id_card'],
  valid_days: 0,
  restrictions: [],
  notice: ''
})

const fetchConfig = async () => {
  try {
    const { data } = await request.get('/admin/settings/certification')
    if (data?.data) {
      Object.assign(formData.value, data.data)
    }
  } catch {
    ElMessage.error('加载配置失败')
  }
}

const saveConfig = async () => {
  try {
    await request.put('/admin/settings/certification', formData.value)
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
.certification-config-page {
  h3 { margin: 0 0 24px; font-size: 16px; font-weight: 600; }
}
</style>
