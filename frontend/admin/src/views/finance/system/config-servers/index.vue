<template>
  <div class="config-servers-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.configServers.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('finance.configServers.addServer') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('finance.configServers.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('finance.configServers.nameOrIp')" clearable />
        </el-form-item>
        <el-form-item :label="$t('finance.configServers.type')">
          <el-select v-model="searchForm.type" :placeholder="$t('finance.configServers.all')" clearable>
            <el-option :label="$t('finance.configServers.physical')" value="physical" />
            <el-option :label="$t('finance.configServers.vm')" value="vm" />
            <el-option :label="$t('finance.configServers.container')" value="container" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('finance.configServers.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('finance.configServers.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" :label="$t('finance.configServers.serverName')" min-width="150" />
        <el-table-column prop="ip" :label="$t('finance.configServers.ipAddress')" width="150" />
        <el-table-column prop="type" :label="$t('finance.configServers.type')" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('finance.configServers.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ row.status === 'active' ? $t('finance.configServers.normal') : $t('finance.configServers.abnormal') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="port" :label="$t('finance.configServers.port')" width="80" />
        <el-table-column prop="username" :label="$t('finance.configServers.username')" width="120" />
        <el-table-column :label="$t('finance.configServers.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('finance.configServers.edit') }}</el-button>
            <el-button type="success" link @click="handleTest(row)">{{ $t('finance.configServers.testConnection') }}</el-button>
            <el-popconfirm :title="$t('finance.configServers.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('finance.configServers.delete') }}</el-button>
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('finance.configServers.serverName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('finance.configServers.enterServerName')" />
        </el-form-item>
        <el-form-item :label="$t('finance.configServers.ipAddress')" prop="ip">
          <el-input v-model="formData.ip" :placeholder="$t('finance.configServers.enterIpAddress')" />
        </el-form-item>
        <el-form-item :label="$t('finance.configServers.port')" prop="port">
          <el-input-number v-model="formData.port" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item :label="$t('finance.configServers.username')" prop="username">
          <el-input v-model="formData.username" :placeholder="$t('finance.configServers.enterUsername')" />
        </el-form-item>
        <el-form-item :label="$t('finance.configServers.password')" prop="password">
          <el-input v-model="formData.password" type="password" :placeholder="$t('finance.configServers.enterPassword')" show-password />
        </el-form-item>
        <el-form-item :label="$t('finance.configServers.type')">
          <el-select v-model="formData.type" :placeholder="$t('finance.configServers.selectType')">
            <el-option :label="$t('finance.configServers.physical')" value="physical" />
            <el-option :label="$t('finance.configServers.vm')" value="vm" />
            <el-option :label="$t('finance.configServers.container')" value="container" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('finance.configServers.description')">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('finance.configServers.enterDescription')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('finance.configServers.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('finance.configServers.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref($t('finance.configServers.addServer'))
const formRef = ref<FormInstance>()

const tableData = ref<any[]>([])

const searchForm = reactive({ keyword: '', type: '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

const formData = reactive({
  id: undefined as number | undefined,
  name: '', ip: '', port: 22, username: '', password: '', type: 'vm', description: ''
})

const formRules: FormRules = {
  name: [{ required: true, message: $t('finance.configServers.enterServerName'), trigger: 'blur' }],
  ip: [{ required: true, message: $t('finance.configServers.enterIpAddress'), trigger: 'blur' }],
  username: [{ required: true, message: $t('finance.configServers.enterUsername'), trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/config/servers',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error($t('finance.configServers.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.type = ''; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleAdd = () => {
  dialogTitle.value = $t('finance.configServers.addServer')
  formData.id = undefined; formData.name = ''; formData.ip = ''; formData.port = 22
  formData.username = ''; formData.password = ''; formData.type = 'vm'; formData.description = ''
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('finance.configServers.editServer')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleTest = async (row: any) => {
  try {
    await request.get({ url: `/api/admin/config-servers/test-link/${row.id}` })
    ElMessage.success($t('finance.configServers.testSuccess'))
  } catch (error) {
    ElMessage.error($t('finance.configServers.testFailed'))
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/config/servers/${row.id}` })
    ElMessage.success($t('finance.configServers.deleteSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('finance.configServers.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/config/servers/${formData.id}`, params: formData })
      } else {
        await request.post({ url: '/api/admin/config/servers', params: formData })
      }
      ElMessage.success(formData.id ? $t('finance.configServers.updateSuccess') : $t('finance.configServers.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('finance.configServers.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.config-servers-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
