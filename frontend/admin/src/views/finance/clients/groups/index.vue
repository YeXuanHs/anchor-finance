<template>
  <div class="customer-groups-page">
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('clientGroup.addGroup') }}</el-button>
      </div>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" :label="$t('clientGroup.groupName')" min-width="150" />
        <el-table-column prop="description" :label="$t('common.description')" min-width="200" />
        <el-table-column prop="clients_count" :label="$t('clientGroup.clientCount')" width="100" align="center" />
        <el-table-column prop="discount" :label="$t('clientGroup.discount')" width="100" align="center">
          <template #default="{ row }">{{ row.discount ? `${row.discount}%` : '-' }}</template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="170" />
        <el-table-column :label="$t('clientGroup.operations')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" @close="formRef?.resetFields()">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item :label="$t('clientGroup.groupName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('clientGroup.enterGroupName')" />
        </el-form-item>
        <el-form-item :label="$t('common.description')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('clientGroup.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('clientGroup.discount')" prop="discount">
          <el-input-number v-model="formData.discount" :min="0" :max="100" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">{{ $t('common.confirm') }}</el-button>
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
import { $t } from '@/locales'

const loading = ref(false)
const tableData = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref($t('clientGroup.addGroup'))
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

const formData = reactive({ name: '', description: '', discount: 0 })
const rules = { name: [{ required: true, message: () => $t('clientGroup.enterGroupName'), trigger: 'blur' }] }

const fetchList = async () => {
  loading.value = true
  try { const data = await request.get({ url: '/api/admin/client-groups' }); tableData.value = data || [] } catch (error) { console.error('fetch client groups failed:', error) } finally { loading.value = false }
}

const handleAdd = () => { isEdit.value = false; dialogTitle.value = $t('clientGroup.addGroup'); editingId.value = null; formData.name = ''; formData.description = ''; formData.discount = 0; dialogVisible.value = true }
const handleEdit = (row: any) => { isEdit.value = true; dialogTitle.value = $t('clientGroup.editGroup'); editingId.value = row.id; Object.assign(formData, { name: row.name, description: row.description || '', discount: row.discount || 0 }); dialogVisible.value = true }

const handleDelete = async (row: any) => {
  try { await ElMessageBox.confirm($t('clientGroup.confirmDelete', { name: row.name }), $t('common.tips'), { type: 'warning' }); await request.del({ url: `/api/admin/client-groups/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchList() } catch (error) { if (error !== 'cancel') console.error('delete failed:', error) }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate(); submitting.value = true
    if (isEdit.value && editingId.value) { await request.put({ url: `/api/admin/client-groups/${editingId.value}`, data: formData }); ElMessage.success($t('common.updateSuccess')) }
    else { await request.post({ url: '/api/admin/client-groups', data: formData }); ElMessage.success($t('common.addSuccess')) }
    dialogVisible.value = false; fetchList()
  } catch (error) { console.error('submit failed:', error) } finally { submitting.value = false }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.customer-groups-page { padding: 16px; }
.action-card { margin-bottom: 16px; }
</style>
