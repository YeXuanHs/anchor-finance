<template>
  <div class="product-servers-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>产品关联服务器 - {{ productName }}</span>
          <el-button type="primary" @click="handleAddServer">
            <el-icon><Plus /></el-icon>
            关联服务器
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="服务器名称/IP" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="在线" value="online" />
            <el-option label="离线" value="offline" />
            <el-option label="维护中" value="maintenance" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="server_name" label="服务器名称" min-width="150" />
        <el-table-column prop="ip_address" label="IP地址" width="150" />
        <el-table-column prop="server_type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ row.server_type || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getServerStatusTag(row.status)" size="small">
              {{ getServerStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cpu_usage" label="CPU使用率" width="120">
          <template #default="{ row }">
            <el-progress
              :percentage="row.cpu_usage || 0"
              :color="getUsageColor(row.cpu_usage)"
              :stroke-width="8"
              :text-inside="false"
            />
          </template>
        </el-table-column>
        <el-table-column prop="memory_usage" label="内存使用率" width="120">
          <template #default="{ row }">
            <el-progress
              :percentage="row.memory_usage || 0"
              :color="getUsageColor(row.memory_usage)"
              :stroke-width="8"
              :text-inside="false"
            />
          </template>
        </el-table-column>
        <el-table-column prop="disk_usage" label="磁盘使用率" width="120">
          <template #default="{ row }">
            <el-progress
              :percentage="row.disk_usage || 0"
              :color="getUsageColor(row.disk_usage)"
              :stroke-width="8"
              :text-inside="false"
            />
          </template>
        </el-table-column>
        <el-table-column prop="max_services" label="最大服务数" width="100" />
        <el-table-column prop="current_services" label="当前服务数" width="100" />
        <el-table-column prop="linked_at" label="关联时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEditConfig(row)">配置</el-button>
            <el-popconfirm title="确定取消关联该服务器吗？" @confirm="handleUnlink(row)">
              <template #reference>
                <el-button type="danger" link>取消关联</el-button>
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
    <el-dialog v-model="linkDialogVisible" title="关联服务器" width="600px">
      <el-transfer
        v-model="selectedServers"
        :data="availableServers"
        :titles="['可选服务器', '已选服务器']"
        :props="{ key: 'id', label: 'name' }"
        filterable
        filter-placeholder="搜索服务器"
      />
      <template #footer>
        <el-button @click="linkDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitLink" :loading="submitLoading">确定关联</el-button>
      </template>
    </el-dialog>

    <!-- 服务器配置对话框 -->
    <el-dialog v-model="configDialogVisible" title="服务器配置" width="500px">
      <el-form :model="configForm" ref="configFormRef" label-width="120px">
        <el-form-item label="服务器名称">
          <el-input :value="configForm.server_name" disabled />
        </el-form-item>
        <el-form-item label="最大服务数" prop="max_services">
          <el-input-number v-model="configForm.max_services" :min="1" :max="10000" style="width: 100%" />
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-input-number v-model="configForm.priority" :min="0" :max="100" style="width: 100%" />
          <div class="form-tip">数值越大优先级越高</div>
        </el-form-item>
        <el-form-item label="自动分配" prop="auto_assign">
          <el-switch v-model="configForm.auto_assign" active-text="开启" inactive-text="关闭" />
        </el-form-item>
        <el-form-item label="权重" prop="weight">
          <el-input-number v-model="configForm.weight" :min="1" :max="100" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="configDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitConfig" :loading="submitLoading">保存</el-button>
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

defineOptions({ name: 'ProductServers' })

const route = useRoute()

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)

// 产品信息
const productId = ref(route.params.id as string)
const productName = ref('')

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: undefined as string | undefined
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref([])

// 关联服务器对话框
const linkDialogVisible = ref(false)
const selectedServers = ref<number[]>([])
const availableServers = ref<any[]>([])

// 配置对话框
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

// 获取服务器状态标签
const getServerStatusTag = (status: string) => {
  const map: Record<string, any> = {
    online: 'success',
    offline: 'danger',
    maintenance: 'warning'
  }
  return map[status] || 'info'
}

// 获取服务器状态文本
const getServerStatusText = (status: string) => {
  const map: Record<string, string> = {
    online: '在线',
    offline: '离线',
    maintenance: '维护中'
  }
  return map[status] || '未知'
}

// 获取使用率颜色
const getUsageColor = (percentage: number) => {
  if (percentage >= 90) return '#f56c6c'
  if (percentage >= 70) return '#e6a23c'
  return '#67c23a'
}

// 获取产品信息
const fetchProductInfo = async () => {
  try {
    const data = await request.get({
      url: `/api/admin/products/${productId.value}`
    })
    productName.value = data.name || '未知产品'
  } catch (error) {
    console.error('获取产品信息失败:', error)
  }
}

// 获取关联服务器列表
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
    ElMessage.error('获取关联服务器失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchServers()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = undefined
  handleSearch()
}

// 关联服务器
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
    ElMessage.error('获取可选服务器失败')
  }
}

// 提交关联
const handleSubmitLink = async () => {
  if (!selectedServers.value.length) {
    ElMessage.warning('请选择要关联的服务器')
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
    ElMessage.error('关联失败')
  } finally {
    submitLoading.value = false
  }
}

// 编辑配置
const handleEditConfig = (row: any) => {
  configForm.id = row.id
  configForm.server_name = row.server_name
  configForm.max_services = row.max_services || 100
  configForm.priority = row.priority || 0
  configForm.auto_assign = row.auto_assign !== false
  configForm.weight = row.weight || 10
  configDialogVisible.value = true
}

// 提交配置
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
    ElMessage.error('保存失败')
  } finally {
    submitLoading.value = false
  }
}

// 取消关联
const handleUnlink = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/products/${productId.value}/servers/${row.id}`
    })
    ElMessage.success('取消关联成功')
    fetchServers()
  } catch (error) {
    ElMessage.error('取消关联失败')
  }
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchServers()
}

// 页码变化
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
