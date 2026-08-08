<template>
  <div class="customer-groups-page">
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>添加客户组</el-button>
      </div>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="组名" min-width="150" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="clients_count" label="客户数" width="100" align="center" />
        <el-table-column prop="discount" label="折扣" width="100" align="center">
          <template #default="{ row }">{{ row.discount ? `${row.discount}%` : '-' }}</template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" @close="formRef?.resetFields()">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="组名" prop="name">
          <el-input v-model="formData.name" placeholder="请输入组名" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="折扣(%)" prop="discount">
          <el-input-number v-model="formData.discount" :min="0" :max="100" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const tableData = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('添加客户组')
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

const formData = reactive({ name: '', description: '', discount: 0 })
const rules = { name: [{ required: true, message: '请输入组名', trigger: 'blur' }] }

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/client-groups' })
    tableData.value = data || []
  } catch (error) {
    console.error('获取客户组列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  dialogTitle.value = '添加客户组'
  editingId.value = null
  formData.name = ''
  formData.description = ''
  formData.discount = 0
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  isEdit.value = true
  dialogTitle.value = '编辑客户组'
  editingId.value = row.id
  Object.assign(formData, { name: row.name, description: row.description || '', discount: row.discount || 0 })
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除客户组 "${row.name}" 吗？`, '确认删除', { type: 'warning' })
    await request.del({ url: `/api/admin/client-groups/${row.id}` })
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') console.error('删除失败:', error)
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    submitting.value = true
    if (isEdit.value && editingId.value) {
      await request.put({ url: `/api/admin/client-groups/${editingId.value}`, data: formData })
      ElMessage.success('更新成功')
    } else {
      await request.post({ url: '/api/admin/client-groups', data: formData })
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    submitting.value = false
  }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.customer-groups-page { padding: 16px; }
.action-card { margin-bottom: 16px; }
.action-bar { display: flex; justify-content: space-between; align-items: center; }
.table-card { :deep(.el-card__body) { padding: 0; } }
</style>
