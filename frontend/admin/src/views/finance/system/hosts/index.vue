<template>
  <div class="hosts-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('hostManage.title') }}</span>
          <el-button type="primary" @click="handleExport">
            <el-icon><Download /></el-icon>
            {{ $t('hostManage.export') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('hostManage.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('hostManage.domainPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('hostManage.userId')">
          <el-input v-model="searchForm.user_id" :placeholder="$t('hostManage.userId')" clearable />
        </el-form-item>
        <el-form-item :label="$t('hostManage.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('hostManage.all')" clearable>
            <el-option :label="$t('hostManage.active')" value="active" />
            <el-option :label="$t('hostManage.suspended')" value="suspended" />
            <el-option :label="$t('hostManage.pending')" value="pending" />
            <el-option :label="$t('hostManage.terminated')" value="terminated" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('hostManage.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('hostManage.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="user_id" :label="$t('hostManage.userId')" width="90" align="center" />
        <el-table-column prop="username" :label="$t('hostManage.username')" width="120" show-overflow-tooltip />
        <el-table-column prop="product_name" :label="$t('hostManage.productName')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="domain" :label="$t('hostManage.domainIp')" min-width="180" show-overflow-tooltip />
        <el-table-column prop="status" :label="$t('hostManage.status')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('hostManage.createdAt')" width="170" />
        <el-table-column prop="due_date" :label="$t('hostManage.dueDate')" width="170">
          <template #default="{ row }">
            <span :class="{ 'danger-text': isExpiringSoon(row.due_date) }">
              {{ row.due_date || '-' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('hostManage.operation')" width="220" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">{{ $t('hostManage.detail') }}</el-button>
            <el-popconfirm
              v-if="row.status === 'active'"
              :title="$t('hostManage.confirmSuspend')"
              @confirm="handleSuspend(row)"
            >
              <template #reference>
                <el-button type="warning" link>{{ $t('hostManage.suspend') }}</el-button>
              </template>
            </el-popconfirm>
            <el-button
              v-if="row.status === 'suspended'"
              type="success"
              link
              @click="handleActivate(row)"
            >
              {{ $t('hostManage.activate') }}
            </el-button>
            <el-popconfirm
              v-if="row.status !== 'terminated'"
              :title="$t('hostManage.confirmTerminate')"
              @confirm="handleTerminate(row)"
            >
              <template #reference>
                <el-button type="danger" link>{{ $t('hostManage.terminate') }}</el-button>
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

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" :title="$t('hostManage.hostDetail')" width="700px" destroy-on-close>
      <el-descriptions :column="2" border>
        <el-descriptions-item :label="$t('hostManage.hostId')">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('hostManage.userId')">{{ detailData.user_id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('hostManage.username')">{{ detailData.username || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('hostManage.productName')">{{ detailData.product_name || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('hostManage.domainIp')" :span="2">{{ detailData.domain || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('hostManage.status')">
          <el-tag :type="getStatusType(detailData.status)" size="small">
            {{ getStatusText(detailData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('hostManage.ipAddress')">{{ detailData.ip || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('hostManage.createdAt')">{{ detailData.created_at || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('hostManage.dueDate')">{{ detailData.due_date || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('hostManage.remark')" :span="2">{{ detailData.remark || $t('hostManage.none') }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">{{ $t('hostManage.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'
import { exportToCSV } from '@/utils/export'

const loading = ref(false)
const detailVisible = ref(false)
const detailData = ref<any>({})

const searchForm = reactive({
  keyword: '',
  user_id: '',
  status: '' as string
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])

const getStatusMap = (): Record<string, { text: string; type: string }> => ({
  active: { text: $t('hostManage.active'), type: 'success' },
  suspended: { text: $t('hostManage.suspended'), type: 'warning' },
  pending: { text: $t('hostManage.pending'), type: 'info' },
  terminated: { text: $t('hostManage.terminated'), type: 'danger' }
})

const getStatusText = (status: string) => getStatusMap()[status]?.text || status

const getStatusType = (status: string) => (getStatusMap()[status]?.type || 'info') as any

const isExpiringSoon = (dueDate: string) => {
  if (!dueDate) return false
  const diff = new Date(dueDate).getTime() - Date.now()
  return diff > 0 && diff < 7 * 24 * 60 * 60 * 1000
}

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/hosts',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        keyword: searchForm.keyword || undefined,
        user_id: searchForm.user_id || undefined,
        status: searchForm.status || undefined
      }
    })
    tableData.value = data.list || data || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取主机列表失败:', error)
    ElMessage.error($t('hostManage.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.user_id = ''
  searchForm.status = ''
  handleSearch()
}

const handleView = (row: any) => {
  detailData.value = { ...row }
  detailVisible.value = true
}

const handleSuspend = async (row: any) => {
  try {
    await request.put({ url: `/api/admin/hosts/${row.id}/status`, params: { status: 'suspended' } })
    ElMessage.success($t('hostManage.hostSuspended'))
    fetchList()
  } catch (error) {
    ElMessage.error($t('hostManage.operationFailed'))
  }
}

const handleActivate = async (row: any) => {
  try {
    await request.put({ url: `/api/admin/hosts/${row.id}/status`, params: { status: 'active' } })
    ElMessage.success($t('hostManage.hostActivated'))
    fetchList()
  } catch (error) {
    ElMessage.error($t('hostManage.operationFailed'))
  }
}

const handleTerminate = async (row: any) => {
  try {
    await request.put({ url: `/api/admin/hosts/${row.id}/status`, params: { status: 'terminated' } })
    ElMessage.success($t('hostManage.hostTerminated'))
    fetchList()
  } catch (error) {
    ElMessage.error($t('hostManage.operationFailed'))
  }
}

const handleExport = async () => {
  try {
    const data = await request.get({ url: '/api/admin/hosts', params: { page: 1, page_size: 9999 } })
    const list = data.list || data || []
    exportToCSV(list, [
      { key: 'id', title: 'ID' },
      { key: 'hostname', title: $t('hostManage.exportHostname') },
      { key: 'ip', title: 'IP' },
      { key: 'user_id', title: $t('hostManage.exportUserId') },
      { key: 'status', title: $t('hostManage.exportStatus') },
      { key: 'created_at', title: $t('hostManage.exportCreatedAt') }
    ], $t('hostManage.exportFilename'))
    ElMessage.success($t('hostManage.exportSuccess'))
  } catch { ElMessage.error($t('hostManage.exportFailed')) }
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchList()
}

const handlePageChange = () => {
  fetchList()
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.hosts-page {
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

.danger-text {
  color: #f56c6c;
  font-weight: 600;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
