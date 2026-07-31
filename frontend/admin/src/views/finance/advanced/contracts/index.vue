<template>
  <div class="contract-list-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>合同管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加合同
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="合同编号">
          <el-input v-model="searchForm.contract_no" placeholder="请输入合同编号" clearable />
        </el-form-item>
        <el-form-item label="客户名称">
          <el-input v-model="searchForm.client_name" placeholder="请输入客户名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="草稿" :value="0" />
            <el-option label="生效中" :value="1" />
            <el-option label="已到期" :value="2" />
            <el-option label="已终止" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
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
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="contract_no" label="合同编号" width="160">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">
              {{ row.contract_no }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="合同标题" min-width="180" show-overflow-tooltip />
        <el-table-column prop="client_name" label="关联客户" width="120" />
        <el-table-column prop="amount" label="合同金额" width="120" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getContractStatusType(row.status)" size="small">
              {{ getContractStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="start_date" label="开始日期" width="120" />
        <el-table-column prop="end_date" label="结束日期" width="120" />
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="250" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="success" link @click="handleGeneratePDF(row)">生成PDF</el-button>
            <el-popconfirm title="确定删除该合同吗？" @confirm="handleDelete(row)">
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="合同编号" prop="contract_no">
              <el-input v-model="formData.contract_no" placeholder="请输入合同编号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="合同标题" prop="title">
              <el-input v-model="formData.title" placeholder="请输入合同标题" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="关联客户" prop="client_id">
              <el-select v-model="formData.client_id" placeholder="请选择客户" filterable style="width: 100%">
                <el-option v-for="client in clientList" :key="client.id" :label="client.username" :value="client.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="合同金额" prop="amount">
              <el-input-number v-model="formData.amount" :min="0" :precision="2" :step="100" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="开始日期" prop="start_date">
              <el-date-picker v-model="formData.start_date" type="date" value-format="YYYY-MM-DD" placeholder="选择开始日期" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束日期" prop="end_date">
              <el-date-picker v-model="formData.end_date" type="date" value-format="YYYY-MM-DD" placeholder="选择结束日期" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="状态" prop="status">
          <el-select v-model="formData.status" placeholder="请选择状态">
            <el-option label="草稿" :value="0" />
            <el-option label="生效中" :value="1" />
            <el-option label="已到期" :value="2" />
            <el-option label="已终止" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="合同内容" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="6" placeholder="请输入合同内容" />
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

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="合同详情" width="750px" destroy-on-close>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="合同编号">{{ detailData.contract_no }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getContractStatusType(detailData.status)" size="small">
            {{ getContractStatusText(detailData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="合同标题" :span="2">{{ detailData.title }}</el-descriptions-item>
        <el-descriptions-item label="关联客户">{{ detailData.client_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="合同金额">
          <span class="amount-text">¥{{ formatAmount(detailData.amount) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="开始日期">{{ detailData.start_date || '-' }}</el-descriptions-item>
        <el-descriptions-item label="结束日期">{{ detailData.end_date || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="合同内容" :span="2">
          <div class="contract-content">{{ detailData.content || '无' }}</div>
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ detailData.remark || '无' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
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
  contract_no: '',
  client_name: '',
  status: undefined as number | undefined,
  date_range: [] as string[]
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref<any[]>([])

// 客户列表（用于下拉选择）
const clientList = ref<any[]>([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('添加合同')
const formRef = ref<FormInstance>()

// 详情对话框
const detailVisible = ref(false)
const detailData = ref<any>({})

// 表单数据
const formData = reactive({
  id: undefined as number | undefined,
  contract_no: '',
  title: '',
  client_id: undefined as number | undefined,
  amount: 0,
  start_date: '',
  end_date: '',
  status: 0,
  content: '',
  remark: ''
})

// 表单验证规则
const formRules: FormRules = {
  contract_no: [
    { required: true, message: '请输入合同编号', trigger: 'blur' }
  ],
  title: [
    { required: true, message: '请输入合同标题', trigger: 'blur' }
  ],
  client_id: [
    { required: true, message: '请选择关联客户', trigger: 'change' }
  ],
  start_date: [
    { required: true, message: '请选择开始日期', trigger: 'change' }
  ],
  end_date: [
    { required: true, message: '请选择结束日期', trigger: 'change' }
  ]
}

// 合同状态映射
const CONTRACT_STATUS_MAP: Record<number, { text: string; type: string }> = {
  0: { text: '草稿', type: 'info' },
  1: { text: '生效中', type: 'success' },
  2: { text: '已到期', type: 'warning' },
  3: { text: '已终止', type: 'danger' }
}

// 获取合同状态文本
const getContractStatusText = (status: number) => {
  return CONTRACT_STATUS_MAP[status]?.text || '未知'
}

// 获取合同状态类型
const getContractStatusType = (status: number) => {
  return (CONTRACT_STATUS_MAP[status]?.type || 'info') as any
}

// 格式化金额
const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

// 获取合同列表
const fetchContracts = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      contract_no: searchForm.contract_no || undefined,
      client_name: searchForm.client_name || undefined,
      status: searchForm.status
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({
      url: '/api/admin/contracts',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取合同列表失败:', error)
    ElMessage.error('获取合同列表失败')
  } finally {
    loading.value = false
  }
}

// 获取客户列表（用于下拉选择）
const fetchClients = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/users',
      params: { page: 1, page_size: 9999 }
    })
    clientList.value = data.list || []
  } catch (error) {
    console.error('获取客户列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchContracts()
}

// 重置
const handleReset = () => {
  searchForm.contract_no = ''
  searchForm.client_name = ''
  searchForm.status = undefined
  searchForm.date_range = []
  handleSearch()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = '添加合同'
  formData.id = undefined
  formData.contract_no = ''
  formData.title = ''
  formData.client_id = undefined
  formData.amount = 0
  formData.start_date = ''
  formData.end_date = ''
  formData.status = 0
  formData.content = ''
  formData.remark = ''
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = '编辑合同'
  Object.assign(formData, row)
  dialogVisible.value = true
}

// 查看详情
const handleViewDetail = (row: any) => {
  detailData.value = { ...row }
  detailVisible.value = true
}

// 生成PDF
const handleGeneratePDF = async (row: any) => {
  try {
    const data = await request.get({
      url: `/api/admin/contracts/${row.id}/pdf`
    })
    if (data.url) {
      window.open(data.url, '_blank')
      ElMessage.success('PDF生成成功')
    }
  } catch (error) {
    ElMessage.error('PDF生成失败')
  }
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/contracts/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchContracts()
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
      const url = formData.id ? `/api/admin/contracts/${formData.id}` : '/api/admin/contracts'

      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchContracts()
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
  fetchContracts()
}

// 页码变化
const handlePageChange = () => {
  fetchContracts()
}

onMounted(() => {
  fetchContracts()
  fetchClients()
})
</script>

<style scoped lang="scss">
.contract-list-page {
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

.contract-content {
  max-height: 200px;
  overflow-y: auto;
  white-space: pre-wrap;
  line-height: 1.6;
}
</style>
