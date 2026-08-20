<template>
  <div class="nav-groups-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('navGroups.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('navGroups.addGroup') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('navGroups.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('navGroups.groupNamePlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('navGroups.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('navGroups.all')" clearable>
            <el-option :label="$t('navGroups.enabled')" :value="1" />
            <el-option :label="$t('navGroups.disabled')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('navGroups.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('navGroups.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border row-key="id">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" :label="$t('navGroups.groupName')" width="180" />
        <el-table-column prop="slug" :label="$t('navGroups.slug')" width="150" />
        <el-table-column prop="description" :label="$t('navGroups.description')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="product_count" :label="$t('navGroups.productCount')" width="110" align="center">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.product_count || 0 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" :label="$t('navGroups.sort')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('navGroups.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('navGroups.enabled') : $t('navGroups.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('navGroups.createdAt')" width="180" />
        <el-table-column :label="$t('navGroups.actions')" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('navGroups.edit') }}</el-button>
            <el-button type="success" link @click="handleManageProducts(row)">{{ $t('navGroups.productAssociation') }}</el-button>
            <el-popconfirm :title="$t('navGroups.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('navGroups.delete') }}</el-button>
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
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('navGroups.groupName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('navGroups.enterGroupName')" />
        </el-form-item>
        <el-form-item :label="$t('navGroups.slug')" prop="slug">
          <el-input v-model="formData.slug" :placeholder="$t('navGroups.slugPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('navGroups.description')">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('navGroups.enterGroupDesc')" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('navGroups.sort')">
              <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('navGroups.status')">
              <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" :active-text="$t('navGroups.enabled')" :inactive-text="$t('navGroups.disabled')" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('navGroups.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('navGroups.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 产品关联对话框 -->
    <el-dialog v-model="productDialogVisible" :title="$t('navGroups.productAssocManage')" width="700px">
      <div v-loading="productLoading">
        <div class="product-header">
          <span>{{ $t('navGroups.group') }}<strong>{{ currentGroup.name }}</strong></span>
          <el-button type="primary" size="small" @click="handleAddProduct">{{ $t('navGroups.addProduct') }}</el-button>
        </div>
        <el-table :data="groupProducts" style="width: 100%" border size="small">
          <el-table-column prop="id" :label="$t('navGroups.productId')" width="80" />
          <el-table-column prop="name" :label="$t('navGroups.productName')" min-width="200" />
          <el-table-column prop="price" :label="$t('navGroups.price')" width="100">
            <template #default="{ row }">¥{{ Number(row.price || 0).toFixed(2) }}</template>
          </el-table-column>
          <el-table-column :label="$t('navGroups.actions')" width="100">
            <template #default="{ row }">
              <el-popconfirm :title="$t('navGroups.confirmRemoveProduct')" @confirm="handleRemoveProduct(row)">
                <template #reference>
                  <el-button type="danger" link size="small">{{ $t('navGroups.remove') }}</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <template #footer>
        <el-button @click="productDialogVisible = false">{{ $t('navGroups.close') }}</el-button>
      </template>
    </el-dialog>

    <!-- 添加产品选择对话框 -->
    <el-dialog v-model="addProductVisible" :title="$t('navGroups.selectProduct')" width="500px">
      <el-form :inline="true" class="search-form" style="margin-bottom: 12px">
        <el-form-item>
          <el-input v-model="productSearchKeyword" :placeholder="$t('navGroups.searchProductName')" clearable @input="handleSearchProducts" />
        </el-form-item>
      </el-form>
      <el-table :data="filteredProducts" style="width: 100%" border size="small" @row-click="handleSelectProduct">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" :label="$t('navGroups.productName')" min-width="200" />
        <el-table-column prop="price" :label="$t('navGroups.price')" width="100">
          <template #default="{ row }">¥{{ Number(row.price || 0).toFixed(2) }}</template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="addProductVisible = false">{{ $t('navGroups.cancel') }}</el-button>
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

defineOptions({ name: 'NavGroupsManage' })

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref($t('navGroups.addNavGroup'))
const formRef = ref<FormInstance>()

const productDialogVisible = ref(false)
const productLoading = ref(false)
const addProductVisible = ref(false)
const productSearchKeyword = ref('')
const groupProducts = ref<any[]>([])
const allProducts = ref<any[]>([])
const currentGroup = reactive({ id: 0, name: '' })

const searchForm = reactive({ keyword: '', status: undefined as number | undefined })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  slug: '',
  description: '',
  sort: 0,
  status: 1
})

const formRules: FormRules = {
  name: [{ required: true, message: $t('navGroups.enterGroupName'), trigger: 'blur' }],
  slug: [{ required: true, message: $t('navGroups.enterSlug'), trigger: 'blur' }]
}

const filteredProducts = ref<any[]>([])

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/nav-groups',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error($t('navGroups.fetchListFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { Object.assign(searchForm, { keyword: '', status: undefined }); handleSearch() }

const resetForm = () => {
  formData.id = undefined
  formData.name = ''
  formData.slug = ''
  formData.description = ''
  formData.sort = 0
  formData.status = 1
}

const handleAdd = () => {
  dialogTitle.value = $t('navGroups.addNavGroup')
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('navGroups.editNavGroup')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/nav-groups/${row.id}` })
    ElMessage.success($t('navGroups.deleteSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('navGroups.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/nav-groups/${formData.id}`, params: { ...formData } })
      } else {
        await request.post({ url: '/api/admin/nav-groups', params: { ...formData } })
      }
      ElMessage.success(formData.id ? $t('navGroups.updateSuccess') : $t('navGroups.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('navGroups.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleManageProducts = async (row: any) => {
  currentGroup.id = row.id
  currentGroup.name = row.name
  productDialogVisible.value = true
  productLoading.value = true
  try {
    const data = await request.get({ url: `/api/admin/nav-groups/${row.id}/products` })
    groupProducts.value = data || []
  } catch (error) {
    ElMessage.error($t('navGroups.fetchAssocProductsFailed'))
  } finally {
    productLoading.value = false
  }
}

const handleAddProduct = async () => {
  addProductVisible.value = true
  productSearchKeyword.value = ''
  try {
    const data = await request.get({ url: '/api/admin/products', params: { page_size: 100 } })
    allProducts.value = data.list || data || []
    filteredProducts.value = allProducts.value
  } catch (error) {
    ElMessage.error($t('navGroups.fetchProductsFailed'))
  }
}

const handleSearchProducts = () => {
  const keyword = productSearchKeyword.value.toLowerCase()
  filteredProducts.value = keyword
    ? allProducts.value.filter((p: any) => p.name?.toLowerCase().includes(keyword))
    : allProducts.value
}

const handleSelectProduct = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/nav-groups/${currentGroup.id}/products`, params: { product_id: row.id } })
    ElMessage.success($t('navGroups.addProductAssocSuccess'))
    addProductVisible.value = false
    handleManageProducts({ id: currentGroup.id, name: currentGroup.name })
  } catch (error) {
    ElMessage.error($t('navGroups.addFailed'))
  }
}

const handleRemoveProduct = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/nav-groups/${currentGroup.id}/products/${row.id}` })
    ElMessage.success($t('navGroups.removeSuccess'))
    handleManageProducts({ id: currentGroup.id, name: currentGroup.name })
  } catch (error) {
    ElMessage.error($t('navGroups.removeFailed'))
  }
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.nav-groups-page {
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

.product-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}
</style>
