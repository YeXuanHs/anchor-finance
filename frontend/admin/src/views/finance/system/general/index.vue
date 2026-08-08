<template>
  <div class="general-settings-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>常规设置</span>
          <el-button type="primary" :loading="saving" @click="handleSave">
            <el-icon><Check /></el-icon>
            保存设置
          </el-button>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="150px"
        size="default"
      >
        <!-- 站点信息 -->
        <el-divider content-position="left">站点信息</el-divider>

        <el-form-item label="站点名称" prop="site_name">
          <el-input v-model="formData.site_name" placeholder="请输入站点名称" style="width: 400px" />
        </el-form-item>

        <el-form-item label="站点URL" prop="site_url">
          <el-input v-model="formData.site_url" placeholder="https://example.com" style="width: 400px" />
        </el-form-item>

        <el-form-item label="站点Logo" prop="site_logo">
          <el-input v-model="formData.site_logo" placeholder="Logo URL" style="width: 400px" />
        </el-form-item>

        <el-form-item label="站点描述" prop="site_description">
          <el-input
            v-model="formData.site_description"
            type="textarea"
            :rows="3"
            placeholder="请输入站点描述"
            style="width: 400px"
          />
        </el-form-item>

        <el-form-item label="关键词" prop="site_keywords">
          <el-input v-model="formData.site_keywords" placeholder="多个关键词用逗号分隔" style="width: 400px" />
        </el-form-item>

        <!-- 联系信息 -->
        <el-divider content-position="left">联系信息</el-divider>

        <el-form-item label="联系邮箱" prop="contact_email">
          <el-input v-model="formData.contact_email" placeholder="admin@example.com" style="width: 400px" />
        </el-form-item>

        <el-form-item label="联系电话" prop="contact_phone">
          <el-input v-model="formData.contact_phone" placeholder="请输入联系电话" style="width: 400px" />
        </el-form-item>

        <el-form-item label="联系地址" prop="contact_address">
          <el-input v-model="formData.contact_address" placeholder="请输入联系地址" style="width: 400px" />
        </el-form-item>

        <!-- 时间设置 -->
        <el-divider content-position="left">时间设置</el-divider>

        <el-form-item label="时区" prop="timezone">
          <el-select v-model="formData.timezone" placeholder="请选择时区" style="width: 400px">
            <el-option label="Asia/Shanghai (UTC+8)" value="Asia/Shanghai" />
            <el-option label="America/New_York (UTC-5)" value="America/New_York" />
            <el-option label="Europe/London (UTC+0)" value="Europe/London" />
            <el-option label="Asia/Tokyo (UTC+9)" value="Asia/Tokyo" />
          </el-select>
        </el-form-item>

        <el-form-item label="日期格式" prop="date_format">
          <el-select v-model="formData.date_format" placeholder="请选择日期格式" style="width: 400px">
            <el-option label="Y-m-d H:i:s" value="Y-m-d H:i:s" />
            <el-option label="Y/m/d H:i:s" value="Y/m/d H:i:s" />
            <el-option label="d/m/Y H:i:s" value="d/m/Y H:i:s" />
            <el-option label="m/d/Y H:i:s" value="m/d/Y H:i:s" />
          </el-select>
        </el-form-item>

        <!-- 注册设置 -->
        <el-divider content-position="left">注册设置</el-divider>

        <el-form-item label="允许注册" prop="allow_registration">
          <el-switch v-model="formData.allow_registration" />
        </el-form-item>

        <el-form-item label="邮箱验证" prop="email_verification">
          <el-switch v-model="formData.email_verification" />
        </el-form-item>

        <el-form-item label="手机验证" prop="phone_verification">
          <el-switch v-model="formData.phone_verification" />
        </el-form-item>

        <el-form-item label="注册验证码" prop="registration_captcha">
          <el-switch v-model="formData.registration_captcha" />
        </el-form-item>

        <!-- 登录设置 -->
        <el-divider content-position="left">登录设置</el-divider>

        <el-form-item label="登录验证码" prop="login_captcha">
          <el-switch v-model="formData.login_captcha" />
        </el-form-item>

        <el-form-item label="登录失败次数" prop="login_attempts">
          <el-input-number v-model="formData.login_attempts" :min="1" :max="10" />
          <span class="form-tip">次失败后需要验证码</span>
        </el-form-item>

        <el-form-item label="会话超时" prop="session_timeout">
          <el-input-number v-model="formData.session_timeout" :min="5" :max="1440" />
          <span class="form-tip">分钟</span>
        </el-form-item>

        <!-- 其他设置 -->
        <el-divider content-position="left">其他设置</el-divider>

        <el-form-item label="默认语言" prop="default_language">
          <el-select v-model="formData.default_language" placeholder="请选择语言" style="width: 200px">
            <el-option label="简体中文" value="zh-CN" />
            <el-option label="English" value="en" />
          </el-select>
        </el-form-item>

        <el-form-item label="默认货币" prop="default_currency">
          <el-select v-model="formData.default_currency" placeholder="请选择货币" style="width: 200px">
            <el-option label="CNY (¥)" value="CNY" />
            <el-option label="USD ($)" value="USD" />
            <el-option label="EUR (€)" value="EUR" />
          </el-select>
        </el-form-item>

        <el-form-item label="维护模式" prop="maintenance_mode">
          <el-switch v-model="formData.maintenance_mode" />
          <span class="form-tip">开启后前台将显示维护页面</span>
        </el-form-item>

        <el-form-item label="维护提示" prop="maintenance_message" v-if="formData.maintenance_mode">
          <el-input
            v-model="formData.maintenance_message"
            type="textarea"
            :rows="2"
            placeholder="维护提示信息"
            style="width: 400px"
          />
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Check } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const formRef = ref<FormInstance>()
const saving = ref(false)

// 表单数据
const formData = reactive({
  site_name: '',
  site_url: '',
  site_logo: '',
  site_description: '',
  site_keywords: '',
  contact_email: '',
  contact_phone: '',
  contact_address: '',
  timezone: 'Asia/Shanghai',
  date_format: 'Y-m-d H:i:s',
  allow_registration: true,
  email_verification: false,
  phone_verification: false,
  registration_captcha: true,
  login_captcha: true,
  login_attempts: 5,
  session_timeout: 120,
  default_language: 'zh-CN',
  default_currency: 'CNY',
  maintenance_mode: false,
  maintenance_message: ''
})

// 表单验证规则
const rules: FormRules = {
  site_name: [
    { required: true, message: '请输入站点名称', trigger: 'blur' }
  ],
  site_url: [
    { required: true, message: '请输入站点URL', trigger: 'blur' }
  ],
  contact_email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

// 获取设置
const fetchSettings = async () => {
  try {
    const data = await request.get({ url: '/api/admin/settings/general' })
    Object.assign(formData, data)
  } catch (error) {
    console.error('获取设置失败:', error)
  }
}

// 保存设置
const handleSave = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    saving.value = true

    await request.post({ url: '/api/admin/settings/general', data: formData })
    ElMessage.success('保存成功')
  } catch (error) {
    console.error('保存设置失败:', error)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchSettings()
})
</script>

<style scoped lang="scss">
.general-settings-page {
  padding: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-tip {
  margin-left: 10px;
  font-size: 12px;
  color: #86909C;
}

:deep(.el-divider__text) {
  font-weight: 600;
  color: #1D2129;
}
</style>
