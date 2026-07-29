<template>
  <div class="dcim-cloud-page page-container">
    <div class="page-header">
      <h2>魔方云对接</h2>
      <el-button type="primary" @click="handleCreate">添加服务器</el-button>
    </div>

    <el-table :data="list" v-loading="listLoading" border fit highlight-current-row>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="服务器名称" />
      <el-table-column prop="hostname" label="主机名" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
            {{ row.status === 'active' ? '正常' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="250">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" type="success" @click="handleTest(row)">测试</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 创建/编辑对话框 -->
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="600px">
      <el-form :model="formData" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="服务器名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入服务器名称" />
        </el-form-item>
        <el-form-item label="主机名" prop="hostname">
          <el-input v-model="formData.hostname" placeholder="请输入主机名" />
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="formData.password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="端口">
          <el-input-number v-model="formData.port" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item label="SSL">
          <el-switch v-model="formData.secure" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="formData.disabled" active-text="禁用" inactive-text="启用" />
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
  hostname: '',
  username: '',
  password: '',
  port: 443,
  secure: true,
  disabled: false
})

const rules = {
  name: [{ required: true, message: '请输入服务器名称', trigger: 'blur' }],
  hostname: [{ required: true, message: '请输入主机名', trigger: 'blur' }],
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const fetchList = async () => {
  listLoading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/dcim-cloud')
    list.value = data.data?.list || data.data || []
  } catch (error) {
    ElMessage.error('获取列表失败')
  } finally {
    listLoading.value = false
  }
}

const handleCreate = () => {
  dialogTitle.value = '添加服务器'
  formData.id = null
  formData.name = ''
  formData.hostname = ''
  formData.username = ''
  formData.password = ''
  formData.port = 443
  formData.secure = true
  formData.disabled = false
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑服务器'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleTest = async (row) => {
  try {
    await request.post(`/admin/api/v1/dcim-cloud/${row.id}/test`)
    ElMessage.success('连接测试成功')
  } catch (error) {
    ElMessage.error('连接测试失败')
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确认删除该服务器？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await request.delete(`/admin/api/v1/dcim-cloud/${row.id}`)
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
      await request.put(`/admin/api/v1/dcim-cloud/${formData.id}`, formData)
    } else {
      await request.post('/admin/api/v1/dcim-cloud', formData)
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
.dcim-cloud-page {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style>
