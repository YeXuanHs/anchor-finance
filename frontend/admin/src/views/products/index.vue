<template>
  <div class="products-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="产品名称/编码" clearable />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchForm.category_id" placeholder="全部" clearable>
            <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="上架" value="active" />
            <el-option label="下架" value="inactive" />
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
        <h3>产品列表</h3>
        <el-button type="primary" @click="showDialog = true; resetForm()">
          <el-icon><Plus /></el-icon>添加产品
        </el-button>
      </div>
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="产品名称" min-width="160" />
        <el-table-column prop="category_name" label="分类" width="120" />
        <el-table-column prop="code" label="编码" width="120" />
        <el-table-column label="价格" width="140">
          <template #default="{ row }">
            <span class="amount">¥{{ row.price?.toFixed(2) }}</span>
            <span class="text-secondary">/{{ row.billing_cycle || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="80" />
        <el-table-column prop="sales" label="销量" width="80" />
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column prop="is_enabled" label="状态" width="100">
          <template #default="{ row }">
            <el-switch v-model="row.is_enabled" @change="toggleStatus(row)" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="editRow(row)">编辑</el-button>
            <el-button link type="primary" @click="configProduct(row)">配置</el-button>
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

    <el-dialog v-model="showDialog" :title="form.id ? '编辑产品' : '添加产品'" width="700px">
      <el-form :model="form" label-width="100px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="产品名称"><el-input v-model="form.name" /></el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="产品编码"><el-input v-model="form.code" :disabled="!!form.id" /></el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="分类">
              <el-select v-model="form.category_id" placeholder="请选择" clearable>
                <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分组">
              <el-select v-model="form.group_id" placeholder="请选择" clearable>
                <el-option v-for="g in groups" :key="g.id" :label="g.name" :value="g.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="价格"><el-input-number v-model="form.price" :min="0" :precision="2" /></el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="计费周期">
              <el-select v-model="form.billing_cycle">
                <el-option label="月付" value="monthly" />
                <el-option label="季付" value="quarterly" />
                <el-option label="半年付" value="semi-annually" />
                <el-option label="年付" value="annually" />
                <el-option label="一次性" value="onetime" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="库存"><el-input-number v-model="form.stock" :min="0" /></el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort_order" :min="0" /></el-form-item>
        <el-form-item label="启用"><el-switch v-model="form.is_enabled" /></el-form-item>
        <el-divider content-position="left">价格配置（JSON）</el-divider>
        <el-form-item label="配置">
          <el-input v-model="form.price_config" type="textarea" :rows="4" placeholder='{"monthly":99,"quarterly":269,"annually":999}' />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveForm">保存</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="showConfig" title="产品配置" size="600px">
      <el-form :model="configForm" label-width="100px">
        <el-divider content-position="left">服务器配置</el-divider>
        <el-form-item label="服务器组">
          <el-select v-model="configForm.server_group_id" placeholder="请选择" clearable>
            <el-option label="默认组" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作系统">
          <el-input v-model="configForm.os_options" type="textarea" :rows="3" placeholder="CentOS, Ubuntu, Debian (每行一个)" />
        </el-form-item>
        <el-divider content-position="left">资源配置</el-divider>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="CPU"><el-input v-model="configForm.cpu" placeholder="如: 2核" /></el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="内存"><el-input v-model="configForm.memory" placeholder="如: 4GB" /></el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="硬盘"><el-input v-model="configForm.disk" placeholder="如: 50GB SSD" /></el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="带宽"><el-input v-model="configForm.bandwidth" placeholder="如: 10Mbps" /></el-form-item>
          </el-col>
        </el-row>
        <el-divider content-position="left">自定义配置（JSON）</el-divider>
        <el-form-item label="配置">
          <el-input v-model="configForm.custom_config" type="textarea" :rows="6" placeholder="{}" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showConfig = false">取消</el-button>
        <el-button type="primary" @click="saveConfig">保存配置</el-button>
      </template>
    </el-drawer>
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
const showConfig = ref(false)
const categories = ref<any[]>([])
const groups = ref<any[]>([])
const searchForm = ref({ keyword: '', category_id: '', status: '' })
const form = ref<any>({ id: 0, name: '', code: '', category_id: '', group_id: '', price: 0, billing_cycle: 'monthly', stock: 999, description: '', sort_order: 0, is_enabled: true, price_config: '{}' })
const configForm = ref<any>({ server_group_id: '', os_options: '', cpu: '', memory: '', disk: '', bandwidth: '', custom_config: '{}' })

const resetForm = () => { form.value = { id: 0, name: '', code: '', category_id: '', group_id: '', price: 0, billing_cycle: 'monthly', stock: 999, description: '', sort_order: 0, is_enabled: true, price_config: '{}' } }
const resetSearch = () => { searchForm.value = { keyword: '', category_id: '', status: '' }; fetchData() }

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/products', { params: { page: page.value, page_size: pageSize.value, ...searchForm.value } })
    list.value = data.data || []; total.value = data.total || 0
  } catch {} finally { loading.value = false }
}

const fetchOptions = async () => {
  try {
    const [catRes, groupRes] = await Promise.all([
      request.get('/admin/api/v1/product-categories', { params: { page_size: 999 } }),
      request.get('/admin/api/v1/product-groups', { params: { page_size: 999 } })
    ])
    categories.value = catRes.data.data || []
    groups.value = groupRes.data.data || []
  } catch {}
}

const editRow = (row: any) => {
  form.value = { ...row, price_config: typeof row.price_config === 'object' ? JSON.stringify(row.price_config) : row.price_config }
  showDialog.value = true
}

const saveForm = async () => {
  try {
    if (form.value.id) { await request.put(`/admin/api/v1/products/${form.value.id}`, form.value) }
    else { await request.post('/admin/api/v1/products', form.value) }
    ElMessage.success('保存成功'); showDialog.value = false; fetchData()
  } catch { ElMessage.error('保存失败') }
}

const toggleStatus = async (row: any) => {
  try { await request.put(`/admin/api/v1/products/${row.id}`, { is_enabled: row.is_enabled }) } catch { row.is_enabled = !row.is_enabled }
}

const deleteRow = async (id: number) => {
  await ElMessageBox.confirm('确定删除该产品？', '确认')
  try { await request.delete(`/admin/api/v1/products/${id}`); ElMessage.success('删除成功'); fetchData() } catch { ElMessage.error('删除失败') }
}

const configProduct = async (row: any) => {
  try {
    const { data } = await request.get(`/admin/api/v1/products/${row.id}/config`)
    configForm.value = {
      product_id: row.id,
      server_group_id: data.server_group_id || '',
      os_options: Array.isArray(data.os_options) ? data.os_options.join('\n') : (data.os_options || ''),
      cpu: data.cpu || '',
      memory: data.memory || '',
      disk: data.disk || '',
      bandwidth: data.bandwidth || '',
      custom_config: typeof data.custom_config === 'object' ? JSON.stringify(data.custom_config) : (data.custom_config || '{}')
    }
  } catch {
    configForm.value = { product_id: row.id, server_group_id: '', os_options: '', cpu: '', memory: '', disk: '', bandwidth: '', custom_config: '{}' }
  }
  showConfig.value = true
}

const saveConfig = async () => {
  try {
    await request.put(`/admin/api/v1/products/${configForm.value.product_id}/config`, configForm.value)
    ElMessage.success('配置保存成功'); showConfig.value = false
  } catch { ElMessage.error('配置保存失败') }
}

onMounted(() => { fetchData(); fetchOptions() })
</script>
<style scoped lang="scss">
.search-bar { margin-bottom: 16px; }
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; h3 { font-size: 18px; font-weight: 600; margin: 0; } }
.amount { color: var(--danger-color); font-weight: 600; }
.text-secondary { color: var(--text-secondary); font-size: 12px; }
</style>
