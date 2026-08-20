<template>
  <div class="contacts-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsContacts.title') }}</span>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('clientsContacts.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('clientsContacts.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('clientsContacts.userId')">
          <el-input v-model="searchForm.user_id" :placeholder="$t('clientsContacts.userId')" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="user_id" :label="$t('clientsContacts.associatedUser')" width="120" />
        <el-table-column prop="name" :label="$t('clientsContacts.name')" width="120" />
        <el-table-column prop="email" :label="$t('common.email')" width="180" />
        <el-table-column prop="phone" :label="$t('common.phone')" width="140" />
        <el-table-column prop="is_default" :label="$t('common.isDefault')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_default ? 'success' : 'info'" size="small">{{ row.is_default ? $t('common.yes') : $t('common.no') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="170" />
        <el-table-column :label="$t('common.action')" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-popconfirm :title="$t('clientsContacts.confirmDelete')" @confirm="handleDelete(row)">
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

    <el-dialog v-model="dialogVisible" :title="$t('clientsContacts.editContact')" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="80px">
        <el-form-item :label="$t('clientsContacts.name')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('clientsContacts.enterName')" />
        </el-form-item>
        <el-form-item :label="$t('common.email')" prop="email">
          <el-input v-model="formData.email" :placeholder="$t('clientsContacts.enterEmail')" />
        </el-form-item>
        <el-form-item :label="$t('common.phone')" prop="phone">
          <el-input v-model="formData.phone" :placeholder="$t('clientsContacts.enterPhone')" />
        </el-form-item>
        <el-form-item :label="$t('clientsContacts.address')">
          <el-input v-model="formData.address" :placeholder="$t('clientsContacts.enterAddress')" />
        </el-form-item>
        <el-form-item :label="$t('common.isDefault')">
          <el-switch v-model="formData.is_default" />
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
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()
const searchForm = reactive({ keyword: '', user_id: '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const dialogVisible = ref(false)
const formData = reactive({ id: undefined as number | undefined, name: '', email: '', phone: '', address: '', is_default: false })

const validateEmail = (_rule: any, value: string, callback: any) => {
  if (!value) { callback(); return }
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(value)) { callback(new Error($t('common.invalidEmail'))) } else { callback() }
}

const formRules: FormRules = {
  name: [{ required: true, message: $t('clientsContacts.enterName'), trigger: 'blur' }],
  email: [{ required: true, message: $t('clientsContacts.enterEmail'), trigger: 'blur' }, { validator: validateEmail, trigger: 'blur' }],
  phone: [{ required: true, message: $t('clientsContacts.enterPhone'), trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/contacts', params: { page: pagination.page, page_size: pagination.page_size, ...searchForm } })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (e) { ElMessage.error($t('clientsContacts.fetchFailed')) } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.user_id = ''; handleSearch() }

const handleEdit = (row: any) => { Object.assign(formData, row); dialogVisible.value = true }

const handleDelete = async (row: any) => {
  try { await request.del({ url: `/api/admin/contacts/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchData() } catch (e) { ElMessage.error($t('common.deleteFailed')) }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/contacts/${formData.id}` : '/api/admin/contacts'
      if (formData.id) { await request.put({ url, params: formData }) } else { await request.post({ url, params: formData }) }
      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (e) { ElMessage.error($t('common.operateFailed')) } finally { submitLoading.value = false }
  })
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }
onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.contacts-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
