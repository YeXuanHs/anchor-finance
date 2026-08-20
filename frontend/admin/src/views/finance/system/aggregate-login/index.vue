<template>
  <div class="aggregate-login-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('aggregateLogin.title') }}</span>
          <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('aggregateLogin.addProvider') }}</el-button>
        </div>
      </template>

      <el-row :gutter="16" class="stats-row">
        <el-col :span="6"><el-card shadow="hover" class="stat-card"><div class="stat-label">{{ $t('aggregateLogin.totalProviders') }}</div><div class="stat-value">{{ stats.total || 0 }}</div></el-card></el-col>
        <el-col :span="6"><el-card shadow="hover" class="stat-card"><div class="stat-label">{{ $t('aggregateLogin.enabled') }}</div><div class="stat-value success">{{ stats.enabled || 0 }}</div></el-card></el-col>
        <el-col :span="6"><el-card shadow="hover" class="stat-card"><div class="stat-label">{{ $t('aggregateLogin.todayLogins') }}</div><div class="stat-value primary">{{ stats.today_logins || 0 }}</div></el-card></el-col>
        <el-col :span="6"><el-card shadow="hover" class="stat-card"><div class="stat-label">{{ $t('aggregateLogin.totalBinds') }}</div><div class="stat-value warning">{{ stats.total_binds || 0 }}</div></el-card></el-col>
      </el-row>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('aggregateLogin.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('aggregateLogin.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('aggregateLogin.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('common.enable')" :value="1" />
            <el-option :label="$t('common.disable')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" :label="$t('aggregateLogin.providerName')" min-width="150" />
        <el-table-column prop="provider" :label="$t('aggregateLogin.providerId')" width="150"><template #default="{ row }"><el-tag size="small">{{ row.provider }}</el-tag></template></el-table-column>
        <el-table-column prop="app_id" label="App ID" min-width="180" show-overflow-tooltip />
        <el-table-column prop="bind_count" :label="$t('aggregateLogin.bindCount')" width="90" align="center" />
        <el-table-column prop="login_count" :label="$t('aggregateLogin.loginCount')" width="100" align="center" />
        <el-table-column prop="status" :label="$t('aggregateLogin.status')" width="80" align="center"><template #default="{ row }"><el-switch v-model="row.status" :active-value="1" :inactive-value="0" @change="handleToggleStatus(row)" /></template></el-table-column>
        <el-table-column :label="$t('aggregateLogin.operations')" width="250" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleConfig(row)">{{ $t('aggregateLogin.config') }}</el-button>
            <el-button type="warning" link @click="handleViewBinds(row)">{{ $t('aggregateLogin.bindManage') }}</el-button>
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-popconfirm :title="$t('aggregateLogin.confirmDelete')" @confirm="handleDelete(row)"><template #reference><el-button type="danger" link>{{ $t('common.delete') }}</el-button></template></el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item :label="$t('aggregateLogin.providerName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('aggregateLogin.providerNamePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('aggregateLogin.providerId')" prop="provider">
          <el-select v-model="formData.provider" :placeholder="$t('aggregateLogin.selectProvider')" style="width: 100%" @change="handleProviderChange">
            <el-option v-for="p in providerOptions" :key="p.value" :label="$t(`aggregateLogin.providers.${p.value}`)" :value="p.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="App ID" prop="app_id">
          <el-input v-model="formData.app_id" :placeholder="$t('aggregateLogin.enterAppId')" />
        </el-form-item>
        <el-form-item label="App Secret" prop="app_secret">
          <el-input v-model="formData.app_secret" type="password" show-password :placeholder="$t('aggregateLogin.enterAppSecret')" />
        </el-form-item>
        <el-form-item :label="$t('aggregateLogin.callbackUrl')">
          <el-input :model-value="formData.callback_url || callbackUrl" disabled>
            <template #append><el-button @click="handleCopyCallback">{{ $t('aggregateLogin.copy') }}</el-button></template>
          </el-input>
          <div class="form-hint">{{ $t('aggregateLogin.callbackHint') }}</div>
        </el-form-item>
        <el-form-item :label="$t('aggregateLogin.extraConfig')">
          <el-input v-model="formData.extra_config" type="textarea" :rows="4" :placeholder="$t('aggregateLogin.extraConfigPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('aggregateLogin.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" :active-text="$t('common.enable')" :inactive-text="$t('common.disable')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="bindsDialogVisible" :title="`${$t('aggregateLogin.bindManage')} - ${currentProviderName}`" width="800px" destroy-on-close>
      <el-table :data="bindsList" v-loading="bindsLoading" border size="small">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="user_id" :label="$t('aggregateLogin.userId')" width="100" />
        <el-table-column prop="username" :label="$t('aggregateLogin.username')" width="120" />
        <el-table-column prop="open_id" label="OpenID" min-width="200" show-overflow-tooltip />
        <el-table-column prop="nickname" :label="$t('aggregateLogin.nickname')" width="120" show-overflow-tooltip />
        <el-table-column prop="created_at" :label="$t('aggregateLogin.bindTime')" width="180" />
        <el-table-column :label="$t('aggregateLogin.operations')" width="100" align="center">
          <template #default="{ row }"><el-popconfirm :title="$t('aggregateLogin.confirmUnbind')" @confirm="handleUnbind(row)"><template #reference><el-button type="danger" link size="small">{{ $t('aggregateLogin.unbind') }}</el-button></template></el-popconfirm></template>
        </el-table-column>
      </el-table>
      <div class="pagination-container">
        <el-pagination v-model:current-page="bindsPagination.page" v-model:page-size="bindsPagination.page_size" :page-sizes="[10, 20, 50]" :total="bindsPagination.total" layout="total, sizes, prev, pager, next" @size-change="fetchBinds" @current-change="fetchBinds" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'AggregateLogin' })

const providerOptions = [
  { label: 'WeChat', value: 'wechat' }, { label: 'QQ', value: 'qq' }, { label: 'Weibo', value: 'weibo' },
  { label: 'Alipay', value: 'alipay' }, { label: 'GitHub', value: 'github' }, { label: 'Google', value: 'google' },
  { label: 'Facebook', value: 'facebook' }, { label: 'Twitter', value: 'twitter' }, { label: 'DingTalk', value: 'dingtalk' },
  { label: 'Feishu', value: 'feishu' }, { label: 'Apple', value: 'apple' }, { label: 'Telegram', value: 'telegram' },
]

const loading = ref(false)
const submitLoading = ref(false)
const bindsLoading = ref(false)
const dialogVisible = ref(false)
const bindsDialogVisible = ref(false)
const dialogTitle = ref($t('aggregateLogin.addProvider'))
const formRef = ref<FormInstance>()
const currentProviderId = ref<number>(0)
const currentProviderName = ref('')

const stats = reactive({ total: 0, enabled: 0, today_logins: 0, total_binds: 0 })
const searchForm = reactive({ keyword: '', status: undefined as number | undefined })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const bindsList = ref<any[]>([])
const bindsPagination = reactive({ page: 1, page_size: 20, total: 0 })

const formData = reactive({ id: undefined as number | undefined, name: '', provider: '', app_id: '', app_secret: '', callback_url: '', extra_config: '', status: 1 })

const formRules: FormRules = {
  name: [{ required: true, message: () => $t('aggregateLogin.enterProviderName'), trigger: 'blur' }],
  provider: [{ required: true, message: () => $t('aggregateLogin.selectProvider'), trigger: 'change' }],
  app_id: [{ required: true, message: () => $t('aggregateLogin.enterAppId'), trigger: 'blur' }],
  app_secret: [{ required: true, message: () => $t('aggregateLogin.enterAppSecret'), trigger: 'blur' }]
}

const callbackUrl = computed(() => `${window.location.origin}/oauth/${formData.provider}/callback`)
const handleProviderChange = () => { const option = providerOptions.find(p => p.value === formData.provider); if (option && !formData.name) formData.name = $t(`aggregateLogin.providers.${option.value}`) + $t('aggregateLogin.loginSuffix') }
const handleCopyCallback = () => { navigator.clipboard.writeText(callbackUrl.value); ElMessage.success($t('aggregateLogin.callbackCopied')) }

const fetchStats = async () => { try { const data = await request.get({ url: '/api/admin/aggregate-login/stats' }); Object.assign(stats, data) } catch {} }
const fetchData = async () => {
  loading.value = true
  try { const data = await request.get({ url: '/api/admin/aggregate-login', params: { page: pagination.page, page_size: pagination.page_size, ...searchForm } }); tableData.value = data.list || data || []; pagination.total = data.total || 0 } catch { ElMessage.error($t('aggregateLogin.fetchFailed')) } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.status = undefined; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleAdd = () => { dialogTitle.value = $t('aggregateLogin.addProvider'); formData.id = undefined; formData.name = ''; formData.provider = ''; formData.app_id = ''; formData.app_secret = ''; formData.callback_url = ''; formData.extra_config = ''; formData.status = 1; dialogVisible.value = true }
const handleEdit = (row: any) => { dialogTitle.value = $t('aggregateLogin.editProvider'); Object.assign(formData, row); dialogVisible.value = true }
const handleConfig = (row: any) => { dialogTitle.value = `${$t('aggregateLogin.config')} - ${row.name}`; Object.assign(formData, row); dialogVisible.value = true }

const handleToggleStatus = async (row: any) => {
  try { await request.put({ url: `/api/admin/aggregate-login/${row.id}`, params: { status: row.status } }); ElMessage.success(row.status ? $t('common.enabled') : $t('common.disabled')); fetchStats() } catch { row.status = row.status ? 0 : 1; ElMessage.error($t('common.operateFailed')) }
}

const handleDelete = async (row: any) => { try { await request.del({ url: `/api/admin/aggregate-login/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchData(); fetchStats() } catch { ElMessage.error($t('common.deleteFailed')) } }

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      const submitData = { ...formData, callback_url: callbackUrl.value }
      if (formData.id) await request.put({ url: `/api/admin/aggregate-login/${formData.id}`, params: submitData })
      else await request.post({ url: '/api/admin/aggregate-login', params: submitData })
      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess')); dialogVisible.value = false; fetchData(); fetchStats()
    } catch { ElMessage.error($t('common.operateFailed')) } finally { submitLoading.value = false }
  })
}

const handleViewBinds = async (row: any) => { currentProviderId.value = row.id; currentProviderName.value = row.name; bindsPagination.page = 1; bindsDialogVisible.value = true; await fetchBinds() }

const fetchBinds = async () => {
  bindsLoading.value = true
  try { const data = await request.get({ url: `/api/admin/aggregate-login/${currentProviderId.value}/binds`, params: { page: bindsPagination.page, page_size: bindsPagination.page_size } }); bindsList.value = data.list || data || []; bindsPagination.total = data.total || 0 } catch { ElMessage.error($t('aggregateLogin.fetchBindsFailed')) } finally { bindsLoading.value = false }
}

const handleUnbind = async (row: any) => { try { await request.del({ url: `/api/admin/aggregate-login/${currentProviderId.value}/binds/${row.id}` }); ElMessage.success($t('aggregateLogin.unbindSuccess')); fetchBinds(); fetchStats() } catch { ElMessage.error($t('aggregateLogin.unbindFailed')) } }

onMounted(() => { fetchStats(); fetchData() })
</script>

<style scoped lang="scss">
.aggregate-login-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.stats-row { margin-bottom: 20px; }
.stat-card { text-align: center; }
.stat-label { font-size: 13px; color: var(--el-text-color-secondary); margin-bottom: 8px; }
.stat-value { font-size: 28px; font-weight: 600; color: var(--el-text-color-primary); &.success { color: var(--el-color-success); } &.primary { color: var(--el-color-primary); } &.warning { color: var(--el-color-warning); } }
.form-hint { color: #909399; font-size: 12px; margin-top: 4px; }
</style>
