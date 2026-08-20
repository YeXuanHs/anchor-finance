<template>
  <div class="ticket-prereply-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('ticketPrereply.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('ticketPrereply.addPrereply') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('ticketPrereply.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('ticketPrereply.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('ticketList.department')">
          <el-select v-model="searchForm.department_id" :placeholder="$t('common.all')" clearable>
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" :label="$t('ticketPrereply.columnTitle')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="department_name" :label="$t('ticketPrereply.department')" width="120" />
        <el-table-column prop="content" :label="$t('ticketPrereply.replyContent')" min-width="300" show-overflow-tooltip />
        <el-table-column prop="sort_order" :label="$t('common.sort')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('common.status')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('common.enabled') : $t('common.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="use_count" :label="$t('ticketPrereply.useCount')" width="100" align="center" />
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180" />
        <el-table-column :label="$t('common.action')" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-popconfirm :title="$t('ticketPrereply.confirmDelete')" @confirm="handleDelete(row)">
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
        <el-form-item :label="$t('ticketPrereply.columnTitle')" prop="title">
          <el-input v-model="formData.title" :placeholder="$t('ticketPrereply.enterTitle')" />
        </el-form-item>
        <el-form-item :label="$t('ticketPrereply.department')" prop="department_id">
          <el-select v-model="formData.department_id" :placeholder="$t('ticketPrereply.selectDepartment')">
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ticketPrereply.replyContent')" prop="content">
          <el-input
            v-model="formData.content"
            type="textarea"
            :rows="6"
            :placeholder="$t('ticketPrereply.enterContent')"
          />
        </el-form-item>
        <el-form-item :label="$t('common.sort')" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="999" />
        </el-form-item>
        <el-form-item :label="$t('common.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
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
import request from '@/utils/http'
import { $t } from '@/locales'

interface PreReply {
  id: number
  title: string
  content: string
  department_id: number
  department_name: string
  sort_order: number
  status: number
  use_count: number
  created_at: string
}

interface Department {
  id: number
  name: string
}

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref($t('ticketPrereply.addPrereply'))
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: '',
  department_id: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<PreReply[]>([])
const departments = ref<Department[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  title: '',
  content: '',
  department_id: undefined as number | undefined,
  sort_order: 0,
  status: 1
})

const formRules: FormRules = {
  title: [
    { required: true, message: () => $t('ticketPrereply.enterTitle'), trigger: 'blur' },
    { min: 2, max: 100, message: () => $t('ticketPrereply.titleLength'), trigger: 'blur' }
  ],
  department_id: [
    { required: true, message: () => $t('ticketPrereply.selectDepartment'), trigger: 'change' }
  ],
  content: [
    { required: true, message: () => $t('ticketPrereply.enterContent'), trigger: 'blur' }
  ]
}

const fetchPreReplies = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/ticket-prereply',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取预回复列表失败:', error)
    ElMessage.error($t('ticketPrereply.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const fetchDepartments = async () => {
  try {
    const data = await request.get({
      url: '/api/admin/ticket-depts'
    })
    departments.value = data || []
  } catch (error) {
    console.error('获取部门列表失败:', error)
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchPreReplies()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.department_id = undefined
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = $t('ticketPrereply.addPrereply')
  formData.id = undefined
  formData.title = ''
  formData.content = ''
  formData.department_id = undefined
  formData.sort_order = 0
  formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('ticketPrereply.editPrereply')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/ticket-prereply/replies/${row.id}`
    })
    ElMessage.success($t('common.deleteSuccess'))
    fetchPreReplies()
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
          url: `/api/admin/ticket-prereply/replies/${formData.id}`,
          params: formData
        })
      } else {
        await request.post({
          url: '/api/admin/ticket-prereply/replies',
          params: formData
        })
      }

      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess'))
      dialogVisible.value = false
      fetchPreReplies()
    } catch (error) {
      ElMessage.error($t('common.operateFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchPreReplies()
}

const handlePageChange = () => {
  fetchPreReplies()
}

onMounted(() => {
  fetchPreReplies()
  fetchDepartments()
})
</script>

<style scoped lang="scss">
.ticket-prereply-page {
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
</style>
