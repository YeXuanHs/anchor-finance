<template>
  <div class="link-knowledge-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>关联知识库管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加关联
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="知识库标题" clearable />
        </el-form-item>
        <el-form-item label="工单类型">
          <el-select v-model="searchForm.ticket_type" placeholder="请选择工单类型" clearable>
            <el-option label="全部" value="" />
            <el-option label="技术问题" value="technical" />
            <el-option label="服务请求" value="service" />
            <el-option label="故障报告" value="incident" />
            <el-option label="变更请求" value="change" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="ticket_type" label="工单类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ getTicketTypeLabel(row.ticket_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="knowledge_title" label="知识库标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="knowledge_id" label="知识库ID" width="100" align="center" />
        <el-table-column prop="sort_order" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该关联吗？" @confirm="handleDelete(row)">
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="工单类型" prop="ticket_type">
          <el-select v-model="formData.ticket_type" placeholder="请选择工单类型">
            <el-option label="技术问题" value="technical" />
            <el-option label="服务请求" value="service" />
            <el-option label="故障报告" value="incident" />
            <el-option label="变更请求" value="change" />
          </el-select>
        </el-form-item>
        <el-form-item label="知识库ID" prop="knowledge_id">
          <el-input v-model="formData.knowledge_id" placeholder="请输入知识库ID" />
        </el-form-item>
        <el-form-item label="知识库标题" prop="knowledge_title">
          <el-input v-model="formData.knowledge_title" placeholder="请输入知识库标题" />
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
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
const dialogTitle = ref('添加关联')
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
    { required: true, message: '请选择工单类型', trigger: 'change' }
  ],
  knowledge_id: [
    { required: true, message: '请输入知识库ID', trigger: 'blur' }
  ],
  knowledge_title: [
    { required: true, message: '请输入知识库标题', trigger: 'blur' }
  ]
}

const ticketTypeOptions = [
  { label: '技术问题', value: 'technical' },
  { label: '服务请求', value: 'service' },
  { label: '故障报告', value: 'incident' },
  { label: '变更请求', value: 'change' }
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
    ElMessage.error('获取知识库关联列表失败')
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
  dialogTitle.value = '添加关联'
  formData.id = undefined
  formData.ticket_type = ''
  formData.knowledge_id = ''
  formData.knowledge_title = ''
  formData.sort_order = 0
  formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: LinkKnowledge) => {
  dialogTitle.value = '编辑关联'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: LinkKnowledge) => {
  try {
    await request.del({
      url: `/api/admin/link-knowledge/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    ElMessage.error('删除失败')
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

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      ElMessage.error('操作失败')
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