<template>
  <div class="add-client-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>添加客户</span>
          <el-button @click="handleBack">
            <el-icon><Back /></el-icon>
            返回
          </el-button>
        </div>
      </template>

      <el-steps :active="currentStep" finish-status="success" align-center class="steps-bar">
        <el-step title="基本信息" />
        <el-step title="详细资料" />
        <el-step title="确认创建" />
      </el-steps>

      <!-- 步骤1: 基本信息 -->
      <div v-show="currentStep === 0" class="step-content">
        <el-form :model="formData" :rules="step1Rules" ref="step1Ref" label-width="120px" style="max-width: 600px;">
          <el-form-item label="用户名" prop="username">
            <el-input v-model="formData.username" placeholder="请输入用户名" />
          </el-form-item>
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="formData.email" placeholder="请输入邮箱" />
          </el-form-item>
          <el-form-item label="手机号" prop="phone">
            <el-input v-model="formData.phone" placeholder="请输入手机号" />
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input v-model="formData.password" type="password" placeholder="请输入密码" show-password />
          </el-form-item>
          <el-form-item label="确认密码" prop="confirm_password">
            <el-input v-model="formData.confirm_password" type="password" placeholder="请再次输入密码" show-password />
          </el-form-item>
          <el-form-item label="客户分组" prop="group_id">
            <el-select v-model="formData.group_id" placeholder="请选择分组">
              <el-option v-for="group in groups" :key="group.id" :label="group.name" :value="group.id" />
            </el-select>
          </el-form-item>
        </el-form>
      </div>

      <!-- 步骤2: 详细资料 -->
      <div v-show="currentStep === 1" class="step-content">
        <el-form :model="formData" label-width="120px" style="max-width: 600px;">
          <el-form-item label="姓名">
            <el-input v-model="formData.fullname" placeholder="请输入真实姓名" />
          </el-form-item>
          <el-form-item label="公司">
            <el-input v-model="formData.company" placeholder="请输入公司名称" />
          </el-form-item>
          <el-form-item label="职位">
            <el-input v-model="formData.position" placeholder="请输入职位" />
          </el-form-item>
          <el-form-item label="国家">
            <el-input v-model="formData.country" placeholder="请输入国家" />
          </el-form-item>
          <el-form-item label="省份">
            <el-input v-model="formData.state" placeholder="请输入省份" />
          </el-form-item>
          <el-form-item label="城市">
            <el-input v-model="formData.city" placeholder="请输入城市" />
          </el-form-item>
          <el-form-item label="地址">
            <el-input v-model="formData.address" type="textarea" :rows="2" placeholder="请输入详细地址" />
          </el-form-item>
          <el-form-item label="邮编">
            <el-input v-model="formData.postcode" placeholder="请输入邮编" />
          </el-form-item>
          <el-form-item label="语言">
            <el-select v-model="formData.language" placeholder="请选择语言">
              <el-option label="简体中文" value="zh-CN" />
              <el-option label="English" value="en" />
            </el-select>
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="客户备注" />
          </el-form-item>
        </el-form>
      </div>

      <!-- 步骤3: 确认创建 -->
      <div v-show="currentStep === 2" class="step-content">
        <el-card shadow="never">
          <template #header>
            <span>确认客户信息</span>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="用户名">{{ formData.username }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ formData.email }}</el-descriptions-item>
            <el-descriptions-item label="手机号">{{ formData.phone || '-' }}</el-descriptions-item>
            <el-descriptions-item label="分组">{{ getGroupName(formData.group_id) }}</el-descriptions-item>
            <el-descriptions-item label="姓名">{{ formData.fullname || '-' }}</el-descriptions-item>
            <el-descriptions-item label="公司">{{ formData.company || '-' }}</el-descriptions-item>
            <el-descriptions-item label="国家">{{ formData.country || '-' }}</el-descriptions-item>
            <el-descriptions-item label="城市">{{ formData.city || '-' }}</el-descriptions-item>
            <el-descriptions-item label="地址" :span="2">{{ formData.address || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </div>

      <!-- 底部按钮 -->
      <div class="step-actions">
        <el-button v-if="currentStep > 0" @click="currentStep--">上一步</el-button>
        <el-button v-if="currentStep < 2" type="primary" @click="handleNext">下一步</el-button>
        <el-button v-if="currentStep === 2" type="primary" @click="handleSubmit" :loading="submitLoading">
          确认创建
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Back } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const router = useRouter()
const currentStep = ref(0)
const submitLoading = ref(false)
const step1Ref = ref<FormInstance>()

const groups = ref<any[]>([])

const formData = reactive({
  username: '',
  email: '',
  phone: '',
  password: '',
  confirm_password: '',
  group_id: undefined as number | undefined,
  fullname: '',
  company: '',
  position: '',
  country: '',
  state: '',
  city: '',
  address: '',
  postcode: '',
  language: 'zh-CN',
  remark: ''
})

const step1Rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '长度在 3 到 50 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string, callback: Function) => {
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

const getGroupName = (id: number | undefined) => {
  if (!id) return '默认分组'
  return groups.value.find(g => g.id === id)?.name || '未知分组'
}

const fetchGroups = async () => {
  try {
    const data = await request.get({ url: '/api/admin/client-groups' })
    groups.value = data || []
  } catch (error) {
    console.error('获取分组列表失败:', error)
  }
}

const handleNext = async () => {
  if (currentStep.value === 0 && step1Ref.value) {
    const valid = await step1Ref.value.validate().catch(() => false)
    if (!valid) return
  }
  currentStep.value++
}

const handleSubmit = async () => {
  submitLoading.value = true
  try {
    await request.post({
      url: '/api/admin/user-manage',
      params: { ...formData }
    })
    ElMessage.success('客户创建成功')
    router.push('/finance/clients/list')
  } catch (error) {
    ElMessage.error('创建失败')
  } finally {
    submitLoading.value = false
  }
}

const handleBack = () => {
  router.back()
}

onMounted(() => {
  fetchGroups()
})
</script>

<style scoped lang="scss">
.add-client-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.steps-bar {
  margin: 20px 0 32px;
}

.step-content {
  min-height: 300px;
  padding: 20px 0;
}

.step-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--el-border-color-lighter);
}
</style>
