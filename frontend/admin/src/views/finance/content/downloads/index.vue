<template>
  <div class="downloads-page">
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>添加下载</el-button>
      </div>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="title" label="标题" min-width="200" />
        <el-table-column prop="category" label="分类" width="100" />
        <el-table-column prop="file_size" label="大小" width="100" />
        <el-table-column prop="downloads" label="下载次数" width="100" align="center" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/http'

const loading = ref(false)
const tableData = ref([])

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/downloads' })
    tableData.value = data || []
  } catch (error) {
    console.error('获取下载列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => { console.log('添加下载') }
const handleEdit = (row: any) => { console.log('编辑:', row.id) }
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除 "${row.title}" 吗？`, '确认删除', { type: 'warning' })
    await request.del({ url: `/api/admin/downloads/${row.id}` })
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') console.error('删除失败:', error)
  }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.downloads-page { padding: 16px; }
.action-card { margin-bottom: 16px; }
.action-bar { display: flex; justify-content: space-between; align-items: center; }
.table-card { :deep(.el-card__body) { padding: 0; } }
</style>
