<template>
  <div class="servers-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="服务器名称/IP" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="物理机" value="physical" />
            <el-option label="云服务器" value="cloud" />
            <el-option label="VPS" value="vps" />
          </el-select>
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
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>服务器配置</h3>
        <el-button type="primary" @click="openAddDialog">
          <el-icon><Plus /></el-icon>
          添加服务器
        </el-button>
      </div>

      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" min-width="140" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ typeTextMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="hostname" label="主机名" width="140" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP 地址" width="140" />
        <el-table-column prop="port" label="端口" width="80" />
        <el-table-column prop="api_key" label="API 密钥" width="160">
          <template #default="{ row }">
            <span v-if="row.api_key">{{ row.api_key.substring(0, 8) }}****</span>
            <span v-else style="color: var(--el-text-color-secondary);">未设置</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTypeMap[row.status]" size="small">
              {{ statusTextMap[row.status] || row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_check" label="最后检测" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="success" link @click="testConnection(row)">连接测试</el-button>
            <el-button type="primary" link @click="openEditDialog(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchData"
          @size-change="fetchData"
        />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑服务器' : '添加服务器'" width="650px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="110px">
        <el-form-item label="服务器名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入服务器名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型">
            <el-option label="物理机" value="physical" />
            <el-option label="云服务器" value="cloud" />
            <el-option label="VPS" value="vps" />
          </el-select>
        </el-form-item>
        <el-form-item label="主机名" prop="hostname">
          <el-input v-model="formData.hostname" placeholder="服务器主机名" />
        </el-form-item>
        <el-form-item label="IP 地址" prop="ip">
          <el-input v-model="formData.ip" placeholder="如: 192.168.1.100" />
        </el-form-item>
        <el-form-item label="SSH 端口" prop="port">
          <el-input-number v-model="formData.port" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item label="SSH 用户名" prop="username">
          <el-input v-model="formData.username" placeholder="root" />
        </el-form-item>
        <el-form-item label="认证方式" prop="auth_type">
          <el-radio-group v-model="formData.auth_type">
            <el-radio value="password">密码</el-radio>
            <el-radio value="key">密钥</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="formData.auth_type === 'password'">
          <el-input v-model="formData.password" type="password" show-password placeholder="SSH 密码" />
        </el-form-item>
        <el-form-item label="私钥" prop="private_key" v-if="formData.auth_type === 'key'">
          <el-input v-model="formData.private_key" type="textarea" :rows="4" placeholder="粘贴 SSH 私钥" />
        </el-form-item>
        <el-form-item label="API 密钥" prop="api_key">
          <el-input v-model="formData.api_key" placeholder="留空则自动生成">
            <template #append>
              <el-button @click="generateApiKey">生成</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="API 端口" prop="api_port">
          <el-input-number v-model="formData.api_port" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item label="数据中心" prop="datacenter_id">
          <el-select v-model="formData.datacenter_id" placeholder="请选择数据中心" clearable>
            <el-option v-for="dc in datacenters" :key="dc.id" :label="dc.name" :value="dc.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
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
const tableData = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const formRef = ref<FormInstance>()
const datacenters = ref<any[]>([])

const typeTextMap: Record<string, string> = { physical: '物理机', cloud: '云服务器', vps: 'VPS' }
const statusTypeMap: Record<string, string> = { online: 'success', offline: 'danger', maintenance: 'warning' }
const statusTextMap: Record<string, string> = { online: '在线', offline: '离线', maintenance: '维护中' }

const searchForm = ref({ keyword: '', type: '', status: '' })

const formData = reactive({
  name: '',
  type: 'cloud' as string,
  hostname: '',
  ip: '',
  port: 22,
  username: 'root',
  auth_type: 'password' as string,
  password: '',
  private_key: '',
  api_key: '',
  api_port: 8080,
  datacenter_id: null as number | null,
  remark: ''
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入服务器名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  ip: [{ required: true, message: '请输入 IP 地址', trigger: 'blur' }],
  port: [{ required: true, message: '请输入端口', trigger: 'blur' }]
}

const generateApiKey = () => {
  const chars = 'abcdef0123456789'
  formData.api_key = Array.from({ length: 32 }, () => chars[Math.floor(Math.random() * chars.length)]).join('')
}

const handleSearch = () => { page.value = 1; fetchData() }

const resetSearch = () => {
  searchForm.value = { keyword: '', type: '', status: '' }
  handleSearch()
}

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/servers', {
      params: { page: page.value, page_size: pageSize.value, ...searchForm.value }
    })
    tableData.value = data.data || []
    total.value = data.total || 0
  } catch {} finally {
    loading.value = false
  }
}

const fetchDatacenters = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/datacenters')
    datacenters.value = data.data || []
  } catch {}
}

const resetForm = () => {
  Object.assign(formData, {
    name: '', type: 'cloud', hostname: '', ip: '', port: 22,
    username: 'root', auth_type: 'password', password: '', private_key: '',
    api_key: '', api_port: 8080, datacenter_id: null, remark: ''
  })
}

const openAddDialog = () => {
  isEdit.value = false
  editId.value = null
  resetForm()
  dialogVisible.value = true
}

const openEditDialog = (row: any) => {
  isEdit.value = true
  editId.value = row.id
  Object.assign(formData, {
    name: row.name, type: row.type, hostname: row.hostname || '', ip: row.ip,
    port: row.port || 22, username: row.username || 'root',
    auth_type: row.auth_type || 'password', password: '', private_key: '',
    api_key: row.api_key || '', api_port: row.api_port || 8080,
    datacenter_id: row.datacenter_id || null, remark: row.remark || ''
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (isEdit.value) {
      await request.put(`/admin/api/v1/servers/${editId.value}`, formData)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/servers', formData)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch {} finally {
    submitLoading.value = false
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm(`确定删除服务器「${row.name}」吗？`, '提示', { type: 'warning' })
  try {
    await request.delete(`/admin/api/v1/servers/${row.id}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch {}
}

const testConnection = async (row: any) => {
  try {
    const { data } = await request.post(`/admin/api/v1/servers/${row.id}/test`)
    if (data.success) {
      ElMessage.success(`连接成功，延迟 ${data.latency}ms`)
    } else {
      ElMessage.error(`连接失败: ${data.message}`)
    }
  } catch {}
}

onMounted(() => {
  fetchData()
  fetchDatacenters()
})
</script>

<style scoped lang="scss">
.servers-page {
  .search-bar { margin-bottom: 16px; }
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 18px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
}
</style>
