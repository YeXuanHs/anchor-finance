<template>
  <div class="add-customer-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clients.add.title') }}</span>
          <el-button @click="$router.back()">{{ $t('common.back') }}</el-button>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px" size="default">
        <!-- 基本信息 -->
        <el-divider content-position="left">{{ $t('clients.add.basicInfo') }}</el-divider>

        <el-form-item :label="$t('clients.add.username')" prop="username">
          <el-input v-model="formData.username" :placeholder="$t('clients.add.usernamePlaceholder')" style="width: 400px" />
        </el-form-item>

        <el-form-item :label="$t('clients.add.email')" prop="email">
          <el-input v-model="formData.email" :placeholder="$t('clients.add.emailPlaceholder')" style="width: 400px" />
        </el-form-item>

        <el-form-item :label="$t('clients.add.password')" prop="password">
          <el-input v-model="formData.password" type="password" :placeholder="$t('clients.add.passwordPlaceholder')" show-password style="width: 400px" />
        </el-form-item>

        <el-form-item :label="$t('clients.add.confirmPassword')" prop="confirm_password">
          <el-input v-model="formData.confirm_password" type="password" :placeholder="$t('clients.add.confirmPasswordPlaceholder')" show-password style="width: 400px" />
        </el-form-item>

        <el-form-item :label="$t('clients.add.phone')" prop="phone">
          <el-input v-model="formData.phone" :placeholder="$t('clients.add.phonePlaceholder')" style="width: 400px" />
        </el-form-item>

        <el-form-item :label="$t('clients.add.company')" prop="company">
          <el-input v-model="formData.company" :placeholder="$t('clients.add.companyPlaceholder')" style="width: 400px" />
        </el-form-item>

        <!-- 账户信息 -->
        <el-divider content-position="left">{{ $t('clients.add.accountInfo') }}</el-divider>

        <el-form-item :label="$t('clients.add.clientGroup')" prop="group_id">
          <el-select v-model="formData.group_id" :placeholder="$t('clients.add.clientGroupPlaceholder')" style="width: 400px">
            <el-option v-for="group in clientGroups" :key="group.id" :label="group.name" :value="group.id" />
          </el-select>
        </el-form-item>

        <el-form-item :label="$t('clients.add.balance')" prop="balance">
          <el-input-number v-model="formData.balance" :min="0" :precision="2" :step="100" />
          <span class="form-tip">{{ $t('clients.add.yuan') }}</span>
        </el-form-item>

        <el-form-item :label="$t('clients.add.credit')" prop="credit">
          <el-input-number v-model="formData.credit" :min="0" :precision="2" :step="100" />
          <span class="form-tip">{{ $t('clients.add.yuan') }}</span>
        </el-form-item>

        <el-form-item :label="$t('clients.add.status')" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio value="active">{{ $t('clients.add.active') }}</el-radio>
            <el-radio value="disabled">{{ $t('clients.add.disabled') }}</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 备注 -->
        <el-divider content-position="left">{{ $t('clients.add.otherInfo') }}</el-divider>

        <el-form-item :label="$t('clients.add.notes')" prop="notes">
          <el-input v-model="formData.notes" type="textarea" :rows="4" :placeholder="$t('clients.add.notesPlaceholder')" style="width: 400px" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">{{ $t('common.submit') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const router = useRouter()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const clientGroups = ref<{ id: number; name: string }[]>([])

const formData = reactive({
  username: '',
  email: '',
  password: '',
  confirm_password: '',
  phone: '',
  company: '',
  group_id: null as number | null,
  balance: 0,
  credit: 0,
  status: 'active',
  notes: ''
})

const rules = computed<FormRules>(() => ({
  username: [
    { required: true, message: $t('clients.add.usernameRequired'), trigger: 'blur' },
    { min: 3, max: 20, message: $t('clients.add.usernameLength'), trigger: 'blur' }
  ],
  email: [
    { required: true, message: $t('clients.add.emailRequired'), trigger: 'blur' },
    { type: 'email', message: $t('clients.add.emailInvalid'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: $t('clients.add.passwordRequired'), trigger: 'blur' },
    { min: 6, message: $t('clients.add.passwordLength'), trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: $t('clients.add.confirmPasswordRequired'), trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (value !== formData.password) {
          callback(new Error($t('clients.add.passwordMismatch')))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}))

const fetchGroups = async () => {
  try {
    const data = await request.get({ url: '/api/admin/client-groups' })
    clientGroups.value = data || []
  } catch (error) {
    console.error('Failed to fetch client groups:', error)
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    submitting.value = true
    await request.post({ url: '/api/admin/clients', data: formData })
    ElMessage.success($t('clients.add.addSuccess'))
    router.push('/customer-list')
  } catch (error) {
    console.error('Failed to add client:', error)
  } finally {
    submitting.value = false
  }
}

const handleReset = () => {
  formRef.value?.resetFields()
}

onMounted(() => {
  fetchGroups()
})
</script>

<style scoped lang="scss">
.add-customer-page {
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
