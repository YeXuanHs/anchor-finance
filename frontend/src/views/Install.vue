<template>
  <div class="install-container">
    <div class="install-card">
      <div class="install-header">
        <img src="/logo.png" alt="Logo" class="logo" />
        <h1>锚点财务 - 安装向导</h1>
        <p class="subtitle">欢迎使用锚点财务系统，请按照以下步骤完成安装</p>
      </div>

      <el-steps :active="currentStep" finish-status="success" class="steps">
        <el-step title="环境检测" />
        <el-step title="数据库配置" />
        <el-step title="管理员设置" />
        <el-step title="完成安装" />
      </el-steps>

      <!-- Step 0: 环境检测 -->
      <div v-if="currentStep === 0" class="step-content">
        <h3>环境检测</h3>
        <div class="env-checks">
          <div v-for="check in envChecks" :key="check.name" class="env-check-item">
            <div class="check-info">
              <el-icon v-if="check.status === 'checking'" class="is-loading"><Loading /></el-icon>
              <el-icon v-else-if="check.status === 'success'" class="success"><CircleCheckFilled /></el-icon>
              <el-icon v-else-if="check.status === 'error'" class="error"><CircleCloseFilled /></el-icon>
              <span class="check-name">{{ check.name }}</span>
            </div>
            <span class="check-message" :class="check.status">{{ check.message }}</span>
          </div>
        </div>
        <el-button type="primary" @click="checkEnvironment" :loading="isChecking">
          {{ isChecking ? '检测中...' : '开始检测' }}
        </el-button>
        <el-button v-if="allEnvChecksPass" type="success" @click="nextStep">
          下一步
        </el-button>
      </div>

      <!-- Step 1: 数据库配置 -->
      <div v-if="currentStep === 1" class="step-content">
        <h3>数据库配置</h3>
        <el-form :model="dbForm" :rules="dbRules" ref="dbFormRef" label-width="120px">
          <el-divider content-position="left">PostgreSQL 数据库</el-divider>
          <el-form-item label="主机" prop="db_host">
            <el-input v-model="dbForm.db_host" placeholder="localhost" />
          </el-form-item>
          <el-form-item label="端口" prop="db_port">
            <el-input v-model="dbForm.db_port" placeholder="5432" />
          </el-form-item>
          <el-form-item label="用户名" prop="db_user">
            <el-input v-model="dbForm.db_user" placeholder="postgres" />
          </el-form-item>
          <el-form-item label="密码" prop="db_password">
            <el-input v-model="dbForm.db_password" type="password" show-password />
          </el-form-item>
          <el-form-item label="数据库名" prop="db_name">
            <el-input v-model="dbForm.db_name" placeholder="anchorfinance" />
          </el-form-item>

          <el-divider content-position="left">Redis (可选)</el-divider>
          <el-form-item label="启用 Redis">
            <el-switch v-model="dbForm.enable_redis" />
          </el-form-item>
          <template v-if="dbForm.enable_redis">
            <el-form-item label="Redis 主机">
              <el-input v-model="dbForm.redis_host" placeholder="localhost" />
            </el-form-item>
            <el-form-item label="Redis 端口">
              <el-input v-model="dbForm.redis_port" placeholder="6379" />
            </el-form-item>
            <el-form-item label="Redis 密码">
              <el-input v-model="dbForm.redis_password" type="password" show-password />
            </el-form-item>
          </template>

          <el-form-item>
            <el-button type="primary" @click="testDB" :loading="isTesting">
              {{ isTesting ? '测试中...' : '测试连接' }}
            </el-button>
          </el-form-item>
        </el-form>

        <div class="step-buttons">
          <el-button @click="prevStep">上一步</el-button>
          <el-button v-if="dbTestPassed" type="success" @click="nextStep">下一步</el-button>
        </div>
      </div>

      <!-- Step 2: 管理员设置 -->
      <div v-if="currentStep === 2" class="step-content">
        <h3>管理员设置</h3>
        <el-form :model="adminForm" :rules="adminRules" ref="adminFormRef" label-width="120px">
          <el-form-item label="站点名称" prop="site_name">
            <el-input v-model="adminForm.site_name" placeholder="锚点财务" />
          </el-form-item>
          <el-form-item label="站点 URL" prop="site_url">
            <el-input v-model="adminForm.site_url" placeholder="http://localhost:8080" />
          </el-form-item>
          <el-divider content-position="left">管理员账号</el-divider>
          <el-form-item label="用户名" prop="admin_username">
            <el-input v-model="adminForm.admin_username" placeholder="admin" />
          </el-form-item>
          <el-form-item label="邮箱" prop="admin_email">
            <el-input v-model="adminForm.admin_email" placeholder="admin@example.com" />
          </el-form-item>
          <el-form-item label="密码" prop="admin_password">
            <el-input v-model="adminForm.admin_password" type="password" show-password placeholder="至少6位" />
          </el-form-item>
          <el-form-item label="确认密码" prop="admin_password_confirm">
            <el-input v-model="adminForm.admin_password_confirm" type="password" show-password />
          </el-form-item>
        </el-form>

        <div class="step-buttons">
          <el-button @click="prevStep">上一步</el-button>
          <el-button type="primary" @click="nextStep">下一步</el-button>
        </div>
      </div>

      <!-- Step 3: 完成安装 -->
      <div v-if="currentStep === 3" class="step-content">
        <h3>确认安装信息</h3>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="数据库主机">{{ dbForm.db_host }}:{{ dbForm.db_port }}</el-descriptions-item>
          <el-descriptions-item label="数据库名">{{ dbForm.db_name }}</el-descriptions-item>
          <el-descriptions-item label="Redis">{{ dbForm.enable_redis ? '已启用' : '未启用' }}</el-descriptions-item>
          <el-descriptions-item label="站点名称">{{ adminForm.site_name || '锚点财务' }}</el-descriptions-item>
          <el-descriptions-item label="管理员">{{ adminForm.admin_username }}</el-descriptions-item>
          <el-descriptions-item label="管理员邮箱">{{ adminForm.admin_email }}</el-descriptions-item>
        </el-descriptions>

        <div class="step-buttons">
          <el-button @click="prevStep">上一步</el-button>
          <el-button type="primary" @click="doInstall" :loading="isInstalling" :disabled="isInstalling">
            {{ isInstalling ? '安装中...' : '开始安装' }}
          </el-button>
        </div>

        <div v-if="installResult" class="install-result" :class="installResult.ok ? 'success' : 'error'">
          <el-icon v-if="installResult.ok"><CircleCheckFilled /></el-icon>
          <el-icon v-else><CircleCloseFilled /></el-icon>
          <span>{{ installResult.message }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Loading, CircleCheckFilled, CircleCloseFilled } from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()
const currentStep = ref(0)
const isChecking = ref(false)
const isTesting = ref(false)
const isInstalling = ref(false)
const dbTestPassed = ref(false)
const installResult = ref<{ ok: boolean; message: string } | null>(null)

// 环境检测
const envChecks = ref([
  { name: 'PostgreSQL 驱动', status: 'pending', message: '等待检测' },
  { name: 'Redis 客户端', status: 'pending', message: '等待检测' },
])

const allEnvChecksPass = computed(() => {
  return envChecks.value.every(c => c.status === 'success')
})

const checkEnvironment = async () => {
  isChecking.value = true
  for (const check of envChecks.value) {
    check.status = 'checking'
    check.message = '检测中...'
  }

  try {
    // Check PostgreSQL
    const pgRes = await request.get('/api/install/check-env', { params: { type: 'pg' } })
    if (pgRes.data.ok) {
      envChecks.value[0].status = 'success'
      envChecks.value[0].message = pgRes.data.message
    } else {
      envChecks.value[0].status = 'error'
      envChecks.value[0].message = pgRes.data.message
    }

    // Check Redis
    const redisRes = await request.get('/api/install/check-env', { params: { type: 'redis' } })
    if (redisRes.data.ok) {
      envChecks.value[1].status = 'success'
      envChecks.value[1].message = redisRes.data.message
    } else {
      envChecks.value[1].status = 'error'
      envChecks.value[1].message = redisRes.data.message
    }
  } catch (error: any) {
    ElMessage.error('环境检测失败: ' + (error.message || '网络错误'))
  } finally {
    isChecking.value = false
  }
}

// 数据库表单
const dbForm = ref({
  db_host: 'localhost',
  db_port: '5432',
  db_user: 'postgres',
  db_password: '',
  db_name: 'anchorfinance',
  enable_redis: false,
  redis_host: 'localhost',
  redis_port: '6379',
  redis_password: '',
})

const dbRules = {
  db_host: [{ required: true, message: '请输入数据库主机', trigger: 'blur' }],
  db_port: [{ required: true, message: '请输入数据库端口', trigger: 'blur' }],
  db_user: [{ required: true, message: '请输入数据库用户名', trigger: 'blur' }],
  db_name: [{ required: true, message: '请输入数据库名', trigger: 'blur' }],
}

const dbFormRef = ref()

const testDB = async () => {
  try {
    await dbFormRef.value?.validate()
  } catch {
    return
  }

  isTesting.value = true
  try {
    const res = await request.post('/api/install/test-db', dbForm.value)
    if (res.data.ok) {
      dbTestPassed.value = true
      ElMessage.success('数据库连接测试成功')
    } else {
      dbTestPassed.value = false
      ElMessage.error(res.data.message)
    }
  } catch (error: any) {
    ElMessage.error('测试失败: ' + (error.message || '网络错误'))
  } finally {
    isTesting.value = false
  }
}

// 管理员表单
const adminForm = ref({
  site_name: '锚点财务',
  site_url: 'http://localhost:8080',
  admin_username: 'admin',
  admin_email: '',
  admin_password: '',
  admin_password_confirm: '',
})

const adminRules = {
  admin_username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, message: '用户名至少3位', trigger: 'blur' },
  ],
  admin_email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email' as const, message: '请输入有效邮箱', trigger: 'blur' },
  ],
  admin_password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' },
  ],
  admin_password_confirm: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (value !== adminForm.value.admin_password) {
          callback(new Error('两次密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
}

const adminFormRef = ref()

// 步骤控制
const nextStep = async () => {
  if (currentStep.value === 2) {
    try {
      await adminFormRef.value?.validate()
    } catch {
      return
    }
  }
  currentStep.value++
}

const prevStep = () => {
  currentStep.value--
}

// 执行安装
const doInstall = async () => {
  isInstalling.value = true
  installResult.value = null

  try {
    const payload = {
      ...dbForm.value,
      admin_username: adminForm.value.admin_username,
      admin_password: adminForm.value.admin_password,
      admin_email: adminForm.value.admin_email,
      site_name: adminForm.value.site_name,
      site_url: adminForm.value.site_url,
    }

    const res = await request.post('/api/install', payload)
    installResult.value = res.data

    if (res.data.ok) {
      ElMessage.success('安装成功！')
      setTimeout(() => {
        router.push('/login')
      }, 2000)
    }
  } catch (error: any) {
    installResult.value = {
      ok: false,
      message: '安装失败: ' + (error.response?.data?.message || error.message || '网络错误'),
    }
  } finally {
    isInstalling.value = false
  }
}
</script>

<style scoped lang="scss">
.install-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1a237e 0%, #0d47a1 50%, #01579b 100%);
  padding: 20px;
}

.install-card {
  background: #fff;
  border-radius: 16px;
  padding: 40px;
  width: 100%;
  max-width: 700px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.install-header {
  text-align: center;
  margin-bottom: 32px;

  .logo {
    width: 80px;
    height: 80px;
    margin-bottom: 16px;
  }

  h1 {
    font-size: 24px;
    color: #1a237e;
    margin: 0 0 8px;
  }

  .subtitle {
    color: #666;
    font-size: 14px;
    margin: 0;
  }
}

.steps {
  margin-bottom: 32px;
}

.step-content {
  h3 {
    color: #1a237e;
    margin: 0 0 20px;
    font-size: 18px;
  }
}

.env-checks {
  margin-bottom: 20px;
}

.env-check-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f5f5f5;
  border-radius: 8px;
  margin-bottom: 8px;

  .check-info {
    display: flex;
    align-items: center;
    gap: 8px;

    .success {
      color: #67c23a;
    }

    .error {
      color: #f56c6c;
    }

    .check-name {
      font-weight: 500;
    }
  }

  .check-message {
    font-size: 13px;

    &.success {
      color: #67c23a;
    }

    &.error {
      color: #f56c6c;
    }

    &.pending {
      color: #909399;
    }
  }
}

.step-buttons {
  margin-top: 24px;
  display: flex;
  justify-content: space-between;
}

.install-result {
  margin-top: 24px;
  padding: 16px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 8px;

  &.success {
    background: #f0f9eb;
    color: #67c23a;
  }

  &.error {
    background: #fef0f0;
    color: #f56c6c;
  }
}
</style>
