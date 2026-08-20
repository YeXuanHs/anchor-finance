<template>
  <div class="promo-plan-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('promoPlan.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('promoPlan.addPlan') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('promoPlan.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('promoPlan.planName')" clearable />
        </el-form-item>
        <el-form-item :label="$t('common.type')">
          <el-select v-model="searchForm.type" :placeholder="$t('common.all')" clearable>
            <el-option label="CPS" value="cps" />
            <el-option label="CPA" value="cpa" />
            <el-option label="CPC" value="cpc" />
            <el-option :label="$t('promoPlan.mixed')" value="mixed" />
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
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="$t('promoPlan.planName')" min-width="180" />
        <el-table-column prop="type" :label="$t('common.type')" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">
              {{ getTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="commission_rate" :label="$t('promoPlan.commissionRate')" width="120">
          <template #default="{ row }">
            {{ (row.commission_rate * 100).toFixed(1) }}%
          </template>
        </el-table-column>
        <el-table-column prop="click_count" :label="$t('promoPlan.clickCount')" width="100" />
        <el-table-column prop="conversion_count" :label="$t('promoPlan.conversionCount')" width="100" />
        <el-table-column prop="total_commission" :label="$t('promoPlan.totalCommission')" width="120">
          <template #default="{ row }">
            <span class="text-primary">¥{{ row.total_commission?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('common.status')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? $t('common.enabled') : $t('common.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180" />
        <el-table-column :label="$t('common.action')" width="280" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="success" link @click="handleGenerateLink(row)">{{ $t('promoPlan.promoLink') }}</el-button>
            <el-button type="info" link @click="handleViewStats(row)">{{ $t('promoPlan.effectStats') }}</el-button>
            <el-popconfirm :title="$t('promoPlan.confirmDelete')" @confirm="handleDelete(row)">
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('promoPlan.planName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('promoPlan.enterPlanName')" />
        </el-form-item>
        <el-form-item :label="$t('promoPlan.promoType')" prop="type">
          <el-select v-model="formData.type" :placeholder="$t('promoPlan.selectPromoType')" style="width: 100%">
            <el-option :label="$t('promoPlan.cpsLabel')" value="cps" />
            <el-option :label="$t('promoPlan.cpaLabel')" value="cpa" />
            <el-option :label="$t('promoPlan.cpcLabel')" value="cpc" />
            <el-option :label="$t('promoPlan.mixedMode')" value="mixed" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('promoPlan.commissionRate')" prop="commission_rate">
          <el-input-number v-model="formData.commission_rate" :min="0" :max="1" :step="0.01" :precision="2" style="width: 100%" />
          <div class="form-tip">{{ $t('promoPlan.rateTip') }}</div>
        </el-form-item>
        <el-form-item :label="$t('promoPlan.planDesc')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('promoPlan.enterPlanDesc')" />
        </el-form-item>
        <el-form-item :label="$t('common.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" :active-text="$t('common.enabled')" :inactive-text="$t('common.disabled')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 推广链接对话框 -->
    <el-dialog v-model="linkDialogVisible" :title="$t('promoPlan.promoLink')" width="600px">
      <el-form label-width="100px">
        <el-form-item :label="$t('promoPlan.planName')">
          <el-input :value="currentPlan.name" disabled />
        </el-form-item>
        <el-form-item :label="$t('promoPlan.promoLink')">
          <el-input v-model="promoLink" readonly>
            <template #append>
              <el-button @click="handleCopyLink">{{ $t('promoPlan.copy') }}</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('promoPlan.shortLink')">
          <el-input v-model="shortLink" readonly>
            <template #append>
              <el-button @click="handleCopyShortLink">{{ $t('promoPlan.copy') }}</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('promoPlan.qrcode')">
          <div class="qrcode-container">
            <el-image :src="qrcodeUrl" style="width: 160px; height: 160px" fit="contain" />
          </div>
        </el-form-item>
      </el-form>
    </el-dialog>

    <!-- 效果统计对话框 -->
    <el-dialog v-model="statsDialogVisible" :title="$t('promoPlan.effectStatsTitle')" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item :label="$t('promoPlan.planName')">{{ currentPlan.name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('promoPlan.promoType')">{{ getTypeText(currentPlan.type) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('promoPlan.totalClicks')">{{ statsData.click_count || 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('promoPlan.totalConversions')">{{ statsData.conversion_count || 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('promoPlan.conversionRate')">{{ statsData.conversion_rate || '0.00' }}%</el-descriptions-item>
        <el-descriptions-item :label="$t('promoPlan.totalCommission')">¥{{ statsData.total_commission?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('promoPlan.monthClicks')">{{ statsData.month_click_count || 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('promoPlan.monthConversions')">{{ statsData.month_conversion_count || 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('promoPlan.monthCommission')">¥{{ statsData.month_commission?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('promoPlan.activePromoters')">{{ statsData.active_promoters || 0 }}</el-descriptions-item>
      </el-descriptions>
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

defineOptions({ name: 'PromoPlan' })

const loading = ref(false)
const submitLoading = ref(false)

const searchForm = reactive({
  keyword: '',
  type: undefined as string | undefined,
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref([])

const dialogVisible = ref(false)
const dialogTitle = ref($t('promoPlan.addPlanTitle'))
const formRef = ref<FormInstance>()

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  type: 'cps',
  commission_rate: 0.1,
  description: '',
  status: 1
})

const formRules: FormRules = {
  name: [
    { required: true, message: $t('promoPlan.enterPlanName'), trigger: 'blur' },
    { min: 2, max: 50, message: $t('promoPlan.nameLength'), trigger: 'blur' }
  ],
  type: [
    { required: true, message: $t('promoPlan.selectPromoType'), trigger: 'change' }
  ],
  commission_rate: [
    { required: true, message: $t('promoPlan.enterCommissionRate'), trigger: 'blur' }
  ]
}

const linkDialogVisible = ref(false)
const currentPlan = ref<any>({})
const promoLink = ref('')
const shortLink = ref('')
const qrcodeUrl = ref('')

const statsDialogVisible = ref(false)
const statsData = ref<any>({})

const getTypeTag = (type: string) => {
  const map: Record<string, any> = {
    cps: 'success',
    cpa: 'primary',
    cpc: 'warning',
    mixed: 'info'
  }
  return map[type] || 'info'
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    cps: 'CPS',
    cpa: 'CPA',
    cpc: 'CPC',
    mixed: $t('promoPlan.mixed')
  }
  return map[type] || $t('common.unknown')
}

const fetchPlans = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/promo-plans',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error($t('promoPlan.fetchFailed'), error)
    ElMessage.error($t('promoPlan.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchPlans()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.type = undefined
  searchForm.status = undefined
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = $t('promoPlan.addPlanTitle')
  formData.id = undefined
  formData.name = ''
  formData.type = 'cps'
  formData.commission_rate = 0.1
  formData.description = ''
  formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('promoPlan.editPlanTitle')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleGenerateLink = async (row: any) => {
  currentPlan.value = row
  try {
    const data = await request.get({
      url: `/api/admin/promo-plans/${row.id}/link`
    })
    promoLink.value = data.link || ''
    shortLink.value = data.short_link || ''
    qrcodeUrl.value = data.qrcode_url || ''
    linkDialogVisible.value = true
  } catch (error) {
    ElMessage.error($t('promoPlan.getLinkFailed'))
  }
}

const handleCopyLink = () => {
  navigator.clipboard.writeText(promoLink.value)
  ElMessage.success($t('promoPlan.linkCopied'))
}

const handleCopyShortLink = () => {
  navigator.clipboard.writeText(shortLink.value)
  ElMessage.success($t('promoPlan.shortLinkCopied'))
}

const handleViewStats = async (row: any) => {
  currentPlan.value = row
  try {
    const data = await request.get({
      url: `/api/admin/promo-plans/${row.id}/stats`
    })
    statsData.value = data || {}
    statsDialogVisible.value = true
  } catch (error) {
    ElMessage.error($t('promoPlan.getStatsFailed'))
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/promo-plans/${row.id}`
    })
    ElMessage.success($t('common.deleteSuccess'))
    fetchPlans()
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
      const url = formData.id ? `/api/admin/promo-plans/${formData.id}` : '/api/admin/promo-plans'

      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }

      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess'))
      dialogVisible.value = false
      fetchPlans()
    } catch (error) {
      ElMessage.error($t('common.operateFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchPlans()
}

const handlePageChange = () => {
  fetchPlans()
}

onMounted(() => {
  fetchPlans()
})
</script>

<style scoped lang="scss">
.promo-plan-page {
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

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.text-primary {
  color: #409eff;
  font-weight: 600;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.qrcode-container {
  text-align: center;
  padding: 10px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}
</style>
