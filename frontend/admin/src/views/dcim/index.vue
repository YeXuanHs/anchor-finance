<template>
  <div class="dcim-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="设备名称/IP/序列号" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="机房">
          <el-select v-model="searchForm.datacenter_id" placeholder="全部" clearable>
            <el-option v-for="dc in datacenters" :key="dc.id" :label="dc.name" :value="dc.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="在线" value="online" />
            <el-option label="离线" value="offline" />
            <el-option label="维护中" value="maintenance" />
            <el-option label="故障" value="fault" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="服务器" value="server" />
            <el-option label="交换机" value="switch" />
            <el-option label="路由器" value="router" />
            <el-option label="防火墙" value="firewall" />
            <el-option label="存储" value="storage" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>搜索
          </el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 统计 -->
    <div class="stats-grid">
      <div class="stat-card" v-for="stat in stats" :key="stat.label">
        <div class="stat-label">{{ stat.label }}</div>
        <div class="stat-value" :class="stat.type">{{ stat.value }}</div>
      </div>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>设备列表</h3>
        <div>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>添加设备
          </el-button>
          <el-button @click="handleBatchAction" :disabled="!selectedDevices.length">
            <el-icon><Operation /></el-icon>批量操作
          </el-button>
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>导出
          </el-button>
        </div>
      </div>

      <el-table :data="devices" style="width: 100%" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="设备名称" min-width="140" />
        <el-table-column prop="sn" label="序列号" width="140" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="90">
          <template #default="{ row }">
            <el-tag size="small">{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="datacenter_name" label="机房" width="100" />
        <el-table-column prop="rack" label="机架" width="90" />
        <el-table-column prop="ip" label="IP地址" width="140" />
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cpu_usage" label="CPU" width="80">
          <template #default="{ row }">
            <el-progress :percentage="row.cpu_usage || 0" :stroke-width="6" :color="getProgressColor(row.cpu_usage)" v-if="row.cpu_usage != null" />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="mem_usage" label="内存" width="80">
          <template #default="{ row }">
            <el-progress :percentage="row.mem_usage || 0" :stroke-width="6" :color="getProgressColor(row.mem_usage)" v-if="row.mem_usage != null" />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-dropdown trigger="click" @command="(cmd: string) => handleRemoteAction(row, cmd)">
              <el-button type="primary" link size="small">远程操作<el-icon class="el-icon--right"><ArrowDown /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="ping">Ping</el-dropdown-item>
                  <el-dropdown-item command="reboot">重启</el-dropdown-item>
                  <el-dropdown-item command="power_off">关机</el-dropdown-item>
                  <el-dropdown-item command="sol">SOL控制台</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="table-footer">
        <div class="batch-info" v-if="selectedDevices.length">已选 {{ selectedDevices.length }} 项</div>
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="650px" :close-on-click-modal="false">
      <el-form :model="formData" :rules="rules" ref="formRef" label-width="100px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="设备名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入设备名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="序列号" prop="sn">
              <el-input v-model="formData.sn" placeholder="请输入序列号" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="类型" prop="type">
              <el-select v-model="formData.type" placeholder="请选择类型" style="width: 100%;">
                <el-option label="服务器" value="server" />
                <el-option label="交换机" value="switch" />
                <el-option label="路由器" value="router" />
                <el-option label="防火墙" value="firewall" />
                <el-option label="存储" value="storage" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="机房" prop="datacenter_id">
              <el-select v-model="formData.datacenter_id" placeholder="请选择机房" style="width: 100%;">
                <el-option v-for="dc in datacenters" :key="dc.id" :label="dc.name" :value="dc.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="机架">
              <el-input v-model="formData.rack" placeholder="如: A01-01" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="IP地址" prop="ip">
              <el-input v-model="formData.ip" placeholder="请输入IP地址" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="管理IP">
              <el-input v-model="formData.mgmt_ip" placeholder="带外管理IP" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="操作系统">
              <el-input v-model="formData.os" placeholder="如: CentOS 7.9" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注">
          <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 远程操作结果 -->
    <el-dialog title="操作结果" v-model="remoteResultVisible" width="500px">
      <pre class="remote-output">{{ remoteResult }}</pre>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Operation, Download, ArrowDown } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const submitting = ref(false)
const devices = ref<any[]>([])
const datacenters = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref()
const selectedDevices = ref<any[]>([])
const remoteResultVisible = ref(false)
const remoteResult = ref('')

const stats = ref([
  { label: '总设备', value: '0', type: '' },
  { label: '在线', value: '0', type: 'success' },
  { label: '离线', value: '0', type: 'danger' },
  { label: '维护中', value: '0', type: 'warning' },
  { label: '故障', value: '0', type: 'fault' }
])

const searchForm = ref({ keyword: '', datacenter_id: '', status: '', type: '' })

const defaultForm = () => ({
  id: null, name: '', sn: '', type: 'server', datacenter_id: null,
  rack: '', ip: '', mgmt_ip: '', os: '', remark: ''
})
const formData = reactive(defaultForm())

const rules = {
  name: [{ required: true, message: '请输入设备名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  datacenter_id: [{ required: true, message: '请选择机房', trigger: 'change' }]
}

const types = [
  { label: '服务器', value: 'server' },
  { label: '交换机', value: 'switch' },
  { label: '路由器', value: 'router' },
  { label: '防火墙', value: 'firewall' },
  { label: '存储', value: 'storage' }
]

const getTypeLabel = (val: string) => types.find(t => t.value === val)?.label || val
const getStatusType = (s: string) => ({ online: 'success', offline: 'danger', maintenance: 'warning', fault: 'danger' } as Record<string, string>)[s] || 'info'
const getStatusLabel = (s: string) => ({ online: '在线', offline: '离线', maintenance: '维护中', fault: '故障' } as Record<string, string>)[s] || s
const getProgressColor = (p: number) => p > 90 ? '#f56c6c' : p > 70 ? '#e6a23c' : '#67c23a'

const fetchDevices = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/dcim', {
      params: { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    })
    devices.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取设备列表失败')
  } finally {
    loading.value = false
  }
}

const fetchDatacenters = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/dcim/datacenters')
    datacenters.value = data.data?.list || []
  } catch {}
}

const fetchStats = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/dcim/stats')
    if (data.data) {
      stats.value[0].value = String(data.data.total || 0)
      stats.value[1].value = String(data.data.online || 0)
      stats.value[2].value = String(data.data.offline || 0)
      stats.value[3].value = String(data.data.maintenance || 0)
      stats.value[4].value = String(data.data.fault || 0)
    }
  } catch {}
}

const handleSearch = () => { currentPage.value = 1; fetchDevices() }
const resetSearch = () => { searchForm.value = { keyword: '', datacenter_id: '', status: '', type: '' }; handleSearch() }
const handleSizeChange = (val: number) => { pageSize.value = val; fetchDevices() }
const handlePageChange = (val: number) => { currentPage.value = val; fetchDevices() }
const handleSelectionChange = (sel: any[]) => { selectedDevices.value = sel }

const handleCreate = () => {
  Object.assign(formData, defaultForm())
  dialogTitle.value = '添加设备'
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  Object.assign(formData, row)
  dialogTitle.value = '编辑设备'
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try { await formRef.value?.validate() } catch { return }
  submitting.value = true
  try {
    if (formData.id) {
      await request.put(`/admin/api/v1/dcim/${formData.id}`, formData)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/dcim', formData)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    fetchDevices()
    fetchStats()
  } catch {
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
  }
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(`确认删除设备「${row.name}」？`, '确认删除', { type: 'warning' }).then(async () => {
    try {
      await request.delete(`/admin/api/v1/dcim/${row.id}`)
      ElMessage.success('删除成功')
      fetchDevices()
      fetchStats()
    } catch {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

const handleRemoteAction = async (device: any, action: string) => {
  const actionMap: Record<string, string> = { ping: 'Ping', reboot: '重启', power_off: '关机', sol: 'SOL控制台' }
  if (['reboot', 'power_off'].includes(action)) {
    try {
      await ElMessageBox.confirm(`确认${actionMap[action]}设备「${device.name}」？`, '确认操作', { type: 'warning' })
    } catch { return }
  }
  try {
    loading.value = true
    const { data } = await request.post(`/admin/api/v1/dcim/${device.id}/remote`, { action })
    remoteResult.value = data.data?.output || '操作已提交'
    remoteResultVisible.value = true
  } catch {
    ElMessage.error(`${actionMap[action]}操作失败`)
  } finally {
    loading.value = false
  }
}

const handleBatchAction = () => {
  ElMessageBox.confirm(`对选中的 ${selectedDevices.value.length} 台设备执行操作？`, '批量操作', {
    confirmButtonText: '批量重启',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await request.post('/admin/api/v1/dcim/batch-action', {
        ids: selectedDevices.value.map((d: any) => d.id),
        action: 'reboot'
      })
      ElMessage.success('批量操作已提交')
      fetchDevices()
    } catch {
      ElMessage.error('批量操作失败')
    }
  }).catch(() => {})
}

const handleExport = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/dcim/export', { responseType: 'blob' })
    const url = URL.createObjectURL(data)
    const a = document.createElement('a')
    a.href = url; a.download = 'dcim-devices.xlsx'; a.click()
    URL.revokeObjectURL(url)
  } catch {
    ElMessage.error('导出失败')
  }
}

onMounted(() => { fetchDevices(); fetchDatacenters(); fetchStats() })
</script>

<style scoped lang="scss">
.dcim-page {
  .stats-grid {
    display: grid; grid-template-columns: repeat(5, 1fr); gap: 16px; margin-bottom: 20px;
    .stat-card {
      background: var(--bg-card); border-radius: var(--border-radius); padding: 16px 20px;
      box-shadow: var(--shadow-sm);
      .stat-label { font-size: 13px; color: var(--text-secondary); margin-bottom: 8px; }
      .stat-value { font-size: 24px; font-weight: 600; color: var(--text-primary);
        &.success { color: var(--success-color, #34c759); }
        &.danger { color: var(--danger-color, #f56c6c); }
        &.warning { color: var(--warning-color, #e6a23c); }
        &.fault { color: #f56c6c; }
      }
    }
  }
  .table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .table-footer { margin-top: 16px; display: flex; justify-content: space-between; align-items: center;
    .batch-info { font-size: 13px; color: var(--text-secondary); }
  }
  .remote-output { background: #1e1e1e; color: #d4d4d4; padding: 16px; border-radius: 8px; font-size: 13px; max-height: 400px; overflow: auto; }
}
</style>
