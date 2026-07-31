<template>
  <div class="coupons-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>优惠券管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加优惠券
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="优惠券名称/编码" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="满减券" value="fixed" />
            <el-option label="折扣券" value="percent" />
            <el-option label="抵扣券" value="deduction" />
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
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="优惠券名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="code" label="优惠码" width="140" />
        <el-table-column prop="type" label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">
              {{ getTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="面值/折扣" width="120" align="center">
          <template #default="{ row }">
            <span v-if="row.type === 'percent'">{{ row.value }}折</span>
            <span v-else class="amount-text">¥{{ formatAmount(row.value) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="min_amount" label="最低消费" width="110" align="center">
          <template #default="{ row }">
            {{ row.min_amount > 0 ? `¥${formatAmount(row.min_amount)}` : '无限制' }}
          </template>
        </el-table-column>
        <el-table-column label="使用情况" width="130" align="center">
          <template #default="{ row }">
            {{ row.used_count }} / {{ row.total_count || '不限' }}
          </template>
        </el-table-column>
        <el-table-column prop="start_time" label="开始时间" width="170" />
        <el-table-column prop="end_time" label="结束时间" width="170" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleToggleStatus(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该优惠券吗？" @confirm="handleDelete(row)">
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="650px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="110px">
        <el-form-item label="优惠券名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入优惠券名称" />
        </el-form-item>
        <el-form-item label="优惠码" prop="code">
          <el-input v-model="formData.code" placeholder="留空则自动生成">
            <template #append>
              <el-button @click="handleGenerateCode">随机生成</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="formData.type">
            <el-radio value="fixed">满减券</el-radio>
            <el-radio value="percent">折扣券</el-radio>
            <el-radio value="deduction">抵扣券</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="formData.type === 'percent' ? '折扣(%)' : '面值(元)'" prop="value">
          <el-input-number
            v-model="formData.value"
            :min="formData.type === 'percent' ? 1 : 0.01"
            :max="formData.type === 'percent' ? 99 : 99999"
            :precision="formData.type === 'percent' ? 0 : 2"
            :step="formData.type === 'percent' ? 1 : 10"
            controls-position="right"
          />
        </el-form-item>
        <el-form-item label="最低消费" prop="min_amount">
          <el-input-number
            v-model="formData.min_amount"
            :min="0"
            :precision="2"
            :step="100"
            controls-position="right"
            placeholder="0 表示无限制"
          />
        </el-form-item>
        <el-form-item label="发放总量" prop="total_count">
          <el-input-number
            v-model="formData.total_count"
            :min="0"
            :step="100"
            controls-position="right"
            placeholder="0 表示不限量"
          />
        </el-form-item>
        <el-form-item label="每人限领" prop="per_user_limit">
          <el-input-number
            v-model="formData.per_user_limit"
            :min="1"
            :max="99"
            controls-position="right"
          />
        </el-form-item>
        <el-form-item label="有效期" prop="date_range">
          <el-date-picker
            v-model="formData.date_range"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="适用产品" prop="product_ids">
          <el-select
            v-model="formData.product_ids"
            multiple
            filterable
            placeholder="留空表示全部产品适用"
            style="width: 100%"
          >
            <el-option
              v-for="product in productOptions"
              :key="product.id"
              :label="product.name"
              :value="product.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入优惠券描述（可选）"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
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

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  keyword: '',
  type: '' as string,
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

// 产品选项
const productOptions = ref<any[]>([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('添加优惠券')
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  code: '',
  type: 'fixed',
  value: 0,
  min_amount: 0,
  total_count: 0,
  per_user_limit: 1,
  date_range: [] as string[],
  product_ids: [] as number[],
  description: '',
  status: 1
})

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入优惠券名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择优惠券类型', trigger: 'change' }
  ],
  value: [
    { required: true, message: '请输入面值或折扣', trigger: 'blur' }
  ]
}

// 类型文本映射
const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    fixed: '满减券',
    percent: '折扣券',
    deduction: '抵扣券'
  }
  return map[type] || type
}

// 类型标签颜色
const getTypeTag = (type: string) => {
  const map: Record<string, string> = {
    fixed: 'primary',
    percent: 'success',
    deduction: 'warning'
  }
  return (map[type] || 'info') as any
}

// 格式化金额
const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

// 生成随机优惠码
const handleGenerateCode = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let code = ''
  for (let i = 0; i < 8; i++) {
    code += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  formData.code = code
}

// 获取优惠券列表
const fetchCoupons = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      keyword: searchForm.keyword || undefined,
      type: searchForm.type || undefined,
      status: searchForm.status
    }
    const data = await request.get({
      url: '/api/admin/coupons',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取优惠券列表失败:', error)
    ElMessage.error('获取优惠券列表失败')
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
    productOptions.value = data.list || data || []
  } catch (error) {
    console.error('获取产品列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchCoupons()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.type = ''
  searchForm.status = undefined
  handleSearch()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = '添加优惠券'
  formData.id = undefined
  formData.name = ''
  formData.code = ''
  formData.type = 'fixed'
  formData.value = 0
  formData.min_amount = 0
  formData.total_count = 0
  formData.per_user_limit = 1
  formData.date_range = []
  formData.product_ids = []
  formData.description = ''
  formData.status = 1
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = '编辑优惠券'
  Object.assign(formData, {
    ...row,
    date_range: row.start_time && row.end_time ? [row.start_time, row.end_time] : [],
    product_ids: row.product_ids || []
  })
  dialogVisible.value = true
}

// 切换状态
const handleToggleStatus = async (row: any) => {
  try {
    await request.put({
      url: `/api/admin/coupons/${row.id}`,
      params: { status: row.status }
    })
    ElMessage.success(row.status === 1 ? '已启用' : '已禁用')
  } catch (error) {
    row.status = row.status === 1 ? 0 : 1
    ElMessage.error('操作失败')
  }
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/coupons/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchCoupons()
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
      const submitData: any = { ...formData }
      if (submitData.date_range?.length === 2) {
        submitData.start_time = submitData.date_range[0]
        submitData.end_time = submitData.date_range[1]
      }
      delete submitData.date_range

      const url = formData.id ? `/api/admin/coupons/${formData.id}` : '/api/admin/coupons'

      if (formData.id) {
        await request.put({ url, params: submitData })
      } else {
        await request.post({ url, params: submitData })
      }

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchCoupons()
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
  fetchCoupons()
}

// 页码变化
const handlePageChange = () => {
  fetchCoupons()
}

onMounted(() => {
  fetchCoupons()
  fetchProducts()
})
</script>

<style scoped lang="scss">
.coupons-page {
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

.amount-text {
  font-weight: 600;
  color: #f56c6c;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
