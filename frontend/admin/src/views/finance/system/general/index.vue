<template>
  <div class="general-settings-page">
    <h2 class="page-title">常规设置</h2>
    
    <el-tabs v-model="activeTab">
      <!-- 基础信息 -->
      <el-tab-pane label="基础信息" name="basic">
        <el-form :model="formData" label-width="120px" label-position="left">
          <el-form-item label="*品牌名">
            <el-input v-model="formData.site_name" placeholder="请输入品牌名" />
          </el-form-item>
          <el-form-item label="*系统链接">
            <el-input v-model="formData.site_url" placeholder="请输入您的安装地址" />
          </el-form-item>
          <el-form-item label="*网站域名">
            <el-input v-model="formData.site_domain" placeholder="请输入您网站主页的链接地址" />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- 本地化 -->
      <el-tab-pane label="本地化" name="locale">
        <el-form :model="formData" label-width="120px" label-position="left">
          <el-form-item label="默认语言">
            <el-select v-model="formData.default_language" style="width: 100%">
              <el-option label="简体中文" value="zh-CN" />
              <el-option label="繁體中文" value="zh-TW" />
              <el-option label="English" value="en" />
            </el-select>
          </el-form-item>
          <el-form-item label="时区">
            <el-select v-model="formData.timezone" style="width: 100%">
              <el-option label="Asia/Shanghai (UTC+8)" value="Asia/Shanghai" />
              <el-option label="America/New_York (UTC-5)" value="America/New_York" />
              <el-option label="Europe/London (UTC+0)" value="Europe/London" />
              <el-option label="Asia/Tokyo (UTC+9)" value="Asia/Tokyo" />
            </el-select>
          </el-form-item>
          <el-form-item label="日期格式">
            <el-select v-model="formData.date_format" style="width: 100%">
              <el-option label="Y-m-d H:i:s" value="Y-m-d H:i:s" />
              <el-option label="Y/m/d H:i:s" value="Y/m/d H:i:s" />
              <el-option label="d/m/Y H:i:s" value="d/m/Y H:i:s" />
            </el-select>
          </el-form-item>
          <el-form-item label="默认货币">
            <el-select v-model="formData.default_currency" style="width: 100%">
              <el-option label="CNY (¥)" value="CNY" />
              <el-option label="USD ($)" value="USD" />
              <el-option label="EUR (€)" value="EUR" />
            </el-select>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- 显示设置 -->
      <el-tab-pane label="显示设置" name="display">
        <el-form :model="formData" label-width="120px" label-position="left">
          <el-form-item label="网站Logo">
            <el-input v-model="formData.site_logo" placeholder="Logo URL" />
          </el-form-item>
          <el-form-item label="网站描述">
            <el-input v-model="formData.site_description" type="textarea" :rows="3" placeholder="请输入网站描述" />
          </el-form-item>
          <el-form-item label="SEO关键词">
            <el-input v-model="formData.site_keywords" placeholder="多个关键词用逗号分隔" />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- 基础安全 -->
      <el-tab-pane label="基础安全" name="security">
        <el-form :model="formData" label-width="120px" label-position="left">
          <el-form-item label="登录验证码">
            <el-switch v-model="formData.login_captcha" />
          </el-form-item>
          <el-form-item label="登录尝试次数">
            <el-input-number v-model="formData.login_attempts" :min="1" :max="10" />
            <span class="form-tip">超过次数后需要验证码</span>
          </el-form-item>
          <el-form-item label="会话超时">
            <el-input-number v-model="formData.session_timeout" :min="5" :max="1440" />
            <span class="form-tip">分钟</span>
          </el-form-item>
          <el-form-item label="注册验证码">
            <el-switch v-model="formData.registration_captcha" />
          </el-form-item>
          <el-form-item label="邮箱验证">
            <el-switch v-model="formData.email_verification" />
          </el-form-item>
          <el-form-item label="手机验证">
            <el-switch v-model="formData.phone_verification" />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- 维护模式 -->
      <el-tab-pane label="维护模式" name="maintenance">
        <el-form :model="formData" label-width="120px" label-position="left">
          <el-form-item label="开启维护">
            <el-switch v-model="formData.maintenance_mode" />
          </el-form-item>
          <el-form-item label="维护提示" v-if="formData.maintenance_mode">
            <el-input v-model="formData.maintenance_message" type="textarea" :rows="3" placeholder="请输入维护提示信息" />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- Debug调试 -->
      <el-tab-pane label="Debug调试" name="debug">
        <el-form :model="formData" label-width="120px" label-position="left">
          <el-form-item label="调试模式">
            <el-switch v-model="formData.debug_mode" />
          </el-form-item>
          <el-form-item label="错误日志">
            <el-switch v-model="formData.error_log" />
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>

    <div class="form-actions">
      <el-button type="primary" :loading="saving" @click="handleSave">保存更改</el-button>
      <el-button @click="fetchSettings">取消更改</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const activeTab = ref('basic')
const saving = ref(false)

const formData = reactive({
  site_name: '',
  site_url: '',
  site_domain: '',
  site_logo: '',
  site_description: '',
  site_keywords: '',
  default_language: 'zh-CN',
  timezone: 'Asia/Shanghai',
  date_format: 'Y-m-d H:i:s',
  default_currency: 'CNY',
  login_captcha: true,
  login_attempts: 5,
  session_timeout: 120,
  registration_captcha: true,
  email_verification: false,
  phone_verification: false,
  maintenance_mode: false,
  maintenance_message: '',
  debug_mode: false,
  error_log: false
})

const fetchSettings = async () => {
  try {
    const data = await request.get({ url: '/api/admin/config/general' })
    Object.assign(formData, data)
  } catch (error) {
    console.error('fetch settings failed:', error)
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await request.put({ url: '/api/admin/config/general', data: formData })
    ElMessage.success('保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
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
  padding: 20px;
}

.page-title {
  margin: 0 0 16px 0;
  font-size: 20px;
  font-weight: 600;
}

.form-tip {
  margin-left: 10px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.form-actions {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--el-border-color-lighter);
}
</style>
