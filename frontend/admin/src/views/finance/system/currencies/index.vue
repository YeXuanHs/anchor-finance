<template>
  <div class="currencies-page">
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('currency.addCurrency') }}</el-button>
      </div>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="code" :label="$t('currency.code')" width="80" />
        <el-table-column prop="name" :label="$t('currency.name')" min-width="120" />
        <el-table-column prop="symbol" :label="$t('currency.symbol')" width="80" align="center" />
        <el-table-column prop="exchange_rate" :label="$t('currency.exchangeRate')" width="100" align="center" />
        <el-table-column prop="is_default" :label="$t('currency.isDefault')" width="80" align="center">
          <template #default="{ row }"><el-tag :type="row.is_default ? 'success' : 'info'" size="small">{{ row.is_default ? $t('common.yes') : $t('common.no') }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('currency.status')" width="80" align="center">
          <template #default="{ row }"><el-switch v-model="row.status" :active-value="'active'" :inactive-value="'disabled'" @change="handleToggleStatus(row)" /></template>
        </el-table-column>
        <el-table-column :label="$t('currency.operations')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="warning" link size="small" @click="handleSetDefault(row)" :disabled="row.is_default">{{ $t('currency.setDefault') }}</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)" :disabled="row.is_default">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" @close="formRef?.resetFields()">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item :label="$t('currency.code')" prop="code">
          <el-input v-model="formData.code" :placeholder="$t('currency.codePlaceholder')" :disabled="isEdit" />
        </el-form-item>
        <el-form-item :label="$t('currency.name')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('currency.namePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('currency.symbol')" prop="symbol">
          <el-input v-model="formData.symbol" :placeholder="$t('currency.symbolPlaceholder')" style="width: 100px" />
        </el-form-item>
        <el-form-item :label="$t('currency.exchangeRate')" prop="exchange_rate">
          <el-input-number v-model="formData.exchange_rate" :min="0" :precision="4" :step="0.01" />
          <span class="form-tip">{{ $t('currency.exchangeRateTip') }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const tableData = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref($t('currency.addCurrency'))
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

const formData = reactive({ code: '', name: '', symbol: '', exchange_rate: 1 })
const rules = {
  code: [{ required: true, message: () => $t('currency.enterCode'), trigger: 'blur' }],
  name: [{ required: true, message: () => $t('currency.enterName'), trigger: 'blur' }],
  symbol: [{ required: true, message: () => $t('currency.enterSymbol'), trigger: 'blur' }]
}

const fetchList = async () => {
  loading.value = true
  try { const data = await request.get({ url: '/api/admin/currencies' }); tableData.value = data || [] } catch (error) { console.error('fetch currencies failed:', error) } finally { loading.value = false }
}

const handleAdd = () => { isEdit.value = false; dialogTitle.value = $t('currency.addCurrency'); editingId.value = null; Object.assign(formData, { code: '', name: '', symbol: '', exchange_rate: 1 }); dialogVisible.value = true }
const handleEdit = (row: any) => { isEdit.value = true; dialogTitle.value = $t('currency.editCurrency'); editingId.value = row.id; Object.assign(formData, { code: row.code, name: row.name, symbol: row.symbol, exchange_rate: row.exchange_rate }); dialogVisible.value = true }

const handleToggleStatus = async (row: any) => {
  try { await request.post({ url: `/api/admin/currencies/${row.id}/status`, data: { status: row.status } }); ElMessage.success($t('common.updateSuccess')) } catch (error) { console.error('update status failed:', error); fetchList() }
}

const handleSetDefault = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('currency.confirmSetDefault', { name: row.name }), $t('common.tips'), { type: 'warning' })
    await request.post({ url: `/api/admin/currencies/${row.id}/set-default` }); ElMessage.success($t('common.operationSuccess')); fetchList()
  } catch (error) { if (error !== 'cancel') console.error('set default failed:', error) }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('currency.confirmDelete', { name: row.name }), $t('common.tips'), { type: 'warning' })
    await request.del({ url: `/api/admin/currencies/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchList()
  } catch (error) { if (error !== 'cancel') console.error('delete failed:', error) }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate(); submitting.value = true
    if (isEdit.value && editingId.value) { await request.put({ url: `/api/admin/currencies/${editingId.value}`, data: formData }); ElMessage.success($t('common.updateSuccess')) }
    else { await request.post({ url: '/api/admin/currencies', data: formData }); ElMessage.success($t('common.addSuccess')) }
    dialogVisible.value = false; fetchList()
  } catch (error) { if (error !== false) console.error('submit failed:', error) } finally { submitting.value = false }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.currencies-page { padding: 16px; }
.action-card { margin-bottom: 16px; }
.form-tip { margin-left: 8px; color: var(--el-text-color-secondary); font-size: 12px; }
</style>
