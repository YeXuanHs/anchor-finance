<template>
  <div class="ticket-password-config-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>工单密码配置</span>
        </div>
      </template>

      <!-- 配置表单 -->
      <el-form :model="configForm" :rules="configRules" ref="configFormRef" label-width="120px" class="config-form">
        <el-form-item label="启用密码" prop="enable_password">
          <el-switch v-model="configForm.enable_password" active-text="启用" inactive-text="禁用" />
          <div class="form-tip">启用后，用户回复工单时需要输入密码</div>
        </el-form-item>

        <el-form-item label="密码规则" prop="password_rule" v-if="configForm.enable_password">
          <el-radio-group v-model="configForm.password_rule">
            <el-radio value="fixed">固定密码</el-radio>
            <el-radio value="random">随机密码</el-radio>
            <el-radio value="custom">自定义规则</el-radio>
          </el-radio-group>
          <div class="form-tip">
            <template v-if="configForm.password_rule === 'fixed'">所有工单使用相同的固定密码</template>
            <template v-else-if="configForm.password_rule === 'random'">系统自动生成随机密码</template>
            <template v-else>根据规则生成密码</template>
          </div>
        </el-form-item>

        <el-form-item label="固定密码" prop="fixed_password" v-if="configForm.enable_password && configForm.password_rule === 'fixed'">
          <el-input v-model="configForm.fixed_password" placeholder="请输入固定密码" show-password />
        </el-form-item>

        <el-form-item label="密码长度" prop="password_length" v-if="configForm.enable_password && configForm.password_rule === 'random'">
          <el-input-number v-model="configForm.password_length" :min="4" :max="20" />
          <div class="form-tip">生成随机密码的长度</div>
        </el-form-item>

        <el-form-item label="密码字符" prop="password_chars" v-if="configForm.enable_password && configForm.password_rule === 'custom'">
          <el-input v-model="configForm.password_chars" placeholder="允许的字符集" />
          <div class="form-tip">例如：0123456789abcdef</div>
        </el-form-item>

        <el-form-item label="密码长度" prop="password_length" v-if="configForm.enable_password && configForm.password_rule === 'custom'">
          <el-input-number v-model="configForm.password_length" :min="4" :max="20" />
          <div class="form-tip">生成密码的长度</div>
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
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

interface PasswordConfig {
  enable_password: boolean
  password_rule: 'fixed' | 'random' | 'custom'
  fixed_password: string
  password_length: number
  password_chars: string
}

const configFormRef = ref<FormInstance>()
const saveLoading = ref(false)

const configForm = reactive<PasswordConfig>({
  enable_password: false,
  password_rule: 'fixed',
  fixed_password: '',
  password_length: 6,
  password_chars: '0123456789'
})

const configRules: FormRules = {
  fixed_password: [
    { required: true, message: '请输入固定密码', trigger: 'blur' }
  ],
  password_length: [
    { required: true, message: '请输入密码长度', trigger: 'blur' }
  ],
  password_chars: [
    { required: true, message: '请输入密码字符', trigger: 'blur' }
  ]
}

const fetchConfig = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/ticket-prereply/password-config'
    })
    if (data) {
      Object.assign(configForm, data)
    }
  } catch (error) {
    console.error('获取密码配置失败:', error)
  }
}

const handleSave = async () => {
  if (!configFormRef.value) return

  await configFormRef.value.validate(async (valid) => {
    if (!valid) return

    saveLoading.value = true
    try {
      await request.post({
        url: '/api/admin/ticket-prereply/password-config',
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
  configForm.enable_password = false
  configForm.password_rule = 'fixed'
  configForm.fixed_password = ''
  configForm.password_length = 6
  configForm.password_chars = '0123456789'
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped lang="scss">
.ticket-password-config-page {
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