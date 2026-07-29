<template>
  <div class="link-cause-page page-container">
    <div class="page-header">
      <h2>关联原因</h2>
      <el-button type="primary" @click="handleCreate">新增原因</el-button>
    </div>

    <el-table :data="list" v-loading="listLoading" border fit highlight-current-row>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="原因名称" />
      <el-table-column prop="link_type" label="类型" width="120" />
      <el-table-column prop="level" label="级别" width="80" />
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 创建/编辑对话框 -->
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="500px">
      <el-form :model="formData" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="原因名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入原因名称" />
        </el-form-item>
        <el-form-item label="类型" prop="link_type">
          <el-select v-model="formData.link_type" placeholder="请选择类型">
            <el-option label="工单" value="ticket" />
            <el-option label="订单" value="order" />
            <el-option label="产品" value="product" />
          </el-select>
        </el-form-item>
        <el-form-item label="级别">
          <el-input-number v-model="formData.level" :min="1" :max="10" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const list = ref([])
const listLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref(null)

const formData = reactive({
  id: null,
  name: '',
  link_type: '',
  level: 1
})

const rules = {
  name: [{ required: true, message: '请输入原因名称', trigger: 'blur' }],
  link_type: [{ required: true, message: '请选择类型', trigger: 'change' }]
}

const fetchList = async () => {
  listLoading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/link-cause')
    list.value = data.data?.list || data.data || []
  } catch (error) {
    ElMessage.error('获取列表失败')
  } finally {
    listLoading.value = false
  }
}

const handleCreate = () => {
  dialogTitle.value = '新增原因'
  formData.id = null
  formData.name = ''
  formData.link_type = ''
  formData.level = 1
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑原因'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确认删除该原因？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await request.delete(`/admin/api/v1/link-cause/${row.id}`)
      ElMessage.success('删除成功')
      fetchList()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    if (formData.id) {
      await request.put(`/admin/api/v1/link-cause/${formData.id}`, formData)
    } else {
      await request.post('/admin/api/v1/link-cause', formData)
    }
    ElMessage.success(formData.id ? '更新成功' : '创建成功')
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('表单验证失败', error)
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.link-cause-page {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style>
