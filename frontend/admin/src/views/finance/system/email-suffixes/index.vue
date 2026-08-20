<template>
  <div class="email-suffixes-page">
    <div class="page-header">
      <div class="header-left">
        <h2>{{ $t('systemEmailSuffixes.title') }}</h2>
        <span class="subtitle">{{ $t('systemEmailSuffixes.subtitle') }}</span>
      </div>
      <div class="header-actions">
        <el-button @click="handleImportDefaults" :loading="importing">{{ $t('systemEmailSuffixes.importDefaults') }}</el-button>
        <el-button type="primary" @click="showAddDialog">{{ $t('systemEmailSuffixes.addSuffix') }}</el-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="filter-bar">
      <el-input v-model="searchKeyword" :placeholder="$t('systemEmailSuffixes.searchPlaceholder')" clearable style="width: 300px" @clear="loadData">
        <template #append>
          <el-button @click="loadData">{{ $t('common.search') }}</el-button>
        </template>
      </el-input>
      <el-checkbox v-model="showInactive" @change="loadData">{{ $t('systemEmailSuffixes.showInactive') }}</el-checkbox>
    </div>

    <!-- 数据表格 -->
    <el-table :data="suffixList" v-loading="loading" stripe>
      <el-table-column prop="suffix" :label="$t('systemEmailSuffixes.suffix')" width="200">
        <template #default="{ row }">
          <span class="suffix-text">@{{ row.suffix }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="remark" :label="$t('common.remark')" min-width="200" />
      <el-table-column prop="is_default" :label="$t('systemEmailSuffixes.isDefault')" width="80" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.is_default" type="success" size="small">{{ $t('systemEmailSuffixes.defaultTag') }}</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="is_active" :label="$t('common.status')" width="100" align="center">
        <template #default="{ row }">
          <el-switch v-model="row.is_active" @change="handleToggleActive(row)" />
        </template>
      </el-table-column>
      <el-table-column :label="$t('common.action')" width="120" align="center">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="showEditDialog(row)">{{ $t('common.edit') }}</el-button>
          <el-popconfirm :title="$t('systemEmailSuffixes.confirmDeleteSuffix')" @confirm="handleDelete(row.id)">
            <template #reference>
              <el-button type="danger" link size="small">{{ $t('common.delete') }}</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? $t('systemEmailSuffixes.editSuffix') : $t('systemEmailSuffixes.addSuffix')" width="500px">
      <el-form :model="formData" label-width="80px">
        <el-form-item :label="$t('systemEmailSuffixes.suffix')" required>
          <el-input v-model="formData.suffix" :placeholder="$t('systemEmailSuffixes.suffixPlaceholder')" :disabled="isEdit">
            <template #prepend>@</template>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('common.remark')">
          <el-input v-model="formData.remark" :placeholder="$t('systemEmailSuffixes.remarkPlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

interface EmailSuffix {
  id: number
  suffix: string
  is_default: boolean
  is_active: boolean
  remark: string
}

const loading = ref(false)
const importing = ref(false)
const submitting = ref(false)
const suffixList = ref<EmailSuffix[]>([])
const searchKeyword = ref('')
const showInactive = ref(true)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)

const formData = ref({
  suffix: '',
  remark: ''
})

const loadData = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/email-suffixes', params: { show_inactive: showInactive.value } })
    suffixList.value = res || []
  } catch (e) {
    ElMessage.error($t('systemEmailSuffixes.loadFailed'))
  } finally {
    loading.value = false
  }
}

const showAddDialog = () => {
  isEdit.value = false
  formData.value = { suffix: '', remark: '' }
  dialogVisible.value = true
}

const showEditDialog = (row: any) => {
  isEdit.value = true
  editId.value = row.id
  formData.value = { suffix: row.suffix, remark: row.remark }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formData.value.suffix) {
    ElMessage.warning($t('systemEmailSuffixes.enterSuffix'))
    return
  }
  submitting.value = true
  try {
    if (isEdit.value) {
      await request.put({ url: `/api/admin/email-suffixes/${editId.value}`, params: { remark: formData.value.remark } })
    } else {
      await request.post({ url: '/api/admin/email-suffixes', params: formData.value })
    }
    ElMessage.success(isEdit.value ? $t('common.updateSuccess') : $t('common.addSuccess'))
    dialogVisible.value = false
    loadData()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.message || $t('common.operateFailed'))
  } finally {
    submitting.value = false
  }
}

const handleToggleActive = async (row: any) => {
  try {
    await request.put({ url: `/api/admin/email-suffixes/${row.id}`, params: { is_active: row.is_active } })
  } catch (e) {
    row.is_active = !row.is_active
    ElMessage.error($t('common.updateFailed'))
  }
}

const handleDelete = async (id: number) => {
  try {
    await request.del({ url: `/api/admin/email-suffixes/${id}` })
    ElMessage.success($t('common.deleteSuccess'))
    loadData()
  } catch (e) {
    ElMessage.error($t('common.deleteFailed'))
  }
}

const handleImportDefaults = async () => {
  importing.value = true
  try {
    await request.post({ url: '/api/admin/email-suffixes/import-defaults' })
    ElMessage.success($t('systemEmailSuffixes.importSuccess'))
    loadData()
  } catch (e) {
    ElMessage.error($t('systemEmailSuffixes.importFailed'))
  } finally {
    importing.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.email-suffixes-page {
  padding: 20px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.header-left h2 {
  margin: 0 0 4px 0;
  font-size: 20px;
}
.subtitle {
  color: #909399;
  font-size: 14px;
}
.filter-bar {
  display: flex;
  gap: 16px;
  align-items: center;
  margin-bottom: 16px;
}
.suffix-text {
  font-family: monospace;
  font-weight: 600;
  color: #409eff;
}
</style>
