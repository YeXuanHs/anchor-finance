<template>
  <div class="upstream-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.upstream.pageTitle') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('finance.upstream.addUpstream') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('finance.upstream.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('finance.upstream.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('finance.upstream.type')">
          <el-select v-model="searchForm.type" :placeholder="$t('finance.upstream.all')" clearable>
            <el-option :label="$t('finance.upstream.typeManual')" value="manual" />
            <el-option :label="$t('finance.upstream.typeV10')" value="v10" />
            <el-option :label="$t('finance.upstream.typeZjmf')" value="zjmf" />
            <el-option :label="$t('finance.upstream.typeAnchor')" value="anchor" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('finance.upstream.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('finance.upstream.all')" clearable>
            <el-option :label="$t('finance.upstream.enabled')" :value="1" />
            <el-option :label="$t('finance.upstream.disabled')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('finance.upstream.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('finance.upstream.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" :label="$t('finance.upstream.upstreamName')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="type" :label="$t('finance.upstream.apiType')" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="typeTagMap[row.type]" size="small">
              {{ typeLabelMap[row.type] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="url" :label="$t('finance.upstream.apiUrl')" min-width="220" show-overflow-tooltip />
        <el-table-column prop="status" :label="$t('finance.upstream.status')" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="test_status" :label="$t('finance.upstream.connectionStatus')" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.test_status === 'success'" type="success" size="small">{{ $t('finance.upstream.normal') }}</el-tag>
            <el-tag v-else-if="row.test_status === 'failed'" type="danger" size="small">{{ $t('finance.upstream.failed') }}</el-tag>
            <el-tag v-else type="info" size="small">{{ $t('finance.upstream.untested') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_sync_at" :label="$t('finance.upstream.lastSync')" width="170">
          <template #default="{ row }">
            {{ row.last_sync_at || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('finance.upstream.createdAt')" width="170" />
        <el-table-column :label="$t('finance.upstream.actions')" width="250" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleTest(row)" :loading="row._testing">
              {{ $t('finance.upstream.testConnection') }}
            </el-button>
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('finance.upstream.edit') }}</el-button>
            <el-popconfirm :title="$t('finance.upstream.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('finance.upstream.delete') }}</el-button>
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="650px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="110px">
        <el-form-item :label="$t('finance.upstream.upstreamName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('finance.upstream.enterUpstreamName')" />
        </el-form-item>
        <el-form-item :label="$t('finance.upstream.apiType')" prop="type">
          <el-select v-model="formData.type" :placeholder="$t('finance.upstream.selectApiType')" style="width: 100%" @change="onTypeChange">
            <el-option :label="$t('finance.upstream.typeManual')" value="manual" />
            <el-option :label="$t('finance.upstream.typeV10')" value="v10" />
            <el-option :label="$t('finance.upstream.typeZjmf')" value="zjmf" />
            <el-option :label="$t('finance.upstream.typeAnchor')" value="anchor" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('finance.upstream.apiUrl')" prop="url" v-if="formData.type !== 'manual'">
          <el-input v-model="formData.url" placeholder="https://example.com/api" />
        </el-form-item>
        <el-form-item :label="$t('finance.upstream.apiKey')" prop="api_key" v-if="formData.type !== 'manual'">
          <el-input v-model="formData.api_key" :placeholder="$t('finance.upstream.enterApiKey')" show-password />
        </el-form-item>
        <el-form-item :label="$t('finance.upstream.apiPassword')" prop="api_password" v-if="formData.type !== 'manual'">
          <el-input v-model="formData.api_password" :placeholder="$t('finance.upstream.enterApiPassword')" show-password />
        </el-form-item>
        <el-form-item :label="$t('finance.upstream.syncInterval')" prop="sync_interval" v-if="formData.type !== 'manual'">
          <el-input-number
            v-model="formData.sync_interval"
            :min="5"
            :max="1440"
            :step="5"
            controls-position="right"
          />
          <span class="form-tip">{{ $t('finance.upstream.minutes') }}</span>
        </el-form-item>
        <el-form-item :label="$t('finance.upstream.remark')" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="3" :placeholder="$t('finance.upstream.enterRemark')" />
        </el-form-item>
        <el-form-item :label="$t('finance.upstream.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('finance.upstream.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('finance.upstream.confirm') }}</el-button>
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
import { $t } from '@/locales'

const typeLabelMap: Record<string, string> = {
  manual: $t('finance.upstream.typeManual'),
  v10: $t('finance.upstream.typeV10'),
  zjmf: $t('finance.upstream.typeZjmf'),
  anchor: $t('finance.upstream.typeAnchor')
}

const typeTagMap: Record<string, any> = {
  manual: 'info',
  v10: 'warning',
  zjmf: '',
  anchor: 'success'
}

const loading = ref(false)
const submitLoading = ref(false)

const searchForm = reactive({
  keyword: '',
  type: '',
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref($t('finance.upstream.addUpstream'))
const formRef = ref<FormInstance>()

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  type: '',
  url: '',
  api_key: '',
  api_password: '',
  sync_interval: 30,
  remark: '',
  status: 1
})

const formRules: FormRules = {
  name: [
    { required: true, message: $t('finance.upstream.enterUpstreamName'), trigger: 'blur' },
    { min: 2, max: 50, message: $t('finance.upstream.nameLength'), trigger: 'blur' }
  ],
  type: [
    { required: true, message: $t('finance.upstream.selectApiType'), trigger: 'change' }
  ],
  url: [
    { required: true, message: $t('finance.upstream.enterApiUrl'), trigger: 'blur' }
  ],
  api_key: [
    { required: true, message: $t('finance.upstream.enterApiKey'), trigger: 'blur' }
  ]
}

const onTypeChange = () => {
  if (formData.type === 'manual') {
    formData.url = ''
    formData.api_key = ''
    formData.api_password = ''
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/upstream/providers',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取上游列表失败:', error)
    ElMessage.error($t('finance.upstream.fetchListFailed'))
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
  searchForm.type = ''
  searchForm.status = undefined
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = $t('finance.upstream.addUpstream')
  formData.id = undefined
  formData.name = ''
  formData.type = ''
  formData.url = ''
  formData.api_key = ''
  formData.api_password = ''
  formData.sync_interval = 30
  formData.remark = ''
  formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('finance.upstream.editUpstream')
  Object.assign(formData, {
    id: row.id,
    name: row.name,
    type: row.type,
    url: row.url || '',
    api_key: row.api_key || '',
    api_password: '',
    sync_interval: row.sync_interval || 30,
    remark: row.remark || '',
    status: row.status
  })
  dialogVisible.value = true
}

const handleTest = async (row: any) => {
  row._testing = true
  try {
    const data = await request.post({
      url: `/api/admin/upstream/providers/${row.id}/test`
    })
    row.test_status = data.status === 'success' ? 'success' : 'failed'
    ElMessage.success(data.message || $t('finance.upstream.testComplete'))
  } catch (error) {
    row.test_status = 'failed'
    ElMessage.error($t('finance.upstream.testFailed'))
  } finally {
    row._testing = false
  }
}

const handleStatusChange = async (row: any) => {
  try {
    await request.put({
      url: `/api/admin/upstream/providers/${row.id}`,
      params: { status: row.status }
    })
    ElMessage.success($t('finance.upstream.statusUpdated'))
  } catch (error) {
    row.status = row.status === 1 ? 0 : 1
    ElMessage.error($t('finance.upstream.updateStatusFailed'))
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/upstream/providers/${row.id}`
    })
    ElMessage.success($t('finance.upstream.deleteSuccess'))
    fetchList()
  } catch (error) {
    ElMessage.error($t('finance.upstream.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const submitData: any = { ...formData }
      if (formData.type === 'manual') {
        delete submitData.url
        delete submitData.api_key
        delete submitData.api_password
        delete submitData.sync_interval
      }
      if (!submitData.api_password) {
        delete submitData.api_password
      }

      const url = formData.id ? `/api/admin/upstream/providers/${formData.id}` : '/api/admin/upstream/providers'

      if (formData.id) {
        await request.put({ url, params: submitData })
      } else {
        await request.post({ url, params: submitData })
      }

      ElMessage.success(formData.id ? $t('finance.upstream.updateSuccess') : $t('finance.upstream.addSuccess'))
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      ElMessage.error($t('finance.upstream.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
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
.upstream-page {
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
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}
</style>
