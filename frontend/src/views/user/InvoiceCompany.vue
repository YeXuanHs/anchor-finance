<template>
  <div class="invoice-company-page">
    <div class="page-header">
      <h1 class="page-title">开票公司信息</h1>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>新增公司
      </el-button>
    </div>

    <!-- 搜索筛选 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="公司名称">
          <el-input v-model="filterForm.companyName" placeholder="搜索公司名称" clearable style="width: 200px;" />
        </el-form-item>
        <el-form-item label="发票类型">
          <el-select v-model="filterForm.invoiceType" placeholder="全部类型" clearable style="width: 160px;">
            <el-option label="全部" value="" />
            <el-option label="增值税普通发票" value="normal" />
            <el-option label="增值税专用发票" value="special" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleResetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table :data="filteredCompanies" stripe style="width: 100%" v-loading="loading" empty-text="暂无公司信息">
        <el-table-column prop="companyName" label="公司名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="taxNo" label="纳税人识别号" width="180">
          <template #default="{ row }">
            <span class="mono-text">{{ row.taxNo }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="invoiceType" label="发票类型" width="160">
          <template #default="{ row }">
            <el-tag :type="row.invoiceType === 'special' ? 'warning' : 'info'" size="small" effect="light">
              {{ row.invoiceType === 'special' ? '增值税专用发票' : '增值税普通发票' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="isDefault" label="默认" width="80" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.isDefault" type="primary" size="small" effect="dark">默认</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="updateTime" label="更新时间" width="120" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click="handleView(row)">查看</el-button>
            <el-button type="primary" size="small" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          background
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="公司信息详情" width="600px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="公司名称">{{ currentCompany?.companyName }}</el-descriptions-item>
        <el-descriptions-item label="纳税人识别号">{{ currentCompany?.taxNo }}</el-descriptions-item>
        <el-descriptions-item label="发票类型">{{ currentCompany?.invoiceType === 'special' ? '增值税专用发票' : '增值税普通发票' }}</el-descriptions-item>
        <el-descriptions-item label="开户银行">{{ currentCompany?.bankName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="银行账号">{{ currentCompany?.bankAccount || '-' }}</el-descriptions-item>
        <el-descriptions-item label="注册地址">{{ currentCompany?.registerAddress || '-' }}</el-descriptions-item>
        <el-descriptions-item label="注册电话">{{ currentCompany?.registerPhone || '-' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="formVisible"
      :title="editingCompany ? '编辑公司信息' : '新增公司信息'"
      width="600px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="公司名称" prop="companyName">
          <el-input v-model="form.companyName" placeholder="请输入公司全称" />
        </el-form-item>
        <el-form-item label="纳税人识别号" prop="taxNo">
          <el-input v-model="form.taxNo" placeholder="请输入纳税人识别号" maxlength="20" />
        </el-form-item>
        <el-form-item label="发票类型" prop="invoiceType">
          <el-radio-group v-model="form.invoiceType">
            <el-radio value="normal">增值税普通发票</el-radio>
            <el-radio value="special">增值税专用发票</el-radio>
          </el-radio-group>
        </el-form-item>
        <template v-if="form.invoiceType === 'special'">
          <el-form-item label="开户银行" prop="bankName">
            <el-input v-model="form.bankName" placeholder="请输入开户银行" />
          </el-form-item>
          <el-form-item label="银行账号" prop="bankAccount">
            <el-input v-model="form.bankAccount" placeholder="请输入银行账号" />
          </el-form-item>
          <el-form-item label="注册地址" prop="registerAddress">
            <el-input v-model="form.registerAddress" placeholder="请输入注册地址" />
          </el-form-item>
          <el-form-item label="注册电话" prop="registerPhone">
            <el-input v-model="form.registerPhone" placeholder="请输入注册电话" />
          </el-form-item>
        </template>
        <el-form-item label="设为默认">
          <el-switch v-model="form.isDefault" />
          <span class="form-tip" style="margin-left: 8px;">设置后申请发票时自动填充</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/request'

interface CompanyInfo {
  id: number
  companyName: string
  taxNo: string
  invoiceType: 'normal' | 'special'
  bankName?: string
  bankAccount?: string
  registerAddress?: string
  registerPhone?: string
  isDefault: boolean
  updateTime: string
}

const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const detailVisible = ref(false)
const formVisible = ref(false)
const editingCompany = ref<CompanyInfo | null>(null)
const currentCompany = ref<CompanyInfo | null>(null)
const formRef = ref<FormInstance>()

const filterForm = reactive({
  companyName: '',
  invoiceType: ''
})

const companies = ref<CompanyInfo[]>([])

const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v2/contacts/default')
    if (data?.data) {
      companies.value = data.data.list || data.data || []
      total.value = companies.value.length
    }
  } catch (e) {
    console.error('Failed to fetch company info:', e)
  } finally {
    loading.value = false
  }
})

const form = reactive({
  companyName: '',
  taxNo: '',
  invoiceType: 'normal' as 'normal' | 'special',
  bankName: '',
  bankAccount: '',
  registerAddress: '',
  registerPhone: '',
  isDefault: false
})

const rules: FormRules = {
  companyName: [{ required: true, message: '请输入公司名称', trigger: 'blur' }],
  taxNo: [{ required: true, message: '请输入纳税人识别号', trigger: 'blur' }],
  invoiceType: [{ required: true, message: '请选择发票类型', trigger: 'change' }],
  bankName: [{ required: true, message: '请输入开户银行', trigger: 'blur' }],
  bankAccount: [{ required: true, message: '请输入银行账号', trigger: 'blur' }],
  registerAddress: [{ required: true, message: '请输入注册地址', trigger: 'blur' }],
  registerPhone: [{ required: true, message: '请输入注册电话', trigger: 'blur' }]
}

const filteredCompanies = computed(() => {
  let result = companies.value
  if (filterForm.companyName) {
    result = result.filter(c => c.companyName.includes(filterForm.companyName))
  }
  if (filterForm.invoiceType) {
    result = result.filter(c => c.invoiceType === filterForm.invoiceType)
  }
  total.value = result.length
  return result
})

function handleSearch() {
  currentPage.value = 1
  ElMessage.success('查询完成')
}

function handleResetFilter() {
  filterForm.companyName = ''
  filterForm.invoiceType = ''
  currentPage.value = 1
}

function handleAdd() {
  editingCompany.value = null
  form.companyName = ''
  form.taxNo = ''
  form.invoiceType = 'normal'
  form.bankName = ''
  form.bankAccount = ''
  form.registerAddress = ''
  form.registerPhone = ''
  form.isDefault = false
  formVisible.value = true
}

function handleView(row: CompanyInfo) {
  currentCompany.value = row
  detailVisible.value = true
}

function handleEdit(row: CompanyInfo) {
  editingCompany.value = row
  form.companyName = row.companyName
  form.taxNo = row.taxNo
  form.invoiceType = row.invoiceType
  form.bankName = row.bankName || ''
  form.bankAccount = row.bankAccount || ''
  form.registerAddress = row.registerAddress || ''
  form.registerPhone = row.registerPhone || ''
  form.isDefault = row.isDefault
  formVisible.value = true
}

function handleDelete(row: CompanyInfo) {
  ElMessageBox.confirm(`确定删除公司"${row.companyName}"的信息吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    companies.value = companies.value.filter(c => c.id !== row.id)
    total.value = companies.value.length
    ElMessage.success('已删除')
  }).catch(() => {})
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate((valid) => {
    if (valid) {
      if (form.isDefault) {
        companies.value.forEach(c => c.isDefault = false)
      }

      const now = new Date()
      const dateStr = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`

      if (editingCompany.value) {
        Object.assign(editingCompany.value, {
          companyName: form.companyName,
          taxNo: form.taxNo,
          invoiceType: form.invoiceType,
          bankName: form.invoiceType === 'special' ? form.bankName : '',
          bankAccount: form.invoiceType === 'special' ? form.bankAccount : '',
          registerAddress: form.invoiceType === 'special' ? form.registerAddress : '',
          registerPhone: form.invoiceType === 'special' ? form.registerPhone : '',
          isDefault: form.isDefault,
          updateTime: dateStr
        })
        ElMessage.success('公司信息已更新')
      } else {
        companies.value.push({
          id: Date.now(),
          companyName: form.companyName,
          taxNo: form.taxNo,
          invoiceType: form.invoiceType,
          bankName: form.invoiceType === 'special' ? form.bankName : '',
          bankAccount: form.invoiceType === 'special' ? form.bankAccount : '',
          registerAddress: form.invoiceType === 'special' ? form.registerAddress : '',
          registerPhone: form.invoiceType === 'special' ? form.registerPhone : '',
          isDefault: form.isDefault,
          updateTime: dateStr
        })
        total.value = companies.value.length
        ElMessage.success('公司信息已添加')
      }
      formVisible.value = false
    }
  })
}
</script>

<style scoped>
.invoice-company-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.filter-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.filter-card :deep(.el-card__body) {
  padding: 16px 20px 0;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
  gap: 0;
}

.table-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.table-card :deep(.el-card__body) { padding: 0; }
.table-card :deep(.el-table th.el-table__cell) { background: #fafbfc; font-weight: 600; }
.mono-text { font-family: 'Monaco', 'Menlo', monospace; font-size: 13px; color: #606266; }

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px 20px;
  border-top: 1px solid #e8ecf1;
}

.form-tip {
  font-size: 12px;
  color: #909399;
}

@media (max-width: 768px) {
  .page-header { flex-direction: column; gap: 12px; align-items: flex-start; }
}
</style>
