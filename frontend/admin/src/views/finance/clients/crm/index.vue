<template>
  <div class="crm-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsCrm.title') }}</span>
          <el-button type="primary" @click="handleAddRecord">
            <el-icon><Plus /></el-icon>
            {{ $t('clientsCrm.addRecord') }}
          </el-button>
        </div>
      </template>

      <el-card shadow="hover" class="client-info-card">
        <el-descriptions :column="4" border>
          <el-descriptions-item :label="$t('clientsCrm.clientName')">{{ clientInfo.username }}</el-descriptions-item>
          <el-descriptions-item :label="$t('common.email')">{{ clientInfo.email }}</el-descriptions-item>
          <el-descriptions-item :label="$t('common.phone')">{{ clientInfo.phone }}</el-descriptions-item>
          <el-descriptions-item :label="$t('clientsCrm.level')">{{ clientInfo.level_name || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('clientsCrm.totalSpent')">&yen;{{ formatAmount(clientInfo.total_spent) }}</el-descriptions-item>
          <el-descriptions-item :label="$t('clientsCrm.orderCount')">{{ clientInfo.order_count || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="$t('clientsCrm.lastLogin')">{{ clientInfo.last_login || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('clientsCrm.source')">{{ clientInfo.source || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <el-tabs v-model="activeTab" class="crm-tabs">
        <el-tab-pane :label="$t('clientsCrm.followRecords')" name="follow">
          <el-form :inline="true" :model="followSearch" class="search-form">
            <el-form-item :label="$t('common.type')">
              <el-select v-model="followSearch.type" :placeholder="$t('common.all')" clearable>
                <el-option :label="$t('clientsCrm.phone')" value="phone" />
                <el-option :label="$t('clientsCrm.email')" value="email" />
                <el-option :label="$t('clientsCrm.visiting')" value="visit" />
                <el-option :label="$t('clientsCrm.other')" value="other" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="fetchFollowRecords">{{ $t('common.search') }}</el-button>
              <el-button @click="followSearch.type = ''; fetchFollowRecords()">{{ $t('common.reset') }}</el-button>
            </el-form-item>
          </el-form>

          <el-timeline>
            <el-timeline-item
              v-for="record in followRecords" :key="record.id"
              :timestamp="record.created_at" placement="top"
              :type="getTimelineType(record.type)"
            >
              <el-card shadow="hover" class="timeline-card">
                <div class="timeline-header">
                  <el-tag :type="getTypeTag(record.type)" size="small">{{ getTypeText(record.type) }}</el-tag>
                  <span class="operator">{{ record.operator }}</span>
                </div>
                <div class="timeline-content">{{ record.content }}</div>
                <div class="timeline-next" v-if="record.next_follow_at">
                  <el-icon><Clock /></el-icon>
                  {{ $t('clientsCrm.nextFollow') }}: {{ record.next_follow_at }}
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>

          <el-empty v-if="followRecords.length === 0" :description="$t('clientsCrm.noFollowRecords')" />
        </el-tab-pane>

        <el-tab-pane :label="$t('clientsCrm.opportunityRecords')" name="opportunity">
          <el-table :data="opportunityData" style="width: 100%" border>
            <el-table-column prop="id" label="ID" width="70" align="center" />
            <el-table-column prop="title" :label="$t('clientsCrm.opportunityName')" min-width="150" />
            <el-table-column prop="amount" :label="$t('clientsCrm.estimatedAmount')" width="120" align="right">
              <template #default="{ row }">&yen;{{ formatAmount(row.amount) }}</template>
            </el-table-column>
            <el-table-column prop="stage" :label="$t('clientsCrm.stage')" width="120" align="center">
              <template #default="{ row }">
                <el-tag :type="getStageTag(row.stage)" size="small">{{ getStageText(row.stage) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="probability" :label="$t('clientsCrm.winProbability')" width="100" align="center">
              <template #default="{ row }">{{ row.probability }}%</template>
            </el-table-column>
            <el-table-column prop="expected_close_date" :label="$t('clientsCrm.expectedClose')" width="120" />
            <el-table-column prop="created_at" :label="$t('common.createdAt')" width="170" />
            <el-table-column :label="$t('common.action')" width="100" fixed="right" align="center">
              <template #default="{ row }">
                <el-button type="primary" link @click="handleEditOpportunity(row)">{{ $t('common.edit') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane :label="$t('clientsCrm.contractRecords')" name="contract">
          <el-table :data="contractData" style="width: 100%" border>
            <el-table-column prop="id" label="ID" width="70" align="center" />
            <el-table-column prop="contract_no" :label="$t('clientsCrm.contractNo')" width="170" />
            <el-table-column prop="title" :label="$t('clientsCrm.contractName')" min-width="150" />
            <el-table-column prop="amount" :label="$t('clientsCrm.contractAmount')" width="120" align="right">
              <template #default="{ row }">&yen;{{ formatAmount(row.amount) }}</template>
            </el-table-column>
            <el-table-column prop="status" :label="$t('common.status')" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getContractStatusTag(row.status)" size="small">{{ getContractStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="start_date" :label="$t('clientsCrm.startDate')" width="120" />
            <el-table-column prop="end_date" :label="$t('clientsCrm.endDate')" width="120" />
            <el-table-column :label="$t('common.action')" width="100" fixed="right" align="center">
              <template #default="{ row }">
                <el-button type="primary" link @click="handleViewContract(row)">{{ $t('common.view') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-dialog v-model="recordVisible" :title="$t('clientsCrm.addFollowRecord')" width="600px">
      <el-form :model="recordForm" :rules="recordRules" ref="recordFormRef" label-width="100px">
        <el-form-item :label="$t('clientsCrm.followMethod')" prop="type">
          <el-select v-model="recordForm.type" :placeholder="$t('common.select')">
            <el-option :label="$t('clientsCrm.phone')" value="phone" />
            <el-option :label="$t('clientsCrm.email')" value="email" />
            <el-option :label="$t('clientsCrm.visiting')" value="visit" />
            <el-option :label="$t('clientsCrm.other')" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('clientsCrm.followContent')" prop="content">
          <el-input v-model="recordForm.content" type="textarea" :rows="4" :placeholder="$t('clientsCrm.enterFollowContent')" />
        </el-form-item>
        <el-form-item :label="$t('clientsCrm.nextFollowDate')">
          <el-date-picker v-model="recordForm.next_follow_at" type="datetime" :placeholder="$t('clientsCrm.selectDateTime')" value-format="YYYY-MM-DD HH:mm:ss" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="recordVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleRecordSubmit" :loading="recordLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Plus, Clock } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const route = useRoute()
const clientId = route.params.id as string

const activeTab = ref('follow')
const recordLoading = ref(false)
const clientInfo = ref<any>({})
const followRecords = ref<any[]>([])
const opportunityData = ref<any[]>([])
const contractData = ref<any[]>([])
const followSearch = reactive({ type: '' })

const recordVisible = ref(false)
const recordFormRef = ref<FormInstance>()
const recordForm = reactive({ type: 'phone', content: '', next_follow_at: '' })

const recordRules: FormRules = {
  type: [{ required: true, message: $t('clientsCrm.selectFollowMethod'), trigger: 'change' }],
  content: [{ required: true, message: $t('clientsCrm.enterFollowContent'), trigger: 'blur' }]
}

const formatAmount = (amount: number | undefined) =>
  amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'

const getTypeText = (type: string) => {
  const map: Record<string, string> = { phone: $t('clientsCrm.phone'), email: $t('clientsCrm.email'), visit: $t('clientsCrm.visiting'), other: $t('clientsCrm.other') }
  return map[type] || type
}
const getTypeTag = (type: string) => {
  const map: Record<string, string> = { phone: 'primary', email: 'success', visit: 'warning', other: 'info' }
  return (map[type] || 'info') as any
}
const getTimelineType = (type: string) => {
  const map: Record<string, string> = { phone: 'primary', email: 'success', visit: 'warning', other: '' }
  return map[type] || '' as any
}
const getStageText = (stage: string) => {
  const map: Record<string, string> = { initial: $t('clientsCrm.initialContact'), negotiation: $t('clientsCrm.negotiation'), proposal: $t('clientsCrm.proposal'), closed_won: $t('clientsCrm.closedWon'), closed_lost: $t('clientsCrm.closedLost') }
  return map[stage] || stage
}
const getStageTag = (stage: string) => {
  const map: Record<string, string> = { initial: 'info', negotiation: 'warning', proposal: 'primary', closed_won: 'success', closed_lost: 'danger' }
  return (map[stage] || 'info') as any
}
const getContractStatusText = (status: number) => {
  const map: Record<number, string> = { 0: $t('clientsCrm.draft'), 1: $t('clientsCrm.active'), 2: $t('clientsCrm.expired'), 3: $t('clientsCrm.terminated') }
  return map[status] || $t('common.unknown')
}
const getContractStatusTag = (status: number) => {
  const map: Record<number, string> = { 0: 'info', 1: 'success', 2: 'warning', 3: 'danger' }
  return (map[status] || 'info') as any
}

const fetchClientInfo = async () => {
  try { clientInfo.value = await request.get({ url: `/api/admin/users/${clientId}` }) } catch (e) { console.error(e) }
}
const fetchFollowRecords = async () => {
  try {
    const params: any = { client_id: clientId }
    if (followSearch.type) params.type = followSearch.type
    followRecords.value = (await request.get({ url: '/api/admin/crm/follow-records', params })).list || []
  } catch (e) { console.error(e) }
}
const fetchOpportunities = async () => {
  try { opportunityData.value = (await request.get({ url: '/api/admin/crm/opportunities', params: { client_id: clientId } })).list || [] } catch (e) { console.error(e) }
}
const fetchContracts = async () => {
  try { contractData.value = (await request.get({ url: '/api/admin/crm/contracts', params: { client_id: clientId } })).list || [] } catch (e) { console.error(e) }
}

const handleAddRecord = () => { recordForm.type = 'phone'; recordForm.content = ''; recordForm.next_follow_at = ''; recordVisible.value = true }

const handleRecordSubmit = async () => {
  if (!recordFormRef.value) return
  await recordFormRef.value.validate(async (valid) => {
    if (!valid) return
    recordLoading.value = true
    try {
      await request.post({ url: '/api/admin/crm/follow-records', params: { ...recordForm, client_id: clientId } })
      ElMessage.success($t('common.addSuccess'))
      recordVisible.value = false
      fetchFollowRecords()
    } catch (e) { ElMessage.error($t('common.addFailed')) } finally { recordLoading.value = false }
  })
}

const handleEditOpportunity = (row: any) => { ElMessage.success(`${$t('common.edit')}: ${row.title || row.id}`) }
const handleViewContract = (row: any) => { window.open(`/finance/advanced/contracts?contract_id=${row.id}`, '_blank') }

onMounted(() => { fetchClientInfo(); fetchFollowRecords(); fetchOpportunities(); fetchContracts() })
</script>

<style scoped lang="scss">
.crm-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.client-info-card { margin-bottom: 20px; }
.crm-tabs { margin-top: 20px; }
.search-form { margin-bottom: 20px; }
.timeline-card { .timeline-header { display: flex; align-items: center; gap: 12px; margin-bottom: 8px; .operator { color: #909399; font-size: 13px; } } .timeline-content { color: #303133; line-height: 1.6; } .timeline-next { margin-top: 8px; color: #e6a23c; font-size: 13px; display: flex; align-items: center; gap: 4px; } }
</style>
