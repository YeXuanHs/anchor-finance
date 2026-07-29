<template>
  <div class="security-page page-container">
    <div class="art-card">
      <div class="table-header">
        <h3>安全设置</h3>
      </div>

      <el-form
        ref="formRef"
        :model="formData"
        label-width="160px"
        style="max-width: 720px"
        v-loading="loading"
      >
        <!-- 登录安全 -->
        <el-divider content-position="left">登录安全</el-divider>

        <el-form-item label="登录验证码">
          <el-switch v-model="formData.captcha_enabled" />
          <span class="form-tip">开启后登录时需要输入图形验证码</span>
        </el-form-item>

        <el-form-item label="登录尝试限制">
          <el-input-number
            v-model="formData.max_login_attempts"
            :min="0"
            :max="20"
            controls-position="right"
          />
          <span class="form-tip">次（0为不限制，超出后账户将被锁定）</span>
        </el-form-item>

        <el-form-item label="账户锁定时间">
          <el-input-number
            v-model="formData.lockout_duration"
            :min="1"
            :max="1440"
            controls-position="right"
          />
          <span class="form-tip">分钟</span>
        </el-form-item>

        <!-- 密码策略 -->
        <el-divider content-position="left">密码策略</el-divider>

        <el-form-item label="密码最小长度">
          <el-input-number
            v-model="formData.password_min_length"
            :min="6"
            :max="32"
            controls-position="right"
          />
          <span class="form-tip">个字符</span>
        </el-form-item>

        <el-form-item label="密码复杂度要求">
          <el-checkbox-group v-model="formData.password_complexity">
            <el-checkbox label="uppercase">大写字母</el-checkbox>
            <el-checkbox label="lowercase">小写字母</el-checkbox>
            <el-checkbox label="number">数字</el-checkbox>
            <el-checkbox label="special">特殊字符</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="密码有效期">
          <el-input-number
            v-model="formData.password_expiry_days"
            :min="0"
            :max="365"
            controls-position="right"
          />
          <span class="form-tip">天（0为永不过期）</span>
        </el-form-item>

        <!-- 会话安全 -->
        <el-divider content-position="left">会话安全</el-divider>

        <el-form-item label="会话超时时间">
          <el-input-number
            v-model="formData.session_timeout"
            :min="5"
            :max="1440"
            controls-position="right"
          />
          <span class="form-tip">分钟（无操作后自动退出）</span>
        </el-form-item>

        <el-form-item label="单点登录">
          <el-switch v-model="formData.single_session" />
          <span class="form-tip">开启后同一账号只允许在一个设备登录</span>
        </el-form-item>

        <!-- 高级安全 -->
        <el-divider content-position="left">高级安全</el-divider>

        <el-form-item label="两步验证">
          <el-switch v-model="formData.two_factor_enabled" />
          <span class="form-tip">开启后登录时需要额外验证（短信/邮箱）</span>
        </el-form-item>

        <el-form-item label="IP白名单">
          <el-input
            v-model="formData.ip_whitelist"
            type="textarea"
            :rows="3"
            placeholder="每行一个IP地址，如：&#10;192.168.1.1&#10;10.0.0.0/24"
          />
          <span class="form-tip">留空表示不限制，仅白名单内的IP可访问后台</span>
        </el-form-item>

        <el-form-item label="登录日志保留">
          <el-input-number
            v-model="formData.log_retention_days"
            :min="7"
            :max="365"
            controls-position="right"
          />
          <span class="form-tip">天</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="saving" @click="handleSave">保存设置</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import request from '@/utils/request'

const formRef = ref<FormInstance>()
const loading = ref(false)
const saving = ref(false)

const formData = reactive({
  captcha_enabled: true,
  max_login_attempts: 5,
  lockout_duration: 30,
  password_min_length: 8,
  password_complexity: ['lowercase', 'number'],
  password_expiry_days: 90,
  session_timeout: 60,
  single_session: false,
  two_factor_enabled: false,
  ip_whitelist: '',
  log_retention_days: 90
})

const fetchConfig = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/admin/system/security')
    if (data?.data) {
      Object.assign(formData, data.data)
    }
  } catch {
    // 使用默认值
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await request.post('/api/admin/system/security', formData)
    ElMessage.success('安全设置保存成功')
  } catch {
    ElMessage.error('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

const handleReset = () => {
  fetchConfig()
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped lang="scss">
.security-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
    }
  }

  .form-tip {
    margin-left: 12px;
    font-size: 13px;
    color: var(--el-text-color-secondary);
  }

  :deep(.el-divider) {
    margin: 28px 0 20px;

    .el-divider__text {
      font-size: 14px;
      font-weight: 600;
      color: var(--el-text-color-primary);
    }
  }
}
</style>
