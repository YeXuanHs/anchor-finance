<template>
  <div class="news-page">
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('contentPage.news.addArticle') }}</el-button>
      </div>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="title" :label="$t('contentPage.news.title')" min-width="200" />
        <el-table-column prop="category" :label="$t('contentPage.news.category')" width="100" />
        <el-table-column prop="author" :label="$t('contentPage.news.author')" width="100" />
        <el-table-column prop="views" :label="$t('contentPage.news.views')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('contentPage.news.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">
              {{ row.status === 'published' ? $t('contentPage.news.published') : $t('contentPage.news.draft') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('contentPage.news.publishTime')" width="170" />
        <el-table-column :label="$t('contentPage.news.action')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">{{ $t('contentPage.news.edit') }}</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">{{ $t('contentPage.news.delete') }}</el-button>
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
    const data = await request.get({ url: '/api/admin/news' })
    tableData.value = data || []
  } catch (error) {
    console.error($t('contentPage.news.fetchFailed') + ':', error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => { console.log('添加文章') }
const handleEdit = (row: any) => { console.log('编辑:', row.id) }
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('contentPage.news.confirmDeleteTitle', { title: row.title }), $t('contentPage.news.confirmDeleteTip'), { type: 'warning' })
    await request.del({ url: `/api/admin/news/${row.id}` })
    ElMessage.success($t('contentPage.news.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') console.error($t('contentPage.news.deleteFailed') + ':', error)
  }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.news-page { padding: 16px; }
.action-card { margin-bottom: 16px; }
.action-bar { display: flex; justify-content: space-between; align-items: center; }
.table-card { :deep(.el-card__body) { padding: 0; } }
</style>
