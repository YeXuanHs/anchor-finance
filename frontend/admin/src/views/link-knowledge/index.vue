<template>
  <div class="link-knowledge-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="标题/内容" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="关联类型">
          <el-select v-model="searchForm.link_type" placeholder="全部" clearable>
            <el-option v-for="t in linkTypes" :key="t.value" :label="t.label" :value="t.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>搜索
          </el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>知识关联管理</h3>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>新增关联
        </el-button>
      </div>

      <el-table :data="list" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="title" label="标题" min-width="180" show-overflow-tooltip />
        <el-table-column prop="link_type" label="关联类型" width="110">
          <template #default="{ row }">
            <el-tag size="small">{{ getLinkTypeLabel(row.link_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="source_type" label="来源类型" width="100">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ getSourceTypeLabel(row.source_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="target_type" label="目标类型" width="100">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ getSourceTypeLabel(row.target_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="target_title" label="关联目标" min-width="180" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-switch v-model="row.status" :active-value="1" :inactive-value="0" @change="handleToggleStatus(row)" />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="650px" :close-on-click-modal="false">
      <el-form :model="formData" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入关联标题" />
        </el-form-item>
        <el-form-item label="关联类型" prop="link_type">
          <el-select v-model="formData.link_type" placeholder="请选择关联类型">
            <el-option v-for="t in linkTypes" :key="t.value" :label="t.label" :value="t.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="来源类型" prop="source_type">
          <el-select v-model="formData.source_type" placeholder="请选择来源类型">
            <el-option v-for="t in sourceTypes" :key="t.value" :label="t.label" :value="t.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="来源对象" prop="source_id">
          <el-select v-model="formData.source_id" filterable remote :remote-method="(q: string) => searchEntities(q, 'source')" placeholder="搜索来源对象" style="width: 100%;">
            <el-option v-for="e in sourceOptions" :key="e.id" :label="e.title" :value="e.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标类型" prop="target_type">
          <el-select v-model="formData.target_type" placeholder="请选择目标类型">
            <el-option v-for="t in sourceTypes" :key="t.value" :label="t.label" :value="t.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标对象" prop="target_id">
          <el-select v-model="formData.target_id" filterable remote :remote-method="(q: string) => searchEntities(q, 'target')" placeholder="搜索目标对象" style="width: 100%;">
            <el-option v-for="e in targetOptions" :key="e.id" :label="e.title" :value="e.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="关联原因">
          <el-input v-model="formData.reason" type="textarea" :rows="2" placeholder="说明关联原因" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const submitting = ref(false)
const list = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref()
const sourceOptions = ref<any[]>([])
const targetOptions = ref<any[]>([])

const linkTypes = [
  { label: '相关推荐', value: 'related' },
  { label: '解决方案', value: 'solution' },
  { label: '前置依赖', value: 'dependency' },
  { label: '替代方案', value: 'alternative' },
  { label: '补充说明', value: 'supplement' }
]

const sourceTypes = [
  { label: '知识文章', value: 'knowledge' },
  { label: '工单', value: 'ticket' },
  { label: '产品', value: 'product' },
  { label: '常见问题', value: 'faq' },
  { label: '文档', value: 'document' }
]

const searchForm = ref({ keyword: '', link_type: '', status: '' as string | number })

const defaultForm = () => ({
  id: null, title: '', link_type: '', source_type: '', source_id: null,
  target_type: '', target_id: null, reason: '', status: 1
})
const formData = reactive(defaultForm())

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  link_type: [{ required: true, message: '请选择关联类型', trigger: 'change' }],
  source_type: [{ required: true, message: '请选择来源类型', trigger: 'change' }],
  source_id: [{ required: true, message: '请选择来源对象', trigger: 'change' }],
  target_type: [{ required: true, message: '请选择目标类型', trigger: 'change' }],
  target_id: [{ required: true, message: '请选择目标对象', trigger: 'change' }]
}

const getLinkTypeLabel = (val: string) => linkTypes.find(t => t.value === val)?.label || val
const getSourceTypeLabel = (val: string) => sourceTypes.find(t => t.value === val)?.label || val

const fetchList = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/link-knowledge', {
      params: { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    })
    list.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取列表失败')
  } finally {
    loading.value = false
  }
}

const searchEntities = async (query: string, field: string) => {
  if (!query) return
  try {
    const type = field === 'source' ? formData.source_type : formData.target_type
    const { data } = await request.get('/admin/api/v1/link-knowledge/search-entities', {
      params: { keyword: query, type }
    })
    if (field === 'source') sourceOptions.value = data.data?.list || []
    else targetOptions.value = data.data?.list || []
  } catch {}
}

const handleSearch = () => { currentPage.value = 1; fetchList() }
const resetSearch = () => { searchForm.value = { keyword: '', link_type: '', status: '' }; handleSearch() }
const handleSizeChange = (val: number) => { pageSize.value = val; fetchList() }
const handlePageChange = (val: number) => { currentPage.value = val; fetchList() }

const handleToggleStatus = async (row: any) => {
  try {
    await request.put(`/admin/api/v1/link-knowledge/${row.id}/status`, { status: row.status })
    ElMessage.success('状态已更新')
  } catch {
    row.status = row.status === 1 ? 0 : 1
    ElMessage.error('更新失败')
  }
}

const handleCreate = () => {
  Object.assign(formData, defaultForm())
  sourceOptions.value = []
  targetOptions.value = []
  dialogTitle.value = '新增关联'
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  Object.assign(formData, row)
  sourceOptions.value = row.source_id ? [{ id: row.source_id, title: row.title }] : []
  targetOptions.value = row.target_id ? [{ id: row.target_id, title: row.target_title }] : []
  dialogTitle.value = '编辑关联'
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try { await formRef.value?.validate() } catch { return }
  submitting.value = true
  try {
    if (formData.id) {
      await request.put(`/admin/api/v1/link-knowledge/${formData.id}`, formData)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/link-knowledge', formData)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchList()
  } catch {
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
  }
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(`确认删除关联「${row.title}」？`, '确认删除', { type: 'warning' }).then(async () => {
    try {
      await request.delete(`/admin/api/v1/link-knowledge/${row.id}`)
      ElMessage.success('删除成功')
      fetchList()
    } catch {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.link-knowledge-page {
  .table-header {
    display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
}
</style>
