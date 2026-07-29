<template>
  <div class="product-groups-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="分组名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="art-card">
      <div class="table-header">
        <h3>产品分组</h3>
        <el-button type="primary" @click="showDialog = true; resetForm()">
          <el-icon><Plus /></el-icon>新增分组
        </el-button>
      </div>
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="分组名称" min-width="160" />
        <el-table-column prop="code" label="标识" width="120" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="product_count" label="产品数" width="100" />
        <el-table-column prop="sort_order" label="排序" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="editRow(row)">编辑</el-button>
            <el-button link type="danger" @click="deleteRow(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @change="fetchData"
        style="margin-top: 16px; justify-content: flex-end"
      />
    </div>

    <el-dialog v-model="showDialog" :title="form.id ? '编辑分组' : '新增分组'" width="550px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称"><el-input v-model="form.name" placeholder="请输入分组名称" /></el-form-item>
        <el-form-item label="标识"><el-input v-model="form.code" placeholder="英文标识" :disabled="!!form.id" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort_order" :min="0" /></el-form-item>
        <el-form-item label="状态"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveForm">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const showDialog = ref(false)
const searchForm = ref({ keyword: '', status: '' as string | number })
const form = ref<any>({ id: 0, name: '', code: '', description: '', sort_order: 0, status: 1 })

const resetForm = () => { form.value = { id: 0, name: '', code: '', description: '', sort_order: 0, status: 1 } }
const resetSearch = () => { searchForm.value = { keyword: '', status: '' }; fetchData() }

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/product-groups', { params: { page: page.value, page_size: pageSize.value, ...searchForm.value } })
    list.value = data.data || []; total.value = data.total || 0
  } catch {} finally { loading.value = false }
}

const editRow = (row: any) => { form.value = { ...row }; showDialog.value = true }

const saveForm = async () => {
  try {
    if (form.value.id) { await request.put(`/admin/api/v1/product-groups/${form.value.id}`, form.value) }
    else { await request.post('/admin/api/v1/product-groups', form.value) }
    ElMessage.success('保存成功'); showDialog.value = false; fetchData()
  } catch { ElMessage.error('保存失败') }
}

const deleteRow = async (id: number) => {
  await ElMessageBox.confirm('确定删除该分组？', '确认')
  try { await request.delete(`/admin/api/v1/product-groups/${id}`); ElMessage.success('删除成功'); fetchData() } catch { ElMessage.error('删除失败') }
}

onMounted(fetchData)
</script>
<style scoped lang="scss">
.search-bar { margin-bottom: 16px; }
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; h3 { font-size: 18px; font-weight: 600; margin: 0; } }
</style>
