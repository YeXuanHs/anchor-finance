<template>
  <div class="config-groups-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>产品配置组</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加配置组
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="配置组名称" clearable />
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
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="配置组名称" min-width="200" />
        <el-table-column prop="code" label="配置组编码" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="product_count" label="关联产品数" width="110" align="center" />
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleManageProducts(row)">管理产品</el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该配置组吗？" @confirm="handleDelete(row)">
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
        <el-form-item label="配置组名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入配置组名称" />
        </el-form-item>
        <el-form-item label="配置组编码" prop="code">
          <el-input v-model="formData.code" placeholder="请输入配置组编码" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 管理产品对话框 -->
    <el-dialog v-model="productDialogVisible" title="管理关联产品" width="700px" destroy-on-close>
      <div class="product-transfer">
        <el-transfer
          v-model="selectedProductIds"
          :data="allProducts"
          :titles="['可选产品', '已关联产品']"
          :props="{ key: 'id', label: 'name' }"
          filterable
          filter-placeholder="搜索产品"
        />
      </div>
      <template #footer>
        <el-button @click="productDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveProducts" :loading="productSaving">保存</el-button>
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

defineOptions({ name: 'ProductConfigGroups' })

const loading = ref(false)
const submitLoading = ref(false)
const productSaving = ref(false)
const dialogVisible = ref(false)
const productDialogVisible = ref(false)
const dialogTitle = ref('添加配置组')
const formRef = ref<FormInstance>()
const currentGroupId = ref<number>(0)

const searchForm = reactive({ keyword: '', status: undefined as number | undefined })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const allProducts = ref<any[]>([])
const selectedProductIds = ref<number[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  code: '',
  description: '',
  sort: 0,
  status: 1
})

const formRules: FormRules = {
  name: [
    { required: true, message: '请输入配置组名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入配置组编码', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '只能包含字母、数字和下划线', trigger: 'blur' }
  ]
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/product-config-groups',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || data || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取配置组列表失败')
  } finally {
    loading.value = false
  }
}

const fetchAllProducts = async () => {
  try {
    const data = await request.get({ url: '/api/admin/products', params: { page_size: 999 } })
    allProducts.value = data.list || data || []
  } catch (error) {
    console.error('获取产品列表失败:', error)
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.status = undefined; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleAdd = () => {
  dialogTitle.value = '添加配置组'
  formData.id = undefined; formData.name = ''; formData.code = ''
  formData.description = ''; formData.sort = 0; formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑配置组'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/product-config-groups/${row.id}` })
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
        await request.put({ url: `/api/admin/product-config-groups/${formData.id}`, params: formData })
      } else {
        await request.post({ url: '/api/admin/product-config-groups', params: formData })
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
  currentGroupId.value = row.id
  try {
    const data = await request.get({ url: `/api/admin/product-config-groups/${row.id}/products` })
    selectedProductIds.value = (data.list || data || []).map((p: any) => p.id || p)
    productDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取关联产品失败')
  }
}

const handleSaveProducts = async () => {
  productSaving.value = true
  try {
    await request.put({
      url: `/api/admin/product-config-groups/${currentGroupId.value}/products`,
      params: { product_ids: selectedProductIds.value }
    })
    ElMessage.success('关联产品更新成功')
    productDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    productSaving.value = false
  }
}

onMounted(() => { fetchData(); fetchAllProducts() })
</script>

<style scoped lang="scss">
.config-groups-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.product-transfer {
  display: flex;
  justify-content: center;
}
</style>
