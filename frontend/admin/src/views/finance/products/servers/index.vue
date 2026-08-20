<template>
  <div class="product-servers-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('servers.title') }} - {{ productName }}</span>
          <el-button type="primary" @click="handleAddServer">
            <el-icon><Plus /></el-icon>
            {{ $t('servers.linkServer') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('servers.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('servers.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('servers.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('servers.all')" clearable>
            <el-option :label="$t('servers.online')" value="online" />
            <el-option :label="$t('servers.offline')" value="offline" />
            <el-option :label="$t('servers.maintenance')" value="maintenance" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('servers.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('servers.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="server_name" :label="$t('servers.serverName')" min-width="150" />
        <el-table-column prop="ip_address" :label="$t('servers.ipAddress')" width="150" />
        <el-table-column prop="server_type" :label="$t('servers.type')" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ row.server_type || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('servers.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="getServerStatusTag(row.status)" size="small">
              {{ getServerStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cpu_usage" :label="$t('servers.cpuUsage')" width="120">
          <template #default="{ row }">
            <el-progress
              :percentage="row.cpu_usage || 0"
              :color="getUsageColor(row.cpu_usage)"
              :stroke-width="8"
              :text-inside="false"
            />
          </template>
        </el-table-column>
        <el-table-column prop="memory_usage" :label="$t('servers.memoryUsage')" width="120">
          <template #default="{ row }">
            <el-progress
              :percentage="row.memory_usage || 0"
              :color="getUsageColor(row.memory_usage)"
              :stroke-width="8"
              :text-inside="false"
            />
          </template>
        </el-table-column>
        <el-table-column prop="disk_usage" :label="$t('servers.diskUsage')" width="120">
          <template #default="{ row }">
            <el-progress
              :percentage="row.disk_usage || 0"
              :color="getUsageColor(row.disk_usage)"
              :stroke-width="8"
              :text-inside="false"
            />
          </template>
        </el-table-column>
        <el-table-column prop="max_services" :label="$t('servers.maxServices')" width="100" />
        <el-table-column prop="current_services" :label="$t('servers.currentServices')" width="100" />
        <el-table-column prop="linked_at" :label="$t('servers.linkedAt')" width="180" />
        <el-table-column :label="$t('servers.action')" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEditConfig(row)">{{ $t('servers.config') }}</el-button>
            <el-popconfirm :title="$t('servers.confirmUnlink')" @confirm="handleUnlink(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('servers.unlink') }}</el-button>
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

    <!-- 关联服务器对话框 -->
    <el-dialog v-model="linkDialogVisible" :title="$t('servers.linkServer')" width="600px">
      <el-transfer
        v-model="selectedServers"
        :data="availableServers"
        :titles="[$t('servers.availableServers'), $t('servers.selectedServers')]"
        :props="{ key: 'id', label: 'name' }"
        filterable
        :filter-placeholder="$t('servers.searchServer')"
      />
      <template #footer>
        <el-button @click="linkDialogVisible = false">{{ $t('servers.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmitLink" :loading="submitLoading">{{ $t('servers.confirmLink') }}</el-button>
      </template>
    </el-dialog>

    <!-- 服务器配置对话框 -->
    <el-dialog v-model="configDialogVisible" :title="$t('servers.serverConfig')" width="500px">
      <el-form :model="configForm" ref="configFormRef" label-width="120px">
        <el-form-item :label="$t('servers.serverName')">
          <el-input :value="configForm.server_name" disabled />
        </el-form-item>
        <el-form-item :label="$t('servers.maxServices')" prop="max_services">
          <el-input-number v-model="configForm.max_services" :min="1" :max="10000" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('servers.priority')" prop="priority">
          <el-input-number v-model="configForm.priority" :min="0" :max="100" style="width: 100%" />
          <div class="form-tip">{{ $t('servers.priorityTip') }}</div>
        </el-form-item>
        <el-form-item :label="$t('servers.autoAssign')" prop="auto_assign">
          <el-switch v-model="configForm.auto_assign" :active-text="$t('servers.on')" :inactive-text="$t('servers.off')" />
        </el-form-item>
        <el-form-item :label="$t('servers.weight')" prop="weight">
          <el-input-number v-model="configForm.weight" :min="1" :max="100" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="configDialogVisible = false">{{ $t('servers.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmitConfig" :loading="submitLoading">{{ $t('servers.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'ProductServers' })

const route = useRoute()

const loading = ref(false)
const submitLoading = ref(false)

const productId = ref(route.params.id as string)
const productName = ref('')

const searchForm = reactive({
  keyword: '',
  status: undefined as string | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref([])

const linkDialogVisible = ref(false)
const selectedServers = ref<number[]>([])
const availableServers = ref<any[]>([])

const configDialogVisible = ref(false)
const configFormRef = ref<FormInstance>()
const configForm = reactive({
  id: 0,
  server_name: '',
  max_services: 100,
  priority: 0,
  auto_assign: true,
  weight: 10
})

const getServerStatusTag = (status: string) => {
  const map: Record<string, any> = {
    online: 'success',
    offline: 'danger',
    maintenance: 'warning'
  }
  return map[status] || 'info'
}

const getServerStatusText = (status: string) => {
  const map: Record<string, string> = {
    online: $t('servers.online'),
    offline: $t('servers.offline'),
    maintenance: $t('servers.maintenance')
  }
  return map[status] || $t('servers.unknown')
}

const getUsageColor = (percentage: number) => {
  if (percentage >= 90) return '#f56c6c'
  if (percentage >= 70) return '#e6a23c'
  return '#67c23a'
}

const fetchProductInfo = async () => {
  try {
    const data = await request.get({
      url: `/api/admin/products/${productId.value}`
    })
    productName.value = data.name || $t('servers.unknownProduct')
  } catch (error) {
    console.error('获取产品信息失败:', error)
  }
}

const fetchServers = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: `/api/admin/products/${productId.value}/servers`,
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取关联服务器失败:', error)
    ElMessage.error($t('servers.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchServers()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = undefined
  handleSearch()
}

const handleAddServer = async () => {
  try {
    const data = await request.get({
      url: `/api/admin/products/${productId.value}/available-servers`
    })
    availableServers.value = (data || []).map((s: any) => ({
      id: s.id,
      name: `${s.name} (${s.ip_address})`
    }))
    selectedServers.value = []
    linkDialogVisible.value = true
  } catch (error) {
    ElMessage.error($t('servers.fetchAvailableFailed'))
  }
}

const handleSubmitLink = async () => {
  if (!selectedServers.value.length) {
    ElMessage.warning($t('servers.selectServerWarning'))
    return
  }

  submitLoading.value = true
  try {
    await request.post({
      url: `/api/admin/products/${productId.value}/servers`,
      params: {
        server_ids: selectedServers.value
      },
      showSuccessMessage: true
    })
    linkDialogVisible.value = false
    fetchServers()
  } catch (error) {
    ElMessage.error($t('servers.linkFailed'))
  } finally {
    submitLoading.value = false
  }
}

const handleEditConfig = (row: any) => {
  configForm.id = row.id
  configForm.server_name = row.server_name
  configForm.max_services = row.max_services || 100
  configForm.priority = row.priority || 0
  configForm.auto_assign = row.auto_assign !== false
  configForm.weight = row.weight || 10
  configDialogVisible.value = true
}

const handleSubmitConfig = async () => {
  submitLoading.value = true
  try {
    await request.put({
      url: `/api/admin/products/${productId.value}/servers/${configForm.id}`,
      params: {
        max_services: configForm.max_services,
        priority: configForm.priority,
        auto_assign: configForm.auto_assign,
        weight: configForm.weight
      },
      showSuccessMessage: true
    })
    configDialogVisible.value = false
    fetchServers()
  } catch (error) {
    ElMessage.error($t('servers.saveFailed'))
  } finally {
    submitLoading.value = false
  }
}

const handleUnlink = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/products/${productId.value}/servers/${row.id}`
    })
    ElMessage.success($t('servers.unlinkSuccess'))
    fetchServers()
  } catch (error) {
    ElMessage.error($t('servers.unlinkFailed'))
  }
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchServers()
}

const handlePageChange = () => {
  fetchServers()
}

onMounted(() => {
  fetchProductInfo()
  fetchServers()
})
</script>

<style scoped lang="scss">
.product-servers-page {
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

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

:deep(.el-transfer) {
  display: flex;
  justify-content: center;
}
</style>
