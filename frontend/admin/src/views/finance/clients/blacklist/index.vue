<template>
  <div class="blacklist-page">
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>添加黑名单</el-button>
      </div>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="type" label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag size="small">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="value" label="值" min-width="200" />
        <el-table-column prop="reason" label="原因" min-width="200" />
        <el-table-column prop="created_at" label="添加时间" width="170" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="添加黑名单" width="500px" @close="formRef?.resetFields()">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="80px">
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型">
            <el-option label="IP地址" value="ip" />
            <el-option label="邮箱" value="email" />
            <el-option label="手机号" value="phone" />
            <el-option label="域名" value="domain" />
          </el-select>
        </el-form-item>
        <el-form-item label="值" prop="value">
          <el-input v-model="formData.value" placeholder="请输入对应的值" />
        </el-form-item>
        <el-form-item label="原因" prop="reason">
          <el-input v-model="formData.reason" type="textarea" :rows="3" placeholder="请输入原因" />
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
const submitting = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({ type: 'ip', value: '', reason: '' })
const rules = {
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  value: [{ required: true, message: '请输入值', trigger: 'blur' }]
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = { ip: 'IP地址', email: '邮箱', phone: '手机号', domain: '域名' }
  return map[type] || type
}

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/blacklist' })
    tableData.value = data || []
  } catch (error) {
    console.error('获取黑名单失败:', error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  Object.assign(formData, { type: 'ip', value: '', reason: '' })
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除此黑名单记录吗？', '确认删除', { type: 'warning' })
    await request.del({ url: `/api/admin/blacklist/${row.id}` })
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
    await request.post({ url: '/api/admin/blacklist', data: formData })
    ElMessage.success('添加成功')
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('添加失败:', error)
  } finally {
    submitting.value = false
  }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.blacklist-page { padding: 16px; }
.action-card { margin-bottom: 16px; }
.action-bar { display: flex; justify-content: space-between; align-items: center; }
.table-card { :deep(.el-card__body) { padding: 0; } }
</style>
