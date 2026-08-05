<template>
  <div class="resource-pool-settings-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>资源池设置</span>
          <el-button type="primary" @click="handleSave" :loading="submitLoading">
            <el-icon><Check /></el-icon>
            保存配置
          </el-button>
        </div>
      </template>

      <el-form :model="settingsForm" :rules="settingsRules" ref="settingsFormRef" label-width="140px" v-loading="loading">
        <!-- 基础配置 -->
        <el-divider content-position="left">基础配置</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="资源池名称" prop="pool_name">
              <el-input v-model="settingsForm.pool_name" placeholder="请输入资源池名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="资源池容量" prop="capacity">
              <el-input-number v-model="settingsForm.capacity" :min="1" :max="10000" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="资源池状态" prop="status">
              <el-switch v-model="settingsForm.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="资源池类型" prop="pool_type">
              <el-select v-model="settingsForm.pool_type" placeholder="请选择资源池类型" style="width: 100%">
                <el-option label="共享型" value="shared" />
                <el-option label="独享型" value="exclusive" />
                <el-option label="混合型" value="hybrid" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="资源池描述" prop="description">
          <el-input v-model="settingsForm.description" type="textarea" :rows="3" placeholder="请输入资源池描述" />
        </el-form-item>

        <!-- 资源分配规则 -->
        <el-divider content-position="left">资源分配规则</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="分配策略" prop="allocation_strategy">
              <el-select v-model="settingsForm.allocation_strategy" placeholder="请选择分配策略" style="width: 100%">
                <el-option label="轮询分配" value="round_robin" />
                <el-option label="权重分配" value="weighted" />
                <el-option label="优先级分配" value="priority" />
                <el-option label="随机分配" value="random" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="最大分配数" prop="max_allocation">
              <el-input-number v-model="settingsForm.max_allocation" :min="1" :max="1000" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="分配阈值(%)" prop="allocation_threshold">
              <el-input-number v-model="settingsForm.allocation_threshold" :min="0" :max="100" style="width: 100%" />
              <div class="form-tip">资源使用率达到此阈值时停止自动分配</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分配冷却(秒)" prop="cooldown_seconds">
              <el-input-number v-model="settingsForm.cooldown_seconds" :min="0" :max="3600" style="width: 100%" />
              <div class="form-tip">两次分配之间的最小间隔时间</div>
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 自动分配开关 -->
        <el-divider content-position="left">自动分配配置</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="自动分配" prop="auto_assign">
              <el-switch v-model="settingsForm.auto_assign" active-text="开启" inactive-text="关闭" />
              <div class="form-tip">开启后系统将按分配策略自动分配资源</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="自动回收" prop="auto_reclaim">
              <el-switch v-model="settingsForm.auto_reclaim" active-text="开启" inactive-text="关闭" />
              <div class="form-tip">开启后空闲资源将自动回收到资源池</div>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="空闲超时(分钟)" prop="idle_timeout">
              <el-input-number v-model="settingsForm.idle_timeout" :min="1" :max="1440" style="width: 100%" />
              <div class="form-tip">资源空闲超过此时间将被自动回收</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通知阈值(%)" prop="notify_threshold">
              <el-input-number v-model="settingsForm.notify_threshold" :min="0" :max="100" style="width: 100%" />
              <div class="form-tip">资源使用率达到此阈值时发送告警通知</div>
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 通知配置 -->
        <el-divider content-position="left">通知配置</el-divider>
        <el-form-item label="通知方式" prop="notify_methods">
          <el-checkbox-group v-model="settingsForm.notify_methods">
            <el-checkbox label="email">邮件通知</el-checkbox>
            <el-checkbox label="sms">短信通知</el-checkbox>
            <el-checkbox label="webhook">Webhook通知</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="Webhook地址" prop="webhook_url" v-if="settingsForm.notify_methods.includes('webhook')">
          <el-input v-model="settingsForm.webhook_url" placeholder="请输入Webhook回调地址" />
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Check } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'ResourcePoolSettings' })

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)

// 表单引用
const settingsFormRef = ref<FormInstance>()

// 表单数据
const settingsForm = reactive({
  pool_name: '',
  capacity: 100,
  status: 1,
  pool_type: 'shared',
  description: '',
  allocation_strategy: 'round_robin',
  max_allocation: 50,
  allocation_threshold: 80,
  cooldown_seconds: 60,
  auto_assign: true,
  auto_reclaim: true,
  idle_timeout: 30,
  notify_threshold: 90,
  notify_methods: ['email'] as string[],
  webhook_url: ''
})

// 表单验证规则
const settingsRules: FormRules = {
  pool_name: [
    { required: true, message: '请输入资源池名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  capacity: [
    { required: true, message: '请输入资源池容量', trigger: 'blur' }
  ],
  pool_type: [
    { required: true, message: '请选择资源池类型', trigger: 'change' }
  ],
  allocation_strategy: [
    { required: true, message: '请选择分配策略', trigger: 'change' }
  ]
}

// 获取资源池设置
const fetchSettings = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/resource-pool/settings'
    })
    Object.assign(settingsForm, data)
  } catch (error) {
    console.error('获取资源池设置失败:', error)
    ElMessage.error('获取资源池设置失败')
  } finally {
    loading.value = false
  }
}

// 保存配置
const handleSave = async () => {
  if (!settingsFormRef.value) return

  await settingsFormRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      await request.put({
        url: '/api/admin/resource-pool/settings',
        params: settingsForm,
        showSuccessMessage: true
      })
      ElMessage.success('保存成功')
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      submitLoading.value = false
    }
  })
}

onMounted(() => {
  fetchSettings()
})
</script>

<style scoped lang="scss">
.resource-pool-settings-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

:deep(.el-divider__text) {
  font-weight: 600;
  color: #303133;
}
</style>
