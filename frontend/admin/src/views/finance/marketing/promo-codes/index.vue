<template>
  <div class="promo-codes-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('promoCode.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('promoCode.addCode') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('promoCode.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('promoCode.code')" clearable />
        </el-form-item>
        <el-form-item :label="$t('common.type')">
          <el-select v-model="searchForm.type" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('promoCode.percentage')" value="percentage" />
            <el-option :label="$t('promoCode.fixedAmount')" value="fixed" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('common.enabled')" :value="1" />
            <el-option :label="$t('common.disabled')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="code" :label="$t('promoCode.code')" width="150">
          <template #default="{ row }">
            <el-tag type="info" effect="plain">{{ row.code }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" :label="$t('common.type')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.type === 'percentage' ? 'primary' : 'success'" size="small">
              {{ row.type === 'percentage' ? $t('promoCode.percentage') : $t('promoCode.fixedAmount') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="value" :label="$t('promoCode.discountValue')" width="120" align="center">
          <template #default="{ row }">
            <span v-if="row.type === 'percentage'">{{ row.value }}%</span>
            <span v-else class="amount-text">¥{{ formatAmount(row.value) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('promoCode.usageStatus')" width="130" align="center">
          <template #default="{ row }">
            {{ row.used_count }} / {{ row.max_uses || $t('promoCode.unlimited') }}
          </template>
        </el-table-column>
        <el-table-column prop="expires_at" :label="$t('promoCode.expiresAt')" width="170">
          <template #default="{ row }">
            <span :class="{ 'danger-text': isExpired(row.expires_at) }">
              {{ row.expires_at || $t('promoCode.neverExpire') }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="is_enabled" :label="$t('common.status')" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_enabled"
              :active-value="1"
              :inactive-value="0"
              @change="handleToggleStatus(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="170" />
        <el-table-column :label="$t('common.action')" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-popconfirm :title="$t('promoCode.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('common.delete') }}</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="110px">
        <el-form-item :label="$t('promoCode.code')" prop="code">
          <el-input v-model="formData.code" :placeholder="$t('promoCode.autoGenerate')">
            <template #append>
              <el-button @click="handleGenerateCode">{{ $t('promoCode.randomGenerate') }}</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('common.type')" prop="type">
          <el-radio-group v-model="formData.type">
            <el-radio value="percentage">{{ $t('promoCode.percentage') }}</el-radio>
            <el-radio value="fixed">{{ $t('promoCode.fixedAmount') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="formData.type === 'percentage' ? $t('promoCode.discountPercent') : $t('promoCode.discountAmount')" prop="value">
          <el-input-number
            v-model="formData.value"
            :min="formData.type === 'percentage' ? 1 : 0.01"
            :max="formData.type === 'percentage' ? 100 : 99999"
            :precision="formData.type === 'percentage' ? 0 : 2"
            controls-position="right"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="$t('promoCode.maxUses')" prop="max_uses">
          <el-input-number
            v-model="formData.max_uses"
            :min="0"
            controls-position="right"
            style="width: 100%"
            :placeholder="$t('promoCode.unlimitedPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="$t('promoCode.expiresAt')" prop="expires_at">
          <el-date-picker
            v-model="formData.expires_at"
            type="datetime"
            :placeholder="$t('promoCode.neverExpirePlaceholder')"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="$t('common.status')" prop="is_enabled">
          <el-switch v-model="formData.is_enabled" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref($t('promoCode.addCodeTitle'))
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: '',
  type: '',
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  code: '',
  type: 'percentage',
  value: 10,
  max_uses: 0,
  expires_at: '',
  is_enabled: 1
})

const formRules: FormRules = {
  type: [
    { required: true, message: $t('promoCode.selectType'), trigger: 'change' }
  ],
  value: [
    { required: true, message: $t('promoCode.enterDiscountValue'), trigger: 'blur' }
  ]
}

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const isExpired = (expiresAt: string) => {
  if (!expiresAt) return false
  return new Date(expiresAt).getTime() < Date.now()
}

const handleGenerateCode = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let code = ''
  for (let i = 0; i < 8; i++) {
    code += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  formData.code = code
}

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/promo-codes',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        keyword: searchForm.keyword || undefined,
        type: searchForm.type || undefined,
        status: searchForm.status
      }
    })
    tableData.value = data.list || data || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error($t('promoCode.fetchFailed'), error)
    ElMessage.error($t('promoCode.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.type = ''
  searchForm.status = undefined
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = $t('promoCode.addCodeTitle')
  formData.id = undefined
  formData.code = ''
  formData.type = 'percentage'
  formData.value = 10
  formData.max_uses = 0
  formData.expires_at = ''
  formData.is_enabled = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('promoCode.editCodeTitle')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleToggleStatus = async (row: any) => {
  try {
    await request.put({
      url: `/api/admin/promo-codes/${row.id}`,
      params: { is_enabled: row.is_enabled }
    })
    ElMessage.success(row.is_enabled ? $t('promoCode.enabled') : $t('promoCode.disabled'))
  } catch (error) {
    row.is_enabled = row.is_enabled ? 0 : 1
    ElMessage.error($t('common.operateFailed'))
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/promo-codes/${row.id}` })
    ElMessage.success($t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    ElMessage.error($t('common.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/promo-codes/${formData.id}`, params: formData })
      } else {
        await request.post({ url: '/api/admin/promo-codes', params: formData })
      }
      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess'))
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      ElMessage.error($t('common.operateFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchList()
}

const handlePageChange = () => {
  fetchList()
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.promo-codes-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.amount-text {
  font-weight: 600;
  color: #f56c6c;
}

.danger-text {
  color: #f56c6c;
  font-weight: 600;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
