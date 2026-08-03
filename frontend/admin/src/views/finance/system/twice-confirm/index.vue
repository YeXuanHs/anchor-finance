<template>
  <div class="twice-confirm-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>二次确认配置</span>
        </div>
      </template>

      <el-form :model="configForm" ref="configFormRef" label-width="160px" class="config-form">
        <el-divider content-position="left">敏感操作确认</el-divider>

        <el-form-item label="删除客户">
          <el-switch v-model="configForm.delete_client" />
          <span class="form-tip">删除客户时需要二次确认</span>
        </el-form-item>

        <el-form-item label="删除订单">
          <el-switch v-model="configForm.delete_order" />
          <span class="form-tip">删除订单时需要二次确认</span>
        </el-form-item>

        <el-form-item label="批量操作">
          <el-switch v-model="configForm.batch_operation" />
          <span class="form-tip">批量删除/修改时需要二次确认</span>
        </el-form-item>

        <el-form-item label="资金操作">
          <el-switch v-model="configForm.fund_operation" />
          <span class="form-tip">充值/退款/扣款时需要二次确认</span>
        </el-form-item>

        <el-form-item label="修改密码">
          <el-switch v-model="configForm.change_password" />
          <span class="form-tip">修改管理员密码时需要二次确认</span>
        </el-form-item>

        <el-form-item label="系统设置">
          <el-switch v-model="configForm.system_settings" />
          <span class="form-tip">修改关键系统设置时需要二次确认</span>
        </el-form-item>

        <el-divider content-position="left">确认方式</el-divider>

        <el-form-item label="确认方式">
          <el-checkbox-group v-model="configForm.methods">
            <el-checkbox value="password">密码确认</el-checkbox>
            <el-checkbox value="email_code">邮箱验证码</el-checkbox>
            <el-checkbox value="sms_code">短信验证码</el-checkbox>
            <el-checkbox value="totp">TOTP动态码</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="验证码有效期">
          <el-input-number v-model="configForm.code_expire" :min="30" :max="300" :step="30" />
          秒
        </el-form-item>

        <el-form-item label="最大重试次数">
          <el-input-number v-model="configForm.max_retries" :min="1" :max="10" />
          <span class="form-tip">超过次数锁定操作</span>
        </el-form-item>

        <el-form-item label="锁定时间">
          <el-input-number v-model="configForm.lock_duration" :min="60" :max="3600" :step="60" />
          秒
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="saveLoading">保存配置</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import request from '@/utils/http'

const configFormRef = ref<FormInstance>()
const saveLoading = ref(false)

const configForm = reactive({
  delete_client: true,
  delete_order: true,
  batch_operation: true,
  fund_operation: true,
  change_password: true,
  system_settings: false,
  methods: ['password'] as string[],
  code_expire: 60,
  max_retries: 3,
  lock_duration: 300
})

const fetchConfig = async () => {
  try {
    const data = await request.get({ url: '/api/admin/config/second-verify' })
    if (data) {
      Object.assign(configForm, data)
    }
  } catch (error) {
    console.error('获取配置失败:', error)
  }
}

const handleSave = async () => {
  saveLoading.value = true
  try {
    await request.put({
      url: '/api/admin/config/second-verify',
      params: { ...configForm }
    })
    ElMessage.success('保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saveLoading.value = false
  }
}

const handleReset = () => {
  fetchConfig()
}

onMounted(() => { fetchConfig() })
</script>

<style scoped lang="scss">
.twice-confirm-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.config-form { max-width: 700px; }
.form-tip { margin-left: 12px; color: var(--el-text-color-secondary); font-size: 13px; }
</style>
