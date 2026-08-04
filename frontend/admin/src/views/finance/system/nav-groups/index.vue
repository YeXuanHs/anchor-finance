<template>
  <div class="nav-groups-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>导航分组管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加分组
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="分组名称" clearable />
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
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border row-key="id">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="分组名称" width="180" />
        <el-table-column prop="slug" label="标识" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="product_count" label="关联产品数" width="110" align="center">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.product_count || 0 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="success" link @click="handleManageProducts(row)">产品关联</el-button>
            <el-popconfirm title="确定删除该导航分组吗？" @confirm="handleDelete(row)">
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
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="分组名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入分组名称" />
        </el-form-item>
        <el-form-item label="标识" prop="slug">
          <el-input v-model="formData.slug" placeholder="英文标识，如: cloud-servers" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入分组描述" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="排序">
              <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 产品关联对话框 -->
    <el-dialog v-model="productDialogVisible" title="产品关联管理" width="700px">
      <div v-loading="productLoading">
        <div class="product-header">
          <span>分组：<strong>{{ currentGroup.name }}</strong></span>
          <el-button type="primary" size="small" @click="handleAddProduct">添加产品</el-button>
        </div>
        <el-table :data="groupProducts" style="width: 100%" border size="small">
          <el-table-column prop="id" label="产品ID" width="80" />
          <el-table-column prop="name" label="产品名称" min-width="200" />
          <el-table-column prop="price" label="价格" width="100">
            <template #default="{ row }">¥{{ Number(row.price || 0).toFixed(2) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-popconfirm title="确定移除该产品关联吗？" @confirm="handleRemoveProduct(row)">
                <template #reference>
                  <el-button type="danger" link size="small">移除</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <template #footer>
        <el-button @click="productDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 添加产品选择对话框 -->
    <el-dialog v-model="addProductVisible" title="选择产品" width="500px">
      <el-form :inline="true" class="search-form" style="margin-bottom: 12px">
        <el-form-item>
          <el-input v-model="productSearchKeyword" placeholder="搜索产品名称" clearable @input="handleSearchProducts" />
        </el-form-item>
      </el-form>
      <el-table :data="filteredProducts" style="width: 100%" border size="small" @row-click="handleSelectProduct">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="产品名称" min-width="200" />
        <el-table-column prop="price" label="价格" width="100">
          <template #default="{ row }">¥{{ Number(row.price || 0).toFixed(2) }}</template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="addProductVisible = false">取消</el-button>
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

defineOptions({ name: 'NavGroupsManage' })

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('添加导航分组')
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
  name: [{ required: true, message: '请输入分组名称', trigger: 'blur' }],
  slug: [{ required: true, message: '请输入标识', trigger: 'blur' }]
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
    ElMessage.error('获取导航分组列表失败')
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
  dialogTitle.value = '添加导航分组'
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑导航分组'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/nav-groups/${row.id}` })
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
        await request.put({ url: `/api/admin/nav-groups/${formData.id}`, params: { ...formData } })
      } else {
        await request.post({ url: '/api/admin/nav-groups', params: { ...formData } })
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

const handleManageProducts = async (row: any) => {
  currentGroup.id = row.id
  currentGroup.name = row.name
  productDialogVisible.value = true
  productLoading.value = true
  try {
    const data = await request.get({ url: `/api/admin/nav-groups/${row.id}/products` })
    groupProducts.value = data || []
  } catch (error) {
    ElMessage.error('获取关联产品失败')
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
    ElMessage.error('获取产品列表失败')
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
    ElMessage.success('添加产品关联成功')
    addProductVisible.value = false
    handleManageProducts({ id: currentGroup.id, name: currentGroup.name })
  } catch (error) {
    ElMessage.error('添加失败')
  }
}

const handleRemoveProduct = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/nav-groups/${currentGroup.id}/products/${row.id}` })
    ElMessage.success('移除成功')
    handleManageProducts({ id: currentGroup.id, name: currentGroup.name })
  } catch (error) {
    ElMessage.error('移除失败')
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
