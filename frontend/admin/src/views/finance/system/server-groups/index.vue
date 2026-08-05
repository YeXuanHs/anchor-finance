<template>
  <div class="server-groups-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>服务器组管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加服务器组
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="组名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
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
        <el-table-column prop="name" label="组名称" min-width="180" />
        <el-table-column prop="code" label="组编码" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="server_count" label="服务器数量" width="110" align="center" />
        <el-table-column prop="load_balance_type" label="负载均衡策略" width="130" align="center">
          <template #default="{ row }">
            <el-tag size="small">{{ loadBalanceMap[row.load_balance_type] || row.load_balance_type || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleManageServers(row)">管理服务器</el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该服务器组吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
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

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="110px">
        <el-form-item label="组名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入服务器组名称" />
        </el-form-item>
        <el-form-item label="组编码" prop="code">
          <el-input v-model="formData.code" placeholder="请输入组编码（英文）" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="负载均衡策略" prop="load_balance_type">
          <el-select v-model="formData.load_balance_type" placeholder="请选择" style="width: 100%">
            <el-option label="轮询" value="round_robin" />
            <el-option label="加权轮询" value="weighted_round_robin" />
            <el-option label="最少连接" value="least_connections" />
            <el-option label="IP哈希" value="ip_hash" />
            <el-option label="随机" value="random" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 管理服务器对话框 -->
    <el-dialog v-model="serverDialogVisible" title="管理关联服务器" width="750px" destroy-on-close>
      <div class="server-transfer">
        <el-transfer
          v-model="selectedServerIds"
          :data="allServers"
          :titles="['可选服务器', '已关联服务器']"
          :props="{ key: 'id', label: 'name' }"
          filterable
          filter-placeholder="搜索服务器"
        />
      </div>
      <template #footer>
        <el-button @click="serverDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveServers" :loading="serverSaving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'ServerGroups' })

const loadBalanceMap: Record<string, string> = {
  round_robin: '轮询',
  weighted_round_robin: '加权轮询',
  least_connections: '最少连接',
  ip_hash: 'IP哈希',
  random: '随机'
}

const loading = ref(false)
const submitLoading = ref(false)
const serverSaving = ref(false)
const dialogVisible = ref(false)
const serverDialogVisible = ref(false)
const dialogTitle = ref('添加服务器组')
const formRef = ref<FormInstance>()
const currentGroupId = ref<number>(0)

const searchForm = reactive({ keyword: '', status: undefined as number | undefined })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const allServers = ref<any[]>([])
const selectedServerIds = ref<number[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  code: '',
  description: '',
  load_balance_type: 'round_robin',
  sort: 0,
  status: 1
})

const formRules: FormRules = {
  name: [
    { required: true, message: '请输入服务器组名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入组编码', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '只能包含字母、数字和下划线', trigger: 'blur' }
  ]
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/server-groups',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || data || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取服务器组列表失败')
  } finally {
    loading.value = false
  }
}

const fetchAllServers = async () => {
  try {
    const data = await request.get({ url: '/api/admin/servers', params: { page_size: 999 } })
    allServers.value = data.list || data || []
  } catch (error) {
    console.error('获取服务器列表失败:', error)
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.status = undefined; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleAdd = () => {
  dialogTitle.value = '添加服务器组'
  formData.id = undefined; formData.name = ''; formData.code = ''
  formData.description = ''; formData.load_balance_type = 'round_robin'
  formData.sort = 0; formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑服务器组'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/server-groups/${row.id}` })
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/server-groups/${formData.id}`, params: formData })
      } else {
        await request.post({ url: '/api/admin/server-groups', params: formData })
      }
      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

const handleManageServers = async (row: any) => {
  currentGroupId.value = row.id
  try {
    const data = await request.get({ url: `/api/admin/server-groups/${row.id}/servers` })
    selectedServerIds.value = (data.list || data || []).map((s: any) => s.id || s)
    serverDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取关联服务器失败')
  }
}

const handleSaveServers = async () => {
  serverSaving.value = true
  try {
    await request.put({
      url: `/api/admin/server-groups/${currentGroupId.value}/servers`,
      params: { server_ids: selectedServerIds.value }
    })
    ElMessage.success('关联服务器更新成功')
    serverDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    serverSaving.value = false
  }
}

onMounted(() => { fetchData(); fetchAllServers() })
</script>

<style scoped lang="scss">
.server-groups-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.server-transfer { display: flex; justify-content: center; }
</style>
