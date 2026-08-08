<template>
  <div class="currencies-page">
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>添加货币</el-button>
      </div>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="code" label="代码" width="80" />
        <el-table-column prop="name" label="名称" min-width="120" />
        <el-table-column prop="symbol" label="符号" width="80" align="center" />
        <el-table-column prop="exchange_rate" label="汇率" width="100" align="center" />
        <el-table-column prop="is_default" label="默认" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_default ? 'success' : 'info'" size="small">
              {{ row.is_default ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-switch v-model="row.status" :active-value="'active'" :inactive-value="'disabled'" @change="handleToggleStatus(row)" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="warning" link size="small" @click="handleSetDefault(row)" :disabled="row.is_default">设为默认</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)" :disabled="row.is_default">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" @close="formRef?.resetFields()">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="货币代码" prop="code">
          <el-input v-model="formData.code" placeholder="如 USD, CNY, EUR" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="货币名称" prop="name">
          <el-input v-model="formData.name" placeholder="如 美元, 人民币, 欧元" />
        </el-form-item>
        <el-form-item label="符号" prop="symbol">
          <el-input v-model="formData.symbol" placeholder="如 $, ¥, €" style="width: 100px" />
        </el-form-item>
        <el-form-item label="汇率" prop="exchange_rate">
          <el-input-number v-model="formData.exchange_rate" :min="0" :precision="4" :step="0.01" />
          <span class="form-tip">相对于默认货币</span>
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
const dialogTitle = ref('添加货币')
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

const formData = reactive({ code: '', name: '', symbol: '', exchange_rate: 1 })
const rules = {
  code: [{ required: true, message: '请输入货币代码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入货币名称', trigger: 'blur' }],
  symbol: [{ required: true, message: '请输入符号', trigger: 'blur' }]
}

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/currencies' })
    tableData.value = data || []
  } catch (error) {
    console.error('获取货币列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  dialogTitle.value = '添加货币'
  editingId.value = null
  Object.assign(formData, { code: '', name: '', symbol: '', exchange_rate: 1 })
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  isEdit.value = true
  dialogTitle.value = '编辑货币'
  editingId.value = row.id
  Object.assign(formData, { code: row.code, name: row.name, symbol: row.symbol, exchange_rate: row.exchange_rate })
  dialogVisible.value = true
}

const handleToggleStatus = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/currencies/${row.id}/status`, data: { status: row.status } })
    ElMessage.success('状态更新成功')
  } catch (error) {
    console.error('更新状态失败:', error)
    fetchList()
  }
}

const handleSetDefault = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要将 ${row.name} 设为默认货币吗？`, '确认操作', { type: 'warning' })
    await request.post({ url: `/api/admin/currencies/${row.id}/set-default` })
    ElMessage.success('设置成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') console.error('设置失败:', error)
  }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除货币 "${row.name}" 吗？`, '确认删除', { type: 'warning' })
    await request.delete({ url: `/api/admin/currencies/${row.id}` })
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
      await request.put({ url: `/api/admin/currencies/${editingId.value}`, data: formData })
      ElMessage.success('更新成功')
    } else {
      await request.post({ url: '/api/admin/currencies', data: formData })
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
.currencies-page { padding: 16px; }
.action-card { margin-bottom: 16px; }
.action-bar { display: flex; justify-content: space-between; align-items: center; }
.table-card { :deep(.el-card__body) { padding: 0; } }
.form-tip { margin-left: 10px; font-size: 12px; color: #86909C; }
</style>
