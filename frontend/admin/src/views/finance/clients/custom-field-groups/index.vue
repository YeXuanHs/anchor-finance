<template>
  <div class="custom-field-groups-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clients.customFieldGroups.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('clients.customFieldGroups.addGroup') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('clients.customFieldGroups.groupName')">
          <el-input v-model="searchForm.name" :placeholder="$t('clients.customFieldGroups.groupNamePlaceholder')" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="filteredTableData" v-loading="loading" style="width: 100%" border row-key="id"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }" default-expand-all>
        <el-table-column prop="name" :label="$t('clients.customFieldGroups.groupName')" min-width="200" />
        <el-table-column prop="code" :label="$t('clients.customFieldGroups.groupCode')" width="150" />
        <el-table-column prop="description" :label="$t('clients.customFieldGroups.description')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="field_count" :label="$t('clients.customFieldGroups.fieldCount')" width="100" align="center" />
        <el-table-column prop="sort" :label="$t('clients.customFieldGroups.sort')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('clients.customFieldGroups.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? $t('common.enable') : $t('common.disable') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('clients.customFieldGroups.operations')" width="280" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleManageFields(row)">{{ $t('clients.customFieldGroups.manageFields') }}</el-button>
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-popconfirm :title="$t('clients.customFieldGroups.deleteGroupConfirm')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('common.delete') }}</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('clients.customFieldGroups.parentGroup')" prop="parent_id">
          <el-tree-select
            v-model="formData.parent_id"
            :data="groupTreeData"
            :props="{ value: 'id', label: 'name', children: 'children' } as any"
            :placeholder="$t('clients.customFieldGroups.topGroup')"
            clearable
            check-strictly
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="$t('clients.customFieldGroups.groupName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('clients.customFieldGroups.groupNamePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('clients.customFieldGroups.groupCode')" prop="code">
          <el-input v-model="formData.code" :placeholder="$t('clients.customFieldGroups.groupCodePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('clients.customFieldGroups.description')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('clients.customFieldGroups.descriptionPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('clients.customFieldGroups.sort')" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('clients.customFieldGroups.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" :active-text="$t('common.enable')" :inactive-text="$t('common.disable')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 管理字段对话框 -->
    <el-dialog v-model="fieldDialogVisible" :title="`${$t('clients.customFieldGroups.manageFields')} - ${currentGroupName}`" width="900px" destroy-on-close>
      <div class="field-toolbar">
        <el-button type="primary" size="small" @click="handleAddField">
          <el-icon><Plus /></el-icon>{{ $t('clients.customFieldGroups.addField') }}
        </el-button>
      </div>
      <el-table :data="fieldList" v-loading="fieldLoading" border size="small">
        <el-table-column prop="field_name" :label="$t('clients.customFieldGroups.fieldName')" min-width="120" />
        <el-table-column prop="field_key" :label="$t('clients.customFieldGroups.fieldKey')" width="120" />
        <el-table-column prop="field_type" :label="$t('clients.customFieldGroups.fieldType')" width="100" align="center">
          <template #default="{ row }">
            <el-tag size="small">{{ fieldTypeMap[row.field_type] || row.field_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_required" :label="$t('clients.customFieldGroups.required')" width="70" align="center">
          <template #default="{ row }">
            <el-icon v-if="row.is_required === 1" color="#67c23a"><Check /></el-icon>
            <el-icon v-else color="#c0c4cc"><Close /></el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="sort" :label="$t('clients.customFieldGroups.sort')" width="70" align="center" />
        <el-table-column :label="$t('clients.customFieldGroups.operations')" width="120" align="center">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEditField(row)">{{ $t('common.edit') }}</el-button>
            <el-popconfirm :title="$t('clients.customFieldGroups.deleteFieldConfirm')" @confirm="handleDeleteField(row)">
              <template #reference>
                <el-button type="danger" link size="small">{{ $t('common.delete') }}</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <!-- 字段编辑对话框 -->
    <el-dialog v-model="fieldEditVisible" :title="fieldEditTitle" width="500px" destroy-on-close>
      <el-form :model="fieldFormData" :rules="fieldFormRules" ref="fieldFormRef" label-width="100px">
        <el-form-item :label="$t('clients.customFieldGroups.fieldName')" prop="field_name">
          <el-input v-model="fieldFormData.field_name" :placeholder="$t('clients.customFieldGroups.fieldNamePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('clients.customFieldGroups.fieldKey')" prop="field_key">
          <el-input v-model="fieldFormData.field_key" :placeholder="$t('clients.customFieldGroups.fieldKeyPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('clients.customFieldGroups.fieldType')" prop="field_type">
          <el-select v-model="fieldFormData.field_type" :placeholder="$t('clients.customFieldGroups.fieldTypePlaceholder')" style="width: 100%">
            <el-option :label="$t('clients.customFieldGroups.fieldTypes.text')" value="text" />
            <el-option :label="$t('clients.customFieldGroups.fieldTypes.number')" value="number" />
            <el-option :label="$t('clients.customFieldGroups.fieldTypes.date')" value="date" />
            <el-option :label="$t('clients.customFieldGroups.fieldTypes.select')" value="select" />
            <el-option :label="$t('clients.customFieldGroups.fieldTypes.multiSelect')" value="multi_select" />
            <el-option :label="$t('clients.customFieldGroups.fieldTypes.textarea')" value="textarea" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('clients.customFieldGroups.required')" prop="is_required">
          <el-switch v-model="fieldFormData.is_required" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item :label="$t('clients.customFieldGroups.sort')" prop="sort">
          <el-input-number v-model="fieldFormData.sort" :min="0" :max="9999" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="fieldEditVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleFieldSubmit" :loading="fieldSubmitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Plus, Check, Close } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'CustomFieldGroups' })

const fieldTypeMap = computed<Record<string, string>>(() => ({
  text: $t('clients.customFieldGroups.fieldTypes.text'),
  number: $t('clients.customFieldGroups.fieldTypes.number'),
  date: $t('clients.customFieldGroups.fieldTypes.date'),
  select: $t('clients.customFieldGroups.fieldTypes.select'),
  multi_select: $t('clients.customFieldGroups.fieldTypes.multiSelect'),
  textarea: $t('clients.customFieldGroups.fieldTypes.textarea')
}))

const loading = ref(false)
const submitLoading = ref(false)
const fieldLoading = ref(false)
const fieldSubmitLoading = ref(false)
const dialogVisible = ref(false)
const fieldDialogVisible = ref(false)
const fieldEditVisible = ref(false)
const dialogTitle = ref($t('clients.customFieldGroups.addGroup'))
const fieldEditTitle = ref($t('clients.customFieldGroups.addField'))
const formRef = ref<FormInstance>()
const fieldFormRef = ref<FormInstance>()
const currentGroupId = ref<number>(0)
const currentGroupName = ref('')

const searchForm = reactive({ name: '' })
const tableData = ref<any[]>([])
const fieldList = ref<any[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  parent_id: undefined as number | undefined,
  name: '',
  code: '',
  description: '',
  sort: 0,
  status: 1
})

const fieldFormData = reactive({
  id: undefined as number | undefined,
  field_name: '',
  field_key: '',
  field_type: 'text',
  is_required: 0,
  sort: 0
})

const formRules = computed<FormRules>(() => ({
  name: [{ required: true, message: $t('clients.customFieldGroups.groupNameRequired'), trigger: 'blur' }],
  code: [
    { required: true, message: $t('clients.customFieldGroups.groupCodeRequired'), trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: $t('clients.customFieldGroups.groupCodePattern'), trigger: 'blur' }
  ]
}))

const fieldFormRules = computed<FormRules>(() => ({
  field_name: [{ required: true, message: $t('clients.customFieldGroups.fieldNameRequired'), trigger: 'blur' }],
  field_key: [{ required: true, message: $t('clients.customFieldGroups.fieldKeyRequired'), trigger: 'blur' }],
  field_type: [{ required: true, message: $t('clients.customFieldGroups.fieldTypeRequired'), trigger: 'change' }]
}))

const filteredTableData = computed(() => {
  if (!searchForm.name) return tableData.value
  return tableData.value.filter(item =>
    item.name.toLowerCase().includes(searchForm.name.toLowerCase())
  )
})

const groupTreeData = computed(() => {
  return buildTreeSelectData(tableData.value, formData.id)
})

const buildTreeSelectData = (data: any[], excludeId?: number): any[] => {
  return data
    .filter(item => item.id !== excludeId)
    .map(item => ({
      id: item.id,
      name: item.name,
      children: item.children ? buildTreeSelectData(item.children, excludeId) : []
    }))
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/custom-field-groups' })
    tableData.value = data.list || data || []
  } catch (error) {
    ElMessage.error($t('clients.customFieldGroups.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { /* 前端过滤 */ }
const handleReset = () => { searchForm.name = '' }

const handleAdd = () => {
  dialogTitle.value = $t('clients.customFieldGroups.addGroup')
  formData.id = undefined; formData.parent_id = undefined; formData.name = ''
  formData.code = ''; formData.description = ''; formData.sort = 0; formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('clients.customFieldGroups.editGroup')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/custom-field-groups/${row.id}` })
    ElMessage.success($t('common.deleteSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('common.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/custom-field-groups/${formData.id}`, params: formData })
      } else {
        await request.post({ url: '/api/admin/custom-field-groups', params: formData })
      }
      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('common.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleManageFields = async (row: any) => {
  currentGroupId.value = row.id
  currentGroupName.value = row.name
  fieldDialogVisible.value = true
  await fetchFields()
}

const fetchFields = async () => {
  fieldLoading.value = true
  try {
    const data = await request.get({ url: `/api/admin/custom-field-groups/${currentGroupId.value}/fields` })
    fieldList.value = data.list || data || []
  } catch (error) {
    ElMessage.error($t('clients.customFieldGroups.fetchFieldsFailed'))
  } finally {
    fieldLoading.value = false
  }
}

const handleAddField = () => {
  fieldEditTitle.value = $t('clients.customFieldGroups.addField')
  fieldFormData.id = undefined; fieldFormData.field_name = ''; fieldFormData.field_key = ''
  fieldFormData.field_type = 'text'; fieldFormData.is_required = 0; fieldFormData.sort = 0
  fieldEditVisible.value = true
}

const handleEditField = (row: any) => {
  fieldEditTitle.value = $t('clients.customFieldGroups.editField')
  Object.assign(fieldFormData, row)
  fieldEditVisible.value = true
}

const handleDeleteField = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/custom-field-groups/${currentGroupId.value}/fields/${row.id}` })
    ElMessage.success($t('common.deleteSuccess'))
    fetchFields()
  } catch (error) {
    ElMessage.error($t('common.deleteFailed'))
  }
}

const handleFieldSubmit = async () => {
  if (!fieldFormRef.value) return
  await fieldFormRef.value.validate(async (valid) => {
    if (!valid) return
    fieldSubmitLoading.value = true
    try {
      if (fieldFormData.id) {
        await request.put({
          url: `/api/admin/custom-field-groups/${currentGroupId.value}/fields/${fieldFormData.id}`,
          params: fieldFormData
        })
      } else {
        await request.post({
          url: `/api/admin/custom-field-groups/${currentGroupId.value}/fields`,
          params: fieldFormData
        })
      }
      ElMessage.success(fieldFormData.id ? $t('common.updateSuccess') : $t('common.addSuccess'))
      fieldEditVisible.value = false
      fetchFields()
    } catch (error) {
      ElMessage.error($t('common.operationFailed'))
    } finally {
      fieldSubmitLoading.value = false
    }
  })
}

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.custom-field-groups-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.field-toolbar { margin-bottom: 16px; }
</style>
