<template>
  <div class="advanced-options-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="选项名称" clearable />
        </el-form-item>
        <el-form-item label="产品">
          <el-select v-model="searchForm.product_id" placeholder="全部" clearable filterable>
            <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
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
        <h3>高级配置选项</h3>
        <div>
          <el-button type="success" @click="showJsonEditor = true">
            <el-icon><EditPen /></el-icon>JSON配置
          </el-button>
          <el-button type="primary" @click="showDialog = true; resetForm()">
            <el-icon><Plus /></el-icon>添加选项
          </el-button>
        </div>
      </div>
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="选项名称" min-width="150" />
        <el-table-column prop="product_name" label="关联产品" width="140" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ typeMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="价格" width="120">
          <template #default="{ row }">
            <span class="amount">¥{{ row.price?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="required" label="必选" width="80">
          <template #default="{ row }">
            <el-tag :type="row.required ? 'danger' : 'info'" size="small">{{ row.required ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column label="操作" width="150" fixed="right">
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

    <el-dialog v-model="showDialog" :title="form.id ? '编辑选项' : '添加选项'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="选项名称"><el-input v-model="form.name" placeholder="请输入选项名称" /></el-form-item>
        <el-form-item label="关联产品">
          <el-select v-model="form.product_id" placeholder="请选择产品" filterable clearable>
            <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.type">
            <el-option label="文本" value="text" />
            <el-option label="数字" value="number" />
            <el-option label="下拉选择" value="select" />
            <el-option label="复选框" value="checkbox" />
            <el-option label="单选框" value="radio" />
          </el-select>
        </el-form-item>
        <el-form-item label="可选值" v-if="form.type === 'select' || form.type === 'radio'">
          <el-input v-model="form.options" type="textarea" :rows="4" placeholder="每行一个选项，格式: value|label" />
        </el-form-item>
        <el-form-item label="默认值"><el-input v-model="form.default_value" placeholder="默认值" /></el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="价格"><el-input-number v-model="form.price" :min="0" :precision="2" /></el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序"><el-input-number v-model="form.sort_order" :min="0" /></el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="必选"><el-switch v-model="form.required" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveForm">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showJsonEditor" title="JSON 批量配置" width="700px">
      <el-alert title="以JSON数组格式编辑高级选项配置，保存后将覆盖当前配置。" type="warning" :closable="false" style="margin-bottom: 16px" />
      <el-input v-model="jsonContent" type="textarea" :rows="16" placeholder='[{"name":"选项名","type":"text","price":0}]' />
      <template #footer>
        <el-button @click="showJsonEditor = false">取消</el-button>
        <el-button type="primary" @click="saveJson">保存配置</el-button>
      </template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, EditPen } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const typeMap: Record<string, string> = { text: '文本', number: '数字', select: '下拉选择', checkbox: '复选框', radio: '单选框' }
const loading = ref(false)
const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const showDialog = ref(false)
const showJsonEditor = ref(false)
const productList = ref<any[]>([])
const jsonContent = ref('[]')
const searchForm = ref({ keyword: '', product_id: '' })
const form = ref<any>({ id: 0, name: '', product_id: '', type: 'text', options: '', default_value: '', price: 0, required: false, sort_order: 0, description: '' })

const resetForm = () => { form.value = { id: 0, name: '', product_id: '', type: 'text', options: '', default_value: '', price: 0, required: false, sort_order: 0, description: '' } }
const resetSearch = () => { searchForm.value = { keyword: '', product_id: '' }; fetchData() }

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/product-options', { params: { page: page.value, page_size: pageSize.value, ...searchForm.value } })
    list.value = data.data || []; total.value = data.total || 0
  } catch {} finally { loading.value = false }
}

const fetchProducts = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/products', { params: { page_size: 999 } })
    productList.value = data.data || []
  } catch {}
}

const editRow = (row: any) => {
  form.value = { ...row, options: Array.isArray(row.options) ? row.options.join('\n') : row.options }
  showDialog.value = true
}

const saveForm = async () => {
  try {
    if (form.value.id) { await request.put(`/admin/api/v1/product-options/${form.value.id}`, form.value) }
    else { await request.post('/admin/api/v1/product-options', form.value) }
    ElMessage.success('保存成功'); showDialog.value = false; fetchData()
  } catch { ElMessage.error('保存失败') }
}

const deleteRow = async (id: number) => {
  await ElMessageBox.confirm('确定删除该选项？', '确认')
  try { await request.delete(`/admin/api/v1/product-options/${id}`); ElMessage.success('删除成功'); fetchData() } catch { ElMessage.error('删除失败') }
}

const saveJson = async () => {
  try {
    const parsed = JSON.parse(jsonContent.value)
    await request.post('/admin/api/v1/product-options/batch', { data: parsed })
    ElMessage.success('JSON配置保存成功'); showJsonEditor.value = false; fetchData()
  } catch { ElMessage.error('JSON格式错误或保存失败') }
}

onMounted(() => { fetchData(); fetchProducts() })
</script>
<style scoped lang="scss">
.search-bar { margin-bottom: 16px; }
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; h3 { font-size: 18px; font-weight: 600; margin: 0; } }
.amount { color: var(--danger-color); font-weight: 600; }
</style>
