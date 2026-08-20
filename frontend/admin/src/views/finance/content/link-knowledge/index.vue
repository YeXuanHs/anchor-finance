<template>
  <div class="link-knowledge-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.linkKnowledge.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('finance.linkKnowledge.addLink') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('finance.linkKnowledge.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('finance.linkKnowledge.knowledgeTitle')" clearable />
        </el-form-item>
        <el-form-item :label="$t('finance.linkKnowledge.ticketType')">
          <el-select v-model="searchForm.ticket_type" :placeholder="$t('finance.linkKnowledge.selectTicketType')" clearable>
            <el-option :label="$t('finance.linkKnowledge.all')" value="" />
            <el-option :label="$t('finance.linkKnowledge.technical')" value="technical" />
            <el-option :label="$t('finance.linkKnowledge.service')" value="service" />
            <el-option :label="$t('finance.linkKnowledge.incident')" value="incident" />
            <el-option :label="$t('finance.linkKnowledge.change')" value="change" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('finance.linkKnowledge.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('finance.linkKnowledge.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="ticket_type" :label="$t('finance.linkKnowledge.ticketType')" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ getTicketTypeLabel(row.ticket_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="knowledge_title" :label="$t('finance.linkKnowledge.knowledgeTitle')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="knowledge_id" :label="$t('finance.linkKnowledge.knowledgeId')" width="100" align="center" />
        <el-table-column prop="sort_order" :label="$t('finance.linkKnowledge.sort')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('finance.linkKnowledge.status')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('finance.linkKnowledge.enabled') : $t('finance.linkKnowledge.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('finance.linkKnowledge.createdAt')" width="180" />
        <el-table-column :label="$t('finance.linkKnowledge.actions')" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('finance.linkKnowledge.edit') }}</el-button>
            <el-popconfirm :title="$t('finance.linkKnowledge.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('finance.linkKnowledge.delete') }}</el-button>
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('finance.linkKnowledge.ticketType')" prop="ticket_type">
          <el-select v-model="formData.ticket_type" :placeholder="$t('finance.linkKnowledge.selectTicketType')">
            <el-option :label="$t('finance.linkKnowledge.technical')" value="technical" />
            <el-option :label="$t('finance.linkKnowledge.service')" value="service" />
            <el-option :label="$t('finance.linkKnowledge.incident')" value="incident" />
            <el-option :label="$t('finance.linkKnowledge.change')" value="change" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('finance.linkKnowledge.knowledgeId')" prop="knowledge_id">
          <el-input v-model="formData.knowledge_id" :placeholder="$t('finance.linkKnowledge.enterKnowledgeId')" />
        </el-form-item>
        <el-form-item :label="$t('finance.linkKnowledge.knowledgeTitle')" prop="knowledge_title">
          <el-input v-model="formData.knowledge_title" :placeholder="$t('finance.linkKnowledge.enterKnowledgeTitle')" />
        </el-form-item>
        <el-form-item :label="$t('finance.linkKnowledge.sort')" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="999" />
        </el-form-item>
        <el-form-item :label="$t('finance.linkKnowledge.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('finance.linkKnowledge.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('finance.linkKnowledge.confirm') }}</el-button>
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

interface LinkKnowledge {
  id: number
  ticket_type: string
  knowledge_id: string
  knowledge_title: string
  sort_order: number
  status: number
  created_at: string
}

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref($t('finance.linkKnowledge.addLink'))
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: '',
  ticket_type: ''
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<LinkKnowledge[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  ticket_type: '',
  knowledge_id: '',
  knowledge_title: '',
  sort_order: 0,
  status: 1
})

const formRules: FormRules = {
  ticket_type: [
    { required: true, message: $t('finance.linkKnowledge.selectTicketType'), trigger: 'change' }
  ],
  knowledge_id: [
    { required: true, message: $t('finance.linkKnowledge.enterKnowledgeId'), trigger: 'blur' }
  ],
  knowledge_title: [
    { required: true, message: $t('finance.linkKnowledge.enterKnowledgeTitle'), trigger: 'blur' }
  ]
}

const ticketTypeOptions = [
  { label: $t('finance.linkKnowledge.technical'), value: 'technical' },
  { label: $t('finance.linkKnowledge.service'), value: 'service' },
  { label: $t('finance.linkKnowledge.incident'), value: 'incident' },
  { label: $t('finance.linkKnowledge.change'), value: 'change' }
]

const getTicketTypeLabel = (value: string) => {
  const option = ticketTypeOptions.find(item => item.value === value)
  return option ? option.label : value
}

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/link-knowledge',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取知识库关联列表失败:', error)
    ElMessage.error($t('finance.linkKnowledge.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.ticket_type = ''
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = $t('finance.linkKnowledge.addLink')
  formData.id = undefined
  formData.ticket_type = ''
  formData.knowledge_id = ''
  formData.knowledge_title = ''
  formData.sort_order = 0
  formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('finance.linkKnowledge.editLink')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/link-knowledge/${row.id}`
    })
    ElMessage.success($t('finance.linkKnowledge.deleteSuccess'))
    fetchList()
  } catch (error) {
    ElMessage.error($t('finance.linkKnowledge.deleteFailed'))
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
          url: `/api/admin/link-knowledge/${formData.id}`,
          params: { ...formData }
        })
      } else {
        await request.post({
          url: '/api/admin/link-knowledge',
          params: { ...formData }
        })
      }

      ElMessage.success(formData.id ? $t('finance.linkKnowledge.updateSuccess') : $t('finance.linkKnowledge.addSuccess'))
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      ElMessage.error($t('finance.linkKnowledge.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchList()
}

const handlePageChange = () => {
  fetchList()
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.link-knowledge-page {
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