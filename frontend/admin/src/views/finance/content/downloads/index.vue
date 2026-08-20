<template>
  <div class="downloads-page">
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('contentPage.download.addDownload') }}</el-button>
      </div>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="title" :label="$t('contentPage.download.title')" min-width="200" />
        <el-table-column prop="category" :label="$t('contentPage.download.category')" width="100" />
        <el-table-column prop="file_size" :label="$t('contentPage.download.fileSize')" width="100" />
        <el-table-column prop="downloads" :label="$t('contentPage.download.downloadCount')" width="100" align="center" />
        <el-table-column prop="status" :label="$t('contentPage.download.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? $t('contentPage.download.enabled') : $t('contentPage.download.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('contentPage.download.createdAt')" width="170" />
        <el-table-column :label="$t('contentPage.download.action')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">{{ $t('contentPage.download.edit') }}</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">{{ $t('contentPage.download.delete') }}</el-button>
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
import { $t } from '@/locales'

const loading = ref(false)
const tableData = ref([])

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/downloads' })
    tableData.value = data || []
  } catch (error) {
    console.error($t('contentPage.download.fetchFailed') + ':', error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => { console.log('添加下载') }
const handleEdit = (row: any) => { console.log('编辑:', row.id) }
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('contentPage.download.confirmDeleteTitle', { title: row.title }), $t('contentPage.download.confirmDeleteTip'), { type: 'warning' })
    await request.del({ url: `/api/admin/downloads/${row.id}` })
    ElMessage.success($t('contentPage.download.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') console.error($t('contentPage.download.deleteFailed') + ':', error)
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
