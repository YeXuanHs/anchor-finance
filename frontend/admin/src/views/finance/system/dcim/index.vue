<template>
  <div class="dcim-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('dcim.title') }}</span>
          <div>
            <el-button @click="handleRefreshAll" :loading="refreshAllLoading">
              <el-icon><Refresh /></el-icon>
              {{ $t('dcim.refreshAll') }}
            </el-button>
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              {{ $t('dcim.addServer') }}
            </el-button>
          </div>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('dcim.keyword')">
          <el-input v-model="searchForm.search" :placeholder="$t('dcim.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" :label="$t('dcim.serverName')" min-width="150" />
        <el-table-column prop="hostname" :label="$t('dcim.hostname')" min-width="180" />
        <el-table-column prop="server_num" :label="$t('dcim.running')" width="80" align="center" />
        <el-table-column prop="api_status" :label="$t('dcim.apiStatus')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.api_status === 1 ? 'success' : 'danger'" size="small">
              {{ row.api_status === 1 ? $t('common.normal') : $t('common.abnormal') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.action')" width="350" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="primary" link @click="handleDetail(row)">{{ $t('common.detail') }}</el-button>
            <el-button type="success" link @click="handleRefreshStatus(row)">{{ $t('dcim.refreshStatus') }}</el-button>
            <el-popconfirm
              v-if="row.removable"
              :title="$t('dcim.confirmDelete')"
              @confirm="handleDelete(row)"
            >
              <template #reference>
                <el-button type="danger" link>{{ $t('common.delete') }}</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.limit"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="650px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item :label="$t('dcim.serverName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('dcim.inputServerName')" />
        </el-form-item>
        <el-form-item :label="$t('dcim.hostnameIp')" prop="hostname">
          <el-input v-model="formData.hostname" :placeholder="$t('dcim.inputHostname')" />
        </el-form-item>
        <el-form-item :label="$t('dcim.username')" prop="username">
          <el-input v-model="formData.username" :placeholder="$t('dcim.inputUsername')" />
        </el-form-item>
        <el-form-item :label="$t('dcim.password')" prop="password">
          <el-input v-model="formData.password" type="password" show-password :placeholder="$t('dcim.inputPassword')" />
        </el-form-item>
        <el-form-item :label="$t('dcim.port')" prop="port">
          <el-input-number v-model="formData.port" :min="0" :max="65535" />
        </el-form-item>
        <el-form-item :label="$t('dcim.useSsl')">
          <el-switch v-model="formData.secure" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-switch v-model="formData.disabled" :active-value="0" :inactive-value="1" :active-text="$t('common.enable')" :inactive-text="$t('common.disable')" />
        </el-form-item>
        <el-form-item :label="$t('dcim.userPrefix')">
          <el-input v-model="formData.user_prefix" :placeholder="$t('dcim.optional')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" :title="$t('dcim.serverDetail')" width="700px">
      <el-descriptions :column="2" border v-if="detailData">
        <el-descriptions-item label="ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('dcim.serverName')">{{ detailData.name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('dcim.hostname')">{{ detailData.hostname }}</el-descriptions-item>
        <el-descriptions-item :label="$t('dcim.username')">{{ detailData.username }}</el-descriptions-item>
        <el-descriptions-item :label="$t('dcim.port')">{{ detailData.port }}</el-descriptions-item>
        <el-descriptions-item label="SSL">{{ detailData.secure ? $t('common.yes') : $t('common.no') }}</el-descriptions-item>
        <el-descriptions-item :label="$t('common.status')">{{ detailData.disabled ? $t('common.disable') : $t('common.enable') }}</el-descriptions-item>
        <el-descriptions-item :label="$t('dcim.billingType')">{{ detailData.bill_type === 'month' ? $t('dcim.monthly') : $t('dcim.traffic') }}</el-descriptions-item>
        <el-descriptions-item :label="$t('dcim.reinstallLimit')">{{ detailData.reinstall_times || $t('common.none') }}</el-descriptions-item>
        <el-descriptions-item :label="$t('dcim.paidReinstall')">{{ detailData.buy_times ? $t('common.enable') : $t('common.disable') }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">{{ $t('common.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

const loading = ref(false)
const submitLoading = ref(false)
const refreshAllLoading = ref(false)

const searchForm = reactive({
  search: ''
})

const pagination = reactive({
  page: 1,
  limit: 10,
  total: 0
})

const tableData = ref<any[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref($t('dcim.addServer'))
const formRef = ref<FormInstance>()

const detailVisible = ref(false)
const detailData = ref<any>(null)

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  hostname: '',
  username: '',
  password: '',
  port: 0,
  secure: 0,
  disabled: 0,
  user_prefix: ''
})

const formRules: FormRules = {
  name: [{ required: true, message: $t('dcim.inputServerName'), trigger: 'blur' }],
  hostname: [{ required: true, message: $t('dcim.inputHostname'), trigger: 'blur' }],
  username: [{ required: true, message: $t('dcim.inputUsername'), trigger: 'blur' }],
  password: [{ required: true, message: $t('dcim.inputPassword'), trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/dcim/servers',
      params: {
        page: pagination.page,
        limit: pagination.limit,
        search: searchForm.search || undefined
      }
    })
    tableData.value = data.list || []
    pagination.total = data.sum || 0
  } catch (error) {
    console.error('获取DCIM服务器列表失败:', error)
    ElMessage.error($t('common.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.search = ''
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = $t('dcim.addServer')
  formData.id = undefined
  formData.name = ''
  formData.hostname = ''
  formData.username = ''
  formData.password = ''
  formData.port = 0
  formData.secure = 0
  formData.disabled = 0
  formData.user_prefix = ''
  dialogVisible.value = true
}

const handleEdit = async (row: any) => {
  dialogTitle.value = $t('dcim.editServer')
  try {
    const data = await request.get({ url: `/api/admin/dcim/servers/${row.id}` })
    Object.assign(formData, data)
    formData.password = ''
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error($t('dcim.fetchDetailFailed'))
  }
}

const handleDetail = async (row: any) => {
  try {
    const data = await request.get({ url: `/api/admin/dcim/servers/${row.id}` })
    detailData.value = data
    detailVisible.value = true
  } catch (error) {
    ElMessage.error($t('dcim.fetchDetailFailed'))
  }
}

const handleRefreshStatus = async (row: any) => {
  try {
    const data = await request.post({ url: `/api/admin/dcim/servers/${row.id}/refresh` })
    ElMessage.success(data.msg || $t('dcim.refreshSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('dcim.refreshFailed'))
  }
}

const handleRefreshAll = async () => {
  refreshAllLoading.value = true
  try {
    // 后端无批量刷新接口，逐个刷新
    const list = tableData.value
    for (const row of list) {
      await request.post({ url: `/api/admin/dcim/servers/${row.id}/refresh` })
    }
    ElMessage.success($t('dcim.refreshComplete'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('dcim.refreshFailed'))
  } finally {
    refreshAllLoading.value = false
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/dcim/servers/${row.id}` })
    ElMessage.success($t('common.deleteSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('common.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/dcim/servers/${formData.id}`, params: { ...formData } })
      } else {
        await request.post({ url: '/api/admin/dcim/servers', params: { ...formData } })
      }
      ElMessage.success(formData.id ? $t('common.editSuccess') : $t('common.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('common.operateFailed'))
    } finally {
      submitLoading.value = false
    }
  })
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
.dcim-page {
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
