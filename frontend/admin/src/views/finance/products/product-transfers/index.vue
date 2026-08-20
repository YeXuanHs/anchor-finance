<template>
  <div class="page-container">
    <art-card :title="$t('productTransfers.title')" shadow="never">
      <!-- 操作栏 -->
      <div class="mb-4">
        <el-button type="primary" @click="handleAdd">{{ $t('productTransfers.addConfig') }}</el-button>
      </div>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="source_product_id" :label="$t('productTransfers.sourceProductId')" width="100" />
        <el-table-column prop="source_product_name" :label="$t('productTransfers.sourceProduct')" min-width="150" />
        <el-table-column prop="target_product_id" :label="$t('productTransfers.targetProductId')" width="100" />
        <el-table-column prop="target_product_name" :label="$t('productTransfers.targetProduct')" min-width="150" />
        <el-table-column prop="transfer_fee" :label="$t('productTransfers.transferFee')" width="120">
          <template #default="{ row }">{{ row.transfer_fee ? `¥${row.transfer_fee}` : $t('productTransfers.free') }}</template>
        </el-table-column>
        <el-table-column prop="enabled" :label="$t('common.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? $t('common.enable') : $t('common.disable') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" :label="$t('common.sort')" width="80" />
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180" />
        <el-table-column :label="$t('common.action')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button size="small" :type="row.enabled ? 'warning' : 'success'" @click="handleToggle(row)">
              {{ row.enabled ? $t('common.disable') : $t('common.enable') }}
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        class="mt-4 justify-end"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchData"
        @current-change="fetchData"
      />
    </art-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? $t('productTransfers.editConfig') : $t('productTransfers.addConfig')" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item :label="$t('productTransfers.sourceProduct')" prop="source_product_id">
          <el-select v-model="formData.source_product_id" :placeholder="$t('productTransfers.selectSourceProduct')" filterable style="width: 100%">
            <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('productTransfers.targetProduct')" prop="target_product_id">
          <el-select v-model="formData.target_product_id" :placeholder="$t('productTransfers.selectTargetProduct')" filterable style="width: 100%">
            <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('productTransfers.transferFee')" prop="transfer_fee">
          <el-input-number v-model="formData.transfer_fee" :min="0" :precision="2" :step="10" style="width: 100%" />
          <div class="text-gray-400 text-sm mt-1">{{ $t('productTransfers.freeTransferTip') }}</div>
        </el-form-item>
        <el-form-item :label="$t('common.sort')" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="9999" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('common.enable')" prop="enabled">
          <el-switch v-model="formData.enabled" />
        </el-form-item>
        <el-form-item :label="$t('common.remark')">
          <el-input v-model="formData.remark" type="textarea" :rows="2" :placeholder="$t('productTransfers.remarkPlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ isEdit ? $t('common.save') : $t('common.add') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'ProductTransfers' })

const loading = ref(false)
const tableData = ref([])
const productList = ref<any[]>([])

const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

const dialogVisible = ref(false)
const isEdit = ref(false)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()
const editId = ref<number | undefined>()
const formData = reactive({
  source_product_id: undefined as number | undefined,
  target_product_id: undefined as number | undefined,
  transfer_fee: 0,
  sort_order: 0,
  enabled: true,
  remark: ''
})
const formRules: FormRules = {
  source_product_id: [{ required: true, message: $t('productTransfers.selectSourceProduct'), trigger: 'change' }],
  target_product_id: [{ required: true, message: $t('productTransfers.selectTargetProduct'), trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get({
      url: '/api/admin/product-transfers',
      params: { page: pagination.page, page_size: pagination.pageSize }
    })
    tableData.value = res?.list || []
    pagination.total = res?.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchProducts = async () => {
  try {
    const res = await request.get({ url: '/api/admin/products' })
    productList.value = res || []
  } catch (error) {
    console.error(error)
  }
}

const handleAdd = () => {
  isEdit.value = false
  editId.value = undefined
  formData.source_product_id = undefined
  formData.target_product_id = undefined
  formData.transfer_fee = 0
  formData.sort_order = 0
  formData.enabled = true
  formData.remark = ''
  dialogVisible.value = true
  fetchProducts()
}

const handleEdit = (row: any) => {
  isEdit.value = true
  editId.value = row.id
  formData.source_product_id = row.source_product_id
  formData.target_product_id = row.target_product_id
  formData.transfer_fee = row.transfer_fee || 0
  formData.sort_order = row.sort_order || 0
  formData.enabled = row.enabled
  formData.remark = row.remark || ''
  dialogVisible.value = true
  fetchProducts()
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (isEdit.value) {
        await request.post({ url: `/api/admin/product-transfers/${editId.value}`, params: formData })
        ElMessage.success($t('common.saveSuccess'))
      } else {
        await request.post({ url: '/api/admin/product-transfers', params: formData })
        ElMessage.success($t('common.addSuccess'))
      }
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error(isEdit.value ? $t('common.saveFailed') : $t('common.addFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleToggle = async (row: any) => {
  const action = row.enabled ? $t('common.disable') : $t('common.enable')
  await ElMessageBox.confirm($t('productTransfers.confirmToggle', { action }), $t('common.tips'), { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/product-transfers/${row.id}/toggle` })
    ElMessage.success($t('common.operateSuccess'))
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm($t('productTransfers.confirmDelete'), $t('common.warning'), { type: 'error' })
  try {
    await request.del({ url: `/api/admin/product-transfers/${row.id}` })
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
