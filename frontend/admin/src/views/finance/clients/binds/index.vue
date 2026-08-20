<template>
  <div class="binds-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsBinds.title') }}</span>
          <div class="header-actions">
            <el-button type="danger" size="small" @click="handleBatchUnbind" :disabled="!selectedRows.length">
              <el-icon><Close /></el-icon>
              {{ $t('clientsBinds.batchUnbind') }}
            </el-button>
          </div>
        </div>
      </template>

      <!-- 统计卡片 -->
      <el-row :gutter="16" class="stats-row">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">{{ $t('clientsBinds.totalBinds') }}</div>
            <div class="stat-value">{{ stats.total || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">{{ $t('clientsBinds.wechatBinds') }}</div>
            <div class="stat-value success">{{ stats.wechat || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">{{ $t('clientsBinds.qqBinds') }}</div>
            <div class="stat-value primary">{{ stats.qq || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">{{ $t('clientsBinds.githubBinds') }}</div>
            <div class="stat-value warning">{{ stats.github || 0 }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('clientsBinds.userId')">
          <el-input v-model="searchForm.user_id" :placeholder="$t('clientsBinds.userId')" clearable />
        </el-form-item>
        <el-form-item :label="$t('common.username')">
          <el-input v-model="searchForm.username" :placeholder="$t('common.username')" clearable />
        </el-form-item>
        <el-form-item :label="$t('clientsBinds.bindType')">
          <el-select v-model="searchForm.provider" :placeholder="$t('common.all')" clearable>
            <el-option label="WeChat" value="wechat" />
            <el-option label="QQ" value="qq" />
            <el-option label="GitHub" value="github" />
            <el-option label="Google" value="google" />
            <el-option label="Weibo" value="weibo" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="user_id" :label="$t('clientsBinds.userId')" width="100" align="center" />
        <el-table-column prop="username" :label="$t('common.username')" width="120" />
        <el-table-column prop="provider" :label="$t('clientsBinds.bindType')" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="providerTagMap[row.provider]" size="small">
              {{ providerTextMap[row.provider] || row.provider }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="open_id" label="OpenID" min-width="200" show-overflow-tooltip />
        <el-table-column prop="nickname" :label="$t('clientsBinds.thirdPartyNickname')" width="150" show-overflow-tooltip />
        <el-table-column prop="created_at" :label="$t('clientsBinds.bindTime')" width="180" />
        <el-table-column :label="$t('common.action')" width="120" fixed="right">
          <template #default="{ row }">
            <el-popconfirm :title="$t('clientsBinds.confirmUnbind')" @confirm="handleUnbind(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('clientsBinds.unbind') }}</el-button>
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Close } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'BindsManage' })

const providerTextMap: Record<string, string> = { wechat: 'WeChat', qq: 'QQ', github: 'GitHub', google: 'Google', weibo: 'Weibo' }
const providerTagMap: Record<string, any> = { wechat: 'success', qq: 'primary', github: 'warning', google: 'danger', weibo: 'info' }

const loading = ref(false)
const selectedRows = ref<any[]>([])

const stats = reactive({ total: 0, wechat: 0, qq: 0, github: 0 })

const searchForm = reactive({
  user_id: '',
  username: '',
  provider: ''
})

const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/binds/stats' })
    Object.assign(stats, data)
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/binds',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error($t('clientsBinds.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { Object.assign(searchForm, { user_id: '', username: '', provider: '' }); handleSearch() }

const handleSelectionChange = (rows: any[]) => { selectedRows.value = rows }

const handleUnbind = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/binds/${row.id}` })
    ElMessage.success($t('clientsBinds.unbindSuccess'))
    fetchData()
    fetchStats()
  } catch (error) {
    ElMessage.error($t('clientsBinds.unbindFailed'))
  }
}

const handleBatchUnbind = async () => {
  try {
    await ElMessageBox.confirm($t('clientsBinds.confirmBatchUnbind', { count: selectedRows.value.length }), $t('clientsBinds.batchUnbindConfirm'), { type: 'warning' })
    const ids = selectedRows.value.map((r) => r.id)
    await request.post({ url: '/api/admin/binds/batch-unbind', params: { ids } })
    ElMessage.success($t('clientsBinds.batchUnbindSuccess'))
    fetchData()
    fetchStats()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error($t('clientsBinds.batchUnbindFailed'))
  }
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchStats(); fetchData() })
</script>

<style scoped lang="scss">
.binds-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.header-actions { display: flex; gap: 8px; }
.stats-row { margin-bottom: 20px; }
.stat-card { text-align: center; .stat-label { font-size: 13px; color: var(--el-text-color-secondary); margin-bottom: 8px; } .stat-value { font-size: 28px; font-weight: 600; color: var(--el-text-color-primary); &.success { color: var(--el-color-success); } &.primary { color: var(--el-color-primary); } &.warning { color: var(--el-color-warning); } } }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
