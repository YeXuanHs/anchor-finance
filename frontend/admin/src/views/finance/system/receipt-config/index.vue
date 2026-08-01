<template>
  <div class="receipt-config-page">
    <art-card title="发票设置" shadow="never">
      <el-tabs v-model="activeTab">
        <!-- 费率设置 -->
        <el-tab-pane label="费率设置" name="rate">
          <el-form :model="rateForm" label-width="140px" style="max-width: 600px" v-loading="loading">
            <el-form-item label="启用发票管理">
              <el-switch v-model="rateForm.voucher_manager" :active-value="1" :inactive-value="0" />
              <div class="form-tip">开启后启用发票管理功能</div>
            </el-form-item>
            <el-form-item label="发票费率(%)">
              <el-input-number v-model="rateForm.rate" :min="0" :max="100" :precision="2" />
              <div class="form-tip">开具发票时额外收取的费率百分比</div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveRate" :loading="saveLoading">保存设置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 快递管理 -->
        <el-tab-pane label="快递管理" name="express">
          <div class="toolbar">
            <el-button type="primary" @click="showExpressDialog()">
              <el-icon><Plus /></el-icon>
              添加快递
            </el-button>
          </div>

          <el-table :data="expressList" v-loading="expressLoading" stripe border>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="快递名称" />
            <el-table-column prop="price" label="快递价格" width="120">
              <template #default="{ row }">¥{{ row.price }}</template>
            </el-table-column>
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="showExpressDialog(row)">编辑</el-button>
                <el-button size="small" type="danger" @click="handleDeleteExpress(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </art-card>

    <!-- 添加/编辑快递对话框 -->
    <el-dialog v-model="expressDialogVisible" :title="isEditExpress ? '编辑快递' : '添加快递'" width="450px">
      <el-form :model="expressForm" :rules="expressRules" ref="expressFormRef" label-width="100px">
        <el-form-item label="快递名称" prop="name">
          <el-input v-model="expressForm.name" placeholder="如：顺丰快递" />
        </el-form-item>
        <el-form-item label="快递价格" prop="price">
          <el-input-number v-model="expressForm.price" :min="0" :precision="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="expressDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveExpress" :loading="saveLoading">确定</el-button>
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
  name: [{ required: true, message: '请输入快递名称', trigger: 'blur' }]
}

const fetchConfig = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/config/receipt' })
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
    const res = await request.get({ url: '/api/admin/config/receipt/express' })
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
    await request.put({ url: '/api/admin/config/receipt', data: rateForm, showSuccessMessage: true })
  } catch (error) {
    ElMessage.error('保存失败')
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
        await request.put({ url: `/api/admin/config/receipt/express/${editingExpressId.value}`, data: expressForm, showSuccessMessage: true })
      } else {
        await request.post({ url: '/api/admin/config/receipt/express', data: expressForm, showSuccessMessage: true })
      }
      expressDialogVisible.value = false
      fetchExpressList()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      saveLoading.value = false
    }
  })
}

const handleDeleteExpress = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定删除快递 "${row.name}" 吗？`, '提示')
    await request.delete({ url: `/api/admin/config/receipt/express/${row.id}`, showSuccessMessage: true })
    fetchExpressList()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
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
