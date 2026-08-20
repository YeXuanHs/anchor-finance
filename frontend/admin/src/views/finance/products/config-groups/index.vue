<template>
  <div class="config-groups-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('configGroup.title') }}</span>
          <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('configGroup.addGroup') }}</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('configGroup.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('configGroup.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('configGroup.status')">
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
        <el-table-column prop="name" :label="$t('configGroup.groupName')" min-width="200" />
        <el-table-column prop="code" :label="$t('configGroup.groupCode')" width="150" />
        <el-table-column prop="description" :label="$t('common.description')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="product_count" :label="$t('configGroup.productCount')" width="110" align="center" />
        <el-table-column prop="sort" :label="$t('configGroup.sort')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('configGroup.status')" width="80" align="center">
          <template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">{{ row.status === 1 ? $t('common.enable') : $t('common.disable') }}</el-tag></template>
        </el-table-column>
        <el-table-column :label="$t('configGroup.operations')" width="250" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleManageProducts(row)">{{ $t('configGroup.manageProducts') }}</el-button>
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-popconfirm :title="$t('configGroup.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference><el-button type="danger" link>{{ $t('common.delete') }}</el-button></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('configGroup.groupName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('configGroup.enterGroupName')" />
        </el-form-item>
        <el-form-item :label="$t('configGroup.groupCode')" prop="code">
          <el-input v-model="formData.code" :placeholder="$t('configGroup.enterGroupCode')" />
        </el-form-item>
        <el-form-item :label="$t('common.description')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('configGroup.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('configGroup.sort')" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('configGroup.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" :active-text="$t('common.enable')" :inactive-text="$t('common.disable')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="productDialogVisible" :title="$t('configGroup.manageProducts')" width="700px" destroy-on-close>
      <div class="product-transfer">
        <el-transfer v-model="selectedProductIds" :data="allProducts" :titles="[$t('configGroup.availableProducts'), $t('configGroup.linkedProducts')]" :props="{ key: 'id', label: 'name' }" filterable :filter-placeholder="$t('configGroup.searchProduct')" />
      </div>
      <template #footer>
        <el-button @click="productDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveProducts" :loading="productSaving">{{ $t('common.save') }}</el-button>
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

defineOptions({ name: 'ProductConfigGroups' })

const loading = ref(false)
const submitLoading = ref(false)
const productSaving = ref(false)
const dialogVisible = ref(false)
const productDialogVisible = ref(false)
const dialogTitle = ref($t('configGroup.addGroup'))
const formRef = ref<FormInstance>()
const currentGroupId = ref<number>(0)

const searchForm = reactive({ keyword: '', status: undefined as number | undefined })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const allProducts = ref<any[]>([])
const selectedProductIds = ref<number[]>([])

const formData = reactive({ id: undefined as number | undefined, name: '', code: '', description: '', sort: 0, status: 1 })

const formRules: FormRules = {
  name: [{ required: true, message: () => $t('configGroup.enterGroupName'), trigger: 'blur' }],
  code: [{ required: true, message: () => $t('configGroup.enterGroupCode'), trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try { const data = await request.get({ url: '/api/admin/product-config-groups', params: { page: pagination.page, page_size: pagination.page_size, ...searchForm } }); tableData.value = data.list || data || []; pagination.total = data.total || 0 } catch { ElMessage.error($t('configGroup.fetchFailed')) } finally { loading.value = false }
}

const fetchAllProducts = async () => { try { const data = await request.get({ url: '/api/admin/products', params: { page_size: 999 } }); allProducts.value = data.list || data || [] } catch {} }

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.status = undefined; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleAdd = () => { dialogTitle.value = $t('configGroup.addGroup'); formData.id = undefined; formData.name = ''; formData.code = ''; formData.description = ''; formData.sort = 0; formData.status = 1; dialogVisible.value = true }
const handleEdit = (row: any) => { dialogTitle.value = $t('configGroup.editGroup'); Object.assign(formData, row); dialogVisible.value = true }

const handleDelete = async (row: any) => {
  try { await request.del({ url: `/api/admin/product-config-groups/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchData() } catch { ElMessage.error($t('common.deleteFailed')) }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) await request.put({ url: `/api/admin/product-config-groups/${formData.id}`, params: formData })
      else await request.post({ url: '/api/admin/product-config-groups', params: formData })
      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess')); dialogVisible.value = false; fetchData()
    } catch { ElMessage.error($t('common.operateFailed')) } finally { submitLoading.value = false }
  })
}

const handleManageProducts = async (row: any) => {
  currentGroupId.value = row.id
  try { const data = await request.get({ url: `/api/admin/product-config-groups/${row.id}/products` }); selectedProductIds.value = (data.list || data || []).map((p: any) => p.id || p); productDialogVisible.value = true } catch { ElMessage.error($t('configGroup.fetchProductsFailed')) }
}

const handleSaveProducts = async () => {
  productSaving.value = true
  try { await request.put({ url: `/api/admin/product-config-groups/${currentGroupId.value}/products`, params: { product_ids: selectedProductIds.value } }); ElMessage.success($t('configGroup.productsUpdated')); productDialogVisible.value = false; fetchData() } catch { ElMessage.error($t('common.saveFailed')) } finally { productSaving.value = false }
}

onMounted(() => { fetchData(); fetchAllProducts() })
</script>

<style scoped lang="scss">
.config-groups-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.product-transfer { display: flex; justify-content: center; }
</style>
