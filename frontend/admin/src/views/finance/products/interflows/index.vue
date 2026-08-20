<template>
  <div class="interflows-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('interflow.title') }}</span>
          <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('interflow.addInterflow') }}</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('interflow.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('interflow.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('interflow.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('common.enable')" :value="1" />
            <el-option :label="$t('common.disable')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" :label="$t('interflow.name')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="source_product_name" :label="$t('interflow.sourceProduct')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="target_product_name" :label="$t('interflow.targetProduct')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="type" :label="$t('interflow.type')" width="120" align="center">
          <template #default="{ row }"><el-tag size="small">{{ $t(`interflow.typeMap.${row.type}`) || row.type }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="description" :label="$t('common.description')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="sort" :label="$t('interflow.sort')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('interflow.status')" width="80" align="center">
          <template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">{{ row.status === 1 ? $t('common.enable') : $t('common.disable') }}</el-tag></template>
        </el-table-column>
        <el-table-column :label="$t('interflow.operations')" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-popconfirm :title="$t('interflow.confirmDelete')" @confirm="handleDelete(row)"><template #reference><el-button type="danger" link>{{ $t('common.delete') }}</el-button></template></el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('interflow.name')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('interflow.enterName')" />
        </el-form-item>
        <el-form-item :label="$t('interflow.sourceProduct')" prop="source_product_id">
          <el-select v-model="formData.source_product_id" :placeholder="$t('interflow.selectSourceProduct')" filterable style="width: 100%">
            <el-option v-for="item in productOptions" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('interflow.targetProduct')" prop="target_product_id">
          <el-select v-model="formData.target_product_id" :placeholder="$t('interflow.selectTargetProduct')" filterable style="width: 100%">
            <el-option v-for="item in productOptions" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('interflow.type')" prop="type">
          <el-select v-model="formData.type" :placeholder="$t('interflow.selectType')" style="width: 100%">
            <el-option :label="$t('interflow.typeMap.auto_provision')" value="auto_provision" />
            <el-option :label="$t('interflow.typeMap.manual_link')" value="manual_link" />
            <el-option :label="$t('interflow.typeMap.sync_config')" value="sync_config" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.description')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('interflow.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('interflow.sort')" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('interflow.status')" prop="status">
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
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'Interflows' })

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref($t('interflow.addInterflow'))
const formRef = ref<FormInstance>()
const productOptions = ref<any[]>([])

const searchForm = reactive({ keyword: '', status: undefined as number | undefined })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const formData = reactive({ id: undefined as number | undefined, name: '', source_product_id: undefined as number | undefined, target_product_id: undefined as number | undefined, type: 'auto_provision', description: '', sort: 0, status: 1 })

const formRules: FormRules = {
  name: [{ required: true, message: () => $t('interflow.enterName'), trigger: 'blur' }],
  source_product_id: [{ required: true, message: () => $t('interflow.selectSourceProduct'), trigger: 'change' }],
  target_product_id: [{ required: true, message: () => $t('interflow.selectTargetProduct'), trigger: 'change' }],
  type: [{ required: true, message: () => $t('interflow.selectType'), trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try { const data = await request.get({ url: '/api/admin/interflows', params: { page: pagination.page, page_size: pagination.page_size, ...searchForm } }); tableData.value = data.list || data || []; pagination.total = data.total || 0 } catch { ElMessage.error($t('interflow.fetchFailed')) } finally { loading.value = false }
}

const fetchProducts = async () => { try { const data = await request.get({ url: '/api/admin/products', params: { page_size: 999 } }); productOptions.value = data.list || data || [] } catch {} }

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.status = undefined; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleAdd = () => { dialogTitle.value = $t('interflow.addInterflow'); formData.id = undefined; formData.name = ''; formData.source_product_id = undefined; formData.target_product_id = undefined; formData.type = 'auto_provision'; formData.description = ''; formData.sort = 0; formData.status = 1; dialogVisible.value = true }
const handleEdit = (row: any) => { dialogTitle.value = $t('interflow.editInterflow'); Object.assign(formData, row); dialogVisible.value = true }

const handleDelete = async (row: any) => { try { await request.del({ url: `/api/admin/interflows/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchData() } catch { ElMessage.error($t('common.deleteFailed')) } }

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) await request.put({ url: `/api/admin/interflows/${formData.id}`, params: formData })
      else await request.post({ url: '/api/admin/interflows', params: formData })
      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess')); dialogVisible.value = false; fetchData()
    } catch { ElMessage.error($t('common.operateFailed')) } finally { submitLoading.value = false }
  })
}

onMounted(() => { fetchData(); fetchProducts() })
</script>

<style scoped lang="scss">
.interflows-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
