<template>
  <div class="zjmf-api-page">
    <!-- 说明文字 -->
    <el-alert type="info" :closable="false" class="page-desc">
      <template #title>
        <span>提供财务系统之间的互相代理，在此添加API接口信息即可成为代理商。</span>
        <el-link type="primary" href="https://bbs.idcsmart.com/forum.php?mod=viewthread&tid=136" target="_blank" :underline="false">帮助文档</el-link>
      </template>
    </el-alert>

    <div class="action-bar">
      <el-button type="primary" @click="handleAdd">添加供应商</el-button>
    </div>

    <!-- 供应商表格 -->
    <el-table 
      :data="supplierList" 
      v-loading="loading" 
      border 
      stripe
      empty-text="暂无数据"
      style="width: 100%"
      row-key="id"
    >
      <el-table-column prop="name" label="名称" min-width="120">
        <template #default="{ row }">
          <el-link type="primary" @click="handleEdit(row)">{{ row.name }}</el-link>
        </template>
      </el-table-column>
      
      <el-table-column prop="type" label="类型" width="120">
        <template #default="{ row }">
          <el-tag>{{ getTypeName(row.type) }}</el-tag>
        </template>
      </el-table-column>
      
      <el-table-column prop="api_url" label="接口地址" min-width="200" show-overflow-tooltip />
      
      <el-table-column label="可售/已设置商品" width="140">
        <template #default="{ row }">
          {{ row.available_products || 0 }}/{{ row.set_products || 0 }}
        </template>
      </el-table-column>
      
      <el-table-column label="产品数量(正常/总)" width="160">
        <template #default="{ row }">
          {{ row.normal_products || 0 }}/{{ row.total_products || 0 }}
        </template>
      </el-table-column>
      
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-switch 
            v-model="row.is_active" 
            @change="handleToggle(row)"
          />
        </template>
      </el-table-column>
      
      <el-table-column prop="balance" label="余额" width="100">
        <template #default="{ row }">
          {{ row.balance ?? '-' }}
        </template>
      </el-table-column>
      
      <el-table-column prop="description" label="描述" min-width="100">
        <template #default="{ row }">
          {{ row.description || '-' }}
        </template>
      </el-table-column>
      
      <el-table-column label="管理" width="140" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" link @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" link @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchList"
        @current-change="fetchList"
      />
    </div>

    <!-- 新增/编辑弹窗 -->
    <el-dialog 
      v-model="editDialogVisible" 
      :title="isEdit ? '编辑供应商' : '添加供应商'" 
      width="600px"
      destroy-on-close
    >
      <el-form :model="formData" label-position="top" :rules="formRules" ref="formRef">
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入供应商名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" style="width: 100%">
            <el-option label="手动" value="manual" />
            <el-option label="智简魔方" value="zjmf" />
            <el-option label="v10" value="v10" />
            <el-option label="锚点财务" value="anchor" />
          </el-select>
        </el-form-item>
        <el-form-item label="接口地址" prop="api_url">
          <el-input v-model="formData.api_url" placeholder="请输入API接口地址" />
        </el-form-item>
        <el-form-item label="接口密钥" prop="api_key">
          <el-input v-model="formData.api_key" type="password" placeholder="请输入接口密钥" show-password />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="formData.description" type="textarea" :rows="2" placeholder="请输入描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const supplierList = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const editDialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()
const saving = ref(false)

const formData = ref({
  name: '',
  type: 'zjmf',
  api_url: '',
  api_key: '',
  description: ''
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入供应商名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  api_url: [{ required: true, message: '请输入接口地址', trigger: 'blur' }]
}

const typeMap: Record<string, string> = {
  manual: '手动',
  zjmf: '智简魔方',
  v10: 'v10',
  anchor: '锚点财务'
}

const getTypeName = (type: string) => typeMap[type] || type

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ 
      url: '/api/admin/suppliers',
      params: { page: currentPage.value, page_size: pageSize.value }
    })
    supplierList.value = data?.list || data || []
    total.value = data?.total || supplierList.value.length
  } catch (error) {
    console.error('fetch suppliers failed:', error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  editingId.value = null
  formData.value = { name: '', type: 'zjmf', api_url: '', api_key: '', description: '' }
  editDialogVisible.value = true
}

const handleEdit = (row: any) => {
  isEdit.value = true
  editingId.value = row.id
  formData.value = { 
    name: row.name, 
    type: row.type, 
    api_url: row.api_url, 
    api_key: row.api_key || '', 
    description: row.description || '' 
  }
  editDialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 "${row.name}" 吗？删除后无法恢复。`, 
      '确认删除', 
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    await request.del({ url: `/api/admin/suppliers/${row.id}` })
    ElMessage.success('已删除')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleToggle = async (row: any) => {
  const original = row.is_active
  try {
    await request.post({ 
      url: `/api/admin/suppliers/${row.id}/toggle`, 
      data: { is_active: row.is_active } 
    })
    ElMessage.success(row.is_active ? '已启用' : '已禁用')
  } catch (error) {
    row.is_active = original
    ElMessage.error('操作失败')
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    saving.value = true
    
    if (isEdit.value && editingId.value) {
      await request.put({ url: `/api/admin/suppliers/${editingId.value}`, data: formData.value })
    } else {
      await request.post({ url: '/api/admin/suppliers', data: formData.value })
    }
    
    ElMessage.success(isEdit.value ? '编辑成功' : '添加成功')
    editDialogVisible.value = false
    fetchList()
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.zjmf-api-page {
  padding: 20px;
}

.page-desc {
  margin-bottom: 16px;
  :deep(.el-alert__title) {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

.action-bar {
  margin-bottom: 16px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
