<template>
  <div class="advanced-options-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="配置项名称" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchForm.group" placeholder="全部" clearable>
            <el-option v-for="g in groups" :key="g.value" :label="g.label" :value="g.value" />
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
        <h3>高级可配置项</h3>
        <div class="header-actions">
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>新增配置项
          </el-button>
          <el-button @click="handleResetAll">
            <el-icon><RefreshLeft /></el-icon>重置默认
          </el-button>
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>导出
          </el-button>
        </div>
      </div>

      <el-table :data="list" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="配置项名称" min-width="150" />
        <el-table-column prop="key" label="键名" min-width="150" show-overflow-tooltip />
        <el-table-column prop="group" label="分组" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ getGroupLabel(row.group) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="value" label="当前值" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            <span v-if="row.type === 'boolean'">{{ row.value ? '是' : '否' }}</span>
            <span v-else-if="row.type === 'json'" class="json-preview">{{ truncateJson(row.value) }}</span>
            <span v-else>{{ row.value ?? '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="说明" min-width="180" show-overflow-tooltip />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="warning" link size="small" @click="handleReset(row)">重置</el-button>
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
      <el-form :model="formData" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="配置项名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入配置项名称" />
        </el-form-item>
        <el-form-item label="键名" prop="key">
          <el-input v-model="formData.key" placeholder="如: system.max_upload_size" :disabled="!!formData.id" />
        </el-form-item>
        <el-form-item label="分组" prop="group">
          <el-select v-model="formData.group" placeholder="请选择分组">
            <el-option v-for="g in groups" :key="g.value" :label="g.label" :value="g.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型" @change="onTypeChange">
            <el-option v-for="t in types" :key="t.value" :label="t.label" :value="t.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="值" prop="value">
          <el-input v-if="formData.type === 'text'" v-model="formData.value" placeholder="请输入值" />
          <el-input-number v-else-if="formData.type === 'number'" v-model="formData.value" style="width: 100%;" />
          <el-switch v-else-if="formData.type === 'boolean'" v-model="formData.value" />
          <el-select v-else-if="formData.type === 'select'" v-model="formData.value" placeholder="请选择">
            <el-option v-for="opt in formData.options_list" :key="opt" :label="opt" :value="opt" />
          </el-select>
          <el-input v-else-if="formData.type === 'json'" v-model="formData.value" type="textarea" :rows="6" placeholder='请输入 JSON，如 {"key": "value"}' />
        </el-form-item>
        <el-form-item label="选项列表" v-if="formData.type === 'select'">
          <el-input v-model="formData.options_str" placeholder="每行一个选项" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="formData.description" type="textarea" :rows="2" placeholder="配置项说明" />
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
import { Search, Plus, RefreshLeft, Download } from '@element-plus/icons-vue'
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

const groups = [
  { label: '系统', value: 'system' },
  { label: '安全', value: 'security' },
  { label: '邮件', value: 'email' },
  { label: '支付', value: 'payment' },
  { label: '界面', value: 'ui' },
  { label: '其他', value: 'other' }
]

const types = [
  { label: '文本', value: 'text' },
  { label: '数字', value: 'number' },
  { label: '布尔', value: 'boolean' },
  { label: '下拉选择', value: 'select' },
  { label: 'JSON', value: 'json' }
]

const searchForm = ref({ keyword: '', group: '' })

const defaultForm = () => ({
  id: null,
  name: '',
  key: '',
  group: 'system',
  type: 'text',
  value: null as any,
  description: '',
  options_str: '',
  options_list: [] as string[]
})

const formData = reactive(defaultForm())

const rules = {
  name: [{ required: true, message: '请输入配置项名称', trigger: 'blur' }],
  key: [{ required: true, message: '请输入键名', trigger: 'blur' }],
  group: [{ required: true, message: '请选择分组', trigger: 'change' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }]
}

const getGroupLabel = (val: string) => groups.find(g => g.value === val)?.label || val
const getTypeLabel = (val: string) => types.find(t => t.value === val)?.label || val

const truncateJson = (val: any) => {
  if (!val) return '-'
  const str = typeof val === 'string' ? val : JSON.stringify(val)
  return str.length > 40 ? str.slice(0, 40) + '...' : str
}

const onTypeChange = () => { formData.value = null }

const fetchList = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/advanced-options', {
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

const handleSearch = () => { currentPage.value = 1; fetchList() }
const resetSearch = () => { searchForm.value = { keyword: '', group: '' }; handleSearch() }
const handleSizeChange = (val: number) => { pageSize.value = val; fetchList() }
const handlePageChange = (val: number) => { currentPage.value = val; fetchList() }

const handleCreate = () => {
  Object.assign(formData, defaultForm())
  dialogTitle.value = '新增配置项'
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  Object.assign(formData, {
    ...row,
    options_str: Array.isArray(row.options_list) ? row.options_list.join('\n') : '',
    value: row.type === 'json' && typeof row.value === 'object' ? JSON.stringify(row.value, null, 2) : row.value
  })
  dialogTitle.value = '编辑配置项'
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch { return }

  submitting.value = true
  try {
    const payload = { ...formData }
    if (payload.type === 'select' && payload.options_str) {
      payload.options_list = payload.options_str.split('\n').map((s: string) => s.trim()).filter(Boolean)
    }
    if (payload.type === 'json' && typeof payload.value === 'string') {
      try { payload.value = JSON.parse(payload.value) } catch { /* keep as string */ }
    }

    if (formData.id) {
      await request.put(`/admin/api/v1/advanced-options/${formData.id}`, payload)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/advanced-options', payload)
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
  ElMessageBox.confirm(`确认删除配置项「${row.name}」？`, '确认删除', { type: 'warning' }).then(async () => {
    try {
      await request.delete(`/admin/api/v1/advanced-options/${row.id}`)
      ElMessage.success('删除成功')
      fetchList()
    } catch {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

const handleReset = (row: any) => {
  ElMessageBox.confirm(`确认将「${row.name}」重置为默认值？`, '确认重置', { type: 'warning' }).then(async () => {
    try {
      await request.post(`/admin/api/v1/advanced-options/${row.id}/reset`)
      ElMessage.success('已重置')
      fetchList()
    } catch {
      ElMessage.error('重置失败')
    }
  }).catch(() => {})
}

const handleResetAll = () => {
  ElMessageBox.confirm('确认将所有配置项重置为默认值？此操作不可恢复。', '确认重置', { type: 'warning' }).then(async () => {
    try {
      await request.post('/admin/api/v1/advanced-options/reset-all')
      ElMessage.success('已重置所有配置')
      fetchList()
    } catch {
      ElMessage.error('重置失败')
    }
  }).catch(() => {})
}

const handleExport = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/advanced-options/export', { responseType: 'blob' })
    const url = URL.createObjectURL(data)
    const a = document.createElement('a')
    a.href = url
    a.download = 'advanced-options.json'
    a.click()
    URL.revokeObjectURL(url)
  } catch {
    ElMessage.error('导出失败')
  }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.advanced-options-page {
  .table-header {
    display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
    .header-actions { display: flex; gap: 8px; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
  .json-preview { font-family: monospace; font-size: 12px; color: var(--text-secondary); }
}
</style>
