<template>
  <div class="ticket-status-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="状态名称/代码" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="开放" value="open" />
            <el-option label="处理中" value="in_progress" />
            <el-option label="已关闭" value="closed" />
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
        <h3>工单状态</h3>
        <el-button type="primary" @click="openAddDialog">
          <el-icon><Plus /></el-icon>
          添加状态
        </el-button>
      </div>

      <el-table :data="statuses" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="状态名称" width="150">
          <template #default="{ row }">
            <div class="status-name-cell">
              <span class="color-dot" :style="{ background: row.color }"></span>
              {{ row.name }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="code" label="状态代码" width="150" />
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="row.type === 'open' ? 'warning' : row.type === 'closed' ? 'success' : 'info'" size="small">
              {{ getTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ticket_count" label="工单数" width="100" align="center" />
        <el-table-column prop="is_system" label="系统内置" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_system ? 'danger' : 'info'" size="small">
              {{ row.is_system ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="80" align="center" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editStatus(row)">编辑</el-button>
            <el-button type="danger" link @click="deleteStatus(row)" :disabled="row.is_system">删除</el-button>
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
          @size-change="fetchStatuses"
          @current-change="fetchStatuses"
        />
      </div>
    </div>

    <el-dialog v-model="showFormDialog" :title="editingItem ? '编辑状态' : '添加状态'" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="状态名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入状态名称" />
        </el-form-item>
        <el-form-item label="状态代码" prop="code">
          <el-input v-model="formData.code" placeholder="如: open, pending, closed" :disabled="editingItem?.is_system" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型">
            <el-option label="开放" value="open" />
            <el-option label="处理中" value="in_progress" />
            <el-option label="已关闭" value="closed" />
          </el-select>
        </el-form-item>
        <el-form-item label="颜色">
          <el-color-picker v-model="formData.color" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="formData.sort_order" :min="0" :max="9999" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFormDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const statuses = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showFormDialog = ref(false)
const editingItem = ref<any>(null)
const formRef = ref<FormInstance>()

const searchForm = ref({ keyword: '', type: '' })
const formData = ref({ name: '', code: '', type: 'open', color: '#409EFF', sort_order: 0 })
const formRules = {
  name: [{ required: true, message: '请输入状态名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入状态代码', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }]
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = { open: '开放', in_progress: '处理中', closed: '已关闭' }
  return map[type] || type
}

const fetchStatuses = async () => {
  loading.value = true
  try {
    const params = { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    const { data } = await request.get('/admin/api/v1/tickets/statuses', { params })
    statuses.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取状态列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { currentPage.value = 1; fetchStatuses() }
const resetSearch = () => { searchForm.value = { keyword: '', type: '' }; handleSearch() }

const openAddDialog = () => {
  editingItem.value = null
  formData.value = { name: '', code: '', type: 'open', color: '#409EFF', sort_order: 0 }
  showFormDialog.value = true
}

const editStatus = (status: any) => {
  editingItem.value = status
  formData.value = { name: status.name, code: status.code, type: status.type, color: status.color, sort_order: status.sort_order }
  showFormDialog.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (editingItem.value) {
      await request.put(`/admin/api/v1/tickets/statuses/${editingItem.value.id}`, formData.value)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/tickets/statuses', formData.value)
      ElMessage.success('添加成功')
    }
    showFormDialog.value = false
    fetchStatuses()
  } catch {
    ElMessage.error(editingItem.value ? '更新失败' : '添加失败')
  } finally {
    submitLoading.value = false
  }
}

const deleteStatus = async (status: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除状态「${status.name}」吗？关联工单将受影响。`, '提示', { type: 'warning' })
    await request.delete(`/admin/api/v1/tickets/statuses/${status.id}`)
    ElMessage.success('删除成功')
    fetchStatuses()
  } catch {}
}

onMounted(fetchStatuses)
</script>

<style scoped lang="scss">
.ticket-status-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
  .status-name-cell { display: flex; align-items: center; gap: 8px; }
  .color-dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    display: inline-block;
    flex-shrink: 0;
  }
}
</style>
