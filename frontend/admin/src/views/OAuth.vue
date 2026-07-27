<template>
  <n-card :bordered="false" rounded>
    <n-space justify="space-between" align="center" style="margin-bottom: 16px">
      <n-text strong style="font-size: 18px">第三方登录管理</n-text>
      <n-button type="primary" @click="openAddModal">
        <template #icon><n-icon><AddOutline /></n-icon></template>
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
    :title="editingId ? '编辑提供商' : '添加提供商'"
    preset="card"
    style="max-width: 560px"
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
      <n-form-item label="提供商" path="provider">
        <n-select
          v-model:value="formData.provider"
          :options="providerOptions"
          placeholder="选择提供商"
          :disabled="!!editingId"
        />
      </n-form-item>
      <n-form-item label="App ID" path="appId">
        <n-input v-model:value="formData.appId" placeholder="输入App ID / Client ID" />
      </n-form-item>
      <n-form-item label="App Secret" path="appSecret">
        <n-input
          v-model:value="formData.appSecret"
          type="password"
          show-password-on="click"
          placeholder="输入App Secret / Client Secret"
        />
      </n-form-item>
      <n-form-item label="回调URL" path="callbackUrl">
        <n-input v-model:value="formData.callbackUrl" placeholder="https://example.com/oauth/callback" />
      </n-form-item>
      <n-form-item label="Scope" path="scope">
        <n-input v-model:value="formData.scope" placeholder="user_info, email" />
      </n-form-item>
      <n-form-item label="启用状态">
        <n-switch v-model:value="formData.enabled" />
      </n-form-item>
    </n-form>
    <template #footer>
      <n-space justify="end">
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="handleSave">保存</n-button>
      </n-space>
    </template>
  </n-modal>

  <!-- 测试结果对话框 -->
  <n-modal
    v-model:show="showTestModal"
    title="连接测试结果"
    preset="card"
    style="max-width: 420px"
    :bordered="false"
  >
    <n-result
      :status="testResult.success ? 'success' : 'error'"
      :title="testResult.success ? '连接成功' : '连接失败'"
      :description="testResult.message"
    />
  </n-modal>
</template>

<script setup lang="ts">
import { ref, reactive, h } from 'vue'
import { useMessage, NTag, NSwitch, NButton, NSpace, NPopconfirm } from 'naive-ui'
import { AddOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'

const message = useMessage()
const formRef = ref<FormInst | null>(null)
const showModal = ref(false)
const showTestModal = ref(false)
const editingId = ref<string | null>(null)

interface OAuthProvider {
  id: string
  provider: string
  name: string
  appId: string
  appSecret: string
  callbackUrl: string
  scope: string
  enabled: boolean
  icon: string
}

const providerOptions = [
  { label: '微信', value: 'wechat' },
  { label: 'QQ', value: 'qq' },
  { label: 'GitHub', value: 'github' },
  { label: 'Google', value: 'google' },
]

const providerNameMap: Record<string, string> = {
  wechat: '微信',
  qq: 'QQ',
  github: 'GitHub',
  google: 'Google',
}

const providerIconMap: Record<string, string> = {
  wechat: '#1AAD19',
  qq: '#12B7F5',
  github: '#333',
  google: '#EA4335',
}

const providerList = ref<OAuthProvider[]>([
  {
    id: '1',
    provider: 'wechat',
    name: '微信',
    appId: 'wx1234567890abcdef',
    appSecret: '••••••••••••••••',
    callbackUrl: 'https://anchorfinance.com/oauth/wechat/callback',
    scope: 'snsapi_login',
    enabled: true,
    icon: '#1AAD19',
  },
  {
    id: '2',
    provider: 'qq',
    name: 'QQ',
    appId: '101234567',
    appSecret: '••••••••••••••••',
    callbackUrl: 'https://anchorfinance.com/oauth/qq/callback',
    scope: 'get_user_info',
    enabled: true,
    icon: '#12B7F5',
  },
  {
    id: '3',
    provider: 'github',
    name: 'GitHub',
    appId: 'Iv1.abc123def456',
    appSecret: '••••••••••••••••',
    callbackUrl: 'https://anchorfinance.com/oauth/github/callback',
    scope: 'user:email',
    enabled: false,
    icon: '#333',
  },
  {
    id: '4',
    provider: 'google',
    name: 'Google',
    appId: '123456789-abc.apps.googleusercontent.com',
    appSecret: '••••••••••••••••',
    callbackUrl: 'https://anchorfinance.com/oauth/google/callback',
    scope: 'openid email profile',
    enabled: false,
    icon: '#EA4335',
  },
])

const formData = reactive({
  provider: '',
  appId: '',
  appSecret: '',
  callbackUrl: '',
  scope: '',
  enabled: true,
})

const formRules: FormRules = {
  provider: { required: true, message: '请选择提供商', trigger: 'change' },
  appId: { required: true, message: '请输入App ID', trigger: 'blur' },
  appSecret: { required: true, message: '请输入App Secret', trigger: 'blur' },
  callbackUrl: { required: true, message: '请输入回调URL', trigger: 'blur' },
}

const testResult = reactive({
  success: true,
  message: '',
})

const columns: DataTableColumns<OAuthProvider> = [
  {
    title: '提供商',
    key: 'name',
    width: 140,
    render(row) {
      return h('div', { style: 'display:flex;align-items:center;gap:8px' }, [
        h('div', {
          style: `width:10px;height:10px;border-radius:50%;background:${row.icon}`,
        }),
        h('span', null, { default: () => row.name }),
      ])
    },
  },
  {
    title: 'App ID',
    key: 'appId',
    ellipsis: { tooltip: true },
  },
  {
    title: '回调URL',
    key: 'callbackUrl',
    ellipsis: { tooltip: true },
  },
  {
    title: 'Scope',
    key: 'scope',
    width: 160,
  },
  {
    title: '状态',
    key: 'enabled',
    width: 100,
    render(row) {
      return h(NTag, {
        type: row.enabled ? 'success' : 'default',
        size: 'small',
        bordered: false,
      }, { default: () => row.enabled ? '已启用' : '已禁用' })
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 240,
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, {
            size: 'small',
            onClick: () => handleEdit(row),
          }, { default: () => '编辑' }),
          h(NButton, {
            size: 'small',
            type: row.enabled ? 'warning' : 'success',
            onClick: () => handleToggle(row),
          }, { default: () => row.enabled ? '禁用' : '启用' }),
          h(NButton, {
            size: 'small',
            type: 'info',
            onClick: () => handleTest(row),
          }, { default: () => '测试' }),
          h(NPopconfirm, {
            onPositiveClick: () => handleDelete(row.id),
          }, {
            trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '删除' }),
            default: () => '确定删除此提供商？',
          }),
        ],
      })
    },
  },
]

function openAddModal() {
  editingId.value = null
  formData.provider = ''
  formData.appId = ''
  formData.appSecret = ''
  formData.callbackUrl = ''
  formData.scope = ''
  formData.enabled = true
  showModal.value = true
}

function handleEdit(row: OAuthProvider) {
  editingId.value = row.id
  formData.provider = row.provider
  formData.appId = row.appId
  formData.appSecret = row.appSecret
  formData.callbackUrl = row.callbackUrl
  formData.scope = row.scope
  formData.enabled = row.enabled
  showModal.value = true
}

function handleSave() {
  formRef.value?.validate((errors) => {
    if (!errors) {
      if (editingId.value) {
        const idx = providerList.value.findIndex((p) => p.id === editingId.value)
        if (idx !== -1) {
          providerList.value[idx] = {
            ...providerList.value[idx],
            ...formData,
            name: providerNameMap[formData.provider] || formData.provider,
            icon: providerIconMap[formData.provider] || '#999',
          }
        }
        message.success('提供商已更新')
      } else {
        providerList.value.push({
          id: String(Date.now()),
          ...formData,
          name: providerNameMap[formData.provider] || formData.provider,
          icon: providerIconMap[formData.provider] || '#999',
        })
        message.success('提供商已添加')
      }
      showModal.value = false
    }
  })
}

function handleToggle(row: OAuthProvider) {
  row.enabled = !row.enabled
  message.success(`已${row.enabled ? '启用' : '禁用'} ${row.name}`)
}

function handleTest(row: OAuthProvider) {
  // TODO: real API test
  testResult.success = row.enabled
  testResult.message = row.enabled
    ? `已成功连接到 ${row.name} OAuth 服务，返回状态码 200`
    : `${row.name} 未启用，请先启用后再测试`
  showTestModal.value = true
}

function handleDelete(id: string) {
  providerList.value = providerList.value.filter((p) => p.id !== id)
  message.success('提供商已删除')
}
</script>

<style scoped>
.n-card {
  border-radius: 12px;
}
</style>
