<template>
  <div class="credit-page art-full-height">
    <!-- 搜索栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.user_id" placeholder="请输入用户ID" clearable />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="searchForm.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
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
          <span>信用额度列表</span>
        </div>
      </template>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="credit_limit" label="信用额度" width="150">
          <template #default="{ row }">
            <span class="text-primary">¥{{ row.credit_limit?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="used_amount" label="已用额度" width="150">
          <template #default="{ row }">
            <span class="text-warning">¥{{ row.used_amount?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="available_amount" label="可用额度" width="150">
          <template #default="{ row }">
            <span class="text-success">¥{{ row.available_amount?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="使用率" width="150">
          <template #default="{ row }">
            <el-progress
              :percentage="getUsagePercent(row)"
              :color="getUsageColor(row)"
              :stroke-width="10"
            />
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="expire_at" label="到期时间" width="180" />
        <el-table-column prop="updated_at" label="更新时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleSetCredit(row)">设置额度</el-button>
            <el-button type="warning" link @click="handleAdjustCredit(row)">调整额度</el-button>
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

    <!-- 设置额度对话框 -->
    <el-dialog v-model="setDialogVisible" title="设置信用额度" width="500px">
      <el-form :model="setForm" :rules="setRules" ref="setFormRef" label-width="100px">
        <el-form-item label="用户">
          <el-input :value="setForm.username" disabled />
        </el-form-item>
        <el-form-item label="当前额度">
          <el-input :value="'¥' + (setForm.current_limit?.toFixed(2) || '0.00')" disabled />
        </el-form-item>
        <el-form-item label="新额度" prop="credit_limit">
          <el-input-number v-model="setForm.credit_limit" :min="0" :precision="2" :step="1000" style="width: 100%" />
        </el-form-item>
        <el-form-item label="到期时间" prop="expire_at">
          <el-date-picker
            v-model="setForm.expire_at"
            type="datetime"
            placeholder="选择到期时间"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="setDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitSet" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 调整额度对话框 -->
    <el-dialog v-model="adjustDialogVisible" title="调整信用额度" width="500px">
      <el-form :model="adjustForm" :rules="adjustRules" ref="adjustFormRef" label-width="100px">
        <el-form-item label="用户">
          <el-input :value="adjustForm.username" disabled />
        </el-form-item>
        <el-form-item label="当前额度">
          <el-input :value="'¥' + (adjustForm.current_limit?.toFixed(2) || '0.00')" disabled />
        </el-form-item>
        <el-form-item label="调整方式">
          <el-radio-group v-model="adjustForm.adjust_type">
            <el-radio label="increase">增加</el-radio>
            <el-radio label="decrease">减少</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="调整金额" prop="amount">
          <el-input-number v-model="adjustForm.amount" :min="0.01" :precision="2" :step="100" style="width: 100%" />
        </el-form-item>
        <el-form-item label="调整原因" prop="reason">
          <el-input v-model="adjustForm.reason" type="textarea" placeholder="请输入调整原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="adjustDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitAdjust" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'CreditManage' })

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  user_id: '',
  username: '',
  status: undefined as number | undefined
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref([])

// 设置额度对话框
const setDialogVisible = ref(false)
const setFormRef = ref<FormInstance>()
const setForm = reactive({
  id: 0,
  username: '',
  current_limit: 0,
  credit_limit: 0,
  expire_at: ''
})

// 调整额度对话框
const adjustDialogVisible = ref(false)
const adjustFormRef = ref<FormInstance>()
const adjustForm = reactive({
  id: 0,
  username: '',
  current_limit: 0,
  adjust_type: 'increase',
  amount: 0,
  reason: ''
})

// 设置表单验证规则
const setRules: FormRules = {
  credit_limit: [
    { required: true, message: '请输入信用额度', trigger: 'blur' }
  ]
}

// 调整表单验证规则
const adjustRules: FormRules = {
  amount: [
    { required: true, message: '请输入调整金额', trigger: 'blur' }
  ]
}

// 获取使用率
const getUsagePercent = (row: any) => {
  if (!row.credit_limit || row.credit_limit === 0) return 0
  return Math.min(100, Math.round((row.used_amount / row.credit_limit) * 100))
}

// 获取使用率颜色
const getUsageColor = (row: any) => {
  const percent = getUsagePercent(row)
  if (percent >= 90) return '#f56c6c'
  if (percent >= 70) return '#e6a23c'
  return '#67c23a'
}

// 获取信用额度列表
const fetchCredits = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size
    }
    if (searchForm.user_id) params.user_id = searchForm.user_id
    if (searchForm.username) params.username = searchForm.username
    if (searchForm.status !== undefined) params.status = searchForm.status

    const data = await request.get({
      url: '/api/admin/credits',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取信用额度列表失败:', error)
    ElMessage.error('获取信用额度列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchCredits()
}

// 重置
const handleReset = () => {
  searchForm.user_id = ''
  searchForm.username = ''
  searchForm.status = undefined
  handleSearch()
}

// 设置额度
const handleSetCredit = (row: any) => {
  setForm.id = row.id
  setForm.username = row.username
  setForm.current_limit = row.credit_limit
  setForm.credit_limit = row.credit_limit
  setForm.expire_at = row.expire_at || ''
  setDialogVisible.value = true
}

// 调整额度
const handleAdjustCredit = (row: any) => {
  adjustForm.id = row.id
  adjustForm.username = row.username
  adjustForm.current_limit = row.credit_limit
  adjustForm.adjust_type = 'increase'
  adjustForm.amount = 0
  adjustForm.reason = ''
  adjustDialogVisible.value = true
}

// 提交设置额度
const handleSubmitSet = async () => {
  if (!setFormRef.value) return

  await setFormRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      await request.put({
        url: `/api/admin/credits/${setForm.id}`,
        params: {
          credit_limit: setForm.credit_limit,
          expire_at: setForm.expire_at
        },
        showSuccessMessage: true
      })
      setDialogVisible.value = false
      fetchCredits()
    } catch (error) {
      console.error('设置额度失败:', error)
    } finally {
      submitLoading.value = false
    }
  })
}

// 提交调整额度
const handleSubmitAdjust = async () => {
  if (!adjustFormRef.value) return

  await adjustFormRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const newLimit = adjustForm.adjust_type === 'increase'
        ? adjustForm.current_limit + adjustForm.amount
        : Math.max(0, adjustForm.current_limit - adjustForm.amount)

      await request.put({
        url: `/api/admin/credits/${adjustForm.id}`,
        params: {
          credit_limit: newLimit,
          reason: adjustForm.reason
        },
        showSuccessMessage: true
      })
      adjustDialogVisible.value = false
      fetchCredits()
    } catch (error) {
      console.error('调整额度失败:', error)
    } finally {
      submitLoading.value = false
    }
  })
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchCredits()
}

// 页码变化
const handlePageChange = () => {
  fetchCredits()
}

onMounted(() => {
  fetchCredits()
})
</script>

<style scoped lang="scss">
.credit-page {
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

.text-warning {
  color: #e6a23c;
  font-weight: 600;
}

.text-success {
  color: #67c23a;
  font-weight: 600;
}
</style>
