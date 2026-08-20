<template>
  <div class="custom-fields-page">
    <art-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('systemCustomFields.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> {{ $t('systemCustomFields.addField') }}
          </el-button>
        </div>
      </template>

      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="60" align="center" />
        <el-table-column prop="field_name" :label="$t('systemCustomFields.fieldName')" min-width="150" />
        <el-table-column prop="field_key" :label="$t('systemCustomFields.fieldKey')" min-width="150" />
        <el-table-column prop="field_type" :label="$t('common.type')" width="100" />
        <el-table-column prop="sort_order" :label="$t('common.sort')" width="80" align="center" />
        <el-table-column :label="$t('common.action')" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </art-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? $t('systemCustomFields.editField') : $t('systemCustomFields.addField')" width="500px" @close="resetForm">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item :label="$t('systemCustomFields.fieldName')" prop="field_name">
          <el-input v-model="formData.field_name" :placeholder="$t('systemCustomFields.enterFieldName')" />
        </el-form-item>
        <el-form-item :label="$t('systemCustomFields.fieldKey')" prop="field_key">
          <el-input v-model="formData.field_key" :placeholder="$t('systemCustomFields.enterFieldKey')" :disabled="isEdit" />
        </el-form-item>
        <el-form-item :label="$t('common.type')" prop="field_type">
          <el-select v-model="formData.field_type" :placeholder="$t('systemCustomFields.selectType')" style="width: 100%">
            <el-option :label="$t('systemCustomFields.typeText')" value="text" />
            <el-option :label="$t('systemCustomFields.typeTextarea')" value="textarea" />
            <el-option :label="$t('systemCustomFields.typeSelect')" value="select" />
            <el-option :label="$t('systemCustomFields.typeRadio')" value="radio" />
            <el-option :label="$t('systemCustomFields.typeCheckbox')" value="checkbox" />
            <el-option :label="$t('systemCustomFields.typeDate')" value="date" />
            <el-option :label="$t('systemCustomFields.typeNumber')" value="number" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.sort')" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="999" />
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
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const tableData = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const formData = reactive({
  field_name: '',
  field_key: '',
  field_type: 'text',
  sort_order: 0
})

const rules: FormRules = {
  field_name: [{ required: true, message: () => $t('systemCustomFields.enterFieldName'), trigger: 'blur' }],
  field_key: [{ required: true, message: () => $t('systemCustomFields.enterFieldKey'), trigger: 'blur' }],
  field_type: [{ required: true, message: () => $t('systemCustomFields.selectType'), trigger: 'change' }]
}

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/custom-fields' })
    tableData.value = data?.list || data || []
  } catch (error) {
    console.error('获取自定义字段失败:', error)
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  formData.field_name = ''
  formData.field_key = ''
  formData.field_type = 'text'
  formData.sort_order = 0
  editingId.value = null
  formRef.value?.resetFields()
}

const handleAdd = () => {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  isEdit.value = true
  editingId.value = row.id
  Object.assign(formData, {
    field_name: row.field_name,
    field_key: row.field_key,
    field_type: row.field_type || 'text',
    sort_order: row.sort_order || 0
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    submitting.value = true
    if (isEdit.value && editingId.value) {
      await request.put({ url: `/api/admin/custom-fields/${editingId.value}`, data: formData })
      ElMessage.success($t('common.updateSuccess'))
    } else {
      await request.post({ url: '/api/admin/custom-fields', data: formData })
      ElMessage.success($t('common.addSuccess'))
    }
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('systemCustomFields.confirmDelete'), $t('common.tips'), { type: 'warning' })
    await request.del({ url: `/api/admin/custom-fields/${row.id}` })
    ElMessage.success($t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') console.error('删除失败:', error)
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.custom-fields-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
