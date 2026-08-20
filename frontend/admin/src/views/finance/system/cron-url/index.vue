<template>
  <div class="page-container">
    <art-card :title="$t('cronUrl.title')" shadow="never">
      <template #header-extra>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          {{ $t('cronUrl.addTask') }}
        </el-button>
      </template>

      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="$t('cronUrl.taskName')" min-width="150" />
        <el-table-column prop="url" label="URL" min-width="250" show-overflow-tooltip />
        <el-table-column prop="cron" :label="$t('cronUrl.cronExpression')" width="150" />
        <el-table-column prop="status" :label="$t('cronUrl.status')" width="100">
          <template #default="{ row }">
            <el-switch v-model="row.status" :active-value="1" :inactive-value="0" @change="handleStatusChange(row)" />
          </template>
        </el-table-column>
        <el-table-column prop="last_run" :label="$t('cronUrl.lastRun')" width="180" />
        <el-table-column prop="next_run" :label="$t('cronUrl.nextRun')" width="180" />
        <el-table-column :label="$t('cronUrl.actions')" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleRun(row)">{{ $t('cronUrl.runNow') }}</el-button>
            <el-button size="small" @click="handleEdit(row)">{{ $t('cronUrl.edit') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">{{ $t('cronUrl.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </art-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" label-width="120px">
        <el-form-item :label="$t('cronUrl.taskName')" required>
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item label="URL" required>
          <el-input v-model="formData.url" />
        </el-form-item>
        <el-form-item :label="$t('cronUrl.cronExpression')" required>
          <el-input v-model="formData.cron" :placeholder="$t('cronUrl.cronPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('cronUrl.requestMethod')">
          <el-select v-model="formData.method">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('cronUrl.status')">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('cronUrl.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ $t('cronUrl.confirm') }}</el-button>
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
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formData = ref({
  id: null,
  name: '',
  url: '',
  cron: '',
  method: 'GET',
  status: 1
})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/cron-url' })
    tableData.value = res || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  dialogTitle.value = $t('cronUrl.addUrlTask')
  formData.value = { id: null, name: '', url: '', cron: '', method: 'GET', status: 1 }
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('cronUrl.editUrlTask')
  formData.value = { ...row }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    if (formData.value.id) {
      await request.put({ url: `/api/admin/cron-url/${formData.value.id}`, params: formData.value })
    } else {
      await request.post({ url: '/api/admin/cron-url', params: formData.value })
    }
    ElMessage.success($t('cronUrl.operationSuccess'))
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleStatusChange = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/cron-url/${row.id}/status`, params: { status: row.status } })
    ElMessage.success($t('cronUrl.statusUpdateSuccess'))
  } catch (error) {
    row.status = row.status === 1 ? 0 : 1
    console.error(error)
  }
}

const handleRun = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/cron-url/${row.id}/run` })
    ElMessage.success($t('cronUrl.taskTriggered'))
  } catch (error) {
    console.error(error)
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm($t('cronUrl.confirmDelete'), $t('cronUrl.tip'), { type: 'warning' })
  try {
    await request.del({ url: `/api/admin/cron-url/${row.id}` })
    ElMessage.success($t('cronUrl.deleteSuccess'))
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
