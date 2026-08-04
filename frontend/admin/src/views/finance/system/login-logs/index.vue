<template>
  <div class="login-logs-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>登录日志</span>
          <div class="header-actions">
            <el-button type="danger" size="small" @click="handleClearOld">
              <el-icon><Delete /></el-icon>
              清理日志
            </el-button>
            <el-button type="success" size="small" @click="handleExport">
              <el-icon><Download /></el-icon>
              导出
            </el-button>
          </div>
        </div>
      </template>

      <!-- 统计卡片 -->
      <el-row :gutter="16" class="stats-row">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">今日登录</div>
            <div class="stat-value">{{ stats.today_total || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">今日成功</div>
            <div class="stat-value success">{{ stats.today_success || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">今日失败</div>
            <div class="stat-value danger">{{ stats.today_failed || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-label">失败IP数</div>
            <div class="stat-value warning">{{ stats.failed_ips || 0 }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="用户名">
          <el-input v-model="searchForm.username" placeholder="用户名" clearable />
        </el-form-item>
        <el-form-item label="IP地址">
          <el-input v-model="searchForm.ip" placeholder="IP地址" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="成功" value="success" />
            <el-option label="失败" value="failed" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="ip" label="IP地址" width="140" />
        <el-table-column prop="location" label="登录地点" width="150" show-overflow-tooltip />
        <el-table-column prop="user_agent" label="浏览器" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="fail_reason" label="失败原因" width="150" show-overflow-tooltip />
        <el-table-column prop="created_at" label="登录时间" width="180" />
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[20, 50, 100, 200]"
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
import { Delete, Download } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'LoginLogsManage' })

const loading = ref(false)

const stats = reactive({
  today_total: 0,
  today_success: 0,
  today_failed: 0,
  failed_ips: 0
})

const searchForm = reactive({
  username: '',
  ip: '',
  status: '',
  date_range: [] as string[]
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])

const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/login-logs/stats' })
    Object.assign(stats, data)
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      username: searchForm.username || undefined,
      ip: searchForm.ip || undefined,
      status: searchForm.status || undefined
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({ url: '/api/admin/login-logs', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取登录日志失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.username = ''
  searchForm.ip = ''
  searchForm.status = ''
  searchForm.date_range = []
  handleSearch()
}

const handleClearOld = async () => {
  try {
    await ElMessageBox.confirm('确定清理30天前的登录日志吗？', '清理确认', { type: 'warning' })
    await request.post({ url: '/api/admin/login-logs/clean', params: { days: 30 } })
    ElMessage.success('清理完成')
    fetchData()
    fetchStats()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error('清理失败')
  }
}

const handleExport = async () => {
  try {
    await request.get({ url: '/api/admin/login-logs/export', params: { ...searchForm } })
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchData()
}

const handlePageChange = () => {
  fetchData()
}

onMounted(() => {
  fetchStats()
  fetchData()
})
</script>

<style scoped lang="scss">
.login-logs-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;

  .stat-label {
    font-size: 13px;
    color: var(--el-text-color-secondary);
    margin-bottom: 8px;
  }

  .stat-value {
    font-size: 28px;
    font-weight: 600;
    color: var(--el-text-color-primary);

    &.success { color: var(--el-color-success); }
    &.danger { color: var(--el-color-danger); }
    &.warning { color: var(--el-color-warning); }
  }
}

.search-form {
  margin-bottom: 20px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
