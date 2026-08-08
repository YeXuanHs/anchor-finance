<template>
  <div class="third-login-page">
    <div class="login-grid" v-loading="loading">
      <el-card v-for="item in loginMethods" :key="item.id" shadow="hover" class="login-card" :class="{ 'is-enabled': item.enabled }">
        <div class="login-header">
          <div class="login-icon">
            <el-icon :size="32" :color="item.enabled ? 'var(--el-color-primary)' : '#86909C'"><Link /></el-icon>
          </div>
          <el-switch v-model="item.enabled" @change="handleToggle(item)" />
        </div>
        <div class="login-body">
          <h3>{{ item.name }}</h3>
          <p>{{ item.description }}</p>
        </div>
        <div class="login-footer">
          <el-button type="primary" size="small" @click="handleConfigure(item)">配置</el-button>
        </div>
      </el-card>
    </div>

    <el-dialog v-model="configDialogVisible" :title="`${currentItem?.name} - 配置`" width="600px">
      <el-form ref="configFormRef" :model="configFormData" label-width="120px">
        <el-form-item label="App ID" prop="app_id">
          <el-input v-model="configFormData.app_id" placeholder="请输入App ID" />
        </el-form-item>
        <el-form-item label="App Secret" prop="app_secret">
          <el-input v-model="configFormData.app_secret" type="password" placeholder="请输入App Secret" show-password />
        </el-form-item>
        <el-form-item label="回调地址" prop="callback_url">
          <el-input v-model="configFormData.callback_url" placeholder="回调地址" disabled />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="configDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="savingConfig" @click="handleSaveConfig">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Link } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const loginMethods = ref<any[]>([])
const configDialogVisible = ref(false)
const currentItem = ref<any>(null)
const configFormData = ref<any>({})
const configFormRef = ref<FormInstance>()
const savingConfig = ref(false)

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/third-login' })
    loginMethods.value = data || []
  } catch (error) {
    console.error('获取第三方登录失败:', error)
  } finally {
    loading.value = false
  }
}

const handleToggle = async (item: any) => {
  try {
    await request.post({ url: `/api/admin/third-login/${item.id}/toggle`, data: { enabled: item.enabled } })
    ElMessage.success(item.enabled ? '已启用' : '已禁用')
  } catch (error) {
    console.error('切换状态失败:', error)
    item.enabled = !item.enabled
  }
}

const handleConfigure = async (item: any) => {
  currentItem.value = item
  try {
    const data = await request.get({ url: `/api/admin/third-login/${item.id}/config` })
    configFormData.value = data || {}
  } catch (error) {
    configFormData.value = {}
  }
  configDialogVisible.value = true
}

const handleSaveConfig = async () => {
  if (!configFormRef.value) return
  try {
    await configFormRef.value.validate()
    savingConfig.value = true
    await request.post({ url: `/api/admin/third-login/${currentItem.value.id}/config`, data: configFormData.value })
    ElMessage.success('配置保存成功')
    configDialogVisible.value = false
  } catch (error) {
    console.error('保存配置失败:', error)
  } finally {
    savingConfig.value = false
  }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.third-login-page { padding: 16px; }
.login-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 16px; }
.login-card { transition: all 0.3s; &:hover { transform: translateY(-4px); } &.is-enabled { border-color: var(--el-color-primary-light-5); } }
.login-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 12px; }
.login-icon { width: 48px; height: 48px; border-radius: 8px; background: #F2F3F5; display: flex; align-items: center; justify-content: center; }
.login-body { margin-bottom: 12px; h3 { margin: 0 0 8px; font-size: 16px; } p { margin: 0; font-size: 13px; color: #86909C; } }
</style>
