<template>
  <div class="first-groups-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('firstGroup.title') }}</span>
          <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('firstGroup.addGroup') }}</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('firstGroup.groupName')">
          <el-input v-model="searchForm.name" :placeholder="$t('firstGroup.enterGroupName')" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="filteredTableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" :label="$t('firstGroup.groupName')" min-width="200" />
        <el-table-column prop="code" :label="$t('firstGroup.groupCode')" width="150" />
        <el-table-column prop="description" :label="$t('common.description')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="product_count" :label="$t('firstGroup.productCount')" width="100" align="center" />
        <el-table-column prop="sort" :label="$t('firstGroup.sort')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('firstGroup.status')" width="80" align="center">
          <template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">{{ row.status === 1 ? $t('common.enable') : $t('common.disable') }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180" />
        <el-table-column :label="$t('firstGroup.operations')" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-popconfirm :title="$t('firstGroup.confirmDelete')" @confirm="handleDelete(row)"><template #reference><el-button type="danger" link>{{ $t('common.delete') }}</el-button></template></el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('firstGroup.groupName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('firstGroup.enterGroupName')" />
        </el-form-item>
        <el-form-item :label="$t('firstGroup.groupCode')" prop="code">
          <el-input v-model="formData.code" :placeholder="$t('firstGroup.enterGroupCode')" />
        </el-form-item>
        <el-form-item :label="$t('common.description')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('firstGroup.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('firstGroup.sort')" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('firstGroup.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" :active-text="$t('common.enable')" :inactive-text="$t('common.disable')" />
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

defineOptions({ name: 'ProductFirstGroups' })

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref($t('firstGroup.addGroup'))
const formRef = ref<FormInstance>()

const searchForm = reactive({ name: '' })
const tableData = ref<any[]>([])

const formData = reactive({ id: undefined as number | undefined, name: '', code: '', description: '', sort: 0, status: 1 })

const formRules: FormRules = {
  name: [{ required: true, message: () => $t('firstGroup.enterGroupName'), trigger: 'blur' }],
  code: [{ required: true, message: () => $t('firstGroup.enterGroupCode'), trigger: 'blur' }]
}

const filteredTableData = computed(() => {
  if (!searchForm.name) return tableData.value
  return tableData.value.filter(item => item.name.toLowerCase().includes(searchForm.name.toLowerCase()))
})

const fetchData = async () => {
  loading.value = true
  try { const data = await request.get({ url: '/api/admin/product-first-groups' }); tableData.value = data.list || data || [] } catch { ElMessage.error($t('firstGroup.fetchFailed')) } finally { loading.value = false }
}

const handleSearch = () => {}
const handleReset = () => { searchForm.name = '' }

const handleAdd = () => { dialogTitle.value = $t('firstGroup.addGroup'); formData.id = undefined; formData.name = ''; formData.code = ''; formData.description = ''; formData.sort = 0; formData.status = 1; dialogVisible.value = true }
const handleEdit = (row: any) => { dialogTitle.value = $t('firstGroup.editGroup'); Object.assign(formData, row); dialogVisible.value = true }

const handleDelete = async (row: any) => { try { await request.del({ url: `/api/admin/product-first-groups/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchData() } catch { ElMessage.error($t('common.deleteFailed')) } }

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) await request.put({ url: `/api/admin/product-first-groups/${formData.id}`, params: formData })
      else await request.post({ url: '/api/admin/product-first-groups', params: formData })
      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess')); dialogVisible.value = false; fetchData()
    } catch { ElMessage.error($t('common.operateFailed')) } finally { submitLoading.value = false }
  })
}

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.first-groups-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
</style>
