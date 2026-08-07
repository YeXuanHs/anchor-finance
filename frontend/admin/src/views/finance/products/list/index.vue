<template>
  <div class="product-list-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>产品列表</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加产品
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="产品名称/描述" clearable />
        </el-form-item>
        <el-form-item label="分组">
          <el-select v-model="searchForm.group_id" placeholder="全部" clearable>
            <el-option v-for="group in groups" :key="group.id" :label="group.name" :value="group.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="上架" :value="1" />
            <el-option label="下架" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="产品名称" min-width="180" />
        <el-table-column prop="group_name" label="所属分组" width="120" />
        <el-table-column prop="price" label="价格" width="120">
          <template #default="{ row }">
            ¥{{ row.price?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="billing_cycle" label="计费周期" width="100">
          <template #default="{ row }">
            {{ getBillingCycleText(row.billing_cycle) }}
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="80" />
        <el-table-column prop="sales" label="销量" width="80" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '上架' : '下架' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button :type="row.status === 1 ? 'warning' : 'success'" link @click="handleToggleStatus(row)">
              {{ row.status === 1 ? '下架' : '上架' }}
            </el-button>
            <el-popconfirm title="确定删除该产品吗？" @confirm="handleDelete(row)">
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="产品名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入产品名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属分组" prop="group_id">
              <el-select v-model="formData.group_id" placeholder="请选择分组" style="width: 100%">
                <el-option v-for="group in groups" :key="group.id" :label="group.name" :value="group.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="价格" prop="price">
              <el-input-number v-model="formData.price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="计费周期" prop="billing_cycle">
              <el-select v-model="formData.billing_cycle" placeholder="请选择计费周期" style="width: 100%">
                <el-option label="月付" value="monthly" />
                <el-option label="季付" value="quarterly" />
                <el-option label="半年付" value="semi_annually" />
                <el-option label="年付" value="annually" />
                <el-option label="一次性" value="onetime" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="库存" prop="stock">
              <el-input-number v-model="formData.stock" :min="-1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="产品描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="4" placeholder="请输入产品描述" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" active-text="上架" inactive-text="下架" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
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

defineOptions({ name: 'ProductList' })

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  keyword: '',
  group_id: undefined as number | undefined,
  status: undefined as number | undefined
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref<any[]>([])

// 分组数据
const groups = ref<any[]>([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('添加产品')
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  group_id: undefined as number | undefined,
  price: 0,
  billing_cycle: 'monthly',
  stock: -1,
  sort: 0,
  description: '',
  status: 1
})

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入产品名称', trigger: 'blur' },
    { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  group_id: [
    { required: true, message: '请选择所属分组', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入价格', trigger: 'blur' }
  ],
  billing_cycle: [
    { required: true, message: '请选择计费周期', trigger: 'change' }
  ]
}

// 获取计费周期文本
const getBillingCycleText = (cycle: string) => {
  const map: Record<string, string> = {
    monthly: '月付',
    quarterly: '季付',
    semi_annually: '半年付',
    annually: '年付',
    onetime: '一次性'
  }
  return map[cycle] || cycle
}

// 获取产品列表
const fetchProducts = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/products',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取产品列表失败:', error)
    ElMessage.error('获取产品列表失败')
  } finally {
    loading.value = false
  }
}

// 获取分组列表
const fetchGroups = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/product-groups'
    })
    groups.value = data || []
  } catch (error) {
    console.error('获取分组列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchProducts()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.group_id = undefined
  searchForm.status = undefined
  handleSearch()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = '添加产品'
  formData.id = undefined
  formData.name = ''
  formData.group_id = undefined
  formData.price = 0
  formData.billing_cycle = 'monthly'
  formData.stock = -1
  formData.sort = 0
  formData.description = ''
  formData.status = 1
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = '编辑产品'
  Object.assign(formData, row)
  dialogVisible.value = true
}

// 切换上下架状态
const handleToggleStatus = async (row: any) => {
  try {
    await request.put({
      url: `/api/admin/products/${row.id}`,
      params: { status: row.status === 1 ? 0 : 1 }
    })
    ElMessage.success(row.status === 1 ? '下架成功' : '上架成功')
    fetchProducts()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/products/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchProducts()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/products/${formData.id}` : '/api/admin/products'

      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchProducts()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchProducts()
}

// 页码变化
const handlePageChange = () => {
  fetchProducts()
}

onMounted(() => {
  fetchProducts()
  fetchGroups()
})
</script>

<style scoped lang="scss">
.product-list-page {
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
</style>
