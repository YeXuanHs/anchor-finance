<template>
  <div class="admin-page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">第三方登录管理</span>
          <el-button type="primary" @click="openAddModal">
            <el-icon><Plus /></el-icon>添加提供商
          </el-button>
        </div>
      </template>

      <el-table :data="providerList" stripe size="small">
        <el-table-column label="提供商" width="140">
          <template #default="{ row }">
            <div style="display: flex; align-items: center; gap: 8px">
              <div :style="{ width: '10px', height: '10px', borderRadius: '50%', background: row.icon }"></div>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="appId" label="App ID" show-overflow-tooltip />
        <el-table-column prop="callbackUrl" label="回调URL" show-overflow-tooltip />
        <el-table-column prop="scope" label="Scope" width="160" />
        <el-table-column prop="enabled" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">{{ row.enabled ? '已启用' : '已禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" :type="row.enabled ? 'warning' : 'success'" @click="handleToggle(row)">
              {{ row.enabled ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" type="info" @click="handleTest(row)">测试</el-button>
            <el-popconfirm title="确定删除此提供商？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Config Dialog -->
    <el-dialog v-model="showModal" :title="editingId ? '编辑提供商' : '添加提供商'" width="560px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="提供商" prop="provider">
          <el-select v-model="formData.provider" placeholder="选择提供商" :disabled="!!editingId" style="width: 100%">
            <el-option v-for="o in providerOptions" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="App ID" prop="appId">
          <el-input v-model="formData.appId" placeholder="输入App ID / Client ID" />
        </el-form-item>
        <el-form-item label="App Secret" prop="appSecret">
          <el-input v-model="formData.appSecret" type="password" show-password placeholder="输入App Secret / Client Secret" />
        </el-form-item>
        <el-form-item label="回调URL" prop="callbackUrl">
          <el-input v-model="formData.callbackUrl" placeholder="https://example.com/oauth/callback" />
        </el-form-item>
        <el-form-item label="Scope">
          <el-input v-model="formData.scope" placeholder="user_info, email" />
        </el-form-item>
        <el-form-item label="启用状态">
          <el-switch v-model="formData.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showModal = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>

    <!-- Test Result Dialog -->
    <el-dialog v-model="showTestModal" title="连接测试结果" width="420px">
      <el-result :status="testResult.success ? 'success' : 'error'" :title="testResult.success ? '连接成功' : '连接失败'" :sub-title="testResult.message" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const formRef = ref<FormInstance>()
const showModal = ref(false)
const showTestModal = ref(false)
const editingId = ref<string | null>(null)

interface OAuthProvider {
  id: string; provider: string; name: string; appId: string; appSecret: string;
  callbackUrl: string; scope: string; enabled: boolean; icon: string;
}

const providerOptions = [
  { label: '微信', value: 'wechat' }, { label: 'QQ', value: 'qq' },
  { label: 'GitHub', value: 'github' }, { label: 'Google', value: 'google' },
]
const providerNameMap: Record<string, string> = { wechat: '微信', qq: 'QQ', github: 'GitHub', google: 'Google' }
const providerIconMap: Record<string, string> = { wechat: '#1AAD19', qq: '#12B7F5', github: '#333', google: '#EA4335' }

const providerList = ref<OAuthProvider[]>([
  { id: '1', provider: 'wechat', name: '微信', appId: 'wx1234567890abcdef', appSecret: '••••••••••••••••', callbackUrl: 'https://anchorfinance.com/oauth/wechat/callback', scope: 'snsapi_login', enabled: true, icon: '#1AAD19' },
  { id: '2', provider: 'qq', name: 'QQ', appId: '101234567', appSecret: '••••••••••••••••', callbackUrl: 'https://anchorfinance.com/oauth/qq/callback', scope: 'get_user_info', enabled: true, icon: '#12B7F5' },
  { id: '3', provider: 'github', name: 'GitHub', appId: 'Iv1.abc123def456', appSecret: '••••••••••••••••', callbackUrl: 'https://anchorfinance.com/oauth/github/callback', scope: 'user:email', enabled: false, icon: '#333' },
  { id: '4', provider: 'google', name: 'Google', appId: '123456789-abc.apps.googleusercontent.com', appSecret: '••••••••••••••••', callbackUrl: 'https://anchorfinance.com/oauth/google/callback', scope: 'openid email profile', enabled: false, icon: '#EA4335' },
])

const formData = reactive({ provider: '', appId: '', appSecret: '', callbackUrl: '', scope: '', enabled: true })
const formRules: FormRules = {
  provider: { required: true, message: '请选择提供商', trigger: 'change' },
  appId: { required: true, message: '请输入App ID', trigger: 'blur' },
  appSecret: { required: true, message: '请输入App Secret', trigger: 'blur' },
  callbackUrl: { required: true, message: '请输入回调URL', trigger: 'blur' },
}

const testResult = reactive({ success: true, message: '' })

function openAddModal() {
  editingId.value = null
  Object.assign(formData, { provider: '', appId: '', appSecret: '', callbackUrl: '', scope: '', enabled: true })
  showModal.value = true
}

function handleEdit(row: OAuthProvider) {
  editingId.value = row.id
  Object.assign(formData, { provider: row.provider, appId: row.appId, appSecret: row.appSecret, callbackUrl: row.callbackUrl, scope: row.scope, enabled: row.enabled })
  showModal.value = true
}

function handleSave() {
  formRef.value?.validate((valid) => {
    if (valid) {
      if (editingId.value) {
        const idx = providerList.value.findIndex((p) => p.id === editingId.value)
        if (idx !== -1) {
          providerList.value[idx] = { ...providerList.value[idx], ...formData, name: providerNameMap[formData.provider] || formData.provider, icon: providerIconMap[formData.provider] || '#999' }
        }
        ElMessage.success('提供商已更新')
      } else {
        providerList.value.push({ id: String(Date.now()), ...formData, name: providerNameMap[formData.provider] || formData.provider, icon: providerIconMap[formData.provider] || '#999' })
        ElMessage.success('提供商已添加')
      }
      showModal.value = false
    }
  })
}

function handleToggle(row: OAuthProvider) {
  row.enabled = !row.enabled
  ElMessage.success(`已${row.enabled ? '启用' : '禁用'} ${row.name}`)
}

function handleTest(row: OAuthProvider) {
  testResult.success = row.enabled
  testResult.message = row.enabled ? `已成功连接到 ${row.name} OAuth 服务，返回状态码 200` : `${row.name} 未启用，请先启用后再测试`
  showTestModal.value = true
}

function handleDelete(id: string) {
  providerList.value = providerList.value.filter((p) => p.id !== id)
  ElMessage.success('提供商已删除')
}
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; }
.card-title { font-size: 16px; font-weight: 600; }
</style>
