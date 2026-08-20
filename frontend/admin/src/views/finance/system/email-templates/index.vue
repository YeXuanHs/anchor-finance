<template>
  <div class="email-templates-page">
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('emailTemplate.addTemplate') }}</el-button>
      </div>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" :label="$t('emailTemplate.templateName')" min-width="150" />
        <el-table-column prop="code" :label="$t('emailTemplate.templateCode')" width="150" />
        <el-table-column prop="subject" :label="$t('emailTemplate.emailSubject')" min-width="200" />
        <el-table-column prop="status" :label="$t('emailTemplate.status')" width="80" align="center">
          <template #default="{ row }"><el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">{{ row.status === 'active' ? $t('common.enable') : $t('common.disable') }}</el-tag></template>
        </el-table-column>
        <el-table-column :label="$t('emailTemplate.operations')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px" @close="formRef?.resetFields()">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item :label="$t('emailTemplate.templateName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('emailTemplate.enterTemplateName')" />
        </el-form-item>
        <el-form-item :label="$t('emailTemplate.templateCode')" prop="code">
          <el-input v-model="formData.code" :placeholder="$t('emailTemplate.enterTemplateCode')" :disabled="isEdit" />
        </el-form-item>
        <el-form-item :label="$t('emailTemplate.emailSubject')" prop="subject">
          <el-input v-model="formData.subject" :placeholder="$t('emailTemplate.enterEmailSubject')" />
        </el-form-item>
        <el-form-item :label="$t('emailTemplate.emailContent')" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="10" :placeholder="$t('emailTemplate.enterEmailContent')" />
          <div class="form-tip">{{ $t('emailTemplate.availableVars') }}</div>
        </el-form-item>
        <el-form-item :label="$t('emailTemplate.status')" prop="status">
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
const dialogTitle = ref($t('emailTemplate.addTemplate'))
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

const formData = reactive({ name: '', code: '', subject: '', content: '', status: 'active' })
const rules = {
  name: [{ required: true, message: () => $t('emailTemplate.enterTemplateName'), trigger: 'blur' }],
  code: [{ required: true, message: () => $t('emailTemplate.enterTemplateCode'), trigger: 'blur' }],
  subject: [{ required: true, message: () => $t('emailTemplate.enterEmailSubject'), trigger: 'blur' }],
  content: [{ required: true, message: () => $t('emailTemplate.enterEmailContent'), trigger: 'blur' }]
}

const fetchList = async () => {
  loading.value = true
  try { const data = await request.get({ url: '/api/admin/email-templates' }); tableData.value = data || [] } catch (error) { console.error('fetch email templates failed:', error) } finally { loading.value = false }
}

const handleAdd = () => { isEdit.value = false; dialogTitle.value = $t('emailTemplate.addTemplate'); editingId.value = null; Object.assign(formData, { name: '', code: '', subject: '', content: '', status: 'active' }); dialogVisible.value = true }
const handleEdit = (row: any) => { isEdit.value = true; dialogTitle.value = $t('emailTemplate.editTemplate'); editingId.value = row.id; Object.assign(formData, { name: row.name, code: row.code, subject: row.subject, content: row.content, status: row.status }); dialogVisible.value = true }

const handleDelete = async (row: any) => {
  try { await ElMessageBox.confirm($t('emailTemplate.confirmDelete', { name: row.name }), $t('common.tips'), { type: 'warning' }); await request.del({ url: `/api/admin/email-templates/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchList() } catch (error) { if (error !== 'cancel') console.error('delete failed:', error) }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate(); submitting.value = true
    if (isEdit.value && editingId.value) { await request.put({ url: `/api/admin/email-templates/${editingId.value}`, data: formData }); ElMessage.success($t('common.updateSuccess')) }
    else { await request.post({ url: '/api/admin/email-templates', data: formData }); ElMessage.success($t('common.addSuccess')) }
    dialogVisible.value = false; fetchList()
  } catch (error) { console.error('submit failed:', error) } finally { submitting.value = false }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.email-templates-page { padding: 16px; }
.action-card { margin-bottom: 16px; }
.form-tip { color: var(--el-text-color-secondary); font-size: 12px; margin-top: 4px; }
</style>
