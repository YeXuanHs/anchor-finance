<template>
  <div class="add-customer-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>添加客户</span>
          <el-button @click="$router.back()">返回</el-button>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px" size="default">
        <!-- 基本信息 -->
        <el-divider content-position="left">基本信息</el-divider>

        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="请输入用户名" style="width: 400px" />
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="请输入邮箱" style="width: 400px" />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="formData.password" type="password" placeholder="请输入密码" show-password style="width: 400px" />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirm_password">
          <el-input v-model="formData.confirm_password" type="password" placeholder="请再次输入密码" show-password style="width: 400px" />
        </el-form-item>

        <el-form-item label="手机号" prop="phone">
          <el-input v-model="formData.phone" placeholder="请输入手机号" style="width: 400px" />
        </el-form-item>

        <el-form-item label="公司" prop="company">
          <el-input v-model="formData.company" placeholder="请输入公司名称" style="width: 400px" />
        </el-form-item>

        <!-- 账户信息 -->
        <el-divider content-position="left">账户信息</el-divider>

        <el-form-item label="客户组" prop="group_id">
          <el-select v-model="formData.group_id" placeholder="请选择客户组" style="width: 400px">
            <el-option v-for="group in clientGroups" :key="group.id" :label="group.name" :value="group.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="初始余额" prop="balance">
          <el-input-number v-model="formData.balance" :min="0" :precision="2" :step="100" />
          <span class="form-tip">元</span>
        </el-form-item>

        <el-form-item label="信用额" prop="credit">
          <el-input-number v-model="formData.credit" :min="0" :precision="2" :step="100" />
          <span class="form-tip">元</span>
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio value="active">正常</el-radio>
            <el-radio value="disabled">禁用</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 备注 -->
        <el-divider content-position="left">其他信息</el-divider>

        <el-form-item label="备注" prop="notes">
          <el-input v-model="formData.notes" type="textarea" :rows="4" placeholder="请输入备注" style="width: 400px" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">提交</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const router = useRouter()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const clientGroups = ref<{ id: number; name: string }[]>([])

// 表单数据
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

// 表单验证规则
const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (value !== formData.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 获取客户组
const fetchGroups = async () => {
  try {
    const data = await request.get({ url: '/api/admin/client-groups' })
    clientGroups.value = data || []
  } catch (error) {
    console.error('获取客户组失败:', error)
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    submitting.value = true
    await request.post({ url: '/api/admin/clients', data: formData })
    ElMessage.success('客户添加成功')
    router.push('/customer-list')
  } catch (error) {
    console.error('添加客户失败:', error)
  } finally {
    submitting.value = false
  }
}

// 重置表单
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

