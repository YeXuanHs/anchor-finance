<template>
  <div class="supplier-list-page">
    <!-- 页面描述 -->
    <el-card shadow="never" class="desc-card">
      <div class="desc-content">
        <span>{{ $t('finance.upstreamProviders.description') }}</span>
        <el-link type="primary" href="https://bbs.idcsmart.com/forum.php?mod=viewthread&tid=136&extra=page%3D1%26filter%3Dtypeid%26typeid%3D7" target="_blank">{{ $t('finance.upstreamProviders.helpDoc') }}</el-link>
      </div>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        {{ $t('finance.upstreamProviders.addSupplier') }}
      </el-button>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="name" :label="$t('finance.upstreamProviders.name')" min-width="120">
          <template #default="{ row }">
            <el-link type="primary" @click="handleEdit(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.upstreamProviders.type')" width="100" align="center">
          <template #default="{ row }">
            {{ getApiTypeText(row.api_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="api_url" :label="$t('finance.upstreamProviders.apiUrl')" min-width="200" show-overflow-tooltip />
        <el-table-column :label="$t('finance.upstreamProviders.availableConfigured')" width="140" align="center">
          <template #default="{ row }">
            {{ row.available_products || 0 }}/{{ row.configured_products || 0 }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.upstreamProviders.productCount')" width="140" align="center">
          <template #default="{ row }">
            {{ row.normal_products || 0 }}/{{ row.total_products || 0 }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.upstreamProviders.status')" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              active-value="active"
              inactive-value="disabled"
              @change="handleToggleStatus(row)"
            />
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.upstreamProviders.balance')" width="100" align="center">
          <template #default="{ row }">
            {{ row.balance ? `¥${row.balance}` : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="description" :label="$t('finance.upstreamProviders.descriptionLabel')" min-width="120" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.description || '-' }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('finance.upstreamProviders.manage')" width="140" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              <el-icon><Edit /></el-icon> {{ $t('finance.upstreamProviders.edit') }}
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">
              <el-icon><Delete /></el-icon> {{ $t('finance.upstreamProviders.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </el-card>

    <!-- 添加/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px">
        <el-form-item :label="$t('finance.upstreamProviders.supplierName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('finance.upstreamProviders.enterSupplierName')" />
        </el-form-item>

        <el-form-item :label="$t('finance.upstreamProviders.apiType')" prop="api_type">
          <el-select v-model="formData.api_type" :placeholder="$t('finance.upstreamProviders.selectApiType')" @change="handleApiTypeChange">
            <el-option :label="$t('finance.upstreamProviders.typeManual')" value="manual" />
            <el-option :label="$t('finance.upstreamProviders.typeZjmf')" value="zjmf" />
            <el-option :label="$t('finance.upstreamProviders.typeV10')" value="v10" />
            <el-option :label="$t('finance.upstreamProviders.typeAnchor')" value="anchor" />
          </el-select>
        </el-form-item>

        <template v-if="formData.api_type !== 'manual'">
          <el-form-item :label="$t('finance.upstreamProviders.apiUrl')" prop="api_url">
            <el-input v-model="formData.api_url" placeholder="https://api.example.com" />
          </el-form-item>

          <el-form-item :label="$t('finance.upstreamProviders.apiKey')" prop="api_key">
            <el-input v-model="formData.api_key" :placeholder="$t('finance.upstreamProviders.enterApiKey')" show-password />
          </el-form-item>

          <el-form-item :label="$t('finance.upstreamProviders.apiPassword')" prop="api_password" v-if="formData.api_type === 'zjmf' || formData.api_type === 'v10'">
            <el-input v-model="formData.api_password" :placeholder="$t('finance.upstreamProviders.enterApiPassword')" show-password />
          </el-form-item>
        </template>

        <el-form-item :label="$t('finance.upstreamProviders.descriptionLabel')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('finance.upstreamProviders.enterDescription')" />
        </el-form-item>

        <el-form-item :label="$t('finance.upstreamProviders.status')" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio value="active">{{ $t('finance.upstreamProviders.enabled') }}</el-radio>
            <el-radio value="disabled">{{ $t('finance.upstreamProviders.disabled') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('finance.upstreamProviders.cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">{{ $t('finance.upstreamProviders.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const tableData = ref([])
const currentPage = ref(1)
const pageSize = ref(100)
const total = ref(0)

// 弹窗
const dialogVisible = ref(false)
const dialogTitle = ref($t('finance.upstreamProviders.addSupplier'))
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

// 表单数据
const formData = reactive({
  name: '',
  api_type: 'manual',
  api_url: '',
  api_key: '',
  api_password: '',
  description: '',
  status: 'active' as string
})

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: $t('finance.upstreamProviders.enterSupplierName'), trigger: 'blur' }
  ],
  api_type: [
    { required: true, message: $t('finance.upstreamProviders.selectApiType'), trigger: 'change' }
  ]
}

// API类型文本
const getApiTypeText = (type: string) => {
  const map: Record<string, string> = {
    manual: $t('finance.upstreamProviders.typeManual'),
    zjmf: $t('finance.upstreamProviders.typeZjmf'),
    v10: 'v10',
    anchor: $t('finance.upstreamProviders.typeAnchor')
  }
  return map[type] || $t('finance.upstreamProviders.unknown')
}

// 获取列表数据
const fetchList = async () => {
  loading.value = true
  try {
    const res = await request.get({
      url: '/api/admin/suppliers',
      params: { page: currentPage.value, page_size: pageSize.value }
    })
    tableData.value = res?.data || res || []
    total.value = res?.total || tableData.value.length
  } catch (error) {
    console.error('获取供应商列表失败:', error)
  } finally {
    loading.value = false
  }
}

// API类型变化
const handleApiTypeChange = (type: string) => {
  if (type === 'manual') {
    formData.api_url = ''
    formData.api_key = ''
    formData.api_password = ''
  }
}

// 添加供应商
const handleAdd = () => {
  isEdit.value = false
  dialogTitle.value = $t('finance.upstreamProviders.addSupplier')
  editingId.value = null
  resetForm()
  dialogVisible.value = true
}

// 编辑供应商
const handleEdit = (row: any) => {
  isEdit.value = true
  dialogTitle.value = $t('finance.upstreamProviders.editSupplier')
  editingId.value = row.id
  Object.assign(formData, {
    name: row.name,
    api_type: row.api_type,
    api_url: row.api_url || '',
    api_key: row.api_key || '',
    api_password: '',
    description: row.description || '',
    status: row.status
  })
  dialogVisible.value = true
}

// 切换状态
const handleToggleStatus = async (row: any) => {
  try {
    await request.put({
      url: `/api/admin/suppliers/${row.id}`,
      data: { status: row.status }
    })
    ElMessage.success($t('finance.upstreamProviders.statusUpdateSuccess'))
  } catch (error) {
    console.error('更新状态失败:', error)
    fetchList()
  }
}

// 删除供应商
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('finance.upstreamProviders.confirmDeleteMessage', { name: row.name }), $t('finance.upstreamProviders.confirmDeleteTitle'), {
      type: 'warning'
    })
    await request.del({ url: `/api/admin/suppliers/${row.id}` })
    ElMessage.success($t('finance.upstreamProviders.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    if (isEdit.value && editingId.value) {
      await request.put({ url: `/api/admin/suppliers/${editingId.value}`, data: formData })
      ElMessage.success($t('finance.upstreamProviders.updateSuccess'))
    } else {
      await request.post({ url: '/api/admin/suppliers', data: formData })
      ElMessage.success($t('finance.upstreamProviders.addSuccess'))
    }

    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    submitting.value = false
  }
}

// 重置表单
const resetForm = () => {
  formData.name = ''
  formData.api_type = 'manual'
  formData.api_url = ''
  formData.api_key = ''
  formData.api_password = ''
  formData.description = ''
  formData.status = 'active'
}

// 弹窗关闭
const handleDialogClose = () => {
  formRef.value?.resetFields()
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.supplier-list-page {
  padding: 16px;
}

.desc-card {
  margin-bottom: 16px;

  :deep(.el-card__body) {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}

.desc-content {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #666;
  font-size: 14px;
}

.table-card {
  :deep(.el-card__body) {
    padding: 0;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px;
}
</style>
