<template>
  <div class="custom-fields-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>客户自定义字段配置</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加字段
          </el-button>
        </div>
      </template>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="field_name" label="字段名称" width="150" />
        <el-table-column prop="field_key" label="字段标识" width="150" />
        <el-table-column prop="field_type_text" label="字段类型" width="120" align="center" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="is_required" label="是否必填" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_required === 1 ? 'danger' : 'info'" size="small">
              {{ row.is_required === 1 ? '必填' : '选填' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handleCopy(row)">复制</el-button>
            <el-popconfirm title="确定删除该字段吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="字段名称" prop="field_name">
          <el-input v-model="formData.field_name" placeholder="请输入字段名称" />
        </el-form-item>
        <el-form-item label="字段标识" prop="field_key">
          <el-input v-model="formData.field_key" placeholder="请输入字段标识（英文）" />
        </el-form-item>
        <el-form-item label="字段类型" prop="field_type">
          <el-select v-model="formData.field_type" placeholder="请选择字段类型">
            <el-option label="文本" value="text" />
            <el-option label="数字" value="number" />
            <el-option label="日期" value="date" />
            <el-option label="单选" value="select" />
            <el-option label="多选" value="multi_select" />
            <el-option label="文本域" value="textarea" />
          </el-select>
        </el-form-item>
        <el-form-item label="选项值" prop="options" v-if="['select', 'multi_select'].includes(formData.field_type)">
          <el-input v-model="formData.options" type="textarea" :rows="3" placeholder="每行一个选项" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="2" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="是否必填" prop="is_required">
          <el-switch v-model="formData.is_required" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
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

const loading = ref(false)
const submitLoading = ref(false)

const tableData = ref<any[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref('添加字段')
const formRef = ref<FormInstance>()

const formData = reactive({
  id: undefined as number | undefined,
  field_name: '',
  field_key: '',
  field_type: 'text',
  options: '',
  description: '',
  is_required: 0,
  sort: 0,
  status: 1
})

const formRules: FormRules = {
  field_name: [{ required: true, message: '请输入字段名称', trigger: 'blur' }],
  field_key: [{ required: true, message: '请输入字段标识', trigger: 'blur' }],
  field_type: [{ required: true, message: '请选择字段类型', trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/custom-fields' })
    tableData.value = data.list || data || []
  } catch (error) {
    console.error('获取数据失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogTitle.value = '添加字段'
  formData.id = undefined
  formData.field_name = ''
  formData.field_key = ''
  formData.field_type = 'text'
  formData.options = ''
  formData.description = ''
  formData.is_required = 0
  formData.sort = 0
  formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑字段'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleCopy = (row: any) => {
  dialogTitle.value = '复制字段'
  Object.assign(formData, row)
  formData.id = undefined
  formData.field_name = row.field_name + '_副本'
  formData.field_key = row.field_key + '_copy'
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/custom-fields/${row.id}` })
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
      const url = formData.id ? `/api/admin/custom-fields/${formData.id}` : '/api/admin/custom-fields'
      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
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

onMounted(() => { fetchData() })
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
