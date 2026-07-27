<template>
  <n-card :bordered="false" rounded>
    <n-space justify="space-between" align="center" style="margin-bottom: 16px">
      <n-h3 prefix="bar" style="margin: 0">第三方登录管理</n-h3>
      <n-button type="primary" @click="openAddModal">
        <template #icon><n-icon><AddIcon /></n-icon></template>
        添加提供商
      </n-button>
    </n-space>

    <n-data-table
      :columns="columns"
      :data="providerList"
      :bordered="false"
      :single-line="false"
      striped
    />
  </n-card>

  <!-- 配置对话框 -->
  <n-modal
    v-model:show="showModal"
    preset="card"
    :title="isEditing ? '编辑 OAuth 提供商' : '添加 OAuth 提供商'"
    style="width: 560px"
    :bordered="false"
    :segmented="{ content: true, footer: true }"
  >
    <n-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-placement="left"
      label-width="100"
    >
      <n-form-item label="提供商名称" path="name">
        <n-select
          v-if="!isEditing"
          v-model:value="formData.name"
          :options="providerOptions"
          placeholder="选择提供商"
          @update:value="handleProviderChange"
        />
        <n-input v-else :value="formData.name" disabled />
      </n-form-item>
      <n-form-item label="显示名称" path="displayName">
        <n-input v-model:value="formData.displayName" placeholder="如：微信登录" />
      </n-form-item>
      <n-form-item label="App ID" path="appId">
        <n-input v-model:value="formData.appId" placeholder="第三方应用 App ID" />
      </n-form-item>
      <n-form-item label="App Secret" path="appSecret">
        <n-input
          v-model:value="formData.appSecret"
          type="password"
          show-password-on="click"
          placeholder="第三方应用 App Secret"
        />
      </n-form-item>
      <n-form-item label="回调 URL" path="callbackUrl">
        <n-input v-model:value="formData.callbackUrl" placeholder="https://example.com/oauth/callback" />
      </n-form-item>
      <n-form-item label="Scope" path="scope">
        <n-input v-model:value="formData.scope" placeholder="如：openid profile email" />
      </n-form-item>
      <n-form-item label="启用状态">
        <n-switch v-model:value="formData.enabled" />
      </n-form-item>
    </n-form>
    <template #footer>
      <n-space justify="end">
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" :loading="saving" @click="handleSave">保存</n-button>
      </n-space>
    </template>
  </n-modal>

  <!-- 测试连接结果 -->
  <n-modal
    v-model:show="showTestResult"
    preset="card"
    title="连接测试结果"
    style="width: 400px"
    :bordered="false"
  >
    <n-result
      :status="testSuccess ? 'success' : 'error'"
      :title="testSuccess ? '连接成功' : '连接失败'"
      :description="testResultMessage"
    />
  </n-modal>
</template>

<script setup lang="ts">
import { ref, reactive, h } from 'vue'
import { useMessage, NTag, NSwitch, NSpace, NButton, NPopconfirm, NIcon } from 'naive-ui'
import { AddOutline as AddIcon } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'

const message = useMessage()

// ---- Provider Options ----
const providerOptions = [
  { label: '微信 (WeChat)', value: 'wechat' },
  { label: 'QQ', value: 'qq' },
  { label: 'GitHub', value: 'github' },
  { label: 'Google', value: 'google' },
]

const providerDefaults: Record<string, { displayName: string; scope: string; callbackSuffix: string }> = {
  wechat: { displayName: '微信登录', scope: 'snsapi_login', callbackSuffix: '/oauth/wechat/callback' },
  qq: { displayName: 'QQ 登录', scope: 'get_user_info', callbackSuffix: '/oauth/qq/callback' },
  github: { displayName: 'GitHub 登录', scope: 'read:user user:email', callbackSuffix: '/oauth/github/callback' },
  google: { displayName: 'Google 登录', scope: 'openid profile email', callbackSuffix: '/oauth/google/callback' },
}

// ---- Types ----
interface OAuthProvider {
  id: number
  name: string
  displayName: string
  appId: string
  appSecret: string
  callbackUrl: string
  scope: string
  enabled: boolean
  status: 'connected' | 'disconnected' | 'error'
  lastTested: string
}

// ---- Mock Data ----
const providerList = ref<OAuthProvider[]>([
  {
    id: 1,
    name: 'wechat',
    displayName: '微信登录',
    appId: 'wx1234567890abcdef',
    appSecret: '••••••••••••••••',
    callbackUrl: 'https://anchorfinance.com/oauth/wechat/callback',
    scope: 'snsapi_login',
    enabled: true,
    status: 'connected',
    lastTested: '2026-07-26 14:30:00',
  },
  {
    id: 2,
    name: 'qq',
    displayName: 'QQ 登录',
    appId: '101234567',
    appSecret: '••••••••••••••••',
    callbackUrl: 'https://anchorfinance.com/oauth/qq/callback',
    scope: 'get_user_info',
    enabled: true,
    status: 'connected',
    lastTested: '2026-07-25 10:15:00',
  },
  {
    id: 3,
    name: 'github',
    displayName: 'GitHub 登录',
    appId: 'Iv1.aabbccddeeff0011',
    appSecret: '••••••••••••••••',
    callbackUrl: 'https://anchorfinance.com/oauth/github/callback',
    scope: 'read:user user:email',
    enabled: false,
    status: 'disconnected',
    lastTested: '—',
  },
  {
    id: 4,
    name: 'google',
    displayName: 'Google 登录',
    appId: '123456789-abcdef.apps.googleusercontent.com',
    appSecret: '••••••••••••••••',
    callbackUrl: 'https://anchorfinance.com/oauth/google/callback',
    scope: 'openid profile email',
    enabled: false,
    status: 'error',
    lastTested: '2026-07-20 09:00:00',
  },
])

// ---- Table Columns ----
const columns: DataTableColumns<OAuthProvider> = [
  {
    title: '提供商',
    key: 'displayName',
    width: 140,
    render(row) {
      const colorMap: Record<string, string> = {
        wechat: '#07C160',
        qq: '#12B7F5',
        github: '#24292E',
        google: '#4285F4',
      }
      return h('span', { style: { fontWeight: 600, color: colorMap[row.name] || '#333' } }, row.displayName)
    },
  },
  { title: 'App ID', key: 'appId', ellipsis: { tooltip: true } },
  {
    title: '回调 URL',
    key: 'callbackUrl',
    ellipsis: { tooltip: true },
    width: 280,
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render(row) {
      const statusMap: Record<string, { label: string; type: 'success' | 'warning' | 'error' }> = {
        connected: { label: '已连接', type: 'success' },
        disconnected: { label: '未启用', type: 'warning' },
        error: { label: '异常', type: 'error' },
      }
      const s = statusMap[row.status] || statusMap.disconnected
      return h(NTag, { type: s.type, size: 'small', round: true }, { default: () => s.label })
    },
  },
  {
    title: '启用',
    key: 'enabled',
    width: 80,
    render(row) {
      return h(NSwitch, {
        value: row.enabled,
        size: 'small',
        onUpdateValue(val: boolean) {
          row.enabled = val
          message.success(`${row.displayName} 已${val ? '启用' : '禁用'}`)
        },
      })
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, { size: 'small', quaternary: true, type: 'info', onClick: () => openEditModal(row) }, { default: () => '编辑' }),
          h(NButton, { size: 'small', quaternary: true, type: 'primary', onClick: () => testConnection(row) }, { default: () => '测试连接' }),
          h(
            NPopconfirm,
            { onPositiveClick: () => deleteProvider(row.id) },
            {
              trigger: () => h(NButton, { size: 'small', quaternary: true, type: 'error' }, { default: () => '删除' }),
              default: () => `确定删除 ${row.displayName}？`,
            }
          ),
        ],
      })
    },
  },
]

// ---- Modal State ----
const showModal = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)
const formRef = ref<FormInst | null>(null)

const formData = reactive({
  name: '',
  displayName: '',
  appId: '',
  appSecret: '',
  callbackUrl: '',
  scope: '',
  enabled: true,
})

const formRules: FormRules = {
  name: { required: true, message: '请选择提供商', trigger: 'change' },
  displayName: { required: true, message: '请输入显示名称', trigger: 'blur' },
  appId: { required: true, message: '请输入 App ID', trigger: 'blur' },
  appSecret: { required: true, message: '请输入 App Secret', trigger: 'blur' },
}

// ---- Test Result ----
const showTestResult = ref(false)
const testSuccess = ref(false)
const testResultMessage = ref('')

// ---- Handlers ----
function openAddModal() {
  isEditing.value = false
  editingId.value = null
  Object.assign(formData, { name: '', displayName: '', appId: '', appSecret: '', callbackUrl: '', scope: '', enabled: true })
  showModal.value = true
}

function openEditModal(row: OAuthProvider) {
  isEditing.value = true
  editingId.value = row.id
  Object.assign(formData, {
    name: row.name,
    displayName: row.displayName,
    appId: row.appId,
    appSecret: '',
    callbackUrl: row.callbackUrl,
    scope: row.scope,
    enabled: row.enabled,
  })
  showModal.value = true
}

function handleProviderChange(val: string) {
  const defaults = providerDefaults[val]
  if (defaults) {
    formData.displayName = defaults.displayName
    formData.scope = defaults.scope
    formData.callbackUrl = `https://anchorfinance.com${defaults.callbackSuffix}`
  }
}

async function handleSave() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  saving.value = true
  setTimeout(() => {
    if (isEditing.value && editingId.value !== null) {
      const idx = providerList.value.findIndex((p) => p.id === editingId.value)
      if (idx !== -1) {
        const existing = providerList.value[idx]
        Object.assign(existing, {
          displayName: formData.displayName,
          appId: formData.appId,
          appSecret: formData.appSecret ? '••••••••••••••••' : existing.appSecret,
          callbackUrl: formData.callbackUrl,
          scope: formData.scope,
          enabled: formData.enabled,
        })
      }
      message.success('提供商配置已更新')
    } else {
      const newId = Math.max(0, ...providerList.value.map((p) => p.id)) + 1
      providerList.value.push({
        id: newId,
        name: formData.name,
        displayName: formData.displayName,
        appId: formData.appId,
        appSecret: '••••••••••••••••',
        callbackUrl: formData.callbackUrl,
        scope: formData.scope,
        enabled: formData.enabled,
        status: 'disconnected',
        lastTested: '—',
      })
      message.success('提供商已添加')
    }
    saving.value = false
    showModal.value = false
  }, 500)
}

function testConnection(row: OAuthProvider) {
  message.loading('正在测试连接...')
  setTimeout(() => {
    const success = row.enabled && row.appId.length > 5
    testSuccess.value = success
    testResultMessage.value = success
      ? `与 ${row.displayName} 的连接测试通过，OAuth 配置有效。`
      : `无法连接到 ${row.displayName}，请检查 App ID 和 App Secret 是否正确。`
    showTestResult.value = true
    row.status = success ? 'connected' : 'error'
    row.lastTested = new Date().toLocaleString('zh-CN')
  }, 1200)
}

function deleteProvider(id: number) {
  providerList.value = providerList.value.filter((p) => p.id !== id)
  message.success('提供商已删除')
}
</script>

<style scoped>
.n-card {
  border-radius: 12px;
}
</style>
