<template>
  <div class="vouchers-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="券码">
          <el-input v-model="searchForm.code" placeholder="券码" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="未使用" value="unused" />
            <el-option label="已使用" value="used" />
            <el-option label="已过期" value="expired" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户">
          <el-input v-model="searchForm.username" placeholder="使用者" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="art-card">
      <div class="table-header">
        <h3>代金券管理</h3>
        <div>
          <el-button type="warning" @click="showBatchDialog = true">
            <el-icon><DocumentCopy /></el-icon>批量生成
          </el-button>
          <el-button type="primary" @click="showDialog = true; resetForm()">
            <el-icon><Plus /></el-icon>添加代金券
          </el-button>
        </div>
      </div>
      <el-table :data="list" v-loading="loading">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="code" label="券码" width="180" />
        <el-table-column prop="name" label="名称" min-width="120" />
        <el-table-column label="面额" width="100">
          <template #default="{ row }">
            <span class="amount">¥{{ row.amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="使用门槛" width="120">
          <template #default="{ row }">
            <span v-if="row.min_amount > 0">满¥{{ row.min_amount?.toFixed(2) }}</span>
            <span v-else>无门槛</span>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="使用者" width="120">
          <template #default="{ row }">{{ row.username || '-' }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status]" size="small">{{ statusMap[row.status] || row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="expire_at" label="过期时间" width="180" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
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

    <el-dialog v-model="showDialog" :title="form.id ? '编辑代金券' : '添加代金券'" width="550px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称"><el-input v-model="form.name" placeholder="代金券名称" /></el-form-item>
        <el-form-item label="券码"><el-input v-model="form.code" placeholder="留空则自动生成" :disabled="!!form.id" /></el-form-item>
        <el-form-item label="面额"><el-input-number v-model="form.amount" :min="0.01" :precision="2" /></el-form-item>
        <el-form-item label="使用门槛"><el-input-number v-model="form.min_amount" :min="0" :precision="2" placeholder="0为无门槛" /></el-form-item>
        <el-form-item label="过期时间">
          <el-date-picker v-model="form.expire_at" type="datetime" placeholder="选择过期时间" value-format="YYYY-MM-DD HH:mm:ss" />
        </el-form-item>
        <el-form-item label="指定用户">
          <el-input v-model="form.bind_username" placeholder="留空则不限制用户" clearable />
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveForm">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showBatchDialog" title="批量生成代金券" width="550px">
      <el-form :model="batchForm" label-width="100px">
        <el-form-item label="名称"><el-input v-model="batchForm.name" placeholder="代金券名称" /></el-form-item>
        <el-form-item label="面额"><el-input-number v-model="batchForm.amount" :min="0.01" :precision="2" /></el-form-item>
        <el-form-item label="使用门槛"><el-input-number v-model="batchForm.min_amount" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="生成数量"><el-input-number v-model="batchForm.count" :min="1" :max="1000" /></el-form-item>
        <el-form-item label="过期时间">
          <el-date-picker v-model="batchForm.expire_at" type="datetime" placeholder="选择过期时间" value-format="YYYY-MM-DD HH:mm:ss" />
        </el-form-item>
        <el-form-item label="券码前缀"><el-input v-model="batchForm.prefix" placeholder="如: VIP (留空随机)" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="batchForm.remark" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showBatchDialog = false">取消</el-button>
        <el-button type="primary" @click="handleBatchGenerate">生成</el-button>
      </template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, DocumentCopy } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const statusMap: Record<string, string> = { unused: '未使用', used: '已使用', expired: '已过期' }
const statusType: Record<string, string> = { unused: 'success', used: 'info', expired: 'danger' }

const loading = ref(false)
const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const showDialog = ref(false)
const showBatchDialog = ref(false)
const searchForm = ref({ code: '', status: '', username: '' })
const form = ref<any>({ id: 0, name: '', code: '', amount: 10, min_amount: 0, expire_at: '', bind_username: '', remark: '' })
const batchForm = ref({ name: '', amount: 10, min_amount: 0, count: 10, expire_at: '', prefix: '', remark: '' })

const resetForm = () => { form.value = { id: 0, name: '', code: '', amount: 10, min_amount: 0, expire_at: '', bind_username: '', remark: '' } }
const resetSearch = () => { searchForm.value = { code: '', status: '', username: '' }; fetchData() }

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/finance-vouchers', { params: { page: page.value, page_size: pageSize.value, ...searchForm.value } })
    list.value = data.data || []; total.value = data.total || 0
  } catch {} finally { loading.value = false }
}

const editRow = (row: any) => { form.value = { ...row }; showDialog.value = true }

const saveForm = async () => {
  try {
    if (form.value.id) { await request.put(`/admin/api/v1/finance-vouchers/${form.value.id}`, form.value) }
    else { await request.post('/admin/api/v1/finance-vouchers', form.value) }
    ElMessage.success('保存成功'); showDialog.value = false; fetchData()
  } catch { ElMessage.error('保存失败') }
}

const deleteRow = async (id: number) => {
  await ElMessageBox.confirm('确定删除该代金券？', '确认')
  try { await request.delete(`/admin/api/v1/finance-vouchers/${id}`); ElMessage.success('删除成功'); fetchData() } catch { ElMessage.error('删除失败') }
}

const handleBatchGenerate = async () => {
  try {
    const { data } = await request.post('/admin/api/v1/finance-vouchers/batch', batchForm.value)
    ElMessage.success(`成功生成 ${data.count || batchForm.value.count} 张代金券`)
    showBatchDialog.value = false; fetchData()
  } catch { ElMessage.error('批量生成失败') }
}

onMounted(fetchData)
</script>
<style scoped lang="scss">
.search-bar { margin-bottom: 16px; }
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; h3 { font-size: 18px; font-weight: 600; margin: 0; } }
.amount { color: var(--danger-color); font-weight: 600; }
</style>
