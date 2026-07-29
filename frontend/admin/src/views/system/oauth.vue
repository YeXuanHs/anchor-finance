<template>
  <div class="oauth-page page-container">
    <div class="art-card">
      <div class="table-header">
        <h3>OAuth配置</h3>
      </div>

      <el-table :data="providers" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="提供商" width="160">
          <template #default="{ row }">
            <div class="provider-name">
              <el-icon :size="20" :style="{ color: row.color }">
                <component :is="row.icon" />
              </el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="app_id" label="App ID" show-overflow-tooltip>
          <template #default="{ row }">
            <span v-if="row.app_id">{{ row.app_id }}</span>
            <span v-else class="text-muted">未配置</span>
          </template>
        </el-table-column>
        <el-table-column prop="callback_url" label="回调URL" show-overflow-tooltip>
          <template #default="{ row }">
            <span v-if="row.callback_url">{{ row.callback_url }}</span>
            <span v-else class="text-muted">未配置</span>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small" effect="light">
              {{ row.enabled ? '已启用' : '已禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button
              :type="row.enabled ? 'danger' : 'success'"
              link
              @click="handleToggle(row)"
            >
              {{ row.enabled ? '禁用' : '启用' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="`${editingProvider?.name} 配置`"
      width="560px"
      destroy-on-close
    >
      <el-form
        ref="dialogFormRef"
        :model="dialogForm"
        :rules="dialogRules"
        label-width="110px"
      >
        <el-form-item label="App ID" prop="app_id">
          <el-input v-model="dialogForm.app_id" placeholder="请输入App ID" clearable />
        </el-form-item>
        <el-form-item label="App Secret" prop="app_secret">
          <el-input
            v-model="dialogForm.app_secret"
            type="password"
            placeholder="请输入App Secret"
            show-password
            clearable
          />
        </el-form-item>
        <el-form-item label="回调URL" prop="callback_url">
          <el-input v-model="dialogForm.callback_url" placeholder="请输入回调URL" clearable />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="dialogForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="dialogLoading" @click="handleDialogSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import request from '@/utils/request'

interface OAuthProvider {
  id: string
  name: string
  icon: string
  color: string
  app_id: string
  app_secret: string
  callback_url: string
  enabled: boolean
}

const loading = ref(false)
const dialogVisible = ref(false)
const dialogLoading = ref(false)
const dialogFormRef = ref<FormInstance>()
const editingProvider = ref<OAuthProvider | null>(null)

const providers = ref<OAuthProvider[]>([])

const dialogForm = reactive({
  app_id: '',
  app_secret: '',
  callback_url: '',
  enabled: false
})

const dialogRules: FormRules = {
  app_id: [{ required: true, message: '请输入App ID', trigger: 'blur' }],
  app_secret: [{ required: true, message: '请输入App Secret', trigger: 'blur' }],
  callback_url: [{ required: true, message: '请输入回调URL', trigger: 'blur' }]
}

const defaultProviders: OAuthProvider[] = [
  { id: 'wechat', name: '微信', icon: 'ChatDotRound', color: '#07C160', app_id: '', app_secret: '', callback_url: '', enabled: false },
  { id: 'qq', name: 'QQ', icon: 'User', color: '#12B7F5', app_id: '', app_secret: '', callback_url: '', enabled: false },
  { id: 'github', name: 'GitHub', icon: 'Link', color: '#333333', app_id: '', app_secret: '', callback_url: '', enabled: false },
  { id: 'google', name: 'Google', icon: 'ChromeFilled', color: '#4285F4', app_id: '', app_secret: '', callback_url: '', enabled: false },
  { id: 'weibo', name: '微博', icon: 'Promotion', color: '#E6162D', app_id: '', app_secret: '', callback_url: '', enabled: false },
  { id: 'dingtalk', name: '钉钉', icon: 'Position', color: '#0089FF', app_id: '', app_secret: '', callback_url: '', enabled: false }
]

const fetchProviders = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/admin/system/oauth')
    if (data?.data?.length) {
      providers.value = defaultProviders.map(dp => {
        const saved = data.data.find((p: any) => p.id === dp.id)
        return saved ? { ...dp, ...saved } : dp
      })
    } else {
      providers.value = [...defaultProviders]
    }
  } catch {
    providers.value = [...defaultProviders]
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: OAuthProvider) => {
  editingProvider.value = row
  dialogForm.app_id = row.app_id
  dialogForm.app_secret = row.app_secret
  dialogForm.callback_url = row.callback_url
  dialogForm.enabled = row.enabled
  dialogVisible.value = true
}

const handleToggle = async (row: OAuthProvider) => {
  const action = row.enabled ? '禁用' : '启用'
  try {
    await ElMessageBox.confirm(`确定要${action}${row.name}登录吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.post(`/api/admin/system/oauth/${row.id}/toggle`, { enabled: !row.enabled })
    row.enabled = !row.enabled
    ElMessage.success(`${row.name}已${action}`)
  } catch {
    // 取消操作
  }
}

const handleDialogSubmit = async () => {
  const valid = await dialogFormRef.value?.validate().catch(() => false)
  if (!valid || !editingProvider.value) return

  dialogLoading.value = true
  try {
    await request.post(`/api/admin/system/oauth/${editingProvider.value.id}`, { ...dialogForm })
    Object.assign(editingProvider.value, {
      app_id: dialogForm.app_id,
      app_secret: dialogForm.app_secret,
      callback_url: dialogForm.callback_url,
      enabled: dialogForm.enabled
    })
    dialogVisible.value = false
    ElMessage.success('配置已保存')
  } catch {
    ElMessage.error('保存失败，请重试')
  } finally {
    dialogLoading.value = false
  }
}

onMounted(() => {
  fetchProviders()
})
</script>

<style scoped lang="scss">
.oauth-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
    }
  }

  .provider-name {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 500;
  }

  .text-muted {
    color: var(--el-text-color-placeholder);
    font-size: 13px;
  }
}
</style>
