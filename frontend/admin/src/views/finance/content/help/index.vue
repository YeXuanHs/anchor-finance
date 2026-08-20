<template>
  <div class="help-center-page">
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('help.addHelp') }}</el-button>
      </div>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="title" :label="$t('help.title')" min-width="200" />
        <el-table-column prop="category" :label="$t('help.category')" width="120" />
        <el-table-column prop="views" :label="$t('help.views')" width="80" align="center" />
        <el-table-column prop="sort_order" :label="$t('help.sortOrder')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('help.status')" width="80" align="center">
          <template #default="{ row }"><el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">{{ row.status === 'published' ? $t('help.published') : $t('help.draft') }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="170" />
        <el-table-column :label="$t('help.operations')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">{{ $t('common.delete') }}</el-button>
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
  try { const data = await request.get({ url: '/api/admin/help' }); tableData.value = data || [] } catch (error) { console.error('fetch help failed:', error) } finally { loading.value = false }
}

const handleAdd = () => { console.log('add help') }
const handleEdit = (row: any) => { console.log('edit:', row.id) }
const handleDelete = async (row: any) => {
  try { await ElMessageBox.confirm($t('help.confirmDelete', { title: row.title }), $t('common.tips'), { type: 'warning' }); await request.del({ url: `/api/admin/help/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchList() } catch (error) { if (error !== 'cancel') console.error('delete failed:', error) }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.help-center-page { padding: 16px; }
.action-card { margin-bottom: 16px; }
</style>
