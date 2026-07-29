<template>
  <div class="promotions-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="活动名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="进行中" value="active" />
            <el-option label="未开始" value="pending" />
            <el-option label="已结束" value="ended" />
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
        <h3>促销活动管理</h3>
        <el-button type="primary" @click="openAddDialog">
          <el-icon><Plus /></el-icon>
          新增活动
        </el-button>
      </div>

      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" min-width="160" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="discount" label="折扣" width="120">
          <template #default="{ row }">
            <span v-if="row.type === 'percent'">{{ row.discount }}%</span>
            <span v-else>¥{{ row.discount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="applicable_products" label="适用产品" min-width="140" show-overflow-tooltip />
        <el-table-column prop="start_time" label="开始时间" width="180" />
        <el-table-column prop="end_time" label="结束时间" width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTypeMap[row.status]" size="small">
              {{ statusTextMap[row.status] || row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑活动' : '新增活动'" width="650px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="110px">
        <el-form-item label="活动名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入活动名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型">
            <el-option label="固定金额" value="fixed" />
            <el-option label="百分比折扣" value="percent" />
          </el-select>
        </el-form-item>
        <el-form-item label="折扣值" prop="discount">
          <el-input-number v-model="formData.discount" :min="0" :precision="2" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">
            {{ formData.type === 'percent' ? '%' : '元' }}
          </span>
        </el-form-item>
        <el-form-item label="最低消费" prop="min_amount">
          <el-input-number v-model="formData.min_amount" :min="0" :precision="2" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">元（0 为无门槛）</span>
        </el-form-item>
        <el-form-item label="适用产品" prop="applicable_products">
          <el-input v-model="formData.applicable_products" placeholder="产品ID，多个用逗号分隔" />
        </el-form-item>
        <el-form-item label="活动时间" prop="time_range">
          <el-date-picker
            v-model="formData.time_range"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入活动描述" />
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

const typeOptions = [
  { label: '固定金额', value: 'fixed' },
  { label: '百分比折扣', value: 'percent' }
]

const statusTypeMap: Record<string, string> = { active: 'success', pending: 'warning', ended: 'info', disabled: 'info' }
const statusTextMap: Record<string, string> = { active: '进行中', pending: '未开始', ended: '已结束', disabled: '禁用' }

const searchForm = ref({ keyword: '', status: '' })

const formData = reactive({
  name: '',
  type: 'fixed' as string,
  discount: 0,
  min_amount: 0,
  applicable_products: '',
  time_range: [] as string[],
  description: '',
  status: 'active' as string
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入活动名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  discount: [{ required: true, message: '请输入折扣值', trigger: 'blur' }],
  time_range: [{ required: true, message: '请选择活动时间', trigger: 'change' }]
}

const getTypeLabel = (val: string) => typeOptions.find(o => o.value === val)?.label || val

const handleSearch = () => { page.value = 1; fetchData() }

const resetSearch = () => {
  searchForm.value = { keyword: '', status: '' }
  handleSearch()
}

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/promotions', {
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
    name: '', type: 'fixed', discount: 0, min_amount: 0,
    applicable_products: '', time_range: [], description: '', status: 'active'
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
    name: row.name, type: row.type, discount: row.discount,
    min_amount: row.min_amount, applicable_products: row.applicable_products,
    time_range: [row.start_time, row.end_time], description: row.description, status: row.status
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  const payload = {
    ...formData,
    start_time: formData.time_range?.[0] || '',
    end_time: formData.time_range?.[1] || ''
  }
  try {
    if (isEdit.value) {
      await request.put(`/admin/api/v1/promotions/${editId.value}`, payload)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/promotions', payload)
      ElMessage.success('新增成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch {} finally {
    submitLoading.value = false
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm(`确定删除活动「${row.name}」吗？`, '提示', { type: 'warning' })
  try {
    await request.delete(`/admin/api/v1/promotions/${row.id}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch {}
}

onMounted(fetchData)
</script>

<style scoped lang="scss">
.promotions-page {
  .search-bar { margin-bottom: 16px; }
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 18px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
}
</style>
