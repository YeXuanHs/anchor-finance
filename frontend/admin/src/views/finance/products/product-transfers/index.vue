<template>
  <div class="page-container">
    <art-card title="产品转移配置" shadow="never">
      <!-- 操作栏 -->
      <div class="mb-4">
        <el-button type="primary" @click="handleAdd">添加配置</el-button>
      </div>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="source_product_id" label="源产品ID" width="100" />
        <el-table-column prop="source_product_name" label="源产品" min-width="150" />
        <el-table-column prop="target_product_id" label="目标产品ID" width="100" />
        <el-table-column prop="target_product_name" label="目标产品" min-width="150" />
        <el-table-column prop="transfer_fee" label="转移费用" width="120">
          <template #default="{ row }">{{ row.transfer_fee ? `¥${row.transfer_fee}` : '免费' }}</template>
        </el-table-column>
        <el-table-column prop="enabled" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" :type="row.enabled ? 'warning' : 'success'" @click="handleToggle(row)">
              {{ row.enabled ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑配置' : '添加配置'" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item label="源产品" prop="source_product_id">
          <el-select v-model="formData.source_product_id" placeholder="请选择源产品" filterable style="width: 100%">
            <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标产品" prop="target_product_id">
          <el-select v-model="formData.target_product_id" placeholder="请选择目标产品" filterable style="width: 100%">
            <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="转移费用" prop="transfer_fee">
          <el-input-number v-model="formData.transfer_fee" :min="0" :precision="2" :step="10" style="width: 100%" />
          <div class="text-gray-400 text-sm mt-1">设置为0表示免费转移</div>
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="9999" style="width: 100%" />
        </el-form-item>
        <el-form-item label="启用" prop="enabled">
          <el-switch v-model="formData.enabled" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="formData.remark" type="textarea" :rows="2" placeholder="备注信息（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ isEdit ? '保存' : '添加' }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'ProductTransfers' })

const loading = ref(false)
const tableData = ref([])
const productList = ref<any[]>([])

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 对话框
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
  source_product_id: [{ required: true, message: '请选择源产品', trigger: 'change' }],
  target_product_id: [{ required: true, message: '请选择目标产品', trigger: 'change' }]
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
        ElMessage.success('保存成功')
      } else {
        await request.post({ url: '/api/admin/product-transfers', params: formData })
        ElMessage.success('添加成功')
      }
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error(isEdit.value ? '保存失败' : '添加失败')
    } finally {
      submitLoading.value = false
    }
  })
}

const handleToggle = async (row: any) => {
  const action = row.enabled ? '禁用' : '启用'
  await ElMessageBox.confirm(`确定${action}该配置？`, '提示', { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/product-transfers/${row.id}/toggle` })
    ElMessage.success(`${action}成功`)
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm('确定删除该配置？', '警告', { type: 'error' })
  try {
    await request.del({ url: `/api/admin/product-transfers/${row.id}` })
    ElMessage.success('删除成功')
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
