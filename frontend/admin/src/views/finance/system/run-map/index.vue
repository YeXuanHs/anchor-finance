<template>
  <div class="run-map-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('runMap.title') }}</span>
          <el-button @click="fetchData">
            <el-icon><Refresh /></el-icon>
            {{ $t('runMap.refresh') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('runMap.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('runMap.routePlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('runMap.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('runMap.all')" clearable>
            <el-option :label="$t('runMap.success')" :value="1" />
            <el-option :label="$t('runMap.failed')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('runMap.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('runMap.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="route" :label="$t('runMap.route')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="controller" :label="$t('runMap.controller')" min-width="150" />
        <el-table-column prop="action" :label="$t('runMap.action')" width="120" />
        <el-table-column prop="status" :label="$t('runMap.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('runMap.success') : $t('runMap.failed') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="run_time" :label="$t('runMap.runTime')" width="120" align="center" />
        <el-table-column prop="created_at" :label="$t('runMap.createdAt')" width="170" />
        <el-table-column :label="$t('runMap.operation')" width="150" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">{{ $t('runMap.detail') }}</el-button>
            <el-button v-if="row.status === 0" type="warning" link @click="handleRetry(row)">{{ $t('runMap.retry') }}</el-button>
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
    <el-dialog v-model="detailVisible" :title="$t('runMap.runDetail')" width="700px">
      <el-descriptions :column="2" border v-if="currentItem">
        <el-descriptions-item label="ID">{{ currentItem.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('runMap.status')">
          <el-tag :type="currentItem.status === 1 ? 'success' : 'danger'" size="small">
            {{ currentItem.status === 1 ? $t('runMap.success') : $t('runMap.failed') }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('runMap.route')" :span="2">{{ currentItem.route }}</el-descriptions-item>
        <el-descriptions-item :label="$t('runMap.controller')">{{ currentItem.controller }}</el-descriptions-item>
        <el-descriptions-item :label="$t('runMap.action')">{{ currentItem.action }}</el-descriptions-item>
        <el-descriptions-item :label="$t('runMap.runTime')">{{ currentItem.run_time }}ms</el-descriptions-item>
        <el-descriptions-item :label="$t('runMap.createdAt')">{{ currentItem.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('runMap.params')" :span="2">
          <pre style="white-space: pre-wrap; word-break: break-all; max-height: 200px; overflow-y: auto;">{{ currentItem.params || $t('runMap.none') }}</pre>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('runMap.errorInfo')" :span="2" v-if="currentItem.error">
          <pre style="white-space: pre-wrap; word-break: break-all; color: #f56c6c;">{{ currentItem.error }}</pre>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">{{ $t('runMap.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { $t } from '@/locales'
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
    const data = await request.get({ url: '/api/admin/run-map', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取运行映射失败:', error)
    ElMessage.error($t('runMap.fetchFailed'))
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
    await request.post({ url: `/api/admin/run-map/${row.id}/repeat` })
    ElMessage.success($t('runMap.retrySuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('runMap.retryFailed'))
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
