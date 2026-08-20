<template>
  <div class="custom-fields-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsCustomFields.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('clientsCustomFields.addField') }}
          </el-button>
        </div>
      </template>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border row-key="id">
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" :label="$t('clientsCustomFields.fieldName')" width="150" />
        <el-table-column prop="field_key" :label="$t('clientsCustomFields.fieldKey')" width="150" />
        <el-table-column prop="field_type" :label="$t('clientsCustomFields.fieldType')" width="120" align="center">
          <template #default="{ row }">
            <el-tag size="small">{{ getFieldTypeText(row.field_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" :label="$t('common.description')" width="200" show-overflow-tooltip />
        <el-table-column prop="required" :label="$t('clientsCustomFields.required')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.required ? 'danger' : 'info'" size="small">{{ row.required ? $t('clientsCustomFields.requiredYes') : $t('clientsCustomFields.requiredNo') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" :label="$t('clientsCustomFields.sort')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('common.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">{{ row.status === 1 ? $t('common.enabled') : $t('common.disabled') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.action')" width="240" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="primary" link @click="handleCopy(row)">{{ $t('clientsCustomFields.copy') }}</el-button>
            <el-popconfirm :title="$t('clientsCustomFields.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('common.delete') }}</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="550px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('clientsCustomFields.fieldName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('clientsCustomFields.enterFieldName')" />
        </el-form-item>
        <el-form-item :label="$t('clientsCustomFields.fieldKey')" prop="field_key">
          <el-input v-model="formData.field_key" :placeholder="$t('clientsCustomFields.enterFieldKey')" :disabled="!!formData.id" />
        </el-form-item>
        <el-form-item :label="$t('clientsCustomFields.fieldType')" prop="field_type">
          <el-select v-model="formData.field_type" :placeholder="$t('common.select')" style="width: 100%">
            <el-option :label="$t('clientsCustomFields.typeText')" value="text" />
            <el-option :label="$t('clientsCustomFields.typeNumber')" value="number" />
            <el-option :label="$t('clientsCustomFields.typeDate')" value="date" />
            <el-option :label="$t('clientsCustomFields.typeRadio')" value="radio" />
            <el-option :label="$t('clientsCustomFields.typeCheckbox')" value="checkbox" />
            <el-option :label="$t('clientsCustomFields.typeTextarea')" value="textarea" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('clientsCustomFields.options')" prop="options" v-if="formData.field_type === 'radio' || formData.field_type === 'checkbox'">
          <el-input v-model="formData.options" type="textarea" :rows="3" :placeholder="$t('clientsCustomFields.optionsPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('common.description')">
          <el-input v-model="formData.description" :placeholder="$t('clientsCustomFields.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('clientsCustomFields.required')">
          <el-switch v-model="formData.required" />
        </el-form-item>
        <el-form-item :label="$t('clientsCustomFields.sort')">
          <el-input-number v-model="formData.sort" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()
const tableData = ref<any[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formData = reactive({ id: undefined as number | undefined, name: '', field_key: '', field_type: '', options: '', description: '', required: false, sort: 0, status: 1 })

const getFieldTypeText = (type: string) => {
  const map: Record<string, string> = { text: $t('clientsCustomFields.typeText'), number: $t('clientsCustomFields.typeNumber'), date: $t('clientsCustomFields.typeDate'), radio: $t('clientsCustomFields.typeRadio'), checkbox: $t('clientsCustomFields.typeCheckbox'), textarea: $t('clientsCustomFields.typeTextarea') }
  return map[type] || type
}

const validateFieldKey = (_rule: any, value: string, callback: any) => {
  if (!value) { callback(new Error($t('clientsCustomFields.enterFieldKey'))); return }
  if (!/^[a-zA-Z_][a-zA-Z0-9_]*$/.test(value)) { callback(new Error($t('clientsCustomFields.fieldKeyFormat'))) } else { callback() }
}

const formRules: FormRules = {
  name: [{ required: true, message: $t('clientsCustomFields.enterFieldName'), trigger: 'blur' }],
  field_key: [{ required: true, message: $t('clientsCustomFields.enterFieldKey'), trigger: 'blur' }, { validator: validateFieldKey, trigger: 'blur' }],
  field_type: [{ required: true, message: $t('common.select'), trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try { const data = await request.get({ url: '/api/admin/custom-fields' }); tableData.value = data.list || data || [] } catch (e) { ElMessage.error($t('common.fetchFailed')) } finally { loading.value = false }
}

const handleAdd = () => { dialogTitle.value = $t('clientsCustomFields.addField'); formData.id = undefined; formData.name = ''; formData.field_key = ''; formData.field_type = ''; formData.options = ''; formData.description = ''; formData.required = false; formData.sort = 0; formData.status = 1; dialogVisible.value = true }

const handleEdit = (row: any) => { dialogTitle.value = $t('common.edit'); Object.assign(formData, row); dialogVisible.value = true }

const handleCopy = (row: any) => {
  dialogTitle.value = $t('clientsCustomFields.copyField')
  Object.assign(formData, { ...row, id: undefined, field_key: row.field_key + '_copy', name: row.name + ' (Copy)' })
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try { await request.del({ url: `/api/admin/custom-fields/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchData() } catch (e) { ElMessage.error($t('common.deleteFailed')) }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/custom-fields/${formData.id}` : '/api/admin/custom-fields'
      if (formData.id) { await request.put({ url, params: formData }) } else { await request.post({ url, params: formData }) }
      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (e) { ElMessage.error($t('common.operateFailed')) } finally { submitLoading.value = false }
  })
}

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.custom-fields-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
