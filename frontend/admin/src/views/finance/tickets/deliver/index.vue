<template>
  <div class="ticket-deliver-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('ticketDeliver.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('ticketDeliver.addRule') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('common.search')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('ticketDeliver.ruleName')" clearable />
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('common.enable')" :value="1" />
            <el-option :label="$t('common.disable')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="$t('ticketDeliver.ruleName')" min-width="150" />
        <el-table-column prop="department_name" :label="$t('ticketDeliver.targetDept')" width="140" />
        <el-table-column prop="priority" :label="$t('ticketDeliver.priority')" width="80" align="center" />
        <el-table-column prop="conditions" :label="$t('ticketDeliver.matchConditions')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="is_enabled" :label="$t('common.status')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_enabled ? 'success' : 'info'" size="small">
              {{ row.is_enabled ? $t('common.enable') : $t('common.disable') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180" />
        <el-table-column :label="$t('common.action')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-popconfirm :title="$t('ticketDeliver.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('common.delete') }}</el-button>
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('ticketDeliver.ruleName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('ticketDeliver.enterRuleName')" />
        </el-form-item>
        <el-form-item :label="$t('ticketDeliver.targetDept')" prop="department_id">
          <el-select v-model="formData.department_id" :placeholder="$t('ticketDeliver.selectDept')">
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ticketDeliver.priority')" prop="priority">
          <el-input-number v-model="formData.priority" :min="0" :max="999" />
          <span class="form-tip">{{ $t('ticketDeliver.priorityTip') }}</span>
        </el-form-item>
        <el-form-item :label="$t('ticketDeliver.matchKeyword')">
          <el-input v-model="formData.match_keyword" :placeholder="$t('ticketDeliver.matchKeywordPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('ticketDeliver.matchDept')">
          <el-select v-model="formData.match_department_id" :placeholder="$t('ticketDeliver.matchDeptPlaceholder')" clearable>
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ticketDeliver.matchPriority')">
          <el-select v-model="formData.match_priority" :placeholder="$t('ticketDeliver.matchPriorityPlaceholder')" clearable>
            <el-option :label="$t('ticketDeliver.priorityLow')" value="low" />
            <el-option :label="$t('ticketDeliver.priorityMedium')" value="medium" />
            <el-option :label="$t('ticketDeliver.priorityHigh')" value="high" />
            <el-option :label="$t('ticketDeliver.priorityUrgent')" value="urgent" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-switch v-model="formData.is_enabled" :active-value="true" :inactive-value="false" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

interface Department {
  id: number
  name: string
}

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref($t('ticketDeliver.addRule'))
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: '',
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])
const departments = ref<Department[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  department_id: undefined as number | undefined,
  priority: 0,
  match_keyword: '',
  match_department_id: undefined as number | undefined,
  match_priority: '',
  is_enabled: true
})

const formRules: FormRules = {
  name: [
    { required: true, message: () => $t('ticketDeliver.enterRuleName'), trigger: 'blur' },
    { min: 2, max: 50, message: () => $t('ticketDeliver.ruleNameLength'), trigger: 'blur' }
  ],
  department_id: [
    { required: true, message: () => $t('ticketDeliver.selectDept'), trigger: 'change' }
  ]
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/ticket-deliver/rules',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error($t('ticketDeliver.fetchFailed'), error)
    ElMessage.error($t('ticketDeliver.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const fetchDepartments = async () => {
  try {
    const data = await request.get({ url: '/api/admin/ticket-depts' })
    departments.value = data || []
  } catch (error) {
    console.error($t('ticketDeliver.fetchDeptFailed'), error)
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = undefined
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = $t('ticketDeliver.addRule')
  formData.id = undefined
  formData.name = ''
  formData.department_id = undefined
  formData.priority = 0
  formData.match_keyword = ''
  formData.match_department_id = undefined
  formData.match_priority = ''
  formData.is_enabled = true
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('ticketDeliver.editRule')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/ticket-deliver/rules/${row.id}` })
    ElMessage.success($t('common.deleteSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('common.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({
          url: `/api/admin/ticket-deliver/rules/${formData.id}`,
          params: { ...formData }
        })
      } else {
        await request.post({
          url: '/api/admin/ticket-deliver/rules',
          params: { ...formData }
        })
      }
      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('common.operateFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchData()
}

const handlePageChange = () => {
  fetchData()
}

onMounted(() => {
  fetchData()
  fetchDepartments()
})
</script>

<style scoped lang="scss">
.ticket-deliver-page {
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

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.form-tip {
  margin-left: 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
</style>
