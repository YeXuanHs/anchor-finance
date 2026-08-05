<template>
  <div class="upstream-providers-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>上游供应商管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增供应商
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="名称">
          <el-input v-model="searchForm.name" placeholder="供应商名称" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="ProxmoxVE" value="pve" />
            <el-option label="VMware" value="vmware" />
            <el-option label="HyperV" value="hyperv" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.is_enabled" placeholder="全部" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="供应商名称" min-width="150" />
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="url" label="API地址" min-width="200" show-overflow-tooltip />
        <el-table-column prop="products_count" label="产品数" width="80" align="center" />
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_enabled ? 'success' : 'info'" size="small">{{ row.is_enabled ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_sync" label="最后同步" width="180" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="success" link @click="handleTest(row)">测试连接</el-button>
            <el-button type="warning" link @click="handleSync(row)">同步</el-button>
            <el-button type="info" link @click="handleViewLog(row)">日志</el-button>
            <el-popconfirm title="确定删除该供应商吗？" @confirm="handleDelete(row)">
              <template #reference><el-button type="danger" link>删除</el-button></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item label="供应商名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入供应商名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型" style="width: 100%">
            <el-option label="ProxmoxVE" value="pve" />
            <el-option label="VMware" value="vmware" />
            <el-option label="HyperV" value="hyperv" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="API地址" prop="url">
          <el-input v-model="formData.url" placeholder="https://example.com:8006/api2/json" />
        </el-form-item>
        <el-form-item label="端口">
          <el-input-number v-model="formData.port" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="formData.username" placeholder="API用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="formData.password" type="password" show-password placeholder="API密码" />
        </el-form-item>
        <el-form-item label="API密钥">
          <el-input v-model="formData.api_key" type="password" show-password placeholder="API密钥（可选）" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="formData.is_enabled" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="formData.remark" type="textarea" :rows="2" placeholder="备注信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 日志对话框 -->
    <el-dialog v-model="logDialogVisible" title="同步日志" width="800px">
      <el-table :data="logData" v-loading="logLoading" border size="small">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="action" label="操作" width="120">
          <template #default="{ row }">
            <el-tag :type="getLogActionType(row.action)" size="small">{{ row.action }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="message" label="消息" min-width="300" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">{{ row.status === 'success' ? '成功' : '失败' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="180" />
      </el-table>
      <div class="pagination-container">
        <el-pagination v-model:current-page="logPagination.page" v-model:page-size="logPagination.page_size"
          :total="logPagination.total" layout="total, prev, pager, next" small
          @current-change="fetchProviderLog" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const typeLabels: Record<string, string> = { pve: 'ProxmoxVE', vmware: 'VMware', hyperv: 'HyperV', custom: '自定义' }
const getTypeLabel = (type: string) => typeLabels[type] || type
const getLogActionType = (action: string) => {
  const map: Record<string, string> = { sync: 'primary', test: 'success', error: 'danger' }
  return (map[action] || 'info') as any
}

const loading = ref(false)
const submitLoading = ref(false)
const logLoading = ref(false)

const searchForm = reactive({ name: '', type: '', is_enabled: undefined as number | undefined })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref('新增供应商')
const formRef = ref<FormInstance>()
const formData = reactive({
  id: null as number | null,
  name: '',
  type: 'pve',
  url: '',
  port: 8006,
  username: '',
  password: '',
  api_key: '',
  is_enabled: 1,
  remark: ''
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入供应商名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  url: [{ required: true, message: '请输入API地址', trigger: 'blur' }]
}

const logDialogVisible = ref(false)
const logData = ref<any[]>([])
const logProviderId = ref(0)
const logPagination = reactive({ page: 1, page_size: 20, total: 0 })

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (searchForm.name) params.name = searchForm.name
    if (searchForm.type) params.type = searchForm.type
    if (searchForm.is_enabled !== undefined) params.is_enabled = searchForm.is_enabled
    const res = await request.get({ url: '/api/admin/upstream/providers', params })
    tableData.value = res.data || res.list || res || []
    pagination.total = res.total || 0
  } catch { ElMessage.error('获取供应商列表失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.name = ''; searchForm.type = ''; searchForm.is_enabled = undefined; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleAdd = () => {
  dialogTitle.value = '新增供应商'
  Object.assign(formData, { id: null, name: '', type: 'pve', url: '', port: 8006, username: '', password: '', api_key: '', is_enabled: 1, remark: '' })
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑供应商'
  Object.assign(formData, { id: row.id, name: row.name, type: row.type, url: row.url, port: row.port || 8006, username: row.username || '', password: '', api_key: '', is_enabled: row.is_enabled, remark: row.remark || '' })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/upstream/providers/${formData.id}`, data: formData, showSuccessMessage: true })
      } else {
        await request.post({ url: '/api/admin/upstream/providers', data: formData, showSuccessMessage: true })
      }
      dialogVisible.value = false; fetchData()
    } catch { ElMessage.error('操作失败') } finally { submitLoading.value = false }
  })
}

const handleTest = async (row: any) => {
  try {
    const res = await request.post({ url: `/api/admin/upstream/providers/${row.id}/test` })
    ElMessage.success(res?.message || '连接测试成功')
  } catch { ElMessage.error('连接测试失败') }
}

const handleSync = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定同步供应商 "${row.name}" 的数据吗？`, '同步确认')
    const res = await request.post({ url: `/api/admin/upstream/providers/${row.id}/sync` })
    ElMessage.success(res?.message || '同步任务已提交')
    fetchData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error('同步失败') }
}

const handleViewLog = async (row: any) => {
  logProviderId.value = row.id
  logPagination.page = 1
  logDialogVisible.value = true
  fetchProviderLog()
}

const fetchProviderLog = async () => {
  logLoading.value = true
  try {
    const data = await request.get({ url: `/api/admin/upstream/providers/${logProviderId.value}/logs`, params: { page: logPagination.page, page_size: logPagination.page_size } })
    logData.value = data.list || data || []
    logPagination.total = data.total || 0
  } catch { ElMessage.error('获取日志失败') } finally { logLoading.value = false }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/upstream/providers/${row.id}` })
    ElMessage.success('删除成功'); fetchData()
  } catch { ElMessage.error('删除失败') }
}

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.upstream-providers-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 16px; .el-form-item { margin-bottom: 0; } }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
