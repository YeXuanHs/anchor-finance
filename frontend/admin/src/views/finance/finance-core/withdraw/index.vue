<template>
  <div class="withdraw-page art-full-height">
    <!-- 搜索栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="申请编号">
          <el-input v-model="searchForm.withdraw_no" placeholder="请输入申请编号" clearable />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="searchForm.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待审核" value="pending" />
            <el-option label="审核通过" value="approved" />
            <el-option label="审核拒绝" value="rejected" />
            <el-option label="已打款" value="paid" />
          </el-select>
        </el-form-item>
        <el-form-item label="提现方式">
          <el-select v-model="searchForm.method" placeholder="全部" clearable>
            <el-option label="支付宝" value="alipay" />
            <el-option label="微信" value="wechat" />
            <el-option label="银行卡" value="bank" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never" class="art-table-card">
      <!-- 表格头部 -->
      <template #header>
        <div class="card-header">
          <span>提现申请列表</span>
          <el-space>
            <el-tag type="warning">待审核: {{ stats.pending || 0 }}</el-tag>
            <el-tag type="primary">今日提现: ¥{{ stats.today_amount?.toFixed(2) || '0.00' }}</el-tag>
          </el-space>
        </div>
      </template>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="withdraw_no" label="申请编号" width="150" />
        <el-table-column prop="user_id" label="用户ID" width="80" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="amount" label="提现金额" width="120">
          <template #default="{ row }">
            <span class="text-primary">¥{{ row.amount?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="fee" label="手续费" width="100">
          <template #default="{ row }">
            ¥{{ row.fee?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="actual_amount" label="实际到账" width="120">
          <template #default="{ row }">
            <span class="text-success">¥{{ row.actual_amount?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="method" label="提现方式" width="100">
          <template #default="{ row }">
            <el-tag :type="getMethodTag(row.method)" size="small">
              {{ getMethodText(row.method) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="account_info" label="收款账号" min-width="180" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTag(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <template v-if="row.status === 'pending'">
              <el-button type="success" link @click="handleReview(row, 'approved')">通过</el-button>
              <el-button type="danger" link @click="handleReview(row, 'rejected')">拒绝</el-button>
            </template>
            <template v-if="row.status === 'approved'">
              <el-button type="primary" link @click="handlePay(row)">打款</el-button>
            </template>
            <el-button type="info" link @click="handleViewDetail(row)">详情</el-button>
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

    <!-- 审核对话框 -->
    <el-dialog v-model="reviewDialogVisible" :title="reviewForm.action === 'approved' ? '审核通过' : '审核拒绝'" width="500px">
      <el-form :model="reviewForm" :rules="reviewRules" ref="reviewFormRef" label-width="100px">
        <el-form-item label="申请编号">
          <el-input :value="reviewForm.withdraw_no" disabled />
        </el-form-item>
        <el-form-item label="提现金额">
          <el-input :value="'¥' + (reviewForm.amount?.toFixed(2) || '0.00')" disabled />
        </el-form-item>
        <el-form-item label="审核意见" prop="reviewRemark">
          <el-input v-model="reviewForm.review_remark" type="textarea" placeholder="请输入审核意见" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reviewDialogVisible = false">取消</el-button>
        <el-button :type="reviewForm.action === 'approved' ? 'success' : 'danger'" @click="handleSubmitReview" :loading="submitLoading">
          {{ reviewForm.action === 'approved' ? '确认通过' : '确认拒绝' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 打款对话框 -->
    <el-dialog v-model="payDialogVisible" title="确认打款" width="500px">
      <el-alert title="请确认已完成线下转账后再进行打款操作" type="warning" :closable="false" show-icon style="margin-bottom: 20px;" />
      <el-form :model="payForm" ref="payFormRef" label-width="100px">
        <el-form-item label="申请编号">
          <el-input :value="payForm.withdraw_no" disabled />
        </el-form-item>
        <el-form-item label="提现金额">
          <el-input :value="'¥' + (payForm.amount?.toFixed(2) || '0.00')" disabled />
        </el-form-item>
        <el-form-item label="收款账号">
          <el-input :value="payForm.account_info" disabled />
        </el-form-item>
        <el-form-item label="打款凭证">
          <el-input v-model="payForm.payment凭证" placeholder="请输入转账单号或凭证" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="payForm.remark" type="textarea" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="payDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitPay" :loading="submitLoading">确认打款</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="提现详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="申请编号">{{ detailData.withdraw_no }}</el-descriptions-item>
        <el-descriptions-item label="用户">{{ detailData.username }} (ID: {{ detailData.user_id }})</el-descriptions-item>
        <el-descriptions-item label="提现金额">¥{{ detailData.amount?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="手续费">¥{{ detailData.fee?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="实际到账">¥{{ detailData.actual_amount?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="提现方式">{{ getMethodText(detailData.method) }}</el-descriptions-item>
        <el-descriptions-item label="收款账号" :span="2">{{ detailData.account_info }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusTag(detailData.status)">{{ getStatusText(detailData.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="申请时间">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="审核人">{{ detailData.reviewer || '-' }}</el-descriptions-item>
        <el-descriptions-item label="审核时间">{{ detailData.reviewed_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="审核意见" :span="2">{{ detailData.review_remark || '-' }}</el-descriptions-item>
        <el-descriptions-item label="打款时间">{{ detailData.paid_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="打款凭证">{{ detailData.payment凭证 || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'WithdrawManage' })

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)

// 统计数据
const stats = ref({
  pending: 0,
  today_amount: 0
})

// 搜索表单
const searchForm = reactive({
  withdraw_no: '',
  username: '',
  status: undefined as string | undefined,
  method: undefined as string | undefined,
  date_range: [] as string[]
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref([])

// 审核对话框
const reviewDialogVisible = ref(false)
const reviewFormRef = ref<FormInstance>()
const reviewForm = reactive({
  id: 0,
  withdraw_no: '',
  amount: 0,
  action: 'approved' as 'approved' | 'rejected',
  review_remark: ''
})

// 打款对话框
const payDialogVisible = ref(false)
const payFormRef = ref<FormInstance>()
const payForm = reactive({
  id: 0,
  withdraw_no: '',
  amount: 0,
  account_info: '',
  payment凭证: '',
  remark: ''
})

// 详情对话框
const detailDialogVisible = ref(false)
const detailData = ref<any>({})

// 审核表单验证规则
const reviewRules: FormRules = {
  review_remark: [
    { required: true, message: '请输入审核意见', trigger: 'blur' }
  ]
}

// 获取提现方式标签
const getMethodTag = (method: string) => {
  const map: Record<string, any> = {
    alipay: 'primary',
    wechat: 'success',
    bank: 'warning'
  }
  return map[method] || 'info'
}

// 获取提现方式文本
const getMethodText = (method: string) => {
  const map: Record<string, string> = {
    alipay: '支付宝',
    wechat: '微信',
    bank: '银行卡'
  }
  return map[method] || '未知'
}

// 获取状态标签
const getStatusTag = (status: string) => {
  const map: Record<string, any> = {
    pending: 'warning',
    approved: 'primary',
    rejected: 'danger',
    paid: 'success'
  }
  return map[status] || 'info'
}

// 获取状态文本
const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待审核',
    approved: '审核通过',
    rejected: '审核拒绝',
    paid: '已打款'
  }
  return map[status] || '未知'
}

// 获取提现列表
const fetchWithdraws = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size
    }
    if (searchForm.withdraw_no) params.withdraw_no = searchForm.withdraw_no
    if (searchForm.username) params.username = searchForm.username
    if (searchForm.status) params.status = searchForm.status
    if (searchForm.method) params.method = searchForm.method
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }

    const data = await request.get({
      url: '/api/admin/affiliate/withdraw-records',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
    stats.value = {
      pending: data.pending || 0,
      today_amount: data.today_amount || 0
    }
  } catch (error) {
    console.error('获取提现列表失败:', error)
    ElMessage.error('获取提现列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchWithdraws()
}

// 重置
const handleReset = () => {
  searchForm.withdraw_no = ''
  searchForm.username = ''
  searchForm.status = undefined
  searchForm.method = undefined
  searchForm.date_range = []
  handleSearch()
}

// 审核
const handleReview = (row: any, action: 'approved' | 'rejected') => {
  reviewForm.id = row.id
  reviewForm.withdraw_no = row.withdraw_no
  reviewForm.amount = row.amount
  reviewForm.action = action
  reviewForm.review_remark = ''
  reviewDialogVisible.value = true
}

// 打款
const handlePay = (row: any) => {
  payForm.id = row.id
  payForm.withdraw_no = row.withdraw_no
  payForm.amount = row.amount
  payForm.account_info = row.account_info
  payForm.payment凭证 = ''
  payForm.remark = ''
  payDialogVisible.value = true
}

// 查看详情
const handleViewDetail = (row: any) => {
  detailData.value = row
  detailDialogVisible.value = true
}

// 提交审核
const handleSubmitReview = async () => {
  if (!reviewFormRef.value) return

  await reviewFormRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      await request.post({
        url: `/api/admin/affiliate/withdraws/${reviewForm.id}/process`,
        params: {
          action: reviewForm.action,
          review_remark: reviewForm.review_remark
        },
        showSuccessMessage: true
      })
      reviewDialogVisible.value = false
      fetchWithdraws()
    } catch (error) {
      console.error('审核失败:', error)
    } finally {
      submitLoading.value = false
    }
  })
}

// 提交打款
const handleSubmitPay = async () => {
  await ElMessageBox.confirm('确认已完成线下转账?', '确认打款', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  })

  submitLoading.value = true
  try {
    await request.post({
      url: `/api/admin/affiliate/withdraws/${payForm.id}/process`,
      params: {
        action: 'paid',
        payment凭证: payForm.payment凭证,
        remark: payForm.remark
      },
      showSuccessMessage: true
    })
    payDialogVisible.value = false
    fetchWithdraws()
  } catch (error) {
    console.error('打款失败:', error)
  } finally {
    submitLoading.value = false
  }
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchWithdraws()
}

// 页码变化
const handlePageChange = () => {
  fetchWithdraws()
}

onMounted(() => {
  fetchWithdraws()
})
</script>

<style scoped lang="scss">
.withdraw-page {
  padding: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  .el-form-item {
    margin-bottom: 0;
  }
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.text-primary {
  color: #409eff;
  font-weight: 600;
}

.text-success {
  color: #67c23a;
  font-weight: 600;
}
</style>
