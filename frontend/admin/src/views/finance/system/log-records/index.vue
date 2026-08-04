<template>
  <div class="log-records-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>日志记录</span>
          <el-button type="danger" @click="handleClearLogs">
            <el-icon><Delete /></el-icon>
            清理日志
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="日志内容/用户" clearable />
        </el-form-item>
        <el-form-item label="日志类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="登录日志" value="login" />
            <el-option label="操作日志" value="operation" />
            <el-option label="错误日志" value="error" />
            <el-option label="系统日志" value="system" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
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
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="type" label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user" label="用户" width="120" />
        <el-table-column prop="content" label="日志内容" min-width="300" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP地址" width="140" />
        <el-table-column prop="created_at" label="时间" width="170" />
        <el-table-column label="操作" width="80" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
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

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="日志详情" width="600px">
      <el-descriptions :column="1" border v-if="currentLog">
        <el-descriptions-item label="ID">{{ currentLog.id }}</el-descriptions-item>
        <el-descriptions-item label="类型">
          <el-tag :type="getTypeTag(currentLog.type)" size="small">{{ getTypeText(currentLog.type) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="用户">{{ currentLog.user }}</el-descriptions-item>
        <el-descriptions-item label="日志内容">{{ currentLog.content }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ currentLog.ip }}</el-descriptions-item>
        <el-descriptions-item label="时间">{{ currentLog.created_at }}</el-descriptions-item>
        <el-descriptions-item label="详细信息">
          <pre style="white-space: pre-wrap; word-break: break-all;">{{ currentLog.detail || '无' }}</pre>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)

const searchForm = reactive({
  keyword: '',
  type: '',
  date_range: [] as string[]
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])
const detailVisible = ref(false)
const currentLog = ref<any>(null)

const getTypeTag = (type: string) => {
  const map: Record<string, string> = {
    login: 'success',
    operation: 'primary',
    error: 'danger',
    system: 'info'
  }
  return (map[type] || 'info') as any
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    login: '登录日志',
    operation: '操作日志',
    error: '错误日志',
    system: '系统日志'
  }
  return map[type] || type
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      keyword: searchForm.keyword || undefined,
      type: searchForm.type || undefined
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({ url: '/api/admin/log-records', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取日志记录失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.type = ''
  searchForm.date_range = []
  handleSearch()
}

const handleViewDetail = (row: any) => {
  currentLog.value = row
  detailVisible.value = true
}

const handleClearLogs = async () => {
  await ElMessageBox.confirm('确定清理所有日志记录吗？此操作不可恢复。', '警告', {
    type: 'warning'
  })
  try {
    await request.post({ url: '/api/admin/log-records/export' })
    ElMessage.success('清理成功')
    fetchData()
  } catch (error) {
    ElMessage.error('清理失败')
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
  fetchData()
})
</script>

<style scoped lang="scss">
.log-records-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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
