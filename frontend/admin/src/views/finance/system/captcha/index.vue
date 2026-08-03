<template>
  <div class="captcha-config-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>验证码配置</span>
        </div>
      </template>

      <el-form :model="configForm" :rules="configRules" ref="configFormRef" label-width="140px" class="config-form">
        <el-form-item label="启用验证码" prop="enabled">
          <el-switch v-model="configForm.enabled" active-text="启用" inactive-text="禁用" />
          <div class="form-tip">启用后，登录和注册时需要输入验证码</div>
        </el-form-item>

        <el-form-item label="验证码类型" prop="type" v-if="configForm.enabled">
          <el-radio-group v-model="configForm.type">
            <el-radio value="image">图片验证码</el-radio>
            <el-radio value="slider">滑块验证码</el-radio>
            <el-radio value="recaptcha">Google reCAPTCHA</el-radio>
            <el-radio value="turnstile">Cloudflare Turnstile</el-radio>
          </el-radio-group>
        </el-form-item>

        <template v-if="configForm.enabled">
          <!-- 图片验证码配置 -->
          <template v-if="configForm.type === 'image'">
            <el-form-item label="验证码长度">
              <el-input-number v-model="configForm.image_length" :min="4" :max="8" />
            </el-form-item>
            <el-form-item label="验证码宽度">
              <el-input-number v-model="configForm.image_width" :min="80" :max="200" :step="10" />
              px
            </el-form-item>
            <el-form-item label="验证码高度">
              <el-input-number v-model="configForm.image_height" :min="30" :max="80" :step="5" />
              px
            </el-form-item>
            <el-form-item label="干扰线数量">
              <el-input-number v-model="configForm.noise_lines" :min="0" :max="10" />
            </el-form-item>
          </template>

          <!-- reCAPTCHA 配置 -->
          <template v-if="configForm.type === 'recaptcha'">
            <el-form-item label="Site Key" prop="recaptcha_site_key">
              <el-input v-model="configForm.recaptcha_site_key" placeholder="请输入 reCAPTCHA Site Key" />
            </el-form-item>
            <el-form-item label="Secret Key" prop="recaptcha_secret_key">
              <el-input v-model="configForm.recaptcha_secret_key" type="password" placeholder="请输入 reCAPTCHA Secret Key" show-password />
            </el-form-item>
            <el-form-item label="版本">
              <el-radio-group v-model="configForm.recaptcha_version">
                <el-radio value="v2">v2</el-radio>
                <el-radio value="v3">v3</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="分数阈值" v-if="configForm.recaptcha_version === 'v3'">
              <el-slider v-model="configForm.recaptcha_score" :min="0" :max="1" :step="0.1" show-input />
            </el-form-item>
          </template>

          <!-- Turnstile 配置 -->
          <template v-if="configForm.type === 'turnstile'">
            <el-form-item label="Site Key" prop="turnstile_site_key">
              <el-input v-model="configForm.turnstile_site_key" placeholder="请输入 Turnstile Site Key" />
            </el-form-item>
            <el-form-item label="Secret Key" prop="turnstile_secret_key">
              <el-input v-model="configForm.turnstile_secret_key" type="password" placeholder="请输入 Turnstile Secret Key" show-password />
            </el-form-item>
          </template>

          <el-form-item label="适用场景">
            <el-checkbox-group v-model="configForm.scenes">
              <el-checkbox value="login">登录</el-checkbox>
              <el-checkbox value="register">注册</el-checkbox>
              <el-checkbox value="reset_password">找回密码</el-checkbox>
              <el-checkbox value="contact">联系我们</el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <el-form-item label="失败锁定次数">
            <el-input-number v-model="configForm.max_failures" :min="3" :max="20" />
            <span class="form-tip">连续失败后锁定IP</span>
          </el-form-item>

          <el-form-item label="锁定时间">
            <el-input-number v-model="configForm.lock_duration" :min="60" :max="3600" :step="60" />
            秒
          </el-form-item>
        </template>

        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="saveLoading">保存配置</el-button>
          <el-button @click="handleReset">重置</el-button>
          <el-button v-if="configForm.enabled" @click="handleTest">测试验证码</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const configFormRef = ref<FormInstance>()
const saveLoading = ref(false)

const configForm = reactive({
  enabled: false,
  type: 'image',
  image_length: 4,
  image_width: 120,
  image_height: 40,
  noise_lines: 3,
  recaptcha_site_key: '',
  recaptcha_secret_key: '',
  recaptcha_version: 'v2',
  recaptcha_score: 0.5,
  turnstile_site_key: '',
  turnstile_secret_key: '',
  scenes: ['login'] as string[],
  max_failures: 5,
  lock_duration: 300
})

const configRules: FormRules = {
  recaptcha_site_key: [{ required: true, message: '请输入 Site Key', trigger: 'blur' }],
  recaptcha_secret_key: [{ required: true, message: '请输入 Secret Key', trigger: 'blur' }],
  turnstile_site_key: [{ required: true, message: '请输入 Site Key', trigger: 'blur' }],
  turnstile_secret_key: [{ required: true, message: '请输入 Secret Key', trigger: 'blur' }]
}

const defaultConfig = { ...configForm }

const fetchConfig = async () => {
  try {
    const data = await request.get({ url: '/api/admin/system/captcha' })
    if (data) {
      Object.assign(configForm, data)
    }
  } catch (error) {
    console.error('获取验证码配置失败:', error)
  }
}

const handleSave = async () => {
  if (!configFormRef.value) return

  await configFormRef.value.validate(async (valid) => {
    if (!valid) return

    saveLoading.value = true
    try {
      await request.post({
        url: '/api/admin/system/captcha',
        params: { ...configForm }
      })
      ElMessage.success('保存成功')
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      saveLoading.value = false
    }
  })
}

const handleReset = () => {
  Object.assign(configForm, defaultConfig)
}

const handleTest = async () => {
  try {
    await request.post({ url: '/api/admin/system/captcha/test' })
    ElMessage.success('验证码测试通过')
  } catch (error) {
    ElMessage.error('验证码测试失败')
  }
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped lang="scss">
.captcha-config-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.config-form {
  max-width: 700px;
}

.form-tip {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 5px;
}
</style>
