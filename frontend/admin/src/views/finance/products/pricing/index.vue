<template>
  <div class="pricing-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>定价管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加定价规则
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="产品">
          <el-select v-model="searchForm.product_id" placeholder="全部产品" clearable filterable style="width: 240px">
            <el-option v-for="product in products" :key="product.id" :label="product.name" :value="product.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="计费周期">
          <el-select v-model="searchForm.billing_cycle" placeholder="全部" clearable>
            <el-option label="月付" value="monthly" />
            <el-option label="季付" value="quarterly" />
            <el-option label="半年付" value="semi_annually" />
            <el-option label="年付" value="annually" />
            <el-option label="一次性" value="onetime" />
          </el-select>
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
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="product_name" label="产品名称" min-width="180" />
        <el-table-column prop="billing_cycle" label="计费周期" width="100">
          <template #default="{ row }">
            {{ getBillingCycleText(row.billing_cycle) }}
          </template>
        </el-table-column>
        <el-table-column prop="price" label="价格" width="120">
          <template #default="{ row }">
            <span class="price-text">¥{{ row.price?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="setup_fee" label="设置费" width="100">
          <template #default="{ row }">
            ¥{{ row.setup_fee?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="currency" label="货币" width="80" />
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button :type="row.status === 1 ? 'warning' : 'success'" link @click="handleToggleStatus(row)">
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-popconfirm title="确定删除该定价规则吗？" @confirm="handleDelete(row)">
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="产品" prop="product_id">
          <el-select v-model="formData.product_id" placeholder="请选择产品" filterable style="width: 100%">
            <el-option v-for="product in products" :key="product.id" :label="product.name" :value="product.id" />
          </el-select>
        </el-form-item>
        <el-row :gutter="20">
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
          <el-col :span="12">
            <el-form-item label="货币" prop="currency">
              <el-select v-model="formData.currency" placeholder="请选择货币" style="width: 100%">
                <el-option label="CNY - 人民币" value="CNY" />
                <el-option label="USD - 美元" value="USD" />
                <el-option label="EUR - 欧元" value="EUR" />
                <el-option label="GBP - 英镑" value="GBP" />
                <el-option label="JPY - 日元" value="JPY" />
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
            <el-form-item label="设置费" prop="setup_fee">
              <el-input-number v-model="formData.setup_fee" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入备注信息" />
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

defineOptions({ name: 'PricingManagement' })

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  product_id: undefined as number | undefined,
  billing_cycle: undefined as string | undefined,
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

// 产品数据
const products = ref<any[]>([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('添加定价规则')
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive({
  id: undefined as number | undefined,
  product_id: undefined as number | undefined,
  billing_cycle: 'monthly',
  price: 0,
  setup_fee: 0,
  currency: 'CNY',
  sort: 0,
  status: 1,
  remark: ''
})

// 表单验证规则
const formRules: FormRules = {
  product_id: [
    { required: true, message: '请选择产品', trigger: 'change' }
  ],
  billing_cycle: [
    { required: true, message: '请选择计费周期', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入价格', trigger: 'blur' }
  ],
  currency: [
    { required: true, message: '请选择货币', trigger: 'change' }
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

// 获取定价列表
const fetchPricingList = async () => {
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
    console.error('获取定价列表失败:', error)
    ElMessage.error('获取定价列表失败')
  } finally {
    loading.value = false
  }
}

// 获取产品列表
const fetchProducts = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/products'
    })
    products.value = data.list || data || []
  } catch (error) {
    console.error('获取产品列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchPricingList()
}

// 重置
const handleReset = () => {
  searchForm.product_id = undefined
  searchForm.billing_cycle = undefined
  searchForm.status = undefined
  handleSearch()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = '添加定价规则'
  formData.id = undefined
  formData.product_id = undefined
  formData.billing_cycle = 'monthly'
  formData.price = 0
  formData.setup_fee = 0
  formData.currency = 'CNY'
  formData.sort = 0
  formData.status = 1
  formData.remark = ''
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = '编辑定价规则'
  Object.assign(formData, row)
  dialogVisible.value = true
}

// 切换状态
const handleToggleStatus = async (row: any) => {
  try {
    await request.put({
      url: `/api/admin/products/${row.id}`,
      params: { status: row.status === 1 ? 0 : 1 }
    })
    ElMessage.success(row.status === 1 ? '已禁用' : '已启用')
    fetchPricingList()
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
    fetchPricingList()
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
      fetchPricingList()
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
  fetchPricingList()
}

// 页码变化
const handlePageChange = () => {
  fetchPricingList()
}

onMounted(() => {
  fetchPricingList()
  fetchProducts()
})
</script>

<style scoped lang="scss">
.pricing-page {
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

.price-text {
  color: #f56c6c;
  font-weight: 600;
}
</style>
