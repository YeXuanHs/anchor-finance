<template>
  <div class="upstream-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="供应商名称" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="物理机" value="physical" />
            <el-option label="云服务器" value="cloud" />
            <el-option label="VPS" value="vps" />
            <el-option label="混合" value="hybrid" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="正常" value="active" />
            <el-option label="异常" value="error" />
            <el-option label="维护中" value="maintenance" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>上游供应商管理</h3>
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          添加供应商
        </el-button>
      </div>

      <el-table :data="list" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="名称" min-width="130" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ typeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="api_url" label="API地址" min-width="180" show-overflow-tooltip />
        <el-table-column prop="server_count" label="服务器数" width="100" />
        <el-table-column prop="sync_status" label="同步状态" width="110">
          <template #default="{ row }">
            <el-tag :type="syncStatusType(row.sync_status)" size="small">
              {{ syncStatusLabel(row.sync_status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : row.status === 'error' ? 'danger' : 'warning'" size="small">
              {{ row.status === 'active' ? '正常' : row.status === 'error' ? '异常' : '维护中' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_sync" label="最后同步" width="170" />
        <el-table-column label="操作" width="230" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="success" link :loading="row._syncing" @click="handleSync(row)">同步</el-button>
            <el-button type="info" link @click="handleTestApi(row)">测试API</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑供应商' : '添加供应商'"
      width="650px"
      destroy-on-close
    >
      <el-form :model="formData" :rules="rules" ref="formRef" label-width="110px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" placeholder="供应商名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型" style="width: 100%">
            <el-option label="物理机" value="physical" />
            <el-option label="云服务器" value="cloud" />
            <el-option label="VPS" value="vps" />
            <el-option label="混合" value="hybrid" />
          </el-select>
        </el-form-item>
        <el-divider content-position="left">API配置</el-divider>
        <el-form-item label="API地址" prop="api_url">
          <el-input v-model="formData.api_url" placeholder="https://api.example.com" />
        </el-form-item>
        <el-form-item label="API密钥" prop="api_key">
          <el-input v-model="formData.api_key" placeholder="请输入API密钥" show-password />
        </el-form-item>
        <el-form-item label="API协议">
          <el-select v-model="formData.api_protocol" placeholder="选择协议" style="width: 100%">
            <el-option label="RESTful" value="rest" />
            <el-option label="SOAP" value="soap" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="同步间隔">
          <el-input-number v-model="formData.sync_interval" :min="5" :max="1440" :step="5" />
          <span style="margin-left: 8px; color: #909399">分钟</span>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="formData.remark" type="textarea" :rows="2" placeholder="备注信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const list = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const searchForm = ref({ keyword: '', type: '', status: '' })

const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

const formData = reactive({
  name: '',
  type: 'cloud',
  api_url: '',
  api_key: '',
  api_protocol: 'rest',
  sync_interval: 30,
  remark: ''
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入供应商名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  api_url: [{ required: true, message: '请输入API地址', trigger: 'blur' }]
}

const typeLabel = (type: string) => {
  const map: Record<string, string> = { physical: '物理机', cloud: '云服务器', vps: 'VPS', hybrid: '混合' }
  return map[type] || type
}

const syncStatusLabel = (status: string) => {
  const map: Record<string, string> = { success: '同步成功', failed: '同步失败', syncing: '同步中', pending: '未同步' }
  return map[status] || status || '未同步'
}

const syncStatusType = (status: string) => {
  const map: Record<string, string> = { success: 'success', failed: 'danger', syncing: 'warning', pending: 'info' }
  return (map[status] as any) || 'info'
}

const fetchList = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/servers/upstream', {
      params: { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    })
    list.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

const handleSearch = () => { currentPage.value = 1; fetchList() }

const resetSearch = () => {
  searchForm.value = { keyword: '', type: '', status: '' }
  handleSearch()
}

const resetForm = () => {
  Object.assign(formData, { name: '', type: 'cloud', api_url: '', api_key: '', api_protocol: 'rest', sync_interval: 30, remark: '' })
  editingId.value = null
}

const handleAdd = () => { isEdit.value = false; resetForm(); dialogVisible.value = true }

const handleEdit = (row: any) => {
  isEdit.value = true
  editingId.value = row.id
  Object.assign(formData, {
    name: row.name, type: row.type, api_url: row.api_url, api_key: row.api_key || '',
    api_protocol: row.api_protocol || 'rest', sync_interval: row.sync_interval || 30, remark: row.remark || ''
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (isEdit.value && editingId.value) {
      await request.put(`/admin/api/v1/servers/upstream/${editingId.value}`, formData)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/servers/upstream', formData)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    fetchList()
  } catch { /* handled */ } finally {
    submitLoading.value = false
  }
}

const handleSync = async (row: any) => {
  row._syncing = true
  try {
    await request.post(`/admin/api/v1/servers/upstream/${row.id}/sync`)
    ElMessage.success('同步任务已触发')
    fetchList()
  } catch { /* handled */ } finally {
    row._syncing = false
  }
}

const handleTestApi = async (row: any) => {
  try {
    const { data } = await request.post(`/admin/api/v1/servers/upstream/${row.id}/test`)
    ElMessage.success(data.message || 'API连接正常')
  } catch {
    ElMessage.error('API连接失败')
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm(`确定要删除供应商「${row.name}」吗？关联的服务器不会被删除。`, '警告', { type: 'error' })
  try {
    await request.delete(`/admin/api/v1/servers/upstream/${row.id}`)
    ElMessage.success('删除成功')
    fetchList()
  } catch { /* handled */ }
}

onMounted(fetchList)
</script>

<style scoped lang="scss">
.upstream-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
}
</style>
