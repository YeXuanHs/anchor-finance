<template>
  <div class="aggregate-login-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>聚合登录管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加提供商
          </el-button>
        </div>
      </template>

      <!-- 统计卡片 -->
      <el-row :gutter="16" class="stats-row">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">总提供商</div>
            <div class="stat-value">{{ stats.total || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">已启用</div>
            <div class="stat-value success">{{ stats.enabled || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">今日登录</div>
            <div class="stat-value primary">{{ stats.today_logins || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">总绑定数</div>
            <div class="stat-value warning">{{ stats.total_binds || 0 }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="提供商名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="提供商名称" min-width="150" />
        <el-table-column prop="provider" label="提供商标识" width="150">
          <template #default="{ row }">
            <el-tag size="small">{{ row.provider }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="app_id" label="App ID" min-width="180" show-overflow-tooltip />
        <el-table-column prop="bind_count" label="绑定数" width="90" align="center" />
        <el-table-column prop="login_count" label="登录次数" width="100" align="center" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-switch v-model="row.status" :active-value="1" :inactive-value="0" @change="handleToggleStatus(row)" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleConfig(row)">配置</el-button>
            <el-button type="warning" link @click="handleViewBinds(row)">绑定管理</el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该提供商吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item label="提供商名称" prop="name">
          <el-input v-model="formData.name" placeholder="如：微信登录" />
        </el-form-item>
        <el-form-item label="提供商标识" prop="provider">
          <el-select v-model="formData.provider" placeholder="请选择提供商" style="width: 100%" @change="handleProviderChange">
            <el-option v-for="p in providerOptions" :key="p.value" :label="p.label" :value="p.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="App ID" prop="app_id">
          <el-input v-model="formData.app_id" placeholder="请输入App ID" />
        </el-form-item>
        <el-form-item label="App Secret" prop="app_secret">
          <el-input v-model="formData.app_secret" type="password" show-password placeholder="请输入App Secret" />
        </el-form-item>
        <el-form-item label="回调地址">
          <el-input :model-value="formData.callback_url || callbackUrl" disabled>
            <template #append>
              <el-button @click="handleCopyCallback">复制</el-button>
            </template>
          </el-input>
          <div class="form-hint">请在第三方平台填写此回调地址</div>
        </el-form-item>
        <el-form-item label="额外配置">
          <el-input v-model="formData.extra_config" type="textarea" :rows="4" placeholder="JSON格式的额外配置参数" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 绑定管理对话框 -->
    <el-dialog v-model="bindsDialogVisible" :title="`绑定管理 - ${currentProviderName}`" width="800px" destroy-on-close>
      <el-table :data="bindsList" v-loading="bindsLoading" border size="small">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="open_id" label="OpenID" min-width="200" show-overflow-tooltip />
        <el-table-column prop="nickname" label="昵称" width="120" show-overflow-tooltip />
        <el-table-column prop="created_at" label="绑定时间" width="180" />
        <el-table-column label="操作" width="100" align="center">
          <template #default="{ row }">
            <el-popconfirm title="确定解绑该账号吗？" @confirm="handleUnbind(row)">
              <template #reference>
                <el-button type="danger" link size="small">解绑</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="bindsPagination.page"
          v-model:page-size="bindsPagination.page_size"
          :page-sizes="[10, 20, 50]"
          :total="bindsPagination.total"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchBinds"
          @current-change="fetchBinds"
        />
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

defineOptions({ name: 'AggregateLogin' })

const providerOptions = [
  { label: '微信', value: 'wechat' },
  { label: 'QQ', value: 'qq' },
  { label: '微博', value: 'weibo' },
  { label: '支付宝', value: 'alipay' },
  { label: 'GitHub', value: 'github' },
  { label: 'Google', value: 'google' },
  { label: 'Facebook', value: 'facebook' },
  { label: 'Twitter', value: 'twitter' },
  { label: '钉钉', value: 'dingtalk' },
  { label: '飞书', value: 'feishu' },
  { label: 'Apple', value: 'apple' },
  { label: 'Telegram', value: 'telegram' },
]

const loading = ref(false)
const submitLoading = ref(false)
const bindsLoading = ref(false)
const dialogVisible = ref(false)
const bindsDialogVisible = ref(false)
const dialogTitle = ref('添加提供商')
const formRef = ref<FormInstance>()
const currentProviderId = ref<number>(0)
const currentProviderName = ref('')

const stats = reactive({ total: 0, enabled: 0, today_logins: 0, total_binds: 0 })
const searchForm = reactive({ keyword: '', status: undefined as number | undefined })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const bindsList = ref<any[]>([])
const bindsPagination = reactive({ page: 1, page_size: 20, total: 0 })

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  provider: '',
  app_id: '',
  app_secret: '',
  callback_url: '',
  extra_config: '',
  status: 1
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入提供商名称', trigger: 'blur' }],
  provider: [{ required: true, message: '请选择提供商', trigger: 'change' }],
  app_id: [{ required: true, message: '请输入App ID', trigger: 'blur' }],
  app_secret: [{ required: true, message: '请输入App Secret', trigger: 'blur' }]
}

const callbackUrl = computed(() => {
  const domain = window.location.origin
  return `${domain}/oauth/${formData.provider}/callback`
})

const handleProviderChange = () => {
  const option = providerOptions.find(p => p.value === formData.provider)
  if (option && !formData.name) {
    formData.name = option.label + '登录'
  }
}

const handleCopyCallback = () => {
  navigator.clipboard.writeText(callbackUrl.value)
  ElMessage.success('回调地址已复制')
}

const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/aggregate-login/stats' })
    Object.assign(stats, data)
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/aggregate-login',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || data || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取提供商列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.status = undefined; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleAdd = () => {
  dialogTitle.value = '添加提供商'
  formData.id = undefined; formData.name = ''; formData.provider = ''
  formData.app_id = ''; formData.app_secret = ''; formData.callback_url = ''
  formData.extra_config = ''; formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑提供商'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleConfig = (row: any) => {
  dialogTitle.value = `配置 - ${row.name}`
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleToggleStatus = async (row: any) => {
  try {
    await request.put({
      url: `/api/admin/aggregate-login/${row.id}`,
      params: { status: row.status }
    })
    ElMessage.success(row.status ? '已启用' : '已禁用')
    fetchStats()
  } catch (error) {
    row.status = row.status ? 0 : 1
    ElMessage.error('操作失败')
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/aggregate-login/${row.id}` })
    ElMessage.success('删除成功')
    fetchData()
    fetchStats()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      const submitData = { ...formData, callback_url: callbackUrl.value }
      if (formData.id) {
        await request.put({ url: `/api/admin/aggregate-login/${formData.id}`, params: submitData })
      } else {
        await request.post({ url: '/api/admin/aggregate-login', params: submitData })
      }
      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchData()
      fetchStats()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

const handleViewBinds = async (row: any) => {
  currentProviderId.value = row.id
  currentProviderName.value = row.name
  bindsPagination.page = 1
  bindsDialogVisible.value = true
  await fetchBinds()
}

const fetchBinds = async () => {
  bindsLoading.value = true
  try {
    const data = await request.get({
      url: `/api/admin/aggregate-login/${currentProviderId.value}/binds`,
      params: { page: bindsPagination.page, page_size: bindsPagination.page_size }
    })
    bindsList.value = data.list || data || []
    bindsPagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取绑定列表失败')
  } finally {
    bindsLoading.value = false
  }
}

const handleUnbind = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/aggregate-login/${currentProviderId.value}/binds/${row.id}` })
    ElMessage.success('解绑成功')
    fetchBinds()
    fetchStats()
  } catch (error) {
    ElMessage.error('解绑失败')
  }
}

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
.stat-value {
  font-size: 28px; font-weight: 600; color: var(--el-text-color-primary);
  &.success { color: var(--el-color-success); }
  &.primary { color: var(--el-color-primary); }
  &.warning { color: var(--el-color-warning); }
}
.form-hint { color: #909399; font-size: 12px; margin-top: 4px; }
</style>
