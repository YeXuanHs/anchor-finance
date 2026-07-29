<template>
  <div class="coupons-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="名称/券码" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="固定金额" value="fixed" />
            <el-option label="百分比折扣" value="percent" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="启用" value="active" />
            <el-option label="禁用" value="disabled" />
            <el-option label="已过期" value="expired" />
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
        <h3>优惠券管理</h3>
        <el-button type="primary" @click="openAddDialog">
          <el-icon><Plus /></el-icon>
          新增优惠券
        </el-button>
      </div>

      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="code" label="券码" width="140">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.code }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" min-width="160" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="110">
          <template #default="{ row }">
            <el-tag size="small">{{ row.type === 'fixed' ? '固定金额' : '百分比折扣' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="value" label="面额" width="120">
          <template #default="{ row }">
            <span v-if="row.type === 'fixed'" class="amount">¥{{ row.value?.toFixed(2) }}</span>
            <span v-else>{{ row.value }}%</span>
          </template>
        </el-table-column>
        <el-table-column prop="min_amount" label="最低消费" width="120">
          <template #default="{ row }">{{ row.min_amount ? `¥${row.min_amount.toFixed(2)}` : '无门槛' }}</template>
        </el-table-column>
        <el-table-column prop="used_count" label="已用" width="80" />
        <el-table-column prop="total_count" label="总量" width="80" />
        <el-table-column prop="expire_at" label="有效期" width="180" />
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑优惠券' : '新增优惠券'" width="650px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="110px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入优惠券名称" />
        </el-form-item>
        <el-form-item label="券码" prop="code">
          <el-input v-model="formData.code" placeholder="留空则自动生成">
            <template #append>
              <el-button @click="generateCode">随机生成</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型">
            <el-option label="固定金额" value="fixed" />
            <el-option label="百分比折扣" value="percent" />
          </el-select>
        </el-form-item>
        <el-form-item label="面额" prop="value">
          <el-input-number v-model="formData.value" :min="0" :precision="2" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">
            {{ formData.type === 'percent' ? '%' : '元' }}
          </span>
        </el-form-item>
        <el-form-item label="最低消费" prop="min_amount">
          <el-input-number v-model="formData.min_amount" :min="0" :precision="2" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">元（0 为无门槛）</span>
        </el-form-item>
        <el-form-item label="最高优惠" prop="max_discount" v-if="formData.type === 'percent'">
          <el-input-number v-model="formData.max_discount" :min="0" :precision="2" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">元（0 为不限制）</span>
        </el-form-item>
        <el-form-item label="总量" prop="total_count">
          <el-input-number v-model="formData.total_count" :min="1" />
        </el-form-item>
        <el-form-item label="每人限用" prop="per_user_limit">
          <el-input-number v-model="formData.per_user_limit" :min="0" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">次（0 为不限制）</span>
        </el-form-item>
        <el-form-item label="有效期" prop="expire_at">
          <el-date-picker
            v-model="formData.expire_at"
            type="datetime"
            placeholder="选择过期时间"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="适用产品" prop="applicable_products">
          <el-input v-model="formData.applicable_products" placeholder="产品ID，多个用逗号分隔" />
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

const statusTypeMap: Record<string, string> = { active: 'success', disabled: 'info', expired: 'danger' }
const statusTextMap: Record<string, string> = { active: '启用', disabled: '禁用', expired: '已过期' }

const searchForm = ref({ keyword: '', type: '', status: '' })

const formData = reactive({
  name: '',
  code: '',
  type: 'fixed' as string,
  value: 0,
  min_amount: 0,
  max_discount: 0,
  total_count: 100,
  per_user_limit: 1,
  expire_at: '',
  applicable_products: '',
  status: 'active' as string
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  value: [{ required: true, message: '请输入面额', trigger: 'blur' }],
  total_count: [{ required: true, message: '请输入总量', trigger: 'blur' }],
  expire_at: [{ required: true, message: '请选择有效期', trigger: 'change' }]
}

const generateCode = () => {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'
  formData.code = Array.from({ length: 8 }, () => chars[Math.floor(Math.random() * chars.length)]).join('')
}

const handleSearch = () => { page.value = 1; fetchData() }

const resetSearch = () => {
  searchForm.value = { keyword: '', type: '', status: '' }
  handleSearch()
}

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/coupons', {
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
    name: '', code: '', type: 'fixed', value: 0, min_amount: 0, max_discount: 0,
    total_count: 100, per_user_limit: 1, expire_at: '', applicable_products: '', status: 'active'
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
    name: row.name, code: row.code, type: row.type, value: row.value,
    min_amount: row.min_amount, max_discount: row.max_discount,
    total_count: row.total_count, per_user_limit: row.per_user_limit,
    expire_at: row.expire_at, applicable_products: row.applicable_products, status: row.status
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (isEdit.value) {
      await request.put(`/admin/api/v1/coupons/${editId.value}`, formData)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/coupons', formData)
      ElMessage.success('新增成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch {} finally {
    submitLoading.value = false
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm(`确定删除优惠券「${row.name}」吗？`, '提示', { type: 'warning' })
  try {
    await request.delete(`/admin/api/v1/coupons/${row.id}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch {}
}

onMounted(fetchData)
</script>

<style scoped lang="scss">
.coupons-page {
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
