<template>
  <div class="page-container">
    <art-card title="用户备注" shadow="never">
      <template #header-extra>
        <el-input v-model="search" placeholder="搜索用户" style="width: 200px; margin-right: 10px" />
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          添加备注
        </el-button>
      </template>

      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="content" label="备注内容" min-width="300" show-overflow-tooltip />
        <el-table-column prop="admin" label="操作人" width="120" />
        <el-table-column prop="created_at" label="时间" width="180" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </art-card>

    <el-dialog v-model="dialogVisible" title="添加备注" width="500px">
      <el-form :model="formData" label-width="100px">
        <el-form-item label="用户ID" required>
          <el-input v-model="formData.user_id" />
        </el-form-item>
        <el-form-item label="备注内容" required>
          <el-input v-model="formData.content" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const tableData = ref([])
const search = ref('')
const dialogVisible = ref(false)
const formData = ref({ user_id: '', content: '' })

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/admin/user-remarks', { params: { search: search.value } })
    tableData.value = data?.data || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  formData.value = { user_id: '', content: '' }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    await request.post('/api/admin/user-remarks', formData.value)
    ElMessage.success('添加成功')
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm('确定删除该备注？', '提示', { type: 'warning' })
  try {
    await request.delete(`/api/admin/user-remarks/${row.id}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => fetchData())
</script>

<style scoped lang="scss">
.page-container {
  padding: 20px;
}
</style>
