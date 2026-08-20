<template>
  <div class="blacklist-page">
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('blacklist.addBlacklist') }}</el-button>
      </div>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="type" :label="$t('blacklist.type')" width="100" align="center">
          <template #default="{ row }"><el-tag size="small">{{ getTypeText(row.type) }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="value" :label="$t('blacklist.value')" min-width="200" />
        <el-table-column prop="reason" :label="$t('blacklist.reason')" min-width="200" />
        <el-table-column prop="created_at" :label="$t('blacklist.addTime')" width="170" />
        <el-table-column :label="$t('blacklist.operations')" width="100" fixed="right">
          <template #default="{ row }"><el-button type="danger" link size="small" @click="handleDelete(row)">{{ $t('common.delete') }}</el-button></template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="$t('blacklist.addBlacklist')" width="500px" @close="formRef?.resetFields()">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="80px">
        <el-form-item :label="$t('blacklist.type')" prop="type">
          <el-select v-model="formData.type" :placeholder="$t('blacklist.selectType')">
            <el-option :label="$t('blacklist.typeIp')" value="ip" />
            <el-option :label="$t('blacklist.typeEmail')" value="email" />
            <el-option :label="$t('blacklist.typePhone')" value="phone" />
            <el-option :label="$t('blacklist.typeDomain')" value="domain" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('blacklist.value')" prop="value">
          <el-input v-model="formData.value" :placeholder="$t('blacklist.enterValue')" />
        </el-form-item>
        <el-form-item :label="$t('blacklist.reason')" prop="reason">
          <el-input v-model="formData.reason" type="textarea" :rows="3" :placeholder="$t('blacklist.enterReason')" />
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
const submitting = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({ type: 'ip', value: '', reason: '' })
const rules = {
  type: [{ required: true, message: () => $t('blacklist.selectType'), trigger: 'change' }],
  value: [{ required: true, message: () => $t('blacklist.enterValue'), trigger: 'blur' }]
}

const getTypeText = (type: string) => { const map: Record<string, () => string> = { ip: () => $t('blacklist.typeIp'), email: () => $t('blacklist.typeEmail'), phone: () => $t('blacklist.typePhone'), domain: () => $t('blacklist.typeDomain') }; return map[type]?.() || type }

const fetchList = async () => { loading.value = true; try { const data = await request.get({ url: '/api/admin/blacklist' }); tableData.value = data || [] } catch {} finally { loading.value = false } }
const handleAdd = () => { Object.assign(formData, { type: 'ip', value: '', reason: '' }); dialogVisible.value = true }

const handleDelete = async (row: any) => {
  try { await ElMessageBox.confirm($t('blacklist.confirmDelete'), $t('common.tips'), { type: 'warning' }); await request.del({ url: `/api/admin/blacklist/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchList() } catch (error) { if (error !== 'cancel') console.error('delete failed:', error) }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try { await formRef.value.validate(); submitting.value = true; await request.post({ url: '/api/admin/blacklist', data: formData }); ElMessage.success($t('common.addSuccess')); dialogVisible.value = false; fetchList() } catch {} finally { submitting.value = false }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.blacklist-page { padding: 16px; }
.action-card { margin-bottom: 16px; }
.action-bar { display: flex; justify-content: space-between; align-items: center; }
.table-card { :deep(.el-card__body) { padding: 0; } }
</style>
