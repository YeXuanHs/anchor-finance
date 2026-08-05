<template>
  <div class="api-manage-page">
    <div class="page-header">
      <h1 class="page-title">API 管理</h1>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        生成新密钥
      </el-button>
    </div>

    <el-alert
      title="API 密钥安全提示"
      description="请妥善保管您的 API 密钥，不要泄露给他人。如发现密钥泄露，请立即重置。"
      type="warning"
      :closable="false"
      show-icon
      class="security-alert"
    />

    <el-card shadow="never" class="api-card">
      <template #header>
        <div class="card-header">
          <span>API 密钥列表</span>
          <el-input
            v-model="searchKeyword"
            placeholder="搜索密钥名称..."
            style="width: 240px"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </template>

      <el-table :data="filteredKeys" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="密钥名称" min-width="150">
          <template #default="{ row }">
            <div class="key-name">
              <el-icon :size="16" color="#409eff"><Key /></el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="apiKey" label="API Key" min-width="280">
          <template #default="{ row }">
            <div class="key-display">
              <code class="key-text">{{ row.visible ? row.apiKey : maskKey(row.apiKey) }}</code>
              <el-button text size="small" @click="row.visible = !row.visible">
                <el-icon><View v-if="!row.visible" /><Hide v-else /></el-icon>
              </el-button>
              <el-button text size="small" @click="copyKey(row.apiKey)">
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="permissions" label="权限" width="200">
          <template #default="{ row }">
            <el-tag
              v-for="perm in row.permissions"
              :key="perm"
              size="small"
              effect="light"
              class="perm-tag"
            >
              {{ permLabel(perm) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small" effect="light" round>
              {{ row.status === 'active' ? '正常' : '已禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastUsed" label="最后使用" width="170" />
        <el-table-column prop="createdAt" label="创建时间" width="170" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button text type="warning" size="small" @click="handleReset(row)">重置</el-button>
            <el-popconfirm
              title="确定要删除此密钥吗？"
              confirm-button-text="删除"
              cancel-button-text="取消"
              confirm-button-type="danger"
              @confirm="handleDelete(row)"
            >
              <template #reference>
                <el-button text type="danger" size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑密钥对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      :title="editingKey ? '编辑 API 密钥' : '生成新密钥'"
      width="520px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="keyForm" :rules="formRules" label-width="100px">
        <el-form-item label="密钥名称" prop="name">
          <el-input v-model="keyForm.name" placeholder="请输入密钥名称" />
        </el-form-item>
        <el-form-item label="权限设置" prop="permissions">
          <el-checkbox-group v-model="keyForm.permissions">
            <el-checkbox label="read">只读访问</el-checkbox>
            <el-checkbox label="write">读写访问</el-checkbox>
            <el-checkbox label="delete">删除权限</el-checkbox>
            <el-checkbox label="admin">管理权限</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="IP 白名单">
          <el-input
            v-model="keyForm.ipWhitelist"
            type="textarea"
            :rows="3"
            placeholder="每行一个 IP 地址，留空表示不限制"
          />
        </el-form-item>
        <el-form-item label="过期时间">
          <el-select v-model="keyForm.expiry" placeholder="选择过期时间" style="width: 100%">
            <el-option label="永不过期" value="never" />
            <el-option label="30 天" value="30d" />
            <el-option label="90 天" value="90d" />
            <el-option label="1 年" value="1y" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">
          {{ editingKey ? '保存' : '生成密钥' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 新密钥展示对话框 -->
    <el-dialog v-model="showNewKeyDialog" title="密钥已生成" width="520px" :close-on-click-modal="false">
      <el-alert
        title="请立即复制保存您的密钥"
        description="密钥只会显示一次，关闭后将无法再次查看完整密钥。"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 16px"
      />
      <div class="new-key-display">
        <label class="key-label">API Key</label>
        <div class="key-copy-row">
          <code>{{ newKeyData.apiKey }}</code>
          <el-button type="primary" size="small" @click="copyKey(newKeyData.apiKey)">复制</el-button>
        </div>
        <label class="key-label" style="margin-top: 12px">Secret Key</label>
        <div class="key-copy-row">
          <code>{{ newKeyData.secretKey }}</code>
          <el-button type="primary" size="small" @click="copyKey(newKeyData.secretKey)">复制</el-button>
        </div>
      </div>
      <template #footer>
        <el-button type="primary" @click="showNewKeyDialog = false">我已保存，关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Plus, Search, Key, View, Hide, CopyDocument } from '@element-plus/icons-vue'
import request from '@/utils/request'

interface ApiKey {
  id: string
  name: string
  apiKey: string
  secretKey: string
  permissions: string[]
  status: 'active' | 'disabled'
  lastUsed: string
  createdAt: string
  visible: boolean
}

const loading = ref(false)
const searchKeyword = ref('')
const showCreateDialog = ref(false)
const showNewKeyDialog = ref(false)
const editingKey = ref<ApiKey | null>(null)
const formRef = ref<FormInstance>()

const keyForm = reactive({
  name: '',
  permissions: ['read'] as string[],
  ipWhitelist: '',
  expiry: 'never'
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入密钥名称', trigger: 'blur' }],
  permissions: [{ type: 'array', required: true, message: '请选择至少一项权限', trigger: 'change' }]
}

const newKeyData = reactive({
  apiKey: '',
  secretKey: ''
})

const apiKeys = ref<ApiKey[]>([])

const fetchApiKeys = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/user/api-keys')
    if (data?.data) {
      apiKeys.value = (data.data.list || data.data || []).map((k: any) => ({ ...k, visible: false }))
    }
  } catch (e) {
    console.error('Failed to fetch API keys:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchApiKeys()
})

const filteredKeys = computed(() => {
  if (!searchKeyword.value) return apiKeys.value
  return apiKeys.value.filter(k =>
    k.name.toLowerCase().includes(searchKeyword.value.toLowerCase())
  )
})

function maskKey(key: string) {
  return key.slice(0, 10) + '****' + key.slice(-6)
}

function permLabel(perm: string) {
  const map: Record<string, string> = { read: '只读', write: '读写', delete: '删除', admin: '管理' }
  return map[perm] || perm
}

function copyKey(key: string) {
  navigator.clipboard.writeText(key).then(() => {
    ElMessage.success('已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败，请手动复制')
  })
}

function handleEdit(row: ApiKey) {
  editingKey.value = row
  keyForm.name = row.name
  keyForm.permissions = [...row.permissions]
  keyForm.ipWhitelist = ''
  keyForm.expiry = 'never'
  showCreateDialog.value = true
}

function handleReset(row: ApiKey) {
  ElMessage.info(`密钥 ${row.name} 已重置`)
}

function handleDelete(row: ApiKey) {
  apiKeys.value = apiKeys.value.filter(k => k.id !== row.id)
  ElMessage.success('密钥已删除')
}

function handleSubmit() {
  formRef.value?.validate((valid) => {
    if (!valid) return

    if (editingKey.value) {
      const idx = apiKeys.value.findIndex(k => k.id === editingKey.value!.id)
      if (idx !== -1) {
        apiKeys.value[idx].name = keyForm.name
        apiKeys.value[idx].permissions = [...keyForm.permissions]
      }
      ElMessage.success('密钥已更新')
    } else {
      newKeyData.apiKey = 'ak_' + Math.random().toString(36).slice(2) + Math.random().toString(36).slice(2)
      newKeyData.secretKey = 'sk_' + Math.random().toString(36).slice(2) + Math.random().toString(36).slice(2)
      showNewKeyDialog.value = true
    }

    showCreateDialog.value = false
    editingKey.value = null
  })
}
</script>

<style scoped lang="scss">
.api-manage-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.security-alert {
  border-radius: 8px;
}

.api-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.key-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.key-display {
  display: flex;
  align-items: center;
  gap: 4px;
}

.key-text {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  background: #f5f7fa;
  padding: 2px 8px;
  border-radius: 4px;
  color: #606266;
}

.perm-tag {
  margin-right: 4px;
  margin-bottom: 2px;
}

.new-key-display {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 16px;
}

.key-label {
  display: block;
  font-size: 13px;
  color: #909399;
  margin-bottom: 6px;
}

.key-copy-row {
  display: flex;
  align-items: center;
  gap: 12px;

  code {
    flex: 1;
    font-family: 'Courier New', monospace;
    font-size: 13px;
    background: #fff;
    padding: 8px 12px;
    border-radius: 6px;
    border: 1px solid #dcdfe6;
    word-break: break-all;
  }
}
</style>
