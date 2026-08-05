<template>
  <div class="custom-field-groups-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>自定义字段组</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加字段组
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="分组名称">
          <el-input v-model="searchForm.name" placeholder="请输入分组名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="filteredTableData" v-loading="loading" style="width: 100%" border row-key="id"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }" default-expand-all>
        <el-table-column prop="name" label="分组名称" min-width="200" />
        <el-table-column prop="code" label="分组编码" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="field_count" label="字段数量" width="100" align="center" />
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleManageFields(row)">管理字段</el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该字段组吗？删除后子分组也会被删除。" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="上级分组" prop="parent_id">
          <el-tree-select
            v-model="formData.parent_id"
            :data="groupTreeData"
            :props="{ value: 'id', label: 'name', children: 'children' }"
            placeholder="无（顶级分组）"
            clearable
            check-strictly
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="分组名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入分组名称" />
        </el-form-item>
        <el-form-item label="分组编码" prop="code">
          <el-input v-model="formData.code" placeholder="请输入分组编码" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 管理字段对话框 -->
    <el-dialog v-model="fieldDialogVisible" :title="`管理字段 - ${currentGroupName}`" width="900px" destroy-on-close>
      <div class="field-toolbar">
        <el-button type="primary" size="small" @click="handleAddField">
          <el-icon><Plus /></el-icon>添加字段
        </el-button>
      </div>
      <el-table :data="fieldList" v-loading="fieldLoading" border size="small">
        <el-table-column prop="field_name" label="字段名称" min-width="120" />
        <el-table-column prop="field_key" label="字段标识" width="120" />
        <el-table-column prop="field_type" label="字段类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag size="small">{{ fieldTypeMap[row.field_type] || row.field_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_required" label="必填" width="70" align="center">
          <template #default="{ row }">
            <el-icon v-if="row.is_required === 1" color="#67c23a"><Check /></el-icon>
            <el-icon v-else color="#c0c4cc"><Close /></el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="70" align="center" />
        <el-table-column label="操作" width="120" align="center">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEditField(row)">编辑</el-button>
            <el-popconfirm title="确定删除该字段吗？" @confirm="handleDeleteField(row)">
              <template #reference>
                <el-button type="danger" link size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <!-- 字段编辑对话框 -->
    <el-dialog v-model="fieldEditVisible" :title="fieldEditTitle" width="500px" destroy-on-close>
      <el-form :model="fieldFormData" :rules="fieldFormRules" ref="fieldFormRef" label-width="100px">
        <el-form-item label="字段名称" prop="field_name">
          <el-input v-model="fieldFormData.field_name" placeholder="请输入字段名称" />
        </el-form-item>
        <el-form-item label="字段标识" prop="field_key">
          <el-input v-model="fieldFormData.field_key" placeholder="请输入字段标识（英文）" />
        </el-form-item>
        <el-form-item label="字段类型" prop="field_type">
          <el-select v-model="fieldFormData.field_type" placeholder="请选择字段类型" style="width: 100%">
            <el-option label="文本" value="text" />
            <el-option label="数字" value="number" />
            <el-option label="日期" value="date" />
            <el-option label="单选" value="select" />
            <el-option label="多选" value="multi_select" />
            <el-option label="文本域" value="textarea" />
          </el-select>
        </el-form-item>
        <el-form-item label="是否必填" prop="is_required">
          <el-switch v-model="fieldFormData.is_required" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="fieldFormData.sort" :min="0" :max="9999" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="fieldEditVisible = false">取消</el-button>
        <el-button type="primary" @click="handleFieldSubmit" :loading="fieldSubmitLoading">确定</el-button>
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

defineOptions({ name: 'CustomFieldGroups' })

const fieldTypeMap: Record<string, string> = {
  text: '文本', number: '数字', date: '日期',
  select: '单选', multi_select: '多选', textarea: '文本域'
}

const loading = ref(false)
const submitLoading = ref(false)
const fieldLoading = ref(false)
const fieldSubmitLoading = ref(false)
const dialogVisible = ref(false)
const fieldDialogVisible = ref(false)
const fieldEditVisible = ref(false)
const dialogTitle = ref('添加字段组')
const fieldEditTitle = ref('添加字段')
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

const formRules: FormRules = {
  name: [{ required: true, message: '请输入分组名称', trigger: 'blur' }],
  code: [
    { required: true, message: '请输入分组编码', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '只能包含字母、数字和下划线', trigger: 'blur' }
  ]
}

const fieldFormRules: FormRules = {
  field_name: [{ required: true, message: '请输入字段名称', trigger: 'blur' }],
  field_key: [{ required: true, message: '请输入字段标识', trigger: 'blur' }],
  field_type: [{ required: true, message: '请选择字段类型', trigger: 'change' }]
}

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
    ElMessage.error('获取字段组列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { /* 前端过滤 */ }
const handleReset = () => { searchForm.name = '' }

const handleAdd = () => {
  dialogTitle.value = '添加字段组'
  formData.id = undefined; formData.parent_id = undefined; formData.name = ''
  formData.code = ''; formData.description = ''; formData.sort = 0; formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑字段组'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/custom-field-groups/${row.id}` })
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    ElMessage.error('删除失败')
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
      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error('操作失败')
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
    ElMessage.error('获取字段列表失败')
  } finally {
    fieldLoading.value = false
  }
}

const handleAddField = () => {
  fieldEditTitle.value = '添加字段'
  fieldFormData.id = undefined; fieldFormData.field_name = ''; fieldFormData.field_key = ''
  fieldFormData.field_type = 'text'; fieldFormData.is_required = 0; fieldFormData.sort = 0
  fieldEditVisible.value = true
}

const handleEditField = (row: any) => {
  fieldEditTitle.value = '编辑字段'
  Object.assign(fieldFormData, row)
  fieldEditVisible.value = true
}

const handleDeleteField = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/custom-field-groups/${currentGroupId.value}/fields/${row.id}` })
    ElMessage.success('删除成功')
    fetchFields()
  } catch (error) {
    ElMessage.error('删除失败')
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
      ElMessage.success(fieldFormData.id ? '更新成功' : '添加成功')
      fieldEditVisible.value = false
      fetchFields()
    } catch (error) {
      ElMessage.error('操作失败')
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
