<template>
  <div class="sales-manage-page">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-card-content">
            <div class="stat-info">
              <div class="stat-title">{{ $t('page.sales.todaySales') }}</div>
              <div class="stat-value">¥{{ formatAmount(statistics.today_sales) }}</div>
            </div>
            <el-icon class="stat-icon sales-icon"><TrendCharts /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-card-content">
            <div class="stat-info">
              <div class="stat-title">{{ $t('page.sales.monthSales') }}</div>
              <div class="stat-value">¥{{ formatAmount(statistics.month_sales) }}</div>
            </div>
            <el-icon class="stat-icon orders-icon"><DataLine /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-card-content">
            <div class="stat-info">
              <div class="stat-title">{{ $t('page.sales.todayOrders') }}</div>
              <div class="stat-value">{{ statistics.today_orders || 0 }}</div>
            </div>
            <el-icon class="stat-icon count-icon"><Document /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-card-content">
            <div class="stat-info">
              <div class="stat-title">{{ $t('page.sales.salesCount') }}</div>
              <div class="stat-value">{{ statistics.sales_count || 0 }}</div>
            </div>
            <el-icon class="stat-icon users-icon"><User /></el-icon>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 主要内容区 -->
    <el-card shadow="never">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <!-- 销售记录 -->
        <el-tab-pane :label="$t('page.sales.salesRecords')" name="records">
          <el-form :inline="true" :model="recordSearchForm" class="search-form">
            <el-form-item :label="$t('page.sales.salesperson')">
              <el-select v-model="recordSearchForm.sales_id" :placeholder="$t('page.sales.all')" clearable filterable>
                <el-option v-for="item in salesUsers" :key="item.id" :label="item.username" :value="item.id" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('page.sales.dateRange')">
              <el-date-picker
                v-model="recordSearchForm.date_range"
                type="daterange"
                :range-separator="$t('page.sales.rangeSeparator')"
                :start-placeholder="$t('page.sales.startDate')"
                :end-placeholder="$t('page.sales.endDate')"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleRecordSearch">{{ $t('page.sales.search') }}</el-button>
              <el-button @click="handleRecordReset">{{ $t('page.sales.reset') }}</el-button>
            </el-form-item>
          </el-form>

          <el-table :data="recordTableData" v-loading="recordLoading" style="width: 100%" border>
            <el-table-column prop="id" label="ID" width="70" align="center" />
            <el-table-column prop="sales_username" :label="$t('page.sales.salesperson')" width="100" />
            <el-table-column prop="client_name" :label="$t('page.sales.client')" width="120" />
            <el-table-column prop="product_name" :label="$t('page.sales.product')" min-width="150" show-overflow-tooltip />
            <el-table-column prop="order_no" :label="$t('page.sales.orderNo')" width="170" />
            <el-table-column prop="amount" :label="$t('page.sales.amount')" width="110" align="right">
              <template #default="{ row }">
                <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="commission" :label="$t('page.sales.commission')" width="100" align="right">
              <template #default="{ row }">
                <span class="commission-text">¥{{ formatAmount(row.commission) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" :label="$t('page.sales.dealTime')" width="170" />
          </el-table>

          <div class="pagination-container">
            <el-pagination
              v-model:current-page="recordPagination.page"
              v-model:page-size="recordPagination.page_size"
              :page-sizes="[10, 20, 50, 100]"
              :total="recordPagination.total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleRecordSizeChange"
              @current-change="handleRecordPageChange"
            />
          </div>
        </el-tab-pane>

        <!-- 销售员管理 -->
        <el-tab-pane :label="$t('page.sales.salesUserManagement')" name="users">
          <div class="toolbar">
            <el-button type="primary" @click="handleAddSalesUser">
              <el-icon><Plus /></el-icon>
              {{ $t('page.sales.addSalesperson') }}
            </el-button>
          </div>

          <el-table :data="salesUserTableData" v-loading="salesUserLoading" style="width: 100%" border>
            <el-table-column prop="id" label="ID" width="70" align="center" />
            <el-table-column prop="username" :label="$t('page.sales.username')" width="120" />
            <el-table-column prop="real_name" :label="$t('page.sales.realName')" width="100" />
            <el-table-column prop="phone" :label="$t('page.sales.phone')" width="130" />
            <el-table-column prop="total_sales" :label="$t('page.sales.totalSales')" width="130" align="right">
              <template #default="{ row }">
                <span class="amount-text">¥{{ formatAmount(row.total_sales) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="order_count" :label="$t('page.sales.orderCount')" width="80" align="center" />
            <el-table-column prop="commission_rate" :label="$t('page.sales.commissionRate')" width="100" align="center">
              <template #default="{ row }">
                <span class="commission-text">{{ row.commission_rate }}%</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" :label="$t('page.sales.status')" width="80" align="center">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                  {{ row.status === 1 ? $t('page.sales.enabled') : $t('page.sales.disabled') }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" :label="$t('page.sales.createTime')" width="170" />
            <el-table-column :label="$t('page.sales.actions')" width="150" fixed="right" align="center">
              <template #default="{ row }">
                <el-button type="primary" link @click="handleEditSalesUser(row)">{{ $t('page.sales.edit') }}</el-button>
                <el-popconfirm :title="$t('page.sales.confirmDeleteSalesUser')" @confirm="handleDeleteSalesUser(row)">
                  <template #reference>
                    <el-button type="danger" link>{{ $t('page.sales.delete') }}</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination-container">
            <el-pagination
              v-model:current-page="salesUserPagination.page"
              v-model:page-size="salesUserPagination.page_size"
              :page-sizes="[10, 20, 50, 100]"
              :total="salesUserPagination.total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSalesUserSizeChange"
              @current-change="handleSalesUserPageChange"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 添加/编辑销售员对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('page.sales.username')" prop="username">
          <el-input v-model="formData.username" :placeholder="$t('page.sales.usernamePlaceholder')" :disabled="!!formData.id" />
        </el-form-item>
        <el-form-item :label="$t('page.sales.realName')" prop="real_name">
          <el-input v-model="formData.real_name" :placeholder="$t('page.sales.realNamePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('page.sales.phone')" prop="phone">
          <el-input v-model="formData.phone" :placeholder="$t('page.sales.phonePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('page.sales.email')" prop="email">
          <el-input v-model="formData.email" :placeholder="$t('page.sales.emailPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('page.sales.password')" prop="password" v-if="!formData.id">
          <el-input v-model="formData.password" type="password" :placeholder="$t('page.sales.passwordPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('page.sales.commissionRate')" prop="commission_rate">
          <el-input-number v-model="formData.commission_rate" :min="0" :max="100" :precision="1" :step="0.5" style="width: 200px" />
          <span style="margin-left: 8px; color: #909399">%</span>
        </el-form-item>
        <el-form-item :label="$t('page.sales.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item :label="$t('page.sales.notes')" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="3" :placeholder="$t('page.sales.notesPlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('page.sales.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('page.sales.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, TrendCharts, DataLine, Document, User } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

// 当前标签
const activeTab = ref('records')

// 统计数据
const statistics = reactive({
  today_sales: 0,
  month_sales: 0,
  today_orders: 0,
  sales_count: 0
})

// 销售记录相关
const recordLoading = ref(false)
const recordSearchForm = reactive({
  sales_id: undefined as number | undefined,
  date_range: [] as string[]
})
const recordPagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})
const recordTableData = ref<any[]>([])

// 销售员用户列表（用于下拉）
const salesUsers = ref<any[]>([])

// 销售员管理相关
const salesUserLoading = ref(false)
const salesUserPagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})
const salesUserTableData = ref<any[]>([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref($t('page.sales.addSalesperson'))
const formRef = ref<FormInstance>()
const submitLoading = ref(false)

// 表单数据
const formData = reactive({
  id: undefined as number | undefined,
  username: '',
  real_name: '',
  phone: '',
  email: '',
  password: '',
  commission_rate: 5,
  status: 1,
  remark: ''
})

// 表单验证规则
const formRules: FormRules = {
  username: [
    { required: true, message: $t('page.sales.msgEnterUsername'), trigger: 'blur' },
    { min: 3, max: 50, message: $t('page.sales.msgUsernameLength'), trigger: 'blur' }
  ],
  real_name: [
    { required: true, message: $t('page.sales.msgEnterRealName'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: $t('page.sales.msgEnterPassword'), trigger: 'blur' },
    { min: 6, message: $t('page.sales.msgPasswordLength'), trigger: 'blur' }
  ],
  commission_rate: [
    { required: true, message: $t('page.sales.msgEnterCommissionRate'), trigger: 'blur' }
  ]
}

// 格式化金额
const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

// 获取统计数据
const fetchStatistics = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/sales/statistics'
    })
    Object.assign(statistics, data)
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 获取销售记录
const fetchRecords = async () => {
  recordLoading.value = true
  try {
    const params: any = {
      page: recordPagination.page,
      page_size: recordPagination.page_size,
      sales_id: recordSearchForm.sales_id
    }
    if (recordSearchForm.date_range?.length === 2) {
      params.start_date = recordSearchForm.date_range[0]
      params.end_date = recordSearchForm.date_range[1]
    }
    const data = await request.get({
      url: '/api/admin/sales/records',
      params
    })
    recordTableData.value = data.list || []
    recordPagination.total = data.total || 0
  } catch (error) {
    console.error('获取销售记录失败:', error)
    ElMessage.error($t('page.sales.msgFetchRecordsFailed'))
  } finally {
    recordLoading.value = false
  }
}

// 获取销售员用户列表（下拉用）
const fetchSalesUsers = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/sales/admin-list',
      params: { page: 1, page_size: 9999 }
    })
    salesUsers.value = data.list || []
  } catch (error) {
    console.error('获取销售员列表失败:', error)
  }
}

// 获取销售员管理列表
const fetchSalesUserTable = async () => {
  salesUserLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/sales',
      params: {
        page: salesUserPagination.page,
        page_size: salesUserPagination.page_size
      }
    })
    salesUserTableData.value = data.list || []
    salesUserPagination.total = data.total || 0
  } catch (error) {
    console.error('获取销售员列表失败:', error)
    ElMessage.error($t('page.sales.msgFetchSalesUsersFailed'))
  } finally {
    salesUserLoading.value = false
  }
}

// 标签切换
const handleTabChange = (tab: string | number) => {
  if (tab === 'records') {
    fetchRecords()
  } else if (tab === 'users') {
    fetchSalesUserTable()
  }
}

// 销售记录搜索
const handleRecordSearch = () => {
  recordPagination.page = 1
  fetchRecords()
}

// 销售记录重置
const handleRecordReset = () => {
  recordSearchForm.sales_id = undefined
  recordSearchForm.date_range = []
  handleRecordSearch()
}

// 销售记录分页
const handleRecordSizeChange = () => {
  recordPagination.page = 1
  fetchRecords()
}

const handleRecordPageChange = () => {
  fetchRecords()
}

// 添加销售员
const handleAddSalesUser = () => {
  dialogTitle.value = $t('page.sales.addSalesperson')
  formData.id = undefined
  formData.username = ''
  formData.real_name = ''
  formData.phone = ''
  formData.email = ''
  formData.password = ''
  formData.commission_rate = 5
  formData.status = 1
  formData.remark = ''
  dialogVisible.value = true
}

// 编辑销售员
const handleEditSalesUser = (row: any) => {
  dialogTitle.value = $t('page.sales.editSalesperson')
  Object.assign(formData, row)
  formData.password = ''
  dialogVisible.value = true
}

// 删除销售员
const handleDeleteSalesUser = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/sales/${row.id}`
    })
    ElMessage.success($t('page.sales.msgDeleteSuccess'))
    fetchSalesUserTable()
    fetchSalesUsers()
    fetchStatistics()
  } catch (error) {
    ElMessage.error($t('page.sales.msgDeleteFailed'))
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/sales/${formData.id}` : '/api/admin/sales'

      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }

      ElMessage.success(formData.id ? $t('page.sales.msgUpdateSuccess') : $t('page.sales.msgAddSuccess'))
      dialogVisible.value = false
      fetchSalesUserTable()
      fetchSalesUsers()
      fetchStatistics()
    } catch (error) {
      ElMessage.error($t('page.sales.msgOperationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

// 销售员分页
const handleSalesUserSizeChange = () => {
  salesUserPagination.page = 1
  fetchSalesUserTable()
}

const handleSalesUserPageChange = () => {
  fetchSalesUserTable()
}

onMounted(() => {
  fetchStatistics()
  fetchRecords()
  fetchSalesUsers()
})
</script>

<style scoped lang="scss">
.sales-manage-page {
  padding: 20px;
}

.stat-cards {
  margin-bottom: 20px;
}

.stat-card {
  .stat-card-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .stat-info {
    .stat-title {
      font-size: 14px;
      color: #909399;
      margin-bottom: 8px;
    }

    .stat-value {
      font-size: 24px;
      font-weight: 600;
      color: #303133;
    }
  }

  .stat-icon {
    font-size: 48px;
    opacity: 0.8;

    &.sales-icon {
      color: #409eff;
    }

    &.orders-icon {
      color: #67c23a;
    }

    &.count-icon {
      color: #e6a23c;
    }

    &.users-icon {
      color: #f56c6c;
    }
  }
}

.search-form {
  margin-bottom: 20px;
}

.toolbar {
  margin-bottom: 16px;
}

.amount-text {
  font-weight: 600;
  color: #f56c6c;
}

.commission-text {
  font-weight: 600;
  color: #409eff;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
