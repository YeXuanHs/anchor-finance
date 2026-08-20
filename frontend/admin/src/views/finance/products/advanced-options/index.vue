<template>
  <div class="page-container">
    <art-card :title="$t('advancedOptions.title')" shadow="never">
      <template #header-extra>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          {{ $t('advancedOptions.addOption') }}
        </el-button>
      </template>

      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="$t('advancedOptions.optionName')" min-width="150" />
        <el-table-column prop="key" :label="$t('advancedOptions.configKey')" min-width="150" />
        <el-table-column prop="type" :label="$t('common.type')" width="100">
          <template #default="{ row }">
            <el-tag>{{ row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="value" :label="$t('advancedOptions.configValue')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" :label="$t('common.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? $t('common.enable') : $t('common.disable') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.action')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </art-card>

    <!-- 编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" label-width="120px">
        <el-form-item :label="$t('advancedOptions.optionName')" required>
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item :label="$t('advancedOptions.configKey')" required>
          <el-input v-model="formData.key" />
        </el-form-item>
        <el-form-item :label="$t('common.type')" required>
          <el-select v-model="formData.type">
            <el-option :label="$t('advancedOptions.typeText')" value="text" />
            <el-option :label="$t('advancedOptions.typeNumber')" value="number" />
            <el-option :label="$t('advancedOptions.typeBoolean')" value="boolean" />
            <el-option :label="$t('advancedOptions.typeSelect')" value="select" />
            <el-option :label="$t('advancedOptions.typeTextarea')" value="textarea" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('advancedOptions.configValue')">
          <el-input v-model="formData.value" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
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
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formData = ref({
  id: null,
  name: '',
  key: '',
  type: 'text',
  value: '',
  status: 1
})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/advanced-options' })
    tableData.value = res || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  dialogTitle.value = $t('advancedOptions.addOption')
  formData.value = { id: null, name: '', key: '', type: 'text', value: '', status: 1 }
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('advancedOptions.editOption')
  formData.value = { ...row }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    if (formData.value.id) {
      await request.put({ url: `/api/admin/advanced-options/${formData.value.id}`, params: formData.value })
    } else {
      await request.post({ url: '/api/admin/advanced-options', params: formData.value })
    }
    ElMessage.success($t('common.operateSuccess'))
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm($t('advancedOptions.confirmDelete'), $t('common.tips'), { type: 'warning' })
  try {
    await request.del({ url: `/api/admin/advanced-options/${row.id}` })
    ElMessage.success($t('common.deleteSuccess'))
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
