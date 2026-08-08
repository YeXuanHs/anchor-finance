<template>
  <div class="system-log-page">
    <!-- 标签页 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="操作日志" name="operation" />
      <el-tab-pane label="登录日志" name="login" />
      <el-tab-pane label="API日志" name="api" />
    </el-tabs>

    <!-- 搜索筛选区域 -->
    <el-card shadow="never" class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="操作内容/用户/IP"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="操作人" v-if="activeTab === 'operation'">
          <el-input
            v-model="searchForm.operator"
            placeholder="操作人"
            clearable
            style="width: 150px"
          />
        </el-form-item>
        <el-form-item label="IP地址">
          <el-input
            v-model="searchForm.ip"
            placeholder="IP地址"
            clearable
            style="width: 150px"
          />
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作栏 -->
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <div class="action-left">
          <el-button type="danger" @click="handleClearLogs">
            <el-icon><Delete /></el-icon>
            清空日志
          </el-button>
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
        <div class="action-right">
          <el-button circle @click="fetchList">
            <el-icon><Refresh /></el-icon>
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <!-- 操作日志 -->
      <el-table
        v-if="activeTab === 'operation'"
        v-loading="loading"
        :data="tableData"
        border
        stripe
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="operator_name" label="操作人" width="120" />
        <el-table-column prop="action" label="操作" width="150" />
        <el-table-column prop="module" label="模块" width="120" />
        <el-table-column prop="detail" label="详情" min-width="250" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP地址" width="130" />
        <el-table-column prop="user_agent" label="浏览器" width="150" show-overflow-tooltip />
        <el-table-column prop="created_at" label="时间" width="170" />
      </el-table>

      <!-- 登录日志 -->
      <el-table
        v-if="activeTab === 'login'"
        v-loading="loading"
        :data="tableData"
        border
        stripe
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="ip" label="IP地址" width="130" />
        <el-table-column prop="location" label="登录地点" width="150" />
        <el-table-column prop="user_agent" label="浏览器" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="failure_reason" label="失败原因" width="150" v-if="searchForm.status === 'failure'" />
        <el-table-column prop="created_at" label="时间" width="170" />
      </el-table>

      <!-- API日志 -->
      <el-table
        v-if="activeTab === 'api'"
        v-loading="loading"
        :data="tableData"
        border
        stripe
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="method" label="方法" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="getMethodType(row.method)" size="small">{{ row.method }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路径" min-width="250" />
        <el-table-column prop="status_code" label="状态码" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status_code)" size="small">{{ row.status_code }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP地址" width="130" />
        <el-table-column prop="duration" label="耗时" width="80" align="center">
          <template #default="{ row }">{{ row.duration }}ms</template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="170" />
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Delete, Download } from '@element-plus/icons-vue'
import request from '@/utils/http'

const loading = ref(false)
const tableData = ref([])
const activeTab = ref('operation')

// 搜索表单
const searchForm = reactive({
  keyword: '',
  operator: '',
  ip: '',
  date_range: null as [Date, Date] | null
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// HTTP 方法类型
const getMethodType = (method: string) => {
  const map: Record<string, string> = {
    GET: 'success',
    POST: 'primary',
    PUT: 'warning',
    DELETE: 'danger',
    PATCH: 'info'
  }
  return map[method] || 'info'
}

// 状态码类型
const getStatusType = (code: number) => {
  if (code >= 200 && code < 300) return 'success'
  if (code >= 300 && code < 400) return 'info'
  if (code >= 400 && code < 500) return 'warning'
  return 'danger'
}

// 标签页切换
const handleTabChange = () => {
  pagination.page = 1
  searchForm.keyword = ''
  searchForm.operator = ''
  searchForm.ip = ''
  searchForm.date_range = null
  fetchList()
}

// 获取列表数据
const fetchList = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      type: activeTab.value
    }

    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.operator) params.operator = searchForm.operator
    if (searchForm.ip) params.ip = searchForm.ip
    if (searchForm.date_range) {
      params.start_date = searchForm.date_range[0].toISOString().split('T')[0]
      params.end_date = searchForm.date_range[1].toISOString().split('T')[0]
    }

    const data = await request.get({ url: '/api/admin/logs', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取日志列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.operator = ''
  searchForm.ip = ''
  searchForm.date_range = null
  pagination.page = 1
  fetchList()
}

// 分页大小变化
const handleSizeChange = (size: number) => {
  pagination.page_size = size
  pagination.page = 1
  fetchList()
}

// 页码变化
const handlePageChange = (page: number) => {
  pagination.page = page
  fetchList()
}

// 清空日志
const handleClearLogs = async () => {
  try {
    await ElMessageBox.confirm('确定要清空所有日志吗？此操作不可恢复。', '确认清空', {
      type: 'warning'
    })
    await request.delete({ url: '/api/admin/logs', params: { type: activeTab.value } })
    ElMessage.success('日志已清空')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('清空日志失败:', error)
    }
  }
}

// 导出
const handleExport = async () => {
  try {
    const params: any = { type: activeTab.value }
    if (searchForm.keyword) params.keyword = searchForm.keyword

    const response = await request.get({
      url: '/api/admin/logs/export',
      params,
      responseType: 'blob'
    })

    const url = window.URL.createObjectURL(new Blob([response]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `${activeTab.value}日志_${new Date().toISOString().split('T')[0]}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('导出失败:', error)
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.system-log-page {
  padding: 16px;
}

.search-card {
  margin-bottom: 16px;

  :deep(.el-card__body) {
    padding-bottom: 0;
  }
}

.action-card {
  margin-bottom: 16px;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.action-left {
  display: flex;
  gap: 8px;
}

.table-card {
  :deep(.el-card__body) {
    padding: 0;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px;
}
</style>
