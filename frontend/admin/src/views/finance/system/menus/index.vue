<template>
  <div class="page-container">
    <art-card title="菜单管理" shadow="never">
      <template #header-extra>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          添加菜单
        </el-button>
      </template>

      <el-table :data="tableData" v-loading="loading" row-key="id" :tree-props="{ children: 'children' }">
        <el-table-column prop="name" label="菜单名称" min-width="200" />
        <el-table-column prop="url" label="链接" min-width="200" show-overflow-tooltip />
        <el-table-column prop="icon" label="图标" width="120" />
        <el-table-column prop="order" label="排序" width="100" />
        <el-table-column prop="menu_type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag>{{ getMenuTypeText(row.menu_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '显示' : '隐藏' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleAddChild(row)">添加子菜单</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </art-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" label-width="120px">
        <el-form-item label="菜单名称" required>
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item label="链接">
          <el-input v-model="formData.url" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="formData.icon" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="formData.order" :min="0" />
        </el-form-item>
        <el-form-item label="菜单类型" required>
          <el-select v-model="formData.menu_type">
            <el-option label="用户中心" :value="1" />
            <el-option label="www头部" :value="2" />
            <el-option label="www尾部" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
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
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formData = ref({
  id: null,
  name: '',
  url: '',
  icon: '',
  order: 0,
  parent_id: 0,
  menu_type: 1,
  status: 1
})

const getMenuTypeText = (type: number) => {
  const map: Record<number, string> = { 1: '用户中心', 2: 'www头部', 3: 'www尾部' }
  return map[type] || '未知'
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/menus/tree' })
    tableData.value = res || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  dialogTitle.value = '添加菜单'
  formData.value = { id: null, name: '', url: '', icon: '', order: 0, parent_id: 0, menu_type: 1, status: 1 }
  dialogVisible.value = true
}

const handleAddChild = (row: any) => {
  dialogTitle.value = '添加子菜单'
  formData.value = { id: null, name: '', url: '', icon: '', order: 0, parent_id: row.id, menu_type: row.menu_type, status: 1 }
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑菜单'
  formData.value = { ...row }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    if (formData.value.id) {
      await request.put({ url: `/api/admin/menus/${formData.value.id}`, params: formData.value })
    } else {
      await request.post({ url: '/api/admin/menus', params: formData.value })
    }
    ElMessage.success('操作成功')
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm('确定删除该菜单及其子菜单？', '提示', { type: 'warning' })
  try {
    await request.del({ url: `/api/admin/menus/${row.id}` })
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
