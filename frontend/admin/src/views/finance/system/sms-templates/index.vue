<template>
  <div class="sms-templates-page">
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('smsTemplate.addTemplate') }}</el-button>
      </div>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" :label="$t('smsTemplate.templateName')" min-width="150" />
        <el-table-column prop="code" :label="$t('smsTemplate.templateCode')" width="150" />
        <el-table-column prop="content" :label="$t('smsTemplate.content')" min-width="300" show-overflow-tooltip />
        <el-table-column prop="status" :label="$t('smsTemplate.status')" width="80" align="center">
          <template #default="{ row }"><el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">{{ row.status === 'active' ? $t('common.enable') : $t('common.disable') }}</el-tag></template>
        </el-table-column>
        <el-table-column :label="$t('smsTemplate.operations')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" @close="formRef?.resetFields()">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item :label="$t('smsTemplate.templateName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('smsTemplate.enterTemplateName')" />
        </el-form-item>
        <el-form-item :label="$t('smsTemplate.templateCode')" prop="code">
          <el-input v-model="formData.code" :placeholder="$t('smsTemplate.enterTemplateCode')" :disabled="isEdit" />
        </el-form-item>
        <el-form-item :label="$t('smsTemplate.smsContent')" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="5" :placeholder="$t('smsTemplate.enterContent')" />
          <div class="form-tip">{{ $t('smsTemplate.availableVars') }}</div>
        </el-form-item>
        <el-form-item :label="$t('smsTemplate.status')" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio value="active">{{ $t('common.enable') }}</el-radio>
            <el-radio value="disabled">{{ $t('common.disable') }}</el-radio>
          </el-radio-group>
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
const dialogTitle = ref($t('smsTemplate.addTemplate'))
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

const formData = reactive({ name: '', code: '', content: '', status: 'active' })
const rules = {
  name: [{ required: true, message: () => $t('smsTemplate.enterTemplateName'), trigger: 'blur' }],
  code: [{ required: true, message: () => $t('smsTemplate.enterTemplateCode'), trigger: 'blur' }],
  content: [{ required: true, message: () => $t('smsTemplate.enterContent'), trigger: 'blur' }]
}

const fetchList = async () => { loading.value = true; try { const data = await request.get({ url: '/api/admin/sms/templates' }); tableData.value = data || [] } catch {} finally { loading.value = false } }

const handleAdd = () => { isEdit.value = false; dialogTitle.value = $t('smsTemplate.addTemplate'); editingId.value = null; Object.assign(formData, { name: '', code: '', content: '', status: 'active' }); dialogVisible.value = true }
const handleEdit = (row: any) => { isEdit.value = true; dialogTitle.value = $t('smsTemplate.editTemplate'); editingId.value = row.id; Object.assign(formData, { name: row.name, code: row.code, content: row.content, status: row.status }); dialogVisible.value = true }

const handleDelete = async (row: any) => {
  try { await ElMessageBox.confirm($t('smsTemplate.confirmDelete', { name: row.name }), $t('common.tips'), { type: 'warning' }); await request.del({ url: `/api/admin/sms/templates/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchList() } catch (error) { if (error !== 'cancel') console.error('delete failed:', error) }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try { await formRef.value.validate(); submitting.value = true
    if (isEdit.value && editingId.value) { await request.put({ url: `/api/admin/sms/templates/${editingId.value}`, data: formData }); ElMessage.success($t('common.updateSuccess')) }
    else { await request.post({ url: '/api/admin/sms/templates', data: formData }); ElMessage.success($t('common.addSuccess')) }
    dialogVisible.value = false; fetchList()
  } catch {} finally { submitting.value = false }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.sms-templates-page { padding: 16px; }
.action-card { margin-bottom: 16px; }
.action-bar { display: flex; justify-content: space-between; align-items: center; }
.table-card { :deep(.el-card__body) { padding: 0; } }
.form-tip { color: #909399; font-size: 12px; margin-top: 4px; }
</style>
