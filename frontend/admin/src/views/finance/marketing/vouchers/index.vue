<template>
  <div class="vouchers-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>代金券管理</span>
          <div>
            <el-button type="success" @click="handleGrant">
              <el-icon><Promotion /></el-icon>
              发放代金券
            </el-button>
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              添加代金券
            </el-button>
          </div>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="代金券名称/编码" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="正常" :value="1" />
            <el-option label="已禁用" :value="0" />
            <el-option label="已过期" :value="2" />
            <el-option label="已用完" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="客户">
          <el-input v-model="searchForm.client_name" placeholder="客户名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="代金券名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="code" label="代金码" width="140" />
        <el-table-column prop="amount" label="面值" width="110" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="balance" label="余额" width="110" align="right">
          <template #default="{ row }">
            <span :class="{ 'balance-empty': row.balance <= 0 }">
              ¥{{ formatAmount(row.balance) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="client_name" label="所属客户" width="120" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.client_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="order_no" label="关联订单" width="150" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.order_no || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="end_time" label="有效期" width="170">
          <template #default="{ row }">
            {{ row.end_time || '永久' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="getVoucherStatusType(row.status)" size="small">
              {{ getVoucherStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button
              v-if="row.status === 1 && row.client_id"
              type="success"
              link
              @click="handleGrantSingle(row)"
            >
              发放
            </el-button>
            <el-popconfirm title="确定删除该代金券吗？" @confirm="handleDelete(row)">
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
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="110px">
        <el-form-item label="代金券名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入代金券名称" />
        </el-form-item>
        <el-form-item label="代金码" prop="code">
          <el-input v-model="formData.code" placeholder="留空则自动生成">
            <template #append>
              <el-button @click="handleGenerateCode">随机生成</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="面值(元)" prop="amount">
          <el-input-number
            v-model="formData.amount"
            :min="0.01"
            :max="99999"
            :precision="2"
            :step="50"
            controls-position="right"
          />
        </el-form-item>
        <el-form-item label="关联客户" prop="client_id">
          <el-select
            v-model="formData.client_id"
            filterable
            remote
            :remote-method="searchClients"
            :loading="clientSearchLoading"
            placeholder="输入搜索客户"
            clearable
            style="width: 100%"
          >
            <el-option
              v-for="client in clientOptions"
              :key="client.id"
              :label="client.username"
              :value="client.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="关联订单" prop="order_id">
          <el-select
            v-model="formData.order_id"
            filterable
            remote
            :remote-method="searchOrders"
            :loading="orderSearchLoading"
            placeholder="输入搜索订单号（可选）"
            clearable
            style="width: 100%"
          >
            <el-option
              v-for="order in orderOptions"
              :key="order.id"
              :label="order.order_no"
              :value="order.id"
            />
          </el-select>
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
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入代金券描述（可选）"
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

    <!-- 批量发放对话框 -->
    <el-dialog v-model="grantDialogVisible" title="发放代金券" width="550px" destroy-on-close>
      <el-form :model="grantForm" :rules="grantFormRules" ref="grantFormRef" label-width="110px">
        <el-form-item label="代金券" prop="voucher_id">
          <el-select v-model="grantForm.voucher_id" placeholder="请选择代金券" style="width: 100%">
            <el-option
              v-for="voucher in availableVouchers"
              :key="voucher.id"
              :label="`${voucher.name} (¥${formatAmount(voucher.amount)})`"
              :value="voucher.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="发放客户" prop="client_ids">
          <el-select
            v-model="grantForm.client_ids"
            multiple
            filterable
            remote
            :remote-method="searchClientsForGrant"
            :loading="clientSearchLoading"
            placeholder="输入搜索客户，可多选"
            style="width: 100%"
          >
            <el-option
              v-for="client in clientOptions"
              :key="client.id"
              :label="client.username"
              :value="client.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="发放数量" prop="count">
          <el-input-number
            v-model="grantForm.count"
            :min="1"
            :max="100"
            controls-position="right"
          />
          <span class="form-tip">每个客户发放的代金券数量</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="grantDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleGrantSubmit" :loading="grantLoading">确认发放</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, Promotion } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)
const grantLoading = ref(false)
const clientSearchLoading = ref(false)
const orderSearchLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: undefined as number | undefined,
  client_name: ''
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref<any[]>([])

// 客户选项
const clientOptions = ref<any[]>([])

// 订单选项
const orderOptions = ref<any[]>([])

// 可用代金券列表（用于发放）
const availableVouchers = ref<any[]>([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('添加代金券')
const formRef = ref<FormInstance>()

// 发放对话框
const grantDialogVisible = ref(false)
const grantFormRef = ref<FormInstance>()

// 表单数据
const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  code: '',
  amount: 0,
  client_id: undefined as number | undefined,
  order_id: undefined as number | undefined,
  date_range: [] as string[],
  description: '',
  status: 1
})

// 发放表单
const grantForm = reactive({
  voucher_id: undefined as number | undefined,
  client_ids: [] as number[],
  count: 1
})

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入代金券名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  amount: [
    { required: true, message: '请输入面值', trigger: 'blur' }
  ]
}

// 发放表单验证规则
const grantFormRules: FormRules = {
  voucher_id: [
    { required: true, message: '请选择代金券', trigger: 'change' }
  ],
  client_ids: [
    { required: true, type: 'array', min: 1, message: '请选择至少一个客户', trigger: 'change' }
  ],
  count: [
    { required: true, message: '请输入发放数量', trigger: 'blur' }
  ]
}

// 代金券状态文本
const getVoucherStatusText = (status: number) => {
  const map: Record<number, string> = {
    0: '已禁用',
    1: '正常',
    2: '已过期',
    3: '已用完'
  }
  return map[status] || '未知'
}

// 代金券状态类型
const getVoucherStatusType = (status: number) => {
  const map: Record<number, string> = {
    0: 'info',
    1: 'success',
    2: 'warning',
    3: 'danger'
  }
  return (map[status] || 'info') as any
}

// 格式化金额
const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

// 生成随机代金码
const handleGenerateCode = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let code = 'V'
  for (let i = 0; i < 9; i++) {
    code += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  formData.code = code
}

// 搜索客户
const searchClients = async (query: string) => {
  if (!query) {
    clientOptions.value = []
    return
  }
  clientSearchLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/users',
      params: { keyword: query, page_size: 20 }
    })
    clientOptions.value = data.list || data || []
  } catch (error) {
    console.error('搜索客户失败:', error)
  } finally {
    clientSearchLoading.value = false
  }
}

// 搜索客户（发放用）
const searchClientsForGrant = async (query: string) => {
  await searchClients(query)
}

// 搜索订单
const searchOrders = async (query: string) => {
  if (!query) {
    orderOptions.value = []
    return
  }
  orderSearchLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/orders',
      params: { order_no: query, page_size: 20 }
    })
    orderOptions.value = data.list || data || []
  } catch (error) {
    console.error('搜索订单失败:', error)
  } finally {
    orderSearchLoading.value = false
  }
}

// 获取代金券列表
const fetchVouchers = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      keyword: searchForm.keyword || undefined,
      status: searchForm.status,
      client_name: searchForm.client_name || undefined
    }
    const data = await request.get({
      url: '/api/admin/vouchers',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取代金券列表失败:', error)
    ElMessage.error('获取代金券列表失败')
  } finally {
    loading.value = false
  }
}

// 获取可用代金券（发放用）
const fetchAvailableVouchers = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/vouchers',
      params: { status: 1, page_size: 100 }
    })
    availableVouchers.value = data.list || data || []
  } catch (error) {
    console.error('获取可用代金券失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchVouchers()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = undefined
  searchForm.client_name = ''
  handleSearch()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = '添加代金券'
  formData.id = undefined
  formData.name = ''
  formData.code = ''
  formData.amount = 0
  formData.client_id = undefined
  formData.order_id = undefined
  formData.date_range = []
  formData.description = ''
  formData.status = 1
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = '编辑代金券'
  Object.assign(formData, {
    ...row,
    date_range: row.start_time && row.end_time ? [row.start_time, row.end_time] : []
  })
  if (row.client_id && row.client_name) {
    clientOptions.value = [{ id: row.client_id, username: row.client_name }]
  }
  if (row.order_id && row.order_no) {
    orderOptions.value = [{ id: row.order_id, order_no: row.order_no }]
  }
  dialogVisible.value = true
}

// 单个发放
const handleGrantSingle = (row: any) => {
  grantForm.voucher_id = row.id
  grantForm.client_ids = row.client_id ? [row.client_id] : []
  grantForm.count = 1
  fetchAvailableVouchers()
  grantDialogVisible.value = true
}

// 批量发放
const handleGrant = () => {
  grantForm.voucher_id = undefined
  grantForm.client_ids = []
  grantForm.count = 1
  fetchAvailableVouchers()
  grantDialogVisible.value = true
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/vouchers/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchVouchers()
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

      const url = formData.id ? `/api/admin/vouchers/${formData.id}` : '/api/admin/vouchers'

      if (formData.id) {
        await request.put({ url, params: submitData })
      } else {
        await request.post({ url, params: submitData })
      }

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchVouchers()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

// 提交发放
const handleGrantSubmit = async () => {
  if (!grantFormRef.value) return

  await grantFormRef.value.validate(async (valid) => {
    if (!valid) return

    grantLoading.value = true
    try {
      await request.post({
        url: '/api/admin/vouchers/grant',
        params: grantForm
      })
      ElMessage.success('发放成功')
      grantDialogVisible.value = false
      fetchVouchers()
    } catch (error) {
      ElMessage.error('发放失败')
    } finally {
      grantLoading.value = false
    }
  })
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchVouchers()
}

// 页码变化
const handlePageChange = () => {
  fetchVouchers()
}

onMounted(() => {
  fetchVouchers()
})
</script>

<style scoped lang="scss">
.vouchers-page {
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

.balance-empty {
  color: #c0c4cc;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.form-tip {
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}
</style>
