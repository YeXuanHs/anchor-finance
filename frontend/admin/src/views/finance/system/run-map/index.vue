<template>
  <div class="run-map-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>运行映射</span>
          <el-button @click="fetchData">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="路由/控制器" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="成功" :value="1" />
            <el-option label="失败" :value="0" />
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
        <el-table-column prop="route" label="路由" min-width="200" show-overflow-tooltip />
        <el-table-column prop="controller" label="控制器" min-width="150" />
        <el-table-column prop="action" label="方法" width="120" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="run_time" label="执行时间(ms)" width="120" align="center" />
        <el-table-column prop="created_at" label="执行时间" width="170" />
        <el-table-column label="操作" width="150" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
            <el-button v-if="row.status === 0" type="warning" link @click="handleRetry(row)">重试</el-button>
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
    <el-dialog v-model="detailVisible" title="运行详情" width="700px">
      <el-descriptions :column="2" border v-if="currentItem">
        <el-descriptions-item label="ID">{{ currentItem.id }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentItem.status === 1 ? 'success' : 'danger'" size="small">
            {{ currentItem.status === 1 ? '成功' : '失败' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="路由" :span="2">{{ currentItem.route }}</el-descriptions-item>
        <el-descriptions-item label="控制器">{{ currentItem.controller }}</el-descriptions-item>
        <el-descriptions-item label="方法">{{ currentItem.action }}</el-descriptions-item>
        <el-descriptions-item label="执行时间">{{ currentItem.run_time }}ms</el-descriptions-item>
        <el-descriptions-item label="执行时间">{{ currentItem.created_at }}</el-descriptions-item>
        <el-descriptions-item label="请求参数" :span="2">
          <pre style="white-space: pre-wrap; word-break: break-all; max-height: 200px; overflow-y: auto;">{{ currentItem.params || '无' }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="错误信息" :span="2" v-if="currentItem.error">
          <pre style="white-space: pre-wrap; word-break: break-all; color: #f56c6c;">{{ currentItem.error }}</pre>
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
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)

const searchForm = reactive({
  keyword: '',
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])
const detailVisible = ref(false)
const currentItem = ref<any>(null)

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      keyword: searchForm.keyword || undefined,
      status: searchForm.status
    }
    const data = await request.get({ url: '/api/admin/system/run-map', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取运行映射失败:', error)
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
  searchForm.status = undefined
  handleSearch()
}

const handleViewDetail = (row: any) => {
  currentItem.value = row
  detailVisible.value = true
}

const handleRetry = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/system/run-map/${row.id}/retry` })
    ElMessage.success('重试成功')
    fetchData()
  } catch (error) {
    ElMessage.error('重试失败')
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
.run-map-page {
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
