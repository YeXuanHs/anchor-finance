<template>
  <div class="page-container">
    <art-card :title="$t('clientsRemarks.title')" shadow="never">
      <template #header-extra>
        <el-input v-model="search" :placeholder="$t('clientsRemarks.searchUser')" style="width: 200px; margin-right: 10px" />
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          {{ $t('clientsRemarks.addRemark') }}
        </el-button>
      </template>

      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" :label="$t('clientsRemarks.userId')" width="100" />
        <el-table-column prop="username" :label="$t('common.username')" width="120" />
        <el-table-column prop="content" :label="$t('clientsRemarks.remarkContent')" min-width="300" show-overflow-tooltip />
        <el-table-column prop="admin" :label="$t('common.operator')" width="120" />
        <el-table-column prop="created_at" :label="$t('clientsRemarks.time')" width="180" />
        <el-table-column :label="$t('common.action')" width="100" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="handleDelete(row)">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </art-card>

    <el-dialog v-model="dialogVisible" :title="$t('clientsRemarks.addRemark')" width="500px">
      <el-form :model="formData" label-width="100px">
        <el-form-item :label="$t('clientsRemarks.userId')" required>
          <el-input v-model="formData.user_id" />
        </el-form-item>
        <el-form-item :label="$t('clientsRemarks.remarkContent')" required>
          <el-input v-model="formData.content" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const tableData = ref([])
const search = ref('')
const dialogVisible = ref(false)
const formData = ref({ user_id: '', content: '' })

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/user-remarks', params: { search: search.value } })
    tableData.value = res || []
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
    await request.post({ url: '/api/admin/user-remarks', params: formData.value })
    ElMessage.success($t('common.addSuccess'))
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm($t('clientsRemarks.confirmDelete'), $t('common.tips'), { type: 'warning' })
  try {
    await request.del({ url: `/api/admin/user-remarks/${row.id}` })
    ElMessage.success($t('common.deleteSuccess'))
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => fetchData())
</script>

<style scoped lang="scss">
.page-container { padding: 20px; }
</style>
