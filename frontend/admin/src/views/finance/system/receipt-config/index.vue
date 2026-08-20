<template>
  <div class="receipt-config-page">
    <art-card :title="$t('receiptConfig.title')" shadow="never">
      <el-tabs v-model="activeTab">
        <el-tab-pane :label="$t('receiptConfig.rateSettings')" name="rate">
          <el-form :model="rateForm" label-width="140px" style="max-width: 600px" v-loading="loading">
            <el-form-item :label="$t('receiptConfig.enableInvoice')">
              <el-switch v-model="rateForm.voucher_manager" :active-value="1" :inactive-value="0" />
              <div class="form-tip">{{ $t('receiptConfig.enableInvoiceTip') }}</div>
            </el-form-item>
            <el-form-item :label="$t('receiptConfig.invoiceRate')">
              <el-input-number v-model="rateForm.rate" :min="0" :max="100" :precision="2" />
              <div class="form-tip">{{ $t('receiptConfig.invoiceRateTip') }}</div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveRate" :loading="saveLoading">{{ $t('receiptConfig.saveSettings') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="$t('receiptConfig.expressManagement')" name="express">
          <div class="toolbar">
            <el-button type="primary" @click="showExpressDialog()">
              <el-icon><Plus /></el-icon>
              {{ $t('receiptConfig.addExpress') }}
            </el-button>
          </div>

          <el-table :data="expressList" v-loading="expressLoading" stripe border>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" :label="$t('receiptConfig.expressName')" />
            <el-table-column prop="price" :label="$t('receiptConfig.expressPrice')" width="120">
              <template #default="{ row }">¥{{ row.price }}</template>
            </el-table-column>
            <el-table-column :label="$t('receiptConfig.actions')" width="160" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="showExpressDialog(row)">{{ $t('receiptConfig.edit') }}</el-button>
                <el-button size="small" type="danger" @click="handleDeleteExpress(row)">{{ $t('receiptConfig.delete') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </art-card>

    <el-dialog v-model="expressDialogVisible" :title="isEditExpress ? $t('receiptConfig.editExpress') : $t('receiptConfig.addExpress')" width="450px">
      <el-form :model="expressForm" :rules="expressRules" ref="expressFormRef" label-width="100px">
        <el-form-item :label="$t('receiptConfig.expressName')" prop="name">
          <el-input v-model="expressForm.name" :placeholder="$t('receiptConfig.expressNamePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('receiptConfig.expressPrice')" prop="price">
          <el-input-number v-model="expressForm.price" :min="0" :precision="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="expressDialogVisible = false">{{ $t('receiptConfig.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveExpress" :loading="saveLoading">{{ $t('receiptConfig.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const activeTab = ref('rate')
const loading = ref(false)
const saveLoading = ref(false)
const expressLoading = ref(false)
const expressDialogVisible = ref(false)
const isEditExpress = ref(false)
const editingExpressId = ref<number | null>(null)
const expressFormRef = ref<FormInstance>()

const rateForm = reactive({
  voucher_manager: 0,
  rate: 0
})

const expressList = ref<Array<{ id: number; name: string; price: number }>>([])
const expressForm = reactive({ name: '', price: 0 })
const expressRules: FormRules = {
  name: [{ required: true, message: $t('receiptConfig.expressNameRequired'), trigger: 'blur' }]
}

const fetchConfig = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/config/invoice' })
    if (res?.data) Object.assign(rateForm, res.data)
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchExpressList = async () => {
  expressLoading.value = true
  try {
    const res = await request.get({ url: '/api/admin/config/invoice/express' })
    expressList.value = res?.data || []
  } catch (error) {
    console.error(error)
  } finally {
    expressLoading.value = false
  }
}

const handleSaveRate = async () => {
  saveLoading.value = true
  try {
    await request.put({ url: '/api/admin/config/invoice', data: rateForm, showSuccessMessage: true })
  } catch (error) {
    ElMessage.error($t('receiptConfig.saveFailed'))
  } finally {
    saveLoading.value = false
  }
}

const showExpressDialog = (row?: any) => {
  if (row) {
    isEditExpress.value = true
    editingExpressId.value = row.id
    expressForm.name = row.name
    expressForm.price = row.price
  } else {
    isEditExpress.value = false
    editingExpressId.value = null
    expressForm.name = ''
    expressForm.price = 0
  }
  expressDialogVisible.value = true
}

const handleSaveExpress = async () => {
  if (!expressFormRef.value) return
  await expressFormRef.value.validate(async (valid) => {
    if (!valid) return
    saveLoading.value = true
    try {
      if (isEditExpress.value && editingExpressId.value) {
        await request.put({ url: `/api/admin/config/invoice/express/${editingExpressId.value}`, data: expressForm, showSuccessMessage: true })
      } else {
        await request.post({ url: '/api/admin/config/invoice/express', data: expressForm, showSuccessMessage: true })
      }
      expressDialogVisible.value = false
      fetchExpressList()
    } catch (error) {
      ElMessage.error($t('receiptConfig.operationFailed'))
    } finally {
      saveLoading.value = false
    }
  })
}

const handleDeleteExpress = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('receiptConfig.confirmDeleteExpress', { name: row.name }), $t('receiptConfig.tip'))
    await request.del({ url: `/api/admin/config/invoice/express/${row.id}`, showSuccessMessage: true })
    fetchExpressList()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error($t('receiptConfig.deleteFailed'))
  }
}

onMounted(() => {
  fetchConfig()
  fetchExpressList()
})
</script>

<style scoped lang="scss">
.receipt-config-page {
  padding: 20px;
}
.form-tip {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
}
.toolbar {
  margin-bottom: 16px;
}
</style>
