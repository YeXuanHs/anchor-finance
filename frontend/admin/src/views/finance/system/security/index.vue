<template>
  <div class="security-settings-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>安全设置</span>
          <el-button type="primary" :loading="saving" @click="handleSave">保存设置</el-button>
        </div>
      </template>

      <el-form :model="formData" label-width="150px">
        <el-divider content-position="left">密码策略</el-divider>
        <el-form-item label="最小密码长度">
          <el-input-number v-model="formData.min_password_length" :min="6" :max="32" />
        </el-form-item>
        <el-form-item label="密码复杂度">
          <el-checkbox-group v-model="formData.password_complexity">
            <el-checkbox value="uppercase">大写字母</el-checkbox>
            <el-checkbox value="lowercase">小写字母</el-checkbox>
            <el-checkbox value="number">数字</el-checkbox>
            <el-checkbox value="special">特殊字符</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="密码有效期">
          <el-input-number v-model="formData.password_expiry_days" :min="0" :max="365" />
          <span class="form-tip">天（0表示永不过期）</span>
        </el-form-item>

        <el-divider content-position="left">登录安全</el-divider>
        <el-form-item label="最大登录失败次数">
          <el-input-number v-model="formData.max_login_attempts" :min="3" :max="20" />
          <span class="form-tip">次</span>
        </el-form-item>
        <el-form-item label="锁定时间">
          <el-input-number v-model="formData.lockout_duration" :min="5" :max="1440" />
          <span class="form-tip">分钟</span>
        </el-form-item>
        <el-form-item label="会话超时">
          <el-input-number v-model="formData.session_timeout" :min="5" :max="1440" />
          <span class="form-tip">分钟</span>
        </el-form-item>
        <el-form-item label="单点登录">
          <el-switch v-model="formData.single_session" />
          <span class="form-tip">同一账号只允许在一个设备登录</span>
        </el-form-item>

        <el-divider content-position="left">IP安全</el-divider>
        <el-form-item label="IP白名单">
          <el-input v-model="formData.ip_whitelist" type="textarea" :rows="3" placeholder="每行一个IP地址，留空表示不限制" style="width: 400px" />
        </el-form-item>
        <el-form-item label="IP黑名单">
          <el-input v-model="formData.ip_blacklist" type="textarea" :rows="3" placeholder="每行一个IP地址" style="width: 400px" />
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
    const data = await request.get({ url: '/api/admin/settings/security' })
    Object.assign(formData, data)
  } catch (error) {
    console.error('获取安全设置失败:', error)
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await request.post({ url: '/api/admin/settings/security', data: formData })
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
.security-settings-page { padding: 16px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.form-tip { margin-left: 10px; font-size: 12px; color: #86909C; }
:deep(.el-divider__text) { font-weight: 600; color: #1D2129; }
</style>
