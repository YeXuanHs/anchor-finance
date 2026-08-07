<template>
  <div class="oauth-bind-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>OAuth账号绑定管理</span>
          <div class="header-actions">
            <el-button type="danger" size="small" @click="handleBatchUnbind" :disabled="!selectedRows.length">
              <el-icon><Close /></el-icon>
              批量解绑
            </el-button>
          </div>
        </div>
      </template>

      <!-- 统计卡片 -->
      <el-row :gutter="16" class="stats-row">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">总绑定数</div>
            <div class="stat-value">{{ stats.total || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">微信绑定</div>
            <div class="stat-value success">{{ stats.wechat || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">QQ绑定</div>
            <div class="stat-value primary">{{ stats.qq || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">GitHub绑定</div>
            <div class="stat-value warning">{{ stats.github || 0 }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.user_id" placeholder="用户ID" clearable />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="searchForm.username" placeholder="用户名" clearable />
        </el-form-item>
        <el-form-item label="绑定类型">
          <el-select v-model="searchForm.provider" placeholder="全部" clearable>
            <el-option v-for="p in providerOptions" :key="p.value" :label="p.label" :value="p.value" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="user_id" label="用户ID" width="90" align="center" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="provider" label="绑定类型" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="providerTagMap[row.provider]" size="small">
              {{ providerTextMap[row.provider] || row.provider }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="open_id" label="OpenID" min-width="200" show-overflow-tooltip />
        <el-table-column prop="nickname" label="第三方昵称" width="150" show-overflow-tooltip />
        <el-table-column prop="avatar" label="头像" width="80" align="center">
          <template #default="{ row }">
            <el-avatar v-if="row.avatar" :src="row.avatar" :size="32" />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="绑定时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
            <el-popconfirm title="确定解绑该账号吗？" @confirm="handleUnbind(row)">
              <template #reference>
                <el-button type="danger" link>解绑</el-button>
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

    <!-- 绑定详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="绑定详情" width="500px" destroy-on-close>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="绑定ID">{{ currentBind.id }}</el-descriptions-item>
        <el-descriptions-item label="用户ID">{{ currentBind.user_id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ currentBind.username }}</el-descriptions-item>
        <el-descriptions-item label="绑定类型">
          <el-tag :type="providerTagMap[currentBind.provider]">
            {{ providerTextMap[currentBind.provider] || currentBind.provider }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="OpenID">{{ currentBind.open_id }}</el-descriptions-item>
        <el-descriptions-item label="第三方昵称">{{ currentBind.nickname || '-' }}</el-descriptions-item>
        <el-descriptions-item label="头像">
          <el-avatar v-if="currentBind.avatar" :src="currentBind.avatar" :size="48" />
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="绑定时间">{{ currentBind.created_at }}</el-descriptions-item>
        <el-descriptions-item label="最后登录">{{ currentBind.last_login_at || '-' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Close } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'OAuthBindManage' })

const providerOptions = [
  { label: '微信', value: 'wechat' },
  { label: 'QQ', value: 'qq' },
  { label: '微博', value: 'weibo' },
  { label: 'GitHub', value: 'github' },
  { label: 'Google', value: 'google' },
  { label: 'Facebook', value: 'facebook' },
  { label: 'Twitter', value: 'twitter' },
  { label: '支付宝', value: 'alipay' },
  { label: '钉钉', value: 'dingtalk' },
  { label: '飞书', value: 'feishu' },
  { label: 'Apple', value: 'apple' },
  { label: 'Telegram', value: 'telegram' }
]

const providerTextMap: Record<string, string> = {
  wechat: '微信', qq: 'QQ', weibo: '微博', github: 'GitHub',
  google: 'Google', facebook: 'Facebook', twitter: 'Twitter',
  alipay: '支付宝', dingtalk: '钉钉', feishu: '飞书',
  apple: 'Apple', telegram: 'Telegram'
}
const providerTagMap: Record<string, any> = {
  wechat: 'success', qq: 'primary', weibo: 'warning', github: 'info',
  google: 'danger', facebook: 'primary', twitter: 'info',
  alipay: 'primary', dingtalk: 'primary', feishu: 'primary',
  apple: 'info', telegram: 'primary'
}

const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedRows = ref<any[]>([])

const stats = reactive({ total: 0, wechat: 0, qq: 0, github: 0 })
const searchForm = reactive({ user_id: '', username: '', provider: '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const currentBind = reactive<any>({})

const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/oauth-bind/stats' })
    Object.assign(stats, data)
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/oauth-bind',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取绑定列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { Object.assign(searchForm, { user_id: '', username: '', provider: '' }); handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }
const handleSelectionChange = (rows: any[]) => { selectedRows.value = rows }

const handleViewDetail = (row: any) => {
  Object.assign(currentBind, row)
  detailDialogVisible.value = true
}

const handleUnbind = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/oauth-bind/${row.id}` })
    ElMessage.success('解绑成功')
    fetchData()
    fetchStats()
  } catch (error) {
    ElMessage.error('解绑失败')
  }
}

const handleBatchUnbind = async () => {
  try {
    await ElMessageBox.confirm(`确定批量解绑选中的 ${selectedRows.value.length} 个绑定吗？`, '批量解绑确认', { type: 'warning' })
    const ids = selectedRows.value.map((r) => r.id)
    await request.post({ url: '/api/admin/oauth-bind/batch-unbind', params: { ids } })
    ElMessage.success('批量解绑成功')
    fetchData()
    fetchStats()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('批量解绑失败')
  }
}

onMounted(() => { fetchStats(); fetchData() })
</script>

<style scoped lang="scss">
.oauth-bind-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.header-actions { display: flex; gap: 8px; }
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
</style>
