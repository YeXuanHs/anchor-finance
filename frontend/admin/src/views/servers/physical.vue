<template>
  <div class="physical-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="服务器名称/IP" clearable />
        </el-form-item>
        <el-form-item label="机房">
          <el-select v-model="searchForm.datacenter_id" placeholder="全部" clearable>
            <el-option v-for="dc in datacenterOptions" :key="dc.id" :label="dc.name" :value="dc.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="运行中" value="online" />
            <el-option label="已关机" value="offline" />
            <el-option label="故障" value="error" />
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
        <h3>物理服务器管理</h3>
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          新增服务器
        </el-button>
      </div>

      <el-table :data="list" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP地址" width="140" />
        <el-table-column prop="cpu" label="CPU" width="140" show-overflow-tooltip />
        <el-table-column prop="memory" label="内存" width="100" />
        <el-table-column prop="disk" label="硬盘" width="100" />
        <el-table-column prop="bandwidth" label="带宽" width="90" />
        <el-table-column prop="datacenter_name" label="机房" width="120" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" size="small">
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button
              :type="row.status === 'online' ? 'warning' : 'success'"
              link
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 'online' ? '关机' : '开机' }}
            </el-button>
            <el-button type="info" link @click="handleMonitor(row)">监控</el-button>
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
      :title="isEdit ? '编辑服务器' : '新增服务器'"
      width="650px"
      destroy-on-close
    >
      <el-form :model="formData" :rules="rules" ref="formRef" label-width="100px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="名称" prop="name">
              <el-input v-model="formData.name" placeholder="服务器名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="IP地址" prop="ip">
              <el-input v-model="formData.ip" placeholder="如: 192.168.1.1" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="CPU" prop="cpu">
              <el-input v-model="formData.cpu" placeholder="如: Intel Xeon E5-2680" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="内存" prop="memory">
              <el-input v-model="formData.memory" placeholder="如: 64GB DDR4" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="硬盘" prop="disk">
              <el-input v-model="formData.disk" placeholder="如: 2TB SSD" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="带宽" prop="bandwidth">
              <el-input v-model="formData.bandwidth" placeholder="如: 100Mbps" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="机房" prop="datacenter_id">
              <el-select v-model="formData.datacenter_id" placeholder="请选择机房" style="width: 100%">
                <el-option v-for="dc in datacenterOptions" :key="dc.id" :label="dc.name" :value="dc.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="formData.status" placeholder="请选择状态" style="width: 100%">
                <el-option label="运行中" value="online" />
                <el-option label="已关机" value="offline" />
                <el-option label="故障" value="error" />
                <el-option label="维护中" value="maintenance" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注">
          <el-input v-model="formData.remark" type="textarea" :rows="2" placeholder="服务器备注" />
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
import { ref, onMounted, reactive } from 'vue'
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

const searchForm = ref({ keyword: '', status: '', datacenter_id: '' })

const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

const datacenterOptions = ref<any[]>([])

const formData = reactive({
  name: '',
  ip: '',
  cpu: '',
  memory: '',
  disk: '',
  bandwidth: '',
  datacenter_id: null as number | null,
  status: 'offline',
  remark: ''
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入服务器名称', trigger: 'blur' }],
  ip: [{ required: true, message: '请输入IP地址', trigger: 'blur' }],
  datacenter_id: [{ required: true, message: '请选择机房', trigger: 'change' }]
}

const statusLabel = (status: string) => {
  const map: Record<string, string> = { online: '运行中', offline: '已关机', error: '故障', maintenance: '维护中' }
  return map[status] || status
}

const statusTagType = (status: string) => {
  const map: Record<string, string> = { online: 'success', offline: 'info', error: 'danger', maintenance: 'warning' }
  return (map[status] as any) || 'info'
}

const fetchList = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/servers/physical', {
      params: { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    })
    list.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

const fetchDatacenters = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/servers/datacenters', { params: { page_size: 999 } })
    datacenterOptions.value = data.data?.list || []
  } catch { /* ignore */ }
}

const handleSearch = () => { currentPage.value = 1; fetchList() }

const resetSearch = () => {
  searchForm.value = { keyword: '', status: '', datacenter_id: '' }
  handleSearch()
}

const resetForm = () => {
  Object.assign(formData, { name: '', ip: '', cpu: '', memory: '', disk: '', bandwidth: '', datacenter_id: null, status: 'offline', remark: '' })
  editingId.value = null
}

const handleAdd = () => { isEdit.value = false; resetForm(); dialogVisible.value = true }

const handleEdit = (row: any) => {
  isEdit.value = true
  editingId.value = row.id
  Object.assign(formData, {
    name: row.name, ip: row.ip, cpu: row.cpu, memory: row.memory,
    disk: row.disk, bandwidth: row.bandwidth || '', datacenter_id: row.datacenter_id,
    status: row.status, remark: row.remark || ''
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (isEdit.value && editingId.value) {
      await request.put(`/admin/api/v1/servers/physical/${editingId.value}`, formData)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/servers/physical', formData)
      ElMessage.success('新增成功')
    }
    dialogVisible.value = false
    fetchList()
  } catch { /* handled */ } finally {
    submitLoading.value = false
  }
}

const handleToggleStatus = async (row: any) => {
  const newStatus = row.status === 'online' ? 'offline' : 'online'
  const action = newStatus === 'online' ? '开机' : '关机'
  await ElMessageBox.confirm(`确定要${action}服务器「${row.name}」吗？`, '提示', { type: 'warning' })
  try {
    await request.put(`/admin/api/v1/servers/physical/${row.id}/status`, { status: newStatus })
    ElMessage.success(`${action}成功`)
    fetchList()
  } catch { /* handled */ }
}

const handleMonitor = (row: any) => {
  ElMessage.info(`查看服务器 ${row.name} 的监控数据`)
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm(`确定要删除服务器「${row.name}」吗？此操作不可恢复。`, '警告', { type: 'error' })
  try {
    await request.delete(`/admin/api/v1/servers/physical/${row.id}`)
    ElMessage.success('删除成功')
    fetchList()
  } catch { /* handled */ }
}

onMounted(() => { fetchList(); fetchDatacenters() })
</script>

<style scoped lang="scss">
.physical-page {
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
