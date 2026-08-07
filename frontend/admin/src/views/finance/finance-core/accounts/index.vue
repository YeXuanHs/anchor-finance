<template>
  <div class="accounts-page">
    <!-- 搜索栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.user_id" placeholder="用户ID" clearable />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="searchForm.username" placeholder="用户名" clearable />
        </el-form-item>
        <el-form-item label="交易类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="充值" value="recharge" />
            <el-option label="消费" value="payment" />
            <el-option label="退款" value="refund" />
            <el-option label="提现" value="withdraw" />
            <el-option label="转入" value="transfer_in" />
            <el-option label="转出" value="transfer_out" />
          </el-select>
        </el-form-item>
        <el-form-item label="支付方式">
          <el-select v-model="searchForm.gateway" placeholder="全部" clearable>
            <el-option label="支付宝" value="alipay" />
            <el-option label="微信" value="wechat" />
            <el-option label="余额" value="balance" />
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
      <template #header>
        <div class="card-header">
          <span>交易流水列表</span>
          <el-space>
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              创建流水
            </el-button>
            <el-button type="success" @click="handleExport">
              <el-icon><Download /></el-icon>
              导出
            </el-button>
          </el-space>
        </div>
      </template>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="transaction_no" label="交易号" width="180" show-overflow-tooltip />
        <el-table-column prop="user_id" label="用户ID" width="90" align="center" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="type" label="交易类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="typeTagMap[row.type]" size="small">{{ typeTextMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            <span :class="isIncome(row.type) ? 'text-green' : 'text-red'">
              {{ isIncome(row.type) ? '+' : '-' }}¥{{ Number(row.amount).toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="gateway" label="支付方式" width="100" align="center">
          <template #default="{ row }">
            {{ gatewayTextMap[row.gateway] || row.gateway || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="balance_after" label="余额" width="110">
          <template #default="{ row }">
            ¥{{ Number(row.balance_after || 0).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="180" show-overflow-tooltip />
        <el-table-column prop="invoice_id" label="关联发票" width="100" align="center">
          <template #default="{ row }">
            <el-button v-if="row.invoice_id" type="primary" link size="small" @click="handleViewInvoice(row)">查看</el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="交易时间" width="180" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleDetail(row)">详情</el-button>
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

    <!-- 创建流水对话框 -->
    <el-dialog v-model="dialogVisible" title="创建交易流水" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="用户ID" prop="user_id">
          <el-input v-model="formData.user_id" placeholder="请输入用户ID" />
        </el-form-item>
        <el-form-item label="交易类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型" style="width: 100%">
            <el-option label="充值" value="recharge" />
            <el-option label="消费" value="payment" />
            <el-option label="退款" value="refund" />
            <el-option label="转入" value="transfer_in" />
            <el-option label="转出" value="transfer_out" />
          </el-select>
        </el-form-item>
        <el-form-item label="金额" prop="amount">
          <el-input-number v-model="formData.amount" :min="0.01" :precision="2" :step="100" style="width: 100%" />
        </el-form-item>
        <el-form-item label="支付方式">
          <el-select v-model="formData.gateway" placeholder="选择支付方式" clearable style="width: 100%">
            <el-option label="支付宝" value="alipay" />
            <el-option label="微信" value="wechat" />
            <el-option label="余额" value="balance" />
            <el-option label="银行卡" value="bank" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="交易详情" width="600px">
      <el-descriptions :column="2" border v-if="detailData">
        <el-descriptions-item label="ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item label="交易号">{{ detailData.transaction_no }}</el-descriptions-item>
        <el-descriptions-item label="用户ID">{{ detailData.user_id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ detailData.username }}</el-descriptions-item>
        <el-descriptions-item label="交易类型">
          <el-tag :type="typeTagMap[detailData.type]" size="small">{{ typeTextMap[detailData.type] || detailData.type }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="金额">¥{{ Number(detailData.amount).toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="支付方式">{{ gatewayTextMap[detailData.gateway] || detailData.gateway || '-' }}</el-descriptions-item>
        <el-descriptions-item label="余额">¥{{ Number(detailData.balance_after || 0).toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="关联发票">{{ detailData.invoice_id || '-' }}</el-descriptions-item>
        <el-descriptions-item label="交易时间">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ detailData.remark || '-' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search, Plus, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'AccountsManage' })

const typeTextMap: Record<string, string> = {
  recharge: '充值', payment: '消费', refund: '退款',
  withdraw: '提现', transfer_in: '转入', transfer_out: '转出'
}
const typeTagMap: Record<string, any> = {
  recharge: 'success', payment: 'warning', refund: 'info',
  withdraw: 'danger', transfer_in: 'success', transfer_out: 'danger'
}
const gatewayTextMap: Record<string, string> = { alipay: '支付宝', wechat: '微信', balance: '余额', bank: '银行卡' }

const isIncome = (type: string) => ['recharge', 'refund', 'transfer_in'].includes(type)

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const detailVisible = ref(false)
const detailData = ref<any>(null)
const formRef = ref<FormInstance>()

const searchForm = reactive({
  user_id: '', username: '', type: '', gateway: '', date_range: [] as string[]
})

const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const formData = reactive({
  user_id: '', type: 'recharge', amount: 0, gateway: '', remark: ''
})

const formRules: FormRules = {
  user_id: [{ required: true, message: '请输入用户ID', trigger: 'blur' }],
  type: [{ required: true, message: '请选择交易类型', trigger: 'change' }],
  amount: [{ required: true, message: '请输入金额', trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page, page_size: pagination.page_size,
      user_id: searchForm.user_id || undefined, username: searchForm.username || undefined,
      type: searchForm.type || undefined, gateway: searchForm.gateway || undefined
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]; params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({ url: '/api/admin/accounts', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取交易流水失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { Object.assign(searchForm, { user_id: '', username: '', type: '', gateway: '', date_range: [] }); handleSearch() }

const handleAdd = () => {
  Object.assign(formData, { user_id: '', type: 'recharge', amount: 0, gateway: '', remark: '' })
  dialogVisible.value = true
}

const handleDetail = (row: any) => { detailData.value = row; detailVisible.value = true }
const handleViewInvoice = (row: any) => { ElMessage.info(`发票ID: ${row.invoice_id}`) }

const handleExport = async () => {
  try {
    await request.get({ url: '/api/admin/accounts/export', params: { ...searchForm } })
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      await request.post({ url: '/api/admin/accounts', params: { ...formData } })
      ElMessage.success('创建成功')
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error('创建失败')
    } finally {
      submitLoading.value = false
    }
  })
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.accounts-page {
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
  .el-form-item { margin-bottom: 0; }
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.text-green { color: #67c23a; font-weight: 600; }
.text-red { color: #f56c6c; font-weight: 600; }
</style>
