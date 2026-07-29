<template>
  <div class="vouchers-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="批次名称/券码" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="进行中" value="active" />
            <el-option label="已结束" value="ended" />
            <el-option label="未开始" value="pending" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>代金券批次</h3>
        <el-button type="primary" @click="openAddDialog">
          <el-icon><Plus /></el-icon>
          新建批次
        </el-button>
      </div>

      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="批次名称" min-width="160" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="110">
          <template #default="{ row }">
            <el-tag size="small">{{ row.type === 'fixed' ? '固定金额' : '百分比折扣' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="面额/折扣" width="120">
          <template #default="{ row }">
            <span v-if="row.type === 'fixed'" class="amount">¥{{ row.amount?.toFixed(2) }}</span>
            <span v-else>{{ row.amount }}%</span>
          </template>
        </el-table-column>
        <el-table-column prop="total_count" label="总数量" width="100" />
        <el-table-column prop="used_count" label="已使用" width="100">
          <template #default="{ row }">
            <span :style="{ color: row.used_count >= row.total_count ? 'var(--el-color-danger)' : '' }">
              {{ row.used_count }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="remaining" label="剩余" width="80">
          <template #default="{ row }">{{ row.total_count - row.used_count }}</template>
        </el-table-column>
        <el-table-column prop="start_time" label="开始时间" width="180" />
        <el-table-column prop="end_time" label="结束时间" width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTypeMap[row.status]" size="small">
              {{ statusTextMap[row.status] || row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewVouchers(row)">券码列表</el-button>
            <el-button type="primary" link @click="openEditDialog(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchData"
          @size-change="fetchData"
        />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑批次' : '新建批次'" width="650px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="110px">
        <el-form-item label="批次名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入批次名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型">
            <el-option label="固定金额" value="fixed" />
            <el-option label="百分比折扣" value="percent" />
          </el-select>
        </el-form-item>
        <el-form-item label="面额/折扣" prop="amount">
          <el-input-number v-model="formData.amount" :min="0" :precision="2" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">
            {{ formData.type === 'percent' ? '%' : '元' }}
          </span>
        </el-form-item>
        <el-form-item label="生成数量" prop="total_count" v-if="!isEdit">
          <el-input-number v-model="formData.total_count" :min="1" :max="10000" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">最多 10000 张</span>
        </el-form-item>
        <el-form-item label="券码前缀" prop="prefix" v-if="!isEdit">
          <el-input v-model="formData.prefix" placeholder="如 VCH，留空自动生成" maxlength="6" />
        </el-form-item>
        <el-form-item label="券码长度" prop="code_length" v-if="!isEdit">
          <el-input-number v-model="formData.code_length" :min="6" :max="16" />
        </el-form-item>
        <el-form-item label="最低消费" prop="min_amount">
          <el-input-number v-model="formData.min_amount" :min="0" :precision="2" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">元（0 为无门槛）</span>
        </el-form-item>
        <el-form-item label="每人限用" prop="per_user_limit">
          <el-input-number v-model="formData.per_user_limit" :min="0" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">次（0 为不限制）</span>
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
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" active-value="active" inactive-value="disabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="voucherListVisible" title="券码列表" width="800px">
      <div style="margin-bottom: 12px; display: flex; justify-content: space-between; align-items: center;">
        <span style="color: var(--el-text-color-secondary);">
          共 {{ voucherTotal }} 张，已使用 {{ voucherUsedCount }} 张
        </span>
        <el-button size="small" @click="exportVouchers">导出券码</el-button>
      </div>
      <el-table :data="voucherList" style="width: 100%" v-loading="voucherLoading" size="small">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="code" label="券码" width="160">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.code }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="used_by" label="使用者" width="120" />
        <el-table-column prop="used_at" label="使用时间" width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'used' ? 'info' : 'success'" size="small">
              {{ row.status === 'used' ? '已使用' : '未使用' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination" style="margin-top: 12px;">
        <el-pagination
          v-model:current-page="voucherPage"
          v-model:page-size="voucherPageSize"
          :page-sizes="[20, 50, 100]"
          :total="voucherTotal"
          layout="total, sizes, prev, pager, next"
          @current-change="fetchVoucherList"
          @size-change="fetchVoucherList"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const formRef = ref<FormInstance>()
const voucherListVisible = ref(false)
const voucherLoading = ref(false)
const voucherList = ref<any[]>([])
const voucherTotal = ref(0)
const voucherUsedCount = ref(0)
const voucherPage = ref(1)
const voucherPageSize = ref(20)
const currentBatchId = ref<number | null>(null)

const statusTypeMap: Record<string, string> = { active: 'success', pending: 'warning', ended: 'info', disabled: 'info' }
const statusTextMap: Record<string, string> = { active: '进行中', pending: '未开始', ended: '已结束', disabled: '禁用' }

const searchForm = ref({ keyword: '', status: '' })

const formData = reactive({
  name: '',
  type: 'fixed' as string,
  amount: 0,
  total_count: 100,
  prefix: '',
  code_length: 10,
  min_amount: 0,
  per_user_limit: 1,
  date_range: [] as string[],
  status: 'active' as string
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入批次名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  amount: [{ required: true, message: '请输入面额', trigger: 'blur' }],
  total_count: [{ required: true, message: '请输入生成数量', trigger: 'blur' }],
  date_range: [{ required: true, message: '请选择有效期', trigger: 'change' }]
}

const handleSearch = () => { page.value = 1; fetchData() }

const resetSearch = () => {
  searchForm.value = { keyword: '', status: '' }
  handleSearch()
}

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/voucher-batches', {
      params: { page: page.value, page_size: pageSize.value, ...searchForm.value }
    })
    tableData.value = data.data || []
    total.value = data.total || 0
  } catch {} finally {
    loading.value = false
  }
}

const resetForm = () => {
  Object.assign(formData, {
    name: '', type: 'fixed', amount: 0, total_count: 100, prefix: '', code_length: 10,
    min_amount: 0, per_user_limit: 1, date_range: [], status: 'active'
  })
}

const openAddDialog = () => {
  isEdit.value = false
  editId.value = null
  resetForm()
  dialogVisible.value = true
}

const openEditDialog = (row: any) => {
  isEdit.value = true
  editId.value = row.id
  Object.assign(formData, {
    name: row.name, type: row.type, amount: row.amount,
    total_count: row.total_count, prefix: '', code_length: 10,
    min_amount: row.min_amount, per_user_limit: row.per_user_limit,
    date_range: [row.start_time, row.end_time], status: row.status
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  const payload = {
    ...formData,
    start_time: formData.date_range?.[0] || '',
    end_time: formData.date_range?.[1] || ''
  }
  try {
    if (isEdit.value) {
      await request.put(`/admin/api/v1/voucher-batches/${editId.value}`, payload)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/voucher-batches', payload)
      ElMessage.success('批次创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch {} finally {
    submitLoading.value = false
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm(`确定删除批次「${row.name}」吗？删除后所有未使用券码将失效。`, '提示', { type: 'warning' })
  try {
    await request.delete(`/admin/api/v1/voucher-batches/${row.id}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch {}
}

const viewVouchers = (row: any) => {
  currentBatchId.value = row.id
  voucherUsedCount.value = row.used_count
  voucherPage.value = 1
  voucherListVisible.value = true
  fetchVoucherList()
}

const fetchVoucherList = async () => {
  if (!currentBatchId.value) return
  voucherLoading.value = true
  try {
    const { data } = await request.get(`/admin/api/v1/voucher-batches/${currentBatchId.value}/vouchers`, {
      params: { page: voucherPage.value, page_size: voucherPageSize.value }
    })
    voucherList.value = data.data || []
    voucherTotal.value = data.total || 0
  } catch {} finally {
    voucherLoading.value = false
  }
}

const exportVouchers = () => {
  ElMessage.info('正在导出券码...')
}

onMounted(fetchData)
</script>

<style scoped lang="scss">
.vouchers-page {
  .search-bar { margin-bottom: 16px; }
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 18px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
  .amount { color: var(--el-color-danger); font-weight: 600; }
}
</style>
