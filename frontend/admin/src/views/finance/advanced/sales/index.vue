<template>
  <div class="sales-manage-page">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-card-content">
            <div class="stat-info">
              <div class="stat-title">今日销售额</div>
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
              <div class="stat-title">本月销售额</div>
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
              <div class="stat-title">今日订单数</div>
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
              <div class="stat-title">销售员数量</div>
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
        <el-tab-pane label="销售记录" name="records">
          <el-form :inline="true" :model="recordSearchForm" class="search-form">
            <el-form-item label="销售员">
              <el-select v-model="recordSearchForm.sales_id" placeholder="全部" clearable filterable>
                <el-option v-for="item in salesUsers" :key="item.id" :label="item.username" :value="item.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="日期范围">
              <el-date-picker
                v-model="recordSearchForm.date_range"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleRecordSearch">搜索</el-button>
              <el-button @click="handleRecordReset">重置</el-button>
            </el-form-item>
          </el-form>

          <el-table :data="recordTableData" v-loading="recordLoading" style="width: 100%" border>
            <el-table-column prop="id" label="ID" width="70" align="center" />
            <el-table-column prop="sales_username" label="销售员" width="100" />
            <el-table-column prop="client_name" label="客户" width="120" />
            <el-table-column prop="product_name" label="产品" min-width="150" show-overflow-tooltip />
            <el-table-column prop="order_no" label="订单号" width="170" />
            <el-table-column prop="amount" label="金额" width="110" align="right">
              <template #default="{ row }">
                <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="commission" label="佣金" width="100" align="right">
              <template #default="{ row }">
                <span class="commission-text">¥{{ formatAmount(row.commission) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="成交时间" width="170" />
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
        <el-tab-pane label="销售员管理" name="users">
          <div class="toolbar">
            <el-button type="primary" @click="handleAddSalesUser">
              <el-icon><Plus /></el-icon>
              添加销售员
            </el-button>
          </div>

          <el-table :data="salesUserTableData" v-loading="salesUserLoading" style="width: 100%" border>
            <el-table-column prop="id" label="ID" width="70" align="center" />
            <el-table-column prop="username" label="用户名" width="120" />
            <el-table-column prop="real_name" label="姓名" width="100" />
            <el-table-column prop="phone" label="手机号" width="130" />
            <el-table-column prop="total_sales" label="累计销售额" width="130" align="right">
              <template #default="{ row }">
                <span class="amount-text">¥{{ formatAmount(row.total_sales) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="order_count" label="订单数" width="80" align="center" />
            <el-table-column prop="commission_rate" label="佣金比例" width="100" align="center">
              <template #default="{ row }">
                <span class="commission-text">{{ row.commission_rate }}%</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80" align="center">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                  {{ row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="170" />
            <el-table-column label="操作" width="150" fixed="right" align="center">
              <template #default="{ row }">
                <el-button type="primary" link @click="handleEditSalesUser(row)">编辑</el-button>
                <el-popconfirm title="确定删除该销售员吗？" @confirm="handleDeleteSalesUser(row)">
                  <template #reference>
                    <el-button type="danger" link>删除</el-button>
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
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="请输入用户名" :disabled="!!formData.id" />
        </el-form-item>
        <el-form-item label="姓名" prop="real_name">
          <el-input v-model="formData.real_name" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="formData.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!formData.id">
          <el-input v-model="formData.password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="佣金比例" prop="commission_rate">
          <el-input-number v-model="formData.commission_rate" :min="0" :max="100" :precision="1" :step="0.5" style="width: 200px" />
          <span style="margin-left: 8px; color: #909399">%</span>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入备注" />
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
import { Plus, TrendCharts, DataLine, Document, User } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

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
const dialogTitle = ref('添加销售员')
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
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '长度在 3 到 50 个字符', trigger: 'blur' }
  ],
  real_name: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' }
  ],
  commission_rate: [
    { required: true, message: '请输入佣金比例', trigger: 'blur' }
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
    ElMessage.error('获取销售记录失败')
  } finally {
    recordLoading.value = false
  }
}

// 获取销售员用户列表（下拉用）
const fetchSalesUsers = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/sales/users',
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
      url: '/api/admin/sales/users',
      params: {
        page: salesUserPagination.page,
        page_size: salesUserPagination.page_size
      }
    })
    salesUserTableData.value = data.list || []
    salesUserPagination.total = data.total || 0
  } catch (error) {
    console.error('获取销售员列表失败:', error)
    ElMessage.error('获取销售员列表失败')
  } finally {
    salesUserLoading.value = false
  }
}

// 标签切换
const handleTabChange = (tab: string) => {
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
  dialogTitle.value = '添加销售员'
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
  dialogTitle.value = '编辑销售员'
  Object.assign(formData, row)
  formData.password = ''
  dialogVisible.value = true
}

// 删除销售员
const handleDeleteSalesUser = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/sales/users/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchSalesUserTable()
    fetchSalesUsers()
    fetchStatistics()
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
      const url = formData.id ? `/api/admin/sales/users/${formData.id}` : '/api/admin/sales/users'

      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchSalesUserTable()
      fetchSalesUsers()
      fetchStatistics()
    } catch (error) {
      ElMessage.error('操作失败')
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
