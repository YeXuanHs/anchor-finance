<template>
  <div class="resources-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsResources.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('clientsResources.addResource') }}
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('clientsResources.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('clientsResources.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('common.type')">
          <el-select v-model="searchForm.type" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('clientsResources.ipAddress')" value="ip" />
            <el-option :label="$t('clientsResources.domain')" value="domain" />
            <el-option :label="$t('clientsResources.sslCert')" value="ssl" />
            <el-option :label="$t('clientsResources.otherType')" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('clientsResources.inUse')" value="in_use" />
            <el-option :label="$t('clientsResources.idle')" value="idle" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" :label="$t('clientsResources.resourceName')" width="150" />
        <el-table-column prop="type" :label="$t('common.type')" width="100" align="center">
          <template #default="{ row }">
            <el-tag size="small">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="value" :label="$t('clientsResources.resourceValue')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="client_name" :label="$t('clientsResources.ownerClient')" width="120" />
        <el-table-column prop="expire_at" :label="$t('clientsResources.expireTime')" width="170" />
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="170" />
        <el-table-column :label="$t('common.action')" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-popconfirm :title="$t('clientsResources.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('common.delete') }}</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? $t('common.edit') : $t('clientsResources.addResource')" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('common.type')" prop="type">
          <el-select v-model="formData.type" :placeholder="$t('common.select')" style="width: 100%">
            <el-option :label="$t('clientsResources.ipAddress')" value="ip" />
            <el-option :label="$t('clientsResources.domain')" value="domain" />
            <el-option :label="$t('clientsResources.sslCert')" value="ssl" />
            <el-option :label="$t('clientsResources.otherType')" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('clientsResources.resourceName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('clientsResources.enterResourceName')" />
        </el-form-item>
        <el-form-item :label="$t('clientsResources.resourceValue')" prop="value">
          <el-input v-model="formData.value" :placeholder="$t('clientsResources.enterResourceValue')" />
        </el-form-item>
        <el-form-item :label="$t('clientsResources.ownerClient')" prop="client_id">
          <el-select v-model="formData.client_id" :placeholder="$t('clientsResources.selectClient')" filterable style="width: 100%">
            <el-option v-for="client in clientList" :key="client.id" :label="client.username" :value="client.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()
const searchForm = reactive({ keyword: '', type: '', status: '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const clientList = ref<any[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formData = reactive({ id: undefined as number | undefined, type: '', name: '', value: '', client_id: undefined as number | undefined })

const formRules: FormRules = {
  type: [{ required: true, message: $t('common.select'), trigger: 'change' }],
  value: [{ required: true, message: $t('clientsResources.enterResourceValue'), trigger: 'blur' }]
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = { ip: $t('clientsResources.ipAddress'), domain: $t('clientsResources.domain'), ssl: $t('clientsResources.sslCert'), other: $t('clientsResources.otherType') }
  return map[type] || type
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/resources', params: { page: pagination.page, page_size: pagination.page_size, ...searchForm } })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (e) { ElMessage.error($t('common.fetchFailed')) } finally { loading.value = false }
}

const fetchClientList = async () => {
  try {
    const data = await request.get({ url: '/api/admin/users', params: { page_size: 9999 } })
    clientList.value = data.list || []
  } catch (e) { ElMessage.error($t('clientsResources.fetchClientsFailed')) }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.type = ''; searchForm.status = ''; handleSearch() }

const handleAdd = () => { isEdit.value = false; formData.id = undefined; formData.type = ''; formData.name = ''; formData.value = ''; formData.client_id = undefined; dialogVisible.value = true }

const handleEdit = (row: any) => { isEdit.value = true; Object.assign(formData, row); dialogVisible.value = true }

const handleDelete = async (row: any) => {
  try { await request.del({ url: `/api/admin/resources/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchData() } catch (e) { ElMessage.error($t('common.deleteFailed')) }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/resources/${formData.id}` : '/api/admin/resources'
      if (formData.id) { await request.put({ url, params: formData }) } else { await request.post({ url, params: formData }) }
      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (e) { ElMessage.error($t('common.operateFailed')) } finally { submitLoading.value = false }
  })
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }
onMounted(() => { fetchData(); fetchClientList() })
</script>

<style scoped lang="scss">
.resources-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
